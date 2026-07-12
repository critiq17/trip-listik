package invites

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/critiq17/tripListik/internal/telegram"
	"github.com/google/uuid"
)

// Sentinel errors — mapped to HTTP status codes in the handler layer.
var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrDuplicate = errors.New("invite already exists")
)

// Service encapsulates all invite business logic including TG notifications.
type Service struct {
	store *store.Store
	bot   *telegram.Bot
	cfg   *config.Config
}

func NewService(s *store.Store, bot *telegram.Bot, cfg *config.Config) *Service {
	return &Service{store: s, bot: bot, cfg: cfg}
}

// ── InviteUser ────────────────────────────────────────────────────────────────

// InviteUser creates an invite for the given trip. The invited user is resolved by
// username or user_id. Returns ErrDuplicate (with the existing invite) if already
// pending, so the handler can return a 200 with the existing invite status.
func (s *Service) InviteUser(ctx context.Context, requesterID, tripID uuid.UUID, username, userIDStr string) (*models.TripInvite, error) {
	trip, err := s.store.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("invites.InviteUser get trip: %w", err)
	}
	if trip == nil {
		return nil, ErrNotFound
	}

	// Permission: owner always allowed; members allowed only on group trips.
	if trip.OwnerID != requesterID {
		if trip.Visibility != "group" {
			return nil, ErrForbidden
		}
		isMember, err := s.store.IsTripMember(ctx, tripID, requesterID)
		if err != nil {
			return nil, fmt.Errorf("invites.InviteUser check member: %w", err)
		}
		if !isMember {
			return nil, ErrForbidden
		}
	}

	// Resolve invited user.
	var invitedUserID uuid.UUID
	if userIDStr != "" {
		id, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, ErrNotFound
		}
		invitedUserID = id
	} else {
		users, err := s.store.SearchUsers(ctx, username, 1)
		if err != nil {
			return nil, fmt.Errorf("invites.InviteUser search: %w", err)
		}
		if len(users) == 0 {
			return nil, ErrNotFound
		}
		invitedUserID = users[0].ID
	}

	// Guard against duplicates. A pending or accepted invite blocks a new one;
	// a declined or cancelled invite is reopened so the user can be re-invited
	// (the table has UNIQUE(trip_id, invited_user_id), so we reuse the row).
	existing, err := s.store.GetInviteByTripAndUser(ctx, tripID, invitedUserID)
	if err != nil {
		return nil, fmt.Errorf("invites.InviteUser check dup: %w", err)
	}

	var invite *models.TripInvite
	switch {
	case existing == nil:
		invite, err = s.store.CreateInvite(ctx, tripID, invitedUserID, requesterID)
		if err != nil {
			return nil, fmt.Errorf("invites.InviteUser create: %w", err)
		}
	case existing.Status == "declined" || existing.Status == "cancelled":
		invite, err = s.store.ReopenInvite(ctx, existing.ID, requesterID)
		if err != nil {
			return nil, fmt.Errorf("invites.InviteUser reopen: %w", err)
		}
	default:
		return existing, ErrDuplicate
	}

	// Fire-and-forget TG notification.
	go s.sendInviteNotification(invite, trip, requesterID, invitedUserID)

	return invite, nil
}

// ── RespondInvite ─────────────────────────────────────────────────────────────

// RespondInvite accepts or declines an invite on behalf of the invited user.
// On accept, the user is added as a trip member.
// Regardless of action, the TG invite message is deleted and the trip owner
// is notified.
func (s *Service) RespondInvite(ctx context.Context, userID, inviteID uuid.UUID, action string, comment *string, alternativeDate *string) error {
	invite, err := s.store.GetInviteByID(ctx, inviteID)
	if err != nil {
		return fmt.Errorf("invites.RespondInvite get: %w", err)
	}
	if invite == nil {
		return ErrNotFound
	}
	if invite.InvitedUserID != userID {
		return ErrForbidden
	}
	// Only pending invites can be answered — a cancelled or already answered
	// invite must not add members or spam the owner again.
	if invite.Status != "pending" {
		return ErrNotFound
	}

	status := map[string]string{"accept": "accepted", "decline": "declined"}[action]

	if err := s.store.UpdateInviteStatus(ctx, inviteID, status, comment, alternativeDate); err != nil {
		return fmt.Errorf("invites.RespondInvite update status: %w", err)
	}

	if status == "accepted" {
		if err := s.store.AddTripMember(ctx, invite.TripID, userID, "member"); err != nil {
			return fmt.Errorf("invites.RespondInvite add member: %w", err)
		}
	}

	// Delete TG invite message and notify owner — non-blocking.
	go s.handleInviteResponse(invite, userID, status, comment, alternativeDate)

	return nil
}

// ── CancelInvite ──────────────────────────────────────────────────────────────

// CancelInvite lets the trip owner rescind a pending invite.
// If a TG notification was already sent to the invitee, it is deleted.
func (s *Service) CancelInvite(ctx context.Context, requesterID, inviteID uuid.UUID) error {
	invite, err := s.store.GetInviteByID(ctx, inviteID)
	if err != nil {
		return fmt.Errorf("invites.CancelInvite get: %w", err)
	}
	if invite == nil {
		return ErrNotFound
	}

	trip, err := s.store.GetTripByID(ctx, invite.TripID)
	if err != nil {
		return fmt.Errorf("invites.CancelInvite get trip: %w", err)
	}
	if trip == nil {
		return ErrNotFound
	}
	if trip.OwnerID != requesterID {
		return ErrForbidden
	}

	if err := s.store.CancelInvite(ctx, inviteID); err != nil {
		return fmt.Errorf("invites.CancelInvite store: %w", err)
	}

	// Delete TG invite message if it was sent.
	if invite.TgMessageID > 0 {
		go func() {
			bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			invitedUser, err := s.store.GetUserByID(bgCtx, invite.InvitedUserID)
			if err == nil && invitedUser != nil {
				s.bot.DeleteMessage(bgCtx, invitedUser.TelegramID, invite.TgMessageID)
			}
		}()
	}

	return nil
}

// ── internal notification helpers ─────────────────────────────────────────────

func (s *Service) sendInviteNotification(invite *models.TripInvite, trip *models.Trip, inviterID, invitedUserID uuid.UUID) {
	bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	invitedUser, err := s.store.GetUserByID(bgCtx, invitedUserID)
	if err != nil || invitedUser == nil {
		return
	}
	if invitedUser.TelegramID == 0 {
		slog.Warn("telegram: invited user has no TelegramID", "user_id", invitedUserID)
		return
	}

	inviter, err := s.store.GetUserByID(bgCtx, inviterID)
	if err != nil || inviter == nil {
		return
	}

	inviterName := inviter.FirstName
	if inviterName == "" {
		inviterName = "@" + inviter.Username
	}

	startDate, endDate := "", ""
	if trip.StartDate != nil {
		startDate = trip.StartDate.Format("Jan 2, 2006")
	}
	if trip.EndDate != nil {
		endDate = trip.EndDate.Format("Jan 2, 2006")
	}

	miniAppURL := ""
	if s.cfg != nil {
		miniAppURL = s.cfg.MiniAppURL
	}
	if miniAppURL == "" {
		slog.Warn("telegram: MINI_APP_URL not set — invite notification will be plain text without WebApp button", "trip_id", trip.ID)
	}

	nctx, ncancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ncancel()

	msgID := s.bot.SendInviteNotification(nctx, invitedUser.TelegramID, inviterName, trip.Title, trip.City, startDate, endDate, miniAppURL)
	if msgID == 0 {
		slog.Error("telegram: invite notification not delivered", "invited_user_id", invitedUserID, "trip_id", trip.ID)
		return
	}
	// Best-effort: store message_id for later deletion.
	_ = s.store.UpdateInviteTgMessageID(context.Background(), invite.ID, msgID)
}

func (s *Service) handleInviteResponse(invite *models.TripInvite, responderID uuid.UUID, status string, comment *string, alternativeDate *string) {
	bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete the invite notification from the invited user's Telegram chat.
	if invite.TgMessageID > 0 {
		invitedUser, err := s.store.GetUserByID(bgCtx, invite.InvitedUserID)
		if err == nil && invitedUser != nil {
			s.bot.DeleteMessage(bgCtx, invitedUser.TelegramID, invite.TgMessageID)
		}
	}

	// Notify the trip owner.
	trip, err := s.store.GetTripByID(bgCtx, invite.TripID)
	if err != nil || trip == nil {
		return
	}
	// Don't notify the owner if they are the one responding (edge case).
	if trip.OwnerID == responderID {
		return
	}

	owner, err := s.store.GetUserByID(bgCtx, trip.OwnerID)
	if err != nil || owner == nil {
		return
	}
	responder, err := s.store.GetUserByID(bgCtx, responderID)
	if err != nil || responder == nil {
		return
	}

	responderName := responder.FirstName
	if responderName == "" {
		responderName = "@" + responder.Username
	}

	nctx, ncancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ncancel()

	if status == "accepted" {
		s.bot.SendInviteAccepted(nctx, owner.TelegramID, responderName, trip.Title)
		return
	}
	s.bot.SendInviteDeclined(nctx, owner.TelegramID, responderName, trip.Title, deref(comment), deref(alternativeDate))
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
