package gliner

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// --- Client Creation Tests ---

func TestNewClientWithDefaultConfig(t *testing.T) {
	// Test creating client with nil config (should use defaults)
	client := NewClient(nil)

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.config == nil {
		t.Fatal("Expected config to be initialized")
	}

	if !client.config.Enabled {
		t.Error("Expected client to be enabled by default")
	}

	if client.config.ServiceURL != "http://localhost:8765" {
		t.Errorf("Expected default service URL, got %s", client.config.ServiceURL)
	}

	if client.config.Timeout != 5 {
		t.Errorf("Expected default timeout of 5 seconds, got %d", client.config.Timeout)
	}

	if client.config.MaxRetries != 2 {
		t.Errorf("Expected default max retries of 2, got %d", client.config.MaxRetries)
	}

	if !client.config.FallbackToPattern {
		t.Error("Expected fallback to be enabled by default")
	}
}

func TestNewClientWithCustomConfig(t *testing.T) {
	// Test creating client with custom config
	config := &Config{
		ServiceURL:        "http://example.com:9000",
		Timeout:           10,
		MaxRetries:        5,
		Enabled:           false,
		FallbackToPattern: false,
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.config.ServiceURL != config.ServiceURL {
		t.Errorf("Expected custom service URL %s, got %s",
			config.ServiceURL, client.config.ServiceURL)
	}

	if client.config.Timeout != config.Timeout {
		t.Errorf("Expected custom timeout %d, got %d",
			config.Timeout, client.config.Timeout)
	}

	if client.config.MaxRetries != config.MaxRetries {
		t.Errorf("Expected custom max retries %d, got %d",
			config.MaxRetries, client.config.MaxRetries)
	}

	if client.config.Enabled {
		t.Error("Expected client to be disabled")
	}

	if client.config.FallbackToPattern {
		t.Error("Expected fallback to be disabled")
	}
}

func TestNewClientHTTPClientConfiguration(t *testing.T) {
	// Test that HTTP client is configured with correct timeout
	config := &Config{
		ServiceURL: "http://localhost:8765",
		Timeout:    15,
		MaxRetries: 2,
		Enabled:    true,
	}

	client := NewClient(config)

	if client.httpClient == nil {
		t.Fatal("Expected HTTP client to be initialized")
	}

	expectedTimeout := time.Duration(config.Timeout) * time.Second
	if client.httpClient.Timeout != expectedTimeout {
		t.Errorf("Expected HTTP client timeout %v, got %v",
			expectedTimeout, client.httpClient.Timeout)
	}
}

// --- Client Disabled Tests ---

func TestClientIsEnabled(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
		want    bool
	}{
		{"Enabled client", true, true},
		{"Disabled client", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				ServiceURL: "http://localhost:8765",
				Enabled:    tt.enabled,
			}

			client := NewClient(config)
			result := client.IsEnabled()

			if result != tt.want {
				t.Errorf("IsEnabled() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestClientDisabledHealthCheck(t *testing.T) {
	// Test that HealthCheck returns error when client is disabled
	config := &Config{
		ServiceURL: "http://localhost:8765",
		Enabled:    false,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.HealthCheck(ctx)

	if err == nil {
		t.Fatal("Expected error when client is disabled")
	}

	if !IsDisabled(err) {
		t.Errorf("Expected ErrDisabled, got %v", err)
	}
}

func TestClientDisabledModelInfo(t *testing.T) {
	// Test that ModelInfo returns error when client is disabled
	config := &Config{
		ServiceURL: "http://localhost:8765",
		Enabled:    false,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.ModelInfo(ctx)

	if err == nil {
		t.Fatal("Expected error when client is disabled")
	}

	if !IsDisabled(err) {
		t.Errorf("Expected ErrDisabled, got %v", err)
	}
}

func TestClientDisabledExtractEntities(t *testing.T) {
	// Test that ExtractEntities returns error when client is disabled
	config := &Config{
		ServiceURL: "http://localhost:8765",
		Enabled:    false,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.ExtractEntities(ctx, "test text", []string{"person"}, 0.3)

	if err == nil {
		t.Fatal("Expected error when client is disabled")
	}

	if !IsDisabled(err) {
		t.Errorf("Expected ErrDisabled, got %v", err)
	}
}

func TestClientDisabledCheckAvailability(t *testing.T) {
	// Test that CheckAvailability returns error when client is disabled
	config := &Config{
		ServiceURL: "http://localhost:8765",
		Enabled:    false,
	}

	client := NewClient(config)
	ctx := context.Background()

	err := client.CheckAvailability(ctx)

	if err == nil {
		t.Fatal("Expected error when client is disabled")
	}

	if !IsDisabled(err) {
		t.Errorf("Expected ErrDisabled, got %v", err)
	}
}

// --- Health Check Tests ---

func TestHealthCheckSuccess(t *testing.T) {
	// Create mock server that returns successful health response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected /health path, got %s", r.URL.Path)
		}

		response := HealthResponse{
			Status:      "healthy",
			ModelLoaded: true,
			ModelName:   "urchade/gliner_small-v2.5",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServiceURL: server.URL,
		Timeout:    5,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	health, err := client.HealthCheck(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if health == nil {
		t.Fatal("Expected health response")
	}

	if health.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", health.Status)
	}

	if !health.ModelLoaded {
		t.Error("Expected model to be loaded")
	}

	if health.ModelName != "urchade/gliner_small-v2.5" {
		t.Errorf("Expected model name 'urchade/gliner_small-v2.5', got %s", health.ModelName)
	}
}

func TestHealthCheckServerError(t *testing.T) {
	// Create mock server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := &Config{
		ServiceURL: server.URL,
		Timeout:    5,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.HealthCheck(ctx)

	if err == nil {
		t.Fatal("Expected error for server error")
	}
}

func TestHealthCheckServiceUnavailable(t *testing.T) {
	// Use invalid URL to simulate service unavailable
	config := &Config{
		ServiceURL: "http://localhost:9999",
		Timeout:    1, // Short timeout
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.HealthCheck(ctx)

	if err == nil {
		t.Fatal("Expected error for unavailable service")
	}

	if !IsUnavailable(err) {
		t.Errorf("Expected ErrServiceUnavailable, got %v", err)
	}
}

// --- Model Info Tests ---

func TestModelInfoSuccess(t *testing.T) {
	// Create mock server that returns model info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/model_info" {
			t.Errorf("Expected /model_info path, got %s", r.URL.Path)
		}

		response := ModelInfoResponse{
			ModelName:  "urchade/gliner_small-v2.5",
			Loaded:     true,
			ModelType:  "GLiNER",
			Parameters: "340M",
			License:    "Apache 2.0",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServiceURL: server.URL,
		Timeout:    5,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	info, err := client.ModelInfo(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if info == nil {
		t.Fatal("Expected model info response")
	}

	if info.ModelName != "urchade/gliner_small-v2.5" {
		t.Errorf("Expected model name 'urchade/gliner_small-v2.5', got %s", info.ModelName)
	}

	if !info.Loaded {
		t.Error("Expected model to be loaded")
	}

	if info.Parameters != "340M" {
		t.Errorf("Expected parameters '340M', got %s", info.Parameters)
	}
}

// --- Extract Entities Tests ---

func TestExtractEntitiesSuccess(t *testing.T) {
	// Create mock server that returns entities
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/extract" {
			t.Errorf("Expected /extract path, got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Parse request
		var req ExtractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		response := ExtractResponse{
			Entities: []Entity{
				{
					Text:  "Alice",
					Label: "person",
					Start: 0,
					End:   5,
					Score: 0.95,
				},
				{
					Text:  "New York",
					Label: "location",
					Start: 15,
					End:   23,
					Score: 0.92,
				},
			},
			EntityCount: 2,
			TextLength:  len(req.Text),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServiceURL: server.URL,
		Timeout:    5,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	entities, err := client.ExtractEntities(ctx, "Alice lives in New York", []string{"person", "location"}, 0.5)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entities) != 2 {
		t.Fatalf("Expected 2 entities, got %d", len(entities))
	}

	// Check first entity
	if entities[0].Text != "Alice" {
		t.Errorf("Expected first entity text 'Alice', got %s", entities[0].Text)
	}

	if entities[0].Label != "person" {
		t.Errorf("Expected first entity label 'person', got %s", entities[0].Label)
	}

	// Check second entity
	if entities[1].Text != "New York" {
		t.Errorf("Expected second entity text 'New York', got %s", entities[1].Text)
	}

	if entities[1].Label != "location" {
		t.Errorf("Expected second entity label 'location', got %s", entities[1].Label)
	}
}

func TestExtractEntitiesEmptyText(t *testing.T) {
	// Test that empty text returns error
	config := &Config{
		ServiceURL: "http://localhost:8765",
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.ExtractEntities(ctx, "", []string{"person"}, 0.3)

	if err == nil {
		t.Fatal("Expected error for empty text")
	}
}

func TestExtractEntitiesNoEntityTypes(t *testing.T) {
	// Test that no entity types returns error
	config := &Config{
		ServiceURL: "http://localhost:8765",
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.ExtractEntities(ctx, "test text", []string{}, 0.3)

	if err == nil {
		t.Fatal("Expected error for no entity types")
	}
}

func TestExtractEntitiesDefaultThreshold(t *testing.T) {
	// Test that default threshold is applied when 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ExtractRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Check that threshold is set to default 0.3
		if req.Threshold != 0.3 {
			t.Errorf("Expected default threshold 0.3, got %f", req.Threshold)
		}

		response := ExtractResponse{
			Entities:    []Entity{},
			EntityCount: 0,
			TextLength:  len(req.Text),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServiceURL: server.URL,
		Timeout:    5,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.ExtractEntities(ctx, "test text", []string{"person"}, 0)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestExtractEntitiesRetry(t *testing.T) {
	// Test retry logic with a server that fails first time then succeeds
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++

		if attemptCount == 1 {
			// First attempt fails
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Second attempt succeeds
		response := ExtractResponse{
			Entities:    []Entity{},
			EntityCount: 0,
			TextLength:  10,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServiceURL: server.URL,
		Timeout:    5,
		MaxRetries: 2,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	_, err := client.ExtractEntities(ctx, "test text", []string{"person"}, 0.3)

	if err != nil {
		t.Fatalf("Expected success after retry, got error: %v", err)
	}

	if attemptCount < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", attemptCount)
	}
}

// --- Check Availability Tests ---

func TestCheckAvailabilityHealthy(t *testing.T) {
	// Create mock server that returns healthy status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:      "healthy",
			ModelLoaded: true,
			ModelName:   "urchade/gliner_small-v2.5",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		ServiceURL: server.URL,
		Timeout:    5,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	err := client.CheckAvailability(ctx)

	if err != nil {
		t.Errorf("Expected no error for healthy service, got %v", err)
	}
}

func TestCheckAvailabilityUnhealthy(t *testing.T) {
	// Use invalid URL to simulate unhealthy service
	config := &Config{
		ServiceURL: "http://localhost:9999",
		Timeout:    1,
		Enabled:    true,
	}

	client := NewClient(config)
	ctx := context.Background()

	err := client.CheckAvailability(ctx)

	if err == nil {
		t.Fatal("Expected error for unhealthy service")
	}
}
