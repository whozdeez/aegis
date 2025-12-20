package tui

import (

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q", "esc":
			m.Action = ActionExit
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}

		case "enter":
			switch m.Cursor {
			case 0:
				m.Action = ActionAdd
			case 1:
				m.Action = ActionGet
			case 2:
				m.Action = ActionEdit
			case 3:
				m.Action = ActionDelete
			case 4:
				m.Action = ActionList
			case 5:
				m.Action = ActionChangeMaster
			case 6:
				m.Action = ActionExit
			}
			return m, tea.Quit
		}
	}
	return m, nil
}
