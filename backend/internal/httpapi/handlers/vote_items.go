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

type VoteItemsHandler struct {
	Store *store.Store
	Hub   *realtime.Hub
}

type createVoteItemRequest struct {
	Category    string `json:"category" validate:"required,oneof=HOTEL AIRLINE ACTIVITY OTHER"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// GET /v1/trips/:id/vote-items
func (h *VoteItemsHandler) List(c *fiber.Ctx) error {
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

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	items, err := h.Store.ListVoteItems(ctx, tripID.String(), userID.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load vote items")
	}
	return c.JSON(fiber.Map{"items": items})
}

// POST /v1/trips/:id/vote-items
func (h *VoteItemsHandler) Create(c *fiber.Ctx) error {
	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	_, allowed, err := ensureTripAccess(ctx, h.Store, tripID, &userID)
	if err != nil || !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	var req createVoteItemRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "category and title are required")
	}

	item := &models.VoteItem{
		TripID:      tripID.String(),
		Category:    req.Category,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		AddedBy:     userID.String(),
	}
	if err := h.Store.CreateVoteItem(ctx, item); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create vote item")
	}

	if h.Hub != nil {
		h.Hub.Publish("trip:"+tripID.String(), realtime.Event{
			Type: "vote_item_created",
			Data: item,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// POST /v1/trips/:id/vote-items/:itemId/vote
func (h *VoteItemsHandler) Vote(c *fiber.Ctx) error {
	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}
	itemID := c.Params("itemId")
	if itemID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid item id")
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	_, allowed, err := ensureTripAccess(ctx, h.Store, tripID, &userID)
	if err != nil || !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	if err := h.Store.AddVoteItemVote(ctx, itemID, userID.String()); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to vote")
	}

	count, err := h.Store.GetVoteItemCount(ctx, itemID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get vote count")
	}

	result := fiber.Map{"item_id": itemID, "vote_count": count}
	if h.Hub != nil {
		h.Hub.Publish("trip:"+tripID.String(), realtime.Event{
			Type: "vote_item_updated",
			Data: result,
		})
	}

	return c.JSON(result)
}

// DELETE /v1/trips/:id/vote-items/:itemId/vote
func (h *VoteItemsHandler) Unvote(c *fiber.Ctx) error {
	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}
	itemID := c.Params("itemId")
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	_, allowed, err := ensureTripAccess(ctx, h.Store, tripID, &userID)
	if err != nil || !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	if err := h.Store.RemoveVoteItemVote(ctx, itemID, userID.String()); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to unvote")
	}

	count, err := h.Store.GetVoteItemCount(ctx, itemID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get vote count")
	}

	result := fiber.Map{"item_id": itemID, "vote_count": count}
	if h.Hub != nil {
		h.Hub.Publish("trip:"+tripID.String(), realtime.Event{
			Type: "vote_item_updated",
			Data: result,
		})
	}

	return c.JSON(result)
}
