package main

import (
	"fmt"
	"os"

	"github.com/Luiz-F3lipe/snapTUI/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app := ui.NewApp()
	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao executar aplicação: %v\n", err)
		os.Exit(1)
	}
}
