package layout

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// MockComponent is a test component that implements the Component interface.
type MockComponent struct {
	id      PaneID
	content string
	updated bool
}

func NewMockComponent(id PaneID, content string) *MockComponent {
	return &MockComponent{
		id:      id,
		content: content,
		updated: false,
	}
}

func (m *MockComponent) Update(msg tea.Msg) (Component, tea.Cmd) {
	m.updated = true
	return m, nil
}

func (m *MockComponent) View(area Rect, focused bool) string {
	focusIndicator := ""
	if focused {
		focusIndicator = "[FOCUSED] "
	}
	return fmt.Sprintf("%s%s (%dx%d at %d,%d)",
		focusIndicator, m.content, area.Width, area.Height, area.X, area.Y)
}

func (m *MockComponent) ID() PaneID {
	return m.id
}

// --- Pane Tests ---

func TestLeafPaneRender(t *testing.T) {
	component := NewMockComponent(PaneEditor, "Editor")
	pane := NewLeafPane(component)

	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}
	view := pane.Render(area, PaneEditor)

	expected := "[FOCUSED] Editor (100x30 at 0,0)"
	if view != expected {
		t.Errorf("Expected '%s', got '%s'", expected, view)
	}
}

func TestLeafPaneRenderUnfocused(t *testing.T) {
	component := NewMockComponent(PaneEditor, "Editor")
	pane := NewLeafPane(component)

	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}
	view := pane.Render(area, PaneSidebar) // Different focus

	expected := "Editor (100x30 at 0,0)"
	if view != expected {
		t.Errorf("Expected '%s', got '%s'", expected, view)
	}
}

func TestLeafPaneUpdate(t *testing.T) {
	component := NewMockComponent(PaneEditor, "Editor")
	pane := NewLeafPane(component)

	// Update with focus
	updatedPane, _ := pane.Update(tea.KeyMsg{}, PaneEditor)

	leafPane := updatedPane.(*LeafPane)
	mockComp := leafPane.component.(*MockComponent)

	if !mockComp.updated {
		t.Error("Component should be updated when focused")
	}
}

func TestLeafPaneFindComponent(t *testing.T) {
	component := NewMockComponent(PaneEditor, "Editor")
	pane := NewLeafPane(component)

	found := pane.FindComponent(PaneEditor)
	if found == nil {
		t.Fatal("Component should be found")
	}

	if found.ID() != PaneEditor {
		t.Errorf("Expected ID %s, got %s", PaneEditor, found.ID())
	}

	notFound := pane.FindComponent(PaneSidebar)
	if notFound != nil {
		t.Error("Should not find non-existent component")
	}
}

func TestLeafPaneAllPaneIDs(t *testing.T) {
	component := NewMockComponent(PaneEditor, "Editor")
	pane := NewLeafPane(component)

	ids := pane.AllPaneIDs()
	if len(ids) != 1 {
		t.Fatalf("Expected 1 pane ID, got %d", len(ids))
	}

	if ids[0] != PaneEditor {
		t.Errorf("Expected ID %s, got %s", PaneEditor, ids[0])
	}
}

func TestSplitPaneHorizontal(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Horizontal,
		0.7, // 70/30 split
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}
	firstArea, secondArea := split.computeSplitAreas(area)

	// First child should get 70% of width (70 cols)
	if firstArea.Width != 70 {
		t.Errorf("Expected first width 70, got %d", firstArea.Width)
	}

	// Second child should get 30% of width (30 cols)
	if secondArea.Width != 30 {
		t.Errorf("Expected second width 30, got %d", secondArea.Width)
	}

	// Both should have full height
	if firstArea.Height != 30 || secondArea.Height != 30 {
		t.Errorf("Both panes should have full height, got %d and %d",
			firstArea.Height, secondArea.Height)
	}

	// Second pane should start where first ends
	if secondArea.X != firstArea.Width {
		t.Errorf("Second pane X should be %d, got %d", firstArea.Width, secondArea.X)
	}
}

func TestSplitPaneVertical(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	diagnostics := NewMockComponent(PaneDiagnostics, "Diagnostics")

	split := NewSplitPane(
		Vertical,
		0.8, // 80/20 split
		NewLeafPane(editor),
		NewLeafPane(diagnostics),
	)

	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}
	firstArea, secondArea := split.computeSplitAreas(area)

	// First child should get 80% of height (24 rows)
	if firstArea.Height != 24 {
		t.Errorf("Expected first height 24, got %d", firstArea.Height)
	}

	// Second child should get 20% of height (6 rows)
	if secondArea.Height != 6 {
		t.Errorf("Expected second height 6, got %d", secondArea.Height)
	}

	// Both should have full width
	if firstArea.Width != 100 || secondArea.Width != 100 {
		t.Errorf("Both panes should have full width, got %d and %d",
			firstArea.Width, secondArea.Width)
	}

	// Second pane should start where first ends
	if secondArea.Y != firstArea.Height {
		t.Errorf("Second pane Y should be %d, got %d", firstArea.Height, secondArea.Y)
	}
}

func TestSplitPaneFindComponent(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	// Find editor
	found := split.FindComponent(PaneEditor)
	if found == nil {
		t.Fatal("Should find editor")
	}
	if found.ID() != PaneEditor {
		t.Errorf("Expected editor ID, got %s", found.ID())
	}

	// Find sidebar
	found = split.FindComponent(PaneSidebar)
	if found == nil {
		t.Fatal("Should find sidebar")
	}
	if found.ID() != PaneSidebar {
		t.Errorf("Expected sidebar ID, got %s", found.ID())
	}

	// Don't find non-existent
	found = split.FindComponent(PaneMemory)
	if found != nil {
		t.Error("Should not find non-existent component")
	}
}

func TestSplitPaneAllPaneIDs(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	ids := split.AllPaneIDs()
	if len(ids) != 2 {
		t.Fatalf("Expected 2 pane IDs, got %d", len(ids))
	}

	// Check both IDs are present (order doesn't matter for this test)
	hasEditor := false
	hasSidebar := false
	for _, id := range ids {
		if id == PaneEditor {
			hasEditor = true
		}
		if id == PaneSidebar {
			hasSidebar = true
		}
	}

	if !hasEditor || !hasSidebar {
		t.Errorf("Expected both editor and sidebar IDs, got %v", ids)
	}
}

func TestSplitPaneRatioClamp(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	// Test ratio < 0.1 (should clamp to 0.1)
	split := NewSplitPane(
		Horizontal,
		0.05,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	if split.ratio != 0.1 {
		t.Errorf("Ratio should be clamped to 0.1, got %f", split.ratio)
	}

	// Test ratio > 0.9 (should clamp to 0.9)
	split = NewSplitPane(
		Horizontal,
		0.95,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	if split.ratio != 0.9 {
		t.Errorf("Ratio should be clamped to 0.9, got %f", split.ratio)
	}
}

// --- Engine Tests ---

func TestEngineRegisterComponent(t *testing.T) {
	engine := NewEngine(LayoutFocus)

	component := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(component)

	found := engine.GetComponent(PaneEditor)
	if found == nil {
		t.Fatal("Component should be registered")
	}

	if found.ID() != PaneEditor {
		t.Errorf("Expected editor ID, got %s", found.ID())
	}
}

func TestEngineUnregisterComponent(t *testing.T) {
	engine := NewEngine(LayoutFocus)

	component := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(component)

	engine.UnregisterComponent(PaneEditor)

	found := engine.GetComponent(PaneEditor)
	if found != nil {
		t.Error("Component should be unregistered")
	}
}

func TestEngineFocusManagement(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(sidebar)

	// Build layout
	root := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	engine.SetRoot(root)

	// Focus should default to first pane (editor)
	if engine.FocusedID() != PaneEditor {
		t.Errorf("Expected initial focus on editor, got %s", engine.FocusedID())
	}

	// Set focus to sidebar
	success := engine.SetFocus(PaneSidebar)
	if !success {
		t.Fatal("SetFocus should succeed for existing pane")
	}

	if engine.FocusedID() != PaneSidebar {
		t.Errorf("Expected focus on sidebar, got %s", engine.FocusedID())
	}

	// Try to focus non-existent pane
	success = engine.SetFocus(PaneMemory)
	if success {
		t.Error("SetFocus should fail for non-existent pane")
	}

	// Focus should remain on sidebar
	if engine.FocusedID() != PaneSidebar {
		t.Errorf("Focus should remain on sidebar, got %s", engine.FocusedID())
	}
}

func TestEngineFocusNext(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(sidebar)

	root := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	engine.SetRoot(root)

	// Should start at editor
	if engine.FocusedID() != PaneEditor {
		t.Fatalf("Expected initial focus on editor, got %s", engine.FocusedID())
	}

	// Next should move to sidebar
	engine.FocusNext()
	if engine.FocusedID() != PaneSidebar {
		t.Errorf("Expected focus on sidebar, got %s", engine.FocusedID())
	}

	// Next should wrap to editor
	engine.FocusNext()
	if engine.FocusedID() != PaneEditor {
		t.Errorf("Expected focus to wrap to editor, got %s", engine.FocusedID())
	}
}

func TestEngineFocusPrev(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(sidebar)

	root := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	engine.SetRoot(root)

	// Should start at editor
	if engine.FocusedID() != PaneEditor {
		t.Fatalf("Expected initial focus on editor, got %s", engine.FocusedID())
	}

	// Prev should wrap to sidebar
	engine.FocusPrev()
	if engine.FocusedID() != PaneSidebar {
		t.Errorf("Expected focus to wrap to sidebar, got %s", engine.FocusedID())
	}

	// Prev should move to editor
	engine.FocusPrev()
	if engine.FocusedID() != PaneEditor {
		t.Errorf("Expected focus on editor, got %s", engine.FocusedID())
	}
}

func TestEngineTerminalResize(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	// Start with normal size
	engine.SetTerminalSize(120, 30)
	width, height := engine.TerminalSize()
	if width != 120 || height != 30 {
		t.Errorf("Expected 120x30, got %dx%d", width, height)
	}

	// Resize to small (should switch to compact mode)
	engine.SetTerminalSize(80, 20)
	if engine.Mode() != LayoutCompact {
		t.Errorf("Expected compact mode for small terminal, got %s", engine.Mode())
	}
}

func TestEngineAllPaneIDs(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	root := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	engine.SetRoot(root)

	ids := engine.AllPaneIDs()
	if len(ids) != 2 {
		t.Fatalf("Expected 2 pane IDs, got %d", len(ids))
	}
}

func TestDirectionString(t *testing.T) {
	tests := []struct {
		direction Direction
		expected  string
	}{
		{Horizontal, "Horizontal"},
		{Vertical, "Vertical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.direction.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestLayoutModeString(t *testing.T) {
	tests := []struct {
		mode     LayoutMode
		expected string
	}{
		{LayoutFocus, "Focus"},
		{LayoutStandard, "Standard"},
		{LayoutAnalysis, "Analysis"},
		{LayoutCompact, "Compact"},
		{LayoutCustom, "Custom"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.mode.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// --- Additional Coverage Tests ---

// TestEngineRoot tests the Root() getter
func TestEngineRoot(t *testing.T) {
	engine := NewEngine(LayoutFocus)

	// Initially, root should be nil
	if engine.Root() != nil {
		t.Fatal("Root should be nil initially")
	}

	// After setting root, should return the same root
	editor := NewMockComponent(PaneEditor, "Editor")
	root := NewLeafPane(editor)
	engine.SetRoot(root)

	retrieved := engine.Root()
	if retrieved == nil {
		t.Fatal("Root should not be nil after SetRoot")
	}
	if retrieved != root {
		t.Error("Root should be the same pane that was set")
	}
}

// TestEngineInit tests the Init() method
func TestEngineInit(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	cmd := engine.Init()
	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

// TestEngineUpdateWindowSizeMsg tests engine Update with WindowSizeMsg
func TestEngineUpdateWindowSizeMsg(t *testing.T) {
	engine := NewEngine(LayoutStandard)
	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)
	engine.SetRoot(NewLeafPane(editor))

	// Send a WindowSizeMsg
	msg := tea.WindowSizeMsg{Width: 150, Height: 40}
	model, cmd := engine.Update(msg)

	if cmd != nil {
		t.Error("Update with WindowSizeMsg should return nil command")
	}

	// Check that engine size was updated
	width, height := model.(*Engine).TerminalSize()
	if width != 150 || height != 40 {
		t.Errorf("Expected 150x40, got %dx%d", width, height)
	}
}

// TestEngineUpdateOtherMessage tests engine Update with other message types
func TestEngineUpdateOtherMessage(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)
	engine.SetRoot(NewLeafPane(editor))

	// Send a key message
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	model, cmd := engine.Update(msg)

	if model == nil {
		t.Error("Update should return a model")
	}
	if cmd != nil {
		t.Error("Update with key message should return nil (MockComponent returns nil)")
	}
}

// TestEngineUpdateWithNoRoot tests Update when root is nil
func TestEngineUpdateWithNoRoot(t *testing.T) {
	engine := NewEngine(LayoutFocus)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	model, cmd := engine.Update(msg)

	if cmd != nil {
		t.Error("Update with no root should return nil command")
	}
	if model == nil {
		t.Error("Update should still return the engine")
	}
}

// TestEngineView tests the View() method
func TestEngineView(t *testing.T) {
	engine := NewEngine(LayoutFocus)

	// Without root, should show placeholder
	view := engine.View()
	if view != "No layout configured" {
		t.Errorf("Expected 'No layout configured', got '%s'", view)
	}

	// With root, should render it
	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)
	engine.SetRoot(NewLeafPane(editor))
	engine.SetTerminalSize(100, 30)

	view = engine.View()
	if view == "No layout configured" {
		t.Error("Should render the layout, not placeholder")
	}
	if view == "" {
		t.Error("Should render something from the component")
	}
}

// TestEngineSetModeRebuild tests that SetMode rebuilds layout
func TestEngineSetModeRebuild(t *testing.T) {
	engine := NewEngine(LayoutFocus)

	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")
	analysis := NewMockComponent(PaneAnalysis, "Analysis")
	diagnostics := NewMockComponent(PaneDiagnostics, "Diagnostics")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(sidebar)
	engine.RegisterComponent(analysis)
	engine.RegisterComponent(diagnostics)

	// Set to Focus mode
	engine.SetMode(LayoutFocus)
	if engine.Mode() != LayoutFocus {
		t.Errorf("Expected Focus mode, got %s", engine.Mode())
	}

	// Switch to Standard mode
	engine.SetMode(LayoutStandard)
	if engine.Mode() != LayoutStandard {
		t.Errorf("Expected Standard mode, got %s", engine.Mode())
	}

	// Switch to Analysis mode
	engine.SetMode(LayoutAnalysis)
	if engine.Mode() != LayoutAnalysis {
		t.Errorf("Expected Analysis mode, got %s", engine.Mode())
	}
}

// TestEngineBuildFocusLayout tests buildFocusLayout with no editor
func TestEngineBuildFocusLayoutNoEditor(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	// Don't register editor - root should stay nil
	engine.buildFocusLayout()
	if engine.Root() != nil {
		t.Error("Root should be nil if editor not registered")
	}
}

// TestEngineBuildStandardLayoutPartial tests buildStandardLayout with only editor
func TestEngineBuildStandardLayoutOnlyEditor(t *testing.T) {
	engine := NewEngine(LayoutStandard)
	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)

	engine.buildStandardLayout()

	if engine.Root() == nil {
		t.Fatal("Root should not be nil")
	}

	// Should be a leaf pane (only editor, no sidebar)
	if !engine.Root().IsLeaf() {
		t.Error("Root should be a leaf pane when only editor is registered")
	}
}

// TestEngineBuildStandardLayoutNoComponents tests buildStandardLayout with no components
func TestEngineBuildStandardLayoutNoComponents(t *testing.T) {
	engine := NewEngine(LayoutStandard)
	engine.buildStandardLayout()

	if engine.Root() != nil {
		t.Error("Root should be nil if no components registered")
	}
}

// TestEngineBuildAnalysisLayoutPartial tests buildAnalysisLayout with only editor
func TestEngineBuildAnalysisLayoutOnlyEditor(t *testing.T) {
	engine := NewEngine(LayoutAnalysis)
	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)

	engine.buildAnalysisLayout()

	if engine.Root() == nil {
		t.Fatal("Root should not be nil")
	}

	if !engine.Root().IsLeaf() {
		t.Error("Root should be a leaf pane when only editor is registered")
	}
}

// TestEngineBuildAnalysisLayoutWithDiagnostics tests buildAnalysisLayout with all components
func TestEngineBuildAnalysisLayoutFull(t *testing.T) {
	engine := NewEngine(LayoutAnalysis)
	editor := NewMockComponent(PaneEditor, "Editor")
	analysis := NewMockComponent(PaneAnalysis, "Analysis")
	diagnostics := NewMockComponent(PaneDiagnostics, "Diagnostics")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(analysis)
	engine.RegisterComponent(diagnostics)

	engine.buildAnalysisLayout()

	if engine.Root() == nil {
		t.Fatal("Root should not be nil")
	}

	// Root should be a split pane
	if engine.Root().IsLeaf() {
		t.Error("Root should be a split pane for full analysis layout")
	}

	// Should be able to find all components
	if engine.Root().FindComponent(PaneEditor) == nil {
		t.Error("Should find editor")
	}
	if engine.Root().FindComponent(PaneAnalysis) == nil {
		t.Error("Should find analysis")
	}
	if engine.Root().FindComponent(PaneDiagnostics) == nil {
		t.Error("Should find diagnostics")
	}
}

// TestEngineBuildAnalysisLayoutNoDiagnostics tests buildAnalysisLayout without diagnostics
func TestEngineBuildAnalysisLayoutNoDiagnostics(t *testing.T) {
	engine := NewEngine(LayoutAnalysis)
	editor := NewMockComponent(PaneEditor, "Editor")
	analysis := NewMockComponent(PaneAnalysis, "Analysis")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(analysis)

	engine.buildAnalysisLayout()

	if engine.Root() == nil {
		t.Fatal("Root should not be nil")
	}

	// Root should be a split pane
	if engine.Root().IsLeaf() {
		t.Error("Root should be a split pane")
	}

	// Should find editor and analysis, but not diagnostics
	if engine.Root().FindComponent(PaneEditor) == nil {
		t.Error("Should find editor")
	}
	if engine.Root().FindComponent(PaneAnalysis) == nil {
		t.Error("Should find analysis")
	}
}

// TestEngineBuildAnalysisLayoutNoComponents tests buildAnalysisLayout with no components
func TestEngineBuildAnalysisLayoutNoComponents(t *testing.T) {
	engine := NewEngine(LayoutAnalysis)
	engine.buildAnalysisLayout()

	if engine.Root() != nil {
		t.Error("Root should be nil if no components registered")
	}
}

// TestEngineFocusNextWithNoRoot tests FocusNext with nil root
func TestEngineFocusNextWithNoRoot(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	// Should not panic
	engine.FocusNext()
	if engine.FocusedID() != "" {
		t.Error("FocusedID should remain empty when no root")
	}
}

// TestEngineFocusNextWithNoPanes tests FocusNext with empty pane list
func TestEngineFocusNextWithEmptyLayout(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)
	// Create a root but it won't have the editor component registered in panes
	root := NewLeafPane(editor)
	engine.SetRoot(root)
	engine.focusedID = "" // Reset focused ID

	// FocusNext should handle empty or missing initial focus
	engine.FocusNext()
	// Should wrap around based on available panes
}

// TestEngineFocusPrevWithNoRoot tests FocusPrev with nil root
func TestEngineFocusPrevWithNoRoot(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	// Should not panic
	engine.FocusPrev()
	if engine.FocusedID() != "" {
		t.Error("FocusedID should remain empty when no root")
	}
}

// TestEngineFocusPrevWithEmptyLayout tests FocusPrev with empty pane list
func TestEngineFocusPrevWithEmptyLayout(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)
	root := NewLeafPane(editor)
	engine.SetRoot(root)
	engine.focusedID = "" // Reset focused ID

	// FocusPrev should handle empty or missing initial focus
	engine.FocusPrev()
}

// TestEngineRegisterNilComponent tests RegisterComponent with nil
func TestEngineRegisterNilComponent(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	// Should not panic when registering nil
	engine.RegisterComponent(nil)
	if len(engine.components) != 0 {
		t.Error("Should not register nil components")
	}
}

// TestEngineSetFocusWithNoRoot tests SetFocus when root is nil
func TestEngineSetFocusWithNoRoot(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	success := engine.SetFocus(PaneEditor)
	if success {
		t.Error("SetFocus should fail when root is nil")
	}
}

// TestLeafPaneRenderNilComponent tests LeafPane.Render with nil component
func TestLeafPaneRenderNilComponent(t *testing.T) {
	pane := NewLeafPane(nil)
	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}
	view := pane.Render(area, PaneEditor)

	if view != "" {
		t.Errorf("Expected empty string for nil component, got '%s'", view)
	}
}

// TestLeafPaneUpdateNilComponent tests LeafPane.Update with nil component
func TestLeafPaneUpdateNilComponent(t *testing.T) {
	pane := NewLeafPane(nil)
	updatedPane, cmd := pane.Update(tea.KeyMsg{}, PaneEditor)

	if updatedPane != pane {
		t.Error("Should return same pane when component is nil")
	}
	if cmd != nil {
		t.Error("Should return nil command when component is nil")
	}
}

// TestLeafPaneUpdateUnfocused tests LeafPane.Update when not focused
func TestLeafPaneUpdateUnfocused(t *testing.T) {
	component := NewMockComponent(PaneEditor, "Editor")
	pane := NewLeafPane(component)

	// Update with different focus
	updatedPane, _ := pane.Update(tea.KeyMsg{}, PaneSidebar)
	leafPane := updatedPane.(*LeafPane)
	mockComp := leafPane.component.(*MockComponent)

	if mockComp.updated {
		t.Error("Component should not be updated when not focused")
	}
}

// TestLeafPaneAllPaneIDsNilComponent tests LeafPane.AllPaneIDs with nil component
func TestLeafPaneAllPaneIDsNilComponent(t *testing.T) {
	pane := NewLeafPane(nil)
	ids := pane.AllPaneIDs()

	if len(ids) != 0 {
		t.Errorf("Expected 0 pane IDs, got %d", len(ids))
	}
}

// TestLeafPaneFindComponentNilComponent tests LeafPane.FindComponent with nil component
func TestLeafPaneFindComponentNilComponent(t *testing.T) {
	pane := NewLeafPane(nil)
	found := pane.FindComponent(PaneEditor)

	if found != nil {
		t.Error("Should not find component when pane's component is nil")
	}
}

// TestLeafPaneIsLeaf tests LeafPane.IsLeaf
func TestLeafPaneIsLeaf(t *testing.T) {
	component := NewMockComponent(PaneEditor, "Editor")
	pane := NewLeafPane(component)

	if !pane.IsLeaf() {
		t.Error("LeafPane should be a leaf")
	}
}

// TestSplitPaneIsLeaf tests SplitPane.IsLeaf
func TestSplitPaneIsLeaf(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	if split.IsLeaf() {
		t.Error("SplitPane should not be a leaf")
	}
}

// TestSplitPaneRenderWithNilChildren tests SplitPane.Render with nil children
func TestSplitPaneRenderWithNilChildren(t *testing.T) {
	split := NewSplitPane(Horizontal, 0.7, nil, nil)
	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}

	// Should not panic
	view := split.Render(area, PaneEditor)
	// With both children nil, render returns "" combined with ""
	if view != "" {
		t.Errorf("Expected empty view with nil children, got '%s'", view)
	}
}

// TestSplitPaneRenderWithOneNilChild tests SplitPane.Render with one nil child
func TestSplitPaneRenderWithOneNilChild(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	split := NewSplitPane(Horizontal, 0.7, NewLeafPane(editor), nil)
	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}

	view := split.Render(area, PaneEditor)
	if view == "" {
		t.Error("Should render the non-nil child")
	}
}

// TestSplitPaneUpdateWithNilChildren tests SplitPane.Update with nil children
func TestSplitPaneUpdateWithNilChildren(t *testing.T) {
	split := NewSplitPane(Horizontal, 0.7, nil, nil)
	updatedPane, cmd := split.Update(tea.KeyMsg{}, PaneEditor)

	if updatedPane == nil {
		t.Error("Should return a pane even with nil children")
	}
	if cmd != nil {
		t.Error("Should batch commands properly")
	}
}

// TestSplitPaneUpdateWithOneNilChild tests SplitPane.Update with one nil child
func TestSplitPaneUpdateWithOneNilChild(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	split := NewSplitPane(Horizontal, 0.7, NewLeafPane(editor), nil)

	updatedPane, _ := split.Update(tea.KeyMsg{}, PaneEditor)
	if updatedPane == nil {
		t.Error("Should return a pane")
	}
}

// TestSplitPaneFindComponentWithNilChildren tests SplitPane.FindComponent with nil children
func TestSplitPaneFindComponentWithNilChildren(t *testing.T) {
	split := NewSplitPane(Horizontal, 0.7, nil, nil)
	found := split.FindComponent(PaneEditor)

	if found != nil {
		t.Error("Should not find component in nil children")
	}
}

// TestSplitPaneFindComponentInSecondChild tests SplitPane.FindComponent in second child
func TestSplitPaneFindComponentInSecondChild(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	// Should find sidebar in second child
	found := split.FindComponent(PaneSidebar)
	if found == nil {
		t.Fatal("Should find sidebar in second child")
	}
	if found.ID() != PaneSidebar {
		t.Errorf("Expected sidebar ID, got %s", found.ID())
	}
}

// TestSplitPaneAllPaneIDsWithNilChildren tests SplitPane.AllPaneIDs with nil children
func TestSplitPaneAllPaneIDsWithNilChildren(t *testing.T) {
	split := NewSplitPane(Horizontal, 0.7, nil, nil)
	ids := split.AllPaneIDs()

	if len(ids) != 0 {
		t.Errorf("Expected 0 pane IDs, got %d", len(ids))
	}
}

// TestSplitPaneAllPaneIDsWithOneNilChild tests SplitPane.AllPaneIDs with one nil child
func TestSplitPaneAllPaneIDsWithOneNilChild(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	split := NewSplitPane(Horizontal, 0.7, NewLeafPane(editor), nil)
	ids := split.AllPaneIDs()

	if len(ids) != 1 {
		t.Fatalf("Expected 1 pane ID, got %d", len(ids))
	}
	if ids[0] != PaneEditor {
		t.Errorf("Expected editor ID, got %s", ids[0])
	}
}

// TestComputeSplitAreasVertical tests computeSplitAreas with vertical direction
func TestComputeSplitAreasVerticalBoundary(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Vertical,
		0.5,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	area := Rect{X: 10, Y: 10, Width: 100, Height: 30}
	firstArea, secondArea := split.computeSplitAreas(area)

	// Check boundaries are correct
	if firstArea.X != secondArea.X {
		t.Errorf("X coordinates should match: %d vs %d", firstArea.X, secondArea.X)
	}
	if firstArea.Width != secondArea.Width {
		t.Errorf("Widths should match: %d vs %d", firstArea.Width, secondArea.Width)
	}
}

// TestComputeSplitAreasUnknownDirection tests computeSplitAreas with unknown direction
func TestComputeSplitAreasUnknownDirection(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Direction(99), // Invalid direction
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}
	_, secondArea := split.computeSplitAreas(area)

	// Should return original area and empty area for unknown direction
	if secondArea.Width != 0 && secondArea.Height != 0 {
		t.Errorf("Expected second area to be empty for unknown direction")
	}
}

// TestCombineVertical tests combineVertical function
func TestCombineVertical(t *testing.T) {
	tests := []struct {
		name     string
		top      string
		bottom   string
		expected string
	}{
		{"both non-empty", "top", "bottom", "top\nbottom"},
		{"empty top", "", "bottom", "bottom"},
		{"empty bottom", "top", "", "top"},
		{"both empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := combineVertical(tt.top, tt.bottom)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestCombineHorizontal tests combineHorizontal function
func TestCombineHorizontal(t *testing.T) {
	tests := []struct {
		name     string
		left     string
		right    string
		expected string
	}{
		{"both non-empty", "left", "right", "leftright"},
		{"empty left", "", "right", "right"},
		{"empty right", "left", "", "left"},
		{"both empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := combineHorizontal(tt.left, tt.right)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestEngineAllPaneIDsWithNoRoot tests AllPaneIDs with nil root
func TestEngineAllPaneIDsWithNoRoot(t *testing.T) {
	engine := NewEngine(LayoutFocus)
	ids := engine.AllPaneIDs()

	if len(ids) != 0 {
		t.Errorf("Expected 0 pane IDs, got %d", len(ids))
	}
}

// TestDirectionStringUnknown tests Direction.String with unknown value
func TestDirectionStringUnknown(t *testing.T) {
	d := Direction(99)
	result := d.String()
	if result != "Unknown" {
		t.Errorf("Expected 'Unknown', got '%s'", result)
	}
}

// TestLayoutModeStringUnknown tests LayoutMode.String with unknown value
func TestLayoutModeStringUnknown(t *testing.T) {
	m := LayoutMode(99)
	result := m.String()
	if result != "Unknown" {
		t.Errorf("Expected 'Unknown', got '%s'", result)
	}
}

// TestEngineFocusNextMultiplePanes tests FocusNext with more than 2 panes
func TestEngineFocusNextMultiplePanes(t *testing.T) {
	engine := NewEngine(LayoutAnalysis)

	editor := NewMockComponent(PaneEditor, "Editor")
	analysis := NewMockComponent(PaneAnalysis, "Analysis")
	diagnostics := NewMockComponent(PaneDiagnostics, "Diagnostics")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(analysis)
	engine.RegisterComponent(diagnostics)

	// Build analysis layout (top: split editor/analysis, bottom: diagnostics)
	topPane := NewSplitPane(
		Horizontal,
		0.5,
		NewLeafPane(editor),
		NewLeafPane(analysis),
	)
	engine.SetRoot(NewSplitPane(
		Vertical,
		0.8,
		topPane,
		NewLeafPane(diagnostics),
	))

	// Should start at first pane
	initial := engine.FocusedID()

	// Move through all panes - we have 3 panes total
	engine.FocusNext()
	second := engine.FocusedID()
	if initial == second {
		t.Error("FocusNext should move to next pane")
	}

	engine.FocusNext()
	third := engine.FocusedID()
	if second == third {
		t.Error("FocusNext should move to next pane")
	}

	// Should wrap back to first after visiting all 3
	engine.FocusNext()
	wrapped := engine.FocusedID()
	if wrapped != initial {
		t.Errorf("FocusNext should wrap around to first pane, got %s, expected %s", wrapped, initial)
	}
}

// TestEngineFocusPrevMultiplePanes tests FocusPrev with more than 2 panes
func TestEngineFocusPrevMultiplePanes(t *testing.T) {
	engine := NewEngine(LayoutAnalysis)

	editor := NewMockComponent(PaneEditor, "Editor")
	analysis := NewMockComponent(PaneAnalysis, "Analysis")
	diagnostics := NewMockComponent(PaneDiagnostics, "Diagnostics")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(analysis)
	engine.RegisterComponent(diagnostics)

	// Build analysis layout
	topPane := NewSplitPane(
		Horizontal,
		0.5,
		NewLeafPane(editor),
		NewLeafPane(analysis),
	)
	engine.SetRoot(NewSplitPane(
		Vertical,
		0.8,
		topPane,
		NewLeafPane(diagnostics),
	))

	// Get to last pane
	ids := engine.AllPaneIDs()
	engine.SetFocus(ids[len(ids)-1])

	lastFocus := engine.FocusedID()

	// Move backwards - should go to previous pane
	engine.FocusPrev()
	secondLast := engine.FocusedID()
	if lastFocus == secondLast {
		t.Error("FocusPrev should move to previous pane")
	}

	// Continue going back
	for i := 0; i < len(ids)-2; i++ {
		engine.FocusPrev()
	}

	// Should wrap to last
	engine.FocusPrev()
	wrapped := engine.FocusedID()
	if wrapped != lastFocus {
		t.Error("FocusPrev should wrap to last pane")
	}
}

// TestSplitPaneRenderVerticalWithChildren tests SplitPane.Render with vertical direction
func TestSplitPaneRenderVerticalWithChildren(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	diagnostics := NewMockComponent(PaneDiagnostics, "Diagnostics")

	split := NewSplitPane(
		Vertical,
		0.8,
		NewLeafPane(editor),
		NewLeafPane(diagnostics),
	)

	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}
	view := split.Render(area, PaneEditor)

	// Should successfully render with vertical split
	if view == "" {
		t.Error("Should render something for vertical split")
	}
}

// TestSplitPaneUpdateWithBothChildren tests SplitPane.Update with both children having cmds
func TestSplitPaneUpdateWithBothChildren(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	split := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)

	// Update with a message - both panes should process
	updatedPane, cmd := split.Update(tea.KeyMsg{}, PaneEditor)

	if updatedPane == nil {
		t.Error("Should return a pane")
	}
	// cmd will be nil because MockComponent returns nil
	if cmd != nil {
		t.Error("Should batch nil commands properly")
	}
}

// TestEngineBuildStandardLayoutFullLayout tests buildStandardLayout with both components
func TestEngineBuildStandardLayoutFullLayout(t *testing.T) {
	engine := NewEngine(LayoutStandard)
	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(sidebar)

	engine.buildStandardLayout()

	if engine.Root() == nil {
		t.Fatal("Root should not be nil")
	}

	// Should be a split pane with both children
	if engine.Root().IsLeaf() {
		t.Error("Root should be a split pane with both editor and sidebar")
	}

	// Verify both components are present
	if engine.Root().FindComponent(PaneEditor) == nil {
		t.Error("Should find editor")
	}
	if engine.Root().FindComponent(PaneSidebar) == nil {
		t.Error("Should find sidebar")
	}
}

// TestEngineFocusNextWithCurrentIndexNotFound tests FocusNext when current focus not in list
func TestEngineFocusNextWhenCurrentNotFound(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(sidebar)

	root := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	engine.SetRoot(root)

	// Set focus to something not in the current layout
	engine.focusedID = PaneMemory

	// FocusNext should handle gracefully (currentIndex will be -1)
	engine.FocusNext()

	// Should set to first pane when current is not found
	paneIDs := engine.AllPaneIDs()
	if len(paneIDs) > 0 && engine.FocusedID() != paneIDs[0] {
		t.Error("Should reset to first pane when current focus not found")
	}
}

// TestEngineFocusPrevWhenCurrentNotFound tests FocusPrev when current focus not in list
func TestEngineFocusPrevWhenCurrentNotFound(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	editor := NewMockComponent(PaneEditor, "Editor")
	sidebar := NewMockComponent(PaneSidebar, "Sidebar")

	engine.RegisterComponent(editor)
	engine.RegisterComponent(sidebar)

	root := NewSplitPane(
		Horizontal,
		0.7,
		NewLeafPane(editor),
		NewLeafPane(sidebar),
	)
	engine.SetRoot(root)

	// Set focus to something not in the current layout
	engine.focusedID = PaneMemory

	// FocusPrev should handle gracefully (currentIndex will be -1)
	engine.FocusPrev()

	// Should set to last pane when current is not found (prevIndex = -1 - 1 = -2, becomes len-1)
	paneIDs := engine.AllPaneIDs()
	if len(paneIDs) > 0 && engine.FocusedID() != paneIDs[len(paneIDs)-1] {
		t.Error("Should reset to last pane when current focus not found")
	}
}

// TestSplitPaneRenderWithSecondNilChild tests SplitPane.Render when only second child is nil
func TestSplitPaneRenderWithSecondNilChild(t *testing.T) {
	editor := NewMockComponent(PaneEditor, "Editor")
	split := NewSplitPane(Horizontal, 0.7, NewLeafPane(editor), nil)
	area := Rect{X: 0, Y: 0, Width: 100, Height: 30}

	view := split.Render(area, PaneEditor)

	// Should render the first child
	if view == "" {
		t.Error("Should render something when only first child is present")
	}
}

// TestEngineCompactModeSmallTerminal tests that small terminal triggers compact mode
func TestEngineCompactModeSmallTerminal(t *testing.T) {
	engine := NewEngine(LayoutStandard)

	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)
	engine.SetRoot(NewLeafPane(editor))

	// Start with standard mode
	if engine.Mode() != LayoutStandard {
		t.Errorf("Expected Standard mode, got %s", engine.Mode())
	}

	// Resize to small - should trigger compact mode
	engine.SetTerminalSize(100, 25)

	if engine.Mode() != LayoutCompact {
		t.Errorf("Expected Compact mode for small terminal, got %s", engine.Mode())
	}
}

// TestEngineCompactModeStayInCompact tests that compact mode persists
func TestEngineCompactModeStayInCompact(t *testing.T) {
	engine := NewEngine(LayoutCompact)

	editor := NewMockComponent(PaneEditor, "Editor")
	engine.RegisterComponent(editor)
	engine.SetRoot(NewLeafPane(editor))

	if engine.Mode() != LayoutCompact {
		t.Errorf("Expected Compact mode, got %s", engine.Mode())
	}

	// Resize small - should stay compact (already in compact mode)
	engine.SetTerminalSize(80, 20)

	if engine.Mode() != LayoutCompact {
		t.Errorf("Expected to stay in Compact mode, got %s", engine.Mode())
	}
}

// TestEngineCustomModePreserved tests that custom mode is preserved
func TestEngineCustomModePreserved(t *testing.T) {
	engine := NewEngine(LayoutCustom)

	if engine.Mode() != LayoutCustom {
		t.Errorf("Expected Custom mode, got %s", engine.Mode())
	}

	// Set mode should work but custom is for custom layouts
	engine.SetMode(LayoutCustom)

	if engine.Mode() != LayoutCustom {
		t.Errorf("Expected Custom mode after SetMode, got %s", engine.Mode())
	}
}
