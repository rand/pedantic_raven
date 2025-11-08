// Package context provides the Context Panel for displaying semantic analysis results.
//
// The Context Panel shows:
// - Entities (with types and occurrence counts)
// - Relationships (subject-predicate-object triples)
// - Typed Holes (with priority and complexity)
// - Dependencies (imports and references)
//
// Features:
// - Real-time updates from semantic analyzer
// - Filtering by type or search query
// - Section toggling (show/hide sections)
// - Scrollable content
package context

import (
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// Section represents a collapsible section in the context panel.
type Section int

const (
	SectionEntities Section = iota
	SectionRelationships
	SectionTypedHoles
	SectionDependencies
	SectionTriples
)

// String returns the string representation of a section.
func (s Section) String() string {
	switch s {
	case SectionEntities:
		return "Entities"
	case SectionRelationships:
		return "Relationships"
	case SectionTypedHoles:
		return "Typed Holes"
	case SectionDependencies:
		return "Dependencies"
	case SectionTriples:
		return "Triples"
	default:
		return "Unknown"
	}
}

// SectionState tracks the visibility of a section.
type SectionState struct {
	Expanded bool
	Filtered bool // Whether section is filtered out
}

// FilterMode determines how filtering is applied.
type FilterMode int

const (
	FilterNone FilterMode = iota // No filtering
	FilterEntity                 // Filter by entity type
	FilterSearch                 // Filter by search query
)

// ContextPanelConfig configures the context panel.
type ContextPanelConfig struct {
	Width          int                    // Panel width
	Height         int                    // Panel height
	ShowLineNumbers bool                   // Show line numbers in results
	DefaultExpanded map[Section]bool       // Which sections are expanded by default
	MaxResults      map[Section]int        // Maximum results per section
}

// DefaultContextPanelConfig returns default configuration.
func DefaultContextPanelConfig() ContextPanelConfig {
	return ContextPanelConfig{
		Width:           40,
		Height:          30,
		ShowLineNumbers: true,
		DefaultExpanded: map[Section]bool{
			SectionEntities:      true,
			SectionRelationships: true,
			SectionTypedHoles:    true,
			SectionDependencies:  false,
			SectionTriples:       false,
		},
		MaxResults: map[Section]int{
			SectionEntities:      50,
			SectionRelationships: 30,
			SectionTypedHoles:    20,
			SectionDependencies:  30,
			SectionTriples:       50,
		},
	}
}

// ContextPanel displays semantic analysis results.
type ContextPanel struct {
	config ContextPanelConfig

	// Analysis results
	analysis *semantic.Analysis

	// Section state
	sections      map[Section]SectionState
	activeSection Section

	// Filtering
	filterMode  FilterMode
	filterQuery string
	filterType  semantic.EntityType

	// Scrolling
	scrollOffset int
	totalLines   int

	// Enhanced holes (for priority/complexity display)
	enhancedHoles []semantic.EnhancedTypedHole
}

// New creates a new context panel.
func New(config ContextPanelConfig) *ContextPanel {
	sections := make(map[Section]SectionState)
	for section, expanded := range config.DefaultExpanded {
		sections[section] = SectionState{
			Expanded: expanded,
			Filtered: false,
		}
	}

	return &ContextPanel{
		config:        config,
		sections:      sections,
		activeSection: SectionEntities,
		filterMode:    FilterNone,
		scrollOffset:  0,
	}
}

// SetAnalysis updates the analysis results.
func (p *ContextPanel) SetAnalysis(analysis *semantic.Analysis) {
	p.analysis = analysis

	// Enhance typed holes for display
	if analysis != nil && len(analysis.TypedHoles) > 0 {
		p.enhancedHoles = make([]semantic.EnhancedTypedHole, len(analysis.TypedHoles))
		for i, hole := range analysis.TypedHoles {
			p.enhancedHoles[i] = semantic.EnhanceTypedHole(hole, analysis.Relationships)
		}
	}
}

// ToggleSection toggles the expansion state of a section.
func (p *ContextPanel) ToggleSection(section Section) {
	if state, ok := p.sections[section]; ok {
		state.Expanded = !state.Expanded
		p.sections[section] = state
	}
}

// ExpandSection expands a section.
func (p *ContextPanel) ExpandSection(section Section) {
	if state, ok := p.sections[section]; ok {
		state.Expanded = true
		p.sections[section] = state
	}
}

// CollapseSection collapses a section.
func (p *ContextPanel) CollapseSection(section Section) {
	if state, ok := p.sections[section]; ok {
		state.Expanded = false
		p.sections[section] = state
	}
}

// IsSectionExpanded returns whether a section is expanded.
func (p *ContextPanel) IsSectionExpanded(section Section) bool {
	if state, ok := p.sections[section]; ok {
		return state.Expanded
	}
	return false
}

// SetFilterMode sets the filter mode.
func (p *ContextPanel) SetFilterMode(mode FilterMode) {
	p.filterMode = mode
	p.updateFilters()
}

// SetFilterQuery sets the search query filter.
func (p *ContextPanel) SetFilterQuery(query string) {
	p.filterQuery = query
	p.filterMode = FilterSearch
	p.updateFilters()
}

// SetFilterType sets the entity type filter.
func (p *ContextPanel) SetFilterType(entityType semantic.EntityType) {
	p.filterType = entityType
	p.filterMode = FilterEntity
	p.updateFilters()
}

// ClearFilter clears all filters.
func (p *ContextPanel) ClearFilter() {
	p.filterMode = FilterNone
	p.filterQuery = ""
	p.updateFilters()
}

// updateFilters applies current filters to sections.
func (p *ContextPanel) updateFilters() {
	for section := range p.sections {
		state := p.sections[section]
		state.Filtered = p.isSectionFiltered(section)
		p.sections[section] = state
	}
}

// isSectionFiltered determines if a section should be filtered out.
func (p *ContextPanel) isSectionFiltered(section Section) bool {
	if p.filterMode == FilterNone {
		return false
	}

	if p.analysis == nil {
		return false
	}

	// For entity type filtering, only show Entities section
	if p.filterMode == FilterEntity && section != SectionEntities {
		return true
	}

	return false
}

// ScrollDown scrolls the view down.
func (p *ContextPanel) ScrollDown(lines int) {
	p.scrollOffset += lines
	// Only clamp if totalLines is set (> 0)
	if p.totalLines > 0 && p.scrollOffset > p.totalLines {
		p.scrollOffset = p.totalLines
	}
}

// ScrollUp scrolls the view up.
func (p *ContextPanel) ScrollUp(lines int) {
	p.scrollOffset -= lines
	if p.scrollOffset < 0 {
		p.scrollOffset = 0
	}
}

// ScrollToTop scrolls to the top.
func (p *ContextPanel) ScrollToTop() {
	p.scrollOffset = 0
}

// ScrollToBottom scrolls to the bottom.
func (p *ContextPanel) ScrollToBottom() {
	p.scrollOffset = p.totalLines
}

// GetScrollOffset returns the current scroll offset.
func (p *ContextPanel) GetScrollOffset() int {
	return p.scrollOffset
}

// NextSection moves to the next section.
func (p *ContextPanel) NextSection() {
	sections := []Section{
		SectionEntities,
		SectionRelationships,
		SectionTypedHoles,
		SectionDependencies,
		SectionTriples,
	}

	for i, section := range sections {
		if section == p.activeSection {
			p.activeSection = sections[(i+1)%len(sections)]
			break
		}
	}
}

// PreviousSection moves to the previous section.
func (p *ContextPanel) PreviousSection() {
	sections := []Section{
		SectionEntities,
		SectionRelationships,
		SectionTypedHoles,
		SectionDependencies,
		SectionTriples,
	}

	for i, section := range sections {
		if section == p.activeSection {
			if i == 0 {
				p.activeSection = sections[len(sections)-1]
			} else {
				p.activeSection = sections[i-1]
			}
			break
		}
	}
}

// GetActiveSection returns the currently active section.
func (p *ContextPanel) GetActiveSection() Section {
	return p.activeSection
}

// SetActiveSection sets the active section.
func (p *ContextPanel) SetActiveSection(section Section) {
	p.activeSection = section
}

// GetAnalysis returns the current analysis.
func (p *ContextPanel) GetAnalysis() *semantic.Analysis {
	return p.analysis
}

// GetConfig returns the panel configuration.
func (p *ContextPanel) GetConfig() ContextPanelConfig {
	return p.config
}

// SetConfig updates the panel configuration.
func (p *ContextPanel) SetConfig(config ContextPanelConfig) {
	p.config = config
}
