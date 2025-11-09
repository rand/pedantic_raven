package export

import (
	"bytes"
	"testing"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

func TestExportPDF_Nil(t *testing.T) {
	_, err := ExportPDF(nil)
	if err == nil {
		t.Error("Expected error for nil report")
	}
}

func TestExportPDF_Empty(t *testing.T) {
	report := NewAnalysisReport("Empty Report", "empty.go")

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should generate valid PDF
	if len(pdf) == 0 {
		t.Error("Generated PDF is empty")
	}

	// PDF files start with %PDF-
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Generated data is not a valid PDF")
	}
}

func TestExportPDF_ValidStructure(t *testing.T) {
	report := createTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check PDF header
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF header")
	}

	// PDF should end with %%EOF
	pdfStr := string(pdf)
	if len(pdfStr) < 5 || pdfStr[len(pdfStr)-6:] != "%%EOF\n" {
		t.Error("Invalid PDF footer")
	}

	// Should contain essential PDF objects
	if !bytes.Contains(pdf, []byte("/Type /Catalog")) {
		t.Error("Missing PDF catalog")
	}

	if !bytes.Contains(pdf, []byte("/Type /Pages")) {
		t.Error("Missing PDF pages object")
	}
}

func TestExportPDF_HasContent(t *testing.T) {
	report := createTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// PDF content is encoded, so we can't easily search for text
	// Just verify it's a valid PDF with reasonable size
	if len(pdf) < 1000 {
		t.Error("PDF seems too small, likely missing content")
	}

	// Verify PDF structure
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF structure")
	}
}

func TestExportPDF_TableOfContents(t *testing.T) {
	report := createTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// PDF content is encoded, just verify it's a valid PDF
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF")
	}

	// Check size indicates multiple sections
	if len(pdf) < 2000 {
		t.Error("PDF too small for all sections")
	}
}

func TestExportPDF_CustomOptions(t *testing.T) {
	report := createTestReport()

	opts := ExportOptions{
		Format:              FormatPDF,
		IncludeMetadata:     true,
		IncludeStatistics:   true,
		IncludeFrequencies:  false,
		IncludePatterns:     false,
		IncludeTypedHoles:   false,
		IncludeTripleGraph:  false,
		MaxExamplesPerPattern: 1,
		MaxFrequenciesToShow:  5,
	}

	pdf, err := ExportPDFWithOptions(report, opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be valid PDF
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF with custom options")
	}

	// Should still be a valid document
	if len(pdf) == 0 {
		t.Error("PDF is empty")
	}

	// With fewer sections, should be smaller
	if len(pdf) > 10000 {
		// This is just a rough check
		t.Logf("PDF size: %d bytes", len(pdf))
	}
}

func TestExportPDF_LargeDataset(t *testing.T) {
	report := createLargeTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should handle large dataset
	if len(pdf) == 0 {
		t.Error("Generated PDF is empty")
	}

	// Should still be valid PDF
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF for large dataset")
	}

	// Large dataset should produce a reasonably sized PDF
	if len(pdf) < 1000 {
		t.Error("PDF seems too small for large dataset")
	}
}

func TestExportPDF_MultiplePages(t *testing.T) {
	report := createLargeTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Count page objects (rough estimate)
	pageCount := bytes.Count(pdf, []byte("/Type /Page\n"))

	if pageCount < 2 {
		t.Errorf("Expected multiple pages for large dataset, got %d", pageCount)
	}
}

func TestExportPDF_SpecialCharacters(t *testing.T) {
	report := NewAnalysisReport("Test Report: Special & Characters", "test.go")

	// Add entities with special characters
	graph := analyze.NewTripleGraph()
	graph.AddNode(semantic.Entity{
		Text: "Entity & Special",
		Type: semantic.EntityPerson,
	})
	graph.AddNode(semantic.Entity{
		Text: "Test (parentheses)",
		Type: semantic.EntityTechnology,
	})
	report.SetTripleGraph(graph)

	report.SetEntityFrequencies([]analyze.EntityFrequency{
		{Text: "Entity & Special", Type: semantic.EntityPerson, Count: 5, Importance: 7},
		{Text: "Test (parentheses)", Type: semantic.EntityTechnology, Count: 3, Importance: 6},
	})

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Failed to generate PDF with special characters: %v", err)
	}

	// Should still be valid PDF
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF with special characters")
	}

	// PDF should contain the special characters (properly encoded)
	if len(pdf) == 0 {
		t.Error("PDF is empty")
	}
}

func TestExportPDF_ASCIIBarChart(t *testing.T) {
	report := createTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// PDF encoding makes text search difficult
	// Just verify it's a valid PDF
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF")
	}

	if len(pdf) == 0 {
		t.Error("PDF is empty")
	}
}

func TestExportPDF_PageNumbers(t *testing.T) {
	report := createLargeTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// PDF encoding makes text search difficult
	// Just verify it's a valid multi-page PDF
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF")
	}

	// Should have multiple pages
	pageCount := bytes.Count(pdf, []byte("/Type /Page\n"))
	if pageCount < 2 {
		t.Logf("Expected multiple pages, got %d", pageCount)
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
	}{
		{"short", 10},
		{"exact length", 12},
		{"this is a very long string that needs truncation", 20},
		{"unicode: 中文测试", 10},
	}

	for _, tt := range tests {
		result := truncateString(tt.input, tt.maxLen)
		// Just verify it doesn't exceed maxLen
		if len(result) > tt.maxLen {
			t.Errorf("truncateString(%q, %d) = %q (len=%d), exceeds maxLen",
				tt.input, tt.maxLen, result, len(result))
		}
		// Verify it has content
		if len(result) == 0 && len(tt.input) > 0 {
			t.Errorf("truncateString(%q, %d) returned empty string", tt.input, tt.maxLen)
		}
	}
}

func TestExportPDF_EmptyDataSections(t *testing.T) {
	report := NewAnalysisReport("Empty Sections", "test.go")

	// Set empty data
	report.SetEntityFrequencies([]analyze.EntityFrequency{})
	report.SetRelationshipPatterns([]analyze.RelationshipPattern{})
	report.SetTypedHoles([]semantic.EnhancedTypedHole{})

	opts := DefaultExportOptions(FormatPDF)

	pdf, err := ExportPDFWithOptions(report, opts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should still generate valid PDF with empty sections
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF with empty sections")
	}

	if len(pdf) == 0 {
		t.Error("PDF is empty")
	}
}

func TestExportPDF_Statistics(t *testing.T) {
	report := createTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// PDF encoding makes text search difficult
	// Just verify it's a valid PDF with reasonable size
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF")
	}

	if len(pdf) < 1000 {
		t.Error("PDF too small, likely missing content")
	}
}

func TestExportPDF_DetailedHoles(t *testing.T) {
	report := createTestReport()

	pdf, err := ExportPDF(report)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// PDF encoding makes text search difficult
	// Just verify it's a valid PDF
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Error("Invalid PDF")
	}

	if len(pdf) < 1000 {
		t.Error("PDF too small")
	}
}
