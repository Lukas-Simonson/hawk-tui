package model

import (
	"strings"

	"hawk-tui/internal/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Global keys
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.FocusSection != types.FocusCommand {
				return m, tea.Quit
			}
		case "?":
			if m.FocusSection != types.FocusCommand {
				m.FocusSection = types.FocusHelp
				return m, nil
			}
		case "tab":
			// Cycle forward through Files -> Diff -> Command -> Files
			switch m.FocusSection {
			case types.FocusFiles:
				m.FocusSection = types.FocusDiff
			case types.FocusDiff:
				m.FocusSection = types.FocusCommand
			case types.FocusCommand:
				m.FocusSection = types.FocusFiles
			}
			return m, nil
		case "shift+tab":
			// Cycle backward through Files -> Command -> Diff -> Files
			switch m.FocusSection {
			case types.FocusFiles:
				m.FocusSection = types.FocusCommand
			case types.FocusCommand:
				m.FocusSection = types.FocusDiff
			case types.FocusDiff:
				m.FocusSection = types.FocusFiles
			}
			return m, nil
		}

		// Section-specific keys
		switch m.FocusSection {
		case types.FocusFiles:
			return m.updateFileList(msg)
		case types.FocusDiff:
			return m.updateDiff(msg)
		case types.FocusCommand:
			return m.updateCommand(msg)
		case types.FocusHelp:
			return m.updateHelp(msg)
		case types.FocusPopup:
			return m.updatePopup(msg)
		}

	case types.FilesMsg:
		m.Files = msg.Files
		if len(m.Files) > 0 && m.Cursor >= len(m.Files) {
			m.Cursor = len(m.Files) - 1
		}
		return m, nil

	case types.DiffMsg:
		m.DiffContent = msg.Content
		m.DiffScroll = 0
		m.DiffHScroll = 0
		return m, nil

	case types.CommandOutputMsg:
		m.PopupOutput = msg.Output
		m.CommandInput = ""
		m.FocusSection = types.FocusPopup
		return m, RefreshGitStatus()

	case types.ErrMsg:
		m.Err = msg.Err
		return m, nil
	}

	return m, nil
}

// updateFileList handles updates for the files section
func (m Model) updateFileList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visibleLines, maxScroll := m.calculateFilesScrollLimits()

	switch msg.String() {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			// Auto-scroll up if cursor goes above visible area
			if m.Cursor < m.FilesVScroll {
				m.FilesVScroll = m.Cursor
			}
			if len(m.Files) > 0 {
				return m, m.ShowDiff()
			}
		}

	case "down", "j":
		if m.Cursor < len(m.Files)-1 {
			m.Cursor++
			// Auto-scroll down if cursor goes below visible area
			if m.Cursor >= m.FilesVScroll+visibleLines {
				m.FilesVScroll = m.Cursor - visibleLines + 1
				if m.FilesVScroll > maxScroll {
					m.FilesVScroll = maxScroll
				}
			}
			if len(m.Files) > 0 {
				return m, m.ShowDiff()
			}
		}

	case "pageup", "pgup":
		m.FilesVScroll -= 10
		if m.FilesVScroll < 0 {
			m.FilesVScroll = 0
		}
		// Move cursor with scroll
		m.Cursor = m.FilesVScroll

	case "pagedown", "pgdown":
		m.FilesVScroll += 10
		if m.FilesVScroll > maxScroll {
			m.FilesVScroll = maxScroll
		}
		// Move cursor with scroll
		m.Cursor = m.FilesVScroll

	case "g":
		// Jump to top
		m.Cursor = 0
		m.FilesVScroll = 0
		if len(m.Files) > 0 {
			return m, m.ShowDiff()
		}

	case "G":
		// Jump to bottom
		m.Cursor = len(m.Files) - 1
		m.FilesVScroll = maxScroll
		if len(m.Files) > 0 {
			return m, m.ShowDiff()
		}

	case "left", "h":
		// Scroll left
		if m.FilesHScroll > 0 {
			m.FilesHScroll -= 5
			if m.FilesHScroll < 0 {
				m.FilesHScroll = 0
			}
		}

	case "right", "l":
		// Scroll right
		m.FilesHScroll += 5

	case "home":
		// Jump to start (horizontal)
		m.FilesHScroll = 0

	case "enter":
		// Add selected file path to command input and switch to command section
		if len(m.Files) > 0 && m.Cursor < len(m.Files) {
			file := m.Files[m.Cursor]
			// Add space before path if there's already content
			if m.CommandInput != "" && !strings.HasSuffix(m.CommandInput, " ") {
				m.CommandInput += " "
			}
			m.CommandInput += file.Path
			m.FocusSection = types.FocusCommand
		}

	case "d":
		// View diff for selected file
		if len(m.Files) > 0 {
			return m, m.ShowDiff()
		}

	case "r":
		m.DiffContent = ""
		m.FilesHScroll = 0
		m.FilesVScroll = 0
		return m, RefreshGitStatus()
	}

	return m, nil
}

// updateDiff handles updates for the diff section
func (m Model) updateDiff(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	_, maxScroll := m.calculateDiffScrollLimits()

	switch msg.String() {
	case "up", "k":
		if m.DiffScroll > 0 {
			m.DiffScroll--
		}

	case "down", "j":
		if m.DiffScroll < maxScroll {
			m.DiffScroll++
		}

	case "left", "h":
		// Scroll left horizontally
		if m.DiffHScroll > 0 {
			m.DiffHScroll -= 5
			if m.DiffHScroll < 0 {
				m.DiffHScroll = 0
			}
		}

	case "right", "l":
		// Scroll right horizontally
		m.DiffHScroll += 5

	case "pageup", "pgup":
		m.DiffScroll -= 10
		if m.DiffScroll < 0 {
			m.DiffScroll = 0
		}

	case "pagedown", "pgdown":
		m.DiffScroll += 10
		if m.DiffScroll > maxScroll {
			m.DiffScroll = maxScroll
		}

	case "home", "g":
		m.DiffScroll = 0

	case "end", "G":
		m.DiffScroll = maxScroll

	case "0":
		// Jump to start horizontally (vim-style)
		m.DiffHScroll = 0
	}

	return m, nil
}

// updateCommand handles updates for the command section
func (m Model) updateCommand(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CommandInput = ""

	case "enter":
		if m.CommandInput != "" {
			return m, m.ExecuteCommand()
		}

	case "backspace":
		if len(m.CommandInput) > 0 {
			m.CommandInput = m.CommandInput[:len(m.CommandInput)-1]
		}

	case "space":
		m.CommandInput += " "

	default:
		// Handle all printable characters
		if len(msg.String()) == 1 && msg.String() >= " " && msg.String() <= "~" {
			m.CommandInput += msg.String()
		}
	}

	return m, nil
}

// updateHelp handles updates for the help screen
func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any key press returns to the files view
	m.FocusSection = types.FocusFiles
	return m, nil
}

// updatePopup handles updates for the popup screen
func (m Model) updatePopup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any key press closes the popup and returns to the files view
	m.FocusSection = types.FocusFiles
	m.PopupOutput = ""
	return m, nil
}
