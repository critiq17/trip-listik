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

func (s *Store) WithTx(tx *gorm.DB) *Store {
	if tx == nil {
		return s
	}
	return &Store{DB: tx}
}

// UpsertTelegramUser creates or refreshes a user by telegram_id. The second
// return value reports whether the user was created by this call, so callers
// can attribute referrals to first-time signups only.
func (s *Store) UpsertTelegramUser(ctx context.Context, u *models.User) (*models.User, bool, error) {
	var existed int64
	if err := s.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("telegram_id = ?", u.TelegramID).
		Count(&existed).Error; err != nil {
		return nil, false, err
	}

	if err := s.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "telegram_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"username", "first_name", "last_name", "photo_url", "updated_at"}),
		}).
		Create(u).Error; err != nil {
		return nil, false, err
	}

	stored, err := s.GetUserByTelegramID(ctx, u.TelegramID)
	if err != nil {
		return nil, false, err
	}
	return stored, existed == 0, nil
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

func (s *Store) UpdateUserProfile(ctx context.Context, userID uuid.UUID, bio string, isPublic bool) (*models.User, error) {
	if err := s.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"bio":       bio,
			"is_public": isPublic,
		}).Error; err != nil {
		return nil, err
	}
	return s.GetUserByID(ctx, userID)
}
