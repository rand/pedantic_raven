package analyze

import (
	"math/rand"
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// BenchmarkTripleGraphAddNode benchmarks adding nodes to the graph.
func BenchmarkTripleGraphAddNode(b *testing.B) {
	graph := NewTripleGraph()
	entity := semantic.Entity{
		Text: "test entity",
		Type: semantic.EntityConcept,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.AddNode(entity)
	}
}

// BenchmarkTripleGraphAddEdge benchmarks adding edges to the graph.
func BenchmarkTripleGraphAddEdge(b *testing.B) {
	graph := NewTripleGraph()
	relation := semantic.Relationship{
		Subject:   "entity1",
		Predicate: "relates_to",
		Object:    "entity2",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.AddEdge(relation)
	}
}

// BenchmarkTripleGraphLayout benchmarks force-directed layout iterations.
func BenchmarkTripleGraphLayout(b *testing.B) {
	sizes := []struct {
		name  string
		nodes int
		edges int
	}{
		{"10nodes", 10, 15},
		{"50nodes", 50, 80},
		{"100nodes", 100, 200},
		{"200nodes", 200, 400},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			graph := createTestGraphWithSize(size.nodes, size.edges)
			graph.InitializeLayout()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.ApplyForceIteration(0.8)
			}
		})
	}
}

// BenchmarkTripleGraphFullLayoutStabilization benchmarks complete layout stabilization.
func BenchmarkTripleGraphFullLayoutStabilization(b *testing.B) {
	sizes := []struct {
		name       string
		nodes      int
		edges      int
		iterations int
	}{
		{"10nodes_50iter", 10, 15, 50},
		{"50nodes_50iter", 50, 80, 50},
		{"100nodes_50iter", 100, 200, 50},
		{"100nodes_100iter", 100, 200, 100},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				graph := createTestGraphWithSize(size.nodes, size.edges)
				graph.InitializeLayout()
				b.StartTimer()

				graph.StabilizeLayout(size.iterations, 0.8)
			}
		})
	}
}

// BenchmarkTripleGraphRepulsion benchmarks repulsion force calculation.
func BenchmarkTripleGraphRepulsion(b *testing.B) {
	sizes := []int{10, 50, 100, 200}

	for _, size := range sizes {
		b.Run(string(rune(size))+"nodes", func(b *testing.B) {
			graph := createTestGraphWithSize(size, size*2)
			graph.InitializeLayout()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.applyRepulsion()
			}
		})
	}
}

// BenchmarkTripleGraphAttraction benchmarks attraction force calculation.
func BenchmarkTripleGraphAttraction(b *testing.B) {
	sizes := []int{10, 50, 100, 200}

	for _, size := range sizes {
		b.Run(string(rune(size))+"nodes", func(b *testing.B) {
			graph := createTestGraphWithSize(size, size*2)
			graph.InitializeLayout()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.applyAttraction()
			}
		})
	}
}

// BenchmarkTripleGraphUpdatePositions benchmarks position updates.
func BenchmarkTripleGraphUpdatePositions(b *testing.B) {
	sizes := []int{10, 50, 100, 200}

	for _, size := range sizes {
		b.Run(string(rune(size))+"nodes", func(b *testing.B) {
			graph := createTestGraphWithSize(size, size*2)
			graph.InitializeLayout()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.updatePositions(0.8)
			}
		})
	}
}

// BenchmarkTripleGraphCalculateImportance benchmarks importance score calculation.
func BenchmarkTripleGraphCalculateImportance(b *testing.B) {
	sizes := []int{10, 50, 100, 200}

	for _, size := range sizes {
		b.Run(string(rune(size))+"nodes", func(b *testing.B) {
			graph := createTestGraphWithSize(size, size*2)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.CalculateImportance()
			}
		})
	}
}

// BenchmarkTripleGraphApplyFilter benchmarks graph filtering operations.
func BenchmarkTripleGraphApplyFilter(b *testing.B) {
	graph := createTestGraphWithSize(100, 200)
	graph.CalculateImportance()

	filters := []struct {
		name   string
		filter Filter
	}{
		{"NoFilter", Filter{}},
		{"MinImportance5", Filter{MinImportance: 5}},
		{"MinConfidence0.5", Filter{MinConfidence: 0.5}},
		{"SearchTerm", Filter{SearchTerm: "entity"}},
		{"Combined", Filter{MinImportance: 3, MinConfidence: 0.3, SearchTerm: "test"}},
	}

	for _, f := range filters {
		b.Run(f.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = graph.ApplyFilter(f.filter)
			}
		})
	}
}

// BenchmarkTripleGraphGetEdgesFrom benchmarks edge lookup by source.
func BenchmarkTripleGraphGetEdgesFrom(b *testing.B) {
	graph := createTestGraphWithSize(100, 300)
	nodeID := "entity_0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.GetEdgesFrom(nodeID)
	}
}

// BenchmarkTripleGraphGetEdgesTo benchmarks edge lookup by target.
func BenchmarkTripleGraphGetEdgesTo(b *testing.B) {
	graph := createTestGraphWithSize(100, 300)
	nodeID := "entity_0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = graph.GetEdgesTo(nodeID)
	}
}

// BenchmarkTripleGraphBuildFromAnalysis benchmarks building graph from semantic analysis.
func BenchmarkTripleGraphBuildFromAnalysis(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(string(rune(size))+"entities", func(b *testing.B) {
			analysis := createTestAnalysis(size, size*2)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = BuildFromAnalysis(analysis)
			}
		})
	}
}

// Helper functions

func createTestGraphWithSize(numNodes, numEdges int) *TripleGraph {
	graph := NewTripleGraph()
	rng := rand.New(rand.NewSource(42))

	// Add nodes
	for i := 0; i < numNodes; i++ {
		entity := semantic.Entity{
			Text: string(rune('a' + (i % 26))) + string(rune(i)),
			Type: semantic.EntityConcept,
		}
		graph.AddNode(entity)
	}

	// Add edges
	nodeIDs := make([]string, 0, numNodes)
	for id := range graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}

	for i := 0; i < numEdges && len(nodeIDs) > 1; i++ {
		sourceIdx := rng.Intn(len(nodeIDs))
		targetIdx := rng.Intn(len(nodeIDs))
		if sourceIdx != targetIdx {
			relation := semantic.Relationship{
				Subject:   nodeIDs[sourceIdx],
				Predicate: "relates_to",
				Object:    nodeIDs[targetIdx],
			}
			graph.AddEdge(relation)
		}
	}

	return graph
}

func createTestAnalysis(numEntities, numRelations int) *semantic.Analysis {
	analysis := &semantic.Analysis{
		Entities:      make([]semantic.Entity, numEntities),
		Relationships: make([]semantic.Relationship, 0, numRelations),
	}

	// Create entities
	for i := 0; i < numEntities; i++ {
		analysis.Entities[i] = semantic.Entity{
			Text: string(rune('a' + (i % 26))) + string(rune(i)),
			Type: semantic.EntityConcept,
		}
	}

	// Create relationships
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < numRelations && numEntities > 1; i++ {
		sourceIdx := rng.Intn(numEntities)
		targetIdx := rng.Intn(numEntities)
		if sourceIdx != targetIdx {
			analysis.Relationships = append(analysis.Relationships, semantic.Relationship{
				Subject:   analysis.Entities[sourceIdx].Text,
				Predicate: "relates_to",
				Object:    analysis.Entities[targetIdx].Text,
			})
		}
	}

	return analysis
}
