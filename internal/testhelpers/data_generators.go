package testhelpers

import (
	"fmt"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
	"github.com/rand/pedantic-raven/internal/memorygraph"
)

// GenerateTestMemories creates a set of test memories with variety
func GenerateTestMemories(count int) []*pb.MemoryNote {
	memories := make([]*pb.MemoryNote, count)

	tags := [][]string{
		{"work", "urgent"},
		{"personal", "health"},
		{"learning", "golang"},
		{"project", "pedantic-raven"},
		{"idea", "future"},
	}

	for i := 0; i < count; i++ {
		importance := uint32((i % 10) + 1) // 1-10
		tagSet := tags[i%len(tags)]

		memories[i] = &pb.MemoryNote{
			Id:          fmt.Sprintf("mem-%d", i+1),
			Content:     fmt.Sprintf("Test memory content %d with various details", i+1),
			Importance:  importance,
			Tags:        tagSet,
			Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Global{Global: &pb.GlobalNamespace{}}},
			CreatedAt:   uint64(time.Now().Add(-time.Duration(i) * time.Hour).Unix()),
			UpdatedAt:   uint64(time.Now().Add(-time.Duration(i/2) * time.Hour).Unix()),
			AccessCount: uint64(i % 20),
		}
	}

	return memories
}

// GenerateTestGraph creates a test graph with specified nodes and edges
func GenerateTestGraph(nodeCount, edgeCount int) *memorygraph.Graph {
	graph := &memorygraph.Graph{
		Nodes: make(map[string]*memorygraph.Node),
		Edges: make([]*memorygraph.Edge, 0),
	}

	// Create nodes
	for i := 0; i < nodeCount; i++ {
		nodeID := fmt.Sprintf("node-%d", i+1)
		node := &memorygraph.Node{
			ID:         nodeID,
			X:          float64(i * 10),
			Y:          float64(i * 10),
			VX:         0,
			VY:         0,
			Mass:       1.0,
			IsExpanded: false,
			Memory: &pb.MemoryNote{
				Id:         nodeID,
				Content:    fmt.Sprintf("Memory for node %d", i+1),
				Importance: uint32((i % 10) + 1),
				Tags:       []string{"test"},
				Namespace:  &pb.Namespace{Namespace: &pb.Namespace_Global{Global: &pb.GlobalNamespace{}}},
				CreatedAt:  uint64(time.Now().Unix()),
				UpdatedAt:  uint64(time.Now().Unix()),
			},
		}
		graph.Nodes[nodeID] = node
	}

	// Create edges (ensuring we don't exceed the requested count)
	edgesCreated := 0
	for i := 0; i < nodeCount-1 && edgesCreated < edgeCount; i++ {
		sourceID := fmt.Sprintf("node-%d", i+1)
		targetID := fmt.Sprintf("node-%d", i+2)

		edge := &memorygraph.Edge{
			SourceID: sourceID,
			TargetID: targetID,
			LinkType: pb.LinkType_LINK_TYPE_REFERENCES,
			Strength: 0.8,
		}
		graph.Edges = append(graph.Edges, edge)
		edgesCreated++
	}

	// Add additional edges if needed (random connections)
	for i := nodeCount - 1; edgesCreated < edgeCount && i < nodeCount*2; i++ {
		source := i % nodeCount
		target := (i + 3) % nodeCount
		if source != target {
			sourceID := fmt.Sprintf("node-%d", source+1)
			targetID := fmt.Sprintf("node-%d", target+1)

			edge := &memorygraph.Edge{
				SourceID: sourceID,
				TargetID: targetID,
				LinkType: pb.LinkType_LINK_TYPE_REFERENCES,
				Strength: 0.5,
			}
			graph.Edges = append(graph.Edges, edge)
			edgesCreated++
		}
	}

	return graph
}

// GenerateLinkedMemories creates memories with links between them
func GenerateLinkedMemories(count int, linksPerMemory int) ([]*pb.MemoryNote, map[string][]*pb.MemoryLink) {
	memories := GenerateTestMemories(count)
	links := make(map[string][]*pb.MemoryLink)

	for i := 0; i < count; i++ {
		memoryID := memories[i].Id
		memoryLinks := make([]*pb.MemoryLink, 0)

		for j := 0; j < linksPerMemory && i+j+1 < count; j++ {
			targetID := memories[i+j+1].Id
			link := &pb.MemoryLink{
				TargetId:  targetID,
				LinkType:  pb.LinkType_LINK_TYPE_REFERENCES,
				Strength:  0.7,
				Reason:    "test link",
				CreatedAt: uint64(time.Now().Unix()),
				UserCreated: true,
			}
			memoryLinks = append(memoryLinks, link)
		}

		if len(memoryLinks) > 0 {
			links[memoryID] = memoryLinks
		}
	}

	return memories, links
}

// GenerateTestMemoryWithLinks creates a single memory with specified links
func GenerateTestMemoryWithLinks(id string, linkCount int) (*pb.MemoryNote, []*pb.MemoryLink) {
	memory := &pb.MemoryNote{
		Id:          id,
		Content:     fmt.Sprintf("Test memory %s with %d links", id, linkCount),
		Importance:  5,
		Tags:        []string{"test"},
		Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Global{Global: &pb.GlobalNamespace{}}},
		CreatedAt:   uint64(time.Now().Unix()),
		UpdatedAt:   uint64(time.Now().Unix()),
		AccessCount: 0,
	}

	links := make([]*pb.MemoryLink, linkCount)
	for i := 0; i < linkCount; i++ {
		links[i] = &pb.MemoryLink{
			TargetId:  fmt.Sprintf("linked-%d", i+1),
			LinkType:  pb.LinkType_LINK_TYPE_REFERENCES,
			Strength:  0.7,
			Reason:    "test link",
			CreatedAt: uint64(time.Now().Unix()),
			UserCreated: true,
		}
	}

	return memory, links
}
