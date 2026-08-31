package registry

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"

	overlockerrors "github.com/web-seven/overlock/pkg/errors"
)

// Catalog returns the repositories stored in the registry.
func (r *Registry) Catalog(ctx context.Context, config *rest.Config, logger *zap.SugaredLogger) ([]string, error) {
	var repos []string
	var err error
	if r.Local {
		repos, err = CatalogLocalRegistry(ctx, config, logger)
	} else {
		repos, err = r.remoteCatalog(ctx)
	}
	if err != nil {
		return nil, r.classifyError("", err)
	}
	return repos, nil
}

// ListTags returns the tags of a repository in the registry.
func (r *Registry) ListTags(ctx context.Context, repository string, config *rest.Config, logger *zap.SugaredLogger) ([]string, error) {
	var tags []string
	var err error
	if r.Local {
		tags, err = ListLocalRegistryTags(ctx, repository, config, logger)
	} else {
		tags, err = r.remoteListTags(ctx, repository)
	}
	if err != nil {
		return nil, r.classifyError(repository, err)
	}
	return tags, nil
}

// remoteAuth returns the credentials stored for the registry, or anonymous access if none are set.
func (r *Registry) remoteAuth() authn.Authenticator {
	for _, auth := range r.Config.Auths {
		if auth.Username != "" || auth.Password != "" {
			return &authn.Basic{Username: auth.Username, Password: auth.Password}
		}
	}
	return authn.Anonymous
}

func (r *Registry) remoteCatalog(ctx context.Context) ([]string, error) {
	host, err := r.Domain()
	if err != nil {
		return nil, err
	}
	reg, err := name.NewRegistry(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse registry host %q: %w", host, err)
	}
	return remote.Catalog(ctx, reg, remote.WithContext(ctx), remote.WithAuth(r.remoteAuth()))
}

func (r *Registry) remoteListTags(ctx context.Context, repository string) ([]string, error) {
	host, err := r.Domain()
	if err != nil {
		return nil, err
	}
	repo, err := name.NewRepository(host + "/" + repository)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repository %q: %w", repository, err)
	}
	return remote.List(repo, remote.WithContext(ctx), remote.WithAuth(r.remoteAuth()))
}

// classifyError wraps a "not found" response (HTTP 404) from a reachable registry in a
// PackageNotFoundError, so callers can tell that case apart from a connection failure.
func (r *Registry) classifyError(repository string, err error) error {
	var terr *transport.Error
	if errors.As(err, &terr) && terr.StatusCode == http.StatusNotFound {
		return overlockerrors.NewPackageNotFoundErrorWithCause(repository, r.GetName(), "", "not found in registry", err)
	}
	return err
}
