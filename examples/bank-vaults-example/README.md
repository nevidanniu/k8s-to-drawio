# Bank-Vaults Integration

This example demonstrates the Bank-Vaults annotation support in k8s-to-drawio. Bank-Vaults is a tool that helps inject HashiCorp Vault secrets into Kubernetes pods using webhook annotations.

## Supported Annotations

The following Bank-Vaults annotations are detected and visualized as dependencies:

### 1. `vault.security.banzaicloud.io/vault-tls-secret`
References a Kubernetes Secret containing the Vault CA certificate for TLS verification.

**Example:**
```yaml
annotations:
  vault.security.banzaicloud.io/vault-tls-secret: "vault-tls-secret"
```
**Visualization:** Creates a connection from the referenced Secret to the Deployment.

### 2. `vault.security.banzaicloud.io/vault-serviceaccount`
References a ServiceAccount that will be used for Vault authentication.

**Example:**
```yaml
annotations:
  vault.security.banzaicloud.io/vault-serviceaccount: "vault-auth-sa"
```
**Visualization:** Creates a connection from the referenced ServiceAccount to the Deployment.

### 3. `vault.security.banzaicloud.io/vault-env-from-path`
References a comma-delimited list of Vault secret paths to pull in all secrets as environment variables.

**Example:**
```yaml
annotations:
  vault.security.banzaicloud.io/vault-env-from-path: "secret/myapp/config,secret/shared/database"
```
**Visualization:** Creates virtual VaultSecret nodes for each path and connects them to the Deployment.

### 4. `vault.security.banzaicloud.io/token-auth-mount`
References a volume mount for Vault tokens in `{volume:file}` format.

**Example:**
```yaml
annotations:
  vault.security.banzaicloud.io/token-auth-mount: "vault-token:token"
```
**Visualization:** Creates a connection from the referenced volume (Secret/ConfigMap) to the Deployment.

## How It Works

1. The parser scans both resource-level annotations and pod template annotations
2. When Bank-Vaults annotations are found, it extracts the referenced resource names
3. Dependencies are created showing the relationship between Vault-related resources and the workloads that use them
4. These dependencies are visualized as arrows in the Draw.io diagram

## Example Resources

This directory contains example manifests showing:

- **Deployments** with Bank-Vaults annotations for secret injection
- **Secrets** for Vault TLS certificates and tokens
- **ServiceAccounts** for Vault authentication
- **Services** exposing the applications

## Running the Example

```bash
go run main.go convert --input examples/bank-vaults-example --output vault-diagram.drawio --layout hierarchical
```

## Generated Connections

The resulting diagram will show:
- `vault-tls-secret` → `vault-aware-app` (TLS certificate dependency)
- `vault-auth-sa` → `vault-aware-app` (ServiceAccount dependency)
- `vault-token` → `vault-aware-app` (Token volume dependency)
- `vault-ca-cert` → `another-vault-app` (CA certificate dependency)
- `another-vault-sa` → `another-vault-app` (ServiceAccount dependency)
- Virtual VaultSecret nodes for `secret/myapp/config` → `vault-aware-app` (Vault secret path dependency)

## Benefits

This integration helps visualize:
- Which applications use Vault secrets
- Dependencies on Vault infrastructure components  
- Security relationships between workloads and authentication resources
- Complete architecture including secret management flow
- Virtual Vault secret paths that applications consume
- Comma-separated secret path dependencies for comprehensive secret mapping