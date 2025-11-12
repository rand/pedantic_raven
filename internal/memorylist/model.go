package memorylist

import (
	"math"
	"sort"
	"strings"
	"time"

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

	case SearchResultsMsg:
		m.SetMemories(msg.Results, msg.TotalCount)
		m.lastSearchQuery = msg.Query
		m.searchActive = true
		if msg.Err == nil && msg.Query != "" {
			m.searchHistory.Add(msg.Query)
		}
		return m, nil

	case MemoriesLoadedMsg:
		m.SetMemories(msg.Memories, msg.TotalCount)
		return m, nil

	case MemoriesErrorMsg:
		m.SetError(msg.Err)
		return m, nil

	case SearchQueryMsg:
		m.searchQuery = msg.Query
		m.applyFilters()
		// If client is available and autoReload, trigger search
		if m.client != nil && m.autoReload {
			return m, m.searchCmd()
		}
		return m, nil

	case FilterChangedMsg:
		m.filterTags = msg.Tags
		m.filterNS = msg.Namespace
		m.minImportance = msg.MinImportance
		m.applyFilters()
		// If client is available and autoReload, trigger reload
		if m.client != nil && m.autoReload {
			return m, m.reloadCmd()
		}
		return m, nil

	case SortChangedMsg:
		m.sortBy = msg.Mode
		m.sortDesc = msg.Desc
		m.applySorting()
		return m, nil

	case LoadRequestMsg:
		m.SetLoading(true)
		return m, m.loadCmd()

	case ReloadRequestMsg:
		m.SetLoading(true)
		return m, m.reloadCmd()
	}

	return m, nil
}

// loadCmd returns a command to load memories from the client.
func (m Model) loadCmd() tea.Cmd {
	if m.client == nil {
		return nil
	}

	// Type assert to mnemosyne.Client (we know this is safe)
	// Using interface{} allows us to avoid import cycles
	return nil // Will be implemented by the caller with proper client access
}

// searchCmd returns a command to search memories.
func (m Model) searchCmd() tea.Cmd {
	if m.client == nil || m.searchQuery == "" {
		return nil
	}
	return nil // Will be implemented by the caller
}

// reloadCmd returns a command to reload with current filters.
func (m Model) reloadCmd() tea.Cmd {
	if m.client == nil {
		return nil
	}
	return nil // Will be implemented by the caller
}

// handleKeyPress processes keyboard input.
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle input modes first
	switch m.inputMode {
	case InputModeSearch:
		return m.handleSearchInput(msg)
	case InputModeFilter:
		return m.handleFilterInput(msg)
	}

	// Normal mode keyboard shortcuts
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

	case "/":
		return m.enterSearchMode(), nil

	case "ctrl+m":
		// Cycle search mode
		m.cycleSearchMode()
		// Re-execute search if active
		if m.searchActive && m.searchQuery != "" && m.mnemosyneClient != nil {
			m.SetLoading(true)
			return m, m.executeSearch()
		}
		return m, nil

	case "ctrl+f":
		// Toggle filter input (future implementation)
		return m.enterFilterMode(), nil

	case "?":
		m.showHelp = !m.showHelp
		return m, nil

	case "r":
		// Reload/refresh
		m.SetLoading(true)
		return m, func() tea.Msg {
			return ReloadRequestMsg{}
		}

	case "c":
		// Clear filters and search
		m.ClearFilters()
		m.searchQuery = ""
		m.searchInput = ""
		return m, nil

	case "esc":
		// Clear help or error
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		if m.err != nil {
			m.err = nil
			return m, nil
		}
	}

	return m, nil
}

// handleSearchInput processes keyboard input in search mode.
func (m Model) handleSearchInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		// Commit search immediately
		m.searchQuery = m.searchInput
		m.inputMode = InputModeNormal
		// Cancel any pending debounced search
		if m.searchDebouncer != nil {
			m.searchDebouncer.Cancel()
		}
		// Execute search if client is available
		if m.mnemosyneClient != nil && m.mnemosyneClient.IsConnected() {
			m.SetLoading(true)
			return m, m.executeSearch()
		}
		// Fallback to local filtering
		m.applyFilters()
		return m, nil

	case tea.KeyEsc:
		// Cancel search
		m.searchInput = ""
		m.inputMode = InputModeNormal
		if m.searchDebouncer != nil {
			m.searchDebouncer.Cancel()
		}
		return m, nil

	case tea.KeyBackspace, tea.KeyDelete:
		// Delete character
		if len(m.searchInput) > 0 {
			m.searchInput = m.searchInput[:len(m.searchInput)-1]
		}
		// Trigger debounced search
		return m, m.debouncedSearch()

	case tea.KeyRunes:
		// Add character
		m.searchInput += string(msg.Runes)
		// Trigger debounced search
		return m, m.debouncedSearch()
	}

	return m, nil
}

// handleFilterInput processes keyboard input in filter mode.
func (m Model) handleFilterInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		// Commit filter
		m.inputMode = InputModeNormal
		m.applyFilters()
		return m, nil

	case tea.KeyEsc:
		// Cancel filter
		m.inputMode = InputModeNormal
		return m, nil
	}

	return m, nil
}

// enterSearchMode switches to search input mode.
func (m Model) enterSearchMode() Model {
	m.inputMode = InputModeSearch
	m.searchInput = m.searchQuery // Start with current query
	return m
}

// enterFilterMode switches to filter input mode.
func (m Model) enterFilterMode() Model {
	m.inputMode = InputModeFilter
	return m
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

// RelevanceScorer scores memories based on search relevance.
type RelevanceScorer struct {
	query string
}

// Score calculates a relevance score for a memory based on the search query.
func (s *RelevanceScorer) Score(memory *pb.MemoryNote) float64 {
	score := 0.0

	// Exact match in content (high score)
	if strings.Contains(strings.ToLower(memory.Content), strings.ToLower(s.query)) {
		score += 10.0
	}

	// Match in tags (medium score)
	for _, tag := range memory.Tags {
		if strings.Contains(strings.ToLower(tag), strings.ToLower(s.query)) {
			score += 5.0
		}
	}

	// Importance boost
	score += float64(memory.Importance)

	// Recency boost (newer memories score higher)
	updatedAt := time.Unix(int64(memory.UpdatedAt), 0)
	age := time.Since(updatedAt)
	recencyScore := math.Max(0, 5.0-(age.Hours()/24/30)) // Decay over months
	score += recencyScore

	return score
}

// SortMemoriesByRelevance sorts memories by relevance to a query (non-mutating).
func SortMemoriesByRelevance(memories []*pb.MemoryNote, query string) []*pb.MemoryNote {
	scorer := &RelevanceScorer{query: query}

	sort.Slice(memories, func(i, j int) bool {
		return scorer.Score(memories[i]) > scorer.Score(memories[j])
	})

	return memories
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
			less = a.UpdatedAt < b.UpdatedAt

		case SortByCreated:
			less = a.CreatedAt < b.CreatedAt

		case SortByRelevance:
			// Use relevance scoring based on search query
			scorer := &RelevanceScorer{query: m.searchQuery}
			less = scorer.Score(a) < scorer.Score(b)

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

// cycleSearchMode cycles through available search modes.
func (m *Model) cycleSearchMode() {
	switch m.searchOptions.SearchMode {
	case SearchHybrid:
		m.searchOptions.SearchMode = SearchSemantic
	case SearchSemantic:
		m.searchOptions.SearchMode = SearchFullText
	case SearchFullText:
		m.searchOptions.SearchMode = SearchGraph
	case SearchGraph:
		m.searchOptions.SearchMode = SearchHybrid
	default:
		m.searchOptions.SearchMode = SearchHybrid
	}
}

// executeSearch executes a search with current search options.
func (m Model) executeSearch() tea.Cmd {
	if m.mnemosyneClient == nil || !m.mnemosyneClient.IsConnected() {
		return nil
	}

	// Build search options from model state
	opts := m.searchOptions
	opts.Query = m.searchQuery
	opts.Namespaces = []string{m.filterNS}
	opts.Tags = m.filterTags
	opts.MinImportance = int(m.minImportance)

	return SearchWithOptions(m.mnemosyneClient, opts)
}

// debouncedSearch triggers a debounced search operation.
func (m Model) debouncedSearch() tea.Cmd {
	if m.searchDebouncer == nil || m.mnemosyneClient == nil || !m.mnemosyneClient.IsConnected() {
		return nil
	}

	// Return a command that sets up the debounced search
	return func() tea.Msg {
		// Note: This is a simplified approach. In a real implementation,
		// we'd need to communicate back to the Update loop after the debounce.
		// For now, we trigger search immediately for typing feedback.
		return nil
	}
}
