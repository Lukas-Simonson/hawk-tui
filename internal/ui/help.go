package ui

import (
	"fmt"
	"strings"

	"hawk-tui/internal/styles"

	"github.com/charmbracelet/x/ansi"
)

// formatCommand formats a command with aligned description
// cmdWidth is the visual width to pad the command to (accounting for ANSI codes)
func formatCommand(command string, description string, cmdWidth int) string {
	styledCmd := styles.Command.Render(command)
	// Calculate the visual width (without ANSI codes)
	visualWidth := ansi.StringWidth(command)
	// Calculate padding needed
	padding := cmdWidth - visualWidth
	if padding < 0 {
		padding = 0
	}
	return fmt.Sprintf("  %s%s  %s", styledCmd, strings.Repeat(" ", padding), description)
}

// RenderHelp renders the help screen with vertical scrolling
func RenderHelp(width, height, scrollOffset int) string {
	// Build help content with styled components
	var helpLines []string

	// Fixed width for command column (visual characters, not including ANSI codes)
	const cmdWidth = 20

	helpLines = append(helpLines, "Hawk TUI - Keyboard Shortcuts")
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Section Navigation:")
	helpLines = append(helpLines, FormatCommand("tab", "Cycle forward through Files → Diff → Command sections", cmdWidth))
	helpLines = append(helpLines, FormatCommand("shift+tab", "Cycle backward through sections", cmdWidth))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Files Section (when focused):")
	helpLines = append(helpLines, FormatCommand("↑/k", "Move cursor up (auto-loads diff)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("↓/j", "Move cursor down (auto-loads diff)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("←/h", "Scroll left (horizontal)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("→/l", "Scroll right (horizontal)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("home", "Jump to start (horizontal)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("pgup/pgdn", "Scroll page up/down", cmdWidth))
	helpLines = append(helpLines, FormatCommand("g", "Jump to top", cmdWidth))
	helpLines = append(helpLines, FormatCommand("G", "Jump to bottom", cmdWidth))
	helpLines = append(helpLines, FormatCommand("enter", "Add file path to command input", cmdWidth))
	helpLines = append(helpLines, FormatCommand("d", "View diff for selected file", cmdWidth))
	helpLines = append(helpLines, FormatCommand("r", "Refresh git status", cmdWidth))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Diff Section (when focused):")
	helpLines = append(helpLines, FormatCommand("↑/k", "Scroll diff up", cmdWidth))
	helpLines = append(helpLines, FormatCommand("↓/j", "Scroll diff down", cmdWidth))
	helpLines = append(helpLines, FormatCommand("←/h", "Scroll left (horizontal)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("→/l", "Scroll right (horizontal)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("0", "Jump to start (horizontal)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("pgup/pgdn", "Scroll page up/down", cmdWidth))
	helpLines = append(helpLines, FormatCommand("g/home", "Jump to top", cmdWidth))
	helpLines = append(helpLines, FormatCommand("G/end", "Jump to bottom", cmdWidth))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Command Section (when focused):")
	helpLines = append(helpLines, FormatCommand("[type]", "Enter git command (e.g., \"add .\", \"commit -m 'msg'\") ('git' is automatically prepended - just type the subcommand)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("help", "Show git command reference", cmdWidth))
	helpLines = append(helpLines, FormatCommand("enter", "Execute command", cmdWidth))
	helpLines = append(helpLines, FormatCommand("backspace", "Delete character", cmdWidth))
	helpLines = append(helpLines, FormatCommand("esc", "Clear command input", cmdWidth))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "General:")
	helpLines = append(helpLines, FormatCommand("?", "Show/hide this help", cmdWidth))
	helpLines = append(helpLines, FormatCommand("q", "Quit application", cmdWidth))
	helpLines = append(helpLines, FormatCommand("ctrl+c", "Quit application", cmdWidth))
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "File Status Indicators:")
	helpLines = append(helpLines, FormatCommand(styles.StatusAdded.Render("Staged"), "Changes added to staging area", cmdWidth))
	helpLines = append(helpLines, FormatCommand(styles.StatusModified.Render("Modified"), "Unstaged modifications", cmdWidth))
	helpLines = append(helpLines, FormatCommand(styles.StatusAdded.Render("Staged")+" +"+styles.StatusModified.Render(" Modified"), "Both staged and unstaged changes", cmdWidth))
	helpLines = append(helpLines, FormatCommand(styles.StatusDeleted.Render("Staged-Del"), "Deleted and staged", cmdWidth))
	helpLines = append(helpLines, FormatCommand(styles.StatusDeleted.Render("Deleted"), "Deleted but not staged", cmdWidth))
	helpLines = append(helpLines, FormatCommand(styles.StatusUntracked.Render("Untracked"), "New file not tracked by git", cmdWidth))
	helpLines = append(helpLines, FormatCommand(styles.Command.Render("Renamed"), "File has been renamed", cmdWidth))
	helpLines = append(helpLines, FormatCommand(styles.Command.Render("Copied"), "File has been copied", cmdWidth))

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
	s += styles.Box.Width(width-4).Render(visibleContent) + "\n"

	// Add scroll indicator if needed
	if totalLines > visibleLines {
		scrollInfo := fmt.Sprintf("↑/↓ scroll [%d-%d/%d] • ", startLine+1, endLine, totalLines)
		s += styles.Help.Render(scrollInfo + "Press any key to return...")
	} else {
		s += styles.Help.Render("Press any key to return...")
	}

	return s
}
