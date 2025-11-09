package testhelpers

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// MockMnemosyneClient is a mock implementation of the mnemosyne client for testing
type MockMnemosyneClient struct {
	mu       sync.RWMutex
	memories map[string]*pb.MemoryNote
	links    map[string][]*pb.MemoryLink
	nextID   int

	// Error simulation
	shouldFail   bool
	failureError error

	// Operation tracking
	RememberCalls int
	RecallCalls   int
	UpdateCalls   int
	DeleteCalls   int
}

// NewMockMnemosyneClient creates a new mock client
func NewMockMnemosyneClient() *MockMnemosyneClient {
	return &MockMnemosyneClient{
		memories: make(map[string]*pb.MemoryNote),
		links:    make(map[string][]*pb.MemoryLink),
		nextID:   1,
	}
}

// Remember stores a new memory
func (m *MockMnemosyneClient) Remember(ctx context.Context, content string, importance uint32, tags []string, namespace string) (*pb.MemoryNote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RememberCalls++

	if m.shouldFail {
		return nil, m.failureError
	}

	id := fmt.Sprintf("mem-%d", m.nextID)
	m.nextID++

	memory := &pb.MemoryNote{
		Id:          id,
		Content:     content,
		Importance:  importance,
		Tags:        tags,
		Namespace:   &pb.Namespace{Namespace: &pb.Namespace_Global{Global: &pb.GlobalNamespace{}}},
		CreatedAt:   uint64(time.Now().Unix()),
		UpdatedAt:   uint64(time.Now().Unix()),
		AccessCount: 0,
	}

	m.memories[id] = memory
	return memory, nil
}

// Recall retrieves memories
func (m *MockMnemosyneClient) Recall(ctx context.Context, query string, limit int) ([]*pb.MemoryNote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.RecallCalls++

	if m.shouldFail {
		return nil, m.failureError
	}

	var results []*pb.MemoryNote
	for _, mem := range m.memories {
		results = append(results, mem)
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

// Update modifies an existing memory
func (m *MockMnemosyneClient) Update(ctx context.Context, id string, content string, importance uint32, tags []string) (*pb.MemoryNote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.UpdateCalls++

	if m.shouldFail {
		return nil, m.failureError
	}

	memory, exists := m.memories[id]
	if !exists {
		return nil, fmt.Errorf("memory not found: %s", id)
	}

	memory.Content = content
	memory.Importance = importance
	memory.Tags = tags
	memory.UpdatedAt = uint64(time.Now().Unix())

	m.memories[id] = memory
	return memory, nil
}

// Delete removes a memory
func (m *MockMnemosyneClient) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DeleteCalls++

	if m.shouldFail {
		return m.failureError
	}

	if _, exists := m.memories[id]; !exists {
		return fmt.Errorf("memory not found: %s", id)
	}

	delete(m.memories, id)
	delete(m.links, id)

	return nil
}

// GetMemory retrieves a single memory by ID
func (m *MockMnemosyneClient) GetMemory(ctx context.Context, id string) (*pb.MemoryNote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, m.failureError
	}

	memory, exists := m.memories[id]
	if !exists {
		return nil, fmt.Errorf("memory not found: %s", id)
	}

	return memory, nil
}

// GetLinks retrieves links for a memory
func (m *MockMnemosyneClient) GetLinks(ctx context.Context, id string) ([]*pb.MemoryLink, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return nil, m.failureError
	}

	links, exists := m.links[id]
	if !exists {
		return []*pb.MemoryLink{}, nil
	}

	return links, nil
}

// AddLink creates a link between memories
func (m *MockMnemosyneClient) AddLink(ctx context.Context, sourceID, targetID string, linkType pb.LinkType) (*pb.MemoryLink, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return nil, m.failureError
	}

	link := &pb.MemoryLink{
		TargetId:  targetID,
		LinkType:  linkType,
		Strength:  0.8,
		Reason:    "test link",
		CreatedAt: uint64(time.Now().Unix()),
		UserCreated: true,
	}

	m.links[sourceID] = append(m.links[sourceID], link)

	return link, nil
}

// SetShouldFail configures the mock to fail operations
func (m *MockMnemosyneClient) SetShouldFail(shouldFail bool, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.shouldFail = shouldFail
	m.failureError = err
}

// GetMemoryCount returns the number of stored memories
func (m *MockMnemosyneClient) GetMemoryCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.memories)
}

// GetAllMemories returns all stored memories
func (m *MockMnemosyneClient) GetAllMemories() []*pb.MemoryNote {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*pb.MemoryNote
	for _, mem := range m.memories {
		results = append(results, mem)
	}

	return results
}

// Reset clears all stored data
func (m *MockMnemosyneClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.memories = make(map[string]*pb.MemoryNote)
	m.links = make(map[string][]*pb.MemoryLink)
	m.nextID = 1
	m.shouldFail = false
	m.failureError = nil
	m.RememberCalls = 0
	m.RecallCalls = 0
	m.UpdateCalls = 0
	m.DeleteCalls = 0
}
