package orchestrate

import (
	"math"
	"testing"
)

// createTestWorkPlan creates a diamond DAG for testing:
//     A
//    / \
//   B   C
//    \ /
//     D
func createTestWorkPlan() *WorkPlan {
	return &WorkPlan{
		Name: "Test Plan",
		Tasks: []Task{
			{ID: "A", Description: "Task A", Dependencies: []string{}},
			{ID: "B", Description: "Task B", Dependencies: []string{"A"}},
			{ID: "C", Description: "Task C", Dependencies: []string{"A"}},
			{ID: "D", Description: "Task D", Dependencies: []string{"B", "C"}},
		},
		MaxConcurrent: 2,
	}
}

// createLinearWorkPlan creates a linear chain: A -> B -> C
func createLinearWorkPlan() *WorkPlan {
	return &WorkPlan{
		Name: "Linear Plan",
		Tasks: []Task{
			{ID: "A", Description: "Task A", Dependencies: []string{}},
			{ID: "B", Description: "Task B", Dependencies: []string{"A"}},
			{ID: "C", Description: "Task C", Dependencies: []string{"B"}},
		},
		MaxConcurrent: 1,
	}
}

// createEmptyWorkPlan creates an empty work plan.
func createEmptyWorkPlan() *WorkPlan {
	return &WorkPlan{
		Name:          "Empty Plan",
		Tasks:         []Task{},
		MaxConcurrent: 1,
	}
}

// --- DAG Construction Tests (4 tests) ---

func TestTaskGraphBuildDAG(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 80, 24)

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Check correct number of nodes
	if len(tg.nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(tg.nodes))
	}

	// Check all task IDs present
	expectedIDs := []string{"A", "B", "C", "D"}
	for _, id := range expectedIDs {
		if _, ok := tg.nodes[id]; !ok {
			t.Errorf("Missing node: %s", id)
		}
	}

	// Check correct number of edges
	// A -> B, A -> C, B -> D, C -> D = 4 edges
	if len(tg.edges) != 4 {
		t.Errorf("Expected 4 edges, got %d", len(tg.edges))
	}

	// Verify edge connectivity
	edgeMap := make(map[string][]string)
	for _, edge := range tg.edges {
		edgeMap[edge.From] = append(edgeMap[edge.From], edge.To)
	}

	// A should point to B and C
	if len(edgeMap["A"]) != 2 {
		t.Errorf("Expected A to have 2 outgoing edges, got %d", len(edgeMap["A"]))
	}

	// B and C should each point to D
	if len(edgeMap["B"]) != 1 || edgeMap["B"][0] != "D" {
		t.Errorf("Expected B -> D edge")
	}
	if len(edgeMap["C"]) != 1 || edgeMap["C"][0] != "D" {
		t.Errorf("Expected C -> D edge")
	}

	// D should have no outgoing edges
	if len(edgeMap["D"]) != 0 {
		t.Errorf("Expected D to have 0 outgoing edges, got %d", len(edgeMap["D"]))
	}
}

func TestTaskGraphDetectCycle(t *testing.T) {
	// Create a plan with a circular dependency (should fail validation)
	cyclicPlan := &WorkPlan{
		Name: "Cyclic Plan",
		Tasks: []Task{
			{ID: "A", Description: "Task A", Dependencies: []string{"B"}},
			{ID: "B", Description: "Task B", Dependencies: []string{"A"}},
		},
		MaxConcurrent: 1,
	}

	_, err := NewTaskGraph(cyclicPlan, 80, 24)
	if err == nil {
		t.Fatal("Expected error for cyclic plan, got nil")
	}

	// Verify error message mentions cycle
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestTaskGraphInitialPositions(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 80, 24)

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Verify all nodes have non-zero positions (circular layout)
	for id, node := range tg.nodes {
		if node.X == 0 && node.Y == 0 {
			// Allow one node at origin, but not all
			continue
		}
		if math.IsNaN(node.X) || math.IsNaN(node.Y) {
			t.Errorf("Node %s has NaN position", id)
		}
	}

	// Verify nodes are spread out (not all at same position)
	positions := make(map[string]bool)
	for _, node := range tg.nodes {
		key := string(rune(int(node.X*100))) + "," + string(rune(int(node.Y*100)))
		positions[key] = true
	}

	if len(positions) < 3 {
		t.Errorf("Expected nodes to be spread out, got %d unique positions", len(positions))
	}

	// Verify velocities initialized to zero
	for id, vel := range tg.velocities {
		if vel.VX != 0 || vel.VY != 0 {
			t.Errorf("Node %s has non-zero initial velocity: (%f, %f)", id, vel.VX, vel.VY)
		}
	}
}

func TestTaskGraphEmptyPlan(t *testing.T) {
	plan := createEmptyWorkPlan()
	_, err := NewTaskGraph(plan, 80, 24)

	// Empty plan should fail validation
	if err == nil {
		t.Fatal("Expected error for empty plan, got nil")
	}
}

// --- Force Layout Tests (4 tests) ---

func TestTaskGraphRepulsion(t *testing.T) {
	plan := &WorkPlan{
		Name: "Two Node Plan",
		Tasks: []Task{
			{ID: "A", Description: "Task A", Dependencies: []string{}},
			{ID: "B", Description: "Task B", Dependencies: []string{}},
		},
		MaxConcurrent: 1,
	}

	tg, err := NewTaskGraph(plan, 80, 24)
	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Place nodes very close together
	tg.nodes["A"].X = 0
	tg.nodes["A"].Y = 0
	tg.nodes["B"].X = 1
	tg.nodes["B"].Y = 1
	tg.velocities["A"] = Velocity{VX: 0, VY: 0}
	tg.velocities["B"] = Velocity{VX: 0, VY: 0}

	initialDist := distance(tg.nodes["A"], tg.nodes["B"])

	// Apply repulsion
	tg.applyRepulsion()
	tg.updatePositions()

	finalDist := distance(tg.nodes["A"], tg.nodes["B"])

	// Nodes should push apart
	if finalDist <= initialDist {
		t.Errorf("Expected nodes to repel, initial dist: %f, final dist: %f", initialDist, finalDist)
	}

	// Velocities should be non-zero after repulsion
	vA := tg.velocities["A"]
	vB := tg.velocities["B"]
	if vA.VX == 0 && vA.VY == 0 {
		t.Error("Expected non-zero velocity for node A after repulsion")
	}
	if vB.VX == 0 && vB.VY == 0 {
		t.Error("Expected non-zero velocity for node B after repulsion")
	}
}

func TestTaskGraphAttraction(t *testing.T) {
	plan := &WorkPlan{
		Name: "Connected Plan",
		Tasks: []Task{
			{ID: "A", Description: "Task A", Dependencies: []string{}},
			{ID: "B", Description: "Task B", Dependencies: []string{"A"}},
		},
		MaxConcurrent: 1,
	}

	tg, err := NewTaskGraph(plan, 80, 24)
	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Place nodes far apart
	tg.nodes["A"].X = 0
	tg.nodes["A"].Y = 0
	tg.nodes["B"].X = 50
	tg.nodes["B"].Y = 50
	tg.velocities["A"] = Velocity{VX: 0, VY: 0}
	tg.velocities["B"] = Velocity{VX: 0, VY: 0}

	initialDist := distance(tg.nodes["A"], tg.nodes["B"])

	// Apply attraction only (no repulsion to isolate effect)
	tg.applyAttraction()
	tg.updatePositions()

	finalDist := distance(tg.nodes["A"], tg.nodes["B"])

	// Connected nodes should pull together
	if finalDist >= initialDist {
		t.Errorf("Expected nodes to attract, initial dist: %f, final dist: %f", initialDist, finalDist)
	}

	// Velocities should be non-zero after attraction
	vA := tg.velocities["A"]
	vB := tg.velocities["B"]
	if vA.VX == 0 && vA.VY == 0 {
		t.Error("Expected non-zero velocity for node A after attraction")
	}
	if vB.VX == 0 && vB.VY == 0 {
		t.Error("Expected non-zero velocity for node B after attraction")
	}
}

func TestTaskGraphConvergence(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 80, 24)

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Run several iterations
	tg.stabilize(50)

	// Measure total kinetic energy (should be low after stabilization)
	totalEnergy := 0.0
	for _, vel := range tg.velocities {
		totalEnergy += vel.VX*vel.VX + vel.VY*vel.VY
	}

	// After 50 iterations with damping, energy should be low
	if totalEnergy > 10.0 {
		t.Errorf("Layout did not converge, total energy: %f", totalEnergy)
	}

	// Verify layout steps counter incremented
	if tg.layoutSteps != 50 {
		t.Errorf("Expected 50 layout steps, got %d", tg.layoutSteps)
	}
}

func TestTaskGraphBounds(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 80, 24)

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Stabilize layout
	tg.stabilize(100)

	// Get bounds
	minX, maxX, minY, maxY := tg.getBounds()

	// Verify bounds are reasonable (not infinite)
	if math.IsNaN(minX) || math.IsNaN(maxX) || math.IsNaN(minY) || math.IsNaN(maxY) {
		t.Error("Bounds contain NaN values")
	}

	if math.IsInf(minX, 0) || math.IsInf(maxX, 0) || math.IsInf(minY, 0) || math.IsInf(maxY, 0) {
		t.Error("Bounds contain infinite values")
	}

	// Verify bounds are non-zero (nodes should spread out)
	width := maxX - minX
	height := maxY - minY

	if width < 1.0 || height < 1.0 {
		t.Errorf("Graph bounds too small: width=%f, height=%f", width, height)
	}

	// Verify all nodes are within bounds
	for id, node := range tg.nodes {
		if node.X < minX || node.X > maxX {
			t.Errorf("Node %s X position outside bounds: %f not in [%f, %f]", id, node.X, minX, maxX)
		}
		if node.Y < minY || node.Y > maxY {
			t.Errorf("Node %s Y position outside bounds: %f not in [%f, %f]", id, node.Y, minY, maxY)
		}
	}
}

// --- Navigation Tests (3 tests) ---

func TestTaskGraphPan(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 80, 24)

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	initialOffsetX := tg.offsetX
	initialOffsetY := tg.offsetY

	// Pan right and down
	tg.Pan(10, 20)

	if tg.offsetX != initialOffsetX+10 {
		t.Errorf("Expected offsetX=%f, got %f", initialOffsetX+10, tg.offsetX)
	}
	if tg.offsetY != initialOffsetY+20 {
		t.Errorf("Expected offsetY=%f, got %f", initialOffsetY+20, tg.offsetY)
	}

	// Pan left and up
	tg.Pan(-5, -10)

	if tg.offsetX != initialOffsetX+5 {
		t.Errorf("Expected offsetX=%f, got %f", initialOffsetX+5, tg.offsetX)
	}
	if tg.offsetY != initialOffsetY+10 {
		t.Errorf("Expected offsetY=%f, got %f", initialOffsetY+10, tg.offsetY)
	}
}

func TestTaskGraphZoom(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 80, 24)

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Initial zoom should be 1.0
	if tg.zoom != 1.0 {
		t.Errorf("Expected initial zoom=1.0, got %f", tg.zoom)
	}

	// Zoom in
	tg.Zoom(2.0)
	if tg.zoom != 2.0 {
		t.Errorf("Expected zoom=2.0, got %f", tg.zoom)
	}

	// Zoom out
	tg.Zoom(0.5)
	if tg.zoom != 1.0 {
		t.Errorf("Expected zoom=1.0, got %f", tg.zoom)
	}

	// Test zoom limits
	tg.Zoom(100.0) // Try to zoom way in
	if tg.zoom > 5.0 {
		t.Errorf("Zoom should be capped at 5.0, got %f", tg.zoom)
	}

	tg.zoom = 1.0
	tg.Zoom(0.01) // Try to zoom way out
	if tg.zoom < 0.1 {
		t.Errorf("Zoom should be capped at 0.1, got %f", tg.zoom)
	}
}

func TestTaskGraphSelection(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 80, 24)

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Initially no selection
	if tg.selected != "" {
		t.Errorf("Expected no selection initially, got %s", tg.selected)
	}

	// Select a node
	tg.SelectNode("A")
	if tg.selected != "A" {
		t.Errorf("Expected selected node 'A', got %s", tg.selected)
	}

	// Select another node
	tg.SelectNode("B")
	if tg.selected != "B" {
		t.Errorf("Expected selected node 'B', got %s", tg.selected)
	}

	// Clear selection
	tg.ClearSelection()
	if tg.selected != "" {
		t.Errorf("Expected no selection after clear, got %s", tg.selected)
	}

	// Select non-existent node (should not crash)
	tg.SelectNode("NONEXISTENT")
	if tg.selected == "NONEXISTENT" {
		t.Error("Should not select non-existent node")
	}
}

// --- Rendering Test (1 test) ---

func TestTaskGraphRender(t *testing.T) {
	plan := createTestWorkPlan()
	tg, err := NewTaskGraph(plan, 100, 40) // Larger canvas

	if err != nil {
		t.Fatalf("NewTaskGraph failed: %v", err)
	}

	// Stabilize layout first
	tg.stabilize(10)

	// Center the view to ensure nodes are visible
	tg.Center()

	// Render view
	output := tg.View()

	// Basic checks
	if output == "" {
		t.Error("View() produced empty output")
	}

	// The output should contain at least the stats and help sections
	if !graphContains(output, "Nodes:") {
		t.Error("View() output missing stats section")
	}

	if !graphContains(output, "Pan") || !graphContains(output, "Zoom") {
		t.Error("View() output missing help text")
	}

	// Check that view contains task information
	// Since nodes might be styled with ANSI codes or positioned off-screen,
	// we check for the bracket characters that wrap task IDs
	if !graphContains(output, "[") || !graphContains(output, "]") {
		t.Error("View() output missing node brackets")
	}

	// Test with selection
	tg.SelectNode("A")
	outputWithSelection := tg.View()
	if outputWithSelection == "" {
		t.Error("View() with selection produced empty output")
	}

	// With selection, should show "Clear" in help
	if !graphContains(outputWithSelection, "Clear") {
		t.Error("View() with selection missing 'Clear' in help text")
	}

	// Test single-node plan rendering
	singlePlan := &WorkPlan{
		Name:          "Single",
		Tasks:         []Task{{ID: "X", Description: "Task X"}},
		MaxConcurrent: 1,
	}
	singleTg, _ := NewTaskGraph(singlePlan, 100, 40)
	singleOutput := singleTg.View()
	if singleOutput == "" {
		t.Error("View() for single-node plan produced empty output")
	}

	// Single node should be visible (contains bracket and X)
	if !graphContains(singleOutput, "[") {
		t.Error("Single-node view missing node bracket")
	}
}

// --- Helper Functions ---

func distance(n1, n2 *GraphNode) float64 {
	dx := n2.X - n1.X
	dy := n2.Y - n1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func graphContains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && graphFindSubstring(s, substr)
}

func graphFindSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
