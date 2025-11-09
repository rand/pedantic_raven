package mnemosyne

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- RetryConfig Tests ---

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()

	if cfg.MaxAttempts != 5 {
		t.Errorf("MaxAttempts = %d, want 5", cfg.MaxAttempts)
	}

	if cfg.InitialBackoff != 1*time.Second {
		t.Errorf("InitialBackoff = %v, want 1s", cfg.InitialBackoff)
	}

	if cfg.MaxBackoff != 30*time.Second {
		t.Errorf("MaxBackoff = %v, want 30s", cfg.MaxBackoff)
	}

	if cfg.Multiplier != 2.0 {
		t.Errorf("Multiplier = %f, want 2.0", cfg.Multiplier)
	}
}

// --- RetryWithBackoff Tests ---

func TestRetryWithBackoffSuccess(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:    3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempt := 0
	operation := func() error {
		attempt++
		if attempt < 3 {
			return status.Error(codes.Unavailable, "service unavailable")
		}
		return nil // Success on 3rd attempt
	}

	ctx := context.Background()
	err := RetryWithBackoff(ctx, cfg, operation)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if attempt != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempt)
	}
}

func TestRetryWithBackoffFailure(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:    3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempt := 0
	expectedErr := status.Error(codes.Unavailable, "service unavailable")
	operation := func() error {
		attempt++
		return expectedErr
	}

	ctx := context.Background()
	err := RetryWithBackoff(ctx, cfg, operation)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if attempt != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempt)
	}
}

func TestRetryWithBackoffNonRetryableError(t *testing.T) {
	cfg := DefaultRetryConfig()

	attempt := 0
	operation := func() error {
		attempt++
		return status.Error(codes.InvalidArgument, "invalid input")
	}

	ctx := context.Background()
	err := RetryWithBackoff(ctx, cfg, operation)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Should fail on first attempt for non-retryable error
	if attempt != 1 {
		t.Errorf("Expected 1 attempt for non-retryable error, got %d", attempt)
	}
}

func TestRetryWithBackoffContextCancelled(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:    5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
		Multiplier:     2.0,
	}

	attempt := 0
	operation := func() error {
		attempt++
		return status.Error(codes.Unavailable, "service unavailable")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := RetryWithBackoff(ctx, cfg, operation)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Should detect cancellation early
	if attempt > 1 {
		t.Errorf("Expected <= 1 attempt with cancelled context, got %d", attempt)
	}
}

func TestRetryWithBackoffBackoffProgression(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:    4,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempts := []time.Time{}
	operation := func() error {
		attempts = append(attempts, time.Now())
		return status.Error(codes.Unavailable, "unavailable")
	}

	ctx := context.Background()
	RetryWithBackoff(ctx, cfg, operation)

	// Verify exponential backoff between attempts
	// Expected: ~10ms, ~20ms, ~40ms
	if len(attempts) != 4 {
		t.Fatalf("Expected 4 attempts, got %d", len(attempts))
	}

	for i := 1; i < len(attempts); i++ {
		delay := attempts[i].Sub(attempts[i-1])
		expectedMin := time.Duration(float64(cfg.InitialBackoff) * float64(uint(1)<<uint(i-1)))

		// Allow some variance (50% below expected)
		minDelay := expectedMin / 2
		if delay < minDelay {
			t.Errorf("Attempt %d: delay %v is less than expected minimum %v", i, delay, minDelay)
		}
	}
}

// --- calculateBackoff Tests ---

func TestCalculateBackoff(t *testing.T) {
	cfg := RetryConfig{
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		Multiplier:     2.0,
	}

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{0, 1 * time.Second},   // 1 * 2^0 = 1
		{1, 2 * time.Second},   // 1 * 2^1 = 2
		{2, 4 * time.Second},   // 1 * 2^2 = 4
		{3, 8 * time.Second},   // 1 * 2^3 = 8
		{4, 16 * time.Second},  // 1 * 2^4 = 16
		{5, 30 * time.Second},  // 1 * 2^5 = 32, capped at 30
		{10, 30 * time.Second}, // Much larger, still capped at 30
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := calculateBackoff(cfg, tt.attempt)
			if got != tt.want {
				t.Errorf("calculateBackoff(attempt=%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}

// --- IsRetryableError (alias) Tests ---

func TestIsRetryableErrorAlias(t *testing.T) {
	// Test that the alias works
	retryableErr := status.Error(codes.Unavailable, "unavailable")
	nonRetryableErr := status.Error(codes.InvalidArgument, "invalid")

	if !IsRetryableError(retryableErr) {
		t.Error("IsRetryableError should return true for retryable error")
	}

	if IsRetryableError(nonRetryableErr) {
		t.Error("IsRetryableError should return false for non-retryable error")
	}

	if IsRetryableError(nil) {
		t.Error("IsRetryableError should return false for nil")
	}
}

// --- Edge Cases ---

func TestRetryWithBackoffImmediateSuccess(t *testing.T) {
	cfg := DefaultRetryConfig()

	attempt := 0
	operation := func() error {
		attempt++
		return nil // Success on first attempt
	}

	ctx := context.Background()
	err := RetryWithBackoff(ctx, cfg, operation)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if attempt != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempt)
	}
}

func TestRetryWithBackoffZeroMaxAttempts(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:    0,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	attempt := 0
	operation := func() error {
		attempt++
		return errors.New("error")
	}

	ctx := context.Background()
	err := RetryWithBackoff(ctx, cfg, operation)

	if err == nil {
		t.Error("Expected error with 0 max attempts")
	}

	if attempt != 0 {
		t.Errorf("Expected 0 attempts with MaxAttempts=0, got %d", attempt)
	}
}
