package buttons

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetAddButton() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Add", "added place to list"),
			tgbotapi.NewInlineKeyboardButtonData("Remove", "removed places from list"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("List", "returned list places"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
