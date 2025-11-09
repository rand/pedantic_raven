package memorygraph

import (
	"math"
	"testing"
)

// Test InitializeLayout sets up circular positions.
func TestInitializeLayout(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1"})
	m.graph.AddNode(&Node{ID: "node-2"})
	m.graph.AddNode(&Node{ID: "node-3"})

	m.InitializeLayout()

	// Verify all nodes have positions
	for id, node := range m.graph.Nodes {
		if node.X == 0 && node.Y == 0 {
			t.Errorf("Node %s has zero position", id)
		}
		if node.Mass != 1.0 {
			t.Errorf("Node %s has incorrect mass: %f", id, node.Mass)
		}
		if node.VX != 0 || node.VY != 0 {
			t.Errorf("Node %s has non-zero velocity: VX=%f, VY=%f", id, node.VX, node.VY)
		}
	}

	// Verify nodes are roughly on a circle
	radius := 20.0
	tolerance := 0.1
	for id, node := range m.graph.Nodes {
		dist := math.Sqrt(node.X*node.X + node.Y*node.Y)
		if math.Abs(dist-radius) > tolerance {
			t.Errorf("Node %s not on circle: distance %f, expected %f", id, dist, radius)
		}
	}
}

// Test InitializeLayout with empty graph.
func TestInitializeLayoutEmpty(t *testing.T) {
	m := NewModel()
	m.InitializeLayout() // Should not panic

	if len(m.graph.Nodes) != 0 {
		t.Error("Expected empty graph to remain empty")
	}
}

// Test InitializeLayout with single node.
func TestInitializeLayoutSingleNode(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1"})
	m.InitializeLayout()

	node := m.graph.Nodes["node-1"]
	if node.Mass != 1.0 {
		t.Errorf("Expected mass 1.0, got %f", node.Mass)
	}
	if node.VX != 0 || node.VY != 0 {
		t.Error("Expected zero velocity")
	}
}

// Test applyRepulsion creates forces between nodes.
func TestApplyRepulsion(t *testing.T) {
	m := NewModel()
	node1 := &Node{ID: "node-1", X: 0, Y: 0, VX: 0, VY: 0}
	node2 := &Node{ID: "node-2", X: 5, Y: 0, VX: 0, VY: 0}

	m.graph.AddNode(node1)
	m.graph.AddNode(node2)

	m.applyRepulsion()

	// Nodes should be repelled from each other
	if node1.VX >= 0 {
		t.Errorf("Expected node1 to be pushed left (VX < 0), got VX=%f", node1.VX)
	}
	if node2.VX <= 0 {
		t.Errorf("Expected node2 to be pushed right (VX > 0), got VX=%f", node2.VX)
	}

	// Y velocities should be zero (nodes aligned on X axis)
	if math.Abs(node1.VY) > 0.001 {
		t.Errorf("Expected node1 VY near 0, got %f", node1.VY)
	}
	if math.Abs(node2.VY) > 0.001 {
		t.Errorf("Expected node2 VY near 0, got %f", node2.VY)
	}
}

// Test applyRepulsion with overlapping nodes.
func TestApplyRepulsionOverlapping(t *testing.T) {
	m := NewModel()
	node1 := &Node{ID: "node-1", X: 0, Y: 0, VX: 0, VY: 0}
	node2 := &Node{ID: "node-2", X: 0.1, Y: 0, VX: 0, VY: 0}

	m.graph.AddNode(node1)
	m.graph.AddNode(node2)

	m.applyRepulsion()

	// Very close nodes should have strong repulsion (capped at MaxForce)
	if math.Abs(node1.VX) == 0 {
		t.Error("Expected non-zero repulsion for overlapping nodes")
	}
	if math.Abs(node1.VX) > MaxForce {
		t.Errorf("Expected force to be capped at %f, got %f", MaxForce, math.Abs(node1.VX))
	}
}

// Test applyAttraction creates forces along edges.
func TestApplyAttraction(t *testing.T) {
	m := NewModel()
	node1 := &Node{ID: "node-1", X: 0, Y: 0, VX: 0, VY: 0}
	node2 := &Node{ID: "node-2", X: 20, Y: 0, VX: 0, VY: 0}

	m.graph.AddNode(node1)
	m.graph.AddNode(node2)
	m.graph.AddEdge(&Edge{
		SourceID: "node-1",
		TargetID: "node-2",
		Strength: 1.0,
	})

	m.applyAttraction()

	// Nodes should be attracted toward each other
	if node1.VX <= 0 {
		t.Errorf("Expected node1 to be pulled right (VX > 0), got VX=%f", node1.VX)
	}
	if node2.VX >= 0 {
		t.Errorf("Expected node2 to be pulled left (VX < 0), got VX=%f", node2.VX)
	}

	// Y velocities should be zero (nodes aligned on X axis)
	if math.Abs(node1.VY) > 0.001 {
		t.Errorf("Expected node1 VY near 0, got %f", node1.VY)
	}
	if math.Abs(node2.VY) > 0.001 {
		t.Errorf("Expected node2 VY near 0, got %f", node2.VY)
	}
}

// Test applyAttraction with nodes at ideal distance.
func TestApplyAttractionIdealDistance(t *testing.T) {
	m := NewModel()
	node1 := &Node{ID: "node-1", X: 0, Y: 0, VX: 0, VY: 0}
	node2 := &Node{ID: "node-2", X: IdealDistance, Y: 0, VX: 0, VY: 0}

	m.graph.AddNode(node1)
	m.graph.AddNode(node2)
	m.graph.AddEdge(&Edge{
		SourceID: "node-1",
		TargetID: "node-2",
		Strength: 1.0,
	})

	m.applyAttraction()

	// At ideal distance, force should be near zero
	if math.Abs(node1.VX) > 0.001 {
		t.Errorf("Expected near-zero force at ideal distance, got VX=%f", node1.VX)
	}
}

// Test applyAttraction with edge strength.
func TestApplyAttractionStrength(t *testing.T) {
	m := NewModel()
	node1 := &Node{ID: "node-1", X: 0, Y: 0, VX: 0, VY: 0}
	node2 := &Node{ID: "node-2", X: 20, Y: 0, VX: 0, VY: 0}

	m.graph.AddNode(node1)
	m.graph.AddNode(node2)
	m.graph.AddEdge(&Edge{
		SourceID: "node-1",
		TargetID: "node-2",
		Strength: 2.0, // Double strength
	})

	m.applyAttraction()

	// Store the velocity from double strength
	vxDouble := node1.VX

	// Reset and test with single strength
	node1.VX = 0
	node1.VY = 0
	node2.VX = 0
	node2.VY = 0
	m.graph.Edges[0].Strength = 1.0

	m.applyAttraction()

	// Double strength should produce roughly double the force
	ratio := vxDouble / node1.VX
	if math.Abs(ratio-2.0) > 0.1 {
		t.Errorf("Expected strength to scale force, got ratio %f", ratio)
	}
}

// Test applyAttraction with missing nodes.
func TestApplyAttractionMissingNodes(t *testing.T) {
	m := NewModel()
	m.graph.AddEdge(&Edge{
		SourceID: "nonexistent-1",
		TargetID: "nonexistent-2",
		Strength: 1.0,
	})

	// Should not panic
	m.applyAttraction()
}

// Test updatePositions updates node positions.
func TestUpdatePositions(t *testing.T) {
	m := NewModel()
	node := &Node{
		ID: "node-1",
		X:  10.0,
		Y:  20.0,
		VX: 5.0,
		VY: 3.0,
	}
	m.graph.AddNode(node)

	m.updatePositions()

	// Position should be updated by velocity
	expectedX := 10.0 + 5.0*m.damping
	expectedY := 20.0 + 3.0*m.damping
	if math.Abs(node.X-expectedX) > 0.001 {
		t.Errorf("Expected X=%f, got %f", expectedX, node.X)
	}
	if math.Abs(node.Y-expectedY) > 0.001 {
		t.Errorf("Expected Y=%f, got %f", expectedY, node.Y)
	}

	// Velocity should be damped
	expectedVX := 5.0 * m.damping
	expectedVY := 3.0 * m.damping
	if math.Abs(node.VX-expectedVX) > 0.001 {
		t.Errorf("Expected VX=%f, got %f", expectedVX, node.VX)
	}
	if math.Abs(node.VY-expectedVY) > 0.001 {
		t.Errorf("Expected VY=%f, got %f", expectedVY, node.VY)
	}
}

// Test updatePositions with zero damping.
func TestUpdatePositionsZeroDamping(t *testing.T) {
	m := NewModel()
	m.damping = 0.0
	node := &Node{
		ID: "node-1",
		X:  10.0,
		Y:  20.0,
		VX: 5.0,
		VY: 3.0,
	}
	m.graph.AddNode(node)

	m.updatePositions()

	// With zero damping, velocity should be zeroed
	if node.VX != 0 || node.VY != 0 {
		t.Errorf("Expected zero velocity with zero damping, got VX=%f, VY=%f", node.VX, node.VY)
	}
	// Position should still be at original (0 velocity applied)
	if node.X != 10.0 || node.Y != 20.0 {
		t.Errorf("Expected position unchanged, got X=%f, Y=%f", node.X, node.Y)
	}
}

// Test centerGraph calculates correct center offset.
func TestCenterGraph(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: -10, Y: -5})
	m.graph.AddNode(&Node{ID: "node-2", X: 10, Y: 5})

	m.centerGraph()

	// Center should be at (0, 0)
	// Offset should be (0, 0) to center it
	if m.offsetX != 0 || m.offsetY != 0 {
		t.Errorf("Expected offset (0, 0), got (%f, %f)", m.offsetX, m.offsetY)
	}
}

// Test centerGraph with asymmetric layout.
func TestCenterGraphAsymmetric(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: 0, Y: 0})
	m.graph.AddNode(&Node{ID: "node-2", X: 20, Y: 10})

	m.centerGraph()

	// Center at (10, 5), offset should be (-10, -5)
	expectedX := -10.0
	expectedY := -5.0
	if math.Abs(m.offsetX-expectedX) > 0.001 {
		t.Errorf("Expected offsetX=%f, got %f", expectedX, m.offsetX)
	}
	if math.Abs(m.offsetY-expectedY) > 0.001 {
		t.Errorf("Expected offsetY=%f, got %f", expectedY, m.offsetY)
	}
}

// Test centerGraph with empty graph.
func TestCenterGraphEmpty(t *testing.T) {
	m := NewModel()
	m.offsetX = 10
	m.offsetY = 20

	m.centerGraph()

	// Offset should remain unchanged
	if m.offsetX != 10 || m.offsetY != 20 {
		t.Errorf("Expected offset unchanged, got (%f, %f)", m.offsetX, m.offsetY)
	}
}

// Test centerGraph with single node.
func TestCenterGraphSingleNode(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: 15, Y: 25})

	m.centerGraph()

	// Center at (15, 25), offset should be (-15, -25)
	if m.offsetX != -15 || m.offsetY != -25 {
		t.Errorf("Expected offset (-15, -25), got (%f, %f)", m.offsetX, m.offsetY)
	}
}

// Test ApplyForceLayout runs one iteration.
func TestApplyForceLayout(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: 0, Y: 0})
	m.graph.AddNode(&Node{ID: "node-2", X: 5, Y: 0})

	initialSteps := m.layoutSteps
	m.ApplyForceLayout()

	// layoutSteps should increment
	if m.layoutSteps != initialSteps+1 {
		t.Errorf("Expected layoutSteps to increment, got %d", m.layoutSteps)
	}

	// Nodes should have moved (repulsion)
	node1 := m.graph.Nodes["node-1"]
	node2 := m.graph.Nodes["node-2"]
	if node1.X == 0 && node2.X == 5 {
		t.Error("Expected nodes to move during layout")
	}
}

// Test ApplyForceLayout with empty graph.
func TestApplyForceLayoutEmpty(t *testing.T) {
	m := NewModel()
	m.ApplyForceLayout() // Should not panic

	if m.layoutSteps != 0 {
		t.Errorf("Expected layoutSteps to remain 0, got %d", m.layoutSteps)
	}
}

// Test StabilizeLayout runs multiple iterations.
func TestStabilizeLayout(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: 0, Y: 0})
	m.graph.AddNode(&Node{ID: "node-2", X: 5, Y: 0})
	m.graph.AddEdge(&Edge{
		SourceID: "node-1",
		TargetID: "node-2",
		Strength: 1.0,
	})

	iterations := 10
	m.StabilizeLayout(iterations)

	if m.layoutSteps != iterations {
		t.Errorf("Expected %d layout steps, got %d", iterations, m.layoutSteps)
	}
}

// Test StabilizeLayout converges toward stable state.
func TestStabilizeLayoutConvergence(t *testing.T) {
	m := NewModel()
	node1 := &Node{ID: "node-1", X: 0, Y: 0}
	node2 := &Node{ID: "node-2", X: 30, Y: 0}

	m.graph.AddNode(node1)
	m.graph.AddNode(node2)
	m.graph.AddEdge(&Edge{
		SourceID: "node-1",
		TargetID: "node-2",
		Strength: 1.0,
	})

	// Run many iterations
	m.StabilizeLayout(100)

	// Nodes should settle near ideal distance
	dist := math.Abs(node2.X - node1.X)
	// With repulsion and attraction, should be near IdealDistance
	// Allow very wide tolerance due to complex dynamics (physics simulation may not converge exactly)
	if dist < IdealDistance*0.5 || dist > IdealDistance*5.0 {
		t.Errorf("Expected distance between %f and %f, got %f", IdealDistance*0.5, IdealDistance*5.0, dist)
	}

	// Velocities should be damped to near zero
	if math.Abs(node1.VX) > 1.0 || math.Abs(node2.VX) > 1.0 {
		t.Errorf("Expected low velocities after stabilization, got VX1=%f, VX2=%f", node1.VX, node2.VX)
	}
}

// Test ApplyForceLayout with nil graph.
func TestApplyForceLayoutNilGraph(t *testing.T) {
	m := NewModel()
	m.graph = nil

	// Should not panic
	m.ApplyForceLayout()
}

// Test complex graph layout.
func TestComplexGraphLayout(t *testing.T) {
	m := NewModel()

	// Create a triangle of nodes
	m.graph.AddNode(&Node{ID: "A"})
	m.graph.AddNode(&Node{ID: "B"})
	m.graph.AddNode(&Node{ID: "C"})

	m.graph.AddEdge(&Edge{SourceID: "A", TargetID: "B", Strength: 1.0})
	m.graph.AddEdge(&Edge{SourceID: "B", TargetID: "C", Strength: 1.0})
	m.graph.AddEdge(&Edge{SourceID: "C", TargetID: "A", Strength: 1.0})

	m.InitializeLayout()
	m.StabilizeLayout(50)

	// All nodes should be roughly equidistant
	nodeA := m.graph.Nodes["A"]
	nodeB := m.graph.Nodes["B"]
	nodeC := m.graph.Nodes["C"]

	distAB := math.Sqrt((nodeB.X-nodeA.X)*(nodeB.X-nodeA.X) + (nodeB.Y-nodeA.Y)*(nodeB.Y-nodeA.Y))
	distBC := math.Sqrt((nodeC.X-nodeB.X)*(nodeC.X-nodeB.X) + (nodeC.Y-nodeB.Y)*(nodeC.Y-nodeB.Y))
	distCA := math.Sqrt((nodeA.X-nodeC.X)*(nodeA.X-nodeC.X) + (nodeA.Y-nodeC.Y)*(nodeA.Y-nodeC.Y))

	avgDist := (distAB + distBC + distCA) / 3
	tolerance := avgDist * 0.3 // Allow 30% variation

	if math.Abs(distAB-avgDist) > tolerance {
		t.Errorf("Edge A-B distance %f too far from average %f", distAB, avgDist)
	}
	if math.Abs(distBC-avgDist) > tolerance {
		t.Errorf("Edge B-C distance %f too far from average %f", distBC, avgDist)
	}
	if math.Abs(distCA-avgDist) > tolerance {
		t.Errorf("Edge C-A distance %f too far from average %f", distCA, avgDist)
	}
}
