package bot

import (
	"github.com/critiq17/tripListik/internal/handlers"
	"github.com/critiq17/tripListik/internal/state"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api          *tgbotapi.BotAPI
	handlers     []handlers.Handler
	stateManager state.StateManager
}

func NewBot(api *tgbotapi.BotAPI, stateManager state.StateManager) *Bot {

	return &Bot{
		api:          api,
		stateManager: stateManager,
		handlers:     handlers.RegisterHandlers(stateManager),
	}
}

func (b *Bot) Start() error {
	return b.RunPolling()
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) error {
	userState := b.stateManager.GetUserState(msg.Chat.ID)

	for _, h := range b.handlers {
		if h.CanHandle(msg, userState) {
			return h.Handle(b.api, msg)
		}
	}
	return nil
}
