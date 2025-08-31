package models

// DiagramElement represents a Draw.io diagram element
type DiagramElement struct {
	XMLContent string
}

// DrawIODiagram represents the complete Draw.io diagram
type DrawIODiagram struct {
	Elements []DiagramElement
	Width    float64
	Height   float64
}
