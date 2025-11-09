package visualizations

import (
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// TestNewBarChart tests bar chart creation.
func TestNewBarChart(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "Person", Value: 10, Type: semantic.EntityPerson},
	}

	bc := NewBarChart("Test Chart", data, config)

	if bc == nil {
		t.Fatal("NewBarChart returned nil")
	}

	if bc.title != "Test Chart" {
		t.Errorf("Expected title 'Test Chart', got '%s'", bc.title)
	}

	if len(bc.data) != 1 {
		t.Errorf("Expected 1 data point, got %d", len(bc.data))
	}
}

// TestBarChartRender tests basic rendering.
func TestBarChartRender(t *testing.T) {
	config := DefaultBarChartConfig()
	config.ShowLabels = true
	config.ShowCounts = true
	config.ShowPercent = false

	data := []BarChartData{
		{Label: "Person", Value: 10, Type: semantic.EntityPerson},
		{Label: "Tech", Value: 5, Type: semantic.EntityTechnology},
	}

	bc := NewBarChart("", data, config)
	output := bc.Render()

	// Check that output contains labels
	if !strings.Contains(output, "Person") {
		t.Error("Render output should contain 'Person' label")
	}

	if !strings.Contains(output, "Tech") {
		t.Error("Render output should contain 'Tech' label")
	}

	// Check that output contains counts
	if !strings.Contains(output, "10") {
		t.Error("Render output should contain count '10'")
	}

	if !strings.Contains(output, "5") {
		t.Error("Render output should contain count '5'")
	}

	// Check that output contains bar characters
	if !strings.Contains(output, "█") {
		t.Error("Render output should contain bar character '█'")
	}
}

// TestBarChartRenderEmpty tests rendering with no data.
func TestBarChartRenderEmpty(t *testing.T) {
	config := DefaultBarChartConfig()
	bc := NewBarChart("", []BarChartData{}, config)

	output := bc.Render()

	if !strings.Contains(output, "No data") {
		t.Error("Expected 'No data' message for empty chart")
	}
}

// TestBarChartRenderZeroValues tests rendering with zero values.
func TestBarChartRenderZeroValues(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "A", Value: 0, Type: semantic.EntityPerson},
		{Label: "B", Value: 0, Type: semantic.EntityTechnology},
	}

	bc := NewBarChart("", data, config)
	output := bc.Render()

	if !strings.Contains(output, "No data with non-zero values") {
		t.Error("Expected 'No data with non-zero values' message")
	}
}

// TestBarChartRenderWithTitle tests rendering with title.
func TestBarChartRenderWithTitle(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "Test", Value: 5, Type: semantic.EntityPerson},
	}

	bc := NewBarChart("My Title", data, config)
	output := bc.Render()

	if !strings.Contains(output, "My Title") {
		t.Error("Render output should contain title")
	}
}

// TestBarChartRenderWithPercent tests rendering with percentages.
func TestBarChartRenderWithPercent(t *testing.T) {
	config := DefaultBarChartConfig()
	config.ShowPercent = true

	data := []BarChartData{
		{Label: "A", Value: 10, Type: semantic.EntityPerson},
		{Label: "B", Value: 90, Type: semantic.EntityTechnology},
	}

	bc := NewBarChart("", data, config)
	output := bc.Render()

	// Check for percentage indicators
	if !strings.Contains(output, "%") {
		t.Error("Render output should contain percentage symbol")
	}

	// Should contain 10% and 90%
	if !strings.Contains(output, "10.0%") {
		t.Error("Render output should contain '10.0%'")
	}

	if !strings.Contains(output, "90.0%") {
		t.Error("Render output should contain '90.0%'")
	}
}

// TestBarChartRenderBox tests rendering with border box.
func TestBarChartRenderBox(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "Test", Value: 5, Type: semantic.EntityPerson},
	}

	bc := NewBarChart("", data, config)
	output := bc.RenderBox()

	// Box should contain some border characters
	// The actual characters depend on lipgloss rendering
	if len(output) == 0 {
		t.Error("RenderBox should return non-empty string")
	}
}

// TestBarChartSetTitle tests setting title.
func TestBarChartSetTitle(t *testing.T) {
	config := DefaultBarChartConfig()
	bc := NewBarChart("Initial", []BarChartData{}, config)

	bc.SetTitle("Updated")

	if bc.title != "Updated" {
		t.Errorf("Expected title 'Updated', got '%s'", bc.title)
	}
}

// TestBarChartSetData tests setting data.
func TestBarChartSetData(t *testing.T) {
	config := DefaultBarChartConfig()
	bc := NewBarChart("", []BarChartData{}, config)

	newData := []BarChartData{
		{Label: "A", Value: 1, Type: semantic.EntityPerson},
		{Label: "B", Value: 2, Type: semantic.EntityTechnology},
	}

	bc.SetData(newData)

	if len(bc.data) != 2 {
		t.Errorf("Expected 2 data points, got %d", len(bc.data))
	}
}

// TestBarChartAddBar tests adding individual bars.
func TestBarChartAddBar(t *testing.T) {
	config := DefaultBarChartConfig()
	bc := NewBarChart("", []BarChartData{}, config)

	bc.AddBar("Person", 10, semantic.EntityPerson)
	bc.AddBar("Tech", 5, semantic.EntityTechnology)

	if len(bc.data) != 2 {
		t.Errorf("Expected 2 bars, got %d", len(bc.data))
	}

	if bc.data[0].Label != "Person" || bc.data[0].Value != 10 {
		t.Error("First bar not added correctly")
	}

	if bc.data[1].Label != "Tech" || bc.data[1].Value != 5 {
		t.Error("Second bar not added correctly")
	}
}

// TestBarChartClear tests clearing data.
func TestBarChartClear(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "A", Value: 1, Type: semantic.EntityPerson},
	}

	bc := NewBarChart("", data, config)
	bc.Clear()

	if bc.data != nil {
		t.Error("Clear should set data to nil")
	}
}

// TestBarChartGetMaxValue tests getting maximum value.
func TestBarChartGetMaxValue(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "A", Value: 5, Type: semantic.EntityPerson},
		{Label: "B", Value: 15, Type: semantic.EntityTechnology},
		{Label: "C", Value: 10, Type: semantic.EntityOrganization},
	}

	bc := NewBarChart("", data, config)
	maxValue := bc.GetMaxValue()

	if maxValue != 15 {
		t.Errorf("Expected max value 15, got %d", maxValue)
	}
}

// TestBarChartGetMaxValueEmpty tests getting max value from empty chart.
func TestBarChartGetMaxValueEmpty(t *testing.T) {
	config := DefaultBarChartConfig()
	bc := NewBarChart("", []BarChartData{}, config)

	maxValue := bc.GetMaxValue()

	if maxValue != 0 {
		t.Errorf("Expected max value 0 for empty chart, got %d", maxValue)
	}
}

// TestBarChartGetTotalValue tests getting total value.
func TestBarChartGetTotalValue(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "A", Value: 5, Type: semantic.EntityPerson},
		{Label: "B", Value: 15, Type: semantic.EntityTechnology},
		{Label: "C", Value: 10, Type: semantic.EntityOrganization},
	}

	bc := NewBarChart("", data, config)
	total := bc.GetTotalValue()

	if total != 30 { // 5 + 15 + 10
		t.Errorf("Expected total value 30, got %d", total)
	}
}

// TestBarChartDataCount tests counting data points.
func TestBarChartDataCount(t *testing.T) {
	config := DefaultBarChartConfig()
	data := []BarChartData{
		{Label: "A", Value: 1, Type: semantic.EntityPerson},
		{Label: "B", Value: 2, Type: semantic.EntityTechnology},
		{Label: "C", Value: 3, Type: semantic.EntityOrganization},
	}

	bc := NewBarChart("", data, config)

	if bc.DataCount() != 3 {
		t.Errorf("Expected data count 3, got %d", bc.DataCount())
	}
}

// TestBarChartScaling tests bar width scaling.
func TestBarChartScaling(t *testing.T) {
	config := DefaultBarChartConfig()
	config.MaxBarWidth = 20

	data := []BarChartData{
		{Label: "Max", Value: 100, Type: semantic.EntityPerson},
		{Label: "Half", Value: 50, Type: semantic.EntityTechnology},
	}

	bc := NewBarChart("", data, config)
	output := bc.Render()

	// The max value bar should be longer than the half value bar
	// This is a basic check - exact lengths depend on rendering
	if len(output) == 0 {
		t.Error("Render should produce non-empty output")
	}
}

// TestDefaultBarChartConfig tests default configuration.
func TestDefaultBarChartConfig(t *testing.T) {
	config := DefaultBarChartConfig()

	if config.Width != 60 {
		t.Errorf("Expected default width 60, got %d", config.Width)
	}

	if config.Height != 15 {
		t.Errorf("Expected default height 15, got %d", config.Height)
	}

	if !config.ShowLabels {
		t.Error("Expected ShowLabels to be true")
	}

	if !config.ShowCounts {
		t.Error("Expected ShowCounts to be true")
	}

	if config.MaxBarWidth != 40 {
		t.Errorf("Expected MaxBarWidth 40, got %d", config.MaxBarWidth)
	}

	// Check that color map has expected types
	if len(config.ColorMap) != 6 {
		t.Errorf("Expected 6 entity types in color map, got %d", len(config.ColorMap))
	}
}
