package state

import (
	"github.com/critiq17/tripListik/internal/handlers"
)

type Manager struct {
	userStates map[int64]handlers.UserState
}

func NewManager() *Manager {
	return &Manager{
		userStates: make(map[int64]handlers.UserState),
	}
}

func (m *Manager) GetUserState(chatID int64) handlers.UserState {
	if state, exists := m.userStates[chatID]; exists {
		return state
	}
	return handlers.StateNormal
}

func (m *Manager) SetUserState(chatID int64, state handlers.UserState) {
	if state == handlers.StateNormal {
		delete(m.userStates, chatID)
	} else {
		m.userStates[chatID] = state
	}
}
