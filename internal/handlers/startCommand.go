package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct{}

func (h *StartHandler) CanHandle(msg *tgbotapi.Message, state UserState) bool {
	return msg.Text == "/start" && state == StateNormal
}

func (h *StartHandler) Handle(bot BotSender, msg *tgbotapi.Message) error {
	welcome := fmt.Sprintf("👋 Hello, %s\n"+
		"I'm your travel wishlist bot\n"+
		"You can save and manage places you want to visit. Just use:\n"+
		"/add - to add a new place\n"+
		"/delete - to delete a place\n"+
		"/list - return your wishlist\n"+
		"/tripcost - calculate trip costs with AI", msg.From.FirstName)

	message := tgbotapi.NewMessage(msg.Chat.ID, welcome)
	_, err := bot.Send(message)
	return err
}
