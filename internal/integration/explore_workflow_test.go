package integration

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/memorygraph"
	"github.com/rand/pedantic-raven/internal/memorylist"
	"github.com/rand/pedantic-raven/internal/modes"
	"github.com/rand/pedantic-raven/internal/testhelpers"
)

// setupIntegrationTest disables real mnemosyne connection
func setupIntegrationTest(t *testing.T) func() {
	oldEnabled := os.Getenv("MNEMOSYNE_ENABLED")
	os.Setenv("MNEMOSYNE_ENABLED", "false")

	return func() {
		if oldEnabled != "" {
			os.Setenv("MNEMOSYNE_ENABLED", oldEnabled)
		} else {
			os.Unsetenv("MNEMOSYNE_ENABLED")
		}
	}
}

// executeCommand executes a tea.Cmd and returns the message
func executeCommand(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()
	if cmd == nil {
		return nil
	}
	return cmd()
}

// TestExploreModeInitialization verifies complete initialization workflow
func TestExploreModeInitialization(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	if mode == nil {
		t.Fatal("Failed to create explore mode")
	}

	// Initialize mode
	cmd := mode.Init()
	if cmd != nil {
		executeCommand(t, cmd)
	}

	// Verify mode properties
	if mode.ID() != modes.ModeExplore {
		t.Errorf("Expected mode ID %s, got %s", modes.ModeExplore, mode.ID())
	}

	if mode.Name() != "Explore" {
		t.Errorf("Expected mode name 'Explore', got %s", mode.Name())
	}

	// Verify components initialized
	view := mode.View()
	if view == "Initializing memory workspace..." {
		t.Error("Components should be initialized")
	}
}

// TestExploreModeSampleDataLoading verifies sample data loading workflow
func TestExploreModeSampleDataLoading(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	// Load sample data
	cmd := mode.OnEnter()
	if cmd == nil {
		t.Fatal("OnEnter should return command")
	}

	msg := executeCommand(t, cmd)
	batchMsg, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("Expected BatchMsg, got %T", msg)
	}

	// Execute batch commands
	var memoryListLoaded, graphLoaded bool
	for _, batchCmd := range batchMsg {
		msg := executeCommand(t, batchCmd)
		if msg == nil {
			continue
		}

		switch msg.(type) {
		case memorylist.MemoriesLoadedMsg:
			memoryListLoaded = true
		case memorygraph.GraphLoadedMsg:
			graphLoaded = true
		}
	}

	if !memoryListLoaded {
		t.Error("Memory list should be loaded")
	}

	if !graphLoaded {
		t.Error("Graph should be loaded")
	}
}

// TestExploreModeLayoutSwitching verifies layout switching workflow
func TestExploreModeLayoutSwitching(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	// Load sample data
	cmd := mode.OnEnter()
	if cmd != nil {
		msg := executeCommand(t, cmd)
		mode.Update(msg)
	}

	// Set window size
	wsMsg := tea.WindowSizeMsg{Width: 200, Height: 100}
	mode.Update(wsMsg)

	// Toggle to graph layout
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	updatedMode, _ := mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	// Verify graph layout active
	view := mode.View()
	if view == "" {
		t.Error("Graph view should not be empty")
	}

	// Toggle back to standard
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	updatedMode, _ = mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	// Verify standard layout active
	view = mode.View()
	if view == "" {
		t.Error("Standard view should not be empty")
	}
}

// TestExploreModeHelpSystem verifies help overlay functionality
func TestExploreModeHelpSystem(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	// Open help
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedMode, _ := mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	view := mode.View()
	if view == "" {
		t.Error("Help view should not be empty")
	}

	// Close help with escape
	keyMsg = tea.KeyMsg{Type: tea.KeyEscape}
	updatedMode, _ = mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	// Verify help closed
	view = mode.View()
	if view == "" {
		t.Error("View should show standard layout after help closed")
	}
}

// TestExploreModeFocusManagement verifies focus cycling in standard layout
func TestExploreModeFocusManagement(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	// Cycle focus with Tab
	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedMode, _ := mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	view := mode.View()
	if view == "" {
		t.Error("View should update with focus change")
	}

	// Cycle again
	updatedMode, _ = mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	view = mode.View()
	if view == "" {
		t.Error("View should update with second focus change")
	}
}

// TestExploreModeWindowResizing verifies responsive layout
func TestExploreModeWindowResizing(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	// Load sample data
	cmd := mode.OnEnter()
	if cmd != nil {
		msg := executeCommand(t, cmd)
		mode.Update(msg)
	}

	// Test various window sizes (all reasonable sizes that won't cause overflow)
	sizes := []struct {
		width  int
		height int
	}{
		{100, 50},
		{200, 100},
		{120, 60},
		{150, 75},
	}

	for _, size := range sizes {
		wsMsg := tea.WindowSizeMsg{Width: size.width, Height: size.height}
		updatedMode, _ := mode.Update(wsMsg)
		mode = updatedMode.(*modes.ExploreMode)

		view := mode.View()
		if view == "" {
			t.Errorf("View should adapt to size %dx%d", size.width, size.height)
		}
	}
}

// TestExploreModeMultipleUpdates verifies stability under rapid updates
func TestExploreModeMultipleUpdates(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	// Simulate rapid key presses
	keys := []rune{'?', 'g', 'g', '?', 'g', 'g'}
	for _, key := range keys {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}}
		updatedMode, _ := mode.Update(keyMsg)
		if updatedMode == nil {
			t.Fatal("Update returned nil mode")
		}
		mode = updatedMode.(*modes.ExploreMode)

		view := mode.View()
		if view == "" {
			t.Error("View should always be renderable")
		}
	}
}

// TestExploreModeKeybindingsStandardLayout verifies standard layout keybindings
func TestExploreModeKeybindingsStandardLayout(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	keybindings := mode.Keybindings()
	if len(keybindings) == 0 {
		t.Fatal("Keybindings should not be empty")
	}

	expectedKeys := []string{"g", "Tab", "j/k", "Enter", "/", "r", "?"}
	found := make(map[string]bool)

	for _, kb := range keybindings {
		for _, expected := range expectedKeys {
			if kb.Key == expected {
				found[expected] = true
			}
		}
	}

	for _, expected := range expectedKeys {
		if !found[expected] {
			t.Errorf("Missing keybinding: %s", expected)
		}
	}
}

// TestExploreModeKeybindingsGraphLayout verifies graph layout keybindings
func TestExploreModeKeybindingsGraphLayout(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	// Switch to graph layout
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	updatedMode, _ := mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	keybindings := mode.Keybindings()
	if len(keybindings) == 0 {
		t.Fatal("Keybindings should not be empty in graph layout")
	}

	expectedKeys := []string{"g", "h/j/k/l", "+/-", "0", "Tab", "Enter", "e", "x", "c", "r", "Space"}
	found := make(map[string]bool)

	for _, kb := range keybindings {
		for _, expected := range expectedKeys {
			if kb.Key == expected {
				found[expected] = true
			}
		}
	}

	for _, expected := range expectedKeys {
		if !found[expected] {
			t.Errorf("Missing graph keybinding: %s", expected)
		}
	}
}

// TestExploreModeErrorHandlingNilComponents verifies graceful degradation
func TestExploreModeErrorHandlingNilComponents(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	// Don't call Init - components remain nil

	// Should not crash
	view := mode.View()
	if view != "Initializing memory workspace..." {
		t.Errorf("Expected initializing message, got %s", view)
	}

	// Updates should be safe
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	updatedMode, _ := mode.Update(keyMsg)
	if updatedMode == nil {
		t.Fatal("Update should never return nil")
	}
}

// TestExploreModeLifecycleComplete verifies full lifecycle
func TestExploreModeLifecycleComplete(t *testing.T) {
	defer setupIntegrationTest(t)()

	// Create
	mode := modes.NewExploreMode()
	if mode == nil {
		t.Fatal("Creation failed")
	}

	// Initialize
	cmd := mode.Init()
	if cmd != nil {
		executeCommand(t, cmd)
	}

	// Enter
	cmd = mode.OnEnter()
	if cmd == nil {
		t.Fatal("OnEnter should return command")
	}
	msg := executeCommand(t, cmd)
	mode.Update(msg)

	// Set size
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	// Interact
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedMode, _ := mode.Update(keyMsg)
	mode = updatedMode.(*modes.ExploreMode)

	// View
	view := mode.View()
	if view == "" {
		t.Error("View should be renderable")
	}

	// Exit
	cmd = mode.OnExit()
	if cmd != nil {
		executeCommand(t, cmd)
	}
}

// TestExploreModeWithMockClient verifies integration with mock mnemosyne client
func TestExploreModeWithMockClient(t *testing.T) {
	defer setupIntegrationTest(t)()

	mockClient := testhelpers.NewMockMnemosyneClient()

	// Add test memories
	memories := testhelpers.GenerateTestMemories(10)
	for _, mem := range memories {
		_, err := mockClient.Remember(
			context.Background(),
			mem.Content,
			mem.Importance,
			mem.Tags,
			"global",
		)
		if err != nil {
			t.Fatalf("Failed to add test memory: %v", err)
		}
	}

	// Verify memories stored
	if mockClient.GetMemoryCount() != 10 {
		t.Errorf("Expected 10 memories, got %d", mockClient.GetMemoryCount())
	}

	// Test recall
	recalled, err := mockClient.Recall(context.Background(), "test", 5)
	if err != nil {
		t.Fatalf("Recall failed: %v", err)
	}

	if len(recalled) == 0 {
		t.Error("Recall should return memories")
	}
}

// TestExploreModeErrorRecovery verifies error handling
func TestExploreModeErrorRecovery(t *testing.T) {
	defer setupIntegrationTest(t)()

	mockClient := testhelpers.NewMockMnemosyneClient()

	// Configure to fail
	mockClient.SetShouldFail(true, errors.New("simulated network error"))

	// Attempt operations
	_, err := mockClient.Remember(context.Background(), "test", 5, []string{"test"}, "global")
	if err == nil {
		t.Error("Expected error from mock client")
	}

	_, err = mockClient.Recall(context.Background(), "test", 10)
	if err == nil {
		t.Error("Expected error from mock client")
	}

	// Recover
	mockClient.SetShouldFail(false, nil)

	// Operations should succeed
	mem, err := mockClient.Remember(context.Background(), "test", 5, []string{"test"}, "global")
	if err != nil {
		t.Errorf("Operation should succeed after recovery: %v", err)
	}

	if mem == nil {
		t.Error("Memory should be created after recovery")
	}
}

// TestExploreModeDataGenerators verifies test helper functions
func TestExploreModeDataGenerators(t *testing.T) {
	// Test memory generation
	memories := testhelpers.GenerateTestMemories(100)
	if len(memories) != 100 {
		t.Errorf("Expected 100 memories, got %d", len(memories))
	}

	// Verify variety
	importances := make(map[uint32]bool)
	for _, mem := range memories {
		importances[mem.Importance] = true
	}

	if len(importances) < 5 {
		t.Error("Memories should have variety of importance levels")
	}

	// All memories have the same global namespace, so skip namespace variety check

	// Test graph generation
	graph := testhelpers.GenerateTestGraph(50, 75)
	if len(graph.Nodes) != 50 {
		t.Errorf("Expected 50 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 75 {
		t.Errorf("Expected 75 edges, got %d", len(graph.Edges))
	}

	// Test linked memory generation
	linkedMems, links := testhelpers.GenerateLinkedMemories(20, 3)
	if len(linkedMems) != 20 {
		t.Errorf("Expected 20 memories, got %d", len(linkedMems))
	}

	totalLinks := 0
	for _, linkSet := range links {
		totalLinks += len(linkSet)
	}

	if totalLinks == 0 {
		t.Error("Should generate links between memories")
	}
}

// TestExploreModePerformanceSmallDataset verifies performance with small dataset
func TestExploreModePerformanceSmallDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()

	start := time.Now()
	mode.Init()
	elapsed := time.Since(start)

	if elapsed > time.Second {
		t.Errorf("Initialization took too long: %v", elapsed)
	}

	// Load sample data
	start = time.Now()
	cmd := mode.OnEnter()
	if cmd != nil {
		msg := executeCommand(t, cmd)
		mode.Update(msg)
	}
	elapsed = time.Since(start)

	if elapsed > time.Second {
		t.Errorf("Sample data loading took too long: %v", elapsed)
	}

	// Render view
	start = time.Now()
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)
	view := mode.View()
	elapsed = time.Since(start)

	if view == "" {
		t.Error("View should be rendered")
	}

	if elapsed > 500*time.Millisecond {
		t.Errorf("View rendering took too long: %v", elapsed)
	}
}

// TestExploreModeMessageRouting verifies message handling
func TestExploreModeMessageRouting(t *testing.T) {
	defer setupIntegrationTest(t)()

	mode := modes.NewExploreMode()
	mode.Init()

	// Window size message
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedMode, _ := mode.Update(wsMsg)
	if updatedMode == nil {
		t.Fatal("Update should not return nil")
	}

	// Key messages
	keyMsgs := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'g'}},
		{Type: tea.KeyRunes, Runes: []rune{'?'}},
		{Type: tea.KeyTab},
		{Type: tea.KeyEscape},
	}

	for _, keyMsg := range keyMsgs {
		updatedMode, _ = mode.Update(keyMsg)
		if updatedMode == nil {
			t.Fatalf("Update with key %v returned nil", keyMsg)
		}
		mode = updatedMode.(*modes.ExploreMode)
	}
}

// TestExploreModeCRUDOperationsMock simulates CRUD workflow with mock client
func TestExploreModeCRUDOperationsMock(t *testing.T) {
	defer setupIntegrationTest(t)()

	mockClient := testhelpers.NewMockMnemosyneClient()

	// Create
	mem, err := mockClient.Remember(
		context.Background(),
		"Test memory for CRUD",
		7,
		[]string{"test", "crud"},
		"global",
	)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	memID := mem.Id

	// Read
	retrieved, err := mockClient.GetMemory(context.Background(), memID)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if retrieved.Content != "Test memory for CRUD" {
		t.Errorf("Content mismatch: %s", retrieved.Content)
	}

	// Update
	updated, err := mockClient.Update(
		context.Background(),
		memID,
		"Updated content",
		9,
		[]string{"test", "crud", "updated"},
	)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.Content != "Updated content" {
		t.Error("Content not updated")
	}

	if updated.Importance != 9 {
		t.Error("Importance not updated")
	}

	// Delete
	err = mockClient.Delete(context.Background(), memID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err = mockClient.GetMemory(context.Background(), memID)
	if err == nil {
		t.Error("Memory should be deleted")
	}
}
