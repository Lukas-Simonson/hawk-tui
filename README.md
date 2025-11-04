# Hawk TUI - Git Repository Viewer

test
A beautiful terminal user interface for viewing and managing git repositories, built with Go and Bubble Tea.

## Features

- **Segmented Layout**: Three-pane interface with Files, Diff, and Command sections always visible
- **File Status Display**: View all modified, added, deleted, and untracked files in your repository
- **2D Scrolling**: Both vertical and horizontal scrolling for files and diffs
- **Colorized Diff Viewer**: View file changes with syntax highlighting - additions in green, deletions in red
- **Full Content Viewing**: Navigate long file paths and diff lines without truncation
- **Git Command Execution**: Run git commands directly from the TUI (manual execution only)
- **Keyboard Navigation**: Vim-style keyboard shortcuts for efficient navigation
- **Color-Coded Status**: Easy-to-read file status indicators

## Installation

```bash
# Build the application
go build -o hawk-tui

# Run the application
./hawk-tui
```

## Usage

The application must be run from within a git repository.

### Keyboard Shortcuts

#### Section Navigation
- `tab` - Cycle forward through Files → Diff → Command sections
- `shift+tab` - Cycle backward through sections
- `?` - Show/hide help screen
- `q` - Quit application
- `ctrl+c` - Quit application

#### Files Section (when focused)
- `↑`/`k` - Move cursor up (auto-loads diff)
- `↓`/`j` - Move cursor down (auto-loads diff)
- `←`/`h` - Scroll left (for long file paths)
- `→`/`l` - Scroll right (for long file paths)
- `Home` - Jump to start (horizontal)
- `PgUp`/`PgDn` - Scroll page up/down (for large file lists)
- `g` - Jump to top of file list
- `G` - Jump to bottom of file list
- `enter` - Add selected file path to command input
- `d` - View diff for selected file
- `r` - Refresh git status

#### Diff Section (when focused)
- `↑`/`k` - Scroll diff up one line
- `↓`/`j` - Scroll diff down one line
- `←`/`h` - Scroll left (for long lines)
- `→`/`l` - Scroll right (for long lines)
- `0` - Jump to start (horizontal)
- `PgUp`/`PgDn` - Scroll page up/down (10 lines)
- `g`/`Home` - Jump to top of diff
- `G`/`End` - Jump to bottom of diff

#### Command Section (when focused)
- Type your git command (e.g., `add .`, `commit -m "message"`)
  - Note: `git` is automatically prepended - just type the subcommand
- `help` - Show git command reference
- `enter` - Execute command
- `backspace` - Delete character
- `esc` - Clear command input

### File Status Indicators

The status column clearly shows whether changes are staged, unstaged, or both:

- `Staged` - Changes added to staging area (green)
- `Modified` - Unstaged modifications (yellow)
- `Staged + Modified` - File has both staged and unstaged changes
- `Staged-Del` - File deleted and staged (red)
- `Deleted` - File deleted but not staged (red)
- `Untracked` - New file not tracked by git (orange)
- `Renamed` - File has been renamed
- `Copied` - File has been copied

### Git Command Reference

Type `help` in the Command section to view a comprehensive git command reference covering:
- Staging and unstaging changes
- Committing (with message, amend, etc.)
- Branch operations (create, switch, delete)
- Remote operations (push, pull, fetch)
- Viewing status and history
- Stashing changes
- Merging and rebasing
- Undoing changes

The reference is scrollable and styled with colored command names for easy reading.

## Examples

1. **View repository status**: Launch the application to see all file changes in the Files pane
2. **View file diff**: Navigate to a file with `↑`/`↓`, press `d` to toggle the diff view
3. **Scroll through a long diff**: When diff is visible, press `tab` to focus it, then use `↑`/`↓` to scroll
4. **Hide diff**: Press `d` again to hide the diff pane, or navigate to another file
5. **View git commands**: In Command section, type `help` and press `enter` to see git command reference
6. **Stage a file**: Navigate to the file, press `enter` to add its path to command, type `add ` before it, press `enter`
7. **Quick staging**: Type `add` in command section, navigate to file in Files section, press `enter` to add the path
8. **Commit changes**: In Command section, type `commit -m "Your message"`, press `enter`
9. **Refresh status**: Press `r` to update the file list

## Requirements

- Go 1.21 or higher
- Git installed and accessible in PATH
- Must be run from within a git repository

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## Safety

The TUI does not run any commands automatically. All git commands must be manually entered and executed by the user, ensuring you have full control over your repository operations.
