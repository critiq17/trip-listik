package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── Invite Links ──────────────────────────────────────────────────────────────
// Shareable multi-use trip invite tokens. Unlike trip_invites they are not
// bound to a user, so they can be sent to people who never opened the bot.

type InviteLink struct {
	Token     string     `gorm:"primaryKey" json:"token"`
	TripID    uuid.UUID  `gorm:"type:uuid" json:"trip_id"`
	CreatedBy uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func (InviteLink) TableName() string { return "invite_links" }

func newInviteToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate invite token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// GetOrCreateInviteLink returns the caller's active link for the trip,
// creating one on first use. One link per (trip, creator) keeps re-sharing
// idempotent and referral attribution unambiguous.
func (s *Store) GetOrCreateInviteLink(ctx context.Context, tripID, createdBy uuid.UUID) (*InviteLink, error) {
	var link InviteLink
	err := s.DB.WithContext(ctx).
		Where("trip_id = ? AND created_by = ? AND revoked_at IS NULL", tripID, createdBy).
		First(&link).Error
	if err == nil {
		return &link, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get invite link: %w", err)
	}

	token, err := newInviteToken()
	if err != nil {
		return nil, err
	}
	link = InviteLink{Token: token, TripID: tripID, CreatedBy: createdBy}
	if err := s.DB.WithContext(ctx).Create(&link).Error; err != nil {
		return nil, fmt.Errorf("create invite link: %w", err)
	}
	return &link, nil
}

// GetInviteLinkByToken returns an active (non-revoked) link or nil.
func (s *Store) GetInviteLinkByToken(ctx context.Context, token string) (*InviteLink, error) {
	var link InviteLink
	err := s.DB.WithContext(ctx).
		Where("token = ? AND revoked_at IS NULL", token).
		First(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get invite link by token: %w", err)
	}
	return &link, nil
}
