package analyze

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// Model represents the triple graph visualization state.
type Model struct {
	// Graph data
	graph *TripleGraph

	// Display state
	selectedNodeID string // Currently selected node ID
	highlightedIDs map[string]bool // Highlighted node IDs

	// Viewport state
	offsetX float64 // Pan offset X
	offsetY float64 // Pan offset Y
	zoom    float64 // Zoom level (1.0 = normal)
	width   int     // Component width
	height  int     // Component height

	// Layout state
	layoutSteps int     // Number of layout iterations performed
	damping     float64 // Velocity damping factor

	// Filter state
	filter Filter

	// UI state
	focused bool
	err     error

	// Analysis source
	analysis *semantic.Analysis
}

// Messages for the triple graph component.
type (
	// GraphLoadedMsg is sent when a graph is loaded.
	GraphLoadedMsg struct {
		Graph *TripleGraph
	}

	// GraphErrorMsg is sent when graph loading fails.
	GraphErrorMsg struct {
		Err error
	}

	// NodeSelectedMsg is sent when a node is selected.
	NodeSelectedMsg struct {
		NodeID string
	}

	// FilterUpdatedMsg is sent when filters change.
	FilterUpdatedMsg struct {
		Filter Filter
	}

	// LayoutStepMsg triggers one layout iteration.
	LayoutStepMsg struct{}
)

// NewModel creates a new triple graph visualization model.
func NewModel() Model {
	return Model{
		graph:          NewTripleGraph(),
		zoom:           1.0,
		width:          80,
		height:         20,
		layoutSteps:    0,
		damping:        0.8,
		focused:        true,
		highlightedIDs: make(map[string]bool),
		filter:         Filter{},
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetSize sets the component size.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocus sets the focus state.
func (m *Model) SetFocus(focused bool) {
	m.focused = focused
}

// IsFocused returns whether the component is focused.
func (m Model) IsFocused() bool {
	return m.focused
}

// SetGraph sets the graph to display.
func (m *Model) SetGraph(graph *TripleGraph) {
	m.graph = graph
	m.layoutSteps = 0
	m.err = nil
	m.selectedNodeID = ""
	m.highlightedIDs = make(map[string]bool)
}

// Graph returns the current graph.
func (m Model) Graph() *TripleGraph {
	return m.graph
}

// SetAnalysis sets the source analysis and builds the graph.
func (m *Model) SetAnalysis(analysis *semantic.Analysis) {
	m.analysis = analysis
	m.graph = BuildFromAnalysis(analysis)
	m.graph.InitializeLayout()
	m.layoutSteps = 0
	m.selectedNodeID = ""
	m.highlightedIDs = make(map[string]bool)
}

// Analysis returns the source analysis.
func (m Model) Analysis() *semantic.Analysis {
	return m.analysis
}

// SetError sets the error state.
func (m *Model) SetError(err error) {
	m.err = err
}

// Error returns the current error, if any.
func (m Model) Error() error {
	return m.err
}

// SelectNode selects a node and highlights connected nodes.
func (m *Model) SelectNode(nodeID string) {
	m.selectedNodeID = nodeID
	m.highlightedIDs = make(map[string]bool)

	if nodeID != "" && m.graph != nil {
		// Highlight the selected node
		m.highlightedIDs[nodeID] = true

		// Highlight all connected nodes
		for _, edge := range m.graph.GetEdgesFrom(nodeID) {
			m.highlightedIDs[edge.TargetID] = true
		}
		for _, edge := range m.graph.GetEdgesTo(nodeID) {
			m.highlightedIDs[edge.SourceID] = true
		}
	}
}

// SelectedNodeID returns the selected node ID.
func (m Model) SelectedNodeID() string {
	return m.selectedNodeID
}

// SelectedNode returns the selected node, or nil.
func (m Model) SelectedNode() *TripleNode {
	if m.selectedNodeID == "" || m.graph == nil {
		return nil
	}
	return m.graph.GetNode(m.selectedNodeID)
}

// IsHighlighted returns whether a node is highlighted.
func (m Model) IsHighlighted(nodeID string) bool {
	return m.highlightedIDs[nodeID]
}

// Pan moves the viewport.
func (m *Model) Pan(dx, dy float64) {
	m.offsetX += dx
	m.offsetY += dy
}

// Zoom adjusts the zoom level.
func (m *Model) Zoom(delta float64) {
	m.zoom += delta
	// Clamp zoom to reasonable range
	if m.zoom < 0.1 {
		m.zoom = 0.1
	}
	if m.zoom > 5.0 {
		m.zoom = 5.0
	}
}

// ResetView resets pan and zoom to defaults.
func (m *Model) ResetView() {
	m.offsetX = 0
	m.offsetY = 0
	m.zoom = 1.0
}

// ApplyFilter applies a filter to the graph.
func (m *Model) ApplyFilter(filter Filter) {
	m.filter = filter
	if m.analysis != nil {
		// Rebuild graph from analysis with filter
		fullGraph := BuildFromAnalysis(m.analysis)
		m.graph = fullGraph.ApplyFilter(filter)
		m.graph.InitializeLayout()
		m.layoutSteps = 0
	}
}

// Filter returns the current filter.
func (m Model) Filter() Filter {
	return m.filter
}

// StabilizeLayout runs layout iterations until stable or max iterations reached.
func (m *Model) StabilizeLayout(maxIterations int) {
	for i := 0; i < maxIterations && m.layoutSteps < 100; i++ {
		m.graph.ApplyForceIteration(m.damping)
		m.layoutSteps++
	}
}

// LayoutSteps returns the number of layout iterations performed.
func (m Model) LayoutSteps() int {
	return m.layoutSteps
}

// GetNodeAtPosition finds a node near the given screen coordinates.
// Returns the node ID or empty string if none found.
func (m Model) GetNodeAtPosition(screenX, screenY int) string {
	if m.graph == nil {
		return ""
	}

	// Convert screen coordinates to graph coordinates
	graphX := (float64(screenX) - float64(m.width)/2 - m.offsetX) / m.zoom
	graphY := (float64(screenY) - float64(m.height)/2 - m.offsetY) / m.zoom

	// Find closest node within threshold
	const clickThreshold = 2.0
	closestID := ""
	closestDist := clickThreshold

	for id, node := range m.graph.Nodes {
		dx := node.X - graphX
		dy := node.Y - graphY
		dist := dx*dx + dy*dy // Squared distance (faster than sqrt)

		if dist < closestDist*closestDist {
			closestDist = dist
			closestID = id
		}
	}

	return closestID
}

// CenterOnNode centers the viewport on a specific node.
func (m *Model) CenterOnNode(nodeID string) {
	if m.graph == nil {
		return
	}

	node := m.graph.GetNode(nodeID)
	if node == nil {
		return
	}

	// Set offset to center this node
	m.offsetX = -node.X * m.zoom
	m.offsetY = -node.Y * m.zoom
}

// AutoCenter centers the viewport on the graph's center of mass.
func (m *Model) AutoCenter() {
	if m.graph == nil || len(m.graph.Nodes) == 0 {
		return
	}

	minX, maxX, minY, maxY := m.graph.GetBounds()

	// Calculate center
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2

	// Set offset to center the graph
	m.offsetX = -centerX * m.zoom
	m.offsetY = -centerY * m.zoom
}

// GetStats returns statistics about the current graph.
type GraphStats struct {
	Nodes          int
	Edges          int
	LayoutSteps    int
	AvgImportance  float64
	AvgConfidence  float64
	FilteredNodes  int
	FilteredEdges  int
}

// GetStats computes current graph statistics.
func (m Model) GetStats() GraphStats {
	stats := GraphStats{
		Nodes:       m.graph.NodeCount(),
		Edges:       m.graph.EdgeCount(),
		LayoutSteps: m.layoutSteps,
	}

	if stats.Nodes > 0 {
		totalImportance := 0
		for _, node := range m.graph.Nodes {
			totalImportance += node.Importance
		}
		stats.AvgImportance = float64(totalImportance) / float64(stats.Nodes)
	}

	if stats.Edges > 0 {
		totalConfidence := 0.0
		for _, edge := range m.graph.Edges {
			totalConfidence += edge.Confidence
		}
		stats.AvgConfidence = totalConfidence / float64(stats.Edges)
	}

	// If we have source analysis, compute filtered counts
	if m.analysis != nil {
		stats.FilteredNodes = len(m.analysis.Entities) - stats.Nodes
		stats.FilteredEdges = len(m.analysis.Relationships) - stats.Edges
	}

	return stats
}
