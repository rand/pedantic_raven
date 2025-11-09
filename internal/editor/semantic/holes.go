package semantic

import (
	"fmt"
	"strings"
)

// HoleConstraint represents a constraint on a typed hole.
type HoleConstraint struct {
	Type        string   // Constraint type (e.g., "implements", "extends", "requires")
	Value       string   // Constraint value
	Description string   // Human-readable description
	Satisfied   bool     // Whether constraint is satisfied
	Dependencies []string // Dependencies required to satisfy
}

// EnhancedTypedHole extends TypedHole with richer metadata.
type EnhancedTypedHole struct {
	TypedHole
	Constraints    []HoleConstraint // Parsed constraints
	Priority       int              // Implementation priority (0-10)
	Complexity     int              // Estimated complexity (1-10)
	SuggestedImpl  string           // Suggested implementation approach
	RelatedHoles   []string         // IDs of related holes
	Dependencies   []string         // What this hole depends on
}

// ParseTypedHoleConstraints extracts constraints from a typed hole.
func ParseTypedHoleConstraints(hole TypedHole) []HoleConstraint {
	var constraints []HoleConstraint

	// Parse constraint from the constraint field
	if hole.Constraint != "" {
		constraint := parseConstraintString(hole.Constraint)
		constraints = append(constraints, constraint)
	}

	return constraints
}

// parseConstraintString parses a constraint string into a structured constraint.
func parseConstraintString(constraint string) HoleConstraint {
	constraint = strings.TrimSpace(constraint)

	// Common constraint patterns (ordered: check longer/more-specific patterns first)
	type patternPair struct {
		pattern     string
		description string
	}
	patterns := []patternPair{
		{"implements", "Interface implementation required"},
		{"idempotent", "Idempotent operation required"},
		{"immutable", "Immutability required"},
		{"concurrent", "Concurrent execution support required"},
		{"async", "Asynchronous execution required"}, // Check before "sync"
		{"extends", "Type extension required"},
		{"requires", "Dependency required"},
		{"thread", "Thread-safety requirement"},
		{"sync", "Synchronous execution required"},
		{"mutable", "Mutability allowed"},
		{"pure", "Pure function required (no side effects)"},
		{"atomic", "Atomic operation required"},
	}

	for _, p := range patterns {
		if strings.Contains(strings.ToLower(constraint), p.pattern) {
			return HoleConstraint{
				Type:        p.pattern,
				Value:       constraint,
				Description: p.description,
				Satisfied:   false,
			}
		}
	}

	// Generic constraint
	return HoleConstraint{
		Type:        "generic",
		Value:       constraint,
		Description: constraint,
		Satisfied:   false,
	}
}

// CalculateHolePriority estimates the priority of implementing a typed hole.
func CalculateHolePriority(hole TypedHole, relationships []Relationship) int {
	priority := 5 // Base priority

	// Higher priority for holes mentioned in relationships
	mentionCount := 0
	for _, rel := range relationships {
		if strings.Contains(rel.Subject, hole.Type) ||
			strings.Contains(rel.Predicate, hole.Type) ||
			strings.Contains(rel.Object, hole.Type) {
			mentionCount++
		}
	}

	// Each mention increases priority
	priority += min(mentionCount, 3)

	// Constraint holes have higher priority
	if hole.Constraint != "" {
		priority += 2
	}

	return min(priority, 10)
}

// CalculateHoleComplexity estimates the complexity of implementing a typed hole.
func CalculateHoleComplexity(hole TypedHole, constraints []HoleConstraint) int {
	complexity := 3 // Base complexity

	// More constraints = more complexity
	complexity += len(constraints)

	// Specific constraint types add complexity
	for _, constraint := range constraints {
		switch constraint.Type {
		case "thread", "concurrent", "atomic":
			complexity += 2 // Concurrency is hard
		case "pure", "idempotent":
			complexity += 1 // Functional requirements add complexity
		case "async":
			complexity += 1 // Async adds complexity
		}
	}

	return min(complexity, 10)
}

// SuggestImplementation provides implementation guidance for a typed hole.
func SuggestImplementation(hole TypedHole, constraints []HoleConstraint) string {
	suggestions := []string{}

	// Type-based suggestions
	typeHints := map[string]string{
		"Function":  "Implement as a pure function with clear input/output",
		"Method":    "Add method to appropriate struct/interface",
		"Interface": "Define interface contract with key methods",
		"Struct":    "Define struct with necessary fields",
		"Type":      "Create type alias or new type definition",
		"Handler":   "Implement handler function with error handling",
		"Service":   "Create service with dependency injection",
		"Manager":   "Implement manager pattern with lifecycle methods",
		"Factory":   "Create factory function for object construction",
		"Builder":   "Implement builder pattern for complex objects",
	}

	if hint, ok := typeHints[hole.Type]; ok {
		suggestions = append(suggestions, hint)
	}

	// Constraint-based suggestions
	for _, constraint := range constraints {
		switch constraint.Type {
		case "thread", "concurrent":
			suggestions = append(suggestions, "Use sync.Mutex or sync.RWMutex for thread safety")
		case "async":
			suggestions = append(suggestions, "Use goroutines and channels for async execution")
		case "pure":
			suggestions = append(suggestions, "Avoid side effects, use immutable data")
		case "atomic":
			suggestions = append(suggestions, "Use atomic operations from sync/atomic")
		}
	}

	if len(suggestions) == 0 {
		return fmt.Sprintf("Implement %s with appropriate error handling and tests", hole.Type)
	}

	return strings.Join(suggestions, "; ")
}

// EnhanceTypedHole enriches a typed hole with additional metadata.
func EnhanceTypedHole(hole TypedHole, relationships []Relationship) EnhancedTypedHole {
	constraints := ParseTypedHoleConstraints(hole)
	priority := CalculateHolePriority(hole, relationships)
	complexity := CalculateHoleComplexity(hole, constraints)
	suggestion := SuggestImplementation(hole, constraints)

	return EnhancedTypedHole{
		TypedHole:     hole,
		Constraints:   constraints,
		Priority:      priority,
		Complexity:    complexity,
		SuggestedImpl: suggestion,
	}
}

// FindRelatedHoles identifies holes that are related through dependencies or mentions.
func FindRelatedHoles(holes []TypedHole, relationships []Relationship) map[string][]string {
	related := make(map[string][]string)

	// Build a unique ID for each hole (using type + position)
	holeIDs := make(map[string]string)
	for i, hole := range holes {
		id := fmt.Sprintf("%s_%d", hole.Type, i)
		holeIDs[hole.Type] = id
		related[id] = []string{}
	}

	// Find relationships mentioning holes
	for _, rel := range relationships {
		// Check if subject or object mentions a hole type
		for holeType, id := range holeIDs {
			if strings.Contains(rel.Subject, holeType) {
				// Find other holes in the same relationship
				for otherType, otherID := range holeIDs {
					if otherType != holeType && strings.Contains(rel.Object, otherType) {
						// Add bidirectional relationship
						if !contains(related[id], otherID) {
							related[id] = append(related[id], otherID)
						}
						if !contains(related[otherID], id) {
							related[otherID] = append(related[otherID], id)
						}
					}
				}
			}
		}
	}

	return related
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// contains checks if a string slice contains a value.
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// HolePrioritizer sorts holes by priority and complexity.
type HolePrioritizer struct {
	Holes []EnhancedTypedHole
}

// NewHolePrioritizer creates a new hole prioritizer.
func NewHolePrioritizer(holes []TypedHole, relationships []Relationship) *HolePrioritizer {
	enhanced := make([]EnhancedTypedHole, len(holes))
	for i, hole := range holes {
		enhanced[i] = EnhanceTypedHole(hole, relationships)
	}

	return &HolePrioritizer{
		Holes: enhanced,
	}
}

// GetByPriority returns holes sorted by priority (highest first).
func (p *HolePrioritizer) GetByPriority() []EnhancedTypedHole {
	// Create a copy to avoid modifying original
	sorted := make([]EnhancedTypedHole, len(p.Holes))
	copy(sorted, p.Holes)

	// Simple bubble sort by priority (descending)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Priority > sorted[i].Priority {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// GetByComplexity returns holes sorted by complexity (lowest first).
func (p *HolePrioritizer) GetByComplexity() []EnhancedTypedHole {
	// Create a copy to avoid modifying original
	sorted := make([]EnhancedTypedHole, len(p.Holes))
	copy(sorted, p.Holes)

	// Simple bubble sort by complexity (ascending)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Complexity < sorted[i].Complexity {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// GetRecommendedOrder returns holes in recommended implementation order.
// Priority: Low complexity, high priority first.
func (p *HolePrioritizer) GetRecommendedOrder() []EnhancedTypedHole {
	// Create a copy to avoid modifying original
	sorted := make([]EnhancedTypedHole, len(p.Holes))
	copy(sorted, p.Holes)

	// Sort by: priority/complexity ratio (higher is better)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			// Calculate scores (priority/complexity)
			scoreI := float64(sorted[i].Priority) / float64(max(sorted[i].Complexity, 1))
			scoreJ := float64(sorted[j].Priority) / float64(max(sorted[j].Complexity, 1))

			if scoreJ > scoreI {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
