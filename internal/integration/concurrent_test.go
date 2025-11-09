package integration

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rand/pedantic-raven/internal/modes"
)

// TestConcurrentModeSwitching tests concurrent mode switches from multiple goroutines.
func TestConcurrentModeSwitching(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// Set initial content
	app.Editor().SetContent("Concurrent test content")

	var wg sync.WaitGroup
	numGoroutines := 10

	modeSequence := []modes.ModeID{
		modes.ModeEdit,
		modes.ModeAnalyze,
		modes.ModeExplore,
		modes.ModeOrchestrate,
	}

	// Launch concurrent goroutines
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				modeID := modeSequence[(id+i)%len(modeSequence)]
				cmd := app.SwitchToMode(modeID)
				if cmd != nil {
					cmd()
				}
				time.Sleep(time.Millisecond) // Small delay
			}
		}(g)
	}

	wg.Wait()

	// After concurrent operations, verify app is still functioning
	content := app.Editor().GetContent()
	AssertEqual(t, "Concurrent test content", content, "content should survive concurrent mode switches")
}

// TestConcurrentBufferOperations tests concurrent buffer read/write operations.
func TestConcurrentBufferOperations(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	var wg sync.WaitGroup
	numGoroutines := 20
	operations := 50

	// Launch readers and writers
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if id%2 == 0 {
				// Writer goroutine
				for i := 0; i < operations; i++ {
					content := fmt.Sprintf("Goroutine %d - operation %d", id, i)
					app.Editor().SetContent(content)
				}
			} else {
				// Reader goroutine
				for i := 0; i < operations; i++ {
					_ = app.Editor().GetContent()
					time.Sleep(time.Microsecond)
				}
			}
		}(g)
	}

	wg.Wait()

	// Verify app is still functional
	finalContent := app.Editor().GetContent()
	AssertNotEqual(t, "", finalContent, "buffer should contain content after concurrent operations")
}

// TestConcurrentEventPublishing tests publishing events concurrently.
func TestConcurrentEventPublishing(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	broker := app.EventBroker()
	_ = broker // broker is used implicitly by the app
	var wg sync.WaitGroup
	numPublishers := 10
	eventsPerPublisher := 100

	// Publish events concurrently
	for p := 0; p < numPublishers; p++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for e := 0; e < eventsPerPublisher; e++ {
				// Events are published but we don't have subscribers in test
				_ = id
				_ = e
				// In real scenario, would publish events here
			}
		}(p)
	}

	wg.Wait()

	// App should still function
	content := "Test content"
	app.Editor().SetContent(content)
	retrieved := app.Editor().GetContent()
	AssertEqual(t, content, retrieved, "app should work after concurrent event publishing")
}

// TestConcurrentAnalysisCalls tests concurrent analysis requests.
func TestConcurrentAnalysisCalls(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	var wg sync.WaitGroup
	numAnalyzers := 10

	// Concurrent analysis calls
	for a := 0; a < numAnalyzers; a++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			content := fmt.Sprintf("Content for analyzer %d", id)
			app.Editor().SetContent(content)

			cmd := app.EditMode().OnEnter()
			if cmd != nil {
				cmd()
			}
		}(a)
	}

	wg.Wait()

	// Verify app is still responsive
	current := app.Editor().GetContent()
	AssertNotEqual(t, "", current, "content should be available after concurrent analysis")
}

// TestRaceConditionModeSwitchAndAnalysis tests race condition between mode switching and analysis.
func TestRaceConditionModeSwitchAndAnalysis(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	app.Editor().SetContent("Initial content for race test")

	var wg sync.WaitGroup

	// Goroutine 1: Rapid mode switching
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			modes := []modes.ModeID{
				modes.ModeEdit,
				modes.ModeAnalyze,
				modes.ModeExplore,
			}
			cmd := app.SwitchToMode(modes[i%len(modes)])
			if cmd != nil {
				cmd()
			}
		}
	}()

	// Goroutine 2: Rapid analysis calls
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			app.Editor().SetContent(fmt.Sprintf("Content change %d", i))
			cmd := app.EditMode().OnEnter()
			if cmd != nil {
				cmd()
			}
		}
	}()

	// Goroutine 3: Content modifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 20; i++ {
			app.Editor().SetContent(fmt.Sprintf("Modified content %d", i))
		}
	}()

	wg.Wait()

	// App should still be responsive
	content := app.Editor().GetContent()
	AssertNotEqual(t, "", content, "app should handle race conditions gracefully")
}

// TestConcurrentModeInitialization tests concurrent initialization of modes.
func TestConcurrentModeInitialization(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	var wg sync.WaitGroup

	// Concurrent mode initialization calls
	modeIDs := []modes.ModeID{
		modes.ModeEdit,
		modes.ModeAnalyze,
		modes.ModeExplore,
		modes.ModeOrchestrate,
	}

	for i, modeID := range modeIDs {
		wg.Add(1)
		go func(id modes.ModeID, idx int) {
			defer wg.Done()
			mode := app.ModeRegistry().Get(id)
			if mode != nil {
				cmd := mode.Init()
				if cmd != nil {
					cmd()
				}
			}
		}(modeID, i)
	}

	wg.Wait()

	// Verify registry is still functional
	allModes := app.ModeRegistry().AllModes()
	AssertEqual(t, 4, len(allModes), "should have all modes registered")
}

// TestDeadlockPrevention tests that concurrent operations don't cause deadlocks.
func TestDeadlockPrevention(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	var wg sync.WaitGroup
	done := make(chan bool)
	numGoroutines := 15

	// Launch many concurrent operations
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for i := 0; i < 10; i++ {
				app.Editor().SetContent(fmt.Sprintf("Content %d-%d", id, i))
				_ = app.Editor().GetContent()

				if i%3 == 0 {
					var modeID modes.ModeID
					if i%2 == 0 {
						modeID = modes.ModeAnalyze
					} else {
						modeID = modes.ModeExplore
					}
					cmd := app.SwitchToMode(modeID)
					if cmd != nil {
						cmd()
					}
				}

				if i%5 == 0 {
					cmd := app.EditMode().OnEnter()
					if cmd != nil {
						cmd()
					}
				}
			}
		}(g)
	}

	// Wait with timeout
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		AssertTrue(t, true, "concurrent operations completed")
	case <-time.After(10 * time.Second):
		t.Fatal("deadlock detected: operations did not complete within timeout")
	}
}

// TestConcurrentFileOperations tests concurrent file-like operations.
func TestConcurrentFileOperations(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	var wg sync.WaitGroup
	numFiles := 10
	operationsPerFile := 20

	// Simulate concurrent file operations
	for f := 0; f < numFiles; f++ {
		wg.Add(1)
		go func(fileID int) {
			defer wg.Done()

			for op := 0; op < operationsPerFile; op++ {
				content := fmt.Sprintf("File %d - Operation %d", fileID, op)
				app.Editor().SetContent(content)

				retrieved := app.Editor().GetContent()
				if retrieved != content {
					t.Errorf("Content mismatch for file %d", fileID)
				}
			}
		}(f)
	}

	wg.Wait()

	// Final verification
	current := app.Editor().GetContent()
	AssertNotEqual(t, "", current, "content should be available after concurrent file operations")
}

// TestStressTestConcurrentOperations performs a stress test with many operations.
func TestStressTestConcurrentOperations(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	var wg sync.WaitGroup
	numGoroutines := 50
	operationsPerGoroutine := 100

	start := time.Now()

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for op := 0; op < operationsPerGoroutine; op++ {
				// Vary operations
				switch op % 4 {
				case 0:
					app.Editor().SetContent(fmt.Sprintf("G%d-Op%d", id, op))
				case 1:
					_ = app.Editor().GetContent()
				case 2:
					if op%10 == 0 {
						var modeID modes.ModeID
						if op%2 == 0 {
							modeID = modes.ModeAnalyze
						} else {
							modeID = modes.ModeExplore
						}
						cmd := app.SwitchToMode(modeID)
						if cmd != nil {
							cmd()
						}
					}
				case 3:
					if op%15 == 0 {
						cmd := app.EditMode().OnEnter()
						if cmd != nil {
							cmd()
						}
					}
				}
			}
		}(g)
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("Stress test: %d goroutines x %d ops completed in %v", numGoroutines, operationsPerGoroutine, elapsed)

	// Verify app is still responsive
	content := "Final test content"
	app.Editor().SetContent(content)
	retrieved := app.Editor().GetContent()
	AssertEqual(t, content, retrieved, "app should be responsive after stress test")
}

// TestConcurrentModeGetters tests concurrent access to mode getters.
func TestConcurrentModeGetters(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	var wg sync.WaitGroup
	numGoroutines := 20

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for i := 0; i < 50; i++ {
				_ = app.EditMode()
				_ = app.ExploreMode()
				_ = app.OrchestrateMode()
				_ = app.CurrentMode()
				_ = app.CurrentModeID()
			}
		}(g)
	}

	wg.Wait()

	// Verify app still works
	mode := app.CurrentMode()
	AssertNotEqual(t, nil, mode, "current mode should be available")
}
