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

	for i, db := range m.Databases {
		prefix := "  "
		isChecked := false

		if i == 0 {
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
		} else {
			// Individual databases
			if _, ok := m.Choices[i]; ok {
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

	s += "\n[↑ ↓ | K J] Navegar   [Espaço] Selecionar   [Enter] Confirmar   [Esc] Voltar\n"

	return s
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
