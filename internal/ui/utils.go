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

// FormatStatus returns a styled status indicator with clear staging information
// Git status format: XY where X=staged, Y=unstaged (space means no change)
func FormatStatus(status string) string {
	if len(status) == 0 {
		return ""
	}

	// Pad status to 2 characters if needed
	for len(status) < 2 {
		status += " "
	}

	// Handle special cases
	if status == "??" {
		return styles.StatusUntracked.Render("Untracked")
	}

	var parts []string

	// Parse first character (staged/index status) - space means NOT staged
	if status[0] != ' ' && status[0] != '?' {
		switch status[0] {
		case 'M':
			parts = append(parts, styles.StatusAdded.Render("Staged"))
		case 'A':
			parts = append(parts, styles.StatusAdded.Render("Added"))
		case 'D':
			parts = append(parts, styles.StatusDeleted.Render("Staged-Del"))
		case 'R':
			parts = append(parts, styles.StatusAdded.Render("Renamed"))
		case 'C':
			parts = append(parts, styles.StatusAdded.Render("Copied"))
		}
	}

	// Parse second character (working tree/unstaged status) - space means clean
	if status[1] != ' ' && status[1] != '?' {
		switch status[1] {
		case 'M':
			parts = append(parts, styles.StatusModified.Render("Modified"))
		case 'D':
			parts = append(parts, styles.StatusDeleted.Render("Deleted"))
		case 'A':
			parts = append(parts, styles.StatusModified.Render("Added-Unstaged"))
		}
	}

	// If no parts, show the raw status
	if len(parts) == 0 {
		return styles.Normal.Render(status)
	}

	// Join parts with " + " if both staged and unstaged
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " + "
		}
		result += part
	}

	return result
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
