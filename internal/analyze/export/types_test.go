package export

import (
	"strings"
	"testing"
	"time"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

func TestNewAnalysisReport(t *testing.T) {
	title := "Test Report"
	source := "test.go"

	report := NewAnalysisReport(title, source)

	if report.Title != title {
		t.Errorf("Expected title %q, got %q", title, report.Title)
	}

	if report.Source != source {
		t.Errorf("Expected source %q, got %q", source, report.Source)
	}

	if report.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}

	if time.Since(report.Timestamp) > time.Second {
		t.Error("Timestamp should be recent")
	}
}

func TestSetTripleGraph(t *testing.T) {
	report := NewAnalysisReport("Test", "test.go")

	// Create a simple graph
	graph := analyze.NewTripleGraph()
	graph.AddNode(semantic.Entity{
		Text: "TestEntity",
		Type: semantic.EntityPerson,
	})
	graph.AddNode(semantic.Entity{
		Text: "AnotherEntity",
		Type: semantic.EntityTechnology,
	})

	report.SetTripleGraph(graph)

	if report.TripleGraph != graph {
		t.Error("Graph not set correctly")
	}

	if report.Stats.TotalEntities != 2 {
		t.Errorf("Expected 2 entities, got %d", report.Stats.TotalEntities)
	}

	if report.Stats.UniqueEntityTypes != 2 {
		t.Errorf("Expected 2 unique types, got %d", report.Stats.UniqueEntityTypes)
	}
}

func TestSetEntityFrequencies(t *testing.T) {
	report := NewAnalysisReport("Test", "test.go")

	frequencies := []analyze.EntityFrequency{
		{Text: "First", Type: semantic.EntityPerson, Count: 10, Importance: 5},
		{Text: "Second", Type: semantic.EntityTechnology, Count: 5, Importance: 3},
	}

	report.SetEntityFrequencies(frequencies)

	if len(report.EntityFrequencies) != 2 {
		t.Errorf("Expected 2 frequencies, got %d", len(report.EntityFrequencies))
	}

	if report.Stats.MostCommonEntity != "First" {
		t.Errorf("Expected most common entity 'First', got %q", report.Stats.MostCommonEntity)
	}

	if report.Stats.MostCommonEntityCount != 10 {
		t.Errorf("Expected count 10, got %d", report.Stats.MostCommonEntityCount)
	}
}

func TestSetRelationshipPatterns(t *testing.T) {
	report := NewAnalysisReport("Test", "test.go")

	patterns := []analyze.RelationshipPattern{
		{
			SubjectType:   semantic.EntityPerson,
			Predicate:     "uses",
			ObjectType:    semantic.EntityTechnology,
			Occurrences:   5,
			Strength:      0.8,
			AvgConfidence: 0.75,
		},
		{
			SubjectType:   semantic.EntityOrganization,
			Predicate:     "develops",
			ObjectType:    semantic.EntityTechnology,
			Occurrences:   3,
			Strength:      0.6,
			AvgConfidence: 0.65,
		},
	}

	report.SetRelationshipPatterns(patterns)

	if len(report.RelationshipPatterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(report.RelationshipPatterns))
	}

	if report.Stats.UniquePatterns != 2 {
		t.Errorf("Expected 2 unique patterns, got %d", report.Stats.UniquePatterns)
	}

	if report.Stats.StrongestPatternScore != 0.8 {
		t.Errorf("Expected strongest pattern score 0.8, got %.2f", report.Stats.StrongestPatternScore)
	}

	// Check that the strongest pattern is formatted correctly
	// EntityType.String() returns capitalized names like "Person"
	if report.Stats.StrongestPattern == "" {
		t.Error("Expected non-empty strongest pattern")
	}

	// Should contain the predicate
	if !strings.Contains(report.Stats.StrongestPattern, "uses") {
		t.Errorf("Expected pattern to contain 'uses', got %q", report.Stats.StrongestPattern)
	}
}

func TestSetTypedHoles(t *testing.T) {
	report := NewAnalysisReport("Test", "test.go")

	holes := []semantic.EnhancedTypedHole{
		{
			TypedHole: semantic.TypedHole{Type: "Cache", Constraint: "concurrent"},
			Priority:  9,
			Complexity: 7,
		},
		{
			TypedHole: semantic.TypedHole{Type: "Logger", Constraint: "async"},
			Priority:  5,
			Complexity: 3,
		},
		{
			TypedHole: semantic.TypedHole{Type: "Database", Constraint: "thread-safe"},
			Priority:  8,
			Complexity: 8,
		},
	}

	report.SetTypedHoles(holes)

	if len(report.TypedHoles) != 3 {
		t.Errorf("Expected 3 typed holes, got %d", len(report.TypedHoles))
	}

	if report.Stats.TotalTypedHoles != 3 {
		t.Errorf("Expected total 3, got %d", report.Stats.TotalTypedHoles)
	}

	if report.Stats.HighestPriority != 9 {
		t.Errorf("Expected highest priority 9, got %d", report.Stats.HighestPriority)
	}

	if report.Stats.HighestPriorityHole != "Cache" {
		t.Errorf("Expected highest priority hole 'Cache', got %q", report.Stats.HighestPriorityHole)
	}

	expectedAvgComplexity := (7.0 + 3.0 + 8.0) / 3.0
	if report.Stats.AvgComplexity != expectedAvgComplexity {
		t.Errorf("Expected avg complexity %.2f, got %.2f", expectedAvgComplexity, report.Stats.AvgComplexity)
	}
}

func TestDefaultExportOptions(t *testing.T) {
	tests := []struct {
		format ExportFormat
	}{
		{FormatMarkdown},
		{FormatHTML},
		{FormatPDF},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			opts := DefaultExportOptions(tt.format)

			if opts.Format != tt.format {
				t.Errorf("Expected format %s, got %s", tt.format, opts.Format)
			}

			// Check all defaults are enabled
			if !opts.IncludeMetadata {
				t.Error("Expected metadata to be included")
			}
			if !opts.IncludeStatistics {
				t.Error("Expected statistics to be included")
			}
			if !opts.IncludeTripleGraph {
				t.Error("Expected triple graph to be included")
			}
			if !opts.IncludeFrequencies {
				t.Error("Expected frequencies to be included")
			}
			if !opts.IncludePatterns {
				t.Error("Expected patterns to be included")
			}
			if !opts.IncludeTypedHoles {
				t.Error("Expected typed holes to be included")
			}

			if opts.MaxExamplesPerPattern != 3 {
				t.Errorf("Expected max examples 3, got %d", opts.MaxExamplesPerPattern)
			}

			if opts.MaxFrequenciesToShow != 20 {
				t.Errorf("Expected max frequencies 20, got %d", opts.MaxFrequenciesToShow)
			}
		})
	}
}

func TestSetEmptyData(t *testing.T) {
	report := NewAnalysisReport("Test", "test.go")

	// Set empty data
	report.SetEntityFrequencies([]analyze.EntityFrequency{})
	report.SetRelationshipPatterns([]analyze.RelationshipPattern{})
	report.SetTypedHoles([]semantic.EnhancedTypedHole{})

	// Should not panic and should have zero stats
	if report.Stats.MostCommonEntity != "" {
		t.Error("Expected empty most common entity")
	}

	if report.Stats.StrongestPattern != "" {
		t.Error("Expected empty strongest pattern")
	}

	if report.Stats.HighestPriorityHole != "" {
		t.Error("Expected empty highest priority hole")
	}
}

func TestCompleteReport(t *testing.T) {
	report := NewAnalysisReport("Complete Test Report", "complete.go")
	report.Description = "A complete test report with all sections"

	// Add graph
	graph := analyze.NewTripleGraph()
	graph.AddNode(semantic.Entity{Text: "Alice", Type: semantic.EntityPerson})
	graph.AddNode(semantic.Entity{Text: "Go", Type: semantic.EntityTechnology})
	graph.AddEdge(semantic.Relationship{
		Subject:   "Alice",
		Predicate: "uses",
		Object:    "Go",
	})
	report.SetTripleGraph(graph)

	// Add frequencies
	report.SetEntityFrequencies([]analyze.EntityFrequency{
		{Text: "Alice", Type: semantic.EntityPerson, Count: 5, Importance: 7},
		{Text: "Go", Type: semantic.EntityTechnology, Count: 3, Importance: 6},
	})

	// Add patterns
	report.SetRelationshipPatterns([]analyze.RelationshipPattern{
		{
			SubjectType:   semantic.EntityPerson,
			Predicate:     "uses",
			ObjectType:    semantic.EntityTechnology,
			Occurrences:   1,
			Strength:      0.9,
			AvgConfidence: 0.85,
		},
	})

	// Add typed holes
	report.SetTypedHoles([]semantic.EnhancedTypedHole{
		{
			TypedHole:  semantic.TypedHole{Type: "UserService", Constraint: "thread-safe"},
			Priority:   8,
			Complexity: 6,
			Dependencies: []string{"Database", "Cache"},
		},
	})

	// Verify complete report structure
	if report.Title != "Complete Test Report" {
		t.Error("Title mismatch")
	}

	if report.Stats.TotalEntities != 2 {
		t.Errorf("Expected 2 entities, got %d", report.Stats.TotalEntities)
	}

	if report.Stats.TotalRelationships != 1 {
		t.Errorf("Expected 1 relationship, got %d", report.Stats.TotalRelationships)
	}

	if report.Stats.UniquePatterns != 1 {
		t.Errorf("Expected 1 pattern, got %d", report.Stats.UniquePatterns)
	}

	if report.Stats.TotalTypedHoles != 1 {
		t.Errorf("Expected 1 typed hole, got %d", report.Stats.TotalTypedHoles)
	}
}
