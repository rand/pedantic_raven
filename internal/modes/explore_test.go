package modes

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/memorygraph"
)

func TestNewExploreMode(t *testing.T) {
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
	mode := NewExploreMode()
	mode.Init()

	cmd := mode.OnExit()

	// Should return BaseMode's OnExit result (nil)
	if cmd != nil {
		t.Error("OnExit should return nil command")
	}
}

func TestExploreModeUpdateWindowSize(t *testing.T) {
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
	mode := NewExploreMode()
	// Don't call Init, so components remain nil

	view := mode.View()

	expected := "Initializing memory workspace..."
	if view != expected {
		t.Errorf("Expected view '%s', got '%s'", expected, view)
	}
}

func TestExploreModeViewWithGraph(t *testing.T) {
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
