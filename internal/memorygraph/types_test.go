package memorygraph

import (
	"errors"
	"testing"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Test NewModel creates a model with correct defaults.
func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.zoom != 1.0 {
		t.Errorf("Expected zoom 1.0, got %f", m.zoom)
	}
	if m.width != 80 {
		t.Errorf("Expected width 80, got %d", m.width)
	}
	if m.height != 20 {
		t.Errorf("Expected height 20, got %d", m.height)
	}
	if m.damping != 0.8 {
		t.Errorf("Expected damping 0.8, got %f", m.damping)
	}
	if !m.focused {
		t.Error("Expected focused to be true")
	}
	if m.layoutSteps != 0 {
		t.Errorf("Expected layoutSteps 0, got %d", m.layoutSteps)
	}
	if m.graph == nil {
		t.Error("Expected graph to be initialized")
	}
}

// Test NewGraph creates an empty graph.
func TestNewGraph(t *testing.T) {
	g := NewGraph()

	if g.Nodes == nil {
		t.Error("Expected Nodes to be initialized")
	}
	if len(g.Nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(g.Nodes))
	}
	if g.Edges == nil {
		t.Error("Expected Edges to be initialized")
	}
	if len(g.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(g.Edges))
	}
}

// Test Graph.AddNode adds a node to the graph.
func TestGraphAddNode(t *testing.T) {
	g := NewGraph()
	node := &Node{
		ID: "test-node",
		Memory: &pb.MemoryNote{
			Id:      "test-node",
			Content: "Test content",
		},
		X: 10.0,
		Y: 20.0,
	}

	g.AddNode(node)

	if len(g.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(g.Nodes))
	}
	if g.Nodes["test-node"] != node {
		t.Error("Expected node to be stored by ID")
	}
}

// Test Graph.AddEdge adds an edge to the graph.
func TestGraphAddEdge(t *testing.T) {
	g := NewGraph()
	edge := &Edge{
		SourceID: "node-1",
		TargetID: "node-2",
		LinkType: pb.LinkType_LINK_TYPE_REFERENCES,
		Strength: 1.0,
	}

	g.AddEdge(edge)

	if len(g.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(g.Edges))
	}
	if g.Edges[0] != edge {
		t.Error("Expected edge to be stored in slice")
	}
}

// Test Graph.GetNode returns the correct node.
func TestGraphGetNode(t *testing.T) {
	g := NewGraph()
	node := &Node{ID: "test-node"}
	g.AddNode(node)

	retrieved := g.GetNode("test-node")
	if retrieved != node {
		t.Error("Expected to retrieve the same node")
	}

	missing := g.GetNode("nonexistent")
	if missing != nil {
		t.Error("Expected nil for nonexistent node")
	}
}

// Test Graph.GetEdgesFrom returns edges originating from a node.
func TestGraphGetEdgesFrom(t *testing.T) {
	g := NewGraph()
	edge1 := &Edge{SourceID: "node-1", TargetID: "node-2"}
	edge2 := &Edge{SourceID: "node-1", TargetID: "node-3"}
	edge3 := &Edge{SourceID: "node-2", TargetID: "node-3"}

	g.AddEdge(edge1)
	g.AddEdge(edge2)
	g.AddEdge(edge3)

	edges := g.GetEdgesFrom("node-1")
	if len(edges) != 2 {
		t.Errorf("Expected 2 edges from node-1, got %d", len(edges))
	}

	edges = g.GetEdgesFrom("node-2")
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge from node-2, got %d", len(edges))
	}

	edges = g.GetEdgesFrom("nonexistent")
	if len(edges) != 0 {
		t.Errorf("Expected 0 edges from nonexistent node, got %d", len(edges))
	}
}

// Test Graph.GetEdgesTo returns edges pointing to a node.
func TestGraphGetEdgesTo(t *testing.T) {
	g := NewGraph()
	edge1 := &Edge{SourceID: "node-1", TargetID: "node-3"}
	edge2 := &Edge{SourceID: "node-2", TargetID: "node-3"}
	edge3 := &Edge{SourceID: "node-1", TargetID: "node-2"}

	g.AddEdge(edge1)
	g.AddEdge(edge2)
	g.AddEdge(edge3)

	edges := g.GetEdgesTo("node-3")
	if len(edges) != 2 {
		t.Errorf("Expected 2 edges to node-3, got %d", len(edges))
	}

	edges = g.GetEdgesTo("node-2")
	if len(edges) != 1 {
		t.Errorf("Expected 1 edge to node-2, got %d", len(edges))
	}

	edges = g.GetEdgesTo("nonexistent")
	if len(edges) != 0 {
		t.Errorf("Expected 0 edges to nonexistent node, got %d", len(edges))
	}
}

// Test Graph.NodeCount returns the correct count.
func TestGraphNodeCount(t *testing.T) {
	g := NewGraph()
	if g.NodeCount() != 0 {
		t.Errorf("Expected 0 nodes, got %d", g.NodeCount())
	}

	g.AddNode(&Node{ID: "node-1"})
	if g.NodeCount() != 1 {
		t.Errorf("Expected 1 node, got %d", g.NodeCount())
	}

	g.AddNode(&Node{ID: "node-2"})
	if g.NodeCount() != 2 {
		t.Errorf("Expected 2 nodes, got %d", g.NodeCount())
	}
}

// Test Graph.EdgeCount returns the correct count.
func TestGraphEdgeCount(t *testing.T) {
	g := NewGraph()
	if g.EdgeCount() != 0 {
		t.Errorf("Expected 0 edges, got %d", g.EdgeCount())
	}

	g.AddEdge(&Edge{SourceID: "node-1", TargetID: "node-2"})
	if g.EdgeCount() != 1 {
		t.Errorf("Expected 1 edge, got %d", g.EdgeCount())
	}

	g.AddEdge(&Edge{SourceID: "node-2", TargetID: "node-3"})
	if g.EdgeCount() != 2 {
		t.Errorf("Expected 2 edges, got %d", g.EdgeCount())
	}
}

// Test Model.SetSize updates width and height.
func TestModelSetSize(t *testing.T) {
	m := NewModel()
	m.SetSize(120, 40)

	if m.width != 120 {
		t.Errorf("Expected width 120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("Expected height 40, got %d", m.height)
	}
}

// Test Model.SetFocus updates focus state.
func TestModelSetFocus(t *testing.T) {
	m := NewModel()
	m.SetFocus(false)
	if m.focused {
		t.Error("Expected focused to be false")
	}

	m.SetFocus(true)
	if !m.focused {
		t.Error("Expected focused to be true")
	}
}

// Test Model.IsFocused returns focus state.
func TestModelIsFocused(t *testing.T) {
	m := NewModel()
	if !m.IsFocused() {
		t.Error("Expected IsFocused to return true")
	}

	m.SetFocus(false)
	if m.IsFocused() {
		t.Error("Expected IsFocused to return false")
	}
}

// Test Model.SetGraph updates the graph.
func TestModelSetGraph(t *testing.T) {
	m := NewModel()
	m.layoutSteps = 10
	m.err = errors.New("test error")

	g := NewGraph()
	g.AddNode(&Node{ID: "test"})

	m.SetGraph(g)

	if m.graph != g {
		t.Error("Expected graph to be updated")
	}
	if m.layoutSteps != 0 {
		t.Errorf("Expected layoutSteps to be reset to 0, got %d", m.layoutSteps)
	}
	if m.err != nil {
		t.Error("Expected error to be cleared")
	}
}

// Test Model.Graph returns the graph.
func TestModelGraph(t *testing.T) {
	m := NewModel()
	g := m.Graph()

	if g != m.graph {
		t.Error("Expected Graph() to return the internal graph")
	}
}

// Test Model.SetError updates error state.
func TestModelSetError(t *testing.T) {
	m := NewModel()
	err := errors.New("test error")
	m.SetError(err)

	if m.err != err {
		t.Error("Expected error to be set")
	}
}

// Test Model.Error returns error state.
func TestModelError(t *testing.T) {
	m := NewModel()
	if m.Error() != nil {
		t.Error("Expected Error() to return nil initially")
	}

	err := errors.New("test error")
	m.SetError(err)
	if m.Error() != err {
		t.Error("Expected Error() to return the set error")
	}
}

// Test Model.SetCenterNode updates center node ID.
func TestModelSetCenterNode(t *testing.T) {
	m := NewModel()
	m.SetCenterNode("test-node")

	if m.centerNodeID != "test-node" {
		t.Errorf("Expected centerNodeID 'test-node', got %s", m.centerNodeID)
	}
}

// Test Model.CenterNodeID returns center node ID.
func TestModelCenterNodeID(t *testing.T) {
	m := NewModel()
	m.SetCenterNode("test-node")

	if m.CenterNodeID() != "test-node" {
		t.Errorf("Expected CenterNodeID 'test-node', got %s", m.CenterNodeID())
	}
}

// Test Model.SelectNode updates selected node ID.
func TestModelSelectNode(t *testing.T) {
	m := NewModel()
	m.SelectNode("test-node")

	if m.selectedNodeID != "test-node" {
		t.Errorf("Expected selectedNodeID 'test-node', got %s", m.selectedNodeID)
	}
}

// Test Model.SelectedNodeID returns selected node ID.
func TestModelSelectedNodeID(t *testing.T) {
	m := NewModel()
	m.SelectNode("test-node")

	if m.SelectedNodeID() != "test-node" {
		t.Errorf("Expected SelectedNodeID 'test-node', got %s", m.SelectedNodeID())
	}
}

// Test Model.SelectedNode returns the selected node.
func TestModelSelectedNode(t *testing.T) {
	m := NewModel()
	node := &Node{ID: "test-node"}
	m.graph.AddNode(node)
	m.SelectNode("test-node")

	selected := m.SelectedNode()
	if selected != node {
		t.Error("Expected SelectedNode to return the selected node")
	}

	m.SelectNode("")
	if m.SelectedNode() != nil {
		t.Error("Expected SelectedNode to return nil when no node selected")
	}

	m.SelectNode("nonexistent")
	if m.SelectedNode() != nil {
		t.Error("Expected SelectedNode to return nil for nonexistent node")
	}
}

// Test Model.SetClient updates the client.
func TestModelSetClient(t *testing.T) {
	m := NewModel()
	client := "mock-client"
	m.SetClient(client)

	if m.client != client {
		t.Error("Expected client to be set")
	}
}

// Test Model.Client returns the client.
func TestModelClient(t *testing.T) {
	m := NewModel()
	client := "mock-client"
	m.SetClient(client)

	if m.Client() != client {
		t.Error("Expected Client() to return the set client")
	}
}

// Test Model.Init returns nil.
func TestModelInit(t *testing.T) {
	m := NewModel()
	cmd := m.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil")
	}
}
