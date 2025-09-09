package bot

import (
	"github.com/critiq17/tripListik/internal/handlers"
	"github.com/critiq17/tripListik/internal/services"
	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewBotWithHandlers(token string, stateManager state.StateManager, service *services.Service) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		api:          api,
		stateManager: stateManager,
	}

	bot.handlers = handlers.RegisterHandlers(stateManager, service)

	return bot, nil
}
