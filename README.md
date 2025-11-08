# Pedantic Raven

> **Interactive Context Engineering Environment**

Pedantic Raven combines the specialized context engineering capabilities of ICS (Integrated Context Studio) with the rich interaction patterns of Crush TUI to create a powerful, production-ready terminal interface for creating, editing, and refining context with semantic analysis and memory integration.

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/dl/)
[![Tests](https://img.shields.io/badge/tests-461%20passing-brightgreen.svg)](#testing)
[![Phase](https://img.shields.io/badge/phase-4.1%20complete-success.svg)](#features)

---

## Features

### Phase 1: Foundation ✅ (Complete)

- **Multi-Pane Layouts**: Hierarchical pane composition with 5 layout modes
- **Mode System**: 5 application modes (Edit, Explore, Analyze, Orchestrate, Collaborate)
- **Event-Driven Architecture**: Decoupled PubSub event system
- **Overlay System**: Modal and non-modal dialogs with flexible positioning
- **Command Palette**: Fuzzy search for command discovery (Ctrl+K)
- **Focus Management**: Keyboard navigation between panes
- **Responsive Design**: Adapts to terminal size

### Phase 2: Semantic Analysis & Edit Mode ✅ (Complete)

- **Semantic Analyzer**: Real-time NLP-style text analysis
  - Entity extraction with 6 types (Person, Place, Thing, Concept, Organization, Technology)
  - Relationship detection (subject-predicate-object patterns)
  - Typed holes (??Type and !!constraint markers)
  - Dependency tracking (imports, requires, references)
  - Triple generation (RDF-style semantic structures)
- **Context Panel**: Rich semantic results display
  - 5 sections (Entities, Relationships, Typed Holes, Dependencies, Triples)
  - Filtering and navigation
  - Priority/complexity scoring
- **Integrated Terminal**: Command execution with history
  - Built-in commands (:help, :clear, :history)
  - Shell command execution
  - mnemosyne CLI integration
  - Scrollable output buffer
- **Edit Mode**: Complete context editing environment
  - Auto-triggered semantic analysis (500ms debounce)
  - Multi-component layout (editor, context panel, terminal)
  - Focus management across components

### Phase 3: Advanced Editor Features ✅ (Complete)

**Phase 3.1** ✅ Complete:
- Buffer Manager integration with undo/redo support
- Undo/redo keybindings (Ctrl+Z, Ctrl+Y)
- Cursor position management

**Phase 3.2** ✅ Complete:
- File I/O operations (OpenFile, SaveFile, SaveFileAs)
- Atomic file saves (temp + rename)
- Error handling and path management
- Dirty flag tracking

**Phase 3.3** ✅ Complete:
- Search engine with literal and regex support
- Case sensitive/insensitive search
- Whole word matching
- Replace current match and replace all
- Multi-line search support
- Undo integration for replacements

**Phase 3.4** ✅ Complete:
- Token-based syntax highlighting system
- Go language full tokenization
- Markdown formatting support
- Automatic language detection (by extension and content)
- Extensible tokenizer architecture
- 12 token types with default color scheme

### Phase 4.1: mnemosyne RPC Client ✅ (Complete)

- **gRPC Client Library**: Full-featured client for mnemosyne memory system
  - Connection management with configurable timeouts
  - Health checks and server statistics
- **CRUD Operations**: Complete memory lifecycle management
  - StoreMemory, GetMemory, UpdateMemory, DeleteMemory, ListMemories
  - Namespace support (Global, Project, Session)
  - Importance scoring and tagging
  - Optional LLM enrichment
- **Search Operations**: Advanced memory retrieval
  - Recall: Hybrid search (semantic + FTS + graph)
  - SemanticSearch: Pure embedding-based search
  - GraphTraverse: Multi-hop graph traversal
  - GetContext: Retrieve memories with linked context
- **Streaming Support**: Progressive result delivery
  - RecallStream, ListMemoriesStream, StoreMemoryStream
  - Progress updates for long-running operations
- **Error Handling**: Comprehensive gRPC status mapping
  - Domain-specific errors (NotFound, InvalidArgument, etc.)
  - Helper functions for error checking

**Next Phase**:
- **Phase 4.2**: Memory List Component - TUI display and interaction
- **Phase 4.3**: Memory Detail View - Rich memory visualization
- **Phase 4.4**: Graph Visualization - Interactive memory graph
- **Phase 4.5**: Explore Mode Integration - Complete memory workspace

---

## Quick Start

### Prerequisites

- Go 1.25 or higher
- Terminal with 256+ colors
- Minimum 120x30 terminal size (80x24 supported with compact layout)

### Installation

```bash
# Clone repository
git clone https://github.com/rand/pedantic_raven.git
cd pedantic_raven

# Build
go build -o pedantic_raven .

# Run demo
./pedantic_raven
```

### Usage

**Quick Reference:**

| Key | Action |
|-----|--------|
| `1`, `2`, `3` | Switch modes (Edit, Explore, Analyze) |
| `Tab` | Cycle focus between panes |
| `Ctrl+K` | Open command palette |
| `?` | Show about dialog |
| `Ctrl+C` | Quit |

**Edit Mode (when editor focused):**
- Type to enter text (triggers automatic semantic analysis)
- `Backspace` - Delete character
- `Enter` - New line
- `Ctrl+Z` - Undo
- `Ctrl+Y` or `Ctrl+Shift+Z` - Redo

**Edit Mode (when context panel focused):**
- `j`/`k` or `↓`/`↑` - Scroll
- `PgUp`/`PgDn` - Scroll by page
- `Enter` - Toggle section

**Edit Mode (when terminal focused):**
- Type commands and press `Enter`
- `↑`/`↓` - Command history
- Built-in commands: `:help`, `:clear`, `:history`, `:exit`

📚 **See [docs/USAGE.md](docs/USAGE.md) for complete keyboard reference and feature guide**

---

## Architecture

Pedantic Raven is built on the **Elm Architecture** using [Bubble Tea](https://github.com/charmbracelet/bubbletea):

- **Immutable state updates**
- **Pure functions**
- **Command-based side effects**
- **Message passing**

### Component Overview

```
┌─────────────────────────────────────────┐
│          Application Model              │
├─────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────────┐ │
│  │ Event Broker│  │  Overlay Manager │ │
│  └─────────────┘  └──────────────────┘ │
│  ┌─────────────┐  ┌──────────────────┐ │
│  │Mode Registry│  │  Layout Engine   │ │
│  └─────────────┘  └──────────────────┘ │
│  ┌─────────────┐                       │
│  │  Palette    │                       │
│  └─────────────┘                       │
└─────────────────────────────────────────┘
```

### Foundation Components

1. **PubSub Event System** ([`internal/app/events/`](internal/app/events/))
   - Thread-safe event broker
   - 40+ domain event types
   - Non-blocking publish

2. **Layout Engine** ([`internal/layout/`](internal/layout/))
   - Hierarchical pane composition
   - 5 layout modes
   - Focus management
   - Responsive design

3. **Mode Registry** ([`internal/modes/`](internal/modes/))
   - 5 application modes
   - Lifecycle hooks (Init, OnEnter, OnExit)
   - Mode switching with history

4. **Overlay System** ([`internal/overlay/`](internal/overlay/))
   - Modal and non-modal overlays
   - Stack-based management
   - Built-in dialogs (Confirm, Message)
   - Flexible positioning strategies

5. **Command Palette** ([`internal/palette/`](internal/palette/))
   - Fuzzy search with scoring
   - Category-based organization
   - Command execution framework

---

## Project Structure

```
pedantic_raven/
├── main.go                      # Demo application
├── spec.md                      # Comprehensive specification
├── README.md                    # This file
├── go.mod                       # Go module definition
├── docs/
│   └── PHASE1_COMPLETE.md       # Phase 1 completion summary
└── internal/
    ├── app/
    │   └── events/              # PubSub event system
    │       ├── types.go         # Event types (40+ events)
    │       ├── broker.go        # Event broker
    │       └── broker_test.go   # 7 tests (53% coverage)
    ├── layout/                  # Multi-pane layout engine
    │   ├── types.go             # Pane hierarchy
    │   ├── engine.go            # Layout engine
    │   └── layout_test.go       # 19 tests (54.5% coverage)
    ├── modes/                   # Mode registry
    │   ├── registry.go          # Mode management
    │   └── registry_test.go     # 17 tests (92% coverage)
    ├── overlay/                 # Overlay system
    │   ├── types.go             # Overlay interface
    │   ├── manager.go           # Stack manager
    │   └── overlay_test.go      # 24 tests (66.7% coverage)
    └── palette/                 # Command palette
        ├── types.go             # Command registry
        ├── palette.go           # Palette overlay
        └── palette_test.go      # 20 tests (88.3% coverage)
```

---

## Testing

### Run Tests

```bash
# All tests
go test ./...

# With coverage
go test ./... -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out
```

### Test Summary

| Package | Tests | Coverage | Status |
|---------|-------|----------|--------|
| app/events  | 18 | ~70% | ✅ |
| context | 25 | ~80% | ✅ |
| editor | 78 | ~85% | ✅ |
| editor/buffer | 52 | ~85% | ✅ |
| editor/search | 35 | ~90% | ✅ |
| editor/semantic | 63 | ~90% | ✅ |
| editor/syntax | 31 | ~85% | ✅ |
| layout  | 34 | ~65% | ✅ |
| mnemosyne | 66 | ~95% | ✅ |
| modes   | 13 | ~92% | ✅ |
| overlay | 25 | ~70% | ✅ |
| palette | 19 | ~88% | ✅ |
| terminal | 38 | ~80% | ✅ |
| **Total** | **461** | **~84%** | **✅** |

All tests passing ✅

📊 **See [docs/PHASE4.1_SUMMARY.md](docs/PHASE4.1_SUMMARY.md) for Phase 4.1 details**
📊 **See [docs/PHASE3_SUMMARY.md](docs/PHASE3_SUMMARY.md) for Phase 3 details**
📊 **See [docs/PHASE2_SUMMARY.md](docs/PHASE2_SUMMARY.md) for Phase 2 details**

---

## Development

### Prerequisites

- Go 1.25+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) v1.2.6+
- [Lipgloss](https://github.com/charmbracelet/lipgloss) v1.0.0+

### Dependencies

```bash
go get github.com/charmbracelet/bubbletea@v1.2.6
go get github.com/charmbracelet/lipgloss@v1.0.0
go get github.com/charmbracelet/bubbles@v0.21.0
```

### Building

```bash
# Development build
go build -o pedantic_raven .

# Production build with optimizations
go build -ldflags="-s -w" -o pedantic_raven .

# Run without building
go run main.go
```

### Code Organization

- **`internal/`**: All internal packages (not importable by external code)
- **Component pattern**: Each package contains `types.go`, implementation, and tests
- **Test co-location**: Tests live alongside the code they test
- **Interface-first**: Define interfaces before implementations

### Design Patterns

1. **Elm Architecture** (Model-Update-View)
2. **Composite Pattern** (Hierarchical panes)
3. **Observer Pattern** (PubSub events)
4. **Registry Pattern** (Modes, commands, components)
5. **Strategy Pattern** (Position strategies, layout modes)

---

## Contributing

### Development Workflow

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Write** tests for your changes
4. **Implement** your feature
5. **Run** tests (`go test ./...`)
6. **Commit** your changes (`git commit -m 'Add amazing feature'`)
7. **Push** to your branch (`git push origin feature/amazing-feature`)
8. **Open** a Pull Request

### Code Standards

- Follow Go conventions (gofmt, golint)
- Write tests for all new code (aim for 70%+ coverage)
- Document all exported types and functions
- Use meaningful variable and function names
- Keep functions focused and short (<50 lines when possible)

### Commit Messages

Follow the pattern:
```
<type>: <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `style`, `chore`

Example:
```
feat: Add command palette fuzzy search

Implement fuzzy matching algorithm with scoring:
- Exact match: +100
- Name contains: +50
- Description contains: +20
- Subsequence match: +30

Tests: 8 new tests for fuzzy matching
```

---

## Documentation

### Comprehensive Docs

- **[USAGE.md](docs/USAGE.md)**: Complete user guide with keyboard shortcuts, features, and examples
- **[PHASE2_SUMMARY.md](docs/PHASE2_SUMMARY.md)**: Phase 2 completion summary with architecture details
- **[PHASE1_COMPLETE.md](docs/PHASE1_COMPLETE.md)**: Phase 1 completion summary
- **[spec.md](spec.md)**: Full specification with requirements and architecture
- **Package docs**: Each package has detailed comments and examples

### External Resources

- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Lipgloss Docs](https://github.com/charmbracelet/lipgloss)
- [Crush TUI Reference](https://github.com/charmbracelet/crush)

---

## Roadmap

### Phase 1: Foundation ✅ (Complete)
- [x] PubSub event system
- [x] Multi-pane layout engine
- [x] Mode registry
- [x] Overlay system
- [x] Command palette

### Phase 2: Semantic Analysis & Edit Mode ✅ (Complete)
- [x] Semantic analyzer with streaming
- [x] Entity extraction (6 types)
- [x] Relationship detection
- [x] Typed holes (??Type, !!constraint)
- [x] Context panel with 5 sections
- [x] Terminal integration
- [x] Edit Mode with auto-analysis
- [x] 291 tests passing

### Phase 3: Advanced Editor ✅ (Complete - 8 days)
- [x] **Phase 3.1**: Buffer Manager integration, undo/redo (Days 1-2) ✅
- [x] **Phase 3.2**: File operations - open, save, save-as (Days 3-4) ✅
- [x] **Phase 3.3**: Search and replace - literal and regex (Days 5-6) ✅
- [x] **Phase 3.4**: Syntax highlighting - Go and Markdown (Days 7-8) ✅

### Phase 4: Explore Mode (3-4 weeks) - In Progress
- [x] **Phase 4.1**: mnemosyne RPC Client (Days 1-3) ✅
  - gRPC client library with full CRUD operations
  - Advanced search (Recall, SemanticSearch, GraphTraverse, GetContext)
  - Streaming support for progressive results
  - 66 tests, comprehensive error handling
- [ ] **Phase 4.2**: Memory List Component (Days 4-6)
  - TUI component for memory display
  - Filtering and navigation
  - Integration with mnemosyne client
- [ ] **Phase 4.3**: Memory Detail View (Days 7-9)
  - Rich memory visualization
  - Linked memories display
  - Editing capabilities
- [ ] **Phase 4.4**: Graph Visualization (Days 10-14)
  - Interactive memory graph
  - Graph traversal UI
  - Visual link exploration
- [ ] **Phase 4.5**: Explore Mode Integration (Days 15-16)
  - Complete workspace integration
  - Namespace navigation
  - Triple exploration

### Phase 5: Analyze Mode (2-3 weeks)
- [ ] Statistical analysis
- [ ] Entity relationship graphs
- [ ] Typed hole prioritization UI
- [ ] Dependency tree visualization

### Phase 6: Orchestrate Mode (4-5 weeks)
- [ ] Agent coordination
- [ ] Task management
- [ ] Progress monitoring
- [ ] Multi-agent workflows

### Phase 7: Collaborate Mode (3-4 weeks)
- [ ] Live multi-user editing
- [ ] Presence awareness
- [ ] Conflict resolution
- [ ] Shared annotations

### Phase 8: Polish & Production (2-3 weeks)
- [ ] Performance optimization
- [ ] Comprehensive documentation
- [ ] Packaging and distribution
- [ ] Release 1.0

**Estimated Timeline**: 4-6 months remaining

---

## Related Projects

- **[mnemosyne](https://github.com/rand/mnemosyne)**: Semantic memory system with RPC server
- **[ICS](https://github.com/rand/mnemosyne/src/ics)**: Original Integrated Context Studio
- **[Crush](https://github.com/charmbracelet/crush)**: Rich TUI inspiration

---

## Performance

### Benchmarks

```
Terminal Size: 120x30
Render Time: <16ms (60 FPS)
Memory Usage: ~10MB
Startup Time: <100ms
```

### Optimization

- Non-blocking event publishing
- Efficient pane rendering (only changed regions)
- Lazy component initialization
- Smart terminal updates (Bubble Tea handles this)

---

## Troubleshooting

### Terminal Too Small

Pedantic Raven automatically switches to compact layout for terminals smaller than 120x30. For best experience, use at least 120x30.

### Rendering Issues

Ensure your terminal supports:
- 256 colors (most modern terminals)
- UTF-8 encoding
- ANSI escape sequences

Test with:
```bash
echo $TERM
# Should show: xterm-256color or similar
```

### Performance Issues

If experiencing lag:
- Close other terminal applications
- Increase terminal buffer size
- Update to latest version of Pedantic Raven

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

This project is part of the mnemosyne ecosystem.

---

## Acknowledgments

- **Bubble Tea**: Excellent TUI framework by Charm
- **ICS**: Original context engineering system
- **Crush**: Rich TUI inspiration and patterns
- **Go Community**: For amazing tooling and libraries

---

## Contact

- **Repository**: https://github.com/rand/pedantic_raven
- **Related**: [mnemosyne](https://github.com/rand/mnemosyne) memory system

---

**Status**: Phase 4.1 Complete ✅ | Next: Phase 4.2 (Memory List Component)

**Stats**: 461 tests passing | ~11,000 lines of code | 84% coverage

Built with ❤️ using [Bubble Tea](https://github.com/charmbracelet/bubbletea)
