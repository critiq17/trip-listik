package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/auth"
	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Store *store.Store
	Cfg   *config.Config
}

type telegramAuthRequest struct {
	InitData string `json:"initData"`
}

type telegramAuthResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

type userResponse struct {
	ID         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	PhotoURL   string `json:"photo_url"`
}

func (h *AuthHandler) TelegramAuth(c *fiber.Ctx) error {
	var req telegramAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if req.InitData == "" {
		return fiber.NewError(fiber.StatusBadRequest, "initData is required")
	}

	tgUser, err := auth.ValidateTelegramInitData(req.InitData, h.Cfg.TelegramWebAppSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid telegram init data")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	user := &models.User{
		TelegramID: tgUser.ID,
		Username:   tgUser.Username,
		FirstName:  tgUser.FirstName,
		LastName:   tgUser.LastName,
		PhotoURL:   tgUser.PhotoURL,
	}

	stored, err := h.Store.UpsertTelegramUser(ctx, user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to upsert user")
	}

	token, err := auth.NewToken(stored.ID, stored.TelegramID, h.Cfg.JWTSecret, 30*24*time.Hour)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	resp := telegramAuthResponse{
		Token: token,
		User: userResponse{
			ID:         stored.ID.String(),
			TelegramID: stored.TelegramID,
			Username:   stored.Username,
			FirstName:  stored.FirstName,
			LastName:   stored.LastName,
			PhotoURL:   stored.PhotoURL,
		},
	}

	return c.JSON(resp)
}
