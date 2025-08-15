package loader

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

	overlockerrors "github.com/web-seven/overlock/pkg/errors"
)

// Load LoadPathArchive from TAR archive path
func LoadPathArchive(path string) (v1.Image, error) {
	image, err := tarball.ImageFromPath(path, nil)
	if err != nil {
		return nil, overlockerrors.NewPackageNotFoundErrorWithCause("", "", "", "failed to load package from archive", err)
	}
	return image, nil
}
