package orchestrate

import (
	"strings"
	"testing"
	"time"
)

// --- Rendering Tests ---

// TestDashboardRenderAgentPanel tests correct agent display formatting.
func TestDashboardRenderAgentPanel(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:          "task1",
				Description: "Task 1",
				Type:        TaskTypeParallel,
			},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Manually set agent statuses
	session.Agents[AgentOrchestrator].Status = "active"
	session.Agents[AgentOrchestrator].CurrentTask = "task1"
	session.Agents[AgentOptimizer].Status = "idle"
	session.Agents[AgentReviewer].Status = "error"
	session.Agents[AgentExecutor].Status = "idle"

	agentPanel := dashboard.renderAgentPanel()

	// Verify output contains all agents
	if !strings.Contains(agentPanel, "Orchestrator") {
		t.Errorf("agent panel missing Orchestrator: %s", agentPanel)
	}
	if !strings.Contains(agentPanel, "Optimizer") {
		t.Errorf("agent panel missing Optimizer: %s", agentPanel)
	}
	if !strings.Contains(agentPanel, "Reviewer") {
		t.Errorf("agent panel missing Reviewer: %s", agentPanel)
	}
	if !strings.Contains(agentPanel, "Executor") {
		t.Errorf("agent panel missing Executor: %s", agentPanel)
	}

	// Verify status indicators are present
	if !strings.Contains(agentPanel, "[●]") {
		t.Error("agent panel missing active indicator [●]")
	}
	if !strings.Contains(agentPanel, "[◐]") {
		t.Error("agent panel missing idle indicator [◐]")
	}
	if !strings.Contains(agentPanel, "[✗]") {
		t.Error("agent panel missing error indicator [✗]")
	}
}

// TestDashboardRenderProgress tests progress bar calculation and formatting.
func TestDashboardRenderProgress(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
			{ID: "task2", Description: "Task 2", Type: TaskTypeParallel},
			{ID: "task3", Description: "Task 3", Type: TaskTypeParallel},
			{ID: "task4", Description: "Task 4", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Simulate completing 2 out of 4 tasks
	dashboard.completedTasks = 2
	dashboard.totalTasks = 4
	dashboard.calculateMetrics()

	progress := dashboard.renderProgress()

	// Verify progress bar is present
	if !strings.Contains(progress, "Progress:") {
		t.Error("progress string missing 'Progress:'")
	}

	// Verify percentage calculation (2/4 = 50%)
	if !strings.Contains(progress, "50%") {
		t.Errorf("progress bar missing correct percentage: %s", progress)
	}

	// Verify task counts
	if !strings.Contains(progress, "2/4") {
		t.Errorf("progress bar missing task count: %s", progress)
	}

	// Verify progress bar characters
	if !strings.Contains(progress, "█") {
		t.Error("progress bar missing filled character [█]")
	}
	if !strings.Contains(progress, "░") {
		t.Error("progress bar missing empty character [░]")
	}
}

// TestDashboardRenderMetrics tests metric formatting.
func TestDashboardRenderMetrics(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Set up metrics
	dashboard.startTime = time.Now().Add(-5 * time.Minute)
	dashboard.completedTasks = 3
	dashboard.failedTasks = 1
	dashboard.calculateMetrics()

	progress := dashboard.renderProgress()

	// Verify elapsed time is shown
	if !strings.Contains(progress, "Elapsed:") {
		t.Error("metrics missing elapsed time")
	}

	// Verify success rate is shown
	if !strings.Contains(progress, "Success Rate:") {
		t.Error("metrics missing success rate")
	}

	// Verify success rate calculation (3/4 = 75%)
	if !strings.Contains(progress, "75.00%") {
		t.Errorf("metrics missing correct success rate: %s", progress)
	}
}

// --- Event Handling Tests ---

// TestDashboardHandleEventStarted tests that event updates agent to active.
func TestDashboardHandleEventStarted(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Initially agent should be idle
	if session.Agents[AgentExecutor].Status != "idle" {
		t.Fatalf("expected agent to be idle, got %s", session.Agents[AgentExecutor].Status)
	}

	// Process event
	event := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventStarted,
		TaskID:    "task1",
		Message:   "Starting task",
	}
	dashboard.handleEvent(event)

	// Verify agent is now active
	if session.Agents[AgentExecutor].Status != "active" {
		t.Errorf("expected agent to be active, got %s", session.Agents[AgentExecutor].Status)
	}
	if session.Agents[AgentExecutor].CurrentTask != "task1" {
		t.Errorf("expected current task to be task1, got %s", session.Agents[AgentExecutor].CurrentTask)
	}
}

// TestDashboardHandleEventCompleted tests that completed event updates progress.
func TestDashboardHandleEventCompleted(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	initialCompleted := dashboard.completedTasks

	// Process completion event
	event := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventCompleted,
		TaskID:    "task1",
		Message:   "Task completed",
	}
	dashboard.handleEvent(event)

	// Verify progress was incremented
	if dashboard.completedTasks != initialCompleted+1 {
		t.Errorf("expected %d completed tasks, got %d", initialCompleted+1, dashboard.completedTasks)
	}

	// Verify agent is idle
	if session.Agents[AgentExecutor].Status != "idle" {
		t.Errorf("expected agent to be idle after completion, got %s", session.Agents[AgentExecutor].Status)
	}
}

// TestDashboardHandleEventFailed tests that failed event tracks failures.
func TestDashboardHandleEventFailed(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	initialFailed := dashboard.failedTasks

	// Process failure event
	event := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentExecutor,
		EventType: EventFailed,
		TaskID:    "task1",
		Message:   "Task failed",
	}
	dashboard.handleEvent(event)

	// Verify failures were tracked
	if dashboard.failedTasks != initialFailed+1 {
		t.Errorf("expected %d failed tasks, got %d", initialFailed+1, dashboard.failedTasks)
	}

	// Verify agent is in error state
	if session.Agents[AgentExecutor].Status != "error" {
		t.Errorf("expected agent to be in error state, got %s", session.Agents[AgentExecutor].Status)
	}
}

// TestDashboardHandleEventHandoff tests agent communication events.
func TestDashboardHandleEventHandoff(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Process handoff event
	event := &AgentEvent{
		Timestamp: time.Now(),
		Agent:     AgentOrchestrator,
		EventType: EventHandoff,
		TaskID:    "task1",
		Message:   "Handing off to Executor",
	}
	dashboard.handleEvent(event)

	// Verify last update was recorded
	if dashboard.lastUpdate.IsZero() {
		t.Error("expected last update to be recorded")
	}
}

// --- Metric Calculation Tests ---

// TestDashboardSuccessRate tests success rate calculation.
func TestDashboardSuccessRate(t *testing.T) {
	tests := []struct {
		name             string
		completed        int
		failed           int
		expectedRate     float64
		expectedRateStr  string
	}{
		{
			name:             "zero tasks",
			completed:        0,
			failed:           0,
			expectedRate:     0.0,
			expectedRateStr:  "0",
		},
		{
			name:             "all successful",
			completed:        10,
			failed:           0,
			expectedRate:     100.0,
			expectedRateStr:  "100",
		},
		{
			name:             "all failed",
			completed:        0,
			failed:           10,
			expectedRate:     0.0,
			expectedRateStr:  "0",
		},
		{
			name:             "75% success",
			completed:        3,
			failed:           1,
			expectedRate:     75.0,
			expectedRateStr:  "75",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &WorkPlan{
				Name:          "Test Plan",
				Description:   "A test plan",
				MaxConcurrent: 1,
				Tasks: []Task{
					{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
				},
			}

			session := NewSessionState("session1", plan)
			events := make(chan *AgentEvent, 10)
			dashboard := NewDashboard(session, events)

			dashboard.completedTasks = tt.completed
			dashboard.failedTasks = tt.failed
			dashboard.calculateMetrics()

			if dashboard.successRate != tt.expectedRate {
				t.Errorf("expected success rate %.2f, got %.2f", tt.expectedRate, dashboard.successRate)
			}

			// Verify it's used in output
			if tt.completed+tt.failed > 0 {
				progress := dashboard.renderProgress()
				if !strings.Contains(progress, tt.expectedRateStr) {
					t.Errorf("expected success rate string to contain %s, got: %s", tt.expectedRateStr, progress)
				}
			}
		})
	}
}

// TestDashboardElapsedTime tests time formatting.
func TestDashboardElapsedTime(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "30s"},
		{5 * time.Minute, "5m"},
		{1*time.Minute + 30*time.Second, "1m 30s"},
		{1*time.Hour + 5*time.Minute, "1h 5m"},
		{2*time.Hour + 15*time.Minute + 45*time.Second, "2h 15m"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatElapsed(tt.duration)
			if !strings.Contains(result, strings.Split(tt.expected, " ")[0]) {
				t.Errorf("formatElapsed(%v) = %s, expected to contain %s", tt.duration, result, tt.expected)
			}
		})
	}
}

// TestDashboardProgressBar tests progress bar width scaling.
func TestDashboardProgressBar(t *testing.T) {
	tests := []struct {
		name             string
		completed        int
		total            int
		expectedPercent  string
		shouldContainBar bool
	}{
		{
			name:             "0% complete",
			completed:        0,
			total:            10,
			expectedPercent:  "0%",
			shouldContainBar: true,
		},
		{
			name:             "50% complete",
			completed:        5,
			total:            10,
			expectedPercent:  "50%",
			shouldContainBar: true,
		},
		{
			name:             "100% complete",
			completed:        10,
			total:            10,
			expectedPercent:  "100%",
			shouldContainBar: true,
		},
		{
			name:             "33% complete",
			completed:        1,
			total:            3,
			expectedPercent:  "33%",
			shouldContainBar: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &WorkPlan{
				Name:          "Test Plan",
				Description:   "A test plan",
				MaxConcurrent: 1,
				Tasks: []Task{
					{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
				},
			}

			session := NewSessionState("session1", plan)
			events := make(chan *AgentEvent, 10)
			dashboard := NewDashboard(session, events)

			dashboard.completedTasks = tt.completed
			dashboard.totalTasks = tt.total

			progressBar := dashboard.renderProgressBar()

			// Verify percentage
			if !strings.Contains(progressBar, tt.expectedPercent) {
				t.Errorf("progress bar missing percentage %s: %s", tt.expectedPercent, progressBar)
			}

			// Verify bar is present
			if tt.shouldContainBar {
				if !strings.Contains(progressBar, "[") || !strings.Contains(progressBar, "]") {
					t.Errorf("progress bar missing brackets: %s", progressBar)
				}
			}

			// Verify bar characters
			if tt.completed > 0 {
				if !strings.Contains(progressBar, "█") {
					t.Errorf("progress bar missing filled character when partially complete: %s", progressBar)
				}
			}
		})
	}
}

// --- Integration and Mock Tests ---

// TestDashboardInitialization tests dashboard creation and initialization.
func TestDashboardInitialization(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 2,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
			{ID: "task2", Description: "Task 2", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Verify initial state
	if dashboard.orchestrating != true {
		t.Error("dashboard should be orchestrating initially")
	}

	if dashboard.totalTasks != len(plan.Tasks) {
		t.Errorf("expected %d total tasks, got %d", len(plan.Tasks), dashboard.totalTasks)
	}

	// Verify all 4 agents are initialized
	expectedAgents := 4
	if len(dashboard.agents) != expectedAgents {
		t.Errorf("expected %d agents, got %d", expectedAgents, len(dashboard.agents))
	}
}

// TestDashboardTaskQueueManagement tests task queue initialization and removal.
func TestDashboardTaskQueueManagement(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
			{ID: "task2", Description: "Task 2", Type: TaskTypeParallel},
			{ID: "task3", Description: "Task 3", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Verify initial queue
	if len(dashboard.taskQueue) != len(plan.Tasks) {
		t.Errorf("expected %d tasks in queue, got %d", len(plan.Tasks), len(dashboard.taskQueue))
	}

	// Remove a task from queue
	dashboard.removeFromQueue("task1")

	if len(dashboard.taskQueue) != len(plan.Tasks)-1 {
		t.Errorf("expected %d tasks after removal, got %d", len(plan.Tasks)-1, len(dashboard.taskQueue))
	}

	// Verify task1 is not in queue
	for _, id := range dashboard.taskQueue {
		if id == "task1" {
			t.Error("task1 should not be in queue after removal")
		}
	}
}

// TestDashboardTeaModelInterface tests that Dashboard implements tea.Model.
func TestDashboardTeaModelInterface(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Test Init() returns a Cmd (should not panic)
	cmd := dashboard.Init()
	if cmd == nil {
		t.Error("Init() should return a command")
	}

	// Test View() returns a string (should not panic)
	view := dashboard.View()
	if view == "" {
		t.Error("View() should return a non-empty string")
	}

	// Verify View contains expected content
	if !strings.Contains(view, "ORCHESTRATE MODE") {
		t.Error("View should contain ORCHESTRATE MODE header")
	}
}

// TestDashboardViewResponsiveness tests that View updates with events.
func TestDashboardViewResponsiveness(t *testing.T) {
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "A test plan",
		MaxConcurrent: 1,
		Tasks: []Task{
			{ID: "task1", Description: "Task 1", Type: TaskTypeParallel},
			{ID: "task2", Description: "Task 2", Type: TaskTypeParallel},
		},
	}

	session := NewSessionState("session1", plan)
	events := make(chan *AgentEvent, 10)
	dashboard := NewDashboard(session, events)

	// Get initial view
	view1 := dashboard.View()
	if !strings.Contains(view1, "0/2") {
		t.Errorf("Initial view should show 0/2 progress, got: %s", view1)
	}

	// Simulate event
	dashboard.completedTasks = 1
	dashboard.calculateMetrics()

	// Get updated view
	view2 := dashboard.View()

	// Verify progress is different
	if !strings.Contains(view2, "1/2") {
		t.Errorf("View should show updated progress, got: %s", view2)
	}
}

// TestFormatElapsedFunction tests the time formatting utility.
func TestFormatElapsedFunction(t *testing.T) {
	tests := []struct {
		duration time.Duration
		contains []string
	}{
		{
			duration: 45 * time.Second,
			contains: []string{"45", "s"},
		},
		{
			duration: 2*time.Minute + 30*time.Second,
			contains: []string{"2", "m", "30", "s"},
		},
		{
			duration: 1*time.Hour + 30*time.Minute + 15*time.Second,
			contains: []string{"1", "h", "30", "m"},
		},
	}

	for _, tt := range tests {
		result := formatElapsed(tt.duration)
		for _, expected := range tt.contains {
			if !strings.Contains(result, expected) {
				t.Errorf("formatElapsed(%v) = %s, expected to contain %s", tt.duration, result, expected)
			}
		}
	}
}
