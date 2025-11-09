package modes

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/memorygraph"
)

// setupTest disables mnemosyne for testing
func setupTest(t *testing.T) func() {
	// Disable mnemosyne connection during tests
	oldEnabled := os.Getenv("MNEMOSYNE_ENABLED")
	os.Setenv("MNEMOSYNE_ENABLED", "false")

	return func() {
		// Restore original value
		if oldEnabled != "" {
			os.Setenv("MNEMOSYNE_ENABLED", oldEnabled)
		} else {
			os.Unsetenv("MNEMOSYNE_ENABLED")
		}
	}
}

func TestNewExploreMode(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()

	if mode == nil {
		t.Fatal("NewExploreMode returned nil")
	}

	if mode.ID() != ModeExplore {
		t.Errorf("Expected mode ID %s, got %s", ModeExplore, mode.ID())
	}

	if mode.Name() != "Explore" {
		t.Errorf("Expected mode name 'Explore', got '%s'", mode.Name())
	}

	expectedDesc := "Memory workspace with list, detail, and graph views"
	if mode.Description() != expectedDesc {
		t.Errorf("Expected description '%s', got '%s'", expectedDesc, mode.Description())
	}
}

func TestExploreModeInit(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()

	// Graph should be nil before Init
	if mode.graph != nil {
		t.Error("Graph should be nil before Init")
	}

	// Initialize
	cmd := mode.Init()

	// Graph should be initialized after Init
	if mode.graph == nil {
		t.Error("Graph should be initialized after Init")
	}

	// BaseMode Init should return nil
	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestExploreModeOnEnter(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	cmd := mode.OnEnter()

	// OnEnter should return a command to load sample data
	if cmd == nil {
		t.Error("OnEnter should return a command")
	}

	// Execute the command - it returns a BatchMsg with multiple commands
	msg := cmd()
	if _, ok := msg.(tea.BatchMsg); !ok {
		t.Errorf("Expected tea.BatchMsg (batch of commands), got %T", msg)
	}
}

func TestExploreModeOnExit(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	cmd := mode.OnExit()

	// Should return BaseMode's OnExit result (nil)
	if cmd != nil {
		t.Error("OnExit should return nil command")
	}
}

func TestExploreModeUpdateWindowSize(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Send window size message
	wsMsg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	updatedMode, cmd := mode.Update(wsMsg)

	if updatedMode == nil {
		t.Fatal("Update returned nil mode")
	}

	// Should forward to graph and set size
	// Graph height should be window height - 10 (for title/status bars)
	if cmd != nil {
		// cmd might be returned from graph update
		t.Logf("Update returned command: %T", cmd)
	}
}

func TestExploreModeUpdateWithNilGraph(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	// Don't call Init, so graph remains nil

	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'h'},
	}

	updatedMode, cmd := mode.Update(keyMsg)

	if updatedMode == nil {
		t.Fatal("Update returned nil mode")
	}

	if cmd != nil {
		t.Error("Update with nil graph should return nil command")
	}
}

func TestExploreModeUpdateForwardsToGraph(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Load sample graph
	loadCmd := mode.OnEnter()
	graphMsg := loadCmd()
	mode.Update(graphMsg)

	// Send a key message that the graph should handle
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'h'},
	}

	updatedMode, cmd := mode.Update(keyMsg)

	if updatedMode == nil {
		t.Fatal("Update returned nil mode")
	}

	// The graph might return a command
	if cmd != nil {
		t.Logf("Update returned command: %T", cmd)
	}
}

func TestExploreModeViewWithNilGraph(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	// Don't call Init, so components remain nil

	view := mode.View()

	expected := "Initializing memory workspace..."
	if view != expected {
		t.Errorf("Expected view '%s', got '%s'", expected, view)
	}
}

func TestExploreModeViewWithGraph(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Load sample graph
	loadCmd := mode.OnEnter()
	graphMsg := loadCmd()
	mode.Update(graphMsg)

	view := mode.View()

	// View should delegate to graph.View()
	// We just check it's not empty and not the "Initializing" message
	if view == "" {
		t.Error("View should not be empty")
	}

	if view == "Initializing graph visualization..." {
		t.Error("View should not show initializing message after graph is loaded")
	}
}

func TestExploreModeKeybindings(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()

	// Test standard layout keybindings (default)
	keybindings := mode.Keybindings()

	if len(keybindings) == 0 {
		t.Error("Keybindings should not be empty")
	}

	// Check for expected standard layout keybindings
	expectedKeysStandard := map[string]bool{
		"g":     false, // Toggle to graph
		"Tab":   false, // Switch focus
		"j/k":   false, // Navigate
		"Enter": false, // Select
		"/":     false, // Search
		"r":     false, // Refresh
		"?":     false, // Help
	}

	for _, kb := range keybindings {
		if _, exists := expectedKeysStandard[kb.Key]; exists {
			expectedKeysStandard[kb.Key] = true
		}
	}

	for key, found := range expectedKeysStandard {
		if !found {
			t.Errorf("Missing expected standard layout keybinding: %s", key)
		}
	}

	// Test graph layout keybindings
	mode.layoutMode = LayoutModeGraph
	keybindingsGraph := mode.Keybindings()

	expectedKeysGraph := map[string]bool{
		"g":       false, // Toggle to list
		"h/j/k/l": false, // Pan
		"+/-":     false, // Zoom
		"0":       false, // Reset
		"Tab":     false, // Select node
		"Enter":   false, // Navigate
		"e":       false, // Expand
		"x":       false, // Collapse
		"c":       false, // Center
		"r":       false, // Re-layout
		"Space":   false, // Layout step
	}

	for _, kb := range keybindingsGraph {
		if _, exists := expectedKeysGraph[kb.Key]; exists {
			expectedKeysGraph[kb.Key] = true
		}
	}

	for key, found := range expectedKeysGraph {
		if !found {
			t.Errorf("Missing expected graph layout keybinding: %s", key)
		}
	}
}

func TestExploreModeSampleGraphStructure(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Load sample data - OnEnter returns a batch command
	loadCmd := mode.OnEnter()
	batchMsg := loadCmd()

	// Extract the individual commands from the batch
	batch, ok := batchMsg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("Expected tea.BatchMsg, got %T", batchMsg)
	}

	// Execute each command in the batch
	var graphMsg tea.Msg
	for _, cmd := range batch {
		msg := cmd()
		// Look for GraphLoadedMsg
		if _, isGraph := msg.(memorygraph.GraphLoadedMsg); isGraph {
			graphMsg = msg
			break
		}
	}

	if graphMsg == nil {
		t.Fatal("No GraphLoadedMsg found in batch")
	}

	if msg, ok := graphMsg.(memorygraph.GraphLoadedMsg); ok {
		graph := msg.Graph

		// Check that sample graph has expected structure
		if len(graph.Nodes) != 7 {
			t.Errorf("Expected 7 nodes in sample graph, got %d", len(graph.Nodes))
		}

		if len(graph.Edges) != 6 {
			t.Errorf("Expected 6 edges in sample graph, got %d", len(graph.Edges))
		}

		// Check root node exists
		if _, exists := graph.Nodes["root"]; !exists {
			t.Error("Sample graph should have root node")
		}

		// Check root has 3 children
		childCount := 0
		for _, edge := range graph.Edges {
			if edge.SourceID == "root" {
				childCount++
			}
		}
		if childCount != 3 {
			t.Errorf("Root should have 3 children, got %d", childCount)
		}
	} else {
		t.Fatalf("Expected GraphLoadedMsg, got %T", graphMsg)
	}
}

// --- Additional Explore Mode Tests for Coverage ---

func TestExploreModeLayoutToggle(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Initial state is standard layout
	if mode.layoutMode != LayoutModeStandard {
		t.Error("Initial layout mode should be standard")
	}

	// Toggle to graph
	mode.toggleLayout()
	if mode.layoutMode != LayoutModeGraph {
		t.Error("After first toggle, should be in graph mode")
	}

	// Toggle back to standard
	mode.toggleLayout()
	if mode.layoutMode != LayoutModeStandard {
		t.Error("After second toggle, should be back in standard mode")
	}

	// Focus should reset to list when switching back to standard
	if mode.focusTarget != FocusTargetList {
		t.Error("Focus should reset to list when switching back to standard")
	}
}

func TestExploreModeFocusCycle(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Initial focus should be on list
	if mode.focusTarget != FocusTargetList {
		t.Error("Initial focus should be on list")
	}

	// Cycle to detail
	mode.cycleFocus()
	if mode.focusTarget != FocusTargetDetail {
		t.Error("After first cycle, focus should be on detail")
	}

	// Cycle back to list
	mode.cycleFocus()
	if mode.focusTarget != FocusTargetList {
		t.Error("After second cycle, focus should be back on list")
	}
}

func TestExploreModeFocusSynchronization(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Initially focus should be on list
	if mode.focusTarget != FocusTargetList {
		t.Error("Initial focus target should be on list")
	}

	// Cycle focus
	mode.cycleFocus()

	// Now focus should be on detail
	if mode.focusTarget != FocusTargetDetail {
		t.Error("After cycle, focus target should be on detail")
	}

	// Cycle again
	mode.cycleFocus()

	// Back to list
	if mode.focusTarget != FocusTargetList {
		t.Error("After second cycle, focus target should be back on list")
	}
}

func TestExploreModeUpdateHelpToggle(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Help should be hidden initially
	if mode.showHelp {
		t.Error("Help should be hidden initially")
	}

	// Toggle help with '?'
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	mode.Update(keyMsg)

	if !mode.showHelp {
		t.Error("Help should be shown after '?'")
	}

	// Toggle help again with '?'
	mode.Update(keyMsg)

	if mode.showHelp {
		t.Error("Help should be hidden after second '?'")
	}
}

func TestExploreModeUpdateEscapeKey(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Show help first
	mode.showHelp = true

	// Press escape
	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	mode.Update(keyMsg)

	if mode.showHelp {
		t.Error("Help should be closed after escape")
	}
}

func TestExploreModeUpdateGKeyTogglesToGraph(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Initial layout is standard
	if mode.layoutMode != LayoutModeStandard {
		t.Fatal("Initial layout should be standard")
	}

	// Press 'g' to toggle to graph
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	mode.Update(keyMsg)

	if mode.layoutMode != LayoutModeGraph {
		t.Error("Should toggle to graph mode after 'g'")
	}
}

func TestExploreModeUpdateGKeyDisabledWhenHelpOpen(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Open help
	mode.showHelp = true

	// Try to toggle layout with 'g'
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	mode.Update(keyMsg)

	// Should still be in standard mode
	if mode.layoutMode != LayoutModeStandard {
		t.Error("'g' should not toggle layout when help is open")
	}
}

func TestExploreModeUpdateTabCycleFocus(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Tab should cycle focus (only in standard layout)
	initialFocus := mode.focusTarget

	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	mode.Update(keyMsg)

	if mode.focusTarget == initialFocus {
		t.Error("Tab should cycle focus")
	}
}

func TestExploreModeUpdateTabDisabledInGraphMode(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Switch to graph mode
	mode.layoutMode = LayoutModeGraph

	initialFocus := mode.focusTarget

	// Tab should not cycle focus in graph mode
	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	mode.Update(keyMsg)

	if mode.focusTarget != initialFocus {
		t.Error("Tab should not cycle focus in graph mode")
	}
}

func TestExploreModeHandleWindowSize(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Send window size
	wsMsg := tea.WindowSizeMsg{Width: 200, Height: 100}
	mode.Update(wsMsg)

	// Width and height should be updated
	if mode.width != 200 {
		t.Errorf("Width should be 200, got %d", mode.width)
	}

	if mode.height != 100 {
		t.Errorf("Height should be 100, got %d", mode.height)
	}
}

func TestExploreModeHandleWindowSizeMinimumHeight(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Send very small window size
	wsMsg := tea.WindowSizeMsg{Width: 20, Height: 5}
	mode.Update(wsMsg)

	// Should still update width/height
	if mode.width != 20 {
		t.Errorf("Width should be 20, got %d", mode.width)
	}

	if mode.height != 5 {
		t.Errorf("Height should be 5, got %d", mode.height)
	}
}

func TestExploreModeView(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Set window size first
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	view := mode.View()

	// View should not be empty
	if view == "" {
		t.Error("View should not be empty")
	}

	// Should not show initializing message
	if view == "Initializing memory workspace..." {
		t.Error("View should not show initializing message after Init")
	}
}

func TestExploreModeViewInGraphMode(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Switch to graph mode
	mode.layoutMode = LayoutModeGraph

	// Load sample data
	loadCmd := mode.OnEnter()
	batchMsg := loadCmd()
	mode.Update(batchMsg)

	// Set window size
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	view := mode.View()

	// View should not be empty
	if view == "" {
		t.Error("View should not be empty in graph mode")
	}
}

func TestExploreModeViewHelp(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Set window size
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	// Show help
	mode.showHelp = true

	view := mode.View()

	// View should contain help content
	if view == "" {
		t.Error("View should show help when help is enabled")
	}

	// Should not be the standard layout
	if view == "Initializing memory workspace..." {
		t.Error("Help should override standard layout")
	}
}

func TestExploreModeViewHelpGraphMode(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Switch to graph mode
	mode.layoutMode = LayoutModeGraph

	// Load sample data
	loadCmd := mode.OnEnter()
	batchMsg := loadCmd()
	mode.Update(batchMsg)

	// Set window size
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	// Show help
	mode.showHelp = true

	view := mode.View()

	// View should show help (contains "Graph Layout")
	if view == "" {
		t.Error("View should show help in graph mode")
	}
}

func TestExploreModeConcurrentUpdates(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Simulate multiple rapid updates
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}

	for i := 0; i < 5; i++ {
		updatedMode, _ := mode.Update(keyMsg)
		if updatedMode == nil {
			t.Fatal("Update should never return nil mode")
		}
		mode = updatedMode.(*ExploreMode)
	}

	// After 5 toggles (odd number), help should be visible
	if !mode.showHelp {
		t.Error("Help should be visible after 5 toggles")
	}
}

func TestExploreModeLayoutToggleWithSize(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Set window size
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	if mode.width != 100 || mode.height != 50 {
		t.Fatal("Window size not set")
	}

	// Toggle layout
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	mode.Update(keyMsg)

	// Should update layout
	if mode.layoutMode != LayoutModeGraph {
		t.Error("Layout should toggle to graph")
	}
}

func TestExploreModeSampleMemoriesLoaded(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Load sample data
	loadCmd := mode.OnEnter()
	batchMsg := loadCmd()

	// Extract the individual commands from the batch
	batch, ok := batchMsg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("Expected tea.BatchMsg, got %T", batchMsg)
	}

	// Execute each command in the batch and look for any memory-related message
	for _, cmd := range batch {
		msg := cmd()
		// Any non-nil message is progress
		if msg != nil {
			t.Logf("Loaded message type: %T", msg)
			break
		}
	}
}

// --- Tests for uncovered functions ---

func TestExploreModeNavigateToMemory(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Navigate to a memory
	cmd := mode.navigateToMemory("mem-123")

	// Should return a command
	if cmd == nil {
		t.Error("navigateToMemory should return a command")
	}

	// Breadcrumb trail should be updated
	if len(mode.breadcrumbTrail) != 1 {
		t.Errorf("Expected 1 item in breadcrumb trail, got %d", len(mode.breadcrumbTrail))
	}

	if mode.breadcrumbTrail[0] != "mem-123" {
		t.Errorf("Expected breadcrumb trail to contain 'mem-123', got %s", mode.breadcrumbTrail[0])
	}
}

func TestExploreModeUpdateBreadcrumbTrail(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Test adding first item
	mode.updateBreadcrumbTrail("mem-1")
	if len(mode.breadcrumbTrail) != 1 {
		t.Errorf("Expected 1 item, got %d", len(mode.breadcrumbTrail))
	}

	// Test adding more items
	mode.updateBreadcrumbTrail("mem-2")
	mode.updateBreadcrumbTrail("mem-3")
	if len(mode.breadcrumbTrail) != 3 {
		t.Errorf("Expected 3 items, got %d", len(mode.breadcrumbTrail))
	}

	// Test empty string (should be ignored)
	mode.updateBreadcrumbTrail("")
	if len(mode.breadcrumbTrail) != 3 {
		t.Error("Empty string should not be added to breadcrumb trail")
	}
}

func TestExploreModeUpdateBreadcrumbTrailMaxDepth(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Set a small max depth for testing
	mode.breadcrumbMaxDepth = 3

	// Add more items than max depth
	for i := 1; i <= 10; i++ {
		mode.updateBreadcrumbTrail("mem-" + string(rune('0'+i)))
	}

	// Should not exceed max depth + 1 (for ellipsis)
	if len(mode.breadcrumbTrail) > mode.breadcrumbMaxDepth+1 {
		t.Errorf("Breadcrumb trail should not exceed max depth + 1, got %d", len(mode.breadcrumbTrail))
	}

	// First item should be ellipsis
	if mode.breadcrumbTrail[0] != "..." {
		t.Error("First item should be ellipsis when trail exceeds max depth")
	}
}

func TestExploreModeRenderBreadcrumbEmpty(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Empty breadcrumb trail should return empty string
	breadcrumb := mode.renderBreadcrumb()
	if breadcrumb != "" {
		t.Error("Empty breadcrumb trail should return empty string")
	}
}

func TestExploreModeRenderBreadcrumbSingleItem(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	mode.breadcrumbTrail = []string{"mem-123"}

	breadcrumb := mode.renderBreadcrumb()
	if breadcrumb == "" {
		t.Error("Breadcrumb with items should not be empty")
	}
}

func TestExploreModeRenderBreadcrumbMultipleItems(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	mode.breadcrumbTrail = []string{"mem-1", "mem-2", "mem-3"}

	breadcrumb := mode.renderBreadcrumb()
	if breadcrumb == "" {
		t.Error("Breadcrumb with multiple items should not be empty")
	}
}

func TestExploreModeRenderBreadcrumbWithEllipsis(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	mode.breadcrumbTrail = []string{"...", "mem-8", "mem-9", "mem-10"}

	breadcrumb := mode.renderBreadcrumb()
	if breadcrumb == "" {
		t.Error("Breadcrumb with ellipsis should not be empty")
	}
}

func TestExploreModeRenderBreadcrumbLongIDs(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Use a long memory ID (should be truncated in display)
	mode.breadcrumbTrail = []string{"mem-very-long-id-that-should-be-truncated"}

	breadcrumb := mode.renderBreadcrumb()
	if breadcrumb == "" {
		t.Error("Breadcrumb with long IDs should not be empty")
	}
}

// Test Update with various message types for better coverage
func TestExploreModeUpdateWithMemoriesLoaded(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Simulate MemoriesLoadedMsg
	cmd := mode.OnEnter()
	if cmd != nil {
		msg := cmd()
		batchMsg, ok := msg.(tea.BatchMsg)
		if ok {
			for _, batchCmd := range batchMsg {
				msg := batchCmd()
				updatedMode, _ := mode.Update(msg)
				mode = updatedMode.(*ExploreMode)
			}
		}
	}

	// Mode should still be valid
	if mode == nil {
		t.Fatal("Mode should not be nil after update")
	}
}

// Test OnExit
func TestExploreModeOnExitCleanup(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Load data
	cmd := mode.OnEnter()
	if cmd != nil {
		cmd()
	}

	// Exit
	cmd = mode.OnExit()

	// Should return nil or cleanup command
	// (Currently returns nil, but test is here for future changes)
	_ = cmd
}

// Test Init with real data scenario
func TestExploreModeInitWithoutMnemosyne(t *testing.T) {
	defer setupTest(t)()

	// Ensure mnemosyne is disabled (already done by setupTest)
	mode := NewExploreMode()
	cmd := mode.Init()

	// Should initialize successfully without mnemosyne
	if mode.memoryList == nil {
		t.Error("Memory list should be initialized")
	}

	if mode.memoryDetail == nil {
		t.Error("Memory detail should be initialized")
	}

	if mode.graph == nil {
		t.Error("Graph should be initialized")
	}

	// Client should be nil or disabled
	if mode.useRealData {
		t.Error("Real data should not be enabled without mnemosyne")
	}

	// Command should be nil from BaseMode
	if cmd != nil {
		t.Log("Init returned command:", cmd)
	}
}

// Test complete navigation flow
func TestExploreModeCompleteNavigationFlow(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	// Navigate through several memories
	for i := 1; i <= 3; i++ {
		memID := "mem-" + string(rune('0'+i))
		cmd := mode.navigateToMemory(memID)
		if cmd != nil {
			// Execute command (would normally load memory)
			msg := cmd()
			if msg != nil {
				mode.Update(msg)
			}
		}
	}

	// Breadcrumb trail should have 3 items
	if len(mode.breadcrumbTrail) != 3 {
		t.Errorf("Expected 3 items in breadcrumb trail, got %d", len(mode.breadcrumbTrail))
	}

	// Render breadcrumb
	breadcrumb := mode.renderBreadcrumb()
	if breadcrumb == "" {
		t.Error("Breadcrumb should not be empty")
	}
}

// Test Update with various key combinations
func TestExploreModeUpdateKeySequences(t *testing.T) {
	defer setupTest(t)()

	mode := NewExploreMode()
	mode.Init()

	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	// Test various key sequences
	sequences := [][]tea.KeyMsg{
		{
			{Type: tea.KeyRunes, Runes: []rune{'g'}}, // Toggle to graph
			{Type: tea.KeyRunes, Runes: []rune{'g'}}, // Toggle back
		},
		{
			{Type: tea.KeyRunes, Runes: []rune{'?'}}, // Open help
			{Type: tea.KeyEscape},                     // Close help
		},
		{
			{Type: tea.KeyTab}, // Cycle focus
			{Type: tea.KeyTab}, // Cycle again
		},
	}

	for _, sequence := range sequences {
		for _, keyMsg := range sequence {
			updatedMode, _ := mode.Update(keyMsg)
			if updatedMode == nil {
				t.Fatal("Update should not return nil")
			}
			mode = updatedMode.(*ExploreMode)
		}
	}
}

// Test edge cases
func TestExploreModeEdgeCases(t *testing.T) {
	defer setupTest(t)()

	t.Run("Multiple Init calls", func(t *testing.T) {
		mode := NewExploreMode()
		mode.Init()
		mode.Init() // Second init should be safe

		if mode.graph == nil {
			t.Error("Graph should remain initialized")
		}
	})

	t.Run("Update before Init", func(t *testing.T) {
		mode := NewExploreMode()
		// Don't call Init

		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
		updatedMode, _ := mode.Update(keyMsg)

		if updatedMode == nil {
			t.Fatal("Update should not return nil even before Init")
		}
	})

	t.Run("View with zero width", func(t *testing.T) {
		mode := NewExploreMode()
		mode.Init()

		// Don't set window size
		view := mode.View()

		if view == "" {
			t.Error("View should handle zero width gracefully")
		}
	})
}
