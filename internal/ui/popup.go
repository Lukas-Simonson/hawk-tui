package ui

import (
	"hawk-tui/internal/styles"

	"github.com/charmbracelet/lipgloss"
)

// RenderPopup renders a popup overlay for command output
func RenderPopup(output string, width, height int) string {
	// Calculate popup dimensions (80% of screen)
	popupWidth := int(float64(width) * 0.8)
	popupHeight := int(float64(height) * 0.8)

	if popupWidth < 40 {
		popupWidth = 40
	}
	if popupHeight < 10 {
		popupHeight = 10
	}

	// Create the popup content
	content := styles.Title.Render("Command Output") + "\n\n"
	content += output

	// Render the popup with a box
	popup := styles.Box.
		Width(popupWidth - 4).
		Height(popupHeight - 4).
		Render(content)

	// Center the popup on screen
	centered := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)

	// Add help text at bottom
	helpText := "\n" + styles.Help.Render("Press any key to close...")

	return centered + helpText
}
