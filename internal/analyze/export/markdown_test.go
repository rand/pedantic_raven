package export

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

func TestExportMarkdown_Nil(t *testing.T) {
	_, err := ExportMarkdown(nil)
	if err == nil {
		t.Error("Expected error for nil report")
	}
}

func TestExportMarkdown_Empty(t *testing.T) {
	report := NewAnalysisReport("Empty Report", "empty.go")

	md, err := ExportMarkdown(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should contain title
	if !strings.Contains(md, "# Empty Report") {
		t.Error("Missing title in markdown")
	}

	// Should contain table of contents
	if !strings.Contains(md, "## Table of Contents") {
		t.Error("Missing table of contents")
	}
}

func TestExportMarkdown_Statistics(t *testing.T) {
	report := createTestReport()

	md, err := ExportMarkdown(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check statistics section
	if !strings.Contains(md, "## Overview Statistics") {
		t.Error("Missing statistics section")
	}

	if !strings.Contains(md, "Total Entities") {
		t.Error("Missing total entities")
	}

	if !strings.Contains(md, "Total Relationships") {
		t.Error("Missing total relationships")
	}

	// Check for table format
	if !strings.Contains(md, "| Metric | Value |") {
		t.Error("Missing statistics table header")
	}
}

func TestExportMarkdown_EntityFrequencies(t *testing.T) {
	report := createTestReport()

	md, err := ExportMarkdown(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check entity frequency section
	if !strings.Contains(md, "## Entity Frequency Analysis") {
		t.Error("Missing entity frequency section")
	}

	if !strings.Contains(md, "### Top Entities by Frequency") {
		t.Error("Missing top entities subsection")
	}

	// Check for table columns
	if !strings.Contains(md, "| Rank | Entity | Type | Count | Importance |") {
		t.Error("Missing entity frequency table header")
	}

	// Check for entity data
	if !strings.Contains(md, "Alice") {
		t.Error("Missing entity 'Alice'")
	}

	// Check type distribution
	if !strings.Contains(md, "### Entity Type Distribution") {
		t.Error("Missing type distribution section")
	}
}

func TestExportMarkdown_RelationshipPatterns(t *testing.T) {
	report := createTestReport()

	md, err := ExportMarkdown(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check patterns section
	if !strings.Contains(md, "## Relationship Patterns") {
		t.Error("Missing relationship patterns section")
	}

	if !strings.Contains(md, "### Discovered Patterns") {
		t.Error("Missing discovered patterns subsection")
	}

	// Check for table
	if !strings.Contains(md, "| Pattern | Occurrences | Strength | Avg Confidence |") {
		t.Error("Missing patterns table header")
	}

	// Check for pattern data
	if !strings.Contains(md, "uses") {
		t.Error("Missing 'uses' predicate")
	}
}

func TestExportMarkdown_TypedHoles(t *testing.T) {
	report := createTestReport()

	md, err := ExportMarkdown(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check typed holes section
	if !strings.Contains(md, "## Typed Holes") {
		t.Error("Missing typed holes section")
	}

	if !strings.Contains(md, "### Implementation Priority Queue") {
		t.Error("Missing priority queue subsection")
	}

	// Check for table
	if !strings.Contains(md, "| Priority | Type | Complexity | Constraint | Dependencies |") {
		t.Error("Missing typed holes table header")
	}

	// Check for hole data
	if !strings.Contains(md, "UserService") {
		t.Error("Missing 'UserService' hole")
	}

	// Check for detailed section
	if !strings.Contains(md, "### Hole Details") {
		t.Error("Missing hole details section")
	}
}

func TestExportMarkdown_TripleGraph(t *testing.T) {
	report := createTestReport()

	md, err := ExportMarkdown(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check graph section
	if !strings.Contains(md, "## Triple Graph Visualization") {
		t.Error("Missing triple graph section")
	}

	// Check for Mermaid diagram
	if !strings.Contains(md, "```mermaid") {
		t.Error("Missing Mermaid diagram")
	}

	if !strings.Contains(md, "graph TD") {
		t.Error("Missing graph TD directive")
	}

	// Check for graph legend
	if !strings.Contains(md, "### Graph Legend") {
		t.Error("Missing graph legend")
	}
}

func TestExportMarkdown_CustomOptions(t *testing.T) {
	report := createTestReport()

	opts := ExportOptions{
		Format:              FormatMarkdown,
		IncludeMetadata:     false,
		IncludeStatistics:   true,
		IncludeFrequencies:  false,
		IncludePatterns:     false,
		IncludeTypedHoles:   false,
		IncludeTripleGraph:  false,
		MaxExamplesPerPattern: 1,
		MaxFrequenciesToShow:  5,
	}

	md, err := ExportMarkdownWithOptions(report, opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should NOT contain metadata
	if strings.Contains(md, "**Generated:**") {
		t.Error("Should not include metadata")
	}

	// Should contain statistics (enabled)
	if !strings.Contains(md, "## Overview Statistics") {
		t.Error("Should include statistics")
	}

	// Should NOT contain frequencies (disabled)
	if strings.Contains(md, "## Entity Frequency Analysis") {
		t.Error("Should not include frequencies")
	}

	// Should NOT contain patterns (disabled)
	if strings.Contains(md, "## Relationship Patterns") {
		t.Error("Should not include patterns")
	}
}

func TestExportMarkdown_MaxLimits(t *testing.T) {
	report := createLargeTestReport()

	opts := DefaultExportOptions(FormatMarkdown)
	opts.MaxFrequenciesToShow = 3

	md, err := ExportMarkdownWithOptions(report, opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Count entity rows (should be limited to 3)
	lines := strings.Split(md, "\n")
	entityRowCount := 0
	inEntityTable := false

	for _, line := range lines {
		if strings.Contains(line, "### Top Entities by Frequency") {
			inEntityTable = true
			continue
		}
		if inEntityTable && strings.HasPrefix(line, "| ") && !strings.Contains(line, "Rank") {
			entityRowCount++
		}
		if inEntityTable && strings.HasPrefix(line, "###") {
			break
		}
	}

	if entityRowCount > opts.MaxFrequenciesToShow {
		t.Errorf("Expected at most %d entity rows, got %d", opts.MaxFrequenciesToShow, entityRowCount)
	}
}

func TestMermaidEscaping(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`Test "quotes"`, `Test &quot;quotes&quot;`},
		{`Test <angle>`, `Test &lt;angle&gt;`},
		{`Test & ampersand`, `Test & ampersand`}, // & is okay in Mermaid
	}

	for _, tt := range tests {
		result := escapeForMermaid(tt.input)
		if result != tt.expected {
			t.Errorf("escapeForMermaid(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestMermaidNodeID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with spaces", "with_spaces"},
		{"with-dashes", "with_dashes"},
		{"123number", "N123number"}, // Should add N prefix
		{"special!@#", "special___"},
	}

	for _, tt := range tests {
		result := nodeIDForMermaid(tt.input)
		if result != tt.expected {
			t.Errorf("nodeIDForMermaid(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestExportMarkdown_LargeDataset(t *testing.T) {
	report := createLargeTestReport()

	md, err := ExportMarkdown(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should handle large dataset without errors
	if len(md) == 0 {
		t.Error("Generated markdown is empty")
	}

	// Should contain all major sections
	requiredSections := []string{
		"## Overview Statistics",
		"## Entity Frequency Analysis",
		"## Relationship Patterns",
		"## Typed Holes",
		"## Triple Graph Visualization",
	}

	for _, section := range requiredSections {
		if !strings.Contains(md, section) {
			t.Errorf("Missing required section: %s", section)
		}
	}
}

// Helper functions

func createTestReport() *AnalysisReport {
	report := NewAnalysisReport("Test Report", "test.go")

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
			Examples: []analyze.PatternExample{
				{Subject: "Alice", Predicate: "uses", Object: "Go", Confidence: 0.85},
			},
		},
	})

	// Add typed holes
	report.SetTypedHoles([]semantic.EnhancedTypedHole{
		{
			TypedHole:    semantic.TypedHole{Type: "UserService", Constraint: "thread-safe"},
			Priority:     8,
			Complexity:   6,
			Dependencies: []string{"Database", "Cache"},
		},
	})

	return report
}

func createLargeTestReport() *AnalysisReport {
	report := NewAnalysisReport("Large Test Report", "large.go")

	// Add many entities to graph
	graph := analyze.NewTripleGraph()
	for i := 0; i < 30; i++ {
		graph.AddNode(semantic.Entity{
			Text: fmt.Sprintf("Entity%d", i),
			Type: semantic.EntityType(i % 6),
		})
	}
	report.SetTripleGraph(graph)

	// Add many frequencies
	var frequencies []analyze.EntityFrequency
	for i := 0; i < 50; i++ {
		frequencies = append(frequencies, analyze.EntityFrequency{
			Text:       fmt.Sprintf("Entity%d", i),
			Type:       semantic.EntityType(i % 6),
			Count:      50 - i,
			Importance: (50 - i) / 5,
		})
	}
	report.SetEntityFrequencies(frequencies)

	// Add many patterns
	var patterns []analyze.RelationshipPattern
	for i := 0; i < 20; i++ {
		patterns = append(patterns, analyze.RelationshipPattern{
			SubjectType:   semantic.EntityType(i % 6),
			Predicate:     fmt.Sprintf("predicate%d", i),
			ObjectType:    semantic.EntityType((i + 1) % 6),
			Occurrences:   20 - i,
			Strength:      0.5 + float64(i)/40,
			AvgConfidence: 0.7,
		})
	}
	report.SetRelationshipPatterns(patterns)

	// Add many holes
	var holes []semantic.EnhancedTypedHole
	for i := 0; i < 25; i++ {
		holes = append(holes, semantic.EnhancedTypedHole{
			TypedHole:  semantic.TypedHole{Type: fmt.Sprintf("Hole%d", i), Constraint: "generic"},
			Priority:   10 - (i / 3),
			Complexity: 5 + (i % 5),
		})
	}
	report.SetTypedHoles(holes)

	return report
}
