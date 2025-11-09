# Pedantic Raven

**Terminal-based context engineering with semantic analysis and memory integration.**

Extract entities, relationships, and typed holes from your text. Edit AI context documents with real-time semantic feedback. Connect to [mnemosyne](https://github.com/rand/mnemosyne) for persistent memory and multi-agent orchestration.

[![Go 1.25+](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/dl/)
[![Tests](https://img.shields.io/badge/tests-754%20passing-brightgreen.svg)](#testing)
[![Coverage](https://img.shields.io/badge/coverage-70%25-yellowgreen.svg)](docs/TESTING.md)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Docs](https://img.shields.io/badge/docs-latest-blue.svg)](docs/)
[![Contributing](https://img.shields.io/badge/contributing-welcome-brightgreen.svg)](docs/CONTRIBUTING.md)

---

## What It Does

Pedantic Raven helps you create structured context documents for AI systems like Claude. It analyzes your text as you type, extracting:

- **Entities** - People, places, concepts, technologies (6 types)
- **Relationships** - Subject-predicate-object patterns with confidence scoring
- **Typed Holes** - `??Type` markers for incomplete sections with priority/complexity scores
- **Dependencies** - Imports, requires, and references across your context
- **Triples** - RDF-style semantic structures for knowledge graphs

All this happens in a **rich terminal UI** with syntax highlighting, integrated terminal, and direct connection to the mnemosyne memory system.

---

## Why Use It?

**Problem**: Creating context for AI systems is tedious. You write documentation, but lose track of entities, relationships, and gaps. Context becomes stale and disorganized.

**Solution**: Pedantic Raven gives you live semantic feedback as you write. See what entities you've mentioned, what relationships exist, what sections need filling in. Store everything in mnemosyne's semantic memory for easy recall across projects.

**Use Cases**:
- Writing structured prompts for Claude with semantic validation
- Building knowledge graphs from documentation
- Tracking architectural decisions with typed hole markers
- Creating AI context that stays organized and searchable

---

## Current Status

**Phase**: 4.4 (Graph Visualization - Complete)
**Project State**: Active Development
**Production Ready**: No (expect breaking changes)

### What Works Now

- ✅ **Edit Mode** - Full-featured text editor
  - Syntax highlighting (Go, Markdown)
  - Search and replace (literal and regex)
  - Undo/redo support
  - File operations (open, save, atomic writes)
- ✅ **Semantic Analysis** - Real-time extraction
  - 6 entity types with occurrence counting
  - Relationship detection with confidence scores
  - Typed hole markers (??Type, !!constraint)
  - Dependency tracking
  - RDF-style triple generation
- ✅ **Context Panel** - Live semantic results
  - 5 sections (entities, relationships, typed holes, dependencies, triples)
  - Filtering and navigation
  - Priority/complexity indicators
- ✅ **Integrated Terminal** - Command execution
  - Built-in commands (`:help`, `:clear`, `:history`)
  - Shell command pass-through
  - Command history (100 entries)
- ✅ **mnemosyne Integration** - Memory system client
  - Full CRUD operations (create, read, update, delete, list)
  - Advanced search (hybrid, semantic, graph traversal)
  - Streaming support for progressive results
  - Namespace management (global, project, session)

### What Works Now (Continued)

- ✅ **Memory List UI** - Browse stored memories
  - Search and filtering
  - Sorting by importance and recency
  - Rich metadata display
- ✅ **Memory Detail View** - Rich memory visualization
  - Full metadata display
  - Link navigation
  - Content scrolling
- ✅ **Graph Visualization** - Visual memory relationships
  - Force-directed layout algorithm
  - ASCII/Unicode rendering
  - Pan, zoom, and node selection
  - Expand/collapse hierarchical nodes
- ✅ **Explore Mode** - Complete memory workspace
  - Graph visualization mode
  - Sample graph for demonstration
  - Full keyboard navigation

### Planned

- 📋 **Analyze Mode** - Statistical analysis and insights
- 📋 **Orchestrate Mode** - Multi-agent coordination
- 📋 **Collaborate Mode** - Live multi-user editing

See [ROADMAP.md](ROADMAP.md) for detailed timeline and [docs/CHANGELOG.md](docs/CHANGELOG.md) for development history.

---

## Quick Start

### Prerequisites

- Go 1.25 or higher
- Terminal with 256+ colors (most modern terminals)
- Minimum 120x30 terminal size recommended (80x24 supported with compact layout)

### Build and Run

```bash
# Clone repository
git clone https://github.com/rand/pedantic_raven.git
cd pedantic_raven

# Build
go build -o pedantic_raven .

# Run
./pedantic_raven
```

### Basic Usage

**Mode Switching:**
- `1` - Edit mode (default)
- `2` - Explore mode (memory workspace)
- `3` - Analyze mode (semantic insights)

**Edit Mode:**
- Type to enter text (semantic analysis runs automatically after 500ms pause)
- `Tab` - Cycle focus between editor, context panel, and terminal
- `Ctrl+K` - Open command palette
- `Ctrl+Z` / `Ctrl+Y` - Undo / Redo
- `Ctrl+F` - Search
- `Ctrl+S` - Save file

**Context Panel** (when focused):
- `j` / `k` or `↓` / `↑` - Scroll results
- `Enter` - Expand/collapse sections
- `PgUp` / `PgDn` - Page navigation

**Terminal** (when focused):
- Type commands and press `Enter`
- `↑` / `↓` - Command history
- `:help` - Show built-in commands

**Explore Mode** (memory workspace):
- `g` - Toggle between list/detail and graph views
- `Tab` - Switch focus (list ↔ detail) in standard layout
- `?` - Show context-aware help overlay
- **Standard Layout:**
  - `j/k` - Navigate list or scroll detail
  - `Enter` - Select memory or navigate link
  - `/` - Search memories (when list focused)
  - `r` - Refresh data
- **Graph Layout:**
  - `h/j/k/l` - Pan graph viewport
  - `+/-` - Zoom in/out, `0` - Reset view
  - `Tab` - Select nodes, `e/x` - Expand/collapse
  - `c` - Center on selected, `r` - Re-layout

📖 **Full guide**: [docs/USAGE.md](docs/USAGE.md)

---

## Documentation

### For Users

- **[USAGE.md](docs/USAGE.md)** - Complete user guide with all features and keybindings
- **[Edit Mode Guide](docs/edit-mode-guide.md)** - Detailed guide to context editing and semantic analysis
- **[Explore Mode Guide](docs/explore-mode-guide.md)** - Memory workspace with browsing, searching, and graph visualization
- **[Analyze Mode Guide](docs/analyze-mode-guide.md)** - Statistical analysis and insights
- **[Orchestrate Mode Guide](docs/orchestrate-mode-guide.md)** - Multi-agent coordination and task management
- **[TESTING.md](docs/TESTING.md)** - Testing framework and coverage

### For Developers

- **[Architecture Guide](docs/architecture.md)** - System design, patterns, and component architecture
- **[Developer Guide](docs/DEVELOPMENT.md)** - Setup, build, test, and debugging instructions
- **[Contributing Guide](docs/CONTRIBUTING.md)** - How to contribute, code style, and PR process
- **[CHANGELOG.md](docs/CHANGELOG.md)** - Version history and breaking changes

### Quick Links

- [Roadmap](ROADMAP.md) - Feature timeline and future plans
- [GitHub Issues](https://github.com/rand/pedantic-raven/issues) - Report bugs or request features
- [GitHub Discussions](https://github.com/rand/pedantic-raven/discussions) - Questions and ideas

---

## mnemosyne Ecosystem

Pedantic Raven is part of the **mnemosyne ecosystem** - a suite of tools for semantic memory and context engineering.

### The Ecosystem

```
┌─────────────────────────────────────────────────────────┐
│                 mnemosyne Ecosystem                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐         ┌──────────────────────┐     │
│  │  mnemosyne   │◄────────┤  Pedantic Raven      │     │
│  │  (server)    │  gRPC   │  (TUI client)        │     │
│  │              │         │  - Context editor    │     │
│  │  - Memory DB │         │  - Semantic analysis │     │
│  │  - Search    │         │  - Memory workspace  │     │
│  │  - Graph     │         └──────────────────────┘     │
│  │  - Agents    │                                      │
│  └──────────────┘         ┌──────────────────────┐     │
│         │                 │  ICS (Legacy)        │     │
│         │                 │  Being replaced by   │     │
│         └─────────────────┤  Pedantic Raven      │     │
│                           └──────────────────────┘     │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Components

**[mnemosyne](https://github.com/rand/mnemosyne)** - Semantic memory system
- gRPC server for memory storage and retrieval
- Graph-based memory relationships with link types
- Semantic search with embedding vectors (768d/1536d)
- Multi-agent orchestration capabilities
- Rust-based for performance and safety

**Pedantic Raven** - Context engineering interface (this project)
- Terminal UI for creating and editing AI context
- Real-time semantic analysis during editing
- Client for mnemosyne memory operations
- Visual memory workspace and graph exploration
- Go-based with Bubble Tea framework

**ICS (Integrated Context Studio)** - Original context tool
- Legacy Python-based context editor
- Being replaced by Pedantic Raven
- Simpler feature set, less interactive

### Data Flow

```
User types in Pedantic Raven
      ↓
Semantic analysis extracts entities/relationships
      ↓
Results shown in context panel
      ↓
User stores to mnemosyne (via gRPC)
      ↓
mnemosyne indexes and enriches memory
      ↓
User recalls memories (hybrid search)
      ↓
Pedantic Raven displays results
```

### Why This Architecture?

- **Separation of concerns**: UI (Pedantic Raven) vs storage/search (mnemosyne)
- **Multiple clients**: CLI, TUI, future web interface all use same mnemosyne backend
- **Language strengths**: Go for TUI responsiveness, Rust for database performance
- **Scalability**: mnemosyne server can handle multiple users and agents
- **Persistence**: Memories survive across Pedantic Raven sessions

---

## Features in Detail

### Semantic Analysis

The analyzer runs automatically as you type (500ms debounce to avoid lag):

**Entity Extraction** - 6 types with context-aware classification
- Person: "Alice", "Dr. Smith"
- Place: "San Francisco", "Building 42"
- Thing: "database", "API"
- Concept: "authentication", "scalability"
- Organization: "Google", "Team Alpha"
- Technology: "PostgreSQL", "React"

**Relationship Detection** - Subject-predicate-object patterns
- Pattern: "PostgreSQL stores user data" → (PostgreSQL, stores, user data)
- Confidence scoring based on pattern strength
- Verb detection for predicate extraction

**Typed Holes** - Mark incomplete sections
- `??Architecture` - Priority 5, Complexity 3 (needs design decisions)
- `??Implementation:high` - Priority 8, Complexity 5 (urgent implementation gap)
- `!!SecurityReview` - Constraint marker (must satisfy security requirements)

**Dependencies** - Track references
- Imports: `import auth from "./auth"`
- Requires: `require database connection`
- References: `See [Architecture Doc]`

### mnemosyne Integration

**CRUD Operations** (via gRPC client):
```go
// Store memory with enrichment
client.StoreMemory(ctx, StoreMemoryOptions{
    Content: "System uses event sourcing for audit trail",
    Namespace: ProjectNamespace("myapp"),
    Importance: 8,
    Tags: []string{"architecture", "patterns"},
    EnrichWithLLM: true,
})

// Recall with hybrid search
results, _ := client.Recall(ctx, RecallOptions{
    Query: "authentication flow",
    MaxResults: 10,
    SemanticWeight: 0.7,  // Prefer semantic similarity
    FtsWeight: 0.2,       // Some full-text matching
    GraphWeight: 0.1,     // Consider graph connections
})
```

**Search Capabilities**:
- **Recall**: Hybrid search (semantic + full-text + graph)
- **SemanticSearch**: Pure embedding similarity (768d or 1536d vectors)
- **GraphTraverse**: Multi-hop graph exploration from seed nodes
- **GetContext**: Retrieve memory with all linked memories (configurable depth)

**Streaming**: Progressive results for long-running operations
- RecallStream: Results arrive as found
- ListMemoriesStream: Batch delivery for large result sets
- StoreMemoryStream: Progress updates ("enriching", "embedding", "indexing")

### Text Editor

**File Operations**:
- Open/save with atomic writes (temp file + rename)
- Dirty flag tracking for unsaved changes
- UTF-8 encoding support
- Error handling for permissions and missing files

**Search & Replace**:
- Literal and regex pattern matching
- Case sensitive/insensitive toggle
- Whole word matching
- Replace current match or all matches
- Full undo support for replacements

**Syntax Highlighting**:
- Token-based system (12 token types)
- Go: Keywords, types, functions, strings, comments, numbers, operators
- Markdown: Headers, code blocks, lists, links, bold, italic
- Extensible: Easy to add new language tokenizers
- Auto-detection by file extension or content

**Editing Features**:
- Undo/redo with full history
- Multi-line operations
- Cursor position tracking
- Line-based editing with buffer manager

---

## 🧠 Entity Extraction with GLiNER

Pedantic Raven now supports **GLiNER2** for ML-based entity extraction, providing significantly higher accuracy than traditional pattern matching while supporting custom entity types.

### Key Features

- **High Accuracy** - 85-95% accuracy vs 60-70% with pattern matching
- **Context-Aware** - Understands ambiguous text (e.g., "Apple" as company vs fruit)
- **Custom Types** - Define unlimited domain-specific entity types (api_endpoint, security_concern, etc.)
- **Zero-Shot** - Works on any domain without training
- **Automatic Fallback** - Gracefully falls back to pattern matcher if service unavailable
- **100% Local** - All processing happens on your machine

### Quick Start

```bash
# Start GLiNER service (one-time setup)
cd services/gliner
python -m venv venv && source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --host 127.0.0.1 --port 8765

# Run Pedantic Raven (in new terminal)
./pedantic_raven
```

GLiNER is optional and requires Python 3.9+, ~1GB RAM, and ~1GB disk space. See [docs/GLINER_INTEGRATION.md](docs/GLINER_INTEGRATION.md) for full documentation including configuration, deployment options, troubleshooting, and performance tuning.

---

## Architecture

### Technology Stack

- **Language**: Go 1.25+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm Architecture)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles)
- **RPC**: gRPC + Protocol Buffers (mnemosyne integration)

### Design Patterns

**Elm Architecture** (Model-Update-View):
- Immutable state updates
- Pure functions for rendering
- Commands for side effects
- Message-based communication

**Event-Driven**:
- PubSub event broker for component decoupling
- 40+ domain event types
- Non-blocking publish/subscribe

**Mode-Based UI**:
- 5 application modes (Edit, Explore, Analyze, Orchestrate, Collaborate)
- Each mode has its own layout and lifecycle
- Mode registry for switching with history

**Component Composition**:
- Hierarchical pane system (composite pattern)
- 5 layout modes (Focus, Standard, Analysis, Compact, Custom)
- Responsive design adapts to terminal size

### Project Structure

```
pedantic_raven/
├── main.go                    # Application entry point
├── internal/
│   ├── app/events/            # PubSub event system (18 tests)
│   ├── context/               # Context panel component (25 tests)
│   ├── editor/                # Text editor (78 tests)
│   │   ├── buffer/            # Buffer manager (52 tests)
│   │   ├── search/            # Search engine (35 tests)
│   │   ├── semantic/          # Semantic analyzer (63 tests)
│   │   └── syntax/            # Syntax highlighting (31 tests)
│   ├── layout/                # Layout engine (34 tests)
│   ├── memorydetail/          # Memory detail view (19 tests)
│   ├── memorygraph/           # Graph visualization (134 tests)
│   ├── memorylist/            # Memory list component (13 tests)
│   ├── mnemosyne/             # mnemosyne RPC client (66 tests)
│   ├── modes/                 # Mode registry and modes (27 tests)
│   ├── overlay/               # Overlay system (25 tests)
│   ├── palette/               # Command palette (19 tests)
│   └── terminal/              # Terminal component (38 tests)
├── proto/mnemosyne/v1/        # Protobuf schemas
├── docs/
│   ├── USAGE.md               # User guide and keyboard reference
│   ├── CHANGELOG.md           # Development history
│   └── PHASE*.md              # Phase completion summaries
├── spec.md                    # Technical specification
└── ROADMAP.md                 # Project timeline and milestones
```

---

## Development

### Running Tests

```bash
# All tests
go test ./...

# With verbose output
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Specific package
go test ./internal/editor/...
```

### Test Coverage

| Package | Tests | Coverage | Focus |
|---------|-------|----------|-------|
| app/events | 18 | ~70% | Event broker, pub/sub |
| context | 25 | ~80% | Context panel rendering |
| editor | 78 | ~85% | Text editing, file ops |
| editor/buffer | 52 | ~85% | Buffer management, undo/redo |
| editor/search | 35 | ~90% | Search and replace |
| editor/semantic | 63 | ~90% | Semantic analysis |
| editor/syntax | 31 | ~85% | Syntax highlighting |
| layout | 34 | ~65% | Layout engine, panes |
| memorydetail | 19 | ~85% | Memory detail view |
| memorygraph | 134 | ~88% | Graph visualization, layout |
| memorylist | 13 | ~85% | Memory list component |
| mnemosyne | 66 | ~95% | gRPC client, CRUD, search |
| modes | 27 | ~90% | Mode registry, switching, ExploreMode |
| overlay | 25 | ~70% | Overlays, dialogs |
| palette | 19 | ~88% | Command palette, fuzzy search |
| terminal | 38 | ~80% | Terminal component, execution |
| **Total** | **754** | **~65%** | **Passing** |

### Building

```bash
# Development build
go build -o pedantic_raven .

# Optimized production build
go build -ldflags="-s -w" -o pedantic_raven .

# Run without building
go run main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o pedantic_raven-linux .
```

### Code Standards

- Go fmt/vet/lint clean
- Tests for all new code (target 70%+ coverage)
- Document exported types and functions
- Meaningful names (no single-letter variables except loop counters)
- Small focused functions (<50 lines preferred)

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Implement your feature
5. Run tests (`go test ./...`)
6. Commit changes (`git commit -m 'Add amazing feature'`)
7. Push to branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

---

## Documentation

### Developer Guides

- **[.claude/CLAUDE.md](.claude/CLAUDE.md)** - Development guide for contributors
  - Quick start and architecture overview
  - Development workflow and testing strategy
  - Common tasks and patterns
  - Release process and repository organization
  - Context-efficient reference for human and AI developers

- **[.claude/AGENT_GUIDE.md](.claude/AGENT_GUIDE.md)** - Guide for autonomous development
  - Decision trees for development workflows
  - Repository organization and tidying protocols
  - Documentation update procedures
  - Release management processes
  - Code generation templates and patterns

### User Documentation

- **[USAGE.md](docs/USAGE.md)** - Complete user guide
  - All keyboard shortcuts
  - Feature walkthroughs
  - Examples and workflows
  - Tips and tricks

- **[ROADMAP.md](ROADMAP.md)** - Project roadmap
  - Phase breakdown and timeline
  - Completed features
  - Planned features
  - Success metrics

### Technical Documentation

- **[spec.md](spec.md)** - Technical specification
  - Requirements (functional and non-functional)
  - Architecture details
  - Interfaces and contracts
  - Design decisions

- **[CHANGELOG.md](docs/CHANGELOG.md)** - Development history
  - Phase summaries
  - Feature additions
  - Technical achievements
  - Lessons learned

### Phase Summaries

- **[PHASE1_COMPLETE.md](docs/PHASE1_COMPLETE.md)** - Foundation (87 tests)
- **[PHASE2_SUMMARY.md](docs/PHASE2_SUMMARY.md)** - Semantic Analysis (291 tests)
- **[PHASE3_SUMMARY.md](docs/PHASE3_SUMMARY.md)** - Advanced Editor (424 tests)
- **[PHASE4.1_SUMMARY.md](docs/PHASE4.1_SUMMARY.md)** - mnemosyne Client (461 tests)

---

## Statistics

**Current Metrics** (as of latest commit):
- **Tests**: 754 passing
- **Code**: ~34,000 lines of Go
- **Coverage**: ~65% average
- **Phases**: 4.4 of 8 complete (~35% of planned features)

**Performance**:
- Render time: <16ms target (60 FPS)
- Semantic analysis: ~500ms for typical files
- Memory usage: ~10MB typical
- Startup time: <100ms

**Commits**: 600+ commits since project start (2025-01-11)

---

## Troubleshooting

### Terminal Too Small

Pedantic Raven auto-switches to compact layout for terminals <120x30. For best experience, use at least 120x30. Check your terminal size:

```bash
echo "Cols: $(tput cols), Rows: $(tput lines)"
```

### Rendering Issues

Ensure your terminal supports:
- 256 colors: `echo $TERM` should show `xterm-256color` or similar
- UTF-8 encoding
- ANSI escape sequences

### Tests Failing

If you see test failures:
1. Update dependencies: `go mod download`
2. Clear test cache: `go clean -testcache`
3. Run specific failing package: `go test -v ./internal/memorygraph`

### Performance Issues

If experiencing lag:
- Close other terminal applications
- Check CPU usage (`top` or `htop`)
- Reduce terminal size temporarily
- Update to latest version

### Connection to mnemosyne

If can't connect to mnemosyne server:
```bash
# Check if mnemosyne is running
netstat -an | grep 50051

# Start mnemosyne server (from mnemosyne project)
cd ../mnemosyne
cargo run --bin mnemosyne-rpc
```

---

## Roadmap Summary

**Completed** (~35%):
- ✅ Foundation (event system, layout, modes, overlays, palette)
- ✅ Edit Mode (semantic analysis, context panel, terminal)
- ✅ Advanced editing (undo/redo, files, search, syntax highlighting)
- ✅ mnemosyne RPC client (CRUD, search, streaming)
- ✅ Explore Mode (memory list, detail view, graph visualization)

**Planned**:
- 📋 Analyze Mode - Statistical analysis and insights
- 📋 Orchestrate Mode - Multi-agent coordination
- 📋 Collaborate Mode - Live multi-user editing
- 📋 Production Release - Performance optimization, docs, packaging

**Timeline**: 6-8 months total (started Jan 2025, targeting Q2 2025)

See [ROADMAP.md](ROADMAP.md) for detailed breakdown.

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Links

**Project**: [github.com/rand/pedantic_raven](https://github.com/rand/pedantic_raven)
**mnemosyne**: [github.com/rand/mnemosyne](https://github.com/rand/mnemosyne)
**Bubble Tea**: [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)

---

**Current Phase**: 4.4 (Graph Visualization - Complete)
**Last Updated**: 2025-11-08
**Status**: Active Development 🚧
