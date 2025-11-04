package ui

import (
	"hawk-tui/internal/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderCommitPopup renders a popup for composing commit messages
func RenderCommitPopup(message string, width, height int) string {
	// Calculate popup dimensions (80% of screen)
	popupWidth := int(float64(width) * 0.8)
	popupHeight := int(float64(height) * 0.6)

	if popupWidth < 50 {
		popupWidth = 50
	}
	if popupHeight < 15 {
		popupHeight = 15
	}

	// Create the popup content
	var content strings.Builder
	content.WriteString(styles.Title.Render("Commit Message") + "\n\n")
	content.WriteString(styles.Help.Render("Enter your commit message below (multi-line supported):") + "\n\n")

	// Show the message with cursor
	content.WriteString(message + "▌\n\n")

	content.WriteString(styles.Help.Render("Ctrl+Enter: Commit  •  Esc: Cancel"))

	// Render the popup with a box
	popup := styles.Box.
		Width(popupWidth - 4).
		Height(popupHeight - 4).
		Render(content.String())

	// Center the popup on screen
	centered := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)

	return centered
}
