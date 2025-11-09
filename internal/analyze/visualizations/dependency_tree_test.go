package visualizations

import (
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// TestNewDependencyTree tests tree creation.
func TestNewDependencyTree(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "TestHole"}},
		},
	}

	config := DefaultDependencyTreeConfig()
	tree := NewDependencyTree(analysis, config)

	if tree == nil {
		t.Fatal("Expected non-nil tree")
	}

	if tree.analysis == nil {
		t.Error("Expected analysis to be set")
	}

	if tree.expanded == nil {
		t.Error("Expected expanded map to be initialized")
	}
}

// TestDependencyTree_Render_Empty tests rendering with no data.
func TestDependencyTree_Render_Empty(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
	}

	config := DefaultDependencyTreeConfig()
	tree := NewDependencyTree(analysis, config)

	result := tree.Render()

	if result == "" {
		t.Error("Expected non-empty result")
	}

	if !strings.Contains(result, "No dependency data") {
		t.Error("Expected 'No dependency data' message")
	}
}

// TestDependencyTree_Render_SingleHole tests rendering with single hole.
func TestDependencyTree_Render_SingleHole(t *testing.T) {
	hole := semantic.EnhancedTypedHole{
		TypedHole:  semantic.TypedHole{Type: "AuthService", Constraint: "thread-safe"},
		Priority:   8,
		Complexity: 5,
	}

	// Build simple tree
	node := &analyze.HoleNode{
		ID:   "AuthService_0",
		Hole: hole,
	}

	root := &analyze.HoleNode{
		ID:       "root",
		Children: []*analyze.HoleNode{node},
	}

	analysis := &analyze.HoleAnalysis{
		Holes:           []semantic.EnhancedTypedHole{hole},
		DependencyTree:  root,
		TotalComplexity: 5,
		AveragePriority: 8.0,
	}

	config := DefaultDependencyTreeConfig()
	tree := NewDependencyTree(analysis, config)

	result := tree.Render()

	if !strings.Contains(result, "AuthService") {
		t.Error("Expected to find 'AuthService' in output")
	}

	if !strings.Contains(result, "Total Holes") {
		t.Error("Expected summary section")
	}
}

// TestDependencyTree_Render_WithChildren tests rendering with hierarchy.
func TestDependencyTree_Render_WithChildren(t *testing.T) {
	parentHole := semantic.EnhancedTypedHole{
		TypedHole:  semantic.TypedHole{Type: "Parent"},
		Priority:   7,
		Complexity: 3,
	}

	childHole := semantic.EnhancedTypedHole{
		TypedHole:  semantic.TypedHole{Type: "Child"},
		Priority:   5,
		Complexity: 4,
	}

	childNode := &analyze.HoleNode{
		ID:   "Child_1",
		Hole: childHole,
	}

	parentNode := &analyze.HoleNode{
		ID:       "Parent_0",
		Hole:     parentHole,
		Children: []*analyze.HoleNode{childNode},
	}

	childNode.Parents = []*analyze.HoleNode{parentNode}

	root := &analyze.HoleNode{
		ID:       "root",
		Children: []*analyze.HoleNode{parentNode},
	}

	analysis := &analyze.HoleAnalysis{
		Holes:           []semantic.EnhancedTypedHole{parentHole, childHole},
		DependencyTree:  root,
		TotalComplexity: 7,
		AveragePriority: 6.0,
	}

	config := DefaultDependencyTreeConfig()
	tree := NewDependencyTree(analysis, config)

	result := tree.Render()

	if !strings.Contains(result, "Parent") {
		t.Error("Expected to find 'Parent' in output")
	}

	if !strings.Contains(result, "Child") {
		t.Error("Expected to find 'Child' in output")
	}

	// Check for tree characters
	if !strings.Contains(result, "└──") && !strings.Contains(result, "├──") {
		t.Error("Expected tree branch characters")
	}
}

// TestDependencyTree_FormatNodeContent tests node formatting.
func TestDependencyTree_FormatNodeContent(t *testing.T) {
	hole := semantic.EnhancedTypedHole{
		TypedHole:  semantic.TypedHole{Type: "Service", Constraint: "async"},
		Priority:   9,
		Complexity: 6,
	}

	node := &analyze.HoleNode{
		ID:   "Service_0",
		Hole: hole,
	}

	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{hole},
	}

	config := DefaultDependencyTreeConfig()
	config.ShowPriority = true
	config.ShowComplexity = true
	config.ShowConstraints = true

	tree := NewDependencyTree(analysis, config)
	content := tree.formatNodeContent(node)

	if !strings.Contains(content, "Service") {
		t.Error("Expected service name in content")
	}

	// Check for metadata (may be formatted/styled)
	if !strings.Contains(content, "P:9") && !strings.Contains(content, "9") {
		t.Error("Expected priority in content")
	}

	if !strings.Contains(content, "C:6") && !strings.Contains(content, "6") {
		t.Error("Expected complexity in content")
	}
}

// TestDependencyTree_GetPriorityColor tests color assignment.
func TestDependencyTree_GetPriorityColor(t *testing.T) {
	config := DefaultDependencyTreeConfig()
	tree := &DependencyTree{config: config}

	// High priority - should be red
	highColor := tree.getPriorityColor(9)
	if highColor != "196" {
		t.Errorf("Expected high priority to be red (196), got %s", highColor)
	}

	// Medium priority - should be yellow
	medColor := tree.getPriorityColor(6)
	if medColor != "226" {
		t.Errorf("Expected medium priority to be yellow (226), got %s", medColor)
	}

	// Low priority - should be green
	lowColor := tree.getPriorityColor(3)
	if lowColor != "34" {
		t.Errorf("Expected low priority to be green (34), got %s", lowColor)
	}
}

// TestDependencyTree_ToggleExpand tests expansion toggling.
func TestDependencyTree_ToggleExpand(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
	}

	config := DefaultDependencyTreeConfig()
	config.ExpandAll = false
	tree := NewDependencyTree(analysis, config)

	nodeID := "test_0"

	// Initially not expanded
	if tree.expanded[nodeID] {
		t.Error("Expected node to start collapsed")
	}

	// Toggle to expanded
	tree.ToggleExpand(nodeID)
	if !tree.expanded[nodeID] {
		t.Error("Expected node to be expanded after toggle")
	}

	// Toggle back to collapsed
	tree.ToggleExpand(nodeID)
	if tree.expanded[nodeID] {
		t.Error("Expected node to be collapsed after second toggle")
	}
}

// TestDependencyTree_ExpandAll tests expanding all nodes.
func TestDependencyTree_ExpandAll(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "A"}},
		{TypedHole: semantic.TypedHole{Type: "B"}},
	}

	analysis := &analyze.HoleAnalysis{
		Holes: holes,
	}

	config := DefaultDependencyTreeConfig()
	config.ExpandAll = false
	tree := NewDependencyTree(analysis, config)

	tree.ExpandAll()

	if !tree.config.ExpandAll {
		t.Error("Expected ExpandAll to be true")
	}
}

// TestDependencyTree_CollapseAll tests collapsing all nodes.
func TestDependencyTree_CollapseAll(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "A"}},
		},
	}

	config := DefaultDependencyTreeConfig()
	config.ExpandAll = true
	tree := NewDependencyTree(analysis, config)

	tree.CollapseAll()

	if tree.config.ExpandAll {
		t.Error("Expected ExpandAll to be false")
	}

	// Check that expanded map is cleared
	for _, expanded := range tree.expanded {
		if expanded {
			t.Error("Expected all nodes to be collapsed")
		}
	}
}

// TestDependencyTree_RenderSummary tests summary generation.
func TestDependencyTree_RenderSummary(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{Complexity: 5},
			{Complexity: 7},
		},
		TotalComplexity: 12,
		AveragePriority: 7.5,
		CriticalPath: []semantic.EnhancedTypedHole{
			{Complexity: 7},
		},
	}

	config := DefaultDependencyTreeConfig()
	tree := NewDependencyTree(analysis, config)

	summary := tree.renderSummary()

	if !strings.Contains(summary, "Total Holes: 2") {
		t.Error("Expected total holes count")
	}

	if !strings.Contains(summary, "Total Complexity: 12") {
		t.Error("Expected total complexity")
	}

	if !strings.Contains(summary, "Average Priority: 7.5") {
		t.Error("Expected average priority")
	}

	if !strings.Contains(summary, "Critical Path") {
		t.Error("Expected critical path info")
	}
}

// TestDependencyTree_RenderSummary_WithCircular tests summary with circular deps.
func TestDependencyTree_RenderSummary_WithCircular(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes:           []semantic.EnhancedTypedHole{},
		TotalComplexity: 0,
		AveragePriority: 0,
		CircularDeps: [][]string{
			{"A", "B", "C"},
		},
	}

	config := DefaultDependencyTreeConfig()
	tree := NewDependencyTree(analysis, config)

	summary := tree.renderSummary()

	if !strings.Contains(summary, "WARNING") {
		t.Error("Expected warning for circular dependencies")
	}

	if !strings.Contains(summary, "circular") {
		t.Error("Expected circular dependency message")
	}
}

// TestDependencyTree_RenderCompact tests compact rendering.
func TestDependencyTree_RenderCompact(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
		ImplementOrder: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "First"}, Priority: 8, Complexity: 3},
			{TypedHole: semantic.TypedHole{Type: "Second"}, Priority: 6, Complexity: 5},
		},
	}

	config := DefaultDependencyTreeConfig()
	tree := NewDependencyTree(analysis, config)

	result := tree.RenderCompact()

	if !strings.Contains(result, "First") {
		t.Error("Expected 'First' in compact view")
	}

	if !strings.Contains(result, "Second") {
		t.Error("Expected 'Second' in compact view")
	}

	if !strings.Contains(result, "→") {
		t.Error("Expected arrow separator")
	}
}

// TestRenderDependencyMatrix tests matrix visualization.
func TestRenderDependencyMatrix(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "A"}},
			{TypedHole: semantic.TypedHole{Type: "B"}},
		},
		Dependencies: []analyze.HoleDependency{
			{From: "A_0", To: "B_1", Relationship: "requires"},
		},
	}

	result := RenderDependencyMatrix(analysis)

	if !strings.Contains(result, "Dependency Matrix") {
		t.Error("Expected matrix header")
	}

	if !strings.Contains(result, "Legend") {
		t.Error("Expected legend")
	}
}

// TestRenderDependencyMatrix_Empty tests matrix with no data.
func TestRenderDependencyMatrix_Empty(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
	}

	result := RenderDependencyMatrix(analysis)

	if !strings.Contains(result, "No dependency data") {
		t.Error("Expected 'No dependency data' message")
	}
}

// TestRenderCircularDependencies tests circular dependency rendering.
func TestRenderCircularDependencies(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		CircularDeps: [][]string{
			{"ServiceA_0", "ServiceB_1", "ServiceC_2"},
		},
	}

	result := RenderCircularDependencies(analysis)

	if !strings.Contains(result, "Warning") && !strings.Contains(result, "WARNING") {
		t.Error("Expected warning message")
	}

	if !strings.Contains(result, "Circular") {
		t.Error("Expected circular dependency mention")
	}

	if !strings.Contains(result, "ServiceA") {
		t.Error("Expected service names in cycle")
	}
}

// TestRenderCircularDependencies_None tests with no circular deps.
func TestRenderCircularDependencies_None(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		CircularDeps: [][]string{},
	}

	result := RenderCircularDependencies(analysis)

	if !strings.Contains(result, "No circular") {
		t.Error("Expected success message for no circular deps")
	}
}

// TestDefaultDependencyTreeConfig tests default config.
func TestDefaultDependencyTreeConfig(t *testing.T) {
	config := DefaultDependencyTreeConfig()

	if config.Width <= 0 {
		t.Error("Expected positive width")
	}

	if config.Height <= 0 {
		t.Error("Expected positive height")
	}

	if !config.ShowComplexity {
		t.Error("Expected complexity to be shown by default")
	}

	if !config.ShowPriority {
		t.Error("Expected priority to be shown by default")
	}
}
