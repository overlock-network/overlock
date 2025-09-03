package resources

import (
	"context"

	"github.com/web-seven/overlock/internal/function"
	"k8s.io/client-go/dynamic"
)

func LoadFunctions(ctx context.Context, dynamicClient dynamic.Interface) ([]ResourceRow, error) {
	if dynamicClient == nil {
		return []ResourceRow{}, nil
	}

	functions := function.GetFunctions(ctx, dynamicClient)
	functionRows := make([]ResourceRow, len(functions))

	for i, fn := range functions {
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
