package palette

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
	"github.com/rand/pedantic-raven/internal/overlay"
)

// --- CommandRegistry Tests ---

func TestCommandRegistryRegister(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{
		ID:          "test.command",
		Name:        "Test Command",
		Description: "A test command",
		Category:    CategoryFile,
	}

	registry.Register(cmd)

	retrieved, ok := registry.Get("test.command")
	if !ok {
		t.Fatal("Command should be registered")
	}

	if retrieved.ID != cmd.ID {
		t.Errorf("Expected ID %s, got %s", cmd.ID, retrieved.ID)
	}

	if registry.Count() != 1 {
		t.Errorf("Expected count 1, got %d", registry.Count())
	}
}

func TestCommandRegistryUnregister(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{ID: "test.command", Name: "Test"}
	registry.Register(cmd)

	// Verify registered
	if _, ok := registry.Get("test.command"); !ok {
		t.Fatal("Command should be registered")
	}

	// Unregister
	registry.Unregister("test.command")

	// Verify gone
	if _, ok := registry.Get("test.command"); ok {
		t.Error("Command should be unregistered")
	}
}

func TestCommandRegistryAll(t *testing.T) {
	registry := NewCommandRegistry()

	cmd1 := Command{ID: "cmd1", Name: "Command 1"}
	cmd2 := Command{ID: "cmd2", Name: "Command 2"}

	registry.Register(cmd1)
	registry.Register(cmd2)

	all := registry.All()
	if len(all) != 2 {
		t.Fatalf("Expected 2 commands, got %d", len(all))
	}

	// Check both are present (order doesn't matter)
	hasCmd1 := false
	hasCmd2 := false

	for _, cmd := range all {
		if cmd.ID == "cmd1" {
			hasCmd1 = true
		}
		if cmd.ID == "cmd2" {
			hasCmd2 = true
		}
	}

	if !hasCmd1 || !hasCmd2 {
		t.Error("Not all commands present")
	}
}

func TestCommandRegistryByCategory(t *testing.T) {
	registry := NewCommandRegistry()

	fileCmd := Command{ID: "file.open", Name: "Open", Category: CategoryFile}
	editCmd := Command{ID: "edit.copy", Name: "Copy", Category: CategoryEdit}
	viewCmd := Command{ID: "view.zoom", Name: "Zoom", Category: CategoryView}

	registry.Register(fileCmd)
	registry.Register(editCmd)
	registry.Register(viewCmd)

	// Get file commands
	fileCommands := registry.ByCategory(CategoryFile)
	if len(fileCommands) != 1 {
		t.Fatalf("Expected 1 file command, got %d", len(fileCommands))
	}

	if fileCommands[0].ID != "file.open" {
		t.Errorf("Expected file.open, got %s", fileCommands[0].ID)
	}

	// Get edit commands
	editCommands := registry.ByCategory(CategoryEdit)
	if len(editCommands) != 1 {
		t.Fatalf("Expected 1 edit command, got %d", len(editCommands))
	}
}

// --- Fuzzy Matching Tests ---

func TestFuzzyMatchEmpty(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{ID: "test", Name: "Test Command"}
	registry.Register(cmd)

	// Empty query should return all commands
	matches := registry.FuzzyMatch("")
	if len(matches) != 1 {
		t.Errorf("Expected 1 match for empty query, got %d", len(matches))
	}

	if matches[0].Score != 0 {
		t.Errorf("Expected score 0 for empty query, got %d", matches[0].Score)
	}
}

func TestFuzzyMatchExactName(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{ID: "open", Name: "Open File", Description: "Opens a file"}
	registry.Register(cmd)

	matches := registry.FuzzyMatch("Open File")
	if len(matches) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matches))
	}

	// Exact match should have highest score
	if matches[0].Score < 100 {
		t.Errorf("Expected score >= 100 for exact match, got %d", matches[0].Score)
	}
}

func TestFuzzyMatchNameContains(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{ID: "save", Name: "Save File", Description: "Saves the current file"}
	registry.Register(cmd)

	matches := registry.FuzzyMatch("save")
	if len(matches) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matches))
	}

	// Name contains should score >= 50
	if matches[0].Score < 50 {
		t.Errorf("Expected score >= 50 for name contains, got %d", matches[0].Score)
	}
}

func TestFuzzyMatchDescription(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{
		ID:          "analyze",
		Name:        "Run Analysis",
		Description: "Performs semantic analysis",
	}
	registry.Register(cmd)

	matches := registry.FuzzyMatch("semantic")
	if len(matches) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(matches))
	}

	// Description match should score >= 20
	if matches[0].Score < 20 {
		t.Errorf("Expected score >= 20 for description match, got %d", matches[0].Score)
	}
}

func TestFuzzyMatchSubsequence(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{ID: "fop", Name: "File Open", Description: "Open a file"}
	registry.Register(cmd)

	// "fop" should match "File Open" as subsequence
	matches := registry.FuzzyMatch("fop")
	if len(matches) != 1 {
		t.Fatalf("Expected 1 match for subsequence, got %d", len(matches))
	}

	if matches[0].Score == 0 {
		t.Error("Expected non-zero score for subsequence match")
	}
}

func TestFuzzyMatchNoMatch(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{ID: "open", Name: "Open File"}
	registry.Register(cmd)

	matches := registry.FuzzyMatch("xyz")
	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
	}
}

func TestFuzzyMatchMultiple(t *testing.T) {
	registry := NewCommandRegistry()

	cmd1 := Command{ID: "save", Name: "Save", Description: "Save file"}
	cmd2 := Command{ID: "save.as", Name: "Save As", Description: "Save with new name"}
	cmd3 := Command{ID: "open", Name: "Open", Description: "Open file"}

	registry.Register(cmd1)
	registry.Register(cmd2)
	registry.Register(cmd3)

	matches := registry.FuzzyMatch("save")

	// Should match cmd1 and cmd2
	if len(matches) < 2 {
		t.Fatalf("Expected at least 2 matches, got %d", len(matches))
	}

	// Results should be sorted by score
	for i := 0; i < len(matches)-1; i++ {
		if matches[i].Score < matches[i+1].Score {
			t.Errorf("Results not sorted by score: %d < %d", matches[i].Score, matches[i+1].Score)
		}
	}
}

func TestFuzzyMatchCaseInsensitive(t *testing.T) {
	registry := NewCommandRegistry()

	cmd := Command{ID: "open", Name: "Open File"}
	registry.Register(cmd)

	matchesLower := registry.FuzzyMatch("open")
	matchesUpper := registry.FuzzyMatch("OPEN")
	matchesMixed := registry.FuzzyMatch("OpEn")

	if len(matchesLower) != 1 || len(matchesUpper) != 1 || len(matchesMixed) != 1 {
		t.Error("Case-insensitive matching failed")
	}
}

// --- Palette Tests ---

func TestPaletteCreation(t *testing.T) {
	registry := NewCommandRegistry()
	palette := NewPalette("palette", registry)

	if palette.ID() != "palette" {
		t.Errorf("Expected ID 'palette', got '%s'", palette.ID())
	}

	if !palette.Modal() {
		t.Error("Palette should be modal")
	}
}

func TestPaletteTypeQuery(t *testing.T) {
	registry := NewCommandRegistry()
	cmd := Command{ID: "test", Name: "Test Command"}
	registry.Register(cmd)

	palette := NewPalette("palette", registry)

	// Type "test"
	updated, _ := palette.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	palette = updated.(*Palette)

	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	palette = updated.(*Palette)

	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	palette = updated.(*Palette)

	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	palette = updated.(*Palette)

	if palette.query != "test" {
		t.Errorf("Expected query 'test', got '%s'", palette.query)
	}

	if len(palette.matches) == 0 {
		t.Error("Expected matches after typing query")
	}
}

func TestPaletteBackspace(t *testing.T) {
	registry := NewCommandRegistry()
	palette := NewPalette("palette", registry)

	// Type "test"
	updated, _ := palette.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}})
	palette = updated.(*Palette)

	if palette.query != "test" {
		t.Fatalf("Expected query 'test', got '%s'", palette.query)
	}

	// Backspace
	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	palette = updated.(*Palette)

	if palette.query != "tes" {
		t.Errorf("Expected query 'tes' after backspace, got '%s'", palette.query)
	}
}

func TestPaletteClearQuery(t *testing.T) {
	registry := NewCommandRegistry()
	palette := NewPalette("palette", registry)

	// Type query
	updated, _ := palette.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}})
	palette = updated.(*Palette)

	if palette.query != "test" {
		t.Fatalf("Expected query 'test', got '%s'", palette.query)
	}

	// Clear with Ctrl+U
	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	palette = updated.(*Palette)

	if palette.query != "" {
		t.Errorf("Expected empty query after Ctrl+U, got '%s'", palette.query)
	}
}

func TestPaletteNavigation(t *testing.T) {
	registry := NewCommandRegistry()

	cmd1 := Command{ID: "cmd1", Name: "Command 1"}
	cmd2 := Command{ID: "cmd2", Name: "Command 2"}
	cmd3 := Command{ID: "cmd3", Name: "Command 3"}

	registry.Register(cmd1)
	registry.Register(cmd2)
	registry.Register(cmd3)

	palette := NewPalette("palette", registry)

	if palette.selected != 0 {
		t.Fatal("Expected initial selection 0")
	}

	// Down arrow
	updated, _ := palette.Update(tea.KeyMsg{Type: tea.KeyDown})
	palette = updated.(*Palette)

	if palette.selected != 1 {
		t.Errorf("Expected selection 1 after down, got %d", palette.selected)
	}

	// Down again
	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyDown})
	palette = updated.(*Palette)

	if palette.selected != 2 {
		t.Errorf("Expected selection 2 after down, got %d", palette.selected)
	}

	// Up arrow
	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyUp})
	palette = updated.(*Palette)

	if palette.selected != 1 {
		t.Errorf("Expected selection 1 after up, got %d", palette.selected)
	}

	// Ctrl+P (up)
	updated, _ = palette.Update(tea.KeyMsg{Type: tea.KeyCtrlP})
	palette = updated.(*Palette)

	if palette.selected != 0 {
		t.Errorf("Expected selection 0 after Ctrl+P, got %d", palette.selected)
	}
}

func TestPaletteExecute(t *testing.T) {
	registry := NewCommandRegistry()

	executed := false
	cmd := Command{
		ID:   "test",
		Name: "Test Command",
		Execute: func() tea.Cmd {
			executed = true
			return nil
		},
	}
	registry.Register(cmd)

	palette := NewPalette("palette", registry)

	// Enter to execute
	_, teaCmd := palette.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !executed {
		t.Error("Expected command to be executed")
	}

	// Should return dismiss command
	if teaCmd == nil {
		t.Error("Expected dismiss command after execution")
	}
}

func TestPaletteDismiss(t *testing.T) {
	registry := NewCommandRegistry()
	palette := NewPalette("palette", registry)

	// Esc to dismiss
	_, teaCmd := palette.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if teaCmd == nil {
		t.Fatal("Expected dismiss command")
	}

	// Execute command to get message
	msg := teaCmd()
	dismissMsg, ok := msg.(overlay.DismissOverlay)
	if !ok {
		t.Fatalf("Expected DismissOverlay message, got %T", msg)
	}

	if dismissMsg.ID != "palette" {
		t.Errorf("Expected dismiss ID 'palette', got '%s'", dismissMsg.ID)
	}
}

func TestPaletteView(t *testing.T) {
	registry := NewCommandRegistry()
	cmd := Command{
		ID:          "test",
		Name:        "Test Command",
		Description: "A test command",
		Keybinding:  "Ctrl+T",
	}
	registry.Register(cmd)

	palette := NewPalette("palette", registry)

	area := layout.Rect{X: 0, Y: 0, Width: 80, Height: 20}
	view := palette.View(area)

	// View should contain command name
	if !contains(view, "Test Command") {
		t.Error("View should contain command name")
	}

	// View should contain keybinding
	if !contains(view, "Ctrl+T") {
		t.Error("View should contain keybinding")
	}

	// View should contain description
	if !contains(view, "A test command") {
		t.Error("View should contain description")
	}
}
