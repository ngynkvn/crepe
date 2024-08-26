package app

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	FocusedBorder lipgloss.Style
	BlurredBorder lipgloss.Style
	InputStyle    lipgloss.Style
}

var DefaultStyles = Styles{
	FocusedBorder: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")),
	BlurredBorder: lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder()),
	InputStyle: lipgloss.NewStyle().Padding(2),
}

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))
)
