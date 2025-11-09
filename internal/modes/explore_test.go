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

	expectedDesc := "Memory workspace with graph visualization"
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

	// OnEnter should return a command to load the sample graph
	if cmd == nil {
		t.Error("OnEnter should return a command")
	}

	// Execute the command to get the GraphLoadedMsg
	msg := cmd()
	if _, ok := msg.(memorygraph.GraphLoadedMsg); !ok {
		t.Errorf("Expected GraphLoadedMsg, got %T", msg)
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
	// Don't call Init, so graph remains nil

	view := mode.View()

	expected := "Initializing graph visualization..."
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

	keybindings := mode.Keybindings()

	if len(keybindings) == 0 {
		t.Error("Keybindings should not be empty")
	}

	// Check for expected keybindings
	expectedKeys := map[string]bool{
		"h/j/k/l": false,
		"+/-":     false,
		"0":       false,
		"Tab":     false,
		"Enter":   false,
		"e":       false,
		"x":       false,
		"c":       false,
		"r":       false,
		"Space":   false,
	}

	for _, kb := range keybindings {
		if _, exists := expectedKeys[kb.Key]; exists {
			expectedKeys[kb.Key] = true
		}
	}

	for key, found := range expectedKeys {
		if !found {
			t.Errorf("Missing expected keybinding: %s", key)
		}
	}
}

func TestExploreModeSampleGraphStructure(t *testing.T) {
	mode := NewExploreMode()
	mode.Init()

	// Load sample graph
	loadCmd := mode.OnEnter()
	graphMsg := loadCmd()

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
