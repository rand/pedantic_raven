package mnemosyne

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestRecallWithNamespace tests namespace filtering in recall operations.
func TestRecallWithNamespace(t *testing.T) {
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer server.Stop()

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    5 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	ctx := context.Background()

	// Store memories in different namespaces
	globalMem, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "global memory content",
		Namespace: GlobalNamespace(),
		Tags:      []string{"global"},
	})
	if err != nil {
		t.Fatalf("StoreMemory (global) failed: %v", err)
	}

	project1Mem, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "project1 memory content",
		Namespace: ProjectNamespace("project1"),
		Tags:      []string{"project1"},
	})
	if err != nil {
		t.Fatalf("StoreMemory (project1) failed: %v", err)
	}

	project2Mem, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "project2 memory content",
		Namespace: ProjectNamespace("project2"),
		Tags:      []string{"project2"},
	})
	if err != nil {
		t.Fatalf("StoreMemory (project2) failed: %v", err)
	}

	// Recall from global namespace
	globalResults, err := client.Recall(ctx, RecallOptions{
		Query:     "memory",
		Namespace: GlobalNamespace(),
	})
	if err != nil {
		t.Fatalf("Recall (global) failed: %v", err)
	}

	// Should only return global memory
	if len(globalResults) != 1 {
		t.Errorf("Expected 1 global result, got %d", len(globalResults))
	} else if globalResults[0].Memory.Id != globalMem.Id {
		t.Errorf("Expected global memory, got %s", globalResults[0].Memory.Id)
	}

	// Recall from project1 namespace
	project1Results, err := client.Recall(ctx, RecallOptions{
		Query:     "memory",
		Namespace: ProjectNamespace("project1"),
	})
	if err != nil {
		t.Fatalf("Recall (project1) failed: %v", err)
	}

	// Should only return project1 memory
	if len(project1Results) != 1 {
		t.Errorf("Expected 1 project1 result, got %d", len(project1Results))
	} else if project1Results[0].Memory.Id != project1Mem.Id {
		t.Errorf("Expected project1 memory, got %s", project1Results[0].Memory.Id)
	}

	// Recall from project2 namespace
	project2Results, err := client.Recall(ctx, RecallOptions{
		Query:     "memory",
		Namespace: ProjectNamespace("project2"),
	})
	if err != nil {
		t.Fatalf("Recall (project2) failed: %v", err)
	}

	// Should only return project2 memory
	if len(project2Results) != 1 {
		t.Errorf("Expected 1 project2 result, got %d", len(project2Results))
	} else if project2Results[0].Memory.Id != project2Mem.Id {
		t.Errorf("Expected project2 memory, got %s", project2Results[0].Memory.Id)
	}
}

// TestRecallPagination tests cursor-based pagination in recall.
func TestRecallPagination(t *testing.T) {
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer server.Stop()

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    5 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	ctx := context.Background()

	// Store multiple memories
	numMemories := 15
	for i := 0; i < numMemories; i++ {
		_, err := client.StoreMemory(ctx, StoreMemoryOptions{
			Content:   fmt.Sprintf("test memory %d", i),
			Namespace: GlobalNamespace(),
		})
		if err != nil {
			t.Fatalf("StoreMemory %d failed: %v", i, err)
		}
	}

	// Test pagination with MaxResults
	page1, err := client.Recall(ctx, RecallOptions{
		Query:      "test",
		Namespace:  GlobalNamespace(),
		MaxResults: 5,
	})
	if err != nil {
		t.Fatalf("Recall page1 failed: %v", err)
	}

	if len(page1) != 5 {
		t.Errorf("Expected 5 results in page1, got %d", len(page1))
	}

	// Test with different page size
	page2, err := client.Recall(ctx, RecallOptions{
		Query:      "test",
		Namespace:  GlobalNamespace(),
		MaxResults: 10,
	})
	if err != nil {
		t.Fatalf("Recall page2 failed: %v", err)
	}

	if len(page2) != 10 {
		t.Errorf("Expected 10 results in page2, got %d", len(page2))
	}

	// Test default pagination (should be 10)
	defaultPage, err := client.Recall(ctx, RecallOptions{
		Query:     "test",
		Namespace: GlobalNamespace(),
	})
	if err != nil {
		t.Fatalf("Recall default page failed: %v", err)
	}

	if len(defaultPage) != 10 {
		t.Errorf("Expected 10 results with default pagination, got %d", len(defaultPage))
	}
}

// TestRecallEmptyResults tests recall when no results are found.
func TestRecallEmptyResults(t *testing.T) {
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer server.Stop()

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    5 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	ctx := context.Background()

	// Store a memory in global namespace
	_, err = client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "test content",
		Namespace: GlobalNamespace(),
	})
	if err != nil {
		t.Fatalf("StoreMemory failed: %v", err)
	}

	// Query for something that doesn't exist in a different namespace
	results, err := client.Recall(ctx, RecallOptions{
		Query:     "nonexistent query",
		Namespace: ProjectNamespace("empty-project"),
	})
	if err != nil {
		t.Fatalf("Recall failed: %v", err)
	}

	// Should return empty results (nil or empty slice), not error
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}

	// Query with tag filter that doesn't match
	results2, err := client.Recall(ctx, RecallOptions{
		Query:     "test",
		Namespace: GlobalNamespace(),
		Tags:      []string{"nonexistent-tag"},
	})
	if err != nil {
		t.Fatalf("Recall with tag filter failed: %v", err)
	}

	if len(results2) != 0 {
		t.Errorf("Expected 0 results with tag filter, got %d", len(results2))
	}

	// Query with high importance filter
	minImportance := uint32(9)
	results3, err := client.Recall(ctx, RecallOptions{
		Query:         "test",
		Namespace:     GlobalNamespace(),
		MinImportance: &minImportance,
	})
	if err != nil {
		t.Fatalf("Recall with importance filter failed: %v", err)
	}

	if len(results3) != 0 {
		t.Errorf("Expected 0 results with high importance filter, got %d", len(results3))
	}
}
