package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const telegramAPIBase = "https://api.telegram.org"

type Bot struct {
	token      string
	httpClient *http.Client
}

func NewBot(token string) *Bot {
	return &Bot{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Enabled returns false if bot token is not configured — allows graceful degradation
func (b *Bot) Enabled() bool {
	return b.token != ""
}

type sendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// SendMessage sends a Telegram message to a chat (user's Telegram ID).
// Errors are logged but not propagated to avoid breaking the main flow.
func (b *Bot) SendMessage(ctx context.Context, telegramID int64, text string) {
	if !b.Enabled() {
		return
	}

	payload := sendMessageRequest{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		slog.Warn("telegram: failed to marshal message", "err", err)
		return
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramAPIBase, b.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		slog.Warn("telegram: failed to create request", "err", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		slog.Warn("telegram: sendMessage failed", "telegram_id", telegramID, "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("telegram: sendMessage non-200", "telegram_id", telegramID, "status", resp.StatusCode)
	}
}
