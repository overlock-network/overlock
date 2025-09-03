package resources

import (
	"context"

	"github.com/web-seven/overlock/internal/provider"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
)

func LoadProviders(ctx context.Context, dynamicClient dynamic.Interface, logger *zap.SugaredLogger) ([]ResourceRow, error) {
	if dynamicClient == nil {
		return []ResourceRow{}, nil
	}

	providers := provider.ListProviders(ctx, dynamicClient, logger)
	providerRows := make([]ResourceRow, len(providers))

	for i, prov := range providers {
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
