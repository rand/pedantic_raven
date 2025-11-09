# Pedantic Raven Roadmap

**Current Phase**: 4.3 (Memory UI Components - In Progress)
**Last Updated**: 2025-11-08
**Project Start**: 2025-01-11

---

## Project Vision

Pedantic Raven is an interactive terminal-based context engineering environment that combines semantic analysis, memory integration, and multi-agent orchestration capabilities. It serves as the next-generation interface to the mnemosyne semantic memory system.

**Core Goals**:
- Rich TUI for creating and editing AI context documents
- Real-time semantic analysis (entities, relationships, typed holes)
- Deep integration with mnemosyne memory system
- Multi-mode interface supporting different workflows
- Production-ready terminal application

---

## Phase Status Overview

| Phase | Focus | Status | Tests | Duration |
|-------|-------|--------|-------|----------|
| **Phase 1** | Foundation | ✅ Complete | 87 | 2 weeks |
| **Phase 2** | Semantic Analysis & Edit Mode | ✅ Complete | 291 | 2 weeks |
| **Phase 3** | Advanced Editor Features | ✅ Complete | 424 | 8 days |
| **Phase 4.1** | mnemosyne RPC Client | ✅ Complete | 461 | 3 days |
| **Phase 4.2** | Memory List Component | 🔄 In Progress | - | 3 days |
| **Phase 4.3** | Memory Detail View | 📋 Planned | - | 3 days |
| **Phase 4.4** | Graph Visualization | 📋 Planned | - | 5 days |
| **Phase 4.5** | Explore Mode Integration | 📋 Planned | - | 2 days |
| **Phase 5** | Analyze Mode | 📋 Planned | - | 2-3 weeks |
| **Phase 6** | Orchestrate Mode | 📋 Planned | - | 4-5 weeks |
| **Phase 7** | Collaborate Mode | 📋 Planned | - | 3-4 weeks |
| **Phase 8** | Polish & Production | 📋 Planned | - | 2-3 weeks |

---

## Completed Phases

### Phase 1: Foundation ✅

**Timeline**: 2 weeks (2025-01-11 to 2025-01-25)
**Tests**: 87 passing
**Status**: Complete

**Deliverables**:
- ✅ PubSub Event System with 40+ domain event types
- ✅ Multi-Pane Layout Engine with 5 layout modes
- ✅ Mode Registry supporting 5 application modes
- ✅ Overlay System for modal and non-modal dialogs
- ✅ Command Palette with fuzzy search

**Key Components**:
- `internal/app/events/` - Thread-safe event broker
- `internal/layout/` - Hierarchical pane composition
- `internal/modes/` - Mode management with lifecycle hooks
- `internal/overlay/` - Stack-based overlay management
- `internal/palette/` - Command discovery interface

**Architecture Patterns**:
- Elm Architecture (Model-Update-View)
- Composite Pattern (pane hierarchy)
- Observer Pattern (PubSub events)
- Registry Pattern (modes, commands)
- Strategy Pattern (position strategies, layouts)

---

### Phase 2: Semantic Analysis & Edit Mode ✅

**Timeline**: 2 weeks (Weeks 2-4)
**Tests**: 291 passing (+204 new)
**Status**: Complete

**Deliverables**:
- ✅ Semantic Analyzer with real-time NLP-style analysis
  - Entity extraction (6 types: Person, Place, Thing, Concept, Organization, Technology)
  - Relationship detection (subject-predicate-object patterns)
  - Typed holes (??Type and !!constraint markers)
  - Dependency tracking (imports, requires, references)
  - Triple generation (RDF-style structures)
- ✅ Context Panel with 5 sections
  - Entities with occurrence counting
  - Relationships with confidence scoring
  - Typed holes with priority/complexity markers
  - Dependencies with type indicators
  - Triples (semantic structures)
- ✅ Integrated Terminal
  - Built-in commands (`:help`, `:clear`, `:history`, `:exit`)
  - Shell command execution
  - mnemosyne CLI integration
  - Command history (100 entries)
  - Scrollable output buffer (1000 lines)
- ✅ Edit Mode integration
  - Auto-triggered semantic analysis (500ms debounce)
  - Multi-component layout
  - Focus management

**Key Components**:
- `internal/editor/semantic/` - Analysis engine
- `internal/context/` - Context panel display
- `internal/terminal/` - Terminal component
- `internal/editor/` - Edit mode implementation

---

### Phase 3: Advanced Editor Features ✅

**Timeline**: 8 days (Weeks 5-6)
**Tests**: 424 passing (+133 new)
**Status**: Complete

**Phase Breakdown**:

**3.1: Buffer Manager Integration** (Days 1-2) ✅
- Full undo/redo support via Buffer interface
- Keybindings: Ctrl+Z (undo), Ctrl+Y (redo)
- Cursor position tracking
- Multi-line editing support

**3.2: File Operations** (Days 3-4) ✅
- OpenFile, SaveFile, SaveFileAs methods
- Atomic file saves (temp + rename pattern)
- Dirty flag tracking
- Error handling and path management
- UTF-8 encoding support

**3.3: Search and Replace** (Days 5-6) ✅
- Search engine with literal and regex support
- Case sensitive/insensitive toggle
- Whole word matching
- Replace current match and replace all
- Multi-line search support
- Undo integration for replacements

**3.4: Syntax Highlighting** (Days 7-8) ✅
- Token-based highlighting system
- Go language full tokenization (23 keywords, 19 types, operators, etc.)
- Markdown formatting support
- Automatic language detection (by extension and content)
- Extensible tokenizer architecture
- 12 token types with default color scheme

**Key Components**:
- `internal/editor/buffer/` - Buffer management (52 tests)
- `internal/editor/search/` - Search engine (35 tests)
- `internal/editor/syntax/` - Syntax highlighting (31 tests)

---

### Phase 4.1: mnemosyne RPC Client ✅

**Timeline**: 3 days
**Tests**: 461 total (66 mnemosyne tests)
**Status**: Complete

**Deliverables**:
- ✅ gRPC Client Library for mnemosyne memory system
  - Connection management with configurable timeouts
  - Health checks and server statistics
- ✅ CRUD Operations
  - StoreMemory, GetMemory, UpdateMemory, DeleteMemory, ListMemories
  - Namespace support (Global, Project, Session)
  - Importance scoring and tagging
  - Optional LLM enrichment
- ✅ Search Operations
  - Recall: Hybrid search (semantic + FTS + graph)
  - SemanticSearch: Pure embedding-based search (768d/1536d)
  - GraphTraverse: Multi-hop graph traversal
  - GetContext: Retrieve memories with linked context
- ✅ Streaming Support
  - RecallStream, ListMemoriesStream, StoreMemoryStream
  - Progress updates for long-running operations
- ✅ Error Handling
  - Domain-specific errors (NotFound, InvalidArgument, etc.)
  - gRPC status mapping
  - Helper functions for error checking

**Key Components**:
- `internal/mnemosyne/client.go` - Client infrastructure
- `internal/mnemosyne/memory.go` - CRUD + search operations
- `internal/mnemosyne/errors.go` - Error handling
- `proto/mnemosyne/v1/` - Protobuf schemas

**Code Statistics**:
- 2,202 lines added (code + tests)
- ~182KB generated protobuf code
- 66 comprehensive tests
- 23 exported functions/methods

---

## Current Phase: Phase 4.2-4.5 (Explore Mode)

### Phase 4.2: Memory List Component 🔄

**Timeline**: 3 days (Days 4-6)
**Status**: In Progress

**Objectives**:
- TUI component for browsing memory list
- Filtering and sorting capabilities
- Search integration with mnemosyne client

**Planned Features**:
- Scrollable memory list with rich metadata display
  - Title, namespace, tags, importance, timestamp, link count
- Sorting: by importance, recency, relevance
- Filtering: by namespace, tags, importance range
- Navigation: j/k keys, Enter to select
- Lazy loading (50-100 memories at a time)
- Integration with mnemosyne client (ListMemories, Recall)

**Success Criteria**:
- Smooth scrolling with 1000+ memories
- Filtering works correctly
- Search returns relevant results
- Responsive to terminal size changes

---

### Phase 4.3: Memory Detail View 📋

**Timeline**: 3 days (Days 7-9)
**Status**: Planned

**Objectives**:
- Rich memory visualization with full metadata
- Linked memories display
- Editing capabilities

**Planned Features**:
- Full metadata display (title, namespace, importance, tags, timestamps)
- Importance visualization (bar chart)
- Tag list with colors
- Content display with syntax highlighting
- Link visualization (inbound and outbound)
- Edit mode for in-place updates
- Actions: edit, delete, export (JSON/Markdown)

**Success Criteria**:
- All memory fields displayed correctly
- Edit mode works without data loss
- Links are clickable and navigable
- Actions functional with proper error handling

---

### Phase 4.4: Graph Visualization 📋

**Timeline**: 5 days (Days 10-14)
**Status**: Planned

**Objectives**:
- Interactive memory graph visualization
- Graph traversal UI
- Visual link exploration

**Planned Features**:
- Layout Algorithms:
  - Force-directed (default): Natural clustering
  - Hierarchical: Top-down tree structure
  - Radial: Center node with concentric circles
- Interaction:
  - Pan with hjkl or arrow keys
  - Zoom with +/- keys
  - Select node with Enter
  - Focus on node with 'f'
- Filtering:
  - By namespace, importance, depth
- ASCII rendering for terminal compatibility

**Success Criteria**:
- Renders up to 500 nodes smoothly
- Force-directed layout converges <2s
- Pan and zoom are smooth
- Performance: <16ms render time

---

### Phase 4.5: Explore Mode Integration 📋

**Timeline**: 2 days (Days 15-16)
**Status**: Planned

**Objectives**:
- Integrate all components into cohesive Explore Mode
- Complete memory workspace

**Planned Features**:
- Layout Modes:
  - Standard: Memory list + detail view
  - Graph: Full-screen graph visualization
- Mode switching with 'g' key
- Focus management (Tab to cycle)
- Keybindings:
  - `/` - Search
  - `n` - New memory
  - `r` - Refresh from server
  - `?` - Help overlay

**Success Criteria**:
- All three components integrated smoothly
- Layout modes switchable
- Focus management works correctly
- All keybindings functional
- 15+ integration tests passing

---

## Planned Phases

### Phase 5: Analyze Mode 📋

**Timeline**: 2-3 weeks
**Status**: Planned

**Objectives**:
- Statistical analysis of semantic data
- Entity relationship visualization
- Typed hole prioritization
- Dependency tree visualization

**Planned Features**:
- Triple graph visualization (entities, relationships)
- Entity frequency charts (bar chart, word cloud)
- Relationship mining (pattern detection)
- Interactive filtering (by type, frequency)
- Export analysis reports (PDF, HTML, Markdown)

---

### Phase 6: Orchestrate Mode 📋

**Timeline**: 4-5 weeks
**Status**: Planned

**Objectives**:
- Multi-agent coordination interface
- Task management
- Progress monitoring
- Agent workflows

**Planned Features**:
- Agent DAG visualization (4 agents + status)
- Task queue UI (pending, in-progress, completed)
- Execution log (streaming updates)
- Agent control (pause, resume, cancel, retry)
- Resource monitoring (CPU, memory, API quota)

---

### Phase 7: Collaborate Mode 📋

**Timeline**: 3-4 weeks
**Status**: Planned

**Objectives**:
- Live multi-user editing
- Presence awareness
- Conflict resolution
- Shared annotations

**Planned Features**:
- Live cursors (colored, labeled)
- Presence sidebar (user list, activity)
- Activity feed (real-time notifications)
- Chat sidebar (quick messages)
- Conflict resolution UI (merge conflicts)

---

### Phase 8: Polish & Production 📋

**Timeline**: 2-3 weeks
**Status**: Planned

**Objectives**:
- Performance optimization
- Comprehensive documentation
- Packaging and distribution
- Release 1.0

**Planned Features**:
- Performance optimization (<16ms render, <50ms events)
- Comprehensive user documentation
- API documentation
- Distribution packages (brew, apt, etc.)
- Release pipeline
- Production monitoring

---

## Overall Timeline

**Total Estimated Duration**: 6-8 months
**Start Date**: 2025-01-11
**Current Progress**: ~25% complete (Phases 1-3 + 4.1 done)
**Estimated Completion**: Q2 2025

**Milestone Breakdown**:
- ✅ Week 6: Foundation complete (architecture proven)
- ✅ Week 10: Edit mode working (ICS parity)
- 🔄 Week 17: Explore mode complete (memory workspace)
- 📋 Week 22: Analyze mode complete
- 📋 Week 30: Orchestrate mode complete
- 📋 Week 34: Collaborate mode complete
- 📋 Week 38: v1.0 production release

---

## Success Metrics

**Functionality** (by v1.0):
- All 5 modes implemented and tested
- mnemosyne Level 3 integration complete
- Multi-agent orchestration working
- Graph visualization with 500+ nodes
- Live collaboration support

**Performance**:
- Render time <16ms (60 FPS)
- Event processing <50ms
- Graph layout <200ms for <1000 nodes
- Memory footprint <100MB

**Quality**:
- 700+ tests passing
- 80%+ code coverage
- Zero critical bugs
- Comprehensive documentation

**User Experience**:
- Keyboard-only navigation
- Command palette for discoverability
- Progressive disclosure
- Graceful error handling
- Responsive design (≥80x24 terminals)

---

## Dependencies

**External**:
- Go 1.25+
- Bubble Tea v1.2.6+ (TUI framework)
- Lipgloss v1.0.0+ (styling)
- gRPC + Protocol Buffers (mnemosyne integration)

**Internal**:
- mnemosyne-rpc server (memory system backend)
- Protobuf schemas from mnemosyne project

---

## Risks & Mitigations

| Risk | Impact | Mitigation | Status |
|------|--------|------------|--------|
| Complexity overwhelms UX | High | Progressive disclosure, user testing | Monitoring |
| Performance degrades | Medium | Profile early, optimize hot paths | Active |
| Event system bottleneck | Medium | Batch events, throttle updates | Monitoring |
| Graph visualizations unclear | Medium | Multiple layouts, zoom/pan/filter | Planned |
| Multi-user conflicts (Phase 7) | High | Conflict resolution UI, presence | Deferred |

---

## Next Steps

**Immediate** (Current Sprint):
1. Complete Phase 4.2: Memory List Component
2. Implement Phase 4.3: Memory Detail View
3. Begin Phase 4.4: Graph Visualization

**Short Term** (Next 2 months):
1. Complete Explore Mode (Phase 4.5)
2. Begin Analyze Mode (Phase 5)
3. User testing and feedback

**Long Term** (Q2 2025):
1. Complete all 5 modes
2. Performance optimization
3. Production release preparation
4. v1.0 launch

---

**Document Version**: 1.0
**Last Updated**: 2025-11-08
**Maintained By**: Development Team
