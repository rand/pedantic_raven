# Phase 9: Explore Mode Completion

**Status**: Planning
**Phase**: 9 of 9 (Final Phase)
**Timeline**: 5-7 days
**Dependencies**: Phase 8 (Refinement & Polish) complete

## Overview

Phase 9 completes the Explore Mode implementation by integrating real mnemosyne operations, implementing CRUD functionality, adding search/filtering, and creating comprehensive documentation. This is the final phase of Pedantic Raven's initial development cycle.

## Current State Assessment

**Existing Implementation** (✅ Complete):
- ✅ Mode structure and lifecycle (OnEnter, OnExit, Init)
- ✅ Two layout modes (Standard: list+detail, Graph: full-screen)
- ✅ Component architecture (memorylist, memorydetail, memorygraph)
- ✅ UI rendering and layout management
- ✅ Help overlays for both layouts
- ✅ Keyboard navigation and focus management
- ✅ Sample data generation (for demonstration)

**Missing Implementation** (❌ Needs Work):
- ❌ Real mnemosyne API integration (currently using sample data)
- ❌ CRUD operations (Create, Update, Delete memories)
- ❌ Search and filtering (regex, tags, importance, namespaces)
- ❌ Link management (create, delete, navigate between memories)
- ❌ Offline mode support (queue operations, sync when online)
- ❌ Error handling and recovery
- ❌ Complete test coverage (unit + integration tests)
- ❌ User documentation (Explore Mode guide)

## Objectives

1. **mnemosyne Integration**: Replace sample data with real API calls
2. **CRUD Operations**: Full create, read, update, delete functionality
3. **Search & Filtering**: Comprehensive search with multiple criteria
4. **Link Management**: Create and navigate memory relationships
5. **Offline Support**: Queue operations and sync gracefully
6. **Testing**: 80%+ coverage with integration tests
7. **Documentation**: Complete user guide (600+ lines)

## Work Streams

### Stream 1: mnemosyne Integration (2 days)

**Goal**: Replace sample data with real mnemosyne API calls

**Current**: Sample data in `loadSampleMemories()` and `loadSampleGraph()`

**Tasks**:
1. **List Operations** (internal/memorylist/realdata.go)
   - Implement `LoadMemories(client, namespace, limit, offset)` using mnemosyne.List RPC
   - Handle pagination (cursor-based or offset-based)
   - Support namespace filtering (project, global, component)
   - Error handling with offline fallback

2. **Detail Operations** (internal/memorydetail/crud.go)
   - Implement `LoadMemory(client, memoryID)` using mnemosyne.Get RPC
   - Support offline cache (read from local cache if offline)
   - Error handling with user-friendly messages

3. **Graph Operations** (internal/memorygraph/loader.go - needs creation)
   - Implement `LoadGraph(client, rootMemoryID, depth)`
   - Traverse links to build graph structure
   - Use mnemosyne.TraverseLinks RPC
   - Support configurable depth (1-5 levels)
   - Handle cycles (prevent infinite loops)

4. **Offline Support**
   - Use existing `internal/mnemosyne/connection_manager.go` offline mode
   - Queue operations when offline
   - Display offline indicator in UI
   - Sync queue when connection restored

**Deliverables**:
- Real data loading from mnemosyne
- Pagination support
- Offline mode support
- Error handling

---

### Stream 2: CRUD Operations (2 days)

**Goal**: Implement full create, update, delete operations

**Current**: Stub implementations in memorydetail/crud.go

**Tasks**:
1. **Create Memory** (memorydetail/crud.go)
   - Form UI for new memory (content, importance, tags, namespace)
   - Validation (content required, importance 0-10, tags format)
   - Call mnemosyne.Store RPC
   - Add to list view after creation
   - Offline queue if not connected

2. **Update Memory** (memorydetail/crud.go)
   - Edit form for existing memory
   - Field validation
   - Call mnemosyne.Update RPC
   - Refresh detail view after update
   - Offline queue if not connected

3. **Delete Memory** (memorydetail/crud.go)
   - Confirmation dialog (prevent accidental deletion)
   - Call mnemosyne.Delete RPC
   - Remove from list view
   - Refresh graph (if memory had links)
   - Offline queue if not connected

4. **Link Management** (memorydetail/links.go)
   - Create link UI (select target memory, set strength, type)
   - Call mnemosyne.CreateLink RPC
   - Delete link with confirmation
   - Call mnemosyne.DeleteLink RPC
   - Update graph visualization

5. **Form Validation**
   - Content: 1-10,000 characters
   - Importance: 0-10 (integer)
   - Tags: alphanumeric + hyphens, comma-separated
   - Namespace: project:name or global format
   - Link strength: 0.0-1.0 (float)

**Deliverables**:
- Create memory form with validation
- Update memory form
- Delete with confirmation
- Link create/delete
- Offline operation queuing

---

### Stream 3: Search & Filtering (1 day)

**Goal**: Implement comprehensive search and filtering

**Current**: Stub in memorylist/search.go

**Tasks**:
1. **Text Search** (memorylist/search.go)
   - Full-text search across memory content
   - Support regex patterns
   - Case-sensitive and case-insensitive modes
   - Use mnemosyne.Search RPC (if available) or client-side filter

2. **Tag Filtering** (memorylist/search.go)
   - Filter by single tag or multiple tags (AND/OR logic)
   - Tag autocomplete from existing tags
   - Show tag counts

3. **Importance Filtering** (memorylist/search.go)
   - Range filter (e.g., importance >= 7)
   - Quick filters (High: 8-10, Medium: 5-7, Low: 0-4)

4. **Namespace Filtering** (memorylist/search.go)
   - Filter by project namespace
   - Filter by global namespace
   - Show namespace counts

5. **Combined Filters** (memorylist/search.go)
   - AND logic for all filters
   - Clear all filters
   - Save filter presets (optional)

6. **Search UI**
   - Search input at top of list
   - Filter pills showing active filters
   - Clear buttons for each filter
   - Results count display

**Deliverables**:
- Text search with regex
- Tag filtering
- Importance filtering
- Namespace filtering
- Combined filter logic
- Search UI

---

### Stream 4: Link Navigation (1 day)

**Goal**: Implement memory link navigation and visualization

**Current**: TODO in explore.go:170

**Tasks**:
1. **Link Navigation** (modes/explore.go)
   - Handle `LinkSelectedMsg` from memorydetail
   - Load linked memory using mnemosyne.Get
   - Update detail view with new memory
   - Update list selection (if memory is in list)
   - Add to navigation history (back/forward)

2. **Navigation History** (modes/explore.go)
   - Track navigation stack (list of memory IDs)
   - Back button (go to previous memory)
   - Forward button (go to next memory in history)
   - Keyboard shortcuts (Alt+Left, Alt+Right)

3. **Graph Navigation** (memorygraph/navigation.go - needs creation)
   - Click on node to view memory
   - Expand node to load and show linked memories
   - Collapse node to hide children
   - Center view on selected node
   - Highlight path from root to selected node

4. **Breadcrumb Trail** (modes/explore.go)
   - Show navigation path at top (e.g., "Root > Concept A > Detail A1")
   - Click breadcrumb to jump to that level
   - Max 5 levels displayed (with ellipsis for deeper)

**Deliverables**:
- Link navigation with history
- Graph node interaction
- Breadcrumb trail
- Keyboard shortcuts

---

### Stream 5: Testing & Quality (2 days)

**Goal**: Comprehensive test coverage for Explore Mode

**Current**: Some tests in memorylist, memorydetail, memorygraph

**Tasks**:
1. **Unit Tests** (target: 80%+ coverage)
   - modes/explore_test.go: Mode lifecycle, layout switching, focus management
   - memorylist/model_test.go: Expand tests for real data loading
   - memorydetail/crud_test.go: CRUD operation tests
   - memorydetail/links_test.go: Link management tests
   - memorygraph/graph_test.go: Graph operations tests

2. **Integration Tests** (internal/integration/)
   - Test Explore mode with real mnemosyne client (use mock server)
   - Test CRUD workflow: Create → Update → Delete
   - Test search and filtering
   - Test link navigation
   - Test offline mode (disconnect, queue operations, reconnect)
   - Test graph visualization and interaction

3. **Error Handling Tests**
   - Test network failures
   - Test invalid input validation
   - Test concurrent operations
   - Test memory not found
   - Test permission errors

4. **Performance Tests** (benchmarks)
   - Benchmark list loading (1000+ memories)
   - Benchmark graph layout (100+ nodes)
   - Benchmark search operations
   - Memory usage profiling

**Deliverables**:
- 80%+ unit test coverage
- 10+ integration tests
- 5+ benchmark tests
- Error handling coverage

---

### Stream 6: Documentation (1 day)

**Goal**: Complete user documentation for Explore Mode

**Tasks**:
1. **User Guide** (docs/explore-mode-guide.md - 600+ lines)
   - Overview and features
   - Getting started (entering Explore mode)
   - Interface overview (Standard layout, Graph layout)
   - Keyboard shortcuts reference (complete table)
   - Memory operations (Create, Edit, Delete)
   - Search and filtering (Text, Tags, Importance, Namespace)
   - Link management (Create, Delete, Navigate)
   - Graph visualization (Pan, Zoom, Expand, Collapse)
   - Offline mode (How it works, Sync behavior)
   - Example workflows (5+ real-world scenarios)
   - Troubleshooting (Common issues and solutions)
   - FAQ (15+ questions)

2. **API Documentation** (if needed)
   - Document public APIs for extending Explore mode
   - Component interfaces
   - Custom renderers

**Deliverables**:
- docs/explore-mode-guide.md (600+ lines)
- Updated README.md

---

## Parallelization Strategy

### Phase 1 (Days 1-2): Foundation

**Parallel Streams** (2 agents):
- **Agent 1 (Sonnet)**: mnemosyne Integration (Stream 1)
- **Agent 2 (Haiku)**: CRUD Operations (Stream 2)

**Dependencies**: None - these can run in parallel

---

### Phase 2 (Days 3-4): Features

**Parallel Streams** (2 agents):
- **Agent 3 (Haiku)**: Search & Filtering (Stream 3)
- **Agent 4 (Haiku)**: Link Navigation (Stream 4)

**Dependencies**: Require Phase 1 completion (need real data loading and CRUD)

---

### Phase 3 (Days 5-7): Quality

**Parallel Streams** (2 agents):
- **Agent 5 (Sonnet)**: Testing & Quality (Stream 5)
- **Agent 6 (Haiku)**: Documentation (Stream 6)

**Dependencies**: Require Phase 1-2 completion (need full implementation for testing and documentation)

---

## Technical Specifications

### mnemosyne RPC Methods Used

```go
// List memories
service.List(ctx, &pb.ListRequest{
    Namespace: namespace,
    Limit:     100,
    Offset:    0,
})

// Get single memory
service.Get(ctx, &pb.GetRequest{
    Id: memoryID,
})

// Store new memory
service.Store(ctx, &pb.StoreRequest{
    Content:    content,
    Importance: importance,
    Tags:       tags,
    Namespace:  namespace,
})

// Update memory
service.Update(ctx, &pb.UpdateRequest{
    Id:         memoryID,
    Content:    content,
    Importance: importance,
    Tags:       tags,
})

// Delete memory
service.Delete(ctx, &pb.DeleteRequest{
    Id: memoryID,
})

// Create link
service.CreateLink(ctx, &pb.CreateLinkRequest{
    SourceId: sourceID,
    TargetId: targetID,
    Strength: strength,
    Type:     linkType,
})

// Delete link
service.DeleteLink(ctx, &pb.DeleteLinkRequest{
    SourceId: sourceID,
    TargetId: targetID,
})

// Traverse links (for graph)
service.TraverseLinks(ctx, &pb.TraverseLinkRequest{
    RootId: rootID,
    Depth:  depth,
})

// Search (if available)
service.Search(ctx, &pb.SearchRequest{
    Query:     query,
    Namespace: namespace,
    Limit:     limit,
})
```

### Offline Operation Queue

```go
type QueuedOperation struct {
    Type      OperationType // Create, Update, Delete, CreateLink, DeleteLink
    MemoryID  string
    Data      interface{}   // Operation-specific data
    Timestamp time.Time
}

// Queue operations when offline
func QueueOperation(op QueuedOperation) error {
    // Add to persistent queue (local storage)
}

// Process queue when online
func ProcessQueue(client MnemosyneClient) error {
    // For each queued operation:
    // 1. Execute RPC call
    // 2. If success, remove from queue
    // 3. If failure, keep in queue for retry
}
```

### Component Communication (Bubble Tea Messages)

```go
// memorylist → modes/explore
type MemorySelectedMsg struct {
    Memory *pb.MemoryNote
}

// memorydetail → modes/explore
type LinkSelectedMsg struct {
    TargetMemoryID string
}

type CloseRequestMsg struct {}

type MemoryUpdatedMsg struct {
    Memory *pb.MemoryNote
}

type MemoryDeletedMsg struct {
    MemoryID string
}

// modes/explore → memorylist
type RefreshListMsg struct {}

// modes/explore → memorygraph
type LoadGraphMsg struct {
    RootMemoryID string
    Depth        int
}
```

---

## Success Criteria

- [ ] All sample data replaced with real mnemosyne API calls
- [ ] CRUD operations fully functional (Create, Update, Delete)
- [ ] Search and filtering working (text, tags, importance, namespace)
- [ ] Link management implemented (create, delete, navigate)
- [ ] Offline mode support (queue operations, sync when online)
- [ ] Navigation history (back/forward)
- [ ] Graph visualization interactive (expand, collapse, navigate)
- [ ] 80%+ test coverage for Explore mode packages
- [ ] 10+ integration tests
- [ ] 5+ benchmark tests
- [ ] User guide complete (600+ lines)
- [ ] Zero compiler warnings
- [ ] Zero lint warnings (golangci-lint)
- [ ] All tests passing

---

## Metrics & Reporting

### Test Coverage Targets

```bash
# Generate coverage report
go test ./internal/modes -cover -coverprofile=explore_coverage.out
go test ./internal/memorylist -cover -coverprofile=list_coverage.out
go test ./internal/memorydetail -cover -coverprofile=detail_coverage.out
go test ./internal/memorygraph -cover -coverprofile=graph_coverage.out

# Target coverage by package
internal/modes (explore):    80%+
internal/memorylist:         80%+
internal/memorydetail:       80%+
internal/memorygraph:        80%+
```

### Performance Targets

- List loading (100 memories): < 100ms
- Memory detail loading: < 50ms
- Graph loading (50 nodes, 3 levels): < 500ms
- Search operation (1000 memories): < 200ms
- Graph layout iteration: < 50ms

---

## Risk Assessment

### Low Risk
- Documentation writing (isolated, no code changes)
- Search implementation (client-side filtering as fallback)
- UI improvements (incremental, testable)

### Medium Risk
- mnemosyne API integration (depends on server availability)
  - Mitigation: Use offline mode, mock server for testing
- Graph traversal (could hit performance issues with large graphs)
  - Mitigation: Depth limits, lazy loading, caching

### High Risk
- None identified (Explore Mode foundation already solid)

---

## Timeline

### Day 1-2: Foundation (Agents 1-2)
- **Morning**: Launch Agent 1 (mnemosyne integration) + Agent 2 (CRUD operations)
- **Afternoon**: Monitor progress, address any blockers

### Day 3-4: Features (Agents 3-4)
- **Morning**: Complete Agents 1-2, launch Agent 3 (search/filtering) + Agent 4 (link navigation)
- **Afternoon**: Monitor progress, verify integration

### Day 5-7: Quality (Agents 5-6)
- **Morning**: Complete Agents 3-4, launch Agent 5 (testing) + Agent 6 (documentation)
- **Afternoon**: Final verification, performance testing

---

## Deliverables Summary

**Code** (10+ files):
- internal/memorylist/realdata.go (mnemosyne integration)
- internal/memorydetail/crud.go (CRUD operations)
- internal/memorydetail/links.go (link management)
- internal/memorylist/search.go (search/filtering)
- internal/memorygraph/loader.go (graph loading - new file)
- internal/memorygraph/navigation.go (graph navigation - new file)
- internal/modes/explore.go (link navigation, history)
- 20+ test files (*_test.go)
- 5+ benchmark files (*_benchmark_test.go)

**Documentation** (600+ lines):
- docs/explore-mode-guide.md (comprehensive user guide)
- Updated README.md

**Tests** (50+ tests):
- 40+ unit tests
- 10+ integration tests
- 5+ benchmarks

---

## Future Enhancements (Post-Phase 9)

After Explore Mode completion:
- **Phase 10**: Production deployment
  - Dockerize application
  - CI/CD pipeline
  - Release process
  - Monitoring and logging

- **Phase 11**: Advanced features
  - Collaborate Mode (multi-user editing)
  - Cloud sync
  - Mobile companion app
  - Plugin system

---

**Last Updated**: 2025-11-09
**Author**: Phase 9 Planning
**Status**: Ready for Execution

