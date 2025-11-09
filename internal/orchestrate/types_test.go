package orchestrate

import (
	"encoding/json"
	"testing"
	"time"
)

// --- WorkPlan Tests ---

// TestWorkPlanValidation tests that a valid work plan passes validation.
func TestWorkPlanValidation(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Sample Plan",
		Description:   "A sample work plan",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "First task",
				Type:        TaskTypeParallel,
				Priority:    5,
			},
			{
				ID:           "task2",
				Description:  "Second task",
				Type:         TaskTypeParallel,
				Dependencies: []string{"task1"},
				Priority:     3,
			},
		},
	}

	if err := plan.Validate(); err != nil {
		t.Fatalf("expected valid plan to pass validation, got error: %v", err)
	}
}

// TestWorkPlanCircularDependency tests that circular dependencies are detected.
func TestWorkPlanCircularDependency(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Circular Plan",
		Description:   "A plan with circular dependencies",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:           "task1",
				Description:  "Task 1",
				Type:         TaskTypeParallel,
				Dependencies: []string{"task2"},
			},
			{
				ID:           "task2",
				Description:  "Task 2",
				Type:         TaskTypeParallel,
				Dependencies: []string{"task1"},
			},
		},
	}

	err := plan.Validate()
	if err == nil {
		t.Fatal("expected circular dependency to be detected")
	}
	if err.Error() != "circular dependency detected: task2 -> task1" &&
		err.Error() != "circular dependency detected: task1 -> task2" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestWorkPlanMissingDependency tests that missing dependencies are detected.
func TestWorkPlanMissingDependency(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Missing Dep Plan",
		Description:   "A plan with missing dependencies",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:           "task1",
				Description:  "Task 1",
				Type:         TaskTypeParallel,
				Dependencies: []string{"nonexistent"},
			},
		},
	}

	err := plan.Validate()
	if err == nil {
		t.Fatal("expected missing dependency error")
	}
	if err.Error() != "task task1 depends on non-existent task nonexistent" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestWorkPlanJSONMarshaling tests JSON serialization and deserialization.
func TestWorkPlanJSONMarshaling(t *testing.T) {
	original := &WorkPlan{
		Name:          "JSON Test Plan",
		Description:   "A plan for JSON testing",
		MaxConcurrent: 3,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "First task",
				Type:        TaskTypeSequential,
				Agent:       AgentExecutor,
				Priority:    7,
			},
			{
				ID:           "task2",
				Description:  "Second task",
				Type:         TaskTypeParallel,
				Dependencies: []string{"task1"},
				Priority:     5,
			},
		},
	}

	// Marshal to JSON
	jsonBytes, err := original.ToJSON()
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	// Unmarshal from JSON
	restored := &WorkPlan{}
	if err := restored.FromJSON(jsonBytes); err != nil {
		t.Fatalf("failed to unmarshal from JSON: %v", err)
	}

	// Verify round-trip
	if restored.Name != original.Name {
		t.Fatalf("name mismatch: expected %s, got %s", original.Name, restored.Name)
	}
	if restored.MaxConcurrent != original.MaxConcurrent {
		t.Fatalf("maxConcurrent mismatch: expected %d, got %d", original.MaxConcurrent, restored.MaxConcurrent)
	}
	if len(restored.Tasks) != len(original.Tasks) {
		t.Fatalf("task count mismatch: expected %d, got %d", len(original.Tasks), len(restored.Tasks))
	}
}

// TestWorkPlanToDependencyGraph tests conversion to dependency graph representation.
func TestWorkPlanToDependencyGraph(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Graph Test Plan",
		Description:   "A plan for graph testing",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "Task 1",
				Type:        TaskTypeParallel,
			},
			{
				ID:           "task2",
				Description:  "Task 2",
				Type:         TaskTypeParallel,
				Dependencies: []string{"task1"},
			},
			{
				ID:           "task3",
				Description:  "Task 3",
				Type:         TaskTypeParallel,
				Dependencies: []string{"task1"},
			},
		},
	}

	graph, err := plan.ToDependencyGraph()
	if err != nil {
		t.Fatalf("failed to generate dependency graph: %v", err)
	}

	// task1 should have task2 and task3 depending on it
	if len(graph["task1"]) != 2 {
		t.Fatalf("expected task1 to have 2 dependents, got %d", len(graph["task1"]))
	}

	// task2 should have no dependents
	if len(graph["task2"]) != 0 {
		t.Fatalf("expected task2 to have 0 dependents, got %d", len(graph["task2"]))
	}
}

// --- Task Tests ---

// TestTaskTypeString tests TaskType string conversion.
func TestTaskTypeString(t *testing.T) {
	tests := []struct {
		taskType TaskType
		expected string
	}{
		{TaskTypeParallel, "parallel"},
		{TaskTypeSequential, "sequential"},
		{TaskTypeBlocking, "blocking"},
		{TaskType(999), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.taskType.String(); got != tt.expected {
			t.Fatalf("TaskType.String() = %s, expected %s", got, tt.expected)
		}
	}
}

// TestTaskStatusString tests TaskStatus string conversion.
func TestTaskStatusString(t *testing.T) {
	tests := []struct {
		status   TaskStatus
		expected string
	}{
		{TaskStatusPending, "pending"},
		{TaskStatusActive, "active"},
		{TaskStatusCompleted, "completed"},
		{TaskStatusFailed, "failed"},
		{TaskStatus(999), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.status.String(); got != tt.expected {
			t.Fatalf("TaskStatus.String() = %s, expected %s", got, tt.expected)
		}
	}
}

// TestTaskValidation tests task validation.
func TestTaskValidation(t *testing.T) {
	tests := []struct {
		name        string
		task        *Task
		shouldError bool
		errMsg      string
	}{
		{
			name: "valid task",
			task: &Task{
				ID:          "task1",
				Description: "Valid task",
				Type:        TaskTypeParallel,
				Priority:    5,
			},
			shouldError: false,
		},
		{
			name: "empty ID",
			task: &Task{
				ID:          "",
				Description: "Task with no ID",
				Type:        TaskTypeParallel,
			},
			shouldError: true,
			errMsg:      "task ID cannot be empty",
		},
		{
			name: "empty description",
			task: &Task{
				ID:          "task1",
				Description: "",
				Type:        TaskTypeParallel,
			},
			shouldError: true,
			errMsg:      "task description cannot be empty",
		},
		{
			name: "invalid priority",
			task: &Task{
				ID:          "task1",
				Description: "Task with high priority",
				Type:        TaskTypeParallel,
				Priority:    15,
			},
			shouldError: true,
			errMsg:      "task priority must be between 0 and 10, got 15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.shouldError {
				if err == nil {
					t.Fatalf("expected error: %s", tt.errMsg)
				}
				if err.Error() != tt.errMsg {
					t.Fatalf("error message mismatch: expected %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// --- AgentEvent Tests ---

// TestAgentEventParseFromLine tests parsing events from log lines.
func TestAgentEventParseFromLine(t *testing.T) {
	// Create a valid event and marshal it to JSON
	event := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventCompleted,
		TaskID:    "task1",
		Message:   "Task completed successfully",
		Metadata:  map[string]interface{}{"duration": 5.5},
	}

	jsonLine, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	parsed, err := ParseFromLine(string(jsonLine))
	if err != nil {
		t.Fatalf("failed to parse event: %v", err)
	}

	if parsed.TaskID != event.TaskID {
		t.Fatalf("TaskID mismatch: expected %s, got %s", event.TaskID, parsed.TaskID)
	}
	if parsed.Message != event.Message {
		t.Fatalf("Message mismatch: expected %s, got %s", event.Message, parsed.Message)
	}
}

// TestAgentEventFormat tests human-readable formatting.
func TestAgentEventFormat(t *testing.T) {
	now := time.Now()
	event := &AgentEvent{
		Timestamp: now,
		Agent:     AgentOrchestrator,
		EventType: EventStarted,
		TaskID:    "task1",
		Message:   "Starting task",
	}

	formatted := event.Format()
	if formatted == "" {
		t.Fatal("expected non-empty formatted string")
	}

	// Check that formatted string contains key components
	if !contains(formatted, "Orchestrator") {
		t.Fatalf("formatted string missing agent name: %s", formatted)
	}
	if !contains(formatted, "started") {
		t.Fatalf("formatted string missing event type: %s", formatted)
	}
	if !contains(formatted, "task1") {
		t.Fatalf("formatted string missing task ID: %s", formatted)
	}
}

// TestAgentEventJSONMarshaling tests JSON serialization round-trip.
func TestAgentEventJSONMarshaling(t *testing.T) {
	original := &AgentEvent{
		Timestamp: time.Now().Round(time.Second), // Round to avoid precision issues
		Agent:     AgentReviewer,
		EventType: EventProgress,
		TaskID:    "task2",
		Message:   "50% complete",
		Metadata:  map[string]interface{}{"progress": 0.5, "eta": 10},
	}

	// Marshal to JSON
	jsonBytes, err := original.ToJSON()
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	// Unmarshal from JSON
	restored := &AgentEvent{}
	if err := restored.FromJSON(jsonBytes); err != nil {
		t.Fatalf("failed to unmarshal from JSON: %v", err)
	}

	// Verify round-trip
	if restored.Agent != original.Agent {
		t.Fatalf("agent mismatch: expected %v, got %v", original.Agent, restored.Agent)
	}
	if restored.TaskID != original.TaskID {
		t.Fatalf("task ID mismatch: expected %s, got %s", original.TaskID, restored.TaskID)
	}
}

// TestEventTypeString tests EventType string conversion.
func TestEventTypeString(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  string
	}{
		{EventStarted, "started"},
		{EventProgress, "progress"},
		{EventCompleted, "completed"},
		{EventFailed, "failed"},
		{EventHandoff, "handoff"},
		{EventLog, "log"},
		{EventType(999), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.eventType.String(); got != tt.expected {
			t.Fatalf("EventType.String() = %s, expected %s", got, tt.expected)
		}
	}
}

// --- AgentStatus Tests ---

// TestAgentStatusIsActive tests active status detection.
func TestAgentStatusIsActive(t *testing.T) {
	tests := []struct {
		status       string
		expectActive bool
	}{
		{"active", true},
		{"idle", false},
		{"error", false},
	}

	for _, tt := range tests {
		agentStatus := &AgentStatus{
			Agent:      AgentExecutor,
			Status:     tt.status,
			LastUpdate: time.Now(),
		}

		if agentStatus.IsActive() != tt.expectActive {
			t.Fatalf("IsActive() for status %q: expected %v, got %v", tt.status, tt.expectActive, agentStatus.IsActive())
		}
	}
}

// TestAgentStatusCanExecute tests task assignment logic.
func TestAgentStatusCanExecute(t *testing.T) {
	executor := &AgentStatus{
		Agent:      AgentExecutor,
		Status:     "idle",
		LastUpdate: time.Now(),
	}

	tests := []struct {
		name    string
		status  string
		task    *Task
		canExec bool
	}{
		{
			name:   "idle agent, any task",
			status: "idle",
			task: &Task{
				ID:          "task1",
				Description: "Task",
				Type:        TaskTypeParallel,
			},
			canExec: true,
		},
		{
			name:   "error agent, cannot execute",
			status: "error",
			task: &Task{
				ID:          "task1",
				Description: "Task",
				Type:        TaskTypeParallel,
			},
			canExec: false,
		},
		{
			name:   "active agent, any task",
			status: "active",
			task: &Task{
				ID:          "task1",
				Description: "Task",
				Type:        TaskTypeParallel,
			},
			canExec: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor.Status = tt.status
			if executor.CanExecute(tt.task) != tt.canExec {
				t.Fatalf("CanExecute() for status %q: expected %v, got %v", tt.status, tt.canExec, executor.CanExecute(tt.task))
			}
		})
	}
}

// --- SessionState Tests ---

// TestSessionStateUpdateProgress tests progress metric calculation.
func TestSessionStateUpdateProgress(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Progress Test Plan",
		Description:   "A plan for progress testing",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "Task 1",
				Type:        TaskTypeParallel,
			},
			{
				ID:          "task2",
				Description: "Task 2",
				Type:        TaskTypeParallel,
			},
			{
				ID:          "task3",
				Description: "Task 3",
				Type:        TaskTypeParallel,
			},
		},
	}

	session := NewSessionState("session1", plan)

	// Verify initial state
	if session.CompletedTasks != 0 {
		t.Fatalf("expected 0 completed tasks, got %d", session.CompletedTasks)
	}
	if session.Progress() != 0 {
		t.Fatalf("expected 0%% progress, got %.1f%%", session.Progress())
	}

	// Simulate task completions
	event1 := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventCompleted,
		TaskID:    "task1",
		Message:   "Task 1 completed",
	}

	if err := session.UpdateProgress(event1); err != nil {
		t.Fatalf("failed to update progress: %v", err)
	}

	if session.CompletedTasks != 1 {
		t.Fatalf("expected 1 completed task, got %d", session.CompletedTasks)
	}
	if session.Progress() != float64(1)/float64(3)*100 {
		t.Fatalf("expected %.1f%% progress, got %.1f%%", float64(1)/float64(3)*100, session.Progress())
	}

	// Complete another task
	event2 := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventCompleted,
		TaskID:    "task2",
		Message:   "Task 2 completed",
	}

	if err := session.UpdateProgress(event2); err != nil {
		t.Fatalf("failed to update progress: %v", err)
	}

	if session.CompletedTasks != 2 {
		t.Fatalf("expected 2 completed tasks, got %d", session.CompletedTasks)
	}
	if session.Progress() != float64(2)/float64(3)*100 {
		t.Fatalf("expected %.1f%% progress, got %.1f%%", float64(2)/float64(3)*100, session.Progress())
	}
}

// TestAgentTypeString tests AgentType string conversion.
func TestAgentTypeString(t *testing.T) {
	tests := []struct {
		agentType AgentType
		expected  string
	}{
		{AgentOrchestrator, "Orchestrator"},
		{AgentOptimizer, "Optimizer"},
		{AgentReviewer, "Reviewer"},
		{AgentExecutor, "Executor"},
		{AgentType(999), "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.agentType.String(); got != tt.expected {
			t.Fatalf("AgentType.String() = %s, expected %s", got, tt.expected)
		}
	}
}

// TestSessionStateToJSON tests JSON serialization of session state.
func TestSessionStateToJSON(t *testing.T) {
	plan := &WorkPlan{
		Name:          "JSON Test Plan",
		Description:   "A plan for JSON testing",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "Task 1",
				Type:        TaskTypeParallel,
			},
			{
				ID:          "task2",
				Description: "Task 2",
				Type:        TaskTypeParallel,
				Dependencies: []string{"task1"},
			},
		},
	}

	session := NewSessionState("session-test", plan)

	// Simulate some activity
	event := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventCompleted,
		TaskID:    "task1",
		Message:   "Task completed",
	}
	session.UpdateProgress(event)

	// Serialize to JSON
	jsonBytes, err := session.ToJSON()
	if err != nil {
		t.Fatalf("failed to serialize session state to JSON: %v", err)
	}

	if len(jsonBytes) == 0 {
		t.Fatal("expected non-empty JSON bytes")
	}

	// Verify it's valid JSON
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		t.Fatalf("generated JSON is not valid: %v", err)
	}

	// Check key fields
	if id, ok := data["id"].(string); !ok || id != "session-test" {
		t.Fatalf("expected session ID 'session-test' in JSON")
	}
	if status, ok := data["status"].(string); !ok || status != "running" {
		t.Fatalf("expected status 'running' in JSON")
	}
}

// TestParseFromLineInvalidJSON tests parsing invalid JSON.
func TestParseFromLineInvalidJSON(t *testing.T) {
	invalidLine := "not a valid json event"
	_, err := ParseFromLine(invalidLine)
	if err == nil {
		t.Fatal("expected error when parsing invalid JSON")
	}
}

// TestWorkPlanEmptyTasks tests validation of plan with no tasks.
func TestWorkPlanEmptyTasks(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Empty Plan",
		Description:   "A plan with no tasks",
		MaxConcurrent: 1,
		Tasks:         []Task{}, // Empty
	}

	err := plan.Validate()
	if err == nil {
		t.Fatal("expected error for plan with no tasks")
	}
	if err.Error() != "work plan must contain at least one task" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestWorkPlanInvalidMaxConcurrent tests validation of invalid maxConcurrent.
func TestWorkPlanInvalidMaxConcurrent(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Invalid Concurrent Plan",
		Description:   "A plan with invalid maxConcurrent",
		MaxConcurrent: 0, // Invalid
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "Task 1",
				Type:        TaskTypeParallel,
			},
		},
	}

	err := plan.Validate()
	if err == nil {
		t.Fatal("expected error for maxConcurrent < 1")
	}
	if err.Error() != "maxConcurrent must be at least 1, got 0" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestWorkPlanDuplicateTaskIDs tests detection of duplicate task IDs.
func TestWorkPlanDuplicateTaskIDs(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Duplicate IDs Plan",
		Description:   "A plan with duplicate task IDs",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "First task",
				Type:        TaskTypeParallel,
			},
			{
				ID:          "task1", // Duplicate ID
				Description: "Second task with same ID",
				Type:        TaskTypeParallel,
			},
		},
	}

	err := plan.Validate()
	if err == nil {
		t.Fatal("expected error for duplicate task IDs")
	}
	if err.Error() != "duplicate task ID: task1" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestSessionStateUpdateProgressNilEvent tests nil event handling.
func TestSessionStateUpdateProgressNilEvent(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Nil Event Plan",
		Description:   "A plan for nil event testing",
		MaxConcurrent: 1,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "Task 1",
				Type:        TaskTypeParallel,
			},
		},
	}

	session := NewSessionState("session-nil", plan)
	err := session.UpdateProgress(nil)
	if err == nil {
		t.Fatal("expected error for nil event")
	}
}

// TestWorkPlanFromJSONInvalid tests invalid JSON deserialization.
func TestWorkPlanFromJSONInvalid(t *testing.T) {
	plan := &WorkPlan{}
	invalidJSON := []byte("{invalid json")

	err := plan.FromJSON(invalidJSON)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// --- Helper Functions ---

// contains checks if a string contains a substring (case-sensitive).
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
