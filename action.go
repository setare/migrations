package migrations

type ActionType string

const (
	ActionTypeDo   ActionType = "do"
	ActionTypeUndo ActionType = "undo"
)

type Action struct {
	Action    ActionType
	Migration Migration
}
