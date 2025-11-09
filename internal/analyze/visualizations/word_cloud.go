package visualizations

import (
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// WordCloudConfig configures word cloud rendering.
type WordCloudConfig struct {
	Width       int                                    // Canvas width
	Height      int                                    // Canvas height
	MaxWords    int                                    // Maximum number of words to display
	MinFontSize int                                    // Minimum font size (number of times to repeat)
	MaxFontSize int                                    // Maximum font size
	ColorMap    map[semantic.EntityType]lipgloss.Color // Color per entity type
	Seed        int64                                  // Random seed for layout
}

// DefaultWordCloudConfig returns default configuration.
func DefaultWordCloudConfig() WordCloudConfig {
	return WordCloudConfig{
		Width:       70,
		Height:      20,
		MaxWords:    30,
		MinFontSize: 1,
		MaxFontSize: 3,
		ColorMap: map[semantic.EntityType]lipgloss.Color{
			semantic.EntityPerson:       lipgloss.Color("39"),  // Blue
			semantic.EntityOrganization: lipgloss.Color("34"),  // Green
			semantic.EntityPlace:        lipgloss.Color("226"), // Yellow
			semantic.EntityTechnology:   lipgloss.Color("196"), // Red
			semantic.EntityConcept:      lipgloss.Color("141"), // Purple
			semantic.EntityThing:        lipgloss.Color("244"), // Gray
		},
		Seed: 42,
	}
}

// WordCloudWord represents a word in the cloud.
type WordCloudWord struct {
	Text      string              // Word text
	Frequency int                 // Word frequency
	Size      int                 // Font size (1-3)
	Type      semantic.EntityType // Entity type
	X         int                 // X position on canvas
	Y         int                 // Y position on canvas
}

// WordCloud renders an ASCII word cloud.
type WordCloud struct {
	config WordCloudConfig
	words  []WordCloudWord
	title  string
	rng    *rand.Rand
}

// NewWordCloud creates a new word cloud.
func NewWordCloud(title string, config WordCloudConfig) *WordCloud {
	return &WordCloud{
		config: config,
		words:  []WordCloudWord{},
		title:  title,
		rng:    rand.New(rand.NewSource(config.Seed)),
	}
}

// AddWord adds a word to the cloud.
func (wc *WordCloud) AddWord(text string, frequency int, entityType semantic.EntityType) {
	// Calculate size based on frequency (logarithmic scale)
	size := wc.calculateSize(frequency)

	wc.words = append(wc.words, WordCloudWord{
		Text:      text,
		Frequency: frequency,
		Size:      size,
		Type:      entityType,
	})
}

// calculateSize computes font size from frequency using logarithmic scale.
func (wc *WordCloud) calculateSize(frequency int) int {
	if frequency <= 0 {
		return wc.config.MinFontSize
	}

	// Small: 1-3 occurrences = size 1
	// Medium: 4-10 occurrences = size 2
	// Large: 11+ occurrences = size 3
	if frequency <= 3 {
		return 1
	} else if frequency <= 10 {
		return 2
	} else {
		return 3
	}
}

// Layout performs word cloud layout algorithm.
// Uses a simple spiral placement strategy to avoid overlaps.
func (wc *WordCloud) Layout() {
	if len(wc.words) == 0 {
		return
	}

	// Limit to max words
	words := wc.words
	if len(words) > wc.config.MaxWords {
		words = words[:wc.config.MaxWords]
	}

	// Sort by size descending (place larger words first)
	sort.Slice(words, func(i, j int) bool {
		if words[i].Size != words[j].Size {
			return words[i].Size > words[j].Size
		}
		return words[i].Frequency > words[j].Frequency
	})

	// Create canvas for collision detection
	canvas := make([][]bool, wc.config.Height)
	for i := range canvas {
		canvas[i] = make([]bool, wc.config.Width)
	}

	// Place words using spiral algorithm
	centerX := wc.config.Width / 2
	centerY := wc.config.Height / 2

	for i := range words {
		placed := false

		// Try center first for first word
		if i == 0 {
			x := centerX - len(words[i].Text)/2
			y := centerY
			if wc.canPlace(&words[i], x, y, canvas) {
				wc.placeWord(&words[i], x, y, canvas)
				placed = true
			}
		}

		// Try spiral positions
		if !placed {
			for radius := 1; radius < max(wc.config.Width, wc.config.Height); radius++ {
				// Try positions along spiral
				for angle := 0; angle < 360; angle += 30 {
					rad := float64(angle) * math.Pi / 180.0
					x := centerX + int(float64(radius)*math.Cos(rad)) - len(words[i].Text)/2
					y := centerY + int(float64(radius)*0.5*math.Sin(rad))

					if wc.canPlace(&words[i], x, y, canvas) {
						wc.placeWord(&words[i], x, y, canvas)
						placed = true
						break
					}
				}
				if placed {
					break
				}
			}
		}

		// If still not placed, skip this word
		if !placed {
			words[i].X = -1
			words[i].Y = -1
		}
	}

	wc.words = words
}

// canPlace checks if a word can be placed at the given position without overlap.
func (wc *WordCloud) canPlace(word *WordCloudWord, x, y int, canvas [][]bool) bool {
	// Check bounds
	if y < 0 || y >= wc.config.Height {
		return false
	}
	if x < 0 || x+len(word.Text) >= wc.config.Width {
		return false
	}

	// Check for collisions
	for i := 0; i < len(word.Text); i++ {
		if canvas[y][x+i] {
			return false
		}
	}

	// Check vertical space for multi-line words (size 2+)
	if word.Size > 1 && y > 0 {
		for i := 0; i < len(word.Text); i++ {
			if canvas[y-1][x+i] {
				return false
			}
		}
	}

	return true
}

// placeWord marks the word's position on the canvas.
func (wc *WordCloud) placeWord(word *WordCloudWord, x, y int, canvas [][]bool) {
	word.X = x
	word.Y = y

	// Mark canvas cells as occupied
	for i := 0; i < len(word.Text); i++ {
		canvas[y][x+i] = true
	}

	// Mark additional space for larger words
	if word.Size > 1 && y > 0 {
		for i := 0; i < len(word.Text); i++ {
			canvas[y-1][x+i] = true
		}
	}
}

// Render renders the word cloud as a string.
func (wc *WordCloud) Render() string {
	// Perform layout if not done
	wc.Layout()

	// Create rendering canvas
	canvas := make([][]string, wc.config.Height)
	for i := range canvas {
		canvas[i] = make([]string, wc.config.Width)
		for j := range canvas[i] {
			canvas[i][j] = " "
		}
	}

	// Render words to canvas
	for _, word := range wc.words {
		if word.X < 0 || word.Y < 0 {
			continue // Word not placed
		}

		// Get color for entity type
		color := wc.config.ColorMap[word.Type]
		if color == "" {
			color = lipgloss.Color("15")
		}

		// Create style based on size
		var style lipgloss.Style
		switch word.Size {
		case 1:
			style = lipgloss.NewStyle().Foreground(color)
		case 2:
			style = lipgloss.NewStyle().Foreground(color).Bold(true)
		case 3:
			style = lipgloss.NewStyle().Foreground(color).Bold(true).Underline(true)
		default:
			style = lipgloss.NewStyle().Foreground(color)
		}

		// Render word text
		text := word.Text
		if word.Size == 3 {
			text = strings.ToUpper(text)
		}

		// Place on canvas
		for i, ch := range text {
			if word.X+i < wc.config.Width && word.Y < wc.config.Height {
				canvas[word.Y][word.X+i] = style.Render(string(ch))
			}
		}
	}

	// Build output string
	var builder strings.Builder

	// Render title
	if wc.title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))
		builder.WriteString(titleStyle.Render(wc.title))
		builder.WriteString("\n\n")
	}

	// Render canvas
	for y := 0; y < wc.config.Height; y++ {
		for x := 0; x < wc.config.Width; x++ {
			builder.WriteString(canvas[y][x])
		}
		builder.WriteString("\n")
	}

	// Render footer
	footerStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("244"))
	builder.WriteString("\n")
	builder.WriteString(footerStyle.Render("Font size = log(frequency)"))

	return builder.String()
}

// RenderBox renders the word cloud with a border box.
func (wc *WordCloud) RenderBox() string {
	content := wc.Render()

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2)

	return boxStyle.Render(content)
}

// SetTitle sets the cloud title.
func (wc *WordCloud) SetTitle(title string) {
	wc.title = title
}

// Clear clears all words from the cloud.
func (wc *WordCloud) Clear() {
	wc.words = nil
}

// WordCount returns the number of words in the cloud.
func (wc *WordCloud) WordCount() int {
	return len(wc.words)
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
