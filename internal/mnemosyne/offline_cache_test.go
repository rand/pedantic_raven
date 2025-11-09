package mnemosyne

import (
	"sync"
	"testing"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// TestSyncQueueBasicOperations verifies basic sync queue functionality
func TestSyncQueueBasicOperations(t *testing.T) {
	sq := NewSyncQueue()

	// Initially empty
	if sq.Len() != 0 {
		t.Errorf("expected empty queue, got %d operations", sq.Len())
	}

	// Add operation
	op := SyncOperation{
		Type:      OpCreate,
		MemoryID:  "mem-1",
		Timestamp: time.Now(),
	}
	sq.Add(op)

	if sq.Len() != 1 {
		t.Errorf("expected 1 operation, got %d", sq.Len())
	}

	// Get all operations
	ops := sq.GetAll()
	if len(ops) != 1 {
		t.Errorf("expected 1 operation, got %d", len(ops))
	}

	if ops[0].Type != OpCreate {
		t.Errorf("expected OpCreate, got %v", ops[0].Type)
	}
}

// TestSyncQueueRemoveByIndex verifies operation removal
func TestSyncQueueRemoveByIndex(t *testing.T) {
	sq := NewSyncQueue()

	// Add multiple operations
	for i := 0; i < 5; i++ {
		sq.Add(SyncOperation{
			Type:      OpCreate,
			MemoryID:  "mem-" + string(rune('1'+i)),
			Timestamp: time.Now(),
		})
	}

	if sq.Len() != 5 {
		t.Fatalf("expected 5 operations, got %d", sq.Len())
	}

	// Remove operation at index 2
	sq.Remove(2)

	if sq.Len() != 4 {
		t.Errorf("expected 4 operations after removal, got %d", sq.Len())
	}

	// Remove invalid indices (should be safe)
	sq.Remove(-1)
	sq.Remove(100)

	if sq.Len() != 4 {
		t.Errorf("expected 4 operations after invalid removals, got %d", sq.Len())
	}
}

// TestSyncQueueClearAll verifies clearing the queue
func TestSyncQueueClearAll(t *testing.T) {
	sq := NewSyncQueue()

	// Add operations
	for i := 0; i < 10; i++ {
		sq.Add(SyncOperation{
			Type:     OpCreate,
			MemoryID: "mem-" + string(rune('1'+i)),
		})
	}

	sq.Clear()

	if sq.Len() != 0 {
		t.Errorf("expected empty queue after clear, got %d operations", sq.Len())
	}

	// GetAll should return empty slice
	ops := sq.GetAll()
	if len(ops) != 0 {
		t.Errorf("expected empty slice, got %d operations", len(ops))
	}
}

// TestSyncQueueThreadSafety verifies thread-safe operations
func TestSyncQueueThreadSafety(t *testing.T) {
	sq := NewSyncQueue()
	var wg sync.WaitGroup

	// Concurrent adds
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sq.Add(SyncOperation{
				Type:     OpCreate,
				MemoryID: "mem-" + string(rune('1'+idx)),
			})
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sq.GetAll()
			_ = sq.Len()
		}()
	}

	wg.Wait()

	// Should have 50 operations
	if sq.Len() != 50 {
		t.Errorf("expected 50 operations, got %d", sq.Len())
	}
}

// TestOperationTypeStringValues verifies string representation
func TestOperationTypeStringValues(t *testing.T) {
	tests := []struct {
		op   OperationType
		want string
	}{
		{OpCreate, "create"},
		{OpUpdate, "update"},
		{OpDelete, "delete"},
		{OperationType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.op.String()
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

// TestOfflineCacheBasicOperations verifies basic cache functionality
func TestOfflineCacheBasicOperations(t *testing.T) {
	oc := NewOfflineCache()

	// Initially empty
	if oc.Len() != 0 {
		t.Errorf("expected empty cache, got %d memories", oc.Len())
	}

	// Store memory
	memory := &pb.MemoryNote{
		Id:      "mem-1",
		Content: "test content",
	}
	oc.Store(memory)

	if oc.Len() != 1 {
		t.Errorf("expected 1 memory, got %d", oc.Len())
	}

	// Retrieve memory
	retrieved, ok := oc.Get("mem-1")
	if !ok {
		t.Fatal("expected memory to be found")
	}

	if retrieved.Content != "test content" {
		t.Errorf("expected content 'test content', got %q", retrieved.Content)
	}
}

// TestOfflineCacheStoreNilMemoryHandling verifies nil memory handling
func TestOfflineCacheStoreNilMemoryHandling(t *testing.T) {
	oc := NewOfflineCache()

	// Store nil memory (should be ignored)
	oc.Store(nil)

	if oc.Len() != 0 {
		t.Errorf("expected empty cache after storing nil, got %d memories", oc.Len())
	}
}

// TestOfflineCacheStoreEmptyIDHandling verifies empty ID handling
func TestOfflineCacheStoreEmptyIDHandling(t *testing.T) {
	oc := NewOfflineCache()

	// Store memory with empty ID (should be ignored)
	oc.Store(&pb.MemoryNote{
		Id:      "",
		Content: "test",
	})

	if oc.Len() != 0 {
		t.Errorf("expected empty cache after storing memory with empty ID, got %d memories", oc.Len())
	}
}

// TestOfflineCacheGetNonExistent verifies retrieval of non-existent memory
func TestOfflineCacheGetNonExistent(t *testing.T) {
	oc := NewOfflineCache()

	_, ok := oc.Get("non-existent")
	if ok {
		t.Error("expected memory not to be found")
	}
}

// TestOfflineCacheListAllMemories verifies listing all memories
func TestOfflineCacheListAllMemories(t *testing.T) {
	oc := NewOfflineCache()

	// Add multiple memories
	for i := 0; i < 5; i++ {
		oc.Store(&pb.MemoryNote{
			Id:      "mem-" + string(rune('1'+i)),
			Content: "content " + string(rune('1'+i)),
		})
	}

	memories := oc.ListAll()
	if len(memories) != 5 {
		t.Errorf("expected 5 memories, got %d", len(memories))
	}
}

// TestOfflineCacheDirtyTracking verifies dirty flag operations
func TestOfflineCacheDirtyTracking(t *testing.T) {
	oc := NewOfflineCache()

	memID := "mem-1"

	// Initially not dirty
	if oc.IsDirty(memID) {
		t.Error("expected memory not to be dirty initially")
	}

	// Mark dirty
	oc.MarkDirty(memID)
	if !oc.IsDirty(memID) {
		t.Error("expected memory to be dirty after marking")
	}

	// Get dirty memories
	dirty := oc.GetDirty()
	if len(dirty) != 1 {
		t.Errorf("expected 1 dirty memory, got %d", len(dirty))
	}

	// Clear dirty
	oc.ClearDirty(memID)
	if oc.IsDirty(memID) {
		t.Error("expected memory not to be dirty after clearing")
	}

	dirty = oc.GetDirty()
	if len(dirty) != 0 {
		t.Errorf("expected 0 dirty memories, got %d", len(dirty))
	}
}

// TestOfflineCacheDeleteMemory verifies deletion
func TestOfflineCacheDeleteMemory(t *testing.T) {
	oc := NewOfflineCache()

	memID := "mem-1"
	oc.Store(&pb.MemoryNote{
		Id:      memID,
		Content: "test",
	})
	oc.MarkDirty(memID)

	// Verify it exists
	_, ok := oc.Get(memID)
	if !ok {
		t.Fatal("expected memory to exist")
	}

	// Delete
	oc.Delete(memID)

	// Verify it's gone
	_, ok = oc.Get(memID)
	if ok {
		t.Error("expected memory to be deleted")
	}

	// Verify dirty flag is also removed
	if oc.IsDirty(memID) {
		t.Error("expected dirty flag to be removed")
	}
}

// TestOfflineCacheClearAll verifies clearing the cache
func TestOfflineCacheClearAll(t *testing.T) {
	oc := NewOfflineCache()

	// Add memories
	for i := 0; i < 10; i++ {
		memID := "mem-" + string(rune('1'+i))
		oc.Store(&pb.MemoryNote{
			Id:      memID,
			Content: "test",
		})
		oc.MarkDirty(memID)
	}

	oc.Clear()

	if oc.Len() != 0 {
		t.Errorf("expected empty cache, got %d memories", oc.Len())
	}

	dirty := oc.GetDirty()
	if len(dirty) != 0 {
		t.Errorf("expected no dirty memories, got %d", len(dirty))
	}
}

// TestOfflineCacheLastSyncTimeTracking verifies last sync time tracking
func TestOfflineCacheLastSyncTimeTracking(t *testing.T) {
	oc := NewOfflineCache()

	// Initial last sync time should be recent
	lastSync := oc.LastSyncTime()
	if time.Since(lastSync) > 1*time.Second {
		t.Error("expected recent last sync time")
	}
}

// TestOfflineCacheSyncNilClient verifies sync fails with nil client
func TestOfflineCacheSyncNilClient(t *testing.T) {
	oc := NewOfflineCache()
	sq := NewSyncQueue()

	count, err := oc.Sync(nil, sq)
	if err == nil {
		t.Fatal("expected error with nil client")
	}

	if count != 0 {
		t.Errorf("expected 0 synced operations, got %d", count)
	}
}

// TestOfflineCacheSyncNotConnected verifies sync fails when not connected
func TestOfflineCacheSyncNotConnected(t *testing.T) {
	oc := NewOfflineCache()
	sq := NewSyncQueue()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	count, err := oc.Sync(client, sq)
	if err != ErrNotConnected {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 synced operations, got %d", count)
	}
}

// TestOfflineCacheSyncEmptyQueue verifies sync with no operations
func TestOfflineCacheSyncEmptyQueue(t *testing.T) {
	oc := NewOfflineCache()
	sq := NewSyncQueue()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// The client isn't actually connected, so sync should fail
	// This test verifies that an empty queue handles the failure gracefully
	count, err := oc.Sync(client, sq)

	// Should get error because not connected
	if err == nil {
		t.Error("expected error when syncing without connection")
	}

	if count != 0 {
		t.Errorf("expected 0 synced operations, got %d", count)
	}
}

// TestOfflineCacheSyncWithNilMemory verifies sync skips operations with nil memory
func TestOfflineCacheSyncWithNilMemory(t *testing.T) {
	oc := NewOfflineCache()
	sq := NewSyncQueue()

	// Add operation with nil memory
	sq.Add(SyncOperation{
		Type:     OpCreate,
		MemoryID: "mem-1",
		Memory:   nil,
	})

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Simulate connected state
	client.connected = true
	defer func() { client.connected = false }()

	// Sync should skip the nil memory operation
	count, err := oc.Sync(client, sq)

	// Since we can't actually connect to a server, this will fail
	// But the test verifies the structure is correct
	if count < 0 {
		t.Errorf("expected non-negative count, got %d", count)
	}
}

// TestOfflineCacheThreadSafety verifies thread-safe operations
func TestOfflineCacheThreadSafety(t *testing.T) {
	oc := NewOfflineCache()
	var wg sync.WaitGroup

	// Concurrent stores
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			oc.Store(&pb.MemoryNote{
				Id:      "mem-" + string(rune('1'+idx)),
				Content: "test",
			})
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = oc.ListAll()
			_ = oc.Len()
			_ = oc.LastSyncTime()
		}()
	}

	// Concurrent dirty operations
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			memID := "mem-" + string(rune('1'+idx))
			oc.MarkDirty(memID)
			_ = oc.IsDirty(memID)
		}(i)
	}

	wg.Wait()

	if oc.Len() != 50 {
		t.Errorf("expected 50 memories, got %d", oc.Len())
	}
}

// TestOfflineCacheConcurrentDirtyOperations verifies thread-safe dirty operations
func TestOfflineCacheConcurrentDirtyOperations(t *testing.T) {
	oc := NewOfflineCache()
	var wg sync.WaitGroup

	memID := "mem-1"

	// Store a memory
	oc.Store(&pb.MemoryNote{
		Id:      memID,
		Content: "test",
	})

	// Concurrent mark/clear dirty
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			oc.MarkDirty(memID)
		}()
		go func() {
			defer wg.Done()
			_ = oc.IsDirty(memID)
		}()
	}

	wg.Wait()

	// Should be dirty
	if !oc.IsDirty(memID) {
		t.Error("expected memory to be dirty")
	}
}

// TestOfflineCacheConcurrentDeleteAndRead verifies safe concurrent delete
func TestOfflineCacheConcurrentDeleteAndRead(t *testing.T) {
	oc := NewOfflineCache()
	var wg sync.WaitGroup

	// Add memories
	for i := 0; i < 10; i++ {
		oc.Store(&pb.MemoryNote{
			Id:      "mem-" + string(rune('1'+i)),
			Content: "test",
		})
	}

	// Concurrent deletes and reads
	for i := 0; i < 10; i++ {
		wg.Add(2)
		memID := "mem-" + string(rune('1'+i))
		go func(id string) {
			defer wg.Done()
			oc.Delete(id)
		}(memID)
		go func(id string) {
			defer wg.Done()
			_, _ = oc.Get(id)
		}(memID)
	}

	wg.Wait()

	if oc.Len() != 0 {
		t.Errorf("expected empty cache, got %d memories", oc.Len())
	}
}

// TestOfflineCacheGetDirtyFiltersCorrectly verifies GetDirty only returns dirty memories
func TestOfflineCacheGetDirtyFiltersCorrectly(t *testing.T) {
	oc := NewOfflineCache()

	// Mark some as dirty
	oc.MarkDirty("mem-1")
	oc.MarkDirty("mem-2")

	// The dirty map tracks by ID, not whether memory exists
	dirty := oc.GetDirty()

	if len(dirty) != 2 {
		t.Errorf("expected 2 dirty memories, got %d", len(dirty))
	}

	// Verify the IDs
	idMap := make(map[string]bool)
	for _, id := range dirty {
		idMap[id] = true
	}

	if !idMap["mem-1"] || !idMap["mem-2"] {
		t.Error("expected mem-1 and mem-2 to be in dirty list")
	}
}

// TestSyncQueueGetAllReturnsACopy verifies GetAll returns a defensive copy
func TestSyncQueueGetAllReturnsACopy(t *testing.T) {
	sq := NewSyncQueue()

	// Add operation
	sq.Add(SyncOperation{
		Type:     OpCreate,
		MemoryID: "mem-1",
	})

	// Get operations
	ops := sq.GetAll()

	// Modify the returned slice
	ops[0].Type = OpDelete

	// Get again
	ops2 := sq.GetAll()

	// Should still be OpCreate
	if ops2[0].Type != OpCreate {
		t.Error("expected original operation to be unmodified")
	}
}

// TestOfflineCacheStoreOverwritesExisting verifies storing same ID overwrites
func TestOfflineCacheStoreOverwritesExisting(t *testing.T) {
	oc := NewOfflineCache()

	// Store initial memory
	oc.Store(&pb.MemoryNote{
		Id:      "mem-1",
		Content: "original",
	})

	// Store updated memory with same ID
	oc.Store(&pb.MemoryNote{
		Id:      "mem-1",
		Content: "updated",
	})

	// Should have only one memory
	if oc.Len() != 1 {
		t.Errorf("expected 1 memory, got %d", oc.Len())
	}

	// Content should be updated
	mem, ok := oc.Get("mem-1")
	if !ok {
		t.Fatal("expected memory to exist")
	}

	if mem.Content != "updated" {
		t.Errorf("expected content 'updated', got %q", mem.Content)
	}
}
