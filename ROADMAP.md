# Overlock v3.0 Roadmap

## Vision
Transform Overlock into a TUI application using the Charm framework that simplifies Crossplane resource management through intelligent forms and visual insights.

## Core Features

### 1. Complete TUI Interface
Convert all existing CLI functionality to interactive TUI interfaces with consistent navigation and keyboard shortcuts.

### 2. Smart Resource Creation
Automatically generate interactive forms from XRD schemas with intelligent dropdowns populated from live cluster data, making resource creation simple form-filling instead of writing YAML.

### 3. Provider Resource Monitoring
Real-time dashboard showing CPU, memory, and storage usage for all providers with cost analysis and optimization recommendations.

### 4. Visual Resource Dependency Graph
Interactive ASCII graph showing relationships between XRs, XRCs, and managed resources with health status and impact analysis. Visualize composition patches and their real-time connection status, showing how managed resources within compositions are linked through patches and connected to other XRs.

## Implementation Phases

### Phase 1: TUI Foundation (Months 1-2)
Set up Charm framework architecture and convert core environment and resource management to TUI.

### Phase 2: Package Management (Months 2-3)
TUI interfaces for configurations, providers, functions, and registries with integrated search and health monitoring.

### Phase 3: Smart Resource Creation (Months 3-4)
XRD-to-form generation engine with live resource dropdown integration and form validation.

### Phase 4: Monitoring & Visualization (Months 4-5)
Real-time provider monitoring dashboard and interactive dependency graph visualization with composition patch mapping and cross-XR connection tracking.

### Phase 5: Polish & Launch (Months 5-6)
Performance optimization, migration tools, and community feedback integration.

## Success Goals
- 90% user preference for TUI over CLI within first week
- 75% faster resource creation for new users
- 60% reduction in configuration errors
- Complete feature parity with existing CLI