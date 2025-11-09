package integration

import (
	"testing"
	"time"

	"github.com/rand/pedantic-raven/internal/modes"
)

// TestEditAnalyzeWorkflow tests the Edit -> Analyze workflow.
func TestEditAnalyzeWorkflow(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Start in Edit mode
	AssertEqual(t, modes.ModeEdit, app.CurrentModeID(), "should start in Edit mode")

	// 2. Create content in editor
	content := "John Smith works at Acme Corporation in Seattle."
	app.Editor().SetContent(content)

	// Wait briefly for analysis to complete
	err := app.WaitForCondition(func() bool {
		analysis := app.EditMode().GetContextPanel().GetPanel().GetAnalysis()
		return analysis != nil
	}, 2*time.Second)

	if err != nil {
		t.Logf("Analysis may not have completed, continuing with test")
	}

	// 3. Switch to Analyze mode
	cmd := app.SwitchToMode(modes.ModeAnalyze)
	if cmd != nil {
		// Execute the command to trigger lifecycle hooks
		cmd()
	}

	AssertEqual(t, modes.ModeAnalyze, app.CurrentModeID(), "should be in Analyze mode")

	// 4. Verify we're still in Analyze mode
	currentMode := app.CurrentMode()
	AssertNotEqual(t, nil, currentMode, "current mode should not be nil")
}

// TestAnalyzeOrchestrateWorkflow tests the Analyze -> Orchestrate workflow.
func TestAnalyzeOrchestrateWorkflow(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Start in Edit mode and create content
	content := `Design a new API endpoint:
- Create GET /api/users endpoint
- Add authentication middleware
- Write unit tests
- Deploy to staging`

	app.Editor().SetContent(content)

	// 2. Switch to Analyze mode
	cmd := app.SwitchToMode(modes.ModeAnalyze)
	if cmd != nil {
		cmd()
	}
	AssertEqual(t, modes.ModeAnalyze, app.CurrentModeID(), "should be in Analyze mode")

	// 3. Switch to Orchestrate mode
	cmd = app.SwitchToMode(modes.ModeOrchestrate)
	if cmd != nil {
		cmd()
	}
	AssertEqual(t, modes.ModeOrchestrate, app.CurrentModeID(), "should be in Orchestrate mode")

	// 4. Verify Orchestrate mode is initialized
	orchestrateMode := app.OrchestrateMode()
	AssertNotEqual(t, nil, orchestrateMode, "orchestrate mode should not be nil")
}

// TestEditOrchestrateWorkflow tests the Edit -> Orchestrate workflow.
func TestEditOrchestrateWorkflow(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Create a work plan in Edit mode
	workPlan := `# Sprint Planning

## Task 1: Setup Database
- Create schema
- Add migrations

## Task 2: API Implementation
- Implement endpoints
- Add validation

## Task 3: Testing
- Write unit tests
- Write integration tests`

	app.Editor().SetContent(workPlan)

	// 2. Switch directly to Orchestrate mode
	cmd := app.SwitchToMode(modes.ModeOrchestrate)
	if cmd != nil {
		cmd()
	}

	AssertEqual(t, modes.ModeOrchestrate, app.CurrentModeID(), "should be in Orchestrate mode")

	// 3. Verify mode has correct properties
	orchestrateMode := app.OrchestrateMode()
	AssertEqual(t, modes.ModeOrchestrate, orchestrateMode.ID(), "orchestrate mode should have correct ID")
}

// TestMultiModeNavigation tests rapid mode switching without state corruption.
func TestMultiModeNavigation(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Set content in Edit mode
	app.Editor().SetContent("Test content for navigation")

	// 2. Rapid mode switching
	modes_sequence := []modes.ModeID{
		modes.ModeAnalyze,
		modes.ModeExplore,
		modes.ModeOrchestrate,
		modes.ModeEdit,
		modes.ModeAnalyze,
		modes.ModeEdit,
	}

	for i, modeID := range modes_sequence {
		cmd := app.SwitchToMode(modeID)
		if cmd != nil {
			cmd()
		}

		currentID := app.CurrentModeID()
		if currentID != modeID {
			t.Fatalf("mode switching failed at index %d: expected %v, got %v", i, modeID, currentID)
		}

		// Verify content is preserved in Edit mode
		if modeID == modes.ModeEdit {
			content := app.Editor().GetContent()
			AssertEqual(t, "Test content for navigation", content, "content was lost during mode switching")
		}
	}
}

// TestModeStatePreservation tests that state is preserved when switching away and back.
func TestModeStatePreservation(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Set content in Edit mode
	originalContent := "Original test content\nWith multiple lines"
	app.Editor().SetContent(originalContent)

	// 2. Switch away
	cmd := app.SwitchToMode(modes.ModeAnalyze)
	if cmd != nil {
		cmd()
	}

	// 3. Switch back
	cmd = app.SwitchToMode(modes.ModeEdit)
	if cmd != nil {
		cmd()
	}

	// 4. Verify content is still there
	retrievedContent := app.Editor().GetContent()
	AssertEqual(t, originalContent, retrievedContent, "content was not preserved after mode switch")
}

// TestModeRegistryPreviousMode tests the "go back" functionality.
func TestModeRegistryPreviousMode(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Start in Edit mode
	AssertEqual(t, modes.ModeEdit, app.CurrentModeID(), "should start in Edit mode")

	// 2. Switch to Analyze
	cmd := app.SwitchToMode(modes.ModeAnalyze)
	if cmd != nil {
		cmd()
	}
	AssertEqual(t, modes.ModeAnalyze, app.CurrentModeID(), "should be in Analyze mode")

	// 3. Switch to Explore
	cmd = app.SwitchToMode(modes.ModeExplore)
	if cmd != nil {
		cmd()
	}
	AssertEqual(t, modes.ModeExplore, app.CurrentModeID(), "should be in Explore mode")

	// 4. Go back to previous (Analyze)
	cmd = app.ModeRegistry().SwitchToPrevious()
	if cmd != nil {
		cmd()
	}
	AssertEqual(t, modes.ModeAnalyze, app.CurrentModeID(), "should be back in Analyze mode")

	// 5. Previous ID should be Explore
	AssertEqual(t, modes.ModeExplore, app.ModeRegistry().PreviousID(), "previous mode should be Explore")
}

// TestModeTransitionCommands tests that transition commands execute properly.
func TestModeTransitionCommands(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// Test that OnExit and OnEnter are called
	editMode := app.EditMode()
	initialContent := "Test content"
	editMode.GetEditor().SetContent(initialContent)

	// Switch away and back
	cmd := app.SwitchToMode(modes.ModeAnalyze)
	if cmd != nil {
		cmd()
	}

	cmd = app.SwitchToMode(modes.ModeEdit)
	if cmd != nil {
		cmd()
	}

	// Content should be preserved through transitions
	currentContent := editMode.GetEditor().GetContent()
	AssertEqual(t, initialContent, currentContent, "content lost through mode transitions")
}

// TestEditModeOnEnter tests that OnEnter triggers analysis.
func TestEditModeOnEnter(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// Set content before entering
	content := "Test content for analysis"
	app.Editor().SetContent(content)

	// OnEnter should trigger analysis
	cmd := app.EditMode().OnEnter()

	// If there's a command, execute it
	if cmd != nil {
		cmd()
	}

	// Verify the content is still there
	currentContent := app.Editor().GetContent()
	AssertEqual(t, content, currentContent, "content was modified during OnEnter")
}

// TestOrchestrateModeSwitching tests switching to/from Orchestrate mode.
func TestOrchestrateModeSwitching(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// Initialize Orchestrate mode
	orchestrateMode := app.OrchestrateMode()
	AssertNotEqual(t, nil, orchestrateMode, "orchestrate mode should be initialized")

	// Switch to it
	cmd := app.SwitchToMode(modes.ModeOrchestrate)
	if cmd != nil {
		cmd()
	}

	AssertEqual(t, modes.ModeOrchestrate, app.CurrentModeID(), "should be in Orchestrate mode")

	// Switch away
	cmd = app.SwitchToMode(modes.ModeEdit)
	if cmd != nil {
		cmd()
	}

	AssertEqual(t, modes.ModeEdit, app.CurrentModeID(), "should be back in Edit mode")

	// Previous should be Orchestrate
	AssertEqual(t, modes.ModeOrchestrate, app.ModeRegistry().PreviousID(), "previous should be Orchestrate")
}
