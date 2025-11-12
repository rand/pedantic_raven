package mnemosyne

import (
	"context"
	"os"
	"testing"
	"time"
)

// --- Client Creation Tests ---

func TestNewClient(t *testing.T) {
	cfg := DefaultConfig()
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.serverAddr != cfg.ServerAddr {
		t.Errorf("Expected server address %s, got %s", cfg.ServerAddr, client.serverAddr)
	}

	if client.connected {
		t.Error("Expected client to not be connected initially")
	}
}

func TestNewClientWithEmptyAddr(t *testing.T) {
	cfg := Config{}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client.serverAddr != "localhost:50051" {
		t.Errorf("Expected default server address localhost:50051, got %s", client.serverAddr)
	}
}

func TestNewClientWithCustomConfig(t *testing.T) {
	cfg := Config{
		ServerAddr: "example.com:9999",
		Timeout:    60 * time.Second,
		MaxRetries: 5,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client.serverAddr != "example.com:9999" {
		t.Errorf("Expected server address example.com:9999, got %s", client.serverAddr)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ServerAddr != "localhost:50051" {
		t.Errorf("Expected default server address localhost:50051, got %s", cfg.ServerAddr)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", cfg.Timeout)
	}

	if cfg.MaxRetries != 3 {
		t.Errorf("Expected default max retries 3, got %d", cfg.MaxRetries)
	}

	if !cfg.Enabled {
		t.Error("Expected mnemosyne to be enabled by default")
	}
}

func TestConfigFromEnv(t *testing.T) {
	// Save original environment
	origAddr := os.Getenv("MNEMOSYNE_ADDR")
	origTimeout := os.Getenv("MNEMOSYNE_TIMEOUT")
	origRetries := os.Getenv("MNEMOSYNE_MAX_RETRIES")
	origEnabled := os.Getenv("MNEMOSYNE_ENABLED")

	// Restore after test
	defer func() {
		os.Setenv("MNEMOSYNE_ADDR", origAddr)
		os.Setenv("MNEMOSYNE_TIMEOUT", origTimeout)
		os.Setenv("MNEMOSYNE_MAX_RETRIES", origRetries)
		os.Setenv("MNEMOSYNE_ENABLED", origEnabled)
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(*testing.T, Config)
	}{
		{
			name:    "no environment variables",
			envVars: map[string]string{},
			validate: func(t *testing.T, c Config) {
				if c.ServerAddr != "localhost:50051" {
					t.Errorf("Expected default addr, got %s", c.ServerAddr)
				}
				if !c.Enabled {
					t.Error("Expected enabled=true by default")
				}
			},
		},
		{
			name: "custom server address",
			envVars: map[string]string{
				"MNEMOSYNE_ADDR": "example.com:8080",
			},
			validate: func(t *testing.T, c Config) {
				if c.ServerAddr != "example.com:8080" {
					t.Errorf("Expected addr example.com:8080, got %s", c.ServerAddr)
				}
			},
		},
		{
			name: "custom timeout",
			envVars: map[string]string{
				"MNEMOSYNE_TIMEOUT": "60",
			},
			validate: func(t *testing.T, c Config) {
				if c.Timeout != 60*time.Second {
					t.Errorf("Expected timeout 60s, got %v", c.Timeout)
				}
			},
		},
		{
			name: "custom max retries",
			envVars: map[string]string{
				"MNEMOSYNE_MAX_RETRIES": "5",
			},
			validate: func(t *testing.T, c Config) {
				if c.MaxRetries != 5 {
					t.Errorf("Expected max retries 5, got %d", c.MaxRetries)
				}
			},
		},
		{
			name: "disable mnemosyne",
			envVars: map[string]string{
				"MNEMOSYNE_ENABLED": "false",
			},
			validate: func(t *testing.T, c Config) {
				if c.Enabled {
					t.Error("Expected mnemosyne to be disabled")
				}
			},
		},
		{
			name: "enable with 1",
			envVars: map[string]string{
				"MNEMOSYNE_ENABLED": "1",
			},
			validate: func(t *testing.T, c Config) {
				if !c.Enabled {
					t.Error("Expected mnemosyne to be enabled")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv("MNEMOSYNE_ADDR")
			os.Unsetenv("MNEMOSYNE_TIMEOUT")
			os.Unsetenv("MNEMOSYNE_MAX_RETRIES")
			os.Unsetenv("MNEMOSYNE_ENABLED")

			// Set test environment
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Get config from environment
			config := ConfigFromEnv()

			// Validate
			tt.validate(t, config)
		})
	}
}

// --- Connection Tests ---

func TestIsConnectedInitially(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client.IsConnected() {
		t.Error("Expected client to not be connected initially")
	}
}

func TestDisconnectWhenNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Should not error when disconnecting an already disconnected client
	err = client.Disconnect()
	if err != nil {
		t.Errorf("Disconnect should not error when not connected: %v", err)
	}
}

func TestCloseAlias(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Close should be an alias for Disconnect
	err = client.Close()
	if err != nil {
		t.Errorf("Close should not error when not connected: %v", err)
	}
}

// --- Error Handling Tests ---

func TestOperationsWhenNotConnected(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	ctx := context.Background()

	// HealthCheck should fail when not connected
	_, err = client.HealthCheck(ctx)
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}

	// GetStats should fail when not connected
	_, err = client.GetStats(ctx, nil)
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}

	// GetMemory should fail when not connected
	_, err = client.GetMemory(ctx, "test-id")
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}

	// DeleteMemory should fail when not connected
	err = client.DeleteMemory(ctx, "test-id")
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}

	// ListMemories should fail when not connected
	_, err = client.ListMemories(ctx, ListMemoriesOptions{})
	if err != ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
}

// --- Namespace Helper Tests ---

func TestGlobalNamespace(t *testing.T) {
	ns := GlobalNamespace()
	if ns == nil {
		t.Fatal("Expected non-nil namespace")
	}

	if ns.GetGlobal() == nil {
		t.Error("Expected global namespace type")
	}
}

func TestProjectNamespace(t *testing.T) {
	ns := ProjectNamespace("myproject")
	if ns == nil {
		t.Fatal("Expected non-nil namespace")
	}

	project := ns.GetProject()
	if project == nil {
		t.Fatal("Expected project namespace type")
	}

	if project.Name != "myproject" {
		t.Errorf("Expected project name 'myproject', got '%s'", project.Name)
	}
}

func TestSessionNamespace(t *testing.T) {
	ns := SessionNamespace("myproject", "session-123")
	if ns == nil {
		t.Fatal("Expected non-nil namespace")
	}

	session := ns.GetSession()
	if session == nil {
		t.Fatal("Expected session namespace type")
	}

	if session.Project != "myproject" {
		t.Errorf("Expected project 'myproject', got '%s'", session.Project)
	}

	if session.SessionId != "session-123" {
		t.Errorf("Expected session ID 'session-123', got '%s'", session.SessionId)
	}
}

// --- Validation Tests ---

func TestStoreMemoryValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Simulate connection (though it will fail to actually connect)
	client.connected = true

	ctx := context.Background()

	// Missing content should fail
	_, err = client.StoreMemory(ctx, StoreMemoryOptions{
		Namespace: GlobalNamespace(),
	})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for missing content, got: %v", err)
	}

	// Missing namespace should fail
	_, err = client.StoreMemory(ctx, StoreMemoryOptions{
		Content: "test content",
	})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for missing namespace, got: %v", err)
	}

	// Reset connection state
	client.connected = false
}

func TestGetMemoryValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	ctx := context.Background()

	// Empty memory ID should fail
	_, err = client.GetMemory(ctx, "")
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty memory ID, got: %v", err)
	}

	client.connected = false
}

func TestDeleteMemoryValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	ctx := context.Background()

	// Empty memory ID should fail
	err = client.DeleteMemory(ctx, "")
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty memory ID, got: %v", err)
	}

	client.connected = false
}

func TestUpdateMemoryValidation(t *testing.T) {
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	client.connected = true
	ctx := context.Background()

	// Empty memory ID should fail
	_, err = client.UpdateMemory(ctx, UpdateMemoryOptions{})
	if !IsInvalidArgument(err) {
		t.Errorf("Expected invalid argument error for empty memory ID, got: %v", err)
	}

	client.connected = false
}

func TestListMemoriesDefaultMaxResults(t *testing.T) {
	// This test just verifies the options are set correctly
	opts := ListMemoriesOptions{
		Namespace: ProjectNamespace("test"),
	}

	if opts.MaxResults != 0 {
		t.Errorf("Expected MaxResults to be 0 (unset), got %d", opts.MaxResults)
	}

	// The actual default of 100 is applied in the ListMemories method
}

// --- New Test Coverage Improvement Tests ---

// TestClientConnectRetry tests connection retry with exponential backoff.
func TestClientConnectRetry(t *testing.T) {
	// Start a test server that fails the first N attempts
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer server.Stop()

	// Make the server fail the first 2 health checks
	server.memory.failAfter = 2

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    5 * time.Second,
		MaxRetries: 3,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	// Connect should succeed after retries
	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed after retries: %v", err)
	}

	if !client.IsConnected() {
		t.Error("Expected client to be connected after retry")
	}
}

// TestClientConnectTimeout tests connection timeout handling.
func TestClientConnectTimeout(t *testing.T) {
	// Use an address that will timeout (blackhole)
	cfg := Config{
		ServerAddr: "192.0.2.1:12345", // TEST-NET-1 address (should timeout)
		Timeout:    1 * time.Second,
		MaxRetries: 0,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	// Connect should timeout
	start := time.Now()
	err = client.Connect()
	duration := time.Since(start)

	if err == nil {
		t.Error("Expected Connect to fail with timeout")
	}

	// Should timeout within reasonable time (add buffer for test overhead)
	if duration > 15*time.Second {
		t.Errorf("Connect took too long: %v (expected ~10s)", duration)
	}

	if client.IsConnected() {
		t.Error("Expected client to not be connected after timeout")
	}
}

// TestClientDisconnectCleanup tests resource cleanup on disconnect.
func TestClientDisconnectCleanup(t *testing.T) {
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer server.Stop()

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    5 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !client.IsConnected() {
		t.Fatal("Expected client to be connected")
	}

	// Disconnect and verify cleanup
	err = client.Disconnect()
	if err != nil {
		t.Errorf("Disconnect failed: %v", err)
	}

	if client.IsConnected() {
		t.Error("Expected client to not be connected after disconnect")
	}

	// Verify connection is nil
	if client.conn != nil {
		t.Error("Expected connection to be nil after disconnect")
	}

	if client.memoryClient != nil {
		t.Error("Expected memoryClient to be nil after disconnect")
	}

	if client.healthClient != nil {
		t.Error("Expected healthClient to be nil after disconnect")
	}

	// Multiple disconnects should be safe
	err = client.Disconnect()
	if err != nil {
		t.Errorf("Second disconnect should not error: %v", err)
	}
}

// TestClientConcurrentRequests tests multiple concurrent RPC calls.
func TestClientConcurrentRequests(t *testing.T) {
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer server.Stop()

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    10 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Store some test memories first
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_, err := client.StoreMemory(ctx, StoreMemoryOptions{
			Content:   fmt.Sprintf("test memory %d", i),
			Namespace: GlobalNamespace(),
		})
		if err != nil {
			t.Fatalf("StoreMemory failed: %v", err)
		}
	}

	// Perform concurrent recalls
	const numConcurrent = 10
	var wg sync.WaitGroup
	errors := make(chan error, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			_, err := client.Recall(ctx, RecallOptions{
				Query:     fmt.Sprintf("test %d", id),
				Namespace: GlobalNamespace(),
			})
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent request failed: %v", err)
	}

	// Verify request count
	if server.memory.requestCount != numConcurrent+5 {
		t.Errorf("Expected %d requests, got %d", numConcurrent+5, server.memory.requestCount)
	}
}

// TestClientContextCancellation tests request cancellation via context.
func TestClientContextCancellation(t *testing.T) {
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer server.Stop()

	// Add delay to store operation
	server.memory.storeDelay = 2 * time.Second

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    10 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// This should be canceled
	_, err = client.StoreMemory(ctx, StoreMemoryOptions{
		Content:   "test memory",
		Namespace: GlobalNamespace(),
	})

	if err == nil {
		t.Error("Expected StoreMemory to fail with context cancellation")
	}

	if !IsCanceled(err) && !IsDeadlineExceeded(err) {
		t.Errorf("Expected cancellation or deadline error, got: %v", err)
	}
}

// TestClientReconnectAfterFailure tests automatic reconnection after failure.
func TestClientReconnectAfterFailure(t *testing.T) {
	server, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}

	cfg := Config{
		ServerAddr: server.address,
		Timeout:    5 * time.Second,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()

	err = client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	ctx := context.Background()

	// Verify initial connection works
	_, err = client.HealthCheck(ctx)
	if err != nil {
		t.Fatalf("Initial health check failed: %v", err)
	}

	// Stop the server to simulate failure
	server.Stop()

	// Wait a bit for connection to fail
	time.Sleep(100 * time.Millisecond)

	// Operations should fail now
	_, err = client.HealthCheck(ctx)
	if err == nil {
		t.Error("Expected health check to fail after server stopped")
	}

	// Restart server
	server, err = newTestServer()
	if err != nil {
		t.Fatalf("Failed to restart test server: %v", err)
	}
	defer server.Stop()

	// Update client with new address (in real scenario, address would be same)
	// For testing, we need to reconnect to new port
	client.serverAddr = server.address

	// Reconnect
	err = client.Connect()
	if err != nil {
		t.Fatalf("Reconnect failed: %v", err)
	}

	// Health check should work again
	_, err = client.HealthCheck(ctx)
	if err != nil {
		t.Errorf("Health check after reconnect failed: %v", err)
	}
}
