package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Luiz-F3lipe/snapTUI/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

// Service handles backup operations
type Service struct{}

// NewService creates a new backup service
func NewService() *Service {
	return &Service{}
}

// FindPgDump locates pg_dump executable
func (s *Service) FindPgDump() (string, error) {
	// Try to find pg_dump in PATH
	pgDumpPath, err := exec.LookPath("pg_dump")
	if err == nil {
		return pgDumpPath, nil
	}

	// Common PostgreSQL paths on Linux
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

	return "", fmt.Errorf("pg_dump not found.\n\nTo install on Ubuntu/Debian: sudo apt install postgresql-client\nTo install on CentOS/RHEL: sudo yum install postgresql\nOr add pg_dump path to system PATH")
}

// BackupDatabase performs backup of a single database
func (s *Service) BackupDatabase(host, port, user, password, dbname string) (string, error) {
	// Find pg_dump
	pgDumpPath, err := s.FindPgDump()
	if err != nil {
		return "", err
	}

	// Get executable directory
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// Create filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.backup", dbname, timestamp)
	backupPath := filepath.Join(exeDir, filename)

	// pg_dump command
	cmd := exec.Command(pgDumpPath,
		"--host", host,
		"--port", port,
		"--username", user,
		"--no-password",
		"--format", "custom",
		"--file", backupPath,
		dbname,
	)

	// Set password environment variable
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))

	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute pg_dump for %s: %w\nOutput: %s", dbname, err, string(output))
	}

	return filename, nil
}

// PerformBackupCmd creates a command to perform backup operation
func (s *Service) PerformBackupCmd(m types.Model) tea.Cmd {
	return func() tea.Msg {
		// Count total databases
		total := 0
		for i := range m.Choices {
			if i > 0 {
				total++
			}
		}

		var errors []string
		var filenames []string
		successCount := 0

		for i, db := range m.Choices {
			if i > 0 { // Skip "All Databases" (index 0)
				filename, err := s.BackupDatabase(m.Inputs[0], m.Inputs[1], m.Inputs[2], m.Inputs[3], db)
				if err != nil {
					errors = append(errors, fmt.Sprintf("Error backing up %s: %v", db, err))
				} else {
					successCount++
					filenames = append(filenames, filename)
				}
			}
		}

		return types.BackupCompleteMsg{
			Success:   successCount,
			Errors:    errors,
			Filenames: filenames,
		}
	}
}
