// Package orchestrate provides types and structures for Orchestrate Mode.
package orchestrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditorMode represents the current editor mode.
type EditorMode int

const (
	ModeEdit EditorMode = iota
	ModeSave
	ModeLoad
)

// PlanEditor is a Bubble Tea model for editing work plan JSON files.
type PlanEditor struct {
	// Content
	content string    // JSON text being edited
	plan    *WorkPlan // Parsed plan (nil if invalid)
	cursor  int       // Cursor position in content

	// Validation
	validationErrors []string
	isValid          bool

	// File operations
	filename string
	dirty    bool // Has unsaved changes

	// UI state
	width      int
	height     int
	viewOffset int // For scrolling (line offset)

	// Mode
	mode EditorMode

	// Dialog
	dialogInput  string // For save/load dialogs
	dialogActive bool
}

// Styling functions
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0066FF")).
			Padding(0, 1)

	lineNumStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00AA00"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)
)

// NewPlanEditor creates a new plan editor instance.
func NewPlanEditor() *PlanEditor {
	return &PlanEditor{
		content:          "",
		plan:             nil,
		cursor:           0,
		validationErrors: []string{},
		isValid:          false,
		filename:         "",
		dirty:            false,
		width:            80,
		height:           24,
		viewOffset:       0,
		mode:             ModeEdit,
		dialogInput:      "",
		dialogActive:     false,
	}
}

// Init initializes the editor model.
func (m *PlanEditor) Init() tea.Cmd {
	return nil
}

// Update processes messages.
func (m *PlanEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress handles key press events.
func (m *PlanEditor) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If a dialog is active, handle dialog-specific input
	if m.dialogActive {
		return m.handleDialogInput(msg)
	}

	switch msg.String() {
	case "ctrl+s":
		// Save file
		m.mode = ModeSave
		m.dialogActive = true
		m.dialogInput = m.filename
		return m, nil

	case "ctrl+o":
		// Load file
		m.mode = ModeLoad
		m.dialogActive = true
		m.dialogInput = ""
		return m, nil

	case "ctrl+n":
		// New plan
		m.content = `{
  "name": "New Plan",
  "description": "A new work plan",
  "maxConcurrent": 2,
  "tasks": []
}`
		m.cursor = 0
		m.filename = ""
		m.dirty = true
		m.validateContent()
		return m, nil

	case "ctrl+q":
		// Quit
		return m, tea.Quit

	case "up":
		m.moveCursorUp()
		return m, nil

	case "down":
		m.moveCursorDown()
		return m, nil

	case "left":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case "right":
		if m.cursor < len(m.content) {
			m.cursor++
		}
		return m, nil

	case "home":
		m.moveCursorToLineStart()
		return m, nil

	case "end":
		m.moveCursorToLineEnd()
		return m, nil

	case "backspace":
		m.deleteChar()
		return m, nil

	case "delete":
		m.deleteCharForward()
		return m, nil

	case "enter":
		m.insertNewline()
		return m, nil

	default:
		// Handle regular character input
		if len(msg.String()) == 1 && msg.Runes[0] >= 32 && msg.Runes[0] < 127 {
			m.insertChar(msg.Runes[0])
		}
		return m, nil
	}
}

// handleDialogInput processes input in dialog mode.
func (m *PlanEditor) handleDialogInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Confirm dialog
		if m.mode == ModeSave {
			if m.dialogInput != "" {
				m.filename = m.dialogInput
				if err := m.save(); err != nil {
					m.validationErrors = []string{fmt.Sprintf("Save error: %v", err)}
				} else {
					m.dirty = false
				}
			}
		} else if m.mode == ModeLoad {
			if m.dialogInput != "" {
				if err := m.load(m.dialogInput); err != nil {
					m.validationErrors = []string{fmt.Sprintf("Load error: %v", err)}
				} else {
					m.filename = m.dialogInput
					m.dirty = false
				}
			}
		}
		m.dialogActive = false
		m.mode = ModeEdit
		return m, nil

	case "esc":
		// Cancel dialog
		m.dialogActive = false
		m.mode = ModeEdit
		return m, nil

	case "backspace":
		if len(m.dialogInput) > 0 {
			m.dialogInput = m.dialogInput[:len(m.dialogInput)-1]
		}
		return m, nil

	default:
		if len(msg.String()) == 1 && msg.Runes[0] >= 32 && msg.Runes[0] < 127 {
			m.dialogInput += string(msg.Runes[0])
		}
		return m, nil
	}
}

// insertChar inserts a character at the cursor position.
func (m *PlanEditor) insertChar(ch rune) {
	before := m.content[:m.cursor]
	after := m.content[m.cursor:]
	m.content = before + string(ch) + after
	m.cursor++
	m.dirty = true
	m.validateContent()
}

// insertNewline inserts a newline at the cursor position.
func (m *PlanEditor) insertNewline() {
	before := m.content[:m.cursor]
	after := m.content[m.cursor:]
	m.content = before + "\n" + after
	m.cursor++
	m.dirty = true
	m.validateContent()
}

// deleteChar deletes the character before the cursor (backspace).
func (m *PlanEditor) deleteChar() {
	if m.cursor > 0 {
		before := m.content[:m.cursor-1]
		after := m.content[m.cursor:]
		m.content = before + after
		m.cursor--
		m.dirty = true
		m.validateContent()
	}
}

// deleteCharForward deletes the character at the cursor position.
func (m *PlanEditor) deleteCharForward() {
	if m.cursor < len(m.content) {
		before := m.content[:m.cursor]
		after := m.content[m.cursor+1:]
		m.content = before + after
		m.dirty = true
		m.validateContent()
	}
}

// moveCursorUp moves the cursor up one line.
func (m *PlanEditor) moveCursorUp() {
	lines := strings.Split(m.content, "\n")

	// Find current line and column
	pos := 0
	var currentLine, targetLine int
	var lineStart int

	for i, line := range lines {
		lineLen := len(line) + 1 // +1 for newline
		if pos+lineLen > m.cursor {
			currentLine = i
			lineStart = pos
			break
		}
		pos += lineLen
	}

	if currentLine == 0 {
		return // Already at first line
	}

	col := m.cursor - lineStart
	targetLine = currentLine - 1

	// Calculate new position
	pos = 0
	for i := 0; i < targetLine; i++ {
		pos += len(lines[i]) + 1
	}

	// Adjust column if target line is shorter
	if col > len(lines[targetLine]) {
		col = len(lines[targetLine])
	}

	m.cursor = pos + col
}

// moveCursorDown moves the cursor down one line.
func (m *PlanEditor) moveCursorDown() {
	lines := strings.Split(m.content, "\n")

	// Find current line and column
	pos := 0
	var currentLine int
	var lineStart int

	for i, line := range lines {
		lineLen := len(line) + 1 // +1 for newline
		if pos+lineLen > m.cursor {
			currentLine = i
			lineStart = pos
			break
		}
		pos += lineLen
	}

	if currentLine >= len(lines)-1 {
		return // Already at last line
	}

	col := m.cursor - lineStart
	targetLine := currentLine + 1

	// Calculate new position
	pos = 0
	for i := 0; i < targetLine; i++ {
		pos += len(lines[i]) + 1
	}

	// Adjust column if target line is shorter
	if col > len(lines[targetLine]) {
		col = len(lines[targetLine])
	}

	m.cursor = pos + col
}

// moveCursorToLineStart moves the cursor to the start of the current line.
func (m *PlanEditor) moveCursorToLineStart() {
	pos := 0
	for i, ch := range m.content {
		if i == m.cursor {
			break
		}
		if ch == '\n' {
			pos = i + 1
		}
	}
	m.cursor = pos
}

// moveCursorToLineEnd moves the cursor to the end of the current line.
func (m *PlanEditor) moveCursorToLineEnd() {
	pos := m.cursor
	for i := m.cursor; i < len(m.content); i++ {
		if m.content[i] == '\n' {
			pos = i
			break
		}
		pos = i + 1
	}
	m.cursor = pos
}

// validateContent validates the JSON content.
func (m *PlanEditor) validateContent() {
	// Try to parse JSON
	var plan WorkPlan
	if err := json.Unmarshal([]byte(m.content), &plan); err != nil {
		m.isValid = false
		m.validationErrors = []string{err.Error()}
		m.plan = nil
		return
	}

	// Validate plan structure
	if err := plan.Validate(); err != nil {
		m.isValid = false
		m.validationErrors = []string{err.Error()}
		m.plan = nil
		return
	}

	m.isValid = true
	m.validationErrors = nil
	m.plan = &plan
}

// save saves the content to a file.
func (m *PlanEditor) save() error {
	if m.filename == "" {
		return fmt.Errorf("no filename specified")
	}

	// Ensure the directory exists
	dir := filepath.Dir(m.filename)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(m.filename, []byte(m.content), 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// load loads content from a file.
func (m *PlanEditor) load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	m.content = string(data)
	m.cursor = 0
	m.dirty = false
	m.validateContent()

	return nil
}

// View renders the editor view.
func (m *PlanEditor) View() string {
	var lines []string

	// Header
	header := "Plan Editor"
	if m.dirty {
		header += " *" // Unsaved changes
	}
	if m.isValid {
		header += " ✓" // Valid
	} else if m.content != "" {
		header += " ✗" // Invalid
	}
	if m.filename != "" {
		header += " (" + filepath.Base(m.filename) + ")"
	}
	lines = append(lines, headerStyle.Render(header))

	// Dialog overlay
	if m.dialogActive {
		prompt := "Save as: "
		if m.mode == ModeLoad {
			prompt = "Load from: "
		}
		lines = append(lines, prompt+m.dialogInput)
		lines = append(lines, helpStyle.Render("Press Enter to confirm, Esc to cancel"))
		return lipgloss.JoinVertical(lipgloss.Left, lines...)
	}

	// Content with line numbers
	contentLines := strings.Split(m.content, "\n")
	maxLines := m.height - 8 // Reserve space for header, footer, validation

	startLine := m.viewOffset
	endLine := startLine + maxLines
	if endLine > len(contentLines) {
		endLine = len(contentLines)
	}

	lines = append(lines, "")

	for i := startLine; i < endLine; i++ {
		line := contentLines[i]
		lineNum := fmt.Sprintf("%3d │ ", i+1)

		// Add cursor indicator
		indicator := "  "
		if i < len(contentLines) {
			lineStart := 0
			for j := 0; j < i; j++ {
				lineStart += len(contentLines[j]) + 1
			}

			if m.cursor >= lineStart && m.cursor <= lineStart+len(line) {
				indicator = "▶ "
			}
		}

		lines = append(lines, indicator+lineNumStyle.Render(lineNum)+line)
	}

	// Validation status
	lines = append(lines, "")
	if len(m.validationErrors) > 0 {
		lines = append(lines, errorStyle.Render("Errors:"))
		for _, err := range m.validationErrors {
			lines = append(lines, errorStyle.Render("  • "+err))
		}
	} else if m.isValid && m.content != "" {
		lines = append(lines, successStyle.Render("Valid JSON ✓"))
	}

	// Footer
	lines = append(lines, "")
	footer := "Ctrl+S: Save  Ctrl+O: Open  Ctrl+N: New  Ctrl+Q: Quit"
	lines = append(lines, helpStyle.Render(footer))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
