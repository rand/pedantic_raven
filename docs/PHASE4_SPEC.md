# Phase 4 Specification: Explore Mode - Memory Workspace

**Status**: Planning
**Timeline**: 3-4 weeks
**Priority**: High
**Dependencies**: Phase 3 (Complete), mnemosyne RPC integration

## Overview

Phase 4 implements the **Explore Mode** - a dedicated interface for browsing, searching, and manipulating the mnemosyne memory graph. This mode transforms Pedantic Raven from a context editor into a full-featured memory workspace, providing visual and interactive access to the semantic memory system.

## Goals

1. **mnemosyne Integration**: Connect to mnemosyne RPC server (Level 1-2 operations)
2. **Memory Browsing**: View and navigate memory notes with rich metadata
3. **Graph Visualization**: Visual representation of memory relationships
4. **In-Place Editing**: Create, update, and delete memories directly
5. **Search & Filter**: Find memories by content, tags, importance, and relationships

## Phase Structure

### Phase 4.1: mnemosyne RPC Client (Days 1-3)

**Objective**: Implement gRPC client for mnemosyne-rpc server

#### Tasks

**Day 1**: Project setup and protobuf integration
- Copy protobuf schemas from mnemosyne-rpc-dev project
- Set up protobuf compilation in build.rs / Makefile
- Create package structure for mnemosyne client
- Verify protobuf compilation

**Day 2**: Implement core client methods
- Connection management (connect, disconnect, health check)
- CRUD operations (Store, Get, Update, Delete, List)
- Error handling and status code mapping
- Client lifecycle management

**Day 3**: Implement search operations
- Recall (semantic search by query)
- ListMemories with filters
- Basic streaming support
- Integration tests with mock server

#### Success Criteria

- ✅ Successfully connects to mnemosyne-rpc server
- ✅ All CRUD operations functional
- ✅ Recall returns relevant memories
- ✅ Graceful handling of connection errors
- ✅ 20+ client tests passing

### Phase 4.2: Memory List Component (Days 4-6)

**Objective**: Create browsable memory list UI

#### UI Design

```
┌─ Memory Workspace ─────────────────────────────┐
│ Search: [                    ] [Filter ▾]      │
├────────────────────────────────────────────────┤
│ > Memory Title Here                   [Imp: 8] │
│   project:myapp, architecture, design          │
│   Updated 2 hours ago • 3 links                │
│                                                 │
│   Another Memory Note                  [Imp: 7]│
│   global:security, patterns                    │
│   Updated yesterday • 1 link                   │
│                                                 │
│   Third Memory Entry                   [Imp: 5]│
│   project:myapp:api                            │
│   Updated 3 days ago • 0 links                 │
├────────────────────────────────────────────────┤
│ 45 memories | Showing 1-10 | j/k: nav, Enter: │
└────────────────────────────────────────────────┘
```

#### Features

- **List View**: Scrollable list of memories
- **Rich Display**: Title, namespace, tags, importance, timestamp, link count
- **Sorting**: By importance, recency, relevance
- **Filtering**: By namespace, tags, importance range
- **Selection**: Navigate with j/k, select with Enter
- **Lazy Loading**: Load memories in batches (50-100 at a time)

#### Tasks

**Day 4 Morning**: Component structure
- Create `MemoryListComponent` type
- Implement basic list rendering
- Add scroll management
- Selection highlighting

**Day 4 Afternoon**: Memory display formatting
- Format memory metadata (importance, timestamp, tags)
- Truncate long content preview
- Color coding by importance
- Namespace hierarchical display

**Day 5 Morning**: Filtering and sorting
- Filter by namespace (prefix matching)
- Filter by tags (multi-select)
- Filter by importance range
- Sort by importance/recency/relevance

**Day 5 Afternoon**: Search integration
- Search input field
- Recall integration for semantic search
- Result highlighting
- Search history

**Day 6**: Polish and testing
- Keyboard shortcuts (j/k, /, Enter, Esc)
- Empty states ("No memories found")
- Loading states with spinners
- Error states with retry
- Component tests (15+)

#### Success Criteria

- ✅ Displays list of memories from mnemosyne
- ✅ Smooth scrolling (no lag with 1000+ memories)
- ✅ Filtering works correctly
- ✅ Search returns relevant results
- ✅ Responsive to terminal size changes

### Phase 4.3: Memory Detail View (Days 7-9)

**Objective**: Display and edit single memory with full metadata

#### UI Design

```
┌─ Memory Detail ────────────────────────────────┐
│ [Edit] [Delete] [Export]                  [x]  │
├────────────────────────────────────────────────┤
│ Title: Architecture Decision Record            │
│ Namespace: project:myapp:architecture          │
│ Importance: ████████░░ 8/10                    │
│ Tags: architecture, patterns, event-sourcing   │
│ Created: 2025-01-15 14:30                      │
│ Updated: 2025-01-16 09:15                      │
│                                                 │
│ Content:                                       │
│ ┌──────────────────────────────────────────┐  │
│ │ We decided to use event sourcing for     │  │
│ │ audit trail requirements...              │  │
│ │                                          │  │
│ │ (full content with syntax highlighting)  │  │
│ └──────────────────────────────────────────┘  │
│                                                 │
│ Links (3):                                     │
│ • system-design-doc (bidirectional)            │
│ • database-schema (forward)                    │
│ • compliance-requirements (backward)           │
└────────────────────────────────────────────────┘
```

#### Features

- **Full Metadata**: All memory fields displayed
- **Editable**: Click [Edit] to modify any field
- **Content Display**: Syntax highlighted markdown
- **Link Visualization**: Show all connected memories
- **Actions**: Edit, delete, export, navigate to links

#### Tasks

**Day 7 Morning**: Display structure
- Create `MemoryDetailComponent` type
- Layout memory metadata
- Content area with scroll
- Link list display

**Day 7 Afternoon**: Metadata display
- Format timestamp (relative and absolute)
- Importance bar chart
- Tag list with colors
- Namespace hierarchical breadcrumb

**Day 8 Morning**: Edit mode
- Toggle between view and edit
- Edit content in integrated editor
- Update metadata (importance, tags, namespace)
- Save changes to mnemosyne

**Day 8 Afternoon**: Link navigation
- Display inbound and outbound links
- Link type indicators
- Click to navigate to linked memory
- Breadcrumb trail for navigation history

**Day 9**: Actions and testing
- Delete confirmation dialog
- Export to markdown/JSON
- Copy memory ID
- Component tests (12+)

#### Success Criteria

- ✅ Displays all memory fields correctly
- ✅ Edit mode works without data loss
- ✅ Links are clickable and navigable
- ✅ Actions (edit, delete, export) functional
- ✅ Proper error handling

### Phase 4.4: Graph Visualization (Days 10-14)

**Objective**: Visual memory graph with interactive navigation

#### UI Design

```
┌─ Memory Graph ─────────────────────────────────┐
│ [Force] [Hierarchical] [Radial]   [Reset View] │
├────────────────────────────────────────────────┤
│                                                 │
│          ●────●                                 │
│         /      \                                │
│        ●        ●───●                           │
│         \      /                                │
│          ●────●                                 │
│            │                                    │
│            ●                                    │
│                                                 │
│ Pan: h/j/k/l | Zoom: +/- | Select: Enter       │
└────────────────────────────────────────────────┘
```

#### Visualization Approaches

**Option 1: ASCII Graph** (Simple, fast)
- Nodes as ● or [Name]
- Edges as lines ─│├┤└┌
- Layout using force-directed algorithm
- Pan and zoom by adjusting viewport

**Option 2: Rich TUI Graph** (More complex)
- Multi-line node boxes with titles
- Curved or angled connections
- Color-coded by importance/namespace
- Interactive hover tooltips

**Recommendation**: Start with Option 1 (ASCII), upgrade to Option 2 in future phase

#### Features

- **Layout Algorithms**:
  - Force-directed (default): Natural clustering
  - Hierarchical: Top-down tree structure
  - Radial: Center node with concentric circles
- **Interaction**:
  - Pan with hjkl or arrow keys
  - Zoom with +/- keys
  - Select node with Enter
  - Focus on node with 'f'
  - Show/hide labels with 'l'
- **Filtering**:
  - By namespace (show only project:myapp)
  - By importance (only >= 7)
  - By depth (N hops from seed node)

#### Tasks

**Day 10-11**: Graph data structure and algorithms
- Graph representation (adjacency list)
- Force-directed layout algorithm (basic physics simulation)
- Hierarchical layout (tree traversal with level assignment)
- Viewport transformation (pan, zoom)

**Day 12-13**: Rendering and interaction
- ASCII node and edge rendering
- Viewport clipping
- Node selection and highlighting
- Keyboard navigation (pan, zoom, select)

**Day 14**: Polish and testing
- Layout algorithm optimization
- Smooth transitions
- Performance testing (100+ nodes)
- Graph component tests (10+)

#### Success Criteria

- ✅ Renders graph of up to 500 nodes
- ✅ Force-directed layout converges <2s
- ✅ Pan and zoom are smooth
- ✅ Node selection works
- ✅ Performance: <16ms render time

### Phase 4.5: Explore Mode Integration (Days 15-16)

**Objective**: Integrate all components into cohesive Explore Mode

#### Layout Structure

**Standard Layout** (default):
```
┌─────────────────────────────────────────────────┐
│ Memory List            │ Memory Detail          │
│                        │                        │
│ > Memory 1      [8]    │ Title: ...             │
│   Memory 2      [7]    │ Namespace: ...         │
│   Memory 3      [6]    │                        │
│                        │ Content:               │
│                        │ ...                    │
│                        │                        │
│                        │ Links: ...             │
├────────────────────────┴────────────────────────┤
│ Status: 45 memories loaded | Help: ?            │
└─────────────────────────────────────────────────┘
```

**Graph Layout** (toggle with 'g'):
```
┌─────────────────────────────────────────────────┐
│ Memory Graph (full screen)                      │
│                                                  │
│          ●────●                                  │
│         /      \                                 │
│        ●        ●───●                            │
│                                                  │
├──────────────────────────────────────────────────┤
│ Pan: hjkl | Zoom: +/- | Toggle: g | Help: ?     │
└──────────────────────────────────────────────────┘
```

#### Features

- **Layout Modes**: Standard (list+detail), Graph (full screen)
- **Mode Switching**: Press 'g' to toggle
- **Focus Management**: Tab cycles between components
- **Keybindings**:
  - `1`/`2`/`3` - Switch application modes (Edit, Explore, Analyze)
  - `g` - Toggle graph view
  - `Tab` - Cycle focus
  - `/` - Search
  - `n` - New memory
  - `r` - Refresh from server
  - `?` - Help overlay

#### Tasks

**Day 15 Morning**: Mode structure
- Create `ExploreMode` type
- Implement `Init()`, `Update()`, `View()`
- Set up layout engine with two modes
- Component initialization

**Day 15 Afternoon**: Component integration
- Memory list on left, detail on right
- Graph full-screen toggle
- Focus management and state passing
- Event handling coordination

**Day 16 Morning**: Actions and commands
- New memory creation flow
- Refresh from server
- Mode switching
- Help overlay with keybindings

**Day 16 Afternoon**: Testing and polish
- Integration tests (mode switching, focus, actions)
- Error handling (server unavailable, no memories)
- Loading states
- Final polish

#### Success Criteria

- ✅ All three components integrated smoothly
- ✅ Layout modes switchable
- ✅ Focus management works correctly
- ✅ All keybindings functional
- ✅ 15+ integration tests passing

## Technical Architecture

### New Packages

```
internal/
├── mnemosyne/                    # NEW: mnemosyne RPC client
│   ├── client.go                 # gRPC client connection
│   ├── memory.go                 # Memory CRUD operations
│   ├── search.go                 # Recall and search
│   ├── streaming.go              # Streaming APIs
│   └── client_test.go            # Client tests
├── explore/                      # NEW: Explore Mode components
│   ├── memory_list.go            # Memory list component
│   ├── memory_detail.go          # Detail view component
│   ├── memory_graph.go           # Graph visualization
│   └── explore_mode.go           # Mode implementation
├── graph/                        # NEW: Graph algorithms
│   ├── layout.go                 # Layout algorithms
│   ├── force_directed.go         # Force-directed layout
│   ├── hierarchical.go           # Tree layout
│   └── render.go                 # ASCII rendering
└── modes/
    └── explore.go                # Explore mode registration
```

### Data Flow

```
ExploreMode
    │
    ├─> MemoryListComponent
    │       │
    │       └─> MnemosyneClient.ListMemories()
    │
    ├─> MemoryDetailComponent
    │       │
    │       ├─> MnemosyneClient.GetMemory()
    │       └─> MnemosyneClient.UpdateMemory()
    │
    └─> MemoryGraphComponent
            │
            ├─> MnemosyneClient.GraphTraverse()
            └─> LayoutAlgorithm.Compute()
```

### Event System Integration

New events for Phase 4:

```go
const (
    // Memory events
    MemoryLoaded EventType = iota + 100
    MemoryCreated
    MemoryUpdated
    MemoryDeleted
    MemoriesListed

    // Search events
    RecallStarted
    RecallProgress
    RecallComplete

    // Connection events
    MnemosyneConnected
    MnemosyneDisconnected
    MnemosyneError
)
```

## Testing Strategy

### Unit Tests (60+ tests)

**mnemosyne Client** (20 tests):
- Connection management
- CRUD operations
- Search operations
- Error handling
- Mock server interactions

**Memory List Component** (15 tests):
- List rendering
- Scrolling
- Filtering
- Sorting
- Search integration

**Memory Detail Component** (12 tests):
- Display formatting
- Edit mode
- Link navigation
- Actions (delete, export)

**Graph Component** (10 tests):
- Layout algorithms
- Rendering
- Pan/zoom
- Node selection

**Explore Mode** (15 tests):
- Mode initialization
- Component coordination
- Layout switching
- Event handling

### Integration Tests (15 tests)

- Connect → List → Select → View detail
- Search → Filter → View results
- Edit memory → Save → Verify update
- Create memory → Save → Verify in list
- Delete memory → Confirm → Verify removed
- Graph view → Select node → View detail
- Mode switch (Edit ↔ Explore)
- Server disconnect → Graceful degradation
- Reconnect → Resume operations

### Performance Tests

- Load 1000 memories: <2s
- Filter 1000 memories: <100ms
- Graph layout 500 nodes: <2s
- Graph render: <16ms (60 FPS)
- Memory search: <500ms
- Recall semantic search: <1s

## Success Metrics

**Functionality**:
- ✅ Can browse all memories from mnemosyne
- ✅ Search returns relevant results
- ✅ Can create, edit, delete memories
- ✅ Graph visualization is readable
- ✅ Navigation between memories works

**Performance**:
- ✅ No lag with 1000+ memories
- ✅ Graph renders smoothly
- ✅ Search completes quickly

**Quality**:
- ✅ 75+ tests passing
- ✅ ~80% code coverage
- ✅ No data loss on errors
- ✅ Graceful server disconnection

**User Experience**:
- ✅ Intuitive navigation
- ✅ Helpful error messages
- ✅ Responsive to all sizes
- ✅ Familiar keybindings

## Timeline Summary

| Days | Focus | Deliverable |
|------|-------|-------------|
| 1-3  | mnemosyne RPC Client | CRUD + Search working |
| 4-6  | Memory List Component | Browsable memory list |
| 7-9  | Memory Detail View | Full memory display + edit |
| 10-14 | Graph Visualization | Interactive graph |
| 15-16 | Explore Mode Integration | Complete Explore Mode |

**Total**: 16 days (~3 weeks)

## Dependencies

**External**:
- mnemosyne-rpc server running
- gRPC client libraries
- Protobuf compiler

**Internal**:
- ✅ Event system (Phase 1)
- ✅ Layout engine (Phase 1)
- ✅ Mode registry (Phase 1)
- ✅ Overlay system (Phase 1)

## Risks & Mitigations

### Risk 1: mnemosyne Server Unavailable
- **Mitigation**: Graceful degradation, offline mode with cached data
- **Validation**: Connection error tests, reconnection logic

### Risk 2: Graph Performance with Large Datasets
- **Mitigation**: Limit initial nodes, lazy loading, viewport culling
- **Validation**: Performance tests with 500+ nodes

### Risk 3: Complex State Synchronization
- **Mitigation**: Event-driven updates, clear data ownership
- **Validation**: Integration tests for all state transitions

### Risk 4: ASCII Graph Readability
- **Mitigation**: Multiple layout algorithms, adjustable zoom
- **Validation**: User testing with real memory graphs

## Future Enhancements (Phase 5+)

Features explicitly deferred:

- Advanced graph layouts (radial, circular)
- Graph export (DOT, SVG)
- Bulk operations (tag multiple, delete multiple)
- Memory templates
- Collaborative features (shared namespaces)
- Full-text search within content
- Relationship type creation UI

## Documentation Updates

After Phase 4 completion:

1. Create `docs/PHASE4_SUMMARY.md`
2. Update `README.md` roadmap
3. Document mnemosyne client API
4. Add Explore Mode user guide
5. Document graph visualization controls

---

## Approval Checklist

Before starting implementation:

- [ ] Specification reviewed
- [ ] mnemosyne-rpc server available for testing
- [ ] Protobuf schemas accessible
- [ ] Timeline realistic
- [ ] Dependencies met
- [ ] Test strategy defined
- [ ] Success criteria clear

**Status**: Ready for review and approval
