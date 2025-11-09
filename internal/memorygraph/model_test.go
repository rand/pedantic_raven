package memorygraph

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Test Update with unfocused model returns unchanged.
func TestUpdateUnfocused(t *testing.T) {
	m := NewModel()
	m.SetFocus(false)

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if cmd != nil {
		t.Error("Expected no command when unfocused")
	}
	if newModel.offsetX != m.offsetX {
		t.Error("Expected state unchanged when unfocused")
	}
}

// Test keyboard pan left (h, left).
func TestKeyboardPanLeft(t *testing.T) {
	m := NewModel()
	m.offsetX = 0

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = newModel

	if m.offsetX != 2.0 {
		t.Errorf("Expected offsetX 2.0, got %f", m.offsetX)
	}

	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = newModel

	if m.offsetX != 4.0 {
		t.Errorf("Expected offsetX 4.0, got %f", m.offsetX)
	}
}

// Test keyboard pan right (l, right).
func TestKeyboardPanRight(t *testing.T) {
	m := NewModel()
	m.offsetX = 0

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = newModel

	if m.offsetX != -2.0 {
		t.Errorf("Expected offsetX -2.0, got %f", m.offsetX)
	}

	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = newModel

	if m.offsetX != -4.0 {
		t.Errorf("Expected offsetX -4.0, got %f", m.offsetX)
	}
}

// Test keyboard pan up (k, up).
func TestKeyboardPanUp(t *testing.T) {
	m := NewModel()
	m.offsetY = 0

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = newModel

	if m.offsetY != 1.0 {
		t.Errorf("Expected offsetY 1.0, got %f", m.offsetY)
	}

	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = newModel

	if m.offsetY != 2.0 {
		t.Errorf("Expected offsetY 2.0, got %f", m.offsetY)
	}
}

// Test keyboard pan down (j, down).
func TestKeyboardPanDown(t *testing.T) {
	m := NewModel()
	m.offsetY = 0

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = newModel

	if m.offsetY != -1.0 {
		t.Errorf("Expected offsetY -1.0, got %f", m.offsetY)
	}

	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newModel

	if m.offsetY != -2.0 {
		t.Errorf("Expected offsetY -2.0, got %f", m.offsetY)
	}
}

// Test keyboard zoom in (+, =).
func TestKeyboardZoomIn(t *testing.T) {
	m := NewModel()
	m.zoom = 1.0

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	m = newModel

	if m.zoom != 1.1 {
		t.Errorf("Expected zoom 1.1, got %f", m.zoom)
	}

	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'='}})
	m = newModel

	if m.zoom < 1.2 || m.zoom > 1.22 {
		t.Errorf("Expected zoom ~1.21, got %f", m.zoom)
	}
}

// Test keyboard zoom out (-, _).
func TestKeyboardZoomOut(t *testing.T) {
	m := NewModel()
	m.zoom = 1.0

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	m = newModel

	if m.zoom != 0.9 {
		t.Errorf("Expected zoom 0.9, got %f", m.zoom)
	}

	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'_'}})
	m = newModel

	if m.zoom < 0.8 || m.zoom > 0.82 {
		t.Errorf("Expected zoom ~0.81, got %f", m.zoom)
	}
}

// Test zoom limits.
func TestZoomLimits(t *testing.T) {
	m := NewModel()
	m.zoom = 2.9

	// Zoom in to max
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	m = newModel

	if m.zoom > 3.0 {
		t.Errorf("Expected zoom capped at 3.0, got %f", m.zoom)
	}

	// Try to zoom past max
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	m = newModel

	if m.zoom != 3.0 {
		t.Errorf("Expected zoom to stay at 3.0, got %f", m.zoom)
	}

	// Zoom out to min
	m.zoom = 0.3
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	m = newModel

	if m.zoom != 0.3 {
		t.Errorf("Expected zoom to stay at 0.3, got %f", m.zoom)
	}
}

// Test keyboard reset (0).
func TestKeyboardReset(t *testing.T) {
	m := NewModel()
	m.zoom = 2.0
	m.offsetX = 10
	m.offsetY = 20

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}})
	m = newModel

	if m.zoom != 1.0 {
		t.Errorf("Expected zoom reset to 1.0, got %f", m.zoom)
	}
	if m.offsetX != 0 {
		t.Errorf("Expected offsetX reset to 0, got %f", m.offsetX)
	}
	if m.offsetY != 0 {
		t.Errorf("Expected offsetY reset to 0, got %f", m.offsetY)
	}
}

// Test keyboard center on node (c).
func TestKeyboardCenter(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node", X: 15, Y: 25})
	m.SelectNode("test-node")

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = newModel

	if m.offsetX != -15 {
		t.Errorf("Expected offsetX -15, got %f", m.offsetX)
	}
	if m.offsetY != -25 {
		t.Errorf("Expected offsetY -25, got %f", m.offsetY)
	}
}

// Test keyboard center with no node selected.
func TestKeyboardCenterNoNode(t *testing.T) {
	m := NewModel()
	m.offsetX = 10
	m.offsetY = 20

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = newModel

	// Should not change offset
	if m.offsetX != 10 || m.offsetY != 20 {
		t.Error("Expected offset unchanged when no node selected")
	}
}

// Test ApplyForceLayout directly.
func TestApplyForceLayoutDirectCall(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: 0, Y: 0, Mass: 1.0})
	m.graph.AddNode(&Node{ID: "node-2", X: 5, Y: 0, Mass: 1.0})

	initialSteps := m.layoutSteps
	m.ApplyForceLayout()

	if m.layoutSteps != initialSteps+1 {
		t.Errorf("Expected layoutSteps to increment from %d, got %d", initialSteps, m.layoutSteps)
	}
}

// Test keyboard space applies single layout iteration.
func TestKeyboardSpace(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: 0, Y: 0, Mass: 1.0})
	m.graph.AddNode(&Node{ID: "node-2", X: 5, Y: 0, Mass: 1.0})
	m.SetFocus(true)

	initialSteps := m.layoutSteps
	// Create a KeyMsg that msg.String() will return "space" for
	keyMsg := tea.KeyMsg{
		Type: tea.KeySpace,
	}

	// Debug: check what String() returns
	keyStr := keyMsg.String()
	if keyStr != "space" {
		t.Logf("WARNING: KeyMsg.String() returned '%s', expected 'space'", keyStr)
	}

	newModel, _ := m.Update(keyMsg)
	m = newModel

	if m.layoutSteps != initialSteps+1 {
		t.Errorf("Expected layoutSteps to increment from %d, got %d (KeyMsg.String() = '%s')", initialSteps, m.layoutSteps, keyStr)
	}
}

// Test keyboard 'r' re-layouts graph.
func TestKeyboardRelayout(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1", X: 100, Y: 200})
	m.layoutSteps = 10

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = newModel

	// Layout should be reinitialized and stabilized (adds 50 to existing 10)
	if m.layoutSteps != 60 {
		t.Errorf("Expected layoutSteps 60 after relayout (10 + 50), got %d", m.layoutSteps)
	}

	// Node position should change from initial (100, 200)
	node := m.graph.Nodes["node-1"]
	if node.X == 100 && node.Y == 200 {
		t.Error("Expected node position to change after relayout")
	}
}

// Test selectNextNode cycles through nodes.
func TestSelectNextNode(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1"})
	m.graph.AddNode(&Node{ID: "node-2"})
	m.graph.AddNode(&Node{ID: "node-3"})

	// First selection
	m.selectNextNode()
	first := m.selectedNodeID
	if first == "" {
		t.Error("Expected a node to be selected")
	}

	// Second selection
	m.selectNextNode()
	second := m.selectedNodeID
	if second == "" || second == first {
		t.Error("Expected next node to be selected")
	}

	// Third selection
	m.selectNextNode()
	third := m.selectedNodeID
	if third == "" || third == first || third == second {
		t.Error("Expected third node to be selected")
	}

	// Fourth selection should wrap to first
	m.selectNextNode()
	if m.selectedNodeID != first {
		t.Error("Expected selection to wrap to first node")
	}
}

// Test selectNextNode with empty graph.
func TestSelectNextNodeEmpty(t *testing.T) {
	m := NewModel()
	m.selectNextNode() // Should not panic

	if m.selectedNodeID != "" {
		t.Error("Expected no selection in empty graph")
	}
}

// Test selectPreviousNode cycles backward.
func TestSelectPreviousNode(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1"})
	m.graph.AddNode(&Node{ID: "node-2"})
	m.graph.AddNode(&Node{ID: "node-3"})

	// Select first node
	m.selectNextNode()
	first := m.selectedNodeID

	// Select previous should wrap to last
	m.selectPreviousNode()
	if m.selectedNodeID == first {
		t.Error("Expected selection to wrap to last node")
	}

	// Keep going backward
	m.selectPreviousNode()
	m.selectPreviousNode()

	// Should be back at first
	if m.selectedNodeID != first {
		t.Error("Expected to cycle back to first node")
	}
}

// Test keyboard Tab selects next node.
func TestKeyboardTab(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1"})
	m.graph.AddNode(&Node{ID: "node-2"})

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = newModel

	if m.selectedNodeID == "" {
		t.Error("Expected node to be selected")
	}
}

// Test keyboard Shift+Tab selects previous node.
func TestKeyboardShiftTab(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "node-1"})
	m.graph.AddNode(&Node{ID: "node-2"})

	// Select first
	m.selectNextNode()
	first := m.selectedNodeID

	// Shift+Tab should go to last
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = newModel

	if m.selectedNodeID == first {
		t.Error("Expected different node to be selected")
	}
}

// Test keyboard Enter emits NodeSelectedMsg.
func TestKeyboardEnter(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node"})
	m.SelectNode("test-node")

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel

	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	nodeMsg, ok := msg.(NodeSelectedMsg)
	if !ok {
		t.Fatal("Expected NodeSelectedMsg")
	}
	if nodeMsg.NodeID != "test-node" {
		t.Errorf("Expected NodeID 'test-node', got %s", nodeMsg.NodeID)
	}
}

// Test keyboard Enter with no selection.
func TestKeyboardEnterNoSelection(t *testing.T) {
	m := NewModel()

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel

	if cmd != nil {
		t.Error("Expected no command when no node selected")
	}
}

// Test keyboard 'e' emits ExpandNodeMsg.
func TestKeyboardExpand(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node"})
	m.SelectNode("test-node")

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m = newModel

	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	expandMsg, ok := msg.(ExpandNodeMsg)
	if !ok {
		t.Fatal("Expected ExpandNodeMsg")
	}
	if expandMsg.NodeID != "test-node" {
		t.Errorf("Expected NodeID 'test-node', got %s", expandMsg.NodeID)
	}
}

// Test keyboard 'e' with no selection.
func TestKeyboardExpandNoSelection(t *testing.T) {
	m := NewModel()

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m = newModel

	if cmd != nil {
		t.Error("Expected no command when no node selected")
	}
}

// Test keyboard 'x' emits CollapseNodeMsg.
func TestKeyboardCollapse(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node"})
	m.SelectNode("test-node")

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m = newModel

	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	collapseMsg, ok := msg.(CollapseNodeMsg)
	if !ok {
		t.Fatal("Expected CollapseNodeMsg")
	}
	if collapseMsg.NodeID != "test-node" {
		t.Errorf("Expected NodeID 'test-node', got %s", collapseMsg.NodeID)
	}
}

// Test keyboard 'x' with no selection.
func TestKeyboardCollapseNoSelection(t *testing.T) {
	m := NewModel()

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m = newModel

	if cmd != nil {
		t.Error("Expected no command when no node selected")
	}
}

// Test GraphLoadedMsg initializes and stabilizes layout.
func TestGraphLoadedMsg(t *testing.T) {
	m := NewModel()

	graph := NewGraph()
	graph.AddNode(&Node{ID: "node-1"})
	graph.AddNode(&Node{ID: "node-2"})

	msg := GraphLoadedMsg{Graph: graph}
	newModel, cmd := m.Update(msg)
	m = newModel

	if cmd != nil {
		t.Error("Expected no command from GraphLoadedMsg")
	}
	if m.graph != graph {
		t.Error("Expected graph to be set")
	}
	if m.layoutSteps != 50 {
		t.Errorf("Expected 50 layout steps after loading, got %d", m.layoutSteps)
	}

	// Nodes should have positions
	for id, node := range m.graph.Nodes {
		if node.X == 0 && node.Y == 0 {
			t.Errorf("Node %s still at origin after layout", id)
		}
	}
}

// Test GraphErrorMsg sets error.
func TestGraphErrorMsg(t *testing.T) {
	m := NewModel()

	err := errors.New("test error")
	msg := GraphErrorMsg{Err: err}
	newModel, cmd := m.Update(msg)
	m = newModel

	if cmd != nil {
		t.Error("Expected no command from GraphErrorMsg")
	}
	if m.err != err {
		t.Error("Expected error to be set")
	}
}

// Test navigateNodeCmd returns correct message.
func TestNavigateNodeCmd(t *testing.T) {
	m := NewModel()
	cmd := m.navigateNodeCmd("test-node")

	msg := cmd()
	nodeMsg, ok := msg.(NodeSelectedMsg)
	if !ok {
		t.Fatal("Expected NodeSelectedMsg")
	}
	if nodeMsg.NodeID != "test-node" {
		t.Errorf("Expected NodeID 'test-node', got %s", nodeMsg.NodeID)
	}
}

// Test expandNodeCmd returns correct message.
func TestExpandNodeCmd(t *testing.T) {
	m := NewModel()
	cmd := m.expandNodeCmd("test-node")

	msg := cmd()
	expandMsg, ok := msg.(ExpandNodeMsg)
	if !ok {
		t.Fatal("Expected ExpandNodeMsg")
	}
	if expandMsg.NodeID != "test-node" {
		t.Errorf("Expected NodeID 'test-node', got %s", expandMsg.NodeID)
	}
}

// Test collapseNodeCmd returns correct message.
func TestCollapseNodeCmd(t *testing.T) {
	m := NewModel()
	cmd := m.collapseNodeCmd("test-node")

	msg := cmd()
	collapseMsg, ok := msg.(CollapseNodeMsg)
	if !ok {
		t.Fatal("Expected CollapseNodeMsg")
	}
	if collapseMsg.NodeID != "test-node" {
		t.Errorf("Expected NodeID 'test-node', got %s", collapseMsg.NodeID)
	}
}

// Test unhandled message returns unchanged.
func TestUnhandledMessage(t *testing.T) {
	m := NewModel()
	m.offsetX = 10

	type UnknownMsg struct{}
	newModel, cmd := m.Update(UnknownMsg{})

	if cmd != nil {
		t.Error("Expected no command for unhandled message")
	}
	if newModel.offsetX != 10 {
		t.Error("Expected state unchanged for unhandled message")
	}
}
