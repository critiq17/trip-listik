package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembersHandler struct {
	Store *store.Store
}

func (h *MembersHandler) ListJoinRequests(c *fiber.Ctx) error {
	ownerID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	status := c.Query("status", "pending")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

	requests, err := h.Store.ListJoinRequestsWithUsers(ctx, tripID, status)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load join requests")
	}

	return c.JSON(fiber.Map{"items": requests})
}

func (h *MembersHandler) ListMembers(c *fiber.Ctx) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	members, err := h.Store.GetTripMembersWithUsers(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load members")
	}

	return c.JSON(fiber.Map{"items": members})
}

func (h *MembersHandler) JoinTrip(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	trip, err := h.Store.GetTripByID(ctx, tripID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if trip == nil {
		return fiber.NewError(fiber.StatusNotFound, "trip not found")
	}

	isMember, err := h.Store.IsTripMember(ctx, tripID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to check membership")
	}
	if isMember {
		return c.JSON(fiber.Map{"status": "already_member"})
	}

	if trip.Visibility == "public" {
		if err := h.Store.AddTripMember(ctx, tripID, userID, "member"); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to join trip")
		}
		return c.JSON(fiber.Map{"status": "joined"})
	}

	jr, err := h.Store.GetJoinRequest(ctx, tripID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load join request")
	}
	if jr != nil {
		return c.JSON(fiber.Map{"status": jr.Status})
	}

	if err := h.Store.CreateJoinRequest(ctx, tripID, userID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create join request")
	}

	return c.JSON(fiber.Map{"status": "pending"})
}

func (h *MembersHandler) ApproveJoin(c *fiber.Ctx) error {
	ownerID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user_id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

	err = h.Store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := h.Store.UpdateJoinRequestStatus(ctx, tripID, userID, "approved"); err != nil {
			return err
		}
		return h.Store.AddTripMember(ctx, tripID, userID, "member")
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to approve join request")
	}

	return c.JSON(fiber.Map{"status": "approved"})
}

func (h *MembersHandler) RejectJoin(c *fiber.Ctx) error {
	ownerID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user_id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

	if err := h.Store.UpdateJoinRequestStatus(ctx, tripID, userID, "rejected"); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to reject join request")
	}

	return c.JSON(fiber.Map{"status": "rejected"})
}

func (h *MembersHandler) RemoveMember(c *fiber.Ctx) error {
	ownerID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}
	memberID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

	if err := h.Store.RemoveTripMember(ctx, tripID, memberID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to remove member")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
