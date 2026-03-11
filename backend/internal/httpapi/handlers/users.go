package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
)

type UsersHandler struct {
	Store *store.Store
}

func (h *UsersHandler) Search(c *fiber.Ctx) error {
	_, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	query := c.Query("q", "")
	if len(query) < 2 {
		return fiber.NewError(fiber.StatusBadRequest, "query is too short")
	}
	limit := c.QueryInt("limit", 10)
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	users, err := h.Store.SearchUsers(ctx, query, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to search users")
	}

	return c.JSON(fiber.Map{"items": users})
}
