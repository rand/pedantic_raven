package editor

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/context"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"github.com/rand/pedantic-raven/internal/modes"
	"github.com/rand/pedantic-raven/internal/overlay"
	"github.com/rand/pedantic-raven/internal/terminal"
)

// EditMode provides context editing with semantic analysis.
type EditMode struct {
	*modes.BaseMode

	// Components
	editor        *EditorComponent
	contextPanel  *ContextPanelComponent
	terminalComp  *TerminalComponent

	// Semantic analysis
	analyzer      semantic.Analyzer
	analyzing     bool
	lastAnalysis  time.Time
	analysisDebounce time.Duration
}

// NewEditMode creates a new Edit mode with integrated components.
func NewEditMode() *EditMode {
	base := modes.NewBaseMode(modes.ModeEdit, "Edit", "Context editing with semantic analysis")

	// Create components
	editor := NewEditorComponent()

	contextConfig := context.DefaultContextPanelConfig()
	contextPanel := context.New(contextConfig)
	contextPanelComp := NewContextPanelComponent(contextPanel)

	termConfig := terminal.DefaultTerminalConfig()
	term := terminal.New(termConfig)
	termComp := NewTerminalComponent(term)

	// Create analyzer
	analyzer := semantic.NewAnalyzer()

	mode := &EditMode{
		BaseMode:         base,
		editor:           editor,
		contextPanel:     contextPanelComp,
		terminalComp:     termComp,
		analyzer:         analyzer,
		analyzing:        false,
		analysisDebounce: 500 * time.Millisecond,
	}

	// Register components with layout engine
	engine := base.Engine()
	engine.RegisterComponent(editor)
	engine.RegisterComponent(contextPanelComp)
	engine.RegisterComponent(termComp)

	// Set up layout: editor (60%) | context panel (40%)
	// with terminal at bottom (20% height)
	mode.configureLayout()

	return mode
}

// configureLayout sets up the default layout structure.
func (m *EditMode) configureLayout() {
	// For now, just register components
	// Layout engine will handle the actual positioning
}

// ID implements modes.Mode.
func (m *EditMode) ID() modes.ModeID {
	return modes.ModeEdit
}

// Init implements modes.Mode.
func (m *EditMode) Init() tea.Cmd {
	return m.BaseMode.Init()
}

// OnEnter implements modes.Mode.
func (m *EditMode) OnEnter() tea.Cmd {
	// Trigger initial analysis if there's content
	if m.editor.GetContent() != "" {
		return m.triggerAnalysis()
	}
	return nil
}

// OnExit implements modes.Mode.
func (m *EditMode) OnExit() tea.Cmd {
	// Stop any ongoing analysis
	if m.analyzing {
		m.analyzer.Stop()
		m.analyzing = false
	}
	return nil
}

// Update implements modes.Mode.
func (m *EditMode) Update(msg tea.Msg) (modes.Mode, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle semantic analysis completion
	if analysisMsg, ok := msg.(SemanticAnalysisMsg); ok {
		m.analyzing = false
		m.lastAnalysis = time.Now()

		// Update context panel with analysis results
		m.contextPanel.GetPanel().SetAnalysis(analysisMsg.Analysis)

		return m, nil
	}

	// Handle file picker result
	if result, ok := msg.(overlay.FilePickerResult); ok {
		if !result.Canceled && result.FilePath != "" {
			// Open the selected file
			err := m.editor.OpenFile(result.FilePath)
			if err != nil {
				// TODO: Show error message overlay
				_ = err
			}
		}
		return m, nil
	}

	// Handle file operations keybindings (before delegating to base)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "ctrl+o":
			// Open file picker
			return m, m.showFilePicker()

		case "ctrl+s":
			// Save file
			if m.editor.GetFilePath() != "" {
				err := m.editor.SaveFile()
				if err != nil {
					// TODO: Show error message
					_ = err
				}
			} else {
				// No file path, show save as dialog
				return m, m.showFilePicker()
			}
			return m, nil

		case "ctrl+shift+s":
			// Save as - show file picker
			return m, m.showFilePicker()
		}
	}

	// Delegate to base mode for layout updates
	_, cmd := m.BaseMode.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Trigger analysis after editor changes (with debounce)
	if !m.analyzing && time.Since(m.lastAnalysis) > m.analysisDebounce {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.Type == tea.KeyRunes || keyMsg.Type == tea.KeyBackspace || keyMsg.Type == tea.KeyEnter {
				cmd := m.triggerAnalysis()
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// triggerAnalysis starts semantic analysis on the current editor content.
func (m *EditMode) triggerAnalysis() tea.Cmd {
	content := m.editor.GetContent()
	if content == "" {
		return nil
	}

	m.analyzing = true

	return func() tea.Msg {
		// Run analysis in background
		updateChan := m.analyzer.Analyze(content)

		// Consume all updates
		for range updateChan {
			// Just drain the channel
		}

		// Get final results
		analysis := m.analyzer.Results()

		return SemanticAnalysisMsg{Analysis: analysis}
	}
}

// View implements modes.Mode.
func (m *EditMode) View() string {
	return m.BaseMode.View()
}

// Keybindings implements modes.Mode.
func (m *EditMode) Keybindings() []modes.Keybinding {
	return []modes.Keybinding{
		{"Tab", "Focus next pane"},
		{"Ctrl+O", "Open file"},
		{"Ctrl+S", "Save file"},
		{"Ctrl+Shift+S", "Save as"},
		{"Ctrl+A", "Trigger analysis"},
		{"Ctrl+T", "Focus terminal"},
		{"Ctrl+E", "Focus editor"},
		{"q", "Quit mode"},
		{"?", "Show help"},
	}
}

// GetEditor returns the editor component.
func (m *EditMode) GetEditor() *EditorComponent {
	return m.editor
}

// GetContextPanel returns the context panel component.
func (m *EditMode) GetContextPanel() *ContextPanelComponent {
	return m.contextPanel
}

// GetTerminal returns the terminal component.
func (m *EditMode) GetTerminal() *TerminalComponent {
	return m.terminalComp
}

// SetAnalyzer sets a custom semantic analyzer.
func (m *EditMode) SetAnalyzer(analyzer semantic.Analyzer) {
	m.analyzer = analyzer
}

// showFilePicker returns a command to display the file picker overlay.
// This is a placeholder - actual implementation requires overlay manager integration.
func (m *EditMode) showFilePicker() tea.Cmd {
	// TODO: This requires integration with the application's overlay manager
	// For now, return nil. The actual implementation will:
	// 1. Create a FilePicker overlay
	// 2. Push it to the overlay manager
	// 3. Handle the FilePickerResult message
	return nil
}
