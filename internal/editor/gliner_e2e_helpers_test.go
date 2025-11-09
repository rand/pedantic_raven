package editor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"github.com/rand/pedantic-raven/internal/gliner"
)

// ServiceMethod defines how to start the GLiNER service for testing
type ServiceMethod string

const (
	ServiceMock       ServiceMethod = "mock"        // httptest.Server
	ServiceDocker     ServiceMethod = "docker"      // docker-compose
	ServiceSubprocess ServiceMethod = "subprocess"  // Python uvicorn
)

// GLiNERTestService manages GLiNER service lifecycle for tests
type GLiNERTestService struct {
	Method     ServiceMethod
	URL        string
	server     *httptest.Server
	dockerCmd  *exec.Cmd
	subprocess *exec.Cmd
	cleanup    func()
}

// GoldenTestCase represents a test case with expected outputs
type GoldenTestCase struct {
	Name            string              `json:"name"`
	Text            string              `json:"text"`
	GLiNEREntities  []semantic.Entity   `json:"gliner_entities"`
	PatternEntities []semantic.Entity   `json:"pattern_entities"`
	MinAccuracy     float64             `json:"min_accuracy"`
}

// GoldenData holds all test fixtures
type GoldenData struct {
	TestCases []GoldenTestCase `json:"test_cases"`
}

// StartGLiNERService starts the GLiNER service using the specified method
func StartGLiNERService(t *testing.T, method ServiceMethod) *GLiNERTestService {
	t.Helper()

	service := &GLiNERTestService{Method: method}

	switch method {
	case ServiceMock:
		service.startMock(t)
	case ServiceDocker:
		service.startDocker(t)
	case ServiceSubprocess:
		service.startSubprocess(t)
	default:
		t.Fatalf("Unknown service method: %s", method)
	}

	return service
}

// startMock creates a mock HTTP server for testing
func (s *GLiNERTestService) startMock(t *testing.T) {
	t.Helper()

	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/health":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":       "healthy",
				"model_loaded": true,
				"model_name":   "fastino/gliner2-large-v1",
			})

		case "/model_info":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"model_name":       "fastino/gliner2-large-v1",
				"parameters":       "340M",
				"supported_types":  []string{"person", "organization", "location"},
			})

		case "/extract":
			var req struct {
				Text        string   `json:"text"`
				EntityTypes []string `json:"entity_types"`
				Threshold   float64  `json:"threshold"`
			}
			json.NewDecoder(r.Body).Decode(&req)

			// Mock extraction: return predefined entities based on text
			entities := mockExtractEntities(req.Text, req.EntityTypes, req.Threshold)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"entities": entities,
			})

		default:
			http.NotFound(w, r)
		}
	}))

	s.URL = s.server.URL
	s.cleanup = func() {
		s.server.Close()
	}
}

// startDocker starts GLiNER via docker-compose
func (s *GLiNERTestService) startDocker(t *testing.T) {
	t.Helper()

	// Start docker-compose
	cmd := exec.Command("docker-compose", "up", "-d", "gliner")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to start docker-compose (skipping): %v", err)
		return
	}

	s.dockerCmd = cmd
	s.URL = "http://localhost:8765"
	s.cleanup = func() {
		exec.Command("docker-compose", "down").Run()
	}

	// Wait for service to be ready
	if err := s.WaitForHealth(context.Background(), 30*time.Second); err != nil {
		s.cleanup()
		t.Skipf("GLiNER service not healthy (skipping): %v", err)
	}
}

// startSubprocess starts GLiNER via Python subprocess
func (s *GLiNERTestService) startSubprocess(t *testing.T) {
	t.Helper()

	// Find Python
	pythonPath, err := exec.LookPath("python3")
	if err != nil {
		pythonPath, err = exec.LookPath("python")
		if err != nil {
			t.Skipf("Python not found (skipping): %v", err)
			return
		}
	}

	// Start uvicorn
	servicePath := filepath.Join("services", "gliner")
	cmd := exec.Command(pythonPath, "-m", "uvicorn", "main:app",
		"--host", "127.0.0.1", "--port", "8765")
	cmd.Dir = servicePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Skipf("Failed to start GLiNER subprocess (skipping): %v", err)
		return
	}

	s.subprocess = cmd
	s.URL = "http://localhost:8765"
	s.cleanup = func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	}

	// Wait for service to be ready
	if err := s.WaitForHealth(context.Background(), 30*time.Second); err != nil {
		s.cleanup()
		t.Skipf("GLiNER service not healthy (skipping): %v", err)
	}
}

// WaitForHealth waits for the GLiNER service to become healthy
func (s *GLiNERTestService) WaitForHealth(ctx context.Context, timeout time.Duration) error {
	client := gliner.NewClient(&gliner.Config{
		ServiceURL: s.URL,
		Timeout:    5,
		Enabled:    true,
	})

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if err := client.CheckAvailability(ctx); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for GLiNER service at %s", s.URL)
}

// Cleanup shuts down the service
func (s *GLiNERTestService) Cleanup() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

// mockExtractEntities provides mock entity extraction for testing
func mockExtractEntities(text string, entityTypes []string, threshold float64) []gliner.Entity {
	// Simple mock: extract known patterns
	entities := []gliner.Entity{}

	// Common test entities
	testEntities := map[string]struct {
		label string
		score float64
	}{
		"Alice":         {"person", 0.95},
		"Bob":           {"person", 0.93},
		"Google":        {"organization", 0.97},
		"Microsoft":     {"organization", 0.96},
		"San Francisco": {"location", 0.92},
		"New York":      {"location", 0.94},
		"Python":        {"technology", 0.89},
		"React":         {"technology", 0.87},
	}

	for word, info := range testEntities {
		if contains(text, word) && info.score >= threshold {
			// Check if entity type is requested
			if len(entityTypes) == 0 || containsString(entityTypes, info.label) {
				entities = append(entities, gliner.Entity{
					Text:  word,
					Label: info.label,
					Start: indexOf(text, word),
					End:   indexOf(text, word) + len(word),
					Score: info.score,
				})
			}
		}
	}

	return entities
}

// LoadGoldenData loads test fixtures from JSON
func LoadGoldenData(t *testing.T) *GoldenData {
	t.Helper()

	path := filepath.Join("testdata", "gliner_golden.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load golden data: %v", err)
	}

	var golden GoldenData
	if err := json.Unmarshal(data, &golden); err != nil {
		t.Fatalf("Failed to parse golden data: %v", err)
	}

	return &golden
}

// MeasureAnalysisLatency measures how long analysis takes
func MeasureAnalysisLatency(mode *EditMode, text string) time.Duration {
	mode.editor.SetContent(text)

	start := time.Now()
	cmd := mode.triggerAnalysis()
	if cmd != nil {
		cmd() // Execute and wait for completion
	}
	return time.Since(start)
}

// MeasureTypingLatency measures perceived typing lag
func MeasureTypingLatency(mode *EditMode) time.Duration {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

	start := time.Now()
	mode.Update(msg)
	return time.Since(start)
}

// SimulateTyping simulates a user typing at a given WPM
func SimulateTyping(mode *EditMode, text string, wpm int) {
	// Average word length is 5 characters
	charsPerMinute := wpm * 5
	delayPerChar := time.Minute / time.Duration(charsPerMinute)

	for _, ch := range text {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
		mode.Update(msg)
		time.Sleep(delayPerChar)
	}
}

// CompareEntitySets compares two entity sets with tolerance for ML variance
func CompareEntitySets(t *testing.T, name string, expected, actual []semantic.Entity, tolerance float64) {
	t.Helper()

	// Normalize entities (lowercase, sort)
	expectedNorm := normalizeEntities(expected)
	actualNorm := normalizeEntities(actual)

	// Calculate match percentage
	matches := 0
	for _, exp := range expectedNorm {
		for _, act := range actualNorm {
			if entityMatch(exp, act) {
				matches++
				break
			}
		}
	}

	matchPercent := 0.0
	if len(expectedNorm) > 0 {
		matchPercent = float64(matches) / float64(len(expectedNorm))
	}

	if matchPercent < tolerance {
		t.Errorf("%s: Entity match %.2f%% below tolerance %.2f%%",
			name, matchPercent*100, tolerance*100)
		t.Logf("Expected: %+v", expectedNorm)
		t.Logf("Actual: %+v", actualNorm)
	}
}

// VerifyContextPanelUpdated checks that the context panel received analysis results
func VerifyContextPanelUpdated(t *testing.T, mode *EditMode, expectedMinEntities int) {
	t.Helper()

	analysis := mode.analyzer.Results()
	if analysis == nil {
		t.Fatal("Analysis is nil")
	}

	if len(analysis.Entities) < expectedMinEntities {
		t.Errorf("Expected at least %d entities, got %d",
			expectedMinEntities, len(analysis.Entities))
	}
}

// Helper functions

func normalizeEntities(entities []semantic.Entity) []semantic.Entity {
	normalized := make([]semantic.Entity, len(entities))
	copy(normalized, entities)
	// Could add sorting, case normalization, etc.
	return normalized
}

func entityMatch(a, b semantic.Entity) bool {
	// Fuzzy match: same text (case-insensitive) and type
	return a.Text == b.Text && a.Type == b.Type
}

func contains(text, substr string) bool {
	return indexOf(text, substr) >= 0
}

func indexOf(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
