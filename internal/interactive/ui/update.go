package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/web-seven/overlock/internal/interactive/ui/theme"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, tea.Batch(
				m.loadConfigurations(),
				m.loadProviders(),
				m.loadFunctions(),
			)
		case "tab":
			m.activeTab = (m.activeTab + 1) % theme.TabCount
			m.updateTable()
			return m, nil
		}

	case configurationsLoadedMsg:
		m.loading = false
		m.err = msg.err
		if msg.err == nil {
			m.configurations = msg.configurations
			if m.activeTab == ConfigurationsTab {
				m.updateTable()
			}
		}
		return m, nil

	case providersLoadedMsg:
		if msg.err == nil {
			m.providers = msg.providers
			if m.activeTab == ProvidersTab {
				m.updateTable()
			}
		}
		return m, nil

	case functionsLoadedMsg:
		if msg.err == nil {
			m.functions = msg.functions
			if m.activeTab == FunctionsTab {
				m.updateTable()
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width

		theme.SetTitleWidth(msg.Width)

		tableWidth := max(theme.MinTerminalWidth, msg.Width-theme.WindowPadding)
		tableHeight := max(theme.MinTerminalHeight, msg.Height-theme.UIOverhead)

		m.table.SetWidth(tableWidth)
		m.table.SetHeight(tableHeight)

		m.updateTable()
		return m, nil
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}
