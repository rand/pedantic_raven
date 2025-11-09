package mnemosyne

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// OperationType defines the type of operation for sync queue.
type OperationType int

const (
	OpCreate OperationType = iota
	OpUpdate
	OpDelete
)

// String returns the string representation of the operation type.
func (ot OperationType) String() string {
	switch ot {
	case OpCreate:
		return "create"
	case OpUpdate:
		return "update"
	case OpDelete:
		return "delete"
	default:
		return "unknown"
	}
}

// SyncOperation represents a pending operation to be synced.
type SyncOperation struct {
	Type      OperationType
	MemoryID  string
	Memory    *pb.MemoryNote
	Timestamp time.Time
}

// SyncQueue tracks pending operations.
type SyncQueue struct {
	operations []SyncOperation
	mu         sync.Mutex
}

// NewSyncQueue creates a new sync queue.
func NewSyncQueue() *SyncQueue {
	return &SyncQueue{
		operations: make([]SyncOperation, 0),
	}
}

// Add adds an operation to the queue.
func (sq *SyncQueue) Add(op SyncOperation) {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	sq.operations = append(sq.operations, op)
}

// GetAll returns all operations in the queue.
func (sq *SyncQueue) GetAll() []SyncOperation {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	// Return a copy to prevent external modification
	result := make([]SyncOperation, len(sq.operations))
	copy(result, sq.operations)
	return result
}

// Remove removes an operation at the specified index.
func (sq *SyncQueue) Remove(index int) {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	if index < 0 || index >= len(sq.operations) {
		return
	}

	sq.operations = append(sq.operations[:index], sq.operations[index+1:]...)
}

// Clear removes all operations from the queue.
func (sq *SyncQueue) Clear() {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	sq.operations = make([]SyncOperation, 0)
}

// Len returns the number of operations in the queue.
func (sq *SyncQueue) Len() int {
	sq.mu.Lock()
	defer sq.mu.Unlock()

	return len(sq.operations)
}

// OfflineCache stores memories when offline.
type OfflineCache struct {
	memories map[string]*pb.MemoryNote
	dirty    map[string]bool // Unsaved changes
	lastSync time.Time
	mu       sync.RWMutex
}

// NewOfflineCache creates a new offline cache.
func NewOfflineCache() *OfflineCache {
	return &OfflineCache{
		memories: make(map[string]*pb.MemoryNote),
		dirty:    make(map[string]bool),
		lastSync: time.Now(),
	}
}

// Store stores a memory in the cache.
func (oc *OfflineCache) Store(memory *pb.MemoryNote) {
	if memory == nil || memory.Id == "" {
		return
	}

	oc.mu.Lock()
	defer oc.mu.Unlock()

	oc.memories[memory.Id] = memory
}

// Get retrieves a memory from the cache.
func (oc *OfflineCache) Get(id string) (*pb.MemoryNote, bool) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	memory, ok := oc.memories[id]
	return memory, ok
}

// ListAll returns all memories in the cache.
func (oc *OfflineCache) ListAll() []*pb.MemoryNote {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	result := make([]*pb.MemoryNote, 0, len(oc.memories))
	for _, memory := range oc.memories {
		result = append(result, memory)
	}
	return result
}

// MarkDirty marks a memory as having unsaved changes.
func (oc *OfflineCache) MarkDirty(id string) {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	oc.dirty[id] = true
}

// IsDirty returns true if the memory has unsaved changes.
func (oc *OfflineCache) IsDirty(id string) bool {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	return oc.dirty[id]
}

// ClearDirty removes the dirty flag for a memory.
func (oc *OfflineCache) ClearDirty(id string) {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	delete(oc.dirty, id)
}

// GetDirty returns all memory IDs with unsaved changes.
func (oc *OfflineCache) GetDirty() []string {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	result := make([]string, 0, len(oc.dirty))
	for id, isDirty := range oc.dirty {
		if isDirty {
			result = append(result, id)
		}
	}
	return result
}

// Delete removes a memory from the cache.
func (oc *OfflineCache) Delete(id string) {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	delete(oc.memories, id)
	delete(oc.dirty, id)
}

// Clear removes all memories from the cache.
func (oc *OfflineCache) Clear() {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	oc.memories = make(map[string]*pb.MemoryNote)
	oc.dirty = make(map[string]bool)
}

// Sync synchronizes the cache with the server using the provided client.
// Returns the number of successfully synced operations and any error.
func (oc *OfflineCache) Sync(client *Client, syncQueue *SyncQueue) (int, error) {
	if client == nil {
		return 0, fmt.Errorf("client is nil")
	}

	if !client.IsConnected() {
		return 0, ErrNotConnected
	}

	// Get all operations from sync queue
	operations := syncQueue.GetAll()
	if len(operations) == 0 {
		// Update last sync time even if no operations
		oc.mu.Lock()
		oc.lastSync = time.Now()
		oc.mu.Unlock()
		return 0, nil
	}

	ctx := context.Background()
	syncedCount := 0

	// Process operations in order (FIFO)
	for i, op := range operations {
		var err error

		switch op.Type {
		case OpCreate:
			if op.Memory == nil {
				continue
			}
			// Convert MemoryNote to StoreMemoryOptions
			importance := uint32(op.Memory.Importance)
			memoryType := op.Memory.MemoryType
			opts := StoreMemoryOptions{
				Content:           op.Memory.Content,
				Namespace:         op.Memory.Namespace,
				Importance:        &importance,
				Context:           op.Memory.Context,
				Tags:              op.Memory.Tags,
				MemoryType:        &memoryType,
				SkipLLMEnrichment: false,
			}
			_, err = client.StoreMemory(ctx, opts)
			if err == nil {
				oc.ClearDirty(op.MemoryID)
			}

		case OpUpdate:
			if op.Memory == nil {
				continue
			}
			// Convert MemoryNote to UpdateMemoryOptions
			content := op.Memory.Content
			importance := uint32(op.Memory.Importance)
			opts := UpdateMemoryOptions{
				MemoryID:   op.Memory.Id,
				Content:    &content,
				Importance: &importance,
				Tags:       op.Memory.Tags,
			}
			_, err = client.UpdateMemory(ctx, opts)
			if err == nil {
				oc.ClearDirty(op.MemoryID)
			}

		case OpDelete:
			err = client.DeleteMemory(ctx, op.MemoryID)
			if err == nil {
				oc.Delete(op.MemoryID)
			}
		}

		if err != nil {
			// If operation fails, return error and number of successful syncs
			return syncedCount, fmt.Errorf("sync failed at operation %d (%s): %w", i, op.Type, err)
		}

		// Remove successfully synced operation
		syncQueue.Remove(i - syncedCount)
		syncedCount++
	}

	// Update last sync time
	oc.mu.Lock()
	oc.lastSync = time.Now()
	oc.mu.Unlock()

	return syncedCount, nil
}

// LastSyncTime returns the last time the cache was synced.
func (oc *OfflineCache) LastSyncTime() time.Time {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	return oc.lastSync
}

// Len returns the number of memories in the cache.
func (oc *OfflineCache) Len() int {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	return len(oc.memories)
}
