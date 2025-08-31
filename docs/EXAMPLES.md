# Examples

This directory contains example Kubernetes manifests and Kustomize configurations to demonstrate the tool's capabilities.

## Simple Application Example

Location: `examples/simple-app/`

A basic web application with:
- Deployment with nginx container
- Service to expose the deployment
- ConfigMap for configuration
- Ingress for external access

### Usage
```bash
k8s-to-drawio convert -i ./examples/simple-app -o simple-app-diagram.drawio
```

### Expected Output
The generated diagram will show:
- webapp Deployment connected to web-config ConfigMap
- web-service Service connected to webapp Deployment
- web-ingress Ingress connected to web-service Service

## Complex Microservices Example

Location: `examples/complex-microservices/`

A comprehensive e-commerce microservices platform demonstrating:
- **37 Kubernetes resources** across multiple namespaces
- **Multiple tiers**: Frontend (API Gateway), Backend (Microservices), Database, Cache
- **Storage**: StatefulSets with PersistentVolumeClaims
- **Security**: RBAC, Secrets, ServiceAccounts
- **Jobs**: Database migration, scheduled backups
- **Monitoring**: Logging agents, metrics collection
- **Networking**: Complex ingress routing, inter-service communication

### Architecture Components

#### Namespaces
- `ecommerce`: Main application namespace
- `monitoring`: Observability stack
- `default`: Kustomization resources

#### Backend Services
- `user-service`: User management and authentication
- `product-service`: Product catalog and inventory
- `order-service`: Order processing and management
- `payment-service`: Payment processing integration
- `api-gateway`: API routing and aggregation

#### Data Layer
- `postgres`: Primary database (StatefulSet)
- `redis`: Cache and session store (StatefulSet)
- PersistentVolumeClaims for data persistence

#### Security & Configuration
- Secrets for database credentials, API keys
- ConfigMaps for service configuration
- RBAC with ServiceAccount, Role, RoleBinding

#### Operations
- Database migration Job
- Scheduled backup CronJob
- Cache cleanup CronJob
- Logging DaemonSet
- Monitoring stack (Prometheus, Grafana)

### Usage Examples

#### Full Architecture
```bash
# Convert all resources (37 total)
k8s-to-drawio convert -i ./examples/complex-microservices -o complex-architecture.drawio

# With Kustomize processing
k8s-to-drawio convert -i ./examples/complex-microservices -o complex-kustomize.drawio --kustomize
```

#### Namespace-Specific Views
```bash
# E-commerce services only (31 resources)
k8s-to-drawio convert -i ./examples/complex-microservices -o ecommerce-services.drawio --namespace ecommerce

# Monitoring stack only (3 resources)
k8s-to-drawio convert -i ./examples/complex-microservices -o monitoring-stack.drawio --namespace monitoring
```

#### Different Layout Algorithms
```bash
# Hierarchical layout (default) - organized by dependencies
k8s-to-drawio convert -i ./examples/complex-microservices -o complex-hierarchical.drawio

# Grid layout - organized grid pattern
k8s-to-drawio convert -i ./examples/complex-microservices -o complex-grid.drawio --layout grid
```

### Expected Output
The generated diagrams will show:
- **Namespace groupings** with proper visual separation
- **Service dependencies** between microservices
- **Data flow** from API Gateway through services to databases
- **Storage relationships** between StatefulSets and PVCs
- **Security context** with RBAC and Secret dependencies
- **Operational components** like jobs and monitoring
- **Network routing** through Ingress to services

## Kustomize Example

Location: `examples/kustomize-example/`

A more complex setup demonstrating Kustomize features:
- Base configuration with API server
- Production overlay with resource limits and replicas
- Generated ConfigMaps

### Structure
```
kustomize-example/
├── base/
│   ├── kustomization.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
└── overlays/
    └── production/
        ├── kustomization.yaml
        └── deployment-patch.yaml
```

### Usage

#### Convert Base Configuration
```bash
k8s-to-drawio convert -i ./examples/kustomize-example/base -o base-diagram.drawio --kustomize
```

#### Convert Production Overlay
```bash
k8s-to-drawio convert -i ./examples/kustomize-example/overlays/production -o prod-diagram.drawio --kustomize
```

### Expected Output
The production diagram will show:
- Scaled deployment (5 replicas instead of 2)
- Additional ConfigMap from generator
- Resource limits applied to containers

## Creating Your Own Examples

### Basic YAML Structure
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: default
spec:
  # deployment spec here
---
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  # service spec here
```

### Kustomize Structure
1. Create a `base/` directory with base resources
2. Add `kustomization.yaml` in base directory
3. Create overlay directories for different environments
4. Add patches and customizations in overlays

### Testing Your Examples
Always test your examples with the validate command first:
```bash
k8s-to-drawio validate -i ./your-example-directory
```

## Advanced Examples

### Multi-Namespace Application
```bash
# Filter by namespace
k8s-to-drawio convert -i ./manifests -o frontend.drawio -n frontend
k8s-to-drawio convert -i ./manifests -o backend.drawio -n backend
```

### Different Layout Algorithms
```bash
# Hierarchical (default)
k8s-to-drawio convert -i ./manifests -o hierarchical.drawio --layout hierarchical

# Grid
k8s-to-drawio convert -i ./manifests -o grid.drawio --layout grid
```

### Complex Dependencies
Create examples with:
- StatefulSets with PersistentVolumeClaims
- Jobs that use ConfigMaps and Secrets
- RBAC resources (ServiceAccounts, Roles, RoleBindings)
- Network policies

These will showcase the tool's ability to detect and visualize complex dependency relationships.