package k8s

import (
	"fmt"
	"k8s-to-drawio/pkg/models"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(collection *models.ResourceCollection) error {
	for _, resource := range collection.Resources {
		if err := v.validateResource(resource); err != nil {
			return fmt.Errorf("validation failed for %s/%s: %w", resource.Kind, resource.Name, err)
		}
	}

	if err := v.validateDependencies(collection); err != nil {
		return fmt.Errorf("dependency validation failed: %w", err)
	}

	return nil
}

func (v *Validator) validateResource(resource models.K8sResource) error {
	if resource.Name == "" {
		return fmt.Errorf("resource name is required")
	}

	if resource.Kind == "" {
		return fmt.Errorf("resource kind is required")
	}

	// Add more specific validations based on resource type
	switch resource.Kind {
	case "Service":
		return v.validateService(resource)
	case "Deployment":
		return v.validateDeployment(resource)
	case "Ingress":
		return v.validateIngress(resource)
	case "Route":
		return v.validateRoute(resource)
	case "ServiceMonitor":
		return v.validateServiceMonitor(resource)
	}

	return nil
}

func (v *Validator) validateService(resource models.K8sResource) error {
	// Validate service-specific requirements
	return nil
}

func (v *Validator) validateDeployment(resource models.K8sResource) error {
	// Validate deployment-specific requirements
	return nil
}

func (v *Validator) validateIngress(resource models.K8sResource) error {
	// Validate ingress-specific requirements
	return nil
}

func (v *Validator) validateRoute(resource models.K8sResource) error {
	// Validate route-specific requirements
	return nil
}

func (v *Validator) validateServiceMonitor(resource models.K8sResource) error {
	// Validate servicemonitor-specific requirements
	return nil
}

func (v *Validator) validateDependencies(collection *models.ResourceCollection) error {
	// For now, skip circular dependency validation as Kubernetes naturally has
	// circular references (Ingress -> Service -> Deployment) that are normal
	// TODO: Implement more sophisticated dependency validation that only catches
	// problematic circular dependencies
	return nil
}

func (v *Validator) hasCycle(node string, graph map[string][]string, visited, recStack map[string]bool) bool {
	visited[node] = true
	recStack[node] = true

	for _, neighbor := range graph[node] {
		if !visited[neighbor] && v.hasCycle(neighbor, graph, visited, recStack) {
			return true
		} else if recStack[neighbor] {
			return true
		}
	}

	recStack[node] = false
	return false
}
