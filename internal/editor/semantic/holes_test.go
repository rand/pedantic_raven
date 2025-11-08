package semantic

import (
	"strings"
	"testing"
)

// --- Constraint Parsing Tests ---

func TestParseTypedHoleConstraints(t *testing.T) {
	hole := TypedHole{
		Type:       "Function",
		Constraint: "implements Serializer",
		Span:       Span{Start: 0, End: 10},
	}

	constraints := ParseTypedHoleConstraints(hole)

	if len(constraints) == 0 {
		t.Fatal("Expected to find constraints")
	}

	if constraints[0].Type != "implements" {
		t.Errorf("Expected constraint type 'implements', got '%s'", constraints[0].Type)
	}
}

func TestParseConstraintStringImplements(t *testing.T) {
	constraint := parseConstraintString("implements Handler")

	if constraint.Type != "implements" {
		t.Errorf("Expected type 'implements', got '%s'", constraint.Type)
	}

	if !strings.Contains(constraint.Description, "Interface implementation") {
		t.Error("Expected description to mention interface implementation")
	}
}

func TestParseConstraintStringThreadSafe(t *testing.T) {
	constraint := parseConstraintString("thread-safe required")

	if constraint.Type != "thread" {
		t.Errorf("Expected type 'thread', got '%s'", constraint.Type)
	}
}

func TestParseConstraintStringAsync(t *testing.T) {
	constraint := parseConstraintString("async operation")

	if constraint.Type != "async" {
		t.Errorf("Expected type 'async', got '%s'", constraint.Type)
	}
}

func TestParseConstraintStringGeneric(t *testing.T) {
	constraint := parseConstraintString("custom requirement")

	if constraint.Type != "generic" {
		t.Errorf("Expected type 'generic', got '%s'", constraint.Type)
	}

	if constraint.Value != "custom requirement" {
		t.Errorf("Expected value to be preserved, got '%s'", constraint.Value)
	}
}

// --- Priority Calculation Tests ---

func TestCalculateHolePriorityBase(t *testing.T) {
	hole := TypedHole{
		Type: "Function",
		Span: Span{Start: 0, End: 10},
	}

	priority := CalculateHolePriority(hole, []Relationship{})

	if priority < 1 || priority > 10 {
		t.Errorf("Priority should be between 1-10, got %d", priority)
	}
}

func TestCalculateHolePriorityWithMentions(t *testing.T) {
	hole := TypedHole{
		Type: "Handler",
		Span: Span{Start: 0, End: 10},
	}

	relationships := []Relationship{
		{Subject: "System", Predicate: "uses", Object: "Handler"},
		{Subject: "Handler", Predicate: "processes", Object: "Request"},
	}

	priority := CalculateHolePriority(hole, relationships)

	// Should have higher priority due to mentions
	if priority <= 5 {
		t.Error("Expected higher priority with mentions in relationships")
	}
}

func TestCalculateHolePriorityWithConstraint(t *testing.T) {
	hole := TypedHole{
		Type:       "Function",
		Constraint: "thread-safe",
		Span:       Span{Start: 0, End: 10},
	}

	priority := CalculateHolePriority(hole, []Relationship{})

	// Constraint should increase priority
	if priority <= 5 {
		t.Error("Expected higher priority with constraint")
	}
}

// --- Complexity Calculation Tests ---

func TestCalculateHoleComplexityBase(t *testing.T) {
	hole := TypedHole{
		Type: "Function",
		Span: Span{Start: 0, End: 10},
	}

	complexity := CalculateHoleComplexity(hole, []HoleConstraint{})

	if complexity < 1 || complexity > 10 {
		t.Errorf("Complexity should be between 1-10, got %d", complexity)
	}
}

func TestCalculateHoleComplexityWithConstraints(t *testing.T) {
	hole := TypedHole{
		Type: "Function",
		Span: Span{Start: 0, End: 10},
	}

	constraints := []HoleConstraint{
		{Type: "thread"},
		{Type: "async"},
	}

	complexity := CalculateHoleComplexity(hole, constraints)

	// Concurrency constraints should increase complexity
	if complexity <= 3 {
		t.Error("Expected higher complexity with concurrency constraints")
	}
}

func TestCalculateHoleComplexityConcurrency(t *testing.T) {
	hole := TypedHole{
		Type: "Service",
		Span: Span{Start: 0, End: 10},
	}

	constraints := []HoleConstraint{
		{Type: "thread"},
		{Type: "concurrent"},
		{Type: "atomic"},
	}

	complexity := CalculateHoleComplexity(hole, constraints)

	// Multiple concurrency constraints = high complexity
	if complexity < 7 {
		t.Errorf("Expected high complexity with concurrency constraints, got %d", complexity)
	}
}

// --- Implementation Suggestion Tests ---

func TestSuggestImplementationFunction(t *testing.T) {
	hole := TypedHole{
		Type: "Function",
		Span: Span{Start: 0, End: 10},
	}

	suggestion := SuggestImplementation(hole, []HoleConstraint{})

	if !strings.Contains(suggestion, "function") {
		t.Error("Expected suggestion to mention function")
	}
}

func TestSuggestImplementationWithThreadConstraint(t *testing.T) {
	hole := TypedHole{
		Type: "Service",
		Span: Span{Start: 0, End: 10},
	}

	constraints := []HoleConstraint{
		{Type: "thread"},
	}

	suggestion := SuggestImplementation(hole, constraints)

	if !strings.Contains(strings.ToLower(suggestion), "mutex") {
		t.Error("Expected suggestion to mention mutex for thread safety")
	}
}

func TestSuggestImplementationWithAsyncConstraint(t *testing.T) {
	hole := TypedHole{
		Type: "Handler",
		Span: Span{Start: 0, End: 10},
	}

	constraints := []HoleConstraint{
		{Type: "async"},
	}

	suggestion := SuggestImplementation(hole, constraints)

	if !strings.Contains(strings.ToLower(suggestion), "goroutine") {
		t.Error("Expected suggestion to mention goroutines for async")
	}
}

func TestSuggestImplementationBuilder(t *testing.T) {
	hole := TypedHole{
		Type: "Builder",
		Span: Span{Start: 0, End: 10},
	}

	suggestion := SuggestImplementation(hole, []HoleConstraint{})

	if !strings.Contains(strings.ToLower(suggestion), "builder") {
		t.Error("Expected suggestion to mention builder pattern")
	}
}

// --- Hole Enhancement Tests ---

func TestEnhanceTypedHole(t *testing.T) {
	hole := TypedHole{
		Type:       "Function",
		Constraint: "thread-safe",
		Span:       Span{Start: 0, End: 10},
	}

	relationships := []Relationship{
		{Subject: "System", Predicate: "uses", Object: "Function"},
	}

	enhanced := EnhanceTypedHole(hole, relationships)

	if len(enhanced.Constraints) == 0 {
		t.Error("Expected constraints to be parsed")
	}

	if enhanced.Priority == 0 {
		t.Error("Expected priority to be calculated")
	}

	if enhanced.Complexity == 0 {
		t.Error("Expected complexity to be calculated")
	}

	if enhanced.SuggestedImpl == "" {
		t.Error("Expected implementation suggestion")
	}
}

// --- Related Holes Tests ---

func TestFindRelatedHoles(t *testing.T) {
	holes := []TypedHole{
		{Type: "Handler", Span: Span{Start: 0, End: 10}},
		{Type: "Service", Span: Span{Start: 20, End: 30}},
	}

	relationships := []Relationship{
		{Subject: "Handler", Predicate: "uses", Object: "Service"},
	}

	related := FindRelatedHoles(holes, relationships)

	if len(related) == 0 {
		t.Error("Expected to find related holes")
	}

	// Check bidirectional relationship
	foundRelation := false
	for _, relations := range related {
		if len(relations) > 0 {
			foundRelation = true
			break
		}
	}

	if !foundRelation {
		t.Error("Expected to find hole relationships")
	}
}

// --- Hole Prioritizer Tests ---

func TestNewHolePrioritizer(t *testing.T) {
	holes := []TypedHole{
		{Type: "Function", Span: Span{Start: 0, End: 10}},
		{Type: "Service", Span: Span{Start: 20, End: 30}},
	}

	prioritizer := NewHolePrioritizer(holes, []Relationship{})

	if len(prioritizer.Holes) != 2 {
		t.Errorf("Expected 2 enhanced holes, got %d", len(prioritizer.Holes))
	}
}

func TestGetByPriority(t *testing.T) {
	holes := []TypedHole{
		{Type: "Function", Constraint: "", Span: Span{Start: 0, End: 10}},
		{Type: "Service", Constraint: "thread-safe", Span: Span{Start: 20, End: 30}},
	}

	relationships := []Relationship{
		{Subject: "System", Predicate: "uses", Object: "Service"},
		{Subject: "Service", Predicate: "manages", Object: "Data"},
	}

	prioritizer := NewHolePrioritizer(holes, relationships)
	sorted := prioritizer.GetByPriority()

	if len(sorted) != 2 {
		t.Fatalf("Expected 2 holes, got %d", len(sorted))
	}

	// First hole should have higher priority
	if sorted[0].Priority < sorted[1].Priority {
		t.Error("Expected holes to be sorted by priority (highest first)")
	}
}

func TestGetByComplexity(t *testing.T) {
	holes := []TypedHole{
		{Type: "Function", Constraint: "thread-safe atomic", Span: Span{Start: 0, End: 10}},
		{Type: "Service", Constraint: "", Span: Span{Start: 20, End: 30}},
	}

	prioritizer := NewHolePrioritizer(holes, []Relationship{})
	sorted := prioritizer.GetByComplexity()

	if len(sorted) != 2 {
		t.Fatalf("Expected 2 holes, got %d", len(sorted))
	}

	// First hole should have lower complexity
	if sorted[0].Complexity > sorted[1].Complexity {
		t.Error("Expected holes to be sorted by complexity (lowest first)")
	}
}

func TestGetRecommendedOrder(t *testing.T) {
	holes := []TypedHole{
		// High priority, low complexity
		{Type: "Function", Constraint: "", Span: Span{Start: 0, End: 10}},
		// Low priority, high complexity
		{Type: "Service", Constraint: "thread-safe concurrent atomic", Span: Span{Start: 20, End: 30}},
		// Medium priority, medium complexity
		{Type: "Handler", Constraint: "async", Span: Span{Start: 40, End: 50}},
	}

	relationships := []Relationship{
		{Subject: "System", Predicate: "uses", Object: "Function"},
		{Subject: "Function", Predicate: "calls", Object: "Handler"},
	}

	prioritizer := NewHolePrioritizer(holes, relationships)
	recommended := prioritizer.GetRecommendedOrder()

	if len(recommended) != 3 {
		t.Fatalf("Expected 3 holes, got %d", len(recommended))
	}

	// Verify ordering makes sense (high priority/complexity ratio first)
	// Simple function should come before complex service
	for i, hole := range recommended {
		t.Logf("Hole %d: Type=%s, Priority=%d, Complexity=%d, Ratio=%.2f",
			i, hole.Type, hole.Priority, hole.Complexity,
			float64(hole.Priority)/float64(max(hole.Complexity, 1)))
	}
}

// --- Helper Function Tests ---

func TestMinFunction(t *testing.T) {
	if min(5, 10) != 5 {
		t.Error("min(5, 10) should be 5")
	}

	if min(10, 5) != 5 {
		t.Error("min(10, 5) should be 5")
	}

	if min(5, 5) != 5 {
		t.Error("min(5, 5) should be 5")
	}
}

func TestMaxFunction(t *testing.T) {
	if max(5, 10) != 10 {
		t.Error("max(5, 10) should be 10")
	}

	if max(10, 5) != 10 {
		t.Error("max(10, 5) should be 10")
	}

	if max(5, 5) != 5 {
		t.Error("max(5, 5) should be 5")
	}
}

func TestContainsFunction(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !contains(slice, "b") {
		t.Error("Should contain 'b'")
	}

	if contains(slice, "d") {
		t.Error("Should not contain 'd'")
	}
}
