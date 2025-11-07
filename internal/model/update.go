package model

import (
	"strconv"
	"strings"

	"hawk-tui/internal/audio"
	"hawk-tui/internal/git"
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
		case types.FocusGitHelp:
			return m.updateGitHelp(msg)
		case types.FocusPopup:
			return m.updatePopup(msg)
		case types.FocusCommitPopup:
			return m.updateCommitPopup(msg)
		}

	case types.FilesMsg:
		m.Files = msg.Files
		if len(m.Files) > 0 && m.Cursor >= len(m.Files) {
			m.Cursor = len(m.Files) - 1
		}
		// Clear diff if there are no files
		if len(m.Files) == 0 {
			m.DiffContent = ""
			return m, nil
		}
		// Auto-load diff for first file
		if len(m.Files) > 0 {
			return m, m.ShowDiff()
		}
		return m, nil

	case types.GitStatusMsg:
		m.GitStatus = msg
		return m, nil

	case types.DiffMsg:
		m.DiffContent = msg.Content
		m.DiffScroll = 0
		m.DiffHScroll = 0
		return m, nil

	case types.CommandOutputMsg:
		m.CommandInput = ""
		// Only show popup if there's output
		if msg.Output != "" {
			m.PopupOutput = msg.Output
			m.FocusSection = types.FocusPopup
		}
		return m, tea.Batch(RefreshGitStatus(), RefreshGitBranchInfo())

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
			// Quote the path if it contains spaces or special characters
			m.CommandInput += quotePathIfNeeded(file.Path)
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
		return m, tea.Batch(RefreshGitStatus(), RefreshGitBranchInfo())
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

// parseCommitMessage extracts the commit message and other flags from a commit command
// Example: "-a -m "my message"" -> ("my message", "-a")
func (m Model) parseCommitMessage(input string) (message string, flags string) {
	// Find the -m flag
	mIndex := strings.Index(input, "-m")
	if mIndex == -1 {
		// Try --message
		mIndex = strings.Index(input, "--message")
		if mIndex == -1 {
			return "", input
		}
		mIndex += len("--message")
	} else {
		mIndex += len("-m")
	}

	// Extract flags before -m
	flagsBefore := strings.TrimSpace(input[:mIndex-len("-m")])

	// Find the message after -m
	rest := strings.TrimSpace(input[mIndex:])
	if rest == "" {
		return "", flagsBefore
	}

	// Check if message is quoted
	if rest[0] == '"' {
		// Find closing quote, handling escaped quotes
		endIdx := 1
		for endIdx < len(rest) {
			if rest[endIdx] == '"' && (endIdx == 1 || rest[endIdx-1] != '\\') {
				message = rest[1:endIdx]
				// Get any flags after the message
				flagsAfter := strings.TrimSpace(rest[endIdx+1:])
				if flagsAfter != "" {
					if flagsBefore != "" {
						flags = flagsBefore + " " + flagsAfter
					} else {
						flags = flagsAfter
					}
				} else {
					flags = flagsBefore
				}
				return message, flags
			}
			endIdx++
		}
		// No closing quote found, return empty
		return "", flagsBefore
	}

	// Not quoted, take until next space or flag
	parts := strings.Fields(rest)
	if len(parts) > 0 {
		message = parts[0]
		if len(parts) > 1 {
			flagsAfter := strings.Join(parts[1:], " ")
			if flagsBefore != "" {
				flags = flagsBefore + " " + flagsAfter
			} else {
				flags = flagsAfter
			}
		} else {
			flags = flagsBefore
		}
	}

	return message, flags
}

// updateCommand handles updates for the command section
func (m Model) updateCommand(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CommandInput = ""

	case "enter":
		if m.CommandInput != "" {
			input := strings.TrimSpace(m.CommandInput)

			// Check if user is based
			if input == "tuah" {
				go audio.PlaySound("assets/hawk_tuah.mp3")
				m.CommandInput = ""
				return m, nil
			}

			// Check if this is the help command
			if input == "help" {
				m.CommandInput = ""
				m.FocusSection = types.FocusGitHelp
				return m, nil
			}

			// Check if this is a commit command
			if strings.HasPrefix(input, "commit") {
				// Extract flags after "commit" (e.g., commit -a, commit --amend)
				rest := strings.TrimPrefix(input, "commit")
				rest = strings.TrimSpace(rest)

				// Check if flags contain -m or --message
				if strings.Contains(rest, "-m") || strings.Contains(rest, "--message") {
					// Parse the message from the command
					message, flags := m.parseCommitMessage(rest)
					if message != "" {
						m.CommandInput = ""
						return m, func() tea.Msg {
							return git.ExecuteCommit(message, flags)
						}
					}
					// If parsing failed, show popup
					m.CommitFlags = flags
					m.CommitMessage = ""
					m.FocusSection = types.FocusCommitPopup
					return m, nil
				}

				// Otherwise, show the commit popup
				m.CommitFlags = rest
				m.CommitMessage = ""
				m.FocusSection = types.FocusCommitPopup
				return m, nil
			}
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
	// Calculate max scroll based on help content
	// Help has ~60 lines, reserve 6 for title/footer
	totalLines := 60
	visibleLines := m.Height - 6
	if visibleLines < 10 {
		visibleLines = 10
	}
	maxScroll := totalLines - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Handle scrolling keys
	switch msg.String() {
	case "up", "k":
		if m.HelpScroll > 0 {
			m.HelpScroll--
		}
		return m, nil

	case "down", "j":
		if m.HelpScroll < maxScroll {
			m.HelpScroll++
		}
		return m, nil

	case "pageup", "pgup":
		m.HelpScroll -= 10
		if m.HelpScroll < 0 {
			m.HelpScroll = 0
		}
		return m, nil

	case "pagedown", "pgdown":
		m.HelpScroll += 10
		if m.HelpScroll > maxScroll {
			m.HelpScroll = maxScroll
		}
		return m, nil

	case "g", "home":
		m.HelpScroll = 0
		return m, nil

	case "G", "end":
		m.HelpScroll = maxScroll
		return m, nil

	default:
		// Any other key press returns to the files view
		m.FocusSection = types.FocusFiles
		m.HelpScroll = 0 // Reset scroll when closing
		return m, nil
	}
}

// updateGitHelp handles updates for the git help screen
func (m Model) updateGitHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Calculate max scroll based on git help content
	// Git help has ~100 lines, reserve 6 for title/footer
	totalLines := 100
	visibleLines := m.Height - 6
	if visibleLines < 10 {
		visibleLines = 10
	}
	maxScroll := totalLines - visibleLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Handle scrolling keys
	switch msg.String() {
	case "up", "k":
		if m.GitHelpScroll > 0 {
			m.GitHelpScroll--
		}
		return m, nil

	case "down", "j":
		if m.GitHelpScroll < maxScroll {
			m.GitHelpScroll++
		}
		return m, nil

	case "pageup", "pgup":
		m.GitHelpScroll -= 10
		if m.GitHelpScroll < 0 {
			m.GitHelpScroll = 0
		}
		return m, nil

	case "pagedown", "pgdown":
		m.GitHelpScroll += 10
		if m.GitHelpScroll > maxScroll {
			m.GitHelpScroll = maxScroll
		}
		return m, nil

	case "g", "home":
		m.GitHelpScroll = 0
		return m, nil

	case "G", "end":
		m.GitHelpScroll = maxScroll
		return m, nil

	default:
		// Any other key press returns to the files view
		m.FocusSection = types.FocusFiles
		m.GitHelpScroll = 0 // Reset scroll when closing
		return m, nil
	}
}

// updatePopup handles updates for the popup screen
func (m Model) updatePopup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any key press closes the popup and returns to the files view
	m.FocusSection = types.FocusFiles
	m.PopupOutput = ""
	return m, nil
}

// updateCommitPopup handles updates for the commit message popup
func (m Model) updateCommitPopup(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel commit
		m.FocusSection = types.FocusFiles
		m.CommitMessage = ""
		m.CommitFlags = ""
		m.CommandInput = ""
		return m, nil

	case "ctrl+s":
		// Execute commit with the message
		if m.CommitMessage != "" {
			message := m.CommitMessage
			flags := m.CommitFlags

			// Clear state
			m.FocusSection = types.FocusFiles
			m.CommitMessage = ""
			m.CommitFlags = ""
			m.CommandInput = ""

			// Execute the commit command using the dedicated function
			return m, func() tea.Msg {
				return git.ExecuteCommit(message, flags)
			}
		}
		return m, nil

	case "backspace":
		if len(m.CommitMessage) > 0 {
			m.CommitMessage = m.CommitMessage[:len(m.CommitMessage)-1]
		}

	case "enter":
		// Add newline
		m.CommitMessage += "\n"

	case "space":
		m.CommitMessage += " "

	default:
		// Handle all printable characters
		if len(msg.String()) == 1 && msg.String() >= " " && msg.String() <= "~" {
			m.CommitMessage += msg.String()
		}
	}

	return m, nil
}

// quotePathIfNeeded quotes a file path if it contains spaces or special characters
func quotePathIfNeeded(path string) string {
	// Check if path contains characters that need quoting
	needsQuoting := false
	for _, char := range path {
		if char == ' ' || char == '\t' || char == '\n' || char == '"' || char == '\'' || char == '\\' {
			needsQuoting = true
			break
		}
	}

	if needsQuoting {
		// Use strconv.Quote to properly escape the path
		return strconv.Quote(path)
	}

	return path
}
