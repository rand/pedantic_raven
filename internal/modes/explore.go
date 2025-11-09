package modes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/memorygraph"
)

// ExploreMode provides memory workspace with graph visualization.
type ExploreMode struct {
	*BaseMode
	graph *memorygraph.Model
}

// NewExploreMode creates a new explore mode with graph visualization.
func NewExploreMode() *ExploreMode {
	return &ExploreMode{
		BaseMode: NewBaseMode(
			ModeExplore,
			"Explore",
			"Memory workspace with graph visualization",
		),
		graph: nil, // Will be initialized in Init
	}
}

// Init initializes the explore mode.
func (m *ExploreMode) Init() tea.Cmd {
	// Initialize graph model
	graphModel := memorygraph.NewModel()
	m.graph = &graphModel

	// Initialize base mode
	if m.BaseMode != nil {
		return m.BaseMode.Init()
	}
	return nil
}

// OnEnter is called when explore mode becomes active.
func (m *ExploreMode) OnEnter() tea.Cmd {
	// TODO: Load initial graph data from mnemosyne
	// For now, create a sample graph
	return m.loadSampleGraph()
}

// OnExit is called when explore mode becomes inactive.
func (m *ExploreMode) OnExit() tea.Cmd {
	if m.BaseMode != nil {
		return m.BaseMode.OnExit()
	}
	return nil
}

// Update processes messages.
func (m *ExploreMode) Update(msg tea.Msg) (Mode, tea.Cmd) {
	if m.graph == nil {
		return m, nil
	}

	var cmds []tea.Cmd

	// Handle window size for graph
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		// Reserve space for title/status bars (about 10 lines)
		graphHeight := wsMsg.Height - 10
		if graphHeight < 5 {
			graphHeight = 5
		}
		m.graph.SetSize(wsMsg.Width, graphHeight)
	}

	// Forward to graph model
	updatedGraph, cmd := m.graph.Update(msg)
	*m.graph = updatedGraph
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Update base mode
	if m.BaseMode != nil {
		_, baseCmd := m.BaseMode.Update(msg)
		if baseCmd != nil {
			cmds = append(cmds, baseCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the explore mode.
func (m *ExploreMode) View() string {
	if m.graph == nil {
		return "Initializing graph visualization..."
	}

	return m.graph.View()
}

// Keybindings returns the keybindings for explore mode.
func (m *ExploreMode) Keybindings() []Keybinding {
	return []Keybinding{
		{Key: "h/j/k/l", Description: "Pan graph"},
		{Key: "+/-", Description: "Zoom in/out"},
		{Key: "0", Description: "Reset view"},
		{Key: "Tab", Description: "Select next node"},
		{Key: "Enter", Description: "Navigate to node"},
		{Key: "e", Description: "Expand node"},
		{Key: "x", Description: "Collapse node"},
		{Key: "c", Description: "Center on selected"},
		{Key: "r", Description: "Re-layout graph"},
		{Key: "Space", Description: "Layout step"},
	}
}

// loadSampleGraph creates a sample graph for demonstration.
func (m *ExploreMode) loadSampleGraph() tea.Cmd {
	return func() tea.Msg {
		// Create a sample graph
		graph := memorygraph.NewGraph()

		// Add root node
		graph.AddNode(&memorygraph.Node{
			ID:         "root",
			IsExpanded: true,
		})

		// Add some child nodes
		graph.AddNode(&memorygraph.Node{
			ID:         "concept-a",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "concept-b",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "concept-c",
			IsExpanded: true,
		})

		// Add deeper nodes
		graph.AddNode(&memorygraph.Node{
			ID:         "detail-a1",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "detail-a2",
			IsExpanded: true,
		})
		graph.AddNode(&memorygraph.Node{
			ID:         "detail-b1",
			IsExpanded: true,
		})

		// Add edges
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "root",
			TargetID: "concept-a",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "root",
			TargetID: "concept-b",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "root",
			TargetID: "concept-c",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "concept-a",
			TargetID: "detail-a1",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "concept-a",
			TargetID: "detail-a2",
			Strength: 1.0,
		})
		graph.AddEdge(&memorygraph.Edge{
			SourceID: "concept-b",
			TargetID: "detail-b1",
			Strength: 1.0,
		})

		return memorygraph.GraphLoadedMsg{Graph: graph}
	}
}
