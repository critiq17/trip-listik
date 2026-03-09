package store

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
)

func (s *Store) ListNotifications(ctx context.Context, userID uuid.UUID, limit int, cursor time.Time) ([]models.Notification, error) {
	var items []models.Notification
	q := s.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if !cursor.IsZero() {
		q = q.Where("created_at < ?", cursor)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
