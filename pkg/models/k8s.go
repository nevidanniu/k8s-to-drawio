package models

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// K8sResource represents a generic Kubernetes resource
type K8sResource struct {
	Object      runtime.Object
	Kind        string
	Name        string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
}

// ResourceCollection holds all parsed Kubernetes resources
type ResourceCollection struct {
	Resources    []K8sResource
	Dependencies map[string][]string // resource name -> dependent resource names
}

// DiagramNode represents a node in the diagram
type DiagramNode struct {
	ID          string
	Label       string
	Kind        string
	Namespace   string
	X           float64
	Y           float64
	Width       float64
	Height      float64
	Style       string
	Connections []Connection
}

// Connection represents a connection between nodes
type Connection struct {
	SourceID string
	TargetID string
	Label    string
	Style    string
}

// Diagram represents the complete diagram structure
type Diagram struct {
	Nodes       []DiagramNode
	Connections []Connection
	Layout      string
	Namespaces  map[string]NamespaceGroup
}

// NamespaceGroup represents a namespace grouping in the diagram
type NamespaceGroup struct {
	Name    string
	X       float64
	Y       float64
	Width   float64
	Height  float64
	NodeIDs []string
}
