package analyze

import (
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"math"
	"testing"
	"time"
)

// Test empty graph creation
func TestNewTripleGraph(t *testing.T) {
	graph := NewTripleGraph()

	if graph == nil {
		t.Fatal("NewTripleGraph returned nil")
	}

	if graph.Nodes == nil {
		t.Fatal("Nodes map is nil")
	}

	if graph.Edges == nil {
		t.Fatal("Edges slice is nil")
	}

	if graph.NodeCount() != 0 {
		t.Errorf("Expected 0 nodes, got %d", graph.NodeCount())
	}

	if graph.EdgeCount() != 0 {
		t.Errorf("Expected 0 edges, got %d", graph.EdgeCount())
	}
}

// Test adding nodes
func TestAddNode(t *testing.T) {
	graph := NewTripleGraph()

	entity1 := semantic.Entity{
		Text: "John",
		Type: semantic.EntityPerson,
	}

	node1 := graph.AddNode(entity1)

	if node1 == nil {
		t.Fatal("AddNode returned nil")
	}

	if node1.ID != "John" {
		t.Errorf("Expected ID 'John', got '%s'", node1.ID)
	}

	if node1.Frequency != 1 {
		t.Errorf("Expected frequency 1, got %d", node1.Frequency)
	}

	if graph.NodeCount() != 1 {
		t.Errorf("Expected 1 node, got %d", graph.NodeCount())
	}

	// Add same entity again - should increment frequency
	node2 := graph.AddNode(entity1)

	if node2.Frequency != 2 {
		t.Errorf("Expected frequency 2, got %d", node2.Frequency)
	}

	if graph.NodeCount() != 1 {
		t.Errorf("Expected still 1 node, got %d", graph.NodeCount())
	}
}

// Test adding edges
func TestAddEdge(t *testing.T) {
	graph := NewTripleGraph()

	// Add nodes first
	graph.AddNode(semantic.Entity{Text: "John", Type: semantic.EntityPerson})
	graph.AddNode(semantic.Entity{Text: "Acme Corp", Type: semantic.EntityOrganization})

	relation := semantic.Relationship{
		Subject:   "John",
		Predicate: "works_at",
		Object:    "Acme Corp",
	}

	edge := graph.AddEdge(relation)

	if edge == nil {
		t.Fatal("AddEdge returned nil")
	}

	if edge.SourceID != "John" {
		t.Errorf("Expected source 'John', got '%s'", edge.SourceID)
	}

	if edge.TargetID != "Acme Corp" {
		t.Errorf("Expected target 'Acme Corp', got '%s'", edge.TargetID)
	}

	if graph.EdgeCount() != 1 {
		t.Errorf("Expected 1 edge, got %d", graph.EdgeCount())
	}
}

// Test building from analysis
func TestBuildFromAnalysis(t *testing.T) {
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "John", Type: semantic.EntityPerson},
			{Text: "Acme Corp", Type: semantic.EntityOrganization},
			{Text: "API", Type: semantic.EntityTechnology},
		},
		Relationships: []semantic.Relationship{
			{Subject: "John", Predicate: "works_at", Object: "Acme Corp"},
			{Subject: "Acme Corp", Predicate: "develops", Object: "API"},
		},
		Timestamp: time.Now(),
	}

	graph := BuildFromAnalysis(analysis)

	if graph.NodeCount() != 3 {
		t.Errorf("Expected 3 nodes, got %d", graph.NodeCount())
	}

	if graph.EdgeCount() != 2 {
		t.Errorf("Expected 2 edges, got %d", graph.EdgeCount())
	}

	// Check specific nodes exist
	if graph.GetNode("John") == nil {
		t.Error("Expected to find 'John' node")
	}

	if graph.GetNode("Acme Corp") == nil {
		t.Error("Expected to find 'Acme Corp' node")
	}

	if graph.GetNode("API") == nil {
		t.Error("Expected to find 'API' node")
	}
}

// Test edge retrieval
func TestGetEdges(t *testing.T) {
	graph := NewTripleGraph()

	graph.AddNode(semantic.Entity{Text: "A", Type: semantic.EntityConcept})
	graph.AddNode(semantic.Entity{Text: "B", Type: semantic.EntityConcept})
	graph.AddNode(semantic.Entity{Text: "C", Type: semantic.EntityConcept})

	graph.AddEdge(semantic.Relationship{Subject: "A", Predicate: "rel1", Object: "B"})
	graph.AddEdge(semantic.Relationship{Subject: "A", Predicate: "rel2", Object: "C"})
	graph.AddEdge(semantic.Relationship{Subject: "B", Predicate: "rel3", Object: "C"})

	// Test GetEdgesFrom
	edgesFromA := graph.GetEdgesFrom("A")
	if len(edgesFromA) != 2 {
		t.Errorf("Expected 2 edges from A, got %d", len(edgesFromA))
	}

	edgesFromB := graph.GetEdgesFrom("B")
	if len(edgesFromB) != 1 {
		t.Errorf("Expected 1 edge from B, got %d", len(edgesFromB))
	}

	// Test GetEdgesTo
	edgesToC := graph.GetEdgesTo("C")
	if len(edgesToC) != 2 {
		t.Errorf("Expected 2 edges to C, got %d", len(edgesToC))
	}

	edgesToA := graph.GetEdgesTo("A")
	if len(edgesToA) != 0 {
		t.Errorf("Expected 0 edges to A, got %d", len(edgesToA))
	}
}

// Test importance calculation
func TestCalculateImportance(t *testing.T) {
	graph := NewTripleGraph()

	// Add entities with different frequencies
	for i := 0; i < 5; i++ {
		graph.AddNode(semantic.Entity{Text: "High", Type: semantic.EntityConcept})
	}
	graph.AddNode(semantic.Entity{Text: "Low", Type: semantic.EntityConcept})

	// Add some edges to test degree centrality
	graph.AddEdge(semantic.Relationship{Subject: "High", Predicate: "rel", Object: "Low"})

	graph.CalculateImportance()

	highNode := graph.GetNode("High")
	lowNode := graph.GetNode("Low")

	if highNode.Importance <= lowNode.Importance {
		t.Errorf("Expected 'High' importance (%d) > 'Low' importance (%d)",
			highNode.Importance, lowNode.Importance)
	}

	if highNode.Importance < 0 || highNode.Importance > 10 {
		t.Errorf("Importance should be 0-10, got %d", highNode.Importance)
	}
}

// Test filtering by entity type
func TestFilterByEntityType(t *testing.T) {
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "John", Type: semantic.EntityPerson},
			{Text: "Acme Corp", Type: semantic.EntityOrganization},
			{Text: "New York", Type: semantic.EntityPlace},
		},
		Relationships: []semantic.Relationship{
			{Subject: "John", Predicate: "works_at", Object: "Acme Corp"},
		},
	}

	graph := BuildFromAnalysis(analysis)

	// Filter to only Person entities
	filter := Filter{
		EntityTypes: map[semantic.EntityType]bool{
			semantic.EntityPerson: true,
		},
	}

	filtered := graph.ApplyFilter(filter)

	if filtered.NodeCount() != 1 {
		t.Errorf("Expected 1 node after filtering, got %d", filtered.NodeCount())
	}

	if filtered.GetNode("John") == nil {
		t.Error("Expected 'John' node in filtered graph")
	}

	if filtered.GetNode("Acme Corp") != nil {
		t.Error("Did not expect 'Acme Corp' in filtered graph")
	}

	// Edges should be removed if nodes are filtered
	if filtered.EdgeCount() != 0 {
		t.Errorf("Expected 0 edges after filtering, got %d", filtered.EdgeCount())
	}
}

// Test filtering by importance
func TestFilterByImportance(t *testing.T) {
	graph := NewTripleGraph()

	// Create nodes with different importance levels
	for i := 0; i < 10; i++ {
		graph.AddNode(semantic.Entity{Text: "High", Type: semantic.EntityConcept})
	}
	graph.AddNode(semantic.Entity{Text: "Low", Type: semantic.EntityConcept})

	graph.CalculateImportance()

	// Filter for high importance only
	filter := Filter{
		MinImportance: 5,
	}

	filtered := graph.ApplyFilter(filter)

	// At least the high-frequency node should pass
	if filtered.NodeCount() == 0 {
		t.Error("Expected at least one node after importance filtering")
	}

	// Check all filtered nodes meet minimum importance
	for _, node := range filtered.Nodes {
		if node.Importance < 5 {
			t.Errorf("Node '%s' has importance %d, below minimum 5",
				node.ID, node.Importance)
		}
	}
}

// Test search term filtering
func TestFilterBySearchTerm(t *testing.T) {
	graph := NewTripleGraph()

	graph.AddNode(semantic.Entity{Text: "JavaScript", Type: semantic.EntityTechnology})
	graph.AddNode(semantic.Entity{Text: "Python", Type: semantic.EntityTechnology})
	graph.AddNode(semantic.Entity{Text: "Java", Type: semantic.EntityTechnology})

	// Filter for "Java" (should match JavaScript and Java)
	filter := Filter{
		SearchTerm: "java",
	}

	filtered := graph.ApplyFilter(filter)

	if filtered.NodeCount() != 2 {
		t.Errorf("Expected 2 nodes matching 'java', got %d", filtered.NodeCount())
	}

	if filtered.GetNode("JavaScript") == nil {
		t.Error("Expected 'JavaScript' in filtered results")
	}

	if filtered.GetNode("Java") == nil {
		t.Error("Expected 'Java' in filtered results")
	}

	if filtered.GetNode("Python") != nil {
		t.Error("Did not expect 'Python' in filtered results")
	}
}

// Test layout initialization
func TestInitializeLayout(t *testing.T) {
	graph := NewTripleGraph()

	// Add several nodes
	for i := 0; i < 5; i++ {
		entity := semantic.Entity{
			Text: string(rune('A' + i)),
			Type: semantic.EntityConcept,
		}
		graph.AddNode(entity)
	}

	graph.InitializeLayout()

	// Check all nodes have positions
	for _, node := range graph.Nodes {
		if node.X == 0 && node.Y == 0 {
			// OK if it's the first node
			continue
		}

		// Check positions are within reasonable bounds (circular layout)
		dist := math.Sqrt(node.X*node.X + node.Y*node.Y)
		if dist > 25.0 { // Radius is 20, allow some margin
			t.Errorf("Node '%s' position (%f, %f) is too far from center",
				node.ID, node.X, node.Y)
		}
	}

	// Check velocities are reset
	for _, node := range graph.Nodes {
		if node.VX != 0 || node.VY != 0 {
			t.Errorf("Node '%s' should have zero velocity after init, got (%f, %f)",
				node.ID, node.VX, node.VY)
		}
	}
}

// Test force-directed layout iteration
func TestApplyForceIteration(t *testing.T) {
	graph := NewTripleGraph()

	// Create two nodes
	graph.AddNode(semantic.Entity{Text: "A", Type: semantic.EntityConcept})
	graph.AddNode(semantic.Entity{Text: "B", Type: semantic.EntityConcept})

	// Initialize positions (will place them at different points on circle)
	graph.InitializeLayout()

	// Add edge to create attraction
	graph.AddEdge(semantic.Relationship{Subject: "A", Predicate: "rel", Object: "B"})

	// Store initial positions
	nodeA := graph.GetNode("A")
	nodeB := graph.GetNode("B")
	initDistX := nodeB.X - nodeA.X
	initDistY := nodeB.Y - nodeA.Y
	initDist := math.Sqrt(initDistX*initDistX + initDistY*initDistY)

	// Apply several force iterations
	for i := 0; i < 10; i++ {
		graph.ApplyForceIteration(0.8)
	}

	// Check that nodes have moved (velocities should be non-zero at some point)
	// After stabilization, they should settle near IdealDistance apart
	finalDistX := nodeB.X - nodeA.X
	finalDistY := nodeB.Y - nodeA.Y
	finalDist := math.Sqrt(finalDistX*finalDistX + finalDistY*finalDistY)

	// Distance should change from initial
	if math.Abs(finalDist-initDist) < 0.1 {
		t.Error("Expected nodes to move during layout iterations")
	}
}

// Test layout stabilization
func TestStabilizeLayout(t *testing.T) {
	graph := NewTripleGraph()

	// Create a simple triangle of nodes
	graph.AddNode(semantic.Entity{Text: "A", Type: semantic.EntityConcept})
	graph.AddNode(semantic.Entity{Text: "B", Type: semantic.EntityConcept})
	graph.AddNode(semantic.Entity{Text: "C", Type: semantic.EntityConcept})

	graph.AddEdge(semantic.Relationship{Subject: "A", Predicate: "rel", Object: "B"})
	graph.AddEdge(semantic.Relationship{Subject: "B", Predicate: "rel", Object: "C"})
	graph.AddEdge(semantic.Relationship{Subject: "C", Predicate: "rel", Object: "A"})

	graph.InitializeLayout()
	graph.StabilizeLayout(50, 0.8)

	// After stabilization, velocities should be very small
	for _, node := range graph.Nodes {
		speed := math.Sqrt(node.VX*node.VX + node.VY*node.VY)
		if speed > 1.0 {
			t.Errorf("Node '%s' still has high velocity %f after stabilization",
				node.ID, speed)
		}
	}
}

// Test bounds calculation
func TestGetBounds(t *testing.T) {
	graph := NewTripleGraph()

	// Empty graph
	minX, maxX, minY, maxY := graph.GetBounds()
	if minX != 0 || maxX != 0 || minY != 0 || maxY != 0 {
		t.Error("Empty graph should have zero bounds")
	}

	// Add nodes at known positions
	graph.AddNode(semantic.Entity{Text: "A", Type: semantic.EntityConcept})
	graph.AddNode(semantic.Entity{Text: "B", Type: semantic.EntityConcept})

	graph.GetNode("A").X = -10
	graph.GetNode("A").Y = -5
	graph.GetNode("B").X = 15
	graph.GetNode("B").Y = 20

	minX, maxX, minY, maxY = graph.GetBounds()

	if minX != -10 {
		t.Errorf("Expected minX -10, got %f", minX)
	}
	if maxX != 15 {
		t.Errorf("Expected maxX 15, got %f", maxX)
	}
	if minY != -5 {
		t.Errorf("Expected minY -5, got %f", minY)
	}
	if maxY != 20 {
		t.Errorf("Expected maxY 20, got %f", maxY)
	}
}

// Test large graph performance
func TestLargeGraph(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large graph test in short mode")
	}

	// Create a graph with 100 nodes and 200 edges
	analysis := &semantic.Analysis{
		Entities:      make([]semantic.Entity, 100),
		Relationships: make([]semantic.Relationship, 200),
	}

	// Generate entities
	for i := 0; i < 100; i++ {
		analysis.Entities[i] = semantic.Entity{
			Text: string(rune('A' + (i % 26))) + string(rune('0' + (i / 26))),
			Type: semantic.EntityType(i % 6),
		}
	}

	// Generate relationships
	for i := 0; i < 200; i++ {
		source := analysis.Entities[i%100].Text
		target := analysis.Entities[(i+1)%100].Text
		analysis.Relationships[i] = semantic.Relationship{
			Subject:   source,
			Predicate: "relates_to",
			Object:    target,
		}
	}

	// Build and stabilize graph
	graph := BuildFromAnalysis(analysis)
	graph.InitializeLayout()

	start := time.Now()
	graph.StabilizeLayout(50, 0.8)
	duration := time.Since(start)

	// Should complete in reasonable time (< 1 second for 100 nodes)
	if duration > time.Second {
		t.Errorf("Large graph stabilization took %v, expected < 1s", duration)
	}

	// Verify graph integrity
	if graph.NodeCount() != 100 {
		t.Errorf("Expected 100 nodes, got %d", graph.NodeCount())
	}

	if graph.EdgeCount() != 200 {
		t.Errorf("Expected 200 edges, got %d", graph.EdgeCount())
	}
}
