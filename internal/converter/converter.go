package converter

import (
	"fmt"
	"os"
	"strings"

	"k8s-to-drawio/internal/drawio"
	"k8s-to-drawio/internal/k8s"
	"k8s-to-drawio/internal/kustomize"
	"k8s-to-drawio/pkg/models"
)

type Config struct {
	InputDir     string
	OutputFile   string
	UseKustomize bool
	Namespace    string
	Layout       string
	NoNamespaces bool
}

type Converter struct {
	config Config
}

func New(config Config) *Converter {
	return &Converter{
		config: config,
	}
}

func (c *Converter) Convert() error {
	// Parse Kubernetes resources
	var collection *models.ResourceCollection
	var err error

	if c.config.UseKustomize {
		processor := kustomize.NewProcessor(c.config.Namespace)
		collection, err = processor.Process(c.config.InputDir)
	} else {
		parser := k8s.NewParser(c.config.Namespace)
		collection, err = parser.ParseDirectory(c.config.InputDir)
	}

	if err != nil {
		return fmt.Errorf("failed to parse resources: %w", err)
	}

	// Validate resources
	validator := k8s.NewValidator()
	if err := validator.Validate(collection); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Convert to diagram
	diagram, err := c.convertToDiagram(collection)
	if err != nil {
		return fmt.Errorf("failed to convert to diagram: %w", err)
	}

	// Generate Draw.io XML
	generator := drawio.NewGenerator(c.config.Layout, c.config.NoNamespaces)
	xml, err := generator.Generate(diagram)
	if err != nil {
		return fmt.Errorf("failed to generate Draw.io XML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(c.config.OutputFile, []byte(xml), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully converted %d resources to %s\n", len(collection.Resources), c.config.OutputFile)
	return nil
}

func (c *Converter) Validate() error {
	// Parse Kubernetes resources
	var collection *models.ResourceCollection
	var err error

	if c.config.UseKustomize {
		processor := kustomize.NewProcessor(c.config.Namespace)
		collection, err = processor.Process(c.config.InputDir)
	} else {
		parser := k8s.NewParser(c.config.Namespace)
		collection, err = parser.ParseDirectory(c.config.InputDir)
	}

	if err != nil {
		return fmt.Errorf("failed to parse resources: %w", err)
	}

	// Validate resources
	validator := k8s.NewValidator()
	if err := validator.Validate(collection); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("Successfully validated %d resources\n", len(collection.Resources))
	return nil
}

func (c *Converter) convertToDiagram(collection *models.ResourceCollection) (*models.Diagram, error) {
	diagram := &models.Diagram{
		Nodes:       make([]models.DiagramNode, 0),
		Connections: make([]models.Connection, 0),
		Layout:      c.config.Layout,
		Namespaces:  make(map[string]models.NamespaceGroup),
	}

	// Debug: Print dependencies
	// fmt.Printf("Dependencies found: %+v\n", collection.Dependencies)

	// Create nodes for each resource (excluding Namespace resources which are represented as containers)
	nodeIndex := 0
	for _, resource := range collection.Resources {
		// Skip Namespace and Kustomization resources as they should be containers/metadata, not nodes
		if resource.Kind == "Namespace" || resource.Kind == "Kustomization" {
			continue
		}

		node := models.DiagramNode{
			ID:        fmt.Sprintf("node-%d", nodeIndex),
			Label:     resource.Name,
			Kind:      resource.Kind,
			Namespace: resource.Namespace,
			X:         0, // Will be set by layout algorithm
			Y:         0, // Will be set by layout algorithm
			Width:     120,
			Height:    60,
		}
		diagram.Nodes = append(diagram.Nodes, node)
		nodeIndex++
	}

	// Collect all virtual Vault secrets referenced in dependencies
	virtualVaultSecrets := make(map[string]bool)
	for _, dependencies := range collection.Dependencies {
		for _, depName := range dependencies {
			if strings.HasPrefix(depName, "vault-secret-") {
				virtualVaultSecrets[depName] = true
			}
		}
	}

	// Create virtual nodes for Vault secrets
	// nodeIndex is already set to the count of non-Namespace resources
	for vaultSecretName := range virtualVaultSecrets {
		// Extract the original path from the virtual name
		// "vault-secret-myapp-config" -> "secret/myapp/config"
		originalPath := strings.TrimPrefix(vaultSecretName, "vault-secret-")
		originalPath = strings.ReplaceAll(originalPath, "-", "/")
		// Re-add the "secret/" prefix if it doesn't already exist
		if !strings.HasPrefix(originalPath, "secret/") {
			originalPath = "secret/" + originalPath
		}

		node := models.DiagramNode{
			ID:        fmt.Sprintf("node-%d", nodeIndex),
			Label:     originalPath,
			Kind:      "VaultSecret",
			Namespace: "vaultstore",
			X:         0, // Will be set by layout algorithm
			Y:         0, // Will be set by layout algorithm
			Width:     140,
			Height:    80,
		}
		diagram.Nodes = append(diagram.Nodes, node)
		nodeIndex++
	}

	// Create connections based on dependencies
	nodeMap := make(map[string]string)     // resource name -> node ID
	resourceMap := make(map[string]string) // "kind/name" -> node ID

	// Map actual Kubernetes resources (excluding Namespace resources)
	nodeIndex = 0
	for _, resource := range collection.Resources {
		// Skip Namespace and Kustomization resources as they are not represented as nodes
		if resource.Kind == "Namespace" || resource.Kind == "Kustomization" {
			continue
		}

		nodeID := fmt.Sprintf("node-%d", nodeIndex)
		// Use composite key to avoid conflicts between resources with same name
		compositeKey := fmt.Sprintf("%s/%s", resource.Kind, resource.Name)
		// For nodeMap, prioritize workload resources (Deployment, StatefulSet, DaemonSet) over Services
		// This ensures Service selectors resolve to the correct workload resources
		if _, exists := nodeMap[resource.Name]; exists {
			// If name collision, check if current resource should take priority
			if resource.Kind == "Deployment" || resource.Kind == "StatefulSet" || resource.Kind == "DaemonSet" {
				nodeMap[resource.Name] = nodeID // Override with workload resource
			}
			// If existing resource was a workload and current is Service, keep existing
		} else {
			nodeMap[resource.Name] = nodeID
		}
		resourceMap[compositeKey] = nodeID
		nodeIndex++
	}

	// Map virtual Vault secrets
	// nodeIndex is already correctly positioned after non-Namespace resources
	vaultNodeIndex := nodeIndex // VaultSecret nodes start after regular resources
	for vaultSecretName := range virtualVaultSecrets {
		nodeID := fmt.Sprintf("node-%d", vaultNodeIndex)
		nodeMap[vaultSecretName] = nodeID
		resourceMap[fmt.Sprintf("VaultSecret/%s", vaultSecretName)] = nodeID
		// fmt.Printf("Mapping VaultSecret: %s -> %s\n", vaultSecretName, nodeID)
		vaultNodeIndex++
	}

	for resourceKey, dependencies := range collection.Dependencies {
		targetID, targetExists := resourceMap[resourceKey]
		if !targetExists {
			continue
		}

		for _, depName := range dependencies {
			// First check if this dependency should be resolved to a specific resource type
			var sourceID string
			var sourceExists bool

			// Special handling for Route, ServiceMonitor, and Ingress dependencies
			resourceParts := strings.Split(resourceKey, "/")
			if len(resourceParts) == 2 {
				resourceKind := resourceParts[0]
				if resourceKind == "Route" {
					// Routes should connect to Services
					serviceKey := fmt.Sprintf("Service/%s", depName)
					if serviceNodeID, exists := resourceMap[serviceKey]; exists {
						sourceID = serviceNodeID
						sourceExists = true
					}
				} else if resourceKind == "ServiceMonitor" {
					// ServiceMonitors should connect to Services
					serviceKey := fmt.Sprintf("Service/%s", depName)
					if serviceNodeID, exists := resourceMap[serviceKey]; exists {
						sourceID = serviceNodeID
						sourceExists = true
					}
				} else if resourceKind == "Ingress" {
					// Ingress should connect to Services
					serviceKey := fmt.Sprintf("Service/%s", depName)
					if serviceNodeID, exists := resourceMap[serviceKey]; exists {
						sourceID = serviceNodeID
						sourceExists = true
					}
				}
			}

			// Fall back to name-based lookup if specific type lookup failed
			if !sourceExists {
				sourceID, sourceExists = nodeMap[depName]
			}

			if !sourceExists {
				continue
			}

			// Avoid self-references by checking if source and target are the same
			if sourceID == targetID {
				continue
			}

			connection := models.Connection{
				SourceID: targetID,
				TargetID: sourceID,
				Label:    "uses",
				Style:    "default",
			}
			diagram.Connections = append(diagram.Connections, connection)
		}
	}

	return diagram, nil
}
