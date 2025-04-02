package registry

import (
	"context"

	"github.com/web-seven/overlock/pkg/registry"
	"go.uber.org/zap"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	fieldErrMsg = "Field validation for '%s' failed on the '%s' tag"
)

type createCmd struct {
	RegistryServer string `help:"is your Private Registry FQDN."`
	Username       string `help:"is your Username."`
	Password       string `help:"is your Password."`
	Email          string `help:"is your Email."`
	Default        bool   `help:"Set registry as default."`
	Local          bool   `help:"Create local registry."`
	Context        string `short:"c" help:"Kubernetes context where registry will be created."`
}

func (c *createCmd) Run(ctx context.Context, client *kubernetes.Clientset, config *rest.Config, logger *zap.SugaredLogger) error {
	reg := registry.New(c.RegistryServer, c.Username, c.Password, c.Email)
	if c.Local {
		reg = registry.NewLocal()
	}
	reg.SetDefault(c.Default)
	reg.SetLocal(c.Local)
	reg.WithContext(c.Context)
	err := reg.Validate(ctx, client, logger)
	if err != nil {
		return err
	}

	err = reg.Create(ctx, config, logger)
	if err != nil {
		return err
	} else {
		logger.Info("Registry created successfully.")
	}
	return nil
}
