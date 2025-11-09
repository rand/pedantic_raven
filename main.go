package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/app/events"
	"github.com/rand/pedantic-raven/internal/config"
	"github.com/rand/pedantic-raven/internal/editor"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"github.com/rand/pedantic-raven/internal/gliner"
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

// Application model integrating all foundation components.
type model struct {
	ready           bool
	eventBroker     *events.Broker
	modeRegistry    *modes.Registry
	overlayManager  *overlay.Manager
	paletteRegistry *palette.CommandRegistry
	width           int
	height          int
	eventLog        []string
}

func initialModel() model {
	// Load configuration
	cfg, err := config.Load("config.toml")
	if err != nil {
		// Fall back to default config if file doesn't exist or can't be loaded
		cfg = config.DefaultConfig()
	}

	// Create semantic analyzer with configured extractor
	analyzer := createAnalyzer(cfg)

	// Create event broker
	broker := events.NewBroker(100)

	// Create mode registry
	modeRegistry := modes.NewRegistry()

	// Register real modes with configured analyzer
	editMode := editor.NewEditModeWithAnalyzer(analyzer)
	exploreMode := modes.NewExploreMode()
	analyzeMode := modes.NewBaseMode(modes.ModeAnalyze, "Analyze", "Semantic analysis")

	modeRegistry.Register(editMode)
	modeRegistry.Register(exploreMode)
	modeRegistry.Register(analyzeMode)

	// Set Edit as default mode
	modeRegistry.SwitchTo(modes.ModeEdit)

	// Create overlay manager
	overlayManager := overlay.NewManager()

	// Create command palette registry
	paletteRegistry := palette.NewCommandRegistry()

	m := model{
		ready:           false,
		eventBroker:     broker,
		modeRegistry:    modeRegistry,
		overlayManager:  overlayManager,
		paletteRegistry: paletteRegistry,
		eventLog:        make([]string, 0),
	}

	// Register commands
	m.registerCommands()

	return m
}

// createAnalyzer creates a semantic analyzer based on configuration.
func createAnalyzer(cfg *config.Config) semantic.Analyzer {
	// Always create pattern extractor as fallback
	patternExtractor := semantic.NewPatternExtractor()

	// If GLiNER is disabled, use pattern extractor only
	if !cfg.GLiNER.Enabled {
		return semantic.NewAnalyzerWithExtractor(patternExtractor)
	}

	// Create GLiNER client with config
	glinerConfig := &gliner.Config{
		ServiceURL:        cfg.GLiNER.ServiceURL,
		Timeout:           cfg.GLiNER.Timeout,
		MaxRetries:        cfg.GLiNER.MaxRetries,
		Enabled:           cfg.GLiNER.Enabled,
		FallbackToPattern: cfg.GLiNER.FallbackToPattern,
	}

	glinerClient := gliner.NewClient(glinerConfig)

	// Combine default and custom entity types
	defaultTypes := cfg.GLiNER.EntityTypes.Default
	if len(cfg.GLiNER.EntityTypes.Custom) > 0 {
		defaultTypes = append(defaultTypes, cfg.GLiNER.EntityTypes.Custom...)
	}

	// Create GLiNER extractor
	glinerExtractor := semantic.NewGLiNERExtractor(
		glinerClient,
		defaultTypes,
		cfg.GLiNER.ScoreThreshold,
	)

	// Check if GLiNER is actually available
	ctx := context.Background()
	if glinerExtractor.IsAvailable(ctx) {
		// GLiNER is available - use hybrid with GLiNER as primary
		if cfg.GLiNER.FallbackToPattern {
			hybridExtractor := semantic.NewHybridExtractor(glinerExtractor, patternExtractor, true)
			return semantic.NewAnalyzerWithExtractor(hybridExtractor)
		}
		// GLiNER only, no fallback
		return semantic.NewAnalyzerWithExtractor(glinerExtractor)
	}

	// GLiNER not available - use pattern extractor
	return semantic.NewAnalyzerWithExtractor(patternExtractor)
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

	// Help commands
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
				"Interactive Context Engineering Environment\n\nPhase 2 Complete:\n✓ Semantic Analysis\n✓ Context Panel\n✓ Terminal Integration\n✓ Edit Mode\n\n291 tests passing",
				func() tea.Cmd {
					m.logEvent("Closed about dialog")
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
	// Initialize current mode
	if currentMode := m.modeRegistry.Current(); currentMode != nil {
		return currentMode.Init()
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle window size
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.ready = true
		m.width = wsMsg.Width
		m.height = wsMsg.Height
		m.overlayManager.Update(wsMsg)

		// Forward to current mode
		if currentMode := m.modeRegistry.Current(); currentMode != nil {
			_, cmd := currentMode.Update(wsMsg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
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
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+k":
			// Open command palette
			p := palette.NewPalette("palette", m.paletteRegistry)
			cmd := m.overlayManager.Push(p)
			cmds = append(cmds, cmd)
			m.logEvent("Opened command palette")
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

		case "?":
			// Show about dialog
			dialog := overlay.NewMessageDialog(
				"about",
				"Pedantic Raven",
				"Interactive Context Engineering Environment\n\nPhase 2 Complete:\n✓ Semantic Analysis\n✓ Context Panel\n✓ Terminal Integration\n✓ Edit Mode\n\n291 tests passing",
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

	// Update current mode
	if currentMode := m.modeRegistry.Current(); currentMode != nil {
		updatedMode, cmd := currentMode.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Update mode in registry (modes are mutable)
		_ = updatedMode
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Title bar
	currentModeID := m.modeRegistry.CurrentID()
	modeName := string(currentModeID)
	if modeName == "" {
		modeName = "None"
	}

	title := titleStyle.Render("🐦 Pedantic Raven - Phase 2")
	status := statusStyle.Render(fmt.Sprintf("Mode: %s | %dx%d", modeName, m.width, m.height))

	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(status)
	b.WriteString("\n\n")

	// Render current mode
	if currentMode := m.modeRegistry.Current(); currentMode != nil {
		modeView := currentMode.View()
		b.WriteString(modeView)
	} else {
		b.WriteString("No mode active")
	}

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
	help := helpStyle.Render("Keys: 1,2,3=modes | Tab=focus | Ctrl+K=palette | ?=about | Ctrl+C=quit")
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
