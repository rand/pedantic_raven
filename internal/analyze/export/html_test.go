package export

import (
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

func TestExportHTML_Nil(t *testing.T) {
	_, err := ExportHTML(nil)
	if err == nil {
		t.Error("Expected error for nil report")
	}
}

func TestExportHTML_Empty(t *testing.T) {
	report := NewAnalysisReport("Empty Report", "empty.go")

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be valid HTML
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Missing DOCTYPE declaration")
	}

	if !strings.Contains(html, "<html") {
		t.Error("Missing html tag")
	}

	if !strings.Contains(html, "</html>") {
		t.Error("Missing closing html tag")
	}

	// Should contain title
	if !strings.Contains(html, "<h1>Empty Report</h1>") {
		t.Error("Missing title heading")
	}
}

func TestExportHTML_HasEmbeddedCSS(t *testing.T) {
	report := NewAnalysisReport("Test", "test.go")

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have embedded CSS
	if !strings.Contains(html, "<style>") {
		t.Error("Missing style tag")
	}

	// Check for some key CSS rules
	if !strings.Contains(html, "font-family") {
		t.Error("Missing CSS font-family")
	}

	if !strings.Contains(html, ".entity-type") {
		t.Error("Missing entity-type CSS class")
	}

	if !strings.Contains(html, "@media print") {
		t.Error("Missing print media query")
	}

	if !strings.Contains(html, "@media (max-width: 768px)") {
		t.Error("Missing mobile responsive media query")
	}
}

func TestExportHTML_HasEmbeddedJS(t *testing.T) {
	report := NewAnalysisReport("Test", "test.go")

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have embedded JavaScript
	if !strings.Contains(html, "<script>") {
		t.Error("Missing script tag")
	}

	// Check for interactive features
	if !strings.Contains(html, "collapsible") {
		t.Error("Missing collapsible functionality")
	}

	if !strings.Contains(html, "sortTable") {
		t.Error("Missing table sorting functionality")
	}

	if !strings.Contains(html, "scrollIntoView") {
		t.Error("Missing smooth scrolling functionality")
	}
}

func TestExportHTML_Metadata(t *testing.T) {
	report := NewAnalysisReport("Test Report", "test.go")
	report.Description = "Test description"

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should contain metadata
	if !strings.Contains(html, "Generated:") {
		t.Error("Missing generation timestamp")
	}

	if !strings.Contains(html, "Source:") {
		t.Error("Missing source")
	}

	if !strings.Contains(html, "<code>test.go</code>") {
		t.Error("Missing source filename")
	}
}

func TestExportHTML_TableOfContents(t *testing.T) {
	report := createTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have table of contents
	if !strings.Contains(html, "<nav class=\"toc\">") {
		t.Error("Missing table of contents nav")
	}

	if !strings.Contains(html, "<h2>Table of Contents</h2>") {
		t.Error("Missing TOC heading")
	}

	// Should have anchor links
	if !strings.Contains(html, "href=\"#statistics\"") {
		t.Error("Missing statistics anchor link")
	}

	if !strings.Contains(html, "href=\"#frequencies\"") {
		t.Error("Missing frequencies anchor link")
	}

	if !strings.Contains(html, "href=\"#patterns\"") {
		t.Error("Missing patterns anchor link")
	}

	if !strings.Contains(html, "href=\"#holes\"") {
		t.Error("Missing holes anchor link")
	}

	if !strings.Contains(html, "href=\"#graph\"") {
		t.Error("Missing graph anchor link")
	}
}

func TestExportHTML_Statistics(t *testing.T) {
	report := createTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have statistics section
	if !strings.Contains(html, "id=\"statistics\"") {
		t.Error("Missing statistics section ID")
	}

	if !strings.Contains(html, "<h2>Overview Statistics</h2>") {
		t.Error("Missing statistics heading")
	}

	// Should have statistics table
	if !strings.Contains(html, "Total Entities") {
		t.Error("Missing total entities statistic")
	}

	if !strings.Contains(html, "Total Relationships") {
		t.Error("Missing total relationships statistic")
	}
}

func TestExportHTML_EntityFrequencies(t *testing.T) {
	report := createTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have frequencies section
	if !strings.Contains(html, "id=\"frequencies\"") {
		t.Error("Missing frequencies section ID")
	}

	if !strings.Contains(html, "<h2>Entity Frequency Analysis</h2>") {
		t.Error("Missing frequencies heading")
	}

	// Should have entity data
	if !strings.Contains(html, "Alice") {
		t.Error("Missing entity 'Alice'")
	}

	// Should have entity type badges
	if !strings.Contains(html, "entity-type") {
		t.Error("Missing entity-type CSS class")
	}
}

func TestExportHTML_RelationshipPatterns(t *testing.T) {
	report := createTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have patterns section
	if !strings.Contains(html, "id=\"patterns\"") {
		t.Error("Missing patterns section ID")
	}

	if !strings.Contains(html, "<h2>Relationship Patterns</h2>") {
		t.Error("Missing patterns heading")
	}

	// Should have pattern data
	if !strings.Contains(html, "uses") {
		t.Error("Missing 'uses' predicate")
	}
}

func TestExportHTML_TypedHoles(t *testing.T) {
	report := createTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have holes section
	if !strings.Contains(html, "id=\"holes\"") {
		t.Error("Missing holes section ID")
	}

	if !strings.Contains(html, "<h2>Typed Holes</h2>") {
		t.Error("Missing holes heading")
	}

	// Should have hole data
	if !strings.Contains(html, "UserService") {
		t.Error("Missing 'UserService' hole")
	}

	// Should have priority classes
	if !strings.Contains(html, "priority-") {
		t.Error("Missing priority CSS classes")
	}
}

func TestExportHTML_TripleGraph(t *testing.T) {
	report := createTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should have graph section
	if !strings.Contains(html, "id=\"graph\"") {
		t.Error("Missing graph section ID")
	}

	if !strings.Contains(html, "<h2>Triple Graph</h2>") {
		t.Error("Missing graph heading")
	}

	// Should have node data
	if !strings.Contains(html, "<h3>Top Nodes by Importance</h3>") {
		t.Error("Missing top nodes subsection")
	}
}

func TestExportHTML_CustomOptions(t *testing.T) {
	report := createTestReport()

	opts := ExportOptions{
		Format:              FormatHTML,
		IncludeMetadata:     false,
		IncludeStatistics:   true,
		IncludeFrequencies:  false,
		IncludePatterns:     false,
		IncludeTypedHoles:   false,
		IncludeTripleGraph:  false,
		MaxExamplesPerPattern: 1,
		MaxFrequenciesToShow:  5,
	}

	html, err := ExportHTMLWithOptions(report, opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should NOT contain metadata
	if strings.Contains(html, "<p class=\"metadata\">Generated:") {
		t.Error("Should not include metadata")
	}

	// Should contain statistics (enabled)
	if !strings.Contains(html, "id=\"statistics\"") {
		t.Error("Should include statistics")
	}

	// Should NOT contain frequencies (disabled)
	if strings.Contains(html, "id=\"frequencies\"") {
		t.Error("Should not include frequencies")
	}

	// Should NOT contain patterns (disabled)
	if strings.Contains(html, "id=\"patterns\"") {
		t.Error("Should not include patterns")
	}
}

func TestExportHTML_EntityTypeClasses(t *testing.T) {
	tests := []struct {
		entityType    semantic.EntityType
		expectedClass string
	}{
		{semantic.EntityPerson, "entity-type-person"},
		{semantic.EntityOrganization, "entity-type-organization"},
		{semantic.EntityTechnology, "entity-type-technology"},
		{semantic.EntityConcept, "entity-type-concept"},
		{semantic.EntityPlace, "entity-type-place"},
		{semantic.EntityThing, "entity-type-thing"},
	}

	for _, tt := range tests {
		result := getEntityTypeClass(tt.entityType)
		if result != tt.expectedClass {
			t.Errorf("getEntityTypeClass(%v) = %q, expected %q", tt.entityType, result, tt.expectedClass)
		}
	}
}

func TestExportHTML_PriorityClasses(t *testing.T) {
	tests := []struct {
		priority      int
		expectedClass string
	}{
		{10, "priority-high"},
		{7, "priority-high"},
		{6, "priority-medium"},
		{4, "priority-medium"},
		{3, "priority-low"},
		{0, "priority-low"},
	}

	for _, tt := range tests {
		result := getPriorityClass(tt.priority)
		if result != tt.expectedClass {
			t.Errorf("getPriorityClass(%d) = %q, expected %q", tt.priority, result, tt.expectedClass)
		}
	}
}

func TestExportHTML_HTMLEscaping(t *testing.T) {
	report := NewAnalysisReport("Test <script>alert('xss')</script>", "test.go")

	// Add entity with special characters
	graph := analyze.NewTripleGraph()
	graph.AddNode(semantic.Entity{
		Text: "<script>alert('xss')</script>",
		Type: semantic.EntityPerson,
	})
	report.SetTripleGraph(graph)

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should escape HTML special characters
	if strings.Contains(html, "<script>alert('xss')</script>") {
		t.Error("HTML not properly escaped - XSS vulnerability!")
	}

	if !strings.Contains(html, "&lt;script&gt;") {
		t.Error("Expected escaped HTML entities")
	}
}

func TestExportHTML_ValidStructure(t *testing.T) {
	report := createTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check basic HTML structure
	requiredTags := []string{
		"<!DOCTYPE html>",
		"<html",
		"<head>",
		"</head>",
		"<body>",
		"</body>",
		"</html>",
		"<header>",
		"</header>",
		"<main>",
		"</main>",
		"<footer>",
		"</footer>",
	}

	for _, tag := range requiredTags {
		if !strings.Contains(html, tag) {
			t.Errorf("Missing required HTML tag: %s", tag)
		}
	}

	// Check meta tags
	if !strings.Contains(html, `<meta charset="UTF-8">`) {
		t.Error("Missing charset meta tag")
	}

	if !strings.Contains(html, `<meta name="viewport"`) {
		t.Error("Missing viewport meta tag")
	}
}

func TestExportHTML_LargeDataset(t *testing.T) {
	report := createLargeTestReport()

	html, err := ExportHTML(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should handle large dataset without errors
	if len(html) == 0 {
		t.Error("Generated HTML is empty")
	}

	// Should still be valid HTML
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Invalid HTML structure for large dataset")
	}
}
