package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/web-seven/overlock/pkg/configuration"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
)

func LoadConfigurations(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger) ([]ResourceRow, error) {
	if dynamicClient == nil {
		if logger != nil {
			logger.Debug("Dynamic client is nil, returning empty configurations")
		}
		return []ResourceRow{}, nil
	}

	if logger != nil {
		logger.Debug("Loading configurations from cluster")
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading configurations: %w", ctx.Err())
	default:
	}

	configs := configuration.GetConfigurations(ctx, dynamicClient)
	if logger != nil {
		logger.Debugf("Found %d configurations", len(configs))
	}

	configRows := make([]ResourceRow, len(configs))

	for i, config := range configs {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while processing configuration %d/%d: %w", i+1, len(configs), ctx.Err())
		default:
		}
		status := ExtractStatusFromConditions(config.Status.Conditions)
		version := ExtractVersion(config.Spec.Package)
		installDate := ExtractInstallDate(config.CreationTimestamp.Time)
		description := ExtractDescription(config.Annotations, config.Spec.Package, "configuration")

		configRows[i] = ResourceRow{
			Name:        config.Name,
			Package:     config.Spec.Package,
			Version:     version,
			Status:      status,
			InstallDate: installDate,
			Description: description,
		}
	}

	return configRows, nil
}

func LoadConfigurationsPaginated(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger, opts PaginationOptions) (*ResourceResult, error) {
	if dynamicClient == nil {
		if logger != nil {
			logger.Debug("Dynamic client is nil, returning empty configurations")
		}
		return &ResourceResult{Items: []ResourceRow{}, Total: 0, ProcessedTime: "0s"}, nil
	}

	if logger != nil {
		logger.Debugf("Loading configurations with pagination (limit: %d, offset: %d)", opts.Limit, opts.Offset)
	}

	start := time.Now()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading configurations: %w", ctx.Err())
	default:
	}

	configs := configuration.GetConfigurations(ctx, dynamicClient)
	total := len(configs)

	if logger != nil {
		logger.Debugf("Found %d total configurations", total)
	}

	var paginatedConfigs []ResourceRow
	hasMore := false

	if opts.Limit > 0 {
		end := opts.Offset + opts.Limit
		if opts.Offset < total {
			if end > total {
				end = total
			} else {
				hasMore = true
			}
			configs = configs[opts.Offset:end]
		} else {
			configs = configs[:0]
		}
	}

	paginatedConfigs = make([]ResourceRow, len(configs))

	for i, config := range configs {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while processing configuration %d/%d: %w", i+1, len(configs), ctx.Err())
		default:
		}

		status := ExtractStatusFromConditions(config.Status.Conditions)
		version := ExtractVersion(config.Spec.Package)
		installDate := ExtractInstallDate(config.CreationTimestamp.Time)
		description := ExtractDescription(config.Annotations, config.Spec.Package, "configuration")

		paginatedConfigs[i] = ResourceRow{
			Name:        config.Name,
			Package:     config.Spec.Package,
			Version:     version,
			Status:      status,
			InstallDate: installDate,
			Description: description,
		}
	}

	processedTime := time.Since(start).String()
	if logger != nil {
		logger.Debugf("Processed %d configurations in %s", len(paginatedConfigs), processedTime)
	}

	return &ResourceResult{
		Items:         paginatedConfigs,
		Total:         total,
		HasMore:       hasMore,
		ProcessedTime: processedTime,
	}, nil
}
