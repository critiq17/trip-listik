package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/httpapi/validate"
	"github.com/critiq17/tripListik/internal/realtime"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommentsHandler struct {
	Store *store.Store
	Hub   *realtime.Hub
}

type createCommentRequest struct {
	Body string `json:"body" validate:"required"`
}

func (h *CommentsHandler) ListComments(c *fiber.Ctx) error {
	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	_, allowed, err := ensureTripAccess(c.Context(), h.Store, tripID, &userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
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

	comments, err := h.Store.ListCommentsWithUsers(ctx, tripID, limit, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load comments")
	}

	return c.JSON(fiber.Map{"items": comments})
}

func (h *CommentsHandler) CreateComment(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	_, allowed, err := ensureTripAccess(c.Context(), h.Store, tripID, &userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	var req createCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "body is required")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	comment := &models.TripComment{
		TripID: tripID,
		UserID: userID,
		Body:   req.Body,
	}
	if err := h.Store.CreateComment(ctx, comment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create comment")
	}

	user, err := h.Store.GetUserByID(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load user")
	}

	resp := fiber.Map{
		"id":         comment.ID,
		"trip_id":    comment.TripID,
		"user_id":    comment.UserID,
		"body":       comment.Body,
		"created_at": comment.CreatedAt,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"photo_url":  user.PhotoURL,
	}

	if h.Hub != nil {
		h.Hub.Publish("trip:"+tripID.String(), realtime.Event{
			Type: "comment_created",
			Data: resp,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
