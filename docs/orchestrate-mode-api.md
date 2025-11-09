# Orchestrate Mode API Reference

**Version**: Phase 7
**Package**: `github.com/rand/pedantic-raven/internal/orchestrate`
**Last Updated**: 2025-11-09

## Overview

This document provides API reference for the Orchestrate Mode implementation. It covers all public types, interfaces, and functions available for integration and extension.

## Package Structure

```
internal/orchestrate/
├── types.go              # Core data structures
├── launcher.go           # Process lifecycle management
├── session.go            # Session persistence
├── plan_editor.go        # Plan editing UI
├── dashboard.go          # Real-time monitoring
├── task_graph.go         # Dependency visualization
├── agent_log.go          # Log viewing
├── orchestrate_mode.go   # Main coordinator
└── mode_adapter.go       # Mode registry integration
```

## Core Types

### WorkPlan

Represents a complete orchestration work plan.

```go
type WorkPlan struct {
    Name          string `json:"name"`
    Description   string `json:"description"`
    MaxConcurrent int    `json:"maxConcurrent"`
    Tasks         []Task `json:"tasks"`
}
```

**Fields**:
- `Name`: Human-readable plan identifier
- `Description`: Brief description of plan purpose
- `MaxConcurrent`: Maximum parallel tasks (1-16)
- `Tasks`: Array of task definitions

**Methods**:

```go
func (wp *WorkPlan) Validate() error
```
Validates plan structure, checks for circular dependencies, verifies task references.

**Returns**: `nil` if valid, error describing validation failure otherwise.

**Example**:
```go
plan := &WorkPlan{
    Name:          "Test Plan",
    Description:   "Run test suite",
    MaxConcurrent: 4,
    Tasks: []Task{
        {
            ID:          "unit-tests",
            Description: "Run unit tests",
            Type:        TaskTypeParallel,
        },
    },
}

if err := plan.Validate(); err != nil {
    log.Fatal("Invalid plan:", err)
}
```

---

### Task

Represents a single task within a work plan.

```go
type Task struct {
    ID           string   `json:"id"`
    Description  string   `json:"description"`
    Type         TaskType `json:"type"`
    Dependencies []string `json:"dependencies,omitempty"`
    Priority     int      `json:"priority,omitempty"`
}
```

**Fields**:
- `ID`: Unique task identifier (required)
- `Description`: Task description (required)
- `Type`: Task execution type (parallel/sequential)
- `Dependencies`: List of task IDs that must complete first
- `Priority`: Execution priority 0-10 (10 = highest, default = 5)

**Methods**:

```go
func (t *Task) Validate() error
```
Validates individual task structure.

**Example**:
```go
task := Task{
    ID:           "build",
    Description:  "Compile application",
    Type:         TaskTypeSequential,
    Dependencies: []string{"checkout"},
    Priority:     9,
}
```

---

### TaskType

Enumeration of task execution types.

```go
type TaskType int

const (
    TaskTypeParallel   TaskType = 0
    TaskTypeSequential TaskType = 1
)
```

**Values**:
- `TaskTypeParallel`: Task can run concurrently with others
- `TaskTypeSequential`: Task blocks until completion

**Methods**:

```go
func (tt TaskType) String() string
```
Returns human-readable task type name.

---

### AgentType

Enumeration of agent types in orchestration system.

```go
type AgentType int

const (
    AgentOrchestrator AgentType = 0
    AgentOptimizer    AgentType = 1
    AgentReviewer     AgentType = 2
    AgentExecutor     AgentType = 3
)
```

**Values**:
- `AgentOrchestrator`: Coordinates handoffs and schedules work
- `AgentOptimizer`: Builds context payloads, applies ACE principles
- `AgentReviewer`: Validates output, enforces quality gates
- `AgentExecutor`: Executes tasks, spawns sub-agents

---

### AgentEvent

Represents an event from an agent during orchestration.

```go
type AgentEvent struct {
    Timestamp time.Time              `json:"timestamp"`
    Agent     AgentType              `json:"agent"`
    EventType EventType              `json:"event_type"`
    TaskID    string                 `json:"task_id"`
    Message   string                 `json:"message"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
```

**Fields**:
- `Timestamp`: Event occurrence time (UTC)
- `Agent`: Agent that generated event
- `EventType`: Type of event (started, progress, completed, etc.)
- `TaskID`: Associated task identifier
- `Message`: Human-readable event message
- `Metadata`: Additional structured data

**Methods**:

```go
func ParseFromLine(line string) (*AgentEvent, error)
```
Parses agent event from JSON log line.

```go
func (e *AgentEvent) Format() string
```
Formats event as human-readable string.

**Example**:
```go
event := &AgentEvent{
    Timestamp: time.Now(),
    Agent:     AgentExecutor,
    EventType: EventStarted,
    TaskID:    "build",
    Message:   "Starting build task",
}

formatted := event.Format()
// Output: "2025-11-09T14:23:45Z  Executor  STARTED  build  Starting build task"
```

---

### EventType

Enumeration of agent event types.

```go
type EventType int

const (
    EventStarted   EventType = 0
    EventProgress  EventType = 1
    EventCompleted EventType = 2
    EventFailed    EventType = 3
    EventHandoff   EventType = 4
    EventLog       EventType = 5
)
```

**Values**:
- `EventStarted`: Task execution started
- `EventProgress`: Task progress update
- `EventCompleted`: Task completed successfully
- `EventFailed`: Task failed with error
- `EventHandoff`: Agent-to-agent communication
- `EventLog`: General log message

---

### SessionState

Represents the state of an orchestration session.

```go
type SessionState struct {
    ID              string                  `json:"id"`
    Plan            *WorkPlan               `json:"plan"`
    StartTime       time.Time               `json:"start_time"`
    EndTime         time.Time               `json:"end_time,omitempty"`
    Status          SessionStatus           `json:"status"`
    TaskStatuses    map[string]TaskStatus   `json:"task_statuses"`
    AgentStatuses   map[AgentType]*AgentStatus `json:"agent_statuses"`
    TotalTasks      int                     `json:"total_tasks"`
    CompletedTasks  int                     `json:"completed_tasks"`
    FailedTasks     int                     `json:"failed_tasks"`
}
```

**Fields**:
- `ID`: Unique session identifier (UUID)
- `Plan`: Associated work plan
- `StartTime`: Session start timestamp
- `EndTime`: Session end timestamp (if complete)
- `Status`: Current session status
- `TaskStatuses`: Per-task execution status
- `AgentStatuses`: Per-agent status
- `TotalTasks`: Total task count
- `CompletedTasks`: Completed task count
- `FailedTasks`: Failed task count

**Methods**:

```go
func (ss *SessionState) UpdateProgress(event *AgentEvent) error
```
Updates session state based on agent event.

---

### SessionStatus

Enumeration of session status values.

```go
type SessionStatus int

const (
    SessionPending   SessionStatus = 0
    SessionRunning   SessionStatus = 1
    SessionPaused    SessionStatus = 2
    SessionCompleted SessionStatus = 3
    SessionFailed    SessionStatus = 4
    SessionCancelled SessionStatus = 5
)
```

---

### TaskStatus

Enumeration of task execution status values.

```go
type TaskStatus int

const (
    TaskStatusPending   TaskStatus = 0
    TaskStatusRunning   TaskStatus = 1
    TaskStatusCompleted TaskStatus = 2
    TaskStatusFailed    TaskStatus = 3
    TaskStatusSkipped   TaskStatus = 4
)
```

## Components

### Launcher

Manages the lifecycle of the `mnemosyne orchestrate` subprocess.

```go
type Launcher struct {
    // private fields
}
```

**Constructor**:

```go
func NewLauncher() *Launcher
```
Creates a new launcher instance.

**Methods**:

```go
func (l *Launcher) Start(plan *WorkPlan, opts LaunchOptions) error
```
Spawns mnemosyne orchestrate process with given plan.

**Parameters**:
- `plan`: Work plan to execute
- `opts`: Launch options (database path, polling interval, etc.)

**Returns**: Error if spawn fails, nil on success.

```go
func (l *Launcher) Stop() error
```
Stops orchestration gracefully (SIGTERM → SIGKILL after timeout).

```go
func (l *Launcher) IsRunning() bool
```
Returns true if orchestration process is running.

```go
func (l *Launcher) Events() <-chan *AgentEvent
```
Returns read-only channel of agent events.

**Example**:
```go
launcher := NewLauncher()

plan := &WorkPlan{...}
opts := LaunchOptions{
    DatabasePath: "~/.mnemosyne/db",
}

if err := launcher.Start(plan, opts); err != nil {
    log.Fatal("Failed to start:", err)
}

// Monitor events
for event := range launcher.Events() {
    log.Printf("Event: %s", event.Format())
}
```

---

### LaunchOptions

Configuration options for launcher.

```go
type LaunchOptions struct {
    DatabasePath   string        `json:"database_path,omitempty"`
    PollingInterval time.Duration `json:"polling_interval,omitempty"`
    Timeout        time.Duration `json:"timeout,omitempty"`
}
```

**Fields**:
- `DatabasePath`: Path to mnemosyne database
- `PollingInterval`: Event polling interval (default: 100ms)
- `Timeout`: Process timeout (default: 1 hour)

---

### Session

Manages session persistence and history.

```go
type Session struct {
    // private fields
}
```

**Constructor**:

```go
func NewSession(plan *WorkPlan) *Session
```
Creates a new session for the given plan.

**Methods**:

```go
func (s *Session) Save() error
```
Saves session state to disk (atomic write).

```go
func (s *Session) Load(id string) error
```
Loads session state from disk by ID.

```go
func (s *Session) UpdateProgress(event *AgentEvent) error
```
Updates session state based on event (auto-saves).

```go
func (s *Session) GetState() *SessionState
```
Returns current session state (read-only).

**Example**:
```go
session := NewSession(plan)

// Update from events
for event := range launcher.Events() {
    if err := session.UpdateProgress(event); err != nil {
        log.Printf("Update error: %v", err)
    }
}

// Final save
if err := session.Save(); err != nil {
    log.Fatal("Save failed:", err)
}
```

---

### PlanEditor

Bubble Tea model for editing work plans.

```go
type PlanEditor struct {
    // private fields
}
```

**Constructor**:

```go
func NewPlanEditor() *PlanEditor
```
Creates a new plan editor instance.

**Methods** (Bubble Tea interface):

```go
func (m *PlanEditor) Init() tea.Cmd
func (m *PlanEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *PlanEditor) View() string
```

**Keyboard Shortcuts**:
- `Ctrl+N`: New plan
- `Ctrl+O`: Open plan
- `Ctrl+S`: Save plan
- `Arrow keys`: Navigate cursor
- `Home/End`: Line start/end
- `Backspace/Delete`: Delete characters
- `Enter`: New line

---

### Dashboard

Bubble Tea model for real-time monitoring.

```go
type Dashboard struct {
    // private fields
}
```

**Constructor**:

```go
func NewDashboard(session *SessionState, events <-chan *AgentEvent) *Dashboard
```
Creates a dashboard monitoring the given session.

**Methods**:

```go
func (d *Dashboard) Init() tea.Cmd
func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (d *Dashboard) View() string
```

---

### TaskGraph

Bubble Tea model for DAG visualization.

```go
type TaskGraph struct {
    // private fields
}
```

**Constructor**:

```go
func NewTaskGraph(plan *WorkPlan) *TaskGraph
```
Creates a task graph from work plan.

**Methods**:

```go
func (tg *TaskGraph) Init() tea.Cmd
func (tg *TaskGraph) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (tg *TaskGraph) View() string
func (tg *TaskGraph) UpdateStatus(taskID string, status TaskStatus)
```

**Keyboard Shortcuts**:
- `h/j/k/l`: Pan left/down/up/right
- `+/-`: Zoom in/out
- `0`: Reset view
- `f`: Fit graph to window

---

### AgentLog

Bubble Tea model for log viewing.

```go
type AgentLog struct {
    // private fields
}
```

**Constructor**:

```go
func NewAgentLog(maxEntries int) *AgentLog
```
Creates a log viewer with circular buffer size.

**Methods**:

```go
func (al *AgentLog) Init() tea.Cmd
func (al *AgentLog) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (al *AgentLog) View() string
func (al *AgentLog) AddEntry(event *AgentEvent)
func (al *AgentLog) Export(filename string) error
```

**Keyboard Shortcuts**:
- `/`: Search (regex)
- `a`: Filter by agent
- `l`: Filter by log level
- `c`: Clear filters
- `e`: Export logs

---

### OrchestrateMode

Main coordinator Bubble Tea model.

```go
type OrchestrateMode struct {
    // private fields
}
```

**Constructor**:

```go
func NewOrchestrateMode() *OrchestrateMode
```
Creates a new orchestrate mode instance.

**Methods**:

```go
func (m *OrchestrateMode) Init() tea.Cmd
func (m *OrchestrateMode) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *OrchestrateMode) View() string
func (m *OrchestrateMode) IsOrchestrating() bool
func (m *OrchestrateMode) GetCurrentPlan() *WorkPlan
func (m *OrchestrateMode) GetSession() *Session
```

**Keyboard Shortcuts**:
- `q`: Exit mode
- `?`: Toggle help
- `Tab`: Next view
- `Shift+Tab`: Previous view
- `1-4`: Direct view navigation
- `Ctrl+L`: Launch orchestration
- `Space`: Pause/Resume
- `r`: Restart
- `x`: Cancel

---

### ModeAdapter

Adapter for mode registry integration.

```go
type ModeAdapter struct {
    // private fields
}
```

**Constructor**:

```go
func NewModeAdapter() modes.Mode
```
Creates a new mode adapter (implements `modes.Mode` interface).

**Methods** (modes.Mode interface):

```go
func (m *ModeAdapter) ID() modes.ModeID
func (m *ModeAdapter) Name() string
func (m *ModeAdapter) Description() string
func (m *ModeAdapter) Init() tea.Cmd
func (m *ModeAdapter) OnEnter() tea.Cmd
func (m *ModeAdapter) OnExit() tea.Cmd
func (m *ModeAdapter) Update(msg tea.Msg) (modes.Mode, tea.Cmd)
func (m *ModeAdapter) View() string
func (m *ModeAdapter) Keybindings() []modes.Keybinding
```

**Integration**:

```go
import (
    "github.com/rand/pedantic-raven/internal/modes"
    "github.com/rand/pedantic-raven/internal/orchestrate"
)

registry := modes.NewRegistry()
orchestrateMode := orchestrate.NewModeAdapter()
registry.Register(orchestrateMode)
```

## Utilities

### Validation Functions

```go
func ValidatePlan(plan *WorkPlan) error
```
Validates work plan structure and dependencies.

```go
func DetectCycles(tasks []Task) error
```
Detects circular dependencies in task graph.

**Algorithm**: Depth-first search with color marking (white/gray/black).

---

### JSON Serialization

All types support JSON marshaling/unmarshaling:

```go
// Marshal
planJSON, err := json.Marshal(plan)

// Unmarshal
var plan WorkPlan
err := json.Unmarshal(planJSON, &plan)
```

---

### Event Streaming

Events are streamed from launcher to components via channels:

```go
eventsChan := launcher.Events()

for event := range eventsChan {
    // Route to components
    dashboard.handleEvent(event)
    agentLog.AddEntry(event)
    taskGraph.UpdateStatus(event.TaskID, status)
    session.UpdateProgress(event)
}
```

## Testing

### Unit Tests

All components have comprehensive unit tests:

```bash
go test ./internal/orchestrate/... -v
```

**Coverage**: 60-90% per component

**Test Files**:
- `types_test.go` (23 tests)
- `launcher_test.go` (11 tests)
- `session_test.go` (15 tests)
- `plan_editor_test.go` (23 tests)
- `dashboard_test.go` (18 tests)
- `task_graph_test.go` (12 tests)
- `agent_log_test.go` (14 tests)
- `orchestrate_mode_test.go` (14 tests)
- `mode_adapter_test.go` (16 tests)

**Total**: 146 tests passing

---

### Integration Tests

Test with mock mnemosyne process:

```go
func TestOrchestrateIntegration(t *testing.T) {
    // Create plan
    plan := &WorkPlan{...}

    // Start launcher
    launcher := NewLauncher()
    err := launcher.Start(plan, LaunchOptions{})
    require.NoError(t, err)

    // Monitor events
    events := []AgentEvent{}
    for event := range launcher.Events() {
        events = append(events, *event)
        if event.EventType == EventCompleted {
            break
        }
    }

    // Verify
    assert.Greater(t, len(events), 0)
}
```

---

### Example Tests

```go
// Validate plan
func TestWorkPlanValidation(t *testing.T) {
    plan := &WorkPlan{
        Name:          "Test",
        MaxConcurrent: 2,
        Tasks: []Task{
            {ID: "a", Description: "A", Type: TaskTypeParallel},
            {ID: "b", Description: "B", Type: TaskTypeParallel, Dependencies: []string{"a"}},
        },
    }

    err := plan.Validate()
    assert.NoError(t, err)
}

// Detect cycle
func TestCycleDetection(t *testing.T) {
    plan := &WorkPlan{
        Name:          "Cyclic",
        MaxConcurrent: 1,
        Tasks: []Task{
            {ID: "a", Dependencies: []string{"b"}},
            {ID: "b", Dependencies: []string{"a"}},
        },
    }

    err := plan.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "circular")
}
```

## Performance

### Benchmarks

```bash
go test ./internal/orchestrate/... -bench=. -benchmem
```

**Results** (example):

```
BenchmarkLauncherStart-8         1000    1.2 ms/op    512 B/op     8 allocs/op
BenchmarkEventParsing-8        100000   15.3 µs/op    256 B/op     4 allocs/op
BenchmarkTaskGraphLayout-8      10000  120.0 µs/op   2048 B/op    32 allocs/op
```

---

### Resource Usage

**Typical session** (50 tasks, maxConcurrent=4):
- Memory: ~10 MB
- CPU: < 5% (event processing)
- Disk: 50-100 KB (session file)

**Large session** (1000 tasks, maxConcurrent=16):
- Memory: ~100 MB
- CPU: 10-20% (graph layout)
- Disk: 1-2 MB (session file)

---

### Optimization

**Event processing**:
- Buffered channels (100 events)
- Batch updates (10ms debounce)
- Selective rendering (viewport culling)

**Graph layout**:
- Force-directed algorithm (O(n²) per iteration)
- Damping coefficient: 0.85
- Max iterations: 50
- Layout cache (reuse when unchanged)

**Log storage**:
- Circular buffer (10,000 entries)
- O(1) append
- O(n) filtered view

## Errors

### Error Types

```go
var (
    ErrInvalidPlan      = errors.New("invalid work plan")
    ErrCircularDep      = errors.New("circular dependency detected")
    ErrMissingTask      = errors.New("task not found")
    ErrLaunchFailed     = errors.New("orchestration launch failed")
    ErrSessionNotFound  = errors.New("session not found")
)
```

---

### Error Handling

All public functions return errors using standard Go error handling:

```go
if err := plan.Validate(); err != nil {
    if errors.Is(err, ErrCircularDep) {
        log.Fatal("Plan has circular dependencies")
    }
    log.Fatal("Validation failed:", err)
}
```

## Extensions

### Custom Agents

Extend with custom agent types:

```go
const (
    AgentCustom AgentType = 100 + iota
    AgentMonitor
    AgentAnalyzer
)
```

---

### Custom Event Types

Define custom event types:

```go
const (
    EventCustom EventType = 100 + iota
    EventMetrics
    EventAlert
)
```

---

### Plugin Architecture

Future support for plugins via Go plugin system:

```go
type Plugin interface {
    OnEvent(event *AgentEvent) error
    OnComplete(session *SessionState) error
}
```

## Examples

### Complete Example

```go
package main

import (
    "log"
    "github.com/rand/pedantic-raven/internal/orchestrate"
)

func main() {
    // Create plan
    plan := &orchestrate.WorkPlan{
        Name:          "Example",
        Description:   "Example workflow",
        MaxConcurrent: 2,
        Tasks: []orchestrate.Task{
            {ID: "setup", Description: "Setup", Type: orchestrate.TaskTypeSequential},
            {ID: "test", Description: "Test", Type: orchestrate.TaskTypeParallel, Dependencies: []string{"setup"}},
            {ID: "deploy", Description: "Deploy", Type: orchestrate.TaskTypeSequential, Dependencies: []string{"test"}},
        },
    }

    // Validate
    if err := plan.Validate(); err != nil {
        log.Fatal("Invalid plan:", err)
    }

    // Create session
    session := orchestrate.NewSession(plan)

    // Launch orchestration
    launcher := orchestrate.NewLauncher()
    if err := launcher.Start(plan, orchestrate.LaunchOptions{}); err != nil {
        log.Fatal("Launch failed:", err)
    }

    // Monitor events
    for event := range launcher.Events() {
        log.Printf("Event: %s", event.Format())

        // Update session
        if err := session.UpdateProgress(event); err != nil {
            log.Printf("Update error: %v", err)
        }

        // Check completion
        if event.EventType == orchestrate.EventCompleted && session.GetState().CompletedTasks == len(plan.Tasks) {
            log.Println("All tasks completed!")
            break
        }
    }

    // Save session
    if err := session.Save(); err != nil {
        log.Fatal("Save failed:", err)
    }
}
```

## Changelog

### Phase 7 (2025-11-09)

- Initial release
- 8 core components implemented
- 146 tests passing (60-90% coverage)
- Full mnemosyne integration
- 4-view TUI (editor, dashboard, graph, logs)
- Session persistence
- Real-time event streaming

---

**Last Updated**: 2025-11-09
**Phase**: 7 (Orchestrate Mode)
**Package Version**: 1.0.0
