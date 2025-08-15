package errors

import (
	"errors"
	"fmt"
)

// InvalidConfigError represents configuration-related errors
type InvalidConfigError struct {
	Field   string
	Value   string
	Message string
	Err     error
}

func (e *InvalidConfigError) Error() string {
	if e.Field != "" && e.Value != "" {
		return fmt.Sprintf("invalid configuration: field '%s' with value '%s': %s", e.Field, e.Value, e.Message)
	}
	if e.Field != "" {
		return fmt.Sprintf("invalid configuration: field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("invalid configuration: %s", e.Message)
}

func (e *InvalidConfigError) Unwrap() error {
	return e.Err
}

// NewInvalidConfigError creates a new InvalidConfigError
func NewInvalidConfigError(field, value, message string) *InvalidConfigError {
	return &InvalidConfigError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// NewInvalidConfigErrorWithCause creates a new InvalidConfigError with an underlying cause
func NewInvalidConfigErrorWithCause(field, value, message string, err error) *InvalidConfigError {
	return &InvalidConfigError{
		Field:   field,
		Value:   value,
		Message: message,
		Err:     err,
	}
}

// KubernetesConnectionError represents cluster connectivity errors
type KubernetesConnectionError struct {
	Context string
	Host    string
	Message string
	Err     error
}

func (e *KubernetesConnectionError) Error() string {
	if e.Context != "" && e.Host != "" {
		return fmt.Sprintf("kubernetes connection error: context '%s' (host: %s): %s", e.Context, e.Host, e.Message)
	}
	if e.Context != "" {
		return fmt.Sprintf("kubernetes connection error: context '%s': %s", e.Context, e.Message)
	}
	if e.Host != "" {
		return fmt.Sprintf("kubernetes connection error: host '%s': %s", e.Host, e.Message)
	}
	return fmt.Sprintf("kubernetes connection error: %s", e.Message)
}

func (e *KubernetesConnectionError) Unwrap() error {
	return e.Err
}

// NewKubernetesConnectionError creates a new KubernetesConnectionError
func NewKubernetesConnectionError(context, host, message string) *KubernetesConnectionError {
	return &KubernetesConnectionError{
		Context: context,
		Host:    host,
		Message: message,
	}
}

// NewKubernetesConnectionErrorWithCause creates a new KubernetesConnectionError with an underlying cause
func NewKubernetesConnectionErrorWithCause(context, host, message string, err error) *KubernetesConnectionError {
	return &KubernetesConnectionError{
		Context: context,
		Host:    host,
		Message: message,
		Err:     err,
	}
}

// PackageNotFoundError represents errors when packages cannot be found
type PackageNotFoundError struct {
	PackageName string
	Registry    string
	Version     string
	Message     string
	Err         error
}

func (e *PackageNotFoundError) Error() string {
	if e.Registry != "" && e.Version != "" {
		return fmt.Sprintf("package not found: '%s' version '%s' in registry '%s': %s", e.PackageName, e.Version, e.Registry, e.Message)
	}
	if e.Registry != "" {
		return fmt.Sprintf("package not found: '%s' in registry '%s': %s", e.PackageName, e.Registry, e.Message)
	}
	if e.Version != "" {
		return fmt.Sprintf("package not found: '%s' version '%s': %s", e.PackageName, e.Version, e.Message)
	}
	return fmt.Sprintf("package not found: '%s': %s", e.PackageName, e.Message)
}

func (e *PackageNotFoundError) Unwrap() error {
	return e.Err
}

// NewPackageNotFoundError creates a new PackageNotFoundError
func NewPackageNotFoundError(packageName, registry, version, message string) *PackageNotFoundError {
	return &PackageNotFoundError{
		PackageName: packageName,
		Registry:    registry,
		Version:     version,
		Message:     message,
	}
}

// NewPackageNotFoundErrorWithCause creates a new PackageNotFoundError with an underlying cause
func NewPackageNotFoundErrorWithCause(packageName, registry, version, message string, err error) *PackageNotFoundError {
	return &PackageNotFoundError{
		PackageName: packageName,
		Registry:    registry,
		Version:     version,
		Message:     message,
		Err:         err,
	}
}

// Helper functions for error checking
func IsInvalidConfigError(err error) bool {
	var invalidConfigErr *InvalidConfigError
	return errors.As(err, &invalidConfigErr)
}

func IsKubernetesConnectionError(err error) bool {
	var k8sErr *KubernetesConnectionError
	return errors.As(err, &k8sErr)
}

func IsPackageNotFoundError(err error) bool {
	var packageErr *PackageNotFoundError
	return errors.As(err, &packageErr)
}
