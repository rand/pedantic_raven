# Phase 5 Technical Specification: Real mnemosyne Integration

**Status**: In Progress
**Timeline**: 1-2 weeks
**Prerequisites**: Phase 4 Complete (Explore Mode with sample data)

---

## Overview

Phase 5 connects the Explore Mode to a live mnemosyne-rpc server, replacing sample data with real memory queries. This phase transforms Pedantic Raven from a prototype into a functional memory management interface.

**Core Objective**: Enable full CRUD operations on real mnemosyne memories with robust error handling and performance optimization.

---

## Architecture

### Current State (Phase 4)

```
┌─────────────────────────────────────┐
│         Explore Mode                │
├─────────────────────────────────────┤
│  MemoryList  │  MemoryDetail        │
│              │                      │
│  [Sample     │  [Sample Memory]     │
│   Data]      │                      │
└──────────────┴──────────────────────┘
           │
           ▼
    MemoryGraph (Sample Data)
```

### Target State (Phase 5)

```
┌─────────────────────────────────────┐
│         Explore Mode                │
├─────────────────────────────────────┤
│  MemoryList  │  MemoryDetail        │
│              │                      │
│  [Live       │  [Live Memory]       │
│   Query]     │  + Edit Mode         │
└──────────────┴──────────────────────┘
           │
           ▼
    MemoryGraph (Live Graph Traverse)
           │
           ▼
┌─────────────────────────────────────┐
│    Mnemosyne Client (existing)      │
│  - Recall()                          │
│  - GetMemory()                       │
│  - StoreMemory()                     │
│  - UpdateMemory()                    │
│  - DeleteMemory()                    │
│  - GraphTraverse()                   │
└─────────────────────────────────────┘
           │
           ▼ gRPC
┌─────────────────────────────────────┐
│    mnemosyne-rpc Server              │
│  - Real vector database              │
│  - Semantic search (embeddings)      │
│  - Graph storage (links)             │
└─────────────────────────────────────┘
```

---

## Components

### 5.1: Connection Management

**File**: `internal/mnemosyne/connection.go` (new)

**Purpose**: Manage persistent connection to mnemosyne-rpc server with health monitoring.

**Interface**:
```go
type ConnectionManager struct {
    client      *Client
    config      *ConnectionConfig
    healthCheck *time.Ticker
    status      ConnectionStatus
    mu          sync.RWMutex
}

type ConnectionConfig struct {
    Host         string
    Port         int
    UseTLS       bool
    Timeout      time.Duration
    RetryPolicy  RetryPolicy
}

type ConnectionStatus int
const (
    StatusDisconnected ConnectionStatus = iota
    StatusConnecting
    StatusConnected
    StatusReconnecting
    StatusFailed
)

// Core methods
func NewConnectionManager(config *ConnectionConfig) *ConnectionManager
func (cm *ConnectionManager) Connect() error
func (cm *ConnectionManager) Disconnect() error
func (cm *ConnectionManager) Status() ConnectionStatus
func (cm *ConnectionManager) HealthCheck() error
```

**Features**:
- Automatic reconnection on failure (exponential backoff)
- Health check monitoring (every 30s)
- Thread-safe status tracking
- Configuration validation
- Connection pooling (if needed)

**Tests**: 15+ tests covering:
- Successful connection
- Connection failure handling
- Reconnection logic
- Health check behavior
- Status transitions
- Configuration validation

---

### 5.2: Real Data Integration

**File**: `internal/memorylist/realdata.go` (new)

**Purpose**: Replace sample data loading with live mnemosyne queries.

**Changes to `internal/memorylist/model.go`**:
```go
type Model struct {
    // ... existing fields ...

    // NEW: Real data support
    mnemosyneClient *mnemosyne.Client
    connectionMgr   *mnemosyne.ConnectionManager
    queryCache      *QueryCache
    currentQuery    string
    currentFilters  Filters
}

type Filters struct {
    Namespace  string
    Tags       []string
    MinImportance int
    MaxImportance int
}

// NEW: Real data commands
func LoadMemoriesFromServer(client *mnemosyne.Client, filters Filters) tea.Cmd
func SearchMemories(client *mnemosyne.Client, query string) tea.Cmd
func RefreshMemories(client *mnemosyne.Client) tea.Cmd
```

**Query Strategy**:
1. **Initial Load**: `Recall("")` with limit 100 (most recent/important)
2. **Search**: `Recall(query)` with semantic + FTS
3. **Filter**: Client-side filtering on loaded results (+ server re-query if needed)
4. **Pagination**: Load more on scroll to bottom

**Caching**:
```go
type QueryCache struct {
    entries map[string]*CacheEntry
    maxAge  time.Duration
    mu      sync.RWMutex
}

type CacheEntry struct {
    memories  []*pb.Memory
    timestamp time.Time
}
```

**Features**:
- Cache results for 5 minutes
- Invalidate on CRUD operations
- Background refresh option
- Loading states (spinner, progress)

**Tests**: 20+ tests covering:
- Loading from server
- Search integration
- Filter application
- Pagination
- Cache hit/miss
- Error handling (server down, timeout)

---

### 5.3: Memory CRUD Operations

**File**: `internal/memorydetail/crud.go` (new)

**Purpose**: Enable creating, editing, updating, and deleting memories.

**Edit Mode**:
```go
type EditState struct {
    isEditing    bool
    editedMemory *pb.Memory
    fieldFocus   EditField
    originalHash string  // For change detection
}

type EditField int
const (
    FieldTitle EditField = iota
    FieldContent
    FieldTags
    FieldImportance
    FieldNamespace
)

// NEW: CRUD operations
func (m *Model) EnterEditMode() tea.Cmd
func (m *Model) SaveChanges() tea.Cmd
func (m *Model) CancelEdit() tea.Cmd
func (m *Model) DeleteMemory() tea.Cmd
func (m *Model) CreateMemory() tea.Cmd
```

**UI Flow**:

**View Mode** (current):
```
┌────────────────────────────────────┐
│ Memory Detail                     │
│                                    │
│ Title: Understanding Force Layout │
│ Namespace: project:pedantic-raven  │
│ Importance: ████████ (8/10)       │
│ Tags: graph, visualization         │
│                                    │
│ Content:                           │
│ Force-directed layouts use...      │
│                                    │
│ Links (3 outbound, 2 inbound)     │
│                                    │
│ [e] Edit  [d] Delete  [n] New     │
└────────────────────────────────────┘
```

**Edit Mode** (new):
```
┌────────────────────────────────────┐
│ Editing Memory                    ⚠│
│                                    │
│ Title: [Understanding Force Layout_]
│ Namespace: [project:pedantic-raven]│
│ Importance: [8] (1-10)             │
│ Tags: [graph,visualization_]       │
│                                    │
│ Content: [Tab to edit]             │
│ ┌────────────────────────────────┐ │
│ │Force-directed layouts use..._  │ │
│ │                                │ │
│ └────────────────────────────────┘ │
│                                    │
│ [Ctrl+S] Save  [Esc] Cancel        │
│ * Unsaved changes                  │
└────────────────────────────────────┘
```

**Delete Confirmation**:
```
┌────────────────────────────────────┐
│  ⚠ Confirm Delete                  │
│                                    │
│  Delete memory "Understanding      │
│  Force Layout"?                    │
│                                    │
│  This action cannot be undone.     │
│  3 outbound links will be removed. │
│                                    │
│  [Enter] Confirm  [Esc] Cancel     │
└────────────────────────────────────┘
```

**Features**:
- In-place editing with multi-field support
- Change detection (warn on unsaved changes)
- Validation (required fields, importance range)
- Optimistic UI updates
- Rollback on server failure
- Confirmation dialogs for destructive actions

**Tests**: 25+ tests covering:
- Enter/exit edit mode
- Field editing (each field type)
- Save changes (success/failure)
- Cancel with/without changes
- Delete with confirmation
- Create new memory
- Validation errors

---

### 5.4: Link Management

**File**: `internal/memorydetail/links.go` (new)

**Purpose**: Create, navigate, and manage bidirectional links between memories.

**Link Operations**:
```go
type LinkManager struct {
    client *mnemosyne.Client
}

// Core methods
func (lm *LinkManager) CreateLink(sourceID, targetID string, linkType pb.LinkType) error
func (lm *LinkManager) DeleteLink(linkID string) error
func (lm *LinkManager) UpdateLinkMetadata(linkID string, metadata *pb.LinkMetadata) error
func (lm *LinkManager) GetLinkedMemories(memoryID string, direction LinkDirection) ([]*pb.Memory, error)

type LinkDirection int
const (
    DirectionOutbound LinkDirection = iota
    DirectionInbound
    DirectionBoth
)
```

**UI for Link Creation**:
```
┌────────────────────────────────────┐
│ Create Link                        │
│                                    │
│ From: Understanding Force Layout   │
│ To: [Search or select memory_]     │
│                                    │
│ Type:                              │
│ ( ) Related                        │
│ (•) Cites                          │
│ ( ) Extends                        │
│ ( ) Contradicts                    │
│                                    │
│ Strength: [5] (1-10)               │
│                                    │
│ [Enter] Create  [Esc] Cancel       │
└────────────────────────────────────┘
```

**Link Navigation**:
- In detail view, links are shown as list
- Press Enter on link → navigate to linked memory
- Update detail view and list selection
- Breadcrumb trail for navigation history

**Features**:
- Bidirectional link creation
- Link type specification (related, cites, extends, contradicts, etc.)
- Link strength/weight metadata
- Link deletion with confirmation
- Navigation history (back/forward)

**Tests**: 20+ tests covering:
- Create link (various types)
- Delete link
- Update link metadata
- Navigate to linked memory
- Get inbound/outbound links
- Bidirectional consistency

---

### 5.5: Search Integration

**File**: `internal/memorylist/search.go` (enhanced)

**Purpose**: Integrate live semantic search from mnemosyne server.

**Enhanced Search**:
```go
type SearchOptions struct {
    Query         string
    Namespaces    []string
    Tags          []string
    MinImportance int
    MaxImportance int
    MaxResults    int
    SearchMode    SearchMode
}

type SearchMode int
const (
    SearchHybrid SearchMode = iota  // Semantic + FTS + Graph (default)
    SearchSemantic                   // Pure embedding search
    SearchFullText                   // FTS only
    SearchGraph                      // Graph traversal
)

// Enhanced search command
func SearchWithOptions(client *mnemosyne.Client, opts SearchOptions) tea.Cmd
```

**Search UI**:
```
┌────────────────────────────────────┐
│ Search: force layout physics____   │
│                                    │
│ Filters:                           │
│ Namespace: [All ▼]                 │
│ Tags: [graph, algorithms]          │
│ Importance: [5] - [10]             │
│ Mode: [Hybrid ▼]                   │
│                                    │
│ Results (42):                      │
│ ┌────────────────────────────────┐ │
│ │▶ Understanding Force Layout  9 │ │
│ │  Physics simulation for...     │ │
│ │  Relevance: ████████ (0.89)    │ │
│ │                                │ │
│ │  Adaptive Force Layout      8  │ │
│ │  Dynamic adjustment of...      │ │
│ │  Relevance: ███████ (0.76)     │ │
│ └────────────────────────────────┘ │
└────────────────────────────────────┘
```

**Features**:
- Debounced search (500ms after typing stops)
- Relevance scoring display
- Search mode selection
- Advanced filters (namespace, tags, importance)
- Search result highlighting
- Search history

**Tests**: 15+ tests covering:
- Hybrid search
- Semantic-only search
- Full-text search
- Graph traversal search
- Filter application
- Debouncing
- Result ranking

---

### 5.6: Error Handling & Offline Mode

**File**: `internal/mnemosyne/errors.go` (enhanced)

**Purpose**: Robust error handling with graceful degradation.

**Error Categories**:
```go
type ErrorCategory int
const (
    ErrCategoryConnection ErrorCategory = iota
    ErrCategoryServer
    ErrCategoryValidation
    ErrCategoryTimeout
    ErrCategoryUnknown
)

type MnemosyneError struct {
    Category ErrorCategory
    Code     string
    Message  string
    Retryable bool
    Underlying error
}

// Error handling utilities
func CategorizeError(err error) ErrorCategory
func IsRetryable(err error) bool
func GetUserMessage(err error) string
```

**Offline Mode**:
```go
type OfflineCache struct {
    memories map[string]*pb.Memory
    lastSync time.Time
    dirty    map[string]bool  // Unsaved changes
}

// Offline operations
func (oc *OfflineCache) Store(memory *pb.Memory)
func (oc *OfflineCache) Get(id string) *pb.Memory
func (oc *OfflineCache) ListAll() []*pb.Memory
func (oc *OfflineCache) Sync(client *mnemosyne.Client) error
```

**Error UI**:
```
┌────────────────────────────────────┐
│ ⚠ Connection Lost                  │
│                                    │
│ Cannot reach mnemosyne server.     │
│                                    │
│ Working in offline mode with       │
│ cached data (last sync: 2m ago).   │
│                                    │
│ Changes will be synced when        │
│ connection is restored.            │
│                                    │
│ [r] Retry Now  [Esc] Dismiss       │
└────────────────────────────────────┘
```

**Features**:
- Automatic retry with exponential backoff
- Offline mode with local cache
- Sync queue for offline changes
- User-friendly error messages
- Connection status indicator
- Manual retry option

**Tests**: 20+ tests covering:
- Error categorization
- Retry logic (success/failure)
- Offline mode activation
- Cache storage/retrieval
- Sync after reconnection
- Error message formatting

---

## Implementation Plan

### Week 1: Core Integration

**Days 1-2: Connection Management & Real Data**
- Implement ConnectionManager with health checks
- Update MemoryList to load from live server
- Add query caching
- Tests: Connection management (15), Real data (20)

**Days 3-4: CRUD Operations**
- Implement edit mode in MemoryDetail
- Add create, update, delete operations
- Add confirmation dialogs
- Tests: CRUD operations (25)

**Day 5: Search Integration**
- Enhance search with live queries
- Add advanced filters
- Implement search modes
- Tests: Search integration (15)

### Week 2: Advanced Features & Polish

**Days 6-7: Link Management**
- Implement link creation UI
- Add link navigation
- Support bidirectional links
- Tests: Link operations (20)

**Days 8-9: Error Handling & Offline Mode**
- Enhanced error handling
- Offline mode with cache
- Sync queue for offline changes
- Tests: Error handling (20)

**Day 10: Integration Testing & Documentation**
- End-to-end integration tests
- Performance testing (1000+ memories)
- Update README and user guide
- Create PHASE5_SUMMARY.md

---

## Success Criteria

### Functionality
- ✅ Connects to live mnemosyne-rpc server
- ✅ Loads real memories (not sample data)
- ✅ Create, edit, update, delete operations work
- ✅ Search returns live results
- ✅ Links can be created and navigated
- ✅ Offline mode activates on connection loss

### Performance
- ✅ Initial load < 2s for 100 memories
- ✅ Search results < 1s
- ✅ CRUD operations < 500ms
- ✅ Smooth scrolling with 1000+ cached memories
- ✅ Memory usage < 150MB with full cache

### Reliability
- ✅ Zero data loss on connection failure
- ✅ Automatic reconnection works
- ✅ Offline changes sync correctly
- ✅ All errors have recovery paths
- ✅ No crashes on invalid server responses

### User Experience
- ✅ Clear loading indicators
- ✅ Intuitive edit mode
- ✅ Helpful error messages
- ✅ Connection status visible
- ✅ Unsaved changes warnings

---

## Testing Strategy

### Unit Tests (~100 new tests)
- Connection management (15)
- Real data loading (20)
- CRUD operations (25)
- Link management (20)
- Search integration (15)
- Error handling (20)

### Integration Tests (~30 new tests)
- Full CRUD workflows
- Search → view → edit → save
- Link creation → navigation
- Connection loss → offline → reconnect → sync
- Concurrent operations

### Manual Testing Scenarios
1. **Happy Path**: Connect → search → view → edit → save
2. **Link Management**: View memory → create link → navigate → delete link
3. **Error Recovery**: Disconnect server → work offline → reconnect → verify sync
4. **Performance**: Load 1000 memories → scroll → search → edit
5. **Edge Cases**: Empty results, very long content, special characters

---

## Configuration

**Server Configuration** (`~/.config/pedantic-raven/config.yaml`):
```yaml
mnemosyne:
  host: localhost
  port: 50051
  tls:
    enabled: false
    cert_path: ""
  connection:
    timeout: 10s
    retry_policy:
      max_attempts: 3
      initial_backoff: 1s
      max_backoff: 30s
  cache:
    max_age: 5m
    max_size: 1000
```

**Environment Variables**:
```bash
MNEMOSYNE_HOST=localhost
MNEMOSYNE_PORT=50051
MNEMOSYNE_TLS=false
```

---

## Dependencies

### Required
- mnemosyne-rpc server (running locally or remote)
- Existing mnemosyne client (`internal/mnemosyne/*`)
- gRPC and Protocol Buffers

### Optional
- mnemosyne server with populated data (for realistic testing)
- TLS certificates (for production)

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Server not available | High | Offline mode with cache, clear error messages |
| Slow network | Medium | Caching, pagination, loading indicators |
| Data corruption | High | Validation, rollback on error, confirmation dialogs |
| Memory leaks (large cache) | Medium | Cache size limits, LRU eviction |
| Concurrent modifications | Medium | Optimistic locking, conflict resolution UI |

---

## Migration from Phase 4

### Code Changes
1. **MemoryList**:
   - Add `mnemosyneClient` field
   - Replace `LoadSampleData()` with `LoadMemoriesFromServer()`
   - Update `Init()` to accept client

2. **MemoryDetail**:
   - Add edit mode state
   - Implement CRUD operation handlers
   - Add link management UI

3. **Explore Mode**:
   - Pass mnemosyne client to components
   - Add connection status indicator
   - Handle connection events

### Data Migration
- No data migration needed (server-side data)
- Sample data removed (keep for testing)

### Backward Compatibility
- Keep sample data mode for offline development
- Environment variable to toggle: `PEDANTIC_RAVEN_USE_SAMPLE_DATA=true`

---

## Documentation Updates

### README.md
- Update setup instructions (mnemosyne server requirement)
- Add configuration section
- Document keyboard shortcuts for edit mode

### User Guide
- Add "Connecting to mnemosyne" section
- Document CRUD operations
- Explain offline mode

### Developer Guide
- mnemosyne client usage patterns
- Error handling best practices
- Testing with mock server

---

## Next Steps (Phase 6 Preview)

After Phase 5, the next logical phase is **Analyze Mode** (statistical analysis):
- Entity frequency analysis
- Relationship pattern mining
- Typed hole prioritization
- Dependency tree visualization

However, Phase 6 depends on having real data in mnemosyne (provided by Phase 5).

---

**Document Version**: 1.0
**Created**: 2025-11-08
**Author**: Development Team
