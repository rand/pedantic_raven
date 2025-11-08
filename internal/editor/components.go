// Package editor provides the Edit Mode for context editing with semantic analysis.
//
// Components:
// - EditorComponent: Text editing with syntax highlighting
// - ContextPanelComponent: Displays semantic analysis results
// - TerminalComponent: Integrated terminal for mnemosyne commands
package editor

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/context"
	"github.com/rand/pedantic-raven/internal/editor/buffer"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"github.com/rand/pedantic-raven/internal/layout"
	"github.com/rand/pedantic-raven/internal/terminal"
)

var (
	editorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1)

	contextStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1)

	terminalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("120")).
			Padding(1)

	focusedBorderColor = lipgloss.Color("170")
	normalBorderColor  = lipgloss.Color("240")
)

// --- EditorComponent ---

// EditorComponent provides text editing functionality with full undo/redo support.
type EditorComponent struct {
	buffer buffer.Buffer
}

// NewEditorComponent creates a new editor component.
func NewEditorComponent() *EditorComponent {
	return &EditorComponent{
		buffer: buffer.NewBuffer("editor-0"),
	}
}

// ID implements layout.Component.
func (e *EditorComponent) ID() layout.PaneID {
	return layout.PaneEditor
}

// Update implements layout.Component.
func (e *EditorComponent) Update(msg tea.Msg) (layout.Component, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+z":
			// Undo
			e.buffer.Undo()
			return e, nil

		case "ctrl+y", "ctrl+shift+z":
			// Redo
			e.buffer.Redo()
			return e, nil

		case "backspace":
			e.deleteChar()

		case "enter":
			e.insertNewline()

		default:
			// Handle regular key presses
			if keyMsg.Type == tea.KeyRunes {
				for _, r := range keyMsg.Runes {
					e.insertRune(r)
				}
			}
		}
	}
	return e, nil
}

// View implements layout.Component.
func (e *EditorComponent) View(area layout.Rect, focused bool) string {
	style := editorStyle.
		Width(area.Width - 4).
		Height(area.Height - 4)

	if focused {
		style = style.BorderForeground(focusedBorderColor)
	} else {
		style = style.BorderForeground(normalBorderColor)
	}

	// Get content from buffer
	content := e.buffer.Content()
	if content == "" {
		content = "(empty)"
	}

	return style.Render(content)
}

// GetContent returns the current editor content.
func (e *EditorComponent) GetContent() string {
	return e.buffer.Content()
}

// SetContent sets the editor content.
func (e *EditorComponent) SetContent(content string) {
	// Clear and recreate buffer with new content
	e.buffer.Clear()
	if content != "" {
		pos := buffer.Position{Line: 0, Column: 0}
		e.buffer.Insert(pos, content)
	}
}

// OpenFile opens a file from disk and loads it into the editor.
// Returns an error if the file cannot be read or if there are encoding issues.
func (e *EditorComponent) OpenFile(path string) error {
	// Read file contents
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// TODO: Detect encoding - for now assume UTF-8
	content := string(data)

	// Set buffer path and content
	e.buffer.SetPath(path)
	e.SetContent(content)

	// Mark as clean since we just loaded from disk
	e.buffer.SetDirty(false)

	return nil
}

// SaveFile saves the current buffer to its associated file path.
// Returns an error if the buffer has no path or if writing fails.
func (e *EditorComponent) SaveFile() error {
	path := e.buffer.Path()
	if path == "" {
		return fmt.Errorf("no file path set, use SaveFileAs")
	}

	return e.SaveFileAs(path)
}

// SaveFileAs saves the current buffer to the specified path using atomic write.
// Uses a temporary file and rename to ensure atomic writes.
func (e *EditorComponent) SaveFileAs(path string) error {
	content := e.buffer.Content()

	// Write to temporary file first (atomic write)
	tmpPath := path + ".tmp"
	err := os.WriteFile(tmpPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Rename temp file to target (atomic on most filesystems)
	err = os.Rename(tmpPath, path)
	if err != nil {
		// Clean up temp file on failure
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	// Update buffer path and mark as clean
	e.buffer.SetPath(path)
	e.buffer.SetDirty(false)

	return nil
}

// GetFilePath returns the current file path associated with the buffer.
func (e *EditorComponent) GetFilePath() string {
	return e.buffer.Path()
}

// IsDirty returns true if the buffer has unsaved changes.
func (e *EditorComponent) IsDirty() bool {
	return e.buffer.IsDirty()
}

func (e *EditorComponent) insertRune(r rune) {
	pos := e.buffer.Cursor()
	e.buffer.Insert(pos, string(r))
	// Move cursor forward after insertion
	e.buffer.SetCursor(buffer.Position{Line: pos.Line, Column: pos.Column + 1})
}

func (e *EditorComponent) deleteChar() {
	pos := e.buffer.Cursor()

	// Delete previous character
	if pos.Column > 0 {
		from := buffer.Position{Line: pos.Line, Column: pos.Column - 1}
		to := pos
		e.buffer.Delete(from, to)
		// Move cursor back after deletion
		e.buffer.SetCursor(from)
	} else if pos.Line > 0 {
		// Delete newline at end of previous line
		prevLine := e.buffer.Line(pos.Line - 1)
		from := buffer.Position{Line: pos.Line - 1, Column: len(prevLine)}
		to := buffer.Position{Line: pos.Line, Column: 0}
		e.buffer.Delete(from, to)
		// Move cursor to end of previous line after deletion
		e.buffer.SetCursor(from)
	}
}

func (e *EditorComponent) insertNewline() {
	pos := e.buffer.Cursor()
	e.buffer.Insert(pos, "\n")
	// Move cursor to start of next line after newline insertion
	e.buffer.SetCursor(buffer.Position{Line: pos.Line + 1, Column: 0})
}

// --- ContextPanelComponent ---

// ContextPanelComponent wraps the context panel for semantic analysis display.
type ContextPanelComponent struct {
	panel *context.ContextPanel
}

// NewContextPanelComponent creates a new context panel component.
func NewContextPanelComponent(panel *context.ContextPanel) *ContextPanelComponent {
	return &ContextPanelComponent{
		panel: panel,
	}
}

// ID implements layout.Component.
func (c *ContextPanelComponent) ID() layout.PaneID {
	return layout.PaneSidebar
}

// Update implements layout.Component.
func (c *ContextPanelComponent) Update(msg tea.Msg) (layout.Component, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "down", "j":
			c.panel.ScrollDown(1)
		case "up", "k":
			c.panel.ScrollUp(1)
		case "pgdown":
			c.panel.ScrollDown(10)
		case "pgup":
			c.panel.ScrollUp(10)
		case "home":
			c.panel.ScrollToTop()
		case "end":
			c.panel.ScrollToBottom()
		case "tab":
			c.panel.NextSection()
		case "shift+tab":
			c.panel.PreviousSection()
		case "enter":
			section := c.panel.GetActiveSection()
			c.panel.ToggleSection(section)
		}
	}
	return c, nil
}

// View implements layout.Component.
func (c *ContextPanelComponent) View(area layout.Rect, focused bool) string {
	style := contextStyle.
		Width(area.Width - 4).
		Height(area.Height - 4)

	if focused {
		style = style.BorderForeground(focusedBorderColor)
	} else {
		style = style.BorderForeground(normalBorderColor)
	}

	// Update panel size
	config := c.panel.GetConfig()
	config.Width = area.Width - 6
	config.Height = area.Height - 6
	c.panel.SetConfig(config)

	// Render panel
	content := c.panel.Render()
	return style.Render(content)
}

// GetPanel returns the underlying context panel.
func (c *ContextPanelComponent) GetPanel() *context.ContextPanel {
	return c.panel
}

// --- TerminalComponent ---

// TerminalComponent wraps the terminal for command execution.
type TerminalComponent struct {
	term *terminal.Terminal
}

// NewTerminalComponent creates a new terminal component.
func NewTerminalComponent(term *terminal.Terminal) *TerminalComponent {
	return &TerminalComponent{
		term: term,
	}
}

// ID implements layout.Component.
func (t *TerminalComponent) ID() layout.PaneID {
	return "terminal"
}

// Update implements layout.Component.
func (t *TerminalComponent) Update(msg tea.Msg) (layout.Component, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyRunes:
			// Insert character
			for _, r := range keyMsg.Runes {
				t.term.InsertChar(r)
			}
		case tea.KeyBackspace:
			t.term.DeleteChar()
		case tea.KeyEnter:
			t.term.Submit()
		case tea.KeyUp:
			t.term.HistoryPrevious()
		case tea.KeyDown:
			t.term.HistoryNext()
		case tea.KeyLeft:
			t.term.MoveCursorLeft()
		case tea.KeyRight:
			t.term.MoveCursorRight()
		case tea.KeyHome:
			t.term.MoveCursorToStart()
		case tea.KeyEnd:
			t.term.MoveCursorToEnd()
		case tea.KeyPgUp:
			t.term.ScrollUp(5)
		case tea.KeyPgDown:
			t.term.ScrollDown(5)
		}
	}
	return t, nil
}

// View implements layout.Component.
func (t *TerminalComponent) View(area layout.Rect, focused bool) string {
	style := terminalStyle.
		Width(area.Width - 4).
		Height(area.Height - 4)

	if focused {
		style = style.BorderForeground(focusedBorderColor)
	} else {
		style = style.BorderForeground(normalBorderColor)
	}

	// Update terminal size
	config := t.term.GetConfig()
	config.Width = area.Width - 6
	config.Height = area.Height - 6
	t.term.SetConfig(config)

	// Build terminal view
	var lines []string

	// Add output lines
	visibleOutput := t.term.GetVisibleOutput()
	lines = append(lines, visibleOutput...)

	// Add input line
	prompt := config.Prompt
	input := t.term.GetInput()
	cursorPos := t.term.GetCursorPos()

	// Simple cursor rendering
	inputLine := prompt + input
	if cursorPos < len(input) {
		// Insert cursor marker
		before := input[:cursorPos]
		at := input[cursorPos : cursorPos+1]
		after := input[cursorPos+1:]
		inputLine = prompt + before + "[" + at + "]" + after
	} else {
		inputLine += "█"
	}

	lines = append(lines, inputLine)

	content := strings.Join(lines, "\n")
	return style.Render(content)
}

// GetTerminal returns the underlying terminal.
func (t *TerminalComponent) GetTerminal() *terminal.Terminal {
	return t.term
}

// --- Semantic Analysis Message ---

// SemanticAnalysisMsg is sent when semantic analysis completes.
type SemanticAnalysisMsg struct {
	Analysis *semantic.Analysis
}
