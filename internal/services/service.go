package services

import (
	"errors"

	"github.com/critiq17/tripListik/internal/models"
	"github.com/critiq17/tripListik/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddPlace(telegramID int64, firstName, userName, placeName string) error {

	user, err := s.repo.GetOrCreateUser(telegramID, firstName, userName)
	if err != nil {
		return err
	}

	return s.repo.AddPlaceToUser(user.ID, placeName)
}

func (s *Service) DeletePlace(telegramID int64, placeName string) error {
	user, err := s.repo.GetUserByUserID(telegramID)

	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	place, err := s.repo.GetPlaceByNameAndUserID(user.ID, placeName)
	if err != nil {
		return err
	}

	if place == nil {
		return errors.New("place not found")
	}
	return s.repo.DeletePlace(place.ID)
}

func (s *Service) GetUserPlaces(telegramID int64) ([]models.Place, error) {

	user, err := s.repo.GetUserByUserID(telegramID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	places, err := s.repo.GetPlacesByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return places, nil
}
