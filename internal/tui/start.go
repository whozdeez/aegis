package tui

import (

	tea "github.com/charmbracelet/bubbletea"
)

func RunStart() (Action, error) {
	p := tea.NewProgram(initialModel())
	model, err := p.Run()
	if err != nil {
		return ActionNone, err
	}

	finalModel := model.(Model)
	return finalModel.Action, nil
}

func initialModel() Model {
	return Model{
		Cursor: 0,
		Action: ActionNone,
		Items: []string{
			"Add new password",
			"Get password",
			"Edit password",
			"Delete password",
			"List services",
			"Change master password",
			"Exit",
		},
	}
}
