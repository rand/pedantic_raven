// Package orchestrate provides task graph visualization for work plan dependencies.
package orchestrate

import (
	"fmt"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TaskGraph visualizes a WorkPlan as a directed acyclic graph with force-directed layout.
type TaskGraph struct {
	// Graph data
	plan  *WorkPlan
	nodes map[string]*GraphNode // task ID -> node
	edges []GraphEdge

	// Force-directed layout state
	velocities map[string]Velocity
	damping    float64 // 0.8 recommended

	// Viewport state
	offsetX float64
	offsetY float64
	zoom    float64 // 1.0 = 100%

	// Canvas dimensions
	width  int
	height int

	// Selection state
	selected string // Selected task ID

	// Task status tracking
	statuses map[string]TaskStatus

	// Layout iteration counter
	layoutSteps int
}

// GraphNode represents a task node in the dependency graph.
type GraphNode struct {
	TaskID      string
	Description string
	Status      TaskStatus
	X           float64
	Y           float64
}

// GraphEdge represents a dependency edge between tasks.
type GraphEdge struct {
	From string // Source task ID (dependency)
	To   string // Target task ID (dependent)
}

// Position represents 2D coordinates.
type Position struct {
	X float64
	Y float64
}

// Velocity represents 2D velocity for force-directed layout.
type Velocity struct {
	VX float64
	VY float64
}

// Layout constants (adapted from analyze/triple_graph.go)
const (
	RepulsionStrength = 100.0 // How strongly nodes repel each other
	AttractionStrength = 0.01 // How strongly edges attract nodes
	MaxForce          = 10.0  // Maximum force applied per iteration
	MinDistance       = 1.0   // Minimum distance to prevent division by zero
	IdealDistance     = 10.0  // Ideal distance between connected nodes
)

// Color scheme for task status
var (
	PendingColor   = lipgloss.Color("8")   // Gray
	ActiveColor    = lipgloss.Color("11")  // Yellow
	CompletedColor = lipgloss.Color("10")  // Green
	FailedColor    = lipgloss.Color("9")   // Red
	SelectedColor  = lipgloss.Color("201") // Magenta
	EdgeColor      = lipgloss.Color("240") // Dark gray
)

// Styles for rendering
var (
	graphNodeStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	graphSelectedNodeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(SelectedColor).
				Padding(0, 1)

	graphEdgeStyle = lipgloss.NewStyle().
			Foreground(EdgeColor)

	graphStatsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)

	graphHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)
)

// NewTaskGraph creates a new task graph from a work plan.
func NewTaskGraph(plan *WorkPlan, width, height int) (*TaskGraph, error) {
	if plan == nil {
		return nil, fmt.Errorf("work plan cannot be nil")
	}

	if err := plan.Validate(); err != nil {
		return nil, fmt.Errorf("invalid work plan: %w", err)
	}

	tg := &TaskGraph{
		plan:       plan,
		nodes:      make(map[string]*GraphNode),
		edges:      []GraphEdge{},
		velocities: make(map[string]Velocity),
		damping:    0.8,
		offsetX:    0,
		offsetY:    0,
		zoom:       1.0,
		width:      width,
		height:     height,
		selected:   "",
		statuses:   make(map[string]TaskStatus),
	}

	// Build DAG
	if err := tg.buildDAG(); err != nil {
		return nil, err
	}

	// Initialize layout
	tg.initializeLayout()

	return tg, nil
}

// buildDAG constructs the dependency graph from the work plan.
func (tg *TaskGraph) buildDAG() error {
	// Create nodes from tasks
	for _, task := range tg.plan.Tasks {
		tg.nodes[task.ID] = &GraphNode{
			TaskID:      task.ID,
			Description: task.Description,
			Status:      TaskStatusPending,
			X:           0,
			Y:           0,
		}
		tg.statuses[task.ID] = TaskStatusPending
	}

	// Create edges from dependencies
	for _, task := range tg.plan.Tasks {
		for _, depID := range task.Dependencies {
			// Edge points from dependency to dependent
			tg.edges = append(tg.edges, GraphEdge{
				From: depID,
				To:   task.ID,
			})
		}
	}

	return nil
}

// initializeLayout sets initial node positions in a circular layout.
func (tg *TaskGraph) initializeLayout() {
	nodeCount := len(tg.nodes)
	if nodeCount == 0 {
		return
	}

	i := 0
	radius := 20.0
	for _, node := range tg.nodes {
		angle := (float64(i) / float64(nodeCount)) * 2 * math.Pi
		node.X = math.Cos(angle) * radius
		node.Y = math.Sin(angle) * radius
		tg.velocities[node.TaskID] = Velocity{VX: 0, VY: 0}
		i++
	}
}

// applyRepulsion applies repulsive forces between all nodes.
func (tg *TaskGraph) applyRepulsion() {
	nodes := make([]*GraphNode, 0, len(tg.nodes))
	for _, node := range tg.nodes {
		nodes = append(nodes, node)
	}

	// Apply repulsion between all pairs
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			n1 := nodes[i]
			n2 := nodes[j]

			// Calculate distance
			dx := n2.X - n1.X
			dy := n2.Y - n1.Y
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist < MinDistance {
				dist = MinDistance
			}

			// Calculate repulsive force (inverse square)
			force := RepulsionStrength / (dist * dist)

			// Limit force magnitude
			if force > MaxForce {
				force = MaxForce
			}

			// Apply force in opposite directions
			fx := (dx / dist) * force
			fy := (dy / dist) * force

			v1 := tg.velocities[n1.TaskID]
			v1.VX -= fx
			v1.VY -= fy
			tg.velocities[n1.TaskID] = v1

			v2 := tg.velocities[n2.TaskID]
			v2.VX += fx
			v2.VY += fy
			tg.velocities[n2.TaskID] = v2
		}
	}
}

// applyAttraction applies attractive forces along edges.
func (tg *TaskGraph) applyAttraction() {
	for _, edge := range tg.edges {
		from := tg.nodes[edge.From]
		to := tg.nodes[edge.To]

		if from == nil || to == nil {
			continue
		}

		// Calculate distance
		dx := to.X - from.X
		dy := to.Y - from.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < MinDistance {
			continue
		}

		// Calculate attractive force (Hooke's law)
		force := (dist - IdealDistance) * AttractionStrength

		// Limit force magnitude
		if math.Abs(force) > MaxForce {
			if force > 0 {
				force = MaxForce
			} else {
				force = -MaxForce
			}
		}

		// Apply force toward each other
		fx := (dx / dist) * force
		fy := (dy / dist) * force

		vFrom := tg.velocities[edge.From]
		vFrom.VX += fx
		vFrom.VY += fy
		tg.velocities[edge.From] = vFrom

		vTo := tg.velocities[edge.To]
		vTo.VX -= fx
		vTo.VY -= fy
		tg.velocities[edge.To] = vTo
	}
}

// updatePositions updates node positions based on velocities.
func (tg *TaskGraph) updatePositions() {
	for id, node := range tg.nodes {
		v := tg.velocities[id]

		// Apply damping
		v.VX *= tg.damping
		v.VY *= tg.damping

		// Update position
		node.X += v.VX
		node.Y += v.VY

		tg.velocities[id] = v
	}
}

// applyForceIteration applies one iteration of force-directed layout.
func (tg *TaskGraph) applyForceIteration() {
	tg.applyRepulsion()
	tg.applyAttraction()
	tg.updatePositions()
	tg.layoutSteps++
}

// stabilize runs multiple layout iterations to converge to stable positions (internal).
func (tg *TaskGraph) stabilize(iterations int) {
	for i := 0; i < iterations; i++ {
		tg.applyForceIteration()
	}
}

// Stabilize runs multiple layout iterations to converge to stable positions (public API).
func (tg *TaskGraph) Stabilize(iterations int) {
	tg.stabilize(iterations)
}

// getBounds returns the bounding box of all nodes.
func (tg *TaskGraph) getBounds() (minX, maxX, minY, maxY float64) {
	if len(tg.nodes) == 0 {
		return 0, 0, 0, 0
	}

	first := true
	for _, node := range tg.nodes {
		if first {
			minX = node.X
			maxX = node.X
			minY = node.Y
			maxY = node.Y
			first = false
		} else {
			if node.X < minX {
				minX = node.X
			}
			if node.X > maxX {
				maxX = node.X
			}
			if node.Y < minY {
				minY = node.Y
			}
			if node.Y > maxY {
				maxY = node.Y
			}
		}
	}

	return minX, maxX, minY, maxY
}

// UpdateStatus updates the status of a task node.
func (tg *TaskGraph) UpdateStatus(taskID string, status TaskStatus) {
	if node, ok := tg.nodes[taskID]; ok {
		node.Status = status
		tg.statuses[taskID] = status
	}
}

// SelectNode selects a node by task ID.
func (tg *TaskGraph) SelectNode(taskID string) {
	if _, ok := tg.nodes[taskID]; ok {
		tg.selected = taskID
	}
}

// ClearSelection clears the current node selection.
func (tg *TaskGraph) ClearSelection() {
	tg.selected = ""
}

// Pan moves the viewport by the given offset.
func (tg *TaskGraph) Pan(dx, dy float64) {
	tg.offsetX += dx
	tg.offsetY += dy
}

// Zoom adjusts the zoom level by the given factor.
func (tg *TaskGraph) Zoom(factor float64) {
	tg.zoom *= factor
	if tg.zoom < 0.1 {
		tg.zoom = 0.1
	}
	if tg.zoom > 5.0 {
		tg.zoom = 5.0
	}
}

// Center resets the viewport to center on the graph.
func (tg *TaskGraph) Center() {
	tg.offsetX = 0
	tg.offsetY = 0
	tg.zoom = 1.0
}

// Resize updates the canvas dimensions.
func (tg *TaskGraph) Resize(width, height int) {
	tg.width = width
	tg.height = height
}

// Init implements tea.Model.
func (tg *TaskGraph) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (tg *TaskGraph) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "left":
			tg.Pan(-5, 0)
		case "j", "down":
			tg.Pan(0, 5)
		case "k", "up":
			tg.Pan(0, -5)
		case "l", "right":
			tg.Pan(5, 0)
		case "+", "=":
			tg.Zoom(1.2)
		case "-", "_":
			tg.Zoom(0.8)
		case "c":
			tg.Center()
		case "r":
			tg.initializeLayout()
			tg.stabilize(50)
		case "enter":
			// Select nearest node (simplified - just toggle selection)
			if tg.selected == "" && len(tg.nodes) > 0 {
				for id := range tg.nodes {
					tg.selected = id
					break
				}
			}
		case "esc":
			tg.ClearSelection()
		}

	case tea.WindowSizeMsg:
		tg.Resize(msg.Width, msg.Height)
	}

	return tg, nil
}

// View implements tea.Model.
func (tg *TaskGraph) View() string {
	if tg.plan == nil || len(tg.nodes) == 0 {
		return graphHelpStyle.Render("No task graph data to display")
	}

	// Create canvas for rendering
	canvas := NewCanvas(tg.width, tg.height-3) // Reserve 3 lines for footer

	// Render edges first (so they appear behind nodes)
	tg.renderEdges(canvas)

	// Render nodes
	tg.renderNodes(canvas)

	// Render selected node details if any
	if tg.selected != "" {
		tg.renderNodeDetails(canvas)
	}

	// Build final view
	var b strings.Builder
	b.WriteString(canvas.Render())
	b.WriteString("\n")
	b.WriteString(tg.renderStats())
	b.WriteString("\n")
	b.WriteString(tg.renderHelp())

	return b.String()
}

// renderNodes renders all nodes to the canvas.
func (tg *TaskGraph) renderNodes(canvas *Canvas) {
	for id, node := range tg.nodes {
		// Transform to screen coordinates
		screenX, screenY := tg.toScreen(node.X, node.Y)

		// Skip if off-screen
		if screenX < 0 || screenX >= canvas.width || screenY < 0 || screenY >= canvas.height {
			continue
		}

		// Choose style based on selection and status
		var style lipgloss.Style
		if id == tg.selected {
			style = graphSelectedNodeStyle
		} else {
			color := tg.getColorForStatus(node.Status)
			style = graphNodeStyle.Copy().Foreground(color)
		}

		// Render node label
		label := tg.formatNodeLabel(node)
		canvas.DrawText(screenX, screenY, style.Render(label))
	}
}

// renderEdges renders all edges to the canvas.
func (tg *TaskGraph) renderEdges(canvas *Canvas) {
	for _, edge := range tg.edges {
		from := tg.nodes[edge.From]
		to := tg.nodes[edge.To]

		if from == nil || to == nil {
			continue
		}

		// Transform to screen coordinates
		x1, y1 := tg.toScreen(from.X, from.Y)
		x2, y2 := tg.toScreen(to.X, to.Y)

		// Skip if both endpoints off-screen
		if (x1 < 0 || x1 >= canvas.width) && (x2 < 0 || x2 >= canvas.width) {
			continue
		}
		if (y1 < 0 || y1 >= canvas.height) && (y2 < 0 || y2 >= canvas.height) {
			continue
		}

		// Draw edge line with arrow
		canvas.DrawLine(x1, y1, x2, y2, graphEdgeStyle)
	}
}

// renderNodeDetails renders details of the selected node.
func (tg *TaskGraph) renderNodeDetails(canvas *Canvas) {
	node := tg.nodes[tg.selected]
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
		fmt.Sprintf("Task: %s", node.TaskID),
		fmt.Sprintf("Status: %s", node.Status.String()),
		fmt.Sprintf("Desc: %s", truncate(node.Description, 25)),
	}

	// Count dependencies
	depCount := 0
	depOnCount := 0
	for _, edge := range tg.edges {
		if edge.To == node.TaskID {
			depCount++
		}
		if edge.From == node.TaskID {
			depOnCount++
		}
	}
	details = append(details, fmt.Sprintf("Deps: %d in, %d out", depCount, depOnCount))

	for i, line := range details {
		if y+i < canvas.height {
			canvas.DrawText(x, y+i, graphHelpStyle.Render(line))
		}
	}
}

// renderStats renders statistics footer.
func (tg *TaskGraph) renderStats() string {
	return graphStatsStyle.Render(fmt.Sprintf(
		"Nodes: %d  Edges: %d  Layout: %d  Zoom: %.1fx  Offset: (%.0f, %.0f)",
		len(tg.nodes), len(tg.edges), tg.layoutSteps, tg.zoom, tg.offsetX, tg.offsetY,
	))
}

// renderHelp renders help text footer.
func (tg *TaskGraph) renderHelp() string {
	if tg.selected != "" {
		return graphHelpStyle.Render("[hjkl] Pan  [+-] Zoom  [esc] Clear  [c] Center  [r] Reset")
	}
	return graphHelpStyle.Render("[hjkl] Pan  [+-] Zoom  [enter] Select  [c] Center  [r] Reset")
}

// toScreen converts graph coordinates to screen coordinates.
func (tg *TaskGraph) toScreen(graphX, graphY float64) (int, int) {
	centerX := float64(tg.width) / 2
	centerY := float64(tg.height-3) / 2 // Account for footer

	screenX := int(graphX*tg.zoom + tg.offsetX + centerX)
	screenY := int(graphY*tg.zoom + tg.offsetY + centerY)

	return screenX, screenY
}

// getColorForStatus returns the color for a task status.
func (tg *TaskGraph) getColorForStatus(status TaskStatus) lipgloss.Color {
	switch status {
	case TaskStatusPending:
		return PendingColor
	case TaskStatusActive:
		return ActiveColor
	case TaskStatusCompleted:
		return CompletedColor
	case TaskStatusFailed:
		return FailedColor
	default:
		return lipgloss.Color("7") // White
	}
}

// formatNodeLabel creates a display label for a node.
func (tg *TaskGraph) formatNodeLabel(node *GraphNode) string {
	// Truncate long labels
	label := node.TaskID
	if len(label) > 10 {
		label = label[:8] + ".."
	}

	// Add brackets
	return "[" + label + "]"
}

// truncate truncates a string to the given length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Canvas represents a 2D character canvas for rendering.
type Canvas struct {
	width  int
	height int
	cells  [][]string // [y][x]
}

// NewCanvas creates a new canvas.
func NewCanvas(width, height int) *Canvas {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

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
// Handles lipgloss styled text by storing it as a single cell.
func (c *Canvas) DrawText(x, y int, text string) {
	if y < 0 || y >= c.height || x < 0 || x >= c.width {
		return
	}

	// For styled text (with ANSI codes), store the whole string in the first cell
	// and skip subsequent cells. For plain text, store character by character.
	if hasANSI(text) {
		c.cells[y][x] = text
	} else {
		for i, ch := range text {
			px := x + i
			if px >= 0 && px < c.width {
				c.cells[y][px] = string(ch)
			}
		}
	}
}

// hasANSI checks if a string contains ANSI escape codes.
func hasANSI(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			return true
		}
	}
	return false
}

// DrawLine draws a line between two points using Bresenham's algorithm.
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
