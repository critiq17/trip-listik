package handlers

import "github.com/critiq17/tripListik/internal/state"

func RegisterHandlers(stateManager state.StateManager) []Handler {
	return []Handler{
		&StartHandler{},
	}
}
