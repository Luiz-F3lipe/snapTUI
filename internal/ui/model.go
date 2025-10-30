package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Luiz-F3lipe/snapTUI/internal/backup"
	"github.com/Luiz-F3lipe/snapTUI/internal/config"
	"github.com/Luiz-F3lipe/snapTUI/internal/database"
	"github.com/Luiz-F3lipe/snapTUI/internal/types"
	"github.com/Luiz-F3lipe/snapTUI/internal/ui/views"
)

// App represents the main application
type App struct {
	model         types.Model
	dbService     *database.Service
	backupService *backup.Service
}

// NewApp creates a new application instance
func NewApp() *App {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(config.ColorPrimary))

	model := types.Model{
		Screen:          types.ScreenConnection,
		Cursor:          0,
		Options:         []string{"Fazer Backup", "Restaurar Backup", "Configurar Conexão", "Sair"},
		Databases:       []string{},
		Choices:         make(map[int]string),
		DbHost:          config.DefaultHost,
		DbPort:          config.DefaultPort,
		DbUser:          "",
		DbPassword:      "",
		DbName:          config.DefaultDatabase,
		InputField:      0,
		Inputs:          []string{config.DefaultHost, config.DefaultPort, "", "", config.DefaultDatabase},
		Spinner:         s,
		BackupCompleted: false,
		BackupErrors:    []string{},
		BackupSuccess:   0,
		BackupFilenames: []string{},
		TotalBackups:    0,
		IsProcessing:    false,
	}

	return &App{
		model:         model,
		dbService:     database.NewService(),
		backupService: backup.NewService(),
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.BackupCompleteMsg:
		a.model.BackupCompleted = true
		a.model.BackupSuccess = msg.Success
		a.model.BackupErrors = msg.Errors
		a.model.BackupFilenames = msg.Filenames
		a.model.IsProcessing = false
		return a, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		a.model.Spinner, cmd = a.model.Spinner.Update(msg)
		return a, cmd
	case tea.KeyMsg:
		return a.handleKeyPress(msg)
	}

	return a, nil
}

// View renders the current view
func (a *App) View() string {
	switch a.model.Screen {
	case types.ScreenConnection:
		return views.RenderConnection(a.model)
	case types.ScreenMenu:
		return views.RenderMenu(a.model)
	case types.ScreenBackupList:
		return views.RenderDatabaseList(a.model)
	case types.ScreenBackupProgress:
		return views.RenderBackupProgress(a.model)
	default:
		return "Tela inválida"
	}
}

// handleKeyPress processes keyboard input for different screens
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch a.model.Screen {
	case types.ScreenConnection:
		return a.handleConnectionKeys(msg)
	case types.ScreenMenu:
		return a.handleMenuKeys(msg)
	case types.ScreenBackupProgress:
		return a.handleBackupProgressKeys(msg)
	case types.ScreenBackupList:
		return a.handleBackupListKeys(msg)
	}
	return a, nil
}

// handleConnectionKeys processes keys for the connection screen
func (a *App) handleConnectionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "up", "k":
		if a.model.InputField > 0 {
			a.model.InputField--
		}
	case "down", "j":
		if a.model.InputField < 4 {
			a.model.InputField++
		}
	case " ":
		a.model.Inputs[a.model.InputField] = ""
	case "tab":
		a.model.InputField = (a.model.InputField + 1) % 5
	case "enter":
		// Try to connect and list databases
		databases, err := a.dbService.ListDatabases(
			a.model.Inputs[0], a.model.Inputs[1], a.model.Inputs[2],
			a.model.Inputs[3], a.model.Inputs[4],
		)
		if err != nil {
			// TODO: Show connection error
			return a, nil
		}
		// Add "All Databases" at the beginning
		a.model.Databases = append([]string{"All Databases"}, databases...)
		a.model.Screen = types.ScreenMenu
		a.model.Cursor = 0
	case "backspace":
		if len(a.model.Inputs[a.model.InputField]) > 0 {
			a.model.Inputs[a.model.InputField] = a.model.Inputs[a.model.InputField][:len(a.model.Inputs[a.model.InputField])-1]
		}
	case "esc":
		a.model.Screen = types.ScreenMenu
		a.model.Cursor = 0
	default:
		// Add characters to current field
		if len(msg.String()) == 1 {
			a.model.Inputs[a.model.InputField] += msg.String()
		}
	}
	return a, nil
}

// handleMenuKeys processes keys for the main menu
func (a *App) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "up", "k":
		if a.model.Cursor > 0 {
			a.model.Cursor--
		}
	case "down", "j":
		if a.model.Cursor < len(a.model.Options)-1 {
			a.model.Cursor++
		}
	case "enter":
		switch a.model.Cursor {
		case 0:
			// Go to database list screen
			if len(a.model.Databases) > 0 {
				a.model.Screen = types.ScreenBackupList
				a.model.Cursor = 0
			}
		case 1:
			// Restore Backup - not implemented yet
			return a, tea.Quit
		case 2:
			// Configure Connection
			a.model.Screen = types.ScreenConnection
			a.model.Cursor = 0
			a.model.InputField = 0
		case 3:
			return a, tea.Quit
		}
	}
	return a, nil
}

// handleBackupProgressKeys processes keys for the backup progress screen
func (a *App) handleBackupProgressKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "enter", "esc":
		if a.model.BackupCompleted {
			// Clear selections and return to menu
			a.model.Choices = make(map[int]string)
			a.model.Screen = types.ScreenMenu
			a.model.Cursor = 0
		}
	}
	return a, nil
}

// handleBackupListKeys processes keys for the database selection screen
func (a *App) handleBackupListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "esc":
		// Return to menu
		a.model.Screen = types.ScreenMenu
		a.model.Cursor = 0
	case "up", "k":
		if a.model.Cursor > 0 {
			a.model.Cursor--
		}
	case "down", "j":
		if a.model.Cursor < len(a.model.Databases)-1 {
			a.model.Cursor++
		}
	case " ":
		return a.handleDatabaseSelection()
	case "enter":
		// Perform backup of selected databases
		if len(a.model.Choices) > 0 {
			a.model.Screen = types.ScreenBackupProgress
			a.model.BackupCompleted = false
			a.model.IsProcessing = true
			// Count total databases
			total := 0
			for i := range a.model.Choices {
				if i > 0 {
					total++
				}
			}
			a.model.TotalBackups = total
			return a, tea.Batch(a.model.Spinner.Tick, a.backupService.PerformBackupCmd(a.model))
		}
		return a, nil
	}
	return a, nil
}

// handleDatabaseSelection handles database selection/deselection logic
func (a *App) handleDatabaseSelection() (tea.Model, tea.Cmd) {
	if a.model.Cursor == 0 { // "All Databases"
		// Check if all individual databases are selected
		allSelected := true
		for i := 1; i < len(a.model.Databases); i++ {
			if _, ok := a.model.Choices[i]; !ok {
				allSelected = false
				break
			}
		}

		if allSelected {
			// Deselect all
			a.model.Choices = make(map[int]string)
		} else {
			// Select all
			a.model.Choices = make(map[int]string)
			a.model.Choices[0] = "All Databases"
			for i := 1; i < len(a.model.Databases); i++ {
				a.model.Choices[i] = a.model.Databases[i]
			}
		}
	} else {
		// Individual database
		db := a.model.Databases[a.model.Cursor]
		if _, ok := a.model.Choices[a.model.Cursor]; ok {
			delete(a.model.Choices, a.model.Cursor)
			// Remove "All Databases" if it was selected
			delete(a.model.Choices, 0)
		} else {
			a.model.Choices[a.model.Cursor] = db

			// Check if all individual databases are now selected
			allIndividualSelected := true
			for i := 1; i < len(a.model.Databases); i++ {
				if _, ok := a.model.Choices[i]; !ok {
					allIndividualSelected = false
					break
				}
			}

			if allIndividualSelected {
				a.model.Choices[0] = "All Databases"
			}
		}
	}
	return a, nil
}
