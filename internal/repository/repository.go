package repository

import (
	"github.com/critiq17/tripListik/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetUserByUserID(userID int64) (*models.User, error) {
	var user models.User
	result := r.DB.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *Repository) AddPlaceToUser(userID uint, placeName string) error {
	newPlace := models.Place{
		Name:   placeName,
		UserID: userID,
	}

	return r.DB.Create(&newPlace).Error
}
