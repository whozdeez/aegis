package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

const Banner = `
    █████╗ ███████╗ ██████╗ ██╗███████╗
   ██╔══██╗██╔════╝██╔════╝ ██║██╔════╝
   ███████║█████╗  ██║  ███╗██║███████╗
   ██╔══██║██╔══╝  ██║   ██║██║╚════██║
   ██║  ██║███████╗╚██████╔╝██║███████║
   ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚══════╝

        Secure • Offline • Zero-Knowledge
`
var bannerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("5")).
	Bold(true)


func (m Model) View() string {
	s := "\n"
	s += bannerStyle.Render(Banner) + "\n"

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
