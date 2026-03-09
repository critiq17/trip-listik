package store

import (
	"context"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) UpsertTelegramUser(ctx context.Context, u *models.User) (*models.User, error) {
	if err := s.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "telegram_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"username", "first_name", "last_name", "photo_url", "updated_at"}),
		}).
		Create(u).Error; err != nil {
		return nil, err
	}

	return s.GetUserByTelegramID(ctx, u.TelegramID)
}

func (s *Store) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	var user models.User
	if err := s.DB.WithContext(ctx).Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
