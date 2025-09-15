package main

import (
	"fmt"
	"log"

	"github.com/critiq17/tripListik/internal/bot"
	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/repository"
	"github.com/critiq17/tripListik/internal/services"
	"github.com/critiq17/tripListik/internal/state"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	cfg := config.LoadConfig()
	dsn := cfg.DSN
	bot_token := cfg.BotToken

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("error init DB: %v", err)
	}

	fmt.Println("DB init successfully")

	repo := repository.NewRepository(db)
	service := services.NewService(repo)
	stateManager := state.NewManager()

	if bot_token == "" {
		log.Fatal("Bot_token not found")
	}

	bot, err := bot.NewBotWithHandlers(cfg.BotToken, stateManager, service)
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	}

	log.Println("Bot started successfully")

	if err := bot.Start(); err != nil {
		log.Fatal("Bot started failed: ", err)
	}
}
