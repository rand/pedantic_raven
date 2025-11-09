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

// --- Link Navigation Tests ---

func TestSelectNextLink(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
		{TargetId: "link-3"},
	}
	m.SetMemory(memory)

	// Initially no link selected
	if m.SelectedLinkIndex() != -1 {
		t.Error("Expected no link selected initially")
	}

	// Select next (should go to index 0)
	m.SelectNextLink()
	if m.SelectedLinkIndex() != 0 {
		t.Errorf("Expected link index 0, got %d", m.SelectedLinkIndex())
	}

	// Select next again (should go to index 1)
	m.SelectNextLink()
	if m.SelectedLinkIndex() != 1 {
		t.Errorf("Expected link index 1, got %d", m.SelectedLinkIndex())
	}

	// Select next again (should go to index 2)
	m.SelectNextLink()
	if m.SelectedLinkIndex() != 2 {
		t.Errorf("Expected link index 2, got %d", m.SelectedLinkIndex())
	}

	// Try to go beyond last link (should stay at 2)
	m.SelectNextLink()
	if m.SelectedLinkIndex() != 2 {
		t.Errorf("Expected to stay at link index 2, got %d", m.SelectedLinkIndex())
	}
}

func TestSelectPreviousLink(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
		{TargetId: "link-3"},
	}
	m.SetMemory(memory)

	// Start at link 2
	m.selectedLinkIndex = 2

	// Previous should go to 1
	m.SelectPreviousLink()
	if m.SelectedLinkIndex() != 1 {
		t.Errorf("Expected link index 1, got %d", m.SelectedLinkIndex())
	}

	// Previous should go to 0
	m.SelectPreviousLink()
	if m.SelectedLinkIndex() != 0 {
		t.Errorf("Expected link index 0, got %d", m.SelectedLinkIndex())
	}

	// Previous should go to -1 (no selection)
	m.SelectPreviousLink()
	if m.SelectedLinkIndex() != -1 {
		t.Errorf("Expected link index -1, got %d", m.SelectedLinkIndex())
	}

	// Try to go before first (should stay at -1)
	m.SelectPreviousLink()
	if m.SelectedLinkIndex() != -1 {
		t.Errorf("Expected to stay at link index -1, got %d", m.SelectedLinkIndex())
	}
}

func TestSelectFirstLink(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
	}
	m.SetMemory(memory)

	m.SelectFirstLink()

	if m.SelectedLinkIndex() != 0 {
		t.Errorf("Expected link index 0, got %d", m.SelectedLinkIndex())
	}
}

func TestClearLinkSelection(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
	}
	m.SetMemory(memory)

	m.selectedLinkIndex = 0

	m.ClearLinkSelection()

	if m.SelectedLinkIndex() != -1 {
		t.Errorf("Expected link index -1 after clear, got %d", m.SelectedLinkIndex())
	}
}

func TestSelectedLink(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
	}
	m.SetMemory(memory)

	// No selection
	if m.SelectedLink() != nil {
		t.Error("Expected nil selected link when index is -1")
	}

	// Select first link
	m.selectedLinkIndex = 0
	link := m.SelectedLink()
	if link == nil {
		t.Fatal("Expected non-nil selected link")
	}
	if link.TargetId != "link-1" {
		t.Errorf("Expected link-1, got %s", link.TargetId)
	}

	// Select second link
	m.selectedLinkIndex = 1
	link = m.SelectedLink()
	if link == nil {
		t.Fatal("Expected non-nil selected link")
	}
	if link.TargetId != "link-2" {
		t.Errorf("Expected link-2, got %s", link.TargetId)
	}
}

func TestHasLinks(t *testing.T) {
	m := NewModel()

	// No memory
	if m.HasLinks() {
		t.Error("Expected no links when memory is nil")
	}

	// Memory with no links
	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	if m.HasLinks() {
		t.Error("Expected no links when links array is empty")
	}

	// Memory with links
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
	}
	m.SetMemory(memory)

	if !m.HasLinks() {
		t.Error("Expected to have links")
	}
}

func TestSetMemoryResetsLinkSelection(t *testing.T) {
	m := NewModel()
	memory1 := createTestMemory("test-1", "Content", 8, nil)
	memory1.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
	}
	m.SetMemory(memory1)
	m.selectedLinkIndex = 0

	// Set new memory should reset selection
	memory2 := createTestMemory("test-2", "New content", 8, nil)
	m.SetMemory(memory2)

	if m.SelectedLinkIndex() != -1 {
		t.Errorf("Expected link selection to be reset, got index %d", m.SelectedLinkIndex())
	}
}

func TestLinkNavigationWithNoLinks(t *testing.T) {
	m := NewModel()
	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	// These should not crash
	m.SelectNextLink()
	m.SelectPreviousLink()
	m.SelectFirstLink()

	// Should stay at -1
	if m.SelectedLinkIndex() != -1 {
		t.Error("Expected link index to remain -1 when there are no links")
	}
}
