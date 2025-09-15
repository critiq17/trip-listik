package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	DSN        string
	API_KEY_AI string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env not found")
	}

	return &Config{
		DSN:        os.Getenv("DSN"),
		BotToken:   os.Getenv("BOT_TOKEN"),
		API_KEY_AI: os.Getenv("API_KEY_AI"),
	}
}
