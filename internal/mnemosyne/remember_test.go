package mnemosyne

import (
	"context"
	"testing"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// TestRememberValidation tests comprehensive input validation for StoreMemory.
func TestRememberValidation(t *testing.T) {
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

	// Test 1: Missing content (should fail)
	_, err = client.StoreMemory(ctx, StoreMemoryOptions{
		Namespace: GlobalNamespace(),
	})
	if err == nil {
		t.Error("Expected error for missing content")
	}
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error, got: %v", err)
	}

	// Test 2: Missing namespace (should fail)
	_, err = client.StoreMemory(ctx, StoreMemoryOptions{
		Content: "test content",
	})
	if err == nil {
		t.Error("Expected error for missing namespace")
	}
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error, got: %v", err)
	}

	// Test 3: Empty content string (should fail)
	_, err = client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "",
		Namespace: GlobalNamespace(),
	})
	if err == nil {
		t.Error("Expected error for empty content")
	}
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error, got: %v", err)
	}

	// Test 4: Valid with minimal options (should succeed)
	memory1, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "minimal valid memory",
		Namespace: GlobalNamespace(),
	})
	if err != nil {
		t.Fatalf("StoreMemory with minimal options failed: %v", err)
	}
	if memory1 == nil {
		t.Fatal("Expected non-nil memory")
	}
	if memory1.Content != "minimal valid memory" {
		t.Errorf("Expected content 'minimal valid memory', got '%s'", memory1.Content)
	}

	// Test 5: Valid with all options (should succeed)
	importance := uint32(8)
	memoryType := pb.MemoryType_INSIGHT
	memory2, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:           "fully specified memory",
		Namespace:         ProjectNamespace("test-project"),
		Importance:        &importance,
		Context:           "test context",
		Tags:              []string{"tag1", "tag2", "tag3"},
		MemoryType:        &memoryType,
		SkipLLMEnrichment: true,
	})
	if err != nil {
		t.Fatalf("StoreMemory with all options failed: %v", err)
	}
	if memory2 == nil {
		t.Fatal("Expected non-nil memory")
	}
	if memory2.Content != "fully specified memory" {
		t.Errorf("Expected content 'fully specified memory', got '%s'", memory2.Content)
	}
	if memory2.Importance != 8 {
		t.Errorf("Expected importance 8, got %d", memory2.Importance)
	}
	if len(memory2.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(memory2.Tags))
	}

	// Test 6: Valid with different namespace types
	sessionMemory, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "session memory",
		Namespace: SessionNamespace("test-project", "session-123"),
	})
	if err != nil {
		t.Fatalf("StoreMemory with session namespace failed: %v", err)
	}
	if sessionMemory == nil {
		t.Fatal("Expected non-nil session memory")
	}

	// Test 7: Verify namespace was set correctly
	sessionNs := sessionMemory.Namespace.GetSession()
	if sessionNs == nil {
		t.Fatal("Expected session namespace")
	}
	if sessionNs.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", sessionNs.Project)
	}
	if sessionNs.SessionId != "session-123" {
		t.Errorf("Expected session ID 'session-123', got '%s'", sessionNs.SessionId)
	}
}

// TestRememberDuplicateHandling tests handling of duplicate memory storage.
func TestRememberDuplicateHandling(t *testing.T) {
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

	// Store first memory
	memory1, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "duplicate test content",
		Namespace: GlobalNamespace(),
		Tags:      []string{"duplicate", "test"},
	})
	if err != nil {
		t.Fatalf("StoreMemory (first) failed: %v", err)
	}

	// Store identical content (should create new memory, not error)
	memory2, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "duplicate test content",
		Namespace: GlobalNamespace(),
		Tags:      []string{"duplicate", "test"},
	})
	if err != nil {
		t.Fatalf("StoreMemory (duplicate) failed: %v", err)
	}

	// Should create different memory IDs
	if memory1.Id == memory2.Id {
		t.Error("Expected different memory IDs for duplicate content")
	}

	// Both should have same content
	if memory1.Content != memory2.Content {
		t.Error("Expected same content for both memories")
	}

	// Store similar but not identical content
	memory3, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "duplicate test content with variation",
		Namespace: GlobalNamespace(),
		Tags:      []string{"duplicate", "test"},
	})
	if err != nil {
		t.Fatalf("StoreMemory (similar) failed: %v", err)
	}

	// Should create different memory ID
	if memory1.Id == memory3.Id || memory2.Id == memory3.Id {
		t.Error("Expected different memory ID for similar content")
	}

	// Store same content in different namespace
	memory4, err := client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "duplicate test content",
		Namespace: ProjectNamespace("different-project"),
		Tags:      []string{"duplicate", "test"},
	})
	if err != nil {
		t.Fatalf("StoreMemory (different namespace) failed: %v", err)
	}

	// Should create different memory ID
	if memory1.Id == memory4.Id {
		t.Error("Expected different memory ID for different namespace")
	}

	// Verify all memories exist
	_, err = client.GetMemory(ctx, memory1.Id)
	if err != nil {
		t.Errorf("Failed to retrieve memory1: %v", err)
	}

	_, err = client.GetMemory(ctx, memory2.Id)
	if err != nil {
		t.Errorf("Failed to retrieve memory2: %v", err)
	}

	_, err = client.GetMemory(ctx, memory3.Id)
	if err != nil {
		t.Errorf("Failed to retrieve memory3: %v", err)
	}

	_, err = client.GetMemory(ctx, memory4.Id)
	if err != nil {
		t.Errorf("Failed to retrieve memory4: %v", err)
	}

	// List memories in global namespace - should have 3 (memory1, memory2, memory3)
	memories, err := client.ListMemories(ctx, ListMemoriesOptions{
		Namespace: GlobalNamespace(),
	})
	if err != nil {
		t.Fatalf("ListMemories failed: %v", err)
	}

	if len(memories) != 3 {
		t.Errorf("Expected 3 memories in global namespace, got %d", len(memories))
	}
}
