package trips

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sentinel errors — mapped to HTTP status codes in the handler layer.
var (
	ErrNotFound   = errors.New("not found")
	ErrForbidden  = errors.New("forbidden")
	ErrValidation = errors.New("validation error")
)

// Service encapsulates all trips business logic.
type Service struct {
	store *store.Store
}

func NewService(s *store.Store) *Service {
	return &Service{store: s}
}

// ── Input types ────────────────────────────────────────────────────────────────

type CreateTripInput struct {
	Title         string
	Description   string
	StartDate     string
	EndDate       string
	Visibility    string
	Status        string
	CountryCode   string
	City          string
	CoverPhotoURL string
}

type UpdateTripInput struct {
	Title         *string
	Description   *string
	StartDate     *string
	EndDate       *string
	Visibility    *string
	Status        *string
	CountryCode   *string
	City          *string
	CoverPhotoURL *string
}

// TripView bundles a trip with its aggregate stats for the detail endpoint.
type TripView struct {
	Trip           *models.Trip `json:"trip"`
	MemberCount    int64        `json:"member_count"`
	VoteCount      int64        `json:"vote_count"`
	VoteAverage    float64      `json:"vote_average"`
	CommentCount   int64        `json:"comment_count"`
	PhotoCount     int64        `json:"photo_count"`
	ViewerIsMember bool         `json:"viewer_is_member"`
}

// ── CreateTrip ────────────────────────────────────────────────────────────────

func (s *Service) CreateTrip(ctx context.Context, ownerID uuid.UUID, in CreateTripInput) (*models.Trip, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrValidation)
	}
	if len(title) > 100 {
		return nil, fmt.Errorf("%w: title must be at most 100 characters", ErrValidation)
	}
	if len(in.Description) > 2000 {
		return nil, fmt.Errorf("%w: description must be at most 2000 characters", ErrValidation)
	}

	startDate, endDate, err := parseDateRange(in.StartDate, in.EndDate)
	if err != nil {
		return nil, err
	}

	trip := &models.Trip{
		OwnerID:       ownerID,
		Title:         title,
		Description:   in.Description,
		StartDate:     startDate,
		EndDate:       endDate,
		Visibility:    defaultStr(in.Visibility, "public"),
		Status:        defaultStr(in.Status, "draft"),
		CountryCode:   in.CountryCode,
		City:          in.City,
		CoverPhotoURL: in.CoverPhotoURL,
	}

	err = s.store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trip).Error; err != nil {
			return err
		}
		return tx.Create(&models.TripMember{TripID: trip.ID, UserID: ownerID, Role: "owner"}).Error
	})
	if err != nil {
		return nil, fmt.Errorf("trips.CreateTrip: %w", err)
	}
	return trip, nil
}

// ── UpdateTrip ────────────────────────────────────────────────────────────────

func (s *Service) UpdateTrip(ctx context.Context, requesterID, tripID uuid.UUID, in UpdateTripInput) (*models.Trip, error) {
	trip, err := s.store.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("trips.UpdateTrip get: %w", err)
	}
	if trip == nil {
		return nil, ErrNotFound
	}
	if trip.OwnerID != requesterID {
		return nil, ErrForbidden
	}

	if in.Title != nil {
		t := strings.TrimSpace(*in.Title)
		if t == "" {
			return nil, fmt.Errorf("%w: title cannot be empty", ErrValidation)
		}
		if len(t) > 100 {
			return nil, fmt.Errorf("%w: title must be at most 100 characters", ErrValidation)
		}
		trip.Title = t
	}
	if in.Description != nil {
		if len(*in.Description) > 2000 {
			return nil, fmt.Errorf("%w: description must be at most 2000 characters", ErrValidation)
		}
		trip.Description = *in.Description
	}
	if in.Visibility != nil {
		trip.Visibility = *in.Visibility
	}
	if in.Status != nil {
		trip.Status = *in.Status
	}
	if in.CountryCode != nil {
		trip.CountryCode = *in.CountryCode
	}
	if in.City != nil {
		trip.City = *in.City
	}
	if in.CoverPhotoURL != nil {
		trip.CoverPhotoURL = *in.CoverPhotoURL
	}
	if in.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *in.StartDate)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid start_date", ErrValidation)
		}
		trip.StartDate = &parsed
	}
	if in.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *in.EndDate)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid end_date", ErrValidation)
		}
		trip.EndDate = &parsed
	}

	if trip.StartDate != nil && trip.EndDate != nil && trip.EndDate.Before(*trip.StartDate) {
		return nil, fmt.Errorf("%w: end_date must be on or after start_date", ErrValidation)
	}

	if err := s.store.UpdateTrip(ctx, trip); err != nil {
		return nil, fmt.Errorf("trips.UpdateTrip save: %w", err)
	}
	return trip, nil
}

// ── DeleteTrip ────────────────────────────────────────────────────────────────

func (s *Service) DeleteTrip(ctx context.Context, requesterID, tripID uuid.UUID) error {
	trip, err := s.store.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("trips.DeleteTrip get: %w", err)
	}
	if trip == nil {
		return ErrNotFound
	}
	if trip.OwnerID != requesterID {
		return ErrForbidden
	}
	if err := s.store.DeleteTrip(ctx, tripID); err != nil {
		return fmt.Errorf("trips.DeleteTrip: %w", err)
	}
	return nil
}

// ── GetTrip ───────────────────────────────────────────────────────────────────

// GetTrip returns enriched trip detail. viewerID is nil for unauthenticated requests.
func (s *Service) GetTrip(ctx context.Context, tripID uuid.UUID, viewerID *uuid.UUID) (*TripView, error) {
	trip, err := s.store.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("trips.GetTrip: %w", err)
	}
	if trip == nil {
		return nil, ErrNotFound
	}

	if trip.Visibility == "private" {
		if viewerID == nil {
			return nil, ErrForbidden
		}
		isMember, err := s.store.IsTripMember(ctx, tripID, *viewerID)
		if err != nil {
			return nil, fmt.Errorf("trips.GetTrip member check: %w", err)
		}
		if !isMember {
			return nil, ErrForbidden
		}
	}

	agg, err := s.store.GetTripAggregates(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("trips.GetTrip aggregates: %w", err)
	}

	viewerIsMember := false
	if viewerID != nil {
		viewerIsMember, _ = s.store.IsTripMember(ctx, tripID, *viewerID)
	}

	return &TripView{
		Trip:           trip,
		MemberCount:    agg.MemberCount,
		VoteCount:      agg.VoteCount,
		VoteAverage:    agg.VoteAverage,
		CommentCount:   agg.CommentCount,
		PhotoCount:     agg.PhotoCount,
		ViewerIsMember: viewerIsMember,
	}, nil
}

// ── ListUserTrips ─────────────────────────────────────────────────────────────

func (s *Service) ListUserTrips(ctx context.Context, userID uuid.UUID, status string) ([]models.Trip, error) {
	trips, err := s.store.ListUserTripsByStatus(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("trips.ListUserTrips: %w", err)
	}
	return trips, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func parseDateRange(startStr, endStr string) (*time.Time, *time.Time, error) {
	var start, end *time.Time
	if startStr != "" {
		t, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: invalid start_date", ErrValidation)
		}
		start = &t
	}
	if endStr != "" {
		t, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: invalid end_date", ErrValidation)
		}
		end = &t
	}
	if start != nil && end != nil && end.Before(*start) {
		return nil, nil, fmt.Errorf("%w: end_date must be on or after start_date", ErrValidation)
	}
	return start, end, nil
}

func defaultStr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
