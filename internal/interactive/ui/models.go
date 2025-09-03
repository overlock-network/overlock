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

type Model struct {
	table          table.Model
	activeTab      Tab
	configurations []ResourceRow
	providers      []ResourceRow
	functions      []ResourceRow
	loading        bool
	err            error
	ctx            context.Context
	dynamicClient  dynamic.Interface
	logger         *zap.SugaredLogger
	windowWidth    int
}

type configurationsLoadedMsg struct {
	configurations []ResourceRow
	err            error
}

type providersLoadedMsg struct {
	providers []ResourceRow
	err       error
}

type functionsLoadedMsg struct {
	functions []ResourceRow
	err       error
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
		table:         t,
		activeTab:     ConfigurationsTab,
		loading:       true,
		ctx:           ctx,
		dynamicClient: dynamicClient,
		logger:        logger,
		windowWidth:   theme.MinTerminalWidth,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadConfigurations(),
		m.loadProviders(),
		m.loadFunctions(),
	)
}

func (m Model) loadConfigurations() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		configs, err := resources.LoadConfigurations(m.ctx, m.dynamicClient)
		return configurationsLoadedMsg{
			configurations: configs,
			err:            err,
		}
	})
}

func (m Model) loadProviders() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		providers, err := resources.LoadProviders(m.ctx, m.dynamicClient, m.logger)
		return providersLoadedMsg{
			providers: providers,
			err:       err,
		}
	})
}

func (m Model) loadFunctions() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		functions, err := resources.LoadFunctions(m.ctx, m.dynamicClient)
		return functionsLoadedMsg{
			functions: functions,
			err:       err,
		}
	})
}
