package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/realtime"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VotesHandler struct {
	Store *store.Store
	Hub   *realtime.Hub
}

type voteRequest struct {
	Vote int16 `json:"vote"`
}

func (h *VotesHandler) CastVote(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	_, allowed, err := ensureTripAccess(h.Store, tripID, &userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	var req voteRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Vote < 1 || req.Vote > 5 {
		return fiber.NewError(fiber.StatusBadRequest, "vote must be between 1 and 5")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.Store.UpsertVote(ctx, tripID, userID, req.Vote); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to cast vote")
	}

	summary, err := h.Store.GetVoteSummary(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load vote summary")
	}

	if h.Hub != nil {
		h.Hub.Publish("trip:"+tripID.String(), realtime.Event{
			Type: "vote_updated",
			Data: summary,
		})
	}

	return c.JSON(summary)
}

func (h *VotesHandler) GetVotes(c *fiber.Ctx) error {
	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	_, allowed, err := ensureTripAccess(h.Store, tripID, &userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	summary, err := h.Store.GetVoteSummary(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load vote summary")
	}

	return c.JSON(summary)
}
