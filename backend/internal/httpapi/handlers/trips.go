package handlers

import (
	"errors"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/trips"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TripsHandler struct {
	Service *trips.Service
}

// ── request bodies ─────────────────────────────────────────────────────────────

type createTripRequest struct {
	Title         string `json:"title"`
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

// ── handlers ───────────────────────────────────────────────────────────────────

func (h *TripsHandler) CreateTrip(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var req createTripRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	trip, err := h.Service.CreateTrip(c.Context(), userID, trips.CreateTripInput{
		Title:         req.Title,
		Description:   req.Description,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Visibility:    req.Visibility,
		Status:        req.Status,
		CountryCode:   req.CountryCode,
		City:          req.City,
		CoverPhotoURL: req.CoverPhotoURL,
	})
	if err != nil {
		return mapTripError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(trip)
}

func (h *TripsHandler) GetTrip(c *fiber.Ctx) error {
	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	var viewerID *uuid.UUID
	if uid, ok := middleware.GetUserID(c); ok {
		viewerID = &uid
	}

	view, err := h.Service.GetTrip(c.Context(), tripID, viewerID)
	if err != nil {
		return mapTripError(err)
	}

	return c.JSON(view)
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

	trip, err := h.Service.UpdateTrip(c.Context(), userID, tripID, trips.UpdateTripInput{
		Title:         req.Title,
		Description:   req.Description,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Visibility:    req.Visibility,
		Status:        req.Status,
		CountryCode:   req.CountryCode,
		City:          req.City,
		CoverPhotoURL: req.CoverPhotoURL,
	})
	if err != nil {
		return mapTripError(err)
	}

	return c.JSON(trip)
}

func (h *TripsHandler) DeleteTrip(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	if err := h.Service.DeleteTrip(c.Context(), userID, tripID); err != nil {
		return mapTripError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
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

	items, err := h.Service.ListUserTrips(c.Context(), userID, status)
	if err != nil {
		if errors.Is(err, gorm.ErrInvalidData) {
			return fiber.NewError(fiber.StatusBadRequest, "invalid status")
		}
		return mapTripError(err)
	}

	return c.JSON(fiber.Map{"items": items})
}

// ── error mapping ──────────────────────────────────────────────────────────────

func mapTripError(err error) error {
	switch {
	case errors.Is(err, trips.ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, trips.ErrForbidden):
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	case errors.Is(err, trips.ErrValidation):
		// unwrap to get the validation message
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
}
