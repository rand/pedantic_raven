package semantic

import (
	"context"
	"errors"
	"testing"
)

// --- Pattern Extractor Tests ---

func TestPatternExtractorBasic(t *testing.T) {
	// Test basic pattern extraction with technology and person entities
	extractor := NewPatternExtractor()
	ctx := context.Background()

	text := "The Developer uses HTTP and PostgreSQL"

	entities, err := extractor.ExtractEntities(ctx, text, nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entities) == 0 {
		t.Fatal("Expected to extract entities")
	}

	// Check for specific entities
	foundDeveloper := false
	foundHTTP := false
	foundPostgreSQL := false

	for _, entity := range entities {
		switch entity.Text {
		case "Developer":
			foundDeveloper = true
			if entity.Type != EntityPerson {
				t.Errorf("Expected Developer to be EntityPerson, got %v", entity.Type)
			}
		case "HTTP":
			foundHTTP = true
			if entity.Type != EntityTechnology {
				t.Errorf("Expected HTTP to be EntityTechnology, got %v", entity.Type)
			}
		case "PostgreSQL":
			foundPostgreSQL = true
			if entity.Type != EntityTechnology {
				t.Errorf("Expected PostgreSQL to be EntityTechnology, got %v", entity.Type)
			}
		}
	}

	if !foundDeveloper {
		t.Error("Expected to find Developer entity")
	}

	if !foundHTTP {
		t.Error("Expected to find HTTP entity")
	}

	if !foundPostgreSQL {
		t.Error("Expected to find PostgreSQL entity")
	}
}

func TestPatternExtractorWithEntityTypes(t *testing.T) {
	// Test pattern extraction with specific entity types requested
	extractor := NewPatternExtractor()
	ctx := context.Background()

	text := "Developer Alice uses PostgreSQL in New York"

	// Only request person entities
	entities, err := extractor.ExtractEntities(ctx, text, []string{"person"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only find person entities
	for _, entity := range entities {
		if entity.Type != EntityPerson {
			t.Errorf("Expected only person entities, got %v: %s", entity.Type, entity.Text)
		}
	}
}

func TestPatternExtractorName(t *testing.T) {
	// Test that pattern extractor returns correct name
	extractor := NewPatternExtractor()

	if extractor.Name() != "Pattern" {
		t.Errorf("Expected name 'Pattern', got %s", extractor.Name())
	}
}

func TestPatternExtractorIsAvailable(t *testing.T) {
	// Test that pattern extractor is always available
	extractor := NewPatternExtractor()
	ctx := context.Background()

	if !extractor.IsAvailable(ctx) {
		t.Error("Expected pattern extractor to always be available")
	}
}

func TestPatternExtractorContextCancellation(t *testing.T) {
	// Test that pattern extractor respects context cancellation
	extractor := NewPatternExtractor()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	text := "Some text to analyze"

	_, err := extractor.ExtractEntities(ctx, text, nil)

	if err == nil {
		t.Error("Expected error for cancelled context")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestPatternExtractorEntityCount(t *testing.T) {
	// Test that duplicate entities are counted correctly
	extractor := NewPatternExtractor()
	ctx := context.Background()

	text := "Developer Alice and Developer Bob use HTTP and HTTP Server"

	entities, err := extractor.ExtractEntities(ctx, text, nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Find Developer entity and check count
	for _, entity := range entities {
		if entity.Text == "Developer" {
			if entity.Count != 2 {
				t.Errorf("Expected Developer count of 2, got %d", entity.Count)
			}
			return
		}
	}

	t.Error("Expected to find Developer entity")
}

// --- Mock GLiNER Extractor for Testing ---

// mockGLiNERExtractor is a mock implementation of EntityExtractor for testing
type mockGLiNERExtractor struct {
	name      string
	enabled   bool
	available bool
	entities  []Entity
	err       error
}

func (m *mockGLiNERExtractor) ExtractEntities(ctx context.Context, text string, entityTypes []string) ([]Entity, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.entities, nil
}

func (m *mockGLiNERExtractor) Name() string {
	return m.name
}

func (m *mockGLiNERExtractor) IsAvailable(ctx context.Context) bool {
	return m.available
}

// --- GLiNER Extractor Tests ---
//
// Note: These tests use the mockGLiNERExtractor since the real GLiNERExtractor
// requires a concrete *gliner.Client (not an interface). For integration tests
// with real HTTP servers, see the client_test.go file.

func TestGLiNERExtractorWithMockClient(t *testing.T) {
	// Test GLiNER-style extractor with mock returning entities
	mockExtractor := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: true,
		entities: []Entity{
			{
				Text: "Alice",
				Type: EntityPerson,
				Span: Span{Start: 0, End: 5},
			},
			{
				Text: "OpenAI",
				Type: EntityOrganization,
				Span: Span{Start: 16, End: 22},
			},
		},
	}

	ctx := context.Background()
	text := "Alice works at OpenAI"

	entities, err := mockExtractor.ExtractEntities(ctx, text, []string{"person", "organization"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entities) != 2 {
		t.Fatalf("Expected 2 entities, got %d", len(entities))
	}

	// Check first entity (Alice)
	if entities[0].Text != "Alice" {
		t.Errorf("Expected first entity text 'Alice', got %s", entities[0].Text)
	}

	if entities[0].Type != EntityPerson {
		t.Errorf("Expected first entity type EntityPerson, got %v", entities[0].Type)
	}

	// Check second entity (OpenAI)
	if entities[1].Text != "OpenAI" {
		t.Errorf("Expected second entity text 'OpenAI', got %s", entities[1].Text)
	}

	if entities[1].Type != EntityOrganization {
		t.Errorf("Expected second entity type EntityOrganization, got %v", entities[1].Type)
	}
}

func TestGLiNERExtractorDefaultTypes(t *testing.T) {
	// Test GLiNER-style extractor uses default behavior
	mockExtractor := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: true,
		entities:  []Entity{},
	}

	ctx := context.Background()
	text := "Some text"

	// Call with empty entity types - should work
	_, err := mockExtractor.ExtractEntities(ctx, text, nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestGLiNERExtractorDisabled(t *testing.T) {
	// Test GLiNER-style extractor returns error when disabled
	mockExtractor := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: false,
		err:       ErrExtractorDisabled,
	}

	ctx := context.Background()
	text := "Some text"

	_, err := mockExtractor.ExtractEntities(ctx, text, []string{"person"})

	if err == nil {
		t.Fatal("Expected error when extractor is disabled")
	}

	if err != ErrExtractorDisabled {
		t.Errorf("Expected ErrExtractorDisabled, got %v", err)
	}
}

func TestGLiNERExtractorUnavailable(t *testing.T) {
	// Test GLiNER-style extractor returns error when unavailable
	mockExtractor := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: false,
		err:       ErrExtractorUnavailable,
	}

	ctx := context.Background()
	text := "Some text"

	_, err := mockExtractor.ExtractEntities(ctx, text, []string{"person"})

	if err == nil {
		t.Fatal("Expected error when extractor unavailable")
	}

	if err != ErrExtractorUnavailable {
		t.Errorf("Expected ErrExtractorUnavailable, got %v", err)
	}
}

func TestGLiNERExtractorName(t *testing.T) {
	// Test that GLiNER-style extractor returns correct name
	mockExtractor := &mockGLiNERExtractor{
		name: "GLiNER",
	}

	if mockExtractor.Name() != "GLiNER" {
		t.Errorf("Expected name 'GLiNER', got %s", mockExtractor.Name())
	}
}

func TestGLiNERExtractorIsAvailable(t *testing.T) {
	tests := []struct {
		name      string
		available bool
		want      bool
	}{
		{"Available", true, true},
		{"Unavailable", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExtractor := &mockGLiNERExtractor{
				name:      "GLiNER",
				available: tt.available,
			}

			ctx := context.Background()
			result := mockExtractor.IsAvailable(ctx)

			if result != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestGLiNERExtractorDeduplication(t *testing.T) {
	// Test entity deduplication behavior
	mockExtractor := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: true,
		entities: []Entity{
			{
				Text:  "alice",
				Type:  EntityPerson,
				Span:  Span{Start: 0, End: 5},
				Count: 1,
			},
			{
				Text:  "Alice",
				Type:  EntityPerson,
				Span:  Span{Start: 20, End: 25},
				Count: 1,
			},
			{
				Text:  "ALICE",
				Type:  EntityPerson,
				Span:  Span{Start: 40, End: 45},
				Count: 1,
			},
		},
	}

	ctx := context.Background()
	text := "alice mentioned Alice and ALICE"

	entities, err := mockExtractor.ExtractEntities(ctx, text, []string{"person"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Mock returns all entities (deduplication would be done by real extractor)
	if len(entities) != 3 {
		t.Errorf("Expected 3 entities from mock, got %d", len(entities))
	}
}

// --- Hybrid Extractor Tests ---

func TestHybridExtractorPrimarySuccess(t *testing.T) {
	// Test hybrid extractor uses primary when available
	mockPrimary := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: true,
		entities: []Entity{
			{
				Text: "Alice",
				Type: EntityPerson,
				Span: Span{Start: 0, End: 5},
			},
		},
	}

	fallbackExtractor := NewPatternExtractor()

	hybrid := NewHybridExtractor(mockPrimary, fallbackExtractor, true)
	ctx := context.Background()

	text := "Alice works here"

	entities, err := hybrid.ExtractEntities(ctx, text, []string{"person"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entities) != 1 {
		t.Fatalf("Expected 1 entity from primary, got %d", len(entities))
	}

	// Check that we got GLiNER result (not pattern result)
	if entities[0].Text != "Alice" {
		t.Errorf("Expected entity from primary extractor")
	}
}

func TestHybridExtractorFallback(t *testing.T) {
	// Test hybrid extractor falls back when primary unavailable
	mockPrimary := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: false, // Not available
	}

	fallbackExtractor := NewPatternExtractor()

	hybrid := NewHybridExtractor(mockPrimary, fallbackExtractor, true)
	ctx := context.Background()

	text := "The Developer uses HTTP"

	entities, err := hybrid.ExtractEntities(ctx, text, nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should get results from fallback (pattern extractor)
	if len(entities) == 0 {
		t.Fatal("Expected entities from fallback extractor")
	}

	// Verify we got pattern results
	foundDeveloper := false
	for _, entity := range entities {
		if entity.Text == "Developer" {
			foundDeveloper = true
			break
		}
	}

	if !foundDeveloper {
		t.Error("Expected to find Developer from pattern extractor")
	}
}

func TestHybridExtractorFallbackDisabled(t *testing.T) {
	// Test hybrid extractor returns error when fallback disabled
	mockPrimary := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: false,
	}

	fallbackExtractor := NewPatternExtractor()

	hybrid := NewHybridExtractor(mockPrimary, fallbackExtractor, false) // Fallback disabled
	ctx := context.Background()

	text := "Some text"

	_, err := hybrid.ExtractEntities(ctx, text, []string{"person"})

	if err == nil {
		t.Fatal("Expected error when fallback disabled and primary unavailable")
	}

	if err != ErrNoExtractorAvailable {
		t.Errorf("Expected ErrNoExtractorAvailable, got %v", err)
	}
}

func TestHybridExtractorName(t *testing.T) {
	tests := []struct {
		name             string
		primaryAvailable bool
		fallbackEnabled  bool
		expectedName     string
	}{
		{"Primary available", true, true, "GLiNER"},
		{"Primary unavailable, fallback enabled", false, true, "Pattern (Fallback)"},
		{"Primary unavailable, fallback disabled", false, false, "Hybrid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrimary := &mockGLiNERExtractor{
				name:      "GLiNER",
				available: tt.primaryAvailable,
			}

			fallbackExtractor := NewPatternExtractor()

			hybrid := NewHybridExtractor(mockPrimary, fallbackExtractor, tt.fallbackEnabled)

			result := hybrid.Name()

			if result != tt.expectedName {
				t.Errorf("Name() = %s, want %s", result, tt.expectedName)
			}
		})
	}
}

func TestHybridExtractorIsAvailable(t *testing.T) {
	tests := []struct {
		name             string
		primaryAvailable bool
		fallbackEnabled  bool
		want             bool
	}{
		{"Primary available", true, false, true},
		{"Primary unavailable, fallback enabled", false, true, true},
		{"Both unavailable", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrimary := &mockGLiNERExtractor{
				name:      "GLiNER",
				available: tt.primaryAvailable,
			}

			fallbackExtractor := NewPatternExtractor()

			hybrid := NewHybridExtractor(mockPrimary, fallbackExtractor, tt.fallbackEnabled)
			ctx := context.Background()

			result := hybrid.IsAvailable(ctx)

			if result != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestHybridExtractorSetPrimary(t *testing.T) {
	// Test setting primary extractor
	fallbackExtractor := NewPatternExtractor()
	hybrid := NewHybridExtractor(nil, fallbackExtractor, true)

	mockPrimary := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: true,
	}

	hybrid.SetPrimary(mockPrimary)

	// Verify by checking name
	if hybrid.Name() != "GLiNER" {
		t.Error("Expected primary to be set")
	}
}

func TestHybridExtractorSetFallback(t *testing.T) {
	// Test setting fallback extractor
	mockPrimary := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: false,
	}

	hybrid := NewHybridExtractor(mockPrimary, nil, true)

	newFallback := NewPatternExtractor()
	hybrid.SetFallback(newFallback)

	// Verify by checking that extraction works with unavailable primary
	ctx := context.Background()
	_, err := hybrid.ExtractEntities(ctx, "test", nil)

	if err != nil {
		t.Error("Expected fallback to be set and working")
	}
}

func TestHybridExtractorEnableFallback(t *testing.T) {
	// Test enabling/disabling fallback
	mockPrimary := &mockGLiNERExtractor{
		name:      "GLiNER",
		available: false,
	}
	fallbackExtractor := NewPatternExtractor()

	hybrid := NewHybridExtractor(mockPrimary, fallbackExtractor, false)
	ctx := context.Background()

	// Should fail with fallback disabled
	_, err := hybrid.ExtractEntities(ctx, "test", nil)
	if err == nil {
		t.Error("Expected error with fallback disabled")
	}

	// Enable fallback
	hybrid.EnableFallback(true)

	// Should succeed with fallback enabled
	_, err = hybrid.ExtractEntities(ctx, "test", nil)
	if err != nil {
		t.Errorf("Expected success with fallback enabled, got %v", err)
	}
}

// --- Interface Compliance Tests ---

func TestExtractorInterface(t *testing.T) {
	// Test that all extractors implement the EntityExtractor interface
	var _ EntityExtractor = (*PatternExtractor)(nil)
	var _ EntityExtractor = (*GLiNERExtractor)(nil)
	var _ EntityExtractor = (*HybridExtractor)(nil)

	// This test passes if it compiles
	t.Log("All extractors implement EntityExtractor interface")
}

func TestExtractorInterfaceMethods(t *testing.T) {
	// Test that all extractors have the required interface methods
	ctx := context.Background()

	extractors := []EntityExtractor{
		NewPatternExtractor(),
		&mockGLiNERExtractor{name: "GLiNER", available: true},
		NewHybridExtractor(
			NewPatternExtractor(),
			NewPatternExtractor(),
			true,
		),
	}

	for i, extractor := range extractors {
		// Test Name() method
		name := extractor.Name()
		if name == "" {
			t.Errorf("Extractor %d: Name() returned empty string", i)
		}

		// Test IsAvailable() method
		_ = extractor.IsAvailable(ctx)

		// Test ExtractEntities() method
		_, err := extractor.ExtractEntities(ctx, "test text", []string{"person"})
		// Error is ok for some extractors, just verify method exists
		_ = err
	}
}
