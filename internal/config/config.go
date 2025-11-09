// Package config handles application configuration loading and parsing.
package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds the application configuration.
type Config struct {
	GLiNER GLiNERConfig `toml:"gliner"`
}

// GLiNERConfig holds GLiNER-specific configuration.
type GLiNERConfig struct {
	Enabled           bool               `toml:"enabled"`
	ServiceURL        string             `toml:"service_url"`
	Timeout           int                `toml:"timeout"`
	MaxRetries        int                `toml:"max_retries"`
	FallbackToPattern bool               `toml:"fallback_to_pattern"`
	ScoreThreshold    float64            `toml:"score_threshold"`
	EntityTypes       EntityTypesConfig  `toml:"entity_types"`
}

// EntityTypesConfig holds entity type configuration.
type EntityTypesConfig struct {
	Default []string `toml:"default"`
	Custom  []string `toml:"custom"`
}

// Load loads configuration from a TOML file with environment variable overrides.
func Load(path string) (*Config, error) {
	// Set defaults
	config := &Config{
		GLiNER: GLiNERConfig{
			Enabled:           true,
			ServiceURL:        "http://localhost:8765",
			Timeout:           5,
			MaxRetries:        2,
			FallbackToPattern: true,
			ScoreThreshold:    0.3,
			EntityTypes: EntityTypesConfig{
				Default: []string{"person", "organization", "location", "technology", "concept", "product"},
				Custom:  []string{},
			},
		},
	}

	// Load from file if it exists
	if _, err := os.Stat(path); err == nil {
		if _, err := toml.DecodeFile(path, config); err != nil {
			return nil, err
		}
	}

	// Override with environment variables
	if enabled := os.Getenv("GLINER_ENABLED"); enabled != "" {
		config.GLiNER.Enabled = strings.ToLower(enabled) == "true"
	}

	if serviceURL := os.Getenv("GLINER_SERVICE_URL"); serviceURL != "" {
		config.GLiNER.ServiceURL = serviceURL
	}

	if timeout := os.Getenv("GLINER_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			config.GLiNER.Timeout = t
		}
	}

	if maxRetries := os.Getenv("GLINER_MAX_RETRIES"); maxRetries != "" {
		if mr, err := strconv.Atoi(maxRetries); err == nil {
			config.GLiNER.MaxRetries = mr
		}
	}

	if fallback := os.Getenv("GLINER_FALLBACK_TO_PATTERN"); fallback != "" {
		config.GLiNER.FallbackToPattern = strings.ToLower(fallback) == "true"
	}

	if threshold := os.Getenv("GLINER_SCORE_THRESHOLD"); threshold != "" {
		if th, err := strconv.ParseFloat(threshold, 64); err == nil {
			config.GLiNER.ScoreThreshold = th
		}
	}

	return config, nil
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		GLiNER: GLiNERConfig{
			Enabled:           true,
			ServiceURL:        "http://localhost:8765",
			Timeout:           5,
			MaxRetries:        2,
			FallbackToPattern: true,
			ScoreThreshold:    0.3,
			EntityTypes: EntityTypesConfig{
				Default: []string{"person", "organization", "location", "technology", "concept", "product"},
				Custom:  []string{},
			},
		},
	}
}
