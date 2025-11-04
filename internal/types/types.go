package types

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

// Screen represents the current screen/view
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenConnection
	ScreenBackupList
	ScreenBackupProgress
)

// BackupCompleteMsg represents a completed backup operation
type BackupCompleteMsg struct {
	Success   int
	Errors    []string
	Filenames []string
}

// Model represents the application state
type Model struct {
	// Screen navigation
	Screen Screen
	Cursor int

	// Menu options
	Options []string

	// Database management
	Databases         []string
	FilteredDatabases []string
	Choices           map[int]string

	// Connection details
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	InputField int
	Inputs     []string

	// UI components
	Spinner     spinner.Model
	SearchInput textinput.Model
	Paginator   paginator.Model
	SearchMode  bool

	// Connection status
	ConnectionError string

	// Backup status
	BackupCompleted bool
	BackupErrors    []string
	BackupSuccess   int
	BackupFilenames []string
	TotalBackups    int
	IsProcessing    bool
}

// DatabaseConnection represents database connection parameters
type DatabaseConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
