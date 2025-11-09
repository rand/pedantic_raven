// Package mnemosyne provides a gRPC client for the mnemosyne RPC server.
package mnemosyne

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Client provides access to mnemosyne RPC services.
type Client struct {
	conn          *grpc.ClientConn
	memoryClient  pb.MemoryServiceClient
	healthClient  pb.HealthServiceClient
	serverAddr    string
	connected     bool
	defaultCtx    context.Context
	defaultCancel context.CancelFunc
}

// Config holds client configuration.
type Config struct {
	// ServerAddr is the mnemosyne server address (host:port)
	ServerAddr string

	// Timeout is the default timeout for operations (0 = no timeout)
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// Enabled determines if mnemosyne integration is enabled
	// If false, applications should fall back to sample/offline mode
	Enabled bool
}

// DefaultConfig returns a default configuration.
func DefaultConfig() Config {
	return Config{
		ServerAddr: "localhost:50051",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		Enabled:    true,
	}
}

// ConfigFromEnv creates a Config from environment variables.
// Falls back to defaults for any missing values.
func ConfigFromEnv() Config {
	config := DefaultConfig()

	// MNEMOSYNE_ENABLED - enable/disable real integration
	if enabled := os.Getenv("MNEMOSYNE_ENABLED"); enabled != "" {
		config.Enabled = enabled == "true" || enabled == "1"
	}

	// MNEMOSYNE_ADDR - server address (host:port)
	if addr := os.Getenv("MNEMOSYNE_ADDR"); addr != "" {
		config.ServerAddr = addr
	}

	// MNEMOSYNE_TIMEOUT - operation timeout in seconds
	if timeoutStr := os.Getenv("MNEMOSYNE_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
			config.Timeout = time.Duration(timeout) * time.Second
		}
	}

	// MNEMOSYNE_MAX_RETRIES - maximum retry attempts
	if retriesStr := os.Getenv("MNEMOSYNE_MAX_RETRIES"); retriesStr != "" {
		if retries, err := strconv.Atoi(retriesStr); err == nil && retries >= 0 {
			config.MaxRetries = retries
		}
	}

	return config
}

// NewClient creates a new mnemosyne client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.ServerAddr == "" {
		cfg.ServerAddr = "localhost:50051"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	client := &Client{
		serverAddr: cfg.ServerAddr,
	}

	// Create default context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	client.defaultCtx = ctx
	client.defaultCancel = cancel

	return client, nil
}

// Connect establishes a connection to the mnemosyne server.
func (c *Client) Connect() error {
	if c.connected {
		return nil // Already connected
	}

	// Set up connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	// Connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, c.serverAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to mnemosyne server at %s: %w", c.serverAddr, err)
	}

	c.conn = conn
	c.memoryClient = pb.NewMemoryServiceClient(conn)
	c.healthClient = pb.NewHealthServiceClient(conn)
	c.connected = true

	return nil
}

// Disconnect closes the connection to the mnemosyne server.
func (c *Client) Disconnect() error {
	if !c.connected {
		return nil // Already disconnected
	}

	if c.defaultCancel != nil {
		c.defaultCancel()
	}

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.memoryClient = nil
		c.healthClient = nil
		c.connected = false
		return err
	}

	c.connected = false
	return nil
}

// IsConnected returns true if the client is connected to the server.
func (c *Client) IsConnected() bool {
	return c.connected && c.conn != nil
}

// HealthCheck performs a basic health check on the server.
func (c *Client) HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	req := &pb.HealthCheckRequest{}
	resp, err := c.healthClient.HealthCheck(ctx, req)
	if err != nil {
		return nil, wrapError(err, "health check")
	}

	return resp, nil
}

// GetStats retrieves server statistics.
func (c *Client) GetStats(ctx context.Context, namespace *pb.Namespace) (*pb.Stats, error) {
	if !c.connected {
		return nil, ErrNotConnected
	}

	req := &pb.GetStatsRequest{
		Namespace: namespace,
	}

	resp, err := c.healthClient.GetStats(ctx, req)
	if err != nil {
		return nil, wrapError(err, "get stats")
	}

	return resp.Stats, nil
}

// Close is an alias for Disconnect for better resource management patterns.
func (c *Client) Close() error {
	return c.Disconnect()
}
