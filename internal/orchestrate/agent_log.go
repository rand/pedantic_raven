package orchestrate

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LogLevel represents the severity level of a log entry.
type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelWarn
	LogLevelError
)

// String returns the string representation of LogLevel.
func (l LogLevel) String() string {
	switch l {
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a single entry in the agent communication log.
type LogEntry struct {
	Timestamp time.Time
	Agent     AgentType
	EventType EventType
	TaskID    string
	Message   string
	Level     LogLevel
}

// AgentLog manages a scrollable, filterable log of agent communications.
// It implements the tea.Model interface for use in TUI applications.
type AgentLog struct {
	// Log buffer (circular buffer, max 10,000 lines)
	entries    []LogEntry
	maxEntries int

	// Filtering
	filterAgent *AgentType
	filterLevel *LogLevel
	searchQuery string
	searchRegex *regexp.Regexp

	// Viewport
	viewOffset   int
	visibleLines int

	// UI state
	width  int
	height int

	// Export
	lastExport string
}

// NewAgentLog creates a new AgentLog with default configuration.
func NewAgentLog(width, height int) *AgentLog {
	return &AgentLog{
		entries:      make([]LogEntry, 0, 100),
		maxEntries:   10000,
		viewOffset:   0,
		visibleLines: height - 3, // Account for header and footer
		width:        width,
		height:       height,
		lastExport:   "",
	}
}

// AddEntry adds a new log entry from an AgentEvent.
func (al *AgentLog) AddEntry(event *AgentEvent) {
	entry := LogEntry{
		Timestamp: event.Timestamp,
		Agent:     event.Agent,
		EventType: event.EventType,
		TaskID:    event.TaskID,
		Message:   event.Message,
		Level:     al.levelFromEventType(event.EventType),
	}

	al.entries = append(al.entries, entry)

	// Trim if exceeds max (keep last maxEntries)
	if len(al.entries) > al.maxEntries {
		al.entries = al.entries[len(al.entries)-al.maxEntries:]
	}

	// Auto-scroll to bottom when new entry added
	al.ScrollToBottom()
}

// levelFromEventType determines the log level based on event type.
func (al *AgentLog) levelFromEventType(et EventType) LogLevel {
	switch et {
	case EventFailed:
		return LogLevelError
	case EventStarted, EventProgress, EventCompleted, EventHandoff, EventLog:
		return LogLevelInfo
	default:
		return LogLevelInfo
	}
}

// SetFilterAgent sets the agent filter (nil = show all agents).
func (al *AgentLog) SetFilterAgent(agent *AgentType) {
	al.filterAgent = agent
	al.viewOffset = 0 // Reset scroll position
}

// SetFilterLevel sets the level filter (nil = show all levels).
func (al *AgentLog) SetFilterLevel(level *LogLevel) {
	al.filterLevel = level
	al.viewOffset = 0
}

// SetSearchQuery sets the search pattern (regex).
func (al *AgentLog) SetSearchQuery(query string) error {
	if query == "" {
		al.searchQuery = ""
		al.searchRegex = nil
		return nil
	}

	re, err := regexp.Compile(query)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	al.searchQuery = query
	al.searchRegex = re
	al.viewOffset = 0
	return nil
}

// ClearFilters removes all active filters.
func (al *AgentLog) ClearFilters() {
	al.filterAgent = nil
	al.filterLevel = nil
	al.searchQuery = ""
	al.searchRegex = nil
	al.viewOffset = 0
}

// getFilteredEntries returns entries matching the current filters.
func (al *AgentLog) getFilteredEntries() []LogEntry {
	var filtered []LogEntry

	for _, entry := range al.entries {
		// Filter by agent
		if al.filterAgent != nil && entry.Agent != *al.filterAgent {
			continue
		}

		// Filter by level
		if al.filterLevel != nil {
			if *al.filterLevel == LogLevelError && entry.Level != LogLevelError {
				continue
			}
			if *al.filterLevel == LogLevelWarn && entry.Level == LogLevelInfo {
				continue
			}
		}

		// Filter by search
		if al.searchRegex != nil && !al.searchRegex.MatchString(entry.Message) {
			continue
		}

		filtered = append(filtered, entry)
	}

	return filtered
}

// ScrollUp moves the viewport up by one line.
func (al *AgentLog) ScrollUp() {
	if al.viewOffset > 0 {
		al.viewOffset--
	}
}

// ScrollDown moves the viewport down by one line.
func (al *AgentLog) ScrollDown() {
	filtered := al.getFilteredEntries()
	maxOffset := len(filtered) - al.visibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}

	if al.viewOffset < maxOffset {
		al.viewOffset++
	}
}

// PageUp moves the viewport up by half a page.
func (al *AgentLog) PageUp() {
	pageSize := al.visibleLines / 2
	if pageSize < 1 {
		pageSize = 1
	}
	for i := 0; i < pageSize; i++ {
		al.ScrollUp()
	}
}

// PageDown moves the viewport down by half a page.
func (al *AgentLog) PageDown() {
	pageSize := al.visibleLines / 2
	if pageSize < 1 {
		pageSize = 1
	}
	for i := 0; i < pageSize; i++ {
		al.ScrollDown()
	}
}

// ScrollToTop moves to the beginning of the log.
func (al *AgentLog) ScrollToTop() {
	al.viewOffset = 0
}

// ScrollToBottom moves to the end of the log.
func (al *AgentLog) ScrollToBottom() {
	filtered := al.getFilteredEntries()
	maxOffset := len(filtered) - al.visibleLines
	if maxOffset < 0 {
		al.viewOffset = 0
	} else {
		al.viewOffset = maxOffset
	}
}

// formatEntry formats a single log entry for display.
func (al *AgentLog) formatEntry(entry LogEntry) string {
	timestamp := entry.Timestamp.Format("15:04:05")

	agent := entry.Agent.String()
	if len(agent) > 12 {
		agent = agent[:12]
	}

	// Level with color
	var levelStr string
	var style lipgloss.Style
	switch entry.Level {
	case LogLevelInfo:
		levelStr = "INFO "
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // Blue
	case LogLevelWarn:
		levelStr = "WARN "
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
	case LogLevelError:
		levelStr = "ERROR"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
	}

	// Message (truncate to fit width)
	message := entry.Message
	maxLen := al.width - 40
	if maxLen < 20 {
		maxLen = 20
	}
	if len(message) > maxLen {
		message = message[:maxLen-3] + "..."
	}

	return fmt.Sprintf("[%s] [%-12s] %s: %s",
		timestamp,
		agent,
		style.Render(levelStr),
		message)
}

// ExportToFile exports the currently visible (filtered) logs to a file.
func (al *AgentLog) ExportToFile(filename string) error {
	filtered := al.getFilteredEntries()

	var lines []string
	for _, entry := range filtered {
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
			entry.Timestamp.Format(time.RFC3339),
			entry.Agent.String(),
			entry.Level.String(),
			entry.TaskID,
			entry.Message)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return err
	}

	al.lastExport = filename
	return nil
}

// LastExport returns the path of the last exported file.
func (al *AgentLog) LastExport() string {
	return al.lastExport
}

// TotalEntries returns the total number of entries in the buffer.
func (al *AgentLog) TotalEntries() int {
	return len(al.entries)
}

// FilteredEntries returns the number of entries matching current filters.
func (al *AgentLog) FilteredEntries() int {
	return len(al.getFilteredEntries())
}

// EntryCount returns the total entries and filtered entries.
func (al *AgentLog) EntryCount() (total, filtered int) {
	return len(al.entries), len(al.getFilteredEntries())
}

// --- tea.Model interface implementation ---

// Init initializes the model.
func (al *AgentLog) Init() tea.Cmd {
	return nil
}

// Update processes messages.
func (al *AgentLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			al.ScrollDown()
		case "k", "up":
			al.ScrollUp()
		case "ctrl+d":
			al.PageDown()
		case "ctrl+u":
			al.PageUp()
		case "g":
			al.ScrollToTop()
		case "G":
			al.ScrollToBottom()
		}
	case tea.WindowSizeMsg:
		al.width = msg.Width
		al.height = msg.Height
		al.visibleLines = msg.Height - 3
	}

	return al, nil
}

// View renders the log viewer.
func (al *AgentLog) View() string {
	var lines []string

	// Header
	headerText := "Agent Communication Log"
	if al.filterAgent != nil {
		headerText += fmt.Sprintf(" (Agent: %s)", al.filterAgent.String())
	}
	if al.searchQuery != "" {
		headerText += fmt.Sprintf(" (Search: %s)", al.searchQuery)
	}
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("4"))
	lines = append(lines, headerStyle.Render(headerText))

	// Log entries (visible viewport)
	filtered := al.getFilteredEntries()
	start := al.viewOffset
	end := start + al.visibleLines
	if end > len(filtered) {
		end = len(filtered)
	}

	for i := start; i < end; i++ {
		entry := filtered[i]
		line := al.formatEntry(entry)
		lines = append(lines, line)
	}

	// Pad with empty lines if needed
	for i := end - start; i < al.visibleLines; i++ {
		lines = append(lines, "")
	}

	// Footer with controls
	_, filteredCount := al.EntryCount()
	footerText := fmt.Sprintf("Logs: %d/%d | j/k: scroll | g/G: top/bottom | Ctrl+U/D: page | e: export | q: quit",
		end-start, filteredCount)
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("0"))
	lines = append(lines, footerStyle.Render(footerText))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
