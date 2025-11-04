package config

import "github.com/charmbracelet/lipgloss"

const Title = `
			_________                 _______________ ___.__ 
		/   _____/ ____ _____  _____\__    ___/    |   \__|
		\_____  \ /    \\__  \ \____ \|    |  |    |   /  |
		/        \   |  \/ __ \|  |_> >    |  |    |  /|  |
		/_______  /___|  (____  /   __/|____|  |______/ |__|
						\/     \/     \/|__|               

	`

// UI Colors
const (
	ColorPrimary  = "#c77dff" // lilás
	ColorWhite    = "#FFFFFF"
	ColorGreen    = "#00ff00"
	ColorRed      = "#ff0000"
	ColorDarkGray = "#444444"
	ColorBlack    = "#000000"
)

// UI Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true)

	MenuStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWhite)).
			Padding(0, 2)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true).
			Padding(0, 2)

	TextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWhite)).
			Padding(0, 2)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGreen)).
			Padding(0, 2)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorRed)).
			Padding(0, 2)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorWhite)).
			Background(lipgloss.Color(ColorDarkGray)).
			Padding(0, 1).
			Width(30)

	SelectedInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorBlack)).
				Background(lipgloss.Color(ColorPrimary)).
				Padding(0, 1).
				Width(30)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGreen)).
			Padding(0, 2)

	CheckedCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorGreen)).
				Bold(true).
				Padding(0, 2)

	// Search input styles (without purple background)
	SearchInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorWhite)).
				Padding(0, 1)

	SearchInputActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorPrimary)).
				Padding(0, 1)
)

// Default values
const (
	DefaultHost     = "localhost"
	DefaultPort     = "5432"
	DefaultDatabase = "postgres"
	TitleWidth      = 80
)
