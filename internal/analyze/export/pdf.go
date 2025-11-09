package export

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// ExportPDF generates a PDF report.
func ExportPDF(report *AnalysisReport) ([]byte, error) {
	if report == nil {
		return nil, fmt.Errorf("report cannot be nil")
	}

	opts := DefaultExportOptions(FormatPDF)
	return ExportPDFWithOptions(report, opts)
}

// ExportPDFWithOptions generates a PDF with custom options.
func ExportPDFWithOptions(report *AnalysisReport, opts ExportOptions) ([]byte, error) {
	if report == nil {
		return nil, fmt.Errorf("report cannot be nil")
	}

	// Initialize PDF (A4 portrait)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)

	// Add first page
	pdf.AddPage()

	// Title page
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 20, report.Title, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(10)
	pdf.CellFormat(0, 10, "Generated: "+report.Timestamp.Format(time.RFC1123), "", 1, "C", false, 0, "")
	if report.Source != "" {
		pdf.CellFormat(0, 10, "Source: "+report.Source, "", 1, "C", false, 0, "")
	}

	pdf.Ln(20)

	// Table of contents
	addPDFTableOfContents(pdf, opts)

	// Statistics section
	if opts.IncludeStatistics {
		pdf.AddPage()
		addPDFStatistics(pdf, report)
	}

	// Entity frequency section
	if opts.IncludeFrequencies && len(report.EntityFrequencies) > 0 {
		pdf.AddPage()
		addPDFEntityFrequencies(pdf, report, opts.MaxFrequenciesToShow)
	}

	// Relationship patterns section
	if opts.IncludePatterns && len(report.RelationshipPatterns) > 0 {
		pdf.AddPage()
		addPDFRelationshipPatterns(pdf, report, opts.MaxExamplesPerPattern)
	}

	// Typed holes section
	if opts.IncludeTypedHoles && len(report.TypedHoles) > 0 {
		pdf.AddPage()
		addPDFTypedHoles(pdf, report)
	}

	// Triple graph section
	if opts.IncludeTripleGraph && report.TripleGraph != nil {
		pdf.AddPage()
		addPDFTripleGraph(pdf, report)
	}

	// Add page numbers
	addPDFFooter(pdf)

	// Get PDF bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// addPDFTableOfContents adds a table of contents to the PDF.
func addPDFTableOfContents(pdf *gofpdf.Fpdf, opts ExportOptions) {
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "Table of Contents", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 12)
	pageNum := 2

	if opts.IncludeStatistics {
		pdf.CellFormat(0, 8, fmt.Sprintf("1. Overview Statistics .............. %d", pageNum), "", 1, "L", false, 0, "")
		pageNum++
	}
	if opts.IncludeFrequencies {
		pdf.CellFormat(0, 8, fmt.Sprintf("2. Entity Frequency Analysis ........ %d", pageNum), "", 1, "L", false, 0, "")
		pageNum++
	}
	if opts.IncludePatterns {
		pdf.CellFormat(0, 8, fmt.Sprintf("3. Relationship Patterns ............ %d", pageNum), "", 1, "L", false, 0, "")
		pageNum++
	}
	if opts.IncludeTypedHoles {
		pdf.CellFormat(0, 8, fmt.Sprintf("4. Typed Holes ...................... %d", pageNum), "", 1, "L", false, 0, "")
		pageNum++
	}
	if opts.IncludeTripleGraph {
		pdf.CellFormat(0, 8, fmt.Sprintf("5. Triple Graph ..................... %d", pageNum), "", 1, "L", false, 0, "")
	}
}

// addPDFStatistics adds the statistics section to the PDF.
func addPDFStatistics(pdf *gofpdf.Fpdf, report *AnalysisReport) {
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "Overview Statistics", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Statistics table
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(100, 126, 234)
	pdf.SetTextColor(255, 255, 255)

	// Table header
	pdf.CellFormat(100, 8, "Metric", "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 8, "Value", "1", 1, "L", true, 0, "")

	// Table rows
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(240, 240, 240)

	addPDFTableRow(pdf, "Total Entities", fmt.Sprintf("%d", report.Stats.TotalEntities), true)
	addPDFTableRow(pdf, "Total Relationships", fmt.Sprintf("%d", report.Stats.TotalRelationships), false)
	addPDFTableRow(pdf, "Total Typed Holes", fmt.Sprintf("%d", report.Stats.TotalTypedHoles), true)
	addPDFTableRow(pdf, "Unique Entity Types", fmt.Sprintf("%d", report.Stats.UniqueEntityTypes), false)
	addPDFTableRow(pdf, "Unique Patterns", fmt.Sprintf("%d", report.Stats.UniquePatterns), true)

	if report.Stats.MostCommonEntity != "" {
		addPDFTableRow(pdf, "Most Common Entity",
			fmt.Sprintf("%s (%d)", report.Stats.MostCommonEntity, report.Stats.MostCommonEntityCount), false)
	}

	if report.Stats.StrongestPattern != "" {
		addPDFTableRow(pdf, "Strongest Pattern",
			fmt.Sprintf("%s (%.2f)", report.Stats.StrongestPattern, report.Stats.StrongestPatternScore), true)
	}

	if report.Stats.HighestPriorityHole != "" {
		addPDFTableRow(pdf, "Highest Priority Hole",
			fmt.Sprintf("%s (priority %d)", report.Stats.HighestPriorityHole, report.Stats.HighestPriority), false)
	}

	if report.Stats.TotalTypedHoles > 0 {
		addPDFTableRow(pdf, "Average Hole Complexity",
			fmt.Sprintf("%.1f", report.Stats.AvgComplexity), true)
	}
}

// addPDFEntityFrequencies adds the entity frequency section.
func addPDFEntityFrequencies(pdf *gofpdf.Fpdf, report *AnalysisReport, maxShow int) {
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "Entity Frequency Analysis", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	frequencies := report.EntityFrequencies
	if len(frequencies) > maxShow {
		frequencies = frequencies[:maxShow]
	}

	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(100, 126, 234)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(15, 8, "Rank", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 8, "Entity", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 8, "Type", "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 8, "Count", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Importance", "1", 1, "C", true, 0, "")

	// Table rows
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(0, 0, 0)

	for i, ef := range frequencies {
		fill := i%2 == 0
		if fill {
			pdf.SetFillColor(240, 240, 240)
		}

		pdf.CellFormat(15, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(70, 7, truncateString(ef.Text, 35), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(40, 7, ef.Type.String(), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(25, 7, fmt.Sprintf("%d", ef.Count), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d/10", ef.Importance), "1", 1, "C", fill, 0, "")
	}

	// ASCII bar chart
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "Type Distribution (ASCII Chart)", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	addPDFASCIIBarChart(pdf, report)
}

// addPDFASCIIBarChart adds an ASCII bar chart for entity type distribution.
func addPDFASCIIBarChart(pdf *gofpdf.Fpdf, report *AnalysisReport) {
	typeMap := make(map[semantic.EntityType]int)
	for _, ef := range report.EntityFrequencies {
		typeMap[ef.Type] += ef.Count
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

	// Find max for scaling
	maxCount := 0
	for _, p := range pairs {
		if p.count > maxCount {
			maxCount = p.count
		}
	}

	pdf.SetFont("Courier", "", 8)
	for _, p := range pairs {
		// Scale to max 40 characters
		barLen := 0
		if maxCount > 0 {
			barLen = (p.count * 40) / maxCount
		}
		if barLen < 1 && p.count > 0 {
			barLen = 1
		}

		bar := strings.Repeat("█", barLen)
		line := fmt.Sprintf("%-15s | %s %d", p.typ.String(), bar, p.count)
		pdf.CellFormat(0, 5, line, "", 1, "L", false, 0, "")
	}
}

// addPDFRelationshipPatterns adds the relationship patterns section.
func addPDFRelationshipPatterns(pdf *gofpdf.Fpdf, report *AnalysisReport, maxExamples int) {
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "Relationship Patterns", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Limit patterns to fit on pages
	patterns := report.RelationshipPatterns
	if len(patterns) > 15 {
		patterns = patterns[:15]
	}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(100, 126, 234)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(80, 8, "Pattern", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Occurrences", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Strength", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Confidence", "1", 1, "C", true, 0, "")

	// Table rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)

	for i, pattern := range patterns {
		fill := i%2 == 0
		if fill {
			pdf.SetFillColor(240, 240, 240)
		}

		patternStr := fmt.Sprintf("%s->%s->%s",
			pattern.SubjectType.String(), pattern.Predicate, pattern.ObjectType.String())

		pdf.CellFormat(80, 7, truncateString(patternStr, 40), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d", pattern.Occurrences), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", pattern.Strength), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", pattern.AvgConfidence), "1", 1, "C", fill, 0, "")
	}
}

// addPDFTypedHoles adds the typed holes section.
func addPDFTypedHoles(pdf *gofpdf.Fpdf, report *AnalysisReport) {
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "Typed Holes", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Sort by priority
	holes := make([]semantic.EnhancedTypedHole, len(report.TypedHoles))
	copy(holes, report.TypedHoles)
	sort.Slice(holes, func(i, j int) bool {
		if holes[i].Priority != holes[j].Priority {
			return holes[i].Priority > holes[j].Priority
		}
		return holes[i].Complexity < holes[j].Complexity
	})

	// Limit to fit on pages
	if len(holes) > 20 {
		holes = holes[:20]
	}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(100, 126, 234)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(30, 8, "Priority", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "Type", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Complexity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "Constraint", "1", 1, "L", true, 0, "")

	// Table rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)

	for i, hole := range holes {
		fill := i%2 == 0
		if fill {
			pdf.SetFillColor(240, 240, 240)
		}

		pdf.CellFormat(30, 7, fmt.Sprintf("%d/10", hole.Priority), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(60, 7, truncateString(hole.Type, 30), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d/10", hole.Complexity), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(60, 7, truncateString(hole.Constraint, 30), "1", 1, "L", fill, 0, "")
	}

	// Detailed information for top holes
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "Top Priority Holes (Detailed)", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	topHoles := holes
	if len(topHoles) > 5 {
		topHoles = topHoles[:5]
	}

	for _, hole := range topHoles {
		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(0, 7, hole.Type, "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 9)
		pdf.CellFormat(0, 6, fmt.Sprintf("Priority: %d/10  |  Complexity: %d/10",
			hole.Priority, hole.Complexity), "", 1, "L", false, 0, "")

		if hole.Constraint != "" {
			pdf.CellFormat(0, 6, "Constraint: "+hole.Constraint, "", 1, "L", false, 0, "")
		}

		if len(hole.Dependencies) > 0 {
			deps := strings.Join(hole.Dependencies, ", ")
			pdf.CellFormat(0, 6, "Dependencies: "+truncateString(deps, 80), "", 1, "L", false, 0, "")
		}

		pdf.Ln(3)
	}
}

// addPDFTripleGraph adds the triple graph section.
func addPDFTripleGraph(pdf *gofpdf.Fpdf, report *AnalysisReport) {
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "Triple Graph", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 8, fmt.Sprintf("Nodes: %d  |  Edges: %d",
		report.TripleGraph.NodeCount(), report.TripleGraph.EdgeCount()), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Top nodes by importance
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "Top Nodes by Importance", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Get top nodes
	type nodePair struct {
		id   string
		node *analyze.TripleNode
	}
	var nodes []nodePair
	for id, node := range report.TripleGraph.Nodes {
		nodes = append(nodes, nodePair{id, node})
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].node.Importance > nodes[j].node.Importance
	})

	maxNodes := 20
	if len(nodes) > maxNodes {
		nodes = nodes[:maxNodes]
	}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(100, 126, 234)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(70, 8, "Entity", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 8, "Type", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Frequency", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Importance", "1", 1, "C", true, 0, "")

	// Table rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(0, 0, 0)

	for i, n := range nodes {
		fill := i%2 == 0
		if fill {
			pdf.SetFillColor(240, 240, 240)
		}

		pdf.CellFormat(70, 7, truncateString(n.id, 35), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(40, 7, n.node.Entity.Type.String(), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d", n.node.Frequency), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d/10", n.node.Importance), "1", 1, "C", fill, 0, "")
	}
}

// addPDFFooter adds page numbers to all pages.
func addPDFFooter(pdf *gofpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
}

// addPDFTableRow adds a table row with alternating fill.
func addPDFTableRow(pdf *gofpdf.Fpdf, label, value string, fill bool) {
	if fill {
		pdf.SetFillColor(240, 240, 240)
	}
	pdf.CellFormat(100, 7, label, "1", 0, "L", fill, 0, "")
	pdf.CellFormat(80, 7, value, "1", 1, "L", fill, 0, "")
}

// truncateString truncates a string to maxLen characters.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
