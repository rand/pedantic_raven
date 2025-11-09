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

	case EditModeEnteredMsg:
		// Enter edit mode
		m.editState = &EditState{
			isEditing:    true,
			editedMemory: msg.Memory,
			fieldFocus:   FieldContent,
			originalHash: hashMemory(msg.Memory),
			hasChanges:   false,
		}
		return m, nil

	case MemorySavedMsg:
		// Handle save result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}

		// Exit edit mode and update memory
		m.editState = nil
		m.isNewMemory = false
		m.SetMemory(msg.Memory)
		return m, nil

	case MemoryCreatedMsg:
		// Handle create result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}

		// Exit edit mode and show created memory
		m.editState = nil
		m.isNewMemory = false
		m.SetMemory(msg.Memory)
		return m, nil

	case MemoryUpdatedMsg:
		// Handle update result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}

		// Exit edit mode and update memory
		m.editState = nil
		m.SetMemory(msg.Memory)
		return m, nil

	case MemoryDeletedMsg:
		// Handle delete result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}

		// Clear memory and close
		m.SetMemory(nil)
		return m, m.closeCmd()

	case DeleteConfirmationRequestMsg:
		// Show delete confirmation dialog
		m.showDeleteConfirm = true
		return m, nil

	case LinkCreatedMsg:
		// Handle link creation result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}

		// Add link to current memory's links
		if m.memory != nil && msg.Link != nil {
			m.memory.Links = append(m.memory.Links, msg.Link)
		}

		// Hide create link dialog
		m.HideCreateLinkDialog()
		return m, nil

	case LinkDeletedMsg:
		// Handle link deletion result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}

		// Remove link from current memory's links
		if m.memory != nil {
			for i, link := range m.memory.Links {
				if link.TargetId == msg.LinkID {
					m.memory.Links = append(m.memory.Links[:i], m.memory.Links[i+1:]...)
					break
				}
			}
		}

		// Clear link selection
		m.ClearLinkSelection()
		return m, nil

	case LinkMetadataUpdatedMsg:
		// Handle link metadata update result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}

		// Update link in current memory's links
		if m.memory != nil && msg.Link != nil {
			for i, link := range m.memory.Links {
				if link.TargetId == msg.Link.TargetId {
					m.memory.Links[i] = msg.Link
					break
				}
			}
		}
		return m, nil

	case LinkedMemoriesLoadedMsg:
		// Handle linked memories result
		if msg.Err != nil {
			m.SetError(msg.Err)
			return m, nil
		}
		// For now, just acknowledge - could be used for link previews
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle delete confirmation dialog
	if m.showDeleteConfirm {
		switch msg.String() {
		case "enter", "y":
			return m, m.ConfirmDelete()
		case "esc", "n":
			m.CancelDelete()
			return m, nil
		}
		return m, nil
	}

	// Handle edit mode
	if m.IsEditing() {
		switch msg.String() {
		case "ctrl+s":
			// Save changes
			return m, m.SaveChanges()

		case "esc":
			// Cancel editing (with unsaved changes warning handled in UI)
			if m.HasUnsavedChanges() {
				// TODO: Show confirmation dialog
				// For now, just cancel
			}
			m.CancelEdit()
			return m, nil

		case "tab":
			// Cycle field focus
			m.CycleFieldFocus()
			return m, nil
		}
		// Other edit mode keys would be handled here
		// (field editing, etc.)
		return m, nil
	}

	// Normal view mode
	switch msg.String() {
	case "j", "down":
		// If a link is selected, move to next link, otherwise scroll
		if m.selectedLinkIndex >= 0 {
			m.SelectNextLink()
			return m, nil
		}
		return m.scrollDown(), nil

	case "k", "up":
		// If a link is selected, move to previous link, otherwise scroll
		if m.selectedLinkIndex >= 0 {
			m.SelectPreviousLink()
			return m, nil
		}
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

	case "e":
		// Enter edit mode
		return m, m.EnterEditMode()

	case "d":
		// Delete memory (with confirmation)
		return m, m.DeleteCurrentMemory()

	case "l":
		// Toggle link selection mode - select first link if not selected
		if m.HasLinks() {
			m.SelectFirstLink()
		}
		return m, nil

	case "c":
		// Show create link dialog
		if !m.ShowingCreateLinkDialog() {
			m.ShowCreateLinkDialog()
		}
		return m, nil

	case "x":
		// Delete selected link
		if m.selectedLinkIndex >= 0 {
			return m, m.DeleteSelectedLink()
		}
		return m, nil

	case "[":
		// Navigate back
		return m, m.NavigateBack()

	case "]":
		// Navigate forward
		return m, m.NavigateForward()

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
			return m, m.NavigateToLinkedMemory(link.TargetId)
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
