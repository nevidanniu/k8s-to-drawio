package drawio

import (
	"fmt"
	"strings"
)

// ShapeTemplates contains Draw.io shape templates for different Kubernetes resources
var ShapeTemplates = map[string]string{
	"Deployment": `<mxCell id="%s" value="%s" style="rounded=1;whiteSpace=wrap;html=1;fillColor=#d5e8d4;strokeColor=#82b366;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"Service": `<mxCell id="%s" value="%s" style="ellipse;whiteSpace=wrap;html=1;fillColor=#fff2cc;strokeColor=#d6b656;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"Ingress": `<mxCell id="%s" value="%s" style="rhombus;whiteSpace=wrap;html=1;fillColor=#f8cecc;strokeColor=#b85450;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"ConfigMap": `<mxCell id="%s" value="%s" style="shape=note;whiteSpace=wrap;html=1;backgroundOutline=1;darkOpacity=0.05;fillColor=#e1d5e7;strokeColor=#9673a6;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"Secret": `<mxCell id="%s" value="%s" style="shape=note;whiteSpace=wrap;html=1;backgroundOutline=1;darkOpacity=0.05;fillColor=#f5f5f5;strokeColor=#666666;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"PersistentVolumeClaim": `<mxCell id="%s" value="%s" style="shape=cylinder3;whiteSpace=wrap;html=1;boundedLbl=1;backgroundOutline=1;size=15;fillColor=#dae8fc;strokeColor=#6c8ebf;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"StatefulSet": `<mxCell id="%s" value="%s" style="rounded=1;whiteSpace=wrap;html=1;fillColor=#d5e8d4;strokeColor=#82b366;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"DaemonSet": `<mxCell id="%s" value="%s" style="rounded=1;whiteSpace=wrap;html=1;fillColor=#d5e8d4;strokeColor=#82b366;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"VaultSecret": `<mxCell id="%s" value="%s" style="shape=hexagon;whiteSpace=wrap;html=1;backgroundOutline=1;darkOpacity=0.05;fillColor=#ffe6cc;strokeColor=#d79b00;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"Route": `<mxCell id="%s" value="%s" style="shape=trapezoid;whiteSpace=wrap;html=1;fillColor=#ffc9c9;strokeColor=#d6536d;size=0.2;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,

	"ServiceMonitor": `<mxCell id="%s" value="%s" style="shape=monitor;whiteSpace=wrap;html=1;fillColor=#e6f3ff;strokeColor=#4a90e2;" vertex="1" parent="1">
		<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
	</mxCell>`,
}

// ConnectionTemplate for drawing connections between resources
var ConnectionTemplate = `<mxCell id="%s" value="%s" style="endArrow=classic;html=1;rounded=0;" edge="1" parent="1" source="%s" target="%s">
	<mxGeometry width="50" height="50" relative="1" as="geometry">
		<mxPoint x="%.1f" y="%.1f" as="sourcePoint"/>
		<mxPoint x="%.1f" y="%.1f" as="targetPoint"/>
	</mxGeometry>
</mxCell>`

// NamespaceGroupTemplate for namespace groupings
var NamespaceGroupTemplate = `<mxCell id="%s" value="%s" style="swimlane;fontStyle=0;childLayout=stackLayout;horizontal=1;startSize=30;horizontalStack=0;resizeParent=1;resizeParentMax=0;resizeLast=0;collapsible=1;marginBottom=0;fillColor=#e1d5e7;strokeColor=#9673a6;" vertex="1" parent="1">
	<mxGeometry x="%.1f" y="%.1f" width="%.1f" height="%.1f" as="geometry"/>
</mxCell>`

func GetShapeTemplate(kind string) string {
	if template, exists := ShapeTemplates[kind]; exists {
		return template
	}
	// Default template for unknown resource types
	return ShapeTemplates["Deployment"]
}

func FormatShape(template, id, label string, x, y, width, height float64) string {
	return fmt.Sprintf(template, id, label, x, y, width, height)
}

func FormatConnection(id, label, sourceID, targetID string, x1, y1, x2, y2 float64) string {
	return fmt.Sprintf(ConnectionTemplate, id, label, sourceID, targetID, x1, y1, x2, y2)
}

func FormatNamespaceGroup(id, name string, x, y, width, height float64) string {
	return fmt.Sprintf(NamespaceGroupTemplate, id, name, x, y, width, height)
}

func EscapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
