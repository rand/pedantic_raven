# Pedantic Raven Development Guide

> **Context-Efficient Guide for Human and AI Developers**
>
> Pedantic Raven is a TUI memory editor for browsing, editing, and reasoning about mnemosyne's semantic memory graph. Built with Go, Bubble Tea, and integrated with mnemosyne's gRPC API.

## Quick Start

```bash
# Clone and build
git clone https://github.com/rand/pedantic_raven.git
cd pedantic_raven
make build

# Run tests
make test                    # All tests
go test ./internal/editor    # Editor tests only
go test -short ./...         # Fast tests only

# Run the editor
./pedantic_raven

# With mnemosyne integration
MNEMOSYNE_ENABLED=true MNEMOSYNE_ADDRESS=localhost:50051 ./pedantic_raven
```

## Project Architecture

### Core Components

```
pedantic_raven/
├── internal/
│   ├── app/           # Main application and coordinator
│   ├── editor/        # Core editor logic and modes
│   │   ├── buffer/    # Text buffer operations
│   │   ├── search/    # Search functionality
│   │   ├── semantic/  # Entity extraction (pattern + GLiNER)
│   │   └── syntax/    # Syntax highlighting
│   ├── gliner/        # GLiNER ML-based NER client
│   ├── layout/        # UI layout management
│   ├── memorydetail/  # Memory detail view
│   ├── memorylist/    # Memory list view
│   ├── mnemosyne/     # mnemosyne gRPC integration
│   ├── modes/         # Mode registry and transitions
│   ├── overlay/       # Overlays (help, confirm, etc.)
│   ├── palette/       # Command palette
│   └── terminal/      # Terminal state management
└── services/gliner/   # Python GLiNER service (FastAPI)
```

### Architecture Patterns

**Elm Architecture (Bubble Tea)**:
```go
type Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```

**Event-Driven**:
- 40+ event types in `internal/app/events/`
- PubSub broker for decoupled communication
- Modes publish events, other components subscribe

**Strategy Pattern**:
- `EntityExtractor` interface with multiple implementations:
  - `PatternExtractor`: Keyword matching (60-70% accuracy, <1ms)
  - `GLiNERExtractor`: ML-based (85-95% accuracy, 100-300ms)
  - `HybridExtractor`: Automatic fallback

**Graceful Degradation**:
- GLiNER unavailable → Falls back to pattern matching
- mnemosyne unavailable → Local-only mode
- Always functional, never crashes from external failures

## Development Workflow

### Branch Strategy

```bash
# Feature development
git checkout -b feature/description

# Bug fixes
git checkout -b fix/description

# Refactoring
git checkout -b refactor/description

# Documentation
git checkout -b docs/description
```

### Commit Guidelines

**Good commits**:
```bash
git commit -m "Add link navigation with keyboard shortcuts"
git commit -m "Fix retry logic to create fresh request on each attempt"
git commit -m "Refactor entity extraction to strategy pattern"
```

**Commit structure**:
1. Logical unit of work
2. Descriptive message (imperative mood)
3. Include context in body for complex changes
4. No AI attribution unless explicitly requested

### Code Quality Gates

Before marking work complete:
- [ ] Intent satisfied
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] No TODO/FIXME/stub comments (create GitHub issues instead)
- [ ] No anti-patterns
- [ ] Type-safe (go vet passes)

## Testing Strategy

### Test Types

**Unit Tests** (fast, isolated):
```go
func TestBufferInsert(t *testing.T) {
    b := buffer.NewBuffer()
    b.Insert([]rune("hello"))
    if b.Content() != "hello" {
        t.Errorf("expected 'hello', got '%s'", b.Content())
    }
}
```

**Integration Tests** (module boundaries):
```go
func TestSemanticAnalyzerWithGLiNER(t *testing.T) {
    // Test GLiNER client + semantic analyzer integration
}
```

**E2E Tests** (full workflows):
```go
func TestGLiNERE2E_UserTypingExperience(t *testing.T) {
    // Simulate user typing and verify real-time analysis
}
```

**Table-Driven Tests** (Go idiom):
```go
tests := []struct {
    name     string
    input    string
    expected Result
}{
    {"case1", "input1", result1},
    {"case2", "input2", result2},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

### Coverage Targets

- Critical path: 90%+
- Business logic: 80%+
- UI layer: 60%+
- Overall: 70%+

### Running Tests

```bash
# All tests
make test
go test ./...

# Fast tests only (skip integration)
go test -short ./...

# Specific package
go test ./internal/editor
go test ./internal/gliner

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Verbose output
go test -v ./...

# Run specific test
go test -run TestName ./internal/editor
```

## Common Tasks

### Adding a New Mode

1. Create mode in `internal/editor/`:
```go
type MyMode struct {
    editor *buffer.Buffer
    // mode-specific state
}

func (m *MyMode) Init() tea.Cmd { return nil }
func (m *MyMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* logic */ }
func (m *MyMode) View() string { /* render */ }
```

2. Register in `internal/modes/registry.go`:
```go
registry.Register("my-mode", func() tea.Model {
    return NewMyMode()
})
```

3. Add mode transitions in mode switcher
4. Write tests in `my_mode_test.go`
5. Update documentation

### Adding Event Types

1. Define event in `internal/app/events/events.go`:
```go
type MyEvent struct {
    Field string
}

func (e MyEvent) Type() string { return "my-event" }
```

2. Publish event:
```go
broker.Publish(events.MyEvent{Field: "value"})
```

3. Subscribe in consumer:
```go
broker.Subscribe("my-event", func(e events.Event) {
    myEvent := e.(events.MyEvent)
    // Handle event
})
```

### Integrating with mnemosyne

```go
// Check if enabled
if !m.mnemosyneClient.IsEnabled() {
    return ErrMnemosyneDisabled
}

// Store memory
resp, err := m.mnemosyneClient.Store(ctx, &pb.StoreRequest{
    Content:    content,
    Namespace:  namespace,
    Importance: importance,
    Tags:       tags,
})

// Search memories
results, err := m.mnemosyneClient.Search(ctx, &pb.SearchRequest{
    Query: query,
    Limit: limit,
})

// Handle gracefully
if err != nil {
    if status.Code(err) == codes.Unavailable {
        // Fall back to local-only mode
    }
    return err
}
```

### Using GLiNER Integration

```go
// Initialize extractor
extractor := semantic.NewHybridExtractor(
    semantic.NewGLiNERExtractor(glinerClient),
    semantic.NewPatternExtractor(),
)

// Extract entities
entities, err := extractor.ExtractEntities(ctx, text, entityTypes)
if err == semantic.ErrExtractorUnavailable {
    // Automatically fell back to pattern matching
}

// Entity types
entityTypes := []string{
    "person",        // Alice, Bob
    "organization",  // Google, Microsoft
    "location",      // San Francisco
    "technology",    // Python, React
    "api_component", // POST, /api/users
}
```

### Configuration Management

Configuration via `config.toml` and environment variables:

```toml
# config.toml
[gliner]
enabled = true
service_url = "http://localhost:8765"
timeout = 5
max_retries = 2
fallback_to_pattern = true
score_threshold = 0.3

[mnemosyne]
enabled = true
address = "localhost:50051"
timeout = 10
```

Environment variables override TOML:
```bash
GLINER_ENABLED=false ./pedantic_raven
MNEMOSYNE_ADDRESS=remote:50051 ./pedantic_raven
```

## Repository Organization

### File Organization

**Go code**: Organized by domain (editor, mnemosyne, gliner, layout, etc.)
**Tests**: Co-located with code (`*_test.go`)
**Docs**: Root-level markdown + `docs/` for detailed specs
**Services**: External services (Python GLiNER service)

### Tidying Guidelines

**Non-destructive tidying**:
- Move deprecated files to `archive/` (preserve git history)
- Update references in documentation
- Never delete files with git history unless truly obsolete
- Always check `git log --follow` before moving

**Documentation tidying**:
- Keep README.md current (high-level overview)
- Move detailed specs to `docs/`
- Consolidate duplicate information
- Update references when moving docs

### Dependency Management

```bash
# Add dependency
go get github.com/package/name

# Update dependencies
go get -u ./...

# Tidy (remove unused)
go mod tidy

# Vendor (optional)
go mod vendor
```

## Release Process

### Semantic Versioning

**Phase-based versioning** (current: Phase 4):
- v0.4.x: Phase 4 (Link Mode) development
- v0.5.0: Phase 5 (Context Mode) milestone
- v0.6.0: Phase 6 (Analyze Mode) milestone
- v1.0.0: All phases complete, production-ready

**Increment rules**:
- v0.X.0: Phase milestone (new mode complete, all tests pass)
- v0.X.Y: Feature additions, bug fixes, improvements within phase
- v1.0.0: Full feature set, production-ready, stable API

### Creating Releases

```bash
# 1. Ensure all tests pass
make test
go test ./...

# 2. Update CHANGELOG.md
# Document all changes since last release

# 3. Tag release
git tag -a v0.4.3 -m "Phase 4.3: Link navigation with keyboard shortcuts"
git push origin v0.4.3

# 4. Create GitHub release
gh release create v0.4.3 \
  --title "v0.4.3: Link Navigation" \
  --notes "$(cat CHANGELOG.md | sed -n '/## v0.4.3/,/## v0.4.2/p')"

# 5. Build and attach binaries (optional)
make build
gh release upload v0.4.3 pedantic_raven
```

### Release Checklist

- [ ] All tests passing (go test ./...)
- [ ] Documentation updated (README, CHANGELOG, relevant docs/)
- [ ] Version bumped in main.go or VERSION file
- [ ] Git tag created with descriptive message
- [ ] GitHub release created with changelog excerpt
- [ ] Binaries built and attached (if distributing)
- [ ] mnemosyne memory stored for release learnings

## Quick Reference

### Key Files

```
main.go                          # Entry point
internal/app/app.go              # Main application coordinator
internal/editor/edit_mode.go     # Primary editor mode
internal/editor/semantic/        # Entity extraction
internal/mnemosyne/client.go     # mnemosyne gRPC client
internal/gliner/client.go        # GLiNER HTTP client
config.toml                      # Configuration template
```

### Keyboard Shortcuts (in Edit Mode)

```
Ctrl+S  - Save memory
Ctrl+Q  - Quit
Ctrl+P  - Command palette
Ctrl+F  - Search mode
Ctrl+L  - Link mode
Ctrl+A  - Analyze mode
Tab     - Next link (Link Mode)
Shift+Tab - Previous link (Link Mode)
Enter   - Follow link (Link Mode)
Esc     - Cancel/back
```

### Environment Variables

```bash
GLINER_ENABLED=true              # Enable GLiNER integration
GLINER_SERVICE_URL=...           # GLiNER service endpoint
MNEMOSYNE_ENABLED=true           # Enable mnemosyne integration
MNEMOSYNE_ADDRESS=localhost:50051 # mnemosyne gRPC address
LOG_LEVEL=debug                  # Logging level
```

### Build Commands

```bash
make build           # Build binary
make test            # Run all tests
make clean           # Clean build artifacts
make install         # Install to $GOPATH/bin
make docker-gliner   # Start GLiNER service
```

### Test Filters

```bash
-short               # Skip slow integration tests
-run TestName        # Run specific test
-v                   # Verbose output
-coverprofile=...    # Generate coverage report
-race                # Enable race detector
```

### Common Errors

**"mnemosyne unavailable"**: Start mnemosyne server or set MNEMOSYNE_ENABLED=false
**"GLiNER service unavailable"**: Start GLiNER service via docker-compose or set GLINER_ENABLED=false
**Test failures in gliner_e2e_test.go**: Ensure docker-compose is available or skip with `-short`

### Documentation Files

```
README.md                        # Project overview
.claude/CLAUDE.md                # This file (development guide)
.claude/AGENT_GUIDE.md           # Agentic development guide
docs/ARCHITECTURE.md             # System architecture
docs/PHASE_*.md                  # Phase specifications
docs/GLINER_INTEGRATION.md       # GLiNER integration details
CHANGELOG.md                     # Release history
```

### mnemosyne Integration

```bash
# Store project decisions
mnemosyne remember -c "Decision: Use strategy pattern for entity extraction" \
  -n "project:pedantic_raven" -i 9 -t "architecture,patterns"

# Recall relevant memories
mnemosyne recall -q "entity extraction patterns" -n "project:pedantic_raven"

# Store release learnings
mnemosyne remember -c "v0.4.3: Link navigation significantly improved UX" \
  -n "project:pedantic_raven" -i 7 -t "release,ux"
```

## Development Principles

1. **Graceful Degradation**: Always functional, even when external services unavailable
2. **Event-Driven**: Publish events, don't call directly across boundaries
3. **Strategy Pattern**: Pluggable implementations for flexibility
4. **Table-Driven Tests**: Go idiom for comprehensive test coverage
5. **Type Safety**: Leverage Go's type system, use interfaces
6. **Context Propagation**: Pass context.Context for cancellation and timeouts
7. **Non-Destructive**: Preserve git history, archive instead of delete
8. **Documentation**: Keep docs current, update on every significant change

## Getting Help

- GitHub Issues: https://github.com/rand/pedantic_raven/issues
- AGENT_GUIDE.md: Detailed guide for autonomous development
- docs/: Detailed specifications and architecture
- mnemosyne: Query project memories for past decisions
