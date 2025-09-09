package handlers

import (
	"github.com/critiq17/tripListik/internal/services"
	"github.com/critiq17/tripListik/internal/state"
)

func RegisterHandlers(stateManager state.StateManager, service *services.Service) []Handler {

	return []Handler{
		&StartHandler{},
		NewAddHandler(stateManager, service),
		NewDeleteHandler(stateManager),
	}
}
