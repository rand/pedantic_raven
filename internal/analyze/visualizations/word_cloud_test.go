package visualizations

import (
	"strings"
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// TestNewWordCloud tests word cloud creation.
func TestNewWordCloud(t *testing.T) {
	config := DefaultWordCloudConfig()
	wc := NewWordCloud("Test Cloud", config)

	if wc == nil {
		t.Fatal("NewWordCloud returned nil")
	}

	if wc.title != "Test Cloud" {
		t.Errorf("Expected title 'Test Cloud', got '%s'", wc.title)
	}

	if wc.WordCount() != 0 {
		t.Errorf("Expected 0 words initially, got %d", wc.WordCount())
	}
}

// TestWordCloudAddWord tests adding words to the cloud.
func TestWordCloudAddWord(t *testing.T) {
	config := DefaultWordCloudConfig()
	wc := NewWordCloud("", config)

	wc.AddWord("API", 10, semantic.EntityTechnology)
	wc.AddWord("Database", 5, semantic.EntityTechnology)

	if wc.WordCount() != 2 {
		t.Errorf("Expected 2 words, got %d", wc.WordCount())
	}

	// Check first word
	if wc.words[0].Text != "API" {
		t.Errorf("Expected first word 'API', got '%s'", wc.words[0].Text)
	}

	if wc.words[0].Frequency != 10 {
		t.Errorf("Expected first word frequency 10, got %d", wc.words[0].Frequency)
	}
}

// TestCalculateSize tests font size calculation.
func TestCalculateSize(t *testing.T) {
	config := DefaultWordCloudConfig()
	wc := NewWordCloud("", config)

	tests := []struct {
		name      string
		frequency int
		wantSize  int
	}{
		{"Zero frequency", 0, 1},
		{"Small (1)", 1, 1},
		{"Small (3)", 3, 1},
		{"Medium (4)", 4, 2},
		{"Medium (10)", 10, 2},
		{"Large (11)", 11, 3},
		{"Large (100)", 100, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := wc.calculateSize(tt.frequency)

			if size != tt.wantSize {
				t.Errorf("calculateSize(%d) = %d, want %d", tt.frequency, size, tt.wantSize)
			}
		})
	}
}

// TestWordCloudLayout tests layout algorithm.
func TestWordCloudLayout(t *testing.T) {
	config := DefaultWordCloudConfig()
	config.Width = 50
	config.Height = 10

	wc := NewWordCloud("", config)

	// Add several words
	wc.AddWord("API", 20, semantic.EntityTechnology)
	wc.AddWord("Database", 15, semantic.EntityTechnology)
	wc.AddWord("Server", 10, semantic.EntityTechnology)
	wc.AddWord("Client", 5, semantic.EntityTechnology)

	// Perform layout
	wc.Layout()

	// Check that words have positions assigned
	placedCount := 0
	for _, word := range wc.words {
		if word.X >= 0 && word.Y >= 0 {
			placedCount++
		}
	}

	if placedCount == 0 {
		t.Error("Layout should place at least one word")
	}

	// Check that positions are within bounds
	for i, word := range wc.words {
		if word.X >= 0 && word.Y >= 0 {
			if word.X < 0 || word.X >= config.Width {
				t.Errorf("Word %d X position %d out of bounds [0, %d)", i, word.X, config.Width)
			}
			if word.Y < 0 || word.Y >= config.Height {
				t.Errorf("Word %d Y position %d out of bounds [0, %d)", i, word.Y, config.Height)
			}
		}
	}
}

// TestWordCloudLayoutSorting tests that layout places larger words first.
func TestWordCloudLayoutSorting(t *testing.T) {
	config := DefaultWordCloudConfig()
	config.Width = 50
	config.Height = 10

	wc := NewWordCloud("", config)

	// Add words in random frequency order
	wc.AddWord("Small", 1, semantic.EntityTechnology)
	wc.AddWord("Large", 100, semantic.EntityTechnology)
	wc.AddWord("Medium", 10, semantic.EntityTechnology)

	wc.Layout()

	// After layout, words should be sorted by size (descending)
	// Large (size 3) should be first
	if wc.words[0].Text != "Large" {
		t.Errorf("Expected 'Large' first after layout, got '%s'", wc.words[0].Text)
	}
}

// TestWordCloudRender tests rendering.
func TestWordCloudRender(t *testing.T) {
	config := DefaultWordCloudConfig()
	config.Width = 50
	config.Height = 10

	wc := NewWordCloud("Test Title", config)

	wc.AddWord("API", 20, semantic.EntityTechnology)
	wc.AddWord("Database", 15, semantic.EntityTechnology)

	output := wc.Render()

	// Check for title
	if !strings.Contains(output, "Test Title") {
		t.Error("Render output should contain title")
	}

	// Check for footer
	if !strings.Contains(output, "Font size = log(frequency)") {
		t.Error("Render output should contain footer text")
	}

	// Output should be non-empty
	if len(output) == 0 {
		t.Error("Render should produce non-empty output")
	}

	// Should contain newlines (multiple lines)
	if !strings.Contains(output, "\n") {
		t.Error("Render output should contain newlines")
	}
}

// TestWordCloudRenderEmpty tests rendering with no words.
func TestWordCloudRenderEmpty(t *testing.T) {
	config := DefaultWordCloudConfig()
	wc := NewWordCloud("", config)

	output := wc.Render()

	// Should still produce output (empty canvas with footer)
	if len(output) == 0 {
		t.Error("Render should produce output even with no words")
	}
}

// TestWordCloudRenderBox tests rendering with border box.
func TestWordCloudRenderBox(t *testing.T) {
	config := DefaultWordCloudConfig()
	wc := NewWordCloud("", config)

	wc.AddWord("Test", 10, semantic.EntityTechnology)

	output := wc.RenderBox()

	// Box should contain content
	if len(output) == 0 {
		t.Error("RenderBox should return non-empty string")
	}
}

// TestWordCloudSetTitle tests setting title.
func TestWordCloudSetTitle(t *testing.T) {
	config := DefaultWordCloudConfig()
	wc := NewWordCloud("Initial", config)

	wc.SetTitle("Updated")

	if wc.title != "Updated" {
		t.Errorf("Expected title 'Updated', got '%s'", wc.title)
	}
}

// TestWordCloudClear tests clearing words.
func TestWordCloudClear(t *testing.T) {
	config := DefaultWordCloudConfig()
	wc := NewWordCloud("", config)

	wc.AddWord("Test", 10, semantic.EntityTechnology)
	wc.Clear()

	if wc.words != nil {
		t.Error("Clear should set words to nil")
	}

	if wc.WordCount() != 0 {
		t.Errorf("Expected 0 words after clear, got %d", wc.WordCount())
	}
}

// TestWordCloudMaxWords tests limiting to max words.
func TestWordCloudMaxWords(t *testing.T) {
	config := DefaultWordCloudConfig()
	config.MaxWords = 3
	config.Width = 100
	config.Height = 20

	wc := NewWordCloud("", config)

	// Add more than MaxWords
	wc.AddWord("Word1", 10, semantic.EntityTechnology)
	wc.AddWord("Word2", 9, semantic.EntityTechnology)
	wc.AddWord("Word3", 8, semantic.EntityTechnology)
	wc.AddWord("Word4", 7, semantic.EntityTechnology)
	wc.AddWord("Word5", 6, semantic.EntityTechnology)

	wc.Layout()

	// After layout, should only have MaxWords
	if len(wc.words) != config.MaxWords {
		t.Errorf("Expected %d words after layout, got %d", config.MaxWords, len(wc.words))
	}
}

// TestWordCloudNoOverlap tests that words don't overlap.
func TestWordCloudNoOverlap(t *testing.T) {
	config := DefaultWordCloudConfig()
	config.Width = 100
	config.Height = 20

	wc := NewWordCloud("", config)

	// Add words
	wc.AddWord("API", 20, semantic.EntityTechnology)
	wc.AddWord("Database", 15, semantic.EntityTechnology)
	wc.AddWord("Server", 10, semantic.EntityTechnology)

	wc.Layout()

	// Check for overlaps
	for i, word1 := range wc.words {
		if word1.X < 0 || word1.Y < 0 {
			continue // Not placed
		}

		for j, word2 := range wc.words {
			if i == j || word2.X < 0 || word2.Y < 0 {
				continue
			}

			// Check if words are on the same line
			if word1.Y == word2.Y {
				// Check horizontal overlap
				word1End := word1.X + len(word1.Text)
				word2End := word2.X + len(word2.Text)

				if !(word1End <= word2.X || word2End <= word1.X) {
					t.Errorf("Words %d and %d overlap: '%s' at (%d,%d) and '%s' at (%d,%d)",
						i, j, word1.Text, word1.X, word1.Y, word2.Text, word2.X, word2.Y)
				}
			}
		}
	}
}

// TestWordCloudCanPlace tests collision detection.
func TestWordCloudCanPlace(t *testing.T) {
	config := DefaultWordCloudConfig()
	config.Width = 20
	config.Height = 5

	wc := NewWordCloud("", config)

	// Create empty canvas
	canvas := make([][]bool, config.Height)
	for i := range canvas {
		canvas[i] = make([]bool, config.Width)
	}

	word := &WordCloudWord{
		Text:      "Test",
		Frequency: 10,
		Size:      1,
		Type:      semantic.EntityTechnology,
	}

	// Should be able to place at 0,0
	if !wc.canPlace(word, 0, 0, canvas) {
		t.Error("Should be able to place word at (0,0) on empty canvas")
	}

	// Mark some cells as occupied
	canvas[0][2] = true

	// Should not be able to place at 0,0 now (overlaps with occupied cell)
	if wc.canPlace(word, 0, 0, canvas) {
		t.Error("Should not be able to place word that overlaps occupied cell")
	}

	// Should be able to place at 5,0 (no overlap)
	if !wc.canPlace(word, 5, 0, canvas) {
		t.Error("Should be able to place word at (5,0)")
	}
}

// TestWordCloudCanPlaceBounds tests boundary checking.
func TestWordCloudCanPlaceBounds(t *testing.T) {
	config := DefaultWordCloudConfig()
	config.Width = 20
	config.Height = 5

	wc := NewWordCloud("", config)

	canvas := make([][]bool, config.Height)
	for i := range canvas {
		canvas[i] = make([]bool, config.Width)
	}

	word := &WordCloudWord{
		Text:      "Test",
		Frequency: 10,
		Size:      1,
		Type:      semantic.EntityTechnology,
	}

	// Test out of bounds positions
	tests := []struct {
		name string
		x    int
		y    int
		want bool
	}{
		{"Negative X", -1, 0, false},
		{"Negative Y", 0, -1, false},
		{"Too far right", config.Width, 0, false},
		{"Too far down", 0, config.Height, false},
		{"Word extends beyond right edge", config.Width - 2, 0, false},
		{"Valid position", 5, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wc.canPlace(word, tt.x, tt.y, canvas)
			if got != tt.want {
				t.Errorf("canPlace(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

// TestDefaultWordCloudConfig tests default configuration.
func TestDefaultWordCloudConfig(t *testing.T) {
	config := DefaultWordCloudConfig()

	if config.Width != 70 {
		t.Errorf("Expected default width 70, got %d", config.Width)
	}

	if config.Height != 20 {
		t.Errorf("Expected default height 20, got %d", config.Height)
	}

	if config.MaxWords != 30 {
		t.Errorf("Expected MaxWords 30, got %d", config.MaxWords)
	}

	if config.MinFontSize != 1 {
		t.Errorf("Expected MinFontSize 1, got %d", config.MinFontSize)
	}

	if config.MaxFontSize != 3 {
		t.Errorf("Expected MaxFontSize 3, got %d", config.MaxFontSize)
	}

	// Check color map
	if len(config.ColorMap) != 6 {
		t.Errorf("Expected 6 entity types in color map, got %d", len(config.ColorMap))
	}
}
