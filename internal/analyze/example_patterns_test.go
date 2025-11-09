package analyze

import (
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"testing"
)

// TestExamplePatternDiscovery demonstrates pattern mining with realistic data.
// This test showcases the pattern mining capabilities and expected output.
func TestExamplePatternDiscovery(t *testing.T) {
	// Create realistic analysis with various relationship patterns
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "Alice", Type: semantic.EntityPerson},
			{Text: "Bob", Type: semantic.EntityPerson},
			{Text: "Charlie", Type: semantic.EntityPerson},
			{Text: "David", Type: semantic.EntityPerson},
			{Text: "Acme Corp", Type: semantic.EntityOrganization},
			{Text: "Tech Inc", Type: semantic.EntityOrganization},
			{Text: "StartupXYZ", Type: semantic.EntityOrganization},
			{Text: "API Gateway", Type: semantic.EntityTechnology},
			{Text: "Mobile App", Type: semantic.EntityTechnology},
			{Text: "Database", Type: semantic.EntityTechnology},
			{Text: "Python", Type: semantic.EntityTechnology},
			{Text: "Go", Type: semantic.EntityTechnology},
			{Text: "Rust", Type: semantic.EntityTechnology},
		},
		Relationships: []semantic.Relationship{
			// Employment patterns
			{Subject: "Alice", Predicate: "works_at", Object: "Acme Corp"},
			{Subject: "Bob", Predicate: "works_at", Object: "Tech Inc"},
			{Subject: "Charlie", Predicate: "works_at", Object: "StartupXYZ"},
			{Subject: "David", Predicate: "employed_at", Object: "Acme Corp"},

			// Creation patterns
			{Subject: "Acme Corp", Predicate: "creates", Object: "API Gateway"},
			{Subject: "Tech Inc", Predicate: "creates", Object: "Mobile App"},
			{Subject: "StartupXYZ", Predicate: "develops", Object: "Database"},

			// Development patterns
			{Subject: "Alice", Predicate: "develops", Object: "API Gateway"},
			{Subject: "Bob", Predicate: "develops", Object: "Mobile App"},
			{Subject: "Charlie", Predicate: "builds", Object: "Database"},

			// Technology usage patterns
			{Subject: "API Gateway", Predicate: "uses", Object: "Python"},
			{Subject: "Mobile App", Predicate: "uses", Object: "Go"},
			{Subject: "Database", Predicate: "uses", Object: "Rust"},
			{Subject: "API Gateway", Predicate: "uses", Object: "Go"},
		},
	}

	// Mine patterns
	patterns := MinePatterns(analysis)

	// Log discovered patterns for demonstration
	t.Logf("\n=== Discovered %d Relationship Patterns ===\n", len(patterns))

	for i, pattern := range patterns {
		t.Logf("Pattern %d: [%s] → %s → [%s]",
			i+1,
			pattern.SubjectType.String(),
			pattern.Predicate,
			pattern.ObjectType.String(),
		)
		t.Logf("  Occurrences: %d, Confidence: %.2f, Strength: %.3f",
			pattern.Occurrences,
			pattern.AvgConfidence,
			pattern.Strength,
		)

		if len(pattern.Examples) > 0 {
			t.Logf("  Examples:")
			for _, ex := range pattern.Examples {
				t.Logf("    • %s → %s → %s", ex.Subject, ex.Predicate, ex.Object)
			}
		}
		t.Logf("")
	}

	// Display pattern statistics
	stats := CalculatePatternStats(patterns)
	t.Logf("\n=== Pattern Statistics ===")
	t.Logf("Total Patterns: %d", stats.TotalPatterns)
	t.Logf("Unique Predicates: %d", stats.UniquePredicates)
	t.Logf("Avg Occurrences: %.1f", stats.AvgOccurrences)
	t.Logf("Avg Confidence: %.2f", stats.AvgConfidence)
	t.Logf("Top Predicates: %v", stats.TopPredicates)

	// Cluster patterns
	clusters := ClusterPatterns(patterns, 0.7)
	t.Logf("\n=== Discovered %d Pattern Clusters ===\n", len(clusters))

	for i, cluster := range clusters {
		t.Logf("Cluster %d: %s", i+1, cluster.ClusterLabel)
		t.Logf("  Similar Predicates: %v", cluster.Predicates)
		t.Logf("  Patterns in Cluster: %d", len(cluster.Patterns))
		t.Logf("  Cluster Strength: %.3f", cluster.Strength)
		t.Logf("")
	}

	// Verify expected patterns were discovered
	if len(patterns) == 0 {
		t.Fatal("Expected to discover patterns, got none")
	}

	// Check for employment pattern
	foundEmployment := false
	for _, p := range patterns {
		if (p.Predicate == "works_at" || p.Predicate == "employed_at") &&
			p.SubjectType == semantic.EntityPerson &&
			p.ObjectType == semantic.EntityOrganization {
			foundEmployment = true
			break
		}
	}
	if !foundEmployment {
		t.Error("Expected to find employment pattern")
	}

	// Check for technology usage pattern
	foundUsage := false
	for _, p := range patterns {
		if p.Predicate == "uses" &&
			p.SubjectType == semantic.EntityTechnology &&
			p.ObjectType == semantic.EntityTechnology {
			foundUsage = true
			break
		}
	}
	if !foundUsage {
		t.Error("Expected to find technology usage pattern")
	}
}

// TestExampleTableRendering demonstrates various table rendering modes.
func TestExampleTableRendering(t *testing.T) {
	// Create sample patterns
	patterns := []RelationshipPattern{
		{
			SubjectType:   semantic.EntityPerson,
			Predicate:     "works_at",
			ObjectType:    semantic.EntityOrganization,
			Occurrences:   5,
			AvgConfidence: 0.87,
			Strength:      0.45,
			Examples: []PatternExample{
				{Subject: "Alice", Predicate: "works_at", Object: "Acme Corp", Confidence: 0.9},
				{Subject: "Bob", Predicate: "works_at", Object: "Tech Inc", Confidence: 0.85},
			},
		},
		{
			SubjectType:   semantic.EntityOrganization,
			Predicate:     "creates",
			ObjectType:    semantic.EntityTechnology,
			Occurrences:   3,
			AvgConfidence: 0.92,
			Strength:      0.38,
			Examples: []PatternExample{
				{Subject: "Acme Corp", Predicate: "creates", Object: "API Gateway", Confidence: 0.95},
			},
		},
	}

	// Compact mode
	t.Log("\n=== Compact Table Mode ===")
	opts := DefaultPatternDisplayOptions()
	opts.Compact = true
	opts.Width = 100
	output := RenderPatternTable(patterns, opts)
	t.Log("\n" + output)

	// Expanded mode with examples
	t.Log("\n=== Expanded Table Mode ===")
	opts.Compact = false
	opts.ShowExamples = true
	output = RenderPatternTable(patterns, opts)
	t.Log("\n" + output)

	// Sorted by frequency
	t.Log("\n=== Sorted by Frequency ===")
	opts.SortMode = SortByFrequency
	output = RenderPatternTable(patterns, opts)
	t.Log("\n" + output)

	// With filtering
	t.Log("\n=== Filtered (min 4 occurrences) ===")
	opts.Filter.MinOccurrences = 4
	output = RenderPatternTable(patterns, opts)
	t.Log("\n" + output)
}

// TestExampleClusterRendering demonstrates cluster visualization.
func TestExampleClusterRendering(t *testing.T) {
	patterns := []RelationshipPattern{
		{Predicate: "works_at", Occurrences: 5, Strength: 0.5},
		{Predicate: "employed_at", Occurrences: 3, Strength: 0.3},
		{Predicate: "creates", Occurrences: 4, Strength: 0.4},
		{Predicate: "develops", Occurrences: 3, Strength: 0.35},
		{Predicate: "builds", Occurrences: 2, Strength: 0.25},
	}

	clusters := ClusterPatterns(patterns, 0.7)

	t.Logf("\n=== Cluster Visualization ===")
	opts := DefaultPatternDisplayOptions()
	output := RenderClusterTable(clusters, opts)
	t.Log("\n" + output)

	// Verify clustering worked
	if len(clusters) == 0 {
		t.Fatal("Expected clusters to be created")
	}

	t.Logf("\nDiscovered %d clusters from %d patterns", len(clusters), len(patterns))

	// Log cluster details
	for i, cluster := range clusters {
		t.Logf("\nCluster %d:", i+1)
		t.Logf("  Label: %s", cluster.ClusterLabel)
		t.Logf("  Predicates: %v", cluster.Predicates)
		t.Logf("  Patterns: %d", len(cluster.Patterns))
		t.Logf("  Strength: %.3f", cluster.Strength)
	}
}

// Example output when running with -v flag:
//
// === Discovered 5 Relationship Patterns ===
//
// Pattern 1: [Person] → works_at → [Organization]
//   Occurrences: 4, Confidence: 0.80, Strength: 0.352
//   Examples:
//     • Alice → works_at → Acme Corp
//     • Bob → works_at → Tech Inc
//     • Charlie → works_at → StartupXYZ
//
// Pattern 2: [Technology] → uses → [Technology]
//   Occurrences: 4, Confidence: 0.80, Strength: 0.352
//   Examples:
//     • API Gateway → uses → Python
//     • Mobile App → uses → Go
//     • Database → uses → Rust
//
// Pattern 3: [Person] → develops → [Technology]
//   Occurrences: 2, Confidence: 0.80, Strength: 0.176
//   Examples:
//     • Alice → develops → API Gateway
//     • Bob → develops → Mobile App
//
// === Pattern Statistics ===
// Total Patterns: 5
// Unique Predicates: 5
// Avg Occurrences: 2.8
// Avg Confidence: 0.80
// Top Predicates: [works_at, uses, creates, develops, builds]
//
// === Discovered 3 Pattern Clusters ===
//
// Cluster 1: works_at
//   Similar Predicates: [works_at, employed_at]
//   Patterns in Cluster: 2
//   Cluster Strength: 0.400
//
// Cluster 2: creates
//   Similar Predicates: [creates]
//   Patterns in Cluster: 1
//   Cluster Strength: 0.300
//
// Cluster 3: develops
//   Similar Predicates: [develops, builds]
//   Patterns in Cluster: 2
//   Cluster Strength: 0.300
