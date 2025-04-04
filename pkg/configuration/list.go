package configuration

import (
	"context"

	crossv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	"github.com/web-seven/overlock/internal/kube"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
)

func GetConfigurations(ctx context.Context, dynamicClient dynamic.Interface) []crossv1.Configuration {
	var params = kube.ResourceParams{
		Dynamic:   dynamicClient,
		Ctx:       ctx,
		Group:     "pkg.crossplane.io",
		Version:   "v1",
		Resource:  "configurations",
		Namespace: "",
	}
	var configurations []crossv1.Configuration
	items, _ := kube.GetKubeResources(params)
	for _, item := range items {
		var configuration crossv1.Configuration
		runtime.DefaultUnstructuredConverter.FromUnstructured(item.UnstructuredContent(), &configuration)
		configurations = append(configurations, configuration)
	}

	return configurations
}
