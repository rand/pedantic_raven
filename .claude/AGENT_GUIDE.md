# Pedantic Raven: Agent Development Guide

> **Comprehensive guide for autonomous development agents (Claude Code, mnemosyne orchestration)**
>
> This guide provides decision trees, workflows, and templates for AI agents to work effectively on Pedantic Raven.

## Project Context

### What is Pedantic Raven?

**Purpose**: TUI memory editor for browsing, editing, and reasoning about mnemosyne's semantic memory graph.

**Core Technology Stack**:
- **Language**: Go 1.21+
- **UI Framework**: Bubble Tea (Elm Architecture)
- **Integration**: mnemosyne gRPC API, GLiNER REST API
- **Testing**: Go testing, table-driven tests, E2E workflows
- **Build**: Make, go build, docker-compose

**Current Phase**: Phase 4 (Link Mode) - 686 tests passing, 17 markdown documentation files.

### Development State

**Completed**:
- Phase 1: Basic editor (buffer, syntax, search)
- Phase 2: Memory integration (list, detail, mnemosyne client)
- Phase 3: Event system (PubSub, 40+ event types)
- Phase 4: Link mode (navigation, keyboard shortcuts, visual highlighting)
- GLiNER Integration (ML-based entity extraction with fallback)

**In Progress**:
- Phase 5: Context Mode (planned)
- Phase 6: Analyze Mode (specified)

**Version**: v0.4.x (Phase 4 development)

### Key Metrics

- **Codebase**: 24,333 lines of Go code
- **Tests**: 686 tests (668 core + 18 GLiNER E2E)
- **Test Coverage**: 70%+ overall
- **Documentation**: 17 markdown files
- **Dependencies**: 15 Go modules, 1 Python service (optional)

## Architecture Quick Reference

### Component Map

```
┌─────────────────────────────────────────────────────┐
│                   Main Application                   │
│              (internal/app/app.go)                   │
└───────────────────┬─────────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │     Event Broker      │
        │   (PubSub pattern)    │
        └───────────┬───────────┘
                    │
    ┌───────────────┼───────────────┐
    │               │               │
┌───▼────┐    ┌────▼─────┐    ┌───▼────┐
│  Edit  │    │   Link   │    │ Analyze│
│  Mode  │    │   Mode   │    │  Mode  │
└───┬────┘    └────┬─────┘    └───┬────┘
    │              │              │
    └──────┬───────┴──────┬───────┘
           │              │
    ┌──────▼──────┐  ┌───▼─────────┐
    │   Semantic  │  │   mnemosyne │
    │  Analyzer   │  │    Client   │
    └──────┬──────┘  └─────────────┘
           │
    ┌──────┴──────────┐
    │                 │
┌───▼────────┐  ┌────▼──────┐
│  Pattern   │  │  GLiNER   │
│ Extractor  │  │ Extractor │
└────────────┘  └───────────┘
```

### Data Flow

**User Input → Mode Update → Event Publish → Subscribers React → View Render**

Example: Entity extraction in Edit Mode
```
1. User types text
2. EditMode.Update() receives KeyMsg
3. Buffer updated
4. triggerAnalysis() called (debounced)
5. Analyzer.Analyze() runs in goroutine
6. ExtractEntities() via HybridExtractor
7. GLiNER attempts extraction
8. Falls back to Pattern if unavailable
9. Results published via AnalysisComplete event
10. UI updates with highlighted entities
```

### Key Interfaces

**Bubble Tea Model**:
```go
type Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```

**Entity Extractor**:
```go
type EntityExtractor interface {
    ExtractEntities(ctx context.Context, text string, entityTypes []string) ([]Entity, error)
    Name() string
    IsAvailable(ctx context.Context) bool
}
```

**Event System**:
```go
type Event interface {
    Type() string
}

type Broker interface {
    Publish(e Event)
    Subscribe(eventType string, handler func(Event))
}
```

## Autonomous Work Guidelines

### Decision Tree: Starting Work

```
New task received
  │
  ├─ Is it a bug fix?
  │   ├─ Yes: git checkout -b fix/description
  │   └─ Run tests to reproduce → Fix → Write regression test → Verify
  │
  ├─ Is it a new feature?
  │   ├─ Yes: git checkout -b feature/description
  │   └─ Research existing patterns → Design → Implement → Test → Document
  │
  ├─ Is it documentation?
  │   ├─ Yes: git checkout -b docs/description
  │   └─ Review existing docs → Update/create → Ensure consistency
  │
  ├─ Is it refactoring?
  │   ├─ Yes: git checkout -b refactor/description
  │   └─ Ensure tests pass → Refactor → Verify tests still pass → No behavior change
  │
  └─ Is it a release?
      └─ Yes: Follow release process (see Release Management section)
```

### Decision Tree: Code Changes

```
Before writing code
  │
  ├─ Read relevant existing code
  │   └─ Use Read tool on files in scope
  │
  ├─ Understand patterns
  │   ├─ Table-driven tests? → Follow pattern
  │   ├─ Event-driven? → Use broker
  │   ├─ Strategy pattern? → Implement interface
  │   └─ Graceful degradation? → Handle errors gracefully
  │
  ├─ Check for existing tests
  │   └─ Add test cases, don't replace existing
  │
  └─ Plan implementation
      ├─ Identify affected files
      ├─ Plan test strategy
      └─ Consider backward compatibility
```

### Decision Tree: Testing

```
After code changes
  │
  ├─ What type of tests?
  │   ├─ New function? → Unit tests
  │   ├─ Module boundary? → Integration tests
  │   ├─ Full workflow? → E2E tests
  │   └─ All of the above? → Layered approach
  │
  ├─ Write tests BEFORE running
  │   └─ Use table-driven pattern for multiple cases
  │
  ├─ Commit changes
  │   └─ git add . && git commit -m "Description"
  │
  ├─ Run tests
  │   ├─ go test ./... (all tests)
  │   ├─ go test -short ./... (fast tests)
  │   └─ go test ./internal/editor (specific package)
  │
  ├─ Tests pass?
  │   ├─ Yes: Verify coverage, update docs
  │   └─ No: Fix → Commit fix → Re-test
  │
  └─ Quality gates
      ├─ [ ] All tests pass
      ├─ [ ] Coverage targets met
      ├─ [ ] go vet passes
      ├─ [ ] Documentation updated
      └─ [ ] No TODO/FIXME comments
```

### Decision Tree: External Service Integration

```
Integrating with external service (mnemosyne, GLiNER)?
  │
  ├─ Is service available?
  │   ├─ Check with health endpoint or IsEnabled()
  │   └─ If unavailable, implement graceful fallback
  │
  ├─ Design interface
  │   └─ Abstract service behind Go interface
  │
  ├─ Implement client
  │   ├─ HTTP client? → Use http.Client with timeout
  │   ├─ gRPC client? → Use generated stubs
  │   └─ Error handling → Wrap errors with context
  │
  ├─ Add retry logic
  │   ├─ Exponential backoff (100ms, 200ms, 400ms)
  │   ├─ Fresh request on each attempt (important!)
  │   └─ Context-aware (respect cancellation)
  │
  ├─ Configuration
  │   ├─ Add to config.toml
  │   ├─ Support environment variable override
  │   └─ Document in CLAUDE.md
  │
  └─ Testing
      ├─ Mock service (httptest.Server for fast tests)
      ├─ Real service (docker-compose for integration)
      └─ Fallback behavior (test unavailability)
```

## Repository Organization

### Tidying Protocol

**CRITICAL**: Always preserve git history and references.

```
Before tidying
  │
  ├─ Identify files to organize
  │   ├─ Deprecated code? → Move to archive/
  │   ├─ Outdated docs? → Update or move to archive/docs/
  │   └─ Temporary files? → Check if still needed
  │
  ├─ Check git history
  │   └─ git log --follow <file>
  │   └─ If significant history, preserve via git mv
  │
  ├─ Find all references
  │   ├─ grep -r "filename" .
  │   ├─ Check documentation (README, docs/, .claude/)
  │   └─ Check import statements (for Go files)
  │
  ├─ Execute moves
  │   ├─ Use git mv (preserves history)
  │   ├─ mkdir -p archive/category
  │   └─ git mv old/path archive/category/
  │
  ├─ Update references
  │   ├─ Documentation links
  │   ├─ Import paths
  │   └─ README.md table of contents
  │
  └─ Commit with descriptive message
      └─ git commit -m "Archive deprecated X, update references"
```

**Example: Archiving deprecated mode**:
```bash
# 1. Create archive structure
mkdir -p archive/deprecated-modes

# 2. Move with git mv (preserves history)
git mv internal/editor/old_mode.go archive/deprecated-modes/
git mv internal/editor/old_mode_test.go archive/deprecated-modes/

# 3. Find and update references
grep -r "old_mode" . --exclude-dir=archive
# Update any references found

# 4. Document in commit
git commit -m "Archive old_mode (replaced by new_mode in v0.4.0)

- Moved to archive/deprecated-modes/
- Updated references in docs/ARCHITECTURE.md
- No functional changes"
```

### Directory Structure Guidelines

**Go code organization**:
```
internal/
  <domain>/           # Domain-specific package
    <domain>.go       # Main logic
    <domain>_test.go  # Tests
    types.go          # Type definitions (if complex)
    errors.go         # Error definitions (if many)
```

**Documentation organization**:
```
README.md             # High-level overview, quick start
.claude/              # Developer guides (this file)
docs/                 # Detailed specifications
  ARCHITECTURE.md     # System architecture
  PHASE_*.md          # Phase specifications
  INTEGRATION_*.md    # Integration guides
CHANGELOG.md          # Release history
```

## Documentation Updates

### When to Update Documentation

**Always update when**:
- Adding new feature → Update README, relevant docs/, add examples
- Changing architecture → Update docs/ARCHITECTURE.md, .claude/CLAUDE.md
- Fixing bug → Update CHANGELOG.md, add note if user-facing
- Changing API → Update docs/, add migration guide if breaking
- Adding configuration → Update config.toml, CLAUDE.md, README
- Creating release → Update CHANGELOG.md, README version

### Documentation Update Protocol

```
Code change committed
  │
  ├─ Identify affected docs
  │   ├─ README.md? (if user-facing change)
  │   ├─ docs/ARCHITECTURE.md? (if architectural change)
  │   ├─ .claude/CLAUDE.md? (if workflow/pattern change)
  │   ├─ .claude/AGENT_GUIDE.md? (if agent workflow change)
  │   ├─ Relevant docs/INTEGRATION_*.md?
  │   └─ CHANGELOG.md? (always for releases)
  │
  ├─ Update each document
  │   ├─ Add new sections if needed
  │   ├─ Update existing sections
  │   ├─ Update examples/code snippets
  │   └─ Update version numbers if applicable
  │
  ├─ Check consistency
  │   ├─ Same terminology across all docs?
  │   ├─ Links still valid?
  │   └─ Table of contents updated?
  │
  ├─ Commit docs separately (optional) or with code
  │   └─ git commit -m "Update docs for X feature"
  │
  └─ Store in mnemosyne
      └─ mnemosyne remember -c "Updated docs for X" -n "project:pedantic_raven"
```

### Documentation Templates

**New integration guide** (docs/INTEGRATION_X.md):
```markdown
# X Integration

## Overview
[What is X, why integrate it]

## Architecture
[How it fits into Pedantic Raven]

## Configuration
[config.toml settings, environment variables]

## API Reference
[Key functions, types, interfaces]

## Usage Examples
[Code snippets showing common use cases]

## Testing
[How to test this integration]

## Troubleshooting
[Common errors and solutions]
```

**README.md structure**:
```markdown
# Pedantic Raven

[One-line description]

## Features
[Bullet list of main features]

## Quick Start
[Installation and basic usage]

## Documentation
[Links to detailed docs]

## Development
[Building, testing, contributing]

## License
[License information]
```

## Release Management

### Semantic Versioning Strategy

**Phase-based versioning** (current state: v0.4.x):

```
v0.X.Y
  │
  ├─ X = Phase number (0.4 = Phase 4)
  │   └─ Increment when phase milestone complete
  │
  └─ Y = Feature/fix within phase
      └─ Increment for features, bug fixes, improvements
```

**Milestone releases**:
- v0.4.0: Phase 4 (Link Mode) complete
- v0.5.0: Phase 5 (Context Mode) complete
- v0.6.0: Phase 6 (Analyze Mode) complete
- v1.0.0: All phases complete, production-ready

### Release Creation Protocol

```
Ready to create release
  │
  ├─ Pre-release checks
  │   ├─ [ ] All tests passing (go test ./...)
  │   ├─ [ ] No uncommitted changes
  │   ├─ [ ] Documentation up to date
  │   ├─ [ ] CHANGELOG.md updated
  │   └─ [ ] Version bumped in code (if applicable)
  │
  ├─ Determine version number
  │   ├─ Phase milestone? → v0.X.0
  │   ├─ Feature addition? → v0.X.Y (increment Y)
  │   ├─ Bug fix? → v0.X.Y (increment Y)
  │   └─ Breaking change before v1.0? → v0.X.0 (increment X)
  │
  ├─ Update CHANGELOG.md
  │   ├─ Add section: ## vX.Y.Z - YYYY-MM-DD
  │   ├─ List changes: ### Added, ### Changed, ### Fixed
  │   └─ Commit: git commit -m "Prepare release vX.Y.Z"
  │
  ├─ Create git tag
  │   ├─ git tag -a vX.Y.Z -m "Version X.Y.Z: Brief description"
  │   └─ git push origin vX.Y.Z
  │
  ├─ Create GitHub release
  │   ├─ gh release create vX.Y.Z --title "vX.Y.Z: Title" --notes "Changelog excerpt"
  │   └─ Or use GitHub web UI
  │
  ├─ Optional: Build and attach binaries
  │   ├─ make build
  │   ├─ tar -czf pedantic_raven-vX.Y.Z-linux-amd64.tar.gz pedantic_raven
  │   └─ gh release upload vX.Y.Z pedantic_raven-vX.Y.Z-*.tar.gz
  │
  └─ Store in mnemosyne
      └─ mnemosyne remember -c "Released vX.Y.Z: Key learnings" -n "project:pedantic_raven"
```

### CHANGELOG.md Format

```markdown
# Changelog

All notable changes to Pedantic Raven will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Feature X with capability Y

### Changed
- Refactored Z for better performance

### Fixed
- Bug in W causing incorrect behavior

## [0.4.3] - 2024-01-15

### Added
- Link navigation with keyboard shortcuts (Tab, Shift+Tab, Enter)
- Visual highlighting for selected links
- Footer display showing link navigation state

### Fixed
- TestViewWithSelectedLink now checks footer instead of indicator

## [0.4.2] - 2024-01-10
...
```

### Release Checklist

Before creating release:
```
[ ] All tests passing locally (go test ./...)
[ ] CI/CD passing (GitHub Actions green)
[ ] Documentation updated (README, docs/, .claude/)
[ ] CHANGELOG.md updated with all changes
[ ] Version number determined (v0.X.Y)
[ ] No TODO/FIXME/XXX comments (create issues instead)
[ ] All quality gates passed
[ ] Breaking changes documented (if any)
[ ] Migration guide written (if breaking changes)
```

After creating release:
```
[ ] Git tag created and pushed
[ ] GitHub release created with changelog
[ ] Binaries built and attached (if distributing)
[ ] mnemosyne memory stored with learnings
[ ] Team notified (if applicable)
[ ] README.md version badge updated (if exists)
```

## Code Generation Templates

### New Mode Template

```go
package editor

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/rand/pedantic-raven/internal/editor/buffer"
)

// MyMode implements a new editor mode.
type MyMode struct {
    editor *buffer.Buffer
    // Add mode-specific state here
}

// NewMyMode creates a new instance of MyMode.
func NewMyMode() *MyMode {
    return &MyMode{
        editor: buffer.NewBuffer(),
    }
}

// Init initializes the mode.
func (m *MyMode) Init() tea.Cmd {
    return nil
}

// Update handles messages and updates mode state.
func (m *MyMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "ctrl+q":
            return m, tea.Quit
        case "esc":
            // Return to previous mode
            return m, nil
        default:
            // Handle mode-specific keys
        }
    }
    return m, nil
}

// View renders the mode.
func (m *MyMode) View() string {
    return m.editor.View()
}
```

### Test Template (Table-Driven)

```go
package editor

import "testing"

func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    invalidInput,
            expected: emptyOutput,
            wantErr:  true,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("MyFunction() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### Integration Test Template

```go
package editor

import (
    "context"
    "testing"
    "time"
)

func TestIntegration_FeatureX(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Initialize components
    component1 := setupComponent1(t)
    defer component1.Cleanup()

    component2 := setupComponent2(t)
    defer component2.Cleanup()

    // Test integration
    result, err := testInteraction(ctx, component1, component2)
    if err != nil {
        t.Fatalf("Integration test failed: %v", err)
    }

    // Verify result
    if !validateResult(result) {
        t.Errorf("Result validation failed: %+v", result)
    }
}
```

## Agent Workflow Examples

### Example 1: Adding GLiNER Integration

**Task**: "Integrate GLiNER for ML-based entity extraction"

**Workflow**:
```
1. Research
   └─ Read GLiNER docs, understand model capabilities
   └─ mnemosyne recall -q "entity extraction patterns"

2. Design
   └─ Create EntityExtractor interface
   └─ Plan HybridExtractor with fallback
   └─ Design configuration system

3. Implement
   ├─ Python service (services/gliner/)
   ├─ Go client (internal/gliner/)
   ├─ Extractor implementations (internal/editor/semantic/)
   └─ Configuration (config.toml, internal/config/)

4. Test
   ├─ Unit tests (client, extractors)
   ├─ Integration tests (docker-compose)
   └─ E2E tests (full workflow)

5. Document
   ├─ docs/GLINER_INTEGRATION.md
   ├─ Update README.md
   └─ Update .claude/CLAUDE.md

6. Release
   └─ Create PR with comprehensive description
   └─ mnemosyne remember -c "GLiNER integration patterns"
```

### Example 2: Fixing Critical Bug

**Task**: "HTTP retry logic sends empty request bodies"

**Workflow**:
```
1. Reproduce
   └─ Write failing test (TestRetryLogic)
   └─ Verify bug exists

2. Debug
   └─ Read internal/gliner/client.go:152-174
   └─ Identify: http.Request body consumed after first attempt

3. Fix
   └─ Create fresh request with new bytes.NewBuffer() each retry
   └─ Commit: "Fix retry logic to create fresh request on each attempt"

4. Verify
   └─ Run tests (go test ./internal/gliner)
   └─ Run E2E tests (go test ./internal/editor -run E2E)

5. Document
   └─ Update CHANGELOG.md
   └─ Add comment explaining retry logic

6. mnemosyne
   └─ mnemosyne remember -c "HTTP request bodies are consumed: create fresh request per retry"
```

### Example 3: Creating Documentation

**Task**: "Create project CLAUDE.md and AGENT_GUIDE.md"

**Workflow**:
```
1. Research
   ├─ Read global CLAUDE.md
   ├─ Read existing project docs
   └─ Understand project patterns

2. Structure
   ├─ Plan CLAUDE.md outline (Quick Start, Architecture, Workflow, etc.)
   └─ Plan AGENT_GUIDE.md outline (Decision trees, templates, etc.)

3. Write
   ├─ Create .claude/CLAUDE.md (300-400 lines, context-efficient)
   └─ Create .claude/AGENT_GUIDE.md (400-500 lines, comprehensive)

4. Review
   ├─ Check consistency with existing docs
   ├─ Verify all workflows covered
   └─ Ensure decision trees are complete

5. Update
   └─ Update README.md with documentation links

6. Store
   └─ mnemosyne remember -c "Project documentation patterns for TUI + gRPC architecture"
```

## mnemosyne Integration

### Storing Project Memories

**When to store**:
- Architecture decisions
- Difficult bugs and solutions
- Performance optimizations
- Release learnings
- Common pitfalls

**How to store**:
```bash
# Architecture decision
mnemosyne remember -c "Decision: Use strategy pattern for entity extraction. Rationale: Allows pluggable implementations (Pattern, GLiNER, Hybrid) with graceful fallback. Proven effective in production." \
  -n "project:pedantic_raven" -i 9 -t "architecture,patterns,entity-extraction"

# Bug fix
mnemosyne remember -c "Bug: HTTP request bodies are consumed after first read. Solution: Create fresh http.Request with new bytes.NewBuffer(data) on each retry attempt. Critical for retry logic." \
  -n "project:pedantic_raven" -i 8 -t "bugs,http,retry-logic"

# Performance insight
mnemosyne remember -c "GLiNER extraction: 100-300ms latency. Pattern matching: <1ms. Hybrid approach with fallback provides best UX: fast when GLiNER available, always functional." \
  -n "project:pedantic_raven" -i 7 -t "performance,gliner,entity-extraction"

# Release learning
mnemosyne remember -c "v0.4.3: Link navigation significantly improved UX. Users praised keyboard shortcuts (Tab/Shift+Tab/Enter). Consider similar patterns for other modes." \
  -n "project:pedantic_raven" -i 7 -t "release,ux,keyboard-shortcuts"
```

### Recalling Project Knowledge

```bash
# Before starting work
mnemosyne recall -q "entity extraction architecture" -n "project:pedantic_raven" -l 5

# When debugging
mnemosyne recall -q "retry logic http bugs" -n "project:pedantic_raven" -l 3

# Before release
mnemosyne recall -q "release checklist learnings" -n "project:pedantic_raven" -l 5

# Cross-project patterns
mnemosyne recall -q "graceful degradation patterns" --min-importance 7
```

## Quick Decision Reference

**Should I create a new file or modify existing?**
→ Modify existing unless new domain or clear separation needed.

**Should I write tests first or after implementation?**
→ After implementation but BEFORE running tests.

**Should I commit before or after running tests?**
→ BEFORE running tests (critical for debugging).

**Should I update docs in same commit or separately?**
→ Same commit if small change, separate if major doc overhaul.

**Should I create a GitHub issue or mnemosyne memory?**
→ Issue for trackable work, memory for knowledge/patterns.

**Should I use pattern matching or GLiNER?**
→ Hybrid approach with automatic fallback (already implemented).

**Should I add environment variable or TOML config?**
→ Both: TOML for defaults, env var for override.

**Should I archive or delete deprecated code?**
→ Archive (git mv to archive/) to preserve history.

**Should I update README or create new doc?**
→ README for high-level changes, new doc for detailed specs.

**Should I create v0.X.0 or v0.X.Y release?**
→ v0.X.0 for phase milestones, v0.X.Y for features/fixes.

## Anti-Patterns to Avoid

**Code**:
- ❌ Reusing http.Request in retry loops (body consumed)
- ❌ Ignoring context.Context (no cancellation/timeout)
- ❌ Panicking instead of returning errors
- ❌ Hardcoding service URLs (use config)
- ❌ Not handling service unavailability (implement fallback)

**Testing**:
- ❌ Running tests before committing (debug stale code)
- ❌ Testing while still writing code (invalid results)
- ❌ Skipping E2E tests (miss integration bugs)
- ❌ Not using table-driven tests (Go idiom)

**Documentation**:
- ❌ Leaving docs outdated after code changes
- ❌ Not updating CHANGELOG.md for releases
- ❌ Documenting "what" instead of "why"
- ❌ Inconsistent terminology across docs

**Git**:
- ❌ Committing directly to main (use feature branches)
- ❌ Vague commit messages ("fix", "update")
- ❌ Deleting files instead of archiving (lose history)
- ❌ Not using git mv (breaks history tracking)

**Releases**:
- ❌ Releasing with failing tests
- ❌ Not updating version numbers
- ❌ Skipping CHANGELOG.md updates
- ❌ Not storing release learnings in mnemosyne

## Summary

This guide provides comprehensive workflows for autonomous development on Pedantic Raven. Key principles:

1. **Follow established patterns** (Bubble Tea, event-driven, strategy pattern)
2. **Test thoroughly** (unit, integration, E2E with table-driven tests)
3. **Document comprehensively** (code changes → doc updates)
4. **Preserve history** (git mv, archive instead of delete)
5. **Release systematically** (semantic versioning, CHANGELOG, mnemosyne)
6. **Store knowledge** (mnemosyne memories for decisions and learnings)

For questions or clarifications, refer to:
- `.claude/CLAUDE.md` for development workflows
- `docs/ARCHITECTURE.md` for system architecture
- `docs/PHASE_*.md` for phase specifications
- mnemosyne for project knowledge and patterns
