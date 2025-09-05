package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/web-seven/overlock/internal/interactive/resources"
	"github.com/web-seven/overlock/internal/interactive/ui/theme"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
)

type Tab int

const (
	ConfigurationsTab Tab = iota
	ProvidersTab
	FunctionsTab
)

func (t Tab) String() string {
	switch t {
	case ConfigurationsTab:
		return "Configurations"
	case ProvidersTab:
		return "Providers"
	case FunctionsTab:
		return "Functions"
	default:
		return "Unknown"
	}
}

type ResourceRow = resources.ResourceRow

type TabState struct {
	items       []ResourceRow
	total       int
	loaded      int
	hasMore     bool
	loading     bool
	initialized bool
}

type Model struct {
	table         table.Model
	activeTab     Tab
	tabStates     map[Tab]*TabState
	loading       bool
	err           error
	ctx           context.Context
	dynamicClient dynamic.Interface
	logger        *zap.SugaredLogger
	windowWidth   int
}

type configurationsLoadedMsg struct {
	result *resources.ResourceResult
	err    error
}

type providersLoadedMsg struct {
	result *resources.ResourceResult
	err    error
}

type functionsLoadedMsg struct {
	result *resources.ResourceResult
	err    error
}

type loadMoreMsg struct {
	tab Tab
}

func NewModel(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger) Model {
	t := table.New(
		table.WithFocused(true),
		table.WithHeight(theme.DefaultTableHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(theme.ColorBorder)).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color(theme.ColorHeaderFg))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(theme.ColorSelected)).
		Background(lipgloss.Color(theme.ColorSelectedBg)).
		Bold(false)
	t.SetStyles(s)

	return Model{
		table:     t,
		activeTab: ConfigurationsTab,
		tabStates: map[Tab]*TabState{
			ConfigurationsTab: {items: []ResourceRow{}, loading: false, initialized: false},
			ProvidersTab:      {items: []ResourceRow{}, loading: false, initialized: false},
			FunctionsTab:      {items: []ResourceRow{}, loading: false, initialized: false},
		},
		loading:       false,
		ctx:           ctx,
		dynamicClient: dynamicClient,
		logger:        logger,
		windowWidth:   theme.MinTerminalWidth,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadActiveTab()
}

func (m Model) loadConfigurations() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		offset := len(m.tabStates[ConfigurationsTab].items)
		opts := resources.PaginationOptions{
			Limit:  theme.DefaultPageSize,
			Offset: offset,
		}

		if m.logger != nil {
			m.logger.Debugw("Loading configurations", "hasClient", m.dynamicClient != nil, "offset", offset, "limit", opts.Limit)
		}

		result, err := resources.LoadConfigurationsPaginated(m.ctx, m.dynamicClient, m.logger, opts)

		if m.logger != nil && result != nil {
			m.logger.Debugw("Configurations loaded", "count", len(result.Items), "hasMore", result.HasMore, "total", result.Total)
		}

		return configurationsLoadedMsg{
			result: result,
			err:    err,
		}
	})
}

func (m Model) loadProviders() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		offset := len(m.tabStates[ProvidersTab].items)
		opts := resources.PaginationOptions{
			Limit:  theme.DefaultPageSize,
			Offset: offset,
		}

		if m.logger != nil {
			m.logger.Debugw("Loading providers", "hasClient", m.dynamicClient != nil, "offset", offset, "limit", opts.Limit)
		}

		result, err := resources.LoadProvidersPaginated(m.ctx, m.dynamicClient, m.logger, opts)

		if m.logger != nil && result != nil {
			m.logger.Debugw("Providers loaded", "count", len(result.Items), "hasMore", result.HasMore, "total", result.Total)
		}

		return providersLoadedMsg{
			result: result,
			err:    err,
		}
	})
}

func (m Model) loadFunctions() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		offset := len(m.tabStates[FunctionsTab].items)
		opts := resources.PaginationOptions{
			Limit:  theme.DefaultPageSize,
			Offset: offset,
		}

		if m.logger != nil {
			m.logger.Debugw("Loading functions", "hasClient", m.dynamicClient != nil, "offset", offset, "limit", opts.Limit)
		}

		result, err := resources.LoadFunctionsPaginated(m.ctx, m.dynamicClient, m.logger, opts)

		if m.logger != nil && result != nil {
			m.logger.Debugw("Functions loaded", "count", len(result.Items), "hasMore", result.HasMore, "total", result.Total)
		}

		return functionsLoadedMsg{
			result: result,
			err:    err,
		}
	})
}

func (m Model) loadActiveTab() tea.Cmd {
	state := m.tabStates[m.activeTab]
	if state.loading || state.initialized {
		return nil
	}

	switch m.activeTab {
	case ConfigurationsTab:
		return m.loadConfigurations()
	case ProvidersTab:
		return m.loadProviders()
	case FunctionsTab:
		return m.loadFunctions()
	default:
		return nil
	}
}

func (m Model) loadMore(tab Tab) tea.Cmd {
	state := m.tabStates[tab]
	if state.loading || !state.hasMore {
		return nil
	}

	switch tab {
	case ConfigurationsTab:
		return m.loadConfigurations()
	case ProvidersTab:
		return m.loadProviders()
	case FunctionsTab:
		return m.loadFunctions()
	default:
		return nil
	}
}
