package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/web-seven/overlock/internal/interactive/ui/theme"
)

func (m Model) View() string {
	var content strings.Builder

	tabs := m.renderTabs()
	content.WriteString(tabs)

	if m.loading {
		status := theme.StatusStyle.Render("Loading data...")
		content.WriteString("\n" + status + "\n\n")
	} else if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ColorRed)).
			Bold(true)
		status := errorStyle.Render(fmt.Sprintf("Error: %s", m.err.Error()))
		content.WriteString("\n" + status + "\n\n")
	}

	if !m.loading && m.err == nil {
		content.WriteString(m.renderConnectedTable())
	}

	help := m.renderHelp()
	content.WriteString("\n" + theme.HelpStyle.Render(help))

	return theme.DocStyle.Render(content.String())
}

func (m Model) renderTabs() string {
	var renderedTabs []string
	tabs := []string{"Configurations", "Providers", "Functions"}

	for i, t := range tabs {
		var style lipgloss.Style
		if i == int(m.activeTab) {
			b := theme.ActiveTabBorder
			if i == 0 {
				b.BottomLeft = "│"
			}
			style = lipgloss.NewStyle().
				Border(b, true).
				BorderForeground(lipgloss.Color(theme.ColorBorder)).
				Foreground(lipgloss.Color(theme.ColorWhite)).
				Padding(0, 1)
		} else {
			b := theme.InactiveTabBorder
			if i == 0 {
				b.BottomLeft = "│"
			}
			style = lipgloss.NewStyle().
				Border(b, true).
				BorderForeground(lipgloss.Color(theme.ColorBorder)).
				Foreground(lipgloss.Color(theme.ColorGray)).
				Padding(0, 1)
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	gapWidth := m.windowWidth - lipgloss.Width(row) - 2
	gap := ""
	if gapWidth > 0 {
		gap = strings.Repeat("─", gapWidth)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		row,
		lipgloss.NewStyle().Foreground(lipgloss.Color(theme.ColorBorder)).Render(gap),
	)
}

func (m Model) renderConnectedTable() string {
	connectedTableStyle := theme.BaseStyle.Copy().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(theme.ColorBorder)).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true)

	return connectedTableStyle.Render(m.table.View())
}

func (m Model) renderHelp() string {
	baseHelp := "Navigation: ↑/↓ select • Tab cycle • r refresh • q quit"

	switch m.activeTab {
	case ConfigurationsTab:
		return baseHelp + " • Showing: Crossplane configurations"
	case ProvidersTab:
		return baseHelp + " • Showing: Crossplane providers"
	case FunctionsTab:
		return baseHelp + " • Showing: Crossplane functions"
	default:
		return baseHelp
	}
}
