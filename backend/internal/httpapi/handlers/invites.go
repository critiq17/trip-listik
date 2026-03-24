package handlers

import (
	"errors"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/httpapi/validate"
	"github.com/critiq17/tripListik/internal/invites"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InvitesHandler struct {
	Service *invites.Service
	Store   *store.Store // kept for read-only queries (GetInvite, ListTripInvites)
}

// ── request bodies ─────────────────────────────────────────────────────────────

type inviteRequest struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

type respondInviteRequest struct {
	Action          string  `json:"action" validate:"required,oneof=accept decline"`
	Comment         *string `json:"comment"`
	AlternativeDate *string `json:"alternative_date"`
}

// ── handlers ───────────────────────────────────────────────────────────────────

func (h *InvitesHandler) InviteUser(c *fiber.Ctx) error {
	requesterID, ok := middleware.GetUserID(c)
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

	invite, err := h.Service.InviteUser(c.Context(), requesterID, tripID, req.Username, req.UserID)
	if err != nil {
		if errors.Is(err, invites.ErrDuplicate) && invite != nil {
			return c.JSON(fiber.Map{"invite_id": invite.ID.String(), "status": invite.Status})
		}
		return mapInviteError(err)
	}

	return c.JSON(fiber.Map{"invite_id": invite.ID.String(), "status": invite.Status})
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

	if err := h.Service.RespondInvite(c.Context(), userID, inviteID, req.Action, req.Comment, req.AlternativeDate); err != nil {
		return mapInviteError(err)
	}

	status := map[string]string{"accept": "accepted", "decline": "declined"}[req.Action]
	return c.JSON(fiber.Map{"status": status})
}

func (h *InvitesHandler) CancelInvite(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	inviteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid invite id")
	}

	if err := h.Service.CancelInvite(c.Context(), userID, inviteID); err != nil {
		return mapInviteError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *InvitesHandler) ListTripInvites(c *fiber.Ctx) error {
	requesterID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	trip, err := h.Store.GetTripByID(c.Context(), tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if trip == nil {
		return fiber.NewError(fiber.StatusNotFound, "trip not found")
	}

	isOwner := trip.OwnerID == requesterID
	if !isOwner {
		if trip.Visibility != "group" {
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		}
		isMember, err := h.Store.IsTripMember(c.Context(), tripID, requesterID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to check membership")
		}
		if !isMember {
			return fiber.NewError(fiber.StatusForbidden, "forbidden")
		}
	}

	items, err := h.Store.ListTripInvites(c.Context(), tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load invites")
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *InvitesHandler) GetInvite(c *fiber.Ctx) error {
	inviteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid invite id")
	}

	invite, err := h.Store.GetInviteByID(c.Context(), inviteID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load invite")
	}
	if invite == nil {
		return fiber.NewError(fiber.StatusNotFound, "invite not found")
	}

	// Return enriched view if available.
	views, err := h.Store.ListUserInvites(c.Context(), invite.InvitedUserID, 20)
	if err == nil {
		for _, v := range views {
			if v.ID == inviteID {
				return c.JSON(fiber.Map{"invite": v})
			}
		}
	}
	return c.JSON(fiber.Map{"invite": invite})
}

// ── error mapping ──────────────────────────────────────────────────────────────

func mapInviteError(err error) error {
	switch {
	case errors.Is(err, invites.ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, "not found")
	case errors.Is(err, invites.ErrForbidden):
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	default:
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
}
