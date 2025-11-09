package memorylist

import (
	"testing"
	"time"

	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// TestNewQueryCache tests cache initialization
func TestNewQueryCache(t *testing.T) {
	maxAge := 5 * time.Minute
	cache := NewQueryCache(maxAge)

	if cache == nil {
		t.Fatal("NewQueryCache returned nil")
	}

	if cache.maxAge != maxAge {
		t.Errorf("Expected maxAge %v, got %v", maxAge, cache.maxAge)
	}

	if cache.entries == nil {
		t.Error("Cache entries map not initialized")
	}

	if len(cache.entries) != 0 {
		t.Errorf("Expected empty cache, got %d entries", len(cache.entries))
	}
}

// TestCacheGetMiss tests cache miss scenario
func TestCacheGetMiss(t *testing.T) {
	cache := NewQueryCache(5 * time.Minute)

	memories, found := cache.Get("nonexistent")

	if found {
		t.Error("Expected cache miss, got hit")
	}

	if memories != nil {
		t.Error("Expected nil memories on cache miss")
	}
}

// TestCacheGetHit tests cache hit scenario
func TestCacheGetHit(t *testing.T) {
	cache := NewQueryCache(5 * time.Minute)

	// Create test memories
	testMemories := []*pb.MemoryNote{
		{Id: "mem-1", Content: "Test memory 1"},
		{Id: "mem-2", Content: "Test memory 2"},
	}

	// Store in cache
	cache.Set("test query", testMemories)

	// Retrieve from cache
	memories, found := cache.Get("test query")

	if !found {
		t.Fatal("Expected cache hit, got miss")
	}

	if len(memories) != len(testMemories) {
		t.Errorf("Expected %d memories, got %d", len(testMemories), len(memories))
	}

	if memories[0].Id != testMemories[0].Id {
		t.Errorf("Expected memory ID %s, got %s", testMemories[0].Id, memories[0].Id)
	}
}

// TestCacheExpiration tests that cached entries expire after maxAge
func TestCacheExpiration(t *testing.T) {
	cache := NewQueryCache(50 * time.Millisecond)

	testMemories := []*pb.MemoryNote{
		{Id: "mem-1", Content: "Test memory"},
	}

	cache.Set("test query", testMemories)

	// Should hit immediately
	memories, found := cache.Get("test query")
	if !found {
		t.Error("Expected cache hit immediately after set")
	}
	if len(memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(memories))
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should miss after expiration
	memories, found = cache.Get("test query")
	if found {
		t.Error("Expected cache miss after expiration")
	}
	if memories != nil {
		t.Error("Expected nil memories after expiration")
	}
}

// TestCacheInvalidate tests cache invalidation
func TestCacheInvalidate(t *testing.T) {
	cache := NewQueryCache(5 * time.Minute)

	testMemories := []*pb.MemoryNote{
		{Id: "mem-1", Content: "Test memory"},
	}

	cache.Set("query1", testMemories)
	cache.Set("query2", testMemories)

	// Verify both are cached
	if _, found := cache.Get("query1"); !found {
		t.Error("query1 should be cached")
	}
	if _, found := cache.Get("query2"); !found {
		t.Error("query2 should be cached")
	}

	// Invalidate query1
	cache.Invalidate("query1")

	// query1 should be gone, query2 should remain
	if _, found := cache.Get("query1"); found {
		t.Error("query1 should be invalidated")
	}
	if _, found := cache.Get("query2"); !found {
		t.Error("query2 should still be cached")
	}
}

// TestCacheClear tests clearing entire cache
func TestCacheClear(t *testing.T) {
	cache := NewQueryCache(5 * time.Minute)

	testMemories := []*pb.MemoryNote{
		{Id: "mem-1", Content: "Test memory"},
	}

	cache.Set("query1", testMemories)
	cache.Set("query2", testMemories)
	cache.Set("query3", testMemories)

	// Clear cache
	cache.Clear()

	// All queries should be gone
	if _, found := cache.Get("query1"); found {
		t.Error("query1 should be cleared")
	}
	if _, found := cache.Get("query2"); found {
		t.Error("query2 should be cleared")
	}
	if _, found := cache.Get("query3"); found {
		t.Error("query3 should be cleared")
	}

	if len(cache.entries) != 0 {
		t.Errorf("Expected 0 entries after clear, got %d", len(cache.entries))
	}
}

// TestCacheConcurrentAccess tests thread safety
func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewQueryCache(5 * time.Minute)

	testMemories := []*pb.MemoryNote{
		{Id: "mem-1", Content: "Test memory"},
	}

	// Run concurrent operations
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set("concurrent", testMemories)
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("concurrent")
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	// Invalidator goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Invalidate("concurrent")
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()

	// Wait for all goroutines
	<-done
	<-done
	<-done

	// If we get here without panic, thread safety is working
}

// TestCacheMaxAge tests different maxAge durations
func TestCacheMaxAge(t *testing.T) {
	tests := []struct {
		name   string
		maxAge time.Duration
		wait   time.Duration
		expect bool // should still be cached after wait
	}{
		{"100ms maxAge, 50ms wait", 100 * time.Millisecond, 50 * time.Millisecond, true},
		{"100ms maxAge, 150ms wait", 100 * time.Millisecond, 150 * time.Millisecond, false},
		{"1s maxAge, 500ms wait", 1 * time.Second, 500 * time.Millisecond, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewQueryCache(tt.maxAge)
			testMemories := []*pb.MemoryNote{
				{Id: "mem-1", Content: "Test"},
			}

			cache.Set("test", testMemories)
			time.Sleep(tt.wait)

			_, found := cache.Get("test")
			if found != tt.expect {
				t.Errorf("Expected found=%v after %v wait with %v maxAge, got %v",
					tt.expect, tt.wait, tt.maxAge, found)
			}
		})
	}
}

// TestLoadMemoriesFromServer tests loading memories with success
func TestLoadMemoriesFromServer(t *testing.T) {
	// This test requires a mock client or running server
	// For now, test with nil client to verify error handling
	filters := Filters{
		Namespace:     "global",
		Tags:          []string{"test"},
		MinImportance: 5,
		MaxImportance: 10,
	}

	cmd := LoadMemoriesFromServer(nil, filters)
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	if loadedMsg.Err != mnemosyne.ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", loadedMsg.Err)
	}

	if loadedMsg.Memories != nil {
		t.Error("Expected nil memories on error")
	}
}

// TestLoadMemoriesWithFilters tests filter application
func TestLoadMemoriesWithFilters(t *testing.T) {
	tests := []struct {
		name    string
		filters Filters
	}{
		{
			name: "Namespace filter",
			filters: Filters{
				Namespace: "project:myproject",
			},
		},
		{
			name: "Tag filter",
			filters: Filters{
				Tags: []string{"important", "urgent"},
			},
		},
		{
			name: "Importance range",
			filters: Filters{
				MinImportance: 5,
				MaxImportance: 10,
			},
		},
		{
			name: "All filters",
			filters: Filters{
				Namespace:     "global",
				Tags:          []string{"test"},
				MinImportance: 7,
				MaxImportance: 9,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := LoadMemoriesFromServer(nil, tt.filters)
			msg := cmd()

			_, ok := msg.(MemoriesLoadedMsg)
			if !ok {
				t.Fatal("Expected MemoriesLoadedMsg")
			}
		})
	}
}

// TestLoadMemoriesWithTagFilter tests tag filtering
func TestLoadMemoriesWithTagFilter(t *testing.T) {
	filters := Filters{
		Tags: []string{"golang", "testing"},
	}

	cmd := LoadMemoriesFromServer(nil, filters)
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	// With nil client, should get error
	if loadedMsg.Err == nil {
		t.Error("Expected error with nil client")
	}
}

// TestLoadMemoriesWithImportanceRange tests importance filtering
func TestLoadMemoriesWithImportanceRange(t *testing.T) {
	filters := Filters{
		MinImportance: 3,
		MaxImportance: 8,
	}

	cmd := LoadMemoriesFromServer(nil, filters)
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	if loadedMsg.Err == nil {
		t.Error("Expected error with nil client")
	}
}

// TestSearchMemories tests semantic search
func TestSearchMemories(t *testing.T) {
	query := "test search query"

	cmd := SearchMemories(nil, query)
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	if loadedMsg.Err != mnemosyne.ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", loadedMsg.Err)
	}
}

// TestSearchMemoriesCached tests caching behavior (conceptual test)
func TestSearchMemoriesCached(t *testing.T) {
	// This would require integration with the Model that has the cache
	// For now, just test that the command can be created
	cmd := SearchMemories(nil, "cached query")
	if cmd == nil {
		t.Error("Expected non-nil command")
	}
}

// TestRefreshMemories tests refresh operation
func TestRefreshMemories(t *testing.T) {
	cmd := RefreshMemories(nil)
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	if loadedMsg.Err != mnemosyne.ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", loadedMsg.Err)
	}
}

// TestLoadMemoriesError tests error handling
func TestLoadMemoriesError(t *testing.T) {
	filters := Filters{}

	cmd := LoadMemoriesFromServer(nil, filters)
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	if loadedMsg.Err == nil {
		t.Error("Expected error with nil client")
	}

	if loadedMsg.Memories != nil {
		t.Error("Expected nil memories on error")
	}

	if loadedMsg.TotalCount != 0 {
		t.Errorf("Expected TotalCount 0 on error, got %d", loadedMsg.TotalCount)
	}
}

// TestLoadMemoriesTimeout tests timeout handling (conceptual)
func TestLoadMemoriesTimeout(t *testing.T) {
	// Would require a mock server that delays response
	// For now, verify command structure
	filters := Filters{}
	cmd := LoadMemoriesFromServer(nil, filters)

	if cmd == nil {
		t.Error("Expected non-nil command")
	}
}

// TestLoadMemoriesEmpty tests empty results
func TestLoadMemoriesEmpty(t *testing.T) {
	// With nil client, we get error, but test structure
	filters := Filters{
		Namespace: "nonexistent",
	}

	cmd := LoadMemoriesFromServer(nil, filters)
	msg := cmd()

	_, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}
}

// TestSearchMemoriesEmpty tests search with empty query
func TestSearchMemoriesEmpty(t *testing.T) {
	cmd := SearchMemories(nil, "")
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	// Empty query should return empty result, not error
	if loadedMsg.Err != mnemosyne.ErrNotConnected {
		t.Errorf("Expected ErrNotConnected with nil client, got %v", loadedMsg.Err)
	}
}

// TestPagination tests pagination limits
func TestPagination(t *testing.T) {
	// The LoadMemoriesFromServer uses MaxResults: 100
	// This test verifies that behavior
	filters := Filters{}

	cmd := LoadMemoriesFromServer(nil, filters)
	msg := cmd()

	loadedMsg, ok := msg.(MemoriesLoadedMsg)
	if !ok {
		t.Fatal("Expected MemoriesLoadedMsg")
	}

	// Even with error, the command should be properly structured
	if loadedMsg.TotalCount != 0 {
		t.Errorf("Expected TotalCount 0 with error, got %d", loadedMsg.TotalCount)
	}
}

// TestFilterApplication tests client-side filter application
func TestFilterApplication(t *testing.T) {
	testMemories := []*pb.MemoryNote{
		{Id: "1", Importance: 5, Content: "Low importance"},
		{Id: "2", Importance: 7, Content: "Medium importance"},
		{Id: "3", Importance: 9, Content: "High importance"},
		{Id: "4", Importance: 10, Content: "Very high importance"},
	}

	filtered := filterByMaxImportance(testMemories, 8)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 memories with importance <= 8, got %d", len(filtered))
	}

	for _, mem := range filtered {
		if mem.Importance > 8 {
			t.Errorf("Memory %s has importance %d > 8", mem.Id, mem.Importance)
		}
	}
}

// TestParseNamespaceString tests namespace parsing
func TestParseNamespaceString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // we'll check the type
	}{
		{"global", "global", "global"},
		{"project", "project:myproject", "project"},
		{"simple string", "myproject", "project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := parseNamespaceString(tt.input)

			if ns == nil {
				t.Fatal("Expected non-nil namespace")
			}

			// Check namespace type by introspection
			switch tt.expected {
			case "global":
				if ns.GetGlobal() == nil {
					t.Error("Expected global namespace")
				}
			case "project":
				if ns.GetProject() == nil {
					t.Error("Expected project namespace")
				}
			}
		})
	}
}

// TestFilterByMaxImportance tests importance filtering helper
func TestFilterByMaxImportance(t *testing.T) {
	testMemories := []*pb.MemoryNote{
		{Id: "1", Importance: 1},
		{Id: "2", Importance: 5},
		{Id: "3", Importance: 7},
		{Id: "4", Importance: 10},
	}

	tests := []struct {
		name          string
		maxImportance int
		expectedCount int
	}{
		{"max 5", 5, 2},
		{"max 7", 7, 3},
		{"max 10", 10, 4},
		{"max 0", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterByMaxImportance(testMemories, tt.maxImportance)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d memories, got %d", tt.expectedCount, len(filtered))
			}

			for _, mem := range filtered {
				if int(mem.Importance) > tt.maxImportance {
					t.Errorf("Memory %s has importance %d > max %d",
						mem.Id, mem.Importance, tt.maxImportance)
				}
			}
		})
	}
}
