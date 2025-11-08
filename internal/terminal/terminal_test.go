package terminal

import (
	"strings"
	"testing"
	"time"
)

// --- Terminal Creation Tests ---

func TestNewTerminal(t *testing.T) {
	config := DefaultTerminalConfig()
	term := New(config)

	if term == nil {
		t.Fatal("Expected terminal to be created")
	}

	if term.GetInput() != "" {
		t.Error("Expected initial input to be empty")
	}

	if term.GetCursorPos() != 0 {
		t.Error("Expected initial cursor position to be 0")
	}

	if len(term.GetOutput()) != 0 {
		t.Error("Expected initial output to be empty")
	}
}

func TestDefaultTerminalConfig(t *testing.T) {
	config := DefaultTerminalConfig()

	if config.Width == 0 {
		t.Error("Expected default width to be set")
	}

	if config.Height == 0 {
		t.Error("Expected default height to be set")
	}

	if config.MaxHistory == 0 {
		t.Error("Expected default max history to be set")
	}

	if config.Prompt == "" {
		t.Error("Expected default prompt to be set")
	}
}

// --- Input Handling Tests ---

func TestSetInput(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test command")

	if term.GetInput() != "test command" {
		t.Errorf("Expected 'test command', got '%s'", term.GetInput())
	}

	if term.GetCursorPos() != len("test command") {
		t.Errorf("Expected cursor at end of input, got %d", term.GetCursorPos())
	}
}

func TestInsertChar(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.InsertChar('h')
	term.InsertChar('i')

	if term.GetInput() != "hi" {
		t.Errorf("Expected 'hi', got '%s'", term.GetInput())
	}

	if term.GetCursorPos() != 2 {
		t.Errorf("Expected cursor at 2, got %d", term.GetCursorPos())
	}
}

func TestInsertCharAtPosition(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test")
	term.MoveCursorToStart()
	term.InsertChar('X')

	if term.GetInput() != "Xtest" {
		t.Errorf("Expected 'Xtest', got '%s'", term.GetInput())
	}

	if term.GetCursorPos() != 1 {
		t.Errorf("Expected cursor at 1, got %d", term.GetCursorPos())
	}
}

func TestDeleteChar(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test")
	term.DeleteChar()

	if term.GetInput() != "tes" {
		t.Errorf("Expected 'tes', got '%s'", term.GetInput())
	}

	if term.GetCursorPos() != 3 {
		t.Errorf("Expected cursor at 3, got %d", term.GetCursorPos())
	}
}

func TestDeleteCharAtStart(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test")
	term.MoveCursorToStart()
	term.DeleteChar()

	// Should not delete anything
	if term.GetInput() != "test" {
		t.Errorf("Expected 'test', got '%s'", term.GetInput())
	}
}

// --- Cursor Movement Tests ---

func TestMoveCursorLeft(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test")
	initialPos := term.GetCursorPos()

	term.MoveCursorLeft()

	if term.GetCursorPos() != initialPos-1 {
		t.Errorf("Expected cursor to move left, got %d", term.GetCursorPos())
	}
}

func TestMoveCursorRight(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test")
	term.MoveCursorToStart()
	initialPos := term.GetCursorPos()

	term.MoveCursorRight()

	if term.GetCursorPos() != initialPos+1 {
		t.Errorf("Expected cursor to move right, got %d", term.GetCursorPos())
	}
}

func TestMoveCursorToStart(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test")
	term.MoveCursorToStart()

	if term.GetCursorPos() != 0 {
		t.Errorf("Expected cursor at 0, got %d", term.GetCursorPos())
	}
}

func TestMoveCursorToEnd(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("test")
	term.MoveCursorToStart()
	term.MoveCursorToEnd()

	if term.GetCursorPos() != 4 {
		t.Errorf("Expected cursor at 4, got %d", term.GetCursorPos())
	}
}

// --- Output Management Tests ---

func TestAddOutput(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.AddOutput("line 1")
	term.AddOutput("line 2")

	output := term.GetOutput()

	if len(output) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(output))
	}

	if output[0] != "line 1" {
		t.Errorf("Expected 'line 1', got '%s'", output[0])
	}

	if output[1] != "line 2" {
		t.Errorf("Expected 'line 2', got '%s'", output[1])
	}
}

func TestAddOutputLines(t *testing.T) {
	term := New(DefaultTerminalConfig())

	lines := []string{"line 1", "line 2", "line 3"}
	term.AddOutputLines(lines)

	output := term.GetOutput()

	if len(output) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(output))
	}
}

func TestClear(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.AddOutput("line 1")
	term.AddOutput("line 2")
	term.Clear()

	output := term.GetOutput()

	if len(output) != 0 {
		t.Errorf("Expected output to be cleared, got %d lines", len(output))
	}

	if term.GetScrollOffset() != 0 {
		t.Error("Expected scroll offset to be reset")
	}
}

// --- Scrolling Tests ---

func TestScrollDown(t *testing.T) {
	term := New(DefaultTerminalConfig())

	for i := 0; i < 50; i++ {
		term.AddOutput("line")
	}

	// Scroll to top first, since AddOutput auto-scrolls to bottom
	term.ScrollToTop()
	initialOffset := term.GetScrollOffset()

	term.ScrollDown(5)

	if term.GetScrollOffset() != initialOffset+5 {
		t.Errorf("Expected scroll offset to be %d, got %d", initialOffset+5, term.GetScrollOffset())
	}
}

func TestScrollUp(t *testing.T) {
	term := New(DefaultTerminalConfig())

	for i := 0; i < 50; i++ {
		term.AddOutput("line")
	}

	term.ScrollDown(10)
	term.ScrollUp(5)

	if term.GetScrollOffset() < 0 {
		t.Error("Expected scroll offset to be non-negative")
	}
}

func TestScrollToTop(t *testing.T) {
	term := New(DefaultTerminalConfig())

	for i := 0; i < 50; i++ {
		term.AddOutput("line")
	}

	term.ScrollToTop()

	if term.GetScrollOffset() != 0 {
		t.Errorf("Expected scroll offset 0, got %d", term.GetScrollOffset())
	}
}

func TestScrollToBottom(t *testing.T) {
	term := New(DefaultTerminalConfig())

	for i := 0; i < 50; i++ {
		term.AddOutput("line")
	}

	term.ScrollToTop()
	term.ScrollToBottom()

	offset := term.GetScrollOffset()
	if offset < 0 {
		t.Error("Expected non-negative scroll offset")
	}
}

func TestGetVisibleOutput(t *testing.T) {
	config := DefaultTerminalConfig()
	config.Height = 10
	term := New(config)

	for i := 0; i < 50; i++ {
		term.AddOutput("line")
	}

	visible := term.GetVisibleOutput()

	// Should return at most height-1 lines (reserve 1 for input)
	if len(visible) > config.Height-1 {
		t.Errorf("Expected at most %d visible lines, got %d", config.Height-1, len(visible))
	}
}

// --- Command Parsing Tests ---

func TestParseCommandShell(t *testing.T) {
	cmd := ParseCommand("ls -la")

	if cmd.Type != CommandShell {
		t.Error("Expected CommandShell type")
	}

	if cmd.Input != "ls -la" {
		t.Errorf("Expected input 'ls -la', got '%s'", cmd.Input)
	}

	if len(cmd.Args) != 2 {
		t.Fatalf("Expected 2 args, got %d", len(cmd.Args))
	}

	if cmd.Args[0] != "ls" || cmd.Args[1] != "-la" {
		t.Error("Expected args ['ls', '-la']")
	}
}

func TestParseCommandBuiltin(t *testing.T) {
	cmd := ParseCommand(":help")

	if cmd.Type != CommandBuiltin {
		t.Error("Expected CommandBuiltin type")
	}

	if len(cmd.Args) != 1 || cmd.Args[0] != "help" {
		t.Error("Expected args ['help']")
	}
}

func TestParseCommandMnemosyne(t *testing.T) {
	cmd := ParseCommand("mnemosyne recall -q test")

	if cmd.Type != CommandMnemosyne {
		t.Error("Expected CommandMnemosyne type")
	}

	if !strings.Contains(cmd.Input, "mnemosyne") {
		t.Error("Expected input to contain 'mnemosyne'")
	}
}

func TestParseCommandEmpty(t *testing.T) {
	cmd := ParseCommand("")

	if cmd.Type != CommandShell {
		t.Error("Expected empty command to be CommandShell type")
	}

	if cmd.Input != "" {
		t.Error("Expected input to be empty")
	}
}

// --- History Tests ---

func TestAddHistory(t *testing.T) {
	term := New(DefaultTerminalConfig())

	cmd := Command{
		Type:  CommandShell,
		Input: "ls",
		Args:  []string{"ls"},
	}

	result := &CommandResult{
		Command:  cmd,
		Output:   "file1\nfile2",
		ExitCode: 0,
	}

	term.AddHistory(HistoryEntry{
		Command: cmd,
		Result:  result,
	})

	history := term.GetHistory()

	if len(history) != 1 {
		t.Fatalf("Expected 1 history entry, got %d", len(history))
	}

	if history[0].Command.Input != "ls" {
		t.Error("Expected history to contain 'ls' command")
	}
}

func TestHistoryPrevious(t *testing.T) {
	term := New(DefaultTerminalConfig())

	// Add some history
	for i := 1; i <= 3; i++ {
		cmd := Command{
			Type:  CommandShell,
			Input: "command" + string(rune('0'+i)),
		}
		term.AddHistory(HistoryEntry{Command: cmd})
	}

	// Navigate to previous
	if !term.HistoryPrevious() {
		t.Error("Expected HistoryPrevious to succeed")
	}

	// Should get the last command
	if term.GetInput() != "command3" {
		t.Errorf("Expected 'command3', got '%s'", term.GetInput())
	}

	// Navigate again
	term.HistoryPrevious()

	if term.GetInput() != "command2" {
		t.Errorf("Expected 'command2', got '%s'", term.GetInput())
	}
}

func TestHistoryNext(t *testing.T) {
	term := New(DefaultTerminalConfig())

	// Add history
	for i := 1; i <= 3; i++ {
		cmd := Command{
			Type:  CommandShell,
			Input: "command" + string(rune('0'+i)),
		}
		term.AddHistory(HistoryEntry{Command: cmd})
	}

	// Navigate to previous twice
	term.HistoryPrevious()
	term.HistoryPrevious()

	// Navigate forward
	if !term.HistoryNext() {
		t.Error("Expected HistoryNext to succeed")
	}

	if term.GetInput() != "command3" {
		t.Errorf("Expected 'command3', got '%s'", term.GetInput())
	}
}

func TestHistoryNavigationEmpty(t *testing.T) {
	term := New(DefaultTerminalConfig())

	// Try to navigate with no history
	if term.HistoryPrevious() {
		t.Error("Expected HistoryPrevious to fail with empty history")
	}

	if term.HistoryNext() {
		t.Error("Expected HistoryNext to fail with empty history")
	}
}

// --- Command Execution Tests ---

func TestExecuteBuiltinHelp(t *testing.T) {
	term := New(DefaultTerminalConfig())

	cmd := ParseCommand(":help")
	result := term.Execute(cmd)

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if result.Output == "" {
		t.Error("Expected help output to be non-empty")
	}

	if !strings.Contains(result.Output, "Built-in Commands") {
		t.Error("Expected help output to contain 'Built-in Commands'")
	}
}

func TestExecuteBuiltinClear(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.AddOutput("line 1")
	term.AddOutput("line 2")

	cmd := ParseCommand(":clear")
	result := term.Execute(cmd)

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if len(term.GetOutput()) != 0 {
		t.Error("Expected output to be cleared")
	}
}

func TestExecuteBuiltinHistory(t *testing.T) {
	term := New(DefaultTerminalConfig())

	// Add some history
	cmd1 := Command{Type: CommandShell, Input: "ls"}
	term.AddHistory(HistoryEntry{Command: cmd1})

	cmd := ParseCommand(":history")
	result := term.Execute(cmd)

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Output, "ls") {
		t.Error("Expected history output to contain 'ls'")
	}
}

func TestExecuteBuiltinUnknown(t *testing.T) {
	term := New(DefaultTerminalConfig())

	cmd := ParseCommand(":unknown")
	result := term.Execute(cmd)

	if result.ExitCode == 0 {
		t.Error("Expected non-zero exit code for unknown command")
	}

	if result.Error == "" {
		t.Error("Expected error message for unknown command")
	}
}

func TestExecuteShellCommand(t *testing.T) {
	term := New(DefaultTerminalConfig())

	// Simple echo command (cross-platform)
	cmd := ParseCommand("echo test")
	result := term.Execute(cmd)

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if !strings.Contains(result.Output, "test") {
		t.Errorf("Expected output to contain 'test', got '%s'", result.Output)
	}
}

func TestExecuteAndDisplay(t *testing.T) {
	term := New(DefaultTerminalConfig())

	initialOutputLen := len(term.GetOutput())

	term.ExecuteAndDisplay("echo hello")

	// Should have added output
	if len(term.GetOutput()) <= initialOutputLen {
		t.Error("Expected output to be added")
	}

	// Should have cleared input
	if term.GetInput() != "" {
		t.Error("Expected input to be cleared after execution")
	}
}

func TestSubmit(t *testing.T) {
	term := New(DefaultTerminalConfig())

	term.SetInput("echo test")
	term.Submit()

	// Should have cleared input
	if term.GetInput() != "" {
		t.Error("Expected input to be cleared after submit")
	}

	// Should have output
	if len(term.GetOutput()) == 0 {
		t.Error("Expected output to be added after submit")
	}
}

// --- Configuration Tests ---

func TestGetSetConfig(t *testing.T) {
	term := New(DefaultTerminalConfig())

	newConfig := DefaultTerminalConfig()
	newConfig.Width = 100

	term.SetConfig(newConfig)

	retrieved := term.GetConfig()
	if retrieved.Width != 100 {
		t.Errorf("Expected width 100, got %d", retrieved.Width)
	}
}

// --- Execution State Tests ---

func TestIsExecuting(t *testing.T) {
	term := New(DefaultTerminalConfig())

	if term.IsExecuting() {
		t.Error("Expected IsExecuting to be false initially")
	}

	term.SetExecuting(true)

	if !term.IsExecuting() {
		t.Error("Expected IsExecuting to be true after setting")
	}
}

func TestGetLastResult(t *testing.T) {
	term := New(DefaultTerminalConfig())

	if term.GetLastResult() != nil {
		t.Error("Expected initial last result to be nil")
	}

	result := &CommandResult{
		ExitCode: 0,
	}

	term.SetLastResult(result)

	retrieved := term.GetLastResult()
	if retrieved == nil {
		t.Fatal("Expected last result to be set")
	}

	if retrieved.ExitCode != 0 {
		t.Error("Expected exit code to match")
	}
}

// --- Command Result Tests ---

func TestCommandTimestamp(t *testing.T) {
	cmd := ParseCommand("test")

	if cmd.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	// Should be recent
	if time.Since(cmd.Timestamp) > time.Second {
		t.Error("Expected timestamp to be recent")
	}
}
