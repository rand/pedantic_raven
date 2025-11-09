# Pedantic Raven Developer Guide

**Version**: 1.0
**Last Updated**: November 9, 2025
**Status**: Complete

## Table of Contents

1. [Getting Started](#getting-started)
2. [Prerequisites](#prerequisites)
3. [Project Structure](#project-structure)
4. [Build Instructions](#build-instructions)
5. [Running Tests](#running-tests)
6. [Debugging](#debugging)
7. [Common Tasks](#common-tasks)
8. [Code Style Guidelines](#code-style-guidelines)
9. [Performance Considerations](#performance-considerations)
10. [Troubleshooting](#troubleshooting)

---

## Getting Started

Welcome to Pedantic Raven development! This guide will help you set up your development environment and understand the codebase.

### Quick Start

```bash
# Clone repository
git clone https://github.com/rand/pedantic-raven.git
cd pedantic-raven

# Install dependencies
go mod download

# Build binary
go build -o pedantic_raven .

# Run application
./pedantic_raven

# Run tests
go test ./...
```

### Project Locations

- **Source Code**: `internal/` - All application packages
- **Main Entry**: `main.go` - Application entry point
- **Tests**: `*_test.go` files throughout codebase
- **Documentation**: `docs/` - User and developer guides
- **Proto Definitions**: `proto/` - Protocol buffer schemas
- **Configuration**: `config.toml` - Application config
- **Examples**: `examples/` - Example files and data

---

## Prerequisites

### Required

- **Go 1.25+** - [Download](https://golang.org/dl/)
  - Check version: `go version`
  - Should output: `go version go1.25.x`

- **git** - For version control
  - Check: `git --version`

### Optional But Recommended

- **mnemosyne Server** - For full memory integration
  - [Installation](https://github.com/rand/mnemosyne)
  - Start with: `mnemosyne serve`
  - Default: `localhost:50051`

- **Protocol Buffer Tools** - For proto regeneration
  ```bash
  # Install protoc compiler
  # macOS: brew install protobuf
  # Linux: apt-get install protobuf-compiler

  # Install Go plugins
  make install-proto-tools
  ```

- **GLiNER Service** - For ML-based entity extraction
  ```bash
  # Optional: Start entity extraction service
  cd services/gliner
  python -m venv venv && source venv/bin/activate
  pip install -r requirements.txt
  uvicorn main:app --host 127.0.0.1 --port 8765
  ```

---

## Project Structure

### Directory Layout

```
pedantic-raven/
├── main.go                          # Application entry point
├── go.mod, go.sum                  # Dependency management
├── Makefile                        # Build targets
├── config.toml                     # Default configuration
│
├── internal/                       # Private packages
│   ├── app/
│   │   └── events/                # Event broker pub/sub system
│   │       ├── broker.go          # Broker implementation
│   │       ├── types.go           # Event type definitions
│   │       ├── broker_test.go     # Unit tests
│   │       └── broker_benchmark_test.go  # Performance tests
│   │
│   ├── analyze/                   # Analyze mode
│   │   ├── analyze_mode.go        # Mode entry point
│   │   ├── triple_graph.go        # RDF-style semantic graph
│   │   ├── entity_analysis.go     # Entity extraction/analysis
│   │   ├── relationship_mining.go # Relationship detection
│   │   ├── hole_prioritization.go # Typed hole scoring
│   │   ├── pattern_display.go     # Visualization rendering
│   │   ├── model.go               # Mode state
│   │   ├── update.go              # Message handling
│   │   ├── view.go                # Rendering
│   │   ├── export/                # Export formats (PDF, etc.)
│   │   ├── visualizations/        # Visualization utilities
│   │   └── *_test.go              # Tests for all above
│   │
│   ├── config/
│   │   └── config.go              # Configuration loading
│   │
│   ├── context/
│   │   ├── types.go               # Analysis result types
│   │   ├── render.go              # Rendering logic
│   │   └── context_test.go        # Tests
│   │
│   ├── editor/                    # Edit mode (main editing interface)
│   │   ├── edit_mode.go           # Mode entry point
│   │   ├── model.go               # Mode state
│   │   ├── update.go              # Message handling
│   │   ├── view.go                # Rendering
│   │   ├── components.go          # UI components
│   │   ├── components_test.go     # Component tests
│   │   │
│   │   ├── buffer/                # Text buffer implementation
│   │   │   ├── buffer.go          # Buffer core
│   │   │   ├── operations.go      # Edit operations (insert, delete)
│   │   │   ├── undo_redo.go       # Undo/redo history
│   │   │   └── *_test.go          # Tests
│   │   │
│   │   ├── search/                # Search and replace
│   │   │   ├── search.go          # Search engine
│   │   │   ├── regex.go           # Regex support
│   │   │   └── *_test.go          # Tests
│   │   │
│   │   ├── semantic/              # Semantic analysis
│   │   │   ├── analyzer.go        # Main analyzer
│   │   │   ├── entity_extractor.go# Entity extraction
│   │   │   ├── relationship_detector.go # Relationship detection
│   │   │   └── *_test.go          # Tests
│   │   │
│   │   ├── syntax/                # Syntax highlighting
│   │   │   ├── highlighter.go     # Main highlighter
│   │   │   ├── lexer.go           # Language lexers
│   │   │   └── *_test.go          # Tests
│   │   │
│   │   ├── testdata/              # Test fixtures
│   │   ├── edit_mode_test.go      # Integration tests
│   │   └── gliner_e2e_test.go     # End-to-end GLiNER tests
│   │
│   ├── gliner/                    # Entity extraction service client
│   │   ├── client.go              # gRPC client
│   │   ├── types.go               # Data types
│   │   ├── errors.go              # Error handling
│   │   └── client_test.go         # Tests
│   │
│   ├── integration/               # End-to-end integration tests
│   │   ├── helpers.go             # Test utilities
│   │   ├── testdata/              # Test data fixtures
│   │   ├── testdata_generator.go  # Test data generation
│   │   ├── workflow_test.go       # Cross-mode workflows
│   │   ├── persistence_test.go    # Session persistence
│   │   ├── error_recovery_test.go # Error recovery
│   │   ├── concurrent_test.go     # Concurrent operations
│   │   └── large_dataset_test.go  # Large dataset handling
│   │
│   ├── layout/                    # Multi-pane layout engine
│   │   ├── engine.go              # Layout computation
│   │   ├── types.go               # Layout types
│   │   └── layout_test.go         # Tests
│   │
│   ├── memorydetail/              # Memory detail view component
│   │   ├── model.go               # Component state
│   │   ├── crud.go                # CRUD operations
│   │   ├── links.go               # Link management
│   │   ├── view.go                # Rendering
│   │   ├── types.go               # Types
│   │   └── *_test.go              # Tests
│   │
│   ├── memorygraph/               # Graph visualization component
│   │   ├── model.go               # Component state
│   │   ├── layout.go              # Force-directed layout
│   │   ├── types.go               # Graph types
│   │   ├── view.go                # ASCII/Unicode rendering
│   │   └── *_test.go              # Tests
│   │
│   ├── memorylist/                # Memory list view component
│   │   ├── model.go               # Component state
│   │   ├── search.go              # Search and filtering
│   │   ├── commands.go            # User commands
│   │   ├── realdata.go            # Real data loading
│   │   ├── types.go               # Types
│   │   ├── view.go                # Rendering
│   │   └── *_test.go              # Tests
│   │
│   ├── mnemosyne/                 # mnemosyne semantic memory client
│   │   ├── client.go              # Main client interface
│   │   ├── memory.go              # Memory operations (CRUD)
│   │   ├── messages.go            # Message marshaling
│   │   ├── connection.go          # Connection handling
│   │   ├── connection_manager.go  # Connection lifecycle
│   │   ├── offline.go             # Offline caching
│   │   ├── offline_cache.go       # Cache implementation
│   │   ├── retry.go               # Retry logic
│   │   ├── errors.go              # Error types
│   │   ├── pb/                    # Generated protobuf code
│   │   │   └── mnemosyne/v1/      # Protobuf messages
│   │   └── *_test.go              # Tests
│   │
│   ├── modes/                     # Mode registry and switching
│   │   ├── registry.go            # Mode registry
│   │   ├── explore.go             # Explore mode implementation
│   │   ├── registry_test.go       # Registry tests
│   │   └── explore_test.go        # Explore mode tests
│   │
│   ├── orchestrate/               # Orchestrate mode (multi-agent)
│   │   ├── orchestrate_mode.go    # Mode entry point
│   │   ├── task_graph.go          # Workflow task graph
│   │   ├── agent_log.go           # Agent execution logging
│   │   ├── session.go             # Session management
│   │   ├── launcher.go            # Task launcher
│   │   ├── plan_editor.go         # Workflow plan editor
│   │   ├── dashboard.go           # Dashboard display
│   │   ├── mode_adapter.go        # Mode integration
│   │   ├── model.go               # Mode state
│   │   ├── types.go               # Mode types
│   │   ├── update.go              # Message handling
│   │   ├── view.go                # Rendering
│   │   └── *_test.go              # Tests
│   │
│   ├── overlay/                   # Overlay UI (search, file picker)
│   │   ├── manager.go             # Overlay lifecycle
│   │   ├── search.go              # Search overlay
│   │   ├── filepicker.go          # File picker overlay
│   │   ├── types.go               # Overlay types
│   │   └── *_test.go              # Tests
│   │
│   ├── palette/                   # Command palette
│   │   ├── palette.go             # Main palette
│   │   ├── types.go               # Palette types
│   │   └── palette_test.go        # Tests
│   │
│   └── terminal/                  # Integrated terminal
│       ├── executor.go            # Command execution
│       ├── types.go               # Terminal types
│       └── terminal_test.go       # Tests
│
├── proto/                         # Protocol buffer definitions
│   └── mnemosyne/
│       └── v1/
│           ├── memory.proto       # Memory service definitions
│           └── ...
│
├── docs/                          # Documentation
│   ├── architecture.md            # Architecture guide (this explains structure)
│   ├── DEVELOPMENT.md             # Developer guide (you are here)
│   ├── CONTRIBUTING.md            # Contribution guidelines
│   ├── TESTING.md                 # Testing guide
│   ├── edit-mode-guide.md         # Edit mode user guide
│   ├── analyze-mode-guide.md      # Analyze mode user guide
│   ├── orchestrate-mode-guide.md  # Orchestrate mode user guide
│   ├── CHANGELOG.md               # Version history
│   └── images/                    # Documentation images
│
├── examples/                      # Example files
│   ├── sample-context.md          # Sample context document
│   └── ...
│
├── services/                      # External services
│   └── gliner/                    # Entity extraction service
│       ├── main.py                # FastAPI service
│       ├── requirements.txt       # Python dependencies
│       └── ...
│
└── Makefile                       # Build automation
```

### Key Files

- **main.go**: Application entry point
  - Creates event broker, mode registry, and main model
  - Initializes all modes and components
  - Starts Bubble Tea event loop

- **go.mod**: Dependency manifest
  - Specifies Go version (1.25)
  - Lists all external dependencies
  - Versions pinned for reproducibility

- **config.toml**: Application configuration
  - GLiNER service settings
  - Entity type definitions
  - Feature flags and timeouts

- **Makefile**: Build targets
  - `make build`: Compile binary
  - `make test`: Run all tests
  - `make proto`: Generate protobuf code
  - `make run`: Start application
  - `make help`: Show all targets

---

## Build Instructions

### Building the Binary

```bash
# Simple build
go build -o pedantic_raven .

# Build with version info
go build -ldflags="-s -w" -o pedantic_raven .

# Cross-platform build
GOOS=linux GOARCH=amd64 go build -o pedantic_raven_linux .
GOOS=darwin GOARCH=amd64 go build -o pedantic_raven_mac .

# Using Makefile
make build
```

### Running from Source

```bash
# Direct execution
go run main.go

# Using Makefile
make run

# With arguments (if application supports them)
go run main.go --config custom.toml
```

### Configuration

Before running, ensure `config.toml` exists in the working directory:

```toml
[gliner]
enabled = true
service_url = "http://localhost:8765"
timeout = 5
max_retries = 2
fallback_to_pattern = true
score_threshold = 0.3

[gliner.entity_types]
default = ["person", "organization", "location", "technology", "concept", "product"]
custom = []
```

### Proto Code Generation

If modifying `.proto` files:

```bash
# Generate protobuf code
make proto

# Clean generated files
make proto-clean

# One-time setup: install protoc tools
make install-proto-tools
```

---

## Running Tests

### Test Organization

Tests are organized by package:
- **Unit Tests**: `*_test.go` files testing individual functions
- **Integration Tests**: `internal/integration/` testing cross-component workflows
- **Benchmarks**: `*_benchmark_test.go` measuring performance

### Running All Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with race detector
go test ./... -race

# Run with coverage
go test ./... -cover

# Combine options
go test ./... -v -race -cover
```

### Running Specific Tests

```bash
# Run tests in a single package
go test ./internal/editor -v

# Run tests matching pattern
go test ./... -v -run TestEdit

# Run specific test
go test ./internal/editor -v -run TestEditMode_Update

# Run only short tests (skip long-running)
go test ./... -short

# Run only long tests
go test ./... -run '.*Long'
```

### Integration Tests

```bash
# Run all integration tests
go test ./internal/integration -v

# Run specific integration test
go test ./internal/integration -v -run TestEditAnalyzeWorkflow

# Run workflow tests only
go test ./internal/integration -v -run TestEdit

# Run with race detection
go test ./internal/integration -v -race
```

### Coverage Reports

```bash
# Generate coverage profile
go test ./... -coverprofile=coverage.out

# View coverage percentage by package
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Coverage by file
go test ./... -coverprofile=coverage.out ./internal/editor
go tool cover -html=coverage.out
```

### Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmark
go test -bench=BenchmarkPublish ./internal/app/events

# Run with different iteration count
go test -bench=. -benchmem -benchtime=10s ./internal/app/events

# Compare benchmark results
go test -bench=. -benchmem ./internal/app/events > new.txt
benchstat old.txt new.txt
```

### Test Coverage Targets

- **Critical Path**: 90%+ coverage
- **Business Logic**: 80%+ coverage
- **UI Components**: 60%+ coverage
- **Overall Target**: 70%+ coverage

### Running Tests with Make

```bash
# Full test suite
make test

# Short tests only
make test-short

# Custom test target
go test ./internal/editor -v -run TestBuffer
```

---

## Debugging

### Using Print Debugging

```go
// Simple print to stderr
fmt.Fprintf(os.Stderr, "DEBUG: value=%v\n", value)

// Use log package (includes timestamp)
import "log"
log.Printf("DEBUG: %v", value)
```

### Using pprof Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./internal/analyze
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./internal/mnemosyne
go tool pprof mem.prof

# Interactive pprof
go tool pprof
> top10        # Show top 10 functions by CPU
> list func    # Show source for specific function
> web          # Generate graph (requires graphviz)
```

### Using Delve Debugger

```bash
# Install debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug tests
dlv test ./internal/editor -- -test.run TestEditMode

# Debug running binary
dlv exec ./pedantic_raven
(dlv) continue     # Run until breakpoint
(dlv) next         # Step to next line
(dlv) print var    # Print variable value
(dlv) breakpoint   # Show breakpoints
```

### Common TUI Debugging Issues

**Problem**: Application crashes or hangs in terminal

**Solution**: Use alternate screen buffer
```bash
# Log to file instead of stdout
go run main.go 2>debug.log

# Or redirect stderr to file
./pedantic_raven 2>debug.log &
tail -f debug.log
```

**Problem**: Testing Bubble Tea components

**Solution**: Test using TestModel
```go
func TestComponent(t *testing.T) {
    m := Model{...}

    // Send message
    updated, cmd := m.Update(tea.KeyMsg{...})

    // Check state
    assert.Equal(t, expected, updated)

    // Execute command
    msg := cmd()
    assert.NotNil(t, msg)
}
```

### Debugging Event Broker

```go
// Add debug logging to broker
func (b *Broker) Publish(event Event) {
    fmt.Fprintf(os.Stderr, "EVENT: type=%s subscribers=%d\n",
        event.Type, len(b.subscribers))
    // ... rest of publish logic
}
```

---

## Common Tasks

### Adding a New Mode

1. **Create mode package**:
```bash
mkdir internal/newmode
touch internal/newmode/{model,update,view}.go
```

2. **Implement Mode interface**:
```go
// internal/newmode/newmode.go
package newmode

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
    // Your state fields
}

func (m Model) ID() modes.ModeID { return "newmode" }
func (m Model) Name() string { return "New Mode" }
func (m Model) Description() string { return "..." }
func (m Model) Init() tea.Cmd { return nil }
func (m Model) OnEnter() tea.Cmd { return nil }
func (m Model) OnExit() tea.Cmd { return nil }
func (m Model) Update(msg tea.Msg) (modes.Mode, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKey(msg)
    }
    return m, nil
}
func (m Model) View() string {
    return "Your mode content"
}
```

3. **Register in main.go**:
```go
// main.go
registry.Register(newmode.New())
```

4. **Add tests**:
```bash
touch internal/newmode/newmode_test.go
```

### Adding a New Component

1. **Define Component interface**:
```go
// internal/component/component.go
package component

import tea "github.com/charmbracelet/bubbletea"

type Component struct {
    width, height int
}

func (c Component) Update(msg tea.Msg) (Component, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return c.handleKey(msg)
    }
    return c, nil
}

func (c Component) View() string {
    return "Component content"
}

func (c Component) SetSize(width, height int) {
    c.width = width
    c.height = height
}
```

2. **Add to mode**:
```go
// internal/mymode/model.go
type Model struct {
    component *component.Component
}

func (m Model) Update(msg tea.Msg) (Mode, tea.Cmd) {
    comp, cmd := m.component.Update(msg)
    m.component = comp
    return m, cmd
}

func (m Model) View() string {
    return m.component.View()
}
```

3. **Test component**:
```bash
touch internal/component/component_test.go
```

### Adding Event Types

1. **Define new event in broker types**:
```go
// internal/app/events/types.go
const (
    CustomEventType EventType = "custom.event"
)
```

2. **Publish event**:
```go
broker.Publish(Event{
    Type: CustomEventType,
    Payload: json.Marshal(data),
})
```

3. **Subscribe to event**:
```go
eventCh := broker.Subscribe(CustomEventType)
defer broker.Unsubscribe(eventCh)

for event := range eventCh {
    // Handle event
}
```

### Integrating with mnemosyne

1. **Create client**:
```go
config := mnemosyne.ConfigFromEnv()
client, err := mnemosyne.NewClient(config)
if err != nil {
    // Handle offline mode
    client = mnemosyne.NewOfflineClient()
}
```

2. **Create memory**:
```go
memory := mnemosyne.Memory{
    Content: "Your content",
    Namespace: "project:myapp",
    Tags: []string{"tag1", "tag2"},
    Metadata: map[string]interface{}{
        "importance": 8,
    },
}

id, err := client.CreateMemory(ctx, memory)
```

3. **Search memories**:
```go
results, err := client.SearchMemories(ctx, "query")
for _, mem := range results {
    fmt.Println(mem.Content)
}
```

---

## Code Style Guidelines

### Go Conventions

**Follow standard Go conventions**:
- Use `gofmt` for formatting
- Use `golint` for linting
- Use `go vet` for static analysis

```bash
# Format code
go fmt ./...

# Lint code
go lint ./...

# Vet code
go vet ./...
```

### Naming

**Packages**: Lowercase, single word preferred
```go
package editor    // Good
package ed        // Bad - unclear
```

**Types**: Exported are PascalCase, unexported are camelCase
```go
type Model struct {}      // Exported
type editorModel struct {} // Unexported
```

**Functions**: Exported are PascalCase, unexported are camelCase
```go
func (m Model) Update() {}  // Exported method
func (m Model) update() {}  // Unexported helper
```

**Constants**: UPPER_CASE for constants
```go
const DefaultTimeout = 30 * time.Second
const maxRetries = 3
```

### Comments

**Package comments**: Describe package purpose
```go
// Package editor provides text editing functionality with semantic analysis.
package editor
```

**Type comments**: Describe exported types
```go
// Model represents the editor mode state.
type Model struct {
    // ... fields
}
```

**Function comments**: Start with function name
```go
// Update processes a Bubble Tea message and returns updated state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ...
}
```

### Error Handling

**Always check errors**:
```go
// Good
result, err := someFunc()
if err != nil {
    return err
}

// Bad
result, _ := someFunc()  // Ignoring error

// Good - when ignoring is intentional, explain why
if _, err := someFunc(); err != nil {
    // Error already logged elsewhere, safe to ignore
}
```

**Wrap errors with context**:
```go
// Good
if err := someFunc(); err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// Bad
return err  // No context
```

### Documentation

**Document public APIs**:
```go
// Create creates a new Model with the specified initial state.
// Returns an error if the state is invalid.
func Create(state State) (*Model, error) {
    // ...
}
```

**Use examples in tests**:
```go
func ExampleCreate() {
    m, _ := Create(DefaultState)
    fmt.Println(m.View())
    // Output: ...
}
```

### Testing

**Test naming**:
```go
func TestCreate(t *testing.T) {}           // Test Create function
func TestUpdate_KeyPress(t *testing.T) {}  // Test Update with KeyPress
func TestUpdate_Edge(t *testing.T) {}      // Test edge case
```

**Table-driven tests** for multiple cases:
```go
func TestUpdate(t *testing.T) {
    cases := []struct {
        name string
        msg  tea.Msg
        want string
    }{
        {"key_press", tea.KeyMsg{Runes: []rune("a")}, "a"},
        {"empty", nil, ""},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            m := Model{}
            _, _ = m.Update(tc.msg)
            // Assert
        })
    }
}
```

### Imports

**Organize imports**:
```go
import (
    // Standard library
    "fmt"
    "os"

    // External packages
    tea "github.com/charmbracelet/bubbletea"

    // Internal packages
    "github.com/rand/pedantic-raven/internal/editor"
)
```

---

## Performance Considerations

### 1. Avoid Allocations in Hot Paths

**Bad**: Allocates on every call
```go
func (m Model) View() string {
    result := ""
    for _, item := range m.items {
        result += item.String() + "\n"
    }
    return result
}
```

**Good**: Uses strings.Builder
```go
func (m Model) View() string {
    var buf strings.Builder
    for _, item := range m.items {
        buf.WriteString(item.String())
        buf.WriteRune('\n')
    }
    return buf.String()
}
```

### 2. Cache Expensive Computations

**Bad**: Recomputes on every View()
```go
func (m Model) View() string {
    sorted := sortItems(m.items)  // O(n log n) every render
    return renderItems(sorted)
}
```

**Good**: Cache in Update()
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ... handle messages ...
    m.sorted = sortItems(m.items)  // Cache once per update
}

func (m Model) View() string {
    return renderItems(m.sorted)
}
```

### 3. Event Broker Buffering

**Configure appropriate buffer size**:
```go
// Small buffer for high-frequency events
broker := events.NewBroker(bufferSize: 10)

// Large buffer for less frequent, important events
broker := events.NewBroker(bufferSize: 100)
```

### 4. mnemosyne Connection Pooling

The gRPC connection automatically multiplexes requests. Don't create multiple clients:

```go
// Bad: Creates new connection per request
func storeMemory() {
    client, _ := mnemosyne.NewClient(config)
    defer client.Close()
    client.CreateMemory(ctx, memory)
}

// Good: Reuse single client
var client *mnemosyne.Client

func init() {
    client, _ = mnemosyne.NewClient(config)
}

func storeMemory() {
    client.CreateMemory(ctx, memory)
}
```

### 5. Lazy Loading

Load data only when needed:
```go
// Bad: Load all memories on startup
func init() {
    allMemories := client.ListMemories(ctx)
}

// Good: Load on-demand
func loadMemories() {
    allMemories := client.ListMemories(ctx)
}
```

---

## Troubleshooting

### Common Issues

**Issue**: Build fails with "go: cannot find module"

**Solution**: Update dependencies
```bash
go mod tidy
go mod download
go build ./...
```

**Issue**: Tests fail with "connection refused" (mnemosyne)

**Solution**: Start mnemosyne or run in offline mode
```bash
# Start mnemosyne
mnemosyne serve

# Or set environment variable to disable mnemosyne
MNEMOSYNE_ENABLED=false go test ./...
```

**Issue**: GLiNER service returns errors

**Solution**: Check service is running and accessible
```bash
# Start GLiNER service
cd services/gliner
python -m venv venv && source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --host 127.0.0.1 --port 8765

# Test endpoint
curl http://localhost:8765/docs
```

**Issue**: Application crashes on startup

**Solution**: Check configuration and ensure mnemosyne is accessible
```bash
# Verify configuration
cat config.toml

# Check mnemosyne server
netstat -an | grep 50051  # Or ss on Linux

# Run with verbose logging
go run main.go 2>debug.log
```

**Issue**: Terminal rendering issues

**Solution**: Check terminal compatibility
```bash
# Ensure TERM is set correctly
echo $TERM

# Try with different terminal
TERM=xterm-256color ./pedantic_raven
```

**Issue**: High memory usage

**Solution**: Profile with pprof
```bash
# Generate memory profile
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# In pprof:
> top10
> list functionName
```

---

## Summary

Development workflow:

1. **Setup**: Install Go 1.25+, clone repo, run `go mod download`
2. **Build**: `make build` creates `pedantic_raven` binary
3. **Run**: `./pedantic_raven` or `make run`
4. **Test**: `go test ./...` runs all tests with coverage
5. **Debug**: Use `fmt.Printf`, `log.Printf`, or Delve debugger
6. **Code**: Follow Go conventions, add comments, write tests
7. **Performance**: Profile with pprof, avoid allocations, cache results

For questions about specific components, see the mode-specific guides:
- [Edit Mode Guide](edit-mode-guide.md)
- [Analyze Mode Guide](analyze-mode-guide.md)
- [Orchestrate Mode Guide](orchestrate-mode-guide.md)

---

*For architecture overview, see [Architecture Guide](architecture.md).*
