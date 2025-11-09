# Pedantic Raven Roadmap

**Current Phase**: 6 (Analyze Mode - Planned)
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
| **Phase 4.2** | Memory List Component | ✅ Complete | 561 | 3 days |
| **Phase 4.3** | Memory Detail View | ✅ Complete | 626 | 3 days |
| **Phase 4.4** | Graph Visualization | ✅ Complete | 706 | 5 days |
| **Phase 4.5** | Explore Mode Integration | ✅ Complete | 754 | 2 days |
| **Phase 5** | Real mnemosyne Integration | ✅ Complete | 934 | 10 days |
| **Phase 6** | Analyze Mode | 📋 Planned | - | 2-3 weeks |
| **Phase 7** | Orchestrate Mode | 📋 Planned | - | 4-5 weeks |
| **Phase 8** | Collaborate Mode | 📋 Planned | - | 3-4 weeks |
| **Phase 9** | Polish & Production | 📋 Planned | - | 2-3 weeks |

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

**Enhancement: GLiNER ML Integration** (Added 2025-11-09):
- ✅ GLiNER2 machine learning for entity extraction (85-95% accuracy vs 60-70% pattern matching)
- ✅ Zero-shot NER with custom entity types
- ✅ Pluggable extractor interface (PatternExtractor, GLiNERExtractor, HybridExtractor)
- ✅ Python GLiNER service (FastAPI + GLiNER2 model, 340M parameters)
- ✅ Go client library with exponential backoff retry
- ✅ Automatic fallback to pattern matcher when service unavailable
- ✅ Configuration system (TOML + environment variables)
- ✅ 60 new tests (994 total tests)
- ✅ Docker Compose deployment
- ✅ Optional feature (disabled by default, zero breaking changes)

Components:
- `internal/gliner/` - GLiNER HTTP client
- `internal/editor/semantic/extractor.go` - Extractor interface
- `internal/editor/semantic/gliner_extractor.go` - ML implementation
- `internal/editor/semantic/pattern_extractor.go` - Pattern matching implementation
- `internal/config/` - Configuration loading
- `services/gliner/` - Python GLiNER service

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

## Completed Phase: Phase 4.2-4.5 (Explore Mode) ✅

### Phase 4.2: Memory List Component ✅

**Timeline**: 3 days (Days 4-6)
**Tests**: 561 total (13 memorylist tests)
**Status**: Complete

**Objectives**:
- TUI component for browsing memory list
- Filtering and sorting capabilities
- Search integration with mnemosyne client

**Deliverables**:
- ✅ Scrollable memory list with rich metadata display
  - Title, namespace, tags, importance, timestamp, link count
- ✅ Sorting: by importance, recency, relevance
- ✅ Filtering: by namespace, tags, importance range
- ✅ Navigation: j/k keys, g/G jumps, Enter to select
- ✅ Sample data integration (50 test memories)
- ✅ Help overlay with keyboard shortcuts
- ✅ Loading, error, and empty states

**Key Components**:
- `internal/memorylist/model.go` - List model and state
- `internal/memorylist/commands.go` - Search and load operations
- `internal/memorylist/view.go` - Rendering logic
- `internal/memorylist/model_test.go` - 13 comprehensive tests

---

### Phase 4.3: Memory Detail View ✅

**Timeline**: 3 days (Days 7-9)
**Tests**: 626 total (19 memorydetail tests)
**Status**: Complete

**Objectives**:
- Rich memory visualization with full metadata
- Linked memories display
- Navigation capabilities

**Deliverables**:
- ✅ Full metadata display (title, namespace, importance, tags, timestamps)
- ✅ Importance visualization (bar chart with Unicode blocks)
- ✅ Tag list with color-coded display
- ✅ Content display with markdown-style formatting
- ✅ Link visualization (inbound and outbound with counts)
- ✅ Link navigation with Enter key
- ✅ Scrollable content view
- ✅ Help overlay with keyboard shortcuts

**Key Components**:
- `internal/memorydetail/model.go` - Detail model and state
- `internal/memorydetail/view.go` - Rich rendering logic
- `internal/memorydetail/model_test.go` - 19 comprehensive tests

---

### Phase 4.4: Graph Visualization ✅

**Timeline**: 5 days (Days 10-14)
**Tests**: 706 total (146 memorygraph tests)
**Status**: Complete

**Objectives**:
- Interactive memory graph visualization
- Graph traversal UI
- Visual link exploration

**Deliverables**:
- ✅ Force-directed graph layout with physics simulation
  - Spring forces (attractive between connected nodes)
  - Repulsion forces (prevents overlap)
  - Damping and adaptive convergence
- ✅ Canvas-based rendering system
  - Efficient diff-based updates
  - Terminal-aware coordinate mapping
  - Box-drawing Unicode characters
- ✅ Interactive controls:
  - Pan with h/j/k/l keys
  - Zoom with +/- keys
  - Node selection with Tab
  - Node expansion/collapse with e/x
  - Center view with 'c', re-layout with 'r'
- ✅ Performance optimizations:
  - Spatial grid for collision detection
  - Lazy layout updates (space bar stepping)
  - Viewport culling

**Key Components**:
- `internal/memorygraph/model.go` - Graph model and state
- `internal/memorygraph/layout.go` - Force-directed algorithm (80 tests)
- `internal/memorygraph/canvas.go` - Terminal rendering (23 tests)
- `internal/memorygraph/physics.go` - Physics simulation (29 tests)

---

### Phase 4.5: Explore Mode Integration ✅

**Timeline**: 2 days (Days 15-16)
**Tests**: 754 total (11 explore mode tests)
**Status**: Complete

**Objectives**:
- Integrate all components into cohesive Explore Mode
- Complete memory workspace with dual layouts

**Deliverables**:
- ✅ Dual-layout system:
  - Standard: Memory list + detail view (side-by-side)
  - Graph: Full-screen graph visualization
- ✅ Layout switching with 'g' key
- ✅ Focus management (Tab to cycle between list/detail)
- ✅ Context-aware help overlay:
  - Standard layout help (list/detail navigation)
  - Graph layout help (pan/zoom/selection)
- ✅ Professional visual presentation:
  - lipgloss borders (blue for list, green for detail)
  - Clean component integration
- ✅ Message-based coordination:
  - Memory selection flows from list → detail
  - Graph selection shows detail
  - Window resize handled gracefully

**Key Components**:
- `internal/modes/explore.go` - Explore Mode implementation (650 lines)
- `internal/modes/explore_test.go` - 11 integration tests
- Complete integration of memorylist, memorydetail, and memorygraph

---

### Phase 5: Real mnemosyne Integration ✅

**Timeline**: 10 days (Days 1-10)
**Tests**: 934 total (194 new)
**Status**: Complete

**Objectives**:
- Connect Explore Mode to live mnemosyne-rpc server
- Replace sample data with real memory queries
- Enable CRUD operations on real memories
- Implement bidirectional link management

**Deliverables**:
- ✅ **Connection Management** (Component 5.1):
  - Health check monitoring (30s intervals)
  - Auto-reconnect with exponential backoff (1s → 30s)
  - 5 connection states (Disconnected, Connecting, Connected, Reconnecting, Failed)
  - Offline detection and mode switching
  - Thread-safe status tracking (RWMutex)
- ✅ **Real Data Integration** (Component 5.2):
  - Server-side memory loading via gRPC
  - 5-minute TTL query cache
  - Client-side filtering (namespace, tags, importance)
  - Environment-based configuration (MNEMOSYNE_ENABLED, MNEMOSYNE_ADDR)
- ✅ **CRUD Operations** (Component 5.3):
  - Create, Edit, Update, Delete with validation
  - SHA256 hash-based change detection
  - Deep cloning to prevent mutations
  - Comprehensive validation (content required, importance 1-10, max 20 tags)
- ✅ **Link Management** (Component 5.4):
  - 8 link types (REFERENCES, EXTENDS, BUILDS_UPON, CONTRADICTS, etc.)
  - Bidirectional link navigation
  - 50-entry LRU navigation history
  - Link metadata and strength tracking
- ✅ **Search Integration** (Component 5.5):
  - 4 search modes (Hybrid, Semantic, Full-Text, Graph)
  - 500ms debounced search (90% server load reduction)
  - Rich filter options (namespaces, tags, importance)
  - 10-entry LRU search history
- ✅ **Error Handling & Offline Mode** (Component 5.6):
  - Multi-layered error categorization
  - FIFO sync queue for offline operations
  - Exponential backoff retry with context cancellation
  - User-friendly error messages

**Key Components**:
- `internal/mnemosyne/connection.go` - Connection management (405 lines, 17 tests)
- `internal/memorylist/realdata.go` - Real data loading (201 lines, 23 tests)
- `internal/memorydetail/crud.go` - CRUD operations (483 lines, 28 tests)
- `internal/memorydetail/links.go` - Link management (340 lines, 35 tests)
- `internal/memorylist/search.go` - Search integration (572 lines, 29 tests)
- `internal/mnemosyne/offline.go` - Offline support (290 lines, 61 tests)

**Code Statistics**:
- ~10,000 lines added (production + tests + documentation)
- 194 new tests (754 → 934 tests passing)
- 20+ new files across 6 major components
- Zero regressions in existing functionality

**Documentation**:
- `docs/PHASE5_SPEC.md` - Complete technical specification
- `docs/PHASE5_COMPLETE.md` - Implementation summary (788 lines)
- `docs/SEARCH_INTEGRATION_SUMMARY.md` - Search component details
- `docs/LINK_MANAGEMENT_SUMMARY.md` - Link component details

---

## Planned Phases

### Phase 6: Analyze Mode 📋

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

### Phase 7: Orchestrate Mode 📋

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

### Phase 8: Collaborate Mode 📋

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

### Phase 9: Polish & Production 📋

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
**Current Progress**: ~55% complete (Phases 1-5 complete)
**Estimated Completion**: Q2 2025

**Milestone Breakdown**:
- ✅ Week 6: Foundation complete (architecture proven)
- ✅ Week 10: Edit mode working (ICS parity)
- ✅ Week 17: Explore mode complete (memory workspace)
- ✅ Week 19: Real mnemosyne integration (Phase 5)
- 📋 Week 22: Analyze mode complete (Phase 6)
- 📋 Week 30: Orchestrate mode complete (Phase 7)
- 📋 Week 34: Collaborate mode complete (Phase 8)
- 📋 Week 38: v1.0 production release (Phase 9)

---

## Success Metrics

**Functionality** (by v1.0):
- All 5 primary modes implemented and tested
- mnemosyne real-time integration complete
- Multi-agent orchestration working
- Graph visualization with 500+ nodes
- Live collaboration support

**Performance**:
- Render time <16ms (60 FPS)
- Event processing <50ms
- Graph layout <200ms for <1000 nodes
- Memory footprint <100MB

**Quality**:
- 994+ tests passing (currently at 994)
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
1. Create Phase 5 technical specification
2. Connect to live mnemosyne-rpc server
3. Replace sample data with real queries

**Short Term** (Next 2 weeks):
1. Complete Real mnemosyne Integration (Phase 5)
2. Implement CRUD operations
3. Enable link management

**Medium Term** (Next 2 months):
1. Begin Analyze Mode (Phase 6)
2. User testing and feedback
3. Performance optimization

**Long Term** (Q2 2025):
1. Complete all 5 modes
2. Performance optimization
3. Production release preparation
4. v1.0 launch

---

**Document Version**: 1.0
**Last Updated**: 2025-11-08
**Maintained By**: Development Team
