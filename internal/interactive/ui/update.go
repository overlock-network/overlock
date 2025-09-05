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
			m.tabStates[m.activeTab].loading = true
			return m, m.loadActiveTab()
		case "R":
			for tab := range m.tabStates {
				m.tabStates[tab].items = []ResourceRow{}
				m.tabStates[tab].loaded = 0
				m.tabStates[tab].hasMore = true
				m.tabStates[tab].initialized = false
			}
			m.tabStates[m.activeTab].loading = true
			return m, m.loadActiveTab()
		case "tab":
			m.activeTab = (m.activeTab + 1) % theme.TabCount
			var cmd tea.Cmd
			if !m.tabStates[m.activeTab].initialized {
				cmd = m.loadActiveTab()
			}
			m.updateTable()
			return m, cmd
		case "m", "ctrl+m":
			if m.tabStates[m.activeTab].hasMore && !m.tabStates[m.activeTab].loading {
				m.tabStates[m.activeTab].loading = true
				return m, m.loadMore(m.activeTab)
			}
			return m, nil
		}

	case configurationsLoadedMsg:
		state := m.tabStates[ConfigurationsTab]
		state.loading = false
		m.err = msg.err
		if msg.err == nil && msg.result != nil {
			if len(state.items) == 0 {
				state.items = msg.result.Items
			} else {
				state.items = append(state.items, msg.result.Items...)
			}
			state.total = msg.result.Total
			state.loaded = len(state.items)
			state.hasMore = msg.result.HasMore
			state.initialized = true
			if m.activeTab == ConfigurationsTab {
				m.updateTable()
			}
		}
		return m, nil

	case providersLoadedMsg:
		state := m.tabStates[ProvidersTab]
		state.loading = false
		if msg.err == nil && msg.result != nil {
			if len(state.items) == 0 {
				state.items = msg.result.Items
			} else {
				state.items = append(state.items, msg.result.Items...)
			}
			state.total = msg.result.Total
			state.loaded = len(state.items)
			state.hasMore = msg.result.HasMore
			state.initialized = true
			if m.activeTab == ProvidersTab {
				m.updateTable()
			}
		}
		return m, nil

	case functionsLoadedMsg:
		state := m.tabStates[FunctionsTab]
		state.loading = false
		if msg.err == nil && msg.result != nil {
			if len(state.items) == 0 {
				state.items = msg.result.Items
			} else {
				state.items = append(state.items, msg.result.Items...)
			}
			state.total = msg.result.Total
			state.loaded = len(state.items)
			state.hasMore = msg.result.HasMore
			state.initialized = true
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
