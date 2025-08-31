package drawio

import (
	"fmt"
	"k8s-to-drawio/pkg/models"
	"strings"
)

type Generator struct {
	layout       *Layout
	noNamespaces bool
}

func NewGenerator(layoutAlgorithm string, noNamespaces bool) *Generator {
	return &Generator{
		layout:       NewLayout(layoutAlgorithm, noNamespaces),
		noNamespaces: noNamespaces,
	}
}

func (g *Generator) Generate(diagram *models.Diagram) (string, error) {
	// Apply layout
	if err := g.layout.ApplyLayout(diagram); err != nil {
		return "", fmt.Errorf("failed to apply layout: %w", err)
	}

	var xmlParts []string

	// Start XML document
	xmlParts = append(xmlParts, `<?xml version="1.0" encoding="UTF-8"?>`)
	xmlParts = append(xmlParts, `<mxfile host="Electron" modified="2023-01-01T00:00:00.000Z" agent="k8s-to-drawio" version="1.0.0" etag="k8s-diagram" type="device">`)
	xmlParts = append(xmlParts, `  <diagram id="k8s-diagram" name="Kubernetes Architecture">`)
	xmlParts = append(xmlParts, `    <mxGraphModel dx="1422" dy="794" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="827" pageHeight="1169" math="0" shadow="0">`)
	xmlParts = append(xmlParts, `      <root>`)
	xmlParts = append(xmlParts, `        <mxCell id="0"/>`)
	xmlParts = append(xmlParts, `        <mxCell id="1" parent="0"/>`)

	// Generate namespace groups
	for _, namespace := range diagram.Namespaces {
		// Format namespace label based on namespace name
		var namespaceLabel string
		if namespace.Name == "vaultstore" {
			// For vaultstore namespace, just show the name without "Namespace:" prefix
			namespaceLabel = namespace.Name
		} else {
			// For other namespaces, show "Namespace: name" format
			namespaceLabel = fmt.Sprintf("Namespace: %s", EscapeXML(namespace.Name))
		}

		namespaceXML := FormatNamespaceGroup(
			fmt.Sprintf("ns-%s", namespace.Name),
			namespaceLabel,
			namespace.X,
			namespace.Y,
			namespace.Width,
			namespace.Height,
		)
		xmlParts = append(xmlParts, "        "+namespaceXML)
	}

	// Generate nodes
	for _, node := range diagram.Nodes {
		template := GetShapeTemplate(node.Kind)

		// Format node label based on kind
		var nodeLabel string
		if node.Kind == "VaultSecret" {
			// For VaultSecret, just show the path since the shape indicates it's a vault secret
			nodeLabel = node.Label
		} else {
			// For other resources, show both kind and name
			nodeLabel = fmt.Sprintf("%s\n%s", node.Kind, node.Label)
		}

		nodeXML := FormatShape(
			template,
			node.ID,
			EscapeXML(nodeLabel),
			node.X,
			node.Y,
			node.Width,
			node.Height,
		)
		xmlParts = append(xmlParts, "        "+nodeXML)
	}

	// Generate connections
	for i, connection := range diagram.Connections {
		connectionID := fmt.Sprintf("conn-%d", i)
		connectionXML := FormatConnection(
			connectionID,
			EscapeXML(connection.Label),
			connection.SourceID,
			connection.TargetID,
			0, 0, 0, 0, // Points will be calculated by Draw.io
		)
		xmlParts = append(xmlParts, "        "+connectionXML)
	}

	// End XML document
	xmlParts = append(xmlParts, `      </root>`)
	xmlParts = append(xmlParts, `    </mxGraphModel>`)
	xmlParts = append(xmlParts, `  </diagram>`)
	xmlParts = append(xmlParts, `</mxfile>`)

	return strings.Join(xmlParts, "\n"), nil
}
