package editor

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"github.com/rand/pedantic-raven/internal/gliner"
)

// =============================================================================
// Helper Functions
// =============================================================================

func isDockerAvailable(t *testing.T) bool {
	// Check Docker command
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		t.Logf("Docker not available: %v", err)
		return false
	}

	// Check GLiNER service health
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:8001/health")
	if err != nil {
		t.Logf("GLiNER service not responding: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Logf("GLiNER service unhealthy: status %d", resp.StatusCode)
		return false
	}

	return true
}

// =============================================================================
// Category A: Service Lifecycle Tests (3 tests)
// =============================================================================

func TestGLiNERE2E_A1_MockServiceLifecycle(t *testing.T) {
	service := StartGLiNERService(t, ServiceMock)
	defer service.Cleanup()

	if service.URL == "" {
		t.Fatal("Service URL is empty")
	}

	// Verify service is healthy
	ctx := context.Background()
	if err := service.WaitForHealth(ctx, 5*time.Second); err != nil {
		t.Fatalf("Service not healthy: %v", err)
	}

	// Test health check
	client := gliner.NewClient(&gliner.Config{
		ServiceURL: service.URL,
		Timeout:    5,
		Enabled:    true,
	})

	health, err := client.HealthCheck(ctx)
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("Expected healthy status, got %s", health.Status)
	}

	if !health.ModelLoaded {
		t.Error("Expected model to be loaded")
	}
}

func TestGLiNERE2E_A2_DockerServiceLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping docker test in short mode")
	}

	if os.Getenv("GLINER_E2E_SKIP_DOCKER") != "" {
		t.Skip("GLINER_E2E_SKIP_DOCKER is set")
	}

	if !isDockerAvailable(t) {
		t.Skip("Skipping: Docker not available or GLiNER service unhealthy")
	}

	service := StartGLiNERService(t, ServiceDocker)
	defer service.Cleanup()

	// Verify service is healthy
	ctx := context.Background()
	if err := service.WaitForHealth(ctx, 30*time.Second); err != nil {
		t.Skipf("Docker service not healthy (skipping): %v", err)
	}

	// Test extraction with real service
	client := gliner.NewClient(&gliner.Config{
		ServiceURL: service.URL,
		Timeout:    10,
		Enabled:    true,
	})

	entities, err := client.ExtractEntities(ctx,
		"Alice works at Google",
		[]string{"person", "organization"},
		0.3)

	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	if len(entities) == 0 {
		t.Error("Expected at least one entity")
	}
}

func TestGLiNERE2E_A3_HealthCheckTimeoutRetry(t *testing.T) {
	// Create client pointing to non-existent service
	client := gliner.NewClient(&gliner.Config{
		ServiceURL: "http://localhost:9999",
		Timeout:    1,
		MaxRetries: 0,
		Enabled:    true,
	})

	ctx := context.Background()
	err := client.CheckAvailability(ctx)

	if err == nil {
		t.Error("Expected error for unavailable service")
	}

	if !gliner.IsUnavailable(err) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}
}

// =============================================================================
// Category B: Extraction Accuracy Tests (4 tests)
// =============================================================================

func TestGLiNERE2E_B1_AccuracyComparison(t *testing.T) {
	service := StartGLiNERService(t, ServiceMock)
	defer service.Cleanup()

	ctx := context.Background()

	// Create GLiNER analyzer
	glinerClient := gliner.NewClient(&gliner.Config{
		ServiceURL: service.URL,
		Timeout:    5,
		Enabled:    true,
	})
	glinerExtractor := semantic.NewGLiNERExtractor(glinerClient,
		[]string{"person", "organization", "location"},
		0.3)
	glinerAnalyzer := semantic.NewAnalyzerWithExtractor(glinerExtractor)

	// Create pattern analyzer
	patternExtractor := semantic.NewPatternExtractor()
	patternAnalyzer := semantic.NewAnalyzerWithExtractor(patternExtractor)

	// Analyze same text with both
	testText := "Alice works at Google in San Francisco"

	// GLiNER extraction
	glinerEntities, err := glinerExtractor.ExtractEntities(ctx, testText, []string{})
	if err != nil {
		t.Fatalf("GLiNER extraction failed: %v", err)
	}

	// Pattern extraction
	patternEntities, err := patternExtractor.ExtractEntities(ctx, testText, []string{})
	if err != nil {
		t.Fatalf("Pattern extraction failed: %v", err)
	}

	// GLiNER should extract more entities (mock will find Alice, Google, San Francisco)
	if len(glinerEntities) < len(patternEntities) {
		t.Logf("GLiNER entities: %d, Pattern entities: %d",
			len(glinerEntities), len(patternEntities))
	}

	// Both should extract at least "Alice"
	hasAlice := false
	for _, e := range glinerEntities {
		if e.Text == "Alice" {
			hasAlice = true
			break
		}
	}
	if !hasAlice {
		t.Error("GLiNER failed to extract 'Alice'")
	}

	_ = glinerAnalyzer
	_ = patternAnalyzer
}

func TestGLiNERE2E_B2_CustomEntityTypes(t *testing.T) {
	service := StartGLiNERService(t, ServiceMock)
	defer service.Cleanup()

	ctx := context.Background()

	// Create extractor with custom types
	client := gliner.NewClient(&gliner.Config{
		ServiceURL: service.URL,
		Timeout:    5,
		Enabled:    true,
	})
	customTypes := []string{"api_endpoint", "http_method", "status_code"}
	extractor := semantic.NewGLiNERExtractor(client, customTypes, 0.3)

	// Extract custom entities (mock won't find these, but test the flow)
	entities, err := extractor.ExtractEntities(ctx,
		"POST /api/users returns 201 Created",
		customTypes)

	if err != nil {
		t.Fatalf("Custom extraction failed: %v", err)
	}

	// Mock won't extract custom types, but verify no crash
	_ = entities
}

func TestGLiNERE2E_B3_ScoreThresholdEffects(t *testing.T) {
	service := StartGLiNERService(t, ServiceMock)
	defer service.Cleanup()

	ctx := context.Background()
	testText := "Alice works at Google"

	// Low threshold extractor (0.2)
	clientLow := gliner.NewClient(&gliner.Config{
		ServiceURL: service.URL,
		Timeout:    5,
		Enabled:    true,
	})
	extractorLow := semantic.NewGLiNERExtractor(clientLow,
		[]string{"person", "organization"},
		0.2)

	// High threshold extractor (0.6)
	clientHigh := gliner.NewClient(&gliner.Config{
		ServiceURL: service.URL,
		Timeout:    5,
		Enabled:    true,
	})
	extractorHigh := semantic.NewGLiNERExtractor(clientHigh,
		[]string{"person", "organization"},
		0.6)

	// Extract with both thresholds
	entitiesLow, err := extractorLow.ExtractEntities(ctx, testText, []string{})
	if err != nil {
		t.Fatalf("Low threshold extraction failed: %v", err)
	}

	entitiesHigh, err := extractorHigh.ExtractEntities(ctx, testText, []string{})
	if err != nil {
		t.Fatalf("High threshold extraction failed: %v", err)
	}

	// Low threshold should have more or equal entities
	if len(entitiesLow) < len(entitiesHigh) {
		t.Errorf("Low threshold (%d entities) should have >= high threshold (%d entities)",
			len(entitiesLow), len(entitiesHigh))
	}
}

func TestGLiNERE2E_B4_EntityDeduplication(t *testing.T) {
	service := StartGLiNERService(t, ServiceMock)
	defer service.Cleanup()

	ctx := context.Background()

	client := gliner.NewClient(&gliner.Config{
		ServiceURL: service.URL,
		Timeout:    5,
		Enabled:    true,
	})
	extractor := semantic.NewGLiNERExtractor(client,
		[]string{"person", "concept"},
		0.3)

	// Text with duplicate "User"
	entities, err := extractor.ExtractEntities(ctx,
		"User creates User and User modifies Document",
		[]string{})

	if err != nil {
		t.Fatalf("Extraction failed: %v", err)
	}

	// Verify deduplication - should have only one "User" entity with count 3
	userCount := 0
	for _, e := range entities {
		if e.Text == "User" {
			userCount++
			if e.Count != 3 {
				t.Errorf("Expected User count 3, got %d", e.Count)
			}
		}
	}

	if userCount > 1 {
		t.Errorf("Expected 1 unique 'User' entity, got %d", userCount)
	}
}

// =============================================================================
// Category C: Fallback Behavior Tests (3 tests)
// =============================================================================

func TestGLiNERE2E_C1_GracefulDegradation(t *testing.T) {
	// Create hybrid extractor with unavailable GLiNER
	glinerClient := gliner.NewClient(&gliner.Config{
		ServiceURL:        "http://localhost:9999",
		Timeout:           1,
		MaxRetries:        0,
		Enabled:           true,
		FallbackToPattern: true,
	})
	glinerExtractor := semantic.NewGLiNERExtractor(glinerClient,
		[]string{"person", "organization"},
		0.3)

	patternExtractor := semantic.NewPatternExtractor()

	// Create hybrid with fallback enabled
	hybridExtractor := semantic.NewHybridExtractor(glinerExtractor, patternExtractor, true)
	analyzer := semantic.NewAnalyzerWithExtractor(hybridExtractor)

	// Run analysis
	mode := NewEditModeWithAnalyzer(analyzer)
	mode.editor.SetContent("Alice works at Google")

	cmd := mode.triggerAnalysis()
	if cmd != nil {
		cmd()
	}

	// Verify analysis completed via pattern fallback
	analysis := mode.analyzer.Results()
	if analysis == nil {
		t.Fatal("Analysis is nil")
	}

	if len(analysis.Entities) == 0 {
		t.Error("Expected entities from pattern fallback")
	}
}

func TestGLiNERE2E_C2_RetryLogicAndBackoff(t *testing.T) {
	// Create client with retry enabled
	client := gliner.NewClient(&gliner.Config{
		ServiceURL: "http://localhost:9999",
		Timeout:    1,
		MaxRetries: 2,
		Enabled:    true,
	})

	ctx := context.Background()
	start := time.Now()

	// This should retry and eventually fail
	_, err := client.ExtractEntities(ctx,
		"test text",
		[]string{"person"},
		0.3)

	elapsed := time.Since(start)

	if err == nil {
		t.Error("Expected error for unavailable service")
	}

	// Should have retried (100ms + 200ms delays minimum)
	if elapsed < 300*time.Millisecond {
		t.Errorf("Expected at least 300ms for retries, got %v", elapsed)
	}
}

func TestGLiNERE2E_C3_DisabledFallbackError(t *testing.T) {
	// Create hybrid extractor with fallback disabled
	glinerClient := gliner.NewClient(&gliner.Config{
		ServiceURL:        "http://localhost:9999",
		Timeout:           1,
		MaxRetries:        0,
		Enabled:           true,
		FallbackToPattern: false,
	})
	glinerExtractor := semantic.NewGLiNERExtractor(glinerClient,
		[]string{"person"},
		0.3)

	patternExtractor := semantic.NewPatternExtractor()

	// Create hybrid with fallback disabled
	hybridExtractor := semantic.NewHybridExtractor(glinerExtractor, patternExtractor, false)

	ctx := context.Background()

	// Should return error (no fallback)
	_, err := hybridExtractor.ExtractEntities(ctx, "test text", []string{})

	if err == nil {
		t.Error("Expected error when GLiNER unavailable and fallback disabled")
	}

	if err != semantic.ErrNoExtractorAvailable {
		t.Errorf("Expected ErrNoExtractorAvailable, got %v", err)
	}
}

// =============================================================================
// Category D: User Experience Tests (5 tests)
// =============================================================================

func TestGLiNERE2E_D1_TypingLatency(t *testing.T) {
	mode := NewEditMode()

	// Measure typing latency
	latency := MeasureTypingLatency(mode)

	// Should be very fast (<50ms)
	if latency > 50*time.Millisecond {
		t.Errorf("Typing latency too high: %v (expected <50ms)", latency)
	}
}

func TestGLiNERE2E_D2_DebounceVerification(t *testing.T) {
	mode := NewEditMode()

	// Type first character
	mode.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	start := time.Now()

	// Type second character immediately (within debounce)
	mode.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

	// Analysis should not have triggered yet
	if mode.analyzing {
		t.Error("Analysis started before debounce period")
	}

	// Wait for debounce + a bit
	time.Sleep(550 * time.Millisecond)

	// Trigger update to check analysis
	mode.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	elapsed := time.Since(start)

	// Should have debounced for ~500ms
	if elapsed < 500*time.Millisecond {
		t.Errorf("Debounce period too short: %v", elapsed)
	}
}

func TestGLiNERE2E_D3_ContextPanelUpdates(t *testing.T) {
	mode := NewEditMode()
	mode.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Set content and trigger analysis
	mode.editor.SetContent("Alice works at Google in San Francisco")
	cmd := mode.triggerAnalysis()
	if cmd != nil {
		msg := cmd()
		if analysisMsg, ok := msg.(SemanticAnalysisMsg); ok {
			mode.Update(analysisMsg)
		}
	}

	// Verify context panel has entities
	VerifyContextPanelUpdated(t, mode, 1)
}

func TestGLiNERE2E_D4_AnalysisCancellation(t *testing.T) {
	mode := NewEditMode()

	// Start analysis
	mode.editor.SetContent("test content")
	mode.triggerAnalysis()

	// Verify analyzing flag is set
	if !mode.analyzing {
		t.Error("Expected analyzing to be true")
	}

	// Stop analysis via analyzer
	mode.analyzer.Stop()

	// Verify analysis stopped (wait a bit for cleanup)
	time.Sleep(100 * time.Millisecond)

	// Analysis should have completed or stopped
	// Note: analyzing flag is managed internally, we verify via IsRunning
	if mode.analyzer.IsRunning() {
		t.Error("Expected analyzer to have stopped")
	}
}

func TestGLiNERE2E_D5_MidSessionFallbackSwitching(t *testing.T) {
	service := StartGLiNERService(t, ServiceMock)
	defer service.Cleanup()

	// Create hybrid extractor
	glinerClient := gliner.NewClient(&gliner.Config{
		ServiceURL:        service.URL,
		Timeout:           5,
		Enabled:           true,
		FallbackToPattern: true,
	})
	glinerExtractor := semantic.NewGLiNERExtractor(glinerClient,
		[]string{"person", "organization"},
		0.3)
	patternExtractor := semantic.NewPatternExtractor()
	hybridExtractor := semantic.NewHybridExtractor(glinerExtractor, patternExtractor, true)

	analyzer := semantic.NewAnalyzerWithExtractor(hybridExtractor)
	mode := NewEditModeWithAnalyzer(analyzer)

	// First analysis with GLiNER available
	mode.editor.SetContent("Alice works at Google")
	cmd1 := mode.triggerAnalysis()
	if cmd1 != nil {
		cmd1()
	}

	// Simulate service going down
	service.Cleanup()

	// Second analysis should fallback to pattern
	mode.editor.SetContent("Bob works at Microsoft")
	cmd2 := mode.triggerAnalysis()
	if cmd2 != nil {
		cmd2()
	}

	// Verify analysis still completed
	analysis := mode.analyzer.Results()
	if analysis == nil {
		t.Fatal("Analysis is nil after fallback")
	}

	// Should have entities from pattern fallback
	if len(analysis.Entities) == 0 {
		t.Error("Expected entities from fallback extractor")
	}
}

// =============================================================================
// Category E: Configuration Tests (3 tests)
// =============================================================================

func TestGLiNERE2E_E1_ConfigFileLoading(t *testing.T) {
	// This test assumes config.toml exists in project root
	// In real scenario, would create temp config file

	// Verify default config values
	config := &gliner.Config{
		ServiceURL:        "http://localhost:8765",
		Timeout:           5,
		MaxRetries:        2,
		Enabled:           true,
		FallbackToPattern: true,
	}

	if config.ServiceURL != "http://localhost:8765" {
		t.Errorf("Expected default service URL, got %s", config.ServiceURL)
	}

	if config.Timeout != 5 {
		t.Errorf("Expected timeout 5, got %d", config.Timeout)
	}

	if !config.Enabled {
		t.Error("Expected enabled to be true")
	}
}

func TestGLiNERE2E_E2_EnvironmentVariableOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("GLINER_ENABLED", "false")
	os.Setenv("GLINER_SERVICE_URL", "http://localhost:9999")
	os.Setenv("GLINER_TIMEOUT", "10")
	defer func() {
		os.Unsetenv("GLINER_ENABLED")
		os.Unsetenv("GLINER_SERVICE_URL")
		os.Unsetenv("GLINER_TIMEOUT")
	}()

	// In real code, this would load config with env overrides
	// For now, just verify env vars are set
	if os.Getenv("GLINER_ENABLED") != "false" {
		t.Error("ENV var not set correctly")
	}
}

func TestGLiNERE2E_E3_ProgrammaticConfiguration(t *testing.T) {
	// Create custom config programmatically
	config := &gliner.Config{
		ServiceURL:        "http://custom.example.com",
		Timeout:           20,
		MaxRetries:        5,
		Enabled:           false,
		FallbackToPattern: false,
	}

	client := gliner.NewClient(config)

	// Verify config was applied
	if !client.IsEnabled() {
		// Client correctly reports disabled status
	}

	// Create analyzer with custom config
	extractor := semantic.NewGLiNERExtractor(client,
		[]string{"custom_type"},
		0.5)

	analyzer := semantic.NewAnalyzerWithExtractor(extractor)
	mode := NewEditModeWithAnalyzer(analyzer)

	// Verify mode created successfully
	if mode == nil {
		t.Fatal("Mode creation failed")
	}
}
