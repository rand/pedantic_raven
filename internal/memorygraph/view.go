package memorygraph

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles for graph rendering.
var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("235")).
			Bold(true).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Background(lipgloss.Color("235"))

	nodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Background(lipgloss.Color("235")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	selectedNodeStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("237")).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("229")).
				Padding(0, 1).
				Bold(true)

	edgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("242"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Background(lipgloss.Color("235")).
			Padding(1, 2)
)

// View implements tea.Model.
func (m Model) View() string {
	if m.err != nil {
		return m.renderError()
	}

	if m.graph == nil || len(m.graph.Nodes) == 0 {
		return m.renderEmpty()
	}

	// Render header
	header := m.renderHeader()

	// Render graph content
	content := m.renderGraph()

	// Render footer
	footer := m.renderFooter()

	// Combine with newlines
	return header + "\n" + content + "\n" + footer
}

// renderHeader renders the graph view header.
func (m Model) renderHeader() string {
	title := "Memory Graph"
	if m.graph != nil {
		title = fmt.Sprintf("Memory Graph (%d nodes, %d edges)", m.graph.NodeCount(), m.graph.EdgeCount())
	}

	stats := fmt.Sprintf("Zoom: %.1fx | Layout steps: %d", m.zoom, m.layoutSteps)

	padding := m.width - lipgloss.Width(title) - lipgloss.Width(stats)
	if padding < 0 {
		padding = 0
	}

	line := title + strings.Repeat(" ", padding) + stats
	return headerStyle.Width(m.width).Render(line)
}

// renderFooter renders the graph view footer with keyboard hints.
func (m Model) renderFooter() string {
	var hints string

	if m.selectedNodeID != "" {
		// Node is selected
		hints = "h/j/k/l: pan | +/-: zoom | c: center | Enter: navigate | e: expand | x: collapse | Tab: next node"
	} else {
		// No node selected
		hints = "h/j/k/l: pan | +/-: zoom | 0: reset | Space: layout step | r: re-layout | Tab: select node"
	}

	if lipgloss.Width(hints) > m.width {
		hints = hints[:m.width]
	}

	return footerStyle.Width(m.width).Render(hints)
}

// renderGraph renders the graph visualization.
func (m Model) renderGraph() string {
	// Calculate content height (total height - header - footer - separators)
	contentHeight := m.height - 3
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Create a canvas for rendering
	canvas := newCanvas(m.width, contentHeight)

	// Draw edges first (so they appear behind nodes)
	// Only draw visible edges
	for _, edge := range m.graph.Edges {
		if m.IsEdgeVisible(edge) {
			m.drawEdge(canvas, edge)
		}
	}

	// Draw nodes on top of edges
	// Only draw visible nodes
	for _, node := range m.graph.Nodes {
		if m.IsNodeVisible(node.ID) {
			m.drawNode(canvas, node)
		}
	}

	return canvas.Render()
}

// renderError renders an error state.
func (m Model) renderError() string {
	errMsg := fmt.Sprintf("Error: %v", m.err)
	return errorStyle.Width(m.width).Height(m.height).Render(errMsg)
}

// renderEmpty renders an empty state.
func (m Model) renderEmpty() string {
	msg := "No graph to display\n\nLoad a memory to view its graph."
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)
	return style.Render(msg)
}

// drawNode draws a node on the canvas.
func (m Model) drawNode(canvas *Canvas, node *Node) {
	// Transform world coordinates to screen coordinates
	screenX, screenY := m.worldToScreen(node.X, node.Y)

	// Get node label (first 10 chars of ID)
	label := node.ID
	if len(label) > 10 {
		label = label[:10] + "…"
	}

	// Add expansion indicator if node has children
	if m.HasChildren(node.ID) {
		if node.IsExpanded {
			label = "[-] " + label
		} else {
			label = "[+] " + label
		}
	}

	// Choose style based on selection
	style := nodeStyle
	if node.ID == m.selectedNodeID {
		style = selectedNodeStyle
	}

	// Render node box
	box := style.Render(label)

	// Draw at position (centered on node coordinates)
	boxLines := strings.Split(box, "\n")
	boxWidth := lipgloss.Width(box)
	boxHeight := len(boxLines)

	// Center the box on the node position
	startX := screenX - boxWidth/2
	startY := screenY - boxHeight/2

	for i, line := range boxLines {
		canvas.DrawText(startX, startY+i, line)
	}
}

// drawEdge draws an edge between two nodes.
func (m Model) drawEdge(canvas *Canvas, edge *Edge) {
	source := m.graph.Nodes[edge.SourceID]
	target := m.graph.Nodes[edge.TargetID]

	if source == nil || target == nil {
		return
	}

	// Transform world coordinates to screen coordinates
	x1, y1 := m.worldToScreen(source.X, source.Y)
	x2, y2 := m.worldToScreen(target.X, target.Y)

	// Draw line between nodes
	canvas.DrawLine(x1, y1, x2, y2, "242")
}

// worldToScreen transforms world coordinates to screen coordinates.
func (m Model) worldToScreen(worldX, worldY float64) (int, int) {
	// Apply offset and zoom
	screenX := int((worldX+m.offsetX)*m.zoom) + m.width/2
	screenY := int((worldY+m.offsetY)*m.zoom) + (m.height-3)/2

	return screenX, screenY
}

// Canvas represents a drawable area.
type Canvas struct {
	width  int
	height int
	cells  [][]rune
	colors [][]string
}

// newCanvas creates a new canvas.
func newCanvas(width, height int) *Canvas {
	cells := make([][]rune, height)
	colors := make([][]string, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]rune, width)
		colors[y] = make([]string, width)
		for x := 0; x < width; x++ {
			cells[y][x] = ' '
			colors[y][x] = ""
		}
	}
	return &Canvas{
		width:  width,
		height: height,
		cells:  cells,
		colors: colors,
	}
}

// DrawText draws text at the given position.
func (c *Canvas) DrawText(x, y int, text string) {
	// Strip ANSI codes and handle styled text
	lines := strings.Split(text, "\n")
	for dy, line := range lines {
		currentY := y + dy
		if currentY < 0 || currentY >= c.height {
			continue
		}

		// Handle ANSI escape sequences by stripping them
		// For simplicity, just get the raw runes
		runes := []rune(stripAnsi(line))
		for dx, r := range runes {
			currentX := x + dx
			if currentX >= 0 && currentX < c.width {
				c.cells[currentY][currentX] = r
			}
		}
	}
}

// DrawLine draws a line between two points using Bresenham's algorithm.
func (c *Canvas) DrawLine(x1, y1, x2, y2 int, color string) {
	// Bresenham's line algorithm
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy

	x, y := x1, y1
	for {
		// Draw point
		if x >= 0 && x < c.width && y >= 0 && y < c.height {
			// Use different characters for different directions
			if dx > dy {
				c.cells[y][x] = '─'
			} else {
				c.cells[y][x] = '│'
			}
			c.colors[y][x] = color
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

// Render renders the canvas to a string.
func (c *Canvas) Render() string {
	var builder strings.Builder
	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			builder.WriteRune(c.cells[y][x])
		}
		if y < c.height-1 {
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}

// Helper functions

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// stripAnsi strips ANSI escape sequences from a string.
func stripAnsi(s string) string {
	// Simple ANSI stripping - remove escape sequences
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}
