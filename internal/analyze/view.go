package analyze

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// Color scheme for entity types
var (
	PersonColor   = lipgloss.Color("39")  // Blue
	OrgColor      = lipgloss.Color("34")  // Green
	PlaceColor    = lipgloss.Color("226") // Yellow
	TechColor     = lipgloss.Color("196") // Red
	ConceptColor  = lipgloss.Color("141") // Purple
	ThingColor    = lipgloss.Color("244") // Gray
	DefaultColor  = lipgloss.Color("15")  // White

	SelectedColor    = lipgloss.Color("201") // Magenta
	HighlightedColor = lipgloss.Color("208") // Orange
	EdgeColor        = lipgloss.Color("240") // Dark gray
)

// Styles for rendering
var (
	nodeStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	selectedNodeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(SelectedColor).
				Padding(0, 1)

	highlightedNodeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(HighlightedColor).
				Padding(0, 1)

	edgeStyle = lipgloss.NewStyle().
			Foreground(EdgeColor)

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)

// View implements tea.Model.
func (m Model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.graph == nil || m.graph.NodeCount() == 0 {
		return helpStyle.Render("No graph data to display")
	}

	// Create canvas for rendering
	canvas := NewCanvas(m.width, m.height-3) // Reserve 3 lines for footer

	// Render edges first (so they appear behind nodes)
	m.renderEdges(canvas)

	// Render nodes
	m.renderNodes(canvas)

	// Render selected node details if any
	if m.selectedNodeID != "" {
		m.renderNodeDetails(canvas)
	}

	// Build final view
	var b strings.Builder
	b.WriteString(canvas.Render())
	b.WriteString("\n")
	b.WriteString(m.renderStats())
	b.WriteString("\n")
	b.WriteString(m.renderHelp())

	return b.String()
}

// renderNodes renders all nodes to the canvas.
func (m Model) renderNodes(canvas *Canvas) {
	for id, node := range m.graph.Nodes {
		// Transform to screen coordinates
		screenX, screenY := m.toScreen(node.X, node.Y)

		// Skip if off-screen
		if screenX < 0 || screenX >= canvas.width || screenY < 0 || screenY >= canvas.height {
			continue
		}

		// Choose style based on selection state
		var style lipgloss.Style
		if id == m.selectedNodeID {
			style = selectedNodeStyle
		} else if m.IsHighlighted(id) {
			style = highlightedNodeStyle
		} else {
			style = nodeStyle.Copy().Foreground(m.getEntityColor(node.Entity.Type))
		}

		// Render node label
		label := m.formatNodeLabel(node)
		canvas.DrawText(screenX, screenY, style.Render(label))
	}
}

// renderEdges renders all edges to the canvas.
func (m Model) renderEdges(canvas *Canvas) {
	for _, edge := range m.graph.Edges {
		source := m.graph.GetNode(edge.SourceID)
		target := m.graph.GetNode(edge.TargetID)

		if source == nil || target == nil {
			continue
		}

		// Transform to screen coordinates
		x1, y1 := m.toScreen(source.X, source.Y)
		x2, y2 := m.toScreen(target.X, target.Y)

		// Skip if both endpoints off-screen
		if (x1 < 0 || x1 >= canvas.width) && (x2 < 0 || x2 >= canvas.width) {
			continue
		}
		if (y1 < 0 || y1 >= canvas.height) && (y2 < 0 || y2 >= canvas.height) {
			continue
		}

		// Draw edge line
		canvas.DrawLine(x1, y1, x2, y2, edgeStyle)

		// Draw edge label at midpoint
		midX := (x1 + x2) / 2
		midY := (y1 + y2) / 2

		if midX >= 0 && midX < canvas.width && midY >= 0 && midY < canvas.height {
			label := m.formatEdgeLabel(edge)
			if label != "" {
				canvas.DrawText(midX, midY, edgeStyle.Render(label))
			}
		}
	}
}

// renderNodeDetails renders details of the selected node.
func (m Model) renderNodeDetails(canvas *Canvas) {
	node := m.SelectedNode()
	if node == nil {
		return
	}

	// Render detail box in top-right corner
	x := canvas.width - 30
	y := 1

	if x < 0 {
		return
	}

	details := []string{
		fmt.Sprintf("Entity: %s", node.Entity.Text),
		fmt.Sprintf("Type: %s", node.Entity.Type),
		fmt.Sprintf("Freq: %d", node.Frequency),
		fmt.Sprintf("Importance: %d/10", node.Importance),
	}

	// Count connections
	outEdges := len(m.graph.GetEdgesFrom(node.ID))
	inEdges := len(m.graph.GetEdgesTo(node.ID))
	details = append(details, fmt.Sprintf("Edges: %d out, %d in", outEdges, inEdges))

	for i, line := range details {
		if y+i < canvas.height {
			canvas.DrawText(x, y+i, helpStyle.Render(line))
		}
	}
}

// renderStats renders statistics footer.
func (m Model) renderStats() string {
	stats := m.GetStats()
	return statsStyle.Render(fmt.Sprintf(
		"Nodes: %d  Edges: %d  Layout: %d  Zoom: %.1fx  Offset: (%.0f, %.0f)",
		stats.Nodes, stats.Edges, stats.LayoutSteps, m.zoom, m.offsetX, m.offsetY,
	))
}

// renderHelp renders help text footer.
func (m Model) renderHelp() string {
	if m.selectedNodeID != "" {
		return helpStyle.Render("[hjkl] Pan  [+-] Zoom  [esc] Clear  [c] Center  [s] Stabilize")
	}
	return helpStyle.Render("[hjkl] Pan  [+-] Zoom  [enter] Select  [c] Center  [r] Reset  [s] Stabilize")
}

// toScreen converts graph coordinates to screen coordinates.
func (m Model) toScreen(graphX, graphY float64) (int, int) {
	centerX := float64(m.width) / 2
	centerY := float64(m.height-3) / 2 // Account for footer

	screenX := int(graphX*m.zoom + m.offsetX + centerX)
	screenY := int(graphY*m.zoom + m.offsetY + centerY)

	return screenX, screenY
}

// getEntityColor returns the color for an entity type.
func (m Model) getEntityColor(entityType semantic.EntityType) lipgloss.Color {
	switch entityType {
	case semantic.EntityPerson:
		return PersonColor
	case semantic.EntityOrganization:
		return OrgColor
	case semantic.EntityPlace:
		return PlaceColor
	case semantic.EntityTechnology:
		return TechColor
	case semantic.EntityConcept:
		return ConceptColor
	case semantic.EntityThing:
		return ThingColor
	default:
		return DefaultColor
	}
}

// formatNodeLabel creates a display label for a node.
func (m Model) formatNodeLabel(node *TripleNode) string {
	// Truncate long labels
	label := node.Entity.Text
	if len(label) > 15 {
		label = label[:12] + "..."
	}

	// Add importance indicator for high-importance nodes
	if node.Importance >= 8 {
		label = "★ " + label
	}

	return label
}

// formatEdgeLabel creates a display label for an edge.
func (m Model) formatEdgeLabel(edge *TripleEdge) string {
	// Only show labels for edges connected to selected node
	if m.selectedNodeID != edge.SourceID && m.selectedNodeID != edge.TargetID {
		return ""
	}

	// Truncate predicate
	pred := edge.Relation.Predicate
	if len(pred) > 10 {
		pred = pred[:7] + "..."
	}

	return pred
}

// Canvas represents a 2D character canvas for rendering.
type Canvas struct {
	width  int
	height int
	cells  [][]string // [y][x]
}

// NewCanvas creates a new canvas.
func NewCanvas(width, height int) *Canvas {
	cells := make([][]string, height)
	for y := range cells {
		cells[y] = make([]string, width)
		for x := range cells[y] {
			cells[y][x] = " "
		}
	}
	return &Canvas{
		width:  width,
		height: height,
		cells:  cells,
	}
}

// DrawText draws text at the given position.
func (c *Canvas) DrawText(x, y int, text string) {
	// Strip ANSI codes for length calculation
	// In real implementation, we'd preserve styling
	for i, ch := range text {
		px := x + i
		if px >= 0 && px < c.width && y >= 0 && y < c.height {
			c.cells[y][px] = string(ch)
		}
	}
}

// DrawLine draws a line between two points.
func (c *Canvas) DrawLine(x1, y1, x2, y2 int, style lipgloss.Style) {
	// Bresenham's line algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	sx := 1
	if x1 > x2 {
		sx = -1
	}

	sy := 1
	if y1 > y2 {
		sy = -1
	}

	err := dx - dy

	x, y := x1, y1

	for {
		// Draw character based on direction
		char := "·"
		if dx > dy {
			char = "─"
		} else if dy > dx {
			char = "│"
		}

		if x >= 0 && x < c.width && y >= 0 && y < c.height {
			c.cells[y][x] = style.Render(char)
		}

		if x == x2 && y == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// Render returns the canvas as a string.
func (c *Canvas) Render() string {
	var b strings.Builder
	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			b.WriteString(c.cells[y][x])
		}
		if y < c.height-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
