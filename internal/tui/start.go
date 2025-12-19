package tui

import (

	tea "github.com/charmbracelet/bubbletea"
)

func RunStart() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}

func initialModel() Model {
	return Model{
		Cursor: 0,
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
