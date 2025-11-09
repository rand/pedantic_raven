# Pedantic Raven Agent Guide

**Last Updated**: 2025-11-09
**Version**: 0.4.4
**For**: Agents working in the Pedantic Raven codebase

---

## Quick Links

- [Project Overview](#project-overview)
- [Quick Start](#quick-start-for-agents)
- [Architecture](#architecture)
- [Repository Structure](#repository-structure)
- [Development Workflows](#development-workflows)
- [Testing Strategy](#testing-strategy)
- [Documentation Management (Zensical)](#documentation-management-zensical)
- [Common Tasks](#common-tasks)
- [Resources](#resources)

---

## Project Overview

### Purpose
Pedantic Raven is a **terminal-based context engineering system** with semantic memory editor and real-time entity extraction. It provides a rich TUI for creating structured AI context documents with live semantic analysis and integration to the mnemosyne memory system.

### Key Capabilities
- **Edit Mode**: Full-featured text editor with syntax highlighting, search/replace, undo/redo
- **Semantic Analysis**: Real-time extraction of entities (6 types), relationships, typed holes, dependencies
- **Context Panel**: Live semantic results with filtering and navigation
- **Integrated Terminal**: Command execution with history and shell pass-through
- **mnemosyne Integration**: Full CRUD, hybrid search, graph traversal via gRPC
- **Explore Mode**: Memory workspace with list, detail, and graph visualization
- **GLiNER Support**: ML-based entity extraction with 85-95% accuracy (optional)

### Current Status
- **Version**: 0.4.4 (Phase 4.4 - Graph Visualization Complete)
- **Test Status**: 754 tests passing (~65% coverage)
- **Build Time**: <1s incremental, ~3-5s clean build
- **Language**: Go 1.25.2, Bubble Tea TUI framework

---

## Quick Start for Agents

### Build Commands
```bash
# Check compilation
go build -o pedantic_raven .

# Run without building
go run main.go

# Optimized production build
go build -ldflags="-s -w" -o pedantic_raven .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o pedantic_raven-linux .

# Run tests
go test ./...                    # All tests
go test ./... -v                 # Verbose output
go test ./... -coverprofile=coverage.out  # With coverage
go test ./internal/editor/...    # Specific package

# Check code quality
go fmt ./...                     # Format code
go vet ./...                     # Static analysis
```

### Common Development Workflows
```bash
# Start pedantic_raven
./pedantic_raven

# Start with GLiNER service (optional)
cd services/gliner
python -m venv venv && source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --host 127.0.0.1 --port 8765

# Start mnemosyne RPC server (for memory integration)
cd ../mnemosyne
cargo run --bin mnemosyne-rpc --features rpc
```

### Where to Find Things
- **Source code**: `internal/`
- **Tests**: `internal/*/` (co-located with source)
- **Documentation**: `docs/`
- **Protobuf schemas**: `proto/mnemosyne/v1/`
- **Examples**: `examples/`
- **Configuration**: `config.toml`, `zensical.toml`

---

## Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────┐
│                     User (Terminal)                      │
│                   (Keyboard/Mouse)                       │
└────────────────────┬────────────────────────────────────┘
                     │ Bubble Tea messages
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   Main Application                       │
│              (Elm Architecture: Model-Update-View)       │
└──────┬──────────────────┬───────────────────┬───────────┘
       │                  │                   │
       ▼                  ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌─────────────────┐
│  Mode System │   │  Event Broker│   │ Layout Engine   │
│  (5 modes)   │   │  (PubSub)    │   │ (5 layouts)     │
└──────────────┘   └──────────────┘   └─────────────────┘
       │                  │                   │
       ▼                  ▼                   ▼
┌─────────────────────────────────────────────────────────┐
│                    Components                            │
│  Editor | Context Panel | Terminal | Memory List/Detail │
│  Overlay | Palette | Graph Visualization                │
└──────┬──────────────────┬───────────────────┬───────────┘
       │                  │                   │
       ▼                  ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌─────────────────┐
│   Services   │   │  Semantic    │   │  mnemosyne      │
│   (GLiNER)   │   │  Analysis    │   │  gRPC Client    │
└──────────────┘   └──────────────┘   └─────────────────┘
```

### Component Diagram

**Key Subsystems**:
- **Modes** (`internal/modes/`): 5 application modes (Edit, Explore, Analyze, Orchestrate, Collaborate)
- **Editor** (`internal/editor/`): Text editing with buffer, search, syntax highlighting, semantic analysis
- **Context Panel** (`internal/context/`): Live semantic results display
- **Terminal** (`internal/terminal/`): Integrated command execution
- **Memory Components** (`internal/memory*/`): List, detail view, graph visualization
- **Layout Engine** (`internal/layout/`): Responsive pane system with 5 layout modes
- **Event System** (`internal/app/events/`): PubSub event broker for component decoupling
- **mnemosyne Client** (`internal/mnemosyne/`): gRPC client for memory operations
- **GLiNER Integration** (`internal/gliner/`): ML-based entity extraction service

### Data Flow Overview

**Editing Flow**:
```
User types in Editor
      ↓
Buffer Manager updates text (undo/redo tracking)
      ↓
Semantic Analyzer runs after 500ms debounce
      ↓
Entities, relationships, typed holes extracted
      ↓
Context Panel displays results (entities, relationships, etc.)
```

**Memory Flow**:
```
User stores content → mnemosyne gRPC Client
                   → StoreMemory RPC
                   → LLM enrichment (summary, tags, keywords)
                   → Embedding generation (768d/1536d vectors)
                   → mnemosyne server stores in LibSQL

User recalls memories → Recall RPC (hybrid search)
                     → FTS5 + vector + graph scoring
                     → Ranked results returned
                     → Memory List displays results
```

---

## Repository Structure

```
pedantic_raven/
├── main.go                    # Application entry point
├── go.mod                     # Go dependencies
├── go.sum                     # Dependency checksums
├── config.toml                # Application configuration
├── zensical.toml              # Documentation site configuration
│
├── internal/                  # Source code (754 tests)
│   ├── app/events/            # PubSub event system (18 tests)
│   ├── analyze/               # Analyze mode components (22 tests)
│   ├── context/               # Context panel component (25 tests)
│   ├── editor/                # Text editor (78 tests)
│   │   ├── buffer/            # Buffer manager (52 tests)
│   │   ├── search/            # Search engine (35 tests)
│   │   ├── semantic/          # Semantic analyzer (63 tests)
│   │   └── syntax/            # Syntax highlighting (31 tests)
│   ├── gliner/                # GLiNER service integration (12 tests)
│   ├── integration/           # Integration tests
│   ├── layout/                # Layout engine (34 tests)
│   ├── memorydetail/          # Memory detail view (19 tests)
│   ├── memorygraph/           # Graph visualization (134 tests)
│   ├── memorylist/            # Memory list component (13 tests)
│   ├── mnemosyne/             # mnemosyne RPC client (66 tests)
│   ├── modes/                 # Mode registry and modes (27 tests)
│   ├── orchestrate/           # Orchestrate mode components (39 tests)
│   ├── overlay/               # Overlay system (25 tests)
│   ├── palette/               # Command palette (19 tests)
│   └── terminal/              # Terminal component (38 tests)
│
├── proto/mnemosyne/v1/        # Protobuf schemas
│   └── mnemosyne.proto        # mnemosyne gRPC service definitions
│
├── docs/                      # Documentation (130+ files)
│   ├── index.md               # Documentation homepage
│   ├── whitepaper.md          # Technical whitepaper
│   ├── USAGE.md               # User guide
│   ├── DEVELOPMENT.md         # Developer guide
│   ├── CONTRIBUTING.md        # Contribution guidelines
│   ├── TESTING.md             # Testing framework
│   ├── PERFORMANCE.md         # Performance benchmarks
│   ├── STYLE_GUIDE.md         # Code style guide
│   ├── architecture.md        # Architecture documentation
│   ├── *-guide.md             # Mode-specific guides
│   ├── PHASE*.md              # Phase completion summaries
│   ├── assets/                # Images, diagrams, icons
│   ├── stylesheets/           # CSS for documentation site
│   ├── javascripts/           # JS for documentation site
│   └── overrides/             # Zensical theme overrides
│
├── .github/workflows/         # GitHub Actions
│   ├── docs.yml               # Zensical documentation deployment
│   ├── test.yml               # Test suite runner
│   └── deploy-pages.yml       # GitHub Pages deployment
│
├── examples/                  # Example files
├── services/gliner/           # GLiNER service (Python)
├── README.md                  # Project overview
├── ROADMAP.md                 # Project roadmap
├── LICENSE                    # MIT license
└── spec.md                    # Technical specification
```

### Key Files

| File | Purpose | Update When |
|------|---------|-------------|
| `README.md` | User documentation, quick start | Features added, status changes |
| `docs/whitepaper.md` | Technical deep dive | Architecture changes, algorithms |
| `AGENT_GUIDE.md` | Agent workflow guidance | Workflow changes, new protocols |
| `docs/CHANGELOG.md` | Version history | Every semantic version change |
| `docs/CONTRIBUTING.md` | Contribution guidelines | Process changes |
| `ROADMAP.md` | Project timeline | Milestones completed/updated |
| `go.mod` | Dependencies | Dependency updates |
| `config.toml` | Application config | Config changes |
| `zensical.toml` | Documentation site config | Site structure/theme changes |

---

## Development Workflows

### Branching Strategy

**ALWAYS use feature branches**:
```bash
git checkout -b feature/gliner-improvements
git checkout -b fix/memory-leak
git checkout -b docs/architecture-update
```

**Branch Naming**:
- `feature/`: New functionality
- `fix/`: Bug fixes
- `refactor/`: Code restructuring
- `docs/`: Documentation updates
- `perf/`: Performance improvements

### Commit Protocol

**CRITICAL**: Commit BEFORE testing, never test uncommitted code.

```bash
# 1. Make changes
# 2. Commit changes
git add .
git commit -m "Add force-directed graph layout algorithm"
git log -1 --oneline

# 3. Run tests
go test ./...

# 4. If tests fail: Fix → Commit → Re-test
```

**Commit Message Format**:
```
<type>: <short summary>

<optional detailed explanation>

<optional breaking changes>
```

**Types**: `feat`, `fix`, `refactor`, `perf`, `docs`, `test`, `chore`

**Examples**:
```
feat: Add semantic search to memory list
fix: Resolve crash when terminal too small
perf: Optimize graph layout algorithm (2x faster)
docs: Update Zensical documentation deployment guide
```

### Pull Requests

```bash
# Push branch
git push -u origin feature/gliner-improvements

# Create PR (GitHub CLI)
gh pr create --title "Add GLiNER integration for entity extraction" \
  --body "Implements ML-based entity extraction with 85-95% accuracy..."
```

**PR Checklist**:
- [ ] All tests pass (`go test ./...`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] Documentation updated (README, docs/, whitepaper if applicable)
- [ ] CHANGELOG.md updated if user-facing changes

---

## Testing Strategy

### Test Organization

**Unit Tests** (co-located with source):
- Test individual functions and methods
- Mock external dependencies
- Fast execution (<1s for most packages)

**Integration Tests** (`internal/integration/`):
- Test component interactions
- Real gRPC client (with mock server)
- Test event system, mode switching, memory operations

**Manual Tests**:
- TUI visual regression testing
- Terminal size compatibility (80x24 to 200x60)
- Color scheme validation (256-color terminals)

### Current Test Status

**Passing**: 754 tests
**Coverage**: ~65% average

**Coverage by Package**:

| Package | Tests | Coverage | Focus |
|---------|-------|----------|-------|
| app/events | 18 | ~70% | Event broker, pub/sub |
| analyze | 22 | ~75% | Statistical analysis |
| context | 25 | ~80% | Context panel rendering |
| editor | 78 | ~85% | Text editing, file ops |
| editor/buffer | 52 | ~85% | Buffer management, undo/redo |
| editor/search | 35 | ~90% | Search and replace |
| editor/semantic | 63 | ~90% | Semantic analysis |
| editor/syntax | 31 | ~85% | Syntax highlighting |
| gliner | 12 | ~80% | GLiNER service integration |
| layout | 34 | ~65% | Layout engine, panes |
| memorydetail | 19 | ~85% | Memory detail view |
| memorygraph | 134 | ~88% | Graph visualization, layout |
| memorylist | 13 | ~85% | Memory list component |
| mnemosyne | 66 | ~95% | gRPC client, CRUD, search |
| modes | 27 | ~90% | Mode registry, ExploreMode |
| orchestrate | 39 | ~82% | Orchestration components |
| overlay | 25 | ~70% | Overlays, dialogs |
| palette | 19 | ~88% | Command palette, fuzzy search |
| terminal | 38 | ~80% | Terminal component, execution |

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

# Specific test
go test -v ./internal/memorygraph -run TestForceDirectedLayout

# With race detection
go test ./... -race

# Clear test cache
go clean -testcache
```

### Test Coverage Targets

- Critical path: 85%+
- Editor/semantic: 90%+
- UI components: 70%+
- Overall: 70%+

---

## Documentation Management (Zensical)

### Documentation Structure

Pedantic Raven uses **Zensical** for documentation site generation with shared infrastructure.

```
docs/
├── index.md                   # Documentation homepage
├── whitepaper.md              # Technical whitepaper
├── USAGE.md                   # User guide
├── DEVELOPMENT.md             # Developer guide
├── CONTRIBUTING.md            # Contribution guidelines
├── TESTING.md                 # Testing framework
├── PERFORMANCE.md             # Performance benchmarks
├── STYLE_GUIDE.md             # Code style guide
├── architecture.md            # Architecture documentation
├── edit-mode-guide.md         # Edit mode guide
├── analyze-mode-guide.md      # Analyze mode guide
├── orchestrate-mode-guide.md  # Orchestrate mode guide
├── CHANGELOG.md               # Version history
├── PHASE*.md                  # Phase completion summaries
│
├── assets/                    # Images, diagrams, icons
│   ├── diagrams/              # Architecture diagrams
│   ├── favicon.svg            # Site favicon (teal glyph ⟡)
│   └── og-image.png           # Social sharing image
│
├── stylesheets/               # CSS for documentation site
│   ├── shared.css             # Shared base styles (from shared-docs-base)
│   └── teal.css               # Teal theme overrides (#16A085)
│
├── javascripts/               # JavaScript for documentation site
│   └── theme.js               # Theme switching logic
│
└── overrides/                 # Zensical theme overrides
```

### Shared Infrastructure

**Location**: `/Users/rand/src/shared-docs-base`

**Shared Resources**:
- **Base Styles**: `shared.css` - Typography, layout, component styles
- **Theme System**: Color schemes, dark/light mode
- **Scripts**: Documentation generation utilities
- **Configuration**: Zensical base configuration

**Updating Shared CSS**:
1. Edit `/Users/rand/src/shared-docs-base/stylesheets/shared.css`
2. Update version number in comment header
3. Copy to project: `cp /Users/rand/src/shared-docs-base/stylesheets/shared.css docs/stylesheets/`
4. Test locally: `zensical serve`
5. Commit changes in **both** repositories

### Building and Serving Documentation

**Prerequisites**:
- Python 3.12+
- UV package manager (`curl -LsSf https://astral.sh/uv/install.sh | sh`)
- Zensical (`uv add zensical`)

**Local Development**:
```bash
# Serve documentation locally (with live reload)
/Users/rand/src/shared-docs-base/.venv/bin/zensical serve

# Build static site
/Users/rand/src/shared-docs-base/.venv/bin/zensical build --clean

# Preview built site
cd site && python -m http.server 8000
```

**Configuration** (`zensical.toml`):
```toml
[project]
site_name = "Pedantic Raven"
site_url = "https://rand.github.io/pedantic_raven/"
site_description = "Terminal-Based Context Engineering - Semantic memory editor with real-time entity extraction"
docs_dir = "docs"
site_dir = "site"

[project.theme]
variant = "modern"
favicon = "assets/favicon.svg"
primary = "custom"    # Teal (#16A085)
accent = "custom"

extra_css = [
  "stylesheets/shared.css",
  "stylesheets/teal.css",
]

extra_javascript = [
  "javascripts/theme.js",
]
```

### Theme

**Primary Color**: Teal (#16A085)
**Glyph**: ⟡ (semantic boundary marker)
**Fonts**:
- Text: Geist
- Code: JetBrains Mono

**Custom CSS** (`docs/stylesheets/teal.css`):
```css
:root {
  --md-primary-fg-color: #16A085;
  --md-accent-fg-color: #1ABC9C;
}

[data-md-color-scheme="slate"] {
  --md-primary-fg-color: #1ABC9C;
  --md-accent-fg-color: #16A085;
}
```

### Deployment

**Automatic via GitHub Actions**:
- **Trigger**: Push to `main` branch
- **Workflow**: `.github/workflows/docs.yml`
- **Build**: Zensical generates static site
- **Deploy**: GitHub Pages serves from `site/` directory
- **URL**: https://rand.github.io/pedantic_raven/
- **Time**: ~2-3 minutes from push to live

**GitHub Actions Workflow** (`.github/workflows/docs.yml`):
```yaml
name: Deploy Zensical Docs

on:
  push:
    branches:
      - main
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: '3.12'
      - name: Install UV and Zensical
        run: |
          curl -LsSf https://astral.sh/uv/install.sh | sh
          echo "$HOME/.cargo/bin" >> $GITHUB_PATH
          uv init --no-workspace || true
          uv add zensical
      - name: Build site
        run: uv run zensical build --clean
      - uses: actions/upload-pages-artifact@v3
        with:
          path: ./site

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - uses: actions/deploy-pages@v4
```

### Documentation Update Workflow

#### Standard Update Process
```bash
# 1. Edit documentation files
vim docs/architecture.md
vim docs/whitepaper.md

# 2. Test locally
/Users/rand/src/shared-docs-base/.venv/bin/zensical serve
# Open http://localhost:8000 in browser

# 3. Verify rendering and navigation
# - Check all links work
# - Verify images load
# - Test search functionality
# - Validate syntax highlighting

# 4. Commit changes
git add docs/
git commit -m "docs: Update architecture documentation

- Added force-directed graph algorithm details
- Updated performance benchmarks
- Added memory workspace screenshots"

# 5. Push to main
git push origin main

# 6. Verify GitHub Pages deployment
# Wait 2-3 minutes, then visit:
# https://rand.github.io/pedantic_raven/
```

#### Quick Documentation Update (Minor Changes)
```bash
# For typos, clarifications, small improvements
vim docs/USAGE.md
git add docs/USAGE.md
git commit -m "docs: Fix typo in keyboard shortcuts"
git push
```

#### Major Documentation Overhaul
```bash
# 1. Create feature branch
git checkout -b docs/major-restructure

# 2. Make all documentation updates
# (multiple files, restructuring, etc.)

# 3. Test locally
/Users/rand/src/shared-docs-base/.venv/bin/zensical serve

# 4. Commit and create PR
git add docs/
git commit -m "docs: Major documentation restructure

- Reorganized guides by user role
- Added troubleshooting section
- Updated all cross-references
- Improved navigation structure"

git push -u origin docs/major-restructure
gh pr create --title "Major Documentation Restructure"
```

### Documentation Update Triggers

**ALWAYS update documentation when**:

| Change Type | Primary Docs | Secondary Docs | Notes |
|-------------|--------------|----------------|-------|
| **New Feature** | README.md, docs/USAGE.md | CHANGELOG.md, ROADMAP.md | Add to feature list |
| **Architecture Change** | docs/architecture.md, docs/whitepaper.md | README.md | Update diagrams |
| **API Modification** | docs/DEVELOPMENT.md | docs/architecture.md | Update code examples |
| **Workflow Change** | AGENT_GUIDE.md, docs/DEVELOPMENT.md | README.md | Update command sequences |
| **Dependency Addition** | README.md, docs/DEVELOPMENT.md | docs/CONTRIBUTING.md | Update installation |
| **Release** | CHANGELOG.md, ROADMAP.md | README.md (version) | Follow release protocol |
| **Configuration Change** | README.md, config.toml | docs/DEVELOPMENT.md | Update examples |
| **Bug Fix** | CHANGELOG.md | docs/troubleshooting (if applicable) | Document fix |
| **Performance Improvement** | README.md, docs/PERFORMANCE.md | CHANGELOG.md | Update benchmarks |

### Troubleshooting

| Issue | Diagnosis | Solution |
|-------|-----------|----------|
| **Site not updating** | GitHub Actions failed | Check Actions tab, review build logs |
| **Broken links** | Relative paths incorrect | Use `[text](file.md)` for same-directory links |
| **CSS not loading** | Path incorrect in zensical.toml | Verify `extra_css` paths relative to `docs/` |
| **Search not working** | Zensical build incomplete | Rebuild with `zensical build --clean` |
| **Images not showing** | Path incorrect or missing | Verify images in `docs/assets/` |
| **Theme not applied** | CSS cache issue | Hard refresh (Ctrl+Shift+R) or clear cache |
| **Local serve fails** | UV/Zensical not installed | Run `uv add zensical` |

### Shared CSS Updates Process

**When updating shared styles**:

1. **Edit base styles**:
```bash
vim /Users/rand/src/shared-docs-base/stylesheets/shared.css
```

2. **Update version number** in CSS header:
```css
/* Shared Documentation Styles v1.2.0 */
```

3. **Copy to pedantic_raven**:
```bash
cp /Users/rand/src/shared-docs-base/stylesheets/shared.css \
   /Users/rand/src/pedantic_raven/docs/stylesheets/
```

4. **Test locally**:
```bash
cd /Users/rand/src/pedantic_raven
/Users/rand/src/shared-docs-base/.venv/bin/zensical serve
```

5. **Commit in both repositories**:
```bash
# In shared-docs-base
cd /Users/rand/src/shared-docs-base
git add stylesheets/shared.css
git commit -m "style: Update typography scale"
git push

# In pedantic_raven
cd /Users/rand/src/pedantic_raven
git add docs/stylesheets/shared.css
git commit -m "docs: Update shared styles to v1.2.0"
git push
```

---

## Common Tasks

### Adding a New Mode

```bash
# 1. Create mode struct in internal/modes/
vim internal/modes/mymode.go

# 2. Implement Mode interface
# - Init() tea.Cmd
# - Update(tea.Msg) (tea.Model, tea.Cmd)
# - View() string
# - Name() string
# - Focus() / Blur()

# 3. Register in mode registry
vim internal/modes/registry.go

# 4. Add mode-specific components (if needed)
mkdir internal/mymode
vim internal/mymode/component.go

# 5. Add tests
vim internal/modes/mymode_test.go

# 6. Update documentation
vim docs/mymode-guide.md
vim README.md  # Add to feature list
vim ROADMAP.md  # Mark milestone complete
```

### Adding Semantic Analysis Features

```bash
# 1. Edit semantic analyzer
vim internal/editor/semantic/analyzer.go

# 2. Add new extraction method
# func (a *Analyzer) ExtractNewFeature(text string) []NewFeature

# 3. Update SemanticResults struct
vim internal/editor/semantic/types.go

# 4. Update context panel rendering
vim internal/context/context.go

# 5. Add tests
vim internal/editor/semantic/analyzer_test.go

# 6. Commit and test
git add .
git commit -m "feat: Add temporal entity extraction"
go test ./internal/editor/semantic/...
```

### Updating mnemosyne Integration

```bash
# 1. Update protobuf schemas (if needed)
vim proto/mnemosyne/v1/mnemosyne.proto

# 2. Regenerate Go code
protoc --go_out=. --go-grpc_out=. proto/mnemosyne/v1/mnemosyne.proto

# 3. Update client methods
vim internal/mnemosyne/client.go

# 4. Add tests
vim internal/mnemosyne/client_test.go

# 5. Update memory components
vim internal/memorylist/memorylist.go
vim internal/memorydetail/detail.go

# 6. Commit and test
git add .
git commit -m "feat: Add graph traversal with depth limit"
go test ./internal/mnemosyne/...
```

### Performance Optimization

```bash
# 1. Profile the application
go test -cpuprofile=cpu.prof -bench=. ./internal/memorygraph/

# 2. Analyze profile
go tool pprof cpu.prof
# (pprof) top20
# (pprof) list FunctionName

# 3. Optimize hot paths
vim internal/memorygraph/layout.go

# 4. Benchmark before/after
go test -bench=BenchmarkForceDirectedLayout ./internal/memorygraph/ > old.txt
# Make changes
go test -bench=BenchmarkForceDirectedLayout ./internal/memorygraph/ > new.txt
benchcmp old.txt new.txt

# 5. Update performance documentation
vim docs/PERFORMANCE.md

# 6. Commit
git add .
git commit -m "perf: Optimize graph layout algorithm (2x faster)"
```

### Adding Tests

```bash
# 1. Create test file (if not exists)
vim internal/mypackage/myfile_test.go

# 2. Write test cases
# func TestMyFunction(t *testing.T) { ... }

# 3. Use table-driven tests for comprehensive coverage
tests := []struct {
    name     string
    input    string
    expected int
    wantErr  bool
}{
    {"empty string", "", 0, false},
    {"single word", "hello", 1, false},
    {"invalid input", "???", 0, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result, err := MyFunction(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("wantErr %v, got %v", tt.wantErr, err)
        }
        if result != tt.expected {
            t.Errorf("expected %d, got %d", tt.expected, result)
        }
    })
}

# 4. Run tests
go test ./internal/mypackage/ -v

# 5. Check coverage
go test ./internal/mypackage/ -coverprofile=coverage.out
go tool cover -func=coverage.out

# 6. Commit
git add .
git commit -m "test: Add comprehensive tests for MyFunction"
```

### Updating Dependencies

```bash
# 1. Check for outdated dependencies
go list -u -m all

# 2. Update specific dependency
go get github.com/charmbracelet/bubbletea@latest
go mod tidy

# 3. Update all dependencies
go get -u ./...
go mod tidy

# 4. Run tests to verify compatibility
go test ./...

# 5. Commit
git add go.mod go.sum
git commit -m "chore: Update Bubble Tea to v1.3.10"
```

---

## Release Management

### Semantic Versioning

Pedantic Raven follows [Semantic Versioning 2.0.0](https://semver.org/):

**Version Format**: `MAJOR.MINOR.PATCH` (e.g., `0.4.4`)

| Component | Increment When | Example |
|-----------|----------------|---------|
| **MAJOR** | Breaking changes, incompatible API | `0.x.x` → `1.0.0` |
| **MINOR** | New features, backward-compatible | `0.4.x` → `0.5.0` |
| **PATCH** | Bug fixes, backward-compatible | `0.4.4` → `0.4.5` |

**Current Version**: 0.4.4 (as of 2025-11-09)

### Release Triggers

#### MINOR Version Release (0.x.0)

**Trigger when**:
- New mode completed (Analyze, Orchestrate, Collaborate)
- Major feature added (graph visualization, GLiNER integration)
- Significant UI improvements
- New mnemosyne capabilities

**Examples**:
- v0.4.0: Explore Mode with memory workspace
- v0.5.0: Analyze Mode (planned)

**Process**:
1. Update CHANGELOG.md with feature descriptions
2. Update README.md feature list and status
3. Update ROADMAP.md milestones
4. Full test suite pass (`go test ./...`)
5. Create release tag and GitHub release

#### PATCH Version Release (0.x.y)

**Trigger when**:
- Bug fixes
- Performance improvements
- Documentation updates
- Minor UI tweaks

**Examples**:
- v0.4.4: Graph visualization bug fixes
- v0.4.5: Performance optimization (planned)

**Process**:
1. Update CHANGELOG.md with fix descriptions
2. Verify tests pass
3. Quick release (minimal ceremony)

### Release Process

```bash
# 1. Verify all tests pass
go test ./...
go fmt ./...
go vet ./...

# 2. Update CHANGELOG.md
vim docs/CHANGELOG.md
# ## [0.5.0] - 2025-11-15
# ### Added
# - Analyze mode with statistical analysis
# - Performance metrics visualization
# ### Fixed
# - Memory leak in graph layout
# ### Changed
# - Improved terminal rendering performance

# 3. Update README.md version
vim README.md
# **Current Phase**: 5.0 (Analyze Mode - Complete)
# **Last Updated**: 2025-11-15

# 4. Update ROADMAP.md
vim ROADMAP.md
# **Completed**:
# - ✅ Phase 5: Analyze Mode

# 5. Commit version bump
git add docs/CHANGELOG.md README.md ROADMAP.md
git commit -m "chore: Bump version to 0.5.0

Prepare for v0.5.0 release with Analyze Mode."

# 6. Create annotated tag
git tag -a v0.5.0 -m "Release v0.5.0: Analyze Mode

Major features:
- Statistical analysis of semantic content
- Performance metrics visualization
- Enhanced memory insights

See CHANGELOG.md for complete list of changes."

# 7. Build release binary
go build -ldflags="-s -w" -o pedantic_raven .

# 8. Verify binary
./pedantic_raven --version  # (if version flag implemented)

# 9. Push tag to GitHub
git push origin main
git push origin v0.5.0

# 10. Create GitHub release
gh release create v0.5.0 \
  --title "v0.5.0: Analyze Mode" \
  --notes "$(cat <<'EOF'
# v0.5.0: Analyze Mode

## Highlights
- **Analyze Mode**: Statistical analysis of semantic content
- **Performance Metrics**: Visual insights into content patterns
- **Memory Insights**: Enhanced understanding of stored memories

## Added
- Statistical analysis engine with 15+ metrics
- Performance visualization with charts
- Memory pattern detection

## Fixed
- Memory leak in graph layout algorithm
- Rendering issues on small terminals

## Changed
- Improved terminal rendering performance (20% faster)

See [CHANGELOG.md](https://github.com/rand/pedantic_raven/blob/main/docs/CHANGELOG.md) for complete details.

---
**Installation**:
```bash
go install github.com/rand/pedantic-raven@v0.5.0
# or download binary from release assets
```
EOF
)" \
  ./pedantic_raven#pedantic_raven-v0.5.0-$(uname -s)-$(uname -m)

# 11. Verify release
gh release view v0.5.0
open https://github.com/rand/pedantic_raven/releases/tag/v0.5.0

# 12. Verify documentation site updated
# Wait 2-3 minutes for GitHub Actions
open https://rand.github.io/pedantic_raven/
```

---

## Troubleshooting

### Common Build Errors

**Error**: `module requires Go 1.25 or later`
**Fix**: Update Go: `brew install go` (macOS) or download from golang.org

**Error**: `cannot find package`
**Fix**: `go mod download && go mod tidy`

**Error**: `protoc: command not found`
**Fix**: Install protobuf compiler: `brew install protobuf` (macOS)

### Runtime Issues

**Terminal too small**:
- Pedantic Raven auto-switches to compact layout for <120x30
- Check size: `echo "Cols: $(tput cols), Rows: $(tput lines)"`
- Resize terminal or run in larger window

**Connection to mnemosyne fails**:
```bash
# Check if mnemosyne-rpc is running
netstat -an | grep 50051

# Start mnemosyne RPC server
cd ../mnemosyne
cargo run --bin mnemosyne-rpc --features rpc
```

**GLiNER service not available**:
```bash
# Check if GLiNER is running
curl http://127.0.0.1:8765/health

# Start GLiNER service
cd services/gliner
source venv/bin/activate
uvicorn main:app --host 127.0.0.1 --port 8765
```

**Tests failing**:
1. Update dependencies: `go mod download`
2. Clear test cache: `go clean -testcache`
3. Run specific package: `go test -v ./internal/memorygraph`

### Performance Issues

**Rendering lag**:
- Close other terminal applications
- Check CPU usage: `top` or `htop`
- Reduce terminal size temporarily
- Disable syntax highlighting in large files

**Memory usage high**:
- Check for memory leaks: `go test -memprofile=mem.prof`
- Profile application: `go tool pprof mem.prof`

---

## Code Quality Standards

### Go Code Standards

- **Formatting**: Always run `go fmt ./...` before committing
- **Linting**: Code should pass `go vet ./...` without warnings
- **Naming**: Use descriptive names (no single-letter vars except loop counters)
- **Functions**: Keep functions small (<50 lines preferred)
- **Comments**: Document all exported types and functions
- **Error handling**: Always handle errors, never ignore

**Example**:
```go
// EntityExtractor extracts semantic entities from text using pattern matching.
// It returns a slice of Entity objects with type classification and position.
type EntityExtractor struct {
    patterns map[EntityType]*regexp.Regexp
}

// Extract analyzes the text and returns all detected entities.
func (e *EntityExtractor) Extract(text string) ([]Entity, error) {
    if text == "" {
        return nil, ErrEmptyText
    }

    var entities []Entity
    for typ, pattern := range e.patterns {
        matches := pattern.FindAllStringIndex(text, -1)
        for _, match := range matches {
            entities = append(entities, Entity{
                Type:  typ,
                Text:  text[match[0]:match[1]],
                Start: match[0],
                End:   match[1],
            })
        }
    }

    return entities, nil
}
```

### Test Standards

- **Coverage**: Aim for 70%+ overall, 85%+ for critical paths
- **Table-driven**: Use table-driven tests for comprehensive coverage
- **Clear names**: Test names should describe what is being tested
- **No sleeps**: Avoid `time.Sleep()` in tests (use synchronization)
- **Isolation**: Tests should not depend on external state

**Example**:
```go
func TestEntityExtractor_Extract(t *testing.T) {
    tests := []struct {
        name      string
        text      string
        wantCount int
        wantTypes []EntityType
        wantErr   bool
    }{
        {
            name:      "empty text",
            text:      "",
            wantCount: 0,
            wantTypes: nil,
            wantErr:   true,
        },
        {
            name:      "single person entity",
            text:      "Alice wrote the code",
            wantCount: 1,
            wantTypes: []EntityType{EntityPerson},
            wantErr:   false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := NewEntityExtractor()
            entities, err := e.Extract(tt.text)

            if (err != nil) != tt.wantErr {
                t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if len(entities) != tt.wantCount {
                t.Errorf("Extract() got %d entities, want %d", len(entities), tt.wantCount)
            }

            // Verify entity types
            for i, entity := range entities {
                if i < len(tt.wantTypes) && entity.Type != tt.wantTypes[i] {
                    t.Errorf("Entity %d type = %v, want %v", i, entity.Type, tt.wantTypes[i])
                }
            }
        })
    }
}
```

---

## Resources

### Essential Documentation

**User-Facing**:
- [README.md](README.md) - Project overview and quick start
- [docs/USAGE.md](docs/USAGE.md) - Complete user guide
- [docs/whitepaper.md](docs/whitepaper.md) - Technical whitepaper
- [ROADMAP.md](ROADMAP.md) - Project roadmap and timeline

**Developer**:
- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) - Developer setup and workflows
- [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) - Contribution guidelines
- [docs/TESTING.md](docs/TESTING.md) - Testing framework and coverage
- [docs/STYLE_GUIDE.md](docs/STYLE_GUIDE.md) - Code style guide
- [docs/architecture.md](docs/architecture.md) - Architecture documentation

**Technical References**:
- [spec.md](spec.md) - Technical specification
- [docs/PERFORMANCE.md](docs/PERFORMANCE.md) - Performance benchmarks
- [docs/GLINER_INTEGRATION.md](docs/GLINER_INTEGRATION.md) - GLiNER setup and usage

**Mode Guides**:
- [docs/edit-mode-guide.md](docs/edit-mode-guide.md) - Edit Mode documentation
- [docs/analyze-mode-guide.md](docs/analyze-mode-guide.md) - Analyze Mode documentation
- [docs/orchestrate-mode-guide.md](docs/orchestrate-mode-guide.md) - Orchestrate Mode documentation

**External Resources**:
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss) - Styling library
- [mnemosyne Documentation](https://github.com/rand/mnemosyne) - Memory system
- [Go Documentation](https://go.dev/doc/) - Go language reference

---

## Principles & Philosophy

### Work Plan Protocol

All development follows the **4-phase Work Plan Protocol**:

1. **Prompt → Spec**: Transform request into clear specification
2. **Spec → Full Spec**: Decompose into components with dependencies
3. **Full Spec → Plan**: Create execution plan with parallelization
4. **Plan → Artifacts**: Execute plan, create code/tests/docs

**Never skip phases**. Each phase has defined exit criteria.

### Quality Gates

From Work Plan Protocol, all work must pass:

- [ ] Intent satisfied (does it solve the problem?)
- [ ] Tests written and passing (70%+ coverage)
- [ ] Documentation complete (README, docs/, CHANGELOG)
- [ ] No anti-patterns (TODOs converted to issues)
- [ ] Code formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)

### Design Principles

**Elm Architecture** (Model-Update-View):
- Immutable state updates
- Pure functions for rendering
- Commands for side effects
- Message-based communication

**Event-Driven Architecture**:
- PubSub event broker for component decoupling
- 40+ domain event types
- Non-blocking publish/subscribe

**Mode-Based UI**:
- 5 application modes with distinct purposes
- Each mode has its own layout and lifecycle
- Mode registry for history and switching

**Responsive Design**:
- Adapts to terminal size (80x24 to 200x60)
- 5 layout modes (Focus, Standard, Analysis, Compact, Custom)
- Graceful degradation for small terminals

---

## Appendix: Key Files Quick Reference

| File | Purpose | Lines |
|------|---------|-------|
| `main.go` | Application entry point | ~300 |
| `internal/modes/registry.go` | Mode system | ~400 |
| `internal/editor/editor.go` | Text editor | ~800 |
| `internal/editor/semantic/analyzer.go` | Semantic analysis | ~600 |
| `internal/memorygraph/layout.go` | Graph layout algorithm | ~700 |
| `internal/mnemosyne/client.go` | gRPC client | ~1200 |
| `internal/layout/layout.go` | Layout engine | ~500 |
| `internal/app/events/broker.go` | Event system | ~300 |
| `config.toml` | Application configuration | ~100 |
| `zensical.toml` | Documentation site config | ~60 |

---

**End of Agent Guide**

For questions or clarifications, consult:
- [docs/architecture.md](docs/architecture.md) for architecture details
- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for development workflows
- [.claude/CLAUDE.md](.claude/CLAUDE.md) for Claude Code guidelines
