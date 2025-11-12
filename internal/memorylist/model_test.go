package memorylist

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Test helpers

func createTestMemory(id, content string, importance uint32, tags []string, updatedAgo time.Duration) *pb.MemoryNote {
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
		UpdatedAt: uint64(now.Add(-updatedAgo).Unix()),
		Links:     []*pb.MemoryLink{},
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

var ErrTestError = errors.New("test error")

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

	if selected.Id != "1" {
		t.Errorf("Expected selected memory ID '1', got '%s'", selected.Id)
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

	if m.filteredMems[0].Id != "2" {
		t.Errorf("Expected first memory to be ID '2' (importance 9), got '%s'", m.filteredMems[0].Id)
	}

	if m.filteredMems[2].Id != "1" {
		t.Errorf("Expected last memory to be ID '1' (importance 5), got '%s'", m.filteredMems[2].Id)
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

	if m.filteredMems[0].Id != "2" {
		t.Errorf("Expected first memory to be ID '2' (most recent), got '%s'", m.filteredMems[0].Id)
	}

	if m.filteredMems[2].Id != "1" {
		t.Errorf("Expected last memory to be ID '1' (least recent), got '%s'", m.filteredMems[2].Id)
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

// --- Input Mode Tests ---

func TestEnterSearchMode(t *testing.T) {
	m := NewModel()
	m.searchQuery = "existing query"

	m = m.enterSearchMode()

	if m.inputMode != InputModeSearch {
		t.Errorf("Expected InputModeSearch, got %v", m.inputMode)
	}

	if m.searchInput != "existing query" {
		t.Errorf("Expected searchInput to be initialized with current query, got '%s'", m.searchInput)
	}
}

func TestSearchInputTyping(t *testing.T) {
	m := NewModel()
	m.inputMode = InputModeSearch

	// Type 'a'
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if m.searchInput != "a" {
		t.Errorf("Expected searchInput 'a', got '%s'", m.searchInput)
	}

	// Type 'u'
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}})
	if m.searchInput != "au" {
		t.Errorf("Expected searchInput 'au', got '%s'", m.searchInput)
	}

	// Type 't'
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	if m.searchInput != "aut" {
		t.Errorf("Expected searchInput 'aut', got '%s'", m.searchInput)
	}

	// Type 'h'
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	if m.searchInput != "auth" {
		t.Errorf("Expected searchInput 'auth', got '%s'", m.searchInput)
	}
}

func TestSearchInputBackspace(t *testing.T) {
	m := NewModel()
	m.inputMode = InputModeSearch
	m.searchInput = "auth"

	// Backspace once
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyBackspace})
	if m.searchInput != "aut" {
		t.Errorf("Expected searchInput 'aut', got '%s'", m.searchInput)
	}

	// Backspace again
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyBackspace})
	if m.searchInput != "au" {
		t.Errorf("Expected searchInput 'au', got '%s'", m.searchInput)
	}

	// Backspace on empty should not crash
	m.searchInput = ""
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyBackspace})
	if m.searchInput != "" {
		t.Errorf("Expected searchInput to remain empty, got '%s'", m.searchInput)
	}
}

func TestSearchInputEnter(t *testing.T) {
	m := NewModel()
	m.inputMode = InputModeSearch
	m.searchInput = "auth"
	memories := []*pb.MemoryNote{
		createTestMemory("1", "Authentication system", 8, nil, time.Hour),
		createTestMemory("2", "Database schema", 6, nil, 2*time.Hour),
	}
	m.SetMemories(memories, 2)

	// Press Enter to commit search
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyEnter})

	if m.inputMode != InputModeNormal {
		t.Errorf("Expected InputModeNormal after Enter, got %v", m.inputMode)
	}

	if m.searchQuery != "auth" {
		t.Errorf("Expected searchQuery 'auth', got '%s'", m.searchQuery)
	}

	// Should filter memories
	if len(m.filteredMems) != 1 {
		t.Errorf("Expected 1 filtered memory, got %d", len(m.filteredMems))
	}
}

func TestSearchInputEscape(t *testing.T) {
	m := NewModel()
	m.inputMode = InputModeSearch
	m.searchInput = "auth"
	m.searchQuery = "existing"

	// Press Escape to cancel search
	m, _ = m.handleSearchInput(tea.KeyMsg{Type: tea.KeyEsc})

	if m.inputMode != InputModeNormal {
		t.Errorf("Expected InputModeNormal after Escape, got %v", m.inputMode)
	}

	if m.searchInput != "" {
		t.Errorf("Expected searchInput to be cleared, got '%s'", m.searchInput)
	}

	if m.searchQuery != "existing" {
		t.Errorf("Expected searchQuery to remain unchanged, got '%s'", m.searchQuery)
	}
}

// --- Keyboard Shortcut Tests ---

func TestSlashKeyEntersSearchMode(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})

	if m.inputMode != InputModeSearch {
		t.Errorf("Expected InputModeSearch after '/', got %v", m.inputMode)
	}
}

func TestQuestionMarkTogglesHelp(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	// First press shows help
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !m.showHelp {
		t.Error("Expected showHelp to be true after first '?'")
	}

	// Second press hides help
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if m.showHelp {
		t.Error("Expected showHelp to be false after second '?'")
	}
}

func TestRKeyReloads(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	if !m.IsLoading() {
		t.Error("Expected loading state after 'r' key")
	}

	if cmd == nil {
		t.Error("Expected non-nil cmd for reload request")
	}

	// Execute cmd to check message type
	msg := cmd()
	if _, ok := msg.(ReloadRequestMsg); !ok {
		t.Errorf("Expected ReloadRequestMsg, got %T", msg)
	}
}

func TestCKeyClearsFilters(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	memories := []*pb.MemoryNote{
		createTestMemory("1", "First", 8, []string{"tag1"}, time.Hour),
		createTestMemory("2", "Second", 6, []string{"tag2"}, 2*time.Hour),
	}
	m.SetMemories(memories, 2)
	m.searchQuery = "first"
	m.searchInput = "first"
	m.SetFilter([]string{"tag1"}, "", 7)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	if m.searchQuery != "" {
		t.Errorf("Expected searchQuery to be cleared, got '%s'", m.searchQuery)
	}

	if m.searchInput != "" {
		t.Errorf("Expected searchInput to be cleared, got '%s'", m.searchInput)
	}

	if len(m.filteredMems) != 2 {
		t.Errorf("Expected all 2 memories after clearing filters, got %d", len(m.filteredMems))
	}
}

func TestEscapeClosesHelp(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.showHelp = true

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if m.showHelp {
		t.Error("Expected showHelp to be false after Escape")
	}
}

func TestEscapeClearsError(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.err = ErrTestError

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if m.err != nil {
		t.Error("Expected error to be cleared after Escape")
	}
}

func TestSearchModeHandlesKeys(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.inputMode = InputModeSearch

	// In search mode, normal keys should not work
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	// Should add 'j' to search input instead of moving down
	if m.searchInput != "j" {
		t.Errorf("Expected searchInput 'j', got '%s'", m.searchInput)
	}

	if m.selectedIndex != 0 {
		t.Error("Expected selectedIndex to remain at 0 in search mode")
	}
}

// --- Relevance Scoring Tests ---

func TestRelevanceScoringExactMatch(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "Authentication system", 5, nil, time.Hour),
		createTestMemory("2", "Database schema design", 8, nil, time.Hour),
		createTestMemory("3", "Authentication middleware", 6, nil, time.Hour),
	}
	m.SetMemories(memories, 3)

	// Set sort mode to relevance and search for "authentication"
	m.SetSort(SortByRelevance, true)
	m.searchQuery = "authentication"
	m.applySorting()

	// Both "1" and "3" have exact matches, but "3" has higher importance (6 vs 5)
	// so it should be first
	if m.filteredMems[0].Id != "3" {
		t.Errorf("Expected first memory to be ID '3' (exact match + higher importance), got '%s'", m.filteredMems[0].Id)
	}

	// Second should be "Authentication system" (exact match but lower importance)
	if m.filteredMems[1].Id != "1" {
		t.Errorf("Expected second memory to be ID '1' (exact match + lower importance), got '%s'", m.filteredMems[1].Id)
	}

	// Last should be "Database schema design" (no match despite higher importance)
	if m.filteredMems[2].Id != "2" {
		t.Errorf("Expected last memory to be ID '2' (no content match), got '%s'", m.filteredMems[2].Id)
	}
}

func TestRelevanceScoringImportance(t *testing.T) {
	m := NewModel()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "Similar content", 3, nil, time.Hour),
		createTestMemory("2", "Similar content", 9, nil, time.Hour),
		createTestMemory("3", "Similar content", 6, nil, time.Hour),
	}
	m.SetMemories(memories, 3)

	// Set sort mode to relevance and search for "similar"
	m.SetSort(SortByRelevance, true)
	m.searchQuery = "similar"
	m.applySorting()

	// First result should be memory with importance 9
	if m.filteredMems[0].Importance != 9 {
		t.Errorf("Expected first memory importance 9, got %d", m.filteredMems[0].Importance)
	}

	// Second should have importance 6
	if m.filteredMems[1].Importance != 6 {
		t.Errorf("Expected second memory importance 6, got %d", m.filteredMems[1].Importance)
	}

	// Last should have importance 3
	if m.filteredMems[2].Importance != 3 {
		t.Errorf("Expected last memory importance 3, got %d", m.filteredMems[2].Importance)
	}
}

func TestRelevanceScoringRecency(t *testing.T) {
	m := NewModel()
	now := time.Now()
	memories := []*pb.MemoryNote{
		createTestMemory("1", "Similar content", 5, nil, 10*time.Hour), // Updated 10 hours ago
		createTestMemory("2", "Similar content", 5, nil, 1*time.Hour),  // Updated 1 hour ago (newer)
		createTestMemory("3", "Similar content", 5, nil, 5*time.Hour),  // Updated 5 hours ago
	}
	m.SetMemories(memories, 3)

	// Manually update timestamps to be more recent
	memories[0].UpdatedAt = uint64(now.Add(-10 * time.Hour).Unix())
	memories[1].UpdatedAt = uint64(now.Add(-1 * time.Hour).Unix())
	memories[2].UpdatedAt = uint64(now.Add(-5 * time.Hour).Unix())

	m.memories = memories
	m.SetMemories(memories, 3)

	// Set sort mode to relevance and search for "similar"
	m.SetSort(SortByRelevance, true)
	m.searchQuery = "similar"
	m.applySorting()

	// First result should be the most recently updated memory (1 hour ago)
	if m.filteredMems[0].Id != "2" {
		t.Errorf("Expected first memory to be ID '2' (most recent), got '%s'", m.filteredMems[0].Id)
	}

	// Second should be 5 hours ago
	if m.filteredMems[1].Id != "3" {
		t.Errorf("Expected second memory to be ID '3' (5 hours ago), got '%s'", m.filteredMems[1].Id)
	}

	// Last should be 10 hours ago
	if m.filteredMems[2].Id != "1" {
		t.Errorf("Expected last memory to be ID '1' (least recent), got '%s'", m.filteredMems[2].Id)
	}
}
