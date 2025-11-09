package mnemosyne

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryConfig configures retry behavior.
type RetryConfig struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Multiplier     float64
}

// DefaultRetryConfig returns sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:    5,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
	}
}

// RetryWithBackoff executes operation with exponential backoff.
// Returns the last error if all attempts fail.
func RetryWithBackoff(ctx context.Context, cfg RetryConfig, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			if lastErr != nil {
				return fmt.Errorf("retry cancelled after %d attempts: %w (last error: %v)", attempt, ctx.Err(), lastErr)
			}
			return ctx.Err()
		default:
		}

		// Execute operation
		err := operation()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryable(err) {
			return fmt.Errorf("non-retryable error after %d attempts: %w", attempt+1, err)
		}

		// Last attempt, don't sleep
		if attempt == cfg.MaxAttempts-1 {
			break
		}

		// Calculate backoff duration
		backoff := calculateBackoff(cfg, attempt)

		// Sleep with context cancellation support
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled during backoff after %d attempts: %w (last error: %v)", attempt+1, ctx.Err(), lastErr)
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("max retry attempts (%d) exceeded: %w", cfg.MaxAttempts, lastErr)
}

// calculateBackoff calculates the backoff duration for a given attempt.
// Uses exponential backoff: min(initialBackoff * multiplier^attempt, maxBackoff)
func calculateBackoff(cfg RetryConfig, attempt int) time.Duration {
	backoff := float64(cfg.InitialBackoff) * math.Pow(cfg.Multiplier, float64(attempt))
	maxBackoff := float64(cfg.MaxBackoff)

	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	return time.Duration(backoff)
}

// IsRetryableError is a deprecated alias for IsRetryable.
// Use IsRetryable instead.
func IsRetryableError(err error) bool {
	return IsRetryable(err)
}
