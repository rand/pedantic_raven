package visualizations

import (
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// TestNewPriorityTable tests table creation.
func TestNewPriorityTable(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
	}

	config := DefaultPriorityTableConfig()
	table := NewPriorityTable(analysis, config)

	if table == nil {
		t.Fatal("Expected non-nil table")
	}

	if table.analysis == nil {
		t.Error("Expected analysis to be set")
	}
}

// TestPriorityTable_Render_Empty tests rendering with no data.
func TestPriorityTable_Render_Empty(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
	}

	config := DefaultPriorityTableConfig()
	table := NewPriorityTable(analysis, config)

	result := table.Render()

	if !strings.Contains(result, "No typed holes") {
		t.Error("Expected 'No typed holes' message")
	}
}

// TestPriorityTable_Render_SingleHole tests rendering with single hole.
func TestPriorityTable_Render_SingleHole(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{
				TypedHole:     semantic.TypedHole{Type: "AuthService", Constraint: "thread-safe"},
				Priority:      8,
				Complexity:    5,
				SuggestedImpl: "Use sync.Mutex",
			},
		},
	}

	config := DefaultPriorityTableConfig()
	table := NewPriorityTable(analysis, config)

	result := table.Render()

	if !strings.Contains(result, "Priority Queue") {
		t.Error("Expected table header")
	}

	if !strings.Contains(result, "AuthService") {
		t.Error("Expected hole type in output")
	}
}

// TestPriorityTable_Render_MultipleHoles tests rendering with multiple holes.
func TestPriorityTable_Render_MultipleHoles(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "High"}, Priority: 9, Complexity: 3},
			{TypedHole: semantic.TypedHole{Type: "Medium"}, Priority: 6, Complexity: 5},
			{TypedHole: semantic.TypedHole{Type: "Low"}, Priority: 3, Complexity: 8},
		},
	}

	config := DefaultPriorityTableConfig()
	config.SortBy = "priority"
	table := NewPriorityTable(analysis, config)

	result := table.Render()

	if !strings.Contains(result, "High") {
		t.Error("Expected high priority hole")
	}

	if !strings.Contains(result, "Medium") {
		t.Error("Expected medium priority hole")
	}

	if !strings.Contains(result, "Low") {
		t.Error("Expected low priority hole")
	}
}

// TestFilterAndSort_ByPriority tests sorting by priority.
func TestFilterAndSort_ByPriority(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "Low"}, Priority: 3, Complexity: 5},
			{TypedHole: semantic.TypedHole{Type: "High"}, Priority: 9, Complexity: 5},
			{TypedHole: semantic.TypedHole{Type: "Medium"}, Priority: 6, Complexity: 5},
		},
	}

	config := DefaultPriorityTableConfig()
	config.SortBy = "priority"
	table := NewPriorityTable(analysis, config)

	sorted := table.filterAndSort()

	if len(sorted) != 3 {
		t.Fatalf("Expected 3 holes, got %d", len(sorted))
	}

	// Should be sorted by priority descending
	if sorted[0].Priority != 9 {
		t.Errorf("Expected first hole to have priority 9, got %d", sorted[0].Priority)
	}

	if sorted[2].Priority != 3 {
		t.Errorf("Expected last hole to have priority 3, got %d", sorted[2].Priority)
	}
}

// TestFilterAndSort_ByComplexity tests sorting by complexity.
func TestFilterAndSort_ByComplexity(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "High"}, Priority: 5, Complexity: 9},
			{TypedHole: semantic.TypedHole{Type: "Low"}, Priority: 5, Complexity: 2},
			{TypedHole: semantic.TypedHole{Type: "Medium"}, Priority: 5, Complexity: 5},
		},
	}

	config := DefaultPriorityTableConfig()
	config.SortBy = "complexity"
	table := NewPriorityTable(analysis, config)

	sorted := table.filterAndSort()

	if len(sorted) != 3 {
		t.Fatalf("Expected 3 holes, got %d", len(sorted))
	}

	// Should be sorted by complexity ascending
	if sorted[0].Complexity != 2 {
		t.Errorf("Expected first hole to have complexity 2, got %d", sorted[0].Complexity)
	}

	if sorted[2].Complexity != 9 {
		t.Errorf("Expected last hole to have complexity 9, got %d", sorted[2].Complexity)
	}
}

// TestFilterAndSort_ByRecommended tests recommended order sorting.
func TestFilterAndSort_ByRecommended(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "A"}, Priority: 8, Complexity: 8}, // ratio: 1.0
			{TypedHole: semantic.TypedHole{Type: "B"}, Priority: 9, Complexity: 3}, // ratio: 3.0 (best)
			{TypedHole: semantic.TypedHole{Type: "C"}, Priority: 4, Complexity: 8}, // ratio: 0.5
		},
	}

	config := DefaultPriorityTableConfig()
	config.SortBy = "recommended"
	table := NewPriorityTable(analysis, config)

	sorted := table.filterAndSort()

	if len(sorted) != 3 {
		t.Fatalf("Expected 3 holes, got %d", len(sorted))
	}

	// Should be sorted by priority/complexity ratio (highest first)
	if sorted[0].Type != "B" {
		t.Errorf("Expected 'B' first (best ratio), got %s", sorted[0].Type)
	}
}

// TestFilterAndSort_WithFilter tests filtering by type.
func TestFilterAndSort_WithFilter(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "Service"}, Priority: 8},
			{TypedHole: semantic.TypedHole{Type: "Handler"}, Priority: 7},
			{TypedHole: semantic.TypedHole{Type: "Service"}, Priority: 6},
		},
	}

	config := DefaultPriorityTableConfig()
	config.FilterType = "Service"
	table := NewPriorityTable(analysis, config)

	filtered := table.filterAndSort()

	if len(filtered) != 2 {
		t.Errorf("Expected 2 Service holes, got %d", len(filtered))
	}

	for _, hole := range filtered {
		if hole.Type != "Service" {
			t.Errorf("Expected only Service type, got %s", hole.Type)
		}
	}
}

// TestFilterAndSort_WithMinPriority tests filtering by minimum priority.
func TestFilterAndSort_WithMinPriority(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "High"}, Priority: 9},
			{TypedHole: semantic.TypedHole{Type: "Medium"}, Priority: 5},
			{TypedHole: semantic.TypedHole{Type: "Low"}, Priority: 2},
		},
	}

	config := DefaultPriorityTableConfig()
	config.MinPriority = 5
	table := NewPriorityTable(analysis, config)

	filtered := table.filterAndSort()

	if len(filtered) != 2 {
		t.Errorf("Expected 2 holes with priority >= 5, got %d", len(filtered))
	}

	for _, hole := range filtered {
		if hole.Priority < 5 {
			t.Errorf("Expected priority >= 5, got %d", hole.Priority)
		}
	}
}

// TestRenderRow tests row rendering.
func TestRenderRow(t *testing.T) {
	hole := semantic.EnhancedTypedHole{
		TypedHole:     semantic.TypedHole{Type: "TestService", Constraint: "async"},
		Priority:      7,
		Complexity:    5,
		SuggestedImpl: "Use goroutines",
	}

	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{hole},
	}

	config := DefaultPriorityTableConfig()
	table := NewPriorityTable(analysis, config)

	row := table.renderRow(1, hole, 4, 20, 10, 12, 20, 30)

	if row == "" {
		t.Error("Expected non-empty row")
	}

	// Check for hole name
	if !strings.Contains(row, "TestService") {
		t.Error("Expected hole type in row")
	}
}

// TestRenderPriorityBar tests priority bar rendering.
func TestRenderPriorityBar(t *testing.T) {
	config := DefaultPriorityTableConfig()
	table := &PriorityTable{config: config}

	bar := table.renderPriorityBar(5)

	if !strings.Contains(bar, "5/10") {
		t.Error("Expected priority value in bar")
	}

	if !strings.Contains(bar, "█") && !strings.Contains(bar, "░") {
		t.Error("Expected bar characters")
	}
}

// TestRenderComplexityBar tests complexity bar rendering.
func TestRenderComplexityBar(t *testing.T) {
	config := DefaultPriorityTableConfig()
	table := &PriorityTable{config: config}

	bar := table.renderComplexityBar(7)

	if !strings.Contains(bar, "7/10") {
		t.Error("Expected complexity value in bar")
	}

	if !strings.Contains(bar, "█") && !strings.Contains(bar, "░") {
		t.Error("Expected bar characters")
	}
}

// TestGetPriorityColor tests color assignment.
func TestGetPriorityColor(t *testing.T) {
	config := DefaultPriorityTableConfig()
	table := &PriorityTable{config: config}

	// High priority
	highColor := table.getPriorityColor(9)
	if highColor != "196" {
		t.Errorf("Expected red for high priority, got %s", highColor)
	}

	// Medium priority
	medColor := table.getPriorityColor(6)
	if medColor != "226" {
		t.Errorf("Expected yellow for medium priority, got %s", medColor)
	}

	// Low priority
	lowColor := table.getPriorityColor(3)
	if lowColor != "34" {
		t.Errorf("Expected green for low priority, got %s", lowColor)
	}
}

// TestGetComplexityColor tests complexity color assignment.
func TestGetComplexityColor(t *testing.T) {
	config := DefaultPriorityTableConfig()
	table := &PriorityTable{config: config}

	// High complexity
	highColor := table.getComplexityColor(9)
	if highColor != "196" {
		t.Errorf("Expected red for high complexity, got %s", highColor)
	}

	// Medium complexity
	medColor := table.getComplexityColor(6)
	if medColor != "226" {
		t.Errorf("Expected yellow for medium complexity, got %s", medColor)
	}

	// Low complexity
	lowColor := table.getComplexityColor(3)
	if lowColor != "34" {
		t.Errorf("Expected green for low complexity, got %s", lowColor)
	}
}

// TestPadRight tests string padding.
func TestPadRight(t *testing.T) {
	config := DefaultPriorityTableConfig()
	table := &PriorityTable{config: config}

	result := table.padRight("test", 10)

	visualLen := len(stripAnsi(result))
	if visualLen < 10 {
		t.Errorf("Expected padded length >= 10, got %d", visualLen)
	}
}

// TestStripAnsi tests ANSI code stripping.
func TestStripAnsi(t *testing.T) {
	// Plain text
	plain := "hello"
	result := stripAnsi(plain)
	if result != plain {
		t.Errorf("Expected '%s', got '%s'", plain, result)
	}

	// Text with ANSI codes (simplified)
	ansi := "\x1b[31mred\x1b[0m"
	result = stripAnsi(ansi)
	if result != "red" {
		t.Errorf("Expected 'red', got '%s'", result)
	}
}

// TestRenderFooter tests footer rendering.
func TestRenderFooter(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "A"}},
			{TypedHole: semantic.TypedHole{Type: "B"}},
			{TypedHole: semantic.TypedHole{Type: "C"}},
		},
	}

	config := DefaultPriorityTableConfig()
	table := NewPriorityTable(analysis, config)

	footer := table.renderFooter(2)

	if !strings.Contains(footer, "2 of 3") {
		t.Error("Expected display count in footer")
	}
}

// TestRenderCompact tests compact rendering.
func TestRenderCompact(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "A"}, Priority: 9},
			{TypedHole: semantic.TypedHole{Type: "B"}, Priority: 6},
			{TypedHole: semantic.TypedHole{Type: "C"}, Priority: 3},
		},
		TotalComplexity: 15,
	}

	config := DefaultPriorityTableConfig()
	table := NewPriorityTable(analysis, config)

	result := table.RenderCompact()

	if !strings.Contains(result, "high") {
		t.Error("Expected high priority count")
	}

	if !strings.Contains(result, "medium") {
		t.Error("Expected medium priority count")
	}

	if !strings.Contains(result, "low") {
		t.Error("Expected low priority count")
	}

	if !strings.Contains(result, "15") {
		t.Error("Expected total complexity")
	}
}

// TestRenderCompact_Empty tests compact rendering with no holes.
func TestRenderCompact_Empty(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
	}

	config := DefaultPriorityTableConfig()
	table := NewPriorityTable(analysis, config)

	result := table.RenderCompact()

	if !strings.Contains(result, "No typed holes") {
		t.Error("Expected 'No typed holes' message")
	}
}

// TestRenderByType tests grouping by type.
func TestRenderByType(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{
			{TypedHole: semantic.TypedHole{Type: "Service"}, Priority: 8},
			{TypedHole: semantic.TypedHole{Type: "Service"}, Priority: 6},
			{TypedHole: semantic.TypedHole{Type: "Handler"}, Priority: 7},
		},
	}

	result := RenderByType(analysis)

	if !strings.Contains(result, "Grouped by Type") {
		t.Error("Expected header")
	}

	if !strings.Contains(result, "Service") {
		t.Error("Expected Service group")
	}

	if !strings.Contains(result, "Handler") {
		t.Error("Expected Handler group")
	}

	if !strings.Contains(result, "2 holes") {
		t.Error("Expected hole count for Service")
	}
}

// TestRenderByType_Empty tests grouping with no holes.
func TestRenderByType_Empty(t *testing.T) {
	analysis := &analyze.HoleAnalysis{
		Holes: []semantic.EnhancedTypedHole{},
	}

	result := RenderByType(analysis)

	if !strings.Contains(result, "No typed holes") {
		t.Error("Expected 'No typed holes' message")
	}
}

// TestDefaultPriorityTableConfig tests default configuration.
func TestDefaultPriorityTableConfig(t *testing.T) {
	config := DefaultPriorityTableConfig()

	if config.Width <= 0 {
		t.Error("Expected positive width")
	}

	if config.MaxRows <= 0 {
		t.Error("Expected positive max rows")
	}

	if !config.ShowConstraints {
		t.Error("Expected constraints to be shown by default")
	}

	if !config.ShowSuggestions {
		t.Error("Expected suggestions to be shown by default")
	}

	if config.SortBy != "recommended" {
		t.Error("Expected default sort to be recommended")
	}
}

// TestMaxRows tests row limiting.
func TestMaxRows(t *testing.T) {
	holes := []semantic.EnhancedTypedHole{}
	for i := 0; i < 30; i++ {
		holes = append(holes, semantic.EnhancedTypedHole{
			TypedHole: semantic.TypedHole{Type: "Hole"},
			Priority:  5,
		})
	}

	analysis := &analyze.HoleAnalysis{
		Holes: holes,
	}

	config := DefaultPriorityTableConfig()
	config.MaxRows = 10
	table := NewPriorityTable(analysis, config)

	result := table.Render()

	// Count rows in output (count "??Hole" which appears once per row)
	count := strings.Count(result, "??Hole")
	if count > config.MaxRows {
		t.Errorf("Expected at most %d rows, got %d", config.MaxRows, count)
	}
}
