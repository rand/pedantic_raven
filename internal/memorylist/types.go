package memorylist

import (
	tea "github.com/charmbracelet/bubbletea"
	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// Model represents the memory list component state.
type Model struct {
	// Memory data
	memories      []*pb.MemoryNote
	filteredMems  []*pb.MemoryNote
	selectedIndex int

	// Viewport state
	scrollOffset int
	height       int
	width        int

	// Filter/search state
	searchQuery   string
	filterTags    []string
	filterNS      string
	minImportance uint32

	// Sort state
	sortBy   SortMode
	sortDesc bool

	// UI state
	loading bool
	err     error
	focused bool

	// Pagination
	pageSize   int
	totalCount uint32
}

// SortMode defines how memories are sorted.
type SortMode int

const (
	SortByImportance SortMode = iota
	SortByUpdated
	SortByCreated
	SortByRelevance
)

// String returns the string representation of the sort mode.
func (s SortMode) String() string {
	switch s {
	case SortByImportance:
		return "Importance"
	case SortByUpdated:
		return "Updated"
	case SortByCreated:
		return "Created"
	case SortByRelevance:
		return "Relevance"
	default:
		return "Unknown"
	}
}

// Messages for the memory list component.
type (
	// MemoriesLoadedMsg is sent when memories are loaded from mnemosyne.
	MemoriesLoadedMsg struct {
		Memories   []*pb.MemoryNote
		TotalCount uint32
	}

	// MemoriesErrorMsg is sent when memory loading fails.
	MemoriesErrorMsg struct {
		Err error
	}

	// MemorySelectedMsg is sent when a memory is selected.
	MemorySelectedMsg struct {
		Memory *pb.MemoryNote
		Index  int
	}

	// SearchQueryMsg is sent when search query changes.
	SearchQueryMsg struct {
		Query string
	}

	// FilterChangedMsg is sent when filters change.
	FilterChangedMsg struct {
		Tags          []string
		Namespace     string
		MinImportance uint32
	}

	// SortChangedMsg is sent when sort mode changes.
	SortChangedMsg struct {
		Mode SortMode
		Desc bool
	}
)

// NewModel creates a new memory list model with default settings.
func NewModel() Model {
	return Model{
		memories:      make([]*pb.MemoryNote, 0),
		filteredMems:  make([]*pb.MemoryNote, 0),
		selectedIndex: 0,
		scrollOffset:  0,
		height:        20,
		width:         80,
		sortBy:        SortByUpdated,
		sortDesc:      true, // Most recent first
		loading:       false,
		focused:       true,
		pageSize:      50,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetSize sets the component size.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocus sets the focus state.
func (m *Model) SetFocus(focused bool) {
	m.focused = focused
}

// IsFocused returns whether the component is focused.
func (m Model) IsFocused() bool {
	return m.focused
}

// SetMemories sets the memory list.
func (m *Model) SetMemories(memories []*pb.MemoryNote, totalCount uint32) {
	m.memories = memories
	m.totalCount = totalCount
	m.loading = false
	m.err = nil
	m.applyFilters()
}

// SetLoading sets the loading state.
func (m *Model) SetLoading(loading bool) {
	m.loading = loading
}

// SetError sets the error state.
func (m *Model) SetError(err error) {
	m.err = err
	m.loading = false
}

// SelectedMemory returns the currently selected memory, or nil if none.
func (m Model) SelectedMemory() *pb.MemoryNote {
	if len(m.filteredMems) == 0 || m.selectedIndex < 0 || m.selectedIndex >= len(m.filteredMems) {
		return nil
	}
	return m.filteredMems[m.selectedIndex]
}

// SelectedIndex returns the current selection index.
func (m Model) SelectedIndex() int {
	return m.selectedIndex
}

// MemoryCount returns the total number of filtered memories.
func (m Model) MemoryCount() int {
	return len(m.filteredMems)
}

// TotalCount returns the total count from the server (before pagination).
func (m Model) TotalCount() uint32 {
	return m.totalCount
}

// IsLoading returns whether the component is loading.
func (m Model) IsLoading() bool {
	return m.loading
}

// Error returns the current error, if any.
func (m Model) Error() error {
	return m.err
}
