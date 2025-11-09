// Package orchestrate provides session management for orchestration workflows.
package orchestrate

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Session represents an orchestration session with persistence and history tracking.
type Session struct {
	state      *SessionState
	historyDir string
	mu         sync.RWMutex
}

// NewSession creates a new session from a work plan.
func NewSession(plan *WorkPlan) *Session {
	id := generateSessionID()
	state := NewSessionState(id, plan)

	return &Session{
		state:      state,
		historyDir: "", // Will be set lazily on first use
	}
}

// LoadSession loads an existing session from disk by ID.
func LoadSession(id string) (*Session, error) {
	historyDir, err := ensureSessionDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get session directory: %w", err)
	}

	filePath := filepath.Join(historyDir, id+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	// Parse session state
	var state SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session state: %w", err)
	}

	return &Session{
		state:      &state,
		historyDir: historyDir,
	}, nil
}

// Save persists the session state to disk atomically.
func (s *Session) Save() error {
	s.mu.RLock()

	historyDir := s.historyDir
	s.mu.RUnlock()

	// If no historyDir is set, use the default
	if historyDir == "" {
		var err error
		historyDir, err = ensureSessionDir()
		if err != nil {
			return fmt.Errorf("failed to ensure session directory: %w", err)
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Marshal to JSON with versioning
	type sessionFile struct {
		Version int            `json:"version"`
		Session *SessionState  `json:"session"`
	}

	stateToSave := sessionFile{
		Version: 1,
		Session: s.state,
	}

	data, err := json.MarshalIndent(stateToSave, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session state: %w", err)
	}

	// Write to temp file
	tmpFile := filepath.Join(historyDir, s.state.ID+".tmp")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary session file: %w", err)
	}

	// Atomic rename to final location
	finalFile := filepath.Join(historyDir, s.state.ID+".json")
	if err := os.Rename(tmpFile, finalFile); err != nil {
		// Clean up temp file if rename fails
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to rename session file: %w", err)
	}

	return nil
}

// UpdateProgress updates the session state based on an agent event and saves.
func (s *Session) UpdateProgress(event *AgentEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Update session state based on event
	if err := s.state.UpdateProgress(event); err != nil {
		return err
	}

	// Auto-save after updates
	if err := s.saveUnlocked(); err != nil {
		return fmt.Errorf("failed to save after progress update: %w", err)
	}

	return nil
}

// saveUnlocked is the internal save method that assumes the lock is already held.
func (s *Session) saveUnlocked() error {
	historyDir := s.historyDir

	// If no historyDir is set, use the default
	if historyDir == "" {
		var err error
		historyDir, err = ensureSessionDir()
		if err != nil {
			return fmt.Errorf("failed to ensure session directory: %w", err)
		}
	}

	type sessionFile struct {
		Version int            `json:"version"`
		Session *SessionState  `json:"session"`
	}

	stateToSave := sessionFile{
		Version: 1,
		Session: s.state,
	}

	data, err := json.MarshalIndent(stateToSave, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session state: %w", err)
	}

	tmpFile := filepath.Join(historyDir, s.state.ID+".tmp")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary session file: %w", err)
	}

	finalFile := filepath.Join(historyDir, s.state.ID+".json")
	if err := os.Rename(tmpFile, finalFile); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to rename session file: %w", err)
	}

	return nil
}

// GetState returns a copy of the current session state.
func (s *Session) GetState() *SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a shallow copy to prevent external mutations
	// Note: For deep copy, would need to marshal/unmarshal
	stateCopy := *s.state
	return &stateCopy
}

// SetStatus updates the session status.
func (s *Session) SetStatus(status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate status
	validStatuses := map[string]bool{
		"running":   true,
		"paused":    true,
		"completed": true,
		"failed":    true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	s.state.Status = status

	// Update end time if session is finishing
	if status == "completed" || status == "failed" || status == "cancelled" {
		now := time.Now()
		s.state.EndTime = &now
	}

	// Save after status change
	if err := s.saveUnlocked(); err != nil {
		return fmt.Errorf("failed to save after status change: %w", err)
	}

	return nil
}

// SessionSummary represents a brief summary of a session from history.
type SessionSummary struct {
	ID             string    `json:"id"`
	StartTime      time.Time `json:"startTime"`
	EndTime        *time.Time
	Status         string    `json:"status"`
	CompletedTasks int       `json:"completedTasks"`
	FailedTasks    int       `json:"failedTasks"`
	TotalTasks     int       `json:"totalTasks"`
	Progress       float64   `json:"progress"`
}

// HistoryOptions controls filtering and pagination for session history.
type HistoryOptions struct {
	Status  string // Filter by status (empty = no filter)
	Limit   int    // Maximum results (0 = no limit)
	Offset  int    // Skip first N results
	Reverse bool   // Reverse sort order (oldest first)
}

// History returns a list of all past sessions with optional filtering.
func History(opts *HistoryOptions) ([]*SessionSummary, error) {
	if opts == nil {
		opts = &HistoryOptions{}
	}

	historyDir, err := ensureSessionDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get session directory: %w", err)
	}

	entries, err := os.ReadDir(historyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read session directory: %w", err)
	}

	var summaries []*SessionSummary

	for _, entry := range entries {
		// Only process .json files (skip .tmp and other files)
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(historyDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			// Skip files that can't be read (e.g., permission issues)
			continue
		}

		// Try to unmarshal session
		var sessionFile struct {
			Version int            `json:"version"`
			Session *SessionState  `json:"session"`
		}

		if err := json.Unmarshal(data, &sessionFile); err != nil {
			// Skip corrupted files
			continue
		}

		if sessionFile.Session == nil {
			continue
		}

		session := sessionFile.Session

		// Apply status filter if specified
		if opts.Status != "" && session.Status != opts.Status {
			continue
		}

		summary := &SessionSummary{
			ID:             session.ID,
			StartTime:      session.StartTime,
			EndTime:        session.EndTime,
			Status:         session.Status,
			CompletedTasks: session.CompletedTasks,
			FailedTasks:    session.FailedTasks,
			TotalTasks:     session.TotalTasks,
			Progress:       session.Progress(),
		}

		summaries = append(summaries, summary)
	}

	// Sort by start time (most recent first)
	sort.Slice(summaries, func(i, j int) bool {
		if opts.Reverse {
			return summaries[i].StartTime.Before(summaries[j].StartTime)
		}
		return summaries[i].StartTime.After(summaries[j].StartTime)
	})

	// Apply pagination
	if opts.Offset > 0 {
		if opts.Offset > len(summaries) {
			return []*SessionSummary{}, nil
		}
		summaries = summaries[opts.Offset:]
	}

	if opts.Limit > 0 && len(summaries) > opts.Limit {
		summaries = summaries[:opts.Limit]
	}

	return summaries, nil
}

// Delete removes a session from disk.
func Delete(id string) error {
	historyDir, err := ensureSessionDir()
	if err != nil {
		return fmt.Errorf("failed to get session directory: %w", err)
	}

	filePath := filepath.Join(historyDir, id+".json")
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("session not found: %s", id)
		}
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// --- Helper Functions ---

// generateSessionID creates a unique session ID based on timestamp with nanosecond precision.
func generateSessionID() string {
	now := time.Now()
	return fmt.Sprintf("session-%s-%d", now.Format("20060102-150405"), now.Nanosecond())
}

// ensureSessionDir creates the session directory if it doesn't exist.
func ensureSessionDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	dir := filepath.Join(home, ".pedantic_raven", "sessions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create session directory: %w", err)
	}

	return dir, nil
}

// ListSessionDirContents returns raw file information for testing purposes.
func ListSessionDirContents() ([]fs.DirEntry, error) {
	historyDir, err := ensureSessionDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(historyDir)
	if err != nil {
		return nil, err
	}

	return entries, nil
}
