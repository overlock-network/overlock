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

```bash
curl -sL "https://raw.githubusercontent.com/overlock-network/overlock/refs/heads/main/scripts/install.sh" | sh
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

# Create with specific engine (kind, k3s, k3d)
overlock environment create <name> --engine=k3d

# Create with port mappings
overlock environment create <name> --http-port=8080 --https-port=8443

# Create with pre-installed packages
overlock environment create <name> \
  --providers=xpkg.upbound.io/crossplane-contrib/provider-gcp:v0.22.0 \
  --configurations=xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31 \
  --functions=xpkg.upbound.io/crossplane-contrib/function-patch-and-transform:v0.8.2

# Create with custom engine configuration
overlock environment create <name> --engine=kind --engine-config=./kind-config.yaml

# Create with mount path
overlock environment create <name> --mount-path=/path/to/storage

# Create with admin service account
overlock environment create <name> --create-admin-service-account

# Create from configuration file
overlock environment create <name> --config=./overlock.yaml

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

**Environment create options:**
| Option | Description |
|--------|-------------|
| `--config` | Path to Overlock configuration file (defaults to ./overlock.yaml) |
| `-p, --http-port` | HTTP host port for mapping (default: 80) |
| `-s, --https-port` | HTTPS host port for mapping (default: 443) |
| `-c, --context` | Kubernetes context to use |
| `-e, --engine` | Kubernetes engine: kind, k3s, k3d (default: kind) |
| `--engine-config` | Path to engine configuration file (kind only) |
| `--mount-path` | Path for mount to /storage host directory |
| `--providers` | Comma-separated list of providers to install |
| `--configurations` | Comma-separated list of configurations to install |
| `--functions` | Comma-separated list of functions to install |
| `--create-admin-service-account` | Create admin service account with cluster-admin privileges |
| `--admin-service-account-name` | Name for admin service account (default: overlock-admin) |

### Provider Management

Install and manage cloud providers (GCP, AWS, Azure, etc.):

```bash
# Install provider from repository
overlock provider install <provider-url>

# List installed providers
overlock provider list

# Load provider from local archive
overlock provider load <name> --path=./provider.xpkg

# Load and apply provider
overlock provider load <name> --path=./provider.xpkg --apply

# Load and upgrade existing provider
overlock provider load <name> --path=./provider.xpkg --apply --upgrade

# Serve provider for development (with live reload)
overlock provider serve <path> <main-path>

# Remove provider
overlock provider delete <provider-url>
```

**Provider load options:**
| Option | Description |
|--------|-------------|
| `--path` | Path to provider package archive |
| `--apply` | Apply provider after loading |
| `--upgrade` | Upgrade existing provider |

### Configuration Management

Manage Crossplane configurations that define infrastructure patterns:

```bash
# Apply configuration from URL
overlock configuration apply <url>

# Apply and wait for installation
overlock configuration apply <url> --wait

# Apply with timeout
overlock configuration apply <url> --wait --timeout=5m

# Apply multiple configurations
overlock configuration apply xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31,xpkg.upbound.io/devops-toolkit/dot-sql:v3.0.31

# List configurations
overlock configuration list

# Load from local archive
overlock configuration load <name> --path=./config.xpkg

# Load from STDIN
cat config.xpkg | overlock configuration load <name> --stdin

# Load and apply configuration
overlock configuration load <name> --path=./config.xpkg --apply

# Load and upgrade existing configuration
overlock configuration load <name> --path=./config.xpkg --apply --upgrade

# Serve for development (with live reload)
overlock configuration serve <path>

# Delete configuration
overlock configuration delete xpkg.upbound.io/devops-toolkit/dot-application:v3.0.31
```

**Configuration apply options:**
| Option | Description |
|--------|-------------|
| `-w, --wait` | Wait until configuration is installed |
| `-t, --timeout` | Timeout for waiting (e.g., 5m, 300s) |

**Configuration load options:**
| Option | Description |
|--------|-------------|
| `--path` | Path to configuration package archive |
| `--stdin` | Load configuration package from STDIN |
| `--apply` | Apply configuration after loading |
| `--upgrade` | Upgrade existing configuration |

### Function Management

Manage Crossplane functions for custom composition logic:

```bash
# Apply function from URL
overlock function apply <url>

# Apply and wait for installation
overlock function apply <url> --wait

# Apply with timeout
overlock function apply <url> --wait --timeout=5m

# Apply multiple functions
overlock function apply <url1>,<url2>

# List functions
overlock function list

# Load from local archive
overlock function load <name> --path=./function.xpkg

# Load from STDIN
cat function.xpkg | overlock function load <name> --stdin

# Load and apply function
overlock function load <name> --path=./function.xpkg --apply

# Load and upgrade existing function
overlock function load <name> --path=./function.xpkg --apply --upgrade

# Serve for development (with live reload)
overlock function serve <path>

# Delete function
overlock function delete <url>
```

**Function apply options:**
| Option | Description |
|--------|-------------|
| `-w, --wait` | Wait until function is installed |
| `-t, --timeout` | Timeout for waiting (e.g., 5m, 300s) |

**Function load options:**
| Option | Description |
|--------|-------------|
| `--path` | Path to function package archive |
| `--stdin` | Load function package from STDIN |
| `--apply` | Apply function after loading |
| `--upgrade` | Upgrade existing function |

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

# Create registry in specific context
overlock registry create --local --context=my-cluster

# List registries
overlock registry list

# Load OCI image to registry
overlock registry load-image <registry> <path> --name=my-image:1.0

# Load image with auto version upgrade
overlock registry load-image <registry> <path> --name=my-image:1.0 --upgrade

# Load Helm chart to registry
overlock registry load-image <registry> <path> --name=my-chart:1.0 --helm

# Delete registry
overlock registry delete
```

**Registry create options:**
| Option | Description |
|--------|-------------|
| `--registry-server` | Private registry FQDN |
| `--username` | Registry username |
| `--password` | Registry password |
| `--email` | Registry email |
| `--default` | Set registry as default |
| `--local` | Create local registry |
| `-c, --context` | Kubernetes context for registry |

**Registry load-image options:**
| Option | Description |
|--------|-------------|
| `-i, --name` | Image name and tag (e.g., my-image:1.0) |
| `--upgrade` | Upgrade patch version if image exists |
| `--helm` | Add Helm chart OCI manifest layers |

### Resource Management

Create and manage custom resources:

```bash
# Create custom resource definition
overlock resource create <type>

# List custom resources
overlock resource list

# Apply resources from file
overlock resource apply --file=<file.yaml>
```

**Resource apply options:**
| Option | Description |
|--------|-------------|
| `-f, --file` | YAML file containing Overlock resources to apply |

### Shell Completions

Install shell completions for your preferred shell:

```bash
# Install shell completions
overlock install-completions
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
overlock resource apply --file=./infrastructure.yaml
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
overlock resource apply --file=./test-resources.yaml

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