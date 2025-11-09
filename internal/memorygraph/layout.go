package memorygraph

import (
	"math"
)

// Layout constants
const (
	RepulsionStrength = 100.0  // How strongly nodes repel each other
	AttractionStrength = 0.01  // How strongly edges attract nodes
	MaxForce          = 10.0   // Maximum force applied per iteration
	MinDistance       = 1.0    // Minimum distance to prevent division by zero
	IdealDistance     = 10.0   // Ideal distance between connected nodes
)

// ApplyForceLayout applies one iteration of force-directed layout.
func (m *Model) ApplyForceLayout() {
	if m.graph == nil || len(m.graph.Nodes) == 0 {
		return
	}

	// Calculate repulsive forces between all nodes
	m.applyRepulsion()

	// Calculate attractive forces along edges
	m.applyAttraction()

	// Update positions based on velocities
	m.updatePositions()

	// Center the graph
	m.centerGraph()

	m.layoutSteps++
}

// applyRepulsion applies repulsive forces between all nodes.
func (m *Model) applyRepulsion() {
	nodes := make([]*Node, 0, len(m.graph.Nodes))
	for _, node := range m.graph.Nodes {
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
func (m *Model) applyAttraction() {
	for _, edge := range m.graph.Edges {
		source := m.graph.Nodes[edge.SourceID]
		target := m.graph.Nodes[edge.TargetID]

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
		// Force is proportional to distance from ideal
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
func (m *Model) updatePositions() {
	for _, node := range m.graph.Nodes {
		// Apply damping
		node.VX *= m.damping
		node.VY *= m.damping

		// Update position
		node.X += node.VX
		node.Y += node.VY
	}
}

// centerGraph centers the graph in the viewport.
func (m *Model) centerGraph() {
	if len(m.graph.Nodes) == 0 {
		return
	}

	// Calculate bounds
	var minX, maxX, minY, maxY float64
	first := true

	for _, node := range m.graph.Nodes {
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

	// Calculate center offset
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2

	// Set offset to center the graph
	m.offsetX = -centerX
	m.offsetY = -centerY
}

// InitializeLayout initializes node positions randomly.
func (m *Model) InitializeLayout() {
	if m.graph == nil {
		return
	}

	// Simple circular layout based on node index
	nodeCount := len(m.graph.Nodes)
	if nodeCount == 0 {
		return
	}

	i := 0
	radius := 20.0
	for _, node := range m.graph.Nodes {
		angle := (float64(i) / float64(nodeCount)) * 2 * math.Pi
		node.X = math.Cos(angle) * radius
		node.Y = math.Sin(angle) * radius
		node.VX = 0
		node.VY = 0
		node.Mass = 1.0
		// Initialize nodes as expanded by default
		node.IsExpanded = true
		i++
	}
}

// StabilizeLayout runs multiple layout iterations.
func (m *Model) StabilizeLayout(iterations int) {
	for i := 0; i < iterations; i++ {
		m.ApplyForceLayout()
	}
}
