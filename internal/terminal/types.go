// Package terminal provides an integrated terminal for executing commands.
//
// The Terminal component supports:
// - Shell commands (ls, git, etc.)
// - mnemosyne CLI integration
// - Built-in commands (:clear, :help)
// - Command history with navigation
// - Output capture and display
// - Scrollable output buffer
package terminal

import (
	"strings"
	"time"
)

// CommandType indicates the type of command.
type CommandType int

const (
	CommandShell CommandType = iota // Shell command
	CommandMnemosyne                 // mnemosyne CLI command
	CommandBuiltin                   // Built-in command
)

// Command represents a command to execute.
type Command struct {
	Type      CommandType // Command type
	Input     string      // Raw input
	Args      []string    // Parsed arguments
	Timestamp time.Time   // When command was entered
}

// CommandResult represents the result of command execution.
type CommandResult struct {
	Command  Command       // Original command
	Output   string        // Command output (stdout)
	Error    string        // Error output (stderr)
	ExitCode int           // Exit code
	Duration time.Duration // Execution duration
}

// HistoryEntry represents a command in history.
type HistoryEntry struct {
	Command Command
	Result  *CommandResult
}

// TerminalConfig configures the terminal.
type TerminalConfig struct {
	Width          int    // Terminal width
	Height         int    // Terminal height
	MaxHistory     int    // Maximum history entries
	MaxOutputLines int    // Maximum output lines to retain
	Prompt         string // Command prompt
	WorkingDir     string // Working directory for shell commands
}

// DefaultTerminalConfig returns default configuration.
func DefaultTerminalConfig() TerminalConfig {
	return TerminalConfig{
		Width:          80,
		Height:         20,
		MaxHistory:     100,
		MaxOutputLines: 1000,
		Prompt:         "$ ",
		WorkingDir:     ".",
	}
}

// Terminal provides command execution and output display.
type Terminal struct {
	config TerminalConfig

	// Command history
	history        []HistoryEntry
	historyIndex   int  // Current position in history (-1 = no navigation)
	historyChanged bool // Whether history was modified

	// Current input
	input       string // Current command input
	cursorPos   int    // Cursor position in input
	completions []string

	// Output buffer
	output       []string // Lines of output
	scrollOffset int      // Scroll position in output

	// Execution state
	executing bool   // Whether a command is executing
	lastResult *CommandResult
}

// New creates a new terminal.
func New(config TerminalConfig) *Terminal {
	return &Terminal{
		config:       config,
		history:      make([]HistoryEntry, 0),
		historyIndex: -1,
		input:        "",
		cursorPos:    0,
		output:       make([]string, 0),
		scrollOffset: 0,
		executing:    false,
	}
}

// GetInput returns the current input.
func (t *Terminal) GetInput() string {
	return t.input
}

// SetInput sets the current input.
func (t *Terminal) SetInput(input string) {
	t.input = input
	t.cursorPos = len(input)
}

// InsertChar inserts a character at the cursor position.
func (t *Terminal) InsertChar(ch rune) {
	before := t.input[:t.cursorPos]
	after := t.input[t.cursorPos:]
	t.input = before + string(ch) + after
	t.cursorPos++
}

// DeleteChar deletes the character before the cursor.
func (t *Terminal) DeleteChar() {
	if t.cursorPos > 0 {
		before := t.input[:t.cursorPos-1]
		after := t.input[t.cursorPos:]
		t.input = before + after
		t.cursorPos--
	}
}

// MoveCursorLeft moves the cursor left.
func (t *Terminal) MoveCursorLeft() {
	if t.cursorPos > 0 {
		t.cursorPos--
	}
}

// MoveCursorRight moves the cursor right.
func (t *Terminal) MoveCursorRight() {
	if t.cursorPos < len(t.input) {
		t.cursorPos++
	}
}

// MoveCursorToStart moves cursor to start of line.
func (t *Terminal) MoveCursorToStart() {
	t.cursorPos = 0
}

// MoveCursorToEnd moves cursor to end of line.
func (t *Terminal) MoveCursorToEnd() {
	t.cursorPos = len(t.input)
}

// GetCursorPos returns the current cursor position.
func (t *Terminal) GetCursorPos() int {
	return t.cursorPos
}

// Clear clears the terminal output.
func (t *Terminal) Clear() {
	t.output = make([]string, 0)
	t.scrollOffset = 0
}

// AddOutput adds a line to the output buffer.
func (t *Terminal) AddOutput(line string) {
	t.output = append(t.output, line)

	// Trim if exceeds max
	if len(t.output) > t.config.MaxOutputLines {
		t.output = t.output[len(t.output)-t.config.MaxOutputLines:]
	}

	// Auto-scroll to bottom
	t.ScrollToBottom()
}

// AddOutputLines adds multiple lines to the output buffer.
func (t *Terminal) AddOutputLines(lines []string) {
	for _, line := range lines {
		t.AddOutput(line)
	}
}

// GetOutput returns the output buffer.
func (t *Terminal) GetOutput() []string {
	return t.output
}

// GetVisibleOutput returns the currently visible output lines.
func (t *Terminal) GetVisibleOutput() []string {
	if len(t.output) == 0 {
		return []string{}
	}

	start := t.scrollOffset
	end := t.scrollOffset + t.config.Height - 1 // Reserve 1 line for input

	if start >= len(t.output) {
		start = len(t.output) - 1
	}

	if end > len(t.output) {
		end = len(t.output)
	}

	if start < 0 {
		start = 0
	}

	return t.output[start:end]
}

// ScrollDown scrolls the output down.
func (t *Terminal) ScrollDown(lines int) {
	t.scrollOffset += lines
	maxOffset := len(t.output) - (t.config.Height - 1)
	if maxOffset < 0 {
		maxOffset = 0
	}
	if t.scrollOffset > maxOffset {
		t.scrollOffset = maxOffset
	}
}

// ScrollUp scrolls the output up.
func (t *Terminal) ScrollUp(lines int) {
	t.scrollOffset -= lines
	if t.scrollOffset < 0 {
		t.scrollOffset = 0
	}
}

// ScrollToTop scrolls to the top of output.
func (t *Terminal) ScrollToTop() {
	t.scrollOffset = 0
}

// ScrollToBottom scrolls to the bottom of output.
func (t *Terminal) ScrollToBottom() {
	maxOffset := len(t.output) - (t.config.Height - 1)
	if maxOffset < 0 {
		maxOffset = 0
	}
	t.scrollOffset = maxOffset
}

// GetScrollOffset returns the current scroll offset.
func (t *Terminal) GetScrollOffset() int {
	return t.scrollOffset
}

// AddHistory adds a command to history.
func (t *Terminal) AddHistory(entry HistoryEntry) {
	t.history = append(t.history, entry)

	// Trim if exceeds max
	if len(t.history) > t.config.MaxHistory {
		t.history = t.history[len(t.history)-t.config.MaxHistory:]
	}

	// Reset history navigation
	t.historyIndex = -1
	t.historyChanged = true
}

// GetHistory returns the command history.
func (t *Terminal) GetHistory() []HistoryEntry {
	return t.history
}

// HistoryPrevious navigates to the previous command in history.
func (t *Terminal) HistoryPrevious() bool {
	if len(t.history) == 0 {
		return false
	}

	// First call: start at end of history
	if t.historyIndex == -1 {
		t.historyIndex = len(t.history) - 1
	} else if t.historyIndex > 0 {
		t.historyIndex--
	} else {
		return false // At beginning
	}

	t.input = t.history[t.historyIndex].Command.Input
	t.cursorPos = len(t.input)
	return true
}

// HistoryNext navigates to the next command in history.
func (t *Terminal) HistoryNext() bool {
	if t.historyIndex == -1 {
		return false // Not navigating
	}

	if t.historyIndex < len(t.history)-1 {
		t.historyIndex++
		t.input = t.history[t.historyIndex].Command.Input
		t.cursorPos = len(t.input)
		return true
	}

	// At end of history, clear input
	t.historyIndex = -1
	t.input = ""
	t.cursorPos = 0
	return true
}

// ResetHistoryNavigation resets history navigation state.
func (t *Terminal) ResetHistoryNavigation() {
	t.historyIndex = -1
}

// IsExecuting returns whether a command is currently executing.
func (t *Terminal) IsExecuting() bool {
	return t.executing
}

// SetExecuting sets the execution state.
func (t *Terminal) SetExecuting(executing bool) {
	t.executing = executing
}

// GetLastResult returns the last command result.
func (t *Terminal) GetLastResult() *CommandResult {
	return t.lastResult
}

// SetLastResult sets the last command result.
func (t *Terminal) SetLastResult(result *CommandResult) {
	t.lastResult = result
}

// GetConfig returns the terminal configuration.
func (t *Terminal) GetConfig() TerminalConfig {
	return t.config
}

// SetConfig updates the terminal configuration.
func (t *Terminal) SetConfig(config TerminalConfig) {
	t.config = config
}

// ParseCommand parses a command string into a Command.
func ParseCommand(input string) Command {
	input = strings.TrimSpace(input)

	cmd := Command{
		Input:     input,
		Timestamp: time.Now(),
	}

	// Empty command
	if input == "" {
		cmd.Type = CommandShell
		return cmd
	}

	// Built-in commands start with ':'
	if strings.HasPrefix(input, ":") {
		cmd.Type = CommandBuiltin
		cmd.Args = strings.Fields(input[1:])
		return cmd
	}

	// mnemosyne commands
	if strings.HasPrefix(input, "mnemosyne ") {
		cmd.Type = CommandMnemosyne
		cmd.Args = strings.Fields(input)
		return cmd
	}

	// Shell command
	cmd.Type = CommandShell
	cmd.Args = strings.Fields(input)
	return cmd
}
