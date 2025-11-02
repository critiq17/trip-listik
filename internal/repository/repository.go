package repository

import (
	"log"

	"github.com/critiq17/tripListik/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetUserByUserID(telegramID int64) (*models.User, error) {
	var user models.User
	result := r.DB.Where("telegram_id = ?", telegramID).First(&user)
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

func (r *Repository) GetPlaceByNameAndUserID(userID uint, placeName string) (*models.Place, error) {
	var place models.Place
	result := r.DB.Where("user_id = ? AND name = ?", userID, placeName).First(&place)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &place, nil
}

func (r *Repository) GetPlacesByUserID(userID uint) ([]models.Place, error) {
	var places []models.Place
	result := r.DB.Where("user_id = ?", userID).Find(&places)
	return places, result.Error
}

func (r *Repository) GetOrCreateUser(telegramID int64, firstName, userName string) (models.User, error) {

	log.Printf("DEBUG: AddPlace called with telegramID=%d, firstName=%s, userName=%s", telegramID, firstName, userName)

	var user models.User

	result := r.DB.
		Where(models.User{TelegramID: telegramID}).
		Attrs(models.User{
			FirstName: firstName,
			UserName:  userName,
		}).
		FirstOrCreate(&user)
	return user, result.Error
}

func (r *Repository) AddPlaceToUser(userID uint, placeName string) error {
	newPlace := models.Place{
		Name:   placeName,
		UserID: userID,
	}

	return r.DB.Create(&newPlace).Error
}

func (r *Repository) DeletePlace(placeID uint) error {
	return r.DB.Delete(&models.Place{}, placeID).Error
}
