package handlers

import (
	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	CanHandle(message *tgbotapi.Message, userState state.UserState) bool
	Handle(bot BotSender, message *tgbotapi.Message) error
}

type BotSender interface {
	Send(msg tgbotapi.Chattable) (tgbotapi.Message, error)
}
