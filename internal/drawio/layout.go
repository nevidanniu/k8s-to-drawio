package drawio

import (
	"k8s-to-drawio/pkg/models"
	"math"
)

type Layout struct {
	algorithm    string
	noNamespaces bool
}

func NewLayout(algorithm string, noNamespaces bool) *Layout {
	return &Layout{
		algorithm:    algorithm,
		noNamespaces: noNamespaces,
	}
}

func (l *Layout) ApplyLayout(diagram *models.Diagram) error {
	switch l.algorithm {
	case "hierarchical":
		return l.applyHierarchicalLayout(diagram)
	case "grid":
		return l.applyGridLayout(diagram)
	case "vertical":
		return l.applyVerticalLayout(diagram)
	default:
		return l.applyHierarchicalLayout(diagram)
	}
}

func (l *Layout) applyHierarchicalLayout(diagram *models.Diagram) error {
	if l.noNamespaces {
		// Flat layout without namespace grouping
		return l.applyFlatHierarchicalLayout(diagram)
	}

	// Group nodes by namespace
	namespaceNodes := make(map[string][]int)
	for i, node := range diagram.Nodes {
		ns := node.Namespace
		if ns == "" {
			ns = "default"
		}
		namespaceNodes[ns] = append(namespaceNodes[ns], i)
	}

	// Layout each namespace separately
	currentY := 80.0

	for ns, nodeIndices := range namespaceNodes {
		nsX := 80.0
		nsY := currentY
		maxX := nsX
		maxY := nsY

		// Sort nodes by dependencies (topological sort)
		sortedIndices := l.topologicalSort(nodeIndices, diagram)

		// Layout nodes in levels
		levels := l.groupIntoLevels(sortedIndices, diagram)

		levelY := nsY + 80 // Space for namespace header
		for _, level := range levels {
			levelX := nsX + 80
			maxLevelHeight := 0.0

			for _, nodeIdx := range level {
				diagram.Nodes[nodeIdx].X = levelX
				diagram.Nodes[nodeIdx].Y = levelY
				diagram.Nodes[nodeIdx].Width = 140
				diagram.Nodes[nodeIdx].Height = 80

				levelX += 220 // Horizontal spacing
				if diagram.Nodes[nodeIdx].Height > maxLevelHeight {
					maxLevelHeight = diagram.Nodes[nodeIdx].Height
				}
				if levelX > maxX {
					maxX = levelX
				}
			}
			levelY += maxLevelHeight + 80 // Vertical spacing
			if levelY > maxY {
				maxY = levelY
			}
		}

		// Create namespace group
		nsWidth := maxX - nsX + 80
		nsHeight := maxY - nsY + 80

		diagram.Namespaces[ns] = models.NamespaceGroup{
			Name:    ns,
			X:       nsX,
			Y:       nsY,
			Width:   nsWidth,
			Height:  nsHeight,
			NodeIDs: l.getNodeIDs(nodeIndices, diagram),
		}

		currentY = maxY + 150 // Space between namespaces
	}

	return nil
}

func (l *Layout) applyFlatHierarchicalLayout(diagram *models.Diagram) error {
	// Layout all nodes in hierarchical fashion without namespace grouping
	if len(diagram.Nodes) == 0 {
		return nil
	}

	// Create a slice of all node indices
	nodeIndices := make([]int, len(diagram.Nodes))
	for i := range diagram.Nodes {
		nodeIndices[i] = i
	}

	// Sort nodes by dependencies
	sortedIndices := l.topologicalSort(nodeIndices, diagram)

	// Layout nodes in levels
	levels := l.groupIntoLevels(sortedIndices, diagram)

	startX := 80.0
	startY := 80.0
	levelY := startY

	for _, level := range levels {
		levelX := startX
		maxLevelHeight := 0.0

		for _, nodeIdx := range level {
			diagram.Nodes[nodeIdx].X = levelX
			diagram.Nodes[nodeIdx].Y = levelY
			diagram.Nodes[nodeIdx].Width = 140
			diagram.Nodes[nodeIdx].Height = 80

			levelX += 220 // Horizontal spacing
			if diagram.Nodes[nodeIdx].Height > maxLevelHeight {
				maxLevelHeight = diagram.Nodes[nodeIdx].Height
			}
		}
		levelY += maxLevelHeight + 80 // Vertical spacing
	}

	return nil
}

func (l *Layout) applyGridLayout(diagram *models.Diagram) error {
	cols := int(math.Ceil(math.Sqrt(float64(len(diagram.Nodes)))))
	cellWidth := 220.0
	cellHeight := 150.0

	for i, _ := range diagram.Nodes {
		row := i / cols
		col := i % cols
		diagram.Nodes[i].X = float64(col)*cellWidth + 80
		diagram.Nodes[i].Y = float64(row)*cellHeight + 80
		diagram.Nodes[i].Width = 140
		diagram.Nodes[i].Height = 80
	}

	return nil
}

func (l *Layout) applyVerticalLayout(diagram *models.Diagram) error {
	if l.noNamespaces {
		// Flat vertical layout without namespace grouping
		return l.applyFlatVerticalLayout(diagram)
	}

	// Group nodes by namespace
	namespaceNodes := make(map[string][]int)
	for i, node := range diagram.Nodes {
		ns := node.Namespace
		if ns == "" {
			ns = "default"
		}
		namespaceNodes[ns] = append(namespaceNodes[ns], i)
	}

	// Layout each namespace separately in vertical columns
	currentX := 80.0
	namespaceSpacing := 300.0 // Space between namespace columns

	for ns, nodeIndices := range namespaceNodes {
		nsX := currentX
		nsY := 80.0
		maxY := nsY

		// Position nodes vertically within the namespace
		nodeY := nsY + 80 // Space for namespace header
		for _, nodeIdx := range nodeIndices {
			diagram.Nodes[nodeIdx].X = nsX + 80 // Offset from namespace border
			diagram.Nodes[nodeIdx].Y = nodeY
			diagram.Nodes[nodeIdx].Width = 140
			diagram.Nodes[nodeIdx].Height = 80

			nodeY += 160 // Vertical spacing between nodes
			if nodeY > maxY {
				maxY = nodeY
			}
		}

		// Create namespace group
		nsWidth := 220.0 // Fixed width for vertical layout
		nsHeight := maxY - nsY + 80

		diagram.Namespaces[ns] = models.NamespaceGroup{
			Name:    ns,
			X:       nsX,
			Y:       nsY,
			Width:   nsWidth,
			Height:  nsHeight,
			NodeIDs: l.getNodeIDs(nodeIndices, diagram),
		}

		currentX += namespaceSpacing // Move to next namespace column
	}

	return nil
}

func (l *Layout) applyFlatVerticalLayout(diagram *models.Diagram) error {
	// Layout all nodes vertically in a single column without namespace grouping
	if len(diagram.Nodes) == 0 {
		return nil
	}

	startX := 80.0
	startY := 80.0
	nodeY := startY

	for i := range diagram.Nodes {
		diagram.Nodes[i].X = startX
		diagram.Nodes[i].Y = nodeY
		diagram.Nodes[i].Width = 140
		diagram.Nodes[i].Height = 80

		nodeY += 160 // Vertical spacing between nodes
	}

	return nil
}

func (l *Layout) topologicalSort(nodeIndices []int, diagram *models.Diagram) []int {
	// Simple topological sort based on dependencies
	// For simplicity, just return the original order
	// In a real implementation, this would perform proper topological sorting
	return nodeIndices
}

func (l *Layout) groupIntoLevels(nodeIndices []int, diagram *models.Diagram) [][]int {
	// Group nodes into levels based on their dependencies
	levels := make([][]int, 0)

	// For simplicity, put all nodes in one level
	// In a real implementation, this would group by dependency depth
	if len(nodeIndices) > 0 {
		levels = append(levels, nodeIndices)
	}

	return levels
}

func (l *Layout) getNodeIDs(nodeIndices []int, diagram *models.Diagram) []string {
	ids := make([]string, len(nodeIndices))
	for i, idx := range nodeIndices {
		ids[i] = diagram.Nodes[idx].ID
	}
	return ids
}

func (l *Layout) findNodeIndex(nodeID string, diagram *models.Diagram) int {
	for i, node := range diagram.Nodes {
		if node.ID == nodeID {
			return i
		}
	}
	return -1
}
