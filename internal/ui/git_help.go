package ui

import (
	"fmt"
	"strings"

	"hawk-tui/internal/styles"
)

// RenderGitHelp renders the git command reference with vertical scrolling
func RenderGitHelp(width, height, scrollOffset int) string {
	// Build git help content with styled components
	var helpLines []string

	// Fixed width for command column (visual characters, not including ANSI codes)
	const cmdWidth = 30

	helpLines = append(helpLines, "Git Command Reference")
	helpLines = append(helpLines, "")
	helpLines = append(helpLines, "Common Git Commands:")
	helpLines = append(helpLines, "")

	// Staging Commands
	helpLines = append(helpLines, styles.Title.Render("Staging Changes:"))
	helpLines = append(helpLines, FormatCommand("add .", "Stage all changes", cmdWidth))
	helpLines = append(helpLines, FormatCommand("add <file>", "Stage specific file", cmdWidth))
	helpLines = append(helpLines, FormatCommand("add -u", "Stage all modified tracked files", cmdWidth))
	helpLines = append(helpLines, FormatCommand("add -p", "Stage interactively (patch mode)", cmdWidth))
	helpLines = append(helpLines, "")

	// Unstaging Commands
	helpLines = append(helpLines, styles.Title.Render("Unstaging Changes:"))
	helpLines = append(helpLines, FormatCommand("reset", "Unstage all changes", cmdWidth))
	helpLines = append(helpLines, FormatCommand("reset <file>", "Unstage specific file", cmdWidth))
	helpLines = append(helpLines, FormatCommand("checkout -- <file>", "Restore file to last commit", cmdWidth))
	helpLines = append(helpLines, "")

	// Commit Commands
	helpLines = append(helpLines, styles.Title.Render("Committing:"))
	helpLines = append(helpLines, FormatCommand("commit", "Open commit with message prompt", cmdWidth))
	helpLines = append(helpLines, FormatCommand("commit -m \"message\"", "Commit with inline message", cmdWidth))
	helpLines = append(helpLines, FormatCommand("commit -a -m \"message\"", "Stage all and commit", cmdWidth))
	helpLines = append(helpLines, FormatCommand("commit --amend", "Amend last commit", cmdWidth))
	helpLines = append(helpLines, "")

	// Branch Commands
	helpLines = append(helpLines, styles.Title.Render("Branches:"))
	helpLines = append(helpLines, FormatCommand("branch", "List branches", cmdWidth))
	helpLines = append(helpLines, FormatCommand("branch <name>", "Create new branch", cmdWidth))
	helpLines = append(helpLines, FormatCommand("checkout <branch>", "Switch to branch", cmdWidth))
	helpLines = append(helpLines, FormatCommand("checkout -b <name>", "Create and switch to new branch", cmdWidth))
	helpLines = append(helpLines, FormatCommand("branch -d <name>", "Delete branch", cmdWidth))
	helpLines = append(helpLines, "")

	// Remote Commands
	helpLines = append(helpLines, styles.Title.Render("Remote Operations:"))
	helpLines = append(helpLines, FormatCommand("push", "Push to remote", cmdWidth))
	helpLines = append(helpLines, FormatCommand("push -u origin <branch>", "Push and set upstream", cmdWidth))
	helpLines = append(helpLines, FormatCommand("pull", "Pull from remote", cmdWidth))
	helpLines = append(helpLines, FormatCommand("fetch", "Fetch from remote", cmdWidth))
	helpLines = append(helpLines, "")

	// Status & History Commands
	helpLines = append(helpLines, styles.Title.Render("Viewing Information:"))
	helpLines = append(helpLines, FormatCommand("status", "Show working tree status", cmdWidth))
	helpLines = append(helpLines, FormatCommand("log", "Show commit history", cmdWidth))
	helpLines = append(helpLines, FormatCommand("log --oneline", "Compact log view", cmdWidth))
	helpLines = append(helpLines, FormatCommand("diff <file>", "View changes in file", cmdWidth))
	helpLines = append(helpLines, FormatCommand("diff --cached", "View staged changes", cmdWidth))
	helpLines = append(helpLines, FormatCommand("show <commit>:<file>", "Show file contents at commit", cmdWidth))
	helpLines = append(helpLines, "")

	// Stash Commands
	helpLines = append(helpLines, styles.Title.Render("Stashing:"))
	helpLines = append(helpLines, FormatCommand("stash", "Stash current changes", cmdWidth))
	helpLines = append(helpLines, FormatCommand("stash save \"message\"", "Stash with message", cmdWidth))
	helpLines = append(helpLines, FormatCommand("stash list", "List stashes", cmdWidth))
	helpLines = append(helpLines, FormatCommand("stash pop", "Apply last stash", cmdWidth))
	helpLines = append(helpLines, FormatCommand("stash apply stash@{n}", "Apply specific stash", cmdWidth))
	helpLines = append(helpLines, "")

	// Merge & Rebase
	helpLines = append(helpLines, styles.Title.Render("Merging & Rebasing:"))
	helpLines = append(helpLines, FormatCommand("merge <branch>", "Merge branch into current", cmdWidth))
	helpLines = append(helpLines, FormatCommand("rebase <branch>", "Rebase onto branch", cmdWidth))
	helpLines = append(helpLines, FormatCommand("rebase --continue", "Continue rebase after conflicts", cmdWidth))
	helpLines = append(helpLines, FormatCommand("rebase --abort", "Abort rebase", cmdWidth))
	helpLines = append(helpLines, "")

	// Undo Commands
	helpLines = append(helpLines, styles.Title.Render("Undoing Changes:"))
	helpLines = append(helpLines, FormatCommand("reset --soft HEAD~1", "Undo last commit (keep changes)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("reset --hard HEAD~1", "Undo last commit (discard changes)", cmdWidth))
	helpLines = append(helpLines, FormatCommand("revert <commit>", "Create inverse commit", cmdWidth))
	helpLines = append(helpLines, "")

	helpLines = append(helpLines, styles.Help.Render("Note: 'git' is automatically prepended - just type the command!"))

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
	s := styles.Title.Render("🦅 Git Command Reference") + "\n"
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
