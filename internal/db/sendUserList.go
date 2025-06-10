package db

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SendUserList(bot *tgbotapi.BotAPI, db *pgxpool.Pool, chatID int64) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.Query(ctx, "SELECT place FROM wishlist WHERE chat_id = $1", chatID)

	if err != nil {
		log.Printf("DB err: %v", err)
		bot.Send(tgbotapi.NewMessage(chatID, "Error loading your list"))
		return
	}
	defer rows.Close()

	var places []string
	for rows.Next() {
		var place string
		if err := rows.Scan(&place); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		places = append(places, place)
	}

	if len(places) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "Your list is empty. Use /add to add place to list"))
		return
	}

	msg := "Your wishlist:\n"
	for i, place := range places {
		msg += fmt.Sprintf("%d. %s\n", i+1, place)
	}

	bot.Send(tgbotapi.NewMessage(chatID, msg))
}
