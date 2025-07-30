package environment

import (
	"context"
	"errors"
	"os"

	"dario.cat/mergo"
	"github.com/web-seven/overlock/pkg/environment"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type createCmd struct {
	Name   string `arg:"" required:"" help:"Name of environment."`
	Config string `optional:"" help:"Path to the Overlock configuration file. Defaults to ./overlock.yaml if present."`
	createOptions
}

type createOptions struct {
	HttpPort       int      `optional:"" short:"p" help:"Http host port for mapping" default:"80"`
	HttpsPort      int      `optional:"" short:"s" help:"Https host port for mapping" default:"443"`
	Context        string   `optional:"" short:"c" help:"Kubernetes context where Environment will be created."`
	Engine         string   `optional:"" short:"e" help:"Specifies the Kubernetes engine to use for the runtime environment." default:"kind"`
	EngineConfig   string   `optional:"" help:"Path to the configuration file for the engine. Currently supported for kind clusters."`
	MountPath      string   `optional:"" help:"Path for mount to /storage host directory. By default no mounts."`
	Providers      []string `optional:"" help:"List of providers to apply to the environment."`
	Configurations []string `optional:"" help:"List of configurations to apply to the environment."`
	Functions      []string `optional:"" help:"List of functions to apply to the environment."`
}

func (c *createCmd) Run(ctx context.Context, logger *zap.SugaredLogger) error {
	configPath := c.Config
	userProvidedConfig := c.Config != ""

	if !userProvidedConfig {
		configPath = "./overlock.yaml"
	}
	cfg, err := loadConfig(configPath)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			if userProvidedConfig {
				logger.Errorf("Configuration file not found at specified path: %s", configPath)
				return err
			}
		} else {
			logger.Errorf("Failed to load or parse configuration file %s: %v", configPath, err)
			return err
		}
	}

	if cfg != nil {
		if err := mergo.MergeWithOverwrite(&c.createOptions, cfg, mergo.WithOverride); err != nil {
			logger.Errorf("Failed to merge configuration: %v", err)
			return err
		}
	}

	return environment.
		New(c.Engine, c.Name).
		WithHttpPort(c.HttpPort).
		WithHttpsPort(c.HttpsPort).
		WithContext(c.Context).
		WithMountPath(c.MountPath).
		WithEngineConfig(c.EngineConfig).
		WithProviders(c.Providers).
		WithConfigurations(c.Configurations).
		WithFunctions(c.Functions).
		Create(ctx, logger)
}

func loadConfig(path string) (*createOptions, error) {
	var cfg createOptions

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
