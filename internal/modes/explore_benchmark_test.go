package modes

import (
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/testhelpers"
)

// setupBenchmark disables mnemosyne for benchmarking
func setupBenchmark(b *testing.B) func() {
	oldEnabled := os.Getenv("MNEMOSYNE_ENABLED")
	os.Setenv("MNEMOSYNE_ENABLED", "false")

	return func() {
		if oldEnabled != "" {
			os.Setenv("MNEMOSYNE_ENABLED", oldEnabled)
		} else {
			os.Unsetenv("MNEMOSYNE_ENABLED")
		}
	}
}

// BenchmarkExploreModeInitialization benchmarks mode initialization
func BenchmarkExploreModeInitialization(b *testing.B) {
	defer setupBenchmark(b)()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mode := NewExploreMode()
		mode.Init()
	}
}

// BenchmarkExploreModeOnEnter benchmarks sample data loading
func BenchmarkExploreModeOnEnter(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd := mode.OnEnter()
		if cmd != nil {
			cmd()
		}
	}
}

// BenchmarkExploreModeLayoutToggle benchmarks layout switching
func BenchmarkExploreModeLayoutToggle(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	// Set window size
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mode.toggleLayout()
	}
}

// BenchmarkExploreModeView benchmarks view rendering
func BenchmarkExploreModeView(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	// Load sample data
	cmd := mode.OnEnter()
	if cmd != nil {
		msg := cmd()
		mode.Update(msg)
	}

	// Set window size
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mode.View()
	}
}

// BenchmarkExploreModeViewGraphLayout benchmarks graph view rendering
func BenchmarkExploreModeViewGraphLayout(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	// Load sample data
	cmd := mode.OnEnter()
	if cmd != nil {
		msg := cmd()
		mode.Update(msg)
	}

	// Set window size and switch to graph
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)
	mode.toggleLayout()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mode.View()
	}
}

// BenchmarkExploreModeUpdate benchmarks update processing
func BenchmarkExploreModeUpdate(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		updatedMode, _ := mode.Update(keyMsg)
		mode = updatedMode.(*ExploreMode)
	}
}

// BenchmarkExploreModeFocusCycle benchmarks focus cycling
func BenchmarkExploreModeFocusCycle(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mode.cycleFocus()
	}
}

// BenchmarkExploreModeKeybindings benchmarks keybindings generation
func BenchmarkExploreModeKeybindings(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mode.Keybindings()
	}
}

// BenchmarkExploreModeHelpView benchmarks help overlay rendering
func BenchmarkExploreModeHelpView(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	mode.showHelp = true
	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mode.View()
	}
}

// BenchmarkExploreModeWindowResize benchmarks window resize handling
func BenchmarkExploreModeWindowResize(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mode.Update(wsMsg)
	}
}

// BenchmarkExploreModeSampleDataLarge benchmarks loading larger dataset
func BenchmarkExploreModeSampleDataLarge(b *testing.B) {
	defer setupBenchmark(b)()

	// Generate large dataset
	memories := testhelpers.GenerateTestMemories(1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mode := NewExploreMode()
		mode.Init()

		// Simulate loading memories (without actual server call)
		_ = memories
	}
}

// BenchmarkExploreModeGraphGeneration benchmarks graph generation
func BenchmarkExploreModeGraphGeneration(b *testing.B) {
	defer setupBenchmark(b)()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = testhelpers.GenerateTestGraph(100, 200)
	}
}

// BenchmarkExploreModeCompleteWorkflow benchmarks full workflow
func BenchmarkExploreModeCompleteWorkflow(b *testing.B) {
	defer setupBenchmark(b)()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create and initialize
		mode := NewExploreMode()
		mode.Init()

		// Load data
		cmd := mode.OnEnter()
		if cmd != nil {
			msg := cmd()
			mode.Update(msg)
		}

		// Set size
		wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
		mode.Update(wsMsg)

		// Toggle layout
		mode.toggleLayout()

		// Render view
		_ = mode.View()

		// Toggle back
		mode.toggleLayout()

		// Final view
		_ = mode.View()
	}
}

// BenchmarkExploreModeRapidUpdates benchmarks rapid update sequences
func BenchmarkExploreModeRapidUpdates(b *testing.B) {
	defer setupBenchmark(b)()

	mode := NewExploreMode()
	mode.Init()

	wsMsg := tea.WindowSizeMsg{Width: 100, Height: 50}
	mode.Update(wsMsg)

	keys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'g'}},
		{Type: tea.KeyRunes, Runes: []rune{'?'}},
		{Type: tea.KeyTab},
		{Type: tea.KeyEscape},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, keyMsg := range keys {
			updatedMode, _ := mode.Update(keyMsg)
			mode = updatedMode.(*ExploreMode)
		}
	}
}
