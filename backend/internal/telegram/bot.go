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

const maxRetries = 3

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

// ── Message types ──────────────────────────────────────────────────────────────

type sendMessageRequest struct {
	ChatID      int64        `json:"chat_id"`
	Text        string       `json:"text"`
	ParseMode   string       `json:"parse_mode,omitempty"`
	ReplyMarkup *replyMarkup `json:"reply_markup,omitempty"`
}

type sendMessageResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}

type replyMarkup struct {
	InlineKeyboard [][]inlineButton `json:"inline_keyboard"`
}

type inlineButton struct {
	Text         string      `json:"text"`
	CallbackData string      `json:"callback_data,omitempty"`
	WebApp       *webAppInfo `json:"web_app,omitempty"`
	URL          string      `json:"url,omitempty"`
}

type webAppInfo struct {
	URL string `json:"url"`
}

// ── SendMessage: plain text ────────────────────────────────────────────────────

// SendMessage sends a Telegram message to a chat (user's Telegram ID).
// Errors are logged but not propagated to avoid breaking the main flow.
func (b *Bot) SendMessage(ctx context.Context, telegramID int64, text string) {
	if !b.Enabled() {
		return
	}
	b.send(ctx, &sendMessageRequest{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "HTML",
	})
}

// ── DeleteMessage: remove a previously sent message ───────────────────────────

// DeleteMessage removes a Telegram message from a user's chat.
// Used to clean up invite notifications once they are acted upon.
func (b *Bot) DeleteMessage(ctx context.Context, chatID int64, messageID int64) {
	if !b.Enabled() || messageID == 0 {
		return
	}

	body, err := json.Marshal(map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
	})
	if err != nil {
		return
	}

	url := fmt.Sprintf("%s/bot%s/deleteMessage", telegramAPIBase, b.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		slog.Warn("telegram: deleteMessage failed", "chat_id", chatID, "message_id", messageID, "err", err)
		return
	}
	defer resp.Body.Close()
}

// ── SendWelcome: /start handler message ───────────────────────────────────────

// SendWelcome sends the bilingual welcome message with an Open Mini App button.
func (b *Bot) SendWelcome(ctx context.Context, telegramID int64, miniAppURL string) {
	if !b.Enabled() {
		return
	}
	text := "Вітаємо у TripListik.\nWelcome to TripListik.\n\nПлануй подорожі. Діліться враженнями.\nPlan trips. Share the journey."
	msg := &sendMessageRequest{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "HTML",
		ReplyMarkup: &replyMarkup{
			InlineKeyboard: [][]inlineButton{
				{
					{
						Text:   "Open TripListik",
						WebApp: &webAppInfo{URL: miniAppURL},
					},
				},
			},
		},
	}
	b.send(ctx, msg)
}

// ── SendInviteNotification: invite with deep link button (returns message_id) ──

// SendInviteNotification sends an invite notification and returns the Telegram
// message_id so it can be stored and later deleted when the invite is acted upon.
// Uses retry logic for reliability.
func (b *Bot) SendInviteNotification(
	ctx context.Context,
	telegramID int64,
	inviterName string,
	tripName string,
	city string,
	startDate string,
	endDate string,
	miniAppBaseURL string,
) int64 {
	if !b.Enabled() {
		return 0
	}

	locationLine := ""
	if city != "" {
		locationLine = "\n" + city
		if startDate != "" && endDate != "" {
			locationLine += " · " + startDate + " – " + endDate
		}
	}

	text := fmt.Sprintf(
		"<b>%s</b> invited you to join <b>%s</b>%s",
		inviterName, tripName, locationLine,
	)

	// When miniAppBaseURL is not configured, send plain text so the notification
	// always arrives. The WebApp button requires a valid HTTPS URL.
	if miniAppBaseURL == "" {
		return b.sendWithRetry(ctx, &sendMessageRequest{
			ChatID:    telegramID,
			Text:      text,
			ParseMode: "HTML",
		})
	}

	// Deep link: opens Mini App and navigates to Notifications tab
	viewURL := miniAppBaseURL + "?startapp=notifications"

	msg := &sendMessageRequest{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "HTML",
		ReplyMarkup: &replyMarkup{
			InlineKeyboard: [][]inlineButton{
				{
					{
						Text:   "View Invite",
						WebApp: &webAppInfo{URL: viewURL},
					},
				},
			},
		},
	}
	return b.sendWithRetry(ctx, msg)
}

// SendReferralNotification notifies a referrer that someone joined via their link.
func (b *Bot) SendReferralNotification(ctx context.Context, telegramID int64, newUserName string, referralCount int64) {
	if !b.Enabled() {
		return
	}
	text := fmt.Sprintf(
		"<b>@%s</b> just joined TripListik via your link! You now have <b>%d</b> referral(s).",
		newUserName, referralCount,
	)
	b.sendWithRetry(ctx, &sendMessageRequest{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "HTML",
	})
}

// ── internal send ─────────────────────────────────────────────────────────────

// sendWithRetry wraps send with up to maxRetries attempts and exponential backoff.
// It logs a structured error on each failed attempt.
func (b *Bot) sendWithRetry(ctx context.Context, payload *sendMessageRequest) int64 {
	delays := []time.Duration{1 * time.Second, 2 * time.Second}
	for attempt := range maxRetries {
		msgID := b.send(ctx, payload)
		if msgID != 0 {
			return msgID
		}
		slog.Error("telegram: send attempt failed",
			"attempt", attempt+1,
			"max_attempts", maxRetries,
			"chat_id", payload.ChatID,
		)
		if attempt < len(delays) {
			select {
			case <-ctx.Done():
				return 0
			case <-time.After(delays[attempt]):
			}
		}
	}
	return 0
}

// send dispatches a message to the Telegram API and returns the message_id (0 on error).
func (b *Bot) send(ctx context.Context, payload *sendMessageRequest) int64 {
	body, err := json.Marshal(payload)
	if err != nil {
		slog.Warn("telegram: failed to marshal message", "err", err)
		return 0
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramAPIBase, b.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		slog.Warn("telegram: failed to create request", "err", err)
		return 0
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		slog.Warn("telegram: sendMessage failed", "telegram_id", payload.ChatID, "err", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("telegram: sendMessage non-200", "telegram_id", payload.ChatID, "status", resp.StatusCode)
		return 0
	}

	var result sendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0
	}
	return result.Result.MessageID
}
