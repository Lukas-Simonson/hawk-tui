package ui

import (
	"strings"

	"hawk-tui/internal/styles"
	"hawk-tui/internal/types"
)

// RenderCommandPane renders the command section
func RenderCommandPane(focusSection types.FocusSection, commandInput string, width, height int) string {
	borderStyle := GetBorderStyle(IsFocused(focusSection, types.FocusCommand))

	var content strings.Builder
	content.WriteString(styles.Title.Render("Command") + "\n\n")

	// Command input
	cursor := ""
	if focusSection == types.FocusCommand {
		cursor = "▌"
	}
	content.WriteString(styles.Normal.Render("git "+commandInput+cursor) + "\n")

	pane := borderStyle.
		Width(width - 4).
		Height(height - 2).
		Render(content.String())

	// Add help text at the bottom
	helpText := styles.Help.Render("tab/shift+tab cycle • ↑/↓ navigate • pgup/pgdn/g/G scroll • enter add to cmd • d view diff • r refresh • ? help • q quit")
	return pane + "\n" + helpText
}
