package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"hawk-tui/internal/types"

	"github.com/charmbracelet/x/ansi"
)

// RefreshStatus fetches the current git status
func RefreshStatus() types.FilesMsg {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return types.FilesMsg{Files: []types.FileState{}}
	}

	lines := strings.Split(string(output), "\n")
	var files []types.FileState

	for _, line := range lines {
		if line == "" {
			continue
		}
		if len(line) < 4 {
			continue
		}
		status := strings.TrimSpace(line[:2])
		path := strings.TrimSpace(line[3:])
		files = append(files, types.FileState{
			Path:   path,
			Status: status,
		})
	}

	return types.FilesMsg{Files: files}
}

// isBinaryFile checks if a file appears to be binary by looking for null bytes
func isBinaryFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 512 bytes to check for binary content
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && n == 0 {
		return false
	}

	// Check for null bytes (common in binary files)
	return bytes.Contains(buf[:n], []byte{0})
}

// GetDiff returns the diff for a given file
func GetDiff(file types.FileState) types.DiffMsg {
	var cmd *exec.Cmd
	shouldFilter := false

	// Handle different file states
	switch file.Status {
	case "??":
		// Untracked file - check if binary first
		if isBinaryFile(file.Path) {
			return types.DiffMsg{Content: "Binary file - diff cannot be displayed"}
		}
		// Show content with line numbers
		cmd = exec.Command("cat", "-n", file.Path)
	case "A", "AM":
		// New file - show colorized diff
		cmd = exec.Command("git", "diff", "--color=always", "--cached", file.Path)
		shouldFilter = true
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			// Check if git reports this as a binary file
			outputStr := string(output)
			if strings.Contains(outputStr, "Binary files") {
				return types.DiffMsg{Content: "Binary file - diff cannot be displayed"}
			}
			filtered := FilterDiffHeaders(outputStr)
			return types.DiffMsg{Content: filtered}
		}
		// If no cached diff, check if binary then show file content with line numbers
		if isBinaryFile(file.Path) {
			return types.DiffMsg{Content: "Binary file - diff cannot be displayed"}
		}
		cmd = exec.Command("cat", "-n", file.Path)
		shouldFilter = false
	default:
		// Modified file - show colorized diff
		cmd = exec.Command("git", "diff", "--color=always", "HEAD", file.Path)
		shouldFilter = true
	}

	output, err := cmd.Output()
	if err != nil {
		return types.DiffMsg{Content: fmt.Sprintf("Error: %v", err)}
	}

	if len(output) == 0 {
		return types.DiffMsg{Content: "No changes to display"}
	}

	result := string(output)

	// Check if git reports this as a binary file
	if strings.Contains(result, "Binary files") {
		return types.DiffMsg{Content: "Binary file - diff cannot be displayed"}
	}

	if shouldFilter {
		result = FilterDiffHeaders(result)
	}

	return types.DiffMsg{Content: result}
}

// ExecuteCommand executes a git command and returns the output
func ExecuteCommand(input string) types.CommandOutputMsg {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return types.CommandOutputMsg{Output: "No command entered"}
	}

	// Automatically prepend "git" to the command
	// User types "add .", we execute "git add ."
	args := parts
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()

	result := string(output)
	if err != nil {
		result += fmt.Sprintf("\nError: %v", err)
	}

	return types.CommandOutputMsg{Output: result}
}

// FilterDiffHeaders removes git diff metadata while preserving color codes
func FilterDiffHeaders(diffOutput string) string {
	lines := strings.Split(diffOutput, "\n")
	var filtered []string

	for _, line := range lines {
		// Strip ANSI codes to check the actual content
		plainLine := ansi.Strip(line)

		// Skip git diff header lines and hunk headers
		if strings.HasPrefix(plainLine, "diff --git") ||
			strings.HasPrefix(plainLine, "index ") ||
			strings.HasPrefix(plainLine, "--- ") ||
			strings.HasPrefix(plainLine, "+++ ") ||
			strings.HasPrefix(plainLine, "@@") {
			continue
		}

		filtered = append(filtered, line)
	}

	return strings.Join(filtered, "\n")
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}
