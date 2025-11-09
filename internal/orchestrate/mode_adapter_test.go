package orchestrate

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/modes"
	"github.com/stretchr/testify/assert"
)

// TestNewModeAdapter tests creation of mode adapter.
func TestNewModeAdapter(t *testing.T) {
	adapter := NewModeAdapter()

	assert.NotNil(t, adapter)
	assert.Equal(t, modes.ModeOrchestrate, adapter.ID())
	assert.Equal(t, "Orchestrate", adapter.Name())
	assert.Contains(t, adapter.Description(), "Multi-agent")
}

// TestModeAdapter_ID tests mode ID.
func TestModeAdapter_ID(t *testing.T) {
	adapter := NewModeAdapter()
	assert.Equal(t, modes.ModeOrchestrate, adapter.ID())
}

// TestModeAdapter_Name tests mode name.
func TestModeAdapter_Name(t *testing.T) {
	adapter := NewModeAdapter()
	assert.Equal(t, "Orchestrate", adapter.Name())
}

// TestModeAdapter_Description tests mode description.
func TestModeAdapter_Description(t *testing.T) {
	adapter := NewModeAdapter()

	desc := adapter.Description()
	assert.NotEmpty(t, desc)
	assert.Contains(t, desc, "Multi-agent")
	assert.Contains(t, desc, "mnemosyne")
}

// TestModeAdapter_Init tests initialization.
func TestModeAdapter_Init(t *testing.T) {
	adapter := NewModeAdapter()

	cmd := adapter.Init()
	// Init may return nil or a command
	_ = cmd
}

// TestModeAdapter_OnEnter tests entering mode.
func TestModeAdapter_OnEnter(t *testing.T) {
	adapter := NewModeAdapter()

	cmd := adapter.OnEnter()
	// OnEnter currently returns nil (future enhancement)
	assert.Nil(t, cmd)
}

// TestModeAdapter_OnExit tests exiting mode.
func TestModeAdapter_OnExit(t *testing.T) {
	adapter := NewModeAdapter()

	// OnExit should stop orchestration if running
	cmd := adapter.OnExit()

	// Should return nil (cleanup is synchronous)
	assert.Nil(t, cmd)
}

// TestModeAdapter_OnExitWithOrchestration tests exit with active orchestration.
func TestModeAdapter_OnExitWithOrchestration(t *testing.T) {
	adapter := NewModeAdapter().(*ModeAdapter)

	// Set plan in model (without actually launching)
	adapter.model.planEditor.content = `{"name":"Test Plan","description":"Test","maxConcurrent":2,"tasks":[{"id":"task1","description":"Task 1","type":1,"dependencies":[]}]}`
	adapter.model.planEditor.validateContent()
	adapter.model.orchestrating = true

	// OnExit should stop orchestration
	cmd := adapter.OnExit()
	assert.Nil(t, cmd)

	// Orchestration should be stopped
	assert.False(t, adapter.model.IsOrchestrating())
}

// TestModeAdapter_Update tests message handling.
func TestModeAdapter_Update(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.Msg
	}{
		{
			name: "Window resize",
			msg:  tea.WindowSizeMsg{Width: 100, Height: 50},
		},
		{
			name: "Key press (Tab)",
			msg:  tea.KeyMsg{Type: tea.KeyTab},
		},
		{
			name: "Key press (1)",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewModeAdapter()

			updatedMode, cmd := adapter.Update(tt.msg)

			assert.NotNil(t, updatedMode)
			assert.Equal(t, modes.ModeOrchestrate, updatedMode.ID())
			_ = cmd
		})
	}
}

// TestModeAdapter_UpdateQuit tests quit handling.
func TestModeAdapter_UpdateQuit(t *testing.T) {
	adapter := NewModeAdapter()

	// Quit when not orchestrating should pass through
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedMode, cmd := adapter.Update(quitMsg)

	assert.NotNil(t, updatedMode)
	assert.Equal(t, modes.ModeOrchestrate, updatedMode.ID())
	_ = cmd
}

// TestModeAdapter_View tests rendering.
func TestModeAdapter_View(t *testing.T) {
	adapter := NewModeAdapter()

	// Initialize size
	adapter.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

	view := adapter.View()

	assert.NotEmpty(t, view)
	// Header shows "ORCHESTRATE MODE" in all caps
	assert.Contains(t, view, "ORCHESTRATE MODE")
}

// TestModeAdapter_Keybindings tests keybinding documentation.
func TestModeAdapter_Keybindings(t *testing.T) {
	adapter := NewModeAdapter()

	keybindings := adapter.Keybindings()

	assert.NotEmpty(t, keybindings)

	// Check for essential keybindings
	keys := make(map[string]string)
	for _, kb := range keybindings {
		keys[kb.Key] = kb.Description
	}

	assert.Contains(t, keys, "q")
	assert.Contains(t, keys, "?")
	assert.Contains(t, keys, "Tab")
	assert.Contains(t, keys, "1")
	assert.Contains(t, keys, "Ctrl+L")
	assert.Contains(t, keys, "Space")
}

// TestModeAdapter_GetModel tests model accessor.
func TestModeAdapter_GetModel(t *testing.T) {
	adapter := NewModeAdapter().(*ModeAdapter)

	model := adapter.GetModel()
	assert.NotNil(t, model)
	assert.IsType(t, &OrchestrateMode{}, model)
}

// TestModeAdapter_InterfaceImplementation tests Mode interface.
func TestModeAdapter_InterfaceImplementation(t *testing.T) {
	var _ modes.Mode = (*ModeAdapter)(nil)
}

// TestModeAdapter_Lifecycle tests full lifecycle.
func TestModeAdapter_Lifecycle(t *testing.T) {
	adapter := NewModeAdapter()

	// Init
	initCmd := adapter.Init()
	_ = initCmd

	// Enter
	enterCmd := adapter.OnEnter()
	_ = enterCmd

	// Simulate resize
	_, _ = adapter.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Simulate key press
	_, _ = adapter.Update(tea.KeyMsg{Type: tea.KeyTab})

	// View
	view := adapter.View()
	assert.NotEmpty(t, view)

	// Exit
	exitCmd := adapter.OnExit()
	_ = exitCmd
}

// TestModeAdapter_SessionPersistence tests session handling on exit.
func TestModeAdapter_SessionPersistence(t *testing.T) {
	adapter := NewModeAdapter().(*ModeAdapter)

	// Create a valid plan
	plan := &WorkPlan{
		Name:          "Test Plan",
		Description:   "Test session persistence",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:           "task1",
				Description:  "Task 1",
				Type:         TaskTypeSequential,
				Dependencies: []string{},
			},
		},
	}

	// Initialize model with session
	session := NewSession(plan)
	session.historyDir = t.TempDir()
	adapter.model.session = session

	// Modify session state
	session.state.CompletedTasks = 5
	session.state.FailedTasks = 1

	// OnExit should save session
	cmd := adapter.OnExit()
	assert.Nil(t, cmd)

	// Session should be saved (file exists)
	// Note: Actual file check would require session ID
}

// TestModeAdapter_ViewSwitching tests view navigation through adapter.
func TestModeAdapter_ViewSwitching(t *testing.T) {
	adapter := NewModeAdapter().(*ModeAdapter)

	// Set orchestrating to allow view switching
	adapter.model.orchestrating = true
	adapter.model.currentView = ViewPlanEditor

	// Switch to dashboard (key '2')
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
	updatedMode, _ := adapter.Update(msg)

	assert.NotNil(t, updatedMode)
	modeAdapter := updatedMode.(*ModeAdapter)
	assert.Equal(t, ViewDashboard, modeAdapter.model.currentView)
}
