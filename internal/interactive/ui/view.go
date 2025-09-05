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

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ColorRed)).
			Bold(true)
		status := errorStyle.Render(fmt.Sprintf("Error: %s", m.err.Error()))
		content.WriteString("\n" + status + "\n\n")
	}

	state := m.tabStates[m.activeTab]
	if len(state.items) > 0 || state.initialized {
		content.WriteString(m.renderConnectedTable())
	} else if state.loading {
		status := theme.StatusStyle.Render("Loading data...")
		content.WriteString("\n" + status + "\n\n")
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
	var tableContent strings.Builder

	paginationInfo := m.renderPaginationInfo()
	if paginationInfo != "" {
		tableContent.WriteString(paginationInfo + "\n")
	}

	tableContent.WriteString(m.table.View())

	connectedTableStyle := theme.BaseStyle.Copy().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(theme.ColorBorder)).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true)

	return connectedTableStyle.Render(tableContent.String())
}

func (m Model) renderPaginationInfo() string {
	state := m.tabStates[m.activeTab]
	if !state.initialized && !state.loading {
		return ""
	}

	var info strings.Builder

	if state.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ColorHighlight)).
			Bold(true)
		info.WriteString(loadingStyle.Render("Loading..."))
		if len(state.items) > 0 {
			info.WriteString(" ")
		}
	}

	if len(state.items) > 0 {
		countStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.ColorWhite))

		if state.total > 0 {
			info.WriteString(countStyle.Render(fmt.Sprintf("Showing %d of %d items", len(state.items), state.total)))
		} else {
			info.WriteString(countStyle.Render(fmt.Sprintf("Showing %d items", len(state.items))))
		}

		if state.hasMore {
			moreStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.ColorHighlight))
			info.WriteString(moreStyle.Render(" • Press 'm' to load more"))
		}
	}

	return info.String()
}

func (m Model) renderHelp() string {
	state := m.tabStates[m.activeTab]
	baseHelp := "↑/↓ select • Tab cycle • r refresh • R refresh all • q quit"

	if state.hasMore && !state.loading {
		baseHelp += " • m load more"
	}

	switch m.activeTab {
	case ConfigurationsTab:
		return baseHelp + " • Crossplane configurations"
	case ProvidersTab:
		return baseHelp + " • Crossplane providers"
	case FunctionsTab:
		return baseHelp + " • Crossplane functions"
	default:
		return baseHelp
	}
}
