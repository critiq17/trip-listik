package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PublicProfileHandler struct {
	Store *store.Store
}

func (h *PublicProfileHandler) GetPublicProfile(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	user, err := h.Store.GetUserByID(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load user")
	}
	if user == nil || !user.IsPublic {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	stats, err := h.Store.ComputeUserStats(ctx, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load stats")
	}

	trips, err := h.Store.ListUserPublicTrips(ctx, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trips")
	}

	visits, err := h.Store.GetUserCountryVisits(ctx, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load map")
	}

	wishlist, err := h.Store.ListWishlist(ctx, user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load wishlist")
	}

	return c.JSON(fiber.Map{
		"user":                   user,
		"stats":                  stats,
		"world_explored_percent": computeWorldPercent(stats.CountriesVisited),
		"public_trips":           trips,
		"country_map":            visits,
		"wishlist":               wishlist,
	})
}
