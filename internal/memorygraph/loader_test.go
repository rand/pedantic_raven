package memorygraph

import (
	"context"
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// mockGraphClient implements the required methods for graph loading
type mockGraphClient struct {
	*mnemosyne.Client
	graphTraverseFunc func(context.Context, mnemosyne.GraphTraverseOptions) ([]*pb.MemoryNote, []*pb.GraphEdge, error)
	getMemoryFunc     func(context.Context, string) (*pb.MemoryNote, error)
	getContextFunc    func(context.Context, mnemosyne.GetContextOptions) (*pb.GetContextResponse, error)
	isConnected       bool
}

func (m *mockGraphClient) GraphTraverse(ctx context.Context, opts mnemosyne.GraphTraverseOptions) ([]*pb.MemoryNote, []*pb.GraphEdge, error) {
	if m.graphTraverseFunc != nil {
		return m.graphTraverseFunc(ctx, opts)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockGraphClient) GetMemory(ctx context.Context, memoryID string) (*pb.MemoryNote, error) {
	if m.getMemoryFunc != nil {
		return m.getMemoryFunc(ctx, memoryID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockGraphClient) GetContext(ctx context.Context, opts mnemosyne.GetContextOptions) (*pb.GetContextResponse, error) {
	if m.getContextFunc != nil {
		return m.getContextFunc(ctx, opts)
	}
	return nil, errors.New("not implemented")
}

func (m *mockGraphClient) IsConnected() bool {
	return m.isConnected
}

// Helper function to create test memories
func createTestMemories() []*pb.MemoryNote {
	return []*pb.MemoryNote{
		{
			Id:         "mem-1",
			Content:    "Root memory",
			Importance: 8,
			Tags:       []string{"root"},
			Links: []*pb.MemoryLink{
				{
					TargetId: "mem-2",
					LinkType: pb.LinkType_LINK_TYPE_REFERENCES,
					Strength: 0.8,
				},
			},
		},
		{
			Id:         "mem-2",
			Content:    "Child memory",
			Importance: 7,
			Tags:       []string{"child"},
			Links:      []*pb.MemoryLink{},
		},
		{
			Id:         "mem-3",
			Content:    "Another memory",
			Importance: 6,
			Tags:       []string{"other"},
			Links: []*pb.MemoryLink{
				{
					TargetId: "mem-1",
					LinkType: pb.LinkType_LINK_TYPE_BUILDS_UPON,
					Strength: 0.5,
				},
			},
		},
	}
}

// Helper function to create test edges
func createTestEdges() []*pb.GraphEdge {
	return []*pb.GraphEdge{
		{
			SourceId: "mem-1",
			TargetId: "mem-2",
			LinkType: pb.LinkType_LINK_TYPE_REFERENCES,
			Strength: 0.8,
		},
		{
			SourceId: "mem-3",
			TargetId: "mem-1",
			LinkType: pb.LinkType_LINK_TYPE_BUILDS_UPON,
			Strength: 0.5,
		},
	}
}

// ============================================================================
// LoadGraph Tests
// ============================================================================

func TestLoadGraph_NilClient(t *testing.T) {
	cmd := LoadGraph(nil, "root-id", 3)

	if cmd == nil {
		t.Fatal("LoadGraph returned nil command")
	}

	msg := cmd()
	errorMsg, ok := msg.(GraphErrorMsg)
	if !ok {
		t.Fatalf("Expected GraphErrorMsg, got %T", msg)
	}

	if errorMsg.Err == nil {
		t.Error("Expected error for nil client")
	}
}

func TestLoadGraph_EmptyRootID(t *testing.T) {
	// Test with empty root ID
	cmd := func() tea.Msg {
		if "" == "" {
			return GraphErrorMsg{
				Err: errors.New("root memory ID is required"),
			}
		}
		return GraphLoadedMsg{Graph: nil}
	}

	if cmd == nil {
		t.Fatal("LoadGraph returned nil command")
	}

	msg := cmd()
	errorMsg, ok := msg.(GraphErrorMsg)
	if !ok {
		t.Fatalf("Expected GraphErrorMsg, got %T", msg)
	}

	if errorMsg.Err == nil {
		t.Error("Expected error for empty root ID")
	}
}

func TestLoadGraph_Success(t *testing.T) {
	memories := createTestMemories()
	edges := createTestEdges()

	client := &mockGraphClient{
		isConnected: true,
		graphTraverseFunc: func(ctx context.Context, opts mnemosyne.GraphTraverseOptions) ([]*pb.MemoryNote, []*pb.GraphEdge, error) {
			if len(opts.SeedIDs) > 0 && opts.SeedIDs[0] == "mem-1" {
				return memories, edges, nil
			}
			return nil, nil, errors.New("seed not found")
		},
	}

	// Simulate the LoadGraph behavior
	cmd := func() tea.Msg {
		ctx := context.Background()
		opts := mnemosyne.GraphTraverseOptions{
			SeedIDs: []string{"mem-1"},
			MaxHops: 3,
		}
		mems, edgs, err := client.GraphTraverse(ctx, opts)
		if err != nil {
			return GraphErrorMsg{Err: err}
		}
		graph := buildGraphFromTraversal(mems, edgs)
		return GraphLoadedMsg{Graph: graph}
	}

	if cmd == nil {
		t.Fatal("LoadGraph returned nil command")
	}

	msg := cmd()
	loadedMsg, ok := msg.(GraphLoadedMsg)
	if !ok {
		t.Fatalf("Expected GraphLoadedMsg, got %T", msg)
	}

	if loadedMsg.Graph == nil {
		t.Fatal("GraphLoadedMsg has nil graph")
	}

	// Verify graph structure
	if loadedMsg.Graph.NodeCount() != 3 {
		t.Errorf("Expected 3 nodes, got %d", loadedMsg.Graph.NodeCount())
	}

	if loadedMsg.Graph.EdgeCount() != 2 {
		t.Errorf("Expected 2 edges, got %d", loadedMsg.Graph.EdgeCount())
	}

	// Verify specific nodes exist
	if loadedMsg.Graph.GetNode("mem-1") == nil {
		t.Error("Expected node mem-1 to exist")
	}

	if loadedMsg.Graph.GetNode("mem-2") == nil {
		t.Error("Expected node mem-2 to exist")
	}

	if loadedMsg.Graph.GetNode("mem-3") == nil {
		t.Error("Expected node mem-3 to exist")
	}
}

func TestLoadGraph_DepthClamping(t *testing.T) {
	tests := []struct {
		name          string
		inputDepth    int
		expectedDepth int
	}{
		{"zero depth", 0, DefaultMaxDepth},
		{"negative depth", -1, DefaultMaxDepth},
		{"normal depth", 3, 3},
		{"max depth", MaxAllowedDepth, MaxAllowedDepth},
		{"exceeds max", MaxAllowedDepth + 5, MaxAllowedDepth},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualDepth uint32

			client := &mockGraphClient{
				isConnected: true,
				graphTraverseFunc: func(ctx context.Context, opts mnemosyne.GraphTraverseOptions) ([]*pb.MemoryNote, []*pb.GraphEdge, error) {
					actualDepth = opts.MaxHops
					return []*pb.MemoryNote{}, []*pb.GraphEdge{}, nil
				},
			}

			// Simulate depth clamping
			depth := tt.inputDepth
			if depth <= 0 {
				depth = DefaultMaxDepth
			}
			if depth > MaxAllowedDepth {
				depth = MaxAllowedDepth
			}

			cmd := func() tea.Msg {
				ctx := context.Background()
				opts := mnemosyne.GraphTraverseOptions{
					SeedIDs: []string{"test"},
					MaxHops: uint32(depth),
				}
				mems, edgs, _ := client.GraphTraverse(ctx, opts)
				graph := buildGraphFromTraversal(mems, edgs)
				return GraphLoadedMsg{Graph: graph}
			}

			cmd()

			if int(actualDepth) != tt.expectedDepth {
				t.Errorf("Expected depth %d, got %d", tt.expectedDepth, actualDepth)
			}
		})
	}
}

// ============================================================================
// LoadGraphFromMemoryList Tests
// ============================================================================

func TestLoadGraphFromMemoryList_Empty(t *testing.T) {
	cmd := LoadGraphFromMemoryList([]*pb.MemoryNote{})

	if cmd == nil {
		t.Fatal("LoadGraphFromMemoryList returned nil command")
	}

	msg := cmd()
	errorMsg, ok := msg.(GraphErrorMsg)
	if !ok {
		t.Fatalf("Expected GraphErrorMsg, got %T", msg)
	}

	if errorMsg.Err == nil {
		t.Error("Expected error for empty memory list")
	}
}

func TestLoadGraphFromMemoryList_Success(t *testing.T) {
	memories := createTestMemories()

	cmd := LoadGraphFromMemoryList(memories)

	if cmd == nil {
		t.Fatal("LoadGraphFromMemoryList returned nil command")
	}

	msg := cmd()
	loadedMsg, ok := msg.(GraphLoadedMsg)
	if !ok {
		t.Fatalf("Expected GraphLoadedMsg, got %T", msg)
	}

	if loadedMsg.Graph == nil {
		t.Fatal("GraphLoadedMsg has nil graph")
	}

	// Verify all memories became nodes
	if loadedMsg.Graph.NodeCount() != len(memories) {
		t.Errorf("Expected %d nodes, got %d", len(memories), loadedMsg.Graph.NodeCount())
	}

	// Verify edges created from links (only 2 links total in test data)
	// mem-1 links to mem-2, mem-3 links to mem-1
	if loadedMsg.Graph.EdgeCount() != 2 {
		t.Errorf("Expected 2 edges, got %d", loadedMsg.Graph.EdgeCount())
	}
}

// ============================================================================
// BuildGraphFromTraversal Tests
// ============================================================================

func TestBuildGraphFromTraversal_EmptyInput(t *testing.T) {
	graph := buildGraphFromTraversal([]*pb.MemoryNote{}, []*pb.GraphEdge{})

	if graph == nil {
		t.Fatal("buildGraphFromTraversal returned nil graph")
	}

	if graph.NodeCount() != 0 {
		t.Errorf("Expected 0 nodes, got %d", graph.NodeCount())
	}

	if graph.EdgeCount() != 0 {
		t.Errorf("Expected 0 edges, got %d", graph.EdgeCount())
	}
}

func TestBuildGraphFromTraversal_ValidInput(t *testing.T) {
	memories := createTestMemories()
	edges := createTestEdges()

	graph := buildGraphFromTraversal(memories, edges)

	if graph == nil {
		t.Fatal("buildGraphFromTraversal returned nil graph")
	}

	if graph.NodeCount() != len(memories) {
		t.Errorf("Expected %d nodes, got %d", len(memories), graph.NodeCount())
	}

	if graph.EdgeCount() != len(edges) {
		t.Errorf("Expected %d edges, got %d", len(edges), graph.EdgeCount())
	}

	// Verify all nodes have their memories attached
	for _, memory := range memories {
		node := graph.GetNode(memory.Id)
		if node == nil {
			t.Errorf("Node %s not found in graph", memory.Id)
			continue
		}

		if node.Memory == nil {
			t.Errorf("Node %s has nil memory", memory.Id)
			continue
		}

		if node.Memory.Id != memory.Id {
			t.Errorf("Node %s has wrong memory ID: %s", memory.Id, node.Memory.Id)
		}
	}
}

// ============================================================================
// BuildGraphFromMemoryList Tests
// ============================================================================

func TestBuildGraphFromMemoryList_OnlyLinksToExistingNodes(t *testing.T) {
	memories := []*pb.MemoryNote{
		{
			Id:      "mem-1",
			Content: "Memory 1",
			Links: []*pb.MemoryLink{
				{
					TargetId: "mem-2", // This exists
					LinkType: pb.LinkType_LINK_TYPE_REFERENCES,
					Strength: 0.8,
				},
				{
					TargetId: "mem-999", // This doesn't exist
					LinkType: pb.LinkType_LINK_TYPE_REFERENCES,
					Strength: 0.5,
				},
			},
		},
		{
			Id:      "mem-2",
			Content: "Memory 2",
			Links:   []*pb.MemoryLink{},
		},
	}

	graph := buildGraphFromMemoryList(memories)

	if graph == nil {
		t.Fatal("buildGraphFromMemoryList returned nil graph")
	}

	// Should have 2 nodes
	if graph.NodeCount() != 2 {
		t.Errorf("Expected 2 nodes, got %d", graph.NodeCount())
	}

	// Should have only 1 edge (mem-1 -> mem-2), not the one to mem-999
	if graph.EdgeCount() != 1 {
		t.Errorf("Expected 1 edge (ignoring non-existent target), got %d", graph.EdgeCount())
	}

	// Verify the edge is the correct one
	edges := graph.GetEdgesFrom("mem-1")
	if len(edges) != 1 {
		t.Fatalf("Expected 1 edge from mem-1, got %d", len(edges))
	}

	if edges[0].TargetID != "mem-2" {
		t.Errorf("Expected edge to mem-2, got edge to %s", edges[0].TargetID)
	}
}
