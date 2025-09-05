package main

import (
	"log"

	"github.com/critiq17/tripListik/internal/bot"
	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/db"
	"github.com/critiq17/tripListik/internal/state"
)

func main() {

	db.Init()
	cfg := config.LoadConfig()
	stateManager := state.NewManager()

	bot_token := cfg.BotToken

	if bot_token == "" {
		log.Fatal("Bot_token not found")
	}

	bot, err := bot.NewBotWithHandlers(cfg.BotToken, stateManager)

	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	}

	log.Println("Bot started successfully")

	if err := bot.Start(); err != nil {
		log.Fatal("Bot started failed: ", err)
	}
}
