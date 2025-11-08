package context

import (
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// --- Panel Creation Tests ---

func TestNewContextPanel(t *testing.T) {
	config := DefaultContextPanelConfig()
	panel := New(config)

	if panel == nil {
		t.Fatal("Expected panel to be created")
	}

	if panel.GetAnalysis() != nil {
		t.Error("Expected initial analysis to be nil")
	}

	if panel.GetScrollOffset() != 0 {
		t.Error("Expected initial scroll offset to be 0")
	}
}

func TestDefaultContextPanelConfig(t *testing.T) {
	config := DefaultContextPanelConfig()

	if config.Width == 0 {
		t.Error("Expected default width to be set")
	}

	if config.Height == 0 {
		t.Error("Expected default height to be set")
	}

	if len(config.DefaultExpanded) == 0 {
		t.Error("Expected default expanded sections")
	}

	if len(config.MaxResults) == 0 {
		t.Error("Expected default max results")
	}
}

// --- Analysis Tests ---

func TestSetAnalysis(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	analysis := &semantic.Analysis{
		Content: "Test content",
		Entities: []semantic.Entity{
			{Text: "User", Type: semantic.EntityPerson, Count: 1},
		},
	}

	panel.SetAnalysis(analysis)

	retrieved := panel.GetAnalysis()
	if retrieved == nil {
		t.Fatal("Expected analysis to be set")
	}

	if retrieved.Content != "Test content" {
		t.Error("Expected analysis content to match")
	}
}

func TestSetAnalysisWithTypedHoles(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	analysis := &semantic.Analysis{
		Content: "Test",
		TypedHoles: []semantic.TypedHole{
			{Type: "Function", Span: semantic.Span{Start: 0, End: 10}},
		},
		Relationships: []semantic.Relationship{},
	}

	panel.SetAnalysis(analysis)

	// Should enhance typed holes
	if len(panel.enhancedHoles) != 1 {
		t.Errorf("Expected 1 enhanced hole, got %d", len(panel.enhancedHoles))
	}
}

// --- Section Tests ---

func TestToggleSection(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Initial state (from default config)
	initialState := panel.IsSectionExpanded(SectionEntities)

	// Toggle
	panel.ToggleSection(SectionEntities)

	newState := panel.IsSectionExpanded(SectionEntities)
	if newState == initialState {
		t.Error("Expected section state to toggle")
	}

	// Toggle back
	panel.ToggleSection(SectionEntities)

	finalState := panel.IsSectionExpanded(SectionEntities)
	if finalState != initialState {
		t.Error("Expected section to return to initial state")
	}
}

func TestExpandCollapseSection(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Collapse
	panel.CollapseSection(SectionEntities)
	if panel.IsSectionExpanded(SectionEntities) {
		t.Error("Expected section to be collapsed")
	}

	// Expand
	panel.ExpandSection(SectionEntities)
	if !panel.IsSectionExpanded(SectionEntities) {
		t.Error("Expected section to be expanded")
	}
}

// --- Filtering Tests ---

func TestSetFilterMode(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	panel.SetFilterMode(FilterEntity)
	if panel.filterMode != FilterEntity {
		t.Error("Expected filter mode to be set to FilterEntity")
	}

	panel.SetFilterMode(FilterNone)
	if panel.filterMode != FilterNone {
		t.Error("Expected filter mode to be set to FilterNone")
	}
}

func TestSetFilterQuery(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	panel.SetFilterQuery("test")

	if panel.filterQuery != "test" {
		t.Errorf("Expected filter query 'test', got '%s'", panel.filterQuery)
	}

	if panel.filterMode != FilterSearch {
		t.Error("Expected filter mode to be FilterSearch")
	}
}

func TestSetFilterType(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	panel.SetFilterType(semantic.EntityPerson)

	if panel.filterType != semantic.EntityPerson {
		t.Error("Expected filter type to be EntityPerson")
	}

	if panel.filterMode != FilterEntity {
		t.Error("Expected filter mode to be FilterEntity")
	}
}

func TestClearFilter(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Set filter
	panel.SetFilterQuery("test")

	// Clear
	panel.ClearFilter()

	if panel.filterMode != FilterNone {
		t.Error("Expected filter mode to be FilterNone after clear")
	}

	if panel.filterQuery != "" {
		t.Error("Expected filter query to be empty after clear")
	}
}

func TestFilterEntities(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	entities := []semantic.Entity{
		{Text: "User", Type: semantic.EntityPerson, Count: 1},
		{Text: "System", Type: semantic.EntityConcept, Count: 1},
		{Text: "Server", Type: semantic.EntityPlace, Count: 1},
	}

	// Test entity type filter
	panel.SetFilterType(semantic.EntityPerson)
	filtered := panel.filterEntities(entities)

	if len(filtered) != 1 {
		t.Fatalf("Expected 1 filtered entity, got %d", len(filtered))
	}

	if filtered[0].Text != "User" {
		t.Errorf("Expected 'User', got '%s'", filtered[0].Text)
	}
}

func TestFilterEntitiesSearch(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	entities := []semantic.Entity{
		{Text: "UserAccount", Type: semantic.EntityPerson, Count: 1},
		{Text: "SystemConfig", Type: semantic.EntityConcept, Count: 1},
		{Text: "ServerHost", Type: semantic.EntityPlace, Count: 1},
	}

	// Test search filter (case insensitive)
	panel.SetFilterQuery("user")
	filtered := panel.filterEntities(entities)

	if len(filtered) != 1 {
		t.Fatalf("Expected 1 filtered entity, got %d", len(filtered))
	}

	if filtered[0].Text != "UserAccount" {
		t.Errorf("Expected 'UserAccount', got '%s'", filtered[0].Text)
	}
}

// --- Scrolling Tests ---

func TestScrollDown(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	initialOffset := panel.GetScrollOffset()

	panel.ScrollDown(5)

	newOffset := panel.GetScrollOffset()
	if newOffset != initialOffset+5 {
		t.Errorf("Expected offset to increase by 5, got %d", newOffset-initialOffset)
	}
}

func TestScrollUp(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Scroll down first
	panel.ScrollDown(10)

	// Then scroll up
	panel.ScrollUp(5)

	offset := panel.GetScrollOffset()
	if offset != 5 {
		t.Errorf("Expected offset 5, got %d", offset)
	}
}

func TestScrollUpBounded(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Try to scroll up beyond 0
	panel.ScrollUp(10)

	offset := panel.GetScrollOffset()
	if offset != 0 {
		t.Errorf("Expected offset to be bounded at 0, got %d", offset)
	}
}

func TestScrollToTop(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	panel.ScrollDown(20)
	panel.ScrollToTop()

	offset := panel.GetScrollOffset()
	if offset != 0 {
		t.Errorf("Expected offset 0 after ScrollToTop, got %d", offset)
	}
}

func TestScrollToBottom(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Set total lines
	panel.totalLines = 100

	panel.ScrollToBottom()

	offset := panel.GetScrollOffset()
	if offset != 100 {
		t.Errorf("Expected offset 100, got %d", offset)
	}
}

// --- Section Navigation Tests ---

func TestNextSection(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	initialSection := panel.GetActiveSection()

	panel.NextSection()

	newSection := panel.GetActiveSection()
	if newSection == initialSection {
		t.Error("Expected section to change")
	}
}

func TestPreviousSection(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Move to next section first
	panel.NextSection()
	currentSection := panel.GetActiveSection()

	// Move back
	panel.PreviousSection()

	newSection := panel.GetActiveSection()
	if newSection == currentSection {
		t.Error("Expected section to change")
	}
}

func TestSectionNavigationWraps(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	// Set to first section
	panel.SetActiveSection(SectionEntities)

	// Previous from first should wrap to last
	panel.PreviousSection()

	section := panel.GetActiveSection()
	if section == SectionEntities {
		t.Error("Expected section to wrap around")
	}
}

func TestSetActiveSection(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	panel.SetActiveSection(SectionTypedHoles)

	section := panel.GetActiveSection()
	if section != SectionTypedHoles {
		t.Errorf("Expected SectionTypedHoles, got %v", section)
	}
}

// --- Rendering Tests ---

func TestRenderEmptyPanel(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	output := panel.Render()

	if output == "" {
		t.Error("Expected non-empty output")
	}

	if !strings.Contains(output, "Context Panel") {
		t.Error("Expected output to contain header")
	}
}

func TestRenderWithAnalysis(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	analysis := &semantic.Analysis{
		Content: "User creates Document",
		Entities: []semantic.Entity{
			{Text: "User", Type: semantic.EntityPerson, Count: 1},
			{Text: "Document", Type: semantic.EntityConcept, Count: 1},
		},
		Relationships: []semantic.Relationship{
			{Subject: "User", Predicate: "creates", Object: "Document"},
		},
	}

	panel.SetAnalysis(analysis)

	output := panel.Render()

	if !strings.Contains(output, "User") {
		t.Error("Expected output to contain 'User'")
	}

	if !strings.Contains(output, "Document") {
		t.Error("Expected output to contain 'Document'")
	}
}

func TestRenderSectionCounts(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "User", Type: semantic.EntityPerson, Count: 1},
			{Text: "Document", Type: semantic.EntityConcept, Count: 1},
		},
	}

	panel.SetAnalysis(analysis)

	output := panel.Render()

	// Should show entity count
	if !strings.Contains(output, "(2)") {
		t.Error("Expected output to show entity count")
	}
}

// --- Configuration Tests ---

func TestGetSetConfig(t *testing.T) {
	panel := New(DefaultContextPanelConfig())

	newConfig := DefaultContextPanelConfig()
	newConfig.Width = 60

	panel.SetConfig(newConfig)

	retrieved := panel.GetConfig()
	if retrieved.Width != 60 {
		t.Errorf("Expected width 60, got %d", retrieved.Width)
	}
}
