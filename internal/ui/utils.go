package ui

import (
	"hawk-tui/internal/styles"
	"hawk-tui/internal/types"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// ApplyHorizontalViewport applies horizontal scrolling to a line
func ApplyHorizontalViewport(line string, scrollOffset, viewportWidth int) string {
	// Get the visual width of the line
	lineWidth := ansi.StringWidth(line)

	// If line fits in viewport and no scroll, return as-is
	if lineWidth <= viewportWidth && scrollOffset == 0 {
		return line
	}

	// Use ansi.Cut to handle ANSI codes properly while cutting
	// Cut from scrollOffset to scrollOffset + viewportWidth
	result := ansi.Cut(line, scrollOffset, scrollOffset+viewportWidth)

	return result
}

// FormatStatus returns a styled status indicator
func FormatStatus(status string) string {
	var style lipgloss.Style
	var label string

	switch status {
	case "M", "MM", "AM":
		style = styles.StatusModified
		label = status
	case "A":
		style = styles.StatusAdded
		label = status
	case "D":
		style = styles.StatusDeleted
		label = status
	case "??":
		style = styles.StatusUntracked
		label = status
	default:
		style = styles.Normal
		label = status
	}

	return style.Render(label)
}

// GetBorderStyle returns the appropriate border style based on focus
func GetBorderStyle(isFocused bool) lipgloss.Style {
	if isFocused {
		return styles.FocusedBorder
	}
	return styles.UnfocusedBorder
}

// IsFocused checks if a specific section is focused
func IsFocused(current, target types.FocusSection) bool {
	return current == target
}
