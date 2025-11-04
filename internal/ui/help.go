package ui

import (
	"fmt"
	"strings"

	"hawk-tui/internal/styles"
)

// RenderHelp renders the help screen with vertical scrolling
func RenderHelp(width, height, scrollOffset int) string {
	helpContent := `Hawk TUI - Keyboard Shortcuts

Section Navigation:
  tab          Cycle forward through Files → Diff → Command sections
  shift+tab    Cycle backward through sections

Files Section (when focused):
  ↑/k          Move cursor up (auto-loads diff)
  ↓/j          Move cursor down (auto-loads diff)
  ←/h          Scroll left (horizontal)
  →/l          Scroll right (horizontal)
  home         Jump to start (horizontal)
  pgup/pgdn    Scroll page up/down
  g            Jump to top
  G            Jump to bottom
  enter        Add file path to command input
  d            View diff for selected file
  r            Refresh git status

Diff Section (when focused):
  ↑/k          Scroll diff up
  ↓/j          Scroll diff down
  ←/h          Scroll left (horizontal)
  →/l          Scroll right (horizontal)
  0            Jump to start (horizontal)
  pgup/pgdn    Scroll page up/down
  g/home       Jump to top
  G/end        Jump to bottom

Command Section (when focused):
  [type]       Enter git command (e.g., "add .", "commit -m 'msg'")
               'git' is automatically prepended - just type the subcommand
  enter        Execute command
  backspace    Delete character
  esc          Clear command input

General:
  ?            Show/hide this help
  q            Quit application
  ctrl+c       Quit application

File Status Indicators:
  Staged             Changes added to staging area (green)
  Modified           Unstaged modifications (yellow)
  Staged + Modified  Both staged and unstaged changes
  Staged-Del         Deleted and staged (red)
  Deleted            Deleted but not staged (red)
  Untracked          New file not tracked by git (orange)
  Renamed            File has been renamed
  Copied             File has been copied`

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
