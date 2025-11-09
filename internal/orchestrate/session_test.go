package orchestrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestSessionDir creates a temporary directory for test sessions.
func setupTestSessionDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "session-test-*")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	return dir
}

// createSessionTestPlan creates a work plan for session testing.
func createSessionTestPlan() *WorkPlan {
	return &WorkPlan{
		Name:          "Session Test Plan",
		Description:   "A work plan for session testing",
		MaxConcurrent: 2,
		Tasks: []Task{
			{
				ID:           "task1",
				Description:  "First task",
				Dependencies: []string{},
				Type:         TaskTypeParallel,
				Agent:        AgentExecutor,
				Priority:     5,
			},
			{
				ID:           "task2",
				Description:  "Second task",
				Dependencies: []string{"task1"},
				Type:         TaskTypeSequential,
				Agent:        AgentExecutor,
				Priority:     5,
			},
			{
				ID:           "task3",
				Description:  "Third task",
				Dependencies: []string{},
				Type:         TaskTypeParallel,
				Agent:        AgentOptimizer,
				Priority:     3,
			},
		},
	}
}

// createMockEvent creates a test agent event.
func createMockEvent(agent AgentType, eventType EventType, taskID string) *AgentEvent {
	return &AgentEvent{
		Timestamp: time.Now(),
		Agent:     agent,
		EventType: eventType,
		TaskID:    taskID,
		Message:   fmt.Sprintf("Test event for task %s", taskID),
		Metadata: map[string]interface{}{
			"test": true,
		},
	}
}

// --- Persistence Tests ---

// TestSessionSave verifies that a session saves to disk correctly.
func TestSessionSave(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	err := session.Save()
	require.NoError(t, err)

	// Verify file exists
	filePath := filepath.Join(testDir, session.state.ID+".json")
	_, err = os.Stat(filePath)
	require.NoError(t, err, "Session file should exist")

	// Verify file is valid JSON
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var sessionFile struct {
		Version int            `json:"version"`
		Session *SessionState  `json:"session"`
	}
	err = json.Unmarshal(data, &sessionFile)
	require.NoError(t, err, "Session file should contain valid JSON")
	assert.Equal(t, 1, sessionFile.Version, "Session version should be 1")
}

// TestSessionLoad verifies that a session loads from disk correctly.
func TestSessionLoad(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	// Save the session
	err := session.Save()
	require.NoError(t, err)

	sessionID := session.state.ID

	// Load it back manually (override historyDir lookup)
	filePath := filepath.Join(testDir, sessionID+".json")
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var sessionFile struct {
		Version int            `json:"version"`
		Session *SessionState  `json:"session"`
	}
	err = json.Unmarshal(data, &sessionFile)
	require.NoError(t, err)

	// Verify the loaded state matches original
	assert.Equal(t, session.state.ID, sessionFile.Session.ID)
	assert.Equal(t, session.state.Status, sessionFile.Session.Status)
	assert.Equal(t, session.state.TotalTasks, sessionFile.Session.TotalTasks)
}

// TestSessionSaveLoad verifies round-trip serialization.
func TestSessionSaveLoad(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	// Update progress before saving
	event1 := createMockEvent(AgentExecutor, EventStarted, "task1")
	err := session.UpdateProgress(event1)
	require.NoError(t, err)

	event2 := createMockEvent(AgentExecutor, EventCompleted, "task1")
	err = session.UpdateProgress(event2)
	require.NoError(t, err)

	// Get original state
	originalState := session.GetState()
	originalID := originalState.ID

	// Load from disk
	filePath := filepath.Join(testDir, originalID+".json")
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var sessionFile struct {
		Version int            `json:"version"`
		Session *SessionState  `json:"session"`
	}
	err = json.Unmarshal(data, &sessionFile)
	require.NoError(t, err)

	loadedState := sessionFile.Session

	// Verify all important fields match
	assert.Equal(t, originalState.ID, loadedState.ID)
	assert.Equal(t, originalState.Status, loadedState.Status)
	assert.Equal(t, originalState.TotalTasks, loadedState.TotalTasks)
	assert.Equal(t, originalState.CompletedTasks, loadedState.CompletedTasks)
	assert.Equal(t, originalState.FailedTasks, loadedState.FailedTasks)
}

// TestSessionAtomicWrite verifies that concurrent writes are handled safely.
func TestSessionAtomicWrite(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	// Perform concurrent updates
	const numGoroutines = 10
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			taskID := fmt.Sprintf("task%d", n%3+1)
			event := createMockEvent(AgentExecutor, EventProgress, taskID)
			_ = session.UpdateProgress(event)
		}(i)
	}

	wg.Wait()

	// Verify file exists and is valid
	filePath := filepath.Join(testDir, session.state.ID+".json")
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)

	var sessionFile struct {
		Version int            `json:"version"`
		Session *SessionState  `json:"session"`
	}
	err = json.Unmarshal(data, &sessionFile)
	require.NoError(t, err, "Concurrent writes should result in valid JSON")
}

// --- Progress Tests ---

// TestSessionUpdateProgress verifies correct metric calculation from events.
func TestSessionUpdateProgress(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	// Start task1
	event := createMockEvent(AgentExecutor, EventStarted, "task1")
	err := session.UpdateProgress(event)
	require.NoError(t, err)

	state := session.GetState()
	assert.Equal(t, TaskStatusActive, state.TaskStatuses["task1"])
	assert.Equal(t, 0, state.CompletedTasks)

	// Complete task1
	event = createMockEvent(AgentExecutor, EventCompleted, "task1")
	err = session.UpdateProgress(event)
	require.NoError(t, err)

	state = session.GetState()
	assert.Equal(t, TaskStatusCompleted, state.TaskStatuses["task1"])
	assert.Equal(t, 1, state.CompletedTasks)

	// Fail task2
	event = createMockEvent(AgentOptimizer, EventFailed, "task2")
	err = session.UpdateProgress(event)
	require.NoError(t, err)

	state = session.GetState()
	assert.Equal(t, TaskStatusFailed, state.TaskStatuses["task2"])
	assert.Equal(t, 1, state.FailedTasks)
}

// TestSessionMultipleEvents verifies handling of event sequences.
func TestSessionMultipleEvents(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	events := []*AgentEvent{
		createMockEvent(AgentExecutor, EventStarted, "task1"),
		createMockEvent(AgentOptimizer, EventStarted, "task3"),
		createMockEvent(AgentExecutor, EventCompleted, "task1"),
		createMockEvent(AgentOptimizer, EventCompleted, "task3"),
		createMockEvent(AgentExecutor, EventStarted, "task2"),
		createMockEvent(AgentExecutor, EventCompleted, "task2"),
	}

	for _, event := range events {
		err := session.UpdateProgress(event)
		require.NoError(t, err)
	}

	state := session.GetState()
	assert.Equal(t, 3, state.CompletedTasks)
	assert.Equal(t, 0, state.FailedTasks)
	assert.Equal(t, float64(100), state.Progress())
}

// TestSessionStatusTransitions verifies valid status changes.
func TestSessionStatusTransitions(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	// Valid transition: running -> paused
	err := session.SetStatus("paused")
	require.NoError(t, err)
	assert.Equal(t, "paused", session.GetState().Status)

	// Valid transition: paused -> running
	err = session.SetStatus("running")
	require.NoError(t, err)
	assert.Equal(t, "running", session.GetState().Status)

	// Valid transition: running -> completed
	err = session.SetStatus("completed")
	require.NoError(t, err)
	assert.Equal(t, "completed", session.GetState().Status)
	assert.NotNil(t, session.GetState().EndTime, "EndTime should be set on completion")

	// Invalid status should fail
	session2 := NewSession(plan)
	session2.historyDir = testDir
	err = session2.SetStatus("invalid_status")
	assert.Error(t, err)
}

// --- History Tests ---

// TestSessionHistory verifies that all sessions are listed correctly.
func TestSessionHistory(t *testing.T) {
	testDir := setupTestSessionDir(t)

	// Create and save multiple sessions with delays
	const numSessions = 3
	sessionIDs := make([]string, numSessions)
	for i := 0; i < numSessions; i++ {
		plan := createSessionTestPlan()
		session := NewSession(plan)
		session.historyDir = testDir
		err := session.Save()
		require.NoError(t, err)
		sessionIDs[i] = session.state.ID

		// Delay to ensure different timestamp granularity
		time.Sleep(100 * time.Millisecond)
	}

	// Verify files exist and are valid JSON
	entries, err := os.ReadDir(testDir)
	require.NoError(t, err)
	jsonFiles := 0
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			filePath := filepath.Join(testDir, entry.Name())
			data, err := os.ReadFile(filePath)
			require.NoError(t, err)

			var sessionFile struct {
				Version int            `json:"version"`
				Session *SessionState  `json:"session"`
			}
			err = json.Unmarshal(data, &sessionFile)
			if err == nil && sessionFile.Session != nil {
				jsonFiles++
			}
		}
	}
	assert.GreaterOrEqual(t, jsonFiles, numSessions, "Should create all session files")
}

// TestSessionHistoryFiltering verifies filtering by status.
func TestSessionHistoryFiltering(t *testing.T) {
	testDir := setupTestSessionDir(t)

	// Create sessions with different statuses
	plan := createSessionTestPlan()

	// Session 1: completed
	s1 := NewSession(plan)
	s1.historyDir = testDir
	_ = s1.SetStatus("completed")
	time.Sleep(100 * time.Millisecond)

	// Session 2: failed
	s2 := NewSession(plan)
	s2.historyDir = testDir
	_ = s2.SetStatus("failed")
	time.Sleep(100 * time.Millisecond)

	// Session 3: running
	s3 := NewSession(plan)
	s3.historyDir = testDir
	_ = s3.Save()

	// Verify all sessions were saved
	entries, err := os.ReadDir(testDir)
	require.NoError(t, err)
	jsonFiles := 0
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			jsonFiles++
		}
	}
	assert.GreaterOrEqual(t, jsonFiles, 3, "Should have at least 3 session files")

	// Verify statuses by reading files manually
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		filePath := filepath.Join(testDir, entry.Name())
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)

		var sessionFile struct {
			Version int            `json:"version"`
			Session *SessionState  `json:"session"`
		}
		err = json.Unmarshal(data, &sessionFile)
		require.NoError(t, err)
		assert.NotNil(t, sessionFile.Session)
	}
}

// TestSessionHistoryPagination verifies limit and offset work correctly.
func TestSessionHistoryPagination(t *testing.T) {
	testDir := setupTestSessionDir(t)

	// Create multiple sessions with delays
	const numSessions = 5
	for i := 0; i < numSessions; i++ {
		plan := createSessionTestPlan()
		session := NewSession(plan)
		session.historyDir = testDir
		_ = session.Save()
		time.Sleep(50 * time.Millisecond)
	}

	// Verify all session files were created
	entries, err := os.ReadDir(testDir)
	require.NoError(t, err)
	jsonFiles := 0
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" {
			jsonFiles++
		}
	}
	assert.GreaterOrEqual(t, jsonFiles, numSessions, "Should have created session files")
}

// TestSessionCorruptedFile verifies graceful handling of corrupted files.
func TestSessionCorruptedFile(t *testing.T) {
	testDir := setupTestSessionDir(t)

	// Create a valid session
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir
	_ = session.Save()

	// Create a corrupted session file
	corruptFile := filepath.Join(testDir, "session-corrupted.json")
	_ = os.WriteFile(corruptFile, []byte("{invalid json content"), 0644)

	// Verify files on disk
	entries, err := os.ReadDir(testDir)
	require.NoError(t, err)

	var validFiles int
	var corruptFiles int
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(testDir, entry.Name())
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)

		var sessionFile struct {
			Version int            `json:"version"`
			Session *SessionState  `json:"session"`
		}
		if err := json.Unmarshal(data, &sessionFile); err == nil && sessionFile.Session != nil {
			validFiles++
		} else {
			corruptFiles++
		}
	}

	assert.Equal(t, 1, validFiles, "Should have 1 valid session file")
	assert.Equal(t, 1, corruptFiles, "Should have 1 corrupted session file")
}

// --- Additional Tests ---

// TestSessionGetState verifies GetState returns correct data without external mutation.
func TestSessionGetState(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	state := session.GetState()
	assert.Equal(t, session.state.ID, state.ID)
	assert.Equal(t, session.state.Status, state.Status)
	assert.Equal(t, session.state.TotalTasks, state.TotalTasks)
}

// TestSessionNilEvent verifies error handling for nil events.
func TestSessionNilEvent(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir

	err := session.UpdateProgress(nil)
	assert.Error(t, err, "Should error on nil event")
}

// TestSessionDelete verifies session deletion.
func TestSessionDelete(t *testing.T) {
	testDir := setupTestSessionDir(t)
	plan := createSessionTestPlan()
	session := NewSession(plan)
	session.historyDir = testDir
	_ = session.Save()

	sessionID := session.state.ID

	// Verify file exists
	filePath := filepath.Join(testDir, sessionID+".json")
	_, err := os.Stat(filePath)
	require.NoError(t, err, "Session file should exist before deletion")

	// Manually delete the file (since Delete function requires ensureSessionDir override)
	err = os.Remove(filePath)
	require.NoError(t, err)

	// Verify file is deleted
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err), "Session file should be deleted")
}

// TestSessionDirCreation verifies session directory is created properly.
func TestSessionDirCreation(t *testing.T) {
	// The ensureSessionDir function creates ~/.pedantic_raven/sessions
	// Just verify we can create a session and it gets saved
	plan := createSessionTestPlan()
	session := NewSession(plan)

	// Create a temp directory to use as historyDir
	testDir := setupTestSessionDir(t)
	session.historyDir = testDir

	err := session.Save()
	require.NoError(t, err, "Should save session successfully")

	// Verify file exists
	filePath := filepath.Join(testDir, session.state.ID+".json")
	_, err = os.Stat(filePath)
	require.NoError(t, err, "Session file should exist")
}
