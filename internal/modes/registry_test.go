package modes

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// MockMode is a test mode that tracks lifecycle calls.
type MockMode struct {
	*BaseMode
	initCalled    bool
	enterCalled   bool
	exitCalled    bool
	updateCalled  bool
}

func NewMockMode(id ModeID, name, description string) *MockMode {
	return &MockMode{
		BaseMode: NewBaseMode(id, name, description),
	}
}

func (m *MockMode) Init() tea.Cmd {
	m.initCalled = true
	return m.BaseMode.Init()
}

func (m *MockMode) OnEnter() tea.Cmd {
	m.enterCalled = true
	return m.BaseMode.OnEnter()
}

func (m *MockMode) OnExit() tea.Cmd {
	m.exitCalled = true
	return m.BaseMode.OnExit()
}

func (m *MockMode) Update(msg tea.Msg) (Mode, tea.Cmd) {
	m.updateCalled = true
	return m, nil
}

// --- Registry Tests ---

func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()

	mode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(mode)

	retrieved := registry.Get(ModeEdit)
	if retrieved == nil {
		t.Fatal("Mode should be registered")
	}

	if retrieved.ID() != ModeEdit {
		t.Errorf("Expected mode ID %s, got %s", ModeEdit, retrieved.ID())
	}
}

func TestRegistryUnregister(t *testing.T) {
	registry := NewRegistry()

	mode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(mode)

	// Verify it's registered
	if registry.Get(ModeEdit) == nil {
		t.Fatal("Mode should be registered")
	}

	// Unregister
	registry.Unregister(ModeEdit)

	// Verify it's gone
	if registry.Get(ModeEdit) != nil {
		t.Error("Mode should be unregistered")
	}
}

func TestRegistryGet(t *testing.T) {
	registry := NewRegistry()

	mode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(mode)

	// Get existing mode
	retrieved := registry.Get(ModeEdit)
	if retrieved == nil {
		t.Fatal("Should retrieve registered mode")
	}

	// Get non-existent mode
	retrieved = registry.Get(ModeExplore)
	if retrieved != nil {
		t.Error("Should return nil for non-existent mode")
	}
}

func TestRegistrySwitchTo(t *testing.T) {
	registry := NewRegistry()

	editMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	exploreMode := NewMockMode(ModeExplore, "Explore", "Explore mode")

	registry.Register(editMode)
	registry.Register(exploreMode)

	// Switch to edit mode
	registry.SwitchTo(ModeEdit)

	if registry.CurrentID() != ModeEdit {
		t.Errorf("Expected current mode %s, got %s", ModeEdit, registry.CurrentID())
	}

	if !editMode.enterCalled {
		t.Error("OnEnter should be called when switching to edit mode")
	}

	// Switch to explore mode
	editMode.exitCalled = false // Reset flag
	registry.SwitchTo(ModeExplore)

	if registry.CurrentID() != ModeExplore {
		t.Errorf("Expected current mode %s, got %s", ModeExplore, registry.CurrentID())
	}

	if !editMode.exitCalled {
		t.Error("OnExit should be called on edit mode when switching away")
	}

	if !exploreMode.enterCalled {
		t.Error("OnEnter should be called when switching to explore mode")
	}

	if registry.PreviousID() != ModeEdit {
		t.Errorf("Expected previous mode %s, got %s", ModeEdit, registry.PreviousID())
	}
}

func TestRegistrySwitchToNonExistent(t *testing.T) {
	registry := NewRegistry()

	editMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(editMode)
	registry.SwitchTo(ModeEdit)

	// Try to switch to non-existent mode
	cmd := registry.SwitchTo(ModeAnalyze)

	// Should be no-op
	if cmd != nil {
		t.Error("Switching to non-existent mode should return nil command")
	}

	// Should remain in edit mode
	if registry.CurrentID() != ModeEdit {
		t.Errorf("Should remain in edit mode, got %s", registry.CurrentID())
	}
}

func TestRegistrySwitchToSameMode(t *testing.T) {
	registry := NewRegistry()

	editMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(editMode)
	registry.SwitchTo(ModeEdit)

	editMode.enterCalled = false // Reset flag

	// Try to switch to same mode
	cmd := registry.SwitchTo(ModeEdit)

	// Should be no-op
	if cmd != nil {
		t.Error("Switching to same mode should return nil command")
	}

	// OnEnter should not be called again
	if editMode.enterCalled {
		t.Error("OnEnter should not be called when already in mode")
	}
}

func TestRegistrySwitchToPrevious(t *testing.T) {
	registry := NewRegistry()

	editMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	exploreMode := NewMockMode(ModeExplore, "Explore", "Explore mode")

	registry.Register(editMode)
	registry.Register(exploreMode)

	// Switch edit → explore
	registry.SwitchTo(ModeEdit)
	registry.SwitchTo(ModeExplore)

	// Go back to edit
	registry.SwitchToPrevious()

	if registry.CurrentID() != ModeEdit {
		t.Errorf("Expected to switch back to %s, got %s", ModeEdit, registry.CurrentID())
	}

	if registry.PreviousID() != ModeExplore {
		t.Errorf("Expected previous mode %s, got %s", ModeExplore, registry.PreviousID())
	}
}

func TestRegistrySwitchToPreviousWhenNone(t *testing.T) {
	registry := NewRegistry()

	editMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(editMode)
	registry.SwitchTo(ModeEdit)

	// No previous mode
	cmd := registry.SwitchToPrevious()

	// Should be no-op
	if cmd != nil {
		t.Error("SwitchToPrevious with no previous should return nil")
	}

	// Should remain in edit mode
	if registry.CurrentID() != ModeEdit {
		t.Errorf("Should remain in edit mode, got %s", registry.CurrentID())
	}
}

func TestRegistryCurrent(t *testing.T) {
	registry := NewRegistry()

	editMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(editMode)

	// No current mode yet
	current := registry.Current()
	if current != nil {
		t.Error("Current should be nil when no mode is active")
	}

	// Switch to edit
	registry.SwitchTo(ModeEdit)

	current = registry.Current()
	if current == nil {
		t.Fatal("Current should return the active mode")
	}

	if current.ID() != ModeEdit {
		t.Errorf("Expected current mode %s, got %s", ModeEdit, current.ID())
	}
}

func TestRegistryAllModes(t *testing.T) {
	registry := NewRegistry()

	editMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	exploreMode := NewMockMode(ModeExplore, "Explore", "Explore mode")
	analyzeMode := NewMockMode(ModeAnalyze, "Analyze", "Analyze mode")

	registry.Register(editMode)
	registry.Register(exploreMode)
	registry.Register(analyzeMode)

	allModes := registry.AllModes()
	if len(allModes) != 3 {
		t.Fatalf("Expected 3 modes, got %d", len(allModes))
	}

	// Check all modes are present (order doesn't matter)
	hasEdit := false
	hasExplore := false
	hasAnalyze := false

	for _, id := range allModes {
		switch id {
		case ModeEdit:
			hasEdit = true
		case ModeExplore:
			hasExplore = true
		case ModeAnalyze:
			hasAnalyze = true
		}
	}

	if !hasEdit || !hasExplore || !hasAnalyze {
		t.Errorf("Not all modes present. Edit: %v, Explore: %v, Analyze: %v",
			hasEdit, hasExplore, hasAnalyze)
	}
}

func TestRegistryCount(t *testing.T) {
	registry := NewRegistry()

	if registry.Count() != 0 {
		t.Errorf("Expected count 0, got %d", registry.Count())
	}

	registry.Register(NewMockMode(ModeEdit, "Edit", "Edit mode"))
	if registry.Count() != 1 {
		t.Errorf("Expected count 1, got %d", registry.Count())
	}

	registry.Register(NewMockMode(ModeExplore, "Explore", "Explore mode"))
	if registry.Count() != 2 {
		t.Errorf("Expected count 2, got %d", registry.Count())
	}

	registry.Unregister(ModeEdit)
	if registry.Count() != 1 {
		t.Errorf("Expected count 1 after unregister, got %d", registry.Count())
	}
}

// --- Base Mode Tests ---

func TestBaseModeAttributes(t *testing.T) {
	mode := NewBaseMode(ModeEdit, "Edit Mode", "Context editing mode")

	if mode.ID() != ModeEdit {
		t.Errorf("Expected ID %s, got %s", ModeEdit, mode.ID())
	}

	if mode.Name() != "Edit Mode" {
		t.Errorf("Expected name 'Edit Mode', got '%s'", mode.Name())
	}

	if mode.Description() != "Context editing mode" {
		t.Errorf("Expected description 'Context editing mode', got '%s'", mode.Description())
	}
}

func TestBaseModeLifecycle(t *testing.T) {
	mode := NewBaseMode(ModeEdit, "Edit", "Edit mode")

	// Init should not panic
	cmd := mode.Init()
	_ = cmd

	// OnEnter should return nil by default
	cmd = mode.OnEnter()
	if cmd != nil {
		t.Error("Default OnEnter should return nil")
	}

	// OnExit should return nil by default
	cmd = mode.OnExit()
	if cmd != nil {
		t.Error("Default OnExit should return nil")
	}
}

func TestBaseModeUpdate(t *testing.T) {
	mode := NewBaseMode(ModeEdit, "Edit", "Edit mode")

	// Update should not panic
	updatedMode, cmd := mode.Update(tea.KeyMsg{})
	_ = cmd

	// Should return same mode
	if updatedMode.ID() != mode.ID() {
		t.Error("Update should return same mode")
	}
}

func TestBaseModeView(t *testing.T) {
	mode := NewBaseMode(ModeEdit, "Edit", "Edit mode")

	// View should not panic
	view := mode.View()
	_ = view
}

func TestBaseModeKeybindings(t *testing.T) {
	mode := NewBaseMode(ModeEdit, "Edit", "Edit mode")

	keybindings := mode.Keybindings()
	if len(keybindings) == 0 {
		t.Error("Should have default keybindings")
	}
}

func TestModeIDString(t *testing.T) {
	tests := []struct {
		id       ModeID
		expected string
	}{
		{ModeEdit, "Edit"},
		{ModeExplore, "Explore"},
		{ModeAnalyze, "Analyze"},
		{ModeOrchestrate, "Orchestrate"},
		{ModeCollaborate, "Collaborate"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.id.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
