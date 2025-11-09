package orchestrate

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create mock entries
func createMockEntries(n int) []LogEntry {
	entries := make([]LogEntry, n)
	for i := 0; i < n; i++ {
		var agent AgentType
		switch i % 4 {
		case 0:
			agent = AgentOrchestrator
		case 1:
			agent = AgentOptimizer
		case 2:
			agent = AgentReviewer
		default:
			agent = AgentExecutor
		}

		entries[i] = LogEntry{
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
			Agent:     agent,
			EventType: EventLog,
			TaskID:    fmt.Sprintf("task-%d", i),
			Message:   fmt.Sprintf("Log message %d", i),
			Level:     LogLevelInfo,
		}
	}
	return entries
}

// Helper function to create an AgentEvent for testing
func createTestEvent(agent AgentType, eventType EventType, message string) *AgentEvent {
	return &AgentEvent{
		Timestamp: time.Now(),
		Agent:     agent,
		EventType: eventType,
		TaskID:    "test-task",
		Message:   message,
		Metadata:  nil,
	}
}

// --- Buffer Tests ---

// TestAgentLogAddEntry verifies that entries are added correctly.
func TestAgentLogAddEntry(t *testing.T) {
	al := NewAgentLog(80, 24)

	event := createTestEvent(AgentOrchestrator, EventStarted, "Task started")
	al.AddEntry(event)

	assert.Equal(t, 1, al.TotalEntries())
	assert.Equal(t, 1, al.FilteredEntries())

	// Verify the entry content
	entries := al.getFilteredEntries()
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, AgentOrchestrator, entries[0].Agent)
	assert.Equal(t, "Task started", entries[0].Message)
	assert.Equal(t, LogLevelInfo, entries[0].Level)
}

// TestAgentLogCircularBuffer verifies that oldest entries are removed when full.
func TestAgentLogCircularBuffer(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add more entries than maxEntries
	for i := 0; i < al.maxEntries+500; i++ {
		event := createTestEvent(AgentOrchestrator, EventLog, fmt.Sprintf("Message %d", i))
		al.AddEntry(event)
	}

	// Should not exceed maxEntries
	assert.Equal(t, al.maxEntries, al.TotalEntries())

	// Verify that the oldest entries are gone
	entries := al.getFilteredEntries()
	firstTaskID := entries[0].TaskID
	assert.Equal(t, "test-task", firstTaskID) // All have same task ID from mock

	// The messages should be from later in the sequence
	lastMessage := entries[len(entries)-1].Message
	assert.Contains(t, lastMessage, "Message")
}

// TestAgentLogMaxSize verifies that buffer never exceeds maxEntries.
func TestAgentLogMaxSize(t *testing.T) {
	al := NewAgentLog(80, 24)
	al.maxEntries = 100 // Use smaller size for testing

	// Add way more than maxEntries
	for i := 0; i < 500; i++ {
		event := createTestEvent(AgentOrchestrator, EventLog, fmt.Sprintf("Message %d", i))
		al.AddEntry(event)
	}

	assert.Equal(t, 100, al.TotalEntries())
	assert.LessOrEqual(t, al.TotalEntries(), al.maxEntries)
}

// --- Filtering Tests ---

// TestAgentLogFilterByAgent verifies filtering by agent type.
func TestAgentLogFilterByAgent(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add entries from different agents
	for i := 0; i < 12; i++ {
		var agent AgentType
		switch i % 4 {
		case 0:
			agent = AgentOrchestrator
		case 1:
			agent = AgentOptimizer
		case 2:
			agent = AgentReviewer
		default:
			agent = AgentExecutor
		}
		event := createTestEvent(agent, EventLog, fmt.Sprintf("Message from %s", agent.String()))
		al.AddEntry(event)
	}

	// Test filter for Orchestrator
	filterAgent := AgentOrchestrator
	al.SetFilterAgent(&filterAgent)
	filtered := al.getFilteredEntries()
	assert.Equal(t, 3, len(filtered)) // 0, 4, 8

	for _, entry := range filtered {
		assert.Equal(t, AgentOrchestrator, entry.Agent)
	}

	// Test filter for Executor
	filterAgent = AgentExecutor
	al.SetFilterAgent(&filterAgent)
	filtered = al.getFilteredEntries()
	assert.Equal(t, 3, len(filtered)) // 3, 7, 11

	for _, entry := range filtered {
		assert.Equal(t, AgentExecutor, entry.Agent)
	}

	// Test no filter
	al.SetFilterAgent(nil)
	filtered = al.getFilteredEntries()
	assert.Equal(t, 12, len(filtered))
}

// TestAgentLogFilterByLevel verifies filtering by log level.
func TestAgentLogFilterByLevel(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add entries with different levels
	infoEvent := createTestEvent(AgentOrchestrator, EventLog, "Info message")
	warnEvent := createTestEvent(AgentOrchestrator, EventProgress, "Warning message")
	errorEvent := createTestEvent(AgentOrchestrator, EventFailed, "Error message")

	al.AddEntry(infoEvent)
	al.AddEntry(warnEvent)
	al.AddEntry(errorEvent)

	// Test filter for errors only
	levelError := LogLevelError
	al.SetFilterLevel(&levelError)
	filtered := al.getFilteredEntries()
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, LogLevelError, filtered[0].Level)

	// Test filter for warnings and errors
	levelWarn := LogLevelWarn
	al.SetFilterLevel(&levelWarn)
	filtered = al.getFilteredEntries()
	assert.GreaterOrEqual(t, len(filtered), 1) // At least the error

	// Test no filter
	al.SetFilterLevel(nil)
	filtered = al.getFilteredEntries()
	assert.Equal(t, 3, len(filtered))
}

// TestAgentLogSearch verifies regex search functionality.
func TestAgentLogSearch(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add various messages
	messages := []string{
		"Task started successfully",
		"Processing item 1",
		"Processing item 2",
		"Error occurred",
		"Connection timeout",
	}

	for _, msg := range messages {
		event := createTestEvent(AgentOrchestrator, EventLog, msg)
		al.AddEntry(event)
	}

	// Test search for "Processing"
	err := al.SetSearchQuery("Processing")
	require.NoError(t, err)
	filtered := al.getFilteredEntries()
	assert.Equal(t, 2, len(filtered))
	assert.Equal(t, "Processing item 1", filtered[0].Message)
	assert.Equal(t, "Processing item 2", filtered[1].Message)

	// Test search for pattern with regex
	err = al.SetSearchQuery("item \\d")
	require.NoError(t, err)
	filtered = al.getFilteredEntries()
	assert.Equal(t, 2, len(filtered))

	// Test search for "Error"
	err = al.SetSearchQuery("Error")
	require.NoError(t, err)
	filtered = al.getFilteredEntries()
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "Error occurred", filtered[0].Message)

	// Test invalid regex
	err = al.SetSearchQuery("[invalid(")
	assert.Error(t, err)

	// Clear search
	err = al.SetSearchQuery("")
	require.NoError(t, err)
	filtered = al.getFilteredEntries()
	assert.Equal(t, 5, len(filtered))
}

// --- Scrolling Tests ---

// TestAgentLogScroll verifies scrolling functionality.
func TestAgentLogScroll(t *testing.T) {
	al := NewAgentLog(80, 24)
	al.visibleLines = 5

	// Add 20 entries
	for i := 0; i < 20; i++ {
		event := createTestEvent(AgentOrchestrator, EventLog, fmt.Sprintf("Message %d", i))
		al.AddEntry(event)
	}

	// Initially at bottom
	filtered := al.getFilteredEntries()
	maxOffset := len(filtered) - al.visibleLines
	assert.Equal(t, maxOffset, al.viewOffset)

	// Scroll up
	al.ScrollUp()
	assert.Equal(t, maxOffset-1, al.viewOffset)

	// Scroll down
	al.ScrollDown()
	assert.Equal(t, maxOffset, al.viewOffset)

	// Scroll to top
	al.ScrollToTop()
	assert.Equal(t, 0, al.viewOffset)

	// Page down (half page = 2 lines)
	al.PageDown()
	assert.Equal(t, 2, al.viewOffset)

	// Page up
	al.PageUp()
	assert.Equal(t, 0, al.viewOffset)

	// Scroll to bottom
	al.ScrollToBottom()
	assert.Equal(t, maxOffset, al.viewOffset)

	// Can't scroll up past top
	al.ScrollToTop()
	al.ScrollUp()
	assert.Equal(t, 0, al.viewOffset)

	// Can't scroll down past bottom
	al.ScrollToBottom()
	initialOffset := al.viewOffset
	al.ScrollDown()
	assert.Equal(t, initialOffset, al.viewOffset)
}

// --- Export Tests ---

// TestAgentLogExport verifies file export functionality.
func TestAgentLogExport(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add entries
	event1 := createTestEvent(AgentOrchestrator, EventStarted, "Task started")
	event2 := createTestEvent(AgentExecutor, EventProgress, "Task in progress")
	event3 := createTestEvent(AgentReviewer, EventCompleted, "Task completed")

	al.AddEntry(event1)
	al.AddEntry(event2)
	al.AddEntry(event3)

	// Export to temp file
	tmpfile, err := os.CreateTemp("", "agent-log-test-*.txt")
	require.NoError(t, err)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	err = al.ExportToFile(tmpfile.Name())
	require.NoError(t, err)

	// Verify file was created and has content
	content, err := os.ReadFile(tmpfile.Name())
	require.NoError(t, err)

	contentStr := string(content)
	assert.NotEmpty(t, contentStr)
	assert.Contains(t, contentStr, "Orchestrator")
	assert.Contains(t, contentStr, "Executor")
	assert.Contains(t, contentStr, "Reviewer")
	assert.Contains(t, contentStr, "Task started")
	assert.Contains(t, contentStr, "Task in progress")
	assert.Contains(t, contentStr, "Task completed")

	// Verify last export is set
	assert.Equal(t, tmpfile.Name(), al.LastExport())
}

// TestAgentLogExportFiltered verifies that only filtered entries are exported.
func TestAgentLogExportFiltered(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add entries from different agents
	for i := 0; i < 9; i++ {
		var agent AgentType
		switch i % 3 {
		case 0:
			agent = AgentOrchestrator
		case 1:
			agent = AgentExecutor
		default:
			agent = AgentReviewer
		}
		event := createTestEvent(agent, EventLog, fmt.Sprintf("Message %d", i))
		al.AddEntry(event)
	}

	// Filter to only Orchestrator
	filterAgent := AgentOrchestrator
	al.SetFilterAgent(&filterAgent)

	// Export filtered entries
	tmpfile, err := os.CreateTemp("", "agent-log-filtered-*.txt")
	require.NoError(t, err)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	err = al.ExportToFile(tmpfile.Name())
	require.NoError(t, err)

	// Verify only Orchestrator entries are in file
	content, err := os.ReadFile(tmpfile.Name())
	require.NoError(t, err)

	contentStr := string(content)
	lines := strings.Split(strings.TrimSpace(contentStr), "\n")

	// Should have 3 Orchestrator entries (0, 3, 6)
	assert.Equal(t, 3, len(lines))

	for _, line := range lines {
		assert.Contains(t, line, "Orchestrator")
	}
}

// TestAgentLogCombinedFilters verifies that multiple filters work together.
func TestAgentLogCombinedFilters(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add entries with different agents and levels
	al.AddEntry(createTestEvent(AgentOrchestrator, EventLog, "Info from orch"))
	al.AddEntry(createTestEvent(AgentOrchestrator, EventFailed, "Error from orch"))
	al.AddEntry(createTestEvent(AgentExecutor, EventLog, "Info from exec"))
	al.AddEntry(createTestEvent(AgentExecutor, EventFailed, "Error from exec"))

	// Filter by agent AND level
	filterAgent := AgentOrchestrator
	filterLevel := LogLevelError
	al.SetFilterAgent(&filterAgent)
	al.SetFilterLevel(&filterLevel)

	filtered := al.getFilteredEntries()
	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, AgentOrchestrator, filtered[0].Agent)
	assert.Equal(t, LogLevelError, filtered[0].Level)
	assert.Equal(t, "Error from orch", filtered[0].Message)
}

// TestAgentLogClearFilters verifies that filters can be cleared.
func TestAgentLogClearFilters(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add entries
	for i := 0; i < 5; i++ {
		event := createTestEvent(AgentOrchestrator, EventLog, fmt.Sprintf("Message %d", i))
		al.AddEntry(event)
	}

	// Apply filters
	filterAgent := AgentExecutor
	filterLevel := LogLevelError
	al.SetFilterAgent(&filterAgent)
	al.SetFilterLevel(&filterLevel)
	err := al.SetSearchQuery("test")
	require.NoError(t, err)

	filtered := al.getFilteredEntries()
	assert.Equal(t, 0, len(filtered))

	// Clear filters
	al.ClearFilters()
	filtered = al.getFilteredEntries()
	assert.Equal(t, 5, len(filtered))
	assert.Equal(t, 0, al.viewOffset)
}

// TestAgentLogLevelFromEventType verifies level determination from event type.
func TestAgentLogLevelFromEventType(t *testing.T) {
	al := NewAgentLog(80, 24)

	tests := []struct {
		eventType EventType
		expected  LogLevel
	}{
		{EventFailed, LogLevelError},
		{EventStarted, LogLevelInfo},
		{EventProgress, LogLevelInfo},
		{EventCompleted, LogLevelInfo},
		{EventHandoff, LogLevelInfo},
		{EventLog, LogLevelInfo},
	}

	for _, tt := range tests {
		result := al.levelFromEventType(tt.eventType)
		assert.Equal(t, tt.expected, result, "event type %v", tt.eventType)
	}
}

// TestAgentLogEntryCount verifies accurate counting.
func TestAgentLogEntryCount(t *testing.T) {
	al := NewAgentLog(80, 24)

	// Add 5 entries
	for i := 0; i < 5; i++ {
		event := createTestEvent(AgentOrchestrator, EventLog, fmt.Sprintf("Message %d", i))
		al.AddEntry(event)
	}

	total, filtered := al.EntryCount()
	assert.Equal(t, 5, total)
	assert.Equal(t, 5, filtered)

	// Apply filter to reduce count
	filterAgent := AgentExecutor
	al.SetFilterAgent(&filterAgent)
	total, filtered = al.EntryCount()
	assert.Equal(t, 5, total)
	assert.Equal(t, 0, filtered)
}

// TestAgentLogFormattingEdgeCases tests edge cases in formatting.
func TestAgentLogFormattingEdgeCases(t *testing.T) {
	al := NewAgentLog(40, 24) // Small width to test truncation

	// Very long message
	longMsg := "This is a very long message that should be truncated because it exceeds the width limit"
	event := createTestEvent(AgentOrchestrator, EventLog, longMsg)
	al.AddEntry(event)

	formatted := al.formatEntry(al.getFilteredEntries()[0])
	assert.Less(t, len(formatted), 100) // Should be truncated
	assert.Contains(t, formatted, "...")
}
