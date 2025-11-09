// Package orchestrate provides the Orchestrate Mode implementation.
package orchestrate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/modes"
)

// ModeAdapter wraps OrchestrateMode to implement the modes.Mode interface.
//
// This adapter bridges between the OrchestrateMode Bubble Tea model and
// the application's mode registry system.
type ModeAdapter struct {
	model *OrchestrateMode
}

// NewModeAdapter creates a new mode adapter for Orchestrate Mode.
func NewModeAdapter() modes.Mode {
	return &ModeAdapter{
		model: NewOrchestrateMode(),
	}
}

// ID implements modes.Mode.
func (m *ModeAdapter) ID() modes.ModeID {
	return modes.ModeOrchestrate
}

// Name implements modes.Mode.
func (m *ModeAdapter) Name() string {
	return "Orchestrate"
}

// Description implements modes.Mode.
func (m *ModeAdapter) Description() string {
	return "Multi-agent coordination and task management using mnemosyne orchestration"
}

// Init implements modes.Mode.
func (m *ModeAdapter) Init() tea.Cmd {
	return m.model.Init()
}

// OnEnter implements modes.Mode.
//
// Called when entering Orchestrate Mode. This can be used to start
// background processes or load session history.
func (m *ModeAdapter) OnEnter() tea.Cmd {
	// Future enhancement: Load last session or start fresh
	return nil
}

// OnExit implements modes.Mode.
//
// Called when exiting Orchestrate Mode. This ensures orchestration
// is stopped cleanly and session state is saved.
func (m *ModeAdapter) OnExit() tea.Cmd {
	// Stop any running orchestration
	if m.model.IsOrchestrating() {
		_ = m.model.stopOrchestration()
	}

	// Save session state
	if session := m.model.GetSession(); session != nil {
		_ = session.Save()
	}

	return nil
}

// Update implements modes.Mode.
func (m *ModeAdapter) Update(msg tea.Msg) (modes.Mode, tea.Cmd) {
	// Handle mode-level quit messages
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "q" && !m.model.IsOrchestrating() {
			// Exit to previous mode (handled by main app)
			return m, nil
		}
	}

	// Delegate to underlying model
	updatedModel, cmd := m.model.Update(msg)
	m.model = updatedModel.(*OrchestrateMode)

	return m, cmd
}

// View implements modes.Mode.
func (m *ModeAdapter) View() string {
	return m.model.View()
}

// Keybindings implements modes.Mode.
func (m *ModeAdapter) Keybindings() []modes.Keybinding {
	return []modes.Keybinding{
		// Global
		{"q", "Exit mode (stops orchestration)"},
		{"?", "Toggle help"},

		// View switching
		{"Tab", "Next view"},
		{"Shift+Tab", "Previous view"},
		{"1", "Plan Editor view"},
		{"2", "Dashboard view"},
		{"3", "Task Graph view"},
		{"4", "Agent Log view"},

		// Plan editor (when in editor view)
		{"Ctrl+N", "New plan"},
		{"Ctrl+O", "Open plan file"},
		{"Ctrl+S", "Save plan"},
		{"Ctrl+L", "Launch orchestration"},

		// Orchestration controls (when running)
		{"Space", "Pause/Resume orchestration"},
		{"r", "Restart orchestration"},
		{"x", "Cancel orchestration"},

		// Dashboard/Graph/Log (when in respective views)
		{"+/-", "Zoom (Task Graph)"},
		{"h/j/k/l", "Pan (Task Graph)"},
		{"/", "Search (Agent Log)"},
		{"e", "Export logs"},
	}
}

// GetModel returns the underlying OrchestrateMode model.
//
// This is useful for testing or advanced integration scenarios
// where direct access to the model is needed.
func (m *ModeAdapter) GetModel() *OrchestrateMode {
	return m.model
}
