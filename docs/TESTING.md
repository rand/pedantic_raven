# Integration Testing Guide for Pedantic Raven

This document provides comprehensive guidance on the integration test framework for Pedantic Raven, including test architecture, execution, and best practices.

## Table of Contents

1. [Overview](#overview)
2. [Test Architecture](#test-architecture)
3. [Running Tests](#running-tests)
4. [Test Categories](#test-categories)
5. [Test Fixtures](#test-fixtures)
6. [Writing New Tests](#writing-new-tests)
7. [Debugging Tests](#debugging-tests)
8. [CI/CD Integration](#cicd-integration)
9. [Performance Benchmarks](#performance-benchmarks)
10. [Best Practices](#best-practices)

## Overview

The Pedantic Raven integration test framework validates end-to-end functionality across the application's five modes:
- **Edit**: Context editing with semantic analysis
- **Explore**: Memory workspace with graph visualization
- **Analyze**: Semantic insights and triple analysis
- **Orchestrate**: Multi-agent coordination and task management
- **Collaborate**: Live multi-user editing (future)

The framework includes 25+ integration tests covering:
- Cross-mode workflows
- Session persistence
- Error recovery
- Large dataset handling
- Concurrent operations
- Race condition prevention

## Test Architecture

### Test App Instance

The `TestApp` provides a complete application instance for testing:

```go
app := NewTestApp(t)
defer app.Cleanup()

// Access components
app.Editor()           // Get editor component
app.EditMode()         // Get Edit mode
app.ExploreMode()      // Get Explore mode
app.OrchestrateMode()  // Get Orchestrate mode
app.ModeRegistry()     // Get mode registry
app.EventBroker()      // Get event broker
```

### Key Components

- **EventBroker**: Pub/sub event system for component communication
- **ModeRegistry**: Manages mode registration and switching
- **EditMode**: Text editing with semantic analysis
- **Semantic Analyzer**: Extracts entities, relationships, holes, triples
- **Mode Lifecycle**: OnEnter/OnExit callbacks for state management

### Test Isolation

Each test:
1. Creates a fresh TestApp instance
2. Uses isolated temporary directories (via `t.TempDir()`)
3. Cleans up resources with `defer app.Cleanup()`
4. Runs independently without shared state

## Running Tests

### Run All Integration Tests

```bash
# Run all tests in integration package
go test ./internal/integration/... -v

# With race detection
go test ./internal/integration/... -v -race

# With coverage
go test ./internal/integration/... -v -cover

# With detailed output
go test ./internal/integration/... -v -run TestEditAnalyzeWorkflow
```

### Run Specific Test Categories

```bash
# Workflow tests
go test ./internal/integration/... -v -run TestEdit

# Persistence tests
go test ./internal/integration/... -v -run TestSession

# Error recovery tests
go test ./internal/integration/... -v -run TestOffline

# Large dataset tests
go test ./internal/integration/... -v -run TestLarge

# Concurrent tests
go test ./internal/integration/... -v -run Concurrent
```

### Run Single Test

```bash
go test ./internal/integration/... -v -run TestEditAnalyzeWorkflow
```

### Benchmark Tests

```bash
# Run benchmarks
go test ./internal/integration/... -bench=. -benchmem

# Run specific benchmark
go test ./internal/integration/... -bench=BenchmarkLargeFileProcessing -benchmem
```

### Race Condition Detection

The `-race` flag enables Go's race detector:

```bash
# Detect race conditions
go test ./internal/integration/... -v -race

# Race detector requires specific environment
CGO_ENABLED=1 go test ./internal/integration/... -v -race
```

## Test Categories

### 1. Workflow Tests (`workflow_test.go`)

Tests cross-mode workflows and state transitions:

- **TestEditAnalyzeWorkflow**: Edit content → Analyze semantic insights
- **TestAnalyzeOrchestrateWorkflow**: Analyze → Create work plan in Orchestrate
- **TestEditOrchestrateWorkflow**: Write plan in Editor → Launch orchestration
- **TestMultiModeNavigation**: Rapid mode switching without state corruption
- **TestModeStatePreservation**: State preserved when switching away/back
- **TestModeRegistryPreviousMode**: "Go back" functionality (previous mode)
- **TestModeTransitionCommands**: OnExit/OnEnter lifecycle hooks execute
- **TestEditModeOnEnter**: OnEnter triggers analysis
- **TestOrchestrateModeSwitching**: Switch to/from Orchestrate mode

**Coverage**: Mode lifecycle, state transitions, content preservation

### 2. Persistence Tests (`persistence_test.go`)

Tests session state preservation across application restarts:

- **TestSessionStatePreservation**: Save and restore session state
- **TestEditBufferPersistence**: Edit buffer state persists across restarts
- **TestModeStatePersistence**: Mode state persists across mode switches
- **TestAnalysisResultsPreservation**: Analysis results preserved
- **TestMultipleFileHandling**: Handle multiple files in session
- **TestContentModificationTracking**: Track content modifications
- **TestSessionRecovery**: Recovery from incomplete state

**Coverage**: File persistence, buffer state, mode state, content tracking

### 3. Error Recovery Tests (`error_recovery_test.go`)

Tests graceful degradation and error handling:

- **TestOfflineModeRecovery**: Handle when external services unavailable
- **TestCorruptContentRecovery**: Recover from malformed content
- **TestInvalidSemanticAnalysisInput**: Handle invalid analysis input (empty, very long)
- **TestModeTransitionUnderLoad**: Mode transitions with significant content
- **TestAnalysisErrorRecovery**: Recovery from analysis errors
- **TestGracefulDegradation**: App degrades gracefully when components unavailable
- **TestEventBrokerResilience**: Event broker resilience
- **TestModeRegistryBoundaryConditions**: Registry edge cases
- **TestRapidContentChanges**: Rapid content modifications (100+ changes)
- **TestAnalysisTimeouts**: Analysis doesn't hang indefinitely

**Coverage**: Error handling, graceful degradation, resilience, edge cases

### 4. Large Dataset Tests (`large_dataset_test.go`)

Tests performance with large files and datasets:

- **TestLargeFileEditing**: Edit 10,000+ line files
- **TestLargeEntityAnalysis**: Analyze content with 1000+ entities
- **TestDeeplyNestedStructure**: Handle deeply nested content
- **TestLargeWorkPlan**: Handle work plans with 100+ tasks
- **TestManyModeSwitches**: Mode switching with large content (5000 lines)
- **TestComplexAnalysisScenario**: Complex real-world analysis
- **TestMultiLanguageContent**: Handle content in multiple formats
- **TestPerformanceRegression**: Operations don't degrade over time

**Coverage**: Performance, scalability, file size limits

### 5. Concurrent Tests (`concurrent_test.go`)

Tests concurrent operations and race condition prevention:

- **TestConcurrentModeSwitching**: Concurrent mode switches (10 goroutines)
- **TestConcurrentBufferOperations**: Concurrent buffer read/write (20 goroutines)
- **TestConcurrentEventPublishing**: Publish events concurrently (10 publishers)
- **TestConcurrentAnalysisCalls**: Concurrent analysis requests (10 analyzers)
- **TestRaceConditionModeSwitchAndAnalysis**: Race condition detection (3 goroutines)
- **TestConcurrentModeInitialization**: Concurrent mode initialization
- **TestDeadlockPrevention**: Deadlock detection with timeout (15 goroutines)
- **TestConcurrentFileOperations**: Concurrent file operations (10 files)
- **TestStressTestConcurrentOperations**: Stress test (50 goroutines × 100 ops)
- **TestConcurrentModeGetters**: Concurrent access to mode getters

**Coverage**: Race conditions, deadlocks, concurrent safety

## Test Fixtures

### Built-in Content Generators

The `testdata_generator.go` provides reusable content:

```go
// Sample markdown content
content := SampleMarkdownContent()

// Sample code content
code := SampleCodeContent()

// Sample work plan
plan := SampleWorkPlan()

// Entity-rich content (1000 entities)
richContent := GenerateEntityRichContent(1000)

// Large content (10000 lines)
large := GenerateLargeContent(10000)

// Nested structure (depth 10, 5 items per level)
nested := GenerateNestedContent(10, 5)

// Code with typed holes
withHoles := GenerateCodeWithTypedHoles()

// Special characters test
special := SampleSpecialCharacters()
```

### Creating Test Files

```go
// Create temporary test file
filePath, err := app.CreateTestFile("test.md", "content")
if err != nil {
    t.Fatalf("Failed to create test file: %v", err)
}

// Access temporary directory
tempDir := app.TempDir()
```

### Mock Analyzer

For predictable test results:

```go
mockAnalyzer := NewMockAnalyzer()
mockAnalyzer.SetAnalysis(&semantic.Analysis{
    Entities: []semantic.Entity{
        {Text: "Alice", Type: semantic.EntityPerson},
        {Text: "Acme", Type: semantic.EntityOrganization},
    },
})
editMode.SetAnalyzer(mockAnalyzer)
```

## Writing New Tests

### Test Template

```go
func TestNewFeature(t *testing.T) {
    // 1. Setup
    app := NewTestApp(t)
    defer app.Cleanup()

    // 2. Execute
    app.Editor().SetContent("test content")
    cmd := app.SwitchToMode(modes.ModeAnalyze)
    if cmd != nil {
        cmd()
    }

    // 3. Verify
    AssertEqual(t, modes.ModeAnalyze, app.CurrentModeID(), "should be in Analyze mode")
}
```

### Assertion Helpers

```go
// Equality assertions
AssertEqual(t, expected, actual, "message")
AssertNotEqual(t, notExpected, actual, "message")

// Boolean assertions
AssertTrue(t, condition, "message")
AssertFalse(t, condition, "message")

// Error assertions
AssertNoError(t, err, "message")
AssertError(t, err, "message")
```

### Condition Waiting

```go
// Wait for condition with timeout
err := app.WaitForCondition(func() bool {
    return app.Editor().GetContent() != ""
}, 2*time.Second)

if err != nil {
    t.Fatalf("Condition not met: %v", err)
}
```

### Testing Concurrent Code

```go
var wg sync.WaitGroup

for i := 0; i < numGoroutines; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // Test concurrent operation
    }(i)
}

wg.Wait()

// Verify results
AssertEqual(t, expectedState, actualState, "concurrent operations should complete safely")
```

### Testing Mode Transitions

```go
// Test state before transition
initialState := app.Editor().GetContent()

// Perform transition
cmd := app.SwitchToMode(modes.ModeAnalyze)
if cmd != nil {
    cmd()
}

// Test state after transition
AssertEqual(t, modes.ModeAnalyze, app.CurrentModeID(), "should be in Analyze mode")

// Return and verify preservation
cmd = app.SwitchToMode(modes.ModeEdit)
if cmd != nil {
    cmd()
}

AssertEqual(t, initialState, app.Editor().GetContent(), "state should be preserved")
```

## Debugging Tests

### Enable Verbose Output

```bash
go test ./internal/integration/... -v
```

### Run Single Test with Output

```bash
go test ./internal/integration/... -v -run TestEditAnalyzeWorkflow -count=1
```

### Add Debug Logging

```go
func TestDebugExample(t *testing.T) {
    app := NewTestApp(t)
    defer app.Cleanup()

    t.Logf("Debug: Starting test")
    t.Logf("Debug: Current mode: %v", app.CurrentModeID())

    app.Editor().SetContent("test")
    t.Logf("Debug: Content set: %s", app.Editor().GetContent())
}
```

### Check Race Conditions

```bash
# Run with race detector
CGO_ENABLED=1 go test ./internal/integration/... -v -race

# Look for "DATA RACE" in output
```

### Use Timeout for Hanging Tests

```bash
# Timeout after 10 seconds
timeout 10s go test ./internal/integration/... -v -run TestName
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.20

      # Run tests
      - name: Run integration tests
        run: go test ./internal/integration/... -v -cover

      # Run with race detection
      - name: Run race detector
        run: go test ./internal/integration/... -v -race

      # Upload coverage
      - name: Upload coverage
        uses: codecov/codecov-action@v2
```

### Local CI Simulation

```bash
#!/bin/bash

# Run all checks
echo "Running tests..."
go test ./internal/integration/... -v

echo "Running race detector..."
CGO_ENABLED=1 go test ./internal/integration/... -v -race

echo "Checking coverage..."
go test ./internal/integration/... -coverprofile=coverage.out
go tool cover -html=coverage.out

echo "All checks complete!"
```

## Performance Benchmarks

### Current Benchmarks

```
TestLargeFileEditing (10,000 lines): < 5s
TestLargeEntityAnalysis (1,000 entities): < 10s
TestManyModeSwitches (50 switches): < 2s
TestStressTest (5000 goroutine ops): < 10s
```

### Running Benchmarks

```bash
go test ./internal/integration/... -bench=. -benchmem -benchtime=3s
```

### Interpreting Results

```
BenchmarkOperation-8    1000    1234567 ns/op    8192 B/op    123 allocs/op
```

- `-8`: Running on 8 CPUs
- `1000`: Number of iterations
- `ns/op`: Nanoseconds per operation
- `B/op`: Bytes allocated per operation
- `allocs/op`: Allocations per operation

## Best Practices

### 1. Test Isolation

```go
// ✓ Good: Each test creates own app
func TestFeatureA(t *testing.T) {
    app := NewTestApp(t)
    defer app.Cleanup()
    // ...
}

// ✗ Bad: Shared global state
var globalApp *TestApp
func init() {
    globalApp = NewTestApp(nil)
}
```

### 2. Resource Cleanup

```go
// ✓ Good: Cleanup in defer
app := NewTestApp(t)
defer app.Cleanup()

// ✓ Good: Explicit resource cleanup
defer os.RemoveAll(app.TempDir())
```

### 3. Deterministic Tests

```go
// ✓ Good: No time-dependent operations
content := "fixed content"

// ✗ Bad: Time-dependent
content := fmt.Sprintf("test-%d", time.Now().Unix())
```

### 4. Meaningful Assertions

```go
// ✓ Good: Descriptive assertion messages
AssertEqual(t, modes.ModeEdit, app.CurrentModeID(),
    "should remain in Edit mode after content change")

// ✗ Bad: Vague assertion messages
AssertEqual(t, modes.ModeEdit, app.CurrentModeID(), "mode check")
```

### 5. Test Documentation

```go
// ✓ Good: Clear test purpose
// TestEditAnalyzeWorkflow tests the Edit -> Analyze cross-mode workflow
// Verifies that semantic analysis results are available when switching
// to Analyze mode after editing content.
func TestEditAnalyzeWorkflow(t *testing.T) {
    // ...
}
```

### 6. Concurrent Test Guidelines

```go
// ✓ Good: Use WaitGroup for goroutine coordination
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // ...
}()
wg.Wait()

// ✓ Good: Use channel for result collection
results := make(chan string)
go func() {
    results <- "value"
}()
value := <-results

// ✗ Bad: Uncoordinated goroutines
go func() {
    // May not complete before test ends
}()
```

### 7. Error Handling

```go
// ✓ Good: Check and handle errors
if err != nil {
    t.Fatalf("Failed to create test file: %v", err)
}

// ✓ Good: Use assertion helpers
AssertNoError(t, err, "should create test file")

// ✗ Bad: Ignore errors
_ = os.WriteFile(path, []byte(content), 0644)
```

### 8. Test Coverage Targets

- **Critical path**: 90%+
- **Business logic**: 80%+
- **UI layer**: 60%+
- **Overall**: 70%+

Run coverage analysis:

```bash
go test ./internal/integration/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Troubleshooting

### Tests Hang

```bash
# Run with timeout
timeout 30s go test ./internal/integration/... -v -run TestName

# Check for deadlocks in code
# Review goroutine coordination (WaitGroup, channels, etc.)
```

### Race Conditions Detected

```bash
# Run with race detector
CGO_ENABLED=1 go test ./internal/integration/... -v -race

# Review concurrent access to shared variables
# Use synchronization primitives (mutex, channels)
```

### Flaky Tests

- Remove time-dependent logic
- Ensure proper synchronization
- Use `WaitForCondition` for asynchronous operations
- Avoid sleeps (use channels/synchronization instead)

### Memory Leaks

```bash
# Use pprof for memory analysis
import _ "net/http/pprof"

# Run tests and analyze memory profile
go test -memprofile=mem.prof ./internal/integration/...
go tool pprof mem.prof
```

## Summary

The integration test framework provides:

- **25+ comprehensive tests** covering all major workflows
- **Multiple test categories** (workflow, persistence, error, performance, concurrent)
- **Reusable test utilities** (TestApp, assertions, content generators)
- **Best practices** for test writing and maintenance
- **CI/CD integration** examples
- **Performance benchmarks** for regression detection

Run the tests regularly during development and before releases to ensure application stability and reliability.
