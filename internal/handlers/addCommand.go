package handlers

import (
	"fmt"

	"github.com/critiq17/tripListik/internal/db"
	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddHandler struct {
	stateManager state.StateManager
}

func NewAddHandler(stateManager state.StateManager) *AddHandler {
	return &AddHandler{stateManager: stateManager}
}

func (h *AddHandler) CanHandle(msg *tgbotapi.Message, userState state.UserState) bool {
	return msg.Text == "/add" || userState == state.StateWaitingAddPlace
}

func (h *AddHandler) Handle(bot BotSender, msg *tgbotapi.Message) error {
	userState := h.stateManager.GetUserState(msg.Chat.ID)

	if msg.Text == "/add" {
		h.stateManager.SetUserState(msg.Chat.ID, state.StateWaitingAddPlace)
		_, err := bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Enter place what you want add in list"))
		return err
	}

	if userState == state.StateWaitingAddPlace {
		err := db.SavePlaceToDB(msg.Chat.ID, msg.Text)
		if err != nil {
			return err
		}
		h.stateManager.SetUserState(msg.Chat.ID, state.StateNormal)

		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Place %s added", msg.Text)))

	}
	return nil

}
