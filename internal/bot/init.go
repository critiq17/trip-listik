package bot

import (
	"github.com/critiq17/tripListik/internal/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BewBotWithHandlers(token string, stateManager handlers.StateManager) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		api:          api,
		stateManager: stateManager,
	}

	bot.handlers = handlers.RegisterHandlers(stateManager)

	return bot, nil
}
