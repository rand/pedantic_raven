# Phase 7 Execution Plan

**Phase**: 7 (Orchestrate Mode)
**Timeline**: 15 days (3 weeks)
**Start Date**: 2025-11-09
**Status**: Ready for Execution

## Critical Path

Based on dependency analysis, the critical path is:

```
Day 1: types.go (2 hours)
  ↓
Day 2-3: launcher.go (12 hours)
  ↓
Day 4-5: dashboard.go (12 hours)
  ↓
Day 6-7: task_graph.go (12 hours)
  ↓
Day 11-12: orchestrate_mode.go (12 hours)
  ↓
Day 14: Integration tests (6 hours)
  ↓
Day 15: Documentation (6 hours)
```

**Total Critical Path**: 62 hours (~8 working days)

## Parallel Execution Streams

### Stream A: Core Foundation (Days 1-5)
**Agent**: Foundation Agent (Haiku for speed)
```markdown
- [ ] Day 1 (2h): types.go + tests (15 tests)
- [ ] Day 2-3 (12h): launcher.go + tests (12 tests)
- [ ] Day 4-5 (12h): dashboard.go + tests (10 tests)
```
**Dependencies**: None → launcher → dashboard
**Output**: Core infrastructure ready for UI integration

### Stream B: Session Management (Days 2-3)
**Agent**: Session Agent (Haiku)
```markdown
- [ ] Day 2-3 (8h): session.go + tests (10 tests)
```
**Dependencies**: types.go (wait for Day 1)
**Output**: Session persistence and lifecycle management
**Can run parallel with**: launcher.go

### Stream C: Plan Editor (Day 4)
**Agent**: Editor Agent (Haiku)
```markdown
- [ ] Day 4 (6h): plan_editor.go + tests (8 tests)
```
**Dependencies**: types.go (wait for Day 1)
**Output**: JSON plan editing UI
**Can run parallel with**: dashboard.go

### Stream D: Task Graph (Days 6-7)
**Agent**: Graph Agent (Sonnet for complexity)
```markdown
- [ ] Day 6-7 (12h): task_graph.go + tests (12 tests)
  - DAG construction
  - Cycle detection
  - Force-directed layout (reuse from analyze/)
  - Rendering with pan/zoom
```
**Dependencies**: types.go, analyze/triple_graph.go (for force layout)
**Output**: Interactive task dependency graph
**Can run parallel with**: agent_log.go

### Stream E: Agent Log (Day 8)
**Agent**: Log Agent (Haiku)
```markdown
- [ ] Day 8 (6h): agent_log.go + tests (8 tests)
```
**Dependencies**: types.go
**Output**: Scrollable, filterable log viewer
**Can run parallel with**: task_graph.go (partial overlap)

### Stream F: Integration (Days 9-12)
**Agent**: Integration Agent (Sonnet)
```markdown
- [ ] Day 9-10 (12h): orchestrate_mode.go - View coordination
- [ ] Day 11-12 (8h): Mode registration and keyboard shortcuts
- [ ] Day 12 (4h): Connect to mnemosyne client
```
**Dependencies**: ALL previous streams
**Output**: Fully integrated Orchestrate Mode

### Stream G: Advanced Features (Day 13)
**Agent**: Features Agent (Haiku)
```markdown
- [ ] Day 13 (6h): Pause/resume/cancel controls
  - Implement control signals
  - Update dashboard UI
  - Add keyboard shortcuts (Space, r, x)
  - Test error recovery
```
**Dependencies**: orchestrate_mode.go complete
**Output**: Advanced orchestration control

### Stream H: Testing (Day 14)
**Agent**: Test Agent (Sonnet for thoroughness)
```markdown
- [ ] Day 14 Morning (4h): Integration tests (20 tests)
  - launcher + dashboard event flow
  - session + mode lifecycle
  - plan editor validation
  - Full view switching
- [ ] Day 14 Afternoon (4h): E2E tests (5 tests)
  - Complete workflow test
  - Error recovery test
  - Cancellation test
  - Parallel execution test
  - Large scale test (100 tasks)
```
**Dependencies**: All code complete
**Output**: 100 total tests passing, 85%+ coverage

### Stream I: Documentation (Day 15)
**Agent**: Docs Agent (Haiku)
```markdown
- [ ] Day 15 Morning (3h): orchestrate-mode-guide.md
  - User guide with screenshots (ASCII art)
  - Workflow examples
  - Keyboard shortcuts reference
  - Troubleshooting

- [ ] Day 15 Afternoon (3h): orchestrate-mode-api.md + examples
  - API documentation
  - Code examples (4 integration patterns)
  - Example work plans (5 JSON files)
  - Architecture diagrams (Mermaid)
```
**Dependencies**: All code complete
**Output**: ~800 lines of documentation

---

## Gantt Chart (ASCII)

```
Week 1: Foundation & Core
Day    1    2    3    4    5
-------|====|====|====|====|====
types  [==]
launch     [========]
session    [========]
editor          [====]
dash            [========]
-------|====|====|====|====|====

Week 2: UI & Integration
Day    6    7    8    9   10
-------|====|====|====|====|====
graph  [========]
log         [====]
mode             [========]
-------|====|====|====|====|====

Week 3: Polish & Ship
Day   11   12   13   14   15
-------|====|====|====|====|====
mode  [====]
feats      [====]
tests           [====]
docs                 [====]
-------|====|====|====|====|====

Legend:
[====] Task duration
|      Day boundary
```

## Parallel Agent Assignments

### Day 1: Solo Work
- **Agent 1** (Haiku): types.go (2h)
  - No parallelization needed (foundation layer)

### Day 2-3: 2 Parallel Agents
- **Agent 2** (Haiku): launcher.go (12h)
- **Agent 3** (Haiku): session.go (8h)
  - **Safe**: No shared files, both depend only on types.go

### Day 4-5: 2 Parallel Agents
- **Agent 4** (Haiku): plan_editor.go (6h) on Day 4
- **Agent 5** (Haiku): dashboard.go (12h) on Days 4-5
  - **Safe**: Different files, different concerns

### Day 6-8: 2 Parallel Agents
- **Agent 6** (Sonnet): task_graph.go (12h) on Days 6-7
- **Agent 7** (Haiku): agent_log.go (6h) on Day 8
  - **Safe**: Independent components, both depend on types.go

### Day 9-12: Solo Work (Integration)
- **Agent 8** (Sonnet): orchestrate_mode.go (20h)
  - No parallelization (integrates all components)

### Day 13: Solo Work (Features)
- **Agent 9** (Haiku): Advanced features (6h)

### Day 14-15: 2 Parallel Agents
- **Agent 10** (Sonnet): Testing (8h) on Day 14
- **Agent 11** (Haiku): Documentation (6h) on Day 15 (can start after Agent 10)
  - **Sequential**: Docs should reflect tested code

---

## Agent Task Definitions

### Agent 1: Foundation (Day 1)
**Model**: Haiku (fast, simple types)
**Task**: Implement types.go
**Prompt**:
```
Implement internal/orchestrate/types.go for Pedantic Raven Phase 7.

Requirements:
- WorkPlan struct with JSON marshaling
- Task struct with dependency tracking
- AgentEvent parsing from mnemosyne stdout
- AgentStatus and AgentType enumerations
- SessionState struct
- 15 unit tests

Reference: docs/PHASE7_SPEC.md, docs/PHASE7_DECOMPOSITION.md

Deliverables:
- internal/orchestrate/types.go (~200 lines)
- internal/orchestrate/types_test.go (~300 lines)
- All tests passing
```

### Agent 2: Launcher (Days 2-3)
**Model**: Haiku
**Task**: Implement launcher.go
**Prompt**:
```
Implement internal/orchestrate/launcher.go for Pedantic Raven Phase 7.

Requirements:
- Spawn mnemosyne orchestrate as subprocess
- Capture stdout/stderr to event channel
- Process lifecycle management (start/stop/restart)
- Error handling and recovery
- 12 unit tests

Dependencies: types.go must be complete

Reference: docs/PHASE7_SPEC.md, docs/PHASE7_DECOMPOSITION.md

Deliverables:
- internal/orchestrate/launcher.go (~250 lines)
- internal/orchestrate/launcher_test.go (~350 lines)
- All tests passing
```

### Agent 3: Session (Days 2-3) [PARALLEL with Agent 2]
**Model**: Haiku
**Task**: Implement session.go
**Prompt**:
```
Implement internal/orchestrate/session.go for Pedantic Raven Phase 7.

Requirements:
- Session persistence (JSON to ~/.pedantic_raven/sessions/)
- Session history tracking
- Progress calculation
- Atomic writes with versioning
- 10 unit tests

Dependencies: types.go must be complete

Reference: docs/PHASE7_SPEC.md, docs/PHASE7_DECOMPOSITION.md

Deliverables:
- internal/orchestrate/session.go (~200 lines)
- internal/orchestrate/session_test.go (~250 lines)
- All tests passing

SAFE TO RUN IN PARALLEL with launcher.go (different files)
```

### Agent 4: Plan Editor (Day 4) [PARALLEL with Agent 5]
**Model**: Haiku
**Task**: Implement plan_editor.go
**Prompt**:
```
Implement internal/orchestrate/plan_editor.go for Pedantic Raven Phase 7.

Requirements:
- Bubble Tea model for JSON plan editing
- Real-time validation with error highlighting
- Save/load file operations
- 3 example plan templates
- 8 unit tests

Dependencies: types.go must be complete

Reference:
- docs/PHASE7_SPEC.md
- docs/PHASE7_DECOMPOSITION.md
- internal/modes/explore.go (for Bubble Tea pattern)

Deliverables:
- internal/orchestrate/plan_editor.go (~300 lines)
- internal/orchestrate/plan_editor_test.go (~250 lines)
- examples/work_plans/*.json (3 templates)
- All tests passing

SAFE TO RUN IN PARALLEL with dashboard.go (different files)
```

### Agent 5: Dashboard (Days 4-5)
**Model**: Haiku
**Task**: Implement dashboard.go
**Prompt**:
```
Implement internal/orchestrate/dashboard.go for Pedantic Raven Phase 7.

Requirements:
- Bubble Tea model for real-time agent monitoring
- 4-agent status panel (Orchestrator, Optimizer, Reviewer, Executor)
- Progress bars and metrics (elapsed time, success rate)
- Task queue visualization
- Event handling with auto-refresh
- 10 unit tests

Dependencies: types.go, launcher.go must be complete

Reference:
- docs/PHASE7_SPEC.md (dashboard layout mockup)
- docs/PHASE7_DECOMPOSITION.md
- internal/analyze/analyze_mode.go (for Bubble Tea + real-time updates)

Deliverables:
- internal/orchestrate/dashboard.go (~350 lines)
- internal/orchestrate/dashboard_test.go (~300 lines)
- All tests passing

Can run in parallel with plan_editor.go
```

### Agent 6: Task Graph (Days 6-7)
**Model**: Sonnet (complex graph algorithms)
**Task**: Implement task_graph.go
**Prompt**:
```
Implement internal/orchestrate/task_graph.go for Pedantic Raven Phase 7.

Requirements:
- DAG construction from WorkPlan
- Cycle detection (DFS with back edges)
- Force-directed layout (REUSE internal/analyze/triple_graph.go)
- Node/edge rendering with status colors
- Pan/zoom controls (h/j/k/l, +/-)
- 12 unit tests

Dependencies: types.go, internal/analyze/triple_graph.go

Reference:
- docs/PHASE7_SPEC.md
- docs/PHASE7_DECOMPOSITION.md
- internal/analyze/triple_graph.go (REUSE force layout algorithm)
- internal/analyze/view.go (REUSE canvas rendering)

Deliverables:
- internal/orchestrate/task_graph.go (~400 lines)
- internal/orchestrate/task_graph_test.go (~350 lines)
- All tests passing

IMPORTANT: Reuse force-directed layout code from analyze package
```

### Agent 7: Agent Log (Day 8)
**Model**: Haiku
**Task**: Implement agent_log.go
**Prompt**:
```
Implement internal/orchestrate/agent_log.go for Pedantic Raven Phase 7.

Requirements:
- Scrollable log viewer with circular buffer (max 10,000 lines)
- Color-coded log levels (info/warn/error)
- Filter by agent type
- Regex search
- Export to file
- 8 unit tests

Dependencies: types.go

Reference:
- docs/PHASE7_SPEC.md
- docs/PHASE7_DECOMPOSITION.md

Deliverables:
- internal/orchestrate/agent_log.go (~250 lines)
- internal/orchestrate/agent_log_test.go (~200 lines)
- All tests passing
```

### Agent 8: Orchestrate Mode (Days 9-12)
**Model**: Sonnet (complex integration)
**Task**: Implement orchestrate_mode.go and integrate with main app
**Prompt**:
```
Implement internal/orchestrate/orchestrate_mode.go and integrate Orchestrate Mode into Pedantic Raven.

Requirements:
- Bubble Tea model coordinating all 4 views
- View switching (Tab/Shift+Tab, 1/2/3/4)
- Keyboard shortcut handling (all shortcuts from PHASE7_SPEC.md)
- Session lifecycle management
- Event routing to child components
- Header/footer rendering with help
- Mode registration in internal/modes/registry.go
- 15 unit tests

Dependencies: ALL previous components (types, launcher, session, plan_editor, dashboard, task_graph, agent_log)

Reference:
- docs/PHASE7_SPEC.md (keyboard shortcuts table)
- docs/PHASE7_DECOMPOSITION.md
- internal/analyze/analyze_mode.go (for view switching pattern)
- internal/modes/explore.go (for mode pattern)
- internal/modes/registry.go (for registration)

Deliverables:
- internal/orchestrate/orchestrate_mode.go (~500 lines)
- internal/orchestrate/orchestrate_mode_test.go (~400 lines)
- Updated internal/modes/registry.go (register ModeOrchestrate)
- All tests passing (90 total in package)

IMPORTANT: This is integration work - no parallelization
```

### Agent 9: Advanced Features (Day 13)
**Model**: Haiku
**Task**: Implement pause/resume/cancel controls
**Prompt**:
```
Add advanced orchestration controls to internal/orchestrate/orchestrate_mode.go.

Requirements:
- Pause/resume orchestration (Space key)
- Restart orchestration (r key)
- Cancel orchestration (x key)
- Error recovery handling
- Updated tests for control signals

Reference:
- docs/PHASE7_SPEC.md (keyboard shortcuts)

Deliverables:
- Updated orchestrate_mode.go with control handlers
- Updated tests
- All tests passing
```

### Agent 10: Testing (Day 14)
**Model**: Sonnet (thoroughness)
**Task**: Comprehensive testing
**Prompt**:
```
Write comprehensive tests for Pedantic Raven Phase 7 Orchestrate Mode.

Requirements:
- 20 integration tests (launcher+dashboard, session+mode, etc.)
- 5 E2E tests (complete workflow, error recovery, cancellation, parallel, large scale)
- Verify 85%+ code coverage
- All 100 total tests passing

Reference:
- docs/PHASE7_DECOMPOSITION.md (test plan)
- All internal/orchestrate/*.go files

Deliverables:
- internal/orchestrate/integration_test.go (~600 lines, 20 tests)
- internal/orchestrate/e2e_test.go (~400 lines, 5 tests)
- Coverage report showing 85%+
- All tests passing
```

### Agent 11: Documentation (Day 15)
**Model**: Haiku
**Task**: Write user guide and API docs
**Prompt**:
```
Write comprehensive documentation for Pedantic Raven Phase 7 Orchestrate Mode.

Requirements:
- User guide (docs/orchestrate-mode-guide.md, ~400 lines)
- API documentation (docs/orchestrate-mode-api.md, ~400 lines)
- 5 example work plans (examples/work_plans/*.json)
- 2 Mermaid architecture diagrams

Reference:
- docs/PHASE7_SPEC.md
- docs/analyze-mode-guide.md (for style/structure)
- docs/analyze-mode-api.md (for API doc pattern)
- All internal/orchestrate/*.go files

Deliverables:
- docs/orchestrate-mode-guide.md (~400 lines)
- docs/orchestrate-mode-api.md (~400 lines)
- examples/work_plans/ci-pipeline.json
- examples/work_plans/data-migration.json
- examples/work_plans/batch-processing.json
- examples/work_plans/test-suite.json
- examples/work_plans/deployment.json
```

---

## Dependencies Matrix

| Component | Depends On | Blocks |
|-----------|-----------|--------|
| types.go | - | ALL |
| launcher.go | types.go | dashboard.go, orchestrate_mode.go |
| session.go | types.go | orchestrate_mode.go |
| plan_editor.go | types.go | orchestrate_mode.go |
| dashboard.go | types.go, launcher.go | orchestrate_mode.go |
| task_graph.go | types.go, analyze/triple_graph.go | orchestrate_mode.go |
| agent_log.go | types.go | orchestrate_mode.go |
| orchestrate_mode.go | ALL above | advanced features, tests, docs |
| Advanced features | orchestrate_mode.go | tests, docs |
| Tests | ALL | docs |
| Docs | ALL | - |

---

## Git Workflow

### Feature Branch
```bash
git checkout main
git pull origin main
git checkout -b feature/phase7-orchestrate-mode
```

### Commit Strategy
Each agent creates a commit after completing their work:
```bash
# Agent 1
git add internal/orchestrate/types.go internal/orchestrate/types_test.go
git commit -m "Implement Orchestrate Mode types and core data structures (Day 1)"

# Agent 2
git add internal/orchestrate/launcher.go internal/orchestrate/launcher_test.go
git commit -m "Implement mnemosyne orchestrate launcher with process management (Days 2-3)"

# ... etc for each agent
```

### Pull Request
After all work complete on Day 15:
```bash
git push origin feature/phase7-orchestrate-mode
gh pr create --title "Phase 7: Orchestrate Mode - Multi-Agent Coordination" \
  --body-file docs/PHASE7_PR_TEMPLATE.md
```

---

## Success Criteria Checklist

Before marking Phase 7 complete, verify:

- [ ] All 100 tests passing (75 unit + 20 integration + 5 E2E)
- [ ] Code coverage ≥ 85%
- [ ] Dashboard updates at 60 FPS
- [ ] Event latency < 10ms
- [ ] Memory usage < 100 MB for 50-task session
- [ ] Supports 1000+ tasks in dependency graph
- [ ] All 4 views functional (Plan Editor, Dashboard, Task Graph, Logs)
- [ ] All keyboard shortcuts working
- [ ] Pause/resume/cancel controls functional
- [ ] Session persistence working
- [ ] Documentation complete (800+ lines)
- [ ] 5 example work plans provided
- [ ] No lint warnings
- [ ] No compile errors
- [ ] Mode registered in registry
- [ ] CI/CD passing

---

## Estimated Time Breakdown

| Stream | Days | Agent | Lines of Code |
|--------|------|-------|---------------|
| A: Foundation | 5 | Haiku | ~1,200 |
| B: Session | 2 | Haiku | ~450 |
| C: Editor | 1 | Haiku | ~550 |
| D: Graph | 2 | Sonnet | ~750 |
| E: Log | 1 | Haiku | ~450 |
| F: Integration | 4 | Sonnet | ~900 |
| G: Features | 1 | Haiku | ~100 |
| H: Testing | 1 | Sonnet | ~1,000 |
| I: Docs | 1 | Haiku | ~800 |
| **Total** | **15** | **Mixed** | **~6,200** |

---

## Risk Mitigation

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| mnemosyne orchestrate API changes | Medium | High | Version lock, extensive error handling |
| Force layout performance issues | Low | Medium | Reuse proven code from analyze mode |
| Event parsing complexity | Medium | Medium | Define clear schema, extensive tests |
| Integration delays | Low | High | Reserve 4 days for integration (Days 9-12) |
| Test failures on CI | Medium | Low | Run tests locally before commit |

---

**Last Updated**: 2025-11-09
**Status**: Ready for Phase 4 (Implementation)
**Next Action**: Launch parallel agents for Week 1 work
