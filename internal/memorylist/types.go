package memorylist

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// InputMode defines the current input mode.
type InputMode int

const (
	InputModeNormal InputMode = iota
	InputModeSearch
	InputModeFilter
)

// Model represents the memory list component state.
type Model struct {
	// Memory data
	memories      []*pb.MemoryNote
	filteredMems  []*pb.MemoryNote
	selectedIndex int

	// Viewport state
	scrollOffset int
	height       int
	width        int

	// Filter/search state
	searchQuery   string
	searchInput   string // Temporary input while typing
	filterTags    []string
	filterNS      string
	minImportance uint32

	// Sort state
	sortBy   SortMode
	sortDesc bool

	// UI state
	loading    bool
	err        error
	focused    bool
	inputMode  InputMode
	showHelp   bool

	// Pagination
	pageSize   int
	totalCount uint32

	// Client integration (optional)
	client     interface{} // Can hold *mnemosyne.Client
	loadOpts   LoadOptions
	autoReload bool

	// NEW: Real data support
	mnemosyneClient *mnemosyne.Client
	queryCache      *QueryCache
	currentQuery    string
	currentFilters  Filters
}

// SortMode defines how memories are sorted.
type SortMode int

const (
	SortByImportance SortMode = iota
	SortByUpdated
	SortByCreated
	SortByRelevance
)

// String returns the string representation of the sort mode.
func (s SortMode) String() string {
	switch s {
	case SortByImportance:
		return "Importance"
	case SortByUpdated:
		return "Updated"
	case SortByCreated:
		return "Created"
	case SortByRelevance:
		return "Relevance"
	default:
		return "Unknown"
	}
}

// Messages for the memory list component.
type (
	// MemoriesLoadedMsg is sent when memories are loaded from mnemosyne.
	MemoriesLoadedMsg struct {
		Memories   []*pb.MemoryNote
		TotalCount uint32
		Err        error
	}

	// MemoriesErrorMsg is sent when memory loading fails.
	MemoriesErrorMsg struct {
		Err error
	}

	// MemorySelectedMsg is sent when a memory is selected.
	MemorySelectedMsg struct {
		Memory *pb.MemoryNote
		Index  int
	}

	// SearchQueryMsg is sent when search query changes.
	SearchQueryMsg struct {
		Query string
	}

	// FilterChangedMsg is sent when filters change.
	FilterChangedMsg struct {
		Tags          []string
		Namespace     string
		MinImportance uint32
	}

	// SortChangedMsg is sent when sort mode changes.
	SortChangedMsg struct {
		Mode SortMode
		Desc bool
	}
)

// NewModel creates a new memory list model with default settings.
func NewModel() Model {
	return Model{
		memories:      make([]*pb.MemoryNote, 0),
		filteredMems:  make([]*pb.MemoryNote, 0),
		selectedIndex: 0,
		scrollOffset:  0,
		height:        20,
		width:         80,
		sortBy:        SortByUpdated,
		sortDesc:      true, // Most recent first
		loading:       false,
		focused:       true,
		pageSize:      50,
		loadOpts:      DefaultLoadOptions(),
		autoReload:    false,
	}
}

// NewModelWithClient creates a new model with a mnemosyne client.
func NewModelWithClient(client interface{}) Model {
	m := NewModel()
	m.client = client
	m.autoReload = true
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	// If we have a client and autoReload is enabled, load initial data
	if m.client != nil && m.autoReload {
		// Type assert to get the actual client
		// This is safe because we control what gets passed in
		return func() tea.Msg {
			return LoadRequestMsg{}
		}
	}
	return nil
}

// LoadRequestMsg requests a data reload.
type LoadRequestMsg struct{}

// ReloadRequestMsg requests a data reload with current filters.
type ReloadRequestMsg struct{}

// SetSize sets the component size.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocus sets the focus state.
func (m *Model) SetFocus(focused bool) {
	m.focused = focused
}

// IsFocused returns whether the component is focused.
func (m Model) IsFocused() bool {
	return m.focused
}

// SetMemories sets the memory list.
func (m *Model) SetMemories(memories []*pb.MemoryNote, totalCount uint32) {
	m.memories = memories
	m.totalCount = totalCount
	m.loading = false
	m.err = nil
	m.applyFilters()
}

// SetLoading sets the loading state.
func (m *Model) SetLoading(loading bool) {
	m.loading = loading
}

// SetError sets the error state.
func (m *Model) SetError(err error) {
	m.err = err
	m.loading = false
}

// SelectedMemory returns the currently selected memory, or nil if none.
func (m Model) SelectedMemory() *pb.MemoryNote {
	if len(m.filteredMems) == 0 || m.selectedIndex < 0 || m.selectedIndex >= len(m.filteredMems) {
		return nil
	}
	return m.filteredMems[m.selectedIndex]
}

// SelectedIndex returns the current selection index.
func (m Model) SelectedIndex() int {
	return m.selectedIndex
}

// MemoryCount returns the total number of filtered memories.
func (m Model) MemoryCount() int {
	return len(m.filteredMems)
}

// TotalCount returns the total count from the server (before pagination).
func (m Model) TotalCount() uint32 {
	return m.totalCount
}

// IsLoading returns whether the component is loading.
func (m Model) IsLoading() bool {
	return m.loading
}

// Error returns the current error, if any.
func (m Model) Error() error {
	return m.err
}

// SetClient sets the mnemosyne client for this model.
func (m *Model) SetClient(client interface{}) {
	m.client = client
}

// Client returns the mnemosyne client, if set.
func (m Model) Client() interface{} {
	return m.client
}

// LoadOptions returns the current load options.
func (m Model) LoadOptions() LoadOptions {
	return m.loadOpts
}

// SetLoadOptions sets the load options.
func (m *Model) SetLoadOptions(opts LoadOptions) {
	m.loadOpts = opts
}

// SearchQuery returns the current search query.
func (m Model) SearchQuery() string {
	return m.searchQuery
}

// QueryCache provides thread-safe caching of query results with TTL.
type QueryCache struct {
	entries map[string]*CacheEntry
	maxAge  time.Duration
	mu      sync.RWMutex
}

// CacheEntry represents a cached query result with timestamp.
type CacheEntry struct {
	memories  []*pb.MemoryNote
	timestamp time.Time
}

// Filters holds filter criteria for loading memories.
type Filters struct {
	Namespace     string
	Tags          []string
	MinImportance int
	MaxImportance int
}

// SetMnemosyneClient sets the mnemosyne client for real data loading.
func (m *Model) SetMnemosyneClient(client *mnemosyne.Client) {
	m.mnemosyneClient = client
	// Initialize query cache with 5 minute TTL
	if m.queryCache == nil {
		m.queryCache = NewQueryCache(5 * time.Minute)
	}
}

// MnemosyneClient returns the mnemosyne client, if set.
func (m Model) MnemosyneClient() *mnemosyne.Client {
	return m.mnemosyneClient
}

// QueryCache returns the query cache.
func (m Model) GetQueryCache() *QueryCache {
	return m.queryCache
}

// CurrentQuery returns the current query string.
func (m Model) CurrentQuery() string {
	return m.currentQuery
}

// CurrentFilters returns the current filters.
func (m Model) CurrentFilters() Filters {
	return m.currentFilters
}

// SetCurrentQuery sets the current query.
func (m *Model) SetCurrentQuery(query string) {
	m.currentQuery = query
}

// SetCurrentFilters sets the current filters.
func (m *Model) SetCurrentFilters(filters Filters) {
	m.currentFilters = filters
}

// NewQueryCache creates a new query cache with the specified TTL.
func NewQueryCache(maxAge time.Duration) *QueryCache {
	return &QueryCache{
		entries: make(map[string]*CacheEntry),
		maxAge:  maxAge,
	}
}

// Get retrieves cached memories for a query if they exist and are not expired.
// Returns the memories and true if found and valid, nil and false otherwise.
func (qc *QueryCache) Get(query string) ([]*pb.MemoryNote, bool) {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	entry, exists := qc.entries[query]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if time.Since(entry.timestamp) > qc.maxAge {
		return nil, false
	}

	return entry.memories, true
}

// Set stores memories for a query in the cache.
func (qc *QueryCache) Set(query string, memories []*pb.MemoryNote) {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	qc.entries[query] = &CacheEntry{
		memories:  memories,
		timestamp: time.Now(),
	}
}

// Invalidate removes a specific query from the cache.
func (qc *QueryCache) Invalidate(query string) {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	delete(qc.entries, query)
}

// Clear removes all entries from the cache.
func (qc *QueryCache) Clear() {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	qc.entries = make(map[string]*CacheEntry)
}
