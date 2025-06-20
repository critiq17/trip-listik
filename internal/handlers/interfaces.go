package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler interface {
	CanHandle(message *tgbotapi.Message, userState UserState) bool
	Handle(bot BotSender, message *tgbotapi.Message) error
}

type BotSender interface {
	Send(msg tgbotapi.Chattable) (tgbotapi.Message, error)
}

type UserState int

const (
	StateNormal UserState = iota
	StateWaitingAddPlace
	StateWaitingDeletePlace
	StatreWaitingForAI
)

type StateManager interface {
	GetUserState(chatID int64) UserState
	SetUserState(chatID int64, state UserState)
}
