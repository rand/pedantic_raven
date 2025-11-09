package memorygraph

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

const (
	// DefaultMaxDepth is the default maximum traversal depth.
	DefaultMaxDepth = 3

	// MaxAllowedDepth is the maximum allowed traversal depth to prevent infinite loops.
	MaxAllowedDepth = 5

	// GraphLoadTimeout is the timeout for graph loading operations.
	GraphLoadTimeout = 30 * time.Second
)

// LoadGraph loads a memory graph by traversing from a root memory ID.
// It performs graph traversal up to the specified depth, building a graph
// structure suitable for visualization.
func LoadGraph(client *mnemosyne.Client, rootID string, depth int) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return GraphErrorMsg{
				Err: fmt.Errorf("mnemosyne client not configured"),
			}
		}

		if !client.IsConnected() {
			return GraphErrorMsg{
				Err: mnemosyne.ErrNotConnected,
			}
		}

		if rootID == "" {
			return GraphErrorMsg{
				Err: fmt.Errorf("root memory ID is required"),
			}
		}

		// Validate and clamp depth
		if depth <= 0 {
			depth = DefaultMaxDepth
		}
		if depth > MaxAllowedDepth {
			depth = MaxAllowedDepth
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), GraphLoadTimeout)
		defer cancel()

		// Load graph using GraphTraverse API
		graph, err := loadGraphFromServer(ctx, client, rootID, depth)
		if err != nil {
			return GraphErrorMsg{
				Err: err,
			}
		}

		return GraphLoadedMsg{
			Graph: graph,
		}
	}
}

// loadGraphFromServer performs the actual graph traversal using the mnemosyne API.
func loadGraphFromServer(ctx context.Context, client *mnemosyne.Client, rootID string, depth int) (*Graph, error) {
	// Use GraphTraverse API to get all connected memories
	opts := mnemosyne.GraphTraverseOptions{
		SeedIDs: []string{rootID},
		MaxHops: uint32(depth),
		// Include all link types
		LinkTypes: []pb.LinkType{
			pb.LinkType_LINK_TYPE_REFERENCES,
			pb.LinkType_LINK_TYPE_REFERENCED_BY,
			pb.LinkType_LINK_TYPE_EXTENDS,
			pb.LinkType_LINK_TYPE_BUILDS_UPON,
			pb.LinkType_LINK_TYPE_CONTRADICTS,
			pb.LinkType_LINK_TYPE_IMPLEMENTS,
			pb.LinkType_LINK_TYPE_CLARIFIES,
			pb.LinkType_LINK_TYPE_SUPERSEDES,
		},
		MinLinkStrength: nil, // Include all strengths
	}

	memories, edges, err := client.GraphTraverse(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("graph traverse failed: %w", err)
	}

	// Build graph structure from traversal results
	graph := buildGraphFromTraversal(memories, edges)

	return graph, nil
}

// buildGraphFromTraversal converts traversal results into a Graph structure.
func buildGraphFromTraversal(memories []*pb.MemoryNote, edges []*pb.GraphEdge) *Graph {
	graph := NewGraph()

	// Create nodes for each memory
	for _, memory := range memories {
		node := &Node{
			ID:         memory.Id,
			Memory:     memory,
			X:          0, // Will be set by layout algorithm
			Y:          0,
			VX:         0,
			VY:         0,
			Mass:       1.0,
			IsExpanded: true,
		}
		graph.AddNode(node)
	}

	// Create edges from graph edges
	for _, graphEdge := range edges {
		edge := &Edge{
			SourceID: graphEdge.SourceId,
			TargetID: graphEdge.TargetId,
			LinkType: graphEdge.LinkType,
			Strength: float64(graphEdge.Strength),
		}
		graph.AddEdge(edge)
	}

	return graph
}

// LoadGraphFromMemoryList loads a graph from a list of memories.
// This is useful when you already have a list of memories and want to
// visualize their connections without server traversal.
func LoadGraphFromMemoryList(memories []*pb.MemoryNote) tea.Cmd {
	return func() tea.Msg {
		if len(memories) == 0 {
			return GraphErrorMsg{
				Err: fmt.Errorf("no memories to visualize"),
			}
		}

		graph := buildGraphFromMemoryList(memories)

		return GraphLoadedMsg{
			Graph: graph,
		}
	}
}

// buildGraphFromMemoryList creates a graph from a flat list of memories.
// It uses the links within each memory to build the edge structure.
func buildGraphFromMemoryList(memories []*pb.MemoryNote) *Graph {
	graph := NewGraph()

	// First pass: create all nodes
	nodeMap := make(map[string]*pb.MemoryNote)
	for _, memory := range memories {
		node := &Node{
			ID:         memory.Id,
			Memory:     memory,
			X:          0,
			Y:          0,
			VX:         0,
			VY:         0,
			Mass:       1.0,
			IsExpanded: true,
		}
		graph.AddNode(node)
		nodeMap[memory.Id] = memory
	}

	// Second pass: create edges from links
	for _, memory := range memories {
		for _, link := range memory.Links {
			// Only add edge if target exists in our memory list
			if _, exists := nodeMap[link.TargetId]; exists {
				edge := &Edge{
					SourceID: memory.Id,
					TargetID: link.TargetId,
					LinkType: link.LinkType,
					Strength: float64(link.Strength),
				}
				graph.AddEdge(edge)
			}
		}
	}

	return graph
}

// ExpandNode requests expansion of a specific node.
// This is used to trigger loading of a node's connections.
func ExpandNode(nodeID string) tea.Cmd {
	return func() tea.Msg {
		return ExpandNodeMsg{
			NodeID: nodeID,
		}
	}
}
