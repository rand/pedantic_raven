package memorydetail

import (
	"testing"
	"time"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Test helpers

func createTestMemory(id, content string, importance uint32, tags []string) *pb.MemoryNote {
	now := time.Now()
	return &pb.MemoryNote{
		Id:         id,
		Content:    content,
		Importance: importance,
		Tags:       tags,
		Namespace: &pb.Namespace{
			Namespace: &pb.Namespace_Global{
				Global: &pb.GlobalNamespace{},
			},
		},
		CreatedAt: uint64(now.Add(-24 * time.Hour).Unix()),
		UpdatedAt: uint64(now.Add(-time.Hour).Unix()),
		Links:     []*pb.MemoryLink{},
	}
}

// --- Model Creation Tests ---

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.memory != nil {
		t.Error("Expected nil memory for new model")
	}

	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset 0, got %d", m.scrollOffset)
	}

	if !m.showMetadata {
		t.Error("Expected metadata to be shown by default")
	}

	if !m.focused {
		t.Error("Expected model to be focused by default")
	}
}

func TestNewModelWithMemory(t *testing.T) {
	memory := createTestMemory("test-1", "Test content", 8, []string{"tag1"})
	m := NewModelWithMemory(memory)

	if m.memory == nil {
		t.Fatal("Expected non-nil memory")
	}

	if m.memory.Id != "test-1" {
		t.Errorf("Expected memory ID 'test-1', got '%s'", m.memory.Id)
	}
}

// --- Size and Focus Tests ---

func TestSetSize(t *testing.T) {
	m := NewModel()
	m.SetSize(100, 30)

	if m.width != 100 {
		t.Errorf("Expected width 100, got %d", m.width)
	}

	if m.height != 30 {
		t.Errorf("Expected height 30, got %d", m.height)
	}
}

func TestSetFocus(t *testing.T) {
	m := NewModel()
	m.SetFocus(false)

	if m.IsFocused() {
		t.Error("Expected unfocused")
	}

	m.SetFocus(true)

	if !m.IsFocused() {
		t.Error("Expected focused")
	}
}

// --- Memory Management Tests ---

func TestSetMemory(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Test content", 8, []string{"tag1"})

	m.SetMemory(memory)

	if m.Memory() == nil {
		t.Fatal("Expected non-nil memory")
	}

	if m.Memory().Id != "test-1" {
		t.Errorf("Expected memory ID 'test-1', got '%s'", m.Memory().Id)
	}

	if m.scrollOffset != 0 {
		t.Error("Expected scroll offset to be reset when setting new memory")
	}

	if m.err != nil {
		t.Error("Expected error to be cleared when setting memory")
	}
}

func TestSetError(t *testing.T) {
	m := NewModel()
	testErr := &testError{msg: "test error"}

	m.SetError(testErr)

	if m.Error() == nil {
		t.Fatal("Expected non-nil error")
	}

	if m.Error().Error() != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", m.Error().Error())
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// --- Metadata Toggle Tests ---

func TestToggleMetadata(t *testing.T) {
	m := NewModel()

	// Should be shown by default
	if !m.ShowMetadata() {
		t.Error("Expected metadata to be shown by default")
	}

	// Toggle off
	m.ToggleMetadata()
	if m.ShowMetadata() {
		t.Error("Expected metadata to be hidden after toggle")
	}

	// Toggle back on
	m.ToggleMetadata()
	if !m.ShowMetadata() {
		t.Error("Expected metadata to be shown after second toggle")
	}
}

// --- Client Management Tests ---

func TestSetClient(t *testing.T) {
	m := NewModel()

	// Mock client
	client := &mockClient{}

	m.SetClient(client)

	if m.Client() == nil {
		t.Error("Expected non-nil client")
	}
}

type mockClient struct{}

// --- Message Type Tests ---

func TestMemoryLoadedMsg(t *testing.T) {
	memory := createTestMemory("test-1", "Content", 8, nil)
	msg := MemoryLoadedMsg{Memory: memory}

	if msg.Memory.Id != "test-1" {
		t.Errorf("Expected memory ID 'test-1', got '%s'", msg.Memory.Id)
	}
}

func TestMemoryErrorMsg(t *testing.T) {
	err := &testError{msg: "load failed"}
	msg := MemoryErrorMsg{Err: err}

	if msg.Err.Error() != "load failed" {
		t.Errorf("Expected error message 'load failed', got '%s'", msg.Err.Error())
	}
}

func TestCloseRequestMsg(t *testing.T) {
	msg := CloseRequestMsg{}

	// Just verify it can be created
	_ = msg
}

func TestLinkSelectedMsg(t *testing.T) {
	msg := LinkSelectedMsg{TargetID: "target-123"}

	if msg.TargetID != "target-123" {
		t.Errorf("Expected target ID 'target-123', got '%s'", msg.TargetID)
	}
}

// --- Init Tests ---

func TestInit(t *testing.T) {
	m := NewModel()

	cmd := m.Init()

	if cmd != nil {
		t.Error("Expected nil cmd from Init")
	}
}
