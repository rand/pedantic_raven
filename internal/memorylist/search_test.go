package memorylist

import (
	"context"
	"testing"
	"time"

	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// mockClient implements a mock mnemosyne client for testing.
type mockClient struct {
	connected     bool
	recallResults []*pb.SearchResult
	recallErr     error
	graphMemories []*pb.MemoryNote
	graphEdges    []*pb.GraphEdge
	graphErr      error
}

func (m *mockClient) IsConnected() bool {
	return m.connected
}

func (m *mockClient) Recall(ctx context.Context, opts mnemosyne.RecallOptions) ([]*pb.SearchResult, error) {
	if m.recallErr != nil {
		return nil, m.recallErr
	}
	return m.recallResults, nil
}

func (m *mockClient) GraphTraverse(ctx context.Context, opts mnemosyne.GraphTraverseOptions) ([]*pb.MemoryNote, []*pb.GraphEdge, error) {
	if m.graphErr != nil {
		return nil, nil, m.graphErr
	}
	return m.graphMemories, m.graphEdges, nil
}

// Test 1: SearchHybrid
func TestSearchHybrid(t *testing.T) {
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-1",
				Content:    "Test memory 1",
				Importance: 8,
			},
			Score: 0.95,
		},
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-2",
				Content:    "Test memory 2",
				Importance: 7,
			},
			Score: 0.85,
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:      "test query",
		SearchMode: SearchHybrid,
		MaxResults: 10,
	}

	// Execute search
	memories, err := searchHybrid(context.Background(), client, opts)

	if err != nil {
		t.Fatalf("searchHybrid failed: %v", err)
	}

	if len(memories) != 2 {
		t.Errorf("Expected 2 memories, got %d", len(memories))
	}

	if memories[0].Id != "mem-1" {
		t.Errorf("Expected first memory ID to be 'mem-1', got %s", memories[0].Id)
	}
}

// Test 2: SearchSemantic
func TestSearchSemantic(t *testing.T) {
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-semantic",
				Content:    "Semantic search result",
				Importance: 9,
			},
			Score: 0.92,
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:      "semantic query",
		SearchMode: SearchSemantic,
		MaxResults: 10,
	}

	memories, err := searchSemantic(context.Background(), client, opts)

	if err != nil {
		t.Fatalf("searchSemantic failed: %v", err)
	}

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(memories))
	}
}

// Test 3: SearchFullText
func TestSearchFullText(t *testing.T) {
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-fts",
				Content:    "Full text search result",
				Importance: 6,
			},
			Score: 0.88,
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:      "full text query",
		SearchMode: SearchFullText,
		MaxResults: 10,
	}

	memories, err := searchFullText(context.Background(), client, opts)

	if err != nil {
		t.Fatalf("searchFullText failed: %v", err)
	}

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(memories))
	}
}

// Test 4: SearchGraph
func TestSearchGraph(t *testing.T) {
	// Mock seed search results
	seedResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:      "seed-1",
				Content: "Seed memory",
			},
		},
	}

	// Mock graph traverse results
	graphMemories := []*pb.MemoryNote{
		{
			Id:         "mem-graph-1",
			Content:    "Graph memory 1",
			Importance: 7,
		},
		{
			Id:         "mem-graph-2",
			Content:    "Graph memory 2",
			Importance: 8,
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: seedResults,
		graphMemories: graphMemories,
	}

	opts := SearchOptions{
		Query:      "graph query",
		SearchMode: SearchGraph,
		MaxResults: 10,
	}

	memories, err := searchGraph(context.Background(), client, opts)

	if err != nil {
		t.Fatalf("searchGraph failed: %v", err)
	}

	if len(memories) != 2 {
		t.Errorf("Expected 2 memories, got %d", len(memories))
	}
}

// Test 5: SearchWithNamespaceFilter
func TestSearchWithNamespaceFilter(t *testing.T) {
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:      "mem-1",
				Content: "Memory in project namespace",
				Namespace: &pb.Namespace{
					Namespace: &pb.Namespace_Project{
						Project: &pb.ProjectNamespace{Name: "test-project"},
					},
				},
				Importance: 8,
			},
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:      "test",
		Namespaces: []string{"project:test-project"},
		SearchMode: SearchHybrid,
		MaxResults: 10,
	}

	memories, err := searchHybrid(context.Background(), client, opts)

	if err != nil {
		t.Fatalf("searchHybrid with namespace filter failed: %v", err)
	}

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(memories))
	}
}

// Test 6: SearchWithTagFilter
func TestSearchWithTagFilter(t *testing.T) {
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-1",
				Content:    "Tagged memory",
				Tags:       []string{"important", "work"},
				Importance: 8,
			},
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:      "test",
		Tags:       []string{"important"},
		SearchMode: SearchHybrid,
		MaxResults: 10,
	}

	memories, err := searchHybrid(context.Background(), client, opts)

	if err != nil {
		t.Fatalf("searchHybrid with tag filter failed: %v", err)
	}

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(memories))
	}
}

// Test 7: SearchWithImportanceRange
func TestSearchWithImportanceRange(t *testing.T) {
	// Server returns all results, but min importance filter is applied server-side
	// Only high importance should be returned by server (importance >= 5)
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-1",
				Content:    "High importance",
				Importance: 9,
			},
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:         "test",
		MinImportance: 5,
		MaxImportance: 10,
		SearchMode:    SearchHybrid,
		MaxResults:    10,
	}

	memories, err := searchHybrid(context.Background(), client, opts)
	if err != nil {
		t.Fatalf("searchHybrid failed: %v", err)
	}

	// Apply client-side max importance filter
	memories = applySearchFilters(memories, opts)

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory after importance filter, got %d", len(memories))
	}

	if memories[0].Importance != 9 {
		t.Errorf("Expected importance 9, got %d", memories[0].Importance)
	}
}

// Test 8: SearchWithCombinedFilters
func TestSearchWithCombinedFilters(t *testing.T) {
	// Server applies namespace, tags, and min importance filters
	// Only mem-1 should be returned by server
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-1",
				Content:    "Filtered memory",
				Tags:       []string{"work"},
				Importance: 8,
				Namespace: &pb.Namespace{
					Namespace: &pb.Namespace_Project{
						Project: &pb.ProjectNamespace{Name: "project1"},
					},
				},
			},
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:         "test",
		Namespaces:    []string{"project:project1"},
		Tags:          []string{"work"},
		MinImportance: 7,
		MaxImportance: 10,
		SearchMode:    SearchHybrid,
		MaxResults:    10,
	}

	memories, err := searchHybrid(context.Background(), client, opts)
	if err != nil {
		t.Fatalf("searchHybrid failed: %v", err)
	}

	// Apply client-side filters
	memories = applySearchFilters(memories, opts)

	if len(memories) != 1 {
		t.Errorf("Expected 1 memory after combined filters, got %d", len(memories))
	}
}

// Test 9: SearchDebouncing
func TestSearchDebouncing(t *testing.T) {
	debouncer := NewSearchDebouncer(100 * time.Millisecond)

	executed := false
	debouncer.Debounce(func() {
		executed = true
	})

	// Should not execute immediately
	if executed {
		t.Error("Search executed immediately, expected debounce")
	}

	// Wait for debounce period
	time.Sleep(150 * time.Millisecond)

	if !executed {
		t.Error("Search not executed after debounce period")
	}
}

// Test 10: SearchEmptyQuery
func TestSearchEmptyQuery(t *testing.T) {
	client := &mockClient{
		connected: true,
	}

	opts := SearchOptions{
		Query:      "",
		SearchMode: SearchHybrid,
		MaxResults: 10,
	}

	cmd := searchWithClient(client, opts)
	msg := cmd()

	searchMsg, ok := msg.(SearchResultsMsg)
	if !ok {
		t.Fatal("Expected SearchResultsMsg")
	}

	if searchMsg.Err != nil {
		t.Errorf("Expected no error for empty query, got %v", searchMsg.Err)
	}

	if len(searchMsg.Results) != 0 {
		t.Errorf("Expected 0 results for empty query, got %d", len(searchMsg.Results))
	}
}

// Test 11: SearchResultRanking
func TestSearchResultRanking(t *testing.T) {
	mockResults := []*pb.SearchResult{
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-1",
				Content:    "First result",
				Importance: 8,
			},
			Score: 0.95,
		},
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-2",
				Content:    "Second result",
				Importance: 7,
			},
			Score: 0.85,
		},
		{
			Memory: &pb.MemoryNote{
				Id:         "mem-3",
				Content:    "Third result",
				Importance: 9,
			},
			Score: 0.75,
		},
	}

	client := &mockClient{
		connected:     true,
		recallResults: mockResults,
	}

	opts := SearchOptions{
		Query:      "test",
		SearchMode: SearchHybrid,
		MaxResults: 10,
	}

	memories, err := searchHybrid(context.Background(), client, opts)
	if err != nil {
		t.Fatalf("searchHybrid failed: %v", err)
	}

	// Results should be in the order returned by server (by score)
	if len(memories) != 3 {
		t.Errorf("Expected 3 memories, got %d", len(memories))
	}

	if memories[0].Id != "mem-1" {
		t.Errorf("Expected first result to be mem-1, got %s", memories[0].Id)
	}
}

// Test 12: SearchHistory
func TestSearchHistory(t *testing.T) {
	history := NewSearchHistory(5)

	// Add queries
	history.Add("query1")
	history.Add("query2")
	history.Add("query3")

	queries := history.Get()
	if len(queries) != 3 {
		t.Errorf("Expected 3 queries, got %d", len(queries))
	}

	if queries[0] != "query3" {
		t.Errorf("Expected most recent query to be 'query3', got %s", queries[0])
	}

	// Add duplicate - should move to front
	history.Add("query1")
	queries = history.Get()
	if queries[0] != "query1" {
		t.Errorf("Expected duplicate query to move to front, got %s", queries[0])
	}

	if len(queries) != 3 {
		t.Errorf("Expected 3 queries after duplicate, got %d", len(queries))
	}
}

// Test 13: ClearSearch
func TestClearSearch(t *testing.T) {
	m := NewModel()
	m.searchQuery = "test query"
	m.searchInput = "test input"
	m.searchActive = true
	m.lastSearchQuery = "last query"

	m.ClearSearch()

	if m.searchQuery != "" {
		t.Errorf("Expected empty search query, got %s", m.searchQuery)
	}

	if m.searchInput != "" {
		t.Errorf("Expected empty search input, got %s", m.searchInput)
	}

	if m.searchActive {
		t.Error("Expected search to be inactive")
	}

	if m.lastSearchQuery != "" {
		t.Errorf("Expected empty last search query, got %s", m.lastSearchQuery)
	}
}

// Test 14: SearchError
func TestSearchError(t *testing.T) {
	client := &mockClient{
		connected: true,
		recallErr: mnemosyne.ErrNotConnected,
	}

	opts := SearchOptions{
		Query:      "test",
		SearchMode: SearchHybrid,
		MaxResults: 10,
	}

	cmd := searchWithClient(client, opts)
	msg := cmd()

	searchMsg, ok := msg.(SearchResultsMsg)
	if !ok {
		t.Fatal("Expected SearchResultsMsg")
	}

	if searchMsg.Err == nil {
		t.Error("Expected error, got nil")
	}

	if searchMsg.Err != mnemosyne.ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", searchMsg.Err)
	}
}

// Test 15: SearchTimeout
func TestSearchTimeout(t *testing.T) {
	// This test simulates a timeout scenario
	// In practice, the timeout is handled by the context
	client := &mockClient{
		connected: false, // Simulate disconnected state
	}

	opts := SearchOptions{
		Query:      "test",
		SearchMode: SearchHybrid,
		MaxResults: 10,
	}

	cmd := searchWithClient(client, opts)
	msg := cmd()

	searchMsg, ok := msg.(SearchResultsMsg)
	if !ok {
		t.Fatal("Expected SearchResultsMsg")
	}

	if searchMsg.Err == nil {
		t.Error("Expected error for disconnected client, got nil")
	}
}

// Test 16: SearchModeString
func TestSearchModeString(t *testing.T) {
	tests := []struct {
		mode     SearchMode
		expected string
	}{
		{SearchHybrid, "Hybrid"},
		{SearchSemantic, "Semantic"},
		{SearchFullText, "Full-Text"},
		{SearchGraph, "Graph"},
	}

	for _, tt := range tests {
		result := tt.mode.String()
		if result != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, result)
		}
	}
}

// Test 17: SetSearchMode
func TestSetSearchMode(t *testing.T) {
	m := NewModel()

	// Initial mode should be Hybrid
	if m.GetSearchMode() != SearchHybrid {
		t.Errorf("Expected initial mode to be Hybrid, got %s", m.GetSearchMode())
	}

	// Set to Semantic
	m.SetSearchMode(SearchSemantic)
	if m.GetSearchMode() != SearchSemantic {
		t.Errorf("Expected mode to be Semantic, got %s", m.GetSearchMode())
	}

	// Set to FullText
	m.SetSearchMode(SearchFullText)
	if m.GetSearchMode() != SearchFullText {
		t.Errorf("Expected mode to be Full-Text, got %s", m.GetSearchMode())
	}
}

// Test 18: SetSearchFilters
func TestSetSearchFilters(t *testing.T) {
	m := NewModel()

	namespaces := []string{"project:test"}
	tags := []string{"tag1", "tag2"}
	minImp := 5
	maxImp := 9

	m.SetSearchFilters(namespaces, tags, minImp, maxImp)

	opts := m.GetSearchOptions()

	if len(opts.Namespaces) != 1 || opts.Namespaces[0] != "project:test" {
		t.Errorf("Expected namespaces to be set correctly, got %v", opts.Namespaces)
	}

	if len(opts.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(opts.Tags))
	}

	if opts.MinImportance != 5 {
		t.Errorf("Expected MinImportance 5, got %d", opts.MinImportance)
	}

	if opts.MaxImportance != 9 {
		t.Errorf("Expected MaxImportance 9, got %d", opts.MaxImportance)
	}
}

// Test 19: CycleSearchMode
func TestCycleSearchMode(t *testing.T) {
	m := NewModel()

	// Start with Hybrid
	if m.GetSearchMode() != SearchHybrid {
		t.Errorf("Expected initial mode to be Hybrid, got %s", m.GetSearchMode())
	}

	// Cycle through modes
	m.cycleSearchMode()
	if m.GetSearchMode() != SearchSemantic {
		t.Errorf("Expected Semantic, got %s", m.GetSearchMode())
	}

	m.cycleSearchMode()
	if m.GetSearchMode() != SearchFullText {
		t.Errorf("Expected Full-Text, got %s", m.GetSearchMode())
	}

	m.cycleSearchMode()
	if m.GetSearchMode() != SearchGraph {
		t.Errorf("Expected Graph, got %s", m.GetSearchMode())
	}

	m.cycleSearchMode()
	if m.GetSearchMode() != SearchHybrid {
		t.Errorf("Expected Hybrid (cycled back), got %s", m.GetSearchMode())
	}
}

// Test 20: SearchHistoryClear
func TestSearchHistoryClear(t *testing.T) {
	history := NewSearchHistory(5)

	history.Add("query1")
	history.Add("query2")
	history.Add("query3")

	if len(history.Get()) != 3 {
		t.Errorf("Expected 3 queries before clear, got %d", len(history.Get()))
	}

	history.Clear()

	if len(history.Get()) != 0 {
		t.Errorf("Expected 0 queries after clear, got %d", len(history.Get()))
	}
}

// Test 21: SearchHistoryMaxSize
func TestSearchHistoryMaxSize(t *testing.T) {
	history := NewSearchHistory(3)

	history.Add("query1")
	history.Add("query2")
	history.Add("query3")
	history.Add("query4") // Should evict query1

	queries := history.Get()
	if len(queries) != 3 {
		t.Errorf("Expected 3 queries (max size), got %d", len(queries))
	}

	// Should have query4, query3, query2 (in that order)
	if queries[0] != "query4" {
		t.Errorf("Expected first query to be 'query4', got %s", queries[0])
	}

	// query1 should not be in history
	for _, q := range queries {
		if q == "query1" {
			t.Error("query1 should have been evicted")
		}
	}
}

// Test 22: DebouncerCancel
func TestDebouncerCancel(t *testing.T) {
	debouncer := NewSearchDebouncer(100 * time.Millisecond)

	executed := false
	debouncer.Debounce(func() {
		executed = true
	})

	// Cancel before execution
	debouncer.Cancel()

	// Wait beyond debounce period
	time.Sleep(150 * time.Millisecond)

	if executed {
		t.Error("Search executed after cancel")
	}
}

// Test 23: MultipleDebounceCalls
func TestMultipleDebounceCalls(t *testing.T) {
	debouncer := NewSearchDebouncer(100 * time.Millisecond)

	count := 0

	// Call multiple times - only last should execute
	for i := 0; i < 5; i++ {
		idx := i
		debouncer.Debounce(func() {
			count = idx
		})
		time.Sleep(20 * time.Millisecond) // Less than debounce period
	}

	// Wait for final execution
	time.Sleep(150 * time.Millisecond)

	if count != 4 {
		t.Errorf("Expected count to be 4 (last call), got %d", count)
	}
}
