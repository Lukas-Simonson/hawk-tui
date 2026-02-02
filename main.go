package main

import (
	"fmt"
	"os"

	"hawk-tui/internal/audio"
	"hawk-tui/internal/git"
	"hawk-tui/internal/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/faiface/beep"
)

func main() {
	// Check if we're in a git repository
	if !git.IsGitRepository() {
		fmt.Println("Error: Not a git repository")
		os.Exit(1)
	}

	// Warm up and initialize speaker
	audio.WarmUpSpeaker(beep.SampleRate(44100))

	// Create and run the program
	p := tea.NewProgram(model.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
