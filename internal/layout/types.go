// Package layout provides a flexible multi-pane layout system for TUI applications.
//
// The layout engine supports hierarchical pane composition with splits (horizontal/vertical),
// focus management, and responsive adaptation to terminal size changes.
//
// Inspired by the crush TUI's layout system and tmux's pane model.
package layout

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Rect represents a rectangular area in the terminal.
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Direction specifies how a split divides space.
type Direction int

const (
	Horizontal Direction = iota // Side by side (left | right)
	Vertical                     // Top and bottom (top / bottom)
)

// String returns the human-readable name of the direction.
func (d Direction) String() string {
	switch d {
	case Horizontal:
		return "Horizontal"
	case Vertical:
		return "Vertical"
	default:
		return "Unknown"
	}
}

// LayoutMode determines the overall layout structure.
type LayoutMode int

const (
	// Focus mode: Single pane (typically editor) takes full screen
	LayoutFocus LayoutMode = iota

	// Standard mode: Main pane + sidebar (≥120x30 terminals)
	LayoutStandard

	// Analysis mode: Split screen for comparing/analyzing (≥120x30 terminals)
	LayoutAnalysis

	// Compact mode: Optimized for smaller terminals (<120x30)
	// Sidebar becomes overlay, simplified layout
	LayoutCompact

	// Custom mode: User-defined layout from configuration
	LayoutCustom
)

// String returns the human-readable name of the layout mode.
func (m LayoutMode) String() string {
	switch m {
	case LayoutFocus:
		return "Focus"
	case LayoutStandard:
		return "Standard"
	case LayoutAnalysis:
		return "Analysis"
	case LayoutCompact:
		return "Compact"
	case LayoutCustom:
		return "Custom"
	default:
		return "Unknown"
	}
}

// PaneID uniquely identifies a pane in the layout.
type PaneID string

// Common pane IDs
const (
	PaneEditor      PaneID = "editor"
	PaneSidebar     PaneID = "sidebar"
	PaneDiagnostics PaneID = "diagnostics"
	PaneMemory      PaneID = "memory"
	PaneAnalysis    PaneID = "analysis"
	PaneOrchestrate PaneID = "orchestrate"
	PaneChat        PaneID = "chat"
	PaneStatus      PaneID = "status"
)

// Component represents a UI component that can be rendered in a pane.
//
// This is the interface that all renderable components must implement.
// Components receive Bubble Tea messages and return updated state + commands.
type Component interface {
	// Update processes a Bubble Tea message and returns updated state.
	Update(msg tea.Msg) (Component, tea.Cmd)

	// View renders the component to a string within the given area.
	// The focused parameter indicates if this component has input focus.
	View(area Rect, focused bool) string

	// ID returns the unique identifier for this component.
	ID() PaneID
}

// Pane represents a node in the layout tree.
//
// A pane can either be:
// - Leaf: Contains a renderable component
// - Split: Contains two child panes (left/right or top/bottom)
type Pane interface {
	// Render generates the visual representation of this pane and its children.
	Render(area Rect, focusedID PaneID) string

	// Update processes a message and returns updated pane state.
	// Messages are routed to the appropriate child component.
	Update(msg tea.Msg, focusedID PaneID) (Pane, tea.Cmd)

	// FindComponent locates a component by ID in this pane's subtree.
	// Returns nil if not found.
	FindComponent(id PaneID) Component

	// AllPaneIDs returns the IDs of all leaf panes in this subtree.
	AllPaneIDs() []PaneID

	// IsLeaf returns true if this is a leaf pane (contains a component).
	IsLeaf() bool
}

// LeafPane contains a single renderable component.
type LeafPane struct {
	component Component
}

// NewLeafPane creates a new leaf pane containing the given component.
func NewLeafPane(component Component) *LeafPane {
	return &LeafPane{
		component: component,
	}
}

// Render implements Pane.
func (p *LeafPane) Render(area Rect, focusedID PaneID) string {
	if p.component == nil {
		return ""
	}
	focused := p.component.ID() == focusedID
	return p.component.View(area, focused)
}

// Update implements Pane.
func (p *LeafPane) Update(msg tea.Msg, focusedID PaneID) (Pane, tea.Cmd) {
	if p.component == nil {
		return p, nil
	}

	// Only update if this pane has focus (optimization)
	if p.component.ID() == focusedID {
		updatedComponent, cmd := p.component.Update(msg)
		return &LeafPane{component: updatedComponent}, cmd
	}

	return p, nil
}

// FindComponent implements Pane.
func (p *LeafPane) FindComponent(id PaneID) Component {
	if p.component != nil && p.component.ID() == id {
		return p.component
	}
	return nil
}

// AllPaneIDs implements Pane.
func (p *LeafPane) AllPaneIDs() []PaneID {
	if p.component != nil {
		return []PaneID{p.component.ID()}
	}
	return []PaneID{}
}

// IsLeaf implements Pane.
func (p *LeafPane) IsLeaf() bool {
	return true
}

// SplitPane divides space between two child panes.
type SplitPane struct {
	direction Direction
	ratio     float32 // 0.0 to 1.0, determines how much space the first child gets
	first     Pane
	second    Pane
}

// NewSplitPane creates a new split pane with the given direction, ratio, and children.
//
// The ratio determines how space is divided:
// - 0.5 = equal split (50/50)
// - 0.7 = first child gets 70%, second gets 30%
// - 0.3 = first child gets 30%, second gets 70%
//
// Ratio is clamped to [0.1, 0.9] to ensure both panes are visible.
func NewSplitPane(direction Direction, ratio float32, first, second Pane) *SplitPane {
	// Clamp ratio to reasonable bounds
	if ratio < 0.1 {
		ratio = 0.1
	}
	if ratio > 0.9 {
		ratio = 0.9
	}

	return &SplitPane{
		direction: direction,
		ratio:     ratio,
		first:     first,
		second:    second,
	}
}

// Render implements Pane.
func (p *SplitPane) Render(area Rect, focusedID PaneID) string {
	firstArea, secondArea := p.computeSplitAreas(area)

	firstView := ""
	if p.first != nil {
		firstView = p.first.Render(firstArea, focusedID)
	}

	secondView := ""
	if p.second != nil {
		secondView = p.second.Render(secondArea, focusedID)
	}

	// Combine views based on direction
	switch p.direction {
	case Horizontal:
		return combineHorizontal(firstView, secondView)
	case Vertical:
		return combineVertical(firstView, secondView)
	default:
		return firstView
	}
}

// Update implements Pane.
func (p *SplitPane) Update(msg tea.Msg, focusedID PaneID) (Pane, tea.Cmd) {
	var cmds []tea.Cmd

	// Update both children (they'll optimize internally)
	if p.first != nil {
		updatedFirst, cmd := p.first.Update(msg, focusedID)
		p.first = updatedFirst
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if p.second != nil {
		updatedSecond, cmd := p.second.Update(msg, focusedID)
		p.second = updatedSecond
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return p, tea.Batch(cmds...)
}

// FindComponent implements Pane.
func (p *SplitPane) FindComponent(id PaneID) Component {
	if p.first != nil {
		if comp := p.first.FindComponent(id); comp != nil {
			return comp
		}
	}

	if p.second != nil {
		if comp := p.second.FindComponent(id); comp != nil {
			return comp
		}
	}

	return nil
}

// AllPaneIDs implements Pane.
func (p *SplitPane) AllPaneIDs() []PaneID {
	ids := []PaneID{}

	if p.first != nil {
		ids = append(ids, p.first.AllPaneIDs()...)
	}

	if p.second != nil {
		ids = append(ids, p.second.AllPaneIDs()...)
	}

	return ids
}

// IsLeaf implements Pane.
func (p *SplitPane) IsLeaf() bool {
	return false
}

// computeSplitAreas calculates the rectangular areas for the two children.
func (p *SplitPane) computeSplitAreas(area Rect) (Rect, Rect) {
	switch p.direction {
	case Horizontal:
		splitX := area.X + int(float32(area.Width)*p.ratio)
		firstWidth := splitX - area.X
		secondWidth := area.Width - firstWidth

		return Rect{
				X:      area.X,
				Y:      area.Y,
				Width:  firstWidth,
				Height: area.Height,
			}, Rect{
				X:      splitX,
				Y:      area.Y,
				Width:  secondWidth,
				Height: area.Height,
			}

	case Vertical:
		splitY := area.Y + int(float32(area.Height)*p.ratio)
		firstHeight := splitY - area.Y
		secondHeight := area.Height - firstHeight

		return Rect{
				X:      area.X,
				Y:      area.Y,
				Width:  area.Width,
				Height: firstHeight,
			}, Rect{
				X:      area.X,
				Y:      splitY,
				Width:  area.Width,
				Height: secondHeight,
			}

	default:
		return area, Rect{}
	}
}

// combineHorizontal merges two views side by side.
func combineHorizontal(left, right string) string {
	// TODO: Implement proper horizontal combination with line-by-line merging
	// For now, return left view (will be implemented in view rendering phase)
	return left + right
}

// combineVertical merges two views top and bottom.
func combineVertical(top, bottom string) string {
	// Vertical is simpler - just concatenate with newline
	if top == "" {
		return bottom
	}
	if bottom == "" {
		return top
	}
	return top + "\n" + bottom
}
