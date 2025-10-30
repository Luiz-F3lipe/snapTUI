package views

import (
	"fmt"

	"github.com/Luiz-F3lipe/snapTUI/internal/config"
	"github.com/Luiz-F3lipe/snapTUI/internal/types"
	"github.com/charmbracelet/lipgloss"
)

// RenderMenu renders the main menu screen
func RenderMenu(m types.Model) string {
	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(config.TitleWidth, lipgloss.Center, config.TitleStyle.Render(config.Title))

	s := centeredTitle + "\n\n"

	for i, option := range m.Options {
		if i == m.Cursor {
			s += config.SelectedStyle.Render("-➤ " + option)
		} else {
			s += config.MenuStyle.Render("  " + option)
		}
		s += "\n"
	}

	s += "\n[↑ ↓] Navegar   [Enter] Selecionar   [Q] Sair\n"

	return s
}

// RenderConnection renders the database connection configuration screen
func RenderConnection(m types.Model) string {
	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(config.TitleWidth, lipgloss.Center, config.TitleStyle.Render(config.Title))

	s := centeredTitle + "\n\n"
	s += config.TextStyle.Render("Configuração de Conexão PostgreSQL") + "\n\n"

	labels := []string{"Host:", "Port:", "User:", "Password:", "Database:"}

	for i, label := range labels {
		s += config.TextStyle.Render(label) + "\n"
		if i == m.InputField {
			s += config.SelectedInputStyle.Render(m.Inputs[i])
		} else {
			// Mask password
			value := m.Inputs[i]
			if i == 3 && value != "" { // Password field
				value = string(make([]byte, len(value)))
				for j := range value {
					value = value[:j] + "*" + value[j+1:]
				}
			}
			s += config.InputStyle.Render(value)
		}
		s += "\n\n"
	}

	s += "\n[↑ ↓ tab] Navegar   [Espaço] Limpar   [Enter] Conectar   [Esc] Menu   [Q] Sair\n"

	return s
}

// RenderDatabaseList renders the database selection screen
func RenderDatabaseList(m types.Model) string {
	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(config.TitleWidth, lipgloss.Center, config.TitleStyle.Render(config.Title))

	s := centeredTitle + "\n\n"

	// Search box
	searchBox := ""
	if m.SearchMode {
		searchBox = config.TextStyle.Render("🔍 ") + config.SelectedInputStyle.Render(m.SearchInput.View()) + config.TextStyle.Render(" [Esc] sair")
	} else {
		searchValue := m.SearchInput.Value()
		if searchValue != "" {
			searchBox = config.TextStyle.Render("🔍 ") + config.InputStyle.Render(searchValue) + config.TextStyle.Render(" [/] editar")
		} else {
			searchBox = config.TextStyle.Render("🔍 [/] pesquisar bancos")
		}
	}
	s += searchBox + "\n\n"

	// Get current page databases
	currentPageDatabases := getCurrentPageDatabases(m)

	// Show current database list
	for i, db := range currentPageDatabases {
		// Find original index for checking selections
		originalIndex := getOriginalDatabaseIndex(m, db)

		prefix := "  "
		isChecked := false

		if originalIndex == 0 {
			// "All Databases" - check if all individual databases are selected
			allSelected := true
			for j := 1; j < len(m.Databases); j++ {
				if _, ok := m.Choices[j]; !ok {
					allSelected = false
					break
				}
			}
			if allSelected && len(m.Databases) > 1 {
				prefix = "[x] "
				isChecked = true
			} else {
				prefix = "[ ] "
			}
		} else if originalIndex > 0 {
			// Individual databases
			if _, ok := m.Choices[originalIndex]; ok {
				prefix = "[x] "
				isChecked = true
			} else {
				prefix = "[ ] "
			}
		}

		if i == m.Cursor {
			if isChecked {
				s += config.CheckedCursorStyle.Render("-➤ " + prefix + db)
			} else {
				s += config.SelectedStyle.Render("-➤ " + prefix + db)
			}
		} else {
			if isChecked {
				s += config.CheckedStyle.Render("  " + prefix + db)
			} else {
				s += config.MenuStyle.Render("  " + prefix + db)
			}
		}
		s += "\n"
	}

	// Pagination info and controls
	s += "\n"
	if len(m.FilteredDatabases) > m.Paginator.PerPage {
		currentStart := m.Paginator.Page*m.Paginator.PerPage + 1
		currentEnd := min(m.Paginator.Page*m.Paginator.PerPage+m.Paginator.PerPage, len(m.FilteredDatabases))

		s += config.TextStyle.Render(fmt.Sprintf("📄 Página %d de %d  |  Mostrando %d-%d de %d bancos",
			m.Paginator.Page+1, m.Paginator.TotalPages, currentStart, currentEnd, len(m.FilteredDatabases))) + "\n\n"
		s += config.TextStyle.Render("[← → ou H L] Páginas   [↑ ↓] Navegar   [Espaço] Selecionar   [/] Pesquisar   [Enter] Confirmar   [Esc] Voltar") + "\n"
	} else {
		s += config.TextStyle.Render(fmt.Sprintf("Total: %d bancos", len(m.FilteredDatabases))) + "\n"
		s += config.TextStyle.Render("[↑ ↓] Navegar   [Espaço] Selecionar   [/] Pesquisar   [Enter] Confirmar   [Esc] Voltar") + "\n"
	}

	return s
}

// getCurrentPageDatabases returns the databases for the current page (helper for views)
func getCurrentPageDatabases(m types.Model) []string {
	totalItems := len(m.FilteredDatabases)
	if totalItems == 0 {
		return []string{}
	}

	start, end := m.Paginator.GetSliceBounds(totalItems)

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

	return m.FilteredDatabases[start:end]
}

// getOriginalDatabaseIndex finds the original index of a database in the main list
func getOriginalDatabaseIndex(m types.Model, db string) int {
	for i, originalDB := range m.Databases {
		if originalDB == db {
			return i
		}
	}
	return -1
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RenderBackupProgress renders the backup progress and results screen
func RenderBackupProgress(m types.Model) string {
	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(config.TitleWidth, lipgloss.Center, config.TitleStyle.Render(config.Title))

	s := centeredTitle + "\n\n"

	if !m.BackupCompleted {
		s += config.TextStyle.Render("Executando backup dos bancos selecionados...") + "\n\n"

		if m.IsProcessing {
			s += config.TextStyle.Render("Aguarde, processando backup...") + "\n\n"
		}

		// Loading spinner
		s += m.Spinner.View() + " Processando backup...\n\n"

		s += config.TextStyle.Render(fmt.Sprintf("Bancos a processar: %d", m.TotalBackups)) + "\n\n"

	} else {
		s += config.SuccessStyle.Render("✓ Backup Concluído!") + "\n\n"

		// Results summary
		s += config.TextStyle.Render("═══════════════════════════════════════") + "\n"
		s += config.TextStyle.Render("            RESUMO DO BACKUP           ") + "\n"
		s += config.TextStyle.Render("═══════════════════════════════════════") + "\n\n"

		s += config.SuccessStyle.Render(fmt.Sprintf("✓ Backups realizados com sucesso: %d", m.BackupSuccess)) + "\n"

		if len(m.BackupFilenames) > 0 {
			s += "\n" + config.TextStyle.Render("Arquivos criados:") + "\n"
			for _, filename := range m.BackupFilenames {
				s += config.TextStyle.Render(fmt.Sprintf("  • %s", filename)) + "\n"
			}
		}

		if len(m.BackupErrors) > 0 {
			s += "\n" + config.ErrorStyle.Render(fmt.Sprintf("✗ Erros encontrados: %d", len(m.BackupErrors))) + "\n"
			for _, err := range m.BackupErrors {
				s += config.ErrorStyle.Render(fmt.Sprintf("  • %s", err)) + "\n"
			}
		}

		s += "\n" + config.TextStyle.Render("═══════════════════════════════════════") + "\n\n"
		s += config.TextStyle.Render("[Enter/Esc] Voltar ao Menu   [Q] Sair") + "\n"
	}

	return s
}
