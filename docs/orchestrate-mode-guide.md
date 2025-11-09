# Orchestrate Mode User Guide

**Version**: Phase 7
**Status**: Complete
**Last Updated**: 2025-11-09

## Overview

Orchestrate Mode provides a real-time dashboard for multi-agent coordination using mnemosyne's orchestration system. It enables you to:

- Define complex work plans with task dependencies
- Launch parallel agents for coordinated execution
- Monitor real-time progress across 4 specialized agents
- Visualize task dependency graphs
- View agent communication logs
- Manage orchestration sessions (pause, resume, cancel)

## Architecture

Orchestrate Mode integrates with mnemosyne's 4-agent orchestration system:

1. **Orchestrator**: Coordinates handoffs, manages state, schedules parallel work
2. **Optimizer**: Constructs context payloads, applies ACE principles
3. **Reviewer**: Validates intent, checks anti-patterns, enforces quality gates
4. **Executor**: Follows Work Plan Protocol, spawns sub-agents, executes tasks

## Getting Started

### Entering Orchestrate Mode

There are three ways to enter Orchestrate Mode:

1. **Keyboard shortcut**: Press `4` from any mode
2. **Command palette**: Press `Ctrl+K`, type "orchestrate", select "Switch to Orchestrate Mode"
3. **Mode switching**: Use the mode switcher UI

### Interface Overview

Orchestrate Mode has 4 views accessible via Tab or number keys:

```
┌─────────────────────────────────────────────────────────────┐
│ ORCHESTRATE MODE - Plan Editor                              │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  [1] Plan Editor  [2] Dashboard  [3] Task Graph  [4] Logs   │
│                                                               │
│  View 1: Create/edit work plans                             │
│  View 2: Monitor agent status and progress                  │
│  View 3: Visualize task dependencies                        │
│  View 4: View agent communication logs                      │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## View 1: Plan Editor

The Plan Editor allows you to create and edit work plans in JSON format.

### Creating a New Plan

1. Press `Ctrl+N` (New plan)
2. A template will load:
```json
{
  "name": "New Plan",
  "description": "A new work plan",
  "maxConcurrent": 2,
  "tasks": []
}
```

### Work Plan Format

A work plan consists of:

- **name**: Human-readable plan name
- **description**: Brief description of the plan's purpose
- **maxConcurrent**: Maximum parallel tasks (1-16)
- **tasks**: Array of task definitions

### Task Structure

Each task has:

```json
{
  "id": "unique-task-id",
  "description": "What this task does",
  "type": 0,
  "dependencies": ["task-id-1", "task-id-2"],
  "priority": 5
}
```

**Fields**:
- `id` (string, required): Unique task identifier
- `description` (string, required): Task description
- `type` (int, required): 0 = parallel, 1 = sequential
- `dependencies` (array, optional): List of task IDs that must complete first
- `priority` (int, optional): Priority 0-10 (10 = highest)

### Example: CI/CD Pipeline

```json
{
  "name": "CI Pipeline",
  "description": "Build, test, and deploy application",
  "maxConcurrent": 2,
  "tasks": [
    {
      "id": "checkout",
      "description": "Checkout source code",
      "type": 1,
      "dependencies": [],
      "priority": 10
    },
    {
      "id": "build",
      "description": "Compile application",
      "type": 1,
      "dependencies": ["checkout"],
      "priority": 9
    },
    {
      "id": "unit-tests",
      "description": "Run unit tests",
      "type": 0,
      "dependencies": ["build"],
      "priority": 8
    },
    {
      "id": "integration-tests",
      "description": "Run integration tests",
      "type": 0,
      "dependencies": ["build"],
      "priority": 8
    },
    {
      "id": "deploy",
      "description": "Deploy to staging",
      "type": 1,
      "dependencies": ["unit-tests", "integration-tests"],
      "priority": 7
    }
  ]
}
```

### Validation

The editor validates your plan in real-time:

- ✓ Valid JSON syntax
- ✓ Required fields present
- ✓ No circular dependencies
- ✓ All dependency IDs exist
- ✓ Priority within 0-10 range
- ✓ maxConcurrent within 1-16

Validation errors appear at the bottom of the editor in red.

### Keyboard Shortcuts (Plan Editor)

| Key | Action |
|-----|--------|
| `Ctrl+N` | New plan (load template) |
| `Ctrl+O` | Open plan from file |
| `Ctrl+S` | Save plan to file |
| `Ctrl+L` | Launch orchestration |
| `Arrow keys` | Navigate cursor |
| `Home/End` | Jump to line start/end |
| `Backspace/Delete` | Delete characters |

## View 2: Dashboard

The Dashboard shows real-time orchestration status.

### Agent Status Panel

```
Agents:
  [●] Orchestrator   (Active)   Task: Coordinate team
  [●] Optimizer      (Active)   Task: Build context
  [◐] Reviewer       (Idle)     Last: Validated API
  [●] Executor       (Active)   Task: Run tests
```

**Agent States**:
- `[●]` Active - Currently executing a task
- `[◐]` Idle - Waiting for work
- `[◯]` Inactive - Not started
- `[✗]` Error - Failed task

### Progress Metrics

```
Progress: [████████████████░░░░] 16/20 tasks (80%)
Success Rate: 15/16 ✓ (93.75%)
Elapsed: 12m 34s
```

**Metrics**:
- **Progress bar**: Visual completion indicator
- **Task count**: Completed/total tasks
- **Success rate**: Successful tasks / completed tasks
- **Elapsed time**: Session duration

### Task Queue

Shows upcoming tasks:

```
Queue:
  1. deploy (waiting on: integration-tests)
  2. cleanup (waiting on: deploy)
```

## View 3: Task Graph

The Task Graph visualizes dependencies using a force-directed layout.

### Graph Elements

**Nodes** (circles with task IDs):
- Gray: Pending (not started)
- Yellow: Active (in progress)
- Green: Completed (successful)
- Red: Failed (error occurred)

**Edges** (arrows):
- Show dependency direction (A → B means B depends on A)

### Navigation

| Key | Action |
|-----|--------|
| `h/j/k/l` | Pan left/down/up/right |
| `+/-` | Zoom in/out |
| `0` | Reset view |
| `f` | Fit graph to window |

### Example Graph

```
    [checkout]
       ↓
    [build]
       ↓
    ┌──┴──┐
    ↓     ↓
[unit]  [integration]
    ↓     ↓
    └──┬──┘
       ↓
    [deploy]
```

## View 4: Agent Log

The Agent Log displays agent communication and events.

### Log Entries

Each entry shows:

```
2025-11-09 14:23:45  Executor  INFO  task1  Started execution
2025-11-09 14:23:50  Executor  INFO  task1  Completed successfully
2025-11-09 14:23:51  Reviewer  INFO  task1  Validation passed
```

**Columns**:
1. Timestamp (ISO 8601)
2. Agent type (Orchestrator, Optimizer, Reviewer, Executor)
3. Log level (INFO, WARN, ERROR)
4. Task ID
5. Message

### Filtering

| Key | Action |
|-----|--------|
| `/` | Search logs (regex) |
| `a` | Filter by agent |
| `l` | Filter by log level |
| `c` | Clear filters |

### Exporting

Press `e` to export filtered logs to a file:

```
logs-2025-11-09-142345.txt
```

Format: Tab-separated values (TSV)

## Orchestration Controls

### Launching Orchestration

1. Create or load a work plan in the Plan Editor
2. Ensure plan is valid (green checkmark)
3. Press `Ctrl+L` to launch

**What happens**:
- mnemosyne orchestrate spawns as subprocess
- 4 agents initialize
- Tasks execute according to dependencies
- Events stream to dashboard in real-time

### Pause/Resume

Press `Space` to pause/resume orchestration:

- **Paused**: Agents stop accepting new tasks (current tasks complete)
- **Resumed**: Agents continue from last state

### Restart

Press `r` to restart orchestration:

1. Current session stops
2. All state resets
3. New session starts with same plan

### Cancel

Press `x` to cancel orchestration:

1. All agents stop gracefully (SIGTERM)
2. Session saves to history
3. Returns to Plan Editor

## Session Management

### Session Persistence

Orchestration sessions are automatically saved to:

```
~/.pedantic_raven/orchestrate/sessions/<session-id>.json
```

**Saved data**:
- Work plan
- Task progress
- Agent statuses
- Start/end times
- Completion metrics

### Session History

View past sessions:

```bash
ls ~/.pedantic_raven/orchestrate/sessions/
```

Load a previous session:

```bash
# Via Plan Editor: Ctrl+O, navigate to session file
```

## Keyboard Shortcuts Reference

### Global Controls

| Key | Action |
|-----|--------|
| `q` | Exit Orchestrate Mode |
| `?` | Toggle help overlay |
| `Tab` | Next view |
| `Shift+Tab` | Previous view |
| `1` | Plan Editor |
| `2` | Dashboard |
| `3` | Task Graph |
| `4` | Agent Log |

### Orchestration Controls

| Key | Action |
|-----|--------|
| `Ctrl+L` | Launch orchestration |
| `Space` | Pause/Resume |
| `r` | Restart |
| `x` | Cancel |

### Plan Editor

| Key | Action |
|-----|--------|
| `Ctrl+N` | New plan |
| `Ctrl+O` | Open plan |
| `Ctrl+S` | Save plan |
| `Arrow keys` | Navigate |
| `Home/End` | Line start/end |

### Task Graph

| Key | Action |
|-----|--------|
| `h/j/k/l` | Pan |
| `+/-` | Zoom |
| `0` | Reset view |
| `f` | Fit to window |

### Agent Log

| Key | Action |
|-----|--------|
| `/` | Search |
| `a` | Filter by agent |
| `l` | Filter by level |
| `e` | Export logs |
| `c` | Clear filters |

## Example Workflows

### Workflow 1: Simple Test Suite

```json
{
  "name": "Test Suite",
  "description": "Run all tests in parallel",
  "maxConcurrent": 4,
  "tasks": [
    {"id": "unit", "description": "Unit tests", "type": 0, "dependencies": []},
    {"id": "integration", "description": "Integration tests", "type": 0, "dependencies": []},
    {"id": "e2e", "description": "E2E tests", "type": 0, "dependencies": []},
    {"id": "lint", "description": "Linting", "type": 0, "dependencies": []},
    {"id": "report", "description": "Generate report", "type": 1, "dependencies": ["unit", "integration", "e2e", "lint"]}
  ]
}
```

### Workflow 2: Database Migration

```json
{
  "name": "Database Migration",
  "description": "Migrate prod database with validation",
  "maxConcurrent": 1,
  "tasks": [
    {"id": "backup", "description": "Backup database", "type": 1, "dependencies": [], "priority": 10},
    {"id": "analyze", "description": "Analyze schema", "type": 1, "dependencies": ["backup"], "priority": 9},
    {"id": "migrate", "description": "Run migration", "type": 1, "dependencies": ["analyze"], "priority": 8},
    {"id": "validate", "description": "Validate data", "type": 1, "dependencies": ["migrate"], "priority": 7},
    {"id": "cutover", "description": "Switch to new DB", "type": 1, "dependencies": ["validate"], "priority": 6}
  ]
}
```

### Workflow 3: Parallel Data Processing

```json
{
  "name": "Data Pipeline",
  "description": "Process data batches in parallel",
  "maxConcurrent": 8,
  "tasks": [
    {"id": "fetch", "description": "Fetch data", "type": 1, "dependencies": []},
    {"id": "batch1", "description": "Process batch 1", "type": 0, "dependencies": ["fetch"]},
    {"id": "batch2", "description": "Process batch 2", "type": 0, "dependencies": ["fetch"]},
    {"id": "batch3", "description": "Process batch 3", "type": 0, "dependencies": ["fetch"]},
    {"id": "batch4", "description": "Process batch 4", "type": 0, "dependencies": ["fetch"]},
    {"id": "aggregate", "description": "Aggregate results", "type": 1, "dependencies": ["batch1", "batch2", "batch3", "batch4"]}
  ]
}
```

## Troubleshooting

### Plan Validation Fails

**Error**: "circular dependency detected involving task X"

**Solution**: Check task dependencies for cycles. Use Task Graph (view 3) to visualize.

**Error**: "task depends on non-existent task Y"

**Solution**: Ensure all dependency IDs match existing task IDs exactly (case-sensitive).

### Orchestration Won't Launch

**Error**: "mnemosyne orchestrate not found"

**Solution**: Ensure mnemosyne is installed and in PATH:
```bash
which mnemosyne
mnemosyne --version
```

**Error**: "plan validation failed"

**Solution**: Check validation errors at bottom of Plan Editor. Fix all errors before launching.

### Agent Stuck/Unresponsive

**Symptom**: Agent shows as "Active" but no progress for > 5 minutes

**Solution**:
1. Check Agent Log (view 4) for error messages
2. Press `r` to restart orchestration
3. If persists, press `x` to cancel and check mnemosyne logs

### Performance Issues

**Symptom**: Dashboard updates slowly, UI laggy

**Solution**:
- Reduce `maxConcurrent` (try 2-4 instead of 8-16)
- Simplify task graph (fewer tasks, fewer dependencies)
- Filter Agent Log (fewer log entries displayed)

## Best Practices

### Plan Design

1. **Start simple**: Begin with 5-10 tasks, add complexity gradually
2. **Use priorities**: Critical path tasks should have priority 8-10
3. **Balance concurrency**: More isn't always faster (context switching overhead)
4. **Name meaningfully**: Use descriptive task IDs and descriptions
5. **Group related tasks**: Use common prefixes (e.g., "test-unit", "test-integration")

### Dependency Management

1. **Minimize dependencies**: Only depend on what's strictly required
2. **Avoid deep chains**: Long sequential chains reduce parallelization
3. **Use fan-out patterns**: One task → many parallel tasks → one aggregation
4. **Test for cycles**: Validate plan before launching

### Monitoring

1. **Watch Dashboard**: Primary view for active orchestration
2. **Check Task Graph**: Understand progress visually
3. **Review logs**: Agent Log shows detailed execution trace
4. **Export logs**: Keep records of important orchestrations

### Session Management

1. **Save plans**: Use Ctrl+S to save reusable plans
2. **Name sessions**: Use descriptive plan names for easy identification
3. **Review history**: Learn from past orchestrations
4. **Clean up**: Periodically archive old sessions

## Advanced Features

### Priority-Based Scheduling

Tasks with higher priority execute first when multiple tasks are ready:

```json
{"id": "critical", "priority": 10, "dependencies": []},
{"id": "normal", "priority": 5, "dependencies": []},
{"id": "low", "priority": 1, "dependencies": []}
```

Execution order (assuming maxConcurrent=1): critical → normal → low

### Dynamic Concurrency

Adjust `maxConcurrent` based on system load:

- **CPU-bound tasks**: maxConcurrent ≈ CPU cores
- **I/O-bound tasks**: maxConcurrent = 2-4× CPU cores
- **Mixed workload**: maxConcurrent = 1.5× CPU cores

### Task Types

- **Parallel (type: 0)**: Can run concurrently with other tasks
- **Sequential (type: 1)**: Blocks other tasks until complete

Use sequential for:
- Database migrations
- Critical path operations
- Resource contention (file writes, etc.)

## Integration with mnemosyne

Orchestrate Mode is a frontend for `mnemosyne orchestrate`. You can also use mnemosyne directly:

```bash
# Export plan from Orchestrate Mode
cd ~/.pedantic_raven/orchestrate/sessions/
cat session-xyz.json

# Run via CLI
mnemosyne orchestrate --plan plan.json --database ~/.mnemosyne/db

# Monitor with Orchestrate Mode dashboard
# Launch Pedantic Raven, enter Orchestrate Mode (key 4)
# Dashboard will show real-time progress
```

## Future Enhancements

Planned features for future releases:

- Visual plan editor (drag-and-drop task builder)
- Plan templates library (common patterns)
- Agent analytics (performance metrics, bottleneck detection)
- Cloud integration (distributed orchestration)
- Notifications (Slack, Discord alerts)
- Agent debugging (breakpoints, step-through)
- Cost tracking (API usage monitoring)

## Support

For issues, questions, or feature requests:

- GitHub Issues: https://github.com/rand/pedantic-raven/issues
- Documentation: https://docs.pedantic-raven.com
- Community: https://discord.gg/pedantic-raven

---

**Last Updated**: 2025-11-09
**Phase**: 7 (Orchestrate Mode)
**Version**: 1.0.0
