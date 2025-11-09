// Package orchestrate provides the main coordinator for Orchestrate Mode.
package orchestrate

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewType identifies the currently active view in Orchestrate Mode.
type ViewType int

const (
	ViewPlanEditor ViewType = iota
	ViewDashboard
	ViewTaskGraph
	ViewAgentLog
)

// String returns the string representation of ViewType.
func (v ViewType) String() string {
	switch v {
	case ViewPlanEditor:
		return "Plan Editor"
	case ViewDashboard:
		return "Dashboard"
	case ViewTaskGraph:
		return "Task Graph"
	case ViewAgentLog:
		return "Agent Logs"
	default:
		return "Unknown"
	}
}

// OrchestrateMode is the main Bubble Tea model that coordinates all 7 components.
type OrchestrateMode struct {
	// Child components (all implemented)
	planEditor *PlanEditor
	dashboard  *Dashboard
	taskGraph  *TaskGraph
	agentLog   *AgentLog
	session    *Session
	launcher   *Launcher

	// View state
	currentView ViewType
	width       int
	height      int

	// Session state
	orchestrating bool
	paused        bool
	helpVisible   bool

	// Errors
	lastError string

	// Launch options (cached for restart)
	launchOpts LaunchOptions
}

// Styling
var (
	orchestrateModeHeaderStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#FFFFFF")).
					Background(lipgloss.Color("#6600FF")).
					Padding(0, 1).
					Width(80)

	orchestrateModeFooterStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#888888")).
					Italic(true)

	orchestrateModeErrorStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FF0000")).
					Bold(true)

	orchestrateModeHelpStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#00AAFF")).
					Background(lipgloss.Color("#222222")).
					Padding(1, 2).
					Border(lipgloss.RoundedBorder())
)

// NewOrchestrateMode creates a new OrchestrateMode instance.
func NewOrchestrateMode() *OrchestrateMode {
	return &OrchestrateMode{
		planEditor:    NewPlanEditor(),
		dashboard:     nil, // Created when orchestration starts
		taskGraph:     nil, // Created when orchestration starts
		agentLog:      nil, // Created when orchestration starts
		session:       nil, // Created when orchestration starts
		launcher:      NewLauncher(),
		currentView:   ViewPlanEditor,
		width:         80,
		height:        24,
		orchestrating: false,
		paused:        false,
		helpVisible:   false,
		lastError:     "",
		launchOpts: LaunchOptions{
			DatabasePath:    "",
			PollingInterval: 100,
			MaxConcurrent:   4,
			EnableDashboard: false,
		},
	}
}

// Init initializes the model and all child components.
func (m *OrchestrateMode) Init() tea.Cmd {
	// Initialize plan editor
	editorCmd := m.planEditor.Init()

	return tea.Batch(editorCmd)
}

// Update handles messages and coordinates child component updates.
func (m *OrchestrateMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Propagate to all active components
		if m.planEditor != nil {
			m.planEditor.width = msg.Width
			m.planEditor.height = msg.Height
		}
		if m.dashboard != nil {
			m.dashboard.width = msg.Width
			m.dashboard.height = msg.Height - 4 // Account for header/footer
		}
		if m.taskGraph != nil {
			m.taskGraph.width = msg.Width
			m.taskGraph.height = msg.Height - 4
		}
		if m.agentLog != nil {
			m.agentLog.width = msg.Width
			m.agentLog.height = msg.Height - 4
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case AgentEventMsg:
		// Route event to all relevant components
		if m.orchestrating {
			if m.dashboard != nil {
				_, dashCmd := m.dashboard.Update(msg)
				cmds = append(cmds, dashCmd)
			}
			if m.agentLog != nil && msg.Event != nil {
				m.agentLog.AddEntry(msg.Event)
			}
			if m.taskGraph != nil && msg.Event != nil {
				// Update task status in graph
				m.taskGraph.UpdateStatus(msg.Event.TaskID, m.session.GetState().TaskStatuses[msg.Event.TaskID])
			}
			if m.session != nil && msg.Event != nil {
				// Update session state and persist
				if err := m.session.UpdateProgress(msg.Event); err != nil {
					m.lastError = fmt.Sprintf("Session update error: %v", err)
				}
			}
		}
		return m, tea.Batch(cmds...)

	case tickMsg:
		// Route tick to active view
		if m.orchestrating {
			switch m.currentView {
			case ViewDashboard:
				if m.dashboard != nil {
					_, tickCmd := m.dashboard.Update(msg)
					cmds = append(cmds, tickCmd)
				}
			}
		}
		return m, tea.Batch(cmds...)

	default:
		// Route to active view
		return m.routeToActiveView(msg)
	}
}

// handleKeyPress processes keyboard shortcuts.
func (m *OrchestrateMode) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global shortcuts
	switch msg.String() {
	case "q":
		if m.orchestrating {
			// Confirm before quitting
			if err := m.stopOrchestration(); err != nil {
				m.lastError = fmt.Sprintf("Stop error: %v", err)
			}
		}
		return m, tea.Quit

	case "?":
		m.helpVisible = !m.helpVisible
		return m, nil

	case "tab":
		m.nextView()
		return m, nil

	case "shift+tab":
		m.prevView()
		return m, nil

	case "1":
		m.currentView = ViewPlanEditor
		return m, nil

	case "2":
		if m.orchestrating {
			m.currentView = ViewDashboard
		}
		return m, nil

	case "3":
		if m.orchestrating {
			m.currentView = ViewTaskGraph
		}
		return m, nil

	case "4":
		if m.orchestrating {
			m.currentView = ViewAgentLog
		}
		return m, nil
	}

	// Orchestration controls (only when not in editor)
	if m.currentView != ViewPlanEditor && m.orchestrating {
		switch msg.String() {
		case " ": // Space
			m.togglePause()
			return m, nil

		case "r":
			if err := m.restartOrchestration(); err != nil {
				m.lastError = fmt.Sprintf("Restart error: %v", err)
			}
			return m, nil

		case "x":
			if err := m.cancelOrchestration(); err != nil {
				m.lastError = fmt.Sprintf("Cancel error: %v", err)
			}
			return m, nil
		}
	}

	// Plan Editor shortcuts (when in editor view)
	if m.currentView == ViewPlanEditor {
		switch msg.String() {
		case "ctrl+s":
			// Delegate to plan editor
			_, cmd := m.planEditor.Update(msg)
			return m, cmd

		case "ctrl+o":
			// Delegate to plan editor
			_, cmd := m.planEditor.Update(msg)
			return m, cmd

		case "ctrl+n":
			// Delegate to plan editor
			_, cmd := m.planEditor.Update(msg)
			return m, cmd

		case "ctrl+l":
			// Launch orchestration if plan is valid
			if err := m.launchOrchestration(); err != nil {
				m.lastError = fmt.Sprintf("Launch error: %v", err)
			} else {
				// Switch to dashboard view
				m.currentView = ViewDashboard
			}
			return m, nil
		}
	}

	// Route all other keys to active view
	return m.routeToActiveView(msg)
}

// routeToActiveView forwards messages to the currently active view's Update method.
func (m *OrchestrateMode) routeToActiveView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewPlanEditor:
		if m.planEditor != nil {
			_, cmd := m.planEditor.Update(msg)
			return m, cmd
		}

	case ViewDashboard:
		if m.dashboard != nil {
			_, cmd := m.dashboard.Update(msg)
			return m, cmd
		}

	case ViewTaskGraph:
		if m.taskGraph != nil {
			_, cmd := m.taskGraph.Update(msg)
			return m, cmd
		}

	case ViewAgentLog:
		if m.agentLog != nil {
			_, cmd := m.agentLog.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// View renders the current view with header and footer.
func (m *OrchestrateMode) View() string {
	if m.width < 40 || m.height < 10 {
		return "Terminal too small. Minimum: 40x10"
	}

	// Show help overlay if active
	if m.helpVisible {
		return m.renderHelp()
	}

	var output []string

	// Header
	output = append(output, m.renderHeader())

	// Error banner (if present)
	if m.lastError != "" {
		output = append(output, orchestrateModeErrorStyle.Render("ERROR: "+m.lastError))
	}

	// Active view content
	output = append(output, m.renderActiveView())

	// Footer
	output = append(output, m.renderFooter())

	return lipgloss.JoinVertical(lipgloss.Left, output...)
}

// renderHeader renders the mode header with session info.
func (m *OrchestrateMode) renderHeader() string {
	sessionInfo := ""
	if m.session != nil {
		state := m.session.GetState()
		sessionInfo = fmt.Sprintf(" | Session: %s | Status: %s | Progress: %.1f%%",
			state.ID[:16],
			state.Status,
			state.Progress())
	}

	pauseIndicator := ""
	if m.paused {
		pauseIndicator = " [PAUSED]"
	}

	headerText := fmt.Sprintf("ORCHESTRATE MODE - %s%s%s",
		m.currentView.String(),
		pauseIndicator,
		sessionInfo)

	style := orchestrateModeHeaderStyle.Width(m.width)
	return style.Render(headerText)
}

// renderActiveView renders the currently active view.
func (m *OrchestrateMode) renderActiveView() string {
	switch m.currentView {
	case ViewPlanEditor:
		if m.planEditor != nil {
			return m.planEditor.View()
		}
		return "Plan Editor not initialized"

	case ViewDashboard:
		if m.dashboard != nil {
			return m.dashboard.View()
		}
		return "Dashboard not available (orchestration not started)"

	case ViewTaskGraph:
		if m.taskGraph != nil {
			return m.taskGraph.View()
		}
		return "Task Graph not available (orchestration not started)"

	case ViewAgentLog:
		if m.agentLog != nil {
			return m.agentLog.View()
		}
		return "Agent Log not available (orchestration not started)"

	default:
		return "Unknown view"
	}
}

// renderFooter renders the keyboard shortcuts footer.
func (m *OrchestrateMode) renderFooter() string {
	shortcuts := []string{}

	if m.currentView == ViewPlanEditor {
		shortcuts = append(shortcuts, "Ctrl+L: Launch", "Ctrl+S: Save", "Ctrl+O: Open", "Ctrl+N: New")
	} else if m.orchestrating {
		shortcuts = append(shortcuts, "Space: Pause/Resume", "R: Restart", "X: Cancel")
	}

	shortcuts = append(shortcuts, "Tab: Next View", "1-4: Switch View", "?: Help", "Q: Quit")

	return orchestrateModeFooterStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left,
		shortcuts...,
	))
}

// renderHelp renders the help overlay.
func (m *OrchestrateMode) renderHelp() string {
	helpText := `
ORCHESTRATE MODE - KEYBOARD SHORTCUTS

Global Controls:
  q           Exit Orchestrate Mode (stops orchestration if running)
  ?           Toggle this help overlay
  Tab         Next view (Editor → Dashboard → Graph → Logs)
  Shift+Tab   Previous view
  1           Plan Editor view
  2           Dashboard view
  3           Task Graph view
  4           Agent Logs view

Plan Editor (when in Editor view):
  Ctrl+S      Save current plan to file
  Ctrl+O      Open plan from file
  Ctrl+N      Create new plan (template)
  Ctrl+L      Launch orchestration (if plan is valid)

Orchestration Controls (when running):
  Space       Pause/Resume orchestration
  r           Restart orchestration (stop + start)
  x           Cancel orchestration (stop + cleanup)

Dashboard View:
  (Automatic updates from mnemosyne orchestrate)

Task Graph View:
  +/-         Zoom in/out
  h/j/k/l     Pan viewport (left/down/up/right)
  Arrow keys  Pan viewport

Agent Logs View:
  e           Export logs to file
  /           Search logs (enter query)
  Up/Down     Scroll through logs
  f           Filter by agent type
  l           Filter by log level

Press ? again to close this help.
`

	return orchestrateModeHelpStyle.Render(helpText)
}

// nextView switches to the next view in the cycle.
func (m *OrchestrateMode) nextView() {
	if !m.orchestrating {
		// Only editor available
		m.currentView = ViewPlanEditor
		return
	}

	switch m.currentView {
	case ViewPlanEditor:
		m.currentView = ViewDashboard
	case ViewDashboard:
		m.currentView = ViewTaskGraph
	case ViewTaskGraph:
		m.currentView = ViewAgentLog
	case ViewAgentLog:
		m.currentView = ViewPlanEditor
	}
}

// prevView switches to the previous view in the cycle.
func (m *OrchestrateMode) prevView() {
	if !m.orchestrating {
		m.currentView = ViewPlanEditor
		return
	}

	switch m.currentView {
	case ViewPlanEditor:
		m.currentView = ViewAgentLog
	case ViewDashboard:
		m.currentView = ViewPlanEditor
	case ViewTaskGraph:
		m.currentView = ViewDashboard
	case ViewAgentLog:
		m.currentView = ViewTaskGraph
	}
}

// launchOrchestration starts a new orchestration session.
func (m *OrchestrateMode) launchOrchestration() error {
	if m.orchestrating {
		return fmt.Errorf("orchestration already running")
	}

	// Validate plan from editor
	if m.planEditor == nil || m.planEditor.plan == nil {
		return fmt.Errorf("no valid plan to launch")
	}

	if !m.planEditor.isValid {
		return fmt.Errorf("plan has validation errors")
	}

	plan := m.planEditor.plan

	// Create session
	m.session = NewSession(plan)

	// Create components
	m.dashboard = NewDashboard(m.session.GetState(), m.launcher.Events())

	taskGraph, err := NewTaskGraph(plan, m.width, m.height-4)
	if err != nil {
		return fmt.Errorf("failed to create task graph: %w", err)
	}
	m.taskGraph = taskGraph

	m.agentLog = NewAgentLog(m.width, m.height-4)

	// Start launcher
	if err := m.launcher.Start(plan, m.launchOpts); err != nil {
		return fmt.Errorf("failed to start launcher: %w", err)
	}

	m.orchestrating = true
	m.paused = false
	m.lastError = ""

	// Save session
	if err := m.session.Save(); err != nil {
		m.lastError = fmt.Sprintf("Session save error: %v", err)
	}

	return nil
}

// stopOrchestration gracefully stops the orchestration.
func (m *OrchestrateMode) stopOrchestration() error {
	if !m.orchestrating {
		return nil
	}

	if err := m.launcher.Stop(); err != nil {
		return fmt.Errorf("failed to stop launcher: %w", err)
	}

	// Update session status
	if m.session != nil {
		if err := m.session.SetStatus("completed"); err != nil {
			m.lastError = fmt.Sprintf("Session status update error: %v", err)
		}
	}

	m.orchestrating = false
	m.paused = false

	return nil
}

// restartOrchestration restarts the orchestration with the current plan.
func (m *OrchestrateMode) restartOrchestration() error {
	if !m.orchestrating {
		return fmt.Errorf("no orchestration to restart")
	}

	plan := m.session.GetState().Plan
	if plan == nil {
		return fmt.Errorf("no plan to restart")
	}

	// Stop current orchestration
	if err := m.stopOrchestration(); err != nil {
		return fmt.Errorf("failed to stop orchestration: %w", err)
	}

	// Small delay to ensure cleanup
	time.Sleep(100 * time.Millisecond)

	// Create new session with same plan
	m.session = NewSession(plan)

	// Recreate components
	m.dashboard = NewDashboard(m.session.GetState(), m.launcher.Events())

	taskGraph, err := NewTaskGraph(plan, m.width, m.height-4)
	if err != nil {
		return fmt.Errorf("failed to create task graph: %w", err)
	}
	m.taskGraph = taskGraph

	m.agentLog = NewAgentLog(m.width, m.height-4)

	// Restart launcher
	if err := m.launcher.Start(plan, m.launchOpts); err != nil {
		return fmt.Errorf("failed to restart launcher: %w", err)
	}

	m.orchestrating = true
	m.paused = false
	m.lastError = ""

	return nil
}

// cancelOrchestration forcefully cancels the orchestration.
func (m *OrchestrateMode) cancelOrchestration() error {
	if !m.orchestrating {
		return nil
	}

	if err := m.launcher.Stop(); err != nil {
		return fmt.Errorf("failed to cancel launcher: %w", err)
	}

	// Update session status
	if m.session != nil {
		if err := m.session.SetStatus("cancelled"); err != nil {
			m.lastError = fmt.Sprintf("Session status update error: %v", err)
		}
	}

	m.orchestrating = false
	m.paused = false

	return nil
}

// togglePause toggles the pause state.
// Note: Actual pause/resume of mnemosyne orchestrate subprocess is not implemented
// in the launcher yet. This is a UI-only toggle for now.
func (m *OrchestrateMode) togglePause() {
	m.paused = !m.paused

	if m.session != nil {
		if m.paused {
			_ = m.session.SetStatus("paused")
		} else {
			_ = m.session.SetStatus("running")
		}
	}
}

// GetCurrentPlan returns the current plan from the editor.
func (m *OrchestrateMode) GetCurrentPlan() *WorkPlan {
	if m.planEditor != nil {
		return m.planEditor.plan
	}
	return nil
}

// IsOrchestrating returns whether orchestration is currently active.
func (m *OrchestrateMode) IsOrchestrating() bool {
	return m.orchestrating
}

// GetSession returns the current session (may be nil).
func (m *OrchestrateMode) GetSession() *Session {
	return m.session
}
