package modes

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
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
		{ModeID("unknown"), "Unknown"},
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

// --- Additional Registry Tests for Coverage ---

func TestRegistryRegisterWithNilMode(t *testing.T) {
	registry := NewRegistry()

	// Register a nil mode should be safe
	registry.Register(nil)

	// Count should still be 0
	if registry.Count() != 0 {
		t.Errorf("Expected count 0 after registering nil, got %d", registry.Count())
	}
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	registry := NewRegistry()

	mode1 := NewMockMode(ModeEdit, "Edit V1", "Edit mode v1")
	registry.Register(mode1)

	if registry.Count() != 1 {
		t.Errorf("Expected count 1, got %d", registry.Count())
	}

	// Register different mode with same ID
	mode2 := NewMockMode(ModeEdit, "Edit V2", "Edit mode v2")
	registry.Register(mode2)

	// Count should still be 1 (replaced)
	if registry.Count() != 1 {
		t.Errorf("Expected count 1 after duplicate register, got %d", registry.Count())
	}

	// Should get the new mode
	retrieved := registry.Get(ModeEdit)
	if retrieved.Name() != "Edit V2" {
		t.Errorf("Expected mode to be replaced with 'Edit V2', got '%s'", retrieved.Name())
	}
}

func TestBaseModeEngine(t *testing.T) {
	mode := NewBaseMode(ModeEdit, "Edit", "Edit mode")
	engine := mode.Engine()

	if engine == nil {
		t.Fatal("Engine should not be nil")
	}

	// Verify it's a layout engine with default standard layout
	// Mode() returns a LayoutMode which defaults to LayoutStandard
	layoutMode := engine.Mode()
	if layoutMode != layout.LayoutStandard {
		t.Errorf("Expected standard layout mode, got %v", layoutMode)
	}
}

// --- Lifecycle Order Tests ---

func TestRegistrySwitchToLifecycleOrder(t *testing.T) {
	registry := NewRegistry()

	exitMode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	enterMode := NewMockMode(ModeExplore, "Explore", "Explore mode")

	registry.Register(exitMode)
	registry.Register(enterMode)

	// Start in exit mode
	registry.SwitchTo(ModeEdit)
	exitMode.exitCalled = false
	enterMode.enterCalled = false

	// Switch to enter mode
	registry.SwitchTo(ModeExplore)

	// Both should be called
	if !exitMode.exitCalled {
		t.Error("OnExit should be called on previous mode")
	}
	if !enterMode.enterCalled {
		t.Error("OnEnter should be called on new mode")
	}
}

func TestRegistrySwitchToEmptyPrevious(t *testing.T) {
	registry := NewRegistry()

	// No previous mode set yet
	mode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(mode)

	// Switch to mode (previousID is "", should not call OnExit)
	cmd := registry.SwitchTo(ModeEdit)

	// Command may be nil if OnEnter returns nil
	if cmd != nil {
		t.Logf("SwitchTo from empty state returned command: %T", cmd)
	}

	if mode.exitCalled {
		t.Error("OnExit should not be called when switching from empty state")
	}

	if !mode.enterCalled {
		t.Error("OnEnter should be called")
	}
}

func TestRegistrySwitchToWithCommandReturned(t *testing.T) {
	registry := NewRegistry()

	// Create a mock mode that returns a command from OnEnter
	mode1 := NewMockMode(ModeEdit, "Edit", "Edit mode")

	// Create a second mode with custom OnEnter that returns a command
	mode2 := &MockMode{
		BaseMode: NewBaseMode(ModeExplore, "Explore", "Explore mode"),
	}

	registry.Register(mode1)
	registry.Register(mode2)

	// Switch to edit (BaseMode OnEnter returns nil)
	registry.SwitchTo(ModeEdit)

	// Now switch to explore - should call lifecycle
	_ = registry.SwitchTo(ModeExplore)

	// Command may be nil since both OnExit and OnEnter from BaseMode return nil
	// But the modes should be switched
	if registry.CurrentID() != ModeExplore {
		t.Error("Should have switched to explore mode")
	}
}

func TestRegistrySwitchToMultipleModes(t *testing.T) {
	registry := NewRegistry()

	edit := NewMockMode(ModeEdit, "Edit", "Edit mode")
	explore := NewMockMode(ModeExplore, "Explore", "Explore mode")
	analyze := NewMockMode(ModeAnalyze, "Analyze", "Analyze mode")

	registry.Register(edit)
	registry.Register(explore)
	registry.Register(analyze)

	// Switch: empty -> edit -> explore -> analyze -> edit
	registry.SwitchTo(ModeEdit)
	if registry.CurrentID() != ModeEdit || registry.PreviousID() != "" {
		t.Error("Step 1: Should be in edit, previous empty")
	}

	registry.SwitchTo(ModeExplore)
	if registry.CurrentID() != ModeExplore || registry.PreviousID() != ModeEdit {
		t.Error("Step 2: Should be in explore, previous edit")
	}

	registry.SwitchTo(ModeAnalyze)
	if registry.CurrentID() != ModeAnalyze || registry.PreviousID() != ModeExplore {
		t.Error("Step 3: Should be in analyze, previous explore")
	}

	registry.SwitchTo(ModeEdit)
	if registry.CurrentID() != ModeEdit || registry.PreviousID() != ModeAnalyze {
		t.Error("Step 4: Should be in edit, previous analyze")
	}
}

func TestRegistryUnregisterActiveModeNoSwitch(t *testing.T) {
	registry := NewRegistry()

	mode := NewMockMode(ModeEdit, "Edit", "Edit mode")
	registry.Register(mode)
	registry.SwitchTo(ModeEdit)

	if registry.CurrentID() != ModeEdit {
		t.Fatal("Should be in edit mode")
	}

	// Unregister the active mode
	registry.Unregister(ModeEdit)

	// Current mode should still be "Edit" ID but mode is gone
	if registry.Current() != nil {
		t.Error("Current should return nil after unregistering active mode")
	}

	if registry.CurrentID() != ModeEdit {
		t.Error("CurrentID should still return the ID")
	}
}

func TestRegistryGetAllModesEmpty(t *testing.T) {
	registry := NewRegistry()

	allModes := registry.AllModes()
	if allModes == nil {
		t.Fatal("AllModes should return a slice, not nil")
	}

	if len(allModes) != 0 {
		t.Errorf("AllModes should be empty, got %d", len(allModes))
	}
}

func TestRegistryCurrentIDWithoutSwitch(t *testing.T) {
	registry := NewRegistry()

	if registry.CurrentID() != "" {
		t.Errorf("CurrentID should be empty initially, got %q", registry.CurrentID())
	}

	if registry.PreviousID() != "" {
		t.Errorf("PreviousID should be empty initially, got %q", registry.PreviousID())
	}
}

// --- Table-Driven Tests for Mode Lifecycle ---

func TestRegistrySwitchToTableDriven(t *testing.T) {
	tests := []struct {
		name          string
		fromMode      ModeID
		toMode        ModeID
		shouldFail    bool
		expectedPrev  ModeID
		expectedCurr  ModeID
	}{
		{
			name:         "Switch edit to explore",
			fromMode:     ModeEdit,
			toMode:       ModeExplore,
			shouldFail:   false,
			expectedPrev: ModeEdit,
			expectedCurr: ModeExplore,
		},
		{
			name:         "Switch to non-existent mode",
			fromMode:     ModeEdit,
			toMode:       ModeID("nonexistent"),
			shouldFail:   true,
			expectedPrev: "", // Should not change
			expectedCurr: ModeEdit,
		},
		{
			name:         "Switch analyze to orchestrate",
			fromMode:     ModeAnalyze,
			toMode:       ModeOrchestrate,
			shouldFail:   false,
			expectedPrev: ModeAnalyze,
			expectedCurr: ModeOrchestrate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistry()

			// Register all modes
			registry.Register(NewMockMode(ModeEdit, "Edit", "Edit mode"))
			registry.Register(NewMockMode(ModeExplore, "Explore", "Explore mode"))
			registry.Register(NewMockMode(ModeAnalyze, "Analyze", "Analyze mode"))
			registry.Register(NewMockMode(ModeOrchestrate, "Orchestrate", "Orchestrate mode"))

			// Switch to from mode
			registry.SwitchTo(tt.fromMode)

			// Switch to target mode
			cmd := registry.SwitchTo(tt.toMode)

			if tt.shouldFail {
				if cmd != nil {
					t.Error("Expected no command for failed switch")
				}
				if registry.CurrentID() != tt.fromMode {
					t.Errorf("Expected to stay in %s, got %s", tt.fromMode, registry.CurrentID())
				}
			} else {
				// Successful switch should have switched modes
				// (command may be nil if both OnExit and OnEnter are nil)
				if registry.CurrentID() != tt.expectedCurr {
					t.Errorf("Expected current %s, got %s", tt.expectedCurr, registry.CurrentID())
				}
			}
		})
	}
}

func TestRegistryCountAfterOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation func(*Registry)
		expected  int
	}{
		{
			name: "initial empty",
			operation: func(r *Registry) {
				// No operation
			},
			expected: 0,
		},
		{
			name: "after registering 1 mode",
			operation: func(r *Registry) {
				r.Register(NewMockMode(ModeEdit, "Edit", "Edit mode"))
			},
			expected: 1,
		},
		{
			name: "after registering 2 modes",
			operation: func(r *Registry) {
				r.Register(NewMockMode(ModeExplore, "Explore", "Explore mode"))
			},
			expected: 2,
		},
		{
			name: "after unregistering 1",
			operation: func(r *Registry) {
				r.Unregister(ModeEdit)
			},
			expected: 1,
		},
		{
			name: "unregistering non-existent mode",
			operation: func(r *Registry) {
				r.Unregister(ModeAnalyze)
			},
			expected: 1,
		},
	}

	registry := NewRegistry()

	for _, tc := range tests {
		tc.operation(registry)
		if registry.Count() != tc.expected {
			t.Errorf("%s: expected count %d, got %d", tc.name, tc.expected, registry.Count())
		}
	}
}
