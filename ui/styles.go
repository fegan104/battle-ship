package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	waterColor    = lipgloss.Color("#1E3A5F")
	shipColor     = lipgloss.Color("#4A5568")
	hitColor      = lipgloss.Color("#E53E3E")
	missColor     = lipgloss.Color("#A0AEC0")
	cursorColor   = lipgloss.Color("#48BB78")
	titleColor    = lipgloss.Color("#63B3ED")
	subtitleColor = lipgloss.Color("#A0AEC0")
	successColor  = lipgloss.Color("#48BB78")
	errorColor    = lipgloss.Color("#FC8181")

	// Base cell style
	cellStyle = lipgloss.NewStyle().
			Width(3).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center)

	// Cell variants
	waterCell = cellStyle.
			Background(waterColor).
			Foreground(lipgloss.Color("#2D4A6F"))

	shipCell = cellStyle.
			Background(shipColor).
			Foreground(lipgloss.Color("#E2E8F0"))

	hitCell = cellStyle.
		Background(hitColor).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true)

	missCell = cellStyle.
			Background(missColor).
			Foreground(lipgloss.Color("#4A5568"))

	cursorCell = cellStyle.
			Background(cursorColor).
			Foreground(lipgloss.Color("#1A202C")).
			Bold(true)

	// Preview styles for ship placement
	validPreviewCell = cellStyle.
				Background(lipgloss.Color("#38A169")).
				Foreground(lipgloss.Color("#FFFFFF"))

	invalidPreviewCell = cellStyle.
				Background(lipgloss.Color("#E53E3E")).
				Foreground(lipgloss.Color("#FFFFFF"))

	// Board styles
	boardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A5568")).
			Padding(0, 1)

	boardTitleStyle = lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true).
			MarginBottom(1)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#718096")).
			MarginRight(0)

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true).
			MarginBottom(1)

	bigTitleStyle = lipgloss.NewStyle().
			Foreground(titleColor).
			Bold(true).
			MarginBottom(2)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtitleColor).
			MarginBottom(2)

	// Message styles
	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E2E8F0")).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Layout styles
	containerStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#718096")).
			MarginTop(2)

	// Status bar style
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0AEC0")).
			MarginTop(1)

	// Menu selection styles
	selectedMenuStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#48BB78")).
				Bold(true).
				MarginLeft(2)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0AEC0")).
			MarginLeft(2)

	// Input field style
	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#63B3ED")).
			Background(lipgloss.Color("#2D3748")).
			Padding(0, 1)
)
