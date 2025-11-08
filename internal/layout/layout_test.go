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
