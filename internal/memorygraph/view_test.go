package memorygraph

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Test View with error state.
func TestViewError(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.SetError(errors.New("test error"))

	view := m.View()

	if !strings.Contains(view, "Error") {
		t.Error("Expected error message in view")
	}
	if !strings.Contains(view, "test error") {
		t.Error("Expected specific error text in view")
	}
}

// Test View with empty graph.
func TestViewEmpty(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	view := m.View()

	if !strings.Contains(view, "No graph to display") {
		t.Error("Expected empty state message")
	}
}

// Test View with normal graph.
func TestViewNormal(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	graph := NewGraph()
	graph.AddNode(&Node{
		ID:     "test-node-1",
		Memory: &pb.MemoryNote{Id: "test-node-1"},
		X:      0,
		Y:      0,
	})
	graph.AddNode(&Node{
		ID:     "test-node-2",
		Memory: &pb.MemoryNote{Id: "test-node-2"},
		X:      10,
		Y:      10,
	})
	graph.AddEdge(&Edge{
		SourceID: "test-node-1",
		TargetID: "test-node-2",
		Strength: 1.0,
	})

	m.SetGraph(graph)
	view := m.View()

	if !strings.Contains(view, "Memory Graph") {
		t.Error("Expected header with 'Memory Graph'")
	}
	if !strings.Contains(view, "2 nodes") {
		t.Error("Expected node count in header")
	}
	if !strings.Contains(view, "1 edges") {
		t.Error("Expected edge count in header")
	}
}

// Test renderHeader with different graph sizes.
func TestRenderHeader(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	// Empty graph
	header := m.renderHeader()
	if !strings.Contains(header, "Memory Graph") {
		t.Error("Expected 'Memory Graph' in header")
	}

	// Graph with nodes
	graph := NewGraph()
	for i := 0; i < 5; i++ {
		graph.AddNode(&Node{ID: string(rune('A' + i))})
	}
	for i := 0; i < 3; i++ {
		graph.AddEdge(&Edge{SourceID: "A", TargetID: string(rune('B' + i))})
	}
	m.SetGraph(graph)

	header = m.renderHeader()
	if !strings.Contains(header, "5 nodes") {
		t.Error("Expected '5 nodes' in header")
	}
	if !strings.Contains(header, "3 edges") {
		t.Error("Expected '3 edges' in header")
	}
	if !strings.Contains(header, "Zoom:") {
		t.Error("Expected zoom level in header")
	}
	if !strings.Contains(header, "Layout steps:") {
		t.Error("Expected layout steps in header")
	}
}

// Test renderFooter with no selection.
func TestRenderFooterNoSelection(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	footer := m.renderFooter()

	if !strings.Contains(footer, "pan") {
		t.Error("Expected pan hint in footer")
	}
	if !strings.Contains(footer, "zoom") {
		t.Error("Expected zoom hint in footer")
	}
	if !strings.Contains(footer, "layout step") {
		t.Error("Expected layout step hint in footer")
	}
	if !strings.Contains(footer, "Tab") {
		t.Errorf("Expected Tab hint in footer, got: %s", footer)
	}
}

// Test renderFooter with selection.
func TestRenderFooterWithSelection(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.SelectNode("test-node")

	footer := m.renderFooter()

	if !strings.Contains(footer, "navigate") {
		t.Error("Expected navigate hint in footer")
	}
	if !strings.Contains(footer, "expand") {
		t.Error("Expected expand hint in footer")
	}
	if !strings.Contains(footer, "collapse") {
		t.Error("Expected collapse hint in footer")
	}
}

// Test renderFooter truncation.
func TestRenderFooterTruncation(t *testing.T) {
	m := NewModel()
	m.SetSize(20, 10) // Very narrow

	footer := m.renderFooter()
	// Should not exceed width
	lines := strings.Split(footer, "\n")
	for _, line := range lines {
		if lipgloss.Width(line) > 20 {
			t.Errorf("Footer line exceeds width: %d > 20", lipgloss.Width(line))
		}
	}
}

// Test worldToScreen transformation.
func TestWorldToScreen(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.zoom = 1.0
	m.offsetX = 0
	m.offsetY = 0

	// Center should map to screen center
	screenX, screenY := m.worldToScreen(0, 0)
	expectedX := 80 / 2
	expectedY := (20 - 3) / 2 // height - header/footer/separators

	if screenX != expectedX {
		t.Errorf("Expected screenX %d, got %d", expectedX, screenX)
	}
	if screenY != expectedY {
		t.Errorf("Expected screenY %d, got %d", expectedY, screenY)
	}
}

// Test worldToScreen with offset.
func TestWorldToScreenWithOffset(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.zoom = 1.0
	m.offsetX = 10
	m.offsetY = 5

	screenX, screenY := m.worldToScreen(0, 0)

	// With offset, world (0,0) should shift
	expectedX := 80/2 + 10
	expectedY := (20-3)/2 + 5

	if screenX != expectedX {
		t.Errorf("Expected screenX %d, got %d", expectedX, screenX)
	}
	if screenY != expectedY {
		t.Errorf("Expected screenY %d, got %d", expectedY, screenY)
	}
}

// Test worldToScreen with zoom.
func TestWorldToScreenWithZoom(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.zoom = 2.0
	m.offsetX = 0
	m.offsetY = 0

	screenX, screenY := m.worldToScreen(10, 5)

	// With 2x zoom, world (10,5) should be scaled
	expectedX := 80/2 + 10*2
	expectedY := (20-3)/2 + 5*2

	if screenX != expectedX {
		t.Errorf("Expected screenX %d, got %d", expectedX, screenX)
	}
	if screenY != expectedY {
		t.Errorf("Expected screenY %d, got %d", expectedY, screenY)
	}
}

// Test newCanvas creates correct dimensions.
func TestNewCanvas(t *testing.T) {
	canvas := newCanvas(80, 20)

	if canvas.width != 80 {
		t.Errorf("Expected width 80, got %d", canvas.width)
	}
	if canvas.height != 20 {
		t.Errorf("Expected height 20, got %d", canvas.height)
	}
	if len(canvas.cells) != 20 {
		t.Errorf("Expected 20 rows, got %d", len(canvas.cells))
	}
	if len(canvas.cells[0]) != 80 {
		t.Errorf("Expected 80 columns, got %d", len(canvas.cells[0]))
	}
}

// Test newCanvas initializes with spaces.
func TestNewCanvasInitialization(t *testing.T) {
	canvas := newCanvas(5, 5)

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if canvas.cells[y][x] != ' ' {
				t.Errorf("Expected space at (%d,%d), got %c", x, y, canvas.cells[y][x])
			}
		}
	}
}

// Test Canvas.DrawText within bounds.
func TestCanvasDrawText(t *testing.T) {
	canvas := newCanvas(10, 5)
	canvas.DrawText(2, 1, "Hello")

	expected := "Hello"
	actual := string(canvas.cells[1][2:7])

	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

// Test Canvas.DrawText multiline.
func TestCanvasDrawTextMultiline(t *testing.T) {
	canvas := newCanvas(10, 5)
	canvas.DrawText(0, 0, "AB\nCD")

	if string(canvas.cells[0][0:2]) != "AB" {
		t.Error("Expected 'AB' on first line")
	}
	if string(canvas.cells[1][0:2]) != "CD" {
		t.Error("Expected 'CD' on second line")
	}
}

// Test Canvas.DrawText clipping.
func TestCanvasDrawTextClipping(t *testing.T) {
	canvas := newCanvas(5, 3)

	// Draw outside bounds
	canvas.DrawText(-5, -5, "XXX")
	canvas.DrawText(10, 10, "YYY")

	// Canvas should remain all spaces
	for y := 0; y < 3; y++ {
		for x := 0; x < 5; x++ {
			if canvas.cells[y][x] != ' ' {
				t.Errorf("Expected space at (%d,%d) after out-of-bounds draw", x, y)
			}
		}
	}

	// Draw partially outside
	canvas.DrawText(3, 1, "ABCDE")
	// Only "AB" should fit
	if canvas.cells[1][3] != 'A' || canvas.cells[1][4] != 'B' {
		t.Error("Expected partial text 'AB' at position (3,1)")
	}
}

// Test Canvas.DrawLine horizontal.
func TestCanvasDrawLineHorizontal(t *testing.T) {
	canvas := newCanvas(10, 5)
	canvas.DrawLine(1, 2, 5, 2, "242")

	// Should draw horizontal line
	for x := 1; x <= 5; x++ {
		if canvas.cells[2][x] != '─' {
			t.Errorf("Expected '─' at (%d,2), got %c", x, canvas.cells[2][x])
		}
	}
}

// Test Canvas.DrawLine vertical.
func TestCanvasDrawLineVertical(t *testing.T) {
	canvas := newCanvas(10, 10)
	canvas.DrawLine(3, 1, 3, 5, "242")

	// Should draw vertical line
	for y := 1; y <= 5; y++ {
		if canvas.cells[y][3] != '│' {
			t.Errorf("Expected '│' at (3,%d), got %c", y, canvas.cells[y][3])
		}
	}
}

// Test Canvas.DrawLine diagonal.
func TestCanvasDrawLineDiagonal(t *testing.T) {
	canvas := newCanvas(10, 10)
	canvas.DrawLine(1, 1, 5, 5, "242")

	// Should draw diagonal line (Bresenham's algorithm)
	// At least the endpoints should be drawn
	if canvas.cells[1][1] == ' ' {
		t.Error("Expected line character at start (1,1)")
	}
	if canvas.cells[5][5] == ' ' {
		t.Error("Expected line character at end (5,5)")
	}
}

// Test Canvas.DrawLine clipping.
func TestCanvasDrawLineClipping(t *testing.T) {
	canvas := newCanvas(5, 5)

	// Draw line outside bounds (should not panic)
	canvas.DrawLine(-10, -10, -5, -5, "242")
	canvas.DrawLine(20, 20, 25, 25, "242")

	// Canvas should remain all spaces
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			if canvas.cells[y][x] != ' ' {
				t.Errorf("Expected space at (%d,%d) after out-of-bounds line draw", x, y)
			}
		}
	}
}

// Test Canvas.Render basic.
func TestCanvasRender(t *testing.T) {
	canvas := newCanvas(5, 3)
	canvas.DrawText(0, 0, "ABCDE")
	canvas.DrawText(0, 1, "FGHIJ")
	canvas.DrawText(0, 2, "KLMNO")

	result := canvas.Render()
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "ABCDE" {
		t.Errorf("Expected 'ABCDE', got '%s'", lines[0])
	}
	if lines[1] != "FGHIJ" {
		t.Errorf("Expected 'FGHIJ', got '%s'", lines[1])
	}
	if lines[2] != "KLMNO" {
		t.Errorf("Expected 'KLMNO', got '%s'", lines[2])
	}
}

// Test Canvas.Render empty.
func TestCanvasRenderEmpty(t *testing.T) {
	canvas := newCanvas(5, 3)
	result := canvas.Render()
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if line != "     " {
			t.Errorf("Expected 5 spaces on line %d, got '%s'", i, line)
		}
	}
}

// Test abs helper function.
func TestAbs(t *testing.T) {
	if abs(5) != 5 {
		t.Error("abs(5) should be 5")
	}
	if abs(-5) != 5 {
		t.Error("abs(-5) should be 5")
	}
	if abs(0) != 0 {
		t.Error("abs(0) should be 0")
	}
}

// Test stripAnsi helper function.
func TestStripAnsi(t *testing.T) {
	// No ANSI codes
	if stripAnsi("Hello") != "Hello" {
		t.Error("stripAnsi should preserve plain text")
	}

	// With ANSI color codes
	text := "\x1b[31mRed\x1b[0m"
	if stripAnsi(text) != "Red" {
		t.Errorf("Expected 'Red', got '%s'", stripAnsi(text))
	}

	// Multiple ANSI codes
	text = "\x1b[1m\x1b[31mBold Red\x1b[0m"
	if stripAnsi(text) != "Bold Red" {
		t.Errorf("Expected 'Bold Red', got '%s'", stripAnsi(text))
	}
}

// Test drawNode renders node at correct position.
func TestDrawNode(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	canvas := newCanvas(80, 17) // height - 3 for header/footer

	node := &Node{
		ID: "test",
		X:  0,
		Y:  0,
	}

	m.drawNode(canvas, node)

	// Node should be drawn near center
	// Check that some non-space characters exist around center
	centerX := 80 / 2
	centerY := 17 / 2

	found := false
	for dy := -3; dy <= 3; dy++ {
		for dx := -10; dx <= 10; dx++ {
			y := centerY + dy
			x := centerX + dx
			if y >= 0 && y < 17 && x >= 0 && x < 80 {
				if canvas.cells[y][x] != ' ' {
					found = true
					break
				}
			}
		}
	}

	if !found {
		t.Error("Expected node to be drawn near center")
	}
}

// Test drawNode with selected node.
func TestDrawNodeSelected(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)
	m.SelectNode("selected")

	canvas := newCanvas(80, 17)

	node := &Node{
		ID: "selected",
		X:  0,
		Y:  0,
	}

	m.drawNode(canvas, node)

	// Should draw something (selected style applied)
	centerX := 80 / 2
	centerY := 17 / 2

	found := false
	for dy := -3; dy <= 3; dy++ {
		for dx := -10; dx <= 10; dx++ {
			y := centerY + dy
			x := centerX + dx
			if y >= 0 && y < 17 && x >= 0 && x < 80 {
				if canvas.cells[y][x] != ' ' {
					found = true
					break
				}
			}
		}
	}

	if !found {
		t.Error("Expected selected node to be drawn")
	}
}

// Test drawNode truncates long IDs.
func TestDrawNodeLongID(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	canvas := newCanvas(80, 17)

	node := &Node{
		ID: "very-long-node-identifier-that-exceeds-limit",
		X:  0,
		Y:  0,
	}

	// Should not panic
	m.drawNode(canvas, node)
}

// Test drawEdge connects two nodes.
func TestDrawEdge(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	canvas := newCanvas(80, 17)

	node1 := &Node{ID: "A", X: -5, Y: 0}
	node2 := &Node{ID: "B", X: 5, Y: 0}
	m.graph.AddNode(node1)
	m.graph.AddNode(node2)

	edge := &Edge{
		SourceID: "A",
		TargetID: "B",
		Strength: 1.0,
	}

	m.drawEdge(canvas, edge)

	// Should draw line between nodes
	// Check that some line characters exist
	found := false
	for y := 0; y < 17; y++ {
		for x := 0; x < 80; x++ {
			if canvas.cells[y][x] == '─' || canvas.cells[y][x] == '│' {
				found = true
				break
			}
		}
	}

	if !found {
		t.Error("Expected edge line to be drawn")
	}
}

// Test drawEdge with missing nodes.
func TestDrawEdgeMissingNodes(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	canvas := newCanvas(80, 17)

	edge := &Edge{
		SourceID: "nonexistent-1",
		TargetID: "nonexistent-2",
		Strength: 1.0,
	}

	// Should not panic
	m.drawEdge(canvas, edge)

	// Canvas should remain empty
	allSpaces := true
	for y := 0; y < 17; y++ {
		for x := 0; x < 80; x++ {
			if canvas.cells[y][x] != ' ' {
				allSpaces = false
				break
			}
		}
	}

	if !allSpaces {
		t.Error("Expected canvas to remain empty when edge nodes don't exist")
	}
}

// Test renderGraph with nodes and edges.
func TestRenderGraph(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	graph := NewGraph()
	graph.AddNode(&Node{ID: "A", X: 0, Y: 0})
	graph.AddNode(&Node{ID: "B", X: 10, Y: 0})
	graph.AddEdge(&Edge{SourceID: "A", TargetID: "B", Strength: 1.0})

	m.SetGraph(graph)
	content := m.renderGraph()

	// Should produce some output
	if len(content) == 0 {
		t.Error("Expected non-empty graph rendering")
	}

	// Should have multiple lines
	lines := strings.Split(content, "\n")
	if len(lines) < 5 {
		t.Error("Expected multiple lines in graph rendering")
	}
}

// Test full View rendering integration.
func TestViewIntegration(t *testing.T) {
	m := NewModel()
	m.SetSize(80, 20)

	graph := NewGraph()
	graph.AddNode(&Node{
		ID:     "node-1",
		Memory: &pb.MemoryNote{Id: "node-1", Content: "Test"},
		X:      0,
		Y:      0,
	})
	graph.AddNode(&Node{
		ID:     "node-2",
		Memory: &pb.MemoryNote{Id: "node-2", Content: "Test"},
		X:      10,
		Y:      5,
	})
	graph.AddEdge(&Edge{SourceID: "node-1", TargetID: "node-2", Strength: 1.0})

	m.SetGraph(graph)
	m.InitializeLayout()
	m.SelectNode("node-1")

	view := m.View()

	// Should have header
	if !strings.Contains(view, "Memory Graph") {
		t.Error("Expected header in integrated view")
	}

	// Should have footer
	if !strings.Contains(view, "navigate") {
		t.Error("Expected footer in integrated view")
	}

	// Should have multiple lines
	lines := strings.Split(view, "\n")
	if len(lines) < 10 {
		t.Errorf("Expected at least 10 lines, got %d", len(lines))
	}
}
