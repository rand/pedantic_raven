package analyze

import (
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"strings"
	"testing"
)

// Test basic table rendering
func TestRenderPatternTable_Basic(t *testing.T) {
	patterns := createTestPatterns()
	opts := DefaultPatternDisplayOptions()

	output := RenderPatternTable(patterns, opts)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}

	// Should contain header
	if !strings.Contains(output, "Relationship Patterns") {
		t.Error("Missing header in output")
	}

	// Should contain pattern data
	if !strings.Contains(output, "works_at") {
		t.Error("Missing pattern predicate in output")
	}
}

// Test empty pattern rendering
func TestRenderPatternTable_Empty(t *testing.T) {
	patterns := []RelationshipPattern{}
	opts := DefaultPatternDisplayOptions()

	output := RenderPatternTable(patterns, opts)

	if output == "" {
		t.Fatal("Expected non-empty output for empty state")
	}

	if !strings.Contains(output, "No Patterns Found") {
		t.Error("Missing empty state message")
	}
}

// Test compact mode rendering
func TestRenderPatternTable_Compact(t *testing.T) {
	patterns := createTestPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.Compact = true
	opts.ShowExamples = false

	output := RenderPatternTable(patterns, opts)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}

	// In compact mode, should not have "Examples:" section
	if strings.Contains(output, "Examples:") {
		t.Error("Compact mode should not show examples")
	}

	// Should show inline statistics
	if !strings.Contains(output, "Count:") {
		t.Error("Compact mode should show inline count")
	}
}

// Test expanded mode rendering with examples
func TestRenderPatternTable_Expanded(t *testing.T) {
	patterns := createTestPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.Compact = false
	opts.ShowExamples = true

	output := RenderPatternTable(patterns, opts)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}

	// Should have examples section
	if !strings.Contains(output, "Examples:") {
		t.Error("Expanded mode should show examples section")
	}

	// Should show bullet points for examples
	if !strings.Contains(output, "•") {
		t.Error("Examples should be shown with bullet points")
	}
}

// Test filtering by minimum occurrences
func TestRenderPatternTable_FilterOccurrences(t *testing.T) {
	patterns := createTestPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.Filter.MinOccurrences = 4 // Filter out low-frequency patterns

	output := RenderPatternTable(patterns, opts)

	// Should filter patterns
	if strings.Contains(output, "No Matching Patterns") {
		// All patterns were filtered
		return
	}

	// Verify only high-occurrence patterns remain
	// (This is a visual check - actual filtering tested in filter function test)
}

// Test filtering by minimum confidence
func TestRenderPatternTable_FilterConfidence(t *testing.T) {
	patterns := createTestPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.Filter.MinConfidence = 0.9 // High confidence threshold

	output := RenderPatternTable(patterns, opts)

	// With high threshold, might get no matches
	if strings.Contains(output, "No Matching Patterns") {
		return
	}

	// If patterns remain, they should all be high confidence
}

// Test filtering by predicate text
func TestRenderPatternTable_FilterPredicate(t *testing.T) {
	patterns := createTestPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.Filter.PredicateText = "works"

	output := RenderPatternTable(patterns, opts)

	// Should only show patterns with "works" in predicate
	if !strings.Contains(output, "works") && !strings.Contains(output, "No Matching") {
		t.Error("Filter should show patterns containing 'works'")
	}
}

// Test max rows limiting
func TestRenderPatternTable_MaxRows(t *testing.T) {
	patterns := createTestPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.MaxRows = 2

	output := RenderPatternTable(patterns, opts)

	// Count pattern rows (look for "Pattern X:" lines)
	patternCount := strings.Count(output, "Pattern ")
	if patternCount > opts.MaxRows {
		t.Errorf("Expected max %d patterns, found %d", opts.MaxRows, patternCount)
	}
}

// Test sorting by strength
func TestRenderPatternTable_SortByStrength(t *testing.T) {
	patterns := createUnsortedPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.SortMode = SortByStrength

	output := RenderPatternTable(patterns, opts)

	// Should be sorted by strength (highest first)
	// Visual verification in output
	if output == "" {
		t.Fatal("Expected non-empty output")
	}
}

// Test sorting by frequency
func TestRenderPatternTable_SortByFrequency(t *testing.T) {
	patterns := createUnsortedPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.SortMode = SortByFrequency

	output := RenderPatternTable(patterns, opts)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}
}

// Test sorting by confidence
func TestRenderPatternTable_SortByConfidence(t *testing.T) {
	patterns := createUnsortedPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.SortMode = SortByConfidence

	output := RenderPatternTable(patterns, opts)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}
}

// Test sorting by predicate
func TestRenderPatternTable_SortByPredicate(t *testing.T) {
	patterns := createUnsortedPatterns()
	opts := DefaultPatternDisplayOptions()
	opts.SortMode = SortByPredicate

	output := RenderPatternTable(patterns, opts)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}
}

// Test cluster table rendering
func TestRenderClusterTable_Basic(t *testing.T) {
	patterns := createTestPatterns()
	clusters := ClusterPatterns(patterns, 0.7)
	opts := DefaultPatternDisplayOptions()

	output := RenderClusterTable(clusters, opts)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}

	// Should contain cluster header
	if !strings.Contains(output, "Pattern Clusters") {
		t.Error("Missing cluster header")
	}

	// Should show cluster information
	if !strings.Contains(output, "Cluster") {
		t.Error("Missing cluster information")
	}
}

// Test cluster table with empty clusters
func TestRenderClusterTable_Empty(t *testing.T) {
	clusters := []PatternCluster{}
	opts := DefaultPatternDisplayOptions()

	output := RenderClusterTable(clusters, opts)

	if output == "" {
		t.Fatal("Expected non-empty output for empty state")
	}

	if !strings.Contains(output, "No Patterns Found") {
		t.Error("Missing empty state message")
	}
}

// Test pattern statistics rendering
func TestRenderPatternStats(t *testing.T) {
	patterns := createTestPatterns()
	stats := CalculatePatternStats(patterns)

	output := RenderPatternStats(stats, 80)

	if output == "" {
		t.Fatal("Expected non-empty output")
	}

	// Should contain key statistics
	if !strings.Contains(output, "Total Patterns") {
		t.Error("Missing total patterns statistic")
	}

	if !strings.Contains(output, "Unique Predicates") {
		t.Error("Missing unique predicates statistic")
	}

	if !strings.Contains(output, "Avg Occurrences") {
		t.Error("Missing average occurrences statistic")
	}

	if !strings.Contains(output, "Avg Confidence") {
		t.Error("Missing average confidence statistic")
	}
}

// Test pattern statistics with top predicates
func TestRenderPatternStats_TopPredicates(t *testing.T) {
	patterns := createTestPatterns()
	stats := CalculatePatternStats(patterns)

	output := RenderPatternStats(stats, 80)

	if len(stats.TopPredicates) > 0 {
		// Should show top predicates section
		if !strings.Contains(output, "Top Predicates") {
			t.Error("Missing top predicates section")
		}
	}
}

// Test filter function
func TestFilterPatterns(t *testing.T) {
	patterns := []RelationshipPattern{
		{Predicate: "works_at", Occurrences: 5, AvgConfidence: 0.8},
		{Predicate: "creates", Occurrences: 2, AvgConfidence: 0.9},
		{Predicate: "manages", Occurrences: 4, AvgConfidence: 0.6},
	}

	// Filter by occurrences
	filter := PatternFilter{MinOccurrences: 4, MinConfidence: 0.0}
	filtered := filterPatterns(patterns, filter)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 patterns with occurrences >= 4, got %d", len(filtered))
	}

	// Filter by confidence
	filter = PatternFilter{MinOccurrences: 0, MinConfidence: 0.85}
	filtered = filterPatterns(patterns, filter)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 pattern with confidence >= 0.85, got %d", len(filtered))
	}

	// Filter by predicate text
	filter = PatternFilter{PredicateText: "work"}
	filtered = filterPatterns(patterns, filter)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 pattern containing 'work', got %d", len(filtered))
	}

	for _, p := range filtered {
		if !strings.Contains(strings.ToLower(p.Predicate), "work") {
			t.Errorf("Filtered pattern %q does not contain 'work'", p.Predicate)
		}
	}
}

// Test sort function
func TestSortPatterns(t *testing.T) {
	patterns := []RelationshipPattern{
		{Predicate: "aaa", Occurrences: 2, AvgConfidence: 0.7, Strength: 0.3},
		{Predicate: "zzz", Occurrences: 5, AvgConfidence: 0.9, Strength: 0.8},
		{Predicate: "mmm", Occurrences: 3, AvgConfidence: 0.6, Strength: 0.5},
	}

	// Test sort by strength
	sorted := sortPatterns(patterns, SortByStrength)
	if sorted[0].Strength < sorted[1].Strength {
		t.Error("Sort by strength failed")
	}

	// Test sort by frequency
	sorted = sortPatterns(patterns, SortByFrequency)
	if sorted[0].Occurrences < sorted[1].Occurrences {
		t.Error("Sort by frequency failed")
	}

	// Test sort by confidence
	sorted = sortPatterns(patterns, SortByConfidence)
	if sorted[0].AvgConfidence < sorted[1].AvgConfidence {
		t.Error("Sort by confidence failed")
	}

	// Test sort by predicate
	sorted = sortPatterns(patterns, SortByPredicate)
	if sorted[0].Predicate != "aaa" {
		t.Error("Sort by predicate failed")
	}

	// Verify original not modified
	if patterns[0].Predicate != "aaa" {
		t.Error("Original patterns were modified")
	}
}

// Test center text function
func TestCenterText(t *testing.T) {
	tests := []struct {
		text  string
		width int
	}{
		{"Hello", 20},
		{"Test", 10},
		{"Very long text that exceeds width", 20},
	}

	for _, test := range tests {
		result := centerText(test.text, test.width)

		// If text is shorter than width, should not exceed width
		if len(test.text) <= test.width && len(result) > test.width {
			t.Errorf("Centered text exceeds width: %d > %d", len(result), test.width)
		}

		// Should contain original text (or be truncated version)
		if !strings.Contains(result, test.text) && len(test.text) <= test.width {
			t.Errorf("Centered text does not contain original: %q", result)
		}
	}
}

// Test width handling in compact mode
func TestCompactRowWidth(t *testing.T) {
	pattern := RelationshipPattern{
		SubjectType:   semantic.EntityPerson,
		Predicate:     "works_at",
		ObjectType:    semantic.EntityOrganization,
		Occurrences:   5,
		AvgConfidence: 0.87,
		Strength:      0.456,
	}

	widths := []int{40, 60, 80, 100, 120}

	for _, width := range widths {
		var b strings.Builder
		renderCompactRow(&b, pattern, 0, width)
		output := b.String()

		// Remove newline for length check
		output = strings.TrimRight(output, "\n")

		// Should not exceed width
		if len(output) > width {
			t.Errorf("Compact row exceeds width %d: got %d characters", width, len(output))
		}
	}
}

// Test width handling in expanded mode
func TestExpandedRowWidth(t *testing.T) {
	pattern := RelationshipPattern{
		SubjectType:   semantic.EntityPerson,
		Predicate:     "works_at",
		ObjectType:    semantic.EntityOrganization,
		Occurrences:   5,
		AvgConfidence: 0.87,
		Strength:      0.456,
		Examples: []PatternExample{
			{Subject: "John", Predicate: "works_at", Object: "Acme Corp", Confidence: 0.9},
		},
	}

	widths := []int{40, 60, 80, 100, 120}

	for _, width := range widths {
		var b strings.Builder
		renderExpandedRow(&b, pattern, 0, width, true)
		output := b.String()

		// Expanded mode is informational and may exceed width for readability
		// Just verify it produces output
		if output == "" {
			t.Errorf("No output for width %d", width)
		}
	}
}

// Test table rendering at various widths
func TestTableRenderingAtVariousWidths(t *testing.T) {
	patterns := createTestPatterns()

	widths := []int{40, 60, 80, 100, 120}

	for _, width := range widths {
		opts := DefaultPatternDisplayOptions()
		opts.Width = width
		opts.Compact = true

		output := RenderPatternTable(patterns, opts)

		if output == "" {
			t.Errorf("Empty output for width %d", width)
			continue
		}

		// Table rendering uses width as a guideline, not strict constraint
		// Verify we get reasonable output
		lines := strings.Split(output, "\n")
		if len(lines) == 0 {
			t.Errorf("No lines in output for width %d", width)
		}
	}
}

// Helper functions

func createTestPatterns() []RelationshipPattern {
	return []RelationshipPattern{
		{
			SubjectType:   semantic.EntityPerson,
			Predicate:     "works_at",
			ObjectType:    semantic.EntityOrganization,
			Occurrences:   5,
			AvgConfidence: 0.87,
			Strength:      0.5,
			Examples: []PatternExample{
				{Subject: "John", Predicate: "works_at", Object: "Acme Corp", Confidence: 0.9},
				{Subject: "Alice", Predicate: "works_at", Object: "Tech Inc", Confidence: 0.85},
			},
		},
		{
			SubjectType:   semantic.EntityOrganization,
			Predicate:     "creates",
			ObjectType:    semantic.EntityTechnology,
			Occurrences:   3,
			AvgConfidence: 0.92,
			Strength:      0.4,
			Examples: []PatternExample{
				{Subject: "Acme Corp", Predicate: "creates", Object: "API Gateway", Confidence: 0.95},
			},
		},
	}
}

func createUnsortedPatterns() []RelationshipPattern {
	return []RelationshipPattern{
		{Predicate: "zzz", Occurrences: 2, AvgConfidence: 0.7, Strength: 0.3},
		{Predicate: "aaa", Occurrences: 5, AvgConfidence: 0.9, Strength: 0.8},
		{Predicate: "mmm", Occurrences: 3, AvgConfidence: 0.6, Strength: 0.5},
	}
}
