package model

import (
	"strings"

	"hawk-tui/internal/git"
	"hawk-tui/internal/types"
	"hawk-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the application state
type Model struct {
	Files        []types.FileState
	Cursor       int
	FocusSection types.FocusSection
	DiffContent  string
	CommandInput string
	PopupOutput  string // Output shown in popup overlay
	Width        int
	Height       int
	Err          error
	DiffScroll   int
	FilesHScroll int // Horizontal scroll offset for files pane
	FilesVScroll int // Vertical scroll offset for files pane
	DiffHScroll  int // Horizontal scroll offset for diff pane
}

// New creates a new model instance
func New() Model {
	return Model{
		Files:        []types.FileState{},
		Cursor:       0,
		FocusSection: types.FocusFiles,
		DiffScroll:   0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return git.RefreshStatus()
	}
}

// View renders the UI
func (m Model) View() string {
	if m.FocusSection == types.FocusHelp {
		return ui.RenderHelp(m.Width)
	}

	if m.FocusSection == types.FocusPopup {
		return ui.RenderPopup(m.PopupOutput, m.Width, m.Height)
	}

	// Calculate dimensions for the panes
	commandHeight := 6
	topHeight := m.Height - commandHeight - 2
	halfWidth := m.Width / 2

	// Render each pane
	filesPane := ui.RenderFilesPane(m.Files, m.Cursor, m.FocusSection, m.FilesHScroll, m.FilesVScroll, m.Err, halfWidth, topHeight, m.calculateFilesScrollLimits)
	diffPane := ui.RenderDiffPane(m.Files, m.Cursor, m.FocusSection, m.DiffContent, m.DiffScroll, m.DiffHScroll, halfWidth, topHeight, m.calculateDiffScrollLimits)
	commandPane := ui.RenderCommandPane(m.FocusSection, m.CommandInput, m.Width, commandHeight)

	// Join horizontally: files | diff
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, filesPane, diffPane)

	// Join vertically: topRow over command
	return lipgloss.JoinVertical(lipgloss.Left, topRow, commandPane)
}

// calculateFilesScrollLimits calculates the visible lines and max scroll for the files pane
func (m Model) calculateFilesScrollLimits() (visibleLines, maxScroll int) {
	if len(m.Files) == 0 {
		return 0, 0
	}

	totalFiles := len(m.Files)

	// Calculate pane height (same as in View())
	commandHeight := 6
	topHeight := m.Height - commandHeight - 2

	// Calculate visible lines
	// Title (2 lines) + file count (2 lines) + borders/padding (4 lines)
	baseOverhead := 8
	scrollIndicatorLines := 2

	visibleLines = topHeight - baseOverhead
	if totalFiles > visibleLines {
		visibleLines -= scrollIndicatorLines
	}
	if visibleLines < 1 {
		visibleLines = 1
	}

	maxScroll = totalFiles - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	return visibleLines, maxScroll
}

// calculateDiffScrollLimits calculates the visible lines and max scroll for the diff pane
func (m Model) calculateDiffScrollLimits() (visibleLines, maxScroll int) {
	if m.DiffContent == "" {
		return 0, 0
	}

	lines := strings.Split(m.DiffContent, "\n")
	totalLines := len(lines)

	// Calculate pane height (same as in View())
	commandHeight := 6
	topHeight := m.Height - commandHeight - 2

	// Calculate visible lines (same logic as renderDiffPane)
	baseOverhead := 8
	scrollIndicatorLines := 3

	visibleLines = topHeight - baseOverhead
	if totalLines > visibleLines {
		visibleLines -= scrollIndicatorLines
	}
	if visibleLines < 1 {
		visibleLines = 1
	}

	maxScroll = totalLines - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	return visibleLines, maxScroll
}

// ShowDiff returns a command to show the diff for the current file
func (m Model) ShowDiff() tea.Cmd {
	return func() tea.Msg {
		if len(m.Files) == 0 || m.Cursor >= len(m.Files) {
			return types.DiffMsg{Content: "No file selected"}
		}

		file := m.Files[m.Cursor]
		return git.GetDiff(file)
	}
}

// ExecuteCommand returns a command to execute the git command
func (m Model) ExecuteCommand() tea.Cmd {
	return func() tea.Msg {
		return git.ExecuteCommand(m.CommandInput)
	}
}

// RefreshGitStatus returns a command to refresh the git status
func RefreshGitStatus() tea.Cmd {
	return func() tea.Msg {
		return git.RefreshStatus()
	}
}
