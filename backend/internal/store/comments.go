package store

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
)

func (s *Store) CreateComment(ctx context.Context, cmt *models.TripComment) error {
	return s.DB.WithContext(ctx).Create(cmt).Error
}

type TripCommentView struct {
	ID        uuid.UUID `json:"id"`
	TripID    uuid.UUID `json:"trip_id"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PhotoURL  string    `json:"photo_url"`
}

func (s *Store) ListCommentsWithUsers(ctx context.Context, tripID uuid.UUID, limit int, cursor time.Time) ([]TripCommentView, error) {
	var comments []TripCommentView
	q := s.DB.WithContext(ctx).
		Table("trip_comments tc").
		Select(`tc.id, tc.trip_id, tc.user_id, tc.body, tc.created_at, u.username, u.first_name, u.last_name, u.photo_url`).
		Joins("JOIN users u ON u.id = tc.user_id").
		Where("tc.trip_id = ?", tripID).
		Order("tc.created_at DESC").
		Limit(limit)

	if !cursor.IsZero() {
		q = q.Where("tc.created_at < ?", cursor)
	}

	if err := q.Scan(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
