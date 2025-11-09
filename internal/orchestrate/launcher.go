package orchestrate

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Launcher manages the lifecycle of a mnemosyne orchestrate subprocess.
type Launcher struct {
	cmd     *exec.Cmd
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	events  chan *AgentEvent
	errors  chan error
	done    chan bool
	running bool
	mu      sync.Mutex
	cancel  chan struct{}
}

// LaunchOptions configures the behavior of the launcher and orchestration.
type LaunchOptions struct {
	DatabasePath    string // Path to mnemosyne database
	PollingInterval int    // milliseconds between polls
	MaxConcurrent   int    // max concurrent agents
	EnableDashboard bool   // enable dashboard output
}

// NewLauncher creates a new process launcher with buffered channels.
func NewLauncher() *Launcher {
	return &Launcher{
		events: make(chan *AgentEvent, 100),
		errors: make(chan error, 10),
		done:   make(chan bool, 1),
		cancel: make(chan struct{}),
	}
}

// Start spawns the mnemosyne orchestrate subprocess with the given plan and options.
func (l *Launcher) Start(plan *WorkPlan, opts LaunchOptions) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.running {
		return fmt.Errorf("launcher is already running")
	}

	if plan == nil {
		return fmt.Errorf("work plan cannot be nil")
	}

	if err := plan.Validate(); err != nil {
		return fmt.Errorf("invalid work plan: %w", err)
	}

	// Serialize plan to JSON
	planJSON, err := plan.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize work plan: %w", err)
	}

	// Build command arguments
	args := []string{"orchestrate", "--plan", string(planJSON)}
	if opts.DatabasePath != "" {
		args = append(args, "--database", opts.DatabasePath)
	}
	if opts.PollingInterval > 0 {
		args = append(args, "--polling-interval", strconv.Itoa(opts.PollingInterval))
	}
	if opts.MaxConcurrent > 0 {
		args = append(args, "--max-concurrent", strconv.Itoa(opts.MaxConcurrent))
	}
	if opts.EnableDashboard {
		args = append(args, "--dashboard")
	}

	// Create command
	l.cmd = exec.Command("mnemosyne", args...)

	// Set up pipes
	var err1, err2 error
	l.stdout, err1 = l.cmd.StdoutPipe()
	l.stderr, err2 = l.cmd.StderrPipe()
	if err1 != nil || err2 != nil {
		return fmt.Errorf("failed to create pipes: stdout=%w stderr=%w", err1, err2)
	}

	// Start the command
	if err := l.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mnemosyne: %w", err)
	}

	l.running = true

	// Start goroutines to read stdout and stderr
	go l.streamOutput(l.stdout, "stdout")
	go l.streamOutput(l.stderr, "stderr")

	// Wait for process to complete in background
	go l.waitForCompletion()

	return nil
}

// Stop gracefully shuts down the orchestration process.
// Sends SIGTERM, waits 5 seconds, then sends SIGKILL if needed.
func (l *Launcher) Stop() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.running || l.cmd == nil || l.cmd.Process == nil {
		return nil
	}

	// Send SIGTERM for graceful shutdown
	if err := l.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		// If process is already gone, that's fine
		if err.Error() != "os: process already finished" {
			return fmt.Errorf("failed to send SIGTERM: %w", err)
		}
	}

	// Wait up to 5 seconds for graceful shutdown
	done := make(chan error, 1)
	go func() {
		done <- l.cmd.Wait()
	}()

	select {
	case <-done:
		// Process exited cleanly
		l.running = false
		l.cleanup()
		return nil
	case <-time.After(5 * time.Second):
		// Force kill if not stopped
		if err := l.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
		// Wait for the kill to complete
		l.cmd.Wait()
		l.running = false
		l.cleanup()
		return nil
	}
}

// Restart stops the current process and starts a new one with the given plan.
func (l *Launcher) Restart(plan *WorkPlan, opts LaunchOptions) error {
	if err := l.Stop(); err != nil {
		return fmt.Errorf("failed to stop launcher: %w", err)
	}

	// Small delay to ensure process is fully cleaned up
	time.Sleep(100 * time.Millisecond)

	return l.Start(plan, opts)
}

// IsRunning returns whether the launcher is currently running.
func (l *Launcher) IsRunning() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.running
}

// Events returns a read-only channel of agent events.
func (l *Launcher) Events() <-chan *AgentEvent {
	return l.events
}

// Errors returns a read-only channel of errors.
func (l *Launcher) Errors() <-chan error {
	return l.errors
}

// Wait blocks until the process completes, returning any error.
func (l *Launcher) Wait() error {
	if l.cmd == nil {
		return fmt.Errorf("launcher not started")
	}
	return l.cmd.Wait()
}

// streamOutput reads from a pipe and sends parsed events to the events channel.
func (l *Launcher) streamOutput(pipe io.ReadCloser, source string) {
	defer pipe.Close()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		select {
		case <-l.cancel:
			return
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		event, err := ParseFromLine(line)
		if err != nil {
			// If we can't parse as JSON, create a generic log event
			event = &AgentEvent{
				Timestamp: time.Now(),
				Agent:     AgentOrchestrator,
				EventType: EventLog,
				Message:   line,
			}
		}

		// Send event with non-blocking semantics
		// If channel is full, drop oldest events
		select {
		case l.events <- event:
		default:
			// Channel is full, drop event
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		select {
		case l.errors <- fmt.Errorf("scanner error on %s: %w", source, err):
		default:
		}
	}
}

// waitForCompletion waits for the process to finish and closes channels.
func (l *Launcher) waitForCompletion() {
	if l.cmd == nil {
		return
	}

	l.cmd.Wait()

	l.mu.Lock()
	l.running = false
	l.mu.Unlock()

	l.cleanup()
}

// cleanup closes all channels and pipes.
func (l *Launcher) cleanup() {
	close(l.cancel)

	if l.stdout != nil {
		l.stdout.Close()
	}
	if l.stderr != nil {
		l.stderr.Close()
	}

	// Signal completion
	select {
	case l.done <- true:
	default:
	}
}
