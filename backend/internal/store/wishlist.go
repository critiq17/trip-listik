package store

import (
	"context"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
)

func (s *Store) ListWishlist(ctx context.Context, userID uuid.UUID) ([]models.UserWishlistItem, error) {
	var items []models.UserWishlistItem
	if err := s.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Store) AddWishlistItem(ctx context.Context, item *models.UserWishlistItem) error {
	return s.DB.WithContext(ctx).Create(item).Error
}

func (s *Store) DeleteWishlistItem(ctx context.Context, userID, itemID uuid.UUID) error {
	return s.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", itemID, userID).
		Delete(&models.UserWishlistItem{}).Error
}
