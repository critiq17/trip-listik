package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) RunPolling() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Chat == nil {
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}

	return nil
}
