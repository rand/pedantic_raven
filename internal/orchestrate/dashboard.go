// Package orchestrate provides types and structures for Orchestrate Mode,
// which enables real-time coordination of multi-agent systems via mnemosyne.
package orchestrate

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Dashboard implements tea.Model for real-time monitoring of mnemosyne orchestration.
// It displays 4 agent statuses, progress metrics, and upcoming tasks.
type Dashboard struct {
	// Data
	session    *SessionState
	agents     map[AgentType]*AgentStatus
	taskQueue  []string // Pending task IDs
	taskMap    map[string]*Task

	// Event stream
	events     <-chan *AgentEvent
	eventMu    sync.Mutex
	lastUpdate time.Time

	// Metrics
	startTime      time.Time
	totalTasks     int
	completedTasks int
	failedTasks    int
	successRate    float64

	// UI state
	width       int
	height      int
	refreshRate time.Duration

	// Status
	orchestrating bool
	mu            sync.RWMutex
}

// NewDashboard creates a new dashboard from a session state and event stream.
func NewDashboard(session *SessionState, events <-chan *AgentEvent) *Dashboard {
	taskMap := make(map[string]*Task)
	if session.Plan != nil {
		for i := range session.Plan.Tasks {
			taskMap[session.Plan.Tasks[i].ID] = &session.Plan.Tasks[i]
		}
	}

	dashboard := &Dashboard{
		session:        session,
		agents:         session.Agents,
		taskQueue:      make([]string, 0),
		taskMap:        taskMap,
		events:         events,
		startTime:      session.StartTime,
		totalTasks:     session.TotalTasks,
		completedTasks: session.CompletedTasks,
		failedTasks:    session.FailedTasks,
		width:          80,
		height:         24,
		refreshRate:    100 * time.Millisecond,
		orchestrating:  true,
		lastUpdate:     time.Now(),
	}

	// Initialize task queue with pending tasks
	dashboard.initializeTaskQueue()

	return dashboard
}

// initializeTaskQueue populates the task queue from session state.
func (d *Dashboard) initializeTaskQueue() {
	if d.session.Plan == nil {
		return
	}

	for _, task := range d.session.Plan.Tasks {
		status, ok := d.session.TaskStatuses[task.ID]
		if !ok || status == TaskStatusPending {
			d.taskQueue = append(d.taskQueue, task.ID)
		}
	}
}

// Init initializes the dashboard (implements tea.Model interface).
func (d *Dashboard) Init() tea.Cmd {
	return d.pollEvents()
}

// Update handles messages and updates the dashboard state (implements tea.Model interface).
func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch msg := msg.(type) {
	case AgentEventMsg:
		if msg.Event != nil {
			d.handleEvent(msg.Event)
		}
		return d, d.pollEvents()

	case tickMsg:
		d.calculateMetrics()
		return d, tick(d.refreshRate)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			d.orchestrating = false
			return d, tea.Quit
		}

	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
	}

	return d, nil
}

// View renders the dashboard UI (implements tea.Model interface).
func (d *Dashboard) View() string {
	if d.width < 40 || d.height < 10 {
		return "Terminal too small. Minimum size: 40x10"
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	var output strings.Builder

	// Title
	title := fmt.Sprintf("ORCHESTRATE MODE - Session: %s", d.session.ID)
	output.WriteString(d.renderBox(title, d.width))
	output.WriteString("\n")

	// Agent panel
	agentPanel := d.renderAgentPanel()
	output.WriteString(d.renderBox("Agents", d.width))
	output.WriteString(agentPanel)
	output.WriteString("\n")

	// Progress section
	output.WriteString(d.renderBox("Progress", d.width))
	output.WriteString(d.renderProgress())
	output.WriteString("\n")

	// Task queue section
	if len(d.taskQueue) > 0 {
		output.WriteString(d.renderBox("Upcoming Tasks", d.width))
		output.WriteString(d.renderTaskQueue())
	}

	return output.String()
}

// renderBox renders a section header box.
func (d *Dashboard) renderBox(title string, width int) string {
	if width < len(title)+4 {
		return title
	}

	padding := width - len(title) - 4
	leftPad := padding / 2
	rightPad := padding - leftPad

	return fmt.Sprintf("╔%s %s %s╗",
		strings.Repeat("═", leftPad),
		title,
		strings.Repeat("═", rightPad))
}

// renderAgentPanel renders the agent status section.
func (d *Dashboard) renderAgentPanel() string {
	var output strings.Builder

	agents := []AgentType{AgentOrchestrator, AgentOptimizer, AgentReviewer, AgentExecutor}
	for _, agentType := range agents {
		if agent, ok := d.agents[agentType]; ok {
			output.WriteString("  ")
			output.WriteString(d.renderAgentStatus(agent))
			output.WriteString("\n")
		}
	}

	return output.String()
}

// renderAgentStatus renders a single agent status line.
func (d *Dashboard) renderAgentStatus(agent *AgentStatus) string {
	var indicator string
	var statusColor lipgloss.Color

	switch agent.Status {
	case "active":
		indicator = "[●]"
		statusColor = lipgloss.Color("10") // Green
	case "idle":
		indicator = "[◐]"
		statusColor = lipgloss.Color("8") // Gray
	case "error":
		indicator = "[✗]"
		statusColor = lipgloss.Color("9") // Red
	default:
		indicator = "[○]"
		statusColor = lipgloss.Color("8")
	}

	indicatorStyle := lipgloss.NewStyle().Foreground(statusColor)
	statusStyle := lipgloss.NewStyle().Foreground(statusColor)

	task := agent.CurrentTask
	if task == "" {
		task = fmt.Sprintf("Last: %s", agent.LastUpdate.Format("15:04:05"))
	} else {
		// Truncate task ID if too long
		if len(task) > 40 {
			task = task[:37] + "..."
		}
	}

	agentName := fmt.Sprintf("%-15s", agent.Agent.String())
	status := fmt.Sprintf("%-7s", agent.Status)

	return fmt.Sprintf("%s %s (%s) Task: %s",
		indicatorStyle.Render(indicator),
		agentName,
		statusStyle.Render(status),
		task)
}

// renderProgress renders the progress section with metrics.
func (d *Dashboard) renderProgress() string {
	var output strings.Builder

	// Progress bar
	progressBar := d.renderProgressBar()
	output.WriteString("  ")
	output.WriteString(progressBar)
	output.WriteString("\n")

	// Success rate
	if d.completedTasks+d.failedTasks > 0 {
		successRateStr := fmt.Sprintf("  Success Rate: %d/%d ✓ (%.2f%%)\n",
			d.completedTasks,
			d.completedTasks+d.failedTasks,
			d.successRate)
		output.WriteString(successRateStr)
	}

	// Elapsed time
	elapsed := time.Since(d.startTime)
	output.WriteString("  ")
	output.WriteString(fmt.Sprintf("Elapsed: %s\n", formatElapsed(elapsed)))

	return output.String()
}

// renderProgressBar renders an ASCII progress bar.
func (d *Dashboard) renderProgressBar() string {
	barWidth := 20
	if d.width > 60 {
		barWidth = (d.width - 30) / 2
	}
	if barWidth < 10 {
		barWidth = 10
	}

	var filled int
	if d.totalTasks > 0 {
		filled = (d.completedTasks * barWidth) / d.totalTasks
		if filled > barWidth {
			filled = barWidth
		}
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	percentage := 0
	if d.totalTasks > 0 {
		percentage = (d.completedTasks * 100) / d.totalTasks
	}

	return fmt.Sprintf("Progress: [%s] %d/%d tasks (%d%%)",
		bar, d.completedTasks, d.totalTasks, percentage)
}

// renderTaskQueue renders the upcoming tasks section.
func (d *Dashboard) renderTaskQueue() string {
	var output strings.Builder

	// Show up to 5 upcoming tasks
	maxTasks := 5
	if len(d.taskQueue) < maxTasks {
		maxTasks = len(d.taskQueue)
	}

	for i := 0; i < maxTasks; i++ {
		taskID := d.taskQueue[i]
		taskDesc := ""
		if task, ok := d.taskMap[taskID]; ok {
			taskDesc = task.Description
		}

		// Truncate description if too long
		if len(taskDesc) > d.width-20 {
			taskDesc = taskDesc[:d.width-23] + "..."
		}

		output.WriteString(fmt.Sprintf("  %d. %s: %s\n", i+1, taskID, taskDesc))
	}

	return output.String()
}

// handleEvent processes an agent event and updates dashboard state.
func (d *Dashboard) handleEvent(event *AgentEvent) {
	d.eventMu.Lock()
	defer d.eventMu.Unlock()

	// Update session state
	if d.session != nil {
		d.session.UpdateProgress(event)
	}

	// Update local metrics
	switch event.EventType {
	case EventStarted:
		// Task started - update queue if present
		d.removeFromQueue(event.TaskID)

	case EventCompleted:
		d.completedTasks++
		d.removeFromQueue(event.TaskID)
		d.calculateMetrics()

	case EventFailed:
		d.failedTasks++
		d.removeFromQueue(event.TaskID)
		d.calculateMetrics()

	case EventProgress:
		// Just update timestamp
	}

	d.lastUpdate = time.Now()
}

// removeFromQueue removes a task from the pending queue.
func (d *Dashboard) removeFromQueue(taskID string) {
	for i, id := range d.taskQueue {
		if id == taskID {
			d.taskQueue = append(d.taskQueue[:i], d.taskQueue[i+1:]...)
			break
		}
	}
}

// calculateMetrics updates derived metrics.
func (d *Dashboard) calculateMetrics() {
	total := d.completedTasks + d.failedTasks
	if total == 0 {
		d.successRate = 0.0
	} else {
		d.successRate = (float64(d.completedTasks) / float64(total)) * 100
	}
}

// pollEvents polls the event channel non-blocking and returns a command.
func (d *Dashboard) pollEvents() tea.Cmd {
	return func() tea.Msg {
		select {
		case event := <-d.events:
			if event != nil {
				return AgentEventMsg{Event: event}
			}
		case <-time.After(10 * time.Millisecond):
			// Timeout, continue
		}
		return tickMsg{}
	}
}

// --- Custom Tea Messages ---

// AgentEventMsg wraps an AgentEvent for tea.Model processing.
type AgentEventMsg struct {
	Event *AgentEvent
}

// tickMsg represents a timer tick.
type tickMsg struct{}

// tick returns a command that sends a tickMsg after the given duration.
func tick(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// formatElapsed formats a duration into human-readable format.
func formatElapsed(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
