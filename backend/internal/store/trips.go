package store

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Store) CreateTrip(ctx context.Context, trip *models.Trip) error {
	return s.DB.WithContext(ctx).Create(trip).Error
}

func (s *Store) UpdateTrip(ctx context.Context, trip *models.Trip) error {
	return s.DB.WithContext(ctx).Save(trip).Error
}

func (s *Store) GetTripByID(ctx context.Context, id uuid.UUID) (*models.Trip, error) {
	var trip models.Trip
	if err := s.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&trip).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &trip, nil
}

func (s *Store) ListPublicTrips(ctx context.Context, limit int, cursor time.Time) ([]models.Trip, error) {
	var trips []models.Trip
	q := s.DB.WithContext(ctx).
		Where("visibility = ? AND deleted_at IS NULL", "public").
		Order("created_at DESC").
		Limit(limit)

	if !cursor.IsZero() {
		q = q.Where("created_at < ?", cursor)
	}

	if err := q.Find(&trips).Error; err != nil {
		return nil, err
	}
	return trips, nil
}

func (s *Store) ListUserTrips(ctx context.Context, userID uuid.UUID) ([]models.Trip, error) {
	var trips []models.Trip
	err := s.DB.WithContext(ctx).
		Joins("JOIN trip_members tm ON tm.trip_id = trips.id").
		Where("tm.user_id = ? AND trips.deleted_at IS NULL", userID).
		Order("trips.start_date DESC NULLS LAST").
		Find(&trips).Error
	if err != nil {
		return nil, err
	}
	return trips, nil
}

func (s *Store) ListFriendTrips(ctx context.Context, userID uuid.UUID, limit int, cursor time.Time) ([]models.Trip, error) {
	var trips []models.Trip
	q := s.DB.WithContext(ctx).
		Joins("JOIN trip_members tm ON tm.trip_id = trips.id").
		Where("tm.user_id = ? AND trips.deleted_at IS NULL", userID).
		Order("trips.created_at DESC").
		Limit(limit)

	if !cursor.IsZero() {
		q = q.Where("trips.created_at < ?", cursor)
	}

	if err := q.Find(&trips).Error; err != nil {
		return nil, err
	}
	return trips, nil
}

func (s *Store) ListPopularTrips(ctx context.Context, limit int, cursor time.Time) ([]models.Trip, error) {
	var trips []models.Trip
	q := s.DB.WithContext(ctx).
		Select("trips.*").
		Joins("LEFT JOIN trip_votes tv ON tv.trip_id = trips.id").
		Where("trips.visibility = ? AND trips.deleted_at IS NULL", "public").
		Group("trips.id").
		Order("COUNT(tv.id) DESC, trips.created_at DESC").
		Limit(limit)

	if !cursor.IsZero() {
		q = q.Where("trips.created_at < ?", cursor)
	}

	if err := q.Find(&trips).Error; err != nil {
		return nil, err
	}
	return trips, nil
}

func (s *Store) ListNearbyTrips(ctx context.Context, countryCode string, limit int, cursor time.Time) ([]models.Trip, error) {
	var trips []models.Trip
	q := s.DB.WithContext(ctx).
		Where("visibility = ? AND country_code = ? AND deleted_at IS NULL", "public", countryCode).
		Order("created_at DESC").
		Limit(limit)

	if !cursor.IsZero() {
		q = q.Where("created_at < ?", cursor)
	}

	if err := q.Find(&trips).Error; err != nil {
		return nil, err
	}
	return trips, nil
}

func (s *Store) SearchPublicTrips(ctx context.Context, query, countryCode string, limit int, cursor time.Time) ([]models.Trip, error) {
	var trips []models.Trip
	q := s.DB.WithContext(ctx).
		Where("visibility = ? AND deleted_at IS NULL", "public")

	if query != "" {
		like := "%" + query + "%"
		q = q.Where("title ILIKE ? OR city ILIKE ?", like, like)
	}
	if countryCode != "" {
		q = q.Where("country_code = ?", countryCode)
	}

	q = q.Order("created_at DESC").Limit(limit)
	if !cursor.IsZero() {
		q = q.Where("created_at < ?", cursor)
	}

	if err := q.Find(&trips).Error; err != nil {
		return nil, err
	}
	return trips, nil
}
