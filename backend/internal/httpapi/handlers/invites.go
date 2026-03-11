package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/httpapi/validate"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InvitesHandler struct {
	Store *store.Store
}

type inviteRequest struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

type inviteResponse struct {
	InviteID string `json:"invite_id"`
	Status   string `json:"status"`
}

type respondInviteRequest struct {
	Action  string  `json:"action" validate:"required,oneof=accept decline"`
	Comment *string `json:"comment"`
}

func (h *InvitesHandler) InviteUser(c *fiber.Ctx) error {
	ownerID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	var req inviteRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.Username == "" && req.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username or user_id is required")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	trip, err := h.Store.GetTripByID(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if trip == nil {
		return fiber.NewError(fiber.StatusNotFound, "trip not found")
	}
	if trip.OwnerID != ownerID {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	var invitedUserID uuid.UUID
	if req.UserID != "" {
		invitedUserID, err = uuid.Parse(req.UserID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid user_id")
		}
	} else {
		users, err := h.Store.SearchUsers(ctx, req.Username, 1)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to search user")
		}
		if len(users) == 0 {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		invitedUserID = users[0].ID
	}

	existing, err := h.Store.GetInviteByTripAndUser(ctx, tripID, invitedUserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to check invites")
	}
	if existing != nil {
		return c.JSON(inviteResponse{InviteID: existing.ID.String(), Status: existing.Status})
	}

	invite, err := h.Store.CreateInvite(ctx, tripID, invitedUserID, ownerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create invite")
	}

	return c.JSON(inviteResponse{InviteID: invite.ID.String(), Status: invite.Status})
}

func (h *InvitesHandler) ListTripInvites(c *fiber.Ctx) error {
	ownerID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	trip, err := h.Store.GetTripByID(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if trip == nil {
		return fiber.NewError(fiber.StatusNotFound, "trip not found")
	}
	if trip.OwnerID != ownerID {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	items, err := h.Store.ListTripInvites(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load invites")
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *InvitesHandler) RespondInvite(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	inviteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid invite id")
	}

	var req respondInviteRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid action")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	invite, err := h.Store.GetInviteByID(ctx, inviteID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load invite")
	}
	if invite == nil {
		return fiber.NewError(fiber.StatusNotFound, "invite not found")
	}
	if invite.InvitedUserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	status := map[string]string{"accept": "accepted", "decline": "declined"}[req.Action]
	if err := h.Store.UpdateInviteStatus(ctx, inviteID, status, req.Comment); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update invite")
	}

	if status == "accepted" {
		if err := h.Store.AddTripMember(ctx, invite.TripID, userID, "member"); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to add member")
		}
	}

	return c.JSON(fiber.Map{"status": status})
}
