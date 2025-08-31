package kustomize

import (
	"fmt"
	"os"

	"k8s-to-drawio/internal/k8s"
	"k8s-to-drawio/pkg/models"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type Processor struct {
	parser *k8s.Parser
}

func NewProcessor(namespace string) *Processor {
	return &Processor{
		parser: k8s.NewParser(namespace),
	}
}

func (p *Processor) Process(dir string) (*models.ResourceCollection, error) {
	// Check if kustomization.yaml exists
	kustomizationPath := fmt.Sprintf("%s/kustomization.yaml", dir)
	if _, err := os.Stat(kustomizationPath); os.IsNotExist(err) {
		kustomizationPath = fmt.Sprintf("%s/kustomization.yml", dir)
		if _, err := os.Stat(kustomizationPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("no kustomization.yaml found in %s", dir)
		}
	}

	// Create file system
	fSys := filesys.MakeFsOnDisk()

	// Build kustomization
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	resMap, err := k.Run(fSys, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to build kustomization: %w", err)
	}

	// Convert to YAML
	yaml, err := resMap.AsYaml()
	if err != nil {
		return nil, fmt.Errorf("failed to convert resources to YAML: %w", err)
	}

	// Create temporary file with generated YAML
	tempFile, err := os.CreateTemp("", "kustomized-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(yaml); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tempFile.Close()

	// Parse the generated YAML
	return p.parser.ParseFile(tempFile.Name())
}
