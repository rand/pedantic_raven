// Package analyze provides visualization components for relationship patterns.
package analyze

import (
	"fmt"
	"sort"
	"strings"
)

// PatternSortMode defines how patterns should be sorted.
type PatternSortMode int

const (
	SortByStrength PatternSortMode = iota
	SortByFrequency
	SortByConfidence
	SortByPredicate
)

// PatternFilter defines filtering criteria for pattern display.
type PatternFilter struct {
	MinOccurrences int     // Minimum occurrences to display
	MinConfidence  float64 // Minimum confidence to display
	PredicateText  string  // Filter by predicate substring (case-insensitive)
}

// PatternDisplayOptions configures pattern table rendering.
type PatternDisplayOptions struct {
	Width        int             // Table width in characters
	MaxRows      int             // Maximum rows to display (0 = all)
	SortMode     PatternSortMode // How to sort patterns
	Filter       PatternFilter   // Filtering criteria
	ShowExamples bool            // Whether to show example instances
	Compact      bool            // Compact mode (no examples)
}

// DefaultPatternDisplayOptions returns default display options.
func DefaultPatternDisplayOptions() PatternDisplayOptions {
	return PatternDisplayOptions{
		Width:        80,
		MaxRows:      20,
		SortMode:     SortByStrength,
		ShowExamples: true,
		Compact:      false,
		Filter: PatternFilter{
			MinOccurrences: 1,
			MinConfidence:  0.0,
		},
	}
}

// RenderPatternTable renders patterns as a formatted table.
func RenderPatternTable(patterns []RelationshipPattern, opts PatternDisplayOptions) string {
	if len(patterns) == 0 {
		return renderEmptyState(opts.Width)
	}

	// Apply filtering
	filtered := filterPatterns(patterns, opts.Filter)
	if len(filtered) == 0 {
		return renderNoMatches(opts.Width)
	}

	// Apply sorting
	sorted := sortPatterns(filtered, opts.SortMode)

	// Limit rows if specified
	if opts.MaxRows > 0 && len(sorted) > opts.MaxRows {
		sorted = sorted[:opts.MaxRows]
	}

	// Build table
	var b strings.Builder

	// Header
	renderHeader(&b, opts.Width)

	// Rows
	for i, pattern := range sorted {
		if opts.Compact {
			renderCompactRow(&b, pattern, i, opts.Width)
		} else {
			renderExpandedRow(&b, pattern, i, opts.Width, opts.ShowExamples)
		}
	}

	// Footer
	renderFooter(&b, opts.Width, len(filtered), len(sorted))

	return b.String()
}

// RenderClusterTable renders pattern clusters as a formatted table.
func RenderClusterTable(clusters []PatternCluster, opts PatternDisplayOptions) string {
	if len(clusters) == 0 {
		return renderEmptyState(opts.Width)
	}

	var b strings.Builder

	// Header
	b.WriteString(strings.Repeat("─", opts.Width))
	b.WriteString("\n")
	b.WriteString(centerText("Relationship Pattern Clusters", opts.Width))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", opts.Width))
	b.WriteString("\n\n")

	// Clusters
	for i, cluster := range clusters {
		renderClusterRow(&b, cluster, i, opts.Width)
		if i < len(clusters)-1 {
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", opts.Width))
	b.WriteString("\n")

	return b.String()
}

// RenderPatternStats renders pattern statistics.
func RenderPatternStats(stats PatternStats, width int) string {
	var b strings.Builder

	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")
	b.WriteString(centerText("Pattern Statistics", width))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  Total Patterns: %d\n", stats.TotalPatterns))
	b.WriteString(fmt.Sprintf("  Unique Predicates: %d\n", stats.UniquePredicates))
	b.WriteString(fmt.Sprintf("  Avg Occurrences: %.1f\n", stats.AvgOccurrences))
	b.WriteString(fmt.Sprintf("  Avg Confidence: %.2f\n", stats.AvgConfidence))

	if len(stats.TopPredicates) > 0 {
		b.WriteString("\n  Top Predicates:\n")
		for i, pred := range stats.TopPredicates {
			b.WriteString(fmt.Sprintf("    %d. %s\n", i+1, pred))
		}
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")

	return b.String()
}

// Helper functions for rendering

func renderHeader(b *strings.Builder, width int) {
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")
	b.WriteString(centerText("Relationship Patterns", width))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n\n")
}

func renderCompactRow(b *strings.Builder, pattern RelationshipPattern, index int, width int) {
	// Format: [Type1] -> predicate -> [Type2] | Count: X | Conf: X.XX | Strength: X.XX
	line := fmt.Sprintf("%2d. [%s] → %s → [%s]",
		index+1,
		pattern.SubjectType.String(),
		pattern.Predicate,
		pattern.ObjectType.String(),
	)

	stats := fmt.Sprintf("Count: %d | Conf: %.2f | Str: %.3f",
		pattern.Occurrences,
		pattern.AvgConfidence,
		pattern.Strength,
	)

	// Ensure line fits within width
	maxLineLen := width - len(stats) - 3
	if len(line) > maxLineLen {
		line = line[:maxLineLen-3] + "..."
	}

	padding := width - len(line) - len(stats) - 2
	if padding < 1 {
		padding = 1
	}

	b.WriteString(line)
	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(stats)
	b.WriteString("\n")
}

func renderExpandedRow(b *strings.Builder, pattern RelationshipPattern, index int, width int, showExamples bool) {
	// Pattern header
	b.WriteString(fmt.Sprintf("Pattern %d: [%s] → %s → [%s]\n",
		index+1,
		pattern.SubjectType.String(),
		pattern.Predicate,
		pattern.ObjectType.String(),
	))

	// Statistics
	b.WriteString(fmt.Sprintf("  Occurrences: %d  Avg Confidence: %.2f  Strength: %.3f\n",
		pattern.Occurrences,
		pattern.AvgConfidence,
		pattern.Strength,
	))

	// Examples
	if showExamples && len(pattern.Examples) > 0 {
		b.WriteString("  Examples:\n")
		for _, ex := range pattern.Examples {
			b.WriteString(fmt.Sprintf("    • %s → %s → %s (conf: %.2f)\n",
				ex.Subject,
				ex.Predicate,
				ex.Object,
				ex.Confidence,
			))
		}
	}

	b.WriteString("\n")
}

func renderClusterRow(b *strings.Builder, cluster PatternCluster, index int, width int) {
	// Cluster header
	b.WriteString(fmt.Sprintf("Cluster %d: %s (Strength: %.3f)\n",
		index+1,
		cluster.ClusterLabel,
		cluster.Strength,
	))

	// Predicates in cluster
	b.WriteString(fmt.Sprintf("  Predicates: %s\n", strings.Join(cluster.Predicates, ", ")))

	// Pattern count
	b.WriteString(fmt.Sprintf("  Patterns: %d\n", len(cluster.Patterns)))

	// Top patterns in cluster
	if len(cluster.Patterns) > 0 {
		b.WriteString("  Top Patterns:\n")
		// Show up to 3 patterns
		maxShow := min(3, len(cluster.Patterns))
		for i := 0; i < maxShow; i++ {
			p := cluster.Patterns[i]
			b.WriteString(fmt.Sprintf("    • [%s] → %s → [%s] (count: %d)\n",
				p.SubjectType.String(),
				p.Predicate,
				p.ObjectType.String(),
				p.Occurrences,
			))
		}
	}
}

func renderFooter(b *strings.Builder, width int, totalFiltered int, displayed int) {
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")

	if displayed < totalFiltered {
		msg := fmt.Sprintf("Showing %d of %d patterns", displayed, totalFiltered)
		b.WriteString(centerText(msg, width))
	} else {
		msg := fmt.Sprintf("Total: %d patterns", displayed)
		b.WriteString(centerText(msg, width))
	}

	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")
}

func renderEmptyState(width int) string {
	var b strings.Builder
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")
	b.WriteString(centerText("No Patterns Found", width))
	b.WriteString("\n")
	b.WriteString(centerText("No relationship patterns were discovered.", width))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")
	return b.String()
}

func renderNoMatches(width int) string {
	var b strings.Builder
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")
	b.WriteString(centerText("No Matching Patterns", width))
	b.WriteString("\n")
	b.WriteString(centerText("Try adjusting your filter criteria.", width))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", width))
	b.WriteString("\n")
	return b.String()
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}

	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

// filterPatterns applies filtering criteria to patterns.
func filterPatterns(patterns []RelationshipPattern, filter PatternFilter) []RelationshipPattern {
	filtered := make([]RelationshipPattern, 0)

	for _, pattern := range patterns {
		// Check occurrences
		if pattern.Occurrences < filter.MinOccurrences {
			continue
		}

		// Check confidence
		if pattern.AvgConfidence < filter.MinConfidence {
			continue
		}

		// Check predicate text filter
		if filter.PredicateText != "" {
			if !strings.Contains(
				strings.ToLower(pattern.Predicate),
				strings.ToLower(filter.PredicateText),
			) {
				continue
			}
		}

		filtered = append(filtered, pattern)
	}

	return filtered
}

// sortPatterns sorts patterns according to the specified mode.
func sortPatterns(patterns []RelationshipPattern, mode PatternSortMode) []RelationshipPattern {
	// Make a copy to avoid modifying original
	sorted := make([]RelationshipPattern, len(patterns))
	copy(sorted, patterns)

	switch mode {
	case SortByStrength:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Strength > sorted[j].Strength
		})
	case SortByFrequency:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Occurrences > sorted[j].Occurrences
		})
	case SortByConfidence:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].AvgConfidence > sorted[j].AvgConfidence
		})
	case SortByPredicate:
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Predicate < sorted[j].Predicate
		})
	}

	return sorted
}
