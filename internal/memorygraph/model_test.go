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

// Test ExpandNode expands a node.
func TestExpandNode(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node", IsExpanded: false})

	m.ExpandNode("test-node")

	node := m.graph.Nodes["test-node"]
	if !node.IsExpanded {
		t.Error("Expected node to be expanded")
	}
}

// Test ExpandNode with nil graph.
func TestExpandNodeNilGraph(t *testing.T) {
	m := NewModel()
	m.graph = nil

	// Should not panic
	m.ExpandNode("test-node")
}

// Test ExpandNode with nonexistent node.
func TestExpandNodeNonexistent(t *testing.T) {
	m := NewModel()

	// Should not panic
	m.ExpandNode("nonexistent")
}

// Test CollapseNode collapses a node.
func TestCollapseNode(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node", IsExpanded: true})

	m.CollapseNode("test-node")

	node := m.graph.Nodes["test-node"]
	if node.IsExpanded {
		t.Error("Expected node to be collapsed")
	}
}

// Test CollapseNode with nil graph.
func TestCollapseNodeNilGraph(t *testing.T) {
	m := NewModel()
	m.graph = nil

	// Should not panic
	m.CollapseNode("test-node")
}

// Test CollapseNode with nonexistent node.
func TestCollapseNodeNonexistent(t *testing.T) {
	m := NewModel()

	// Should not panic
	m.CollapseNode("nonexistent")
}

// Test ExpandNodeMsg updates node state.
func TestExpandNodeMsg(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node", IsExpanded: false})

	msg := ExpandNodeMsg{NodeID: "test-node"}
	newModel, cmd := m.Update(msg)
	m = newModel

	if cmd != nil {
		t.Error("Expected no command from ExpandNodeMsg")
	}

	node := m.graph.Nodes["test-node"]
	if !node.IsExpanded {
		t.Error("Expected node to be expanded after ExpandNodeMsg")
	}
}

// Test CollapseNodeMsg updates node state.
func TestCollapseNodeMsg(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "test-node", IsExpanded: true})

	msg := CollapseNodeMsg{NodeID: "test-node"}
	newModel, cmd := m.Update(msg)
	m = newModel

	if cmd != nil {
		t.Error("Expected no command from CollapseNodeMsg")
	}

	node := m.graph.Nodes["test-node"]
	if node.IsExpanded {
		t.Error("Expected node to be collapsed after CollapseNodeMsg")
	}
}

// Test IsNodeVisible for root nodes.
func TestIsNodeVisibleRoot(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "root", IsExpanded: true})

	// Root nodes (no incoming edges) are always visible
	if !m.IsNodeVisible("root") {
		t.Error("Expected root node to be visible")
	}
}

// Test IsNodeVisible with expanded parent.
func TestIsNodeVisibleExpandedParent(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "parent", IsExpanded: true})
	m.graph.AddNode(&Node{ID: "child", IsExpanded: true})
	m.graph.AddEdge(&Edge{SourceID: "parent", TargetID: "child"})

	// Child should be visible when parent is expanded
	if !m.IsNodeVisible("child") {
		t.Error("Expected child to be visible when parent is expanded")
	}
}

// Test IsNodeVisible with collapsed parent.
func TestIsNodeVisibleCollapsedParent(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "parent", IsExpanded: false})
	m.graph.AddNode(&Node{ID: "child", IsExpanded: true})
	m.graph.AddEdge(&Edge{SourceID: "parent", TargetID: "child"})

	// Child should be hidden when parent is collapsed
	if m.IsNodeVisible("child") {
		t.Error("Expected child to be hidden when parent is collapsed")
	}
}

// Test IsNodeVisible with multiple parents.
func TestIsNodeVisibleMultipleParents(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "parent1", IsExpanded: true})
	m.graph.AddNode(&Node{ID: "parent2", IsExpanded: false})
	m.graph.AddNode(&Node{ID: "child", IsExpanded: true})
	m.graph.AddEdge(&Edge{SourceID: "parent1", TargetID: "child"})
	m.graph.AddEdge(&Edge{SourceID: "parent2", TargetID: "child"})

	// Child should be hidden if ANY parent is collapsed
	if m.IsNodeVisible("child") {
		t.Error("Expected child to be hidden when any parent is collapsed")
	}
}

// Test IsNodeVisible with nested hierarchy.
func TestIsNodeVisibleNestedHierarchy(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "grandparent", IsExpanded: true})
	m.graph.AddNode(&Node{ID: "parent", IsExpanded: false})
	m.graph.AddNode(&Node{ID: "child", IsExpanded: true})
	m.graph.AddEdge(&Edge{SourceID: "grandparent", TargetID: "parent"})
	m.graph.AddEdge(&Edge{SourceID: "parent", TargetID: "child"})

	// Parent is visible (grandparent expanded)
	if !m.IsNodeVisible("parent") {
		t.Error("Expected parent to be visible")
	}

	// Child is hidden (parent collapsed)
	if m.IsNodeVisible("child") {
		t.Error("Expected child to be hidden when parent is collapsed")
	}
}

// Test IsNodeVisible with nil graph.
func TestIsNodeVisibleNilGraph(t *testing.T) {
	m := NewModel()
	m.graph = nil

	if m.IsNodeVisible("test") {
		t.Error("Expected false for nil graph")
	}
}

// Test IsNodeVisible with nonexistent node.
func TestIsNodeVisibleNonexistent(t *testing.T) {
	m := NewModel()

	if m.IsNodeVisible("nonexistent") {
		t.Error("Expected false for nonexistent node")
	}
}

// Test IsEdgeVisible with expanded source.
func TestIsEdgeVisibleExpanded(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "A", IsExpanded: true})
	m.graph.AddNode(&Node{ID: "B", IsExpanded: true})
	edge := &Edge{SourceID: "A", TargetID: "B"}
	m.graph.AddEdge(edge)

	// Edge should be visible when source is expanded
	if !m.IsEdgeVisible(edge) {
		t.Error("Expected edge to be visible when source is expanded")
	}
}

// Test IsEdgeVisible with collapsed source.
func TestIsEdgeVisibleCollapsed(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "A", IsExpanded: false})
	m.graph.AddNode(&Node{ID: "B", IsExpanded: true})
	edge := &Edge{SourceID: "A", TargetID: "B"}
	m.graph.AddEdge(edge)

	// Edge should be hidden when source is collapsed
	if m.IsEdgeVisible(edge) {
		t.Error("Expected edge to be hidden when source is collapsed")
	}
}

// Test IsEdgeVisible with hidden target.
func TestIsEdgeVisibleHiddenTarget(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "root", IsExpanded: true})
	m.graph.AddNode(&Node{ID: "parent", IsExpanded: false})
	m.graph.AddNode(&Node{ID: "child", IsExpanded: true})
	m.graph.AddEdge(&Edge{SourceID: "root", TargetID: "parent"})
	edge := &Edge{SourceID: "parent", TargetID: "child"}
	m.graph.AddEdge(edge)

	// Edge should be hidden when target is not visible
	if m.IsEdgeVisible(edge) {
		t.Error("Expected edge to be hidden when target is not visible")
	}
}

// Test IsEdgeVisible with nil edge.
func TestIsEdgeVisibleNilEdge(t *testing.T) {
	m := NewModel()

	if m.IsEdgeVisible(nil) {
		t.Error("Expected false for nil edge")
	}
}

// Test IsEdgeVisible with missing nodes.
func TestIsEdgeVisibleMissingNodes(t *testing.T) {
	m := NewModel()
	edge := &Edge{SourceID: "nonexistent-1", TargetID: "nonexistent-2"}

	if m.IsEdgeVisible(edge) {
		t.Error("Expected false when edge nodes don't exist")
	}
}

// Test HasChildren returns true for nodes with edges.
func TestHasChildren(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "parent"})
	m.graph.AddNode(&Node{ID: "child"})
	m.graph.AddEdge(&Edge{SourceID: "parent", TargetID: "child"})

	if !m.HasChildren("parent") {
		t.Error("Expected parent to have children")
	}
}

// Test HasChildren returns false for leaf nodes.
func TestHasChildrenLeaf(t *testing.T) {
	m := NewModel()
	m.graph.AddNode(&Node{ID: "leaf"})

	if m.HasChildren("leaf") {
		t.Error("Expected leaf node to have no children")
	}
}

// Test HasChildren with nil graph.
func TestHasChildrenNilGraph(t *testing.T) {
	m := NewModel()
	m.graph = nil

	if m.HasChildren("test") {
		t.Error("Expected false for nil graph")
	}
}

// Test expand/collapse workflow.
func TestExpandCollapseWorkflow(t *testing.T) {
	m := NewModel()

	// Create hierarchy: root -> parent -> child
	m.graph.AddNode(&Node{ID: "root", IsExpanded: true})
	m.graph.AddNode(&Node{ID: "parent", IsExpanded: true})
	m.graph.AddNode(&Node{ID: "child", IsExpanded: true})
	m.graph.AddEdge(&Edge{SourceID: "root", TargetID: "parent"})
	m.graph.AddEdge(&Edge{SourceID: "parent", TargetID: "child"})

	// Initially all visible
	if !m.IsNodeVisible("root") || !m.IsNodeVisible("parent") || !m.IsNodeVisible("child") {
		t.Error("Expected all nodes visible initially")
	}

	// Collapse parent
	m.CollapseNode("parent")

	// Root and parent still visible, child hidden
	if !m.IsNodeVisible("root") {
		t.Error("Expected root to remain visible")
	}
	if !m.IsNodeVisible("parent") {
		t.Error("Expected parent to remain visible")
	}
	if m.IsNodeVisible("child") {
		t.Error("Expected child to be hidden after collapsing parent")
	}

	// Expand parent again
	m.ExpandNode("parent")

	// All visible again
	if !m.IsNodeVisible("child") {
		t.Error("Expected child to be visible after expanding parent")
	}
}
