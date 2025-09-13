package services

import (
	"errors"

	"github.com/critiq17/tripListik/internal/models"
	"github.com/critiq17/tripListik/internal/repository"
)

type Service struct {
	Repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) AddPlace(telegramID int64, firstName, userName, placeName string) error {
	var user models.User

	// get user for TgID or create if user not found
	result := s.Repo.DB.Where("telegram_id = ?", telegramID).Attrs(models.User{
		FirstName: firstName,
		UserName:  userName,
	}).FirstOrCreate(&user)

	if result.Error != nil {
		return result.Error
	}

	return s.Repo.AddPlaceToUser(user.ID, placeName)
}

func (s *Service) DeletePlace(telegramID int64, placeName string) error {
	user, err := s.Repo.GetUserByUserID(telegramID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// place, err := s.Repo.GetPlaceByNameAndUserID( placeName)
	if err != nil {
		return err
	}

	//if place == nil {
	//		return errors.New("place not found")
	//	}
	//return s.Repo.DeletePlace(place.ID)

	return nil
}
