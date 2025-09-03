package interactive

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/web-seven/overlock/internal/interactive/ui"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
)

func Run(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger) error {
	m := ui.NewModel(ctx, dynamicClient, logger)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running interactive mode: %w", err)
	}

	return nil
}
