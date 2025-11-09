// Package analyze provides statistical analysis and visualization for semantic data.
package analyze

import (
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"math"
	"strings"
)

// TripleNode represents an entity node in the triple graph.
type TripleNode struct {
	ID         string             // Unique identifier (entity text)
	Entity     semantic.Entity    // The semantic entity
	X          float64            // Screen X position
	Y          float64            // Screen Y position
	VX         float64            // Velocity X (for force-directed layout)
	VY         float64            // Velocity Y
	Mass       float64            // Node mass (affects layout)
	Frequency  int                // Number of occurrences
	Importance int                // Calculated importance (0-10)
}

// TripleEdge represents a relationship edge between entities.
type TripleEdge struct {
	SourceID   string                 // Source node ID
	TargetID   string                 // Target node ID
	Relation   semantic.Relationship  // The semantic relationship
	Strength   float64                // Edge strength (0.0-1.0)
	Confidence float64                // Confidence score (0.0-1.0)
}

// TripleGraph represents the entity-relationship graph for visualization.
type TripleGraph struct {
	Nodes map[string]*TripleNode // Entity text -> node
	Edges []*TripleEdge          // All relationship edges
}

// NewTripleGraph creates an empty triple graph.
func NewTripleGraph() *TripleGraph {
	return &TripleGraph{
		Nodes: make(map[string]*TripleNode),
		Edges: make([]*TripleEdge, 0),
	}
}

// AddNode adds a node to the graph or updates frequency if it exists.
func (g *TripleGraph) AddNode(entity semantic.Entity) *TripleNode {
	if node, exists := g.Nodes[entity.Text]; exists {
		node.Frequency++
		return node
	}

	node := &TripleNode{
		ID:         entity.Text,
		Entity:     entity,
		X:          0,
		Y:          0,
		VX:         0,
		VY:         0,
		Mass:       1.0,
		Frequency:  1,
		Importance: 5, // Default mid-range importance
	}
	g.Nodes[entity.Text] = node
	return node
}

// AddEdge adds an edge to the graph.
func (g *TripleGraph) AddEdge(relation semantic.Relationship) *TripleEdge {
	edge := &TripleEdge{
		SourceID:   relation.Subject,
		TargetID:   relation.Object,
		Relation:   relation,
		Strength:   1.0,
		Confidence: 0.8, // Default confidence
	}
	g.Edges = append(g.Edges, edge)
	return edge
}

// GetNode returns a node by ID.
func (g *TripleGraph) GetNode(id string) *TripleNode {
	return g.Nodes[id]
}

// GetEdgesFrom returns all edges from a node.
func (g *TripleGraph) GetEdgesFrom(nodeID string) []*TripleEdge {
	edges := make([]*TripleEdge, 0)
	for _, edge := range g.Edges {
		if edge.SourceID == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetEdgesTo returns all edges to a node.
func (g *TripleGraph) GetEdgesTo(nodeID string) []*TripleEdge {
	edges := make([]*TripleEdge, 0)
	for _, edge := range g.Edges {
		if edge.TargetID == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// NodeCount returns the number of nodes.
func (g *TripleGraph) NodeCount() int {
	return len(g.Nodes)
}

// EdgeCount returns the number of edges.
func (g *TripleGraph) EdgeCount() int {
	return len(g.Edges)
}

// BuildFromAnalysis builds a triple graph from semantic analysis results.
func BuildFromAnalysis(analysis *semantic.Analysis) *TripleGraph {
	graph := NewTripleGraph()

	// Add all entities as nodes
	for _, entity := range analysis.Entities {
		graph.AddNode(entity)
	}

	// Add all relationships as edges
	for _, relation := range analysis.Relationships {
		// Ensure both subject and object nodes exist
		if graph.GetNode(relation.Subject) != nil && graph.GetNode(relation.Object) != nil {
			graph.AddEdge(relation)
		}
	}

	// Calculate importance scores based on frequency and connectivity
	graph.CalculateImportance()

	return graph
}

// CalculateImportance computes importance scores for all nodes.
func (g *TripleGraph) CalculateImportance() {
	if len(g.Nodes) == 0 {
		return
	}

	// Find max frequency for normalization
	maxFreq := 0
	for _, node := range g.Nodes {
		if node.Frequency > maxFreq {
			maxFreq = node.Frequency
		}
	}

	// Calculate importance based on frequency and degree centrality
	for _, node := range g.Nodes {
		// Frequency component (normalized to 0-5)
		freqScore := 0
		if maxFreq > 0 {
			freqScore = int((float64(node.Frequency) / float64(maxFreq)) * 5)
		}

		// Degree centrality component (number of connections, capped at 5)
		outEdges := len(g.GetEdgesFrom(node.ID))
		inEdges := len(g.GetEdgesTo(node.ID))
		degreeScore := min(outEdges+inEdges, 5)

		// Combined importance (0-10 scale)
		node.Importance = min(freqScore+degreeScore, 10)
	}
}

// Filter filters the graph by criteria.
type Filter struct {
	EntityTypes      map[semantic.EntityType]bool // Allowed entity types (nil = all)
	MinConfidence    float64                      // Minimum edge confidence (0.0-1.0)
	MinImportance    int                          // Minimum node importance (0-10)
	SearchTerm       string                       // Text search filter (case-insensitive)
}

// ApplyFilter returns a new graph with filtered nodes and edges.
func (g *TripleGraph) ApplyFilter(filter Filter) *TripleGraph {
	filtered := NewTripleGraph()

	// Filter nodes
	for id, node := range g.Nodes {
		// Check entity type
		if filter.EntityTypes != nil {
			if !filter.EntityTypes[node.Entity.Type] {
				continue
			}
		}

		// Check importance
		if node.Importance < filter.MinImportance {
			continue
		}

		// Check search term
		if filter.SearchTerm != "" {
			if !strings.Contains(strings.ToLower(node.Entity.Text), strings.ToLower(filter.SearchTerm)) {
				continue
			}
		}

		// Add filtered node (preserve position and velocity for smooth transitions)
		filtered.Nodes[id] = &TripleNode{
			ID:         node.ID,
			Entity:     node.Entity,
			X:          node.X,
			Y:          node.Y,
			VX:         node.VX,
			VY:         node.VY,
			Mass:       node.Mass,
			Frequency:  node.Frequency,
			Importance: node.Importance,
		}
	}

	// Filter edges (only include if both nodes are in filtered set)
	for _, edge := range g.Edges {
		// Check if both nodes exist in filtered graph
		if filtered.GetNode(edge.SourceID) == nil || filtered.GetNode(edge.TargetID) == nil {
			continue
		}

		// Check confidence
		if edge.Confidence < filter.MinConfidence {
			continue
		}

		filtered.Edges = append(filtered.Edges, &TripleEdge{
			SourceID:   edge.SourceID,
			TargetID:   edge.TargetID,
			Relation:   edge.Relation,
			Strength:   edge.Strength,
			Confidence: edge.Confidence,
		})
	}

	return filtered
}

// Layout constants (adapted from memorygraph)
const (
	RepulsionStrength = 100.0 // How strongly nodes repel each other
	AttractionStrength = 0.01 // How strongly edges attract nodes
	MaxForce          = 10.0  // Maximum force applied per iteration
	MinDistance       = 1.0   // Minimum distance to prevent division by zero
	IdealDistance     = 10.0  // Ideal distance between connected nodes
)

// InitializeLayout initializes node positions in a circle.
func (g *TripleGraph) InitializeLayout() {
	nodeCount := len(g.Nodes)
	if nodeCount == 0 {
		return
	}

	i := 0
	radius := 20.0
	for _, node := range g.Nodes {
		angle := (float64(i) / float64(nodeCount)) * 2 * math.Pi
		node.X = math.Cos(angle) * radius
		node.Y = math.Sin(angle) * radius
		node.VX = 0
		node.VY = 0
		node.Mass = 1.0
		i++
	}
}

// ApplyForceIteration applies one iteration of force-directed layout.
func (g *TripleGraph) ApplyForceIteration(damping float64) {
	g.applyRepulsion()
	g.applyAttraction()
	g.updatePositions(damping)
}

// applyRepulsion applies repulsive forces between all nodes.
func (g *TripleGraph) applyRepulsion() {
	nodes := make([]*TripleNode, 0, len(g.Nodes))
	for _, node := range g.Nodes {
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

			n1.VX -= fx
			n1.VY -= fy
			n2.VX += fx
			n2.VY += fy
		}
	}
}

// applyAttraction applies attractive forces along edges.
func (g *TripleGraph) applyAttraction() {
	for _, edge := range g.Edges {
		source := g.Nodes[edge.SourceID]
		target := g.Nodes[edge.TargetID]

		if source == nil || target == nil {
			continue
		}

		// Calculate distance
		dx := target.X - source.X
		dy := target.Y - source.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist < MinDistance {
			continue
		}

		// Calculate attractive force (Hooke's law)
		force := (dist - IdealDistance) * AttractionStrength * edge.Strength

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

		source.VX += fx
		source.VY += fy
		target.VX -= fx
		target.VY -= fy
	}
}

// updatePositions updates node positions based on velocities.
func (g *TripleGraph) updatePositions(damping float64) {
	for _, node := range g.Nodes {
		// Apply damping
		node.VX *= damping
		node.VY *= damping

		// Update position
		node.X += node.VX
		node.Y += node.VY
	}
}

// StabilizeLayout runs multiple layout iterations.
func (g *TripleGraph) StabilizeLayout(iterations int, damping float64) {
	for i := 0; i < iterations; i++ {
		g.ApplyForceIteration(damping)
	}
}

// GetBounds returns the bounding box of all nodes.
func (g *TripleGraph) GetBounds() (minX, maxX, minY, maxY float64) {
	if len(g.Nodes) == 0 {
		return 0, 0, 0, 0
	}

	first := true
	for _, node := range g.Nodes {
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

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
