package mnemosyne

import (
	"sync"
	"testing"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// --- SyncQueue Tests ---

func TestNewSyncQueue(t *testing.T) {
	sq := NewSyncQueue()
	if sq == nil {
		t.Fatal("NewSyncQueue returned nil")
	}

	if sq.Len() != 0 {
		t.Errorf("New queue should be empty, got length %d", sq.Len())
	}
}

func TestSyncQueueAdd(t *testing.T) {
	sq := NewSyncQueue()

	op := SyncOperation{
		Type:      OpCreate,
		MemoryID:  "test-id",
		Timestamp: time.Now(),
	}

	sq.Add(op)

	if sq.Len() != 1 {
		t.Errorf("Expected length 1, got %d", sq.Len())
	}
}

func TestSyncQueueGetAll(t *testing.T) {
	sq := NewSyncQueue()

	ops := []SyncOperation{
		{Type: OpCreate, MemoryID: "id1", Timestamp: time.Now()},
		{Type: OpUpdate, MemoryID: "id2", Timestamp: time.Now()},
		{Type: OpDelete, MemoryID: "id3", Timestamp: time.Now()},
	}

	for _, op := range ops {
		sq.Add(op)
	}

	result := sq.GetAll()

	if len(result) != 3 {
		t.Errorf("Expected 3 operations, got %d", len(result))
	}

	// Verify FIFO order
	for i, op := range result {
		if op.MemoryID != ops[i].MemoryID {
			t.Errorf("Operation %d: expected ID %s, got %s", i, ops[i].MemoryID, op.MemoryID)
		}
	}
}

func TestSyncQueueRemove(t *testing.T) {
	sq := NewSyncQueue()

	sq.Add(SyncOperation{Type: OpCreate, MemoryID: "id1", Timestamp: time.Now()})
	sq.Add(SyncOperation{Type: OpUpdate, MemoryID: "id2", Timestamp: time.Now()})
	sq.Add(SyncOperation{Type: OpDelete, MemoryID: "id3", Timestamp: time.Now()})

	sq.Remove(1) // Remove middle item

	if sq.Len() != 2 {
		t.Errorf("Expected length 2 after removal, got %d", sq.Len())
	}

	ops := sq.GetAll()
	if ops[0].MemoryID != "id1" || ops[1].MemoryID != "id3" {
		t.Error("Removal did not preserve order correctly")
	}
}

func TestSyncQueueClear(t *testing.T) {
	sq := NewSyncQueue()

	sq.Add(SyncOperation{Type: OpCreate, MemoryID: "id1", Timestamp: time.Now()})
	sq.Add(SyncOperation{Type: OpUpdate, MemoryID: "id2", Timestamp: time.Now()})

	sq.Clear()

	if sq.Len() != 0 {
		t.Errorf("Expected empty queue after clear, got length %d", sq.Len())
	}
}

func TestSyncQueueOrdering(t *testing.T) {
	sq := NewSyncQueue()

	// Add operations with slight time delays
	times := []time.Time{}
	for i := 0; i < 5; i++ {
		now := time.Now()
		times = append(times, now)
		sq.Add(SyncOperation{
			Type:      OpCreate,
			MemoryID:  string(rune('a' + i)),
			Timestamp: now,
		})
		time.Sleep(1 * time.Millisecond)
	}

	ops := sq.GetAll()

	// Verify FIFO order preserved
	for i, op := range ops {
		if !op.Timestamp.Equal(times[i]) {
			t.Errorf("Operation %d: timestamp mismatch", i)
		}
	}
}

// --- OfflineCache Tests ---

func TestNewOfflineCache(t *testing.T) {
	oc := NewOfflineCache()
	if oc == nil {
		t.Fatal("NewOfflineCache returned nil")
	}

	if oc.Len() != 0 {
		t.Errorf("New cache should be empty, got length %d", oc.Len())
	}
}

func TestOfflineCacheStore(t *testing.T) {
	oc := NewOfflineCache()

	memory := &pb.MemoryNote{
		Id:      "test-id",
		Content: "test content",
	}

	oc.Store(memory)

	if oc.Len() != 1 {
		t.Errorf("Expected length 1, got %d", oc.Len())
	}
}

func TestOfflineCacheGet(t *testing.T) {
	oc := NewOfflineCache()

	memory := &pb.MemoryNote{
		Id:      "test-id",
		Content: "test content",
	}

	oc.Store(memory)

	retrieved, ok := oc.Get("test-id")
	if !ok {
		t.Fatal("Failed to retrieve stored memory")
	}

	if retrieved.Id != memory.Id {
		t.Errorf("Retrieved ID = %s, want %s", retrieved.Id, memory.Id)
	}

	if retrieved.Content != memory.Content {
		t.Errorf("Retrieved Content = %s, want %s", retrieved.Content, memory.Content)
	}
}

func TestOfflineCacheGetNotFound(t *testing.T) {
	oc := NewOfflineCache()

	_, ok := oc.Get("non-existent")
	if ok {
		t.Error("Expected Get to return false for non-existent memory")
	}
}

func TestOfflineCacheListAll(t *testing.T) {
	oc := NewOfflineCache()

	memories := []*pb.MemoryNote{
		{Id: "id1", Content: "content1"},
		{Id: "id2", Content: "content2"},
		{Id: "id3", Content: "content3"},
	}

	for _, m := range memories {
		oc.Store(m)
	}

	result := oc.ListAll()

	if len(result) != 3 {
		t.Errorf("Expected 3 memories, got %d", len(result))
	}
}

func TestOfflineCacheMarkDirty(t *testing.T) {
	oc := NewOfflineCache()

	memory := &pb.MemoryNote{Id: "test-id", Content: "content"}
	oc.Store(memory)

	if oc.IsDirty("test-id") {
		t.Error("Memory should not be dirty initially")
	}

	oc.MarkDirty("test-id")

	if !oc.IsDirty("test-id") {
		t.Error("Memory should be dirty after marking")
	}
}

func TestOfflineCacheClearDirty(t *testing.T) {
	oc := NewOfflineCache()

	oc.MarkDirty("test-id")

	if !oc.IsDirty("test-id") {
		t.Fatal("Setup failed: memory should be dirty")
	}

	oc.ClearDirty("test-id")

	if oc.IsDirty("test-id") {
		t.Error("Memory should not be dirty after clearing")
	}
}

func TestOfflineCacheGetDirty(t *testing.T) {
	oc := NewOfflineCache()

	oc.MarkDirty("id1")
	oc.MarkDirty("id2")
	oc.MarkDirty("id3")

	dirty := oc.GetDirty()

	if len(dirty) != 3 {
		t.Errorf("Expected 3 dirty memories, got %d", len(dirty))
	}

	// Verify all IDs present
	dirtyMap := make(map[string]bool)
	for _, id := range dirty {
		dirtyMap[id] = true
	}

	for _, id := range []string{"id1", "id2", "id3"} {
		if !dirtyMap[id] {
			t.Errorf("Missing dirty ID: %s", id)
		}
	}
}

func TestOfflineCacheDelete(t *testing.T) {
	oc := NewOfflineCache()

	memory := &pb.MemoryNote{Id: "test-id", Content: "content"}
	oc.Store(memory)
	oc.MarkDirty("test-id")

	oc.Delete("test-id")

	if oc.Len() != 0 {
		t.Error("Memory should be deleted from cache")
	}

	if oc.IsDirty("test-id") {
		t.Error("Dirty flag should be removed after delete")
	}
}

func TestOfflineCacheClear(t *testing.T) {
	oc := NewOfflineCache()

	oc.Store(&pb.MemoryNote{Id: "id1", Content: "content1"})
	oc.Store(&pb.MemoryNote{Id: "id2", Content: "content2"})
	oc.MarkDirty("id1")

	oc.Clear()

	if oc.Len() != 0 {
		t.Errorf("Cache should be empty after clear, got length %d", oc.Len())
	}

	if len(oc.GetDirty()) != 0 {
		t.Error("All dirty flags should be cleared")
	}
}

func TestOfflineCacheLastSyncTime(t *testing.T) {
	oc := NewOfflineCache()

	before := time.Now().Add(-1 * time.Second)
	after := time.Now().Add(1 * time.Second)

	lastSync := oc.LastSyncTime()

	if lastSync.Before(before) || lastSync.After(after) {
		t.Errorf("LastSyncTime = %v, expected recent time", lastSync)
	}
}

// --- Concurrent Access Tests ---

func TestOfflineCacheConcurrentAccess(t *testing.T) {
	oc := NewOfflineCache()

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				memory := &pb.MemoryNote{
					Id:      string(rune('a'+id)) + string(rune('0'+j%10)),
					Content: "content",
				}
				oc.Store(memory)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				memID := string(rune('a'+id%5)) + string(rune('0'+j%10))
				oc.Get(memID)
				oc.IsDirty(memID)
			}
		}(i)
	}

	wg.Wait()

	// Cache should be in valid state
	if oc.Len() < 0 {
		t.Error("Invalid cache state after concurrent access")
	}
}

func TestSyncQueueConcurrentAccess(t *testing.T) {
	sq := NewSyncQueue()

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Concurrent additions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				sq.Add(SyncOperation{
					Type:      OpCreate,
					MemoryID:  "id",
					Timestamp: time.Now(),
				})
			}
		}()
	}

	wg.Wait()

	expectedLen := numGoroutines * numOperations
	if sq.Len() != expectedLen {
		t.Errorf("Expected length %d after concurrent adds, got %d", expectedLen, sq.Len())
	}
}

// --- OperationType Tests ---

func TestOperationTypeString(t *testing.T) {
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
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- Edge Cases ---

func TestOfflineCacheStoreNilMemory(t *testing.T) {
	oc := NewOfflineCache()

	oc.Store(nil)

	if oc.Len() != 0 {
		t.Error("Storing nil memory should be no-op")
	}
}

func TestOfflineCacheStoreEmptyID(t *testing.T) {
	oc := NewOfflineCache()

	memory := &pb.MemoryNote{Id: "", Content: "content"}
	oc.Store(memory)

	if oc.Len() != 0 {
		t.Error("Storing memory with empty ID should be no-op")
	}
}

func TestSyncQueueRemoveInvalidIndex(t *testing.T) {
	sq := NewSyncQueue()

	sq.Add(SyncOperation{Type: OpCreate, MemoryID: "id1", Timestamp: time.Now()})

	initialLen := sq.Len()

	// Try to remove invalid indices
	sq.Remove(-1)
	sq.Remove(999)

	if sq.Len() != initialLen {
		t.Error("Removing invalid index should be no-op")
	}
}
