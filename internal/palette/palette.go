package palette

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rand/pedantic-raven/internal/layout"
	"github.com/rand/pedantic-raven/internal/overlay"
)

// Palette is an overlay that provides command search and execution.
type Palette struct {
	*overlay.BaseOverlay
	registry *CommandRegistry
	query    string
	matches  []MatchScore
	selected int
	maxItems int // Maximum items to display
}

// NewPalette creates a new command palette overlay.
func NewPalette(id overlay.OverlayID, registry *CommandRegistry) *Palette {
	return &Palette{
		BaseOverlay: overlay.NewBaseOverlay(id, true, overlay.CenterPosition{}, 80, 20),
		registry:    registry,
		query:       "",
		matches:     registry.FuzzyMatch(""),
		selected:    0,
		maxItems:    10,
	}
}

// Init implements overlay.Overlay.
func (p *Palette) Init() tea.Cmd {
	return nil
}

// Update implements overlay.Overlay.
func (p *Palette) Update(msg tea.Msg) (overlay.Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return p.handleKey(msg)
	}
	return p, nil
}

func (p *Palette) handleKey(msg tea.KeyMsg) (overlay.Overlay, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Dismiss palette
		return p, func() tea.Msg {
			return overlay.DismissOverlay{ID: p.ID()}
		}

	case "enter":
		// Execute selected command
		if p.selected >= 0 && p.selected < len(p.matches) {
			cmd := p.matches[p.selected].Command
			dismissCmd := func() tea.Msg {
				return overlay.DismissOverlay{ID: p.ID()}
			}

			if cmd.Execute != nil {
				return p, tea.Batch(dismissCmd, cmd.Execute())
			}
			return p, dismissCmd
		}
		return p, nil

	case "up", "ctrl+p":
		// Move selection up
		if p.selected > 0 {
			p.selected--
		}
		return p, nil

	case "down", "ctrl+n":
		// Move selection down
		if p.selected < len(p.matches)-1 {
			p.selected++
		}
		return p, nil

	case "backspace":
		// Delete last character
		if len(p.query) > 0 {
			p.query = p.query[:len(p.query)-1]
			p.matches = p.registry.FuzzyMatch(p.query)
			p.selected = 0
		}
		return p, nil

	case "ctrl+u":
		// Clear query
		p.query = ""
		p.matches = p.registry.FuzzyMatch("")
		p.selected = 0
		return p, nil

	default:
		// Add character to query
		if msg.Type == tea.KeyRunes {
			p.query += string(msg.Runes)
			p.matches = p.registry.FuzzyMatch(p.query)
			p.selected = 0
		}
		return p, nil
	}
}

// View implements overlay.Overlay.
func (p *Palette) View(area layout.Rect) string {
	var b strings.Builder

	// Title
	b.WriteString("╭─ Command Palette ─╮\n")

	// Query input
	b.WriteString("│ > ")
	b.WriteString(p.query)
	b.WriteString("█") // Cursor
	b.WriteString("\n")

	// Separator
	b.WriteString("├───────────────────┤\n")

	// Results (limited to maxItems)
	displayCount := min(len(p.matches), p.maxItems)

	if displayCount == 0 {
		b.WriteString("│ No matches found\n")
	} else {
		for i := 0; i < displayCount; i++ {
			match := p.matches[i]

			// Selection indicator
			if i == p.selected {
				b.WriteString("│ > ")
			} else {
				b.WriteString("│   ")
			}

			// Command name
			b.WriteString(match.Command.Name)

			// Keybinding (if present)
			if match.Command.Keybinding != "" {
				b.WriteString(" [")
				b.WriteString(match.Command.Keybinding)
				b.WriteString("]")
			}

			b.WriteString("\n")

			// Description (indented)
			if match.Command.Description != "" {
				b.WriteString("│     ")
				b.WriteString(match.Command.Description)
				b.WriteString("\n")
			}
		}

		// Show count if there are more
		if len(p.matches) > displayCount {
			remaining := len(p.matches) - displayCount
			b.WriteString(fmt.Sprintf("│ ... %d more\n", remaining))
		}
	}

	// Bottom border
	b.WriteString("╰───────────────────╯")

	return b.String()
}

// OnDismiss implements overlay.Overlay.
func (p *Palette) OnDismiss() tea.Cmd {
	return nil
}

// --- Helper ---

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
