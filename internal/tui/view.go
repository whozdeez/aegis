package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

func (m Model) View() string {
	s := "\n"
	s += titleStyle.Render("🛡️  Aegis Password Manager") + "\n\n"

	for i, item := range m.Items {
		if m.Cursor == i {
			s += cursorStyle.Render("❯ ") + item + "\n"
		} else {
			s += "  " + item + "\n"
		}
	}

	s += "\n↑/↓ navigate • Enter select • q quit\n"
	return s
}
