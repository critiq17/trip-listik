package store

import (
	"context"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
)

func (s *Store) CreatePhoto(ctx context.Context, p *models.TripPhoto) error {
	return s.DB.WithContext(ctx).Create(p).Error
}

func (s *Store) ListPhotos(ctx context.Context, tripID uuid.UUID, limit int) ([]models.TripPhoto, error) {
	var photos []models.TripPhoto
	q := s.DB.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (s *Store) DeletePhoto(ctx context.Context, photoID uuid.UUID, userID uuid.UUID) error {
	return s.DB.WithContext(ctx).Where("id = ? AND user_id = ?", photoID, userID).Delete(&models.TripPhoto{}).Error
}
