package layout

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Engine manages the layout, focus, and mode of the application.
//
// The engine maintains a tree of panes and handles:
// - Mode switching (Focus, Standard, Analysis, etc.)
// - Focus management (which pane receives keyboard input)
// - Responsive layout (adapting to terminal size)
// - Message routing to focused components
type Engine struct {
	mode        LayoutMode
	root        Pane
	focusedID   PaneID
	termWidth   int
	termHeight  int
	components  map[PaneID]Component // Registry of all components
}

// NewEngine creates a new layout engine with the specified initial mode.
func NewEngine(mode LayoutMode) *Engine {
	return &Engine{
		mode:       mode,
		root:       nil,
		focusedID:  "",
		termWidth:  120,
		termHeight: 30,
		components: make(map[PaneID]Component),
	}
}

// SetMode changes the layout mode and rebuilds the pane tree.
//
// This is called when the user switches modes (e.g., Focus → Standard)
// or when the terminal is resized to a size that requires a different mode.
func (e *Engine) SetMode(mode LayoutMode) {
	e.mode = mode
	e.rebuildLayout()
}

// Mode returns the current layout mode.
func (e *Engine) Mode() LayoutMode {
	return e.mode
}

// RegisterComponent adds a component to the engine's registry.
//
// Components must be registered before they can be added to the layout.
// If a component with the same ID already exists, it will be replaced.
func (e *Engine) RegisterComponent(component Component) {
	if component != nil {
		e.components[component.ID()] = component
	}
}

// UnregisterComponent removes a component from the engine's registry.
func (e *Engine) UnregisterComponent(id PaneID) {
	delete(e.components, id)
}

// GetComponent retrieves a registered component by ID.
// Returns nil if not found.
func (e *Engine) GetComponent(id PaneID) Component {
	return e.components[id]
}

// SetRoot sets the root pane of the layout tree.
//
// This is typically called during initialization or when switching modes
// to create a new layout structure.
func (e *Engine) SetRoot(root Pane) {
	e.root = root

	// If no focus is set, focus the first leaf pane
	if e.focusedID == "" && root != nil {
		paneIDs := root.AllPaneIDs()
		if len(paneIDs) > 0 {
			e.focusedID = paneIDs[0]
		}
	}
}

// Root returns the root pane of the layout tree.
func (e *Engine) Root() Pane {
	return e.root
}

// SetFocus changes which pane has input focus.
//
// The focused pane will receive keyboard input and be visually highlighted.
// Returns false if the specified pane ID doesn't exist in the layout.
func (e *Engine) SetFocus(id PaneID) bool {
	if e.root == nil {
		return false
	}

	// Verify the pane exists
	if e.root.FindComponent(id) != nil {
		e.focusedID = id
		return true
	}

	return false
}

// FocusedID returns the ID of the currently focused pane.
func (e *Engine) FocusedID() PaneID {
	return e.focusedID
}

// FocusNext moves focus to the next pane in the layout.
//
// Panes are ordered left-to-right, top-to-bottom.
// Wraps around to the first pane after the last.
func (e *Engine) FocusNext() {
	if e.root == nil {
		return
	}

	paneIDs := e.root.AllPaneIDs()
	if len(paneIDs) == 0 {
		return
	}

	// Find current focus index
	currentIndex := -1
	for i, id := range paneIDs {
		if id == e.focusedID {
			currentIndex = i
			break
		}
	}

	// Move to next (wrap around)
	nextIndex := (currentIndex + 1) % len(paneIDs)
	e.focusedID = paneIDs[nextIndex]
}

// FocusPrev moves focus to the previous pane in the layout.
//
// Wraps around to the last pane when at the first.
func (e *Engine) FocusPrev() {
	if e.root == nil {
		return
	}

	paneIDs := e.root.AllPaneIDs()
	if len(paneIDs) == 0 {
		return
	}

	// Find current focus index
	currentIndex := -1
	for i, id := range paneIDs {
		if id == e.focusedID {
			currentIndex = i
			break
		}
	}

	// Move to previous (wrap around)
	prevIndex := currentIndex - 1
	if prevIndex < 0 {
		prevIndex = len(paneIDs) - 1
	}
	e.focusedID = paneIDs[prevIndex]
}

// SetTerminalSize updates the engine's knowledge of terminal dimensions.
//
// This is called when a WindowSizeMsg is received from Bubble Tea.
// The engine will adapt the layout to the new size.
func (e *Engine) SetTerminalSize(width, height int) {
	e.termWidth = width
	e.termHeight = height

	// Check if we need to switch to compact mode
	if width < 120 || height < 30 {
		if e.mode != LayoutCompact && e.mode != LayoutFocus {
			e.SetMode(LayoutCompact)
		}
	}
}

// TerminalSize returns the current terminal dimensions.
func (e *Engine) TerminalSize() (width, height int) {
	return e.termWidth, e.termHeight
}

// Init implements tea.Model.
func (e *Engine) Init() tea.Cmd {
	return nil
}

// Update processes a Bubble Tea message and routes it to the focused component.
//
// The engine handles:
// - WindowSizeMsg: Updates terminal size, potentially switching modes
// - Other messages: Routed to the focused component
func (e *Engine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.SetTerminalSize(msg.Width, msg.Height)
		return e, nil

	default:
		if e.root != nil {
			updatedRoot, cmd := e.root.Update(msg, e.focusedID)
			e.root = updatedRoot
			return e, cmd
		}
		return e, nil
	}
}

// View renders the entire layout to a string.
func (e *Engine) View() string {
	if e.root == nil {
		return "No layout configured"
	}

	area := Rect{
		X:      0,
		Y:      0,
		Width:  e.termWidth,
		Height: e.termHeight,
	}

	return e.root.Render(area, e.focusedID)
}

// rebuildLayout reconstructs the pane tree based on the current mode.
//
// This is called when the mode changes or when the terminal is resized.
// Each mode has a predefined layout structure.
func (e *Engine) rebuildLayout() {
	switch e.mode {
	case LayoutFocus:
		e.buildFocusLayout()
	case LayoutStandard:
		e.buildStandardLayout()
	case LayoutAnalysis:
		e.buildAnalysisLayout()
	case LayoutCompact:
		e.buildCompactLayout()
	default:
		// Keep current layout for LayoutCustom
	}
}

// buildFocusLayout creates a single-pane layout (typically editor).
func (e *Engine) buildFocusLayout() {
	if editor := e.components[PaneEditor]; editor != nil {
		e.root = NewLeafPane(editor)
		e.focusedID = PaneEditor
	}
}

// buildStandardLayout creates editor + sidebar layout.
//
// Layout: [Editor 70%] | [Sidebar 30%]
func (e *Engine) buildStandardLayout() {
	editor := e.components[PaneEditor]
	sidebar := e.components[PaneSidebar]

	if editor != nil && sidebar != nil {
		e.root = NewSplitPane(
			Horizontal,
			0.7, // Editor gets 70%
			NewLeafPane(editor),
			NewLeafPane(sidebar),
		)
		if e.focusedID == "" {
			e.focusedID = PaneEditor
		}
	} else if editor != nil {
		e.root = NewLeafPane(editor)
		e.focusedID = PaneEditor
	}
}

// buildAnalysisLayout creates split-screen layout for analysis.
//
// Layout:
//   [Editor 50%] | [Analysis 50%]
//   [Diagnostics (bottom bar)]
func (e *Engine) buildAnalysisLayout() {
	editor := e.components[PaneEditor]
	analysis := e.components[PaneAnalysis]
	diagnostics := e.components[PaneDiagnostics]

	if editor != nil && analysis != nil {
		topPane := NewSplitPane(
			Horizontal,
			0.5, // Equal split
			NewLeafPane(editor),
			NewLeafPane(analysis),
		)

		if diagnostics != nil {
			e.root = NewSplitPane(
				Vertical,
				0.8, // Top gets 80%, diagnostics gets 20%
				topPane,
				NewLeafPane(diagnostics),
			)
		} else {
			e.root = topPane
		}

		if e.focusedID == "" {
			e.focusedID = PaneEditor
		}
	} else if editor != nil {
		e.root = NewLeafPane(editor)
		e.focusedID = PaneEditor
	}
}

// buildCompactLayout creates simplified layout for small terminals.
//
// Layout: Just the editor (sidebar becomes overlay)
func (e *Engine) buildCompactLayout() {
	e.buildFocusLayout()
}

// AllPaneIDs returns the IDs of all panes in the current layout.
func (e *Engine) AllPaneIDs() []PaneID {
	if e.root == nil {
		return []PaneID{}
	}
	return e.root.AllPaneIDs()
}
