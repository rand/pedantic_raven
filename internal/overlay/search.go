package overlay

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/editor/search"
	"github.com/rand/pedantic-raven/internal/layout"
)

var (
	searchOverlayStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1, 2)

	searchTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170"))

	searchLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("244"))

	searchInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	searchActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Bold(true)

	searchOptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39"))

	searchHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)
)

// SearchAction represents the action performed in the search overlay.
type SearchAction int

const (
	SearchActionNone SearchAction = iota
	SearchActionFind
	SearchActionFindNext
	SearchActionFindPrevious
	SearchActionReplace
	SearchActionReplaceAll
)

// SearchResult is sent when a search action is performed.
type SearchResult struct {
	Action      SearchAction
	Query       string
	Replacement string
	Options     search.SearchOptions
	Canceled    bool
}

// searchField represents which field is currently active.
type searchField int

const (
	fieldQuery searchField = iota
	fieldReplacement
)

// SearchOverlay is an overlay for performing search and replace operations.
type SearchOverlay struct {
	*BaseOverlay
	mode            string // "search" or "replace"
	activeField     searchField
	queryText       string
	replacementText string
	opts            search.SearchOptions
	matchCount      int
	currentMatch    int
}

// NewSearchOverlay creates a new search overlay.
// mode can be "search" or "replace".
func NewSearchOverlay(id OverlayID, mode string) *SearchOverlay {
	if mode != "search" && mode != "replace" {
		mode = "search"
	}

	return &SearchOverlay{
		BaseOverlay:     NewBaseOverlay(id, true, CenterPosition{}, 60, 15),
		mode:            mode,
		activeField:     fieldQuery,
		queryText:       "",
		replacementText: "",
		opts:            search.DefaultSearchOptions(),
		matchCount:      -1, // -1 indicates no search performed yet
		currentMatch:    -1,
	}
}

// SetMatchInfo sets the match count and current match index for display.
func (so *SearchOverlay) SetMatchInfo(matchCount, currentMatch int) {
	so.matchCount = matchCount
	so.currentMatch = currentMatch
}

// Update implements Overlay.
func (so *SearchOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel and dismiss
			dismissCmd := func() tea.Msg {
				return DismissOverlay{ID: so.id}
			}
			resultCmd := func() tea.Msg {
				return SearchResult{Canceled: true}
			}
			return so, tea.Batch(resultCmd, dismissCmd)

		case "tab":
			// Switch between fields (in replace mode)
			if so.mode == "replace" {
				if so.activeField == fieldQuery {
					so.activeField = fieldReplacement
				} else {
					so.activeField = fieldQuery
				}
			}

		case "ctrl+c":
			// Toggle case sensitive
			so.opts.CaseSensitive = !so.opts.CaseSensitive
			return so, so.emitSearchAction(SearchActionFind)

		case "ctrl+w":
			// Toggle whole word
			so.opts.WholeWord = !so.opts.WholeWord
			return so, so.emitSearchAction(SearchActionFind)

		case "ctrl+r":
			// Toggle regex
			so.opts.Regex = !so.opts.Regex
			return so, so.emitSearchAction(SearchActionFind)

		case "enter":
			// Perform action based on mode and field
			if so.mode == "search" {
				return so, so.emitSearchAction(SearchActionFind)
			} else {
				// Replace mode
				if so.activeField == fieldQuery {
					// On query field, just do a find
					return so, so.emitSearchAction(SearchActionFind)
				} else {
					// On replacement field, do replace current
					return so, so.emitSearchAction(SearchActionReplace)
				}
			}

		case "f3":
			// Find next
			return so, so.emitSearchAction(SearchActionFindNext)

		case "shift+f3":
			// Find previous
			return so, so.emitSearchAction(SearchActionFindPrevious)

		case "ctrl+a":
			// Replace all (only in replace mode)
			if so.mode == "replace" {
				return so, so.emitSearchAction(SearchActionReplaceAll)
			}

		case "backspace":
			// Delete character from active field
			if so.activeField == fieldQuery {
				if len(so.queryText) > 0 {
					so.queryText = so.queryText[:len(so.queryText)-1]
				}
			} else {
				if len(so.replacementText) > 0 {
					so.replacementText = so.replacementText[:len(so.replacementText)-1]
				}
			}

		default:
			// Add character to active field
			if msg.Type == tea.KeyRunes {
				if so.activeField == fieldQuery {
					so.queryText += string(msg.Runes)
				} else {
					so.replacementText += string(msg.Runes)
				}
			}
		}
	}

	return so, nil
}

// emitSearchAction creates a command that emits a SearchResult message.
func (so *SearchOverlay) emitSearchAction(action SearchAction) tea.Cmd {
	return func() tea.Msg {
		return SearchResult{
			Action:      action,
			Query:       so.queryText,
			Replacement: so.replacementText,
			Options:     so.opts,
			Canceled:    false,
		}
	}
}

// View implements Overlay.
func (so *SearchOverlay) View(area layout.Rect) string {
	var content strings.Builder

	// Title
	title := "Search"
	if so.mode == "replace" {
		title = "Search and Replace"
	}
	content.WriteString(searchTitleStyle.Render(title) + "\n\n")

	// Query field
	queryLabel := searchLabelStyle.Render("Find: ")
	queryCursor := ""
	if so.activeField == fieldQuery {
		queryCursor = "█"
	}
	queryDisplay := so.queryText + queryCursor
	if so.activeField == fieldQuery {
		queryDisplay = searchActiveStyle.Render(queryDisplay)
	} else {
		queryDisplay = searchInputStyle.Render(queryDisplay)
	}
	content.WriteString(queryLabel + queryDisplay + "\n")

	// Replacement field (only in replace mode)
	if so.mode == "replace" {
		replaceLabel := searchLabelStyle.Render("Replace: ")
		replaceCursor := ""
		if so.activeField == fieldReplacement {
			replaceCursor = "█"
		}
		replaceDisplay := so.replacementText + replaceCursor
		if so.activeField == fieldReplacement {
			replaceDisplay = searchActiveStyle.Render(replaceDisplay)
		} else {
			replaceDisplay = searchInputStyle.Render(replaceDisplay)
		}
		content.WriteString(replaceLabel + replaceDisplay + "\n")
	}

	content.WriteString("\n")

	// Search options
	content.WriteString(searchLabelStyle.Render("Options:") + "\n")

	caseOpt := "[ ]"
	if so.opts.CaseSensitive {
		caseOpt = "[✓]"
	}
	content.WriteString(searchOptionStyle.Render(fmt.Sprintf("  %s Case Sensitive (Ctrl+C)", caseOpt)) + "\n")

	wordOpt := "[ ]"
	if so.opts.WholeWord {
		wordOpt = "[✓]"
	}
	content.WriteString(searchOptionStyle.Render(fmt.Sprintf("  %s Whole Word (Ctrl+W)", wordOpt)) + "\n")

	regexOpt := "[ ]"
	if so.opts.Regex {
		regexOpt = "[✓]"
	}
	content.WriteString(searchOptionStyle.Render(fmt.Sprintf("  %s Regex (Ctrl+R)", regexOpt)) + "\n")

	content.WriteString("\n")

	// Match count
	if so.matchCount >= 0 {
		matchInfo := ""
		if so.matchCount == 0 {
			matchInfo = "No matches found"
		} else if so.currentMatch >= 0 {
			matchInfo = fmt.Sprintf("Match %d of %d", so.currentMatch+1, so.matchCount)
		} else {
			matchInfo = fmt.Sprintf("%d matches", so.matchCount)
		}
		content.WriteString(searchLabelStyle.Render(matchInfo) + "\n\n")
	}

	// Help text
	content.WriteString(strings.Repeat("─", area.Width-6) + "\n")

	if so.mode == "search" {
		content.WriteString(searchHelpStyle.Render("Enter: Find | F3: Next | Shift+F3: Prev | Esc: Close"))
	} else {
		content.WriteString(searchHelpStyle.Render("Tab: Switch | Enter: Replace | Ctrl+A: Replace All\n"))
		content.WriteString(searchHelpStyle.Render("F3: Next | Shift+F3: Prev | Esc: Close"))
	}

	return searchOverlayStyle.
		Width(area.Width - 4).
		Height(area.Height - 4).
		Render(content.String())
}
