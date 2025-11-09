package memorydetail

import (
	"context"
	"errors"
	"testing"

	"github.com/rand/pedantic-raven/internal/mnemosyne"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// mockMemoryClient implements MemoryClient for testing
type mockMemoryClient struct {
	storeFunc  func(context.Context, mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error)
	updateFunc func(context.Context, mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error)
	deleteFunc func(context.Context, string) error
}

func (m *mockMemoryClient) StoreMemory(ctx context.Context, opts mnemosyne.StoreMemoryOptions) (*pb.MemoryNote, error) {
	if m.storeFunc != nil {
		return m.storeFunc(ctx, opts)
	}
	return nil, errors.New("not implemented")
}

func (m *mockMemoryClient) UpdateMemory(ctx context.Context, opts mnemosyne.UpdateMemoryOptions) (*pb.MemoryNote, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, opts)
	}
	return nil, errors.New("not implemented")
}

func (m *mockMemoryClient) DeleteMemory(ctx context.Context, memoryID string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, memoryID)
	}
	return errors.New("not implemented")
}

// TestNavigationHistory tests the navigation history functionality.
func TestNavigationHistory(t *testing.T) {
	nh := NewNavigationHistory()

	// Test initial state
	if nh.CanGoBack() {
		t.Error("Should not be able to go back initially")
	}
	if nh.CanGoForward() {
		t.Error("Should not be able to go forward initially")
	}
	if nh.Current() != "" {
		t.Error("Current should be empty initially")
	}

	// Test pushing entries
	nh.Push("memory1")
	// After first push, current is 0, so we can't go back yet
	if nh.CanGoBack() {
		t.Error("Should not be able to go back after first push (at index 0)")
	}
	if nh.Current() != "memory1" {
		t.Errorf("Expected current to be 'memory1', got '%s'", nh.Current())
	}

	nh.Push("memory2")
	// After second push, current is 1, so we can go back to 0
	if !nh.CanGoBack() {
		t.Error("Should be able to go back after second push")
	}
	if nh.Current() != "memory2" {
		t.Errorf("Expected current to be 'memory2', got '%s'", nh.Current())
	}

	nh.Push("memory3")
	if nh.Current() != "memory3" {
		t.Errorf("Expected current to be 'memory3', got '%s'", nh.Current())
	}
}

func TestNavigationHistoryBack(t *testing.T) {
	nh := NewNavigationHistory()
	nh.Push("memory1")
	nh.Push("memory2")
	nh.Push("memory3")

	// Go back
	id, ok := nh.Back()
	if !ok {
		t.Error("Should be able to go back")
	}
	if id != "memory2" {
		t.Errorf("Expected 'memory2', got '%s'", id)
	}

	// Go back again
	id, ok = nh.Back()
	if !ok {
		t.Error("Should be able to go back again")
	}
	if id != "memory1" {
		t.Errorf("Expected 'memory1', got '%s'", id)
	}

	// Try to go back beyond start
	_, ok = nh.Back()
	if ok {
		t.Error("Should not be able to go back beyond start")
	}
}

func TestNavigationHistoryForward(t *testing.T) {
	nh := NewNavigationHistory()
	nh.Push("memory1")
	nh.Push("memory2")
	nh.Push("memory3")

	// Go back twice
	nh.Back()
	nh.Back()

	// Go forward
	id, ok := nh.Forward()
	if !ok {
		t.Error("Should be able to go forward")
	}
	if id != "memory2" {
		t.Errorf("Expected 'memory2', got '%s'", id)
	}

	// Go forward again
	id, ok = nh.Forward()
	if !ok {
		t.Error("Should be able to go forward again")
	}
	if id != "memory3" {
		t.Errorf("Expected 'memory3', got '%s'", id)
	}

	// Try to go forward beyond end
	_, ok = nh.Forward()
	if ok {
		t.Error("Should not be able to go forward beyond end")
	}
}

func TestNavigationHistoryMaxSize(t *testing.T) {
	nh := NewNavigationHistory()
	nh.maxSize = 3 // Set small max size for testing

	// Push 5 entries
	nh.Push("memory1")
	nh.Push("memory2")
	nh.Push("memory3")
	nh.Push("memory4")
	nh.Push("memory5")

	// Should only keep last 3
	if len(nh.history) != 3 {
		t.Errorf("Expected history length 3, got %d", len(nh.history))
	}

	if nh.Current() != "memory5" {
		t.Errorf("Expected current to be 'memory5', got '%s'", nh.Current())
	}

	// Go back twice should give us memory4 and memory3
	id, _ := nh.Back()
	if id != "memory4" {
		t.Errorf("Expected 'memory4', got '%s'", id)
	}

	id, _ = nh.Back()
	if id != "memory3" {
		t.Errorf("Expected 'memory3', got '%s'", id)
	}
}

func TestNavigationHistoryTruncate(t *testing.T) {
	nh := NewNavigationHistory()
	nh.Push("memory1")
	nh.Push("memory2")
	nh.Push("memory3")

	// Go back
	nh.Back()
	nh.Back()

	// Push new entry - should truncate forward history
	nh.Push("memory4")

	// Should not be able to go forward
	if nh.CanGoForward() {
		t.Error("Should not be able to go forward after push")
	}

	// History should be memory1, memory4
	if len(nh.history) != 2 {
		t.Errorf("Expected history length 2, got %d", len(nh.history))
	}
}

func TestNavigationHistoryClear(t *testing.T) {
	nh := NewNavigationHistory()
	nh.Push("memory1")
	nh.Push("memory2")

	nh.Clear()

	if len(nh.history) != 0 {
		t.Error("History should be empty after clear")
	}
	if nh.current != -1 {
		t.Error("Current should be -1 after clear")
	}
}

func TestNavigationHistoryBoundaries(t *testing.T) {
	nh := NewNavigationHistory()

	// Test with empty history
	if nh.CanGoBack() {
		t.Error("Should not be able to go back with empty history")
	}
	if nh.CanGoForward() {
		t.Error("Should not be able to go forward with empty history")
	}

	_, ok := nh.Back()
	if ok {
		t.Error("Back should fail with empty history")
	}

	_, ok = nh.Forward()
	if ok {
		t.Error("Forward should fail with empty history")
	}
}

// TestCreateLink tests link creation.
func TestCreateLink(t *testing.T) {
	client := &mockMemoryClient{}
	sourceID := "source-123"
	targetID := "target-456"
	linkType := pb.LinkType_LINK_TYPE_REFERENCES
	strength := float32(0.8)
	reason := "test link"

	cmd := CreateLink(client, sourceID, targetID, linkType, strength, reason)
	msg := cmd()

	linkMsg, ok := msg.(LinkCreatedMsg)
	if !ok {
		t.Fatalf("Expected LinkCreatedMsg, got %T", msg)
	}

	if linkMsg.Err != nil {
		t.Errorf("Expected no error, got %v", linkMsg.Err)
	}

	if linkMsg.Link == nil {
		t.Fatal("Expected link to be created")
	}

	if linkMsg.Link.TargetId != targetID {
		t.Errorf("Expected target ID %s, got %s", targetID, linkMsg.Link.TargetId)
	}

	if linkMsg.Link.LinkType != linkType {
		t.Errorf("Expected link type %v, got %v", linkType, linkMsg.Link.LinkType)
	}

	if linkMsg.Link.Strength != strength {
		t.Errorf("Expected strength %f, got %f", strength, linkMsg.Link.Strength)
	}

	if !linkMsg.Link.UserCreated {
		t.Error("Expected user_created to be true")
	}
}

func TestCreateLinkWithoutClient(t *testing.T) {
	cmd := CreateLink(nil, "source", "target", pb.LinkType_LINK_TYPE_REFERENCES, 0.5, "")
	msg := cmd()

	linkMsg, ok := msg.(LinkCreatedMsg)
	if !ok {
		t.Fatalf("Expected LinkCreatedMsg, got %T", msg)
	}

	if linkMsg.Err != ErrNoClient {
		t.Errorf("Expected ErrNoClient, got %v", linkMsg.Err)
	}
}

func TestCreateLinkInvalidInputs(t *testing.T) {
	client := &mockMemoryClient{}

	tests := []struct {
		name      string
		sourceID  string
		targetID  string
		linkType  pb.LinkType
		strength  float32
		expectErr error
	}{
		{
			name:      "empty source",
			sourceID:  "",
			targetID:  "target",
			linkType:  pb.LinkType_LINK_TYPE_REFERENCES,
			strength:  0.5,
			expectErr: ErrSourceNotFound,
		},
		{
			name:      "empty target",
			sourceID:  "source",
			targetID:  "",
			linkType:  pb.LinkType_LINK_TYPE_REFERENCES,
			strength:  0.5,
			expectErr: ErrTargetNotFound,
		},
		{
			name:      "self reference",
			sourceID:  "same",
			targetID:  "same",
			linkType:  pb.LinkType_LINK_TYPE_REFERENCES,
			strength:  0.5,
			expectErr: ErrLinkToSelf,
		},
		{
			name:      "strength too low",
			sourceID:  "source",
			targetID:  "target",
			linkType:  pb.LinkType_LINK_TYPE_REFERENCES,
			strength:  -0.1,
			expectErr: ErrInvalidStrength,
		},
		{
			name:      "strength too high",
			sourceID:  "source",
			targetID:  "target",
			linkType:  pb.LinkType_LINK_TYPE_REFERENCES,
			strength:  1.5,
			expectErr: ErrInvalidStrength,
		},
		{
			name:      "unspecified link type",
			sourceID:  "source",
			targetID:  "target",
			linkType:  pb.LinkType_LINK_TYPE_UNSPECIFIED,
			strength:  0.5,
			expectErr: ErrInvalidLinkType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateLink(client, tt.sourceID, tt.targetID, tt.linkType, tt.strength, "")
			msg := cmd()

			linkMsg, ok := msg.(LinkCreatedMsg)
			if !ok {
				t.Fatalf("Expected LinkCreatedMsg, got %T", msg)
			}

			if !errors.Is(linkMsg.Err, tt.expectErr) {
				t.Errorf("Expected error %v, got %v", tt.expectErr, linkMsg.Err)
			}
		})
	}
}

func TestDeleteLink(t *testing.T) {
	client := &mockMemoryClient{}
	sourceID := "source-123"
	targetID := "target-456"

	cmd := DeleteLink(client, sourceID, targetID)
	msg := cmd()

	delMsg, ok := msg.(LinkDeletedMsg)
	if !ok {
		t.Fatalf("Expected LinkDeletedMsg, got %T", msg)
	}

	if delMsg.Err != nil {
		t.Errorf("Expected no error, got %v", delMsg.Err)
	}

	if delMsg.LinkID != targetID {
		t.Errorf("Expected link ID %s, got %s", targetID, delMsg.LinkID)
	}
}

func TestDeleteLinkWithoutClient(t *testing.T) {
	cmd := DeleteLink(nil, "source", "target")
	msg := cmd()

	delMsg, ok := msg.(LinkDeletedMsg)
	if !ok {
		t.Fatalf("Expected LinkDeletedMsg, got %T", msg)
	}

	if delMsg.Err != ErrNoClient {
		t.Errorf("Expected ErrNoClient, got %v", delMsg.Err)
	}
}

func TestDeleteLinkInvalidInputs(t *testing.T) {
	client := &mockMemoryClient{}

	tests := []struct {
		name     string
		sourceID string
		targetID string
	}{
		{"empty source", "", "target"},
		{"empty target", "source", ""},
		{"both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := DeleteLink(client, tt.sourceID, tt.targetID)
			msg := cmd()

			delMsg, ok := msg.(LinkDeletedMsg)
			if !ok {
				t.Fatalf("Expected LinkDeletedMsg, got %T", msg)
			}

			if delMsg.Err != ErrLinkNotFound {
				t.Errorf("Expected ErrLinkNotFound, got %v", delMsg.Err)
			}
		})
	}
}

func TestUpdateLinkMetadata(t *testing.T) {
	client := &mockMemoryClient{}
	sourceID := "source-123"
	targetID := "target-456"
	strength := float32(0.9)
	reason := "updated reason"

	cmd := UpdateLinkMetadata(client, sourceID, targetID, &strength, &reason)
	msg := cmd()

	updateMsg, ok := msg.(LinkMetadataUpdatedMsg)
	if !ok {
		t.Fatalf("Expected LinkMetadataUpdatedMsg, got %T", msg)
	}

	if updateMsg.Err != nil {
		t.Errorf("Expected no error, got %v", updateMsg.Err)
	}

	if updateMsg.Link == nil {
		t.Fatal("Expected link to be updated")
	}

	if updateMsg.Link.Strength != strength {
		t.Errorf("Expected strength %f, got %f", strength, updateMsg.Link.Strength)
	}

	if updateMsg.Link.Reason != reason {
		t.Errorf("Expected reason %s, got %s", reason, updateMsg.Link.Reason)
	}
}

func TestUpdateLinkMetadataWithoutClient(t *testing.T) {
	strength := float32(0.5)
	cmd := UpdateLinkMetadata(nil, "source", "target", &strength, nil)
	msg := cmd()

	updateMsg, ok := msg.(LinkMetadataUpdatedMsg)
	if !ok {
		t.Fatalf("Expected LinkMetadataUpdatedMsg, got %T", msg)
	}

	if updateMsg.Err != ErrNoClient {
		t.Errorf("Expected ErrNoClient, got %v", updateMsg.Err)
	}
}

func TestUpdateLinkMetadataInvalidStrength(t *testing.T) {
	client := &mockMemoryClient{}

	tests := []struct {
		name     string
		strength float32
	}{
		{"strength too low", -0.1},
		{"strength too high", 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := UpdateLinkMetadata(client, "source", "target", &tt.strength, nil)
			msg := cmd()

			updateMsg, ok := msg.(LinkMetadataUpdatedMsg)
			if !ok {
				t.Fatalf("Expected LinkMetadataUpdatedMsg, got %T", msg)
			}

			if updateMsg.Err != ErrInvalidStrength {
				t.Errorf("Expected ErrInvalidStrength, got %v", updateMsg.Err)
			}
		})
	}
}

func TestGetLinkedMemoriesOutbound(t *testing.T) {
	client := &mockMemoryClient{}
	memoryID := "memory-123"

	cmd := GetLinkedMemories(client, memoryID, DirectionOutbound)
	msg := cmd()

	loadedMsg, ok := msg.(LinkedMemoriesLoadedMsg)
	if !ok {
		t.Fatalf("Expected LinkedMemoriesLoadedMsg, got %T", msg)
	}

	if loadedMsg.Err != nil {
		t.Errorf("Expected no error, got %v", loadedMsg.Err)
	}

	if loadedMsg.Direction != DirectionOutbound {
		t.Errorf("Expected direction %v, got %v", DirectionOutbound, loadedMsg.Direction)
	}
}

func TestGetLinkedMemoriesInbound(t *testing.T) {
	client := &mockMemoryClient{}
	memoryID := "memory-123"

	cmd := GetLinkedMemories(client, memoryID, DirectionInbound)
	msg := cmd()

	loadedMsg, ok := msg.(LinkedMemoriesLoadedMsg)
	if !ok {
		t.Fatalf("Expected LinkedMemoriesLoadedMsg, got %T", msg)
	}

	if loadedMsg.Err != nil {
		t.Errorf("Expected no error, got %v", loadedMsg.Err)
	}

	if loadedMsg.Direction != DirectionInbound {
		t.Errorf("Expected direction %v, got %v", DirectionInbound, loadedMsg.Direction)
	}
}

func TestGetLinkedMemoriesWithoutClient(t *testing.T) {
	cmd := GetLinkedMemories(nil, "memory-123", DirectionBoth)
	msg := cmd()

	loadedMsg, ok := msg.(LinkedMemoriesLoadedMsg)
	if !ok {
		t.Fatalf("Expected LinkedMemoriesLoadedMsg, got %T", msg)
	}

	if loadedMsg.Err != ErrNoClient {
		t.Errorf("Expected ErrNoClient, got %v", loadedMsg.Err)
	}
}

func TestGetLinkedMemoriesInvalidInput(t *testing.T) {
	client := &mockMemoryClient{}

	cmd := GetLinkedMemories(client, "", DirectionBoth)
	msg := cmd()

	loadedMsg, ok := msg.(LinkedMemoriesLoadedMsg)
	if !ok {
		t.Fatalf("Expected LinkedMemoriesLoadedMsg, got %T", msg)
	}

	if loadedMsg.Err != ErrSourceNotFound {
		t.Errorf("Expected ErrSourceNotFound, got %v", loadedMsg.Err)
	}
}

// Test Model integration with link management

func TestModelShowCreateLinkDialog(t *testing.T) {
	m := NewModel()

	if m.ShowingCreateLinkDialog() {
		t.Error("Should not show create link dialog initially")
	}

	m.ShowCreateLinkDialog()

	if !m.ShowingCreateLinkDialog() {
		t.Error("Should show create link dialog after ShowCreateLinkDialog()")
	}

	if m.GetLinkType() != pb.LinkType_LINK_TYPE_REFERENCES {
		t.Error("Default link type should be REFERENCES")
	}

	if m.GetLinkStrength() != 0.7 {
		t.Errorf("Expected default strength 0.7, got %f", m.GetLinkStrength())
	}
}

func TestModelHideCreateLinkDialog(t *testing.T) {
	m := NewModel()
	m.ShowCreateLinkDialog()
	m.SetLinkTargetSearch("test search")

	m.HideCreateLinkDialog()

	if m.ShowingCreateLinkDialog() {
		t.Error("Should not show create link dialog after HideCreateLinkDialog()")
	}

	if m.GetLinkTargetSearch() != "" {
		t.Error("Link target search should be cleared")
	}
}

func TestModelSetLinkType(t *testing.T) {
	m := NewModel()

	linkTypes := []pb.LinkType{
		pb.LinkType_LINK_TYPE_EXTENDS,
		pb.LinkType_LINK_TYPE_BUILDS_UPON,
		pb.LinkType_LINK_TYPE_CONTRADICTS,
		pb.LinkType_LINK_TYPE_IMPLEMENTS,
	}

	for _, lt := range linkTypes {
		m.SetLinkType(lt)
		if m.GetLinkType() != lt {
			t.Errorf("Expected link type %v, got %v", lt, m.GetLinkType())
		}
	}
}

func TestModelSetLinkStrength(t *testing.T) {
	m := NewModel()

	tests := []struct {
		input    float32
		expected float32
	}{
		{0.5, 0.5},
		{0.0, 0.0},
		{1.0, 1.0},
		{-0.1, 0.0}, // Should clamp to 0
		{1.5, 1.0},  // Should clamp to 1
	}

	for _, tt := range tests {
		m.SetLinkStrength(tt.input)
		if m.GetLinkStrength() != tt.expected {
			t.Errorf("Input %f: expected %f, got %f", tt.input, tt.expected, m.GetLinkStrength())
		}
	}
}

func TestModelNavigationHistory(t *testing.T) {
	m := NewModel()

	if m.CanNavigateBack() {
		t.Error("Should not be able to navigate back initially")
	}

	if m.CanNavigateForward() {
		t.Error("Should not be able to navigate forward initially")
	}

	// Set a memory
	memory1 := &pb.MemoryNote{Id: "memory1"}
	m.SetMemory(memory1)

	// Navigate to another memory (simulating link navigation)
	m.navigationHistory.Push("memory1")
	m.navigationHistory.Push("memory2") // Push the second memory to create history
	memory2 := &pb.MemoryNote{Id: "memory2"}
	m.SetMemory(memory2)

	if !m.CanNavigateBack() {
		t.Error("Should be able to navigate back after navigation")
	}

	// Navigate back
	m.navigationHistory.Back()
	if !m.CanNavigateForward() {
		t.Error("Should be able to navigate forward after going back")
	}
}
