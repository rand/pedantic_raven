package analyze

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Keyboard shortcuts (simple string-based matching)
const (
	KeyPanLeft  = "h"
	KeyPanRight = "l"
	KeyPanUp    = "k"
	KeyPanDown  = "j"

	KeyZoomIn  = "+"
	KeyZoomOut = "-"

	KeySelect = "enter"
	KeyClear  = "esc"

	KeyCenter = "c"
	KeyReset  = "r"

	KeyStabilize = "s"
)

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case GraphLoadedMsg:
		m.SetGraph(msg.Graph)
		if m.graph != nil {
			m.graph.InitializeLayout()
			// Trigger automatic layout stabilization
			return m, func() tea.Msg {
				return LayoutStepMsg{}
			}
		}
		return m, nil

	case GraphErrorMsg:
		m.SetError(msg.Err)
		return m, nil

	case NodeSelectedMsg:
		m.SelectNode(msg.NodeID)
		return m, nil

	case FilterUpdatedMsg:
		m.ApplyFilter(msg.Filter)
		if m.graph != nil && m.graph.NodeCount() > 0 {
			// Trigger layout after filter change
			return m, func() tea.Msg {
				return LayoutStepMsg{}
			}
		}
		return m, nil

	case LayoutStepMsg:
		// Apply one layout iteration
		if m.graph != nil && m.layoutSteps < 100 {
			m.graph.ApplyForceIteration(m.damping)
			m.layoutSteps++

			// Continue layout if not stable yet
			if m.layoutSteps < 50 {
				return m, func() tea.Msg {
					return LayoutStepMsg{}
				}
			}
		}
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg.String() {
	// Pan controls
	case KeyPanLeft, "left":
		m.Pan(2.0, 0)
		return m, nil

	case KeyPanRight, "right":
		m.Pan(-2.0, 0)
		return m, nil

	case KeyPanUp, "up":
		m.Pan(0, 2.0)
		return m, nil

	case KeyPanDown, "down":
		m.Pan(0, -2.0)
		return m, nil

	// Zoom controls
	case KeyZoomIn, "=":
		m.Zoom(0.1)
		return m, nil

	case KeyZoomOut, "_":
		m.Zoom(-0.1)
		return m, nil

	// Selection controls
	case KeySelect, " ":
		// Toggle through nodes in order
		if m.graph != nil && len(m.graph.Nodes) > 0 {
			// Get next node in iteration order
			found := false
			var firstID string
			for id := range m.graph.Nodes {
				if firstID == "" {
					firstID = id
				}
				if found {
					m.SelectNode(id)
					return m, nil
				}
				if id == m.selectedNodeID {
					found = true
				}
			}
			// Wrap around to first
			if firstID != "" {
				m.SelectNode(firstID)
			}
		}
		return m, nil

	case KeyClear:
		m.SelectNode("")
		return m, nil

	// View controls
	case KeyCenter:
		m.AutoCenter()
		return m, nil

	case KeyReset:
		m.ResetView()
		m.AutoCenter()
		return m, nil

	// Layout controls
	case KeyStabilize:
		m.StabilizeLayout(10)
		return m, nil
	}

	return m, nil
}

// Help returns help text for keyboard shortcuts.
func (m Model) Help() []string {
	return []string{
		"h/l or ←/→: Pan left/right",
		"k/j or ↑/↓: Pan up/down",
		"+/-: Zoom in/out",
		"enter: Select node",
		"esc: Clear selection",
		"c: Center view, r: Reset",
		"s: Stabilize layout",
	}
}
