package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ── Pending Invites ───────────────────────────────────────────────────────────
// For users who haven't started the bot yet — delivered on first /start

type PendingInvite struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID      string    `gorm:"type:uuid" json:"trip_id"`
	TelegramID  int64     `json:"telegram_id"`
	InvitedBy   string    `gorm:"type:uuid" json:"invited_by"`
	CreatedAt   time.Time `json:"created_at"`
}

func (PendingInvite) TableName() string { return "pending_invites" }

func (s *Store) CreatePendingInvite(ctx context.Context, tripID string, invitedByUserID string, telegramID int64) error {
	pi := PendingInvite{
		TripID:     tripID,
		TelegramID: telegramID,
		InvitedBy:  invitedByUserID,
	}
	return s.DB.WithContext(ctx).Create(&pi).Error
}

type PendingInviteWithTrip struct {
	ID        string    `json:"id"`
	TripID    string    `json:"trip_id"`
	TripTitle string    `json:"trip_title"`
	InvitedBy string    `json:"invited_by"`
	InviterName string  `json:"inviter_name"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Store) GetAndClearPendingInvites(ctx context.Context, telegramID int64) ([]PendingInviteWithTrip, error) {
	var items []PendingInviteWithTrip
	err := s.DB.WithContext(ctx).
		Table("pending_invites pi").
		Select(`pi.id, pi.trip_id, pi.invited_by, pi.created_at,
			t.title AS trip_title,
			COALESCE(u.first_name || ' ' || u.last_name, u.username, 'Someone') AS inviter_name`).
		Joins("JOIN trips t ON t.id::text = pi.trip_id").
		Joins("JOIN users u ON u.id::text = pi.invited_by").
		Where("pi.telegram_id = ?", telegramID).
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return items, nil
	}
	// Delete after fetching
	if err := s.DB.WithContext(ctx).
		Exec("DELETE FROM pending_invites WHERE telegram_id = ?", telegramID).Error; err != nil {
		return items, err
	}
	return items, nil
}

// ── Referrals ────────────────────────────────────────────────────────────────

type Referral struct {
	ReferrerID string    `gorm:"type:uuid;primaryKey" json:"referrer_id"`
	ReferredID string    `gorm:"type:uuid;primaryKey" json:"referred_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Referral) TableName() string { return "referrals" }

func (s *Store) RecordReferral(ctx context.Context, referrerUserID, referredUserID string) error {
	return s.DB.WithContext(ctx).
		Exec("INSERT INTO referrals (referrer_id, referred_id, created_at) VALUES (?, ?, NOW()) ON CONFLICT DO NOTHING",
			referrerUserID, referredUserID).Error
}

// ReferralEntry is a user who joined via a referral link.
type ReferralEntry struct {
	Username string    `json:"username"`
	JoinedAt time.Time `json:"joined_at"`
}

// GetReferralEntries returns users who joined via the given referrer's link.
func (s *Store) GetReferralEntries(ctx context.Context, referrerID uuid.UUID) ([]ReferralEntry, error) {
	var entries []ReferralEntry
	err := s.DB.WithContext(ctx).
		Table("referrals r").
		Select("COALESCE(NULLIF(u.username, ''), u.first_name, 'User') AS username, r.created_at AS joined_at").
		Joins("JOIN users u ON u.id = r.referred_id").
		Where("r.referrer_id = ?", referrerID).
		Order("r.created_at DESC").
		Scan(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// GetReferralCount returns the total number of users referred by the given user.
func (s *Store) GetReferralCount(ctx context.Context, referrerID uuid.UUID) (int64, error) {
	var count int64
	err := s.DB.WithContext(ctx).
		Table("referrals").
		Where("referrer_id = ?", referrerID).
		Count(&count).Error
	return count, err
}
