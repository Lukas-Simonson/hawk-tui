package ui

import (
	"hawk-tui/internal/styles"
)

// RenderHelp renders the help screen
func RenderHelp(width int) string {
	help := `
Hawk TUI - Keyboard Shortcuts

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
  M            Modified
  A            Added (staged)
  D            Deleted
  ??           Untracked
  MM           Modified (staged and unstaged)
  AM           Added then modified
`

	s := styles.Title.Render("🦅 Hawk TUI - Help") + "\n"
	s += styles.Box.Width(width - 4).Render(help) + "\n"
	s += styles.Help.Render("\nPress any key to return...")

	return s
}
