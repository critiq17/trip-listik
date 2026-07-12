package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"strings"
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

// apiResponse covers both success and error payloads from the Telegram API.
type apiResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  struct {
		RetryAfter int `json:"retry_after"`
	} `json:"parameters"`
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

// esc escapes user-provided text for Telegram HTML parse mode.
// Without this, a trip title like "Sea & Sun <3" makes the API reject the
// whole message with 400 and the notification never reaches the chat.
func esc(s string) string {
	return html.EscapeString(s)
}

// ── SendMessage: plain text ────────────────────────────────────────────────────

// SendMessage sends a Telegram message to a chat (user's Telegram ID).
// The text must already be valid Telegram HTML — escape user input with esc().
// Errors are logged but not propagated to avoid breaking the main flow.
func (b *Bot) SendMessage(ctx context.Context, telegramID int64, text string) {
	if !b.Enabled() {
		return
	}
	b.sendWithRetry(ctx, &sendMessageRequest{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "HTML",
	})
}

// SendInviteAccepted notifies the trip owner that an invite was accepted.
func (b *Bot) SendInviteAccepted(ctx context.Context, ownerTelegramID int64, responderName, tripTitle string) {
	text := fmt.Sprintf("✅ <b>%s</b> accepted your invite to <b>%s</b>!", esc(responderName), esc(tripTitle))
	b.SendMessage(ctx, ownerTelegramID, text)
}

// SendInviteDeclined notifies the trip owner that an invite was declined,
// optionally with a reason and a suggested alternative date.
func (b *Bot) SendInviteDeclined(ctx context.Context, ownerTelegramID int64, responderName, tripTitle, comment, alternativeDate string) {
	text := fmt.Sprintf("❌ <b>%s</b> declined your invite to <b>%s</b>.", esc(responderName), esc(tripTitle))
	if comment != "" {
		text += "\nReason: " + esc(comment)
	}
	if alternativeDate != "" {
		text += "\nSuggested date: " + esc(alternativeDate)
	}
	b.SendMessage(ctx, ownerTelegramID, text)
}

// SendPendingInviteReminder tells a user (on /start) about an invite that was
// created before they opened the bot.
func (b *Bot) SendPendingInviteReminder(ctx context.Context, telegramID int64, tripTitle, inviterName string) {
	text := fmt.Sprintf("You have a pending trip invite for <b>%s</b> from <b>%s</b>.", esc(tripTitle), esc(inviterName))
	b.SendMessage(ctx, telegramID, text)
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
	}
	if miniAppURL != "" {
		msg.ReplyMarkup = &replyMarkup{
			InlineKeyboard: [][]inlineButton{
				{
					{
						Text:   "Open TripListik",
						WebApp: &webAppInfo{URL: miniAppURL},
					},
				},
			},
		}
	}
	b.sendWithRetry(ctx, msg)
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
		locationLine = "\n" + esc(city)
		if startDate != "" && endDate != "" {
			locationLine += " · " + esc(startDate) + " – " + esc(endDate)
		}
	}

	text := fmt.Sprintf(
		"<b>%s</b> invited you to join <b>%s</b>%s",
		esc(inviterName), esc(tripName), locationLine,
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

	// Deep link: opens the Mini App directly on the Inbox page.
	viewURL := strings.TrimRight(miniAppBaseURL, "/") + "/inbox"

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
		esc(newUserName), referralCount,
	)
	b.sendWithRetry(ctx, &sendMessageRequest{
		ChatID:    telegramID,
		Text:      text,
		ParseMode: "HTML",
	})
}

// ── internal send ─────────────────────────────────────────────────────────────

// sendWithRetry wraps send with up to maxRetries attempts and backoff.
// Permanent API errors (4xx other than 429) are not retried — resending the
// same broken payload can never succeed.
func (b *Bot) sendWithRetry(ctx context.Context, payload *sendMessageRequest) int64 {
	delays := []time.Duration{1 * time.Second, 2 * time.Second}
	for attempt := range maxRetries {
		msgID, retryable, wait := b.send(ctx, payload)
		if msgID != 0 {
			return msgID
		}
		if !retryable {
			return 0
		}
		slog.Error("telegram: send attempt failed",
			"attempt", attempt+1,
			"max_attempts", maxRetries,
			"chat_id", payload.ChatID,
		)
		if attempt < len(delays) {
			delay := delays[attempt]
			if wait > delay {
				delay = wait
			}
			select {
			case <-ctx.Done():
				return 0
			case <-time.After(delay):
			}
		}
	}
	return 0
}

// send dispatches a message to the Telegram API.
// Returns the message_id (0 on error), whether the failure is worth retrying,
// and a minimum wait imposed by rate limiting (429 retry_after).
func (b *Bot) send(ctx context.Context, payload *sendMessageRequest) (msgID int64, retryable bool, wait time.Duration) {
	body, err := json.Marshal(payload)
	if err != nil {
		slog.Warn("telegram: failed to marshal message", "err", err)
		return 0, false, 0
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramAPIBase, b.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		slog.Warn("telegram: failed to create request", "err", err)
		return 0, false, 0
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.httpClient.Do(req)
	if err != nil {
		slog.Warn("telegram: sendMessage failed", "telegram_id", payload.ChatID, "err", err)
		return 0, true, 0
	}
	defer resp.Body.Close()

	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Warn("telegram: failed to decode response", "telegram_id", payload.ChatID, "status", resp.StatusCode, "err", err)
		return 0, resp.StatusCode >= 500, 0
	}

	if !result.OK {
		slog.Warn("telegram: sendMessage rejected",
			"telegram_id", payload.ChatID,
			"error_code", result.ErrorCode,
			"description", result.Description,
		)
		switch {
		case result.ErrorCode == http.StatusTooManyRequests:
			return 0, true, time.Duration(result.Parameters.RetryAfter) * time.Second
		case result.ErrorCode >= 500:
			return 0, true, 0
		default:
			// 400 bad request, 403 bot blocked by user, etc. — permanent.
			return 0, false, 0
		}
	}

	return result.Result.MessageID, false, 0
}
