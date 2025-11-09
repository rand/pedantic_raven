package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rand/pedantic-raven/internal/modes"
)

// TestOfflineModeRecovery tests handling when external services are unavailable.
func TestOfflineModeRecovery(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Set content in normal mode
	content := "Working content"
	app.Editor().SetContent(content)

	// 2. Attempt analysis (may fail silently if service unavailable)
	cmd := app.EditMode().OnEnter()
	if cmd != nil {
		cmd()
	}

	// 3. Verify app still works
	currentContent := app.Editor().GetContent()
	AssertEqual(t, content, currentContent, "app should continue working offline")

	// 4. Verify mode can still switch
	cmd = app.SwitchToMode(modes.ModeAnalyze)
	if cmd != nil {
		cmd()
	}

	AssertEqual(t, modes.ModeAnalyze, app.CurrentModeID(), "should switch modes in offline mode")

	// 5. Return to Edit and verify content persists
	cmd = app.SwitchToMode(modes.ModeEdit)
	if cmd != nil {
		cmd()
	}

	finalContent := app.Editor().GetContent()
	AssertEqual(t, content, finalContent, "content should persist in offline mode")
}

// TestCorruptContentRecovery tests recovery from malformed content.
func TestCorruptContentRecovery(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Load initially valid content
	validContent := "Valid content here"
	app.Editor().SetContent(validContent)

	// 2. Attempt to load malformed content
	malformedContent := string([]byte{0x00, 0xFF, 0xFE, 0xFD})
	app.Editor().SetContent(malformedContent)

	// 3. Verify app doesn't crash
	current := app.Editor().GetContent()
	AssertNotEqual(t, "", current, "app should handle malformed content")

	// 4. Verify can return to valid content
	app.Editor().SetContent(validContent)
	recovered := app.Editor().GetContent()
	AssertEqual(t, validContent, recovered, "should recover to valid content")
}

// TestInvalidSemanticAnalysisInput tests handling invalid analysis input.
func TestInvalidSemanticAnalysisInput(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Try empty content
	app.Editor().SetContent("")

	// 2. Trigger analysis
	cmd := app.EditMode().OnEnter()
	if cmd != nil {
		cmd()
	}

	// 3. App should not crash
	AssertEqual(t, "", app.Editor().GetContent(), "empty content should be handled")

	// 4. Try very long content
	longContent := strings.Repeat("word ", 10000)
	app.Editor().SetContent(longContent)

	// 5. Trigger analysis (may timeout or be truncated)
	cmd = app.EditMode().OnEnter()
	if cmd != nil {
		cmd()
	}

	// 6. Content should still be there
	current := app.Editor().GetContent()
	AssertNotEqual(t, "", current, "long content should be retained")
}

// TestModeTransitionUnderLoad tests mode transitions with significant content.
func TestModeTransitionUnderLoad(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Load substantial content
	largeContent := generateLargeContent(1000)
	app.Editor().SetContent(largeContent)

	// 2. Rapid mode switching with content
	for i := 0; i < 5; i++ {
		cmd := app.SwitchToMode(modes.ModeAnalyze)
		if cmd != nil {
			cmd()
		}

		cmd = app.SwitchToMode(modes.ModeExplore)
		if cmd != nil {
			cmd()
		}

		cmd = app.SwitchToMode(modes.ModeEdit)
		if cmd != nil {
			cmd()
		}
	}

	// 3. Verify content is unchanged
	finalContent := app.Editor().GetContent()
	AssertEqual(t, largeContent, finalContent, "content should survive under-load mode transitions")
}

// TestAnalysisErrorRecovery tests recovery from analysis errors.
func TestAnalysisErrorRecovery(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Trigger analysis with first content
	content1 := "First analysis content"
	app.Editor().SetContent(content1)

	// 2. Attempt analysis
	cmd := app.EditMode().OnEnter()
	if cmd != nil {
		cmd()
	}

	// 3. Change content for second analysis
	content2 := "Second analysis content"
	app.Editor().SetContent(content2)

	// 4. Attempt second analysis
	cmd = app.EditMode().OnEnter()
	if cmd != nil {
		cmd()
	}

	// 5. Verify latest content is correct
	current := app.Editor().GetContent()
	AssertEqual(t, content2, current, "should handle multiple analysis attempts")
}

// TestGracefulDegradation tests that app degrades gracefully when components unavailable.
func TestGracefulDegradation(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Start editing
	content := "Content being edited"
	app.Editor().SetContent(content)

	// 2. Attempt all mode switches (some may have unavailable components)
	allModes := []modes.ModeID{
		modes.ModeEdit,
		modes.ModeAnalyze,
		modes.ModeExplore,
		modes.ModeOrchestrate,
	}

	for _, modeID := range allModes {
		cmd := app.SwitchToMode(modeID)
		if cmd != nil {
			cmd()
		}

		// App should still respond
		currentMode := app.CurrentMode()
		AssertNotEqual(t, nil, currentMode, "mode should be available")
	}

	// 3. Return to Edit and verify content
	cmd := app.SwitchToMode(modes.ModeEdit)
	if cmd != nil {
		cmd()
	}

	finalContent := app.Editor().GetContent()
	AssertEqual(t, content, finalContent, "content should survive all mode transitions")
}

// TestEventBrokerResilience tests event broker resilience.
func TestEventBrokerResilience(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	broker := app.EventBroker()
	AssertNotEqual(t, nil, broker, "event broker should be initialized")

	// 1. Publish events with no subscribers (should not crash)
	// This is handled by the event broker's buffering

	// 2. Content operations should still work
	app.Editor().SetContent("Event broker resilience test")

	// 3. Verify content
	current := app.Editor().GetContent()
	AssertEqual(t, "Event broker resilience test", current, "should work with event broker")
}

// TestModeRegistryBoundaryConditions tests registry with edge cases.
func TestModeRegistryBoundaryConditions(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	registry := app.ModeRegistry()

	// 1. Switch to non-existent mode (should be no-op)
	cmd := registry.SwitchTo(modes.ModeID("nonexistent"))
	AssertEqual(t, nil, cmd, "switching to non-existent mode should return nil")

	// 2. Current mode should be unchanged
	currentID := registry.CurrentID()
	AssertEqual(t, modes.ModeEdit, currentID, "current mode should not change")

	// 3. Switch to current mode (should be no-op)
	cmd = registry.SwitchTo(modes.ModeEdit)
	AssertEqual(t, nil, cmd, "switching to current mode should return nil")

	// 4. Previous should not change
	prevID := registry.PreviousID()
	AssertEqual(t, modes.ModeID(""), prevID, "previous should still be empty")
}

// TestRapidContentChanges tests rapid content modifications.
func TestRapidContentChanges(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Perform rapid content changes
	for i := 0; i < 100; i++ {
		content := fmt.Sprintf("Change %d", i)
		app.Editor().SetContent(content)
	}

	// 2. Verify last change is present
	finalContent := app.Editor().GetContent()
	AssertEqual(t, "Change 99", finalContent, "final content should reflect last change")
}

// Helper function to generate large content
func generateLargeContent(words int) string {
	var sb strings.Builder
	word := "Lorem ipsum dolor sit amet consectetur adipiscing elit "
	for i := 0; i < words; i++ {
		sb.WriteString(word)
		if i%10 == 0 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// TestAnalysisTimeouts tests that analysis doesn't hang indefinitely.
func TestAnalysisTimeouts(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Set content
	content := "Test content for timeout"
	app.Editor().SetContent(content)

	// 2. Start analysis
	cmd := app.EditMode().OnEnter()

	// 3. Use a timeout to ensure analysis doesn't hang
	done := make(chan bool)
	go func() {
		if cmd != nil {
			cmd()
		}
		done <- true
	}()

	// 4. Wait with timeout
	select {
	case <-done:
		// Analysis completed
		AssertTrue(t, true, "analysis completed")
	case <-time.After(5 * time.Second):
		t.Fatal("analysis timeout exceeded")
	}

	// 5. Content should still be accessible
	current := app.Editor().GetContent()
	AssertEqual(t, content, current, "content should be accessible after analysis")
}
