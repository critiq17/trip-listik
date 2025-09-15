package handlers

import (
	"log"

	"github.com/critiq17/tripListik/internal/api"
	"github.com/critiq17/tripListik/internal/services"
	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TripCostHandler struct {
	stateManager state.StateManager
	Service      *services.Service
}

func NewTripCostHandler(stateManager state.StateManager, service *services.Service) *TripCostHandler {
	return &TripCostHandler{stateManager: stateManager, Service: service}
}

func (h *TripCostHandler) CanHandle(msg *tgbotapi.Message, userState state.UserState) bool {
	return msg.Text == "/tripcost" || userState == state.StateWaitingForAI
}

func (h *TripCostHandler) Handle(bot BotSender, msg *tgbotapi.Message) error {

	userState := h.stateManager.GetUserState(msg.Chat.ID)

	if msg.Text == "/tripcost" {
		h.stateManager.SetUserState(msg.Chat.ID, state.StateWaitingForAI)
		_, err := bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Send your trip plan"))
		return err
	}

	if userState == state.StateWaitingForAI {
		string, err := api.SendToAi(msg.Text)
		if err != nil {
			return err
		}
		log.Printf("user: %s send to AI %s", msg.From.UserName, msg.Text)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, string))
	}

	return nil
}
