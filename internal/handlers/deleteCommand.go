package handlers

import (
	"fmt"

	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeleteHandler struct {
	stateManager state.StateManager
}

func NewDeleteHandler(stateManager state.StateManager) *DeleteHandler {
	return &DeleteHandler{stateManager: stateManager}
}

func (h *DeleteHandler) CanHandle(msg *tgbotapi.Message, userState state.UserState) bool {
	return msg.Text == "/delete" || userState == state.StateWaitingDeletePlace
}

func (h *DeleteHandler) Handle(bot BotSender, msg *tgbotapi.Message) error {

	userState := h.stateManager.GetUserState(msg.Chat.ID)

	if msg.Text == "/delete" {
		h.stateManager.SetUserState(msg.Chat.ID, state.StateWaitingDeletePlace)
		_, err := bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Send place what you want to delete from list"))
		return err
	}

	if userState == state.StateWaitingDeletePlace {

		h.stateManager.SetUserState(msg.Chat.ID, state.StateNormal)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Place %s delete from list", msg.Text)))
	}

	return nil
}
