package mnemosyne

import (
	"context"
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
