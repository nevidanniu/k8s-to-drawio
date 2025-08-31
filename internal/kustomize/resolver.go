package kustomize

import (
	"fmt"
	"path/filepath"

	"k8s-to-drawio/pkg/models"
)

type Resolver struct{}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) ResolveDependencies(collection *models.ResourceCollection, baseDir string) error {
	// Resolve base references
	for i, resource := range collection.Resources {
		if err := r.resolveResourcePaths(resource, baseDir); err != nil {
			return fmt.Errorf("failed to resolve paths for %s: %w", resource.Name, err)
		}
		collection.Resources[i] = resource
	}

	return nil
}

func (r *Resolver) resolveResourcePaths(resource models.K8sResource, baseDir string) error {
	// Resolve any relative paths in the resource
	// This is a simplified implementation
	return nil
}

func (r *Resolver) findBasePath(overlay, base string) (string, error) {
	if filepath.IsAbs(base) {
		return base, nil
	}

	// Resolve relative to overlay directory
	overlayDir := filepath.Dir(overlay)
	resolvedPath := filepath.Join(overlayDir, base)

	// Clean the path
	return filepath.Clean(resolvedPath), nil
}
