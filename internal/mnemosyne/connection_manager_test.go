package mnemosyne

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestConnectionManagerOfflineMode verifies offline mode operations
func TestConnectionManagerOfflineMode(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        59999, // Invalid port to trigger offline mode
		Timeout:     500 * time.Millisecond,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Attempt connection - should fail and enter offline mode
	err = cm.Connect()
	if err == nil {
		t.Fatal("expected connection to fail")
	}

	// Verify offline mode is active
	if !cm.IsOffline() {
		t.Error("expected offline mode to be active after connection failure")
	}

	// Verify last error is set
	if cm.GetLastError() == nil {
		t.Error("expected last error to be set")
	}
}

// TestConnectionManagerSyncQueue verifies sync queue operations
func TestConnectionManagerSyncQueue(t *testing.T) {
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

	// Get sync queue
	syncQueue := cm.GetSyncQueue()
	if syncQueue == nil {
		t.Fatal("expected sync queue to be non-nil")
	}

	// Initially should be empty
	if syncQueue.Len() != 0 {
		t.Errorf("expected empty sync queue, got %d operations", syncQueue.Len())
	}
}

// TestConnectionManagerOfflineCache verifies offline cache operations
func TestConnectionManagerOfflineCache(t *testing.T) {
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

	// Get offline cache
	cache := cm.GetOfflineCache()
	if cache == nil {
		t.Fatal("expected offline cache to be non-nil")
	}

	// Initially should be empty
	if cache.Len() != 0 {
		t.Errorf("expected empty offline cache, got %d memories", cache.Len())
	}
}

// TestConnectionManagerErrorCallback verifies error callback mechanism
func TestConnectionManagerErrorCallback(t *testing.T) {
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

	// Set up error callback
	var callbackCalled bool
	var callbackErr error
	var mu sync.Mutex

	cm.SetErrorCallback(func(err error) {
		mu.Lock()
		defer mu.Unlock()
		callbackCalled = true
		callbackErr = err
	})

	// Trigger offline mode which should call the callback
	cm.enterOfflineMode(ErrConnection)

	// Give callback time to execute
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if !callbackCalled {
		t.Error("expected error callback to be called")
	}

	if callbackErr == nil {
		t.Error("expected error to be passed to callback")
	}
}

// TestConnectionManagerTriggerSyncWhileOffline verifies sync fails when offline
func TestConnectionManagerTriggerSyncWhileOffline(t *testing.T) {
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

	// Enter offline mode
	cm.enterOfflineMode(ErrConnection)

	// Attempt sync - should fail
	count, err := cm.TriggerSync()
	if err == nil {
		t.Fatal("expected sync to fail while offline")
	}

	if count != 0 {
		t.Errorf("expected 0 synced operations, got %d", count)
	}
}

// TestConnectionManagerTriggerSyncNotConnected verifies sync fails when not connected
func TestConnectionManagerTriggerSyncNotConnected(t *testing.T) {
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

	// Don't connect, but also not in offline mode
	count, err := cm.TriggerSync()
	if err != ErrNotConnected {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 synced operations, got %d", count)
	}
}

// TestConnectionManagerExitOfflineMode verifies offline mode exit
func TestConnectionManagerExitOfflineMode(t *testing.T) {
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

	// Enter offline mode
	cm.enterOfflineMode(ErrConnection)
	if !cm.IsOffline() {
		t.Fatal("expected offline mode to be active")
	}

	// Exit offline mode (this will trigger async sync in background)
	cm.exitOfflineMode()

	// Give it a moment to process
	time.Sleep(50 * time.Millisecond)

	if cm.IsOffline() {
		t.Error("expected offline mode to be inactive after exit")
	}
}

// TestConnectionManagerExitOfflineModeWhenNotOffline verifies exit is safe when not offline
func TestConnectionManagerExitOfflineModeWhenNotOffline(t *testing.T) {
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

	// Exit offline mode when not in offline mode (should be safe)
	cm.exitOfflineMode()

	if cm.IsOffline() {
		t.Error("expected offline mode to remain inactive")
	}
}

// TestConnectionManagerEnterOfflineModeMultipleTimes verifies idempotency
func TestConnectionManagerEnterOfflineModeMultipleTimes(t *testing.T) {
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

	// Enter offline mode multiple times (should be safe)
	cm.enterOfflineMode(ErrConnection)
	cm.enterOfflineMode(ErrConnection)
	cm.enterOfflineMode(ErrConnection)

	if !cm.IsOffline() {
		t.Error("expected offline mode to be active")
	}
}

// TestConnectionManagerHealthCheckNotConnected verifies health check when not connected
func TestConnectionManagerHealthCheckNotConnected(t *testing.T) {
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
}

// TestConnectionManagerHealthCheckNilClient verifies health check with nil client
func TestConnectionManagerHealthCheckNilClient(t *testing.T) {
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

	// Manually set status to connected but client is nil
	cm.mu.Lock()
	cm.status = StatusConnected
	cm.client = nil
	cm.mu.Unlock()

	// Health check should fail with nil client
	err = cm.HealthCheck()
	if err == nil {
		t.Fatal("expected health check to fail with nil client")
	}

	// Reset status
	cm.mu.Lock()
	cm.status = StatusDisconnected
	cm.mu.Unlock()
}

// TestConnectionManagerConnectWhenAlreadyConnecting verifies concurrent connect prevention
func TestConnectionManagerConnectWhenAlreadyConnecting(t *testing.T) {
	config := &ConnectionConfig{
		Host:        "localhost",
		Port:        50051,
		Timeout:     5 * time.Second,
		RetryPolicy: DefaultRetryPolicy(),
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Set status to connecting
	cm.mu.Lock()
	cm.status = StatusConnecting
	cm.mu.Unlock()

	// Attempt to connect while already connecting (should return immediately)
	err = cm.Connect()
	if err != nil {
		t.Errorf("expected no error when connecting while already connecting, got: %v", err)
	}

	// Verify status is still connecting
	if cm.Status() != StatusConnecting {
		t.Errorf("expected status to remain connecting, got %s", cm.Status())
	}

	// Reset status
	cm.mu.Lock()
	cm.status = StatusDisconnected
	cm.mu.Unlock()
}

// TestConnectionManagerConnectWhenAlreadyConnected verifies concurrent connect prevention
func TestConnectionManagerConnectWhenAlreadyConnected(t *testing.T) {
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

	// Set status to connected
	cm.mu.Lock()
	cm.status = StatusConnected
	cm.mu.Unlock()

	// Attempt to connect while already connected (should return immediately)
	err = cm.Connect()
	if err != nil {
		t.Errorf("expected no error when connecting while already connected, got: %v", err)
	}

	// Reset status
	cm.mu.Lock()
	cm.status = StatusDisconnected
	cm.mu.Unlock()
}

// TestConnectionManagerStartHealthCheckMultipleTimes verifies health check ticker idempotency
func TestConnectionManagerStartHealthCheckMultipleTimes(t *testing.T) {
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

	// Start health check monitoring
	cm.startHealthCheckMonitoring()

	// Verify ticker is running
	cm.mu.RLock()
	firstTicker := cm.healthCheck
	cm.mu.RUnlock()

	if firstTicker == nil {
		t.Fatal("expected health check ticker to be running")
	}

	// Start again (should be idempotent)
	cm.startHealthCheckMonitoring()

	// Verify ticker is still the same
	cm.mu.RLock()
	secondTicker := cm.healthCheck
	cm.mu.RUnlock()

	if secondTicker != firstTicker {
		t.Error("expected health check ticker to remain the same")
	}

	// Clean up
	cm.stopHealthCheckMonitoring()
}

// TestConnectionManagerStopHealthCheckWhenNotRunning verifies safe stop
func TestConnectionManagerStopHealthCheckWhenNotRunning(t *testing.T) {
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

	// Stop when not running (should be safe)
	cm.stopHealthCheckMonitoring()

	cm.mu.RLock()
	ticker := cm.healthCheck
	stopChan := cm.stopHealth
	cm.mu.RUnlock()

	if ticker != nil {
		t.Error("expected health check ticker to be nil")
	}

	if stopChan != nil {
		t.Error("expected stop channel to be nil")
	}
}

// TestConnectionManagerAttemptReconnectStopsOnDisconnected verifies reconnect loop stops
func TestConnectionManagerAttemptReconnectStopsOnDisconnected(t *testing.T) {
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    59999, // Invalid port
		Timeout: 200 * time.Millisecond,
		RetryPolicy: RetryPolicy{
			MaxAttempts:    10,
			InitialBackoff: 50 * time.Millisecond,
			MaxBackoff:     500 * time.Millisecond,
		},
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Set to reconnecting and start reconnect attempt
	cm.mu.Lock()
	cm.status = StatusReconnecting
	cm.mu.Unlock()

	go cm.attemptReconnect()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Change status to disconnected (simulating manual disconnect)
	cm.mu.Lock()
	cm.status = StatusDisconnected
	cm.mu.Unlock()

	// Wait and verify it stopped
	time.Sleep(300 * time.Millisecond)

	if cm.Status() != StatusDisconnected {
		t.Errorf("expected status to be disconnected, got %s", cm.Status())
	}
}

// TestConnectionManagerAttemptReconnectMaxRetries verifies max retries enforcement
func TestConnectionManagerAttemptReconnectMaxRetries(t *testing.T) {
	config := &ConnectionConfig{
		Host:    "localhost",
		Port:    59999, // Invalid port
		Timeout: 100 * time.Millisecond,
		RetryPolicy: RetryPolicy{
			MaxAttempts:    2, // Very low
			InitialBackoff: 50 * time.Millisecond,
			MaxBackoff:     200 * time.Millisecond,
		},
	}

	cm, err := NewConnectionManager(config)
	if err != nil {
		t.Fatalf("failed to create connection manager: %v", err)
	}

	// Set to reconnecting
	cm.mu.Lock()
	cm.status = StatusReconnecting
	cm.mu.Unlock()

	// Attempt reconnect (should fail after max retries)
	cm.attemptReconnect()

	// Should be in failed state
	if cm.Status() != StatusFailed {
		t.Errorf("expected status to be failed after max retries, got %s", cm.Status())
	}
}

// TestConnectionManagerClientAccessThreadSafety verifies concurrent client access
func TestConnectionManagerClientAccessThreadSafety(t *testing.T) {
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
	iterations := 100

	// Concurrent client access
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cm.Client()
		}()
	}

	// Concurrent status access
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cm.Status()
		}()
	}

	// Concurrent offline check
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cm.IsOffline()
		}()
	}

	wg.Wait()
}

// TestConnectionManagerConcurrentOfflineModeOperations verifies thread-safe offline operations
func TestConnectionManagerConcurrentOfflineModeOperations(t *testing.T) {
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

	var wg sync.WaitGroup

	// Concurrent enter offline mode
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cm.enterOfflineMode(ErrConnection)
		}()
	}

	// Concurrent offline check
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = cm.IsOffline()
		}()
	}

	wg.Wait()

	// Should be in offline mode
	if !cm.IsOffline() {
		t.Error("expected offline mode to be active")
	}
}
