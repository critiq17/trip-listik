package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/httpapi/validate"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripsHandler struct {
	Store *store.Store
}

type createTripRequest struct {
	Title         string `json:"title" validate:"required"`
	Description   string `json:"description"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Visibility    string `json:"visibility"`
	Status        string `json:"status"`
	CountryCode   string `json:"country_code"`
	City          string `json:"city"`
	CoverPhotoURL string `json:"cover_photo_url"`
}

type updateTripRequest struct {
	Title         *string `json:"title"`
	Description   *string `json:"description"`
	StartDate     *string `json:"start_date"`
	EndDate       *string `json:"end_date"`
	Visibility    *string `json:"visibility"`
	Status        *string `json:"status"`
	CountryCode   *string `json:"country_code"`
	City          *string `json:"city"`
	CoverPhotoURL *string `json:"cover_photo_url"`
}

func (h *TripsHandler) CreateTrip(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var req createTripRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "title is required")
	}

	var startDate *time.Time
	var endDate *time.Time
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid start_date")
		}
		startDate = &parsed
	}
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid end_date")
		}
		endDate = &parsed
	}

	trip := &models.Trip{
		OwnerID:       userID,
		Title:         req.Title,
		Description:   req.Description,
		StartDate:     startDate,
		EndDate:       endDate,
		Visibility:    defaultString(req.Visibility, "public"),
		Status:        defaultString(req.Status, "draft"),
		CountryCode:   req.CountryCode,
		City:          req.City,
		CoverPhotoURL: req.CoverPhotoURL,
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	err := h.Store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(trip).Error; err != nil {
			return err
		}
		member := models.TripMember{TripID: trip.ID, UserID: userID, Role: "owner"}
		return tx.Create(&member).Error
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create trip")
	}

	return c.Status(fiber.StatusCreated).JSON(trip)
}

func (h *TripsHandler) GetTrip(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	trip, err := h.Store.GetTripByID(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if trip == nil {
		return fiber.NewError(fiber.StatusNotFound, "trip not found")
	}

	if trip.Visibility == "private" {
		userID, ok := middleware.GetUserID(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		isMember, err := h.Store.IsTripMember(ctx, trip.ID, userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to check membership")
		}
		if !isMember {
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		}
	}

	agg, err := h.Store.GetTripAggregates(ctx, trip.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load stats")
	}

	return c.JSON(fiber.Map{
		"trip":          trip,
		"member_count":  agg.MemberCount,
		"vote_count":    agg.VoteCount,
		"vote_average":  agg.VoteAverage,
		"comment_count": agg.CommentCount,
		"photo_count":   agg.PhotoCount,
	})
}

func (h *TripsHandler) UpdateTrip(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	var req updateTripRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	trip, err := h.Store.GetTripByID(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if trip == nil {
		return fiber.NewError(fiber.StatusNotFound, "trip not found")
	}
	if trip.OwnerID != userID {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	if req.Title != nil {
		trip.Title = *req.Title
	}
	if req.Description != nil {
		trip.Description = *req.Description
	}
	if req.Visibility != nil {
		trip.Visibility = *req.Visibility
	}
	if req.Status != nil {
		trip.Status = *req.Status
	}
	if req.CountryCode != nil {
		trip.CountryCode = *req.CountryCode
	}
	if req.City != nil {
		trip.City = *req.City
	}
	if req.CoverPhotoURL != nil {
		trip.CoverPhotoURL = *req.CoverPhotoURL
	}
	if req.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid start_date")
		}
		trip.StartDate = &parsed
	}
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid end_date")
		}
		trip.EndDate = &parsed
	}

	if err := h.Store.UpdateTrip(ctx, trip); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update trip")
	}

	return c.JSON(trip)
}

func (h *TripsHandler) ListMyTrips(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	scope := c.Query("scope", "mine")
	if scope != "mine" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid scope")
	}
	status := c.Query("status", "")

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	trips, err := h.Store.ListUserTripsByStatus(ctx, userID, status)
	if err != nil {
		if err == gorm.ErrInvalidData {
			return fiber.NewError(fiber.StatusBadRequest, "invalid status")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to list trips")
	}

	return c.JSON(fiber.Map{"items": trips})
}

func defaultString(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
