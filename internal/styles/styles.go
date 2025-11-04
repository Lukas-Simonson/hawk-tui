package styles

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	Selected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	Normal = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	StatusModified = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Bold(true)

	StatusAdded = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Bold(true)

	StatusDeleted = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	StatusUntracked = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)

	Help = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("170")).
			Padding(0, 1).
			Align(lipgloss.Left)

	UnfocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Align(lipgloss.Left)

	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Align(lipgloss.Left)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)
)
