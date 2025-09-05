package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/web-seven/overlock/internal/provider"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
)

func LoadProviders(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger) ([]ResourceRow, error) {
	if dynamicClient == nil {
		if logger != nil {
			logger.Debug("Dynamic client is nil, returning empty providers")
		}
		return []ResourceRow{}, nil
	}

	if logger != nil {
		logger.Debug("Loading providers from cluster")
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading providers: %w", ctx.Err())
	default:
	}

	providers := provider.ListProviders(ctx, dynamicClient, logger)
	if logger != nil {
		logger.Debugf("Found %d providers", len(providers))
	}

	providerRows := make([]ResourceRow, len(providers))

	for i, prov := range providers {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while processing provider %d/%d: %w", i+1, len(providers), ctx.Err())
		default:
		}
		status := ExtractStatusFromConditions(prov.Status.Conditions)
		version := ExtractVersion(prov.Spec.Package)
		installDate := ExtractInstallDate(prov.CreationTimestamp.Time)
		description := ExtractDescription(prov.Annotations, prov.Spec.Package, "provider")

		providerRows[i] = ResourceRow{
			Name:        prov.Name,
			Package:     prov.Spec.Package,
			Version:     version,
			Status:      status,
			InstallDate: installDate,
			Description: description,
		}
	}

	return providerRows, nil
}

func LoadProvidersPaginated(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger, opts PaginationOptions) (*ResourceResult, error) {
	if dynamicClient == nil {
		if logger != nil {
			logger.Debug("Dynamic client is nil, returning empty providers")
		}
		return &ResourceResult{Items: []ResourceRow{}, Total: 0, ProcessedTime: "0s"}, nil
	}

	if logger != nil {
		logger.Debugf("Loading providers with pagination (limit: %d, offset: %d)", opts.Limit, opts.Offset)
	}

	start := time.Now()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading providers: %w", ctx.Err())
	default:
	}

	providers := provider.ListProviders(ctx, dynamicClient, logger)
	total := len(providers)

	if logger != nil {
		logger.Debugf("Found %d total providers", total)
	}

	var paginatedProviders []ResourceRow
	hasMore := false

	if opts.Limit > 0 {
		end := opts.Offset + opts.Limit
		if opts.Offset < total {
			if end > total {
				end = total
			} else {
				hasMore = true
			}
			providers = providers[opts.Offset:end]
		} else {
			providers = providers[:0]
		}
	}

	paginatedProviders = make([]ResourceRow, len(providers))

	for i, prov := range providers {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while processing provider %d/%d: %w", i+1, len(providers), ctx.Err())
		default:
		}

		status := ExtractStatusFromConditions(prov.Status.Conditions)
		version := ExtractVersion(prov.Spec.Package)
		installDate := ExtractInstallDate(prov.CreationTimestamp.Time)
		description := ExtractDescription(prov.Annotations, prov.Spec.Package, "provider")

		paginatedProviders[i] = ResourceRow{
			Name:        prov.Name,
			Package:     prov.Spec.Package,
			Version:     version,
			Status:      status,
			InstallDate: installDate,
			Description: description,
		}
	}

	processedTime := time.Since(start).String()
	if logger != nil {
		logger.Debugf("Processed %d providers in %s", len(paginatedProviders), processedTime)
	}

	return &ResourceResult{
		Items:         paginatedProviders,
		Total:         total,
		HasMore:       hasMore,
		ProcessedTime: processedTime,
	}, nil
}
