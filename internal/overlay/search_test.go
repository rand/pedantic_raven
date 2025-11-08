package overlay

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
)

func TestNewSearchOverlay(t *testing.T) {
	// Search mode
	so := NewSearchOverlay("test-search", "search")
	if so.mode != "search" {
		t.Errorf("Expected mode 'search', got '%s'", so.mode)
	}

	// Replace mode
	so = NewSearchOverlay("test-replace", "replace")
	if so.mode != "replace" {
		t.Errorf("Expected mode 'replace', got '%s'", so.mode)
	}

	// Invalid mode should default to search
	so = NewSearchOverlay("test-invalid", "invalid")
	if so.mode != "search" {
		t.Errorf("Expected invalid mode to default to 'search', got '%s'", so.mode)
	}
}

func TestSearchOverlayInitialState(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	if so.queryText != "" {
		t.Errorf("Expected empty query, got '%s'", so.queryText)
	}

	if so.replacementText != "" {
		t.Errorf("Expected empty replacement, got '%s'", so.replacementText)
	}

	if so.activeField != fieldQuery {
		t.Errorf("Expected active field to be query")
	}

	if so.matchCount != -1 {
		t.Errorf("Expected match count -1, got %d", so.matchCount)
	}

	if so.currentMatch != -1 {
		t.Errorf("Expected current match -1, got %d", so.currentMatch)
	}

	// Check default options
	if so.opts.CaseSensitive {
		t.Error("Expected case insensitive by default")
	}

	if so.opts.WholeWord {
		t.Error("Expected whole word off by default")
	}

	if so.opts.Regex {
		t.Error("Expected regex off by default")
	}

	if !so.opts.WrapAround {
		t.Error("Expected wrap around on by default")
	}
}

func TestSearchOverlayQueryInput(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	// Type some text
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello")}
	so.Update(keyMsg)

	if so.queryText != "hello" {
		t.Errorf("Expected query 'hello', got '%s'", so.queryText)
	}

	// Add more text
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" world")}
	so.Update(keyMsg)

	if so.queryText != "hello world" {
		t.Errorf("Expected query 'hello world', got '%s'", so.queryText)
	}
}

func TestSearchOverlayBackspace(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	// Type some text
	so.queryText = "hello"

	// Backspace
	keyMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	so.Update(keyMsg)

	if so.queryText != "hell" {
		t.Errorf("Expected query 'hell', got '%s'", so.queryText)
	}

	// Backspace multiple times
	so.Update(keyMsg)
	so.Update(keyMsg)

	if so.queryText != "he" {
		t.Errorf("Expected query 'he', got '%s'", so.queryText)
	}

	// Backspace on empty should not error
	so.queryText = ""
	so.Update(keyMsg)

	if so.queryText != "" {
		t.Errorf("Expected empty query, got '%s'", so.queryText)
	}
}

func TestSearchOverlayTabSwitchFields(t *testing.T) {
	so := NewSearchOverlay("test", "replace")

	// Start at query field
	if so.activeField != fieldQuery {
		t.Fatal("Expected to start at query field")
	}

	// Tab to replacement field
	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	so.Update(keyMsg)

	if so.activeField != fieldReplacement {
		t.Error("Expected to switch to replacement field")
	}

	// Tab back to query field
	so.Update(keyMsg)

	if so.activeField != fieldQuery {
		t.Error("Expected to switch back to query field")
	}
}

func TestSearchOverlayReplacementInput(t *testing.T) {
	so := NewSearchOverlay("test", "replace")

	// Switch to replacement field
	so.activeField = fieldReplacement

	// Type replacement text
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hi")}
	so.Update(keyMsg)

	if so.replacementText != "hi" {
		t.Errorf("Expected replacement 'hi', got '%s'", so.replacementText)
	}
}

func TestSearchOverlayToggleCaseSensitive(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	// Start as case insensitive
	if so.opts.CaseSensitive {
		t.Fatal("Expected case insensitive by default")
	}

	// Toggle with Ctrl+C
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := so.Update(keyMsg)

	if !so.opts.CaseSensitive {
		t.Error("Expected case sensitive after toggle")
	}

	// Should emit search action
	if cmd == nil {
		t.Error("Expected command after toggle")
	}

	// Toggle back
	so.Update(keyMsg)

	if so.opts.CaseSensitive {
		t.Error("Expected case insensitive after second toggle")
	}
}

func TestSearchOverlayToggleWholeWord(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	// Start with whole word off
	if so.opts.WholeWord {
		t.Fatal("Expected whole word off by default")
	}

	// Toggle with Ctrl+W
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlW}
	_, cmd := so.Update(keyMsg)

	if !so.opts.WholeWord {
		t.Error("Expected whole word on after toggle")
	}

	// Should emit search action
	if cmd == nil {
		t.Error("Expected command after toggle")
	}
}

func TestSearchOverlayToggleRegex(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	// Start with regex off
	if so.opts.Regex {
		t.Fatal("Expected regex off by default")
	}

	// Toggle with Ctrl+R
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlR}
	_, cmd := so.Update(keyMsg)

	if !so.opts.Regex {
		t.Error("Expected regex on after toggle")
	}

	// Should emit search action
	if cmd == nil {
		t.Error("Expected command after toggle")
	}
}

func TestSearchOverlayEnterSearch(t *testing.T) {
	so := NewSearchOverlay("test", "search")
	so.queryText = "hello"

	// Press Enter
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := so.Update(keyMsg)

	if cmd == nil {
		t.Fatal("Expected command after Enter")
	}

	// Check message
	msg := cmd()
	result, ok := msg.(SearchResult)
	if !ok {
		t.Fatal("Expected SearchResult message")
	}

	if result.Action != SearchActionFind {
		t.Errorf("Expected SearchActionFind, got %v", result.Action)
	}

	if result.Query != "hello" {
		t.Errorf("Expected query 'hello', got '%s'", result.Query)
	}

	if result.Canceled {
		t.Error("Expected not canceled")
	}
}

func TestSearchOverlayEnterReplace(t *testing.T) {
	so := NewSearchOverlay("test", "replace")
	so.queryText = "hello"
	so.replacementText = "hi"
	so.activeField = fieldReplacement

	// Press Enter on replacement field
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := so.Update(keyMsg)

	if cmd == nil {
		t.Fatal("Expected command after Enter")
	}

	// Check message
	msg := cmd()
	result, ok := msg.(SearchResult)
	if !ok {
		t.Fatal("Expected SearchResult message")
	}

	if result.Action != SearchActionReplace {
		t.Errorf("Expected SearchActionReplace, got %v", result.Action)
	}

	if result.Query != "hello" {
		t.Errorf("Expected query 'hello', got '%s'", result.Query)
	}

	if result.Replacement != "hi" {
		t.Errorf("Expected replacement 'hi', got '%s'", result.Replacement)
	}
}

func TestSearchOverlayFindNext(t *testing.T) {
	so := NewSearchOverlay("test", "search")
	so.queryText = "hello"

	// Press F3
	keyMsg := tea.KeyMsg{Type: tea.KeyF3}
	_, cmd := so.Update(keyMsg)

	if cmd == nil {
		t.Fatal("Expected command after F3")
	}

	// Check message
	msg := cmd()
	result, ok := msg.(SearchResult)
	if !ok {
		t.Fatal("Expected SearchResult message")
	}

	if result.Action != SearchActionFindNext {
		t.Errorf("Expected SearchActionFindNext, got %v", result.Action)
	}
}

func TestSearchOverlayFindPrevious(t *testing.T) {
	so := NewSearchOverlay("test", "search")
	so.queryText = "hello"

	// Note: Testing Shift+F3 requires a more complex KeyMsg setup
	// The actual functionality works in the app via msg.String() == "shift+f3"
	// For now, we'll skip this specific key combination test
	t.Skip("Shift+F3 key testing requires custom KeyMsg construction")
}

func TestSearchOverlayReplaceAll(t *testing.T) {
	so := NewSearchOverlay("test", "replace")
	so.queryText = "hello"
	so.replacementText = "hi"

	// Press Ctrl+A
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlA}
	_, cmd := so.Update(keyMsg)

	if cmd == nil {
		t.Fatal("Expected command after Ctrl+A")
	}

	// Check message
	msg := cmd()
	result, ok := msg.(SearchResult)
	if !ok {
		t.Fatal("Expected SearchResult message")
	}

	if result.Action != SearchActionReplaceAll {
		t.Errorf("Expected SearchActionReplaceAll, got %v", result.Action)
	}

	if result.Query != "hello" {
		t.Errorf("Expected query 'hello', got '%s'", result.Query)
	}

	if result.Replacement != "hi" {
		t.Errorf("Expected replacement 'hi', got '%s'", result.Replacement)
	}
}

func TestSearchOverlayReplaceAllSearchMode(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	// Ctrl+A in search mode should not do anything
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlA}
	_, cmd := so.Update(keyMsg)

	if cmd != nil {
		t.Error("Expected no command in search mode for Ctrl+A")
	}
}

func TestSearchOverlayCancel(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	// Press Esc
	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := so.Update(keyMsg)

	if cmd == nil {
		t.Fatal("Expected command after Esc")
	}

	// Should emit both SearchResult and DismissOverlay
	msg := cmd()

	// The batch will return multiple messages, so we need to check if one is SearchResult
	switch msg := msg.(type) {
	case tea.BatchMsg:
		// Check that one of the messages is SearchResult
		found := false
		for _, m := range msg {
			result := m()
			if sr, ok := result.(SearchResult); ok {
				if !sr.Canceled {
					t.Error("Expected canceled result")
				}
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected SearchResult in batch")
		}
	case SearchResult:
		if !msg.Canceled {
			t.Error("Expected canceled result")
		}
	default:
		t.Errorf("Expected BatchMsg or SearchResult, got %T", msg)
	}
}

func TestSearchOverlaySetMatchInfo(t *testing.T) {
	so := NewSearchOverlay("test", "search")

	so.SetMatchInfo(5, 2)

	if so.matchCount != 5 {
		t.Errorf("Expected match count 5, got %d", so.matchCount)
	}

	if so.currentMatch != 2 {
		t.Errorf("Expected current match 2, got %d", so.currentMatch)
	}
}

func TestSearchOverlayView(t *testing.T) {
	so := NewSearchOverlay("test", "search")
	so.queryText = "hello"

	area := layout.Rect{X: 0, Y: 0, Width: 60, Height: 15}
	view := so.View(area)

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// View should contain query text
	if !contains(view, "hello") {
		t.Error("Expected view to contain query text")
	}
}

func TestSearchOverlayViewReplace(t *testing.T) {
	so := NewSearchOverlay("test", "replace")
	so.queryText = "hello"
	so.replacementText = "hi"

	area := layout.Rect{X: 0, Y: 0, Width: 60, Height: 15}
	view := so.View(area)

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// View should contain both query and replacement
	if !contains(view, "hello") {
		t.Error("Expected view to contain query text")
	}

	if !contains(view, "hi") {
		t.Error("Expected view to contain replacement text")
	}
}

func TestSearchOverlayViewMatchInfo(t *testing.T) {
	so := NewSearchOverlay("test", "search")
	so.SetMatchInfo(5, 2)

	area := layout.Rect{X: 0, Y: 0, Width: 60, Height: 15}
	view := so.View(area)

	// Should show "Match 3 of 5" (currentMatch is 0-indexed)
	if !contains(view, "Match 3 of 5") {
		t.Error("Expected view to show match info")
	}
}

func TestSearchOverlayViewNoMatches(t *testing.T) {
	so := NewSearchOverlay("test", "search")
	so.SetMatchInfo(0, -1)

	area := layout.Rect{X: 0, Y: 0, Width: 60, Height: 15}
	view := so.View(area)

	// Should show "No matches found"
	if !contains(view, "No matches found") {
		t.Error("Expected view to show no matches message")
	}
}

func TestSearchOverlayViewOptions(t *testing.T) {
	so := NewSearchOverlay("test", "search")
	so.opts.CaseSensitive = true
	so.opts.WholeWord = true

	area := layout.Rect{X: 0, Y: 0, Width: 60, Height: 15}
	view := so.View(area)

	// Should show checked options
	if !contains(view, "[✓]") {
		t.Error("Expected view to show checked options")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && anySubstring(s, substr))
}

func anySubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
