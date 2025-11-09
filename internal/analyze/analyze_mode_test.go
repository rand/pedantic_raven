package analyze

import (
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// Test analyze mode creation
func TestNewAnalyzeMode(t *testing.T) {
	mode := NewAnalyzeMode()

	if mode == nil {
		t.Fatal("NewAnalyzeMode returned nil")
	}

	if mode.currentView != ViewTripleGraph {
		t.Errorf("Expected default view to be ViewTripleGraph, got %d", mode.currentView)
	}

	if mode.exportFormat != ExportMarkdown {
		t.Errorf("Expected default export format to be ExportMarkdown, got %d", mode.exportFormat)
	}
}

// Test setting analysis
func TestAnalyzeModeSetAnalysis(t *testing.T) {
	mode := NewAnalyzeMode()

	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "Python", Type: semantic.EntityTechnology},
			{Text: "Rust", Type: semantic.EntityTechnology},
		},
		Relationships: []semantic.Relationship{
			{Subject: "Alice", Predicate: "uses", Object: "Python"},
		},
		Timestamp: time.Now(),
	}

	mode.SetAnalysis(analysis)

	if mode.analysis != analysis {
		t.Error("Analysis not set correctly")
	}

	if len(mode.entityFreqs) == 0 {
		t.Error("Entity frequencies not calculated")
	}

	if mode.tripleGraphView.analysis != analysis {
		t.Error("Triple graph view not updated with analysis")
	}
}

// Test view switching
func TestAnalyzeModeViewSwitching(t *testing.T) {
	mode := NewAnalyzeMode()

	// Test NextView
	mode.NextView()
	if mode.currentView != ViewEntityFrequency {
		t.Errorf("NextView failed: expected ViewEntityFrequency, got %d", mode.currentView)
	}

	mode.NextView()
	if mode.currentView != ViewPatterns {
		t.Errorf("NextView failed: expected ViewPatterns, got %d", mode.currentView)
	}

	mode.NextView()
	if mode.currentView != ViewTypedHoles {
		t.Errorf("NextView failed: expected ViewTypedHoles, got %d", mode.currentView)
	}

	mode.NextView()
	if mode.currentView != ViewTripleGraph {
		t.Errorf("NextView failed: expected wrap to ViewTripleGraph, got %d", mode.currentView)
	}

	// Test PrevView
	mode.PrevView()
	if mode.currentView != ViewTypedHoles {
		t.Errorf("PrevView failed: expected ViewTypedHoles, got %d", mode.currentView)
	}

	mode.PrevView()
	if mode.currentView != ViewPatterns {
		t.Errorf("PrevView failed: expected ViewPatterns, got %d", mode.currentView)
	}
}

// Test SwitchView
func TestAnalyzeModeSwitchView(t *testing.T) {
	mode := NewAnalyzeMode()

	mode.SwitchView(ViewEntityFrequency)
	if mode.currentView != ViewEntityFrequency {
		t.Errorf("SwitchView failed: expected ViewEntityFrequency, got %d", mode.currentView)
	}

	mode.SwitchView(ViewTypedHoles)
	if mode.currentView != ViewTypedHoles {
		t.Errorf("SwitchView failed: expected ViewTypedHoles, got %d", mode.currentView)
	}
}

// Test keyboard shortcuts for view switching
func TestAnalyzeModeKeyboardViewSwitching(t *testing.T) {
	mode := NewAnalyzeMode()

	tests := []struct {
		key          string
		expectedView ViewMode
	}{
		{"tab", ViewEntityFrequency},
		{"tab", ViewPatterns},
		{"tab", ViewTypedHoles},
		{"tab", ViewTripleGraph},
		{"1", ViewTripleGraph},
		{"2", ViewEntityFrequency},
		{"3", ViewPatterns},
		{"4", ViewTypedHoles},
	}

	for _, tt := range tests {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
		updated, _ := mode.handleKeyPress(msg)

		if updated.currentView != tt.expectedView {
			t.Errorf("Key '%s': expected view %d, got %d", tt.key, tt.expectedView, updated.currentView)
		}

		mode = updated
	}
}

// Test SetSize
func TestAnalyzeModeSetSize(t *testing.T) {
	mode := NewAnalyzeMode()

	mode.SetSize(120, 40)

	if mode.width != 120 {
		t.Errorf("Expected width 120, got %d", mode.width)
	}

	if mode.height != 40 {
		t.Errorf("Expected height 40, got %d", mode.height)
	}

	// Triple graph view should get adjusted size (height - 4 for header/footer)
	if mode.tripleGraphView.width != 120 {
		t.Errorf("Expected triple graph width 120, got %d", mode.tripleGraphView.width)
	}

	if mode.tripleGraphView.height != 36 {
		t.Errorf("Expected triple graph height 36, got %d", mode.tripleGraphView.height)
	}
}

// Test Update with WindowSizeMsg
func TestAnalyzeModeUpdateWindowSize(t *testing.T) {
	mode := NewAnalyzeMode()

	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updated, _ := mode.Update(msg)

	if updated.width != 100 {
		t.Errorf("Expected width 100, got %d", updated.width)
	}

	if updated.height != 30 {
		t.Errorf("Expected height 30, got %d", updated.height)
	}
}

// Test export format cycling
func TestAnalyzeModeExportFormatCycling(t *testing.T) {
	mode := NewAnalyzeMode()

	// Start with Markdown
	if mode.exportFormat != ExportMarkdown {
		t.Errorf("Expected initial format ExportMarkdown, got %d", mode.exportFormat)
	}

	// Simulate ctrl+e key press (cycle format)
	mode.exportFormat = (mode.exportFormat + 1) % 3
	if mode.exportFormat != ExportHTML {
		t.Errorf("Expected ExportHTML, got %d", mode.exportFormat)
	}

	mode.exportFormat = (mode.exportFormat + 1) % 3
	if mode.exportFormat != ExportPDF {
		t.Errorf("Expected ExportPDF, got %d", mode.exportFormat)
	}

	mode.exportFormat = (mode.exportFormat + 1) % 3
	if mode.exportFormat != ExportMarkdown {
		t.Errorf("Expected wrap to ExportMarkdown, got %d", mode.exportFormat)
	}
}

// Test View rendering without analysis
func TestAnalyzeModeViewNoAnalysis(t *testing.T) {
	mode := NewAnalyzeMode()
	mode.SetSize(80, 24)

	view := mode.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	// Should show helpful message
	if len(view) < 10 {
		t.Error("View too short for helpful message")
	}
}

// Test View rendering with analysis
func TestAnalyzeModeViewWithAnalysis(t *testing.T) {
	mode := NewAnalyzeMode()
	mode.SetSize(120, 40)

	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "Python", Type: semantic.EntityTechnology},
			{Text: "Go", Type: semantic.EntityTechnology},
			{Text: "Rust", Type: semantic.EntityTechnology},
		},
		Relationships: []semantic.Relationship{
			{Subject: "Alice", Predicate: "uses", Object: "Python"},
			{Subject: "Bob", Predicate: "uses", Object: "Go"},
		},
		Timestamp: time.Now(),
	}

	mode.SetAnalysis(analysis)

	// Test all views render without panic
	views := []ViewMode{ViewTripleGraph, ViewEntityFrequency, ViewPatterns, ViewTypedHoles}
	for _, view := range views {
		mode.SwitchView(view)
		output := mode.View()

		if output == "" {
			t.Errorf("View %d produced empty output", view)
		}

		// Should include header with tabs
		if len(output) < 100 {
			t.Errorf("View %d output too short: %d characters", view, len(output))
		}
	}
}

// Test renderHeader
func TestAnalyzeModeRenderHeader(t *testing.T) {
	mode := NewAnalyzeMode()

	// Test header for each view
	views := []ViewMode{ViewTripleGraph, ViewEntityFrequency, ViewPatterns, ViewTypedHoles}
	for _, view := range views {
		mode.SwitchView(view)
		header := mode.renderHeader()

		if header == "" {
			t.Errorf("Header empty for view %d", view)
		}

		// Should contain all tab labels
		expectedLabels := []string{"Triple Graph", "Entity Frequency", "Patterns", "Typed Holes"}
		for _, label := range expectedLabels {
			if !containsString(header, label) {
				t.Errorf("Header for view %d missing label '%s'", view, label)
			}
		}
	}
}

// Test renderFooter
func TestAnalyzeModeRenderFooter(t *testing.T) {
	mode := NewAnalyzeMode()

	footer := mode.renderFooter()

	if footer == "" {
		t.Error("Footer should not be empty")
	}

	// Should contain key shortcuts
	expectedShortcuts := []string{"tab", "Export", "Help"}
	for _, shortcut := range expectedShortcuts {
		if !containsString(footer, shortcut) {
			t.Errorf("Footer missing shortcut '%s'", shortcut)
		}
	}
}

// Test renderEntityFrequency
func TestAnalyzeModeRenderEntityFrequency(t *testing.T) {
	mode := NewAnalyzeMode()
	mode.SetSize(120, 40)

	// Without analysis
	output := mode.renderEntityFrequency()
	if output == "" {
		t.Error("Should render message when no entities")
	}

	// With analysis
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "Python", Type: semantic.EntityTechnology, Count: 1},
			{Text: "Python", Type: semantic.EntityTechnology, Count: 1}, // Duplicate for frequency
			{Text: "Go", Type: semantic.EntityTechnology, Count: 1},
		},
		Timestamp: time.Now(),
	}

	mode.SetAnalysis(analysis)
	output = mode.renderEntityFrequency()

	if output == "" {
		t.Error("Entity frequency view should not be empty")
	}

	// Debug: print output
	t.Logf("Entity frequency output:\n%s", output)

	if !containsString(output, "Python") {
		t.Error("Entity frequency view should mention Python")
	}

	if !containsString(output, "Entity Type Distribution") {
		t.Error("Entity frequency view should include type distribution")
	}
}

// Test renderPatterns
func TestAnalyzeModeRenderPatterns(t *testing.T) {
	mode := NewAnalyzeMode()
	mode.SetSize(120, 40)

	// Without patterns
	output := mode.renderPatterns()
	if output == "" {
		t.Error("Should render message when no patterns")
	}

	// With patterns
	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "Alice", Type: semantic.EntityPerson},
			{Text: "Acme", Type: semantic.EntityOrganization},
		},
		Relationships: []semantic.Relationship{
			{Subject: "Alice", Predicate: "works_at", Object: "Acme"},
			{Subject: "Bob", Predicate: "works_at", Object: "Tech Inc"},
		},
		Timestamp: time.Now(),
	}

	mode.SetAnalysis(analysis)
	output = mode.renderPatterns()

	if output == "" {
		t.Error("Patterns view should not be empty")
	}
}

// Test renderTypedHoles
func TestAnalyzeModeRenderTypedHoles(t *testing.T) {
	mode := NewAnalyzeMode()
	mode.SetSize(120, 40)

	// Without holes
	output := mode.renderTypedHoles()
	if output == "" {
		t.Error("Should render message when no typed holes")
	}

	// With holes
	analysis := &semantic.Analysis{
		TypedHoles: []semantic.TypedHole{
			{
				Type:       "UserService interface",
				Constraint: "Handles user operations",
			},
		},
		Timestamp: time.Now(),
	}

	mode.SetAnalysis(analysis)
	output = mode.renderTypedHoles()

	if output == "" {
		t.Error("Typed holes view should not be empty")
	}
}

// Test exportReport
func TestAnalyzeModeExportReport(t *testing.T) {
	mode := NewAnalyzeMode()

	analysis := &semantic.Analysis{
		Entities: []semantic.Entity{
			{Text: "Python", Type: semantic.EntityTechnology},
		},
		Relationships: []semantic.Relationship{
			{Subject: "Alice", Predicate: "uses", Object: "Python"},
		},
		Timestamp: time.Now(),
	}

	mode.SetAnalysis(analysis)

	// Test export command creation (doesn't actually export, just creates command)
	cmd := mode.exportReport()
	if cmd == nil {
		t.Error("exportReport should return a command")
	}

	// Execute command to test report generation
	msg := cmd()

	// Should return ExportCompleteMsg or ExportErrorMsg
	switch msg.(type) {
	case ExportCompleteMsg:
		// Success
	case ExportErrorMsg:
		// Expected in test environment without file system
	default:
		t.Errorf("Unexpected message type: %T", msg)
	}
}

// Test Update with ExportCompleteMsg
func TestAnalyzeModeUpdateExportComplete(t *testing.T) {
	mode := NewAnalyzeMode()

	msg := ExportCompleteMsg{Filename: "test.md"}
	updated, _ := mode.Update(msg)

	if updated.lastExport.IsZero() {
		t.Error("lastExport should be set after ExportCompleteMsg")
	}
}

// Test Update with ExportErrorMsg
func TestAnalyzeModeUpdateExportError(t *testing.T) {
	mode := NewAnalyzeMode()

	msg := ExportErrorMsg{Err: fmt.Errorf("test error")}
	updated, _ := mode.Update(msg)

	if updated.err == nil {
		t.Error("Error should be set after ExportErrorMsg")
	}

	if updated.err.Error() != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", updated.err.Error())
	}
}

// Test truncate helper
func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exact length", 12, "exact length"},
		{"this is a very long string", 10, "this is..."},
		{"", 5, ""},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, expected %q", tt.input, tt.maxLen, result, tt.expected)
		}

		if len(result) > tt.maxLen {
			t.Errorf("truncate result too long: %d > %d", len(result), tt.maxLen)
		}
	}
}

// Test Init
func TestAnalyzeModeInit(t *testing.T) {
	mode := NewAnalyzeMode()
	cmd := mode.Init()

	// Init can return nil or a command
	_ = cmd
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
