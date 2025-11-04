package types

// FileState represents a file in the git repository
type FileState struct {
	Path   string
	Status string // M, A, D, ??, etc.
}

// FocusSection represents which section is currently focused
type FocusSection int

const (
	FocusFiles FocusSection = iota
	FocusDiff
	FocusCommand
	FocusHelp        // Help is an overlay
	FocusPopup       // Command output popup is an overlay
	FocusCommitPopup // Commit message popup is an overlay
)

// Messages for the Bubble Tea update cycle

type ErrMsg struct {
	Err error
}

func (e ErrMsg) Error() string { return e.Err.Error() }

type FilesMsg struct {
	Files []FileState
}

type DiffMsg struct {
	Content string
}

type CommandOutputMsg struct {
	Output string
}
