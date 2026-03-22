package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/telegram"
	"github.com/gofiber/fiber/v2"
)

type BotHandler struct {
	Store *store.Store
	Bot   *telegram.Bot
	Cfg   *config.Config
}

// TelegramUpdate is a minimal struct for deserializing Telegram webhook payloads.
type TelegramUpdate struct {
	Message *struct {
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		From *struct {
			ID int64 `json:"id"`
		} `json:"from"`
		Text string `json:"text"`
	} `json:"message"`
}

// HandleUpdate processes incoming Telegram bot webhook updates.
// Returns HTTP 200 always to prevent Telegram from retrying.
func (h *BotHandler) HandleUpdate(c *fiber.Ctx) error {
	var update TelegramUpdate
	if err := c.BodyParser(&update); err != nil {
		return c.SendStatus(200)
	}
	if update.Message == nil {
		return c.SendStatus(200)
	}

	text := strings.TrimSpace(update.Message.Text)
	chatID := update.Message.Chat.ID

	if strings.HasPrefix(text, "/start") {
		parts := strings.Fields(text)
		payload := ""
		if len(parts) > 1 {
			payload = parts[1]
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		// Check for pending invites and notify the user
		pendingInvites, err := h.Store.GetAndClearPendingInvites(ctx, chatID)
		if err == nil && len(pendingInvites) > 0 {
			for _, pi := range pendingInvites {
				msg := "You have a pending trip invite for: " + pi.TripTitle + " from " + pi.InviterName
				h.Bot.SendMessage(ctx, chatID, msg)
			}
		}

		// Note: referral tracking from start_param happens in the auth handler
		// when the user authenticates via the Mini App — not duplicated here.
		_ = payload

		miniAppURL := ""
		if h.Cfg != nil {
			miniAppURL = h.Cfg.MiniAppURL
		}
		h.Bot.SendWelcome(ctx, chatID, miniAppURL)
	}

	return c.SendStatus(200)
}
