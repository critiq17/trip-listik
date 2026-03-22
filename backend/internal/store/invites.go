package store

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InviteView struct {
	ID              uuid.UUID `json:"id"`
	TripID          uuid.UUID `json:"trip_id"`
	InvitedUserID   uuid.UUID `json:"invited_user_id"`
	InvitedByUserID uuid.UUID `json:"invited_by_user_id"`
	Status          string    `json:"status"`
	Comment         *string   `json:"comment"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	TripTitle       string    `json:"trip_title"`
	TripCover       string    `json:"trip_cover_photo_url"`
	InviterUsername string    `json:"inviter_username"`
	InviterFirst    string    `json:"inviter_first_name"`
	InviterLast     string    `json:"inviter_last_name"`
	InviterPhoto    string    `json:"inviter_photo_url"`
}

func (s *Store) CreateInvite(ctx context.Context, tripID, invitedUserID, invitedByUserID uuid.UUID) (*models.TripInvite, error) {
	invite := &models.TripInvite{
		TripID:          tripID,
		InvitedUserID:   invitedUserID,
		InvitedByUserID: invitedByUserID,
		Status:          "pending",
	}
	if err := s.DB.WithContext(ctx).Create(invite).Error; err != nil {
		return nil, err
	}
	return invite, nil
}

func (s *Store) GetInviteByID(ctx context.Context, inviteID uuid.UUID) (*models.TripInvite, error) {
	var invite models.TripInvite
	if err := s.DB.WithContext(ctx).Where("id = ?", inviteID).First(&invite).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &invite, nil
}

func (s *Store) GetInviteByTripAndUser(ctx context.Context, tripID, invitedUserID uuid.UUID) (*models.TripInvite, error) {
	var invite models.TripInvite
	if err := s.DB.WithContext(ctx).
		Where("trip_id = ? AND invited_user_id = ?", tripID, invitedUserID).
		First(&invite).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &invite, nil
}

func (s *Store) ListTripInvites(ctx context.Context, tripID uuid.UUID) ([]InviteView, error) {
	var items []InviteView
	err := s.DB.WithContext(ctx).
		Table("trip_invites ti").
		Select(`ti.id, ti.trip_id, ti.invited_user_id, ti.invited_by_user_id, ti.status, ti.comment,
			ti.created_at, ti.updated_at,
			t.title as trip_title, t.cover_photo_url as trip_cover_photo_url,
			u.username as inviter_username, u.first_name as inviter_first_name, u.last_name as inviter_last_name, u.photo_url as inviter_photo_url`).
		Joins("JOIN trips t ON t.id = ti.trip_id").
		Joins("JOIN users u ON u.id = ti.invited_by_user_id").
		Where("ti.trip_id = ?", tripID).
		Order("ti.created_at DESC").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Store) ListUserInvites(ctx context.Context, userID uuid.UUID, limit int) ([]InviteView, error) {
	var items []InviteView
	q := s.DB.WithContext(ctx).
		Table("trip_invites ti").
		Select(`ti.id, ti.trip_id, ti.invited_user_id, ti.invited_by_user_id, ti.status, ti.comment,
			ti.created_at, ti.updated_at,
			t.title as trip_title, t.cover_photo_url as trip_cover_photo_url,
			u.username as inviter_username, u.first_name as inviter_first_name, u.last_name as inviter_last_name, u.photo_url as inviter_photo_url`).
		Joins("JOIN trips t ON t.id = ti.trip_id").
		Joins("JOIN users u ON u.id = ti.invited_by_user_id").
		Where("ti.invited_user_id = ?", userID).
		Order("ti.created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Store) UpdateInviteTgMessageID(ctx context.Context, inviteID uuid.UUID, msgID int64) error {
	return s.DB.WithContext(ctx).
		Model(&models.TripInvite{}).
		Where("id = ?", inviteID).
		Update("tg_message_id", msgID).Error
}

func (s *Store) CancelInvite(ctx context.Context, inviteID uuid.UUID) error {
	return s.DB.WithContext(ctx).
		Model(&models.TripInvite{}).
		Where("id = ?", inviteID).
		Update("status", "cancelled").Error
}

func (s *Store) UpdateInviteStatus(ctx context.Context, inviteID uuid.UUID, status string, comment *string) error {
	updates := map[string]any{
		"status": status,
	}
	if comment != nil {
		updates["comment"] = *comment
	}
	return s.DB.WithContext(ctx).
		Model(&models.TripInvite{}).
		Where("id = ?", inviteID).
		Updates(updates).Error
}
