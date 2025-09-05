package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/web-seven/overlock/internal/function"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
)

func LoadFunctions(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger) ([]ResourceRow, error) {
	if dynamicClient == nil {
		if logger != nil {
			logger.Debug("Dynamic client is nil, returning empty functions")
		}
		return []ResourceRow{}, nil
	}

	if logger != nil {
		logger.Debug("Loading functions from cluster")
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading functions: %w", ctx.Err())
	default:
	}

	functions := function.GetFunctions(ctx, dynamicClient)
	if logger != nil {
		logger.Debugf("Found %d functions", len(functions))
	}

	functionRows := make([]ResourceRow, len(functions))

	for i, fn := range functions {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while processing function %d/%d: %w", i+1, len(functions), ctx.Err())
		default:
		}
		status := ExtractStatusFromConditions(fn.Status.Conditions)
		version := ExtractVersion(fn.Spec.Package)
		installDate := ExtractInstallDate(fn.CreationTimestamp.Time)
		description := ExtractDescription(fn.Annotations, fn.Spec.Package, "function")

		functionRows[i] = ResourceRow{
			Name:        fn.Name,
			Package:     fn.Spec.Package,
			Version:     version,
			Status:      status,
			InstallDate: installDate,
			Description: description,
		}
	}

	return functionRows, nil
}

func LoadFunctionsPaginated(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger, opts PaginationOptions) (*ResourceResult, error) {
	if dynamicClient == nil {
		if logger != nil {
			logger.Debug("Dynamic client is nil, returning empty functions")
		}
		return &ResourceResult{Items: []ResourceRow{}, Total: 0, ProcessedTime: "0s"}, nil
	}

	if logger != nil {
		logger.Debugf("Loading functions with pagination (limit: %d, offset: %d)", opts.Limit, opts.Offset)
	}

	start := time.Now()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled while loading functions: %w", ctx.Err())
	default:
	}

	functions := function.GetFunctions(ctx, dynamicClient)
	total := len(functions)

	if logger != nil {
		logger.Debugf("Found %d total functions", total)
	}

	var paginatedFunctions []ResourceRow
	hasMore := false

	if opts.Limit > 0 {
		end := opts.Offset + opts.Limit
		if opts.Offset < total {
			if end > total {
				end = total
			} else {
				hasMore = true
			}
			functions = functions[opts.Offset:end]
		} else {
			functions = functions[:0]
		}
	}

	paginatedFunctions = make([]ResourceRow, len(functions))

	for i, fn := range functions {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while processing function %d/%d: %w", i+1, len(functions), ctx.Err())
		default:
		}

		status := ExtractStatusFromConditions(fn.Status.Conditions)
		version := ExtractVersion(fn.Spec.Package)
		installDate := ExtractInstallDate(fn.CreationTimestamp.Time)
		description := ExtractDescription(fn.Annotations, fn.Spec.Package, "function")

		paginatedFunctions[i] = ResourceRow{
			Name:        fn.Name,
			Package:     fn.Spec.Package,
			Version:     version,
			Status:      status,
			InstallDate: installDate,
			Description: description,
		}
	}

	processedTime := time.Since(start).String()
	if logger != nil {
		logger.Debugf("Processed %d functions in %s", len(paginatedFunctions), processedTime)
	}

	return &ResourceResult{
		Items:         paginatedFunctions,
		Total:         total,
		HasMore:       hasMore,
		ProcessedTime: processedTime,
	}, nil
}
