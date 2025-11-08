// Package modes provides the mode registry and application mode management.
//
// Pedantic Raven has 5 primary modes:
// - Edit: Context editing with semantic analysis (ICS-like)
// - Explore: Memory workspace with graph visualization
// - Analyze: Semantic insights and triple analysis
// - Orchestrate: Multi-agent coordination and task management
// - Collaborate: Live multi-user editing
//
// Each mode has its own UI layout, keybindings, and functionality.
package modes

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
)

// ModeID uniquely identifies an application mode.
type ModeID string

const (
	ModeEdit        ModeID = "edit"
	ModeExplore     ModeID = "explore"
	ModeAnalyze     ModeID = "analyze"
	ModeOrchestrate ModeID = "orchestrate"
	ModeCollaborate ModeID = "collaborate"
)

// String returns the human-readable name of the mode.
func (m ModeID) String() string {
	switch m {
	case ModeEdit:
		return "Edit"
	case ModeExplore:
		return "Explore"
	case ModeAnalyze:
		return "Analyze"
	case ModeOrchestrate:
		return "Orchestrate"
	case ModeCollaborate:
		return "Collaborate"
	default:
		return "Unknown"
	}
}

// Mode represents an application mode with its own state and behavior.
//
// Each mode manages its own:
// - State (model data)
// - Layout (which components are visible and how they're arranged)
// - Event handling (mode-specific key bindings and interactions)
// - Lifecycle (initialization, activation, deactivation)
type Mode interface {
	// ID returns the unique identifier for this mode.
	ID() ModeID

	// Name returns the human-readable name for this mode.
	Name() string

	// Description returns a brief description of what this mode does.
	Description() string

	// Init initializes the mode and returns an initial command.
	// Called once when the mode is first created.
	Init() tea.Cmd

	// OnEnter is called when this mode becomes active.
	// Use this to start background tasks, subscribe to events, etc.
	OnEnter() tea.Cmd

	// OnExit is called when this mode becomes inactive.
	// Use this to clean up resources, unsubscribe from events, etc.
	OnExit() tea.Cmd

	// Update processes a Bubble Tea message and returns updated state.
	Update(msg tea.Msg) (Mode, tea.Cmd)

	// View renders the mode to a string.
	// The mode is responsible for rendering its layout and components.
	View() string

	// Keybindings returns a description of available keybindings for this mode.
	// Used for help/documentation.
	Keybindings() []Keybinding
}

// Keybinding describes a keyboard shortcut and its action.
type Keybinding struct {
	Key         string
	Description string
}

// Registry manages the collection of available modes and handles mode switching.
type Registry struct {
	modes      map[ModeID]Mode
	currentID  ModeID
	previousID ModeID
}

// NewRegistry creates a new mode registry.
func NewRegistry() *Registry {
	return &Registry{
		modes:      make(map[ModeID]Mode),
		currentID:  "",
		previousID: "",
	}
}

// Register adds a mode to the registry.
//
// If a mode with the same ID already exists, it will be replaced.
func (r *Registry) Register(mode Mode) {
	if mode != nil {
		r.modes[mode.ID()] = mode
	}
}

// Unregister removes a mode from the registry.
func (r *Registry) Unregister(id ModeID) {
	delete(r.modes, id)
}

// Get retrieves a mode by ID.
// Returns nil if the mode doesn't exist.
func (r *Registry) Get(id ModeID) Mode {
	return r.modes[id]
}

// Current returns the currently active mode.
// Returns nil if no mode is active.
func (r *Registry) Current() Mode {
	return r.modes[r.currentID]
}

// CurrentID returns the ID of the currently active mode.
func (r *Registry) CurrentID() ModeID {
	return r.currentID
}

// PreviousID returns the ID of the previously active mode.
// Useful for "go back" functionality.
func (r *Registry) PreviousID() ModeID {
	return r.previousID
}

// SwitchTo changes the active mode.
//
// This calls OnExit on the current mode and OnEnter on the new mode.
// Returns a command that batches both lifecycle commands.
//
// If the specified mode doesn't exist, this is a no-op and returns nil.
func (r *Registry) SwitchTo(id ModeID) tea.Cmd {
	// Check if mode exists
	newMode := r.modes[id]
	if newMode == nil {
		return nil
	}

	// No-op if already in this mode
	if r.currentID == id {
		return nil
	}

	var cmds []tea.Cmd

	// Call OnExit on current mode
	if r.currentID != "" {
		if currentMode := r.modes[r.currentID]; currentMode != nil {
			if cmd := currentMode.OnExit(); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	// Update state
	r.previousID = r.currentID
	r.currentID = id

	// Call OnEnter on new mode
	if cmd := newMode.OnEnter(); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// SwitchToPrevious switches back to the previously active mode.
//
// This is useful for "go back" or "escape to previous mode" functionality.
func (r *Registry) SwitchToPrevious() tea.Cmd {
	if r.previousID != "" {
		return r.SwitchTo(r.previousID)
	}
	return nil
}

// AllModes returns a list of all registered mode IDs.
func (r *Registry) AllModes() []ModeID {
	ids := make([]ModeID, 0, len(r.modes))
	for id := range r.modes {
		ids = append(ids, id)
	}
	return ids
}

// Count returns the number of registered modes.
func (r *Registry) Count() int {
	return len(r.modes)
}

// --- Base Mode Implementation ---

// BaseMode provides a default implementation of the Mode interface.
//
// Concrete modes can embed this to get default implementations
// and only override the methods they need to customize.
type BaseMode struct {
	id          ModeID
	name        string
	description string
	engine      *layout.Engine
}

// NewBaseMode creates a new base mode with the given parameters.
func NewBaseMode(id ModeID, name, description string) *BaseMode {
	return &BaseMode{
		id:          id,
		name:        name,
		description: description,
		engine:      layout.NewEngine(layout.LayoutStandard),
	}
}

// ID implements Mode.
func (m *BaseMode) ID() ModeID {
	return m.id
}

// Name implements Mode.
func (m *BaseMode) Name() string {
	return m.name
}

// Description implements Mode.
func (m *BaseMode) Description() string {
	return m.description
}

// Init implements Mode.
func (m *BaseMode) Init() tea.Cmd {
	return m.engine.Init()
}

// OnEnter implements Mode.
func (m *BaseMode) OnEnter() tea.Cmd {
	return nil
}

// OnExit implements Mode.
func (m *BaseMode) OnExit() tea.Cmd {
	return nil
}

// Update implements Mode.
func (m *BaseMode) Update(msg tea.Msg) (Mode, tea.Cmd) {
	_, cmd := m.engine.Update(msg)
	return m, cmd
}

// View implements Mode.
func (m *BaseMode) View() string {
	return m.engine.View()
}

// Keybindings implements Mode.
func (m *BaseMode) Keybindings() []Keybinding {
	return []Keybinding{
		{"q", "Quit"},
		{"?", "Show help"},
	}
}

// Engine returns the layout engine for this mode.
// This allows concrete modes to customize their layout.
func (m *BaseMode) Engine() *layout.Engine {
	return m.engine
}
