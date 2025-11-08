# Pedantic Raven - Specification Document

**Version**: 1.0
**Date**: 2025-11-08
**Status**: Phase 1 - Foundation

---

## 1. Vision & Purpose

**Pedantic Raven** is an Interactive Context Engineering Environment - an improved version of ICS (Integrated Context Studio) that provides a rich, dynamic TUI for creating, editing, and refining context documents with deep integration to mnemosyne's memory system.

### Core Philosophy

Combine the best of:
1. **ICS's specialized context engineering** - Typed holes, semantic analysis, memory integration
2. **Crush's rich interactive architecture** - Multi-pane layouts, pubsub events, real-time streaming
3. **New capabilities** - Orchestration view, memory workspace, graph visualization, command palette

### Success Criteria

- [ ] Edit mode provides ICS parity with streaming improvements
- [ ] Level 3 mnemosyne integration (multi-agent orchestration) fully functional
- [ ] All 5 modes (Edit, Explore, Analyze, Orchestrate, Collaborate) implemented
- [ ] Command palette enables discoverability
- [ ] Production-ready: tested, documented, performant, accessible

---

## 2. Requirements

### Functional Requirements

#### FR-1: Rich Interactive TUI
- **FR-1.1**: Multi-pane layouts (Focus, Standard, Analysis modes)
- **FR-1.2**: Keyboard and mouse navigation
- **FR-1.3**: Responsive to terminal resize (≥120x30 normal, <120x30 compact)
- **FR-1.4**: Overlay system for dialogs and popups
- **FR-1.5**: Command palette (Ctrl+K) with fuzzy search

#### FR-2: Edit Mode (ICS Parity + Improvements)
- **FR-2.1**: Rope-based text editor with O(log n) operations
- **FR-2.2**: Syntax highlighting (Markdown, TOML, JSON)
- **FR-2.3**: Streaming semantic analysis (triples, holes, entities)
- **FR-2.4**: Typed holes navigator with quick-fill
- **FR-2.5**: Suggestions sidebar with ranked proposals
- **FR-2.6**: Diagnostics pane with quick-fix actions

#### FR-3: Mnemosyne Integration (Level 1-3)
- **FR-3.1**: gRPC connection to mnemosyne-rpc server
- **FR-3.2**: Memory CRUD operations (Store, Get, Update, Delete, List)
- **FR-3.3**: Search operations (Recall, SemanticSearch, GraphTraverse)
- **FR-3.4**: Streaming APIs (progressive results)
- **FR-3.5**: Multi-agent orchestration (bidirectional events)

#### FR-4: Explore Mode (Memory Workspace)
- **FR-4.1**: Force-directed memory graph visualization
- **FR-4.2**: In-place memory editing (create, update, delete)
- **FR-4.3**: Search and filter UI (fuzzy search, importance, tags)
- **FR-4.4**: Link creation UI (drag-and-drop relationships)
- **FR-4.5**: Export subgraphs (JSON, Markdown, DOT)

#### FR-5: Analyze Mode (Semantic Insights)
- **FR-5.1**: Triple graph visualization (entities, relationships)
- **FR-5.2**: Entity frequency charts (bar chart, word cloud)
- **FR-5.3**: Relationship mining (pattern detection)
- **FR-5.4**: Interactive filtering (by type, frequency)
- **FR-5.5**: Export analysis reports (PDF, HTML, Markdown)

#### FR-6: Orchestrate Mode (Agent Coordination)
- **FR-6.1**: Agent DAG visualization (4 agents + status)
- **FR-6.2**: Task queue UI (pending, in-progress, completed)
- **FR-6.3**: Execution log (streaming updates)
- **FR-6.4**: Agent control (pause, resume, cancel, retry)
- **FR-6.5**: Resource monitoring (CPU, memory, API quota)

#### FR-7: Collaborate Mode (Live Multi-User)
- **FR-7.1**: Live cursors (colored, labeled)
- **FR-7.2**: Presence sidebar (user list, activity)
- **FR-7.3**: Activity feed (real-time notifications)
- **FR-7.4**: Chat sidebar (quick messages)
- **FR-7.5**: Conflict resolution UI (merge conflicts)

### Non-Functional Requirements

#### NFR-1: Performance
- **NFR-1.1**: Render time <16ms (60 FPS) for smooth UI
- **NFR-1.2**: Event processing latency <50ms
- **NFR-1.3**: Graph layout update <200ms for <1000 nodes
- **NFR-1.4**: Memory footprint <100MB for typical usage

#### NFR-2: Usability
- **NFR-2.1**: Keyboard-only navigation supported
- **NFR-2.2**: Mouse interactions optional but supported
- **NFR-2.3**: Command palette for discoverability
- **NFR-2.4**: Progressive disclosure (start simple, reveal complexity)

#### NFR-3: Accessibility
- **NFR-3.1**: Screen reader support (announce mode changes)
- **NFR-3.2**: High contrast mode for visibility
- **NFR-3.3**: Configurable keybinds
- **NFR-3.4**: Clear visual focus indicators

#### NFR-4: Reliability
- **NFR-4.1**: Graceful degradation (mnemosyne server unavailable)
- **NFR-4.2**: Auto-save (prevent data loss)
- **NFR-4.3**: Crash recovery (restore session state)
- **NFR-4.4**: Error messages with recovery suggestions

#### NFR-5: Maintainability
- **NFR-5.1**: Clean architecture (Elm pattern)
- **NFR-5.2**: Comprehensive tests (unit, integration, UI)
- **NFR-5.3**: API documentation for all public interfaces
- **NFR-5.4**: Component isolation (easy to modify/extend)

---

## 3. Technical Architecture

### Tech Stack

- **Language**: Go 1.25+
- **TUI Framework**: Bubble Tea (Elm Architecture)
- **Styling**: Lipgloss
- **Components**: Bubbles
- **RPC**: gRPC + Protocol Buffers
- **Graph**: Custom force-directed layout algorithms
- **State**: Immutable Model-Update-View pattern
- **Events**: Channel-based pubsub broker

### Component Architecture

```
pedantic_raven/
├── main.go                          # Entry point
├── internal/
│   ├── app/                         # Application core
│   │   ├── model.go                 # Root application model
│   │   ├── update.go                # Root update function
│   │   ├── view.go                  # Root view function
│   │   └── events/                  # Event system
│   │       ├── broker.go            # PubSub event broker
│   │       └── types.go             # Event type definitions
│   ├── layout/                      # Layout engine
│   │   ├── engine.go                # Layout computation
│   │   ├── pane.go                  # Pane hierarchy
│   │   └── modes.go                 # Layout modes
│   ├── modes/                       # 5 application modes
│   │   ├── registry.go              # Mode registry
│   │   ├── edit/                    # Edit mode
│   │   ├── explore/                 # Explore mode
│   │   ├── analyze/                 # Analyze mode
│   │   ├── orchestrate/             # Orchestrate mode
│   │   └── collaborate/             # Collaborate mode
│   ├── components/                  # Reusable UI components
│   │   ├── editor/                  # Text editor
│   │   ├── tree/                    # File tree
│   │   ├── palette/                 # Command palette
│   │   ├── overlay/                 # Overlay system
│   │   └── graph/                   # Graph visualization
│   ├── mnemosyne/                   # Mnemosyne RPC client
│   │   ├── client.go                # gRPC client connection
│   │   ├── memory.go                # Memory operations
│   │   ├── orchestration.go         # Agent orchestration
│   │   └── streaming.go             # Streaming APIs
│   ├── semantic/                    # Semantic analysis
│   │   ├── analyzer.go              # Analysis engine
│   │   ├── triples.go               # Triple extraction
│   │   ├── holes.go                 # Typed hole detection
│   │   └── entities.go              # Entity tracking
│   └── utils/                       # Utilities
│       ├── fuzzy.go                 # Fuzzy search
│       ├── rope.go                  # Rope data structure
│       └── layout_algorithms.go     # Graph layout
├── proto/                           # Protobuf definitions
│   └── mnemosyne/                   # (Copied from mnemosyne)
│       └── v1/
│           ├── types.proto
│           ├── memory.proto
│           └── health.proto
└── docs/                            # Documentation
    ├── architecture.md
    ├── user-guide.md
    └── api.md
```

### Key Design Patterns

#### 1. Elm Architecture (Bubble Tea)
```go
type Model struct {
    currentMode    Mode
    layout         *layout.Engine
    eventBroker    *events.Broker
    mnemosyneClient *mnemosyne.Client
    // ... other fields
}

func (m Model) Init() tea.Cmd { /* ... */ }
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { /* ... */ }
func (m Model) View() string { /* ... */ }
```

#### 2. PubSub Event System
```go
type EventBroker struct {
    subscribers map[EventType][]chan Event
}

func (b *EventBroker) Publish(event Event)
func (b *EventBroker) Subscribe(eventType EventType) <-chan Event
```

#### 3. Mode Registry
```go
type Mode interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Mode, tea.Cmd)
    View(area tea.Rect) string
    OnEnter()
    OnExit()
}

type ModeRegistry struct {
    modes map[string]Mode
    current string
}
```

#### 4. Layout Engine
```go
type Pane interface {
    Render(area tea.Rect, focus bool) string
}

type Split struct {
    direction Direction
    ratio     float32
    left      Pane
    right     Pane
}

type LayoutEngine struct {
    mode LayoutMode
    root Pane
}
```

---

## 4. Data Structures

### Core Types

#### Event Types
```go
type EventType int

const (
    SemanticAnalysisStarted EventType = iota
    SemanticAnalysisProgress
    SemanticAnalysisComplete
    MemoryRecalled
    MemoryCreated
    MemoryUpdated
    MemoryDeleted
    AgentStarted
    AgentProgress
    AgentCompleted
    AgentFailed
    ProposalGenerated
    ProposalAccepted
    ProposalRejected
    // ... more event types
)
```

#### Layout Modes
```go
type LayoutMode int

const (
    LayoutFocus    LayoutMode = iota  // Editor only
    LayoutStandard                    // Editor + sidebar
    LayoutAnalysis                    // Split screen
)
```

#### Application Modes
```go
type ApplicationMode int

const (
    ModeEdit       ApplicationMode = iota
    ModeExplore
    ModeAnalyze
    ModeOrchestrate
    ModeCollaborate
)
```

---

## 5. Interfaces & Contracts

### Component Interface
```go
type Component interface {
    Update(msg tea.Msg) (Component, tea.Cmd)
    View(area tea.Rect, focused bool) string
}
```

### Mode Interface
```go
type Mode interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Mode, tea.Cmd)
    View(area tea.Rect) string
    OnEnter()  // Called when mode becomes active
    OnExit()   // Called when mode becomes inactive
}
```

### Mnemosyne Client Interface
```go
type MemoryClient interface {
    // CRUD operations
    StoreMemory(ctx context.Context, req *StoreMemoryRequest) (*Memory, error)
    GetMemory(ctx context.Context, id string) (*Memory, error)
    UpdateMemory(ctx context.Context, id string, updates *MemoryUpdates) (*Memory, error)
    DeleteMemory(ctx context.Context, id string) error
    ListMemories(ctx context.Context, filters *MemoryFilters) ([]*Memory, error)

    // Search operations
    Recall(ctx context.Context, query string, opts *RecallOptions) ([]*SearchResult, error)
    SemanticSearch(ctx context.Context, embedding []float32, opts *SearchOptions) ([]*SearchResult, error)
    GraphTraverse(ctx context.Context, seedIDs []string, opts *TraverseOptions) ([]*Memory, error)

    // Streaming
    RecallStream(ctx context.Context, query string, opts *RecallOptions) (<-chan *SearchResult, error)
    ListMemoriesStream(ctx context.Context, filters *MemoryFilters) (<-chan *Memory, error)

    // Orchestration
    SubscribeToAgentEvents(ctx context.Context) (<-chan *AgentEvent, error)
    CreateTask(ctx context.Context, task *Task) (*TaskResult, error)
}
```

---

## 6. Constraints & Assumptions

### Constraints

1. **Terminal Support**: Requires ANSI-capable terminal (xterm-256color or better)
2. **Go Version**: Minimum Go 1.21 for generics and improved error handling
3. **Mnemosyne Server**: Requires mnemosyne-rpc server running (graceful degradation if unavailable)
4. **Screen Size**: Normal mode requires ≥120x30, compact mode for smaller

### Assumptions

1. **Single User**: Phase 1 focuses on single-user experience (Collaborate mode in Phase 7)
2. **Local Network**: Mnemosyne server assumed on localhost or LAN (low latency)
3. **UTF-8 Terminal**: All text rendering assumes UTF-8 support
4. **Modern Terminal**: Mouse support, 256+ colors, Unicode box-drawing characters

---

## 7. Success Metrics

### Phase 1 (Foundation)
- [ ] Event system functional (publish/subscribe working)
- [ ] Layout engine renders 3 modes correctly
- [ ] Mode registry switches between 5 modes
- [ ] Overlay system stacks dialogs correctly
- [ ] Command palette fuzzy search functional

### Phase 2 (Edit Mode)
- [ ] Editor handles files >10MB without lag
- [ ] Semantic analysis completes <2s for 1000-line file
- [ ] Typed holes detected accurately (>95% precision)
- [ ] Suggestions sidebar shows AI proposals
- [ ] Diagnostics pane shows validation errors

### Phase 3 (Mnemosyne Integration)
- [ ] All 13 MemoryService methods working
- [ ] Streaming APIs deliver progressive results
- [ ] Orchestration client receives agent events
- [ ] Graceful degradation when server unavailable
- [ ] Reconnection logic works after network issues

### Phases 4-7
(Detailed success metrics defined in phase-specific specs)

---

## 8. Dependencies

### External Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/charmbracelet/bubbles` - Pre-built components
- `google.golang.org/grpc` - gRPC client
- `google.golang.org/protobuf` - Protocol Buffers

### Internal Dependencies

- Mnemosyne RPC server (Level 3 orchestration)
- Protobuf schemas from `/Users/rand/src/mnemosyne/proto/`

---

## 9. Risks & Mitigations

### Risk 1: Complexity Overwhelms UX
- **Mitigation**: Progressive disclosure, start with Edit mode (familiar), add modes incrementally
- **Validation**: User testing after each mode implementation

### Risk 2: Performance Degrades with Rich UI
- **Mitigation**: Profile early, optimize rendering hot paths, lazy rendering
- **Validation**: Benchmark tests, target <16ms render time

### Risk 3: Event System Becomes Bottleneck
- **Mitigation**: Batch events, throttle high-frequency updates, prioritize critical events
- **Validation**: Load testing with 1000+ events/second

### Risk 4: Graph Visualizations Hard to Read
- **Mitigation**: Multiple layout algorithms, interactive controls (zoom/pan/filter)
- **Validation**: User research with real projects

### Risk 5: Multi-User Conflicts (Phase 7)
- **Mitigation**: Conflict resolution UI, presence awareness, manual merge for complex cases
- **Validation**: Stress testing with 10+ simultaneous users

---

## 10. Timeline & Milestones

See main plan for detailed phase breakdown (30-44 weeks total).

**Key Milestones**:
- Week 6: Foundation complete (architecture proven)
- Week 10: Edit mode working (ICS parity)
- Week 17: Level 3 mnemosyne integration
- Week 40: v1.0 production release

---

## 11. Open Questions

1. **Q**: Should Edit mode support multiple buffers simultaneously?
   - **Status**: TBD - decide in Phase 2 based on user feedback

2. **Q**: What graph layout algorithm performs best for memory graphs?
   - **Status**: TBD - benchmark force-directed vs hierarchical in Phase 4

3. **Q**: How to handle very large memory graphs (>10,000 nodes)?
   - **Status**: TBD - implement pagination/clustering in Phase 4

4. **Q**: Should Collaborate mode require authentication?
   - **Status**: TBD - security review in Phase 7

---

**Document Status**: ✅ Complete
**Next Phase**: Phase 2 - Full Spec (Component decomposition, test plan)
