package memorygraph

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case GraphLoadedMsg:
		m.SetGraph(msg.Graph)
		m.InitializeLayout()
		m.StabilizeLayout(50) // Run 50 iterations initially
		return m, nil

	case GraphErrorMsg:
		m.SetError(msg.Err)
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "h", "left":
		m.offsetX += 2.0
		return m, nil

	case "l", "right":
		m.offsetX -= 2.0
		return m, nil

	case "k", "up":
		m.offsetY += 1.0
		return m, nil

	case "j", "down":
		m.offsetY -= 1.0
		return m, nil

	case "+", "=":
		// Zoom in
		m.zoom *= 1.1
		if m.zoom > 3.0 {
			m.zoom = 3.0
		}
		return m, nil

	case "-", "_":
		// Zoom out
		m.zoom *= 0.9
		if m.zoom < 0.3 {
			m.zoom = 0.3
		}
		return m, nil

	case "0":
		// Reset zoom
		m.zoom = 1.0
		m.offsetX = 0
		m.offsetY = 0
		return m, nil

	case "c":
		// Center on selected node
		if m.selectedNodeID != "" {
			node := m.graph.Nodes[m.selectedNodeID]
			if node != nil {
				m.offsetX = -node.X
				m.offsetY = -node.Y
			}
		}
		return m, nil

	case " ":
		// Run layout iteration
		m.ApplyForceLayout()
		return m, nil

	case "r":
		// Re-layout
		m.InitializeLayout()
		m.StabilizeLayout(50)
		return m, nil

	case "tab":
		// Select next node
		m.selectNextNode()
		return m, nil

	case "shift+tab":
		// Select previous node
		m.selectPreviousNode()
		return m, nil

	case "enter":
		// Navigate to selected node
		if m.selectedNodeID != "" {
			return m, m.navigateNodeCmd(m.selectedNodeID)
		}
		return m, nil

	case "e":
		// Expand selected node
		if m.selectedNodeID != "" {
			return m, m.expandNodeCmd(m.selectedNodeID)
		}
		return m, nil

	case "x":
		// Collapse selected node
		if m.selectedNodeID != "" {
			return m, m.collapseNodeCmd(m.selectedNodeID)
		}
		return m, nil
	}

	return m, nil
}

// selectNextNode selects the next node in the graph.
func (m *Model) selectNextNode() {
	if len(m.graph.Nodes) == 0 {
		return
	}

	// Get all node IDs in order
	nodeIDs := make([]string, 0, len(m.graph.Nodes))
	for id := range m.graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}

	if len(nodeIDs) == 0 {
		return
	}

	// Find current index
	currentIndex := -1
	for i, id := range nodeIDs {
		if id == m.selectedNodeID {
			currentIndex = i
			break
		}
	}

	// Select next
	nextIndex := (currentIndex + 1) % len(nodeIDs)
	m.selectedNodeID = nodeIDs[nextIndex]
}

// selectPreviousNode selects the previous node in the graph.
func (m *Model) selectPreviousNode() {
	if len(m.graph.Nodes) == 0 {
		return
	}

	// Get all node IDs in order
	nodeIDs := make([]string, 0, len(m.graph.Nodes))
	for id := range m.graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}

	if len(nodeIDs) == 0 {
		return
	}

	// Find current index
	currentIndex := -1
	for i, id := range nodeIDs {
		if id == m.selectedNodeID {
			currentIndex = i
			break
		}
	}

	// Select previous
	prevIndex := currentIndex - 1
	if prevIndex < 0 {
		prevIndex = len(nodeIDs) - 1
	}
	m.selectedNodeID = nodeIDs[prevIndex]
}

// navigateNodeCmd returns a command to navigate to a node.
func (m Model) navigateNodeCmd(nodeID string) tea.Cmd {
	return func() tea.Msg {
		return NodeSelectedMsg{
			NodeID: nodeID,
		}
	}
}

// expandNodeCmd returns a command to expand a node.
func (m Model) expandNodeCmd(nodeID string) tea.Cmd {
	return func() tea.Msg {
		return ExpandNodeMsg{
			NodeID: nodeID,
		}
	}
}

// collapseNodeCmd returns a command to collapse a node.
func (m Model) collapseNodeCmd(nodeID string) tea.Cmd {
	return func() tea.Msg {
		return CollapseNodeMsg{
			NodeID: nodeID,
		}
	}
}
