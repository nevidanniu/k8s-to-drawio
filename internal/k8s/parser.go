package k8s

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilYaml "k8s.io/apimachinery/pkg/util/yaml"

	"k8s-to-drawio/pkg/models"
)

type Parser struct {
	namespace string
}

func NewParser(namespace string) *Parser {
	return &Parser{
		namespace: namespace,
	}
}

func (p *Parser) ParseDirectory(dir string) (*models.ResourceCollection, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	yamlFiles, err := filepath.Glob(filepath.Join(dir, "*.yml"))
	if err != nil {
		return nil, err
	}
	files = append(files, yamlFiles...)

	collection := &models.ResourceCollection{
		Resources:    make([]models.K8sResource, 0),
		Dependencies: make(map[string][]string),
	}

	for _, file := range files {
		resources, err := p.parseFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file %s: %w", file, err)
		}
		collection.Resources = append(collection.Resources, resources...)
	}

	p.buildDependencies(collection)
	return collection, nil
}

// ParseFile parses a single YAML file and returns a ResourceCollection
func (p *Parser) ParseFile(filename string) (*models.ResourceCollection, error) {
	collection := &models.ResourceCollection{
		Resources:    make([]models.K8sResource, 0),
		Dependencies: make(map[string][]string),
	}

	resources, err := p.parseFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filename, err)
	}
	collection.Resources = append(collection.Resources, resources...)

	p.buildDependencies(collection)
	return collection, nil
}

func (p *Parser) parseFile(filename string) ([]models.K8sResource, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var resources []models.K8sResource
	decoder := utilYaml.NewYAMLOrJSONDecoder(strings.NewReader(string(content)), 4096)

	for {
		var obj unstructured.Unstructured
		err := decoder.Decode(&obj)
		if err != nil {
			break // End of file or error
		}

		if obj.Object == nil {
			continue
		}

		// Filter by namespace if specified
		if p.namespace != "" && obj.GetNamespace() != p.namespace {
			continue
		}

		resource := models.K8sResource{
			Object:      &obj,
			Kind:        obj.GetKind(),
			Name:        obj.GetName(),
			Namespace:   obj.GetNamespace(),
			Labels:      obj.GetLabels(),
			Annotations: obj.GetAnnotations(),
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (p *Parser) buildDependencies(collection *models.ResourceCollection) {
	for _, resource := range collection.Resources {
		deps := p.findDependencies(resource, collection.Resources)
		if len(deps) > 0 {
			// Use kind+name as key to avoid conflicts between resources with same name
			key := fmt.Sprintf("%s/%s", resource.Kind, resource.Name)
			collection.Dependencies[key] = deps
		}
	}
}

func (p *Parser) findDependencies(resource models.K8sResource, allResources []models.K8sResource) []string {
	var dependencies []string

	// fmt.Printf("Finding dependencies for %s (kind: %s)\n", resource.Name, resource.Kind)

	switch resource.Kind {
	case "Service":
		// Services depend on Deployments/StatefulSets via selectors
		if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
			if selector, found, _ := unstructured.NestedMap(obj.Object, "spec", "selector"); found {
				for _, other := range allResources {
					// fmt.Printf("Checking if %s matches selector %+v\n", other.Name, selector)
					if p.matchesSelector(other, selector) {
						// fmt.Printf("MATCH: %s matches selector for service %s\n", other.Name, resource.Name)
						dependencies = append(dependencies, other.Name)
					}
				}
			}
		}

	case "Ingress":
		// Ingress depends on Services
		if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
			if rules, found, _ := unstructured.NestedSlice(obj.Object, "spec", "rules"); found {
				for _, rule := range rules {
					if ruleMap, ok := rule.(map[string]interface{}); ok {
						if http, found, _ := unstructured.NestedMap(ruleMap, "http"); found {
							if paths, found, _ := unstructured.NestedSlice(http, "paths"); found {
								for _, path := range paths {
									if pathMap, ok := path.(map[string]interface{}); ok {
										if backend, found, _ := unstructured.NestedMap(pathMap, "backend"); found {
											if service, found, _ := unstructured.NestedMap(backend, "service"); found {
												if serviceName, found, _ := unstructured.NestedString(service, "name"); found {
													dependencies = append(dependencies, serviceName)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

	case "Route":
		// Routes depend on Services (OpenShift Routes)
		if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
			if to, found, _ := unstructured.NestedMap(obj.Object, "spec", "to"); found {
				if kind, found, _ := unstructured.NestedString(to, "kind"); found && kind == "Service" {
					if serviceName, found, _ := unstructured.NestedString(to, "name"); found {
						dependencies = append(dependencies, serviceName)
					}
				}
			}
		}

	case "ServiceMonitor":
		// ServiceMonitors depend on Services via selectors (Prometheus Operator)
		if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
			if selector, found, _ := unstructured.NestedMap(obj.Object, "spec", "selector"); found {
				if matchLabels, found, _ := unstructured.NestedMap(selector, "matchLabels"); found {
					for _, other := range allResources {
						if other.Kind == "Service" && p.matchesServiceLabels(other, matchLabels) {
							dependencies = append(dependencies, other.Name)
						}
					}
				}
			}
		}

	case "Deployment", "StatefulSet", "DaemonSet":
		// These depend on ConfigMaps, Secrets, PVCs, and ServiceAccounts
		volumeDeps := p.findVolumeDependencies(resource)
		envDeps := p.findEnvDependencies(resource)
		serviceAccountDeps := p.findServiceAccountDependencies(resource)
		bankVaultsDeps := p.findBankVaultsDependencies(resource)
		// Also check annotations on the pod template for Bank-Vaults
		bankVaultsTemplateDeps := p.findBankVaultsTemplateAnnotations(resource)
		// fmt.Printf("%s volume deps: %+v\n", resource.Name, volumeDeps)
		// fmt.Printf("%s env deps: %+v\n", resource.Name, envDeps)
		// fmt.Printf("%s serviceAccount deps: %+v\n", resource.Name, serviceAccountDeps)
		// fmt.Printf("%s bank-vaults deps: %+v\n", resource.Name, bankVaultsDeps)
		// fmt.Printf("%s bank-vaults template deps: %+v\n", resource.Name, bankVaultsTemplateDeps)
		dependencies = append(dependencies, volumeDeps...)
		dependencies = append(dependencies, envDeps...)
		dependencies = append(dependencies, serviceAccountDeps...)
		dependencies = append(dependencies, bankVaultsDeps...)
		dependencies = append(dependencies, bankVaultsTemplateDeps...)

	case "RoleBinding", "ClusterRoleBinding":
		// RoleBindings depend on ServiceAccounts and Roles
		rbacDeps := p.findRoleBindingDependencies(resource)
		dependencies = append(dependencies, rbacDeps...)

	case "ServiceAccount":
		// ServiceAccounts are used by workloads (reverse dependency)
		workloadDeps := p.findServiceAccountUsers(resource, allResources)
		dependencies = append(dependencies, workloadDeps...)
	}

	// fmt.Printf("Final dependencies for %s: %+v\n", resource.Name, dependencies)
	// fmt.Printf("Final dependencies for %s: %+v\n", resource.Name, dependencies)
	return dependencies
}

func (p *Parser) matchesSelector(resource models.K8sResource, selector map[string]interface{}) bool {
	if resource.Kind != "Deployment" && resource.Kind != "StatefulSet" && resource.Kind != "DaemonSet" {
		return false
	}

	if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
		if labels, found, _ := unstructured.NestedMap(obj.Object, "spec", "template", "metadata", "labels"); found {
			for key, value := range selector {
				if labelValue, exists := labels[key]; !exists || labelValue != value {
					return false
				}
			}
			return true
		}
	}
	return false
}

// matchesServiceLabels checks if a Service's labels match the given selector
func (p *Parser) matchesServiceLabels(resource models.K8sResource, selector map[string]interface{}) bool {
	if resource.Kind != "Service" {
		return false
	}

	if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
		labels := obj.GetLabels()
		if labels != nil {
			for key, value := range selector {
				if labelValue, exists := labels[key]; !exists || labelValue != value {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (p *Parser) findVolumeDependencies(resource models.K8sResource) []string {
	var dependencies []string

	if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
		if volumes, found, _ := unstructured.NestedSlice(obj.Object, "spec", "template", "spec", "volumes"); found {
			for _, vol := range volumes {
				if volMap, ok := vol.(map[string]interface{}); ok {
					if configMap, found, _ := unstructured.NestedMap(volMap, "configMap"); found {
						if name, found, _ := unstructured.NestedString(configMap, "name"); found {
							dependencies = append(dependencies, name)
						}
					}
					if secret, found, _ := unstructured.NestedMap(volMap, "secret"); found {
						if name, found, _ := unstructured.NestedString(secret, "secretName"); found {
							dependencies = append(dependencies, name)
						}
					}
					if pvc, found, _ := unstructured.NestedMap(volMap, "persistentVolumeClaim"); found {
						if name, found, _ := unstructured.NestedString(pvc, "claimName"); found {
							dependencies = append(dependencies, name)
						}
					}
				}
			}
		}
	}

	return dependencies
}

func (p *Parser) findEnvDependencies(resource models.K8sResource) []string {
	var dependencies []string

	// fmt.Printf("Checking env dependencies for %s\n", resource.Name)

	if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
		if containers, found, _ := unstructured.NestedSlice(obj.Object, "spec", "template", "spec", "containers"); found {
			for _, container := range containers {
				if containerMap, ok := container.(map[string]interface{}); ok {
					// Check envFrom for ConfigMap and Secret references
					if envFrom, found, _ := unstructured.NestedSlice(containerMap, "envFrom"); found {
						// fmt.Printf("Found envFrom for %s: %+v\n", resource.Name, envFrom)
						for _, envSource := range envFrom {
							if envSourceMap, ok := envSource.(map[string]interface{}); ok {
								if configMapRef, found, _ := unstructured.NestedMap(envSourceMap, "configMapRef"); found {
									if name, found, _ := unstructured.NestedString(configMapRef, "name"); found {
										// fmt.Printf("Found configMapRef: %s\n", name)
										dependencies = append(dependencies, name)
									}
								}
								if secretRef, found, _ := unstructured.NestedMap(envSourceMap, "secretRef"); found {
									if name, found, _ := unstructured.NestedString(secretRef, "name"); found {
										dependencies = append(dependencies, name)
									}
								}
							}
						}
					}

					// Check individual env entries for ConfigMap and Secret references
					if env, found, _ := unstructured.NestedSlice(containerMap, "env"); found {
						for _, envVar := range env {
							if envMap, ok := envVar.(map[string]interface{}); ok {
								if valueFrom, found, _ := unstructured.NestedMap(envMap, "valueFrom"); found {
									if configMapRef, found, _ := unstructured.NestedMap(valueFrom, "configMapKeyRef"); found {
										if name, found, _ := unstructured.NestedString(configMapRef, "name"); found {
											dependencies = append(dependencies, name)
										}
									}
									if secretRef, found, _ := unstructured.NestedMap(valueFrom, "secretKeyRef"); found {
										if name, found, _ := unstructured.NestedString(secretRef, "name"); found {
											dependencies = append(dependencies, name)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// fmt.Printf("Final env dependencies for %s: %+v\n", resource.Name, dependencies)
	return dependencies
}

// findServiceAccountDependencies finds ServiceAccount dependencies in workload specifications
func (p *Parser) findServiceAccountDependencies(resource models.K8sResource) []string {
	var dependencies []string

	if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
		// Check serviceAccountName in pod template spec
		if serviceAccountName, found, _ := unstructured.NestedString(obj.Object, "spec", "template", "spec", "serviceAccountName"); found && serviceAccountName != "" {
			dependencies = append(dependencies, serviceAccountName)
		}
	}

	return dependencies
}

// findBankVaultsDependencies finds dependencies based on Bank-Vaults annotations
func (p *Parser) findBankVaultsDependencies(resource models.K8sResource) []string {
	var dependencies []string

	// Check if resource has Bank-Vaults annotations
	if resource.Annotations == nil {
		return dependencies
	}

	// fmt.Printf("Checking Bank-Vaults annotations for %s: %+v\n", resource.Name, resource.Annotations)

	// Look for Bank-Vaults specific annotations that reference Kubernetes resources
	for key, value := range resource.Annotations {
		switch {
		// vault.security.banzaicloud.io/vault-tls-secret references a Kubernetes Secret
		case key == "vault.security.banzaicloud.io/vault-tls-secret" && value != "":
			// fmt.Printf("Found vault-tls-secret reference: %s\n", value)
			dependencies = append(dependencies, value)

		// vault.security.banzaicloud.io/vault-serviceaccount references a ServiceAccount
		case key == "vault.security.banzaicloud.io/vault-serviceaccount" && value != "":
			// fmt.Printf("Found vault-serviceaccount reference: %s\n", value)
			dependencies = append(dependencies, value)

		// vault.security.banzaicloud.io/token-auth-mount can reference volumes/secrets
		// Format: {volume:file} where volume might be a Secret or ConfigMap
		case key == "vault.security.banzaicloud.io/token-auth-mount" && value != "":
			// fmt.Printf("Found token-auth-mount reference: %s\n", value)
			// Parse the volume:file format
			if strings.Contains(value, ":") {
				parts := strings.Split(value, ":")
				if len(parts) >= 1 && parts[0] != "" {
					// The volume name is the first part
					// fmt.Printf("Extracted volume name: %s\n", parts[0])
					dependencies = append(dependencies, parts[0])
				}
			}
		}
	}

	// fmt.Printf("Final Bank-Vaults dependencies for %s: %+v\n", resource.Name, dependencies)
	return dependencies
}

// findBankVaultsTemplateAnnotations finds Bank-Vaults dependencies from pod template annotations
func (p *Parser) findBankVaultsTemplateAnnotations(resource models.K8sResource) []string {
	var dependencies []string

	if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
		// Check annotations on the pod template
		if annotations, found, _ := unstructured.NestedStringMap(obj.Object, "spec", "template", "metadata", "annotations"); found {
			// fmt.Printf("Found pod template annotations for %s: %+v\n", resource.Name, annotations)

			// Look for Bank-Vaults specific annotations that reference Kubernetes resources
			for key, value := range annotations {
				// fmt.Printf("Processing annotation: %s = %s\n", key, value)
				switch {
				// vault.security.banzaicloud.io/vault-tls-secret references a Kubernetes Secret
				case key == "vault.security.banzaicloud.io/vault-tls-secret" && value != "":
					// fmt.Printf("Found vault-tls-secret template reference: %s\n", value)
					dependencies = append(dependencies, value)

				// vault.security.banzaicloud.io/vault-serviceaccount references a ServiceAccount
				case key == "vault.security.banzaicloud.io/vault-serviceaccount" && value != "":
					// fmt.Printf("Found vault-serviceaccount template reference: %s\n", value)
					dependencies = append(dependencies, value)

				// vault.security.banzaicloud.io/vault-env-from-path references Vault secret paths
				// This creates a virtual dependency to represent the Vault secret access
				case key == "vault.security.banzaicloud.io/vault-env-from-path" && value != "":
					// Parse comma-delimited list of vault paths
					paths := strings.Split(value, ",")
					for _, path := range paths {
						path = strings.TrimSpace(path)
						if path != "" {
							// Create a virtual Vault secret node name from the path
							// Convert "secret/myapp/config" to "vault-secret-myapp-config"
							// Remove the "secret/" prefix if it exists, then replace remaining slashes with dashes
							cleanPath := strings.TrimPrefix(path, "secret/")
							virtualSecretName := "vault-secret-" + strings.ReplaceAll(strings.ReplaceAll(cleanPath, "/", "-"), ":", "-")
							// fmt.Printf("Found vault-env-from-path: %s -> %s\n", path, virtualSecretName)
							dependencies = append(dependencies, virtualSecretName)
						}
					}

				// vault.security.banzaicloud.io/token-auth-mount can reference volumes/secrets
				// Format: {volume:file} where volume might be a Secret or ConfigMap
				case key == "vault.security.banzaicloud.io/token-auth-mount" && value != "":
					// fmt.Printf("Found token-auth-mount template reference: %s\n", value)
					// Parse the volume:file format
					if strings.Contains(value, ":") {
						parts := strings.Split(value, ":")
						if len(parts) >= 1 && parts[0] != "" {
							// The volume name is the first part
							// fmt.Printf("Extracted template volume name: %s\n", parts[0])
							dependencies = append(dependencies, parts[0])
						}
					}
				}
			}
		}
	}

	// fmt.Printf("Final Bank-Vaults template dependencies for %s: %+v\n", resource.Name, dependencies)
	return dependencies
}

// findRoleBindingDependencies finds dependencies for RoleBinding and ClusterRoleBinding resources
func (p *Parser) findRoleBindingDependencies(resource models.K8sResource) []string {
	var dependencies []string

	if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
		// Check subjects (ServiceAccounts, Users, Groups)
		if subjects, found, _ := unstructured.NestedSlice(obj.Object, "subjects"); found {
			for _, subject := range subjects {
				if subjectMap, ok := subject.(map[string]interface{}); ok {
					if kind, found, _ := unstructured.NestedString(subjectMap, "kind"); found && kind == "ServiceAccount" {
						if name, found, _ := unstructured.NestedString(subjectMap, "name"); found {
							dependencies = append(dependencies, name)
						}
					}
				}
			}
		}

		// Check roleRef (Role or ClusterRole)
		if roleRef, found, _ := unstructured.NestedMap(obj.Object, "roleRef"); found {
			if name, found, _ := unstructured.NestedString(roleRef, "name"); found {
				dependencies = append(dependencies, name)
			}
		}
	}

	return dependencies
}

// findServiceAccountUsers finds workloads that use a specific ServiceAccount
func (p *Parser) findServiceAccountUsers(serviceAccount models.K8sResource, allResources []models.K8sResource) []string {
	var dependencies []string

	// Look for workloads that reference this ServiceAccount
	for _, resource := range allResources {
		if resource.Kind == "Deployment" || resource.Kind == "StatefulSet" || resource.Kind == "DaemonSet" || resource.Kind == "Job" || resource.Kind == "CronJob" {
			if obj, ok := resource.Object.(*unstructured.Unstructured); ok {
				// Check serviceAccountName in pod template spec
				var serviceAccountName string
				var found bool

				if resource.Kind == "CronJob" {
					// CronJob has nested jobTemplate.spec.template.spec structure
					serviceAccountName, found, _ = unstructured.NestedString(obj.Object, "spec", "jobTemplate", "spec", "template", "spec", "serviceAccountName")
				} else if resource.Kind == "Job" {
					// Job has spec.template.spec structure
					serviceAccountName, found, _ = unstructured.NestedString(obj.Object, "spec", "template", "spec", "serviceAccountName")
				} else {
					// Deployment, StatefulSet, DaemonSet have spec.template.spec structure
					serviceAccountName, found, _ = unstructured.NestedString(obj.Object, "spec", "template", "spec", "serviceAccountName")
				}

				if found && serviceAccountName == serviceAccount.Name {
					dependencies = append(dependencies, resource.Name)
				}
			}
		}
	}

	return dependencies
}
