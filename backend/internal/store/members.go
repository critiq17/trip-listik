package store

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Store) GetTripMembers(ctx context.Context, tripID uuid.UUID) ([]models.TripMember, error) {
	var members []models.TripMember
	if err := s.DB.WithContext(ctx).Where("trip_id = ?", tripID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

type TripMemberView struct {
	TripID    uuid.UUID `json:"trip_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PhotoURL  string    `json:"photo_url"`
}

func (s *Store) GetTripMembersWithUsers(ctx context.Context, tripID uuid.UUID) ([]TripMemberView, error) {
	var members []TripMemberView
	err := s.DB.WithContext(ctx).
		Table("trip_members tm").
		Select(`tm.trip_id, tm.user_id, tm.role, tm.joined_at, u.username, u.first_name, u.last_name, u.photo_url`).
		Joins("JOIN users u ON u.id = tm.user_id").
		Where("tm.trip_id = ?", tripID).
		Order("tm.joined_at ASC").
		Scan(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s *Store) GetTripMembersWithUsersPaged(ctx context.Context, tripID uuid.UUID, limit int, cursor time.Time) ([]TripMemberView, error) {
	var members []TripMemberView
	q := s.DB.WithContext(ctx).
		Table("trip_members tm").
		Select(`tm.trip_id, tm.user_id, tm.role, tm.joined_at, u.username, u.first_name, u.last_name, u.photo_url`).
		Joins("JOIN users u ON u.id = tm.user_id").
		Where("tm.trip_id = ?", tripID)
	if !cursor.IsZero() {
		q = q.Where("tm.joined_at > ?", cursor)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Order("tm.joined_at ASC").Scan(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (s *Store) IsTripMember(ctx context.Context, tripID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.TripMember{}).
		Where("trip_id = ? AND user_id = ?", tripID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) AddTripMember(ctx context.Context, tripID, userID uuid.UUID, role string) error {
	member := models.TripMember{TripID: tripID, UserID: userID, Role: role}
	return s.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "trip_id"}, {Name: "user_id"}},
			DoNothing: true,
		}).
		Create(&member).Error
}

func (s *Store) RemoveTripMember(ctx context.Context, tripID, userID uuid.UUID) error {
	return s.DB.WithContext(ctx).Where("trip_id = ? AND user_id = ?", tripID, userID).Delete(&models.TripMember{}).Error
}

func (s *Store) CreateJoinRequest(ctx context.Context, tripID, userID uuid.UUID) error {
	jr := models.TripJoinRequest{TripID: tripID, UserID: userID, Status: "pending"}
	return s.DB.WithContext(ctx).Create(&jr).Error
}

func (s *Store) UpdateJoinRequestStatus(ctx context.Context, tripID, userID uuid.UUID, status string) error {
	return s.DB.WithContext(ctx).
		Model(&models.TripJoinRequest{}).
		Where("trip_id = ? AND user_id = ?", tripID, userID).
		Updates(map[string]any{"status": status}).Error
}

func (s *Store) GetJoinRequest(ctx context.Context, tripID, userID uuid.UUID) (*models.TripJoinRequest, error) {
	var jr models.TripJoinRequest
	if err := s.DB.WithContext(ctx).Where("trip_id = ? AND user_id = ?", tripID, userID).First(&jr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jr, nil
}

func (s *Store) ListJoinRequests(ctx context.Context, tripID uuid.UUID, status string) ([]models.TripJoinRequest, error) {
	var items []models.TripJoinRequest
	q := s.DB.WithContext(ctx).Where("trip_id = ?", tripID)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

type TripJoinRequestView struct {
	TripID    uuid.UUID `json:"trip_id"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PhotoURL  string    `json:"photo_url"`
}

func (s *Store) ListJoinRequestsWithUsers(ctx context.Context, tripID uuid.UUID, status string) ([]TripJoinRequestView, error) {
	var items []TripJoinRequestView
	q := s.DB.WithContext(ctx).
		Table("trip_join_requests tjr").
		Select(`tjr.trip_id, tjr.user_id, tjr.status, tjr.created_at, u.username, u.first_name, u.last_name, u.photo_url`).
		Joins("JOIN users u ON u.id = tjr.user_id").
		Where("tjr.trip_id = ?", tripID)
	if status != "" {
		q = q.Where("tjr.status = ?", status)
	}
	if err := q.Order("tjr.created_at DESC").Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
