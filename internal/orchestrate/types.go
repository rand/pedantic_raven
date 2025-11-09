// Package orchestrate provides types and structures for Orchestrate Mode,
// which enables real-time coordination of multi-agent systems via mnemosyne.
//
// Orchestrate Mode provides:
// - WorkPlan: High-level task decomposition with dependencies
// - Task: Individual work units with type and agent assignment
// - AgentEvent: Real-time events from orchestration agents
// - SessionState: Tracking of orchestration progress and metrics
package orchestrate

import (
	"encoding/json"
	"fmt"
	"time"
)

// --- Task Types and Status ---

// TaskType categorizes the execution mode of a task.
type TaskType int

const (
	TaskTypeParallel TaskType = iota
	TaskTypeSequential
	TaskTypeBlocking
)

// String returns the string representation of TaskType.
func (t TaskType) String() string {
	switch t {
	case TaskTypeParallel:
		return "parallel"
	case TaskTypeSequential:
		return "sequential"
	case TaskTypeBlocking:
		return "blocking"
	default:
		return "unknown"
	}
}

// TaskStatus represents the current state of a task.
type TaskStatus int

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusActive
	TaskStatusCompleted
	TaskStatusFailed
)

// String returns the string representation of TaskStatus.
func (s TaskStatus) String() string {
	switch s {
	case TaskStatusPending:
		return "pending"
	case TaskStatusActive:
		return "active"
	case TaskStatusCompleted:
		return "completed"
	case TaskStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// --- Agent Types ---

// AgentType identifies the specialized role of an agent in orchestration.
type AgentType int

const (
	AgentOrchestrator AgentType = iota // Coordinates overall execution
	AgentOptimizer                     // Optimizes resource allocation
	AgentReviewer                      // Reviews and validates work
	AgentExecutor                      // Executes tasks
)

// String returns the string representation of AgentType.
func (a AgentType) String() string {
	switch a {
	case AgentOrchestrator:
		return "Orchestrator"
	case AgentOptimizer:
		return "Optimizer"
	case AgentReviewer:
		return "Reviewer"
	case AgentExecutor:
		return "Executor"
	default:
		return "Unknown"
	}
}

// --- Event Types ---

// EventType categorizes orchestration events.
type EventType int

const (
	EventStarted EventType = iota
	EventProgress
	EventCompleted
	EventFailed
	EventHandoff // Agent-to-agent communication
	EventLog     // General log message
)

// String returns the string representation of EventType.
func (e EventType) String() string {
	switch e {
	case EventStarted:
		return "started"
	case EventProgress:
		return "progress"
	case EventCompleted:
		return "completed"
	case EventFailed:
		return "failed"
	case EventHandoff:
		return "handoff"
	case EventLog:
		return "log"
	default:
		return "unknown"
	}
}

// --- Core Structures ---

// Task represents a single unit of work in a WorkPlan.
type Task struct {
	ID           string    `json:"id"`
	Description  string    `json:"description"`
	Dependencies []string  `json:"dependencies"`       // Task IDs this task depends on
	Type         TaskType  `json:"type"`               // parallel, sequential, blocking
	Agent        AgentType `json:"agent,omitempty"`    // preferred agent
	Priority     int       `json:"priority,omitempty"` // 0-10, higher = more urgent
}

// Validate checks that a task is well-formed.
func (t *Task) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if t.Description == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	if t.Priority < 0 || t.Priority > 10 {
		return fmt.Errorf("task priority must be between 0 and 10, got %d", t.Priority)
	}
	return nil
}

// WorkPlan represents a structured decomposition of work into tasks with dependencies.
type WorkPlan struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Tasks         []Task `json:"tasks"`
	MaxConcurrent int    `json:"maxConcurrent"`
}

// Validate checks that the work plan is well-formed and acyclic.
func (w *WorkPlan) Validate() error {
	if w.Name == "" {
		return fmt.Errorf("work plan name cannot be empty")
	}
	if len(w.Tasks) == 0 {
		return fmt.Errorf("work plan must contain at least one task")
	}

	// Validate individual tasks
	taskMap := make(map[string]bool)
	for _, task := range w.Tasks {
		if err := task.Validate(); err != nil {
			return err
		}
		if taskMap[task.ID] {
			return fmt.Errorf("duplicate task ID: %s", task.ID)
		}
		taskMap[task.ID] = true
	}

	// Validate that all dependencies exist
	for _, task := range w.Tasks {
		for _, depID := range task.Dependencies {
			if !taskMap[depID] {
				return fmt.Errorf("task %s depends on non-existent task %s", task.ID, depID)
			}
		}
	}

	// Check for circular dependencies using DFS
	if err := w.detectCycles(); err != nil {
		return err
	}

	if w.MaxConcurrent < 1 {
		return fmt.Errorf("maxConcurrent must be at least 1, got %d", w.MaxConcurrent)
	}

	return nil
}

// detectCycles uses depth-first search to detect circular dependencies.
// Uses three colors: white (unvisited), gray (visiting), black (visited).
func (w *WorkPlan) detectCycles() error {
	// Build adjacency map
	graph := make(map[string][]string)
	for _, task := range w.Tasks {
		if _, ok := graph[task.ID]; !ok {
			graph[task.ID] = []string{}
		}
		graph[task.ID] = append(graph[task.ID], task.Dependencies...)
	}

	colors := make(map[string]int) // 0 = white, 1 = gray, 2 = black
	for _, task := range w.Tasks {
		colors[task.ID] = 0
	}

	var visit func(string) error
	visit = func(taskID string) error {
		colors[taskID] = 1 // Mark as gray (visiting)

		for _, depID := range graph[taskID] {
			if colors[depID] == 1 {
				// Gray node found = cycle
				return fmt.Errorf("circular dependency detected: %s -> %s", taskID, depID)
			}
			if colors[depID] == 0 {
				// White node, recurse
				if err := visit(depID); err != nil {
					return err
				}
			}
		}

		colors[taskID] = 2 // Mark as black (visited)
		return nil
	}

	for _, task := range w.Tasks {
		if colors[task.ID] == 0 {
			if err := visit(task.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// ToDependencyGraph converts the work plan to a dependency graph representation.
// Returns a map where key is task ID and value is the list of tasks that depend on it.
func (w *WorkPlan) ToDependencyGraph() (map[string][]string, error) {
	if err := w.Validate(); err != nil {
		return nil, err
	}

	graph := make(map[string][]string)

	// Initialize all tasks in the graph
	for _, task := range w.Tasks {
		if _, ok := graph[task.ID]; !ok {
			graph[task.ID] = []string{}
		}
	}

	// Build reverse dependency graph (who depends on whom)
	for _, task := range w.Tasks {
		for _, depID := range task.Dependencies {
			graph[depID] = append(graph[depID], task.ID)
		}
	}

	return graph, nil
}

// ToJSON serializes the work plan to JSON bytes.
func (w *WorkPlan) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

// FromJSON deserializes a work plan from JSON bytes.
func (w *WorkPlan) FromJSON(data []byte) error {
	if err := json.Unmarshal(data, w); err != nil {
		return fmt.Errorf("failed to unmarshal work plan: %w", err)
	}
	return w.Validate()
}

// --- Agent Event ---

// AgentEvent represents a real-time event from an orchestration agent.
type AgentEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Agent     AgentType              `json:"agent"`
	EventType EventType              `json:"eventType"`
	TaskID    string                 `json:"taskId"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Format returns a human-readable string representation of the event.
func (e *AgentEvent) Format() string {
	return fmt.Sprintf("[%s] %s/%s: task=%s msg=%q",
		e.Timestamp.Format("15:04:05"),
		e.Agent.String(),
		e.EventType.String(),
		e.TaskID,
		e.Message,
	)
}

// ToJSON serializes the event to JSON bytes.
func (e *AgentEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON deserializes an event from JSON bytes.
func (e *AgentEvent) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// ParseFromLine attempts to parse an AgentEvent from a log line.
// Supports JSON format and structured plaintext format.
func ParseFromLine(line string) (*AgentEvent, error) {
	// Try JSON format first
	var event AgentEvent
	if err := json.Unmarshal([]byte(line), &event); err == nil {
		// Successfully parsed as JSON
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}
		return &event, nil
	}

	// Could not parse as JSON
	// Return a generic error
	return nil, fmt.Errorf("could not parse event from line: %s", line)
}

// --- Agent Status ---

// AgentStatus represents the current state of an agent.
type AgentStatus struct {
	Agent       AgentType
	Status      string // "active", "idle", "error"
	CurrentTask string // Task ID currently being executed
	LastUpdate  time.Time
}

// IsActive returns whether the agent is currently active.
func (a *AgentStatus) IsActive() bool {
	return a.Status == "active"
}

// CanExecute checks whether the agent can execute a given task.
// Returns true if the agent is idle or if no agent is specified for the task.
func (a *AgentStatus) CanExecute(task *Task) bool {
	if a.Status != "idle" && a.Status != "active" {
		return false
	}
	// Can execute if no agent preference or agent matches
	return task.Agent == a.Agent || task.Agent == 0 // 0 = AgentOrchestrator, no preference
}

// --- Session State ---

// SessionState represents the overall state of an orchestration session.
type SessionState struct {
	ID             string
	Plan           *WorkPlan
	StartTime      time.Time
	EndTime        *time.Time
	Status         string                // "running", "paused", "completed", "failed", "cancelled"
	TaskStatuses   map[string]TaskStatus // task ID -> status
	Agents         map[AgentType]*AgentStatus
	TotalTasks     int
	CompletedTasks int
	FailedTasks    int
}

// NewSessionState creates a new session state from a work plan.
func NewSessionState(id string, plan *WorkPlan) *SessionState {
	taskStatuses := make(map[string]TaskStatus)
	for _, task := range plan.Tasks {
		taskStatuses[task.ID] = TaskStatusPending
	}

	agents := make(map[AgentType]*AgentStatus)
	for i := AgentOrchestrator; i <= AgentExecutor; i++ {
		agents[i] = &AgentStatus{
			Agent:      i,
			Status:     "idle",
			LastUpdate: time.Now(),
		}
	}

	return &SessionState{
		ID:             id,
		Plan:           plan,
		StartTime:      time.Now(),
		Status:         "running",
		TaskStatuses:   taskStatuses,
		Agents:         agents,
		TotalTasks:     len(plan.Tasks),
		CompletedTasks: 0,
		FailedTasks:    0,
	}
}

// UpdateProgress updates session state based on an agent event.
func (s *SessionState) UpdateProgress(event *AgentEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Update agent status
	if agent, ok := s.Agents[event.Agent]; ok {
		agent.LastUpdate = time.Now()
		switch event.EventType {
		case EventStarted:
			agent.Status = "active"
			agent.CurrentTask = event.TaskID
		case EventCompleted:
			agent.Status = "idle"
			agent.CurrentTask = ""
		case EventFailed:
			agent.Status = "error"
		case EventProgress:
			agent.Status = "active"
		}
	}

	// Update task status
	if event.TaskID != "" {
		switch event.EventType {
		case EventStarted:
			s.TaskStatuses[event.TaskID] = TaskStatusActive
		case EventCompleted:
			s.TaskStatuses[event.TaskID] = TaskStatusCompleted
			s.CompletedTasks++
		case EventFailed:
			s.TaskStatuses[event.TaskID] = TaskStatusFailed
			s.FailedTasks++
		}
	}

	return nil
}

// Progress returns the completion percentage (0-100).
func (s *SessionState) Progress() float64 {
	if s.TotalTasks == 0 {
		return 0
	}
	return float64(s.CompletedTasks) / float64(s.TotalTasks) * 100
}

// ToJSON serializes the session state to JSON bytes.
func (s *SessionState) ToJSON() ([]byte, error) {
	// Create a serializable version
	type agentStatusJSON struct {
		Agent       string `json:"agent"`
		Status      string `json:"status"`
		CurrentTask string `json:"currentTask"`
		LastUpdate  string `json:"lastUpdate"`
	}

	type taskStatusJSON struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	type sessionJSON struct {
		ID             string            `json:"id"`
		StartTime      string            `json:"startTime"`
		EndTime        *string           `json:"endTime,omitempty"`
		Status         string            `json:"status"`
		TaskStatuses   []taskStatusJSON  `json:"taskStatuses"`
		Agents         []agentStatusJSON `json:"agents"`
		TotalTasks     int               `json:"totalTasks"`
		CompletedTasks int               `json:"completedTasks"`
		FailedTasks    int               `json:"failedTasks"`
		Progress       float64           `json:"progress"`
	}

	taskStatuses := make([]taskStatusJSON, 0, len(s.TaskStatuses))
	for id, status := range s.TaskStatuses {
		taskStatuses = append(taskStatuses, taskStatusJSON{
			ID:     id,
			Status: status.String(),
		})
	}

	agents := make([]agentStatusJSON, 0, len(s.Agents))
	for _, agent := range s.Agents {
		agents = append(agents, agentStatusJSON{
			Agent:       agent.Agent.String(),
			Status:      agent.Status,
			CurrentTask: agent.CurrentTask,
			LastUpdate:  agent.LastUpdate.Format(time.RFC3339),
		})
	}

	var endTimeStr *string
	if s.EndTime != nil {
		endTimeStr = new(string)
		*endTimeStr = s.EndTime.Format(time.RFC3339)
	}

	data := sessionJSON{
		ID:             s.ID,
		StartTime:      s.StartTime.Format(time.RFC3339),
		EndTime:        endTimeStr,
		Status:         s.Status,
		TaskStatuses:   taskStatuses,
		Agents:         agents,
		TotalTasks:     s.TotalTasks,
		CompletedTasks: s.CompletedTasks,
		FailedTasks:    s.FailedTasks,
		Progress:       s.Progress(),
	}

	return json.Marshal(data)
}

// FromJSON deserializes a session state from JSON bytes.
func (s *SessionState) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}
