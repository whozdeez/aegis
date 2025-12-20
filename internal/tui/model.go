package tui

type Model struct {
	Cursor int
	Items  []string
	Action Action
}

type Action int

const (
	ActionNone Action = iota
	ActionAdd
	ActionGet
	ActionEdit
	ActionDelete
	ActionList
	ActionChangeMaster
	ActionExit
)