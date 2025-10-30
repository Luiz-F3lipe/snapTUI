package types

import (
	"github.com/charmbracelet/bubbles/spinner"
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
	Screen          Screen
	Cursor          int
	Options         []string
	Databases       []string
	Choices         map[int]string
	DbHost          string
	DbPort          string
	DbUser          string
	DbPassword      string
	DbName          string
	InputField      int // 0=host, 1=port, 2=user, 3=password, 4=dbname
	Inputs          []string
	Spinner         spinner.Model
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
