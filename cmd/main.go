package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/critiq17/tripListik/db"
	"github.com/critiq17/tripListik/internal/buttons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userWaitingForPlace = make(map[int64]bool)

func main() {

	db.Init()

	botToken := os.Getenv("BOT_TOKEN_TRIP_LISTIK")
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
			chatID := update.CallbackQuery.Message.Chat.ID

			userWaitingForPlace[chatID] = true

			msg := tgbotapi.NewMessage(chatID, "Write the place you want to visit")
			bot.Send(msg)

			bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Okay"))
		}

		if update.Message != nil && update.Message.Chat != nil {
			chatID := update.Message.Chat.ID

			if userWaitingForPlace[chatID] {
				place := update.Message.Text

				err := savePlaceToDB(chatID, place)
				if err != nil {
					log.Printf("Error saving to DB: %v", err)
					bot.Send(tgbotapi.NewMessage(chatID, "Error saving place"))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Place \"%s\" saved", place)))
				}

				userWaitingForPlace[chatID] = false
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Click the button to add a place."))
			}
		}

		if update.Message != nil {

			user := update.Message.From
			userName := user.UserName
			firstName := user.FirstName
			lastName := user.LastName
			text := update.Message.Text

			log.Printf("@%s  %s %s : %s", userName, firstName, lastName, text)

			switch update.Message.Text {
			case "/start":
				msgText := "Hello, " + firstName + "\n" +
					"I'll help you keep track of all the places you want to visit"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyMarkup = buttons.GetAddButton()
				bot.Send(msg)
				continue
			}
		}
	}
}

func savePlaceToDB(chatID int64, place string) error {
	_, err := db.DB.Exec(context.Background(),
		"INSERT INTO userlist (chat_id, place) VALUES ($1, $2)", chatID, place)
	return err
}
