# Kubernetes to Draw.io Converter

A CLI tool that converts Kubernetes manifests to Draw.io diagrams with Kustomize support.

## Features

- Parse Kubernetes YAML manifests
- Support for Kustomize overlays and bases
- Generate Draw.io diagrams with dependency relationships
- Bank-Vaults annotation support for Vault secret injection visualization
- Multiple layout algorithms (hierarchical, grid, vertical)
- Namespace grouping (can be disabled with --no-namespaces)
- Comprehensive resource support

## Installation

```bash
go mod download
go build -o k8s-to-drawio
```

## Usage

### Basic Conversion
```bash
k8s-to-drawio convert -i ./manifests -o diagram.drawio
```

### With Kustomize Support
```bash
k8s-to-drawio convert -i ./kustomize-app -o diagram.drawio --kustomize
```

### Custom Layout
```bash
k8s-to-drawio convert -i ./manifests -o diagram.drawio --layout grid
```

### Vertical Layout
```bash
k8s-to-drawio convert -i ./manifests -o diagram.drawio --layout vertical
```

### Disable Namespace Grouping
```bash
k8s-to-drawio convert -i ./manifests -o diagram.drawio --no-namespaces
```

### Combined Options
```bash
k8s-to-drawio convert -i ./manifests -o diagram.drawio --layout vertical --no-namespaces
```

### Validate Manifests
```bash
k8s-to-drawio validate -i ./manifests
```

## Supported Resources

- Deployments, StatefulSets, DaemonSets
- Services (ClusterIP, NodePort, LoadBalancer)
- Ingress
- ConfigMaps, Secrets
- PersistentVolumes, PersistentVolumeClaims
- Namespaces
- ServiceAccounts
- Bank-Vaults annotations for secret injection dependencies

## Layout Algorithms

- **hierarchical**: Organizes resources in dependency-based hierarchy with namespace grouping
- **grid**: Simple grid-based arrangement  
- **vertical**: Aligns resources vertically in namespace columns

All layout algorithms support the `--no-namespaces` flag to create flat diagrams without namespace containers.

## Examples

### Simple Application
```bash
k8s-to-drawio convert -i ./examples/simple-app -o diagram.drawio
```

### Complex Microservices (37 resources)
```bash
# Full e-commerce platform with microservices, databases, monitoring
k8s-to-drawio convert -i ./examples/complex-microservices -o architecture.drawio

# Namespace-specific views
k8s-to-drawio convert -i ./examples/complex-microservices -o ecommerce.drawio --namespace ecommerce

# Different layouts
k8s-to-drawio convert -i ./examples/complex-microservices -o grid-layout.drawio --layout grid
k8s-to-drawio convert -i ./examples/complex-microservices -o vertical-layout.drawio --layout vertical
```

See the `examples/` directory for sample Kubernetes manifests and Kustomize configurations.

### Bank-Vaults Integration
```bash
# Visualize applications with Vault secret injection
k8s-to-drawio convert -i ./examples/bank-vaults-example -o vault-diagram.drawio
```

The Bank-Vaults integration automatically detects and visualizes dependencies created by Vault annotations:
- `vault.security.banzaicloud.io/vault-tls-secret` → Creates connections to referenced TLS Secrets
- `vault.security.banzaicloud.io/vault-serviceaccount` → Creates connections to authentication ServiceAccounts  
- `vault.security.banzaicloud.io/token-auth-mount` → Creates connections to token volumes

See `examples/bank-vaults-example/` for detailed examples and documentation.