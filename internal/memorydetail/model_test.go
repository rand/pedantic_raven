package memorydetail

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// --- Update Tests ---

func TestUpdateUnfocused(t *testing.T) {
	m := NewModel()
	m.SetFocus(false)

	memory := createTestMemory("test-1", "Content", 8, nil)
	msg := MemoryLoadedMsg{Memory: memory}

	m, _ = m.Update(msg)

	// Should not process messages when unfocused
	if m.Memory() != nil {
		t.Error("Expected unfocused model to ignore messages")
	}
}

func TestUpdateMemoryLoadedMsg(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Test content", 8, []string{"tag1"})
	msg := MemoryLoadedMsg{Memory: memory}

	m, _ = m.Update(msg)

	if m.Memory() == nil {
		t.Fatal("Expected memory to be set")
	}

	if m.Memory().Id != "test-1" {
		t.Errorf("Expected memory ID 'test-1', got '%s'", m.Memory().Id)
	}
}

func TestUpdateMemoryErrorMsg(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	err := &testError{msg: "load failed"}
	msg := MemoryErrorMsg{Err: err}

	m, _ = m.Update(msg)

	if m.Error() == nil {
		t.Fatal("Expected error to be set")
	}

	if m.Error().Error() != "load failed" {
		t.Errorf("Expected error 'load failed', got '%s'", m.Error().Error())
	}
}

// --- Keyboard Handling Tests ---

func TestScrollDown(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.SetSize(80, 10)

	// Create memory with multiple lines
	content := strings.Repeat("Line of content\n", 20)
	memory := createTestMemory("test-1", content, 8, nil)
	m.SetMemory(memory)

	// Scroll down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if m.scrollOffset != 1 {
		t.Errorf("Expected scrollOffset 1 after 'j', got %d", m.scrollOffset)
	}

	// Scroll down with arrow
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	if m.scrollOffset != 2 {
		t.Errorf("Expected scrollOffset 2 after down arrow, got %d", m.scrollOffset)
	}
}

func TestScrollUp(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.SetSize(80, 10)

	content := strings.Repeat("Line\n", 20)
	memory := createTestMemory("test-1", content, 8, nil)
	m.SetMemory(memory)

	m.scrollOffset = 5

	// Scroll up
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if m.scrollOffset != 4 {
		t.Errorf("Expected scrollOffset 4 after 'k', got %d", m.scrollOffset)
	}

	// Scroll up with arrow
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})

	if m.scrollOffset != 3 {
		t.Errorf("Expected scrollOffset 3 after up arrow, got %d", m.scrollOffset)
	}
}

func TestScrollUpBoundary(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)
	m.scrollOffset = 0

	// Try to scroll up past top
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0, got %d", m.scrollOffset)
	}
}

func TestScrollDownBoundary(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.SetSize(80, 10)

	// Short content that fits in viewport
	memory := createTestMemory("test-1", "Line 1\nLine 2\nLine 3", 8, nil)
	m.SetMemory(memory)

	maxScroll := m.maxScrollOffset()

	// Try to scroll past max
	for i := 0; i < 20; i++ {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	}

	if m.scrollOffset > maxScroll {
		t.Errorf("Expected scrollOffset <= %d, got %d", maxScroll, m.scrollOffset)
	}
}

func TestScrollToTop(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)
	m.scrollOffset = 10

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})

	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset 0 after 'g', got %d", m.scrollOffset)
	}
}

func TestScrollToBottom(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.SetSize(80, 10)

	content := strings.Repeat("Line\n", 20)
	memory := createTestMemory("test-1", content, 8, nil)
	m.SetMemory(memory)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})

	maxScroll := m.maxScrollOffset()
	if m.scrollOffset != maxScroll {
		t.Errorf("Expected scrollOffset %d after 'G', got %d", maxScroll, m.scrollOffset)
	}
}

func TestPageDown(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.SetSize(80, 10)

	content := strings.Repeat("Line\n", 50)
	memory := createTestMemory("test-1", content, 8, nil)
	m.SetMemory(memory)

	visibleLines := m.visibleLines()

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

	if m.scrollOffset != visibleLines {
		t.Errorf("Expected scrollOffset %d after Ctrl+D, got %d", visibleLines, m.scrollOffset)
	}
}

func TestPageUp(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)
	m.SetSize(80, 10)

	content := strings.Repeat("Line\n", 50)
	memory := createTestMemory("test-1", content, 8, nil)
	m.SetMemory(memory)

	visibleLines := m.visibleLines()
	m.scrollOffset = visibleLines * 2

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlU})

	expectedOffset := visibleLines
	if m.scrollOffset != expectedOffset {
		t.Errorf("Expected scrollOffset %d after Ctrl+U, got %d", expectedOffset, m.scrollOffset)
	}
}

func TestToggleMetadataKey(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	// Should be shown by default
	if !m.ShowMetadata() {
		t.Error("Expected metadata shown by default")
	}

	// Press 'm' to toggle
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

	if m.ShowMetadata() {
		t.Error("Expected metadata hidden after 'm'")
	}

	// Press 'm' again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

	if !m.ShowMetadata() {
		t.Error("Expected metadata shown after second 'm'")
	}
}

func TestCloseKeys(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	// Test 'q' key
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if cmd == nil {
		t.Error("Expected non-nil cmd for 'q' key")
	}

	msg := cmd()
	if _, ok := msg.(CloseRequestMsg); !ok {
		t.Errorf("Expected CloseRequestMsg, got %T", msg)
	}

	// Test 'esc' key
	m, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Error("Expected non-nil cmd for 'esc' key")
	}

	msg = cmd()
	if _, ok := msg.(CloseRequestMsg); !ok {
		t.Errorf("Expected CloseRequestMsg, got %T", msg)
	}
}

// --- Helper Function Tests ---

func TestVisibleLines(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	visible := m.visibleLines()

	// Height 10 - header (1) - footer (1) = 8
	if visible != 8 {
		t.Errorf("Expected 8 visible lines, got %d", visible)
	}
}

func TestContentLines(t *testing.T) {
	m := NewModel()

	// Test with no memory
	if m.contentLines() != 0 {
		t.Error("Expected 0 content lines when memory is nil")
	}

	// Test with single line
	memory := createTestMemory("test-1", "Single line", 8, nil)
	m.SetMemory(memory)

	if m.contentLines() != 1 {
		t.Errorf("Expected 1 content line, got %d", m.contentLines())
	}

	// Test with multiple lines
	memory2 := createTestMemory("test-2", "Line 1\nLine 2\nLine 3", 8, nil)
	m.SetMemory(memory2)

	if m.contentLines() != 3 {
		t.Errorf("Expected 3 content lines, got %d", m.contentLines())
	}
}

func TestMaxScrollOffset(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 10)

	// No memory
	if m.maxScrollOffset() != 0 {
		t.Error("Expected 0 max scroll when memory is nil")
	}

	// Content shorter than viewport
	memory := createTestMemory("test-1", "Line 1\nLine 2", 8, nil)
	m.SetMemory(memory)

	if m.maxScrollOffset() != 0 {
		t.Errorf("Expected 0 max scroll for short content, got %d", m.maxScrollOffset())
	}

	// Content longer than viewport
	content := strings.Repeat("Line\n", 20)
	memory2 := createTestMemory("test-2", content, 8, nil)
	m.SetMemory(memory2)

	maxScroll := m.maxScrollOffset()
	visibleLines := m.visibleLines()
	totalLines := m.contentLines()
	expectedMax := totalLines - visibleLines

	if maxScroll != expectedMax {
		t.Errorf("Expected max scroll %d, got %d", expectedMax, maxScroll)
	}
}

func TestScrollWithNoMemory(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	// Try scrolling with no memory set
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	if m.scrollOffset != 0 {
		t.Error("Expected scrollOffset to remain 0 when memory is nil")
	}
}

// --- Link Navigation Keyboard Tests ---

func TestLKeySelectsFirstLink(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
	}
	m.SetMemory(memory)

	// Press 'l' to select first link
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if m.SelectedLinkIndex() != 0 {
		t.Errorf("Expected link index 0 after 'l', got %d", m.SelectedLinkIndex())
	}
}

func TestLKeyWithNoLinks(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	// Press 'l' when there are no links
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	if m.SelectedLinkIndex() != -1 {
		t.Errorf("Expected link index to remain -1, got %d", m.SelectedLinkIndex())
	}
}

func TestTabKeyNextLink(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
		{TargetId: "link-3"},
	}
	m.SetMemory(memory)

	// Press tab to select next link
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if m.SelectedLinkIndex() != 0 {
		t.Errorf("Expected link index 0 after first tab, got %d", m.SelectedLinkIndex())
	}

	// Press tab again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if m.SelectedLinkIndex() != 1 {
		t.Errorf("Expected link index 1 after second tab, got %d", m.SelectedLinkIndex())
	}
}

func TestNKeyNextLink(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
	}
	m.SetMemory(memory)

	// Press 'n' to select next link
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if m.SelectedLinkIndex() != 0 {
		t.Errorf("Expected link index 0 after 'n', got %d", m.SelectedLinkIndex())
	}
}

func TestPKeyPreviousLink(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
		{TargetId: "link-2"},
	}
	m.SetMemory(memory)

	m.selectedLinkIndex = 1

	// Press 'p' to select previous link
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})

	if m.SelectedLinkIndex() != 0 {
		t.Errorf("Expected link index 0 after 'p', got %d", m.SelectedLinkIndex())
	}
}

func TestEscClearsLinkSelection(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "link-1"},
	}
	m.SetMemory(memory)

	m.selectedLinkIndex = 0

	// Press Esc to clear link selection
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if m.SelectedLinkIndex() != -1 {
		t.Errorf("Expected link selection cleared, got index %d", m.SelectedLinkIndex())
	}

	// Should not close the view
	if cmd != nil {
		t.Error("Expected nil cmd when clearing link selection")
	}
}

func TestEscClosesWhenNoLinkSelected(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	// Press Esc when no link is selected
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Error("Expected non-nil cmd to close view")
	}

	msg := cmd()
	if _, ok := msg.(CloseRequestMsg); !ok {
		t.Errorf("Expected CloseRequestMsg, got %T", msg)
	}
}

func TestEnterNavigatesToSelectedLink(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	memory.Links = []*pb.MemoryLink{
		{TargetId: "target-memory-123"},
	}
	m.SetMemory(memory)

	m.selectedLinkIndex = 0

	// Press Enter to navigate
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Error("Expected non-nil cmd for link navigation")
	}

	msg := cmd()
	linkMsg, ok := msg.(LinkSelectedMsg)
	if !ok {
		t.Fatalf("Expected LinkSelectedMsg, got %T", msg)
	}

	if linkMsg.TargetID != "target-memory-123" {
		t.Errorf("Expected target ID 'target-memory-123', got '%s'", linkMsg.TargetID)
	}
}

func TestEnterWithNoLinkSelected(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	// Press Enter when no link is selected
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd != nil {
		t.Error("Expected nil cmd when no link is selected")
	}
}

// --- Edit Confirmation Dialog Tests ---

func TestConfirmationDialogYes(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	// Create a memory and set it
	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	// Enter edit mode
	m.editState = &EditState{
		isEditing:    true,
		editedMemory: memory,
		fieldFocus:   FieldContent,
		hasChanges:   true,
	}

	// Press Esc to trigger confirmation dialog
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if !m.showEditConfirm {
		t.Error("Expected showEditConfirm to be true after Esc with unsaved changes")
	}

	// Press 'y' to confirm cancellation
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	if m.showEditConfirm {
		t.Error("Expected showEditConfirm to be false after 'y' confirmation")
	}

	if m.IsEditing() {
		t.Error("Expected editing to be cancelled after 'y' confirmation")
	}

	if cmd != nil {
		t.Error("Expected nil cmd when confirming cancellation")
	}
}

func TestConfirmationDialogNo(t *testing.T) {
	m := NewModel()
	m.SetFocus(true)

	// Create a memory and set it
	memory := createTestMemory("test-1", "Content", 8, nil)
	m.SetMemory(memory)

	// Enter edit mode
	m.editState = &EditState{
		isEditing:    true,
		editedMemory: memory,
		fieldFocus:   FieldContent,
		hasChanges:   true,
	}

	// Press Esc to trigger confirmation dialog
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if !m.showEditConfirm {
		t.Error("Expected showEditConfirm to be true after Esc with unsaved changes")
	}

	// Press 'n' to cancel the cancellation
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	if m.showEditConfirm {
		t.Error("Expected showEditConfirm to be false after 'n' response")
	}

	if !m.IsEditing() {
		t.Error("Expected editing to continue after 'n' response")
	}

	if cmd != nil {
		t.Error("Expected nil cmd when cancelling the cancellation")
	}
}
