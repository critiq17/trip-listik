package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	DB_URL   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env not found")
	}

	return &Config{
		DB_URL:   os.Getenv("DATABASE_URL"),
		BotToken: os.Getenv("BOT_TOKEN_TRIP_LISTIK"),
	}
}
