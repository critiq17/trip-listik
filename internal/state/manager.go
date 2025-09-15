package state

type Manager struct {
	userStates map[int64]UserState
}

func NewManager() *Manager {
	return &Manager{
		userStates: make(map[int64]UserState),
	}
}

func (m *Manager) GetUserState(chatID int64) UserState {
	if state, exists := m.userStates[chatID]; exists {
		return state
	}
	return StateNormal
}

func (m *Manager) SetUserState(chatID int64, state UserState) {
	if state == StateNormal {
		delete(m.userStates, chatID)
	} else {
		m.userStates[chatID] = state
	}
}
