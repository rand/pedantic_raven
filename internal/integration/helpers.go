// Package integration provides end-to-end integration tests for Pedantic Raven.
package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/app/events"
	"github.com/rand/pedantic-raven/internal/editor"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"github.com/rand/pedantic-raven/internal/modes"
	"github.com/rand/pedantic-raven/internal/orchestrate"
	"github.com/rand/pedantic-raven/internal/overlay"
	"github.com/rand/pedantic-raven/internal/palette"
)

// TestApp represents a test instance of the Pedantic Raven application.
type TestApp struct {
	t               *testing.T
	eventBroker     *events.Broker
	modeRegistry    *modes.Registry
	overlayManager  *overlay.Manager
	paletteRegistry *palette.CommandRegistry
	tempDir         string
	analyzer        semantic.Analyzer
}

// NewTestApp creates a new test application instance.
func NewTestApp(t *testing.T) *TestApp {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "pedantic-raven-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create semantic analyzer with pattern-based extractor for predictable testing
	analyzer := semantic.NewAnalyzer()

	// Create event broker
	broker := events.NewBroker(100)

	// Create mode registry
	modeRegistry := modes.NewRegistry()

	// Register modes
	editMode := editor.NewEditModeWithAnalyzer(analyzer)
	exploreMode := modes.NewExploreMode()
	analyzeMode := modes.NewBaseMode(modes.ModeAnalyze, "Analyze", "Semantic analysis")
	orchestrateMode := orchestrate.NewModeAdapter()

	modeRegistry.Register(editMode)
	modeRegistry.Register(exploreMode)
	modeRegistry.Register(analyzeMode)
	modeRegistry.Register(orchestrateMode)

	// Set Edit as default mode
	modeRegistry.SwitchTo(modes.ModeEdit)

	// Create overlay manager
	overlayManager := overlay.NewManager()

	// Create command palette registry
	paletteRegistry := palette.NewCommandRegistry()

	return &TestApp{
		t:               t,
		eventBroker:     broker,
		modeRegistry:    modeRegistry,
		overlayManager:  overlayManager,
		paletteRegistry: paletteRegistry,
		tempDir:         tempDir,
		analyzer:        analyzer,
	}
}

// Cleanup cleans up test resources.
func (ta *TestApp) Cleanup() {
	if err := os.RemoveAll(ta.tempDir); err != nil {
		ta.t.Logf("Warning: Failed to clean up temp directory: %v", err)
	}
}

// EventBroker returns the event broker.
func (ta *TestApp) EventBroker() *events.Broker {
	return ta.eventBroker
}

// ModeRegistry returns the mode registry.
func (ta *TestApp) ModeRegistry() *modes.Registry {
	return ta.modeRegistry
}

// OverlayManager returns the overlay manager.
func (ta *TestApp) OverlayManager() *overlay.Manager {
	return ta.overlayManager
}

// SwitchToMode switches to the specified mode.
func (ta *TestApp) SwitchToMode(modeID modes.ModeID) tea.Cmd {
	return ta.modeRegistry.SwitchTo(modeID)
}

// CurrentMode returns the current mode.
func (ta *TestApp) CurrentMode() modes.Mode {
	return ta.modeRegistry.Current()
}

// CurrentModeID returns the current mode ID.
func (ta *TestApp) CurrentModeID() modes.ModeID {
	return ta.modeRegistry.CurrentID()
}

// Editor returns the Edit mode editor component.
func (ta *TestApp) Editor() *editor.EditorComponent {
	editMode, ok := ta.modeRegistry.Get(modes.ModeEdit).(*editor.EditMode)
	if !ok {
		ta.t.Fatal("Failed to cast to EditMode")
	}
	return editMode.GetEditor()
}

// EditMode returns the Edit mode.
func (ta *TestApp) EditMode() *editor.EditMode {
	editMode, ok := ta.modeRegistry.Get(modes.ModeEdit).(*editor.EditMode)
	if !ok {
		ta.t.Fatal("Failed to cast to EditMode")
	}
	return editMode
}

// ExploreMode returns the Explore mode.
func (ta *TestApp) ExploreMode() *modes.ExploreMode {
	exploreMode, ok := ta.modeRegistry.Get(modes.ModeExplore).(*modes.ExploreMode)
	if !ok {
		ta.t.Fatal("Failed to cast to ExploreMode")
	}
	return exploreMode
}

// OrchestrateMode returns the Orchestrate mode.
func (ta *TestApp) OrchestrateMode() modes.Mode {
	mode := ta.modeRegistry.Get(modes.ModeOrchestrate)
	if mode == nil {
		ta.t.Fatal("Orchestrate mode not found in registry")
	}
	return mode
}

// TempDir returns the temporary directory for test files.
func (ta *TestApp) TempDir() string {
	return ta.tempDir
}

// CreateTestFile creates a test file with the given content.
func (ta *TestApp) CreateTestFile(name, content string) (string, error) {
	path := filepath.Join(ta.tempDir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	return path, nil
}

// WaitForAnalysis waits for semantic analysis to complete.
func (ta *TestApp) WaitForAnalysis(timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("analysis timeout after %v", timeout)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// WaitForCondition waits for a condition to be true.
func (ta *TestApp) WaitForCondition(condition func() bool, timeout time.Duration) error {
	start := time.Now()
	for {
		if condition() {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("condition timeout after %v", timeout)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// MockAnalyzer is a test analyzer that returns predictable results.
type MockAnalyzer struct {
	analysis *semantic.Analysis
	stopped  bool
	running  bool
}

// NewMockAnalyzer creates a new mock analyzer.
func NewMockAnalyzer() *MockAnalyzer {
	return &MockAnalyzer{
		analysis: &semantic.Analysis{
			Content:       "",
			Entities:      []semantic.Entity{},
			Relationships: []semantic.Relationship{},
			TypedHoles:    []semantic.TypedHole{},
			Dependencies:  []semantic.Dependency{},
			Triples:       []semantic.Triple{},
			Duration:      0,
		},
		stopped: false,
		running: false,
	}
}

// Analyze returns a channel that immediately closes.
func (ma *MockAnalyzer) Analyze(content string) <-chan semantic.AnalysisUpdate {
	ma.running = true
	ma.analysis.Content = content
	ch := make(chan semantic.AnalysisUpdate)
	close(ch)
	ma.running = false
	return ch
}

// Results returns the mock analysis results.
func (ma *MockAnalyzer) Results() *semantic.Analysis {
	return ma.analysis
}

// Stop marks the analyzer as stopped.
func (ma *MockAnalyzer) Stop() {
	ma.stopped = true
	ma.running = false
}

// IsRunning returns whether analysis is running.
func (ma *MockAnalyzer) IsRunning() bool {
	return ma.running
}

// SetAnalysis sets the mock analysis results.
func (ma *MockAnalyzer) SetAnalysis(analysis *semantic.Analysis) {
	ma.analysis = analysis
}

// IsStopped returns whether the analyzer has been stopped.
func (ma *MockAnalyzer) IsStopped() bool {
	return ma.stopped
}

// ContextWithTimeout creates a context with a timeout.
func ContextWithTimeout(t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatalf("%s: expected error, got nil", msg)
	}
}

// AssertEqual fails the test if actual != expected.
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	if expected != actual {
		t.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNotEqual fails the test if actual == expected.
func AssertNotEqual(t *testing.T, expected, actual interface{}, msg string) {
	if expected == actual {
		t.Fatalf("%s: expected not %v, got %v", msg, expected, actual)
	}
}

// AssertTrue fails the test if condition is false.
func AssertTrue(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Fatalf("%s", msg)
	}
}

// AssertFalse fails the test if condition is true.
func AssertFalse(t *testing.T, condition bool, msg string) {
	if condition {
		t.Fatalf("%s", msg)
	}
}

// TestFixture represents a reusable test fixture.
type TestFixture struct {
	Name    string
	Content string
	Expected interface{}
}

// LoadTestFixture loads a test fixture from testdata directory.
func LoadTestFixture(t *testing.T, filename string) (string, error) {
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	return string(data), err
}
