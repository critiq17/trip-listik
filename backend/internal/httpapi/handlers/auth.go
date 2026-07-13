package handlers

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/critiq17/tripListik/internal/auth"
	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/httpapi/validate"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/critiq17/tripListik/internal/telegram"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Store *store.Store
	Cfg   *config.Config
	Bot   *telegram.Bot
}

type telegramAuthRequest struct {
	InitData   string `json:"initData"   validate:"required"`
	StartParam string `json:"start_param"`
}

type telegramAuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
}

type userResponse struct {
	ID         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PhotoURL   string `json:"photo_url"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type refreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	accessTokenTTL  = 2 * time.Hour
	refreshTokenTTL = 30 * 24 * time.Hour
)

func (h *AuthHandler) TelegramAuth(c *fiber.Ctx) error {
	var req telegramAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "initData is required")
	}

	tgUser, err := auth.ValidateTelegramInitData(req.InitData, h.Cfg.TelegramWebAppSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid telegram init data")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	user := &models.User{
		TelegramID: tgUser.ID,
		Username:   tgUser.Username,
		FirstName:  tgUser.FirstName,
		LastName:   tgUser.LastName,
		PhotoURL:   tgUser.PhotoURL,
	}

	stored, created, err := h.Store.UpsertTelegramUser(ctx, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to upsert user")
	}

	// Referral attribution: only first-time signups count, and the deep-link
	// payload comes from the signed initData (client body is a fallback for
	// older Telegram clients that omit start_param there).
	startParam := tgUser.StartParam
	if startParam == "" {
		startParam = req.StartParam
	}
	if created && startParam != "" {
		go h.recordReferral(startParam, stored)
	}

	token, err := auth.NewToken(stored.ID, stored.TelegramID, h.Cfg.JWTSecret, accessTokenTTL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	refreshToken, err := h.Store.CreateRefreshToken(ctx, stored.ID, refreshTokenTTL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate refresh token")
	}

	resp := telegramAuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: userResponse{
			ID:         stored.ID.String(),
			TelegramID: stored.TelegramID,
			Username:   stored.Username,
			FirstName:  stored.FirstName,
			LastName:   stored.LastName,
			PhotoURL:   stored.PhotoURL,
		},
	}

	return c.JSON(resp)
}

// recordReferral resolves a startapp deep-link payload to a referrer and
// credits them for the new signup. Supported payloads: profile_<uuid>
// (profile share link) and join_<token> (trip invite link — the link
// creator is the referrer). Runs in a goroutine: failures only cost the
// referral notification, never the login.
func (h *AuthHandler) recordReferral(startParam string, newUser *models.User) {
	rCtx, rCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer rCancel()

	var referrerID uuid.UUID
	switch {
	case strings.HasPrefix(startParam, "profile_"):
		id, err := uuid.Parse(strings.TrimPrefix(startParam, "profile_"))
		if err != nil {
			return
		}
		referrerID = id
	case strings.HasPrefix(startParam, "join_"):
		link, err := h.Store.GetInviteLinkByToken(rCtx, strings.TrimPrefix(startParam, "join_"))
		if err != nil || link == nil {
			return
		}
		referrerID = link.CreatedBy
	default:
		return
	}

	if referrerID == newUser.ID {
		return
	}

	if err := h.Store.RecordReferral(rCtx, referrerID.String(), newUser.ID.String()); err != nil {
		return // referral already recorded or error — skip notification
	}

	referrer, err := h.Store.GetUserByID(rCtx, referrerID)
	if err != nil || referrer == nil || h.Bot == nil {
		return
	}
	if referrer.TelegramID == 0 {
		slog.Warn("referral: referrer has no TelegramID", "referrer_id", referrerID)
		return
	}
	count, err := h.Store.GetReferralCount(rCtx, referrerID)
	if err != nil {
		return
	}

	newUserName := newUser.Username
	if newUserName == "" {
		newUserName = newUser.FirstName
	}
	nCtx, nCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer nCancel()
	h.Bot.SendReferralNotification(nCtx, referrer.TelegramID, newUserName, count)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "refresh_token is required")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	userID, newRefreshToken, err := h.Store.RotateRefreshToken(ctx, req.RefreshToken, refreshTokenTTL)
	if err != nil {
		if err == store.ErrInvalidRefreshToken {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to refresh token")
	}

	user, err := h.Store.GetUserByID(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load user")
	}

	token, err := auth.NewToken(user.ID, user.TelegramID, h.Cfg.JWTSecret, accessTokenTTL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(refreshResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
	})
}
