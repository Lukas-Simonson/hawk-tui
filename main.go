package main

import (
	"fmt"
	"os"

	"hawk-tui/internal/git"
	"hawk-tui/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check if we're in a git repository
	if !git.IsGitRepository() {
		fmt.Println("Error: Not a git repository")
		os.Exit(1)
	}

	// Create and run the program
	p := tea.NewProgram(model.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
