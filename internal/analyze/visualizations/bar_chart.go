package visualizations

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// BarChartConfig configures bar chart rendering.
type BarChartConfig struct {
	Width       int             // Total width of chart
	Height      int             // Total height of chart
	ShowLabels  bool            // Show entity type labels
	ShowCounts  bool            // Show count values
	ShowPercent bool            // Show percentage values
	MaxBarWidth int             // Maximum width for bars
	ColorMap    map[semantic.EntityType]lipgloss.Color // Color per entity type
}

// DefaultBarChartConfig returns default configuration.
func DefaultBarChartConfig() BarChartConfig {
	return BarChartConfig{
		Width:       60,
		Height:      15,
		ShowLabels:  true,
		ShowCounts:  true,
		ShowPercent: true,
		MaxBarWidth: 40,
		ColorMap: map[semantic.EntityType]lipgloss.Color{
			semantic.EntityPerson:       lipgloss.Color("39"),  // Blue
			semantic.EntityOrganization: lipgloss.Color("34"),  // Green
			semantic.EntityPlace:        lipgloss.Color("226"), // Yellow
			semantic.EntityTechnology:   lipgloss.Color("196"), // Red
			semantic.EntityConcept:      lipgloss.Color("141"), // Purple
			semantic.EntityThing:        lipgloss.Color("244"), // Gray
		},
	}
}

// BarChartData represents data for a bar chart.
type BarChartData struct {
	Label string // Bar label
	Value int    // Bar value
	Type  semantic.EntityType // Entity type for coloring
}

// BarChart renders a horizontal bar chart.
type BarChart struct {
	config BarChartConfig
	data   []BarChartData
	title  string
}

// NewBarChart creates a new bar chart.
func NewBarChart(title string, data []BarChartData, config BarChartConfig) *BarChart {
	return &BarChart{
		config: config,
		data:   data,
		title:  title,
	}
}

// Render renders the bar chart as a string.
func (bc *BarChart) Render() string {
	if len(bc.data) == 0 {
		return "No data to display"
	}

	// Find maximum value for scaling
	maxValue := 0
	totalValue := 0
	for _, d := range bc.data {
		if d.Value > maxValue {
			maxValue = d.Value
		}
		totalValue += d.Value
	}

	if maxValue == 0 {
		return "No data with non-zero values"
	}

	var builder strings.Builder

	// Render title
	if bc.title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))
		builder.WriteString(titleStyle.Render(bc.title))
		builder.WriteString("\n\n")
	}

	// Find max label width for alignment
	maxLabelWidth := 0
	for _, d := range bc.data {
		if len(d.Label) > maxLabelWidth {
			maxLabelWidth = len(d.Label)
		}
	}

	// Render each bar
	for _, d := range bc.data {
		// Calculate bar width
		barWidth := 0
		if maxValue > 0 {
			barWidth = (d.Value * bc.config.MaxBarWidth) / maxValue
		}

		// Get color for this entity type
		color := bc.config.ColorMap[d.Type]
		if color == "" {
			color = lipgloss.Color("15") // Default white
		}

		barStyle := lipgloss.NewStyle().Foreground(color)
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
		countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

		// Render label (aligned)
		if bc.config.ShowLabels {
			label := d.Label
			padding := maxLabelWidth - len(label)
			builder.WriteString(labelStyle.Render(label))
			builder.WriteString(strings.Repeat(" ", padding))
			builder.WriteString("  ")
		}

		// Render bar
		bar := strings.Repeat("█", barWidth)
		builder.WriteString(barStyle.Render(bar))
		builder.WriteString(" ")

		// Render count
		if bc.config.ShowCounts {
			builder.WriteString(countStyle.Render(fmt.Sprintf("%d", d.Value)))
		}

		// Render percentage
		if bc.config.ShowPercent && totalValue > 0 {
			percent := float64(d.Value) * 100.0 / float64(totalValue)
			builder.WriteString(countStyle.Render(fmt.Sprintf(" (%.1f%%)", percent)))
		}

		builder.WriteString("\n")
	}

	// Render footer with total
	if totalValue > 0 {
		builder.WriteString("\n")
		footerStyle := lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("244"))
		builder.WriteString(footerStyle.Render(fmt.Sprintf("Total: %d entities", totalValue)))
	}

	return builder.String()
}

// RenderBox renders the bar chart with a border box.
func (bc *BarChart) RenderBox() string {
	content := bc.Render()

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2)

	return boxStyle.Render(content)
}

// SetTitle sets the chart title.
func (bc *BarChart) SetTitle(title string) {
	bc.title = title
}

// SetData sets the chart data.
func (bc *BarChart) SetData(data []BarChartData) {
	bc.data = data
}

// AddBar adds a single bar to the chart.
func (bc *BarChart) AddBar(label string, value int, entityType semantic.EntityType) {
	bc.data = append(bc.data, BarChartData{
		Label: label,
		Value: value,
		Type:  entityType,
	})
}

// Clear clears all chart data.
func (bc *BarChart) Clear() {
	bc.data = nil
}

// GetMaxValue returns the maximum value in the data.
func (bc *BarChart) GetMaxValue() int {
	maxValue := 0
	for _, d := range bc.data {
		if d.Value > maxValue {
			maxValue = d.Value
		}
	}
	return maxValue
}

// GetTotalValue returns the sum of all values.
func (bc *BarChart) GetTotalValue() int {
	total := 0
	for _, d := range bc.data {
		total += d.Value
	}
	return total
}

// DataCount returns the number of data points.
func (bc *BarChart) DataCount() int {
	return len(bc.data)
}
