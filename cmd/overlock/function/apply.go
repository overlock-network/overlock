package function

import (
	"context"
	"time"

	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/web-seven/overlock/internal/function"
)

type applyCmd struct {
	Link    string `arg:"" required:"" help:"Link URL (or multiple comma separated) to Crossplane function to be applied to Environment."`
	Wait    bool   `optional:"" short:"w" help:"Wait until function is installed."`
	Timeout string `optional:"" short:"t" help:"Timeout is used to set how much to wait until function is installed (valid time units are ns, us, ms, s, m, h)"`
}

func (c *applyCmd) Run(ctx context.Context, dc *dynamic.DynamicClient, config *rest.Config, logger *zap.SugaredLogger) error {
	if err := function.ApplyFunction(ctx, c.Link, config, logger); err != nil {
		return err
	}
	if !c.Wait {
		return nil
	}

	var timeoutChan <-chan time.Time
	if c.Timeout != "" {
		timeout, err := time.ParseDuration(c.Timeout)
		if err != nil {
			return err
		}
		timeoutChan = time.After(timeout)
	}
	return function.HealthCheck(ctx, dc, c.Link, c.Wait, timeoutChan, logger)
}
