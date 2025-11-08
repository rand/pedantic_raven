package buffer

import (
	"fmt"
	"io/ioutil"
	"sync"
)

// Manager manages multiple buffers and tracks the active buffer.
//
// The manager provides:
// - Buffer registry and lifecycle management
// - Active buffer tracking and switching
// - Buffer creation from files or scratch
// - Thread-safe operations
type Manager struct {
	mu      sync.RWMutex
	buffers map[BufferID]Buffer
	active  BufferID
	nextID  int
}

// NewManager creates a new buffer manager.
func NewManager() *Manager {
	return &Manager{
		buffers: make(map[BufferID]Buffer),
		active:  "",
		nextID:  1,
	}
}

// Create creates a new empty buffer and makes it active.
//
// If id is empty, a unique ID will be generated.
// Returns the created buffer.
func (m *Manager) Create(id BufferID) Buffer {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate ID if not provided
	if id == "" {
		id = m.generateID()
	}

	// Create buffer
	buf := NewBuffer(id)
	m.buffers[id] = buf
	m.active = id

	return buf
}

// Open creates a buffer from a file and makes it active.
//
// If a buffer for this path already exists, it will be made active
// and returned without reloading.
// Returns the buffer and any error from reading the file.
func (m *Manager) Open(path string) (Buffer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already open
	for _, buf := range m.buffers {
		if buf.Path() == path {
			m.active = buf.ID()
			return buf, nil
		}
	}

	// Read file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	// Create buffer
	id := m.generateID()
	buf := NewBufferFromContent(id, string(content))
	buf.SetPath(path)
	buf.MarkClean() // File content is clean

	m.buffers[id] = buf
	m.active = id

	return buf, nil
}

// Close removes a buffer from the manager.
//
// If the buffer is active, the active buffer will be cleared.
// Returns true if the buffer was closed, false if not found.
func (m *Manager) Close(id BufferID) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.buffers[id]; !exists {
		return false
	}

	delete(m.buffers, id)

	// Clear active if this was the active buffer
	if m.active == id {
		m.active = ""
		// Try to activate another buffer
		for bufID := range m.buffers {
			m.active = bufID
			break
		}
	}

	return true
}

// Get retrieves a buffer by ID.
// Returns nil if the buffer doesn't exist.
func (m *Manager) Get(id BufferID) Buffer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.buffers[id]
}

// Active returns the currently active buffer.
// Returns nil if no buffer is active.
func (m *Manager) Active() Buffer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.active == "" {
		return nil
	}

	return m.buffers[m.active]
}

// ActiveID returns the ID of the active buffer.
// Returns empty string if no buffer is active.
func (m *Manager) ActiveID() BufferID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.active
}

// SwitchTo makes the specified buffer active.
//
// Returns true if the switch was successful, false if the buffer doesn't exist.
func (m *Manager) SwitchTo(id BufferID) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.buffers[id]; !exists {
		return false
	}

	m.active = id
	return true
}

// All returns all buffers.
// Returns a slice of buffers in arbitrary order.
func (m *Manager) All() []Buffer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	buffers := make([]Buffer, 0, len(m.buffers))
	for _, buf := range m.buffers {
		buffers = append(buffers, buf)
	}

	return buffers
}

// AllIDs returns all buffer IDs.
// Returns a slice of IDs in arbitrary order.
func (m *Manager) AllIDs() []BufferID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]BufferID, 0, len(m.buffers))
	for id := range m.buffers {
		ids = append(ids, id)
	}

	return ids
}

// Count returns the number of buffers.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.buffers)
}

// HasUnsaved returns true if any buffer has unsaved changes.
func (m *Manager) HasUnsaved() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, buf := range m.buffers {
		if buf.IsDirty() {
			return true
		}
	}

	return false
}

// UnsavedBuffers returns all buffers with unsaved changes.
func (m *Manager) UnsavedBuffers() []Buffer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var unsaved []Buffer
	for _, buf := range m.buffers {
		if buf.IsDirty() {
			unsaved = append(unsaved, buf)
		}
	}

	return unsaved
}

// CloseAll closes all buffers.
//
// This clears the active buffer and removes all buffers from the registry.
func (m *Manager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.buffers = make(map[BufferID]Buffer)
	m.active = ""
}

// Next switches to the next buffer in the registry.
//
// If there's only one buffer, it remains active.
// If no buffer is active, the first buffer becomes active.
// Returns true if a buffer is now active, false if there are no buffers.
func (m *Manager) Next() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.buffers) == 0 {
		return false
	}

	// Get all IDs in a consistent order
	ids := make([]BufferID, 0, len(m.buffers))
	for id := range m.buffers {
		ids = append(ids, id)
	}

	// Find current position
	currentIdx := -1
	for i, id := range ids {
		if id == m.active {
			currentIdx = i
			break
		}
	}

	// Move to next
	nextIdx := (currentIdx + 1) % len(ids)
	m.active = ids[nextIdx]

	return true
}

// Previous switches to the previous buffer in the registry.
//
// If there's only one buffer, it remains active.
// If no buffer is active, the first buffer becomes active.
// Returns true if a buffer is now active, false if there are no buffers.
func (m *Manager) Previous() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.buffers) == 0 {
		return false
	}

	// Get all IDs in a consistent order
	ids := make([]BufferID, 0, len(m.buffers))
	for id := range m.buffers {
		ids = append(ids, id)
	}

	// Find current position
	currentIdx := -1
	for i, id := range ids {
		if id == m.active {
			currentIdx = i
			break
		}
	}

	// Move to previous
	prevIdx := currentIdx - 1
	if prevIdx < 0 {
		prevIdx = len(ids) - 1
	}
	m.active = ids[prevIdx]

	return true
}

// --- Internal helpers ---

// generateID generates a unique buffer ID.
// Caller must hold the lock.
func (m *Manager) generateID() BufferID {
	for {
		id := BufferID(fmt.Sprintf("buffer-%d", m.nextID))
		m.nextID++

		// Ensure uniqueness (very unlikely to collide)
		if _, exists := m.buffers[id]; !exists {
			return id
		}
	}
}
