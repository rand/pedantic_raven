// Package editor provides the Edit Mode for context editing with semantic analysis.
//
// Components:
// - EditorComponent: Text editing with syntax highlighting
// - ContextPanelComponent: Displays semantic analysis results
// - TerminalComponent: Integrated terminal for mnemosyne commands
package editor

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/context"
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

// EditorComponent provides text editing functionality.
type EditorComponent struct {
	content string
	cursor  int
	lines   []string
}

// NewEditorComponent creates a new editor component.
func NewEditorComponent() *EditorComponent {
	return &EditorComponent{
		content: "",
		cursor:  0,
		lines:   []string{""},
	}
}

// ID implements layout.Component.
func (e *EditorComponent) ID() layout.PaneID {
	return layout.PaneEditor
}

// Update implements layout.Component.
func (e *EditorComponent) Update(msg tea.Msg) (layout.Component, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyRunes:
			// Insert character
			e.insertRune(keyMsg.Runes[0])
		case tea.KeyBackspace:
			e.deleteChar()
		case tea.KeyEnter:
			e.insertNewline()
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

	// Build content view
	var lines []string
	if len(e.lines) == 0 {
		lines = []string{"(empty)"}
	} else {
		lines = e.lines
	}

	content := strings.Join(lines, "\n")
	if content == "" {
		content = "(empty)"
	}

	return style.Render(content)
}

// GetContent returns the current editor content.
func (e *EditorComponent) GetContent() string {
	return strings.Join(e.lines, "\n")
}

// SetContent sets the editor content.
func (e *EditorComponent) SetContent(content string) {
	e.content = content
	e.lines = strings.Split(content, "\n")
	if len(e.lines) == 0 {
		e.lines = []string{""}
	}
}

func (e *EditorComponent) insertRune(r rune) {
	// Simple single-line insertion for now
	if len(e.lines) == 0 {
		e.lines = []string{string(r)}
	} else {
		e.lines[len(e.lines)-1] += string(r)
	}
}

func (e *EditorComponent) deleteChar() {
	if len(e.lines) == 0 {
		return
	}
	lastLine := e.lines[len(e.lines)-1]
	if len(lastLine) > 0 {
		e.lines[len(e.lines)-1] = lastLine[:len(lastLine)-1]
	}
}

func (e *EditorComponent) insertNewline() {
	e.lines = append(e.lines, "")
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
