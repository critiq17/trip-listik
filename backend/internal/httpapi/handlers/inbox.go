package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
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

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	items, err := h.Store.ListNotifications(ctx, userID, limit, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load inbox")
	}

	invites, err := h.Store.ListUserInvites(ctx, userID, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load invites")
	}

	return c.JSON(fiber.Map{"items": items, "invites": invites})
}

// UnreadCount returns the number of unread notifications for badge polling
func (h *InboxHandler) UnreadCount(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 1*time.Second)
	defer cancel()

	var count int64
	err := h.Store.DB.WithContext(ctx).
		Table("notifications").
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to count notifications")
	}

	return c.JSON(fiber.Map{"unread": count})
}

// MarkAllRead marks all unread notifications as read for the current user.
func (h *InboxHandler) MarkAllRead(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	now := time.Now()
	err := h.Store.DB.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", now).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to mark read")
	}
	return c.JSON(fiber.Map{"ok": true})
}
