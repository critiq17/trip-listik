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

type updateProfileRequest struct {
	Bio      string `json:"bio"`
	IsPublic *bool  `json:"is_public"`
}

func (h *ProfileHandler) Me(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
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

func (h *ProfileHandler) Update(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	user, err := h.Store.GetUserByID(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load user")
	}

	isPublic := user.IsPublic
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	updated, err := h.Store.UpdateUserProfile(ctx, userID, req.Bio, isPublic)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update profile")
	}

	return c.JSON(updated)
}

func (h *ProfileHandler) Stats(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
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

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
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

func (h *ProfileHandler) Map(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	visits, err := h.Store.GetUserCountryVisits(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load map data")
	}

	totalCountries := int64(len(visits))
	return c.JSON(fiber.Map{
		"countries":              visits,
		"total_countries":        totalCountries,
		"world_explored_percent": computeWorldPercent(totalCountries),
	})
}

func computeWorldPercent(countriesVisited int64) float64 {
	if countriesVisited == 0 {
		return 3.5
	}
	return (float64(countriesVisited) / 195.0) * 100.0
}

// Referrals returns the list of users who joined via this user's referral link.
func (h *ProfileHandler) Referrals(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	entries, err := h.Store.GetReferralEntries(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load referrals")
	}
	count, err := h.Store.GetReferralCount(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to count referrals")
	}

	return c.JSON(fiber.Map{
		"referral_count": count,
		"referrals":      entries,
	})
}

// Friends returns the list of users who share at least one trip with the authenticated user.
func (h *ProfileHandler) Friends(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	friends, err := h.Store.GetUserFriends(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load friends")
	}

	return c.JSON(fiber.Map{"items": friends})
}
