package mnemosyne

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

// ConnectionStatus represents the current state of the connection.
type ConnectionStatus int

const (
	StatusDisconnected ConnectionStatus = iota
	StatusConnecting
	StatusConnected
	StatusReconnecting
	StatusFailed
)

// String returns the string representation of the connection status.
func (s ConnectionStatus) String() string {
	switch s {
	case StatusDisconnected:
		return "disconnected"
	case StatusConnecting:
		return "connecting"
	case StatusConnected:
		return "connected"
	case StatusReconnecting:
		return "reconnecting"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// RetryPolicy defines the retry behavior for connection failures.
type RetryPolicy struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// DefaultRetryPolicy returns a sensible default retry policy.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:    5,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
	}
}

// ConnectionConfig holds configuration for connection management.
type ConnectionConfig struct {
	Host        string
	Port        int
	UseTLS      bool
	Timeout     time.Duration
	RetryPolicy RetryPolicy
}

// Validate checks if the configuration is valid.
func (cfg *ConnectionConfig) Validate() error {
	if cfg.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", cfg.Port)
	}
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got %v", cfg.Timeout)
	}
	if cfg.RetryPolicy.MaxAttempts < 0 {
		return fmt.Errorf("max retry attempts must be non-negative, got %d", cfg.RetryPolicy.MaxAttempts)
	}
	if cfg.RetryPolicy.InitialBackoff <= 0 {
		return fmt.Errorf("initial backoff must be positive, got %v", cfg.RetryPolicy.InitialBackoff)
	}
	if cfg.RetryPolicy.MaxBackoff < cfg.RetryPolicy.InitialBackoff {
		return fmt.Errorf("max backoff must be >= initial backoff")
	}
	return nil
}

// ConnectionManager manages persistent connection to mnemosyne-rpc server.
type ConnectionManager struct {
	client      *Client
	config      *ConnectionConfig
	status      ConnectionStatus
	mu          sync.RWMutex
	healthCheck *time.Ticker
	stopHealth  chan struct{}
	retryCount  int
}

// NewConnectionManager creates a new connection manager.
func NewConnectionManager(config *ConnectionConfig) (*ConnectionManager, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	cm := &ConnectionManager{
		config: config,
		status: StatusDisconnected,
	}

	return cm, nil
}

// Connect establishes a connection to the mnemosyne server.
func (cm *ConnectionManager) Connect() error {
	cm.mu.Lock()

	// Prevent concurrent connections
	if cm.status == StatusConnecting || cm.status == StatusConnected {
		cm.mu.Unlock()
		return nil
	}

	cm.status = StatusConnecting
	cm.retryCount = 0
	cm.mu.Unlock()

	// Create client
	serverAddr := fmt.Sprintf("%s:%d", cm.config.Host, cm.config.Port)
	clientCfg := Config{
		ServerAddr: serverAddr,
		Timeout:    cm.config.Timeout,
	}

	client, err := NewClient(clientCfg)
	if err != nil {
		cm.mu.Lock()
		cm.status = StatusFailed
		cm.mu.Unlock()
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Attempt connection
	if err := client.Connect(); err != nil {
		cm.mu.Lock()
		cm.status = StatusFailed
		cm.mu.Unlock()
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Verify connection with health check
	ctx, cancel := context.WithTimeout(context.Background(), cm.config.Timeout)
	defer cancel()

	if _, err := client.HealthCheck(ctx); err != nil {
		client.Disconnect()
		cm.mu.Lock()
		cm.status = StatusFailed
		cm.mu.Unlock()
		return fmt.Errorf("health check failed: %w", err)
	}

	// Connection successful
	cm.mu.Lock()
	cm.client = client
	cm.status = StatusConnected
	cm.mu.Unlock()

	// Start health check monitoring
	cm.startHealthCheckMonitoring()

	return nil
}

// Disconnect closes the connection to the server.
func (cm *ConnectionManager) Disconnect() error {
	cm.mu.Lock()

	if cm.status == StatusDisconnected {
		cm.mu.Unlock()
		return nil
	}

	// Capture state before releasing lock
	client := cm.client
	hasHealthCheck := cm.healthCheck != nil

	cm.client = nil
	cm.status = StatusDisconnected
	cm.mu.Unlock()

	// Stop health check monitoring (must be done without lock to avoid deadlock)
	if hasHealthCheck {
		cm.stopHealthCheckMonitoring()
	}

	// Disconnect client
	if client != nil {
		return client.Disconnect()
	}

	return nil
}

// Status returns the current connection status.
func (cm *ConnectionManager) Status() ConnectionStatus {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.status
}

// HealthCheck performs a health check on the connection.
func (cm *ConnectionManager) HealthCheck() error {
	cm.mu.RLock()
	client := cm.client
	status := cm.status
	cm.mu.RUnlock()

	if status != StatusConnected && status != StatusReconnecting {
		return fmt.Errorf("not connected (status: %s)", status)
	}

	if client == nil {
		return fmt.Errorf("client is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cm.config.Timeout)
	defer cancel()

	_, err := client.HealthCheck(ctx)
	if err != nil {
		// Health check failed, trigger reconnection
		go cm.attemptReconnect()
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

// Client returns the underlying mnemosyne client.
// Returns nil if not connected.
func (cm *ConnectionManager) Client() *Client {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.client
}

// startHealthCheckMonitoring starts periodic health checks.
func (cm *ConnectionManager) startHealthCheckMonitoring() {
	cm.mu.Lock()
	if cm.healthCheck != nil {
		// Already running
		cm.mu.Unlock()
		return
	}

	cm.healthCheck = time.NewTicker(30 * time.Second)
	cm.stopHealth = make(chan struct{})
	ticker := cm.healthCheck
	stopChan := cm.stopHealth
	cm.mu.Unlock()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := cm.HealthCheck(); err != nil {
					// Health check failed, reconnection will be attempted
					// by HealthCheck() method
				}
			case <-stopChan:
				return
			}
		}
	}()
}

// stopHealthCheckMonitoring stops the health check ticker.
func (cm *ConnectionManager) stopHealthCheckMonitoring() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.healthCheck != nil {
		cm.healthCheck.Stop()
		cm.healthCheck = nil
	}

	if cm.stopHealth != nil {
		close(cm.stopHealth)
		cm.stopHealth = nil
	}
}

// attemptReconnect attempts to reconnect with exponential backoff.
func (cm *ConnectionManager) attemptReconnect() {
	cm.mu.Lock()

	// Check if already reconnecting or disconnected
	if cm.status == StatusReconnecting || cm.status == StatusDisconnected {
		cm.mu.Unlock()
		return
	}

	// Set status to reconnecting
	cm.status = StatusReconnecting

	// Stop health check during reconnection
	if cm.healthCheck != nil {
		cm.mu.Unlock()
		cm.stopHealthCheckMonitoring()
		cm.mu.Lock()
	}

	// Close existing client
	if cm.client != nil {
		cm.client.Disconnect()
		cm.client = nil
	}

	cm.retryCount = 0
	cm.mu.Unlock()

	// Retry with exponential backoff
	for {
		cm.mu.RLock()
		status := cm.status
		retryCount := cm.retryCount
		cm.mu.RUnlock()

		// Stop if status changed (e.g., manual disconnect)
		if status != StatusReconnecting {
			return
		}

		// Check if max retries exceeded
		if cm.config.RetryPolicy.MaxAttempts > 0 && retryCount >= cm.config.RetryPolicy.MaxAttempts {
			cm.mu.Lock()
			cm.status = StatusFailed
			cm.mu.Unlock()
			return
		}

		// Calculate backoff duration
		backoff := cm.calculateBackoff(retryCount)
		time.Sleep(backoff)

		// Attempt connection
		serverAddr := fmt.Sprintf("%s:%d", cm.config.Host, cm.config.Port)
		clientCfg := Config{
			ServerAddr: serverAddr,
			Timeout:    cm.config.Timeout,
		}

		client, err := NewClient(clientCfg)
		if err != nil {
			cm.mu.Lock()
			cm.retryCount++
			cm.mu.Unlock()
			continue
		}

		if err := client.Connect(); err != nil {
			cm.mu.Lock()
			cm.retryCount++
			cm.mu.Unlock()
			continue
		}

		// Verify with health check
		ctx, cancel := context.WithTimeout(context.Background(), cm.config.Timeout)
		_, err = client.HealthCheck(ctx)
		cancel()

		if err != nil {
			client.Disconnect()
			cm.mu.Lock()
			cm.retryCount++
			cm.mu.Unlock()
			continue
		}

		// Reconnection successful
		cm.mu.Lock()
		cm.client = client
		cm.status = StatusConnected
		cm.retryCount = 0
		cm.mu.Unlock()

		// Restart health check monitoring
		cm.startHealthCheckMonitoring()
		return
	}
}

// calculateBackoff calculates the backoff duration for a given attempt.
// Uses exponential backoff: min(initialBackoff * 2^attempt, maxBackoff)
func (cm *ConnectionManager) calculateBackoff(attempt int) time.Duration {
	backoff := float64(cm.config.RetryPolicy.InitialBackoff) * math.Pow(2, float64(attempt))
	maxBackoff := float64(cm.config.RetryPolicy.MaxBackoff)

	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	return time.Duration(backoff)
}
