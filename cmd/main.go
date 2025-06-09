package main

import (
	"fmt"
	"log"
	"os"

	"github.com/critiq17/tripListik/internal/api"
	"github.com/critiq17/tripListik/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userWaitingAddingPlace = make(map[int64]bool)
var userWaitingDeletingPlace = make(map[int64]bool)
var userWaitingForAI = make(map[int64]bool)

func main() {

	db.Init()

	botToken := os.Getenv("BOT_TOKEN_TRIP_LISTIK")
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Println(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		user := update.Message.From
		userName := user.UserName
		firstName := user.FirstName
		lastName := user.LastName
		text := update.Message.Text
		chatID := update.Message.Chat.ID

		log.Printf("@%s  %s %s : %s", userName, firstName, lastName, text)

		/*

				if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
				chatID := update.CallbackQuery.Message.Chat.ID

				userWaitingForPlace[chatID] = true

				msg := tgbotapi.NewMessage(chatID, "Write the place you want to visit")
				bot.Send(msg)

				bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Okay"))
			}

		*/

		if update.Message != nil && update.Message.Chat != nil {

			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(chatID, "👋 Hello, "+firstName+"\n"+
					"I'm your travel wishlist bot\n"+
					"You can save and manage places you want to visit  Just use:\n"+
					"/add - to add a new place\n"+
					"/delete - to delete a place\n"+
					"/list - return your wishlist\n"+
					"/tripcost")

				bot.Send(msg)
			}

			if update.Message.Text == "/add" {
				userWaitingAddingPlace[chatID] = true

				msg := tgbotapi.NewMessage(chatID, "Write place you want to visit")
				bot.Send(msg)
				continue
			}

			if update.Message.Text == "/delete" {
				userWaitingDeletingPlace[chatID] = true
				msg := tgbotapi.NewMessage(chatID, "Write the place you want tot delete")
				bot.Send(msg)
				continue
			}

			if update.Message.Text == "/list" {
				db.SendUserList(bot, db.DB, chatID)
			}

			if update.Message.Text == "/tripcost" {
				userWaitingForAI[chatID] = true

				bot.Send(tgbotapi.NewMessage(chatID, "Send me your trip details (e.g, city, days, preferences), and i calculate an appriximate cost"))
				continue
			}

			if userWaitingForAI[chatID] {
				prompt := update.Message.Text

				response, err := api.SendToAi(prompt)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Error getting answer from AI"))
				} else {
					msg := tgbotapi.NewMessage(chatID, response)
					msg.ParseMode = "Markdown"
					bot.Send(msg)
				}

				userWaitingForAI[chatID] = false
				continue
			}

			if userWaitingAddingPlace[chatID] {
				place := update.Message.Text

				err := db.SavePlaceToDB(chatID, place)

				if err != nil {
					log.Printf("Error saving to DB: %v", err)
					bot.Send(tgbotapi.NewMessage(chatID, "Error saving place"))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Place \"%s\" saved", place)))
				}

				userWaitingAddingPlace[chatID] = false
				continue
			} else {
				// bot.Send(tgbotapi.NewMessage(chatID, "Click the button to add a place."))
			}

			if userWaitingDeletingPlace[chatID] {
				place := update.Message.Text

				err := db.DeletePlaceFromDB(chatID, place)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Error deleting place"))
				} else {
					bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Place \"%s\" deleted", place)))
				}

				userWaitingDeletingPlace[chatID] = false
				continue
			}
		}
	}
}
