package memorylist

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Styles for the memory list.
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")). // Blue
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")). // Yellow
			Background(lipgloss.Color("237")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")) // Light gray

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")) // Darker gray

	importanceColors = map[uint32]lipgloss.Color{
		9: lipgloss.Color("196"), // Red (critical)
		8: lipgloss.Color("208"), // Orange
		7: lipgloss.Color("226"), // Yellow
		6: lipgloss.Color("118"), // Green
		5: lipgloss.Color("81"),  // Cyan
	}

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	loadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true)
)

// View implements tea.Model.
func (m Model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Search input (if in search mode)
	if m.inputMode == InputModeSearch {
		b.WriteString(m.renderSearchInput())
		b.WriteString("\n")
	}

	// Content area
	if m.showHelp {
		b.WriteString(m.renderHelp())
	} else if m.loading {
		b.WriteString(m.renderLoading())
	} else if m.err != nil {
		b.WriteString(m.renderError())
	} else if len(m.filteredMems) == 0 {
		b.WriteString(m.renderEmpty())
	} else {
		b.WriteString(m.renderMemories())
	}

	// Footer
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderHeader renders the header bar.
func (m Model) renderHeader() string {
	title := "Memory Workspace"
	if m.searchQuery != "" {
		title += fmt.Sprintf(" - Search: %q", m.searchQuery)
	}

	sortIndicator := fmt.Sprintf(" [Sort: %s", m.sortBy.String())
	if m.sortDesc {
		sortIndicator += " ↓]"
	} else {
		sortIndicator += " ↑]"
	}

	titleWidth := m.width - len(sortIndicator) - 4
	if titleWidth < len(title) {
		title = title[:titleWidth] + "..."
	}

	header := title + strings.Repeat(" ", titleWidth-len(title)) + sortIndicator

	return headerStyle.Width(m.width).Render(header)
}

// renderMemories renders the list of memories.
func (m Model) renderMemories() string {
	var b strings.Builder

	visibleLines := m.visibleLines()
	start := m.scrollOffset
	end := m.scrollOffset + visibleLines
	if end > len(m.filteredMems) {
		end = len(m.filteredMems)
	}

	for i := start; i < end; i++ {
		mem := m.filteredMems[i]
		isSelected := i == m.selectedIndex

		b.WriteString(m.renderMemoryRow(mem, isSelected, i))
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	// Pad with empty lines if needed
	renderedLines := end - start
	for i := renderedLines; i < visibleLines; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

// renderMemoryRow renders a single memory row (3 lines).
func (m Model) renderMemoryRow(mem *pb.MemoryNote, selected bool, index int) string {
	var lines []string

	// Line 1: Title + Importance
	titleLine := m.renderTitleLine(mem, selected)
	lines = append(lines, titleLine)

	// Line 2: Namespace + Tags
	metaLine := m.renderMetaLine(mem, selected)
	lines = append(lines, metaLine)

	// Line 3: Updated time + Link count
	footerLine := m.renderMemoryFooter(mem, selected)
	lines = append(lines, footerLine)

	return strings.Join(lines, "\n")
}

// renderTitleLine renders the title and importance.
func (m Model) renderTitleLine(mem *pb.MemoryNote, selected bool) string {
	// Extract title from content (first line, max 50 chars)
	title := extractTitle(mem.Content)

	// Importance indicator
	impStr := m.renderImportance(mem.Importance)

	// Selector indicator
	selector := "  "
	if selected {
		selector = "> "
	}

	// Calculate available width for title
	availableWidth := m.width - len(selector) - len(impStr) - 2

	if len(title) > availableWidth {
		title = title[:availableWidth-3] + "..."
	}

	// Pad title
	title = title + strings.Repeat(" ", availableWidth-len(title))

	line := selector + title + " " + impStr

	if selected {
		return selectedStyle.Width(m.width).Render(line)
	}
	return normalStyle.Width(m.width).Render(line)
}

// renderMetaLine renders namespace and tags.
func (m Model) renderMetaLine(mem *pb.MemoryNote, selected bool) string {
	parts := []string{}

	// Namespace
	ns := formatNamespace(mem.Namespace)
	if ns != "" {
		parts = append(parts, ns)
	}

	// Tags
	if len(mem.Tags) > 0 {
		tags := strings.Join(mem.Tags, ", ")
		if len(tags) > 50 {
			tags = tags[:47] + "..."
		}
		parts = append(parts, tags)
	}

	line := "  " + strings.Join(parts, " • ")

	// Truncate if too long
	if len(line) > m.width {
		line = line[:m.width-3] + "..."
	}

	style := metaStyle
	if selected {
		style = selectedStyle.Foreground(lipgloss.Color("243"))
	}

	return style.Width(m.width).Render(line)
}

// renderMemoryFooter renders the timestamp and link count.
func (m Model) renderMemoryFooter(mem *pb.MemoryNote, selected bool) string {
	parts := []string{}

	// Updated time
	updated := formatRelativeTime(time.Unix(int64(mem.UpdatedAt), 0))
	parts = append(parts, "Updated "+updated)

	// Link count
	linkCount := len(mem.Links)
	if linkCount == 1 {
		parts = append(parts, "1 link")
	} else if linkCount > 0 {
		parts = append(parts, fmt.Sprintf("%d links", linkCount))
	}

	line := "  " + strings.Join(parts, " • ")

	style := metaStyle
	if selected {
		style = selectedStyle.Foreground(lipgloss.Color("243"))
	}

	return style.Width(m.width).Render(line)
}

// renderImportance renders the importance indicator.
func (m Model) renderImportance(imp uint32) string {
	if imp < 1 || imp > 10 {
		return "[Imp: -]"
	}

	color, ok := importanceColors[imp]
	if !ok {
		color = lipgloss.Color("250") // Default gray
	}

	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	return style.Render(fmt.Sprintf("[Imp: %d]", imp))
}

// renderLoading renders the loading state.
func (m Model) renderLoading() string {
	lines := m.visibleLines()
	padTop := (lines - 1) / 2

	var b strings.Builder
	for i := 0; i < padTop; i++ {
		b.WriteString("\n")
	}

	loading := loadingStyle.Render("Loading memories...")
	padding := (m.width - len("Loading memories...")) / 2

	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(loading)

	for i := padTop + 1; i < lines; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

// renderError renders the error state.
func (m Model) renderError() string {
	lines := m.visibleLines()
	padTop := (lines - 2) / 2

	var b strings.Builder
	for i := 0; i < padTop; i++ {
		b.WriteString("\n")
	}

	errLine1 := errorStyle.Render("Error loading memories")
	errLine2 := metaStyle.Render(m.err.Error())

	padding1 := (m.width - len("Error loading memories")) / 2
	padding2 := (m.width - len(m.err.Error())) / 2

	b.WriteString(strings.Repeat(" ", padding1))
	b.WriteString(errLine1)
	b.WriteString("\n")
	b.WriteString(strings.Repeat(" ", padding2))
	b.WriteString(errLine2)

	for i := padTop + 2; i < lines; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

// renderEmpty renders the empty state.
func (m Model) renderEmpty() string {
	lines := m.visibleLines()
	padTop := (lines - 1) / 2

	var b strings.Builder
	for i := 0; i < padTop; i++ {
		b.WriteString("\n")
	}

	message := "No memories found"
	if m.searchQuery != "" {
		message = fmt.Sprintf("No memories matching %q", m.searchQuery)
	}

	empty := emptyStyle.Render(message)
	padding := (m.width - len(message)) / 2

	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(empty)

	for i := padTop + 1; i < lines; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

// renderFooter renders the footer status bar.
func (m Model) renderFooter() string {
	if len(m.filteredMems) == 0 {
		footer := fmt.Sprintf("0 memories | j/k: nav, /: search, ?: help")
		return footerStyle.Width(m.width).Render(footer)
	}

	// Calculate visible range
	start := m.scrollOffset + 1
	end := m.scrollOffset + m.visibleLines()
	if end > len(m.filteredMems) {
		end = len(m.filteredMems)
	}

	showing := fmt.Sprintf("Showing %d-%d of %d", start, end, len(m.filteredMems))
	if m.totalCount > uint32(len(m.filteredMems)) {
		showing += fmt.Sprintf(" (total: %d)", m.totalCount)
	}

	keys := " | j/k: nav, Enter: view, /: search"

	footer := showing + keys

	// Truncate if too long
	if len(footer) > m.width {
		footer = showing
	}

	return footerStyle.Width(m.width).Render(footer)
}

// Helper functions

// extractTitle extracts a title from memory content (first line, trimmed).
func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "(Empty)"
	}

	title := strings.TrimSpace(lines[0])
	if title == "" {
		return "(Untitled)"
	}

	// Remove markdown heading markers
	title = strings.TrimPrefix(title, "# ")
	title = strings.TrimPrefix(title, "## ")
	title = strings.TrimPrefix(title, "### ")

	return title
}

// renderSearchInput renders the search input bar.
func (m Model) renderSearchInput() string {
	prompt := "Search: "
	cursor := "█"

	input := m.searchInput + cursor

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("237"))

	line := prompt + input

	// Pad to full width
	if len(line) < m.width {
		line += strings.Repeat(" ", m.width-len(line))
	}

	return style.Width(m.width).Render(line)
}

// renderHelp renders the help overlay.
func (m Model) renderHelp() string {
	lines := m.visibleLines()
	helpText := []string{
		"Keyboard Shortcuts",
		"",
		"Navigation:",
		"  j/↓       Move down",
		"  k/↑       Move up",
		"  g         Go to top",
		"  G         Go to bottom",
		"  Ctrl+D    Page down",
		"  Ctrl+U    Page up",
		"  Enter     Select memory",
		"",
		"Search & Filters:",
		"  /         Enter search mode",
		"  c         Clear search & filters",
		"  r         Reload/refresh",
		"",
		"Other:",
		"  ?         Toggle this help",
		"  Esc       Close help/error",
		"",
		"Press ? to close",
	}

	var b strings.Builder

	// Calculate padding
	startLine := (lines - len(helpText)) / 2
	if startLine < 0 {
		startLine = 0
	}

	// Top padding
	for i := 0; i < startLine; i++ {
		b.WriteString("\n")
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	for i, line := range helpText {
		padding := (m.width - len(line)) / 2
		if padding < 0 {
			padding = 0
		}

		b.WriteString(strings.Repeat(" ", padding))

		if i == 0 {
			b.WriteString(titleStyle.Render(line))
		} else {
			b.WriteString(helpStyle.Render(line))
		}

		if i < len(helpText)-1 {
			b.WriteString("\n")
		}
	}

	// Bottom padding
	for i := startLine + len(helpText); i < lines; i++ {
		b.WriteString("\n")
	}

	return b.String()
}

// formatRelativeTime formats a timestamp as relative time.
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	case diff < 365*24*time.Hour:
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(diff.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
