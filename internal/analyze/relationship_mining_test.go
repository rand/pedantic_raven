package analyze

import (
	"fmt"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"testing"
	"time"
)

// Test basic pattern mining with simple dataset
func TestMinePatterns_Basic(t *testing.T) {
	analysis := createSampleAnalysis()

	patterns := MinePatterns(analysis)

	if len(patterns) == 0 {
		t.Fatal("Expected patterns to be discovered, got none")
	}

	// Verify patterns are sorted by strength
	for i := 1; i < len(patterns); i++ {
		if patterns[i-1].Strength < patterns[i].Strength {
			t.Errorf("Patterns not sorted by strength: pattern[%d]=%.3f < pattern[%d]=%.3f",
				i-1, patterns[i-1].Strength, i, patterns[i].Strength)
		}
	}
}

// Test pattern mining with empty analysis
func TestMinePatterns_EmptyAnalysis(t *testing.T) {
	analysis := &semantic.Analysis{
		Entities:      []semantic.Entity{},
		Relationships: []semantic.Relationship{},
	}

	patterns := MinePatterns(analysis)

	if len(patterns) != 0 {
		t.Errorf("Expected 0 patterns from empty analysis, got %d", len(patterns))
	}
}

// Test pattern mining with nil analysis
func TestMinePatterns_NilAnalysis(t *testing.T) {
	patterns := MinePatterns(nil)

	if len(patterns) != 0 {
		t.Errorf("Expected 0 patterns from nil analysis, got %d", len(patterns))
	}
}

// Test pattern mining with custom options
func TestMinePatternsWithOptions(t *testing.T) {
	analysis := createSampleAnalysis()

	opts := MiningOptions{
		MinOccurrences:   3, // Higher threshold
		MinConfidence:    0.7,
		MaxExamples:      2,
		ClusterThreshold: 0.7,
	}

	patterns := MinePatternsWithOptions(analysis, opts)

	// Should filter out low-occurrence patterns
	for _, pattern := range patterns {
		if pattern.Occurrences < opts.MinOccurrences {
			t.Errorf("Pattern has occurrences %d, below minimum %d",
				pattern.Occurrences, opts.MinOccurrences)
		}

		if pattern.AvgConfidence < opts.MinConfidence {
			t.Errorf("Pattern has confidence %.2f, below minimum %.2f",
				pattern.AvgConfidence, opts.MinConfidence)
		}

		if len(pattern.Examples) > opts.MaxExamples {
			t.Errorf("Pattern has %d examples, exceeding maximum %d",
				len(pattern.Examples), opts.MaxExamples)
		}
	}
}

// Test pattern structure and fields
func TestPatternStructure(t *testing.T) {
	analysis := createSampleAnalysis()
	patterns := MinePatterns(analysis)

	if len(patterns) == 0 {
		t.Skip("No patterns to test")
	}

	pattern := patterns[0]

	// Verify all fields are populated
	if pattern.Predicate == "" {
		t.Error("Pattern predicate is empty")
	}

	if pattern.Occurrences <= 0 {
		t.Error("Pattern occurrences should be > 0")
	}

	if pattern.AvgConfidence < 0 || pattern.AvgConfidence > 1.0 {
		t.Errorf("Pattern confidence %.2f out of range [0.0, 1.0]", pattern.AvgConfidence)
	}

	if pattern.Strength < 0 || pattern.Strength > 1.0 {
		t.Errorf("Pattern strength %.3f out of range [0.0, 1.0]", pattern.Strength)
	}

	if len(pattern.Examples) == 0 {
		t.Error("Pattern should have at least one example")
	}
}

// Test pattern grouping by predicate
func TestPatternGrouping(t *testing.T) {
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "John", Type: semantic.EntityPerson},
			{Text: "Alice", Type: semantic.EntityPerson},
			{Text: "Acme Corp", Type: semantic.EntityOrganization},
			{Text: "Tech Inc", Type: semantic.EntityOrganization},
		},
		Relationships: []semantic.Relationship{
			{Subject: "John", Predicate: "works_at", Object: "Acme Corp"},
			{Subject: "Alice", Predicate: "works_at", Object: "Tech Inc"},
			{Subject: "John", Predicate: "works at", Object: "Acme Corp"}, // Space variation
		},
	}

	patterns := MinePatterns(analysis)

	// Should group "works_at" and "works at" together
	if len(patterns) != 1 {
		t.Errorf("Expected 1 pattern (normalized), got %d", len(patterns))
	}

	if patterns[0].Occurrences != 3 {
		t.Errorf("Expected 3 occurrences (grouped), got %d", patterns[0].Occurrences)
	}
}

// Test strength calculation
func TestStrengthCalculation(t *testing.T) {
	analysis := createSampleAnalysis()
	patterns := MinePatterns(analysis)

	if len(patterns) == 0 {
		t.Skip("No patterns to test")
	}

	// Verify strength is computed correctly
	// Strength = (occurrences × avg_confidence × diversity_factor) / total_relationships
	totalRelationships := float64(len(analysis.Relationships))

	for _, pattern := range patterns {
		// Manually calculate expected strength
		diversityFactor := float64(len(pattern.Examples)) / 3.0 // MaxExamples = 3
		if diversityFactor > 1.0 {
			diversityFactor = 1.0
		}

		expectedStrength := (float64(pattern.Occurrences) * pattern.AvgConfidence * (1.0 + diversityFactor)) / totalRelationships

		// Allow small floating point differences
		diff := pattern.Strength - expectedStrength
		if diff < -0.001 || diff > 0.001 {
			t.Errorf("Strength calculation mismatch: got %.3f, expected %.3f",
				pattern.Strength, expectedStrength)
		}
	}
}

// Test pattern clustering
func TestClusterPatterns_Basic(t *testing.T) {
	patterns := []RelationshipPattern{
		{Predicate: "works_at", Occurrences: 5, Strength: 0.5},
		{Predicate: "employed_at", Occurrences: 3, Strength: 0.3},
		{Predicate: "creates", Occurrences: 4, Strength: 0.4},
		{Predicate: "develops", Occurrences: 2, Strength: 0.2},
	}

	clusters := ClusterPatterns(patterns, 0.7)

	if len(clusters) == 0 {
		t.Fatal("Expected clusters to be created, got none")
	}

	// Verify clusters are sorted by strength
	for i := 1; i < len(clusters); i++ {
		if clusters[i-1].Strength < clusters[i].Strength {
			t.Errorf("Clusters not sorted by strength")
		}
	}

	// Verify each cluster has predicates and patterns
	for i, cluster := range clusters {
		if len(cluster.Predicates) == 0 {
			t.Errorf("Cluster %d has no predicates", i)
		}
		if len(cluster.Patterns) == 0 {
			t.Errorf("Cluster %d has no patterns", i)
		}
		if cluster.ClusterLabel == "" {
			t.Errorf("Cluster %d has no label", i)
		}
	}
}

// Test clustering with high similarity threshold
func TestClusterPatterns_HighThreshold(t *testing.T) {
	patterns := []RelationshipPattern{
		{Predicate: "works_at", Occurrences: 5, Strength: 0.5},
		{Predicate: "employed_at", Occurrences: 3, Strength: 0.3},
	}

	// High threshold means fewer clusters (each predicate separate)
	clusters := ClusterPatterns(patterns, 0.95)

	// Should have at least 1 cluster
	if len(clusters) == 0 {
		t.Fatal("Expected at least 1 cluster")
	}
}

// Test clustering with similar predicates
func TestClusterPatterns_SimilarPredicates(t *testing.T) {
	patterns := []RelationshipPattern{
		{Predicate: "create", Occurrences: 5, Strength: 0.5},
		{Predicate: "creates", Occurrences: 4, Strength: 0.4},
		{Predicate: "created", Occurrences: 3, Strength: 0.3},
	}

	clusters := ClusterPatterns(patterns, 0.6)

	// Similar predicates should cluster together
	if len(clusters) > 2 {
		t.Errorf("Expected similar predicates to cluster, got %d clusters", len(clusters))
	}
}

// Test Levenshtein distance calculation
func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "abd", 1},
		{"abc", "def", 3},
		{"kitten", "sitting", 3},
		{"works_at", "employed_at", 7},
	}

	for _, test := range tests {
		distance := levenshteinDistance(test.s1, test.s2)
		if distance != test.expected {
			t.Errorf("levenshteinDistance(%q, %q) = %d, expected %d",
				test.s1, test.s2, distance, test.expected)
		}
	}
}

// Test predicate similarity calculation
func TestCalculatePredicateSimilarity(t *testing.T) {
	tests := []struct {
		pred1       string
		pred2       string
		minExpected float64
	}{
		{"works_at", "works_at", 1.0},
		{"works_at", "WORKS_AT", 1.0}, // Case insensitive
		{"works", "works_at", 0.8},    // Substring match
		{"create", "creates", 0.7},    // Similar
	}

	for _, test := range tests {
		similarity := calculatePredicateSimilarity(test.pred1, test.pred2)
		if similarity < test.minExpected {
			t.Errorf("calculatePredicateSimilarity(%q, %q) = %.2f, expected >= %.2f",
				test.pred1, test.pred2, similarity, test.minExpected)
		}
	}
}

// Test pattern statistics calculation
func TestCalculatePatternStats(t *testing.T) {
	patterns := []RelationshipPattern{
		{Predicate: "works_at", Occurrences: 5, AvgConfidence: 0.8},
		{Predicate: "creates", Occurrences: 3, AvgConfidence: 0.9},
		{Predicate: "works_at", Occurrences: 2, AvgConfidence: 0.7},
	}

	stats := CalculatePatternStats(patterns)

	if stats.TotalPatterns != 3 {
		t.Errorf("Expected 3 total patterns, got %d", stats.TotalPatterns)
	}

	if stats.UniquePredicates != 2 {
		t.Errorf("Expected 2 unique predicates, got %d", stats.UniquePredicates)
	}

	expectedAvgOccurrences := (5.0 + 3.0 + 2.0) / 3.0
	epsilon := 0.001
	if stats.AvgOccurrences < expectedAvgOccurrences-epsilon || stats.AvgOccurrences > expectedAvgOccurrences+epsilon {
		t.Errorf("Expected avg occurrences %.2f, got %.2f",
			expectedAvgOccurrences, stats.AvgOccurrences)
	}

	expectedAvgConfidence := (0.8 + 0.9 + 0.7) / 3.0
	if stats.AvgConfidence < expectedAvgConfidence-epsilon || stats.AvgConfidence > expectedAvgConfidence+epsilon {
		t.Errorf("Expected avg confidence %.2f, got %.2f",
			expectedAvgConfidence, stats.AvgConfidence)
	}

	if len(stats.TopPredicates) == 0 {
		t.Error("Expected top predicates to be populated")
	}
}

// Test pattern statistics with empty patterns
func TestCalculatePatternStats_Empty(t *testing.T) {
	stats := CalculatePatternStats([]RelationshipPattern{})

	if stats.TotalPatterns != 0 {
		t.Errorf("Expected 0 total patterns, got %d", stats.TotalPatterns)
	}

	if stats.UniquePredicates != 0 {
		t.Errorf("Expected 0 unique predicates, got %d", stats.UniquePredicates)
	}
}

// Test predicate normalization
func TestNormalizePredicateForKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"works_at", "works_at"},
		{"works at", "works_at"},
		{"  WORKS AT  ", "works_at"},
		{"Creates", "creates"},
	}

	for _, test := range tests {
		result := normalizePredicateForKey(test.input)
		if result != test.expected {
			t.Errorf("normalizePredicateForKey(%q) = %q, expected %q",
				test.input, result, test.expected)
		}
	}
}

// Test building entity type lookup
func TestBuildEntityTypeLookup(t *testing.T) {
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "John", Type: semantic.EntityPerson},
			{Text: "Acme Corp", Type: semantic.EntityOrganization},
			{Text: "Python", Type: semantic.EntityTechnology},
		},
	}

	lookup := buildEntityTypeLookup(analysis)

	if len(lookup) != 3 {
		t.Errorf("Expected 3 entries in lookup, got %d", len(lookup))
	}

	if lookup["John"] != semantic.EntityPerson {
		t.Errorf("Expected John to be Person type")
	}

	if lookup["Acme Corp"] != semantic.EntityOrganization {
		t.Errorf("Expected Acme Corp to be Organization type")
	}

	if lookup["Python"] != semantic.EntityTechnology {
		t.Errorf("Expected Python to be Technology type")
	}
}

// Test with large dataset (performance)
func TestMinePatterns_LargeDataset(t *testing.T) {
	// Create a large analysis with 1000+ relationships
	analysis := createLargeAnalysis(1000)

	start := time.Now()
	patterns := MinePatterns(analysis)
	duration := time.Since(start)

	if len(patterns) == 0 {
		t.Fatal("Expected patterns from large dataset")
	}

	// Should complete in reasonable time (< 1 second for 1000 relationships)
	if duration > time.Second {
		t.Logf("Warning: Mining took %v for 1000 relationships", duration)
	}

	t.Logf("Mined %d patterns from %d relationships in %v",
		len(patterns), len(analysis.Relationships), duration)
}

// Helper function to create sample analysis
func createSampleAnalysis() *semantic.Analysis {
	return &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "John", Type: semantic.EntityPerson},
			{Text: "Alice", Type: semantic.EntityPerson},
			{Text: "Bob", Type: semantic.EntityPerson},
			{Text: "Acme Corp", Type: semantic.EntityOrganization},
			{Text: "Tech Inc", Type: semantic.EntityOrganization},
			{Text: "API Gateway", Type: semantic.EntityTechnology},
			{Text: "Mobile App", Type: semantic.EntityTechnology},
		},
		Relationships: []semantic.Relationship{
			{Subject: "John", Predicate: "works_at", Object: "Acme Corp"},
			{Subject: "Alice", Predicate: "works_at", Object: "Tech Inc"},
			{Subject: "Bob", Predicate: "works_at", Object: "Acme Corp"},
			{Subject: "Acme Corp", Predicate: "creates", Object: "API Gateway"},
			{Subject: "Tech Inc", Predicate: "creates", Object: "Mobile App"},
			{Subject: "John", Predicate: "develops", Object: "API Gateway"},
		},
	}
}

// Helper function to create large analysis
func createLargeAnalysis(numRelationships int) *semantic.Analysis {
	entities := make([]semantic.Entity, 0)
	relationships := make([]semantic.Relationship, 0, numRelationships)

	// Create diverse entities
	for i := 0; i < numRelationships/2; i++ {
		entities = append(entities, semantic.Entity{
			Text: fmt.Sprintf("Person%d", i),
			Type: semantic.EntityPerson,
		})
		entities = append(entities, semantic.Entity{
			Text: fmt.Sprintf("Org%d", i),
			Type: semantic.EntityOrganization,
		})
	}

	// Create relationships with various predicates
	predicates := []string{"works_at", "creates", "develops", "manages", "uses"}
	for i := 0; i < numRelationships; i++ {
		relationships = append(relationships, semantic.Relationship{
			Subject:   fmt.Sprintf("Person%d", i%100),
			Predicate: predicates[i%len(predicates)],
			Object:    fmt.Sprintf("Org%d", i%50),
		})
	}

	return &semantic.Analysis{
		Entities:      entities,
		Relationships: relationships,
	}
}
