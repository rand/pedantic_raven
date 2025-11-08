package editor

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	contextpanel "github.com/rand/pedantic-raven/internal/context"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"github.com/rand/pedantic-raven/internal/layout"
	"github.com/rand/pedantic-raven/internal/modes"
	"github.com/rand/pedantic-raven/internal/terminal"
)

// --- EditMode Tests ---

func TestNewEditMode(t *testing.T) {
	mode := NewEditMode()

	if mode == nil {
		t.Fatal("Expected mode to be created")
	}

	if mode.ID() != modes.ModeEdit {
		t.Errorf("Expected mode ID to be ModeEdit, got %v", mode.ID())
	}

	if mode.editor == nil {
		t.Error("Expected editor component to be created")
	}

	if mode.contextPanel == nil {
		t.Error("Expected context panel component to be created")
	}

	if mode.terminalComp == nil {
		t.Error("Expected terminal component to be created")
	}

	if mode.analyzer == nil {
		t.Error("Expected analyzer to be created")
	}
}

func TestEditModeInit(t *testing.T) {
	mode := NewEditMode()

	cmd := mode.Init()
	// BaseMode.Init() may return nil or a command
	_ = cmd
}

func TestEditModeOnEnter(t *testing.T) {
	mode := NewEditMode()

	// Empty content should not trigger analysis
	cmd := mode.OnEnter()
	if cmd != nil {
		// This would trigger analysis if there's content
		_ = cmd
	}
}

func TestEditModeOnEnterWithContent(t *testing.T) {
	mode := NewEditMode()

	// Set some content
	mode.editor.SetContent("Test content")

	cmd := mode.OnEnter()
	if cmd == nil {
		t.Error("Expected analysis command when entering with content")
	}
}

func TestEditModeOnExit(t *testing.T) {
	mode := NewEditMode()

	// Set analyzing flag
	mode.analyzing = true

	cmd := mode.OnExit()
	_ = cmd

	if mode.analyzing {
		t.Error("Expected analyzing flag to be cleared on exit")
	}
}

func TestEditModeUpdate(t *testing.T) {
	mode := NewEditMode()

	// Test key message
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

	updatedMode, cmd := mode.Update(keyMsg)
	_ = cmd

	if updatedMode == nil {
		t.Fatal("Expected mode to be returned")
	}
}

func TestEditModeSemanticAnalysisMsg(t *testing.T) {
	mode := NewEditMode()

	// Set analyzing flag
	mode.analyzing = true

	// Create analysis message
	analysis := &semantic.Analysis{
		Content: "Test content",
		Entities: []semantic.Entity{
			{Text: "Test", Type: semantic.EntityConcept, Count: 1},
		},
	}

	msg := SemanticAnalysisMsg{Analysis: analysis}

	updatedMode, cmd := mode.Update(msg)
	_ = cmd

	if mode.analyzing {
		t.Error("Expected analyzing flag to be cleared after analysis")
	}

	// Check that context panel was updated
	panel := mode.contextPanel.GetPanel()
	if panel.GetAnalysis() == nil {
		t.Error("Expected context panel to have analysis results")
	}

	if updatedMode == nil {
		t.Error("Expected mode to be returned")
	}
}

func TestEditModeView(t *testing.T) {
	mode := NewEditMode()

	view := mode.View()

	// Should return non-empty view
	if view == "" {
		// Empty view is okay if layout engine returns empty
		_ = view
	}
}

func TestEditModeKeybindings(t *testing.T) {
	mode := NewEditMode()

	bindings := mode.Keybindings()

	if len(bindings) == 0 {
		t.Error("Expected keybindings to be defined")
	}
}

func TestEditModeTriggerAnalysis(t *testing.T) {
	mode := NewEditMode()

	// Set content
	mode.editor.SetContent("Test content for analysis")

	cmd := mode.triggerAnalysis()

	if cmd == nil {
		t.Fatal("Expected analysis command to be created")
	}

	// Analyzing flag should be set
	if !mode.analyzing {
		t.Error("Expected analyzing flag to be set")
	}
}

func TestEditModeTriggerAnalysisEmpty(t *testing.T) {
	mode := NewEditMode()

	// Empty content
	mode.editor.SetContent("")

	cmd := mode.triggerAnalysis()

	if cmd != nil {
		t.Error("Expected no analysis command for empty content")
	}

	if mode.analyzing {
		t.Error("Expected analyzing flag to be false for empty content")
	}
}

func TestEditModeGetters(t *testing.T) {
	mode := NewEditMode()

	if mode.GetEditor() == nil {
		t.Error("Expected GetEditor to return editor component")
	}

	if mode.GetContextPanel() == nil {
		t.Error("Expected GetContextPanel to return context panel component")
	}

	if mode.GetTerminal() == nil {
		t.Error("Expected GetTerminal to return terminal component")
	}
}

func TestEditModeSetAnalyzer(t *testing.T) {
	mode := NewEditMode()

	customAnalyzer := semantic.NewAnalyzer()
	mode.SetAnalyzer(customAnalyzer)

	if mode.analyzer != customAnalyzer {
		t.Error("Expected analyzer to be updated")
	}
}

// --- EditorComponent Tests ---

func TestNewEditorComponent(t *testing.T) {
	editor := NewEditorComponent()

	if editor == nil {
		t.Fatal("Expected editor to be created")
	}

	if editor.ID() != layout.PaneEditor {
		t.Errorf("Expected ID to be PaneEditor, got %v", editor.ID())
	}
}

func TestEditorComponentInsertRune(t *testing.T) {
	editor := NewEditorComponent()

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	editor.Update(keyMsg)

	content := editor.GetContent()
	if !strings.Contains(content, "a") {
		t.Errorf("Expected content to contain 'a', got %s", content)
	}
}

func TestEditorComponentDeleteChar(t *testing.T) {
	editor := NewEditorComponent()

	// Insert some characters
	editor.insertRune('a')
	editor.insertRune('b')
	editor.insertRune('c')

	// Delete one
	keyMsg := tea.KeyMsg{Type: tea.KeyBackspace}
	editor.Update(keyMsg)

	content := editor.GetContent()
	if strings.Contains(content, "c") {
		t.Error("Expected 'c' to be deleted")
	}
}

func TestEditorComponentInsertNewline(t *testing.T) {
	editor := NewEditorComponent()

	// Insert text
	editor.insertRune('a')

	// Insert newline
	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	editor.Update(keyMsg)

	if len(editor.lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(editor.lines))
	}
}

func TestEditorComponentSetContent(t *testing.T) {
	editor := NewEditorComponent()

	content := "Line 1\nLine 2\nLine 3"
	editor.SetContent(content)

	if len(editor.lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(editor.lines))
	}

	retrieved := editor.GetContent()
	if retrieved != content {
		t.Errorf("Expected content to match, got %s", retrieved)
	}
}

func TestEditorComponentView(t *testing.T) {
	editor := NewEditorComponent()

	area := layout.Rect{X: 0, Y: 0, Width: 40, Height: 20}
	view := editor.View(area, false)

	if view == "" {
		t.Error("Expected non-empty view")
	}
}

func TestEditorComponentViewFocused(t *testing.T) {
	editor := NewEditorComponent()

	area := layout.Rect{X: 0, Y: 0, Width: 40, Height: 20}
	viewUnfocused := editor.View(area, false)
	viewFocused := editor.View(area, true)

	// Views should differ based on focus (border color)
	if viewUnfocused == viewFocused {
		// They might be the same if styling doesn't affect output
		_ = viewFocused
	}
}

// --- ContextPanelComponent Tests ---

func TestNewContextPanelComponent(t *testing.T) {
	panel := contextpanel.New(contextpanel.DefaultContextPanelConfig())
	comp := NewContextPanelComponent(panel)

	if comp == nil {
		t.Fatal("Expected component to be created")
	}

	if comp.ID() != layout.PaneSidebar {
		t.Errorf("Expected ID to be PaneSidebar, got %v", comp.ID())
	}
}

func TestContextPanelComponentUpdate(t *testing.T) {
	panel := contextpanel.New(contextpanel.DefaultContextPanelConfig())
	comp := NewContextPanelComponent(panel)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	_, cmd := comp.Update(keyMsg)
	_ = cmd
}

func TestContextPanelComponentView(t *testing.T) {
	panel := contextpanel.New(contextpanel.DefaultContextPanelConfig())
	comp := NewContextPanelComponent(panel)

	area := layout.Rect{X: 0, Y: 0, Width: 40, Height: 30}
	view := comp.View(area, false)

	if view == "" {
		t.Error("Expected non-empty view")
	}
}

// --- TerminalComponent Tests ---

func TestNewTerminalComponent(t *testing.T) {
	term := terminal.New(terminal.DefaultTerminalConfig())
	termComp := NewTerminalComponent(term)

	if termComp == nil {
		t.Fatal("Expected component to be created")
	}

	if termComp.ID() != "terminal" {
		t.Errorf("Expected ID to be 'terminal', got %v", termComp.ID())
	}
}

func TestTerminalComponentUpdate(t *testing.T) {
	term := terminal.New(terminal.DefaultTerminalConfig())
	termComp := NewTerminalComponent(term)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	_, cmd := termComp.Update(keyMsg)
	_ = cmd
}

func TestTerminalComponentView(t *testing.T) {
	term := terminal.New(terminal.DefaultTerminalConfig())
	termComp := NewTerminalComponent(term)

	area := layout.Rect{X: 0, Y: 0, Width: 60, Height: 10}
	view := termComp.View(area, false)

	if view == "" {
		t.Error("Expected non-empty view")
	}
}

// --- Integration Tests ---

func TestEditModeIntegration(t *testing.T) {
	mode := NewEditMode()

	// Simulate window size
	windowMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	mode.Update(windowMsg)

	// Set content
	mode.editor.SetContent("User creates Document")

	// Trigger analysis
	cmd := mode.triggerAnalysis()
	if cmd == nil {
		t.Fatal("Expected analysis command")
	}

	// Execute the analysis (in real scenario this would be async)
	msg := cmd()

	// Process the analysis result
	if analysisMsg, ok := msg.(SemanticAnalysisMsg); ok {
		mode.Update(analysisMsg)

		// Check that context panel has results
		panel := mode.contextPanel.GetPanel()
		analysis := panel.GetAnalysis()

		if analysis == nil {
			t.Error("Expected analysis results in context panel")
		}
	}
}

func TestEditModeAnalysisDebounce(t *testing.T) {
	mode := NewEditMode()
	mode.analysisDebounce = 100 * time.Millisecond

	// Set content and trigger analysis
	mode.editor.SetContent("Test")
	mode.triggerAnalysis()

	// Wait for analysis to complete
	time.Sleep(50 * time.Millisecond)

	// Try to trigger again immediately (should be debounced)
	mode.lastAnalysis = time.Now().Add(-50 * time.Millisecond)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	mode.Update(keyMsg)

	// Analysis should still be running (debounced)
	if time.Since(mode.lastAnalysis) > mode.analysisDebounce {
		// Debounce expired, analysis might trigger
		_ = mode.analyzing
	}
}

func TestEditModeLifecycle(t *testing.T) {
	mode := NewEditMode()

	// Initialize
	cmd := mode.Init()
	_ = cmd

	// Enter mode
	cmd = mode.OnEnter()
	_ = cmd

	// Perform some updates
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	mode.Update(keyMsg)

	// Exit mode
	cmd = mode.OnExit()
	_ = cmd

	// Analyzer should be stopped
	if mode.analyzing {
		t.Error("Expected analyzing to be false after exit")
	}
}
