package memorylist

import (
	"testing"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Helper function to create test memories for filter tests
func createTestMemoryWithNamespace(id string, content string, tags []string, importance uint32, ns *pb.Namespace) *pb.MemoryNote {
	return &pb.MemoryNote{
		Id:         id,
		Content:    content,
		Tags:       tags,
		Importance: importance,
		Namespace:  ns,
		CreatedAt:  0,
		UpdatedAt:  0,
	}
}

// Helper to create global namespace
func globalNamespace() *pb.Namespace {
	return &pb.Namespace{
		Namespace: &pb.Namespace_Global{
			Global: &pb.GlobalNamespace{},
		},
	}
}

// Helper to create project namespace
func projectNamespace(name string) *pb.Namespace {
	return &pb.Namespace{
		Namespace: &pb.Namespace_Project{
			Project: &pb.ProjectNamespace{Name: name},
		},
	}
}

// ===== Text Search Tests =====

// Test 1: Search with simple substring match (client-side local search)
func TestLocalSearchMemoriesSimple(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "hello world", nil, 5, nil),
		createTestMemoryWithNamespace("2", "goodbye world", nil, 5, nil),
		createTestMemoryWithNamespace("3", "hello there", nil, 5, nil),
	}

	results := SearchMemoriesLocal(memories, "hello", false, false)
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

// Test 2: Search case insensitive
func TestLocalSearchMemoriesCaseInsensitive(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "Hello World", nil, 5, nil),
		createTestMemoryWithNamespace("2", "goodbye world", nil, 5, nil),
	}

	results := SearchMemoriesLocal(memories, "hello", false, false)
	if len(results) != 1 {
		t.Errorf("Expected 1 result (case insensitive), got %d", len(results))
	}
}

// Test 3: Search case sensitive
func TestLocalSearchMemoriesCaseSensitive(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "Hello World", nil, 5, nil),
		createTestMemoryWithNamespace("2", "hello world", nil, 5, nil),
	}

	results := SearchMemoriesLocal(memories, "Hello", true, false)
	if len(results) != 1 {
		t.Errorf("Expected 1 result (case sensitive), got %d", len(results))
	}
}

// Test 4: Search with regex
func TestLocalSearchMemoriesRegex(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "memory123", nil, 5, nil),
		createTestMemoryWithNamespace("2", "memoryx", nil, 5, nil),
		createTestMemoryWithNamespace("3", "memory456", nil, 5, nil),
	}

	results := SearchMemoriesLocal(memories, "memory\\d+", false, true)
	if len(results) != 2 {
		t.Errorf("Expected 2 regex matches, got %d", len(results))
	}
}

// Test 5: Search with invalid regex
func TestLocalSearchMemoriesInvalidRegex(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "test", nil, 5, nil),
	}

	results := SearchMemoriesLocal(memories, "[invalid(", false, true)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for invalid regex, got %d", len(results))
	}
}

// Test 6: Empty search returns all
func TestLocalSearchMemoriesEmpty(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "test", nil, 5, nil),
		createTestMemoryWithNamespace("2", "test2", nil, 5, nil),
	}

	results := SearchMemoriesLocal(memories, "", false, false)
	if len(results) != len(memories) {
		t.Errorf("Expected all memories for empty search, got %d/%d", len(results), len(memories))
	}
}

// ===== Tag Filter Tests =====

// Test 7: Filter by tags with OR logic
func TestFilterByTagsOR(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", []string{"api", "bug"}, 5, nil),
		createTestMemoryWithNamespace("2", "mem2", []string{"frontend"}, 5, nil),
		createTestMemoryWithNamespace("3", "mem3", []string{"api"}, 5, nil),
	}

	results := FilterByTags(memories, []string{"api", "frontend"}, false)
	if len(results) != 3 {
		t.Errorf("Expected 3 results with OR logic (1 has api, 2 has frontend, 3 has api), got %d", len(results))
	}
}

// Test 8: Filter by tags with AND logic
func TestFilterByTagsAND(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", []string{"api", "bug"}, 5, nil),
		createTestMemoryWithNamespace("2", "mem2", []string{"api"}, 5, nil),
		createTestMemoryWithNamespace("3", "mem3", []string{"api", "bug", "urgent"}, 5, nil),
	}

	results := FilterByTags(memories, []string{"api", "bug"}, true)
	if len(results) != 2 {
		t.Errorf("Expected 2 results with AND logic, got %d", len(results))
	}

	// Check that both api and bug are present
	for _, mem := range results {
		hasApi := false
		hasBug := false
		for _, tag := range mem.Tags {
			if tag == "api" {
				hasApi = true
			}
			if tag == "bug" {
				hasBug = true
			}
		}
		if !hasApi || !hasBug {
			t.Errorf("Memory %s doesn't have both required tags", mem.Id)
		}
	}
}

// Test 9: Filter by empty tags returns all
func TestFilterByTagsEmpty(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", []string{"api"}, 5, nil),
		createTestMemoryWithNamespace("2", "mem2", []string{"frontend"}, 5, nil),
	}

	results := FilterByTags(memories, []string{}, true)
	if len(results) != len(memories) {
		t.Errorf("Expected all memories for empty tags, got %d/%d", len(results), len(memories))
	}
}

// Test 10: Filter by tags with no matches
func TestFilterByTagsNoMatches(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", []string{"api"}, 5, nil),
		createTestMemoryWithNamespace("2", "mem2", []string{"frontend"}, 5, nil),
	}

	results := FilterByTags(memories, []string{"nonexistent"}, false)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-existent tags, got %d", len(results))
	}
}

// ===== Importance Filter Tests =====

// Test 11: Filter by importance range
func TestFilterByImportanceRange(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", nil, 3, nil),
		createTestMemoryWithNamespace("2", "mem2", nil, 5, nil),
		createTestMemoryWithNamespace("3", "mem3", nil, 7, nil),
		createTestMemoryWithNamespace("4", "mem4", nil, 9, nil),
	}

	results := FilterByImportance(memories, 5, 8)
	if len(results) != 2 {
		t.Errorf("Expected 2 results in importance range [5,8], got %d", len(results))
	}

	// Verify the results
	for _, mem := range results {
		if mem.Importance < 5 || mem.Importance > 8 {
			t.Errorf("Memory %s is outside the importance range", mem.Id)
		}
	}
}

// Test 12: Filter by importance inclusive boundaries
func TestFilterByImportanceBoundaries(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", nil, 5, nil),
		createTestMemoryWithNamespace("2", "mem2", nil, 6, nil),
		createTestMemoryWithNamespace("3", "mem3", nil, 7, nil),
	}

	results := FilterByImportance(memories, 5, 7)
	if len(results) != 3 {
		t.Errorf("Expected 3 results (inclusive boundaries), got %d", len(results))
	}
}

// Test 13: Filter by importance with zero importance
func TestFilterByImportanceZero(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", nil, 0, nil),
		createTestMemoryWithNamespace("2", "mem2", nil, 5, nil),
	}

	results := FilterByImportance(memories, 0, 3)
	if len(results) != 1 {
		t.Errorf("Expected 1 result with low importance, got %d", len(results))
	}
}

// ===== Namespace Filter Tests =====

// Test 14: Filter by global namespace
func TestFilterByNamespaceGlobal(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", nil, 5, globalNamespace()),
		createTestMemoryWithNamespace("2", "mem2", nil, 5, projectNamespace("test")),
		createTestMemoryWithNamespace("3", "mem3", nil, 5, globalNamespace()),
	}

	results := FilterByNamespace(memories, "global", "")
	if len(results) != 2 {
		t.Errorf("Expected 2 global namespace memories, got %d", len(results))
	}
}

// Test 15: Filter by project namespace
func TestFilterByNamespaceProject(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", nil, 5, projectNamespace("project1")),
		createTestMemoryWithNamespace("2", "mem2", nil, 5, projectNamespace("project2")),
		createTestMemoryWithNamespace("3", "mem3", nil, 5, globalNamespace()),
	}

	results := FilterByNamespace(memories, "project", "")
	if len(results) != 2 {
		t.Errorf("Expected 2 project namespace memories, got %d", len(results))
	}
}

// Test 16: Filter by specific project name
func TestFilterByNamespaceSpecificProject(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", nil, 5, projectNamespace("project1")),
		createTestMemoryWithNamespace("2", "mem2", nil, 5, projectNamespace("project1")),
		createTestMemoryWithNamespace("3", "mem3", nil, 5, projectNamespace("project2")),
	}

	results := FilterByNamespace(memories, "project", "project1")
	if len(results) != 2 {
		t.Errorf("Expected 2 memories in project1, got %d", len(results))
	}
}

// Test 17: Filter by namespace with no filter
func TestFilterByNamespaceNoFilter(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "mem1", nil, 5, globalNamespace()),
		createTestMemoryWithNamespace("2", "mem2", nil, 5, projectNamespace("test")),
	}

	results := FilterByNamespace(memories, "", "")
	if len(results) != len(memories) {
		t.Errorf("Expected all memories with empty namespace filter, got %d/%d", len(results), len(memories))
	}
}

// ===== Combined Filter Tests =====

// Test 18: Combined search and tag filter
func TestCombinedSearchAndTags(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "hello api", []string{"important"}, 5, nil),
		createTestMemoryWithNamespace("2", "goodbye api", []string{"trivial"}, 5, nil),
		createTestMemoryWithNamespace("3", "hello world", []string{"important"}, 5, nil),
	}

	opts := FilterOptions{
		SearchQuery: "api",
		Tags:        []string{"important"},
		TagsAND:     false,
		MinImportance: 0,
		MaxImportance: 10,
	}

	results := FilterMemories(memories, opts)
	if len(results) != 1 {
		t.Errorf("Expected 1 result with search and tag filter, got %d", len(results))
	}

	if results[0].Id != "1" {
		t.Errorf("Expected memory 1, got %s", results[0].Id)
	}
}

// Test 19: Combined all filters
func TestCombinedAllFilters(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "api bug fix", []string{"api", "bug"}, 8, projectNamespace("proj1")),
		createTestMemoryWithNamespace("2", "api feature", []string{"api"}, 7, projectNamespace("proj1")),
		createTestMemoryWithNamespace("3", "database issue", []string{"db"}, 8, projectNamespace("proj2")),
	}

	opts := FilterOptions{
		SearchQuery:   "api",
		CaseSensitive: false,
		UseRegex:      false,
		Tags:          []string{"api"},
		TagsAND:       false,
		MinImportance: 7,
		MaxImportance: 10,
		NamespaceType: "project",
		NamespaceName: "proj1",
	}

	results := FilterMemories(memories, opts)
	if len(results) != 2 {
		t.Errorf("Expected 2 results with all filters, got %d", len(results))
	}
}

// Test 20: Combined filters with no matches
func TestCombinedFiltersNoMatches(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "api test", []string{"frontend"}, 5, nil),
	}

	opts := FilterOptions{
		SearchQuery: "backend",
		Tags:        []string{"api"},
		TagsAND:     false,
	}

	results := FilterMemories(memories, opts)
	if len(results) != 0 {
		t.Errorf("Expected 0 results with incompatible filters, got %d", len(results))
	}
}

// ===== Fluent Filter Builder Tests =====

// Test 21: Filter builder with search
func TestFilterBuilderSearch(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "hello world", nil, 5, nil),
		createTestMemoryWithNamespace("2", "goodbye world", nil, 5, nil),
	}

	filter := NewFilter().WithSearch("hello", false, false)
	results := filter.Apply(memories)

	if len(results) != 1 {
		t.Errorf("Expected 1 result from filter builder, got %d", len(results))
	}
}

// Test 22: Filter builder chaining
func TestFilterBuilderChaining(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "api test", []string{"important"}, 8, projectNamespace("proj1")),
		createTestMemoryWithNamespace("2", "api test", []string{"trivial"}, 8, projectNamespace("proj1")),
		createTestMemoryWithNamespace("3", "api test", []string{"important"}, 3, projectNamespace("proj1")),
	}

	filter := NewFilter().
		WithSearch("api", false, false).
		WithTags([]string{"important"}, false).
		WithImportance(5, 10).
		WithNamespace("project", "proj1")

	results := filter.Apply(memories)

	if len(results) != 1 {
		t.Errorf("Expected 1 result from chained filters, got %d", len(results))
	}

	if results[0].Id != "1" {
		t.Errorf("Expected memory 1, got %s", results[0].Id)
	}
}

// Test 23: Filter builder reset
func TestFilterBuilderReset(t *testing.T) {
	memories := []*pb.MemoryNote{
		createTestMemoryWithNamespace("1", "test", []string{"tag1"}, 5, nil),
		createTestMemoryWithNamespace("2", "test", []string{"tag2"}, 5, nil),
	}

	filter := NewFilter().WithTags([]string{"tag1"}, false)
	filter.Reset()

	results := filter.Apply(memories)

	if len(results) != len(memories) {
		t.Errorf("Expected all memories after reset, got %d/%d", len(results), len(memories))
	}
}

// Test 24: Filter builder options getter
func TestFilterBuilderOptions(t *testing.T) {
	filter := NewFilter().
		WithSearch("test", true, false).
		WithImportance(5, 8)

	opts := filter.Options()

	if opts.SearchQuery != "test" {
		t.Errorf("Expected search query 'test', got %s", opts.SearchQuery)
	}

	if !opts.CaseSensitive {
		t.Error("Expected case sensitive to be true")
	}

	if opts.MinImportance != 5 || opts.MaxImportance != 8 {
		t.Errorf("Expected importance [5,8], got [%d,%d]", opts.MinImportance, opts.MaxImportance)
	}
}

// Test 25: Large dataset performance
func TestLargeDatasetPerformance(t *testing.T) {
	// Create 1000 test memories
	memories := make([]*pb.MemoryNote, 1000)
	for i := 0; i < 1000; i++ {
		importance := uint32((i % 10) + 1)
		tags := []string{"tag1", "tag2"}
		if i%2 == 0 {
			tags = []string{"tag1", "tag3"}
		}
		memories[i] = createTestMemoryWithNamespace(string(rune(i)), "content "+string(rune(i%100)), tags, importance, nil)
	}

	opts := FilterOptions{
		SearchQuery:   "content",
		Tags:          []string{"tag1"},
		TagsAND:       false,
		MinImportance: 5,
		MaxImportance: 10,
	}

	results := FilterMemories(memories, opts)

	// Should filter to roughly 25% of dataset (importance 5-10) * 100% (has tag1)
	// With search for "content", should match all since all have "content" in them
	if len(results) == 0 {
		t.Error("Expected non-zero results for large dataset")
	}

	if len(results) > len(memories) {
		t.Errorf("Results exceed input size: %d > %d", len(results), len(memories))
	}
}
