# Usage Guide

## Command Line Interface

### Basic Commands

#### Convert Command
The `convert` command is the main functionality that transforms Kubernetes manifests into Draw.io diagrams.

```bash
k8s-to-drawio convert [flags]
```

**Required Flags:**
- `-i, --input`: Directory containing Kubernetes manifests
- `-o, --output`: Output file path for the Draw.io diagram

**Optional Flags:**
- `-k, --kustomize`: Enable Kustomize processing
- `-n, --namespace`: Filter resources by namespace
- `-l, --layout`: Choose layout algorithm (hierarchical/grid/vertical)

#### Validate Command
The `validate` command checks the syntax and structure of Kubernetes manifests without generating a diagram.

```bash
k8s-to-drawio validate [flags]
```

**Required Flags:**
- `-i, --input`: Directory containing Kubernetes manifests

**Optional Flags:**
- `-k, --kustomize`: Enable Kustomize processing
- `-n, --namespace`: Filter resources by namespace

#### Version Command
Shows the version information of the tool.

```bash
k8s-to-drawio version
```

## Examples

### Basic Usage

Convert a directory of standard Kubernetes YAML files:
```bash
k8s-to-drawio convert -i ./k8s-manifests -o infrastructure.drawio
```

### Kustomize Support

Process Kustomize configuration with base and overlays:
```bash
k8s-to-drawio convert -i ./kustomize/overlays/production -o prod-infrastructure.drawio --kustomize
```

### Namespace Filtering

Generate diagram for resources in a specific namespace:
```bash
k8s-to-drawio convert -i ./manifests -o frontend.drawio -n frontend-namespace
```

### Layout Options

#### Hierarchical Layout (Default)
Organizes resources based on their dependencies in a top-down hierarchy:
```bash
k8s-to-drawio convert -i ./manifests -o diagram.drawio --layout hierarchical
```

#### Grid Layout
Arranges resources in a simple grid pattern:
```bash
k8s-to-drawio convert -i ./manifests -o diagram.drawio --layout grid
```

## Input Requirements

### Supported File Formats
- YAML files (`.yaml`, `.yml`)
- Multi-document YAML files
- JSON files (`.json`)

### Supported Kubernetes Resources
- **Workloads**: Deployment, StatefulSet, DaemonSet, Pod, Job, CronJob
- **Services**: Service, Ingress
- **Configuration**: ConfigMap, Secret
- **Storage**: PersistentVolume, PersistentVolumeClaim
- **RBAC**: ServiceAccount, Role, RoleBinding, ClusterRole, ClusterRoleBinding
- **Cluster**: Namespace

### Directory Structure
The tool can process:
- Flat directory with YAML files
- Kustomize directory structure with base and overlays
- Mixed YAML and JSON files

## Output Format

The tool generates Draw.io (`.drawio`) files that can be opened with:
- [Draw.io](https://app.diagrams.net/) web application
- Draw.io desktop application
- VS Code with Draw.io integration extension

### Generated Elements
- **Resource Shapes**: Different shapes for different Kubernetes resource types
- **Connections**: Arrows showing dependencies between resources
- **Namespace Groups**: Visual grouping of resources by namespace
- **Labels**: Resource names and types

## Troubleshooting

### Common Issues

#### "No kustomization.yaml found"
Ensure your directory contains a valid `kustomization.yaml` or `kustomization.yml` file when using the `--kustomize` flag.

#### "Failed to parse file"
Check that your YAML files are valid and contain proper Kubernetes resource definitions.

#### "Validation failed"
The tool detected issues in your Kubernetes manifests. Review the error message for specific problems.

### Getting Help
```bash
k8s-to-drawio --help
k8s-to-drawio convert --help
k8s-to-drawio validate --help
```