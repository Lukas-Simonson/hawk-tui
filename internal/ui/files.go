package ui

import (
	"fmt"
	"strings"

	"hawk-tui/internal/styles"
	"hawk-tui/internal/types"

	"github.com/charmbracelet/x/ansi"
)

// RenderFilesPane renders the files section
func RenderFilesPane(files []types.FileState, cursor int, focusSection types.FocusSection, filesHScroll, filesVScroll int, err error, width, height int, calculateScrollLimits func() (int, int)) string {
	borderStyle := GetBorderStyle(IsFocused(focusSection, types.FocusFiles))

	var content strings.Builder
	content.WriteString(styles.Title.Render("📁 Files") + "\n\n")

	if err != nil {
		content.WriteString(styles.Error.Render(fmt.Sprintf("Error: %v", err)) + "\n\n")
	}

	if len(files) == 0 {
		content.WriteString(styles.Normal.Render("No changes detected.\nWorking tree clean."))
	} else {
		// Calculate viewport width (account for borders and padding)
		viewportWidth := width - 10
		if viewportWidth < 10 {
			viewportWidth = 10
		}

		totalFiles := len(files)
		visibleLines, _ := calculateScrollLimits()

		// Calculate visible range (vertical)
		startFile := filesVScroll
		endFile := startFile + visibleLines

		if endFile > totalFiles {
			endFile = totalFiles
		}
		if startFile >= totalFiles {
			startFile = totalFiles - 1
			if startFile < 0 {
				startFile = 0
			}
		}

		hasHScrollableContent := false
		content.WriteString(fmt.Sprintf("%d file(s) changed\n", len(files)))

		for i := startFile; i < endFile; i++ {
			file := files[i]
			isSelected := cursor == i

			// Build components
			var cursorIcon string
			if isSelected {
				cursorIcon = "▶"
			} else {
				cursorIcon = " "
			}

			// Build plain text line for width calculation
			plainLine := fmt.Sprintf("%s [%s] %s", cursorIcon, file.Status, file.Path)
			lineWidth := ansi.StringWidth(plainLine)
			if lineWidth > viewportWidth {
				hasHScrollableContent = true
			}

			// Now build the styled version by applying colors to individual parts
			// Get the colored status
			statusStyled := FormatStatus(file.Status)

			// Build the final styled line by combining parts
			var styledLine string
			if isSelected {
				// For selected lines, apply selected color to cursor and path only
				styledLine = fmt.Sprintf("%s [%s] %s",
					styles.Selected.Render(cursorIcon),
					statusStyled,
					styles.Selected.Render(file.Path))
			} else {
				// For unselected lines, only color the status
				styledLine = fmt.Sprintf("%s [%s] %s", cursorIcon, statusStyled, file.Path)
			}

			// Apply horizontal scroll to the styled line
			visibleStyledLine := ApplyHorizontalViewport(styledLine, filesHScroll, viewportWidth)

			content.WriteString(visibleStyledLine + "\n")
		}

		// Add scroll indicators
		var indicators []string
		if totalFiles > visibleLines {
			indicators = append(indicators, fmt.Sprintf("↑/↓ scroll [%d-%d/%d]",
				startFile+1,
				endFile,
				totalFiles))
		}
		if hasHScrollableContent {
			indicators = append(indicators, "←/→ scroll horizontally")
		}
		if len(indicators) > 0 {
			content.WriteString("\n" + styles.Help.Render(strings.Join(indicators, " • ")))
		}
	}

	return borderStyle.
		Width(width - 4).
		Height(height - 2).
		Render(content.String())
}
