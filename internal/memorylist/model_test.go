package memorylist

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Test helpers

func createTestMemory(id, content string, importance uint32, tags []string, updatedAgo time.Duration) *pb.MemoryNote {
	return &pb.MemoryNote{
		MemoryId:   id,
		Content:    content,
		Importance: importance,
		Tags:       tags,
		Namespace: &pb.Namespace{
			Namespace: &pb.Namespace_Global{
				Global: &pb.GlobalNamespace{},
			},
		},
		CreatedAt: timestamppb.New(time.Now().Add(-24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now().Add(-updatedAgo)),
		Metadata: &pb.MemoryMetadata{
			Links: []*pb.MemoryLink{},
		},
	}
}

// --- Model Creation Tests ---

func TestNewModel(t *testing.T) {
	m := NewModel()

	if m.memories == nil {
		t.Error("Expected non-nil memories slice")
	}

	if m.filteredMems == nil {
		t.Error("Expected non-nil filtered memories slice")
	}

	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex 0, got %d", m.selectedIndex)
	}

	if m.sortBy != SortByUpdated {
		t.Errorf("Expected sort by updated, got %v", m.sortBy)
	}

	if !m.sortDesc {
		t.Error("Expected descending sort by default")
	}

	if !m.focused {
		t.Error("Expected focused by default")
	}
}

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

// --- Memory Loading Tests ---

func TestSetMemories(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First memory", 8, []string{"tag1"}, time.Hour),
		createTestMemory("2", "Second memory", 6, []string{"tag2"}, 2*time.Hour),
	}

	m.SetMemories(memories, 2)

	if len(m.memories) != 2 {
		t.Errorf("Expected 2 memories, got %d", len(m.memories))
	}

	if m.totalCount != 2 {
		t.Errorf("Expected total count 2, got %d", m.totalCount)
	}

	if m.loading {
		t.Error("Expected loading to be false after setting memories")
	}

	if m.err != nil {
		t.Errorf("Expected no error, got %v", m.err)
	}
}

func TestSetLoading(t *testing.T) {
	m := NewModel()
	m.SetLoading(true)

	if !m.IsLoading() {
		t.Error("Expected loading state")
	}

	m.SetLoading(false)

	if m.IsLoading() {
		t.Error("Expected not loading")
	}
}

func TestSetError(t *testing.T) {
	m := NewModel()
	err := ErrTestError
	m.SetError(err)

	if m.Error() != err {
		t.Errorf("Expected error %v, got %v", err, m.Error())
	}

	if m.loading {
		t.Error("Expected loading to be false after error")
	}
}

var ErrTestError = tea.Quit().(error)

// --- Selection Tests ---

func TestSelectedMemory(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 2)

	// First memory selected by default
	selected := m.SelectedMemory()
	if selected == nil {
		t.Fatal("Expected selected memory")
	}

	if selected.MemoryId != "1" {
		t.Errorf("Expected selected memory ID '1', got '%s'", selected.MemoryId)
	}
}

func TestSelectedMemoryEmpty(t *testing.T) {
	m := NewModel()

	selected := m.SelectedMemory()
	if selected != nil {
		t.Error("Expected nil selected memory when list is empty")
	}
}

func TestMemoryCount(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 10)

	if m.MemoryCount() != 2 {
		t.Errorf("Expected memory count 2, got %d", m.MemoryCount())
	}

	if m.TotalCount() != 10 {
		t.Errorf("Expected total count 10, got %d", m.TotalCount())
	}
}

// --- Navigation Tests ---

func TestMoveDown(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
		createTestMemory("3", "Third", 5, nil, 3*time.Hour),
	}
	m.SetMemories(memories, 3)

	// Move down once
	m = m.moveDown()
	if m.selectedIndex != 1 {
		t.Errorf("Expected selectedIndex 1, got %d", m.selectedIndex)
	}

	// Move down again
	m = m.moveDown()
	if m.selectedIndex != 2 {
		t.Errorf("Expected selectedIndex 2, got %d", m.selectedIndex)
	}

	// Try to move down past end (should stay at end)
	m = m.moveDown()
	if m.selectedIndex != 2 {
		t.Errorf("Expected selectedIndex to stay at 2, got %d", m.selectedIndex)
	}
}

func TestMoveUp(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
		createTestMemory("3", "Third", 5, nil, 3*time.Hour),
	}
	m.SetMemories(memories, 3)

	m.selectedIndex = 2

	// Move up once
	m = m.moveUp()
	if m.selectedIndex != 1 {
		t.Errorf("Expected selectedIndex 1, got %d", m.selectedIndex)
	}

	// Move up again
	m = m.moveUp()
	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex 0, got %d", m.selectedIndex)
	}

	// Try to move up past start (should stay at start)
	m = m.moveUp()
	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex to stay at 0, got %d", m.selectedIndex)
	}
}

func TestMoveToTop(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 2)
	m.selectedIndex = 1

	m = m.moveToTop()

	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex 0, got %d", m.selectedIndex)
	}

	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset 0, got %d", m.scrollOffset)
	}
}

func TestMoveToBottom(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 2)

	m = m.moveToBottom()

	if m.selectedIndex != 1 {
		t.Errorf("Expected selectedIndex 1, got %d", m.selectedIndex)
	}
}

// --- Filtering Tests ---

func TestSearchFilter(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "Authentication system", 8, nil, time.Hour),
		createTestMemory("2", "Database schema", 6, nil, 2*time.Hour),
		createTestMemory("3", "Auth middleware", 7, nil, 3*time.Hour),
	}
	m.SetMemories(memories, 3)

	m.SetSearchQuery("auth")

	if len(m.filteredMems) != 2 {
		t.Errorf("Expected 2 filtered memories, got %d", len(m.filteredMems))
	}
}

func TestTagFilter(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, []string{"backend", "api"}, time.Hour),
		createTestMemory("2", "Second", 6, []string{"frontend"}, 2*time.Hour),
		createTestMemory("3", "Third", 7, []string{"backend", "database"}, 3*time.Hour),
	}
	m.SetMemories(memories, 3)

	m.SetFilter([]string{"backend"}, "", 0)

	if len(m.filteredMems) != 2 {
		t.Errorf("Expected 2 filtered memories (with backend tag), got %d", len(m.filteredMems))
	}
}

func TestImportanceFilter(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
		createTestMemory("3", "Third", 9, nil, 3*time.Hour),
	}
	m.SetMemories(memories, 3)

	m.SetFilter(nil, "", 7)

	if len(m.filteredMems) != 2 {
		t.Errorf("Expected 2 filtered memories (importance >= 7), got %d", len(m.filteredMems))
	}
}

func TestClearFilters(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, []string{"tag1"}, time.Hour),
		createTestMemory("2", "Second", 6, []string{"tag2"}, 2*time.Hour),
	}
	m.SetMemories(memories, 2)

	m.SetSearchQuery("first")
	m.SetFilter([]string{"tag1"}, "", 7)

	if len(m.filteredMems) != 1 {
		t.Errorf("Expected 1 filtered memory before clear, got %d", len(m.filteredMems))
	}

	m.ClearFilters()

	if len(m.filteredMems) != 2 {
		t.Errorf("Expected 2 memories after clearing filters, got %d", len(m.filteredMems))
	}

	if m.searchQuery != "" {
		t.Error("Expected empty search query after clear")
	}
}

// --- Sorting Tests ---

func TestSortByImportance(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 5, nil, time.Hour),
		createTestMemory("2", "Second", 9, nil, 2*time.Hour),
		createTestMemory("3", "Third", 7, nil, 3*time.Hour),
	}
	m.SetMemories(memories, 3)

	m.SetSort(SortByImportance, true) // Descending

	if m.filteredMems[0].MemoryId != "2" {
		t.Errorf("Expected first memory to be ID '2' (importance 9), got '%s'", m.filteredMems[0].MemoryId)
	}

	if m.filteredMems[2].MemoryId != "1" {
		t.Errorf("Expected last memory to be ID '1' (importance 5), got '%s'", m.filteredMems[2].MemoryId)
	}
}

func TestSortByUpdated(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, 3*time.Hour),
		createTestMemory("2", "Second", 6, nil, time.Hour),
		createTestMemory("3", "Third", 7, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 3)

	m.SetSort(SortByUpdated, true) // Most recent first

	if m.filteredMems[0].MemoryId != "2" {
		t.Errorf("Expected first memory to be ID '2' (most recent), got '%s'", m.filteredMems[0].MemoryId)
	}

	if m.filteredMems[2].MemoryId != "1" {
		t.Errorf("Expected last memory to be ID '1' (least recent), got '%s'", m.filteredMems[2].MemoryId)
	}
}

// --- Message Handling Tests ---

func TestHandleMemoriesLoadedMsg(t *testing.T) {
	m := NewModel()
	m.SetLoading(true)

	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
	}

	msg := MemoriesLoadedMsg{
		Memories:   memories,
		TotalCount: 10,
	}

	m, _ = m.Update(msg)

	if m.IsLoading() {
		t.Error("Expected loading to be false after receiving memories")
	}

	if len(m.memories) != 1 {
		t.Errorf("Expected 1 memory, got %d", len(m.memories))
	}

	if m.TotalCount() != 10 {
		t.Errorf("Expected total count 10, got %d", m.TotalCount())
	}
}

func TestHandleMemoriesErrorMsg(t *testing.T) {
	m := NewModel()
	m.SetLoading(true)

	msg := MemoriesErrorMsg{
		Err: ErrTestError,
	}

	m, _ = m.Update(msg)

	if m.IsLoading() {
		t.Error("Expected loading to be false after error")
	}

	if m.Error() == nil {
		t.Error("Expected error to be set")
	}
}

func TestHandleKeyPress(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
		createTestMemory("2", "Second", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 2)

	// Test 'j' key (move down)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.selectedIndex != 1 {
		t.Errorf("Expected selectedIndex 1 after 'j', got %d", m.selectedIndex)
	}

	// Test 'k' key (move up)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.selectedIndex != 0 {
		t.Errorf("Expected selectedIndex 0 after 'k', got %d", m.selectedIndex)
	}
}

func TestUnfocusedIgnoresKeys(t *testing.T) {
	m := NewModel()
	m.SetFocus(false)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, nil, time.Hour),
	}
	m.SetMemories(memories, 1)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if m.selectedIndex != 0 {
		t.Error("Expected unfocused component to ignore key presses")
	}
}

// --- SortMode String Tests ---

func TestSortModeString(t *testing.T) {
	tests := []struct {
		mode     SortMode
		expected string
	}{
		{SortByImportance, "Importance"},
		{SortByUpdated, "Updated"},
		{SortByCreated, "Created"},
		{SortByRelevance, "Relevance"},
	}

	for _, tt := range tests {
		if tt.mode.String() != tt.expected {
			t.Errorf("Expected SortMode %d to be '%s', got '%s'", tt.mode, tt.expected, tt.mode.String())
		}
	}
}
