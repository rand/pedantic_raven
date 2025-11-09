package visualizations

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/analyze"
)

// DependencyTreeConfig configures dependency tree rendering.
type DependencyTreeConfig struct {
	Width           int  // Total width of visualization
	Height          int  // Maximum height
	ShowComplexity  bool // Show complexity scores
	ShowPriority    bool // Show priority scores
	ShowConstraints bool // Show constraints
	ColorByPriority bool // Color nodes by priority
	ExpandAll       bool // Expand all nodes by default
	MaxDepth        int  // Maximum depth to display (0 = unlimited)
}

// DefaultDependencyTreeConfig returns default configuration.
func DefaultDependencyTreeConfig() DependencyTreeConfig {
	return DependencyTreeConfig{
		Width:           80,
		Height:          40,
		ShowComplexity:  true,
		ShowPriority:    true,
		ShowConstraints: false,
		ColorByPriority: true,
		ExpandAll:       true,
		MaxDepth:        0,
	}
}

// DependencyTree renders a dependency tree for typed holes.
type DependencyTree struct {
	config   DependencyTreeConfig
	analysis *analyze.HoleAnalysis
	expanded map[string]bool // Track expanded/collapsed nodes
}

// NewDependencyTree creates a new dependency tree renderer.
func NewDependencyTree(analysis *analyze.HoleAnalysis, config DependencyTreeConfig) *DependencyTree {
	expanded := make(map[string]bool)

	// Initialize expansion state
	if config.ExpandAll {
		// Mark all nodes as expanded
		for _, hole := range analysis.Holes {
			id := fmt.Sprintf("%s_%d", hole.Type, 0)
			expanded[id] = true
		}
	}

	return &DependencyTree{
		config:   config,
		analysis: analysis,
		expanded: expanded,
	}
}

// Render generates the ASCII tree visualization.
func (dt *DependencyTree) Render() string {
	if dt.analysis == nil || dt.analysis.DependencyTree == nil {
		return "No dependency data available"
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	sb.WriteString(headerStyle.Render("Typed Hole Dependencies"))
	sb.WriteString("\n\n")

	// Render tree from root
	dt.renderNode(dt.analysis.DependencyTree, &sb, "", true, 0)

	// Summary statistics
	sb.WriteString("\n")
	sb.WriteString(dt.renderSummary())

	return sb.String()
}

// renderNode recursively renders a node and its children.
func (dt *DependencyTree) renderNode(node *analyze.HoleNode, sb *strings.Builder, prefix string, isLast bool, depth int) {
	if node == nil {
		return
	}

	// Check max depth
	if dt.config.MaxDepth > 0 && depth > dt.config.MaxDepth {
		return
	}

	// Skip root node (virtual)
	if node.ID == "root" {
		// Render all children
		for i, child := range node.Children {
			isLastChild := i == len(node.Children)-1
			dt.renderNode(child, sb, "", isLastChild, depth)
		}
		return
	}

	// Determine branch characters
	branch := "├── "
	if isLast {
		branch = "└── "
	}

	// Render current node
	sb.WriteString(prefix)
	sb.WriteString(branch)

	// Node content
	nodeStr := dt.formatNodeContent(node)
	sb.WriteString(nodeStr)
	sb.WriteString("\n")

	// Prepare prefix for children
	var childPrefix string
	if isLast {
		childPrefix = prefix + "    "
	} else {
		childPrefix = prefix + "│   "
	}

	// Render children if expanded
	if dt.isExpanded(node) && len(node.Children) > 0 {
		for i, child := range node.Children {
			isLastChild := i == len(node.Children)-1
			dt.renderNode(child, sb, childPrefix, isLastChild, depth+1)
		}
	} else if len(node.Children) > 0 {
		// Show collapsed indicator
		sb.WriteString(childPrefix)
		collapsedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)
		sb.WriteString(collapsedStyle.Render(fmt.Sprintf("[%d collapsed]\n", len(node.Children))))
	}
}

// formatNodeContent formats the display text for a node.
func (dt *DependencyTree) formatNodeContent(node *analyze.HoleNode) string {
	hole := node.Hole

	// Base node text
	nodeText := fmt.Sprintf("??%s", hole.Type)

	// Color by priority
	var style lipgloss.Style
	if dt.config.ColorByPriority {
		color := dt.getPriorityColor(hole.Priority)
		style = lipgloss.NewStyle().Foreground(color)

		// Bold if on critical path
		if node.CriticalPath {
			style = style.Bold(true)
		}
	}

	nodeText = style.Render(nodeText)

	// Add metadata
	metadata := []string{}

	if dt.config.ShowPriority {
		metadata = append(metadata, fmt.Sprintf("P:%d", hole.Priority))
	}

	if dt.config.ShowComplexity {
		metadata = append(metadata, fmt.Sprintf("C:%d", hole.Complexity))
	}

	if dt.config.ShowConstraints && hole.Constraint != "" {
		constraint := hole.Constraint
		if len(constraint) > 20 {
			constraint = constraint[:17] + "..."
		}
		metadata = append(metadata, fmt.Sprintf("[%s]", constraint))
	}

	// Critical path indicator
	if node.CriticalPath {
		criticalStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
		metadata = append(metadata, criticalStyle.Render("CRITICAL"))
	}

	if len(metadata) > 0 {
		metaStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))
		nodeText += " " + metaStyle.Render(fmt.Sprintf("(%s)", strings.Join(metadata, ", ")))
	}

	return nodeText
}

// getPriorityColor returns color based on priority level.
func (dt *DependencyTree) getPriorityColor(priority int) lipgloss.Color {
	switch {
	case priority >= 8:
		return lipgloss.Color("196") // Red - High priority
	case priority >= 5:
		return lipgloss.Color("226") // Yellow - Medium priority
	default:
		return lipgloss.Color("34") // Green - Low priority
	}
}

// isExpanded checks if a node is expanded.
func (dt *DependencyTree) isExpanded(node *analyze.HoleNode) bool {
	if dt.config.ExpandAll {
		return true
	}
	return dt.expanded[node.ID]
}

// ToggleExpand toggles the expansion state of a node.
func (dt *DependencyTree) ToggleExpand(nodeID string) {
	dt.expanded[nodeID] = !dt.expanded[nodeID]
}

// ExpandAll expands all nodes.
func (dt *DependencyTree) ExpandAll() {
	dt.config.ExpandAll = true
	for _, hole := range dt.analysis.Holes {
		id := fmt.Sprintf("%s_%d", hole.Type, 0)
		dt.expanded[id] = true
	}
}

// CollapseAll collapses all nodes.
func (dt *DependencyTree) CollapseAll() {
	dt.config.ExpandAll = false
	for id := range dt.expanded {
		dt.expanded[id] = false
	}
}

// renderSummary generates summary statistics.
func (dt *DependencyTree) renderSummary() string {
	var sb strings.Builder

	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true)

	sb.WriteString(summaryStyle.Render("Summary:"))
	sb.WriteString("\n")

	// Total holes
	sb.WriteString(fmt.Sprintf("  Total Holes: %d\n", len(dt.analysis.Holes)))

	// Total complexity
	sb.WriteString(fmt.Sprintf("  Total Complexity: %d\n", dt.analysis.TotalComplexity))

	// Average priority
	sb.WriteString(fmt.Sprintf("  Average Priority: %.1f\n", dt.analysis.AveragePriority))

	// Critical path
	if len(dt.analysis.CriticalPath) > 0 {
		criticalComplexity := 0
		for _, hole := range dt.analysis.CriticalPath {
			criticalComplexity += hole.Complexity
		}

		criticalStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

		sb.WriteString(fmt.Sprintf("  %s: %d holes, complexity %d\n",
			criticalStyle.Render("Critical Path"),
			len(dt.analysis.CriticalPath),
			criticalComplexity))
	}

	// Circular dependencies warning
	if len(dt.analysis.CircularDeps) > 0 {
		warningStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("202")).
			Bold(true)

		sb.WriteString(fmt.Sprintf("\n  %s: %d circular dependencies detected!\n",
			warningStyle.Render("WARNING"),
			len(dt.analysis.CircularDeps)))
	}

	return sb.String()
}

// RenderCompact generates a compact horizontal tree view.
func (dt *DependencyTree) RenderCompact() string {
	if dt.analysis == nil || len(dt.analysis.ImplementOrder) == 0 {
		return "No holes to display"
	}

	var sb strings.Builder

	// Sort holes by implementation order
	for i, hole := range dt.analysis.ImplementOrder {
		// Arrow separator
		if i > 0 {
			sb.WriteString(" → ")
		}

		// Color by priority
		color := dt.getPriorityColor(hole.Priority)
		style := lipgloss.NewStyle().Foreground(color)

		// Hole name
		sb.WriteString(style.Render(fmt.Sprintf("??%s", hole.Type)))

		// Complexity in parentheses
		metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
		sb.WriteString(metaStyle.Render(fmt.Sprintf("(%d)", hole.Complexity)))
	}

	return sb.String()
}

// RenderDependencyMatrix generates a dependency matrix visualization.
func RenderDependencyMatrix(analysis *analyze.HoleAnalysis) string {
	if analysis == nil || len(analysis.Holes) == 0 {
		return "No dependency data available"
	}

	var sb strings.Builder

	// Create matrix
	n := len(analysis.Holes)
	matrix := make([][]bool, n)
	for i := range matrix {
		matrix[i] = make([]bool, n)
	}

	// Build ID to index map
	idxMap := make(map[string]int)
	for i, hole := range analysis.Holes {
		id := fmt.Sprintf("%s_%d", hole.Type, i)
		idxMap[id] = i
	}

	// Fill matrix from dependencies
	for _, dep := range analysis.Dependencies {
		if dep.Relationship == "requires" || dep.Relationship == "extends" || dep.Relationship == "implements" {
			if fromIdx, ok := idxMap[dep.From]; ok {
				if toIdx, ok := idxMap[dep.To]; ok {
					matrix[fromIdx][toIdx] = true
				}
			}
		}
	}

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	sb.WriteString(headerStyle.Render("Dependency Matrix"))
	sb.WriteString("\n\n")

	// Column headers
	sb.WriteString("     ")
	for i := 0; i < n && i < 10; i++ {
		sb.WriteString(fmt.Sprintf(" %d ", i))
	}
	sb.WriteString("\n")

	// Matrix rows
	for i := 0; i < n && i < 10; i++ {
		sb.WriteString(fmt.Sprintf("%3d: ", i))

		for j := 0; j < n && j < 10; j++ {
			if matrix[i][j] {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
				sb.WriteString(style.Render(" ■ "))
			} else {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				sb.WriteString(style.Render(" · "))
			}
		}
		sb.WriteString("\n")
	}

	// Legend
	sb.WriteString("\n")
	legendStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	sb.WriteString(legendStyle.Render("Legend: ■ = dependency exists, · = no dependency"))
	sb.WriteString("\n")

	// Hole index reference
	sb.WriteString("\nHoles:\n")
	for i := 0; i < n && i < 10; i++ {
		hole := analysis.Holes[i]
		sb.WriteString(fmt.Sprintf("  %d: ??%s\n", i, hole.Type))
	}

	if n > 10 {
		sb.WriteString(fmt.Sprintf("  ... and %d more\n", n-10))
	}

	return sb.String()
}

// RenderCircularDependencies highlights circular dependency chains.
func RenderCircularDependencies(analysis *analyze.HoleAnalysis) string {
	if analysis == nil || len(analysis.CircularDeps) == 0 {
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")).
			Bold(true)
		return successStyle.Render("✓ No circular dependencies detected")
	}

	var sb strings.Builder

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("202")).
		Bold(true)

	sb.WriteString(warningStyle.Render(fmt.Sprintf("⚠ Warning: %d Circular Dependencies Detected", len(analysis.CircularDeps))))
	sb.WriteString("\n\n")

	for i, cycle := range analysis.CircularDeps {
		sb.WriteString(fmt.Sprintf("Cycle %d:\n", i+1))

		// Render cycle with arrows
		cycleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

		for j, nodeID := range cycle {
			if j > 0 {
				sb.WriteString(" → ")
			}
			// Extract type from ID
			parts := strings.Split(nodeID, "_")
			if len(parts) > 0 {
				sb.WriteString(cycleStyle.Render(fmt.Sprintf("??%s", parts[0])))
			}
		}

		// Close the cycle
		if len(cycle) > 0 {
			sb.WriteString(" → ")
			parts := strings.Split(cycle[0], "_")
			if len(parts) > 0 {
				sb.WriteString(cycleStyle.Render(fmt.Sprintf("??%s", parts[0])))
			}
		}

		sb.WriteString("\n\n")
	}

	recommendStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true)

	sb.WriteString(recommendStyle.Render("Recommendation: Break circular dependencies by introducing abstractions or interfaces"))
	sb.WriteString("\n")

	return sb.String()
}
