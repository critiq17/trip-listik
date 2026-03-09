package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
)

type InboxHandler struct {
	Store *store.Store
}

func (h *InboxHandler) ListInbox(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	limit := c.QueryInt("limit", 20)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	cursorStr := c.Query("cursor", "")
	var cursor time.Time
	if cursorStr != "" {
		parsed, err := time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid cursor")
		}
		cursor = parsed
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	items, err := h.Store.ListNotifications(ctx, userID, limit, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load inbox")
	}

	return c.JSON(fiber.Map{"items": items})
}
