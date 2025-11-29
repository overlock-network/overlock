package registry

import (
	"context"
	"fmt"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	"github.com/google/go-containerregistry/pkg/name"
	regv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/web-seven/overlock/internal/loader"
	"github.com/web-seven/overlock/pkg/registry"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const tagDelim = ":"

type loadImageCmd struct {
	Registry string `arg:"" help:"Name of the registry to load the image to."`
	Path     string `arg:"" help:"Path to OCI image TAR archive."`
	Name     string `required:"" short:"i" help:"Image name and tag (e.g., my-image:1.0)."`
	Upgrade  bool   `help:"Upgrade patch version if image exists."`
}

func (c *loadImageCmd) Run(ctx context.Context, client *kubernetes.Clientset, config *rest.Config, logger *zap.SugaredLogger) error {
	logger.Debugf("Loading image from: %s", c.Path)

	// Check if local registry exists
	isLocal, err := registry.IsLocalRegistry(ctx, client)
	if !isLocal || err != nil {
		if err != nil {
			logger.Debug(err)
		}
		reg := registry.NewLocal()
		reg.SetDefault(true)
		err := reg.Create(ctx, config, logger)
		if err != nil {
			return fmt.Errorf("failed to create local registry: %w", err)
		}
	}

	// Load image from TAR archive
	image, err := loader.LoadPathArchive(c.Path)
	if err != nil {
		return fmt.Errorf("failed to load image from archive: %w", err)
	}

	imageName := c.Name
	if c.Upgrade {
		logger.Debug("Upgrading image version")
		imageName, err = c.upgradeImageVersion(ctx, config, logger)
		if err != nil {
			return fmt.Errorf("failed to upgrade image version: %w", err)
		}
	}

	// Push to local registry
	logger.Debugf("Pushing image to local registry as: %s", imageName)
	err = registry.PushLocalRegistry(ctx, imageName, image, config, logger)
	if err != nil {
		return fmt.Errorf("failed to push image to registry: %w", err)
	}

	logger.Infof("Image %s loaded to local registry.", imageName)
	return nil
}

// upgradeImageVersion finds existing versions and increments patch version
func (c *loadImageCmd) upgradeImageVersion(ctx context.Context, config *rest.Config, logger *zap.SugaredLogger) (string, error) {
	pRef, err := name.ParseReference(c.Name, name.WithDefaultRegistry(""))
	if err != nil {
		return "", fmt.Errorf("failed to parse image reference: %w", err)
	}

	requestedVersion, err := semver.NewVersion(pRef.Identifier())
	if err != nil {
		return "", fmt.Errorf("failed to parse version: %w", err)
	}

	// Get existing tags from local registry
	existingTags, err := registry.ListLocalRegistryTags(ctx, pRef.Context().Name(), config, logger)
	if err != nil {
		logger.Debugf("Could not list existing tags: %v", err)
		// If we can't list tags, start with patch 0
		newVersion := semver.New(requestedVersion.Major(), requestedVersion.Minor(), 0, "", "")
		return strings.Join([]string{pRef.Context().Name(), newVersion.String()}, tagDelim), nil
	}

	// Find the highest patch version for the requested minor version
	requestedMinorVersion := semver.New(requestedVersion.Major(), requestedVersion.Minor(), 0, "", "")
	highestPatch := uint64(0)
	found := false

	for _, tag := range existingTags {
		existingVersion, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}
		existingMinorVersion := semver.New(existingVersion.Major(), existingVersion.Minor(), 0, "", "")
		if requestedMinorVersion.Equal(existingMinorVersion) {
			found = true
			if existingVersion.Patch() >= highestPatch {
				highestPatch = existingVersion.Patch() + 1
			}
		}
	}

	var newVersion *semver.Version
	if found {
		newVersion = semver.New(requestedVersion.Major(), requestedVersion.Minor(), highestPatch, "", "")
	} else {
		newVersion = semver.New(requestedVersion.Major(), requestedVersion.Minor(), 0, "", "")
	}

	return strings.Join([]string{pRef.Context().Name(), newVersion.String()}, tagDelim), nil
}

// Image wraps regv1.Image for registry operations
type Image struct {
	regv1.Image
}
