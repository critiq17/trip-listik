package state

type UserState int

const (
	StateNormal UserState = iota
	StateWaitingAddPlace
	StateWaitingDeletePlace
	StateWaitingForAI
)

type StateManager interface {
	GetUserState(chatID int64) UserState
	SetUserState(chatID int64, state UserState)
}
