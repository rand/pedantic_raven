package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/app/events"
	"github.com/rand/pedantic-raven/internal/layout"
	"github.com/rand/pedantic-raven/internal/modes"
	"github.com/rand/pedantic-raven/internal/overlay"
	"github.com/rand/pedantic-raven/internal/palette"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Padding(1)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(0, 2)

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("120")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)
)

// DemoComponent is a simple component for testing.
type DemoComponent struct {
	id      layout.PaneID
	content string
}

func NewDemoComponent(id layout.PaneID, content string) *DemoComponent {
	return &DemoComponent{
		id:      id,
		content: content,
	}
}

func (c *DemoComponent) ID() layout.PaneID {
	return c.id
}

func (c *DemoComponent) Update(msg tea.Msg) (layout.Component, tea.Cmd) {
	return c, nil
}

func (c *DemoComponent) View(area layout.Rect, focused bool) string {
	style := lipgloss.NewStyle().
		Width(area.Width).
		Height(area.Height).
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	if focused {
		style = style.BorderForeground(lipgloss.Color("170"))
	} else {
		style = style.BorderForeground(lipgloss.Color("240"))
	}

	return style.Render(c.content)
}

// Application model integrating all foundation components.
type model struct {
	ready           bool
	eventBroker     *events.Broker
	modeRegistry    *modes.Registry
	layoutEngine    *layout.Engine
	overlayManager  *overlay.Manager
	paletteRegistry *palette.CommandRegistry
	width           int
	height          int
	eventLog        []string
}

func initialModel() model {
	// Create event broker
	broker := events.NewBroker(100)

	// Create mode registry
	modeRegistry := modes.NewRegistry()

	// Register demo modes
	editMode := modes.NewBaseMode(modes.ModeEdit, "Edit", "Context editing mode")
	exploreMode := modes.NewBaseMode(modes.ModeExplore, "Explore", "Memory workspace")
	analyzeMode := modes.NewBaseMode(modes.ModeAnalyze, "Analyze", "Semantic analysis")

	modeRegistry.Register(editMode)
	modeRegistry.Register(exploreMode)
	modeRegistry.Register(analyzeMode)

	// Create layout engine
	layoutEngine := layout.NewEngine(layout.LayoutStandard)

	// Add demo components
	editor := NewDemoComponent("editor", "Editor Pane\n\nEdit your context here...")
	sidebar := NewDemoComponent("sidebar", "Sidebar\n\nMemory notes\nTriples\nAgents")
	terminal := NewDemoComponent("terminal", "Terminal\n\n$ mnemosyne recall \"test\"")

	layoutEngine.RegisterComponent(editor)
	layoutEngine.RegisterComponent(sidebar)
	layoutEngine.RegisterComponent(terminal)

	// Create overlay manager
	overlayManager := overlay.NewManager()

	// Create command palette registry
	paletteRegistry := palette.NewCommandRegistry()

	m := model{
		ready:           false,
		eventBroker:     broker,
		modeRegistry:    modeRegistry,
		layoutEngine:    layoutEngine,
		overlayManager:  overlayManager,
		paletteRegistry: paletteRegistry,
		eventLog:        make([]string, 0),
	}

	// Register commands
	m.registerCommands()

	return m
}

func (m *model) registerCommands() {
	// Mode switching commands
	m.paletteRegistry.Register(palette.Command{
		ID:          "mode.edit",
		Name:        "Switch to Edit Mode",
		Description: "Context editing with semantic analysis",
		Keybinding:  "1",
		Category:    palette.CategoryMode,
		Execute: func() tea.Cmd {
			m.modeRegistry.SwitchTo(modes.ModeEdit)
			m.logEvent("Switched to Edit mode")
			return nil
		},
	})

	m.paletteRegistry.Register(palette.Command{
		ID:          "mode.explore",
		Name:        "Switch to Explore Mode",
		Description: "Memory workspace with graph visualization",
		Keybinding:  "2",
		Category:    palette.CategoryMode,
		Execute: func() tea.Cmd {
			m.modeRegistry.SwitchTo(modes.ModeExplore)
			m.logEvent("Switched to Explore mode")
			return nil
		},
	})

	m.paletteRegistry.Register(palette.Command{
		ID:          "mode.analyze",
		Name:        "Switch to Analyze Mode",
		Description: "Semantic insights and triple analysis",
		Keybinding:  "3",
		Category:    palette.CategoryMode,
		Execute: func() tea.Cmd {
			m.modeRegistry.SwitchTo(modes.ModeAnalyze)
			m.logEvent("Switched to Analyze mode")
			return nil
		},
	})

	// Layout commands
	m.paletteRegistry.Register(palette.Command{
		ID:          "layout.focus",
		Name:        "Focus Layout",
		Description: "Single large editor pane",
		Keybinding:  "F",
		Category:    palette.CategoryView,
		Execute: func() tea.Cmd {
			m.layoutEngine.SetMode(layout.LayoutFocus)
			m.logEvent("Layout: Focus")
			return nil
		},
	})

	m.paletteRegistry.Register(palette.Command{
		ID:          "layout.standard",
		Name:        "Standard Layout",
		Description: "Editor + sidebar + terminal",
		Keybinding:  "S",
		Category:    palette.CategoryView,
		Execute: func() tea.Cmd {
			m.layoutEngine.SetMode(layout.LayoutStandard)
			m.logEvent("Layout: Standard")
			return nil
		},
	})

	// Overlay commands
	m.paletteRegistry.Register(palette.Command{
		ID:          "help.about",
		Name:        "About Pedantic Raven",
		Description: "Show information about this application",
		Keybinding:  "?",
		Category:    palette.CategoryHelp,
		Execute: func() tea.Cmd {
			dialog := overlay.NewMessageDialog(
				"about",
				"Pedantic Raven",
				"Interactive Context Engineering Environment\n\nPhase 1 Foundation Complete:\n✓ PubSub Events\n✓ Layout Engine\n✓ Mode Registry\n✓ Overlay System\n✓ Command Palette",
				func() tea.Cmd {
					m.logEvent("Closed about dialog")
					return nil
				},
			)
			return m.overlayManager.Push(dialog)
		},
	})

	m.paletteRegistry.Register(palette.Command{
		ID:          "help.confirm",
		Name:        "Test Confirm Dialog",
		Description: "Show a confirmation dialog",
		Keybinding:  "C",
		Category:    palette.CategoryHelp,
		Execute: func() tea.Cmd {
			dialog := overlay.NewConfirmDialog(
				"test-confirm",
				"Confirmation",
				"This is a test confirmation dialog.\n\nDo you want to proceed?",
				func() tea.Cmd {
					m.logEvent("User clicked Yes")
					return nil
				},
				func() tea.Cmd {
					m.logEvent("User clicked No")
					return nil
				},
			)
			return m.overlayManager.Push(dialog)
		},
	})
}

func (m *model) logEvent(msg string) {
	m.eventLog = append(m.eventLog, msg)
	if len(m.eventLog) > 5 {
		m.eventLog = m.eventLog[len(m.eventLog)-5:]
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.layoutEngine.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle window size
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.ready = true
		m.width = wsMsg.Width
		m.height = wsMsg.Height
		m.layoutEngine.SetTerminalSize(wsMsg.Width, wsMsg.Height)
		m.overlayManager.Update(wsMsg)
		return m, nil
	}

	// Handle overlays first (if present)
	if !m.overlayManager.IsEmpty() {
		cmd := m.overlayManager.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// If modal overlay present, block input to underlying UI
		if m.overlayManager.HasModal() {
			return m, tea.Batch(cmds...)
		}
	}

	// Handle global keys
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+k":
			// Open command palette
			p := palette.NewPalette("palette", m.paletteRegistry)
			cmd := m.overlayManager.Push(p)
			cmds = append(cmds, cmd)
			m.logEvent("Opened command palette")
			return m, tea.Batch(cmds...)

		case "tab":
			// Focus next pane
			m.layoutEngine.FocusNext()
			m.logEvent("Focus next pane")
			return m, tea.Batch(cmds...)

		case "1", "2", "3":
			// Quick mode switching
			switch keyMsg.String() {
			case "1":
				cmd := m.modeRegistry.SwitchTo(modes.ModeEdit)
				cmds = append(cmds, cmd)
				m.logEvent("Switched to Edit mode")
			case "2":
				cmd := m.modeRegistry.SwitchTo(modes.ModeExplore)
				cmds = append(cmds, cmd)
				m.logEvent("Switched to Explore mode")
			case "3":
				cmd := m.modeRegistry.SwitchTo(modes.ModeAnalyze)
				cmds = append(cmds, cmd)
				m.logEvent("Switched to Analyze mode")
			}
			return m, tea.Batch(cmds...)

		case "f", "F":
			m.layoutEngine.SetMode(layout.LayoutFocus)
			m.logEvent("Layout: Focus")
			return m, tea.Batch(cmds...)

		case "s", "S":
			m.layoutEngine.SetMode(layout.LayoutStandard)
			m.logEvent("Layout: Standard")
			return m, tea.Batch(cmds...)

		case "?":
			// Show about dialog
			dialog := overlay.NewMessageDialog(
				"about",
				"Pedantic Raven",
				"Interactive Context Engineering Environment\n\nPhase 1 Foundation Complete:\n✓ PubSub Events\n✓ Layout Engine\n✓ Mode Registry\n✓ Overlay System\n✓ Command Palette",
				func() tea.Cmd {
					m.logEvent("Closed about dialog")
					return nil
				},
			)
			cmd := m.overlayManager.Push(dialog)
			cmds = append(cmds, cmd)
			m.logEvent("Opened about dialog")
			return m, tea.Batch(cmds...)
		}
	}

	// Update layout engine
	_, cmd := m.layoutEngine.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Title bar
	currentMode := m.modeRegistry.CurrentID()
	if currentMode == "" {
		currentMode = "None"
	}

	title := titleStyle.Render("🐦 Pedantic Raven - Phase 1 Demo")
	status := statusStyle.Render(fmt.Sprintf("Mode: %s | Layout: %s", currentMode, m.layoutEngine.Mode()))

	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(status)
	b.WriteString("\n\n")

	// Main layout
	layoutView := m.layoutEngine.View()
	b.WriteString(layoutView)
	b.WriteString("\n\n")

	// Event log
	if len(m.eventLog) > 0 {
		b.WriteString(helpStyle.Render("Recent Events:"))
		b.WriteString("\n")
		for _, event := range m.eventLog {
			b.WriteString(helpStyle.Render("  • " + event))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Help
	help := helpStyle.Render("Keys: 1,2,3=modes | f,s=layout | Tab=focus | Ctrl+K=palette | ?=about | q=quit")
	b.WriteString(help)

	// Render overlays on top
	if !m.overlayManager.IsEmpty() {
		b.WriteString("\n\n")
		b.WriteString(m.overlayManager.View())
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
