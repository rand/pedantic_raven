package terminal

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Execute executes a command and returns the result.
func (t *Terminal) Execute(cmd Command) *CommandResult {
	startTime := time.Now()
	result := &CommandResult{
		Command:  cmd,
		ExitCode: 0,
	}

	// Set executing state
	t.SetExecuting(true)
	defer t.SetExecuting(false)

	switch cmd.Type {
	case CommandBuiltin:
		t.executeBuiltin(cmd, result)
	case CommandMnemosyne:
		t.executeMnemosyne(cmd, result)
	case CommandShell:
		t.executeShell(cmd, result)
	}

	result.Duration = time.Since(startTime)
	t.SetLastResult(result)

	// Add to history
	t.AddHistory(HistoryEntry{
		Command: cmd,
		Result:  result,
	})

	return result
}

// executeBuiltin executes a built-in command.
func (t *Terminal) executeBuiltin(cmd Command, result *CommandResult) {
	if len(cmd.Args) == 0 {
		result.Error = "No built-in command specified"
		result.ExitCode = 1
		return
	}

	builtin := cmd.Args[0]

	switch builtin {
	case "clear":
		t.Clear()
		result.Output = "Terminal cleared"

	case "help":
		result.Output = t.getHelp()

	case "history":
		result.Output = t.getHistoryOutput()

	case "exit", "quit":
		result.Output = "Exit command received"
		result.ExitCode = 0

	default:
		result.Error = fmt.Sprintf("Unknown built-in command: %s", builtin)
		result.ExitCode = 1
	}
}

// executeMnemosyne executes a mnemosyne CLI command.
func (t *Terminal) executeMnemosyne(cmd Command, result *CommandResult) {
	// Execute as shell command
	t.executeShell(cmd, result)
}

// executeShell executes a shell command.
func (t *Terminal) executeShell(cmd Command, result *CommandResult) {
	if len(cmd.Args) == 0 {
		result.ExitCode = 0
		return
	}

	// Create command
	shellCmd := exec.Command(cmd.Args[0], cmd.Args[1:]...)
	shellCmd.Dir = t.config.WorkingDir

	// Capture output
	var stdout, stderr bytes.Buffer
	shellCmd.Stdout = &stdout
	shellCmd.Stderr = &stderr

	// Execute
	err := shellCmd.Run()

	result.Output = stdout.String()
	result.Error = stderr.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
			if result.Error == "" {
				result.Error = err.Error()
			}
		}
	}
}

// getHelp returns help text for built-in commands.
func (t *Terminal) getHelp() string {
	return `Built-in Commands:
  :clear     - Clear terminal output
  :help      - Show this help message
  :history   - Show command history
  :exit      - Exit terminal (or quit)

Shell Commands:
  Any command not starting with ':' will be executed as a shell command.
  Examples: ls, git status, echo "hello"

mnemosyne Commands:
  Commands starting with 'mnemosyne' are passed to the mnemosyne CLI.
  Examples:
    mnemosyne recall -q "search query"
    mnemosyne remember -c "content to remember"

Navigation:
  Up/Down    - Navigate command history
  Ctrl+A     - Move cursor to start
  Ctrl+E     - Move cursor to end
  Ctrl+L     - Clear screen
`
}

// getHistoryOutput returns formatted command history.
func (t *Terminal) getHistoryOutput() string {
	if len(t.history) == 0 {
		return "No command history"
	}

	var lines []string
	for i, entry := range t.history {
		status := "✓"
		if entry.Result != nil && entry.Result.ExitCode != 0 {
			status = "✗"
		}

		line := fmt.Sprintf("%3d %s %s", i+1, status, entry.Command.Input)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// ExecuteAndDisplay executes a command and displays the result.
func (t *Terminal) ExecuteAndDisplay(input string) {
	// Parse command
	cmd := ParseCommand(input)

	// Add command echo to output
	t.AddOutput(t.config.Prompt + input)

	// Execute
	result := t.Execute(cmd)

	// Display output
	if result.Output != "" {
		lines := strings.Split(strings.TrimRight(result.Output, "\n"), "\n")
		t.AddOutputLines(lines)
	}

	// Display errors
	if result.Error != "" {
		lines := strings.Split(strings.TrimRight(result.Error, "\n"), "\n")
		for _, line := range lines {
			t.AddOutput("ERROR: " + line)
		}
	}

	// Display exit code if non-zero
	if result.ExitCode != 0 {
		t.AddOutput(fmt.Sprintf("Exit code: %d", result.ExitCode))
	}

	// Clear input
	t.SetInput("")
	t.ResetHistoryNavigation()
}

// Submit submits the current input for execution.
func (t *Terminal) Submit() {
	input := t.GetInput()
	if input == "" {
		return
	}

	t.ExecuteAndDisplay(input)
}
