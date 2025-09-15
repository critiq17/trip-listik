package handlers

import (
	"fmt"
	"log"

	"github.com/critiq17/tripListik/internal/services"
	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddHandler struct {
	stateManager state.StateManager
	Service      *services.Service
}

func NewAddHandler(stateManager state.StateManager, service *services.Service) *AddHandler {
	return &AddHandler{stateManager: stateManager, Service: service}
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
		err := h.Service.AddPlace(msg.From.ID, msg.From.FirstName, msg.From.UserName, msg.Text)
		if err != nil {
			return err
		}
		h.stateManager.SetUserState(msg.Chat.ID, state.StateNormal)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Place %s added", msg.Text)))
		log.Printf("user: %s add a place: %s", msg.From.UserName, msg.Text)
	}
	return nil

}
