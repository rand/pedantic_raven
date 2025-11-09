package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rand/pedantic-raven/internal/modes"
)

// TestLargeFileEditing tests editing very large files (10000+ lines).
func TestLargeFileEditing(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Generate large file content
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString(fmt.Sprintf("Line %d: Some content here\n", i))
	}
	largeContent := sb.String()

	// 2. Load large file
	start := time.Now()
	app.Editor().SetContent(largeContent)
	duration := time.Since(start)

	t.Logf("Loading 10000 lines took %v", duration)

	// 3. Verify content is loaded
	current := app.Editor().GetContent()
	AssertEqual(t, largeContent, current, "large file should be loaded completely")

	// 4. Verify operations are responsive
	AssertTrue(t, duration < 5*time.Second, "loading 10000 lines should complete in reasonable time")
}

// TestLargeEntityAnalysis tests semantic analysis on content with 1000+ entities.
func TestLargeEntityAnalysis(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Generate content with many entities
	var sb strings.Builder
	sb.WriteString("# Entity List\n\n")

	// Create many entity-like lines
	entities := []string{
		"Alice", "Bob", "Charlie", "David", "Eve",
		"Frank", "Grace", "Henry", "Iris", "Jack",
	}

	for i := 0; i < 100; i++ {
		for _, entity := range entities {
			sb.WriteString(fmt.Sprintf("%s works at Company%d in Location%d\n", entity, i%10, i%5))
		}
	}

	largeContent := sb.String()

	// 2. Analyze large content
	start := time.Now()
	app.Editor().SetContent(largeContent)
	duration := time.Since(start)

	// 3. Trigger analysis
	cmd := app.EditMode().OnEnter()
	if cmd != nil {
		cmd()
	}

	t.Logf("Large entity analysis took %v", duration)

	// 4. Verify content integrity
	current := app.Editor().GetContent()
	AssertEqual(t, largeContent, current, "large content with many entities should be preserved")
}

// TestDeeplyNestedStructure tests handling deeply nested content.
func TestDeeplyNestedStructure(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Generate deeply nested structure
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		for j := 0; j < i%10; j++ {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("- Item at depth %d\n", i))
	}

	nestedContent := sb.String()

	// 2. Load nested structure
	app.Editor().SetContent(nestedContent)

	// 3. Verify content integrity
	current := app.Editor().GetContent()
	AssertEqual(t, nestedContent, current, "deeply nested content should be preserved")
}

// TestLargeWorkPlan tests handling work plans with 100+ tasks.
func TestLargeWorkPlan(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Generate large work plan
	var sb strings.Builder
	sb.WriteString("# Large Work Plan\n\n")

	for i := 0; i < 100; i++ {
		sb.WriteString(fmt.Sprintf("## Sprint %d\n", i))
		for j := 0; j < 5; j++ {
			sb.WriteString(fmt.Sprintf("- [ ] Task %d.%d: Implementation work\n", i, j))
			sb.WriteString(fmt.Sprintf("  - Subtask: Details\n"))
		}
		sb.WriteString("\n")
	}

	largePlan := sb.String()

	// 2. Load work plan
	start := time.Now()
	app.Editor().SetContent(largePlan)
	loadDuration := time.Since(start)

	// 3. Verify loading performance
	AssertTrue(t, loadDuration < 5*time.Second, "loading large work plan should be fast")

	// 4. Verify content integrity
	current := app.Editor().GetContent()
	AssertEqual(t, largePlan, current, "large work plan should be complete")
}

// TestManyModeSwitches tests mode switching with large content.
func TestManyModeSwitches(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Load large content
	var sb strings.Builder
	for i := 0; i < 5000; i++ {
		sb.WriteString(fmt.Sprintf("Content line %d\n", i))
	}
	largeContent := sb.String()
	app.Editor().SetContent(largeContent)

	// 2. Perform many mode switches
	start := time.Now()
	for i := 0; i < 50; i++ {
		var modeID modes.ModeID
		if (i % 2) == 0 {
			modeID = modes.ModeAnalyze
		} else {
			modeID = modes.ModeExplore
		}
		cmd := app.SwitchToMode(modeID)
		if cmd != nil {
			cmd()
		}
	}
	switchDuration := time.Since(start)

	// 3. Verify performance
	t.Logf("50 mode switches with 5000 lines took %v", switchDuration)

	// 4. Verify content integrity
	final := app.Editor().GetContent()
	AssertEqual(t, largeContent, final, "content should survive many mode switches")
}

// TestComplexAnalysisScenario tests a complex real-world analysis scenario.
func TestComplexAnalysisScenario(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Create complex content mixing multiple patterns
	content := `# Project: Enterprise System Design

## Participants
- Alice Chen (Product Manager) at ACME Corp
- Bob Rodriguez (Lead Architect) at ACME Corp
- Charlie Zhang (Senior Developer) at TechFlow Inc
- Diana Park (QA Lead) at TechFlow Inc

## Requirements
- System must process 10,000 transactions per second
- 99.99% uptime requirement (?? Database: PostgreSQL or DynamoDB)
- Multi-region deployment (US, EU, APAC)
- ?? Authentication: OAuth2 or SAML integration

## Architecture Components

### API Layer
- REST endpoints with GraphQL fallback
- Rate limiting: 1000 req/min per user
- ?? Caching: Redis or Memcached

### Processing Layer
- Message queue (RabbitMQ/Kafka)
- Worker pool (50-200 workers)
- ?? Circuit breaker pattern implementation

### Data Layer
- Master-slave replication
- Sharding key: user_id
- Backup: S3 + daily snapshots

## Timeline
- Week 1-2: Design review
- Week 3-4: Core implementation
- Week 5-6: Integration testing
- Week 7: Deployment preparation

## Risks
- Integration complexity with legacy systems
- ?? Risk mitigation for scale testing
- Training overhead for new team members`

	// 2. Load complex content
	start := time.Now()
	app.Editor().SetContent(content)
	loadTime := time.Since(start)

	// 3. Trigger analysis
	cmd := app.EditMode().OnEnter()
	if cmd != nil {
		cmd()
	}

	// 4. Verify content
	current := app.Editor().GetContent()
	AssertEqual(t, content, current, "complex content should be preserved")

	t.Logf("Complex analysis scenario completed in %v", loadTime)
}

// TestMultiLanguageContent tests handling content in multiple formats.
func TestMultiLanguageContent(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	// 1. Create mixed-language content
	content := `# Multilingual Documentation

## English Section
Alice works with Bob on Project Alpha.

## Code Section
` + "```go" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

## Configuration Section
name: "ProjectAlpha"
version: "1.0.0"
author: "Team"

## Mathematical Notation
Area = π * r²
Distance = √(x₁-x₂)² + (y₁-y₂)²`

	// 2. Load mixed content
	app.Editor().SetContent(content)

	// 3. Verify integrity
	current := app.Editor().GetContent()
	AssertEqual(t, content, current, "multilingual content should be preserved")

	// 4. Test mode switching with mixed content
	cmd := app.SwitchToMode("analyze")
	if cmd != nil {
		cmd()
	}

	cmd = app.SwitchToMode("edit")
	if cmd != nil {
		cmd()
	}

	// 5. Verify content after switches
	final := app.Editor().GetContent()
	AssertEqual(t, content, final, "mixed-language content should survive mode switches")
}

// TestPerformanceRegression tests that operations don't degrade over time.
func TestPerformanceRegression(t *testing.T) {
	app := NewTestApp(t)
	defer app.Cleanup()

	baseContent := "Initial content"
	iterations := 100

	// 1. Perform repeated operations and measure time
	times := make([]time.Duration, iterations)

	for i := 0; i < iterations; i++ {
		content := fmt.Sprintf("%s - iteration %d", baseContent, i)

		start := time.Now()
		app.Editor().SetContent(content)
		_ = app.Editor().GetContent()
		times[i] = time.Since(start)
	}

	// 2. Calculate average time
	var totalTime time.Duration
	for _, duration := range times {
		totalTime += duration
	}

	// 3. Ensure first iteration isn't dramatically slower than last
	firstTime := times[0]
	lastTime := times[iterations-1]

	// Allow last to be up to 2x slower (due to accumulation)
	maxAllowedRatio := float64(2.0)
	actualRatio := float64(lastTime) / float64(firstTime)

	t.Logf("Performance: first=%v, last=%v, ratio=%.2f", firstTime, lastTime, actualRatio)

	AssertTrue(t, actualRatio < maxAllowedRatio, "performance degradation should be minimal")
}
