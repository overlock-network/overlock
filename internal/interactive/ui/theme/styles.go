package theme

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorHeaderFg)).
			Background(lipgloss.Color(ColorHeaderBg)).
			Width(MinTerminalWidth).
			Align(lipgloss.Center)

	StatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorHeaderFg)).
			Background(lipgloss.Color(ColorHeaderBg)).
			PaddingLeft(1).
			PaddingRight(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorHelpText))

	BaseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(ColorBorder))

	DocStyle = lipgloss.NewStyle().Margin(1, 2)

	ActiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	InactiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	activeTabStyle = lipgloss.NewStyle().
			Border(ActiveTabBorder, true).
			BorderForeground(lipgloss.Color(ColorBorder)).
			Foreground(lipgloss.Color(ColorWhite)).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Border(InactiveTabBorder, true).
				BorderForeground(lipgloss.Color(ColorBorder)).
				Foreground(lipgloss.Color(ColorGray)).
				Padding(0, 1)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorHighlight))

	tabGap = inactiveTabStyle.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		Render(" ")
)

func SetTitleWidth(width int) {
	titleStyle = titleStyle.Width(width)
}
