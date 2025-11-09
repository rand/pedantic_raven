# Phase 7: Orchestrate Mode - Multi-Agent Coordination

**Status**: Planning
**Phase**: 7 of 9
**Timeline**: 15 days (3 weeks)
**Dependencies**: Phase 5 (mnemosyne integration) complete

## Overview

Orchestrate Mode provides a real-time dashboard for multi-agent coordination using mnemosyne's orchestration system. It allows users to define work plans, launch parallel agents, monitor progress, and manage complex multi-step tasks with dependency tracking.

## Objectives

1. **Work Plan Management**: Create, edit, and save work plan definitions (JSON/prompt format)
2. **Agent Orchestration**: Launch and control `mnemosyne orchestrate` with customizable concurrency
3. **Real-Time Dashboard**: Live monitoring of agent status, progress, and task completion
4. **Task Dependency Visualization**: Graph view of task dependencies and execution flow
5. **Agent Communication**: View agent logs, errors, and handoffs between agents
6. **Session Management**: Pause, resume, and cancel orchestration sessions

## Architecture

### Components

```
internal/orchestrate/
├── orchestrate_mode.go      // Main mode coordinator (Bubble Tea model)
├── plan_editor.go            // Work plan creation/editing UI
├── dashboard.go              // Real-time agent monitoring dashboard
├── task_graph.go             // Dependency graph visualization
├── agent_log.go              // Agent communication log viewer
├── session.go                // Orchestration session management
├── launcher.go               // mnemosyne orchestrate process launcher
├── types.go                  // Core types (WorkPlan, AgentStatus, Task, etc.)
└── *_test.go                 // Test files
```

### Data Flow

```
User Input → Orchestrate Mode → mnemosyne orchestrate
                ↓                        ↓
         Plan Editor                  Orchestrator
                ↓                        ↓
         Work Plan (JSON)         4 Specialized Agents
                ↓                   ↓    ↓    ↓    ↓
         Session Manager         [Orchestrator][Optimizer][Reviewer][Executor]
                ↓                        ↓
         Dashboard ←──────────────── Agent Events
```

### mnemosyne Agents

The orchestration system uses 4 specialized agents (defined in CLAUDE.md):

1. **Orchestrator**: Coordinate handoffs, manage state, schedule parallel work, prevent deadlocks
2. **Optimizer**: Construct context payloads, apply ACE principles, monitor context budget
3. **Reviewer**: Validate intent satisfaction, check anti-patterns, enforce quality gates
4. **Executor**: Follow Work Plan Protocol, spawn sub-agents, execute tasks

## Features

### Week 1: Core Infrastructure (Days 1-5)

#### Day 1-2: Work Plan Editor
- JSON schema for work plans
- Simple text editor for plan creation
- Plan validation and syntax highlighting
- Save/load plan files
- Example plan templates

**Types**:
```go
type WorkPlan struct {
    Name        string
    Description string
    Tasks       []Task
    MaxConcurrent int
}

type Task struct {
    ID           string
    Description  string
    Dependencies []string // Task IDs
    Type         TaskType // parallel, sequential, blocking
}
```

#### Day 3-4: Orchestration Launcher
- Spawn `mnemosyne orchestrate` as subprocess
- Parse command-line arguments
- Environment setup (DB path, polling interval, etc.)
- Process lifecycle management (start, stop, restart)
- Capture stdout/stderr for logs

**Interface**:
```go
type Launcher interface {
    Start(plan WorkPlan, opts LaunchOptions) error
    Stop() error
    Restart() error
    IsRunning() bool
    Events() <-chan AgentEvent
}
```

#### Day 5: Session Management
- Session state tracking
- Persistence across restarts
- Session history
- Session metadata (start time, duration, task count)

### Week 2: Dashboard & Visualization (Days 6-10)

#### Day 6-7: Real-Time Dashboard
- Agent status panel (4 agents: Orchestrator, Optimizer, Reviewer, Executor)
- Task progress indicators
- Completion metrics (tasks done/total, success rate)
- Timeline view of agent activity
- Live updates via polling

**Dashboard Layout**:
```
╔══════════════════════════════════════════════════════════╗
║ ORCHESTRATE MODE - Session: Deploy Pipeline             ║
╠══════════════════════════════════════════════════════════╣
║ Agents:                                                  ║
║   [●] Orchestrator   (Active)   Task: Coordinate team   ║
║   [●] Optimizer      (Active)   Task: Build context     ║
║   [◐] Reviewer       (Idle)     Last: Validated API     ║
║   [●] Executor       (Active)   Task: Run tests         ║
╠══════════════════════════════════════════════════════════╣
║ Progress: [████████████████░░░░] 16/20 tasks (80%)      ║
║ Success Rate: 15/16 ✓ (93.75%)                          ║
║ Elapsed: 12m 34s                                         ║
╠══════════════════════════════════════════════════════════╣
║ [1] Plan Editor  [2] Dashboard  [3] Task Graph  [4] Logs║
╚══════════════════════════════════════════════════════════╝
```

#### Day 8-9: Task Dependency Graph
- Directed acyclic graph (DAG) visualization
- Node colors: pending (gray), active (yellow), done (green), failed (red)
- Edge types: sequential, parallel, blocking
- Force-directed layout (reuse from Analyze Mode)
- Pan/zoom navigation

**Graph Format**:
```
    [Task A]
       ↓
    [Task B] ──→ [Task C]
       ↓            ↓
    [Task D] ←─────┘
```

#### Day 10: Agent Communication Log
- Scrollable log viewer
- Filtering by agent type
- Color-coded log levels (info, warn, error)
- Search functionality
- Export logs to file

### Week 3: Integration & Polish (Days 11-15)

#### Day 11-12: Mode Integration
- Register OrchestrateMod in mode registry
- Keyboard shortcuts (Ctrl+O to enter, q to exit, Tab to switch views)
- Integrate with main application layout
- Connect to mnemosyne client from Phase 5
- Test with real orchestration examples

#### Day 13: Advanced Features
- Pause/resume orchestration
- Cancel running tasks
- Retry failed tasks
- Agent concurrency adjustment (live)
- Dashboard refresh rate control

#### Day 14: Testing
- Unit tests for all components
- Integration tests with mock mnemosyne process
- End-to-end test with real mnemosyne orchestrate
- Error handling tests (agent crashes, timeout, invalid plans)
- Stress test (100+ tasks, high concurrency)

#### Day 15: Documentation & Polish
- User guide for Orchestrate Mode
- Example work plans (CI/CD pipeline, data migration, batch processing)
- API documentation
- Keyboard shortcuts reference
- Error message improvements

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Ctrl+O` | Enter Orchestrate Mode |
| `q` | Exit mode |
| `Tab` | Switch view (Plan Editor → Dashboard → Task Graph → Logs) |
| `Shift+Tab` | Previous view |
| `1` | Plan Editor view |
| `2` | Dashboard view |
| `3` | Task Graph view |
| `4` | Logs view |
| `Space` | Pause/Resume orchestration |
| `r` | Restart orchestration |
| `x` | Cancel orchestration |
| `e` | Export logs |
| `n` | New plan |
| `o` | Open plan file |
| `s` | Save plan |
| `/` | Search logs |
| `+/-` | Zoom in/out (Task Graph) |
| `h/j/k/l` | Pan (Task Graph) |
| `?` | Show help |

## Technical Specifications

### Work Plan JSON Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "name": {"type": "string"},
    "description": {"type": "string"},
    "maxConcurrent": {"type": "integer", "minimum": 1, "maximum": 16},
    "tasks": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "id": {"type": "string"},
          "description": {"type": "string"},
          "dependencies": {"type": "array", "items": {"type": "string"}},
          "type": {"enum": ["parallel", "sequential", "blocking"]},
          "agent": {"enum": ["orchestrator", "optimizer", "reviewer", "executor"]},
          "priority": {"type": "integer", "minimum": 0, "maximum": 10}
        },
        "required": ["id", "description"]
      }
    }
  },
  "required": ["name", "tasks"]
}
```

### Agent Event Stream

Events from `mnemosyne orchestrate` (parsed from stdout):

```go
type AgentEvent struct {
    Timestamp   time.Time
    Agent       AgentType  // orchestrator, optimizer, reviewer, executor
    EventType   EventType  // started, progress, completed, failed, handoff
    TaskID      string
    Message     string
    Metadata    map[string]interface{}
}

type EventType int

const (
    EventStarted EventType = iota
    EventProgress
    EventCompleted
    EventFailed
    EventHandoff  // Agent-to-agent communication
    EventLog      // General log message
)
```

### Performance Targets

- **UI Responsiveness**: 60 FPS dashboard updates
- **Event Processing**: < 10ms latency from event to UI
- **Max Tasks**: Support 1000+ tasks in dependency graph
- **Max Agents**: Support 16 concurrent agents
- **Memory**: < 100 MB for typical session (50 tasks)
- **Log Retention**: Last 10,000 log lines in memory

## Testing Strategy

### Unit Tests (50+ tests)
- Plan validation
- Task dependency resolution
- DAG cycle detection
- Event parsing
- Agent status tracking
- Session persistence

### Integration Tests (20+ tests)
- Mock orchestration process
- Event stream simulation
- Dashboard updates
- View switching
- Keyboard input handling

### End-to-End Tests (5+ tests)
- Real mnemosyne orchestrate execution
- Complete workflow (create plan → launch → monitor → complete)
- Error recovery (agent crash, network failure)
- Session resume after restart

## Dependencies

### Existing (from Phase 5)
- `internal/mnemosyne` - mnemosyne client
- `internal/modes` - mode registry
- `internal/layout` - layout engine
- `internal/overlay` - overlay system (for dialogs)

### New (Go modules)
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- DAG library (consider `github.com/yourbasic/graph` or custom)

## Success Criteria

- [ ] Users can create and edit work plans in JSON format
- [ ] Orchestration launches successfully with `mnemosyne orchestrate`
- [ ] Dashboard shows real-time agent status and progress
- [ ] Task graph visualizes dependencies correctly
- [ ] Logs display agent communication and errors
- [ ] Pause/resume/cancel controls work reliably
- [ ] All tests pass (unit, integration, e2e)
- [ ] Documentation complete (user guide, API reference)
- [ ] Performance targets met (60 FPS, < 10ms latency)

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| mnemosyne orchestrate API instability | High | Version lock, comprehensive error handling |
| Event parsing complexity | Medium | Well-defined event schema, extensive tests |
| Performance with 1000+ tasks | Medium | Lazy rendering, viewport culling for graph |
| Dashboard refresh rate | Medium | Debouncing, configurable refresh rate |
| Session state corruption | High | Atomic writes, versioned schema, backups |

## Future Enhancements (Post-Phase 7)

- **Plan Templates**: Library of common orchestration patterns
- **Visual Plan Editor**: Drag-and-drop task graph builder
- **Agent Analytics**: Performance metrics, bottleneck detection
- **Cloud Integration**: Distributed orchestration across machines
- **Slack/Discord Notifications**: Real-time alerts on task completion/failure
- **Agent Debugging**: Breakpoints, step-through execution
- **Cost Tracking**: API usage and cost per orchestration session

## Deliverables

### Code (est. ~3,000 lines)
- `internal/orchestrate/` package with 9+ files
- 75+ tests (unit + integration + e2e)

### Documentation (est. ~800 lines)
- `docs/orchestrate-mode-guide.md` (user guide)
- `docs/orchestrate-mode-api.md` (API reference)
- `examples/work_plans/` (example JSON plans)

### Diagrams
- Architecture diagram (Mermaid)
- Data flow diagram
- State machine diagram

## Open Questions

1. **Event Format**: Does `mnemosyne orchestrate` output structured JSON events or plaintext logs?
   - **Action**: Test mnemosyne orchestrate with --dashboard flag, capture output format

2. **Plan Format**: Is there a standard JSON schema for work plans in mnemosyne?
   - **Action**: Check mnemosyne docs/source code for plan schema

3. **Agent Communication**: How do agents communicate handoffs? Via stdout, files, or API?
   - **Action**: Review mnemosyne orchestration internals

4. **Dashboard API**: Does mnemosyne provide an HTTP API for real-time dashboard data?
   - **Action**: Check for `api-server` integration with orchestrate command

## Next Steps

1. **Validate Assumptions**: Test `mnemosyne orchestrate` with sample plan, capture output
2. **Prototype Event Parser**: Build minimal parser for mnemosyne events
3. **Create Plan Schema**: Define JSON schema based on mnemosyne requirements
4. **Build MVP Dashboard**: Simple 4-agent status display
5. **Iterate**: Add features incrementally following 3-week timeline

---

**Last Updated**: 2025-11-09
**Author**: Phase 7 Planning
**Status**: Ready for Review → Implementation
