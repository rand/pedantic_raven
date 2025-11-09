package analyze

import (
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// TestAnalyzeTypedHoles_Empty tests with no holes.
func TestAnalyzeTypedHoles_Empty(t *testing.T) {
	analysis := &semantic.Analysis{
		TypedHoles:    []semantic.TypedHole{},
		Relationships: []semantic.Relationship{},
	}

	result := AnalyzeTypedHoles(analysis)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Holes) != 0 {
		t.Errorf("Expected 0 holes, got %d", len(result.Holes))
	}

	if len(result.Dependencies) != 0 {
		t.Errorf("Expected 0 dependencies, got %d", len(result.Dependencies))
	}

	if result.TotalComplexity != 0 {
		t.Errorf("Expected 0 total complexity, got %d", result.TotalComplexity)
	}

	if result.AveragePriority != 0 {
		t.Errorf("Expected 0 average priority, got %f", result.AveragePriority)
	}
}

// TestAnalyzeTypedHoles_Nil tests with nil analysis.
func TestAnalyzeTypedHoles_Nil(t *testing.T) {
	result := AnalyzeTypedHoles(nil)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result.Holes) != 0 {
		t.Errorf("Expected 0 holes, got %d", len(result.Holes))
	}
}

// TestAnalyzeTypedHoles_SingleHole tests with a single hole.
func TestAnalyzeTypedHoles_SingleHole(t *testing.T) {
	analysis := &semantic.Analysis{
		TypedHoles: []semantic.TypedHole{
			{
				Type:       "AuthService",
				Constraint: "thread-safe",
			},
		},
		Relationships: []semantic.Relationship{},
	}

	result := AnalyzeTypedHoles(analysis)

	if len(result.Holes) != 1 {
		t.Fatalf("Expected 1 hole, got %d", len(result.Holes))
	}

	hole := result.Holes[0]
	if hole.Type != "AuthService" {
		t.Errorf("Expected type AuthService, got %s", hole.Type)
	}

	if hole.Priority == 0 {
		t.Error("Expected non-zero priority")
	}

	if hole.Complexity == 0 {
		t.Error("Expected non-zero complexity")
	}

	if result.TotalComplexity == 0 {
		t.Error("Expected non-zero total complexity")
	}

	if result.AveragePriority == 0 {
		t.Error("Expected non-zero average priority")
	}
}

// TestAnalyzeTypedHoles_MultipleHoles tests with multiple holes.
func TestAnalyzeTypedHoles_MultipleHoles(t *testing.T) {
	analysis := &semantic.Analysis{
		TypedHoles: []semantic.TypedHole{
			{Type: "AuthService", Constraint: "thread-safe"},
			{Type: "DatabaseLayer", Constraint: "atomic"},
			{Type: "ConfigLoader", Constraint: "immutable"},
		},
		Relationships: []semantic.Relationship{
			{
				Subject:   "AuthService",
				Predicate: "requires",
				Object:    "DatabaseLayer",
			},
			{
				Subject:   "DatabaseLayer",
				Predicate: "requires",
				Object:    "ConfigLoader",
			},
		},
	}

	result := AnalyzeTypedHoles(analysis)

	if len(result.Holes) != 3 {
		t.Fatalf("Expected 3 holes, got %d", len(result.Holes))
	}

	// Check that dependencies were built
	if len(result.Dependencies) == 0 {
		t.Error("Expected dependencies to be built")
	}

	// Check implementation order
	if len(result.ImplementOrder) != 3 {
		t.Errorf("Expected 3 holes in implementation order, got %d", len(result.ImplementOrder))
	}

	// Verify topological order: ConfigLoader should come before DatabaseLayer
	orderMap := make(map[string]int)
	for i, hole := range result.ImplementOrder {
		orderMap[hole.Type] = i
	}

	if orderMap["ConfigLoader"] > orderMap["DatabaseLayer"] {
		t.Error("Expected ConfigLoader before DatabaseLayer in implementation order")
	}

	if orderMap["DatabaseLayer"] > orderMap["AuthService"] {
		t.Error("Expected DatabaseLayer before AuthService in implementation order")
	}
}

// TestBuildDependencies tests dependency extraction.
func TestBuildDependencies(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "ServiceA"}},
		{TypedHole: semantic.TypedHole{Type: "ServiceB"}},
	}

	relationships := []semantic.Relationship{
		{
			Subject:   "ServiceA",
			Predicate: "requires",
			Object:    "ServiceB",
		},
	}

	deps := buildDependencies(holes, relationships)

	if len(deps) == 0 {
		t.Fatal("Expected dependencies to be found")
	}

	found := false
	for _, dep := range deps {
		if dep.Relationship == "requires" {
			found = true
			if dep.Strength < 5 {
				t.Errorf("Expected high strength for 'requires', got %d", dep.Strength)
			}
		}
	}

	if !found {
		t.Error("Expected to find 'requires' dependency")
	}
}

// TestBuildDependencies_Extends tests extends relationship.
func TestBuildDependencies_Extends(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "BaseService"}},
		{TypedHole: semantic.TypedHole{Type: "ExtendedService"}},
	}

	relationships := []semantic.Relationship{
		{
			Subject:   "ExtendedService",
			Predicate: "extends",
			Object:    "BaseService",
		},
	}

	deps := buildDependencies(holes, relationships)

	if len(deps) == 0 {
		t.Fatal("Expected dependencies to be found")
	}

	found := false
	for _, dep := range deps {
		if dep.Relationship == "extends" {
			found = true
		}
	}

	if !found {
		t.Error("Expected to find 'extends' dependency")
	}
}

// TestBuildDependencies_Implements tests implements relationship.
func TestBuildDependencies_Implements(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "Interface"}},
		{TypedHole: semantic.TypedHole{Type: "Implementation"}},
	}

	relationships := []semantic.Relationship{
		{
			Subject:   "Implementation",
			Predicate: "implements",
			Object:    "Interface",
		},
	}

	deps := buildDependencies(holes, relationships)

	if len(deps) == 0 {
		t.Fatal("Expected dependencies to be found")
	}

	found := false
	for _, dep := range deps {
		if dep.Relationship == "implements" {
			found = true
		}
	}

	if !found {
		t.Error("Expected to find 'implements' dependency")
	}
}

// TestBuildDependencyTree tests tree construction.
func TestBuildDependencyTree(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "Root"}},
		{TypedHole: semantic.TypedHole{Type: "Child"}},
	}

	deps := []HoleDependency{
		{
			From:         "Child_1",
			To:           "Root_0",
			Relationship: "requires",
			Strength:     8,
		},
	}

	tree := buildDependencyTree(holes, deps)

	if tree == nil {
		t.Fatal("Expected non-nil tree")
	}

	if tree.ID != "root" {
		t.Errorf("Expected root node ID to be 'root', got %s", tree.ID)
	}

	// Root should have children
	if len(tree.Children) == 0 {
		t.Error("Expected root to have children")
	}
}

// TestDetectCircularDependencies_NoCycle tests no circular dependencies.
func TestDetectCircularDependencies_NoCycle(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "A"}},
		{TypedHole: semantic.TypedHole{Type: "B"}},
		{TypedHole: semantic.TypedHole{Type: "C"}},
	}

	deps := []HoleDependency{
		{From: "A_0", To: "B_1", Relationship: "requires"},
		{From: "B_1", To: "C_2", Relationship: "requires"},
	}

	cycles := detectCircularDependencies(holes, deps)

	if len(cycles) != 0 {
		t.Errorf("Expected 0 cycles, got %d", len(cycles))
	}
}

// TestDetectCircularDependencies_WithCycle tests circular dependency detection.
func TestDetectCircularDependencies_WithCycle(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "A"}},
		{TypedHole: semantic.TypedHole{Type: "B"}},
		{TypedHole: semantic.TypedHole{Type: "C"}},
	}

	deps := []HoleDependency{
		{From: "A_0", To: "B_1", Relationship: "requires"},
		{From: "B_1", To: "C_2", Relationship: "requires"},
		{From: "C_2", To: "A_0", Relationship: "requires"},
	}

	cycles := detectCircularDependencies(holes, deps)

	if len(cycles) == 0 {
		t.Error("Expected to find circular dependency")
	}
}

// TestCalculateImplementationOrder tests topological sort.
func TestCalculateImplementationOrder(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{
			TypedHole:  semantic.TypedHole{Type: "Config"},
			Priority:   8,
			Complexity: 3,
		},
		{
			TypedHole:  semantic.TypedHole{Type: "Database"},
			Priority:   9,
			Complexity: 7,
		},
		{
			TypedHole:  semantic.TypedHole{Type: "Auth"},
			Priority:   10,
			Complexity: 8,
		},
	}

	deps := []HoleDependency{
		{From: "Database_1", To: "Config_0", Relationship: "requires"},
		{From: "Auth_2", To: "Database_1", Relationship: "requires"},
	}

	order := calculateImplementationOrder(holes, deps)

	if len(order) != 3 {
		t.Fatalf("Expected 3 holes in order, got %d", len(order))
	}

	// Config should come first (no dependencies)
	if order[0].Type != "Config" {
		t.Errorf("Expected Config first, got %s", order[0].Type)
	}

	// Database should come before Auth
	dbIdx := -1
	authIdx := -1
	for i, hole := range order {
		if hole.Type == "Database" {
			dbIdx = i
		}
		if hole.Type == "Auth" {
			authIdx = i
		}
	}

	if dbIdx > authIdx {
		t.Error("Expected Database before Auth in implementation order")
	}
}

// TestIdentifyCriticalPath tests critical path identification.
func TestIdentifyCriticalPath(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{
			TypedHole:  semantic.TypedHole{Type: "Simple"},
			Complexity: 2,
		},
		{
			TypedHole:  semantic.TypedHole{Type: "Complex1"},
			Complexity: 8,
		},
		{
			TypedHole:  semantic.TypedHole{Type: "Complex2"},
			Complexity: 7,
		},
	}

	deps := []HoleDependency{
		{From: "Complex1_1", To: "Simple_0", Relationship: "requires"},
		{From: "Complex2_2", To: "Complex1_1", Relationship: "requires"},
	}

	tree := buildDependencyTree(holes, deps)
	order := calculateImplementationOrder(holes, deps)

	criticalPath := identifyCriticalPath(tree, order)

	// Critical path should prioritize higher complexity
	if len(criticalPath) == 0 {
		t.Error("Expected non-empty critical path")
	}

	// Check that critical path includes complex holes
	totalComplexity := calculatePathComplexity(criticalPath)
	if totalComplexity == 0 {
		t.Error("Expected non-zero critical path complexity")
	}
}

// TestGenerateImplementationRoadmap tests roadmap generation.
func TestGenerateImplementationRoadmap(t *testing.T) {
	analysis := &HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{
				TypedHole:  semantic.TypedHole{Type: "ServiceA", Constraint: "thread-safe"},
				Priority:   8,
				Complexity: 5,
			},
			{
				TypedHole:  semantic.TypedHole{Type: "ServiceB"},
				Priority:   6,
				Complexity: 3,
			},
		},
		ImplementOrder:  []semantic.EnhancedTypedHole{},
		CriticalPath:    []semantic.EnhancedTypedHole{},
		TotalComplexity: 8,
		AveragePriority: 7.0,
		CircularDeps:    [][]string{},
	}

	roadmap := GenerateImplementationRoadmap(analysis)

	if roadmap == "" {
		t.Error("Expected non-empty roadmap")
	}

	// Check for key sections
	if !strings.Contains(roadmap, "Total Holes") {
		t.Error("Expected 'Total Holes' in roadmap")
	}

	if !strings.Contains(roadmap, "Total Complexity") {
		t.Error("Expected 'Total Complexity' in roadmap")
	}

	if !strings.Contains(roadmap, "Average Priority") {
		t.Error("Expected 'Average Priority' in roadmap")
	}
}

// TestGenerateImplementationRoadmap_WithCircular tests roadmap with circular deps.
func TestGenerateImplementationRoadmap_WithCircular(t *testing.T) {
	analysis := &HoleAnalysis{
		Holes:           []semantic.EnhancedTypedHole{},
		ImplementOrder:  []semantic.EnhancedTypedHole{},
		TotalComplexity: 0,
		AveragePriority: 0,
		CircularDeps: [][]string{
			{"A_0", "B_1", "C_2"},
		},
	}

	roadmap := GenerateImplementationRoadmap(analysis)

	if !strings.Contains(roadmap, "WARNING") {
		t.Error("Expected WARNING for circular dependencies")
	}

	if !strings.Contains(roadmap, "Circular dependencies") {
		t.Error("Expected circular dependency message")
	}
}

// TestGroupIntoMilestones tests milestone grouping.
func TestGroupIntoMilestones(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "A"}, Complexity: 5},
		{TypedHole: semantic.TypedHole{Type: "B"}, Complexity: 5},
		{TypedHole: semantic.TypedHole{Type: "C"}, Complexity: 5},
		{TypedHole: semantic.TypedHole{Type: "D"}, Complexity: 5},
		{TypedHole: semantic.TypedHole{Type: "E"}, Complexity: 5},
		{TypedHole: semantic.TypedHole{Type: "F"}, Complexity: 5},
	}

	milestones := groupIntoMilestones(holes)

	if len(milestones) == 0 {
		t.Fatal("Expected at least one milestone")
	}

	// Verify all holes are included
	totalHoles := 0
	for _, milestone := range milestones {
		totalHoles += len(milestone)
	}

	if totalHoles != len(holes) {
		t.Errorf("Expected %d holes total, got %d", len(holes), totalHoles)
	}
}

// TestGroupIntoMilestones_ByComplexity tests grouping by complexity threshold.
func TestGroupIntoMilestones_ByComplexity(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{
		{TypedHole: semantic.TypedHole{Type: "HighComplex"}, Complexity: 30},
		{TypedHole: semantic.TypedHole{Type: "LowComplex"}, Complexity: 2},
	}

	milestones := groupIntoMilestones(holes)

	// High complexity hole should trigger new milestone
	if len(milestones) < 2 {
		t.Error("Expected high complexity to trigger separate milestones")
	}
}

// TestHasDependency tests dependency checking helper.
func TestHasDependency(t *testing.T) {
	deps := []HoleDependency{
		{From: "A", To: "B", Relationship: "requires"},
		{From: "B", To: "C", Relationship: "extends"},
	}

	if !hasDependency(deps, "A", "B") {
		t.Error("Expected to find dependency A->B")
	}

	if hasDependency(deps, "C", "A") {
		t.Error("Expected not to find dependency C->A")
	}
}

// TestContains tests contains helper.
func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !contains(slice, "b") {
		t.Error("Expected to find 'b' in slice")
	}

	if contains(slice, "d") {
		t.Error("Expected not to find 'd' in slice")
	}
}

// TestCalculatePathComplexity tests complexity calculation.
func TestCalculatePathComplexity(t *testing.T) {
	path := []semantic.EnhancedTypedHole{
		{Complexity: 5},
		{Complexity: 7},
		{Complexity: 3},
	}

	total := calculatePathComplexity(path)

	expected := 15
	if total != expected {
		t.Errorf("Expected complexity %d, got %d", expected, total)
	}
}

// TestCalculatePathComplexity_Empty tests with empty path.
func TestCalculatePathComplexity_Empty(t *testing.T) {
	path := []semantic.EnhancedTypedHole{}

	total := calculatePathComplexity(path)

	if total != 0 {
		t.Errorf("Expected complexity 0, got %d", total)
	}
}
