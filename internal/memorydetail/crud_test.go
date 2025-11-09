package memorydetail

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Mock client interface for CRUD operations
type crudClient interface {
	StoreMemory(context.Context, mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error)
	UpdateMemory(context.Context, mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error)
	DeleteMemory(context.Context, string) error
}

// mockCRUDClient implements crudClient for testing
type mockCRUDClient struct {
	storeFunc  func(context.Context, mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error)
	updateFunc func(context.Context, mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error)
	deleteFunc func(context.Context, string) error
}

func (m *mockCRUDClient) StoreMemory(ctx context.Context, opts mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error) {
	if m.storeFunc != nil {
		return m.storeFunc(ctx, opts)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCRUDClient) UpdateMemory(ctx context.Context, opts mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, opts)
	}
	return nil, errors.New("not implemented")
}

func (m *mockCRUDClient) DeleteMemory(ctx context.Context, memoryID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, memoryID)
	}
	return errors.New("not implemented")
}

// Helper function to create a simple test memory for CRUD operations
func createCRUDTestMemory() *pb.MemoryNote {
	return &pb.MemoryNote{
		Id:         "test-memory-1",
		Content:    "Test content",
		Importance: 5,
		Tags:       []string{"test", "memory"},
		Namespace: &pb.Namespace{
			Namespace: &pb.Namespace_Global{
				Global: &pb.GlobalNamespace{},
			},
		},
	}
}

// ============================================================================
// Edit Mode Tests
// ============================================================================

func TestEnterEditMode(t *testing.T) {
	memory := createCRUDTestMemory()
	cmd := EnterEditMode(memory)

	if cmd == nil {
		t.Fatal("EnterEditMode returned nil command")
	}

	msg := cmd()
	editMsg, ok := msg.(EditModeEnteredMsg)
	if !ok {
		t.Fatalf("Expected EditModeEnteredMsg, got %T", msg)
	}

	if editMsg.Memory == nil {
		t.Fatal("EditModeEnteredMsg has nil memory")
	}

	if editMsg.Memory.Id != memory.Id {
		t.Errorf("Expected memory ID %s, got %s", memory.Id, editMsg.Memory.Id)
	}

	// Verify it's a copy, not the same instance
	if editMsg.Memory == memory {
		t.Error("EditModeEnteredMsg should contain a copy, not the original memory")
	}
}

func TestExitEditMode(t *testing.T) {
	m := NewModel()
	memory := createCRUDTestMemory()

	// Enter edit mode
	m.editState = &EditState{
		isEditing:    true,
		editedMemory: memory,
		originalHash: hashMemory(memory),
	}

	if !m.IsEditing() {
		t.Error("Model should be in edit mode")
	}

	// Exit edit mode
	m.CancelEdit()

	if m.IsEditing() {
		t.Error("Model should not be in edit mode after cancel")
	}

	if m.editState != nil {
		t.Error("Edit state should be nil after cancel")
	}
}

func TestCancelEdit(t *testing.T) {
	m := NewModel()
	memory := createCRUDTestMemory()

	m.editState = &EditState{
		isEditing:    true,
		editedMemory: memory,
		originalHash: hashMemory(memory),
	}
	m.isNewMemory = true

	m.CancelEdit()

	if m.editState != nil {
		t.Error("Edit state should be nil after cancel")
	}

	if m.isNewMemory {
		t.Error("isNewMemory should be false after cancel")
	}
}

func TestEditModeWithoutMemory(t *testing.T) {
	cmd := EnterEditMode(nil)

	if cmd == nil {
		t.Fatal("EnterEditMode returned nil command")
	}

	msg := cmd()
	editMsg, ok := msg.(EditModeEnteredMsg)
	if !ok {
		t.Fatalf("Expected EditModeEnteredMsg, got %T", msg)
	}

	if editMsg.Memory != nil {
		t.Error("EditModeEnteredMsg should have nil memory when called with nil")
	}
}

func TestFieldFocusCycling(t *testing.T) {
	m := NewModel()
	memory := createCRUDTestMemory()

	m.editState = &EditState{
		isEditing:    true,
		editedMemory: memory,
		fieldFocus:   FieldContent,
		originalHash: hashMemory(memory),
	}

	// Test cycling through all fields
	tests := []struct {
		current  EditField
		expected EditField
	}{
		{FieldContent, FieldTags},
		{FieldTags, FieldImportance},
		{FieldImportance, FieldNamespace},
		{FieldNamespace, FieldContent}, // Wraps around
	}

	for _, tt := range tests {
		m.editState.fieldFocus = tt.current
		m.CycleFieldFocus()

		if m.editState.fieldFocus != tt.expected {
			t.Errorf("After cycling from %v, expected %v, got %v",
				tt.current, tt.expected, m.editState.fieldFocus)
		}
	}
}

func TestEditModeKeybindings(t *testing.T) {
	m := NewModel()
	memory := createCRUDTestMemory()
	m.memory = memory
	m.SetFocus(true)

	// Enter edit mode
	cmd := m.EnterEditMode()
	msg := cmd()
	m, _ = m.Update(msg)

	if !m.IsEditing() {
		t.Fatal("Model should be in edit mode")
	}

	// Test Tab key (cycle field focus)
	initialFocus := m.GetFieldFocus()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if m.GetFieldFocus() == initialFocus {
		t.Error("Tab key should cycle field focus")
	}

	// Test Esc key (cancel edit)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.IsEditing() {
		t.Error("Esc key should cancel edit mode")
	}
}

func TestUnsavedChangesDetection(t *testing.T) {
	m := NewModel()
	memory := createCRUDTestMemory()

	// Enter edit mode
	m.editState = &EditState{
		isEditing:    true,
		editedMemory: cloneMemory(memory),
		originalHash: hashMemory(memory),
	}

	// Initially no changes
	if m.HasUnsavedChanges() {
		t.Error("Should have no unsaved changes initially")
	}

	// Make a change
	m.SetEditedContent("Modified content")

	// Should detect changes
	if !m.HasUnsavedChanges() {
		t.Error("Should detect unsaved changes after modification")
	}
}

func TestEditStateInitialization(t *testing.T) {
	memory := createCRUDTestMemory()
	cmd := EnterEditMode(memory)
	msg := cmd()

	editMsg := msg.(EditModeEnteredMsg)

	// Create edit state as the Update handler would
	state := &EditState{
		isEditing:    true,
		editedMemory: editMsg.Memory,
		fieldFocus:   FieldContent,
		originalHash: hashMemory(editMsg.Memory),
		hasChanges:   false,
	}

	if !state.isEditing {
		t.Error("Edit state should be in editing mode")
	}

	if state.fieldFocus != FieldContent {
		t.Error("Initial field focus should be FieldContent")
	}

	if state.originalHash == "" {
		t.Error("Original hash should be set")
	}

	if state.hasChanges {
		t.Error("Should not have changes initially")
	}
}

// ============================================================================
// CRUD Operation Tests
// ============================================================================

func TestSaveChangesSuccess(t *testing.T) {
	memory := createCRUDTestMemory()
	savedMemory := cloneMemory(memory)
	savedMemory.Content = "Updated content"

	client := &mockCRUDClient{
		updateFunc: func(ctx context.Context, opts mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error) {
			return savedMemory, nil
		},
	}

	cmd := SaveChanges(client, memory, false)
	msg := cmd()

	saveMsg, ok := msg.(MemorySavedMsg)
	if !ok {
		t.Fatalf("Expected MemorySavedMsg, got %T", msg)
	}

	if saveMsg.Err != nil {
		t.Errorf("Expected no error, got %v", saveMsg.Err)
	}

	if saveMsg.Memory == nil {
		t.Fatal("Saved memory should not be nil")
	}

	if saveMsg.Memory.Content != "Updated content" {
		t.Errorf("Expected updated content, got %s", saveMsg.Memory.Content)
	}
}

func TestSaveChangesError(t *testing.T) {
	memory := createCRUDTestMemory()
	expectedErr := errors.New("server error")

	client := &mockCRUDClient{
		updateFunc: func(ctx context.Context, opts mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error) {
			return nil, expectedErr
		},
	}

	cmd := SaveChanges(client, memory, false)
	msg := cmd()

	saveMsg := msg.(MemorySavedMsg)
	if saveMsg.Err == nil {
		t.Error("Expected error, got nil")
	}

	if saveMsg.Memory != nil {
		t.Error("Memory should be nil on error")
	}
}

func TestCreateMemorySuccess(t *testing.T) {
	memory := createCRUDTestMemory()
	memory.Id = "" // New memory has no ID yet

	createdMemory := cloneMemory(memory)
	createdMemory.Id = "new-memory-id"

	client := &mockCRUDClient{
		storeFunc: func(ctx context.Context, opts mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error) {
			return createdMemory, nil
		},
	}

	cmd := CreateMemory(client, memory)
	msg := cmd()

	createMsg, ok := msg.(MemoryCreatedMsg)
	if !ok {
		t.Fatalf("Expected MemoryCreatedMsg, got %T", msg)
	}

	if createMsg.Err != nil {
		t.Errorf("Expected no error, got %v", createMsg.Err)
	}

	if createMsg.Memory == nil {
		t.Fatal("Created memory should not be nil")
	}

	if createMsg.Memory.Id != "new-memory-id" {
		t.Errorf("Expected new memory ID, got %s", createMsg.Memory.Id)
	}
}

func TestCreateMemoryError(t *testing.T) {
	memory := createCRUDTestMemory()
	expectedErr := errors.New("creation failed")

	client := &mockCRUDClient{
		storeFunc: func(ctx context.Context, opts mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error) {
			return nil, expectedErr
		},
	}

	cmd := CreateMemory(client, memory)
	msg := cmd()

	createMsg := msg.(MemoryCreatedMsg)
	if createMsg.Err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestUpdateMemorySuccess(t *testing.T) {
	memory := createCRUDTestMemory()
	updatedMemory := cloneMemory(memory)
	updatedMemory.Content = "Updated"

	client := &mockCRUDClient{
		updateFunc: func(ctx context.Context, opts mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error) {
			return updatedMemory, nil
		},
	}

	cmd := UpdateMemory(client, memory)
	msg := cmd()

	updateMsg := msg.(MemoryUpdatedMsg)
	if updateMsg.Err != nil {
		t.Errorf("Expected no error, got %v", updateMsg.Err)
	}

	if updateMsg.Memory.Content != "Updated" {
		t.Error("Memory content should be updated")
	}
}

func TestUpdateMemoryError(t *testing.T) {
	memory := createCRUDTestMemory()
	expectedErr := errors.New("update failed")

	client := &mockCRUDClient{
		updateFunc: func(ctx context.Context, opts mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error) {
			return nil, expectedErr
		},
	}

	cmd := UpdateMemory(client, memory)
	msg := cmd()

	updateMsg := msg.(MemoryUpdatedMsg)
	if updateMsg.Err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDeleteMemorySuccess(t *testing.T) {
	client := &mockCRUDClient{
		deleteFunc: func(ctx context.Context, memoryID string) error {
			if memoryID != "test-memory-1" {
				t.Errorf("Expected memory ID test-memory-1, got %s", memoryID)
			}
			return nil
		},
	}

	cmd := DeleteMemory(client, "test-memory-1")
	msg := cmd()

	deleteMsg, ok := msg.(MemoryDeletedMsg)
	if !ok {
		t.Fatalf("Expected MemoryDeletedMsg, got %T", msg)
	}

	if deleteMsg.Err != nil {
		t.Errorf("Expected no error, got %v", deleteMsg.Err)
	}

	if deleteMsg.MemoryID != "test-memory-1" {
		t.Errorf("Expected memory ID test-memory-1, got %s", deleteMsg.MemoryID)
	}
}

func TestDeleteMemoryError(t *testing.T) {
	expectedErr := errors.New("delete failed")

	client := &mockCRUDClient{
		deleteFunc: func(ctx context.Context, memoryID string) error {
			return expectedErr
		},
	}

	cmd := DeleteMemory(client, "test-memory-1")
	msg := cmd()

	deleteMsg := msg.(MemoryDeletedMsg)
	if deleteMsg.Err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDeleteWithoutClient(t *testing.T) {
	cmd := DeleteMemory(nil, "test-memory-1")
	msg := cmd()

	deleteMsg := msg.(MemoryDeletedMsg)
	if deleteMsg.Err != ErrNoClient {
		t.Errorf("Expected ErrNoClient, got %v", deleteMsg.Err)
	}
}

func TestSaveWithoutClient(t *testing.T) {
	memory := createCRUDTestMemory()
	cmd := SaveChanges(nil, memory, false)
	msg := cmd()

	saveMsg := msg.(MemorySavedMsg)
	if saveMsg.Err != ErrNoClient {
		t.Errorf("Expected ErrNoClient, got %v", saveMsg.Err)
	}
}

// ============================================================================
// Validation Tests
// ============================================================================

func TestValidateMemoryAllValid(t *testing.T) {
	memory := createCRUDTestMemory()
	err := validateMemory(memory)

	if err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestValidateMissingContent(t *testing.T) {
	memory := createCRUDTestMemory()
	memory.Content = ""

	err := validateMemory(memory)
	if err != ErrContentRequired {
		t.Errorf("Expected ErrContentRequired, got %v", err)
	}

	// Test whitespace-only content
	memory.Content = "   \n\t  "
	err = validateMemory(memory)
	if err != ErrContentRequired {
		t.Errorf("Expected ErrContentRequired for whitespace, got %v", err)
	}
}

func TestValidateContentTooLong(t *testing.T) {
	memory := createCRUDTestMemory()
	memory.Content = strings.Repeat("a", MaxContentLength+1)

	err := validateMemory(memory)
	if err != ErrContentTooLong {
		t.Errorf("Expected ErrContentTooLong, got %v", err)
	}
}

func TestValidateInvalidImportance(t *testing.T) {
	tests := []struct {
		name       string
		importance uint32
	}{
		{"too low", 0},
		{"too high", 11},
		{"way too high", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memory := createCRUDTestMemory()
			memory.Importance = tt.importance

			err := validateMemory(memory)
			if err != ErrImportanceInvalid {
				t.Errorf("Expected ErrImportanceInvalid for importance %d, got %v",
					tt.importance, err)
			}
		})
	}

	// Test valid importance values
	for i := uint32(MinImportance); i <= uint32(MaxImportance); i++ {
		memory := createCRUDTestMemory()
		memory.Importance = i

		err := validateMemory(memory)
		if err != nil {
			t.Errorf("Importance %d should be valid, got error: %v", i, err)
		}
	}
}

func TestValidateTooManyTags(t *testing.T) {
	memory := createCRUDTestMemory()
	memory.Tags = make([]string, MaxTags+1)
	for i := range memory.Tags {
		memory.Tags[i] = "tag"
	}

	err := validateMemory(memory)
	if err != ErrTooManyTags {
		t.Errorf("Expected ErrTooManyTags, got %v", err)
	}

	// Test exactly at limit
	memory.Tags = make([]string, MaxTags)
	for i := range memory.Tags {
		memory.Tags[i] = "tag"
	}

	err = validateMemory(memory)
	if err != nil {
		t.Errorf("Expected no error for %d tags, got %v", MaxTags, err)
	}
}

func TestValidateNilMemory(t *testing.T) {
	err := validateMemory(nil)
	if err != ErrContentRequired {
		t.Errorf("Expected ErrContentRequired for nil memory, got %v", err)
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestHashMemory(t *testing.T) {
	memory := createCRUDTestMemory()

	hash1 := hashMemory(memory)
	if hash1 == "" {
		t.Error("Hash should not be empty")
	}

	// Same memory should produce same hash
	hash2 := hashMemory(memory)
	if hash1 != hash2 {
		t.Error("Same memory should produce same hash")
	}

	// Different content should produce different hash
	memory2 := cloneMemory(memory)
	memory2.Content = "Different content"
	hash3 := hashMemory(memory2)
	if hash1 == hash3 {
		t.Error("Different content should produce different hash")
	}

	// Nil memory should produce empty hash
	nilHash := hashMemory(nil)
	if nilHash != "" {
		t.Error("Nil memory should produce empty hash")
	}
}

func TestCloneMemory(t *testing.T) {
	original := createCRUDTestMemory()
	original.Links = []*pb.MemoryLink{
		{
			TargetId: "link-1",
			LinkType: pb.LinkType_LINK_TYPE_EXTENDS,
			Strength: 0.8,
		},
	}

	clone := cloneMemory(original)

	if clone == nil {
		t.Fatal("Clone should not be nil")
	}

	// Verify it's a different instance
	if clone == original {
		t.Error("Clone should be a different instance")
	}

	// Verify fields are copied
	if clone.Id != original.Id {
		t.Error("Clone ID should match original")
	}

	if clone.Content != original.Content {
		t.Error("Clone content should match original")
	}

	if clone.Importance != original.Importance {
		t.Error("Clone importance should match original")
	}

	// Verify deep copy of tags
	if &clone.Tags == &original.Tags {
		t.Error("Tags should be a deep copy")
	}

	if len(clone.Tags) != len(original.Tags) {
		t.Error("Tags length should match")
	}

	// Verify deep copy of links
	if len(clone.Links) != len(original.Links) {
		t.Error("Links length should match")
	}

	if clone.Links[0] == original.Links[0] {
		t.Error("Links should be deep copied")
	}

	// Modify clone and verify original is unchanged
	clone.Content = "Modified"
	clone.Tags[0] = "modified-tag"

	if original.Content == "Modified" {
		t.Error("Modifying clone should not affect original")
	}

	if original.Tags[0] == "modified-tag" {
		t.Error("Modifying clone tags should not affect original")
	}
}

func TestCloneMemoryNil(t *testing.T) {
	clone := cloneMemory(nil)
	if clone != nil {
		t.Error("Cloning nil memory should return nil")
	}
}

func TestDetectChanges(t *testing.T) {
	memory := createCRUDTestMemory()

	state := &EditState{
		isEditing:    true,
		editedMemory: cloneMemory(memory),
		originalHash: hashMemory(memory),
	}

	// No changes initially
	if state.detectChanges() {
		t.Error("Should not detect changes initially")
	}

	// Modify content
	state.editedMemory.Content = "Modified"
	if !state.detectChanges() {
		t.Error("Should detect changes after modification")
	}

	// Modify tags
	state.editedMemory = cloneMemory(memory)
	state.editedMemory.Tags = append(state.editedMemory.Tags, "new-tag")
	if !state.detectChanges() {
		t.Error("Should detect changes after tag modification")
	}

	// Modify importance
	state.editedMemory = cloneMemory(memory)
	state.editedMemory.Importance = 10
	if !state.detectChanges() {
		t.Error("Should detect changes after importance modification")
	}
}
