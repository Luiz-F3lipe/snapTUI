package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(config.ColorPrimary))

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "pesquise aqui..."
	ti.CharLimit = 40
	ti.Width = 40

	// Initialize paginator
	p := paginator.New()
	p.Type = paginator.Arabic
	p.PerPage = 15
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	model := types.Model{
		Screen:            types.ScreenConnection,
		Cursor:            0,
		Options:           []string{"Fazer Backup", "Restaurar Backup", "Configurar Conexão", "Sair"},
		Databases:         []string{},
		FilteredDatabases: []string{},
		Choices:           make(map[int]string),
		DbHost:            config.DefaultHost,
		DbPort:            config.DefaultPort,
		DbUser:            "",
		DbPassword:        "",
		DbName:            config.DefaultDatabase,
		InputField:        0,
		Inputs:            []string{config.DefaultHost, config.DefaultPort, "", "", config.DefaultDatabase},
		Spinner:           s,
		SearchInput:       ti,
		Paginator:         p,
		SearchMode:        false,
		BackupCompleted:   false,
		BackupErrors:      []string{},
		BackupSuccess:     0,
		BackupFilenames:   []string{},
		TotalBackups:      0,
		IsProcessing:      false,
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
		// Handle search input updates when in search mode and backup list screen
		if a.model.Screen == types.ScreenBackupList && a.model.SearchMode {
			oldValue := a.model.SearchInput.Value()
			a.model.SearchInput, _ = a.model.SearchInput.Update(msg)
			// Update filtered databases if search changed
			if a.model.SearchInput.Value() != oldValue {
				a.updateFilteredDatabases()
			}
		}
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
		a.model.FilteredDatabases = a.model.Databases

		// Initialize paginator properly
		a.model.Paginator.SetTotalPages(len(a.model.Databases))
		a.model.Paginator.Page = 0

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
	// If in search mode, handle search-specific keys
	if a.model.SearchMode {
		switch msg.String() {
		case "esc":
			// Exit search mode
			a.model.SearchMode = false
			a.model.SearchInput.Blur()
			return a, nil
		case "enter":
			// Exit search mode and apply filter
			a.model.SearchMode = false
			a.model.SearchInput.Blur()
			return a, nil
		}
		// Let the search input handle other keys
		return a, nil
	}

	// Regular navigation mode
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
	case "esc":
		// Return to menu
		a.model.Screen = types.ScreenMenu
		a.model.Cursor = 0
	case "/":
		// Enter search mode
		a.model.SearchMode = true
		a.model.SearchInput.Focus()
		return a, nil
	case "up", "k":
		if a.model.Cursor > 0 {
			a.model.Cursor--
		}
	case "down", "j":
		currentPageDatabases := a.getCurrentPageDatabases()
		if len(currentPageDatabases) > 0 && a.model.Cursor < len(currentPageDatabases)-1 {
			a.model.Cursor++
		}
	case "left", "h", "pgup":
		// Previous page
		a.model.Paginator.PrevPage()
		a.model.Cursor = 0
	case "right", "l", "pgdown":
		// Next page
		a.model.Paginator.NextPage()
		a.model.Cursor = 0
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
	currentPageDatabases := a.getCurrentPageDatabases()
	if len(currentPageDatabases) == 0 {
		return a, nil
	}

	// Get the actual database from the current page
	selectedDB := currentPageDatabases[a.model.Cursor]

	// Find the index in the original database list
	selectedIndex := -1
	for i, db := range a.model.Databases {
		if db == selectedDB {
			selectedIndex = i
			break
		}
	}

	if selectedIndex == -1 {
		return a, nil
	}

	if selectedIndex == 0 { // "All Databases"
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
		if _, ok := a.model.Choices[selectedIndex]; ok {
			delete(a.model.Choices, selectedIndex)
			// Remove "All Databases" if it was selected
			delete(a.model.Choices, 0)
		} else {
			a.model.Choices[selectedIndex] = selectedDB

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

// updateFilteredDatabases updates the filtered database list based on search query
func (a *App) updateFilteredDatabases() {
	query := strings.ToLower(a.model.SearchInput.Value())
	if query == "" {
		a.model.FilteredDatabases = a.model.Databases
	} else {
		filtered := make([]string, 0)
		for _, db := range a.model.Databases {
			if strings.Contains(strings.ToLower(db), query) {
				filtered = append(filtered, db)
			}
		}
		a.model.FilteredDatabases = filtered
	}
	a.updatePaginator()
}

// updatePaginator updates the paginator based on filtered database count
func (a *App) updatePaginator() {
	totalItems := len(a.model.FilteredDatabases)
	if totalItems == 0 {
		totalItems = 1 // Avoid division by zero
	}

	// SetTotalPages expects the total number of items, not pages
	// It will calculate pages automatically based on PerPage
	a.model.Paginator.SetTotalPages(totalItems)

	// Ensure current page is valid
	if a.model.Paginator.Page >= a.model.Paginator.TotalPages {
		a.model.Paginator.Page = a.model.Paginator.TotalPages - 1
	}
	if a.model.Paginator.Page < 0 {
		a.model.Paginator.Page = 0
	}

	// Reset cursor to valid position for current page
	currentPageItems := a.getCurrentPageDatabases()
	if a.model.Cursor >= len(currentPageItems) {
		a.model.Cursor = 0
	}
}

// getCurrentPageDatabases returns the databases for the current page
func (a *App) getCurrentPageDatabases() []string {
	totalItems := len(a.model.FilteredDatabases)
	if totalItems == 0 {
		return []string{}
	}

	start, end := a.model.Paginator.GetSliceBounds(totalItems)

	// Safety check for bounds
	if start < 0 {
		start = 0
	}
	if end > totalItems {
		end = totalItems
	}
	if start >= end {
		return []string{}
	}

	return a.model.FilteredDatabases[start:end]
}
