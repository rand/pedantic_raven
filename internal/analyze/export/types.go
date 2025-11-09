// Package export provides multi-format export capabilities for analysis reports.
package export

import (
	"time"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// AnalysisReport aggregates all analysis data for export.
type AnalysisReport struct {
	// Metadata
	Timestamp   time.Time
	Source      string
	Title       string
	Description string

	// Core Analysis Data
	TripleGraph         *analyze.TripleGraph
	EntityFrequencies   []analyze.EntityFrequency
	RelationshipPatterns []analyze.RelationshipPattern
	TypedHoles          []semantic.EnhancedTypedHole

	// Statistics
	Stats ReportStatistics
}

// ReportStatistics contains summary statistics for the report.
type ReportStatistics struct {
	TotalEntities      int
	TotalRelationships int
	TotalTypedHoles    int
	UniqueEntityTypes  int
	UniquePatterns     int

	// Entity statistics
	MostCommonEntity    string
	MostCommonEntityCount int

	// Pattern statistics
	StrongestPattern    string
	StrongestPatternScore float64

	// Typed hole statistics
	HighestPriorityHole string
	HighestPriority     int
	AvgComplexity       float64
}

// NewAnalysisReport creates a new analysis report with timestamp.
func NewAnalysisReport(title, source string) *AnalysisReport {
	return &AnalysisReport{
		Timestamp:   time.Now(),
		Title:       title,
		Source:      source,
		Stats:       ReportStatistics{},
	}
}

// SetTripleGraph sets the triple graph and updates statistics.
func (r *AnalysisReport) SetTripleGraph(graph *analyze.TripleGraph) {
	r.TripleGraph = graph
	if graph != nil {
		r.Stats.TotalEntities = graph.NodeCount()
		r.Stats.TotalRelationships = graph.EdgeCount()

		// Count unique entity types
		types := make(map[semantic.EntityType]bool)
		for _, node := range graph.Nodes {
			types[node.Entity.Type] = true
		}
		r.Stats.UniqueEntityTypes = len(types)
	}
}

// SetEntityFrequencies sets entity frequencies and updates statistics.
func (r *AnalysisReport) SetEntityFrequencies(frequencies []analyze.EntityFrequency) {
	r.EntityFrequencies = frequencies
	if len(frequencies) > 0 {
		r.Stats.MostCommonEntity = frequencies[0].Text
		r.Stats.MostCommonEntityCount = frequencies[0].Count
	}
}

// SetRelationshipPatterns sets relationship patterns and updates statistics.
func (r *AnalysisReport) SetRelationshipPatterns(patterns []analyze.RelationshipPattern) {
	r.RelationshipPatterns = patterns
	r.Stats.UniquePatterns = len(patterns)

	if len(patterns) > 0 {
		// Find strongest pattern
		strongest := patterns[0]
		for _, p := range patterns {
			if p.Strength > strongest.Strength {
				strongest = p
			}
		}
		r.Stats.StrongestPattern = formatPattern(strongest)
		r.Stats.StrongestPatternScore = strongest.Strength
	}
}

// SetTypedHoles sets typed holes and updates statistics.
func (r *AnalysisReport) SetTypedHoles(holes []semantic.EnhancedTypedHole) {
	r.TypedHoles = holes
	r.Stats.TotalTypedHoles = len(holes)

	if len(holes) > 0 {
		// Find highest priority hole
		highest := holes[0]
		totalComplexity := 0
		for _, h := range holes {
			if h.Priority > highest.Priority {
				highest = h
			}
			totalComplexity += h.Complexity
		}
		r.Stats.HighestPriorityHole = highest.Type
		r.Stats.HighestPriority = highest.Priority
		r.Stats.AvgComplexity = float64(totalComplexity) / float64(len(holes))
	}
}

// formatPattern formats a relationship pattern as a readable string.
func formatPattern(p analyze.RelationshipPattern) string {
	return p.SubjectType.String() + " " + p.Predicate + " " + p.ObjectType.String()
}

// ExportFormat represents supported export formats.
type ExportFormat string

const (
	FormatMarkdown ExportFormat = "markdown"
	FormatHTML     ExportFormat = "html"
	FormatPDF      ExportFormat = "pdf"
)

// ExportOptions configures export behavior.
type ExportOptions struct {
	Format              ExportFormat
	IncludeMetadata     bool
	IncludeStatistics   bool
	IncludeTripleGraph  bool
	IncludeFrequencies  bool
	IncludePatterns     bool
	IncludeTypedHoles   bool
	MaxExamplesPerPattern int
	MaxFrequenciesToShow  int
}

// DefaultExportOptions returns default export options for a format.
func DefaultExportOptions(format ExportFormat) ExportOptions {
	return ExportOptions{
		Format:              format,
		IncludeMetadata:     true,
		IncludeStatistics:   true,
		IncludeTripleGraph:  true,
		IncludeFrequencies:  true,
		IncludePatterns:     true,
		IncludeTypedHoles:   true,
		MaxExamplesPerPattern: 3,
		MaxFrequenciesToShow:  20,
	}
}
