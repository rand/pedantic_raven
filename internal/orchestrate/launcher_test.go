package orchestrate

import (
	"os"
	"testing"
	"time"
)

// Helper: Create a mock mnemosyne script for testing
func createMockMnemosyne(t *testing.T) string {
	// Create a temporary script that mimics mnemosyne orchestrate output
	script := `#!/bin/bash
# Mock mnemosyne orchestrate script

# Emit a started event as JSON
echo '{"timestamp":"2025-11-09T12:00:00Z","agent":0,"eventType":0,"taskId":"task1","message":"Starting orchestration"}'

# Emit a progress event
sleep 0.1
echo '{"timestamp":"2025-11-09T12:00:01Z","agent":3,"eventType":1,"taskId":"task1","message":"Executing task1"}'

# Emit a plaintext log line (fallback format)
sleep 0.1
echo "Some log output from task1"

# Emit a completion event
sleep 0.1
echo '{"timestamp":"2025-11-09T12:00:02Z","agent":3,"eventType":2,"taskId":"task1","message":"Completed task1"}'

# Emit second task events
sleep 0.1
echo '{"timestamp":"2025-11-09T12:00:03Z","agent":3,"eventType":0,"taskId":"task2","message":"Starting task2"}'

sleep 0.1
echo '{"timestamp":"2025-11-09T12:00:04Z","agent":3,"eventType":2,"taskId":"task2","message":"Completed task2"}'

exit 0
`

	tmpfile, err := os.CreateTemp("", "mock-mnemosyne-*.sh")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpfile.WriteString(script); err != nil {
		tmpfile.Close()
		t.Fatalf("failed to write mock script: %v", err)
	}
	tmpfile.Close()

	// Make script executable
	if err := os.Chmod(tmpfile.Name(), 0755); err != nil {
		t.Fatalf("failed to chmod mock script: %v", err)
	}

	return tmpfile.Name()
}

// --- Launcher Lifecycle Tests ---

// TestLauncherNewCreatesInstance tests that NewLauncher creates a launcher instance.
func TestLauncherNewCreatesInstance(t *testing.T) {
	launcher := NewLauncher()

	if launcher == nil {
		t.Fatal("expected NewLauncher to return non-nil launcher")
	}

	if launcher.IsRunning() {
		t.Fatal("expected new launcher to not be running")
	}
}

// TestLauncherIsRunningInitiallyFalse tests that IsRunning() is false initially.
func TestLauncherIsRunningInitiallyFalse(t *testing.T) {
	launcher := NewLauncher()

	if launcher.IsRunning() {
		t.Fatal("expected IsRunning() to return false for new launcher")
	}
}

// TestLauncherEventsReturnsChannel tests that Events() returns a channel.
func TestLauncherEventsReturnsChannel(t *testing.T) {
	launcher := NewLauncher()

	ch := launcher.Events()
	if ch == nil {
		t.Fatal("expected Events() to return non-nil channel")
	}
}

// TestLauncherErrorsReturnsChannel tests that Errors() returns a channel.
func TestLauncherErrorsReturnsChannel(t *testing.T) {
	launcher := NewLauncher()

	ch := launcher.Errors()
	if ch == nil {
		t.Fatal("expected Errors() to return non-nil channel")
	}
}

// --- Validation Tests ---

// TestLauncherStartWithNilPlan tests that Start with nil plan returns error.
func TestLauncherStartWithNilPlan(t *testing.T) {
	launcher := NewLauncher()
	opts := LaunchOptions{}

	err := launcher.Start(nil, opts)
	if err == nil {
		t.Fatal("expected error for nil plan")
	}

	if launcher.IsRunning() {
		t.Fatal("expected launcher to not be running after failed start")
	}
}

// TestLauncherStartWithEmptyPlan tests that Start with empty plan returns error.
func TestLauncherStartWithEmptyPlan(t *testing.T) {
	launcher := NewLauncher()

	// Plan with no tasks
	invalidPlan := &WorkPlan{
		Name:          "Empty Plan",
		Description:   "A plan with no tasks",
		MaxConcurrent: 1,
		Tasks:         []Task{},
	}

	opts := LaunchOptions{DatabasePath: "/tmp/test.db"}

	err := launcher.Start(invalidPlan, opts)
	if err == nil {
		t.Fatal("expected error for plan with no tasks")
	}

	if launcher.IsRunning() {
		t.Fatal("expected launcher to not be running after failed start")
	}
}

// TestLauncherStopWhenNotRunning tests that Stop is safe when not running.
func TestLauncherStopWhenNotRunning(t *testing.T) {
	launcher := NewLauncher()

	// Try to stop without starting
	err := launcher.Stop()
	if err != nil {
		t.Fatalf("expected no error when stopping non-running launcher, got: %v", err)
	}

	if launcher.IsRunning() {
		t.Fatal("expected launcher to still not be running")
	}
}

// TestLauncherWaitWhenNotStarted tests that Wait returns error when not started.
func TestLauncherWaitWhenNotStarted(t *testing.T) {
	launcher := NewLauncher()

	err := launcher.Wait()
	if err == nil {
		t.Fatal("expected error when waiting on non-started launcher")
	}
}

// TestLauncherMultipleStopsAreSafe tests that multiple Stop calls are safe.
func TestLauncherMultipleStopsAreSafe(t *testing.T) {
	launcher := NewLauncher()

	// First stop when not running
	err1 := launcher.Stop()
	if err1 != nil {
		t.Fatalf("first stop failed: %v", err1)
	}

	// Second stop when not running
	err2 := launcher.Stop()
	if err2 != nil {
		t.Fatalf("second stop failed: %v", err2)
	}

	// Third stop when not running
	err3 := launcher.Stop()
	if err3 != nil {
		t.Fatalf("third stop failed: %v", err3)
	}
}

// Note: ParseFromLine tests are in types_test.go (they test the types package functionality)

// --- AgentEvent Tests ---

// TestAgentEventFormatContainsKey tests that Format() includes key information.
func TestAgentEventFormatContainsKey(t *testing.T) {
	event := &AgentEvent{
		Timestamp: time.Date(2025, 11, 9, 12, 0, 0, 0, time.UTC),
		Agent:     AgentExecutor,
		EventType: EventCompleted,
		TaskID:    "task1",
		Message:   "Task completed successfully",
	}

	formatted := event.Format()
	if formatted == "" {
		t.Fatal("expected non-empty formatted string")
	}

	// Verify key information is present
	if !stringContains(formatted, "task1") {
		t.Fatalf("expected format to contain task ID, got: %s", formatted)
	}

	if !stringContains(formatted, "Executor") {
		t.Fatalf("expected format to contain agent name, got: %s", formatted)
	}

	if !stringContains(formatted, "completed") {
		t.Fatalf("expected format to contain event type, got: %s", formatted)
	}
}

// TestAgentEventJSONRoundTrip tests JSON serialization/deserialization.
func TestAgentEventJSONRoundTrip(t *testing.T) {
	event := &AgentEvent{
		Timestamp: time.Date(2025, 11, 9, 12, 0, 0, 0, time.UTC),
		Agent:     AgentReviewer,
		EventType: EventProgress,
		TaskID:    "task2",
		Message:   "In progress",
		Metadata: map[string]interface{}{
			"progress": 50,
		},
	}

	// Serialize
	jsonBytes, err := event.ToJSON()
	if err != nil {
		t.Fatalf("failed to serialize: %v", err)
	}

	// Deserialize
	restored := &AgentEvent{}
	if err := restored.FromJSON(jsonBytes); err != nil {
		t.Fatalf("failed to deserialize: %v", err)
	}

	// Verify round-trip
	if restored.TaskID != event.TaskID {
		t.Fatalf("taskId mismatch: expected %s, got %s", event.TaskID, restored.TaskID)
	}

	if restored.Message != event.Message {
		t.Fatalf("message mismatch: expected %s, got %s", event.Message, restored.Message)
	}

	if restored.Agent != event.Agent {
		t.Fatalf("agent mismatch: expected %s, got %s", event.Agent.String(), restored.Agent.String())
	}
}

// TestAgentEventFromJSONInvalid tests that invalid JSON is rejected.
func TestAgentEventFromJSONInvalid(t *testing.T) {
	event := &AgentEvent{}

	err := event.FromJSON([]byte("not valid json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// TestAgentEventToJSONAndFromJSON tests full JSON cycle.
func TestAgentEventToJSONAndFromJSON(t *testing.T) {
	original := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentOrchestrator,
		EventType: EventStarted,
		TaskID:    "workflow-1",
		Message:   "Starting orchestration",
	}

	data, err := original.ToJSON()
	if err != nil {
		t.Fatalf("failed to convert to JSON: %v", err)
	}

	restored := &AgentEvent{}
	err = restored.FromJSON(data)
	if err != nil {
		t.Fatalf("failed to convert from JSON: %v", err)
	}

	if restored.TaskID != original.TaskID {
		t.Fatalf("task ID mismatch after JSON round-trip")
	}
}

// --- Concurrency Tests ---

// TestLauncherConcurrentIsRunning tests IsRunning is safe for concurrent access.
func TestLauncherConcurrentIsRunning(t *testing.T) {
	launcher := NewLauncher()
	done := make(chan bool, 5)

	// Start 5 goroutines calling IsRunning concurrently
	for i := 0; i < 5; i++ {
		go func() {
			_ = launcher.IsRunning()
			done <- true
		}()
	}

	// All should complete without panic
	for i := 0; i < 5; i++ {
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for goroutines")
		}
	}
}

// TestLauncherConcurrentChannelAccess tests multiple consumers of event channel.
func TestLauncherConcurrentChannelAccess(t *testing.T) {
	launcher := NewLauncher()

	// Both calls should return the same underlying channel
	ch1 := launcher.Events()
	ch2 := launcher.Events()

	if ch1 != ch2 {
		t.Fatal("expected Events() to return the same channel")
	}
}

// --- Benchmark Tests ---

// BenchmarkParseFromLineJSON benchmarks JSON parsing performance.
func BenchmarkParseFromLineJSON(b *testing.B) {
	line := `{"timestamp":"2025-11-09T12:00:00Z","agent":0,"eventType":0,"taskId":"task1","message":"Test event"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseFromLine(line)
	}
}

// BenchmarkAgentEventFormat benchmarks the Format method.
func BenchmarkAgentEventFormat(b *testing.B) {
	event := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventCompleted,
		TaskID:    "task1",
		Message:   "Task completed",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = event.Format()
	}
}

// BenchmarkNewLauncher benchmarks launcher creation.
func BenchmarkNewLauncher(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewLauncher()
	}
}

// BenchmarkLauncherIsRunning benchmarks IsRunning method.
func BenchmarkLauncherIsRunning(b *testing.B) {
	launcher := NewLauncher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = launcher.IsRunning()
	}
}

// --- Helper Functions ---

// stringContains checks if a string contains a substring (case-sensitive).
func stringContains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
