package converter

import (
	"k8s-to-drawio/pkg/models"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) MapResourceToDiagramNode(resource models.K8sResource, nodeID string) models.DiagramNode {
	return models.DiagramNode{
		ID:        nodeID,
		Label:     resource.Name,
		Kind:      resource.Kind,
		Namespace: resource.Namespace,
		X:         0, // Will be set by layout algorithm
		Y:         0, // Will be set by layout algorithm
		Width:     120,
		Height:    60,
		Style:     m.getStyleForKind(resource.Kind),
	}
}

func (m *Mapper) MapDependenciesToConnections(dependencies map[string][]string, nodeMap map[string]string) []models.Connection {
	var connections []models.Connection

	for resourceName, deps := range dependencies {
		sourceID, sourceExists := nodeMap[resourceName]
		if !sourceExists {
			continue
		}

		for _, depName := range deps {
			targetID, targetExists := nodeMap[depName]
			if !targetExists {
				continue
			}

			connection := models.Connection{
				SourceID: sourceID,
				TargetID: targetID,
				Label:    m.getLabelForConnection(resourceName, depName),
				Style:    "default",
			}
			connections = append(connections, connection)
		}
	}

	return connections
}

func (m *Mapper) getStyleForKind(kind string) string {
	switch kind {
	case "Deployment", "StatefulSet", "DaemonSet":
		return "workload"
	case "Service":
		return "service"
	case "Ingress":
		return "ingress"
	case "ConfigMap", "Secret":
		return "config"
	case "PersistentVolumeClaim", "PersistentVolume":
		return "storage"
	default:
		return "default"
	}
}

func (m *Mapper) getLabelForConnection(source, target string) string {
	return "uses"
}
