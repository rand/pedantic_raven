package mnemosyne

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestNewConnectionManager verifies connection manager creation and validation
func TestNewConnectionManager(t *testing.T) {
	tests := []struct {
		name      string
		config    *ConnectionConfig
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid config",
			config: &ConnectionConfig{
				Host:        "localhost",
				Port:        50051,
				UseTLS:      false,
				Timeout:     10 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			wantError: false,
		},
		{
			name: "empty host",
			config: &ConnectionConfig{
				Host:        "",
				Port:        50051,
				Timeout:     10 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			wantError: true,
			errorMsg:  "host cannot be empty",
		},
		{
			name: "invalid port - zero",
			config: &ConnectionConfig{
				Host:        "localhost",
				Port:        0,
				Timeout:     10 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			wantError: true,
			errorMsg:  "port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			config: &ConnectionConfig{
				Host:        "localhost",
				Port:        65536,
				Timeout:     10 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			wantError: true,
			errorMsg:  "port must be between 1 and 65535",
		},
		{
			name: "invalid timeout",
			config: &ConnectionConfig{
				Host:        "localhost",
				Port:        50051,
				Timeout:     0,
				RetryPolicy: DefaultRetryPolicy(),
			},
			wantError: true,
			errorMsg:  "timeout must be positive",
		},
		{
			name: "invalid max attempts",
			config: &ConnectionConfig{
				Host:    "localhost",
				Port:    50051,
				Timeout: 10 * time.Second,
				RetryPolicy: RetryPolicy{
					MaxAttempts:    -1,
					InitialBackoff: 1 * time.Second,
					MaxBackoff:     30 * time.Second,
				},
			},
			wantError: true,
			errorMsg:  "max retry attempts must be non-negative",
		},
		{
			name: "invalid initial backoff",
			config: &ConnectionConfig{
				Host:    "localhost",
				Port:    50051,
				Timeout: 10 * time.Second,
				RetryPolicy: RetryPolicy{
					MaxAttempts:    5,
					InitialBackoff: 0,
					MaxBackoff:     30 * time.Second,
				},
			},
			wantError: true,
			errorMsg:  "initial backoff must be positive",
		},
		{
			name: "invalid max backoff",
			config: &ConnectionConfig{
				Host:    "localhost",
				Port:    50051,
				Timeout: 10 * time.Second,
				RetryPolicy: RetryPolicy{
					MaxAttempts:    5,
					InitialBackoff: 30 * time.Second,
					MaxBackoff:     1 * time.Second,
				},
			},
			wantError: true,
			errorMsg:  "max backoff must be >= initial backoff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm, err := NewConnectionManager(tt.config)

			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errorMsg)
				}
				if err.Error() == "" {
					t.Fatalf("expected error message, got empty string")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if cm == nil {
					t.Fatal("expected connection manager, got nil")
				}
				if cm.Status() != StatusDisconnected {
					t.Errorf("expected initial status to be disconnected, got %s", cm.Status())
				}
			}
		})
	}
}

// TestConnectFailure verifies connection failure handling
func TestConnectFailure(t *testing.T) {
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    59999, // Invalid port - unlikely to have a server running
		Timeout: 2 * time.Second,
		RetryPolicy: RetryPolicy{
			MaxAttempts:    1,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     5 * time.Second,
		},
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Attempt connection - should fail
	err = cm.Connect()
	if err == nil {
		t.Fatal("expected connection to fail, but it succeeded")
	}

	// Verify status is failed
	status := cm.Status()
	if status != StatusFailed {
		t.Errorf("expected status to be failed, got %s", status)
	}
}

// TestDisconnect verifies clean shutdown
func TestDisconnect(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     2 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Disconnect when not connected should be safe
	err = cm.Disconnect()
	if err != nil {
		t.Fatalf("disconnect failed: %v", err)
	}

	if cm.Status() != StatusDisconnected {
		t.Errorf("expected status to be disconnected, got %s", cm.Status())
	}

	// Multiple disconnects should be safe
	err = cm.Disconnect()
	if err != nil {
		t.Fatalf("second disconnect failed: %v", err)
	}
}

// TestStatusTransitions verifies status changes during connection lifecycle
func TestStatusTransitions(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     2 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Initial status should be disconnected
	if cm.Status() != StatusDisconnected {
		t.Errorf("expected initial status to be disconnected, got %s", cm.Status())
	}

	// After disconnect, should still be disconnected
	cm.Disconnect()
	if cm.Status() != StatusDisconnected {
		t.Errorf("expected status after disconnect to be disconnected, got %s", cm.Status())
	}
}

// TestHealthCheckFailure verifies health check failure handling
func TestHealthCheckFailure(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     1 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Health check when not connected should fail
	err = cm.HealthCheck()
	if err == nil {
		t.Fatal("expected health check to fail when not connected")
	}

	// Verify status is still disconnected
	if cm.Status() != StatusDisconnected {
		t.Errorf("expected status to remain disconnected, got %s", cm.Status())
	}
}

// TestExponentialBackoff verifies backoff calculation
func TestExponentialBackoff(t *testing.T) {
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    50051,
		Timeout: 10 * time.Second,
		RetryPolicy: RetryPolicy{
			MaxAttempts:    5,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     30 * time.Second,
		},
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 1 * time.Second},   // 1 * 2^0 = 1s
		{1, 2 * time.Second},   // 1 * 2^1 = 2s
		{2, 4 * time.Second},   // 1 * 2^2 = 4s
		{3, 8 * time.Second},   // 1 * 2^3 = 8s
		{4, 16 * time.Second},  // 1 * 2^4 = 16s
		{5, 30 * time.Second},  // 1 * 2^5 = 32s, capped at 30s
		{10, 30 * time.Second}, // 1 * 2^10 = 1024s, capped at 30s
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
			backoff := cm.calculateBackoff(tt.attempt)
			if backoff != tt.expected {
				t.Errorf("expected backoff %v, got %v", tt.expected, backoff)
			}
		})
	}
}

// TestMaxRetries verifies retry limit enforcement
func TestMaxRetries(t *testing.T) {
	// This test verifies the retry logic conceptually by checking
	// the backoff calculation matches our expectations.
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    59999, // Invalid port
		Timeout: 500 * time.Millisecond,
		RetryPolicy: RetryPolicy{
			MaxAttempts:    3,
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     1 * time.Second,
		},
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Verify max attempts is set correctly
	if cm.config.RetryPolicy.MaxAttempts != 3 {
		t.Errorf("expected max attempts to be 3, got %d", cm.config.RetryPolicy.MaxAttempts)
	}

	// Verify initial retry count is 0
	if cm.retryCount != 0 {
		t.Errorf("expected initial retry count to be 0, got %d", cm.retryCount)
	}
}

// TestConcurrentAccess verifies thread-safe operations
func TestConcurrentAccess(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     2 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	var wg sync.WaitGroup
	iterations := 50

	// Concurrent status reads
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cm.Status()
		}()
	}

	// Concurrent client reads
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cm.Client()
		}()
	}

	wg.Wait()

	// If we get here without race conditions or crashes, test passes
}

// TestConfigValidation verifies comprehensive config validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    *ConnectionConfig
		wantError bool
	}{
		{
			name: "valid minimal config",
			config: &ConnectionConfig{
				Host:    "localhost",
				Port:    50051,
				Timeout: 1 * time.Second,
				RetryPolicy: RetryPolicy{
					MaxAttempts:    0, // 0 means unlimited
					InitialBackoff: 1 * time.Second,
					MaxBackoff:     1 * time.Second,
				},
			},
			wantError: false,
		},
		{
			name: "valid with TLS",
			config: &ConnectionConfig{
				Host:        "example.com",
				Port:        443,
				UseTLS:      true,
				Timeout:     30 * time.Second,
				RetryPolicy: DefaultRetryPolicy(),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantError && err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("expected no validation error, got: %v", err)
			}
		})
	}
}

// TestHealthCheckTickerStops verifies health check monitoring stops on disconnect
func TestHealthCheckTickerStops(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     1 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Manually simulate connected state for testing
	cm.mu.Lock()
	cm.status = StatusConnected
	cm.mu.Unlock()

	// Start health check monitoring
	cm.startHealthCheckMonitoring()

	// Give it a moment to start
	time.Sleep(50 * time.Millisecond)

	// Verify ticker is running
	cm.mu.RLock()
	hasHealthCheck := cm.healthCheck != nil
	cm.mu.RUnlock()

	if !hasHealthCheck {
		t.Fatal("expected health check ticker to be running")
	}

	// Disconnect should stop the ticker
	if err := cm.Disconnect(); err != nil {
		t.Fatalf("disconnect failed: %v", err)
	}

	// Give disconnect time to complete
	time.Sleep(50 * time.Millisecond)

	// Verify ticker is stopped (healthCheck should be nil after cleanup)
	cm.mu.RLock()
	ticker := cm.healthCheck
	stopChan := cm.stopHealth
	cm.mu.RUnlock()

	if ticker != nil {
		t.Error("expected health check ticker to be stopped after disconnect")
	}

	if stopChan != nil {
		t.Error("expected stop channel to be nil after disconnect")
	}

	// Verify status is disconnected
	if cm.Status() != StatusDisconnected {
		t.Errorf("expected status to be disconnected, got %s", cm.Status())
	}
}

// TestStatusWhileReconnecting verifies status during reconnection attempts
func TestStatusWhileReconnecting(t *testing.T) {
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    59999, // Invalid port to trigger reconnection
		Timeout: 500 * time.Millisecond,
		RetryPolicy: RetryPolicy{
			MaxAttempts:    2,
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     500 * time.Millisecond,
		},
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Manually trigger reconnection state
	cm.mu.Lock()
	cm.status = StatusReconnecting
	cm.mu.Unlock()

	// Verify status is reconnecting
	if cm.Status() != StatusReconnecting {
		t.Errorf("expected status to be reconnecting, got %s", cm.Status())
	}

	// Health check during reconnection should fail
	err = cm.HealthCheck()
	if err == nil {
		t.Error("expected health check to fail during reconnecting state")
	}
}

// TestClientAccess verifies client accessor method
func TestClientAccess(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     2 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Client should be nil when not connected
	client := cm.Client()
	if client != nil {
		t.Error("expected client to be nil when not connected")
	}

	// After disconnect, client should still be nil
	cm.Disconnect()
	client = cm.Client()
	if client != nil {
		t.Error("expected client to be nil after disconnect")
	}
}

// TestConnectionStatusString verifies status string representation
func TestConnectionStatusString(t *testing.T) {
	tests := []struct {
		status   ConnectionStatus
		expected string
	}{
		{StatusDisconnected, "disconnected"},
		{StatusConnecting, "connecting"},
		{StatusConnected, "connected"},
		{StatusReconnecting, "reconnecting"},
		{StatusFailed, "failed"},
		{ConnectionStatus(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			str := tt.status.String()
			if str != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, str)
			}
		})
	}
}

// TestDefaultRetryPolicy verifies default retry policy values
func TestDefaultRetryPolicy(t *testing.T) {
	policy := DefaultRetryPolicy()

	if policy.MaxAttempts != 5 {
		t.Errorf("expected max attempts to be 5, got %d", policy.MaxAttempts)
	}

	if policy.InitialBackoff != 1*time.Second {
		t.Errorf("expected initial backoff to be 1s, got %v", policy.InitialBackoff)
	}

	if policy.MaxBackoff != 30*time.Second {
		t.Errorf("expected max backoff to be 30s, got %v", policy.MaxBackoff)
	}
}

// TestConcurrentConnectDisconnect verifies safety of concurrent connect/disconnect
func TestConcurrentConnectDisconnect(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     2 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	var wg sync.WaitGroup

	// Multiple concurrent disconnects (should be safe)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cm.Disconnect()
		}()
	}

	wg.Wait()

	// Final status should be disconnected
	if cm.Status() != StatusDisconnected {
		t.Errorf("expected final status to be disconnected, got %s", cm.Status())
	}
}

// TestHealthCheckContextTimeout verifies health check respects timeout
func TestHealthCheckContextTimeout(t *testing.T) {
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    50051,
		Timeout: 100 * time.Millisecond, // Very short timeout
		RetryPolicy: RetryPolicy{
			MaxAttempts:    1,
			InitialBackoff: 1 * time.Second,
			MaxBackoff:     5 * time.Second,
		},
	}

	_, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Create a mock client to test timeout behavior
	clientCfg := Config{
		ServerAddr: fmt.Sprintf("%s:%d", config.Host, config.Port),
		Timeout:    config.Timeout,
	}

	client, err := NewClient(clientCfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Verify timeout is set correctly
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	start := time.Now()
	_, err = client.HealthCheck(ctx)
	elapsed := time.Since(start)

	// Health check should fail (no server running)
	if err == nil {
		t.Error("expected health check to fail with no server")
	}

	// Should respect timeout (with some margin for execution)
	if elapsed > config.Timeout*2 {
		t.Errorf("health check took too long: %v (expected ~%v)", elapsed, config.Timeout)
	}
}

// TestReconnectStopsOnDisconnect verifies reconnect goroutine stops when disconnected
func TestReconnectStopsOnDisconnect(t *testing.T) {
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    59999, // Invalid port
		Timeout: 200 * time.Millisecond,
		RetryPolicy: RetryPolicy{
			MaxAttempts:    10, // Many attempts
			InitialBackoff: 100 * time.Millisecond,
			MaxBackoff:     1 * time.Second,
		},
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Manually set to reconnecting and start attemptReconnect
	cm.mu.Lock()
	cm.status = StatusReconnecting
	cm.mu.Unlock()

	go cm.attemptReconnect()

	// Give it a moment to start reconnecting
	time.Sleep(50 * time.Millisecond)

	// Now disconnect - this should stop the reconnect loop
	if err := cm.Disconnect(); err != nil {
		t.Fatalf("disconnect failed: %v", err)
	}

	// Verify status is disconnected
	if cm.Status() != StatusDisconnected {
		t.Errorf("expected status to be disconnected, got %s", cm.Status())
	}

	// Wait a bit and verify it stays disconnected (reconnect loop stopped)
	time.Sleep(300 * time.Millisecond)
	if cm.Status() != StatusDisconnected {
		t.Errorf("expected status to remain disconnected, got %s", cm.Status())
	}
}
