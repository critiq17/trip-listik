package services

import (
	"github.com/critiq17/tripListik/internal/models"
	"github.com/critiq17/tripListik/internal/repository"
)

type Service struct {
	Repo repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Repo: *repo}
}

func (s *Service) AddPlace(telegramID int64, firstName, userName, placeName string) error {
	user, err := s.Repo.GetUserByUserID(telegramID)
	if err != nil {
		return err
	}

	if user == nil {
		user = &models.User{
			TelegramID: telegramID,
			FirstName:  firstName,
			UserName:   userName,
		}

		if err := s.Repo.CreateUser(user); err != nil {
			return err
		}
	}

	return s.Repo.AddPlaceToUser(user.ID, placeName)
}
