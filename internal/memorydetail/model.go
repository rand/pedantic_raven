package memorydetail

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case MemoryLoadedMsg:
		m.SetMemory(msg.Memory)
		return m, nil

	case MemoryErrorMsg:
		m.SetError(msg.Err)
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		return m.scrollDown(), nil

	case "k", "up":
		return m.scrollUp(), nil

	case "g":
		return m.scrollToTop(), nil

	case "G":
		return m.scrollToBottom(), nil

	case "ctrl+d":
		return m.pageDown(), nil

	case "ctrl+u":
		return m.pageUp(), nil

	case "m":
		m.ToggleMetadata()
		return m, nil

	case "l":
		// Enter link navigation mode (select first link)
		if m.HasLinks() {
			m.SelectFirstLink()
		}
		return m, nil

	case "tab", "n":
		// Navigate to next link
		m.SelectNextLink()
		return m, nil

	case "shift+tab", "p":
		// Navigate to previous link
		m.SelectPreviousLink()
		return m, nil

	case "q":
		return m, m.closeCmd()

	case "esc":
		// Clear link selection if active, otherwise close
		if m.selectedLinkIndex >= 0 {
			m.ClearLinkSelection()
			return m, nil
		}
		return m, m.closeCmd()

	case "enter":
		// Navigate to selected link
		if link := m.SelectedLink(); link != nil {
			return m, m.navigateLinkCmd(link.TargetId)
		}
		return m, nil
	}

	return m, nil
}

// Scrolling methods

func (m Model) scrollDown() Model {
	if m.memory == nil {
		return m
	}

	maxScroll := m.maxScrollOffset()
	if m.scrollOffset < maxScroll {
		m.scrollOffset++
	}

	return m
}

func (m Model) scrollUp() Model {
	if m.scrollOffset > 0 {
		m.scrollOffset--
	}

	return m
}

func (m Model) scrollToTop() Model {
	m.scrollOffset = 0
	return m
}

func (m Model) scrollToBottom() Model {
	if m.memory == nil {
		return m
	}

	m.scrollOffset = m.maxScrollOffset()
	return m
}

func (m Model) pageDown() Model {
	if m.memory == nil {
		return m
	}

	visibleLines := m.visibleLines()
	m.scrollOffset += visibleLines
	maxScroll := m.maxScrollOffset()
	if m.scrollOffset > maxScroll {
		m.scrollOffset = maxScroll
	}

	return m
}

func (m Model) pageUp() Model {
	visibleLines := m.visibleLines()
	m.scrollOffset -= visibleLines
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	return m
}

// visibleLines returns the number of visible content lines.
func (m Model) visibleLines() int {
	// Subtract header (1) and footer (1)
	return m.height - 2
}

// maxScrollOffset returns the maximum scroll offset.
func (m Model) maxScrollOffset() int {
	if m.memory == nil {
		return 0
	}

	totalLines := m.contentLines()
	visibleLines := m.visibleLines()

	if totalLines <= visibleLines {
		return 0
	}

	return totalLines - visibleLines
}

// contentLines returns the total number of content lines.
func (m Model) contentLines() int {
	if m.memory == nil {
		return 0
	}

	// Count lines in content
	lines := 1
	for _, c := range m.memory.Content {
		if c == '\n' {
			lines++
		}
	}

	return lines
}

// closeCmd returns a command to close the detail view.
func (m Model) closeCmd() tea.Cmd {
	return func() tea.Msg {
		return CloseRequestMsg{}
	}
}

// navigateLinkCmd returns a command to navigate to a linked memory.
func (m Model) navigateLinkCmd(targetID string) tea.Cmd {
	return func() tea.Msg {
		return LinkSelectedMsg{
			TargetID: targetID,
		}
	}
}
