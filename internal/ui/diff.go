package ui

import (
	"fmt"
	"strings"

	"hawk-tui/internal/styles"
	"hawk-tui/internal/types"

	"github.com/charmbracelet/x/ansi"
)

// RenderDiffPane renders the diff section
func RenderDiffPane(files []types.FileState, cursor int, focusSection types.FocusSection, diffContent string, diffScroll, diffHScroll int, width, height int, calculateScrollLimits func() (int, int)) string {
	borderStyle := GetBorderStyle(IsFocused(focusSection, types.FocusDiff))

	var content strings.Builder
	content.WriteString(styles.Title.Render("Diff") + "\n\n")

	if len(files) > 0 && cursor < len(files) {
		file := files[cursor]
		fileInfo := fmt.Sprintf("%s [%s]", file.Path, FormatStatus(file.Status))
		content.WriteString(styles.Help.Render(fileInfo) + "\n\n")
	}

	if diffContent == "" {
		content.WriteString("Select a file to view diff\n↑/↓ to navigate")
	} else {
		// Calculate viewport width (account for borders and padding)
		viewportWidth := width - 10
		if viewportWidth < 10 {
			viewportWidth = 10
		}

		// Split content into lines and handle scrolling
		lines := strings.Split(diffContent, "\n")
		totalLines := len(lines)

		// Use the scroll limits calculation
		visibleLines, _ := calculateScrollLimits()

		// Calculate visible range (vertical)
		startLine := diffScroll
		endLine := startLine + visibleLines

		if endLine > totalLines {
			endLine = totalLines
		}
		if startLine >= totalLines {
			startLine = totalLines - 1
			if startLine < 0 {
				startLine = 0
			}
		}

		// Ensure we don't scroll past the end
		if startLine > 0 && endLine == totalLines {
			maxStart := totalLines - visibleLines
			if maxStart < 0 {
				maxStart = 0
			}
			if startLine > maxStart {
				startLine = maxStart
				endLine = totalLines
			}
		}

		// Apply horizontal viewport to each visible line
		hasHScrollableContent := false
		var viewportLines []string
		for i := startLine; i < endLine; i++ {
			line := lines[i]
			lineWidth := ansi.StringWidth(line)
			if lineWidth > viewportWidth {
				hasHScrollableContent = true
			}
			visibleLine := ApplyHorizontalViewport(line, diffHScroll, viewportWidth)
			viewportLines = append(viewportLines, visibleLine)
		}

		// Join and render
		visibleContent := strings.Join(viewportLines, "\n")
		content.WriteString(visibleContent)

		// Add scroll indicators
		var indicators []string
		if totalLines > visibleLines {
			indicators = append(indicators, fmt.Sprintf("↑/↓ scroll [%d-%d/%d]",
				startLine+1,
				endLine,
				totalLines))
		}
		if hasHScrollableContent {
			indicators = append(indicators, "←/→ scroll horizontally • 0 to reset")
		}
		if len(indicators) > 0 {
			content.WriteString("\n\n" + strings.Join(indicators, " • "))
		}
	}

	return borderStyle.
		Width(width - 4).
		Height(height - 2).
		Render(content.String())
}
