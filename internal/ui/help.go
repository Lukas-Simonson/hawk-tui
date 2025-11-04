package ui

import (
	"fmt"
	"strings"

	"hawk-tui/internal/styles"
)

// RenderHelp renders the help screen with vertical scrolling
func RenderHelp(width, height, scrollOffset int) string {
	// Build help content with styled components
	var helpLines []string

	helpLines = append(helpLines, "Hawk TUI - Keyboard Shortcuts")
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Section Navigation:")
	helpLines = append(helpLines, fmt.Sprintf("  %s          Cycle forward through Files → Diff → Command sections", styles.Command.Render("tab")))
	helpLines = append(helpLines, fmt.Sprintf("  %s    Cycle backward through sections", styles.Command.Render("shift+tab")))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Files Section (when focused):")
	helpLines = append(helpLines, fmt.Sprintf("  %s          Move cursor up (auto-loads diff)", styles.Command.Render("↑/k")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          Move cursor down (auto-loads diff)", styles.Command.Render("↓/j")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          Scroll left (horizontal)", styles.Command.Render("←/h")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          Scroll right (horizontal)", styles.Command.Render("→/l")))
	helpLines = append(helpLines, fmt.Sprintf("  %s         Jump to start (horizontal)", styles.Command.Render("home")))
	helpLines = append(helpLines, fmt.Sprintf("  %s    Scroll page up/down", styles.Command.Render("pgup/pgdn")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            Jump to top", styles.Command.Render("g")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            Jump to bottom", styles.Command.Render("G")))
	helpLines = append(helpLines, fmt.Sprintf("  %s        Add file path to command input", styles.Command.Render("enter")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            View diff for selected file", styles.Command.Render("d")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            Refresh git status", styles.Command.Render("r")))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Diff Section (when focused):")
	helpLines = append(helpLines, fmt.Sprintf("  %s          Scroll diff up", styles.Command.Render("↑/k")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          Scroll diff down", styles.Command.Render("↓/j")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          Scroll left (horizontal)", styles.Command.Render("←/h")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          Scroll right (horizontal)", styles.Command.Render("→/l")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            Jump to start (horizontal)", styles.Command.Render("0")))
	helpLines = append(helpLines, fmt.Sprintf("  %s    Scroll page up/down", styles.Command.Render("pgup/pgdn")))
	helpLines = append(helpLines, fmt.Sprintf("  %s       Jump to top", styles.Command.Render("g/home")))
	helpLines = append(helpLines, fmt.Sprintf("  %s         Jump to bottom", styles.Command.Render("G/end")))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Command Section (when focused):")
	helpLines = append(helpLines, fmt.Sprintf("  %s       Enter git command (e.g., \"add .\", \"commit -m 'msg'\")", styles.Command.Render("[type]")))
	helpLines = append(helpLines, "               'git' is automatically prepended - just type the subcommand")
	helpLines = append(helpLines, fmt.Sprintf("  %s        Execute command", styles.Command.Render("enter")))
	helpLines = append(helpLines, fmt.Sprintf("  %s    Delete character", styles.Command.Render("backspace")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          Clear command input", styles.Command.Render("esc")))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "General:")
	helpLines = append(helpLines, fmt.Sprintf("  %s            Show/hide this help", styles.Command.Render("?")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            Quit application", styles.Command.Render("q")))
	helpLines = append(helpLines, fmt.Sprintf("  %s       Quit application", styles.Command.Render("ctrl+c")))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "File Status Indicators:")
	helpLines = append(helpLines, fmt.Sprintf("  %s             Changes added to staging area", styles.StatusAdded.Render("Staged")))
	helpLines = append(helpLines, fmt.Sprintf("  %s           Unstaged modifications", styles.StatusModified.Render("Modified")))
	helpLines = append(helpLines, fmt.Sprintf("  %s  Both staged and unstaged changes", styles.StatusAdded.Render("Staged")+" +"+styles.StatusModified.Render(" Modified")))
	helpLines = append(helpLines, fmt.Sprintf("  %s         Deleted and staged", styles.StatusDeleted.Render("Staged-Del")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            Deleted but not staged", styles.StatusDeleted.Render("Deleted")))
	helpLines = append(helpLines, fmt.Sprintf("  %s          New file not tracked by git", styles.StatusUntracked.Render("Untracked")))
	helpLines = append(helpLines, fmt.Sprintf("  %s            File has been renamed", styles.Command.Render("Renamed")))
	helpLines = append(helpLines, fmt.Sprintf("  %s             File has been copied", styles.Command.Render("Copied")))

	helpContent := strings.Join(helpLines, "\n")

	// Split into lines
	lines := strings.Split(helpContent, "\n")
	totalLines := len(lines)

	// Calculate visible area (accounting for title and help text)
	visibleLines := height - 6 // Reserve space for title and footer

	if visibleLines < 10 {
		visibleLines = 10
	}

	// Calculate viewport
	startLine := scrollOffset
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

	// Get visible lines
	visibleContent := strings.Join(lines[startLine:endLine], "\n")

	// Build the display
	s := styles.Title.Render("🦅 Hawk TUI - Help") + "\n"
	s += styles.Box.Width(width - 4).Render(visibleContent) + "\n"

	// Add scroll indicator if needed
	if totalLines > visibleLines {
		scrollInfo := fmt.Sprintf("↑/↓ scroll [%d-%d/%d] • ", startLine+1, endLine, totalLines)
		s += styles.Help.Render(scrollInfo + "Press any key to return...")
	} else {
		s += styles.Help.Render("Press any key to return...")
	}

	return s
}
