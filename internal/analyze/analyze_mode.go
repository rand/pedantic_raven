package analyze

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// ViewMode represents the different analysis views.
type ViewMode int

const (
	ViewTripleGraph ViewMode = iota
	ViewEntityFrequency
	ViewPatterns
	ViewTypedHoles
)

// AnalyzeMode coordinates all analysis views.
type AnalyzeMode struct {
	// Current view mode
	currentView ViewMode

	// Shared analysis data
	analysis *semantic.Analysis

	// View components
	tripleGraphView Model // The existing triple graph model

	// Analysis results (cached)
	entityFreqs    []EntityFrequency
	patterns       []RelationshipPattern
	holeAnalysis   *HoleAnalysis

	// UI state
	width  int
	height int
	err    error

	// Export state
	exportFormat ExportFormat
	lastExport   time.Time
}

// ExportFormat represents export file format.
type ExportFormat int

const (
	ExportMarkdown ExportFormat = iota
	ExportHTML
	ExportPDF
)

// NewAnalyzeMode creates a new analyze mode coordinator.
func NewAnalyzeMode() *AnalyzeMode {
	return &AnalyzeMode{
		currentView:     ViewTripleGraph,
		tripleGraphView: NewModel(),
		exportFormat:    ExportMarkdown,
	}
}

// Init implements tea.Model.
func (m *AnalyzeMode) Init() tea.Cmd {
	return m.tripleGraphView.Init()
}

// SetAnalysis sets the analysis data and updates all views.
func (m *AnalyzeMode) SetAnalysis(analysis *semantic.Analysis) {
	m.analysis = analysis

	// Update triple graph view
	m.tripleGraphView.SetAnalysis(analysis)

	// Cache analysis results
	m.entityFreqs = CalculateEntityFrequency(analysis)
	m.patterns = MinePatterns(analysis)
	m.holeAnalysis = AnalyzeTypedHoles(analysis)
}

// SetSize sets the component size.
func (m *AnalyzeMode) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.tripleGraphView.SetSize(width, height-4) // Reserve space for header/footer
}

// SwitchView changes the current view mode.
func (m *AnalyzeMode) SwitchView(view ViewMode) {
	m.currentView = view
}

// NextView cycles to the next view.
func (m *AnalyzeMode) NextView() {
	m.currentView = (m.currentView + 1) % 4
}

// PrevView cycles to the previous view.
func (m *AnalyzeMode) PrevView() {
	m.currentView = (m.currentView + 3) % 4 // +3 mod 4 = -1 mod 4
}

// Update implements tea.Model.
func (m *AnalyzeMode) Update(msg tea.Msg) (*AnalyzeMode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case ExportCompleteMsg:
		m.lastExport = time.Now()
		return m, nil

	case ExportErrorMsg:
		m.err = msg.Err
		return m, nil
	}

	// Forward messages to active view
	if m.currentView == ViewTripleGraph {
		updated, cmd := m.tripleGraphView.Update(msg)
		m.tripleGraphView = updated
		return m, cmd
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m *AnalyzeMode) handleKeyPress(msg tea.KeyMsg) (*AnalyzeMode, tea.Cmd) {
	switch msg.String() {
	// View switching
	case "tab":
		m.NextView()
		return m, nil

	case "shift+tab":
		m.PrevView()
		return m, nil

	case "1":
		m.SwitchView(ViewTripleGraph)
		return m, nil

	case "2":
		m.SwitchView(ViewEntityFrequency)
		return m, nil

	case "3":
		m.SwitchView(ViewPatterns)
		return m, nil

	case "4":
		m.SwitchView(ViewTypedHoles)
		return m, nil

	// Export commands
	case "e":
		return m, m.exportReport()

	case "ctrl+e":
		// Cycle export format
		m.exportFormat = (m.exportFormat + 1) % 3
		return m, nil

	// Help
	case "?":
		// TODO: Show help overlay
		return m, nil
	}

	// Forward to active view
	if m.currentView == ViewTripleGraph {
		updated, cmd := m.tripleGraphView.Update(msg)
		m.tripleGraphView = updated
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model.
func (m *AnalyzeMode) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.analysis == nil {
		return helpStyle.Render("No analysis data available")
	}

	// Render header with tabs
	header := m.renderHeader()

	// Render active view
	var content string
	switch m.currentView {
	case ViewTripleGraph:
		content = m.tripleGraphView.View()
	case ViewEntityFrequency:
		content = m.renderEntityFrequency()
	case ViewPatterns:
		content = m.renderPatterns()
	case ViewTypedHoles:
		content = m.renderTypedHoles()
	}

	// Render footer with shortcuts
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		content,
		footer,
	)
}

// renderHeader renders the tab navigation header.
func (m *AnalyzeMode) renderHeader() string {
	tabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(lipgloss.Color("244"))

	activeTabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("62")).
		Bold(true)

	tabs := []string{
		"[1] Triple Graph",
		"[2] Entity Frequency",
		"[3] Patterns",
		"[4] Typed Holes",
	}

	var renderedTabs []string
	for i, tab := range tabs {
		if ViewMode(i) == m.currentView {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, tabStyle.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// renderFooter renders the keyboard shortcuts footer.
func (m *AnalyzeMode) renderFooter() string {
	var shortcuts string

	exportFmt := "MD"
	switch m.exportFormat {
	case ExportHTML:
		exportFmt = "HTML"
	case ExportPDF:
		exportFmt = "PDF"
	}

	if m.lastExport.IsZero() {
		shortcuts = fmt.Sprintf("[tab] Switch  [e] Export (%s)  [ctrl+e] Change format  [?] Help", exportFmt)
	} else {
		elapsed := time.Since(m.lastExport)
		shortcuts = fmt.Sprintf("[tab] Switch  [e] Export (%s)  [ctrl+e] Change format  Last: %s ago", exportFmt, elapsed.Round(time.Second))
	}

	return helpStyle.Render(shortcuts)
}

// renderEntityFrequency renders the entity frequency analysis view.
func (m *AnalyzeMode) renderEntityFrequency() string {
	if len(m.entityFreqs) == 0 {
		return helpStyle.Render("No entities found")
	}

	// Sort by frequency
	sorted := make([]EntityFrequency, len(m.entityFreqs))
	copy(sorted, m.entityFreqs)
	FrequencyList(sorted).SortByFrequency()

	// Render top entities as bar chart
	topN := 15
	if len(sorted) < topN {
		topN = len(sorted)
	}

	var lines []string
	lines = append(lines, "\nTop Entities by Frequency:\n")

	// Simple text-based rendering (bar chart visualization is complex, so keep it simple)
	maxCount := sorted[0].Count
	if maxCount == 0 {
		lines = append(lines, "No entity data available")
		return lipgloss.JoinVertical(lipgloss.Left, lines...)
	}

	barWidth := m.width - 35 // Leave room for label and count
	if barWidth > 50 {
		barWidth = 50
	}

	for _, freq := range sorted[:topN] {
		// Calculate bar length
		length := (freq.Count * barWidth) / maxCount
		if length < 1 && freq.Count > 0 {
			length = 1
		}

		// Create bar
		bar := strings.Repeat("█", length) + strings.Repeat("░", barWidth-length)
		label := fmt.Sprintf("%-20s %s %4d", truncate(freq.Text, 20), bar, freq.Count)
		lines = append(lines, label)
	}

	// Add type distribution
	lines = append(lines, "\n\nEntity Type Distribution:\n")
	typeCounts := make(map[semantic.EntityType]int)
	for _, freq := range m.entityFreqs {
		typeCounts[freq.Type]++
	}

	// Sort types by count
	var types []semantic.EntityType
	for t := range typeCounts {
		types = append(types, t)
	}
	// Simple sort by count (descending)
	for i := 0; i < len(types); i++ {
		for j := i + 1; j < len(types); j++ {
			if typeCounts[types[i]] < typeCounts[types[j]] {
				types[i], types[j] = types[j], types[i]
			}
		}
	}

	for _, entityType := range types {
		count := typeCounts[entityType]
		label := fmt.Sprintf("%-15s: %4d", entityType, count)
		lines = append(lines, label)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderPatterns renders the relationship patterns view.
func (m *AnalyzeMode) renderPatterns() string {
	if len(m.patterns) == 0 {
		return helpStyle.Render("No patterns discovered")
	}

	// Render pattern table
	opts := PatternDisplayOptions{
		Width:        m.width,
		MaxRows:      m.height - 8,
		SortMode:     SortByStrength,
		ShowExamples: true,
		Compact:      false,
	}
	return RenderPatternTable(m.patterns, opts)
}

// renderTypedHoles renders the typed holes prioritization view.
func (m *AnalyzeMode) renderTypedHoles() string {
	if m.holeAnalysis == nil || len(m.holeAnalysis.Holes) == 0 {
		return helpStyle.Render("No typed holes found")
	}

	// Render priority table (simple inline rendering to avoid import cycle)
	var lines []string
	lines = append(lines, "\nTyped Holes Priority List:\n")

	// Sort holes by priority (descending)
	holes := make([]semantic.EnhancedTypedHole, len(m.holeAnalysis.ImplementOrder))
	copy(holes, m.holeAnalysis.ImplementOrder)
	sort.Slice(holes, func(i, j int) bool {
		return holes[i].Priority > holes[j].Priority
	})

	// Render top holes
	topN := 20
	if len(holes) < topN {
		topN = len(holes)
	}

	for i, hole := range holes[:topN] {
		priority := strings.Repeat("█", hole.Priority) + strings.Repeat("░", 10-hole.Priority)
		name := truncate(hole.Type, 30) // TypedHole has Type field, not Name
		complexity := hole.Complexity
		line := fmt.Sprintf("%2d. [%s] %-30s (Complexity: %d)", i+1, priority, name, complexity)
		lines = append(lines, line)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Add circular dependency warnings
	if len(m.holeAnalysis.CircularDeps) > 0 {
		warning := fmt.Sprintf("\n⚠️  Warning: %d circular dependencies detected", len(m.holeAnalysis.CircularDeps))
		content += errorStyle.Render(warning)
	}

	return content
}

// exportReport generates and exports the analysis report.
func (m *AnalyzeMode) exportReport() tea.Cmd {
	return func() tea.Msg {
		// TODO: Wire up actual export functionality
		// For now, this is a stub that will be completed during integration
		// The export package already has all the export functions (ExportMarkdown, ExportHTML, ExportPDF)
		// but we can't import it here due to circular dependency.
		// Solution: Create an export helper function in export package that takes analyze data.

		var filename string
		switch m.exportFormat {
		case ExportMarkdown:
			filename = fmt.Sprintf("analysis_%s.md", time.Now().Format("20060102_150405"))
		case ExportHTML:
			filename = fmt.Sprintf("analysis_%s.html", time.Now().Format("20060102_150405"))
		case ExportPDF:
			filename = fmt.Sprintf("analysis_%s.pdf", time.Now().Format("20060102_150405"))
		}

		// In real implementation, would call export functions here
		// For now, return success to allow UI testing
		return ExportCompleteMsg{Filename: filename}
	}
}

// Messages for export operations.
type (
	ExportCompleteMsg struct {
		Filename string
	}

	ExportErrorMsg struct {
		Err error
	}
)

// Helper functions.

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
