package auth

import (
	"crypto/subtle"
	"os"
)

// TokenAuth provides simple token-based authentication for single-user deployment.
// It uses constant-time comparison to prevent timing attacks.
type TokenAuth struct {
	token   string
	enabled bool
}

// NewTokenAuth creates a new TokenAuth instance by reading the PEDANTIC_RAVEN_TOKEN
// environment variable. If the variable is not set, authentication is disabled.
func NewTokenAuth() *TokenAuth {
	token := os.Getenv("PEDANTIC_RAVEN_TOKEN")
	return &TokenAuth{
		token:   token,
		enabled: token != "",
	}
}

// IsEnabled returns true if authentication is enabled (token is set).
func (a *TokenAuth) IsEnabled() bool {
	return a.enabled
}

// Validate checks if the provided token matches the configured token.
// If authentication is disabled, it always returns true.
// Uses constant-time comparison to prevent timing attacks.
func (a *TokenAuth) Validate(providedToken string) bool {
	if !a.enabled {
		return true // Auth disabled, allow all
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(
		[]byte(a.token),
		[]byte(providedToken),
	) == 1
}
