# Phase 8: Refinement & Polish - Execution Plan

**Date**: 2025-11-09
**Phase**: 8 Execution Plan
**Timeline**: 7 days (can be compressed to 5 if needed)

## Overview

This plan executes Phase 8 using parallel agents to maximize efficiency. All work is parallelizable with minimal conflict risk.

## Agent Execution Strategy

### Days 1-2: Foundation Testing (4 Parallel Agents)

Launch 4 agents simultaneously to improve test coverage in priority packages.

---

#### Agent 1 (Sonnet): `mnemosyne` Package Testing
**Model**: Sonnet (complex retry logic, connection pools)
**Target**: 52.2% → 75% (+22.8%)

**Task Description**:
```markdown
You are Agent 1 improving test coverage for internal/mnemosyne package.

**Current Coverage**: 52.2%
**Target Coverage**: 75%+
**Gap**: +22.8% (need ~25-30 new tests)

**Files to Test** (priority order):
1. client.go - mnemosyne client wrapper, connection management
2. connection_pool.go - Connection pooling, health checks
3. recall.go - Semantic search queries
4. remember.go - Memory storage operations
5. evolve.go - Memory evolution and consolidation

**Test Scenarios (Required)**:
- Connection failure and automatic retry logic
- Connection pool exhaustion handling
- Concurrent recall operations (race detector)
- Malformed server responses
- Network timeouts
- Health check failures
- Memory evolution edge cases

**Guidelines**:
- Use table-driven tests
- Test both success and error paths
- Add concurrency tests (use t.Parallel())
- Mock HTTP responses using httptest
- Test retry backoff behavior
- Verify context cancellation

**Deliverable**:
- Create/update internal/mnemosyne/*_test.go
- 25-30 new passing tests
- Coverage report showing 75%+ coverage
- Commit with message format:
  "Improve mnemosyne package test coverage to 75%+ (Agent 1)"

**Do NOT**:
- Modify production code (client.go, etc.)
- Add new features
- Refactor existing tests
```

**Expected Output**: 25-30 tests, coverage 75%+

---

#### Agent 2 (Haiku): `modes` Package Testing
**Model**: Haiku (simpler mode switching logic)
**Target**: 54.7% → 75% (+20.3%)

**Task Description**:
```markdown
You are Agent 2 improving test coverage for internal/modes package.

**Current Coverage**: 54.7%
**Target Coverage**: 75%+
**Gap**: +20.3% (need ~15-20 new tests)

**Files to Test**:
1. registry.go - Mode registration and switching
2. explore_mode.go - Explore mode implementation (stub)

**Test Scenarios (Required)**:
- Mode registration (happy path + duplicate registration)
- Mode retrieval (existing + non-existent)
- Mode switching (OnExit → OnEnter lifecycle)
- Previous mode navigation (SwitchToPrevious)
- Concurrent mode operations
- Invalid mode IDs
- Mode count and AllModes

**Guidelines**:
- Use table-driven tests
- Test mode lifecycle (Init, OnEnter, OnExit)
- Test mode registry state management
- Verify OnExit called before OnEnter
- Test concurrent mode switches

**Deliverable**:
- Create/update internal/modes/*_test.go
- 15-20 new passing tests
- Coverage report showing 75%+ coverage
- Commit with message format:
  "Improve modes package test coverage to 75%+ (Agent 2)"

**Do NOT**:
- Modify mode implementations
- Add new modes
- Change mode registry API
```

**Expected Output**: 15-20 tests, coverage 75%+

---

#### Agent 3 (Haiku): `layout` Package Testing
**Model**: Haiku (layout engine logic)
**Target**: 54.5% → 75% (+20.5%)

**Task Description**:
```markdown
You are Agent 3 improving test coverage for internal/layout package.

**Current Coverage**: 54.5%
**Target Coverage**: 75%+
**Gap**: +20.5% (need ~20-25 new tests)

**Files to Test**:
1. engine.go - Layout engine lifecycle
2. layouts.go - Layout types (Standard, Split, etc.)
3. component.go - Component management
4. focus.go - Focus navigation (Tab/Shift+Tab)

**Test Scenarios (Required)**:
- Layout switching (Standard ↔ Split)
- Component registration and rendering
- Focus traversal (Tab wraps around)
- Window resize handling
- Component visibility toggling
- Multiple components in layout
- Focus navigation edge cases (no components, single component)

**Guidelines**:
- Use table-driven tests
- Test all layout types
- Test focus wrap-around behavior
- Verify component rendering order
- Test window resize propagation

**Deliverable**:
- Create/update internal/layout/*_test.go
- 20-25 new passing tests
- Coverage report showing 75%+ coverage
- Commit with message format:
  "Improve layout package test coverage to 75%+ (Agent 3)"

**Do NOT**:
- Modify layout engine logic
- Add new layout types
- Change component API
```

**Expected Output**: 20-25 tests, coverage 75%+

---

#### Agent 4 (Haiku): `app/events` Package Testing
**Model**: Haiku (event broker logic)
**Target**: 53.0% → 75% (+22.0%)

**Task Description**:
```markdown
You are Agent 4 improving test coverage for internal/app/events package.

**Current Coverage**: 53.0%
**Target Coverage**: 75%+
**Gap**: +22.0% (need ~15-20 new tests)

**Files to Test**:
1. broker.go - Event broker, subscriptions
2. types.go - Event types and handlers

**Test Scenarios (Required)**:
- Subscription management (subscribe + unsubscribe)
- Event publishing and delivery
- Multiple subscribers
- Subscriber removal during publish
- Concurrent publish/subscribe (race detector)
- Event queue overflow
- Nil handler handling
- Event delivery ordering

**Guidelines**:
- Use table-driven tests
- Test concurrent operations (t.Parallel())
- Verify event delivery order
- Test subscriber isolation (one subscriber error doesn't affect others)
- Test cleanup (unsubscribe removes handler)

**Deliverable**:
- Create/update internal/app/events/*_test.go
- 15-20 new passing tests
- Coverage report showing 75%+ coverage
- Commit with message format:
  "Improve app/events package test coverage to 75%+ (Agent 4)"

**Do NOT**:
- Modify broker implementation
- Add new event types
- Change event API
```

**Expected Output**: 15-20 tests, coverage 75%+

---

### Days 2-3: Documentation (2 Parallel Agents)

After testing foundation is solid, launch documentation agents.

---

#### Agent 5 (Haiku): Edit Mode Guide
**Model**: Haiku (documentation writing)
**Target**: 600+ lines

**Task Description**:
```markdown
You are Agent 5 writing the comprehensive Edit Mode user guide.

**Target**: 600+ lines in docs/edit-mode-guide.md

**Required Sections**:
1. Overview and Features
2. Getting Started (entering Edit mode)
3. Keyboard Shortcuts Reference (complete table)
4. Buffer Management
   - Multi-buffer editing
   - Buffer switching
   - Buffer creation/deletion
5. Search Functionality
   - Regex search
   - Case sensitivity
   - Search navigation
6. Semantic Analysis Integration
   - Entity extraction in editor
   - Context panel usage
7. Syntax Highlighting
8. Terminal Integration
9. Example Workflows
   - Writing code with semantic analysis
   - Multi-file editing
   - Search and replace
10. Troubleshooting

**Style**:
- Clear, concise language
- Code examples for each feature
- Screenshots (ASCII art representations)
- Practical workflows
- FAQ section

**Deliverable**:
- docs/edit-mode-guide.md (600+ lines)
- Well-structured with ToC
- Practical examples
- Commit with message format:
  "Add comprehensive Edit Mode user guide (Agent 5)"
```

**Expected Output**: Complete user guide, 600+ lines

---

#### Agent 6 (Haiku): Analyze Mode Guide
**Model**: Haiku (documentation writing)
**Target**: 600+ lines

**Task Description**:
```markdown
You are Agent 6 writing the comprehensive Analyze Mode user guide.

**Target**: 600+ lines in docs/analyze-mode-guide.md

**Required Sections**:
1. Overview and Triple Analysis
2. Getting Started (entering Analyze mode)
3. Entity Extraction
   - Person, Organization, Location
   - Custom entity types
4. Relationship Identification
   - Triple structure (Subject-Predicate-Object)
   - Relationship types
5. Graph Visualization
   - Force-directed layout
   - Navigation (pan/zoom)
   - Node/edge styling
6. Filtering and Search
7. Export Capabilities
   - JSON format
   - CSV format
   - GraphML format
8. Integration with mnemosyne
9. Example Analyses
   - Code analysis
   - Document analysis
   - Conversation analysis
10. Best Practices and Troubleshooting

**Style**:
- Clear, concise language
- Examples for each feature
- Visual representations (ASCII art graphs)
- Practical use cases
- FAQ section

**Deliverable**:
- docs/analyze-mode-guide.md (600+ lines)
- Well-structured with ToC
- Practical examples
- Commit with message format:
  "Add comprehensive Analyze Mode user guide (Agent 6)"
```

**Expected Output**: Complete user guide, 600+ lines

---

### Days 3-5: Performance & Integration (2 Parallel Agents)

Run performance optimization and integration testing simultaneously.

---

#### Agent 7 (Sonnet): Performance Optimization
**Model**: Sonnet (complex profiling and optimization)
**Target**: 15+ benchmarks, performance improvements

**Task Description**:
```markdown
You are Agent 7 profiling and optimizing performance bottlenecks.

**Tasks**:
1. CPU Profiling
   - Profile graph layout algorithms
   - Profile event processing
   - Profile semantic analysis
   - Identify hot paths

2. Memory Profiling
   - Identify allocation hot spots
   - Reduce allocations in event loops
   - Optimize string operations

3. Create Benchmark Suite
   - Graph layout benchmarks (4+)
   - Event processing benchmarks (3+)
   - Semantic analysis benchmarks (3+)
   - Memory operation benchmarks (3+)
   - JSON parsing benchmarks (2+)

4. Implement Optimizations
   - Optimize identified bottlenecks
   - Document before/after metrics
   - Ensure no regressions (all tests still pass)

5. Documentation
   - Write docs/PERFORMANCE.md
   - Document profiling methodology
   - Document optimization results
   - Add performance regression tests

**Performance Targets**:
- Dashboard updates: 55-60 FPS (currently 30-40 FPS)
- Event processing: < 5ms latency (currently 10-15ms)
- Graph layout: < 100ms for 100 nodes (currently 200-300ms)
- Semantic analysis: < 50ms for 1000 words (currently 100-150ms)

**Deliverable**:
- 15+ benchmarks in *_bench_test.go files
- Optimizations applied
- docs/PERFORMANCE.md (300+ lines)
- Before/after metrics documented
- Commit with message format:
  "Add performance benchmarks and optimizations (Agent 7)"

**Guidelines**:
- Use pprof for profiling
- Benchmark before optimizing
- Verify no regressions
- Document all changes
```

**Expected Output**: 15+ benchmarks, measurable improvements, comprehensive docs

---

#### Agent 8 (Haiku): Integration Testing Framework
**Model**: Haiku (integration testing)
**Target**: 20+ integration tests

**Task Description**:
```markdown
You are Agent 8 creating the integration test framework.

**Tasks**:
1. Create Integration Test Package
   - internal/integration/ directory
   - Test helpers and utilities
   - Mock mnemosyne server

2. Test Scenarios (20+ tests)
   - Edit → Analyze workflow (3 tests)
   - Analyze → Orchestrate workflow (3 tests)
   - Session persistence across restarts (3 tests)
   - Error recovery (crash, corrupt data) (4 tests)
   - Large dataset handling (1000+ entities) (3 tests)
   - Concurrent mode operations (4 tests)

3. Test Fixtures
   - Example content files
   - Sample work plans
   - Mock data generators

4. Documentation
   - Write docs/TESTING.md
   - Document integration test strategy
   - Document test fixtures
   - Add CI/CD integration examples

**Deliverable**:
- internal/integration/ package (20+ tests)
- Test fixtures (data files)
- docs/TESTING.md (300+ lines)
- All tests passing
- Commit with message format:
  "Add integration test framework with 20+ E2E tests (Agent 8)"

**Guidelines**:
- Use httptest for mock servers
- Use t.TempDir() for test isolation
- Clean up resources (defer cleanup)
- Make tests deterministic
```

**Expected Output**: Integration test framework, 20+ tests, complete docs

---

### Days 5-7: Polish & Architecture (2 Parallel Agents)

Final polish and architecture documentation.

---

#### Agent 9 (Haiku): UI/UX Polish
**Model**: Haiku (UI improvements)
**Target**: Style guide, consistent UX, 50+ improvements

**Task Description**:
```markdown
You are Agent 9 polishing UI/UX across all modes.

**Tasks**:
1. Create Style Guide
   - Write docs/STYLE_GUIDE.md
   - Define color palette
   - Define typography standards
   - Define spacing/padding standards
   - Document component patterns

2. UI Consistency Audit
   - Audit all modes for inconsistencies
   - Standardize colors across modes
   - Standardize spacing/padding
   - Ensure consistent error message format

3. Error Message Improvements
   - Review all error messages (50+)
   - Make messages clear and actionable
   - Add suggestions for resolution
   - Standardize error format

4. Contextual Help
   - Add help overlay for Edit Mode (`?` key)
   - Add help overlay for Analyze Mode
   - Update Orchestrate Mode help
   - Ensure all modes have help

5. Loading States
   - Add loading indicators to async operations
   - Add progress bars where appropriate
   - Ensure smooth transitions

**Deliverable**:
- docs/STYLE_GUIDE.md (200+ lines)
- UI improvements applied to all modes
- 50+ improved error messages
- Help overlays for all modes
- Commit with message format:
  "Polish UI/UX with style guide and consistency improvements (Agent 9)"

**Guidelines**:
- Use lipgloss for styling
- Maintain existing functionality
- Document all style decisions
```

**Expected Output**: Style guide, polished UI, comprehensive help system

---

#### Agent 10 (Haiku): Architecture Documentation
**Model**: Haiku (technical writing)
**Target**: 1,200+ lines of architecture docs

**Task Description**:
```markdown
You are Agent 10 writing comprehensive architecture documentation.

**Documents to Create**:

1. docs/architecture.md (400+ lines)
   - System overview
   - Component architecture
   - Mode architecture
   - Data flow diagrams (Mermaid)
   - Package dependencies
   - Design patterns used

2. docs/DEVELOPMENT.md (500+ lines)
   - Developer onboarding
   - Project structure
   - Build instructions
   - Testing strategy
   - Debugging tips
   - Common tasks

3. docs/CONTRIBUTING.md (300+ lines)
   - How to contribute
   - Code style guidelines
   - PR process
   - Issue guidelines
   - Testing requirements

**Deliverable**:
- docs/architecture.md (400+ lines)
- docs/DEVELOPMENT.md (500+ lines)
- docs/CONTRIBUTING.md (300+ lines)
- Updated README.md with badges
- Commit with message format:
  "Add comprehensive architecture and development documentation (Agent 10)"

**Guidelines**:
- Use Mermaid for diagrams
- Include code examples
- Reference actual code paths
- Keep docs maintainable
```

**Expected Output**: Complete architecture and developer docs

---

## Execution Timeline

```
Day 1-2:  Agents 1-4  (Testing in parallel)       4 agents
Day 2-3:  Agents 5-6  (Documentation in parallel) 2 agents
Day 3-5:  Agents 7-8  (Perf + Integration)        2 agents
Day 5-7:  Agents 9-10 (Polish + Architecture)     2 agents
```

**Total Agent-Days**: 18
**Wall-Clock Days**: 7 (can be compressed to 5 with overlap)

---

## Safety & Conflict Avoidance

**Zero Conflict Guarantee**:
- Agents 1-4: Different packages (mnemosyne, modes, layout, app/events)
- Agents 5-6: Different files (edit-mode-guide.md, analyze-mode-guide.md)
- Agents 7-8: Different focus (benchmarks vs. integration tests)
- Agents 9-10: Code vs. docs (low conflict)

**Merge Strategy**:
- Each agent commits independently
- Sequential merges (no parallel PR merges)
- Test suite runs after each merge

---

## Success Criteria Verification

After all agents complete:

```bash
# 1. Coverage check
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total

# Expected: 75%+ total coverage

# 2. Test count
go test ./... -v | grep "^PASS" | wc -l

# Expected: 400+ tests

# 3. Benchmark check
go test ./... -bench=. -run=^$ | grep "^Benchmark"

# Expected: 15+ benchmarks

# 4. Documentation check
wc -l docs/*.md

# Expected: 4,500+ lines total

# 5. Build verification
go build ./...

# Expected: Clean build, zero warnings
```

---

## Rollback Plan

If any agent fails:

1. **Testing agents (1-4)**: Partial coverage improvement still valuable
2. **Documentation agents (5-6)**: Can complete sequentially if parallel fails
3. **Performance agent (7)**: Benchmarks valuable even without optimizations
4. **Integration agent (8)**: Framework valuable even with fewer tests
5. **Polish agents (9-10)**: Can be deferred to Phase 8.5 if needed

---

**Last Updated**: 2025-11-09
**Phase**: 8 Execution Plan
**Status**: Ready to Execute - Launch Agents 1-4
