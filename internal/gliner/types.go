// Package gliner provides a client for the GLiNER NER service.
//
// GLiNER (Generalist Model for Named Entity Recognition) is a zero-shot NER model
// that can extract any entity types specified at runtime. This package provides
// a Go client for communicating with the Python GLiNER service via HTTP REST API.
package gliner

// Entity represents an extracted named entity.
type Entity struct {
	Text  string  `json:"text"`  // Entity text
	Label string  `json:"label"` // Entity type/label
	Start int     `json:"start"` // Start character index
	End   int     `json:"end"`   // End character index
	Score float64 `json:"score"` // Confidence score (0.0-1.0)
}

// ExtractRequest represents a request to extract entities.
type ExtractRequest struct {
	Text         string   `json:"text"`          // Text to analyze
	EntityTypes  []string `json:"entity_types"`  // Entity types to extract
	Threshold    float64  `json:"threshold"`     // Confidence threshold (default: 0.3)
}

// ExtractResponse represents the response from entity extraction.
type ExtractResponse struct {
	Entities    []Entity `json:"entities"`     // Extracted entities
	EntityCount int      `json:"entity_count"` // Number of entities found
	TextLength  int      `json:"text_length"`  // Length of input text
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status      string `json:"status"`       // Service status ("healthy")
	ModelLoaded bool   `json:"model_loaded"` // Whether model is loaded
	ModelName   string `json:"model_name"`   // Model identifier
}

// ModelInfoResponse represents model metadata.
type ModelInfoResponse struct {
	ModelName  string `json:"model_name"`  // Model identifier
	Loaded     bool   `json:"loaded"`      // Whether model is loaded
	ModelType  string `json:"model_type"`  // Model architecture
	Parameters string `json:"parameters"`  // Model size (e.g., "340M")
	License    string `json:"license"`     // Model license
}

// Config holds configuration for the GLiNER client.
type Config struct {
	// ServiceURL is the base URL of the GLiNER service
	// Default: "http://localhost:8765"
	ServiceURL string

	// Timeout for HTTP requests
	// Default: 5 seconds
	Timeout int // seconds

	// MaxRetries for failed requests
	// Default: 2
	MaxRetries int

	// Enabled determines whether to use GLiNER
	// If false, client methods return ErrDisabled
	// Default: true
	Enabled bool

	// FallbackToPattern determines whether to fall back to pattern matcher
	// when GLiNER service is unavailable
	// Default: true
	FallbackToPattern bool
}

// DefaultConfig returns the default client configuration.
func DefaultConfig() *Config {
	return &Config{
		ServiceURL:        "http://localhost:8765",
		Timeout:           5,
		MaxRetries:        2,
		Enabled:           true,
		FallbackToPattern: true,
	}
}
