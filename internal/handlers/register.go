package handlers

func RegisterHandlers(stateManager StateManager) []Handler {
	return []Handler{
		&StartHandler{},
	}
}
