package export

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// ExportMarkdown generates a GitHub-flavored markdown report.
func ExportMarkdown(report *AnalysisReport) (string, error) {
	if report == nil {
		return "", fmt.Errorf("report cannot be nil")
	}

	opts := DefaultExportOptions(FormatMarkdown)
	return ExportMarkdownWithOptions(report, opts)
}

// ExportMarkdownWithOptions generates markdown with custom options.
func ExportMarkdownWithOptions(report *AnalysisReport, opts ExportOptions) (string, error) {
	if report == nil {
		return "", fmt.Errorf("report cannot be nil")
	}

	var sb strings.Builder

	// Title and metadata
	if opts.IncludeMetadata {
		sb.WriteString(fmt.Sprintf("# %s\n\n", report.Title))
		sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", report.Timestamp.Format(time.RFC3339)))
		if report.Source != "" {
			sb.WriteString(fmt.Sprintf("**Source:** `%s`\n\n", report.Source))
		}
		if report.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", report.Description))
		}
		sb.WriteString("---\n\n")
	}

	// Table of contents
	sb.WriteString("## Table of Contents\n\n")
	if opts.IncludeStatistics {
		sb.WriteString("- [Overview Statistics](#overview-statistics)\n")
	}
	if opts.IncludeFrequencies {
		sb.WriteString("- [Entity Frequency Analysis](#entity-frequency-analysis)\n")
	}
	if opts.IncludePatterns {
		sb.WriteString("- [Relationship Patterns](#relationship-patterns)\n")
	}
	if opts.IncludeTypedHoles {
		sb.WriteString("- [Typed Holes](#typed-holes)\n")
	}
	if opts.IncludeTripleGraph {
		sb.WriteString("- [Triple Graph Visualization](#triple-graph-visualization)\n")
	}
	sb.WriteString("\n---\n\n")

	// Statistics section
	if opts.IncludeStatistics {
		sb.WriteString(formatStatistics(report))
	}

	// Entity frequency section
	if opts.IncludeFrequencies && len(report.EntityFrequencies) > 0 {
		sb.WriteString(formatEntityFrequencies(report, opts.MaxFrequenciesToShow))
	}

	// Relationship patterns section
	if opts.IncludePatterns && len(report.RelationshipPatterns) > 0 {
		sb.WriteString(formatRelationshipPatterns(report, opts.MaxExamplesPerPattern))
	}

	// Typed holes section
	if opts.IncludeTypedHoles && len(report.TypedHoles) > 0 {
		sb.WriteString(formatTypedHoles(report))
	}

	// Triple graph section
	if opts.IncludeTripleGraph && report.TripleGraph != nil {
		sb.WriteString(formatTripleGraph(report))
	}

	return sb.String(), nil
}

// formatStatistics formats the statistics section.
func formatStatistics(report *AnalysisReport) string {
	var sb strings.Builder

	sb.WriteString("## Overview Statistics\n\n")
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| **Total Entities** | %d |\n", report.Stats.TotalEntities))
	sb.WriteString(fmt.Sprintf("| **Total Relationships** | %d |\n", report.Stats.TotalRelationships))
	sb.WriteString(fmt.Sprintf("| **Total Typed Holes** | %d |\n", report.Stats.TotalTypedHoles))
	sb.WriteString(fmt.Sprintf("| **Unique Entity Types** | %d |\n", report.Stats.UniqueEntityTypes))
	sb.WriteString(fmt.Sprintf("| **Unique Patterns** | %d |\n", report.Stats.UniquePatterns))

	if report.Stats.MostCommonEntity != "" {
		sb.WriteString(fmt.Sprintf("| **Most Common Entity** | `%s` (%d occurrences) |\n",
			report.Stats.MostCommonEntity, report.Stats.MostCommonEntityCount))
	}

	if report.Stats.StrongestPattern != "" {
		sb.WriteString(fmt.Sprintf("| **Strongest Pattern** | %s (%.2f) |\n",
			report.Stats.StrongestPattern, report.Stats.StrongestPatternScore))
	}

	if report.Stats.HighestPriorityHole != "" {
		sb.WriteString(fmt.Sprintf("| **Highest Priority Hole** | `%s` (priority %d) |\n",
			report.Stats.HighestPriorityHole, report.Stats.HighestPriority))
	}

	if report.Stats.TotalTypedHoles > 0 {
		sb.WriteString(fmt.Sprintf("| **Average Hole Complexity** | %.1f |\n", report.Stats.AvgComplexity))
	}

	sb.WriteString("\n")
	return sb.String()
}

// formatEntityFrequencies formats the entity frequency section.
func formatEntityFrequencies(report *AnalysisReport, maxShow int) string {
	var sb strings.Builder

	sb.WriteString("## Entity Frequency Analysis\n\n")
	sb.WriteString("### Top Entities by Frequency\n\n")

	// Limit display
	frequencies := report.EntityFrequencies
	if len(frequencies) > maxShow {
		frequencies = frequencies[:maxShow]
	}

	sb.WriteString("| Rank | Entity | Type | Count | Importance |\n")
	sb.WriteString("|------|--------|------|-------|------------|\n")

	for i, ef := range frequencies {
		sb.WriteString(fmt.Sprintf("| %d | `%s` | %s | %d | %d/10 |\n",
			i+1, ef.Text, ef.Type.String(), ef.Count, ef.Importance))
	}

	// Entity type distribution
	sb.WriteString("\n### Entity Type Distribution\n\n")
	typeMap := make(map[semantic.EntityType]int)
	totalCount := 0
	for _, ef := range report.EntityFrequencies {
		typeMap[ef.Type] += ef.Count
		totalCount += ef.Count
	}

	// Sort types by count
	type typePair struct {
		typ   semantic.EntityType
		count int
	}
	var pairs []typePair
	for typ, count := range typeMap {
		pairs = append(pairs, typePair{typ, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	sb.WriteString("| Type | Count | Percentage |\n")
	sb.WriteString("|------|-------|------------|\n")
	for _, p := range pairs {
		pct := 0.0
		if totalCount > 0 {
			pct = float64(p.count) / float64(totalCount) * 100
		}
		sb.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n", p.typ.String(), p.count, pct))
	}

	sb.WriteString("\n")
	return sb.String()
}

// formatRelationshipPatterns formats the relationship patterns section.
func formatRelationshipPatterns(report *AnalysisReport, maxExamples int) string {
	var sb strings.Builder

	sb.WriteString("## Relationship Patterns\n\n")
	sb.WriteString("### Discovered Patterns\n\n")

	sb.WriteString("| Pattern | Occurrences | Strength | Avg Confidence |\n")
	sb.WriteString("|---------|-------------|----------|----------------|\n")

	for _, pattern := range report.RelationshipPatterns {
		patternStr := fmt.Sprintf("%s → %s → %s",
			pattern.SubjectType.String(), pattern.Predicate, pattern.ObjectType.String())
		sb.WriteString(fmt.Sprintf("| %s | %d | %.2f | %.2f |\n",
			patternStr, pattern.Occurrences, pattern.Strength, pattern.AvgConfidence))
	}

	// Pattern examples
	if maxExamples > 0 && len(report.RelationshipPatterns) > 0 {
		sb.WriteString("\n### Pattern Examples\n\n")

		for _, pattern := range report.RelationshipPatterns {
			if len(pattern.Examples) == 0 {
				continue
			}

			patternStr := fmt.Sprintf("%s → %s → %s",
				pattern.SubjectType.String(), pattern.Predicate, pattern.ObjectType.String())

			sb.WriteString(fmt.Sprintf("#### %s\n\n", patternStr))

			examples := pattern.Examples
			if len(examples) > maxExamples {
				examples = examples[:maxExamples]
			}

			for _, ex := range examples {
				sb.WriteString(fmt.Sprintf("- `%s` %s `%s` (confidence: %.2f)\n",
					ex.Subject, ex.Predicate, ex.Object, ex.Confidence))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

// formatTypedHoles formats the typed holes section.
func formatTypedHoles(report *AnalysisReport) string {
	var sb strings.Builder

	sb.WriteString("## Typed Holes\n\n")
	sb.WriteString("### Implementation Priority Queue\n\n")

	// Sort by priority (descending)
	holes := make([]semantic.EnhancedTypedHole, len(report.TypedHoles))
	copy(holes, report.TypedHoles)
	sort.Slice(holes, func(i, j int) bool {
		if holes[i].Priority != holes[j].Priority {
			return holes[i].Priority > holes[j].Priority
		}
		return holes[i].Complexity < holes[j].Complexity
	})

	sb.WriteString("| Priority | Type | Complexity | Constraint | Dependencies |\n")
	sb.WriteString("|----------|------|------------|------------|-------------|\n")

	for _, hole := range holes {
		deps := strings.Join(hole.Dependencies, ", ")
		if deps == "" {
			deps = "none"
		}

		sb.WriteString(fmt.Sprintf("| %d/10 | `%s` | %d/10 | %s | %s |\n",
			hole.Priority, hole.Type, hole.Complexity, hole.Constraint, deps))
	}

	// Detailed hole information
	sb.WriteString("\n### Hole Details\n\n")

	for _, hole := range holes {
		sb.WriteString(fmt.Sprintf("#### `%s`\n\n", hole.Type))
		sb.WriteString(fmt.Sprintf("**Priority:** %d/10 | **Complexity:** %d/10\n\n",
			hole.Priority, hole.Complexity))

		if hole.Constraint != "" {
			sb.WriteString(fmt.Sprintf("**Constraint:** `%s`\n\n", hole.Constraint))
		}

		if hole.SuggestedImpl != "" {
			sb.WriteString(fmt.Sprintf("**Suggested Implementation:**\n```\n%s\n```\n\n", hole.SuggestedImpl))
		}

		if len(hole.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("**Dependencies:** %s\n\n", strings.Join(hole.Dependencies, ", ")))
		}

		if len(hole.RelatedHoles) > 0 {
			sb.WriteString(fmt.Sprintf("**Related Holes:** %s\n\n", strings.Join(hole.RelatedHoles, ", ")))
		}

		if len(hole.Constraints) > 0 {
			sb.WriteString("**Parsed Constraints:**\n\n")
			for _, c := range hole.Constraints {
				sb.WriteString(fmt.Sprintf("- **%s**: %s\n", c.Type, c.Description))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

// formatTripleGraph formats the triple graph section with Mermaid diagram.
func formatTripleGraph(report *AnalysisReport) string {
	var sb strings.Builder

	sb.WriteString("## Triple Graph Visualization\n\n")

	// Graph statistics
	sb.WriteString(fmt.Sprintf("**Nodes:** %d | **Edges:** %d\n\n",
		report.TripleGraph.NodeCount(), report.TripleGraph.EdgeCount()))

	// Mermaid diagram (limit to top nodes for readability)
	sb.WriteString("```mermaid\ngraph TD\n")

	// Get top nodes by importance
	type nodePair struct {
		id         string
		importance int
	}
	var nodes []nodePair
	for id, node := range report.TripleGraph.Nodes {
		nodes = append(nodes, nodePair{id, node.Importance})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].importance > nodes[j].importance
	})

	// Limit to top 15 nodes for clarity
	maxNodes := 15
	if len(nodes) > maxNodes {
		nodes = nodes[:maxNodes]
	}

	// Create node ID mapping
	nodeSet := make(map[string]bool)
	for _, n := range nodes {
		nodeSet[n.id] = true
	}

	// Output nodes with labels
	for _, n := range nodes {
		node := report.TripleGraph.GetNode(n.id)
		label := escapeForMermaid(n.id)
		sb.WriteString(fmt.Sprintf("    %s[\"%s<br/>(%s)\"]\n",
			nodeIDForMermaid(n.id), label, node.Entity.Type.String()))
	}

	// Output edges (only between displayed nodes)
	for _, edge := range report.TripleGraph.Edges {
		if nodeSet[edge.SourceID] && nodeSet[edge.TargetID] {
			label := escapeForMermaid(edge.Relation.Predicate)
			sb.WriteString(fmt.Sprintf("    %s -->|%s| %s\n",
				nodeIDForMermaid(edge.SourceID), label, nodeIDForMermaid(edge.TargetID)))
		}
	}

	sb.WriteString("```\n\n")

	// Node legend
	sb.WriteString("### Graph Legend\n\n")
	sb.WriteString("Node format: `Entity Name (Type)`\n\n")
	sb.WriteString(fmt.Sprintf("Showing top %d nodes by importance. Full graph contains %d nodes.\n\n",
		len(nodes), report.TripleGraph.NodeCount()))

	return sb.String()
}

// nodeIDForMermaid converts an entity text to a valid Mermaid node ID.
func nodeIDForMermaid(text string) string {
	// Replace special characters with underscores
	id := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, text)

	// Ensure it starts with a letter
	if len(id) > 0 && id[0] >= '0' && id[0] <= '9' {
		id = "N" + id
	}

	return id
}

// escapeForMermaid escapes special characters for Mermaid labels.
func escapeForMermaid(text string) string {
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}
