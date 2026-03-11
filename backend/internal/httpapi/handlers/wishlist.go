package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/httpapi/validate"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WishlistHandler struct {
	Store *store.Store
}

type wishlistRequest struct {
	CountryCode string `json:"country_code" validate:"required,len=2"`
	City        string `json:"city"`
	Note        string `json:"note"`
}

func (h *WishlistHandler) List(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	items, err := h.Store.ListWishlist(ctx, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load wishlist")
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *WishlistHandler) Create(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var req wishlistRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid country_code")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	item := &models.UserWishlistItem{
		UserID:      userID,
		CountryCode: req.CountryCode,
		City:        req.City,
		Note:        req.Note,
	}
	if err := h.Store.AddWishlistItem(ctx, item); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to add wishlist item")
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *WishlistHandler) Delete(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	itemID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid item id")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := h.Store.DeleteWishlistItem(ctx, userID, itemID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete wishlist item")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
