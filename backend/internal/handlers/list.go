package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/critiq17/tripListik/internal/services"
	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ListHandler struct {
	stateManager state.StateManager
	Service      *services.Service
}

func NewListHandler(stateManager state.StateManager, service *services.Service) *ListHandler {
	return &ListHandler{stateManager: stateManager, Service: service}
}

func (h *ListHandler) CanHandle(msg *tgbotapi.Message, userState state.UserState) bool {
	return msg.Text == "/list" && userState == state.StateNormal
}

func (h *ListHandler) Handle(bot BotSender, msg *tgbotapi.Message) error {

	var messageText string

	places, err := h.Service.GetUserPlaces(msg.From.ID)

	if err != nil {
		messageText = "Error"
	} else if len(places) == 0 {
		messageText = "You dont have wishlist"
	} else {
		var placeNames []string
		for _, place := range places {
			placeNames = append(placeNames, place.Name)
		}
		messageText = fmt.Sprintf("Your wishlist:\n- %s", strings.Join(placeNames, "\n- "))
	}
	message := tgbotapi.NewMessage(msg.Chat.ID, messageText)
	if _, err := bot.Send(message); err != nil {
		return err
	}

	log.Printf("user: %s wishlist requested", msg.From.UserName)

	return nil
}
