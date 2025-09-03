package resources

import (
	"context"

	"github.com/web-seven/overlock/pkg/configuration"
	"k8s.io/client-go/dynamic"
)

type ResourceRow struct {
	Name        string
	Package     string
	Version     string
	Status      string
	InstallDate string
	Description string
}

func LoadConfigurations(ctx context.Context, dynamicClient dynamic.Interface) ([]ResourceRow, error) {
	if dynamicClient == nil {
		return []ResourceRow{}, nil
	}

	configs := configuration.GetConfigurations(ctx, dynamicClient)
	configRows := make([]ResourceRow, len(configs))

	for i, config := range configs {
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
