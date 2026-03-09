package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	Store *store.Store
}

func (h *ProfileHandler) Me(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user, err := h.Store.GetUserByID(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load user")
	}

	stats, err := h.Store.ComputeUserStats(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load stats")
	}

	worldPercent := computeWorldPercent(stats.CountriesVisited)

	return c.JSON(fiber.Map{
		"user":                   user,
		"stats":                  stats,
		"world_explored_percent": worldPercent,
	})
}

func (h *ProfileHandler) Stats(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stats, err := h.Store.ComputeUserStats(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load stats")
	}

	return c.JSON(stats)
}

func (h *ProfileHandler) World(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stats, err := h.Store.ComputeUserStats(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load stats")
	}

	return c.JSON(fiber.Map{
		"countries_visited":      stats.CountriesVisited,
		"world_explored_percent": computeWorldPercent(stats.CountriesVisited),
	})
}

func computeWorldPercent(countriesVisited int64) float64 {
	if countriesVisited == 0 {
		return 3.5
	}
	return (float64(countriesVisited) / 195.0) * 100.0
}
