package configuration

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/web-seven/overlock/internal/engine"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	crossv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const apiName = "configurations.pkg.crossplane.io"

// RunConfigurationHealthCheck performs a health check on configurations defined by the links string.
func HealthCheck(ctx context.Context, dc dynamic.Interface, links string, wait bool, timeoutChan <-chan time.Time, logger *zap.SugaredLogger) error {

	linkSet := make(map[string]struct{})
	for _, link := range strings.Split(links, ",") {
		linkSet[link] = struct{}{}
	}
	cfgs := GetConfigurations(ctx, dc)

	for {
		select {
		case <-timeoutChan:
			logger.Error("Timeout reached.")
			return nil
		default:
			allHealthy := true
			for _, cfg := range cfgs {
				if _, linkMatched := linkSet[cfg.Spec.Package]; linkMatched {
					if !CheckHealthStatus(cfg.Status.Conditions) {
						allHealthy = false
						break
					}
				}
			}
			if allHealthy {
				logger.Info("Configuration(s) are healthy.")
				return nil
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *Configuration) Apply(ctx context.Context, config *rest.Config, logger *zap.SugaredLogger) error {

	_, err := engine.VerifyApi(ctx, config, apiName)
	if err != nil {
		logger.Debug(err)
		logger.Infoln("Crossplane not installed in current context.")
		logger.Infoln("Configuration not applied.")
		return nil
	}

	scheme := runtime.NewScheme()
	crossv1.AddToScheme(scheme)
	if kube, err := client.New(config, client.Options{Scheme: scheme}); err == nil {
		for _, link := range strings.Split(c.Name, ",") {
			cfg := &crossv1.Configuration{}
			logger.Debugf("Building package %s", link)
			engine.BuildPack(cfg, link, map[string]string{})
			pa := resource.NewAPIPatchingApplicator(kube)

			if err := pa.Apply(ctx, cfg); err != nil {
				return errors.Wrap(err, "Error apply configuration(s).")
			}
		}
	} else {
		return err
	}

	logger.Info("Configuration(s) applied successfully.")
	return nil
}
