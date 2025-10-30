package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/lib/pq"
)

type screen int

const (
	screenMenu screen = iota
	screenConnection
	screenBackupList
	screenBackupProgress
)

const title = `
			_________                 _______________ ___.__ 
		/   _____/ ____ _____  _____\__    ___/    |   \__|
		\_____  \ /    \\__  \ \____ \|    |  |    |   /  |
		/        \   |  \/ __ \|  |_> >    |  |    |  /|  |
		/_______  /___|  (____  /   __/|____|  |______/ |__|
						\/     \/     \/|__|               

	`

type backupCompleteMsg struct {
	success   int
	errors    []string
	filenames []string
}

type model struct {
	screen          screen
	cursor          int
	options         []string
	databases       []string
	choices         map[int]string
	dbHost          string
	dbPort          string
	dbUser          string
	dbPassword      string
	dbName          string
	inputField      int // 0=host, 1=port, 2=user, 3=password, 4=dbname
	inputs          []string
	spinner         spinner.Model
	backupCompleted bool
	backupErrors    []string
	backupSuccess   int
	backupFilenames []string
	totalBackups    int
	isProcessing    bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#c77dff"))

	return model{
		screen:          screenConnection,
		cursor:          0,
		options:         []string{"Fazer Backup", "Restaurar Backup", "Configurar Conexão", "Sair"},
		databases:       []string{},
		choices:         make(map[int]string),
		dbHost:          "localhost",
		dbPort:          "5432",
		dbUser:          "",
		dbPassword:      "",
		dbName:          "postgres",
		inputField:      0,
		inputs:          []string{"localhost", "5432", "", "", "postgres"},
		spinner:         s,
		backupCompleted: false,
		backupErrors:    []string{},
		backupSuccess:   0,
		backupFilenames: []string{},
		totalBackups:    0,
		isProcessing:    false,
	}
}

func findPgDump() (string, error) {
	// Tenta encontrar pg_dump no PATH
	pgDumpPath, err := exec.LookPath("pg_dump")
	if err == nil {
		return pgDumpPath, nil
	}

	// Caminhos comuns do PostgreSQL no Linux
	commonPaths := []string{
		"/usr/bin/pg_dump",
		"/usr/local/bin/pg_dump",
		"/usr/pgsql-15/bin/pg_dump",
		"/usr/pgsql-14/bin/pg_dump",
		"/usr/pgsql-13/bin/pg_dump",
		"/usr/pgsql-12/bin/pg_dump",
		"/opt/postgresql/bin/pg_dump",
		"/snap/bin/pg_dump",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("pg_dump não encontrado.\n\nPara instalar no Ubuntu/Debian: sudo apt install postgresql-client\nPara instalar no CentOS/RHEL: sudo yum install postgresql\nOu adicione o caminho do pg_dump ao PATH do sistema")
}

func backupDatabase(host, port, user, password, dbname string) error {
	// Encontrar o pg_dump
	pgDumpPath, err := findPgDump()
	if err != nil {
		return err
	}

	// Obter o diretório do executável
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("erro ao obter caminho do executável: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Criar nome do arquivo com timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.backup", dbname, timestamp)
	backupPath := filepath.Join(exeDir, filename)

	// Comando pg_dump
	cmd := exec.Command(pgDumpPath,
		"--host", host,
		"--port", port,
		"--username", user,
		"--no-password",
		"--format", "custom",
		"--file", backupPath,
		dbname,
	)

	// Definir variável de ambiente para senha
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	// Executar o comando
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao executar pg_dump para %s: %v\nOutput: %s", dbname, err, string(output))
	}

	// Removido o fmt.Printf para não interferir com a TUI
	return nil
}

func performBackupCmd(m model) tea.Cmd {
	return func() tea.Msg {
		// Contar total de bancos
		total := 0
		for i := range m.choices {
			if i > 0 {
				total++
			}
		}

		var errors []string
		var filenames []string
		successCount := 0
		current := 0

		for i, db := range m.choices {
			if i > 0 { // Ignora "All Databases" (índice 0)
				current++

				// Criar nome do arquivo
				timestamp := time.Now().Format("20060102_150405")
				filename := fmt.Sprintf("%s_%s.backup", db, timestamp)

				err := backupDatabase(m.inputs[0], m.inputs[1], m.inputs[2], m.inputs[3], db)
				if err != nil {
					errors = append(errors, fmt.Sprintf("Erro no backup de %s: %v", db, err))
				} else {
					successCount++
					filenames = append(filenames, filename)
				}
			}
		}

		return backupCompleteMsg{
			success:   successCount,
			errors:    errors,
			filenames: filenames,
		}
	}
}

func listDatabases(host, port, user, password, dbname string) ([]string, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case backupCompleteMsg:
		m.backupCompleted = true
		m.backupSuccess = msg.success
		m.backupErrors = msg.errors
		m.backupFilenames = msg.filenames
		m.isProcessing = false
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch m.screen {
		case screenConnection:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.inputField > 0 {
					m.inputField--
				}
			case "down", "j":
				if m.inputField < 4 {
					m.inputField++
				}
			case " ":
				m.inputs[m.inputField] = ""
			case "tab":
				m.inputField = (m.inputField + 1) % 5
			case "enter":
				// Tenta conectar e listar bancos
				databases, err := listDatabases(m.inputs[0], m.inputs[1], m.inputs[2], m.inputs[3], m.inputs[4])
				if err != nil {
					// TODO: Mostrar erro de conexão
					return m, nil
				}
				// Adiciona "All Databases" no início
				m.databases = append([]string{"All Databases"}, databases...)
				m.screen = screenMenu
				m.cursor = 0
			case "backspace":
				if len(m.inputs[m.inputField]) > 0 {
					m.inputs[m.inputField] = m.inputs[m.inputField][:len(m.inputs[m.inputField])-1]
				}
			case "esc":
				m.screen = screenMenu
				m.cursor = 0
			default:
				// Adiciona caracteres ao campo atual
				if len(msg.String()) == 1 {
					m.inputs[m.inputField] += msg.String()
				}
			}

		case screenMenu:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.options)-1 {
					m.cursor++
				}
			case "enter":
				switch m.cursor {
				case 0:
					// Vai para a tela de lista de bancos
					if len(m.databases) > 0 {
						m.screen = screenBackupList
						m.cursor = 0
					}
				case 1:
					// Restaurar Backup - ainda não implementado
					return m, tea.Quit
				case 2:
					// Configurar Conexão
					m.screen = screenConnection
					m.cursor = 0
					m.inputField = 0
				case 3:
					return m, tea.Quit
				}

			}

		case screenBackupProgress:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter", "esc":
				if m.backupCompleted {
					// Limpa as seleções e volta para o menu
					m.choices = make(map[int]string)
					m.screen = screenMenu
					m.cursor = 0
				}
			}

		case screenBackupList:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				// Volta para o menu
				m.screen = screenMenu
				m.cursor = 0
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.databases)-1 {
					m.cursor++
				}
			case " ":
				if m.cursor == 0 { // "All Databases"
					// Verifica se todos os bancos individuais estão selecionados
					allSelected := true
					for i := 1; i < len(m.databases); i++ {
						if _, ok := m.choices[i]; !ok {
							allSelected = false
							break
						}
					}

					if allSelected {
						// Deseleciona todos
						m.choices = make(map[int]string)
					} else {
						// Seleciona todos
						m.choices = make(map[int]string)
						m.choices[0] = "All Databases"
						for i := 1; i < len(m.databases); i++ {
							m.choices[i] = m.databases[i]
						}
					}
				} else {
					// Banco individual
					db := m.databases[m.cursor]
					if _, ok := m.choices[m.cursor]; ok {
						delete(m.choices, m.cursor)
						// Remove "All Databases" se estava selecionado
						delete(m.choices, 0)
					} else {
						m.choices[m.cursor] = db

						// Verifica se todos os bancos individuais estão selecionados
						allIndividualSelected := true
						for i := 1; i < len(m.databases); i++ {
							if _, ok := m.choices[i]; !ok {
								allIndividualSelected = false
								break
							}
						}

						if allIndividualSelected {
							m.choices[0] = "All Databases"
						}
					}
				}
			case "enter":
				// Fazer backup dos bancos selecionados
				if len(m.choices) > 0 {
					m.screen = screenBackupProgress
					m.backupCompleted = false
					m.isProcessing = true
					// Contar total de bancos
					total := 0
					for i := range m.choices {
						if i > 0 {
							total++
						}
					}
					m.totalBackups = total
					return m, tea.Batch(m.spinner.Tick, performBackupCmd(m))
				}
				return m, nil
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case screenConnection:
		return renderConnection(m)
	case screenMenu:
		return renderMenu(m)
	case screenBackupList:
		return renderDatabaseList(m)
	case screenBackupProgress:
		return renderBackupProgress(m)
	default:
		return "Tela inválida"
	}
}

func renderMenu(m model) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c77dff")). // lilás
		Bold(true)

	menuStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c77dff")).
		Bold(true).
		Padding(0, 2)

	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(80, lipgloss.Center, titleStyle.Render(title))

	s := centeredTitle + "\n\n"

	for i, option := range m.options {
		if i == m.cursor {
			s += selectedStyle.Render("-➤ " + option)
		} else {
			s += menuStyle.Render("  " + option)
		}
		s += "\n"
	}

	s += "\n[↑ ↓] Navegar   [Enter] Selecionar   [Q] Sair\n"

	return s
}

func renderDatabaseList(m model) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c77dff")). // lilás
		Bold(true)

	dbStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c77dff")).
		Bold(true).
		Padding(0, 2)

	checkedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")). // verde para selecionados
		Padding(0, 2)

	checkedCursorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")). // verde para selecionados com cursor
		Bold(true).
		Padding(0, 2)

	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(80, lipgloss.Center, titleStyle.Render(title))

	s := centeredTitle + "\n\n"

	for i, db := range m.databases {
		prefix := "  "
		isChecked := false

		if i == 0 {
			// "All Databases" - verifica se todos os bancos individuais estão selecionados
			allSelected := true
			for j := 1; j < len(m.databases); j++ {
				if _, ok := m.choices[j]; !ok {
					allSelected = false
					break
				}
			}
			if allSelected && len(m.databases) > 1 {
				prefix = "[x] "
				isChecked = true
			} else {
				prefix = "[ ] "
			}
		} else {
			// Bancos individuais
			if _, ok := m.choices[i]; ok {
				prefix = "[x] "
				isChecked = true
			} else {
				prefix = "[ ] "
			}
		}

		if i == m.cursor {
			if isChecked {
				s += checkedCursorStyle.Render("-➤ " + prefix + db)
			} else {
				s += selectedStyle.Render("-➤ " + prefix + db)
			}
		} else {
			if isChecked {
				s += checkedStyle.Render("  " + prefix + db)
			} else {
				s += dbStyle.Render("  " + prefix + db)
			}
		}
		s += "\n"
	}

	s += "\n[↑ ↓ | K J] Navegar   [Espaço] Selecionar   [Enter] Confirmar   [Esc] Voltar\n"

	return s
}

func renderConnection(m model) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c77dff")). // lilás
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2)

	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#444444")).
		Padding(0, 1).
		Width(30)

	selectedInputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#c77dff")).
		Padding(0, 1).
		Width(30)

	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(80, lipgloss.Center, titleStyle.Render(title))

	s := centeredTitle + "\n\n"
	s += labelStyle.Render("Configuração de Conexão PostgreSQL") + "\n\n"

	labels := []string{"Host:", "Port:", "User:", "Password:", "Database:"}

	for i, label := range labels {
		s += labelStyle.Render(label) + "\n"
		if i == m.inputField {
			s += selectedInputStyle.Render(m.inputs[i])
		} else {
			// Mascarar senha
			value := m.inputs[i]
			if i == 3 && value != "" { // Password field
				value = string(make([]byte, len(value)))
				for j := range value {
					value = value[:j] + "*" + value[j+1:]
				}
			}
			s += inputStyle.Render(value)
		}
		s += "\n\n"
	}

	s += "\n[↑ ↓ tab] Navegar   [Espaço] Selecionar   [Enter] Conectar   [Q] Sair\n"

	return s
}

func renderBackupProgress(m model) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c77dff")). // lilás
		Bold(true)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 2)

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00ff00")). // verde
		Padding(0, 2)

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0000")). // vermelho
		Padding(0, 2)

	// Título centralizado
	centeredTitle := lipgloss.PlaceHorizontal(80, lipgloss.Center, titleStyle.Render(title))

	s := centeredTitle + "\n\n"

	if !m.backupCompleted {
		s += textStyle.Render("Executando backup dos bancos selecionados...") + "\n\n"

		if m.isProcessing {
			s += textStyle.Render("Aguarde, processando backup...") + "\n\n"
		}

		// Spinner de carregamento
		s += m.spinner.View() + " Processando backup...\n\n"

		s += textStyle.Render(fmt.Sprintf("Bancos a processar: %d", m.totalBackups)) + "\n\n"

	} else {
		s += successStyle.Render("✓ Backup Concluído!") + "\n\n"

		// Resumo dos resultados
		s += textStyle.Render("═══════════════════════════════════════") + "\n"
		s += textStyle.Render("            RESUMO DO BACKUP           ") + "\n"
		s += textStyle.Render("═══════════════════════════════════════") + "\n\n"

		s += successStyle.Render(fmt.Sprintf("✓ Backups realizados com sucesso: %d", m.backupSuccess)) + "\n"

		if len(m.backupFilenames) > 0 {
			s += "\n" + textStyle.Render("Arquivos criados:") + "\n"
			for _, filename := range m.backupFilenames {
				s += textStyle.Render(fmt.Sprintf("  • %s", filename)) + "\n"
			}
		}

		if len(m.backupErrors) > 0 {
			s += "\n" + errorStyle.Render(fmt.Sprintf("✗ Erros encontrados: %d", len(m.backupErrors))) + "\n"
			for _, err := range m.backupErrors {
				s += errorStyle.Render(fmt.Sprintf("  • %s", err)) + "\n"
			}
		}

		s += "\n" + textStyle.Render("═══════════════════════════════════════") + "\n\n"
		s += textStyle.Render("[Enter/Esc] Voltar ao Menu   [Q] Sair") + "\n"
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
}
