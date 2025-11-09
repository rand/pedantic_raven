package memorylist

import (
	"context"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// searchClient defines the interface needed for search operations.
// This allows for easier testing with mock implementations.
type searchClient interface {
	IsConnected() bool
	Recall(ctx context.Context, opts mnemosyne.RecallOptions) ([]*pb.SearchResult, error)
	GraphTraverse(ctx context.Context, opts mnemosyne.GraphTraverseOptions) ([]*pb.MemoryNote, []*pb.GraphEdge, error)
}

// SearchMode defines the type of search to perform.
type SearchMode int

const (
	SearchHybrid   SearchMode = iota // Semantic + FTS + Graph (default)
	SearchSemantic                    // Pure embedding search
	SearchFullText                    // FTS only
	SearchGraph                       // Graph traversal
)

// String returns the string representation of the search mode.
func (s SearchMode) String() string {
	switch s {
	case SearchHybrid:
		return "Hybrid"
	case SearchSemantic:
		return "Semantic"
	case SearchFullText:
		return "Full-Text"
	case SearchGraph:
		return "Graph"
	default:
		return "Unknown"
	}
}

// SearchOptions configures search behavior.
type SearchOptions struct {
	Query         string
	Namespaces    []string
	Tags          []string
	MinImportance int
	MaxImportance int
	MaxResults    int
	SearchMode    SearchMode
}

// DefaultSearchOptions returns default search options.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Query:         "",
		Namespaces:    nil,
		Tags:          nil,
		MinImportance: 0,
		MaxImportance: 10,
		MaxResults:    100,
		SearchMode:    SearchHybrid,
	}
}

// SearchResultsMsg is sent when search results are available.
type SearchResultsMsg struct {
	Results    []*pb.MemoryNote
	Query      string
	TotalCount uint32
	Err        error
}

// SearchDebouncer manages debounced search requests.
type SearchDebouncer struct {
	timer    *time.Timer
	mu       sync.Mutex
	delay    time.Duration
	callback func()
}

// NewSearchDebouncer creates a new search debouncer with the specified delay.
func NewSearchDebouncer(delay time.Duration) *SearchDebouncer {
	return &SearchDebouncer{
		delay: delay,
	}
}

// Debounce cancels any pending search and schedules a new one.
func (sd *SearchDebouncer) Debounce(fn func()) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	// Cancel existing timer
	if sd.timer != nil {
		sd.timer.Stop()
	}

	// Create new timer
	sd.callback = fn
	sd.timer = time.AfterFunc(sd.delay, func() {
		if sd.callback != nil {
			sd.callback()
		}
	})
}

// Cancel stops any pending search.
func (sd *SearchDebouncer) Cancel() {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.timer != nil {
		sd.timer.Stop()
		sd.timer = nil
	}
	sd.callback = nil
}

// SearchHistory tracks recent searches.
type SearchHistory struct {
	queries []string
	maxSize int
	mu      sync.RWMutex
}

// NewSearchHistory creates a new search history with the specified max size.
func NewSearchHistory(maxSize int) *SearchHistory {
	return &SearchHistory{
		queries: make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add adds a query to the search history.
func (sh *SearchHistory) Add(query string) {
	if query == "" {
		return
	}

	sh.mu.Lock()
	defer sh.mu.Unlock()

	// Check if query already exists
	for i, q := range sh.queries {
		if q == query {
			// Move to front
			sh.queries = append(sh.queries[:i], sh.queries[i+1:]...)
			break
		}
	}

	// Add to front
	sh.queries = append([]string{query}, sh.queries...)

	// Trim to max size
	if len(sh.queries) > sh.maxSize {
		sh.queries = sh.queries[:sh.maxSize]
	}
}

// Get returns all queries in the search history.
func (sh *SearchHistory) Get() []string {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	// Return a copy
	result := make([]string, len(sh.queries))
	copy(result, sh.queries)
	return result
}

// Clear removes all queries from the search history.
func (sh *SearchHistory) Clear() {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.queries = make([]string, 0, sh.maxSize)
}

// SearchWithOptions creates a Bubble Tea command that performs search with the specified options.
func SearchWithOptions(client *mnemosyne.Client, opts SearchOptions) tea.Cmd {
	return searchWithClient(client, opts)
}

// searchWithClient is an internal function that works with the searchClient interface.
func searchWithClient(client searchClient, opts SearchOptions) tea.Cmd {
	return func() tea.Msg {
		if client == nil || !client.IsConnected() {
			return SearchResultsMsg{
				Results:    nil,
				Query:      opts.Query,
				TotalCount: 0,
				Err:        mnemosyne.ErrNotConnected,
			}
		}

		if opts.Query == "" {
			return SearchResultsMsg{
				Results:    []*pb.MemoryNote{},
				Query:      opts.Query,
				TotalCount: 0,
				Err:        nil,
			}
		}

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Execute search based on mode
		var results []*pb.MemoryNote
		var err error

		switch opts.SearchMode {
		case SearchHybrid:
			results, err = searchHybrid(ctx, client, opts)
		case SearchSemantic:
			results, err = searchSemantic(ctx, client, opts)
		case SearchFullText:
			results, err = searchFullText(ctx, client, opts)
		case SearchGraph:
			results, err = searchGraph(ctx, client, opts)
		default:
			results, err = searchHybrid(ctx, client, opts)
		}

		if err != nil {
			return SearchResultsMsg{
				Results:    nil,
				Query:      opts.Query,
				TotalCount: 0,
				Err:        err,
			}
		}

		// Apply client-side filters
		results = applySearchFilters(results, opts)

		return SearchResultsMsg{
			Results:    results,
			Query:      opts.Query,
			TotalCount: uint32(len(results)),
			Err:        nil,
		}
	}
}

// searchHybrid performs hybrid search (semantic + FTS + graph).
func searchHybrid(ctx context.Context, client searchClient, opts SearchOptions) ([]*pb.MemoryNote, error) {
	// Build namespace filter if specified
	var namespace *pb.Namespace
	if len(opts.Namespaces) > 0 {
		// Use first namespace for now
		namespace = parseNamespaceString(opts.Namespaces[0])
	}

	// Build minimum importance filter
	var minImportance *uint32
	if opts.MinImportance > 0 {
		imp := uint32(opts.MinImportance)
		minImportance = &imp
	}

	// Create recall options
	recallOpts := mnemosyne.RecallOptions{
		Query:         opts.Query,
		Namespace:     namespace,
		MaxResults:    uint32(opts.MaxResults),
		MinImportance: minImportance,
		Tags:          opts.Tags,
	}

	// Perform hybrid search
	searchResults, err := client.Recall(ctx, recallOpts)
	if err != nil {
		return nil, err
	}

	// Extract memories from search results
	memories := make([]*pb.MemoryNote, 0, len(searchResults))
	for _, result := range searchResults {
		if result.Memory != nil {
			memories = append(memories, result.Memory)
		}
	}

	return memories, nil
}

// searchSemantic performs pure semantic search.
// Note: This requires generating embeddings from the query.
// For now, we use Recall with high semantic weight.
func searchSemantic(ctx context.Context, client searchClient, opts SearchOptions) ([]*pb.MemoryNote, error) {
	// Build namespace filter
	var namespace *pb.Namespace
	if len(opts.Namespaces) > 0 {
		namespace = parseNamespaceString(opts.Namespaces[0])
	}

	// Build minimum importance filter
	var minImportance *uint32
	if opts.MinImportance > 0 {
		imp := uint32(opts.MinImportance)
		minImportance = &imp
	}

	// Use high semantic weight for semantic-focused search
	semanticWeight := float32(0.95)
	ftsWeight := float32(0.05)
	graphWeight := float32(0.0)

	recallOpts := mnemosyne.RecallOptions{
		Query:          opts.Query,
		Namespace:      namespace,
		MaxResults:     uint32(opts.MaxResults),
		MinImportance:  minImportance,
		Tags:           opts.Tags,
		SemanticWeight: &semanticWeight,
		FtsWeight:      &ftsWeight,
		GraphWeight:    &graphWeight,
	}

	// Perform semantic search via Recall
	searchResults, err := client.Recall(ctx, recallOpts)
	if err != nil {
		return nil, err
	}

	// Extract memories
	memories := make([]*pb.MemoryNote, 0, len(searchResults))
	for _, result := range searchResults {
		if result.Memory != nil {
			memories = append(memories, result.Memory)
		}
	}

	return memories, nil
}

// searchFullText performs full-text search only.
func searchFullText(ctx context.Context, client searchClient, opts SearchOptions) ([]*pb.MemoryNote, error) {
	// Build namespace filter
	var namespace *pb.Namespace
	if len(opts.Namespaces) > 0 {
		namespace = parseNamespaceString(opts.Namespaces[0])
	}

	// Build minimum importance filter
	var minImportance *uint32
	if opts.MinImportance > 0 {
		imp := uint32(opts.MinImportance)
		minImportance = &imp
	}

	// Use high FTS weight for full-text search
	semanticWeight := float32(0.0)
	ftsWeight := float32(1.0)
	graphWeight := float32(0.0)

	recallOpts := mnemosyne.RecallOptions{
		Query:          opts.Query,
		Namespace:      namespace,
		MaxResults:     uint32(opts.MaxResults),
		MinImportance:  minImportance,
		Tags:           opts.Tags,
		SemanticWeight: &semanticWeight,
		FtsWeight:      &ftsWeight,
		GraphWeight:    &graphWeight,
	}

	// Perform FTS search via Recall
	searchResults, err := client.Recall(ctx, recallOpts)
	if err != nil {
		return nil, err
	}

	// Extract memories
	memories := make([]*pb.MemoryNote, 0, len(searchResults))
	for _, result := range searchResults {
		if result.Memory != nil {
			memories = append(memories, result.Memory)
		}
	}

	return memories, nil
}

// searchGraph performs graph traversal search.
// This requires seed nodes, so we first do a quick search to find seeds,
// then traverse from those seeds.
func searchGraph(ctx context.Context, client searchClient, opts SearchOptions) ([]*pb.MemoryNote, error) {
	// First, get some seed nodes via a quick hybrid search
	var namespace *pb.Namespace
	if len(opts.Namespaces) > 0 {
		namespace = parseNamespaceString(opts.Namespaces[0])
	}

	seedOpts := mnemosyne.RecallOptions{
		Query:      opts.Query,
		Namespace:  namespace,
		MaxResults: 5, // Just need a few seeds
	}

	seedResults, err := client.Recall(ctx, seedOpts)
	if err != nil {
		return nil, err
	}

	// Extract seed IDs
	seedIDs := make([]string, 0, len(seedResults))
	for _, result := range seedResults {
		if result.Memory != nil && result.Memory.Id != "" {
			seedIDs = append(seedIDs, result.Memory.Id)
		}
	}

	if len(seedIDs) == 0 {
		// No seeds found
		return []*pb.MemoryNote{}, nil
	}

	// Perform graph traversal
	traverseOpts := mnemosyne.GraphTraverseOptions{
		SeedIDs: seedIDs,
		MaxHops: 2,
	}

	memories, _, err := client.GraphTraverse(ctx, traverseOpts)
	if err != nil {
		return nil, err
	}

	return memories, nil
}

// applySearchFilters applies client-side filters to search results.
func applySearchFilters(memories []*pb.MemoryNote, opts SearchOptions) []*pb.MemoryNote {
	filtered := make([]*pb.MemoryNote, 0, len(memories))

	for _, mem := range memories {
		// Apply max importance filter
		if opts.MaxImportance > 0 && int(mem.Importance) > opts.MaxImportance {
			continue
		}

		// Apply namespace filter (if multiple namespaces specified)
		if len(opts.Namespaces) > 1 {
			matched := false
			for _, ns := range opts.Namespaces {
				if matchesNamespace(mem.Namespace, ns) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		filtered = append(filtered, mem)
	}

	return filtered
}

// matchesNamespace checks if a memory's namespace matches the filter.
func matchesNamespace(memNS *pb.Namespace, filter string) bool {
	if memNS == nil {
		return filter == ""
	}

	nsStr := formatNamespace(memNS)
	return nsStr == filter || (len(filter) > 0 && len(nsStr) > len(filter) && nsStr[:len(filter)] == filter)
}

// SetSearchMode sets the search mode in the model.
func (m *Model) SetSearchMode(mode SearchMode) {
	if m.searchOptions.SearchMode != mode {
		m.searchOptions.SearchMode = mode
	}
}

// GetSearchMode returns the current search mode.
func (m Model) GetSearchMode() SearchMode {
	return m.searchOptions.SearchMode
}

// SetSearchFilters sets the search filters.
func (m *Model) SetSearchFilters(namespaces []string, tags []string, minImp, maxImp int) {
	m.searchOptions.Namespaces = namespaces
	m.searchOptions.Tags = tags
	m.searchOptions.MinImportance = minImp
	m.searchOptions.MaxImportance = maxImp
}

// GetSearchOptions returns the current search options.
func (m Model) GetSearchOptions() SearchOptions {
	return m.searchOptions
}

// ClearSearch clears the search query and results.
func (m *Model) ClearSearch() {
	m.searchQuery = ""
	m.searchInput = ""
	m.searchActive = false
	m.lastSearchQuery = ""
	if m.searchDebouncer != nil {
		m.searchDebouncer.Cancel()
	}
	m.applyFilters()
}

// IsSearchActive returns whether search is currently active.
func (m Model) IsSearchActive() bool {
	return m.searchActive
}

// LastSearchQuery returns the last executed search query.
func (m Model) LastSearchQuery() string {
	return m.lastSearchQuery
}
