package gliner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a client for the GLiNER NER service.
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new GLiNER client with the given configuration.
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// IsEnabled returns true if the GLiNER client is enabled.
func (c *Client) IsEnabled() bool {
	return c.config.Enabled
}

// HealthCheck performs a health check on the GLiNER service.
func (c *Client) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	if !c.config.Enabled {
		return nil, wrapError("HealthCheck", ErrDisabled)
	}

	url := fmt.Sprintf("%s/health", c.config.ServiceURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, wrapError("HealthCheck", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapError("HealthCheck", ErrServiceUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, wrapError("HealthCheck", fmt.Errorf("unexpected status: %d", resp.StatusCode))
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return nil, wrapError("HealthCheck", err)
	}

	return &healthResp, nil
}

// ModelInfo retrieves information about the loaded model.
func (c *Client) ModelInfo(ctx context.Context) (*ModelInfoResponse, error) {
	if !c.config.Enabled {
		return nil, wrapError("ModelInfo", ErrDisabled)
	}

	url := fmt.Sprintf("%s/model_info", c.config.ServiceURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, wrapError("ModelInfo", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, wrapError("ModelInfo", ErrServiceUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, wrapError("ModelInfo", fmt.Errorf("unexpected status: %d", resp.StatusCode))
	}

	var info ModelInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, wrapError("ModelInfo", err)
	}

	return &info, nil
}

// ExtractEntities extracts named entities from text.
//
// Parameters:
//   - ctx: Context for cancellation
//   - text: Text to analyze
//   - entityTypes: Entity types to extract (e.g., ["person", "organization", "location"])
//   - threshold: Confidence threshold (0.0-1.0), default 0.3
//
// Returns:
//   - List of extracted entities
//   - Error if extraction fails or service unavailable
func (c *Client) ExtractEntities(
	ctx context.Context,
	text string,
	entityTypes []string,
	threshold float64,
) ([]Entity, error) {
	if !c.config.Enabled {
		return nil, wrapError("ExtractEntities", ErrDisabled)
	}

	// Validate input
	if text == "" {
		return nil, wrapError("ExtractEntities", fmt.Errorf("%w: empty text", ErrInvalidRequest))
	}
	if len(entityTypes) == 0 {
		return nil, wrapError("ExtractEntities", fmt.Errorf("%w: no entity types specified", ErrInvalidRequest))
	}

	// Set default threshold
	if threshold == 0 {
		threshold = 0.3
	}

	// Prepare request
	reqBody := ExtractRequest{
		Text:        text,
		EntityTypes: entityTypes,
		Threshold:   threshold,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, wrapError("ExtractEntities", err)
	}

	url := fmt.Sprintf("%s/extract", c.config.ServiceURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, wrapError("ExtractEntities", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Retry logic
	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, lastErr = c.httpClient.Do(req)
		if lastErr == nil && resp.StatusCode == http.StatusOK {
			break
		}

		// Close body if request failed
		if resp != nil {
			resp.Body.Close()
		}

		// Wait before retry (exponential backoff)
		if attempt < c.config.MaxRetries {
			time.Sleep(time.Duration(100*(1<<attempt)) * time.Millisecond)
		}
	}

	if lastErr != nil {
		return nil, wrapError("ExtractEntities", ErrServiceUnavailable)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, wrapError("ExtractEntities", fmt.Errorf("%w: %s", ErrExtractionFailed, string(body)))
	}

	var extractResp ExtractResponse
	if err := json.NewDecoder(resp.Body).Decode(&extractResp); err != nil {
		return nil, wrapError("ExtractEntities", err)
	}

	return extractResp.Entities, nil
}

// CheckAvailability checks if the GLiNER service is available and returns any error.
// This is useful for determining whether to fall back to pattern matching.
func (c *Client) CheckAvailability(ctx context.Context) error {
	if !c.config.Enabled {
		return ErrDisabled
	}

	_, err := c.HealthCheck(ctx)
	return err
}
