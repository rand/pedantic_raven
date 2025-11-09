package memorygraph

import (
	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Node represents a memory node in the graph.
type Node struct {
	ID         string
	Memory     *pb.MemoryNote
	X          float64 // Screen X position
	Y          float64 // Screen Y position
	VX         float64 // Velocity X (for force-directed layout)
	VY         float64 // Velocity Y
	Mass       float64 // Node mass (affects layout)
	IsExpanded bool    // Whether node children are shown
}

// Edge represents a link between two nodes.
type Edge struct {
	SourceID string
	TargetID string
	LinkType pb.LinkType
	Strength float64
}

// Graph represents the memory graph.
type Graph struct {
	Nodes map[string]*Node // ID -> Node
	Edges []*Edge
}

// Model represents the graph visualization component state.
type Model struct {
	// Graph data
	graph *Graph

	// Display state
	centerNodeID  string // ID of the central/focused node
	selectedNodeID string // ID of the currently selected node

	// Viewport state
	offsetX      float64 // Pan offset X
	offsetY      float64 // Pan offset Y
	zoom         float64 // Zoom level (1.0 = normal)
	width        int
	height       int

	// Layout state
	layoutSteps  int     // Number of layout iterations
	damping      float64 // Velocity damping

	// UI state
	focused bool
	err     error

	// Client integration
	client interface{} // Can hold *mnemosyne.Client
}

// Messages for the graph component.
type (
	// GraphLoadedMsg is sent when a graph is loaded.
	GraphLoadedMsg struct {
		Graph *Graph
	}

	// GraphErrorMsg is sent when graph loading fails.
	GraphErrorMsg struct {
		Err error
	}

	// NodeSelectedMsg is sent when a node is selected.
	NodeSelectedMsg struct {
		NodeID string
	}

	// ExpandNodeMsg is sent when a node should be expanded.
	ExpandNodeMsg struct {
		NodeID string
	}

	// CollapseNodeMsg is sent when a node should be collapsed.
	CollapseNodeMsg struct {
		NodeID string
	}
)

// NewModel creates a new graph visualization model.
func NewModel() Model {
	return Model{
		graph:        NewGraph(),
		zoom:         1.0,
		width:        80,
		height:       20,
		layoutSteps:  0,
		damping:      0.8,
		focused:      true,
	}
}

// NewGraph creates an empty graph.
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make([]*Edge, 0),
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
func (m *Model) SetGraph(graph *Graph) {
	m.graph = graph
	m.layoutSteps = 0
	m.err = nil
}

// Graph returns the current graph.
func (m Model) Graph() *Graph {
	return m.graph
}

// SetError sets the error state.
func (m *Model) SetError(err error) {
	m.err = err
}

// Error returns the current error, if any.
func (m Model) Error() error {
	return m.err
}

// SetCenterNode sets the central node ID.
func (m *Model) SetCenterNode(nodeID string) {
	m.centerNodeID = nodeID
}

// CenterNodeID returns the central node ID.
func (m Model) CenterNodeID() string {
	return m.centerNodeID
}

// SelectNode selects a node.
func (m *Model) SelectNode(nodeID string) {
	m.selectedNodeID = nodeID
}

// SelectedNodeID returns the selected node ID.
func (m Model) SelectedNodeID() string {
	return m.selectedNodeID
}

// SelectedNode returns the selected node, or nil.
func (m Model) SelectedNode() *Node {
	if m.selectedNodeID == "" {
		return nil
	}
	return m.graph.Nodes[m.selectedNodeID]
}

// SetClient sets the mnemosyne client for this model.
func (m *Model) SetClient(client interface{}) {
	m.client = client
}

// Client returns the mnemosyne client, if set.
func (m Model) Client() interface{} {
	return m.client
}

// Graph helper methods

// AddNode adds a node to the graph.
func (g *Graph) AddNode(node *Node) {
	g.Nodes[node.ID] = node
}

// AddEdge adds an edge to the graph.
func (g *Graph) AddEdge(edge *Edge) {
	g.Edges = append(g.Edges, edge)
}

// GetNode returns a node by ID, or nil if not found.
func (g *Graph) GetNode(id string) *Node {
	return g.Nodes[id]
}

// GetEdgesFrom returns all edges originating from a node.
func (g *Graph) GetEdgesFrom(nodeID string) []*Edge {
	edges := make([]*Edge, 0)
	for _, edge := range g.Edges {
		if edge.SourceID == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetEdgesTo returns all edges pointing to a node.
func (g *Graph) GetEdgesTo(nodeID string) []*Edge {
	edges := make([]*Edge, 0)
	for _, edge := range g.Edges {
		if edge.TargetID == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// NodeCount returns the number of nodes in the graph.
func (g *Graph) NodeCount() int {
	return len(g.Nodes)
}

// EdgeCount returns the number of edges in the graph.
func (g *Graph) EdgeCount() int {
	return len(g.Edges)
}
