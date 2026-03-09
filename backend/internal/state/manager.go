package state

import "sync"

type Manager struct {
	mu         sync.RWMutex
	userStates map[int64]UserState
}

func NewManager() *Manager {
	return &Manager{
		userStates: make(map[int64]UserState),
	}
}

func (m *Manager) GetUserState(chatID int64) UserState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if state, exists := m.userStates[chatID]; exists {
		return state
	}
	return StateNormal
}

func (m *Manager) SetUserState(chatID int64, state UserState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if state == StateNormal {
		delete(m.userStates, chatID)
	} else {
		m.userStates[chatID] = state
	}
}
