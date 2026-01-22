[![Discord](https://img.shields.io/badge/discord-join-7289DA.svg?logo=discord&longCache=true&style=flat)](https://discord.gg/W7AsrUb5GC)

<p align="center">
  <img width="500" src="https://raw.githubusercontent.com/overlock-network/overlock/refs/heads/main/docs/overlock_galaxy_text.png"/>
</p>

# Overlock

Overlock is a CLI tool that simplifies Crossplane development and testing. It handles the complexity of setting up Crossplane environments, making it easy for developers to build, test, and deploy infrastructure-as-code solutions.

## Key Features

- **Quick Environment Setup**: Create fully configured Crossplane environments with a single command
- **Multi-Engine Support**: Works with KinD, K3s, and K3d Kubernetes distributions
- **Package Management**: Install and manage Crossplane configurations, providers, and functions
- **Development Workflow**: Live-reload support for local package development
- **Registry Integration**: Support for both local and remote package registries

## Installation

### Prerequisites

- Docker (required for creating Kubernetes clusters)
- One of: KinD, K3s, or K3d (choose based on your preference)

### Install Overlock

Install the latest version:
```bash
curl -sL "https://raw.githubusercontent.com/overlock-network/overlock/refs/heads/main/scripts/install.sh" | sh
sudo mv overlock /usr/local/bin/
```

Install a specific version:
```bash
curl -sL "https://raw.githubusercontent.com/overlock-network/overlock/refs/heads/main/scripts/install.sh" | sh -s -- -v 0.11.0-beta.11
sudo mv overlock /usr/local/bin/
```

Verify installation:
```bash
overlock --version
```

## Quick Start

```bash
# Create a new Crossplane environment
overlock environment create my-dev-env

# Install GCP provider
overlock provider install xpkg.upbound.io/crossplane-contrib/provider-gcp:v0.22.0

# Apply a configuration
overlock configuration apply xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31

# List your environments
overlock environment list
```

## Command Reference

### Environment Management

Create and manage Crossplane-enabled Kubernetes environments:

```bash
# Create new environment
overlock environment create <name>

# List all environments  
overlock environment list

# Start/stop environments
overlock environment start <name>
overlock environment stop <name>

# Upgrade environment to latest Crossplane
overlock environment upgrade <name>

# Delete environment
overlock environment delete <name>
```

### Provider Management

Install and manage cloud providers (GCP, AWS, Azure, etc.):

```bash
# Install provider from repository
overlock provider install <provider-url>

# List installed providers
overlock provider list

# Load provider from local file
overlock provider load <name>

# Serve provider for development (with live reload)
overlock provider serve <path> <main-path>

# Remove provider
overlock provider delete <provider-url>
```

### Configuration Management

Manage Crossplane configurations that define infrastructure patterns:

```bash
# Apply configuration from URL
overlock configuration apply <url>

# Apply multiple configurations
overlock configuration apply xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31,xpkg.upbound.io/devops-toolkit/dot-sql:v3.0.31

# List configurations
overlock configuration list

# Load from local file
overlock configuration load <name>

# Serve for development (with live reload)
overlock configuration serve <path>

# Delete configuration
overlock configuration delete xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31
```

### Function Management

Manage Crossplane functions for custom composition logic:

```bash
# Apply function from URL
overlock function apply <url>

# Apply multiple functions
overlock function apply <url1>,<url2>

# List functions
overlock function list

# Load from local file
overlock function load <name>

# Serve for development (with live reload)
overlock function serve <path>

# Delete function
overlock function delete <url>
```

### Registry Management

Configure package registries for storing and distributing Crossplane packages:

```bash
# Create local registry
overlock registry create --local --default

# Create remote registry connection
overlock registry create --registry-server=<url> \
                        --username=<user> \
                        --password=<pass> \
                        --email=<email>

# List registries
overlock registry list

# Delete registry
overlock registry delete
```

### Resource Management

Create and manage custom resources:

```bash
# Create custom resource definition
overlock resource create <type>

# List custom resources
overlock resource list

# Apply resources from file
overlock resource apply <file.yaml>
```

## Configuration

### Global Options

```bash
overlock [global-options] <command>

Global Options:
  -D, --debug                    Enable debug mode
  -n, --namespace=STRING         Namespace for cluster resources  
  -r, --engine-release=STRING    Crossplane Helm release name
  -v, --engine-version=STRING    Crossplane version (default: 1.19.0)
      --plugin-path=STRING       Path to plugin directory
```

### Environment Variables

- `OVERLOCK_ENGINE_NAMESPACE`: Default namespace for resources
- `OVERLOCK_ENGINE_RELEASE`: Default Helm release name  
- `OVERLOCK_ENGINE_VERSION`: Default Crossplane version

### Command Aliases

All commands support short aliases:
- `environment` → `env`
- `configuration` → `cfg`
- `provider` → `prv`
- `function` → `fnc`
- `registry` → `reg`
- `resource` → `res`

## Usage Examples

### Basic Development Setup

```bash
# Create development environment
overlock environment create crossplane-dev

# Set up local registry
overlock registry create --local --default

# Install commonly used providers
overlock provider install xpkg.upbound.io/crossplane-contrib/provider-gcp:v0.22.0
overlock provider install xpkg.upbound.io/crossplane-contrib/provider-kubernetes:v0.14.0

# Apply base configurations
overlock configuration apply xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31
```

### Working with GCP Infrastructure

```bash
# Create GCP-focused environment
overlock environment create gcp-project

# Install GCP provider and configurations
overlock provider install xpkg.upbound.io/crossplane-contrib/provider-gcp:v0.22.0
overlock configuration apply xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31

# Verify setup
overlock provider list
overlock configuration list

# Apply your infrastructure definitions
overlock resource apply ./infrastructure.yaml
```

### Local Package Development

```bash
# Create development environment
overlock environment create package-dev

# Start live development servers (run in separate terminals)
overlock configuration serve ./my-config-package &
overlock provider serve ./my-provider ./cmd/provider &
overlock function serve ./my-function &

# Test your packages
overlock resource apply ./test-resources.yaml

# Packages automatically reload when you modify code
```

### Multi-Environment Workflow

```bash
# Development environment
overlock environment create dev
overlock environment start dev
# ... do development work ...
overlock environment stop dev

# Testing environment  
overlock environment create test
overlock environment start test
# ... run tests ...
overlock environment stop test

# Staging environment
overlock environment create staging
# ... deploy to staging ...
```

## Troubleshooting

### Common Issues

**Environment creation fails:**
- Ensure Docker is running
- Check that your chosen Kubernetes engine (KinD/K3s/K3d) is installed
- Verify you have sufficient system resources

**Package installation fails:**
- Check internet connectivity for remote packages
- Verify package URLs are correct and accessible
- Use `--debug` flag for detailed error information

**Provider not working:**
- Ensure provider is properly installed: `overlock provider list`
- Check provider configuration and credentials (e.g., GCP service account keys)
- Verify Crossplane version compatibility

### Getting Help

```bash
# General help
overlock --help

# Command-specific help
overlock environment --help
overlock configuration --help

# Enable debug output
overlock --debug <command>
```

## Community

- **Discord**: [Join our Discord](https://discord.gg/W7AsrUb5GC) for questions and community support
- **GitHub**: [Report issues and contribute](https://github.com/overlock-network/overlock)
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.