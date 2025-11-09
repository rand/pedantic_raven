package memorydetail

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Styles for the memory detail view.
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")). // Blue
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	contentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")) // Light gray

	metaLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // Blue
			Bold(true)

	metaValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")) // Light gray

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // Red
			Bold(true)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")). // Gray
			Italic(true)
)

// View implements tea.Model.
func (m Model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Content area
	if m.err != nil {
		b.WriteString(m.renderError())
	} else if m.memory == nil {
		b.WriteString(m.renderEmpty())
	} else {
		b.WriteString(m.renderContent())
	}

	// Footer
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderHeader renders the header bar.
func (m Model) renderHeader() string {
	title := "Memory Detail"
	if m.memory != nil {
		title += fmt.Sprintf(" - %s", extractTitle(m.memory.Content))
	}

	// Truncate if too long
	maxWidth := m.width - 4
	if len(title) > maxWidth {
		title = title[:maxWidth-3] + "..."
	}

	// Pad to full width
	padding := m.width - len(title) - 2
	if padding < 0 {
		padding = 0
	}
	title += strings.Repeat(" ", padding)

	return headerStyle.Width(m.width).Render(title)
}

// renderContent renders the memory content.
func (m Model) renderContent() string {
	if m.memory == nil {
		return ""
	}

	var b strings.Builder
	contentWidth := m.width
	metadataWidth := 0

	// Calculate widths if metadata is shown
	if m.showMetadata {
		metadataWidth = m.width / 3
		if metadataWidth < 20 {
			metadataWidth = 20
		}
		contentWidth = m.width - metadataWidth - 1 // -1 for separator
	}

	visibleLines := m.visibleLines()
	lines := strings.Split(m.memory.Content, "\n")

	// Render content lines
	start := m.scrollOffset
	end := m.scrollOffset + visibleLines
	if end > len(lines) {
		end = len(lines)
	}

	for i := 0; i < visibleLines; i++ {
		lineIdx := start + i

		if m.showMetadata {
			// Content on left
			var contentLine string
			if lineIdx < len(lines) {
				contentLine = lines[lineIdx]
				if len(contentLine) > contentWidth {
					contentLine = contentLine[:contentWidth-3] + "..."
				}
			}
			contentLine = padRight(contentLine, contentWidth)

			// Metadata on right
			metaLine := m.renderMetadataLine(i)
			metaLine = padRight(metaLine, metadataWidth)

			// Combine with separator
			b.WriteString(contentStyle.Render(contentLine))
			b.WriteString(" ")
			b.WriteString(metaValueStyle.Render(metaLine))
		} else {
			// Full width content
			var contentLine string
			if lineIdx < len(lines) {
				contentLine = lines[lineIdx]
				if len(contentLine) > contentWidth {
					contentLine = contentLine[:contentWidth-3] + "..."
				}
			}
			b.WriteString(contentStyle.Render(contentLine))
		}

		if i < visibleLines-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderMetadataLine renders a line of the metadata panel.
func (m Model) renderMetadataLine(lineNum int) string {
	if m.memory == nil {
		return ""
	}

	// Metadata items to display
	switch lineNum {
	case 0:
		return metaLabelStyle.Render("Metadata")
	case 1:
		return ""
	case 2:
		return metaLabelStyle.Render("ID: ") + metaValueStyle.Render(m.memory.Id)
	case 3:
		return metaLabelStyle.Render("Importance: ") + metaValueStyle.Render(fmt.Sprintf("%d", m.memory.Importance))
	case 4:
		return ""
	case 5:
		return metaLabelStyle.Render("Namespace:")
	case 6:
		ns := formatNamespace(m.memory.Namespace)
		return "  " + metaValueStyle.Render(ns)
	case 7:
		return ""
	case 8:
		if len(m.memory.Tags) > 0 {
			return metaLabelStyle.Render("Tags:")
		}
	case 9:
		if len(m.memory.Tags) > 0 {
			tags := strings.Join(m.memory.Tags, ", ")
			if len(tags) > 25 {
				tags = tags[:22] + "..."
			}
			return "  " + metaValueStyle.Render(tags)
		}
	case 10:
		return ""
	case 11:
		return metaLabelStyle.Render("Created:")
	case 12:
		created := time.Unix(int64(m.memory.CreatedAt), 0)
		return "  " + metaValueStyle.Render(formatRelativeTime(created))
	case 13:
		return ""
	case 14:
		return metaLabelStyle.Render("Updated:")
	case 15:
		updated := time.Unix(int64(m.memory.UpdatedAt), 0)
		return "  " + metaValueStyle.Render(formatRelativeTime(updated))
	case 16:
		return ""
	case 17:
		if len(m.memory.Links) > 0 {
			return metaLabelStyle.Render(fmt.Sprintf("Links (%d):", len(m.memory.Links)))
		}
	default:
		if lineNum >= 18 && lineNum < 18+len(m.memory.Links) {
			link := m.memory.Links[lineNum-18]
			return "  → " + metaValueStyle.Render(link.TargetId)
		}
	}

	return ""
}

// renderError renders the error state.
func (m Model) renderError() string {
	lines := m.visibleLines()
	padTop := (lines - 2) / 2

	var b strings.Builder
	for i := 0; i < padTop; i++ {
		b.WriteString("\n")
	}

	errLine1 := errorStyle.Render("Error loading memory")
	errLine2 := metaValueStyle.Render(m.err.Error())

	padding1 := (m.width - len("Error loading memory")) / 2
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

	message := "No memory selected"
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
	if m.memory == nil {
		footer := "No memory | q/Esc: close, m: toggle metadata"
		return footerStyle.Width(m.width).Render(footer)
	}

	totalLines := m.contentLines()
	visibleLines := m.visibleLines()

	var scrollInfo string
	if totalLines > visibleLines {
		scrollInfo = fmt.Sprintf("Lines %d-%d of %d", m.scrollOffset+1, m.scrollOffset+visibleLines, totalLines)
	} else {
		scrollInfo = fmt.Sprintf("%d lines", totalLines)
	}

	keys := " | j/k: scroll, m: metadata, q: close"

	footer := scrollInfo + keys

	// Truncate if too long
	if len(footer) > m.width {
		footer = scrollInfo
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

	// Limit length
	if len(title) > 50 {
		title = title[:47] + "..."
	}

	return title
}

// formatNamespace converts a protobuf namespace to a string.
func formatNamespace(ns *pb.Namespace) string {
	if ns == nil {
		return ""
	}

	switch ns := ns.Namespace.(type) {
	case *pb.Namespace_Global:
		return "global"

	case *pb.Namespace_Project:
		return "project:" + ns.Project.Name

	case *pb.Namespace_Session:
		return "project:" + ns.Session.Project + ":session:" + ns.Session.SessionId

	default:
		return ""
	}
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

// padRight pads a string to the right with spaces.
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
