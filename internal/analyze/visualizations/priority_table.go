package visualizations

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// PriorityTableConfig configures priority table rendering.
type PriorityTableConfig struct {
	Width           int    // Total width of table
	MaxRows         int    // Maximum rows to display
	ShowConstraints bool   // Show constraint column
	ShowSuggestions bool   // Show suggestion column
	SortBy          string // Sort column: "priority", "complexity", "recommended"
	FilterType      string // Filter by type (empty = all)
	MinPriority     int    // Minimum priority to display
}

// DefaultPriorityTableConfig returns default configuration.
func DefaultPriorityTableConfig() PriorityTableConfig {
	return PriorityTableConfig{
		Width:           100,
		MaxRows:         20,
		ShowConstraints: true,
		ShowSuggestions: true,
		SortBy:          "recommended",
		FilterType:      "",
		MinPriority:     0,
	}
}

// PriorityTable renders a table of typed holes with priorities.
type PriorityTable struct {
	config   PriorityTableConfig
	analysis *analyze.HoleAnalysis
}

// NewPriorityTable creates a new priority table renderer.
func NewPriorityTable(analysis *analyze.HoleAnalysis, config PriorityTableConfig) *PriorityTable {
	return &PriorityTable{
		config:   config,
		analysis: analysis,
	}
}

// Render generates the priority table visualization.
func (pt *PriorityTable) Render() string {
	if pt.analysis == nil || len(pt.analysis.Holes) == 0 {
		return "No typed holes to display"
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	sb.WriteString(headerStyle.Render("Typed Hole Priority Queue"))
	sb.WriteString("\n\n")

	// Filter and sort holes
	holes := pt.filterAndSort()

	// Limit rows
	if len(holes) > pt.config.MaxRows {
		holes = holes[:pt.config.MaxRows]
	}

	// Render table
	sb.WriteString(pt.renderTable(holes))

	// Footer with statistics
	sb.WriteString("\n")
	sb.WriteString(pt.renderFooter(len(holes)))

	return sb.String()
}

// filterAndSort filters and sorts holes based on configuration.
func (pt *PriorityTable) filterAndSort() []semantic.EnhancedTypedHole {
	filtered := []semantic.EnhancedTypedHole{}

	for _, hole := range pt.analysis.Holes {
		// Apply filters
		if pt.config.FilterType != "" && hole.Type != pt.config.FilterType {
			continue
		}

		if hole.Priority < pt.config.MinPriority {
			continue
		}

		filtered = append(filtered, hole)
	}

	// Sort based on configuration
	switch pt.config.SortBy {
	case "priority":
		sort.Slice(filtered, func(i, j int) bool {
			if filtered[i].Priority != filtered[j].Priority {
				return filtered[i].Priority > filtered[j].Priority
			}
			return filtered[i].Complexity < filtered[j].Complexity
		})

	case "complexity":
		sort.Slice(filtered, func(i, j int) bool {
			if filtered[i].Complexity != filtered[j].Complexity {
				return filtered[i].Complexity < filtered[j].Complexity
			}
			return filtered[i].Priority > filtered[j].Priority
		})

	case "recommended":
		// Sort by priority/complexity ratio
		sort.Slice(filtered, func(i, j int) bool {
			scoreI := float64(filtered[i].Priority) / float64(max(filtered[i].Complexity, 1))
			scoreJ := float64(filtered[j].Priority) / float64(max(filtered[j].Complexity, 1))
			return scoreI > scoreJ
		})

	default:
		// Default to recommended order
		sort.Slice(filtered, func(i, j int) bool {
			scoreI := float64(filtered[i].Priority) / float64(max(filtered[i].Complexity, 1))
			scoreJ := float64(filtered[j].Priority) / float64(max(filtered[j].Complexity, 1))
			return scoreI > scoreJ
		})
	}

	return filtered
}

// renderTable renders the table with holes.
func (pt *PriorityTable) renderTable(holes []semantic.EnhancedTypedHole) string {
	var sb strings.Builder

	// Column widths
	rankWidth := 4
	typeWidth := 20
	priorityWidth := 10
	complexityWidth := 12
	constraintWidth := 20
	suggestionWidth := 30

	// Adjust based on config
	if !pt.config.ShowConstraints {
		constraintWidth = 0
	}
	if !pt.config.ShowSuggestions {
		suggestionWidth = 0
	}

	// Table header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("244")).
		Underline(true)

	header := []string{
		pt.padRight("#", rankWidth),
		pt.padRight("Type", typeWidth),
		pt.padRight("Priority", priorityWidth),
		pt.padRight("Complexity", complexityWidth),
	}

	if pt.config.ShowConstraints {
		header = append(header, pt.padRight("Constraint", constraintWidth))
	}

	if pt.config.ShowSuggestions {
		header = append(header, pt.padRight("Suggestion", suggestionWidth))
	}

	sb.WriteString(headerStyle.Render(strings.Join(header, " ")))
	sb.WriteString("\n")

	// Separator line
	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	totalWidth := rankWidth + typeWidth + priorityWidth + complexityWidth
	if pt.config.ShowConstraints {
		totalWidth += constraintWidth + 1
	}
	if pt.config.ShowSuggestions {
		totalWidth += suggestionWidth + 1
	}

	sb.WriteString(separatorStyle.Render(strings.Repeat("─", totalWidth)))
	sb.WriteString("\n")

	// Table rows
	for i, hole := range holes {
		sb.WriteString(pt.renderRow(i+1, hole, rankWidth, typeWidth, priorityWidth, complexityWidth, constraintWidth, suggestionWidth))
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderRow renders a single table row.
func (pt *PriorityTable) renderRow(rank int, hole semantic.EnhancedTypedHole, rankW, typeW, priorityW, complexityW, constraintW, suggestionW int) string {
	row := []string{}

	// Rank
	rankStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))
	row = append(row, rankStyle.Render(pt.padRight(fmt.Sprintf("%d.", rank), rankW)))

	// Type (colored by priority)
	typeColor := pt.getPriorityColor(hole.Priority)
	typeStyle := lipgloss.NewStyle().
		Foreground(typeColor).
		Bold(true)
	typeText := fmt.Sprintf("??%s", hole.Type)
	if len(typeText) > typeW {
		typeText = typeText[:typeW-3] + "..."
	}
	row = append(row, typeStyle.Render(pt.padRight(typeText, typeW)))

	// Priority
	priorityStyle := lipgloss.NewStyle().
		Foreground(typeColor)
	priorityBar := pt.renderPriorityBar(hole.Priority)
	row = append(row, priorityStyle.Render(pt.padRight(priorityBar, priorityW)))

	// Complexity
	complexityStyle := lipgloss.NewStyle().
		Foreground(pt.getComplexityColor(hole.Complexity))
	complexityBar := pt.renderComplexityBar(hole.Complexity)
	row = append(row, complexityStyle.Render(pt.padRight(complexityBar, complexityW)))

	// Constraint
	if pt.config.ShowConstraints {
		constraintText := hole.Constraint
		if constraintText == "" {
			constraintText = "-"
		}
		if len(constraintText) > constraintW {
			constraintText = constraintText[:constraintW-3] + "..."
		}
		constraintStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
		row = append(row, constraintStyle.Render(pt.padRight(constraintText, constraintW)))
	}

	// Suggestion
	if pt.config.ShowSuggestions {
		suggestion := hole.SuggestedImpl
		if suggestion == "" {
			suggestion = "-"
		}
		if len(suggestion) > suggestionW {
			suggestion = suggestion[:suggestionW-3] + "..."
		}
		suggestionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)
		row = append(row, suggestionStyle.Render(pt.padRight(suggestion, suggestionW)))
	}

	return strings.Join(row, " ")
}

// renderPriorityBar creates a visual bar for priority.
func (pt *PriorityTable) renderPriorityBar(priority int) string {
	maxBars := 10
	bars := (priority * maxBars) / 10
	if bars > maxBars {
		bars = maxBars
	}

	return fmt.Sprintf("%s %d/10", strings.Repeat("█", bars)+strings.Repeat("░", maxBars-bars), priority)
}

// renderComplexityBar creates a visual bar for complexity.
func (pt *PriorityTable) renderComplexityBar(complexity int) string {
	maxBars := 10
	bars := (complexity * maxBars) / 10
	if bars > maxBars {
		bars = maxBars
	}

	return fmt.Sprintf("%s %d/10", strings.Repeat("█", bars)+strings.Repeat("░", maxBars-bars), complexity)
}

// getPriorityColor returns color based on priority level.
func (pt *PriorityTable) getPriorityColor(priority int) lipgloss.Color {
	switch {
	case priority >= 8:
		return lipgloss.Color("196") // Red - High priority
	case priority >= 5:
		return lipgloss.Color("226") // Yellow - Medium priority
	default:
		return lipgloss.Color("34") // Green - Low priority
	}
}

// getComplexityColor returns color based on complexity level.
func (pt *PriorityTable) getComplexityColor(complexity int) lipgloss.Color {
	switch {
	case complexity >= 8:
		return lipgloss.Color("196") // Red - High complexity
	case complexity >= 5:
		return lipgloss.Color("226") // Yellow - Medium complexity
	default:
		return lipgloss.Color("34") // Green - Low complexity
	}
}

// padRight pads a string to the right.
func (pt *PriorityTable) padRight(s string, width int) string {
	// Remove ANSI codes for length calculation
	visualLen := len(stripAnsi(s))
	if visualLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visualLen)
}

// stripAnsi removes ANSI escape codes for length calculation.
func stripAnsi(s string) string {
	// Simple implementation - remove common ANSI patterns
	result := s
	inEscape := false
	var sb strings.Builder

	for _, r := range result {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		sb.WriteRune(r)
	}

	return sb.String()
}

// renderFooter generates footer statistics.
func (pt *PriorityTable) renderFooter(displayCount int) string {
	var sb strings.Builder

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true)

	sb.WriteString(footerStyle.Render(fmt.Sprintf("Displaying %d of %d typed holes", displayCount, len(pt.analysis.Holes))))
	sb.WriteString("\n")

	// Sort info
	sortLabel := map[string]string{
		"priority":    "Priority (High to Low)",
		"complexity":  "Complexity (Low to High)",
		"recommended": "Recommended Order (Priority/Complexity Ratio)",
	}

	if label, ok := sortLabel[pt.config.SortBy]; ok {
		sb.WriteString(footerStyle.Render(fmt.Sprintf("Sorted by: %s", label)))
		sb.WriteString("\n")
	}

	return sb.String()
}

// RenderCompact generates a compact single-line summary.
func (pt *PriorityTable) RenderCompact() string {
	if pt.analysis == nil || len(pt.analysis.Holes) == 0 {
		return "No typed holes"
	}

	// Count by priority
	highPriority := 0   // 8-10
	mediumPriority := 0 // 5-7
	lowPriority := 0    // 0-4

	for _, hole := range pt.analysis.Holes {
		if hole.Priority >= 8 {
			highPriority++
		} else if hole.Priority >= 5 {
			mediumPriority++
		} else {
			lowPriority++
		}
	}

	highStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	mediumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	lowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("34"))

	return fmt.Sprintf("Holes: %s high, %s medium, %s low | Total complexity: %d",
		highStyle.Render(fmt.Sprintf("%d", highPriority)),
		mediumStyle.Render(fmt.Sprintf("%d", mediumPriority)),
		lowStyle.Render(fmt.Sprintf("%d", lowPriority)),
		pt.analysis.TotalComplexity)
}

// RenderByType generates a table grouped by hole type.
func RenderByType(analysis *analyze.HoleAnalysis) string {
	if analysis == nil || len(analysis.Holes) == 0 {
		return "No typed holes to display"
	}

	var sb strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	sb.WriteString(headerStyle.Render("Typed Holes Grouped by Type"))
	sb.WriteString("\n\n")

	// Group by type
	typeGroups := make(map[string][]semantic.EnhancedTypedHole)
	for _, hole := range analysis.Holes {
		typeGroups[hole.Type] = append(typeGroups[hole.Type], hole)
	}

	// Sort types by total priority
	types := make([]string, 0, len(typeGroups))
	for t := range typeGroups {
		types = append(types, t)
	}

	sort.Slice(types, func(i, j int) bool {
		totalI := 0
		totalJ := 0
		for _, h := range typeGroups[types[i]] {
			totalI += h.Priority
		}
		for _, h := range typeGroups[types[j]] {
			totalJ += h.Priority
		}
		return totalI > totalJ
	})

	// Render each group
	for _, t := range types {
		holes := typeGroups[t]

		typeStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226"))

		sb.WriteString(typeStyle.Render(fmt.Sprintf("??%s (%d holes)", t, len(holes))))
		sb.WriteString("\n")

		// Sort holes in group by priority
		sort.Slice(holes, func(i, j int) bool {
			return holes[i].Priority > holes[j].Priority
		})

		// Show top 3 holes in each group
		for i, hole := range holes {
			if i >= 3 {
				sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(holes)-3))
				break
			}

			priorityColor := getPriorityColorFunc(hole.Priority)
			priorityStyle := lipgloss.NewStyle().Foreground(priorityColor)

			sb.WriteString(fmt.Sprintf("  • Priority: %s, Complexity: %d\n",
				priorityStyle.Render(fmt.Sprintf("%d", hole.Priority)),
				hole.Complexity))
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// getPriorityColorFunc returns color based on priority (standalone function).
func getPriorityColorFunc(priority int) lipgloss.Color {
	switch {
	case priority >= 8:
		return lipgloss.Color("196")
	case priority >= 5:
		return lipgloss.Color("226")
	default:
		return lipgloss.Color("34")
	}
}
