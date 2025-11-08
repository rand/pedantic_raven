package memorylist

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case MemoriesLoadedMsg:
		m.SetMemories(msg.Memories, msg.TotalCount)
		return m, nil

	case MemoriesErrorMsg:
		m.SetError(msg.Err)
		return m, nil

	case SearchQueryMsg:
		m.searchQuery = msg.Query
		m.applyFilters()
		return m, nil

	case FilterChangedMsg:
		m.filterTags = msg.Tags
		m.filterNS = msg.Namespace
		m.minImportance = msg.MinImportance
		m.applyFilters()
		return m, nil

	case SortChangedMsg:
		m.sortBy = msg.Mode
		m.sortDesc = msg.Desc
		m.applySorting()
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		return m.moveDown(), nil

	case "k", "up":
		return m.moveUp(), nil

	case "g":
		return m.moveToTop(), nil

	case "G":
		return m.moveToBottom(), nil

	case "ctrl+d":
		return m.pageDown(), nil

	case "ctrl+u":
		return m.pageUp(), nil

	case "enter":
		return m, m.selectCurrent()
	}

	return m, nil
}

// Navigation methods

func (m Model) moveDown() Model {
	if len(m.filteredMems) == 0 {
		return m
	}

	m.selectedIndex++
	if m.selectedIndex >= len(m.filteredMems) {
		m.selectedIndex = len(m.filteredMems) - 1
	}

	// Auto-scroll
	visibleStart := m.scrollOffset
	visibleEnd := m.scrollOffset + m.visibleLines()

	if m.selectedIndex >= visibleEnd {
		m.scrollOffset = m.selectedIndex - m.visibleLines() + 1
	}

	return m
}

func (m Model) moveUp() Model {
	if len(m.filteredMems) == 0 {
		return m
	}

	m.selectedIndex--
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}

	// Auto-scroll
	if m.selectedIndex < m.scrollOffset {
		m.scrollOffset = m.selectedIndex
	}

	return m
}

func (m Model) moveToTop() Model {
	m.selectedIndex = 0
	m.scrollOffset = 0
	return m
}

func (m Model) moveToBottom() Model {
	if len(m.filteredMems) == 0 {
		return m
	}

	m.selectedIndex = len(m.filteredMems) - 1
	visibleLines := m.visibleLines()
	if len(m.filteredMems) > visibleLines {
		m.scrollOffset = len(m.filteredMems) - visibleLines
	}
	return m
}

func (m Model) pageDown() Model {
	if len(m.filteredMems) == 0 {
		return m
	}

	visibleLines := m.visibleLines()
	m.selectedIndex += visibleLines
	if m.selectedIndex >= len(m.filteredMems) {
		m.selectedIndex = len(m.filteredMems) - 1
	}

	m.scrollOffset += visibleLines
	if m.scrollOffset+visibleLines > len(m.filteredMems) {
		m.scrollOffset = len(m.filteredMems) - visibleLines
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	}

	return m
}

func (m Model) pageUp() Model {
	if len(m.filteredMems) == 0 {
		return m
	}

	visibleLines := m.visibleLines()
	m.selectedIndex -= visibleLines
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}

	m.scrollOffset -= visibleLines
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	return m
}

// visibleLines returns the number of visible memory rows.
func (m Model) visibleLines() int {
	// Subtract header (1) and footer (1)
	return m.height - 2
}

// selectCurrent returns a command to notify selection.
func (m Model) selectCurrent() tea.Cmd {
	selected := m.SelectedMemory()
	if selected == nil {
		return nil
	}

	return func() tea.Msg {
		return MemorySelectedMsg{
			Memory: selected,
			Index:  m.selectedIndex,
		}
	}
}

// Filtering and sorting

// applyFilters applies current filters and re-sorts.
func (m *Model) applyFilters() {
	m.filteredMems = make([]*pb.MemoryNote, 0, len(m.memories))

	for _, mem := range m.memories {
		if m.matchesFilters(mem) {
			m.filteredMems = append(m.filteredMems, mem)
		}
	}

	m.applySorting()

	// Reset selection if out of bounds
	if m.selectedIndex >= len(m.filteredMems) {
		m.selectedIndex = len(m.filteredMems) - 1
		if m.selectedIndex < 0 {
			m.selectedIndex = 0
		}
	}
}

// matchesFilters returns true if the memory matches all active filters.
func (m Model) matchesFilters(mem *pb.MemoryNote) bool {
	// Search query filter (simple contains for now)
	if m.searchQuery != "" {
		query := strings.ToLower(m.searchQuery)
		content := strings.ToLower(mem.Content)
		if !strings.Contains(content, query) {
			return false
		}
	}

	// Namespace filter
	if m.filterNS != "" {
		nsStr := formatNamespace(mem.Namespace)
		if !strings.HasPrefix(nsStr, m.filterNS) {
			return false
		}
	}

	// Tag filter (all tags must be present)
	if len(m.filterTags) > 0 {
		memTags := make(map[string]bool)
		for _, tag := range mem.Tags {
			memTags[tag] = true
		}
		for _, filterTag := range m.filterTags {
			if !memTags[filterTag] {
				return false
			}
		}
	}

	// Importance filter
	if m.minImportance > 0 && mem.Importance < m.minImportance {
		return false
	}

	return true
}

// applySorting sorts the filtered memories by the current sort mode.
func (m *Model) applySorting() {
	sort.SliceStable(m.filteredMems, func(i, j int) bool {
		a, b := m.filteredMems[i], m.filteredMems[j]

		var less bool
		switch m.sortBy {
		case SortByImportance:
			less = a.Importance < b.Importance

		case SortByUpdated:
			less = a.UpdatedAt.AsTime().Before(b.UpdatedAt.AsTime())

		case SortByCreated:
			less = a.CreatedAt.AsTime().Before(b.CreatedAt.AsTime())

		case SortByRelevance:
			// TODO: Implement relevance scoring based on search query
			less = a.Importance < b.Importance

		default:
			less = false
		}

		if m.sortDesc {
			return !less
		}
		return less
	})
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

// SetSearchQuery sets the search query and re-filters.
func (m *Model) SetSearchQuery(query string) {
	m.searchQuery = query
	m.applyFilters()
}

// SetFilter sets the filter options and re-filters.
func (m *Model) SetFilter(tags []string, namespace string, minImportance uint32) {
	m.filterTags = tags
	m.filterNS = namespace
	m.minImportance = minImportance
	m.applyFilters()
}

// SetSort sets the sort mode and re-sorts.
func (m *Model) SetSort(sortBy SortMode, desc bool) {
	m.sortBy = sortBy
	m.sortDesc = desc
	m.applySorting()
}

// ClearFilters clears all filters and shows all memories.
func (m *Model) ClearFilters() {
	m.searchQuery = ""
	m.filterTags = nil
	m.filterNS = ""
	m.minImportance = 0
	m.applyFilters()
}
