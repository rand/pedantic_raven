package context

import (
	"fmt"
	"strings"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// Render renders the context panel to a string.
func (p *ContextPanel) Render() string {
	var lines []string

	// Header
	lines = append(lines, p.renderHeader())
	lines = append(lines, "")

	// Filter info
	if p.filterMode != FilterNone {
		lines = append(lines, p.renderFilterInfo())
		lines = append(lines, "")
	}

	// Sections
	if p.analysis != nil {
		lines = append(lines, p.renderSections()...)
	} else {
		lines = append(lines, "No analysis available")
	}

	// Apply scrolling
	p.totalLines = len(lines)
	visibleLines := p.applyScrolling(lines)

	// Join and ensure width
	return p.formatOutput(visibleLines)
}

// renderHeader renders the panel header.
func (p *ContextPanel) renderHeader() string {
	title := "Context Panel"
	if p.analysis != nil {
		stats := p.analysis.GetStatistics()
		title = fmt.Sprintf("Context (%d entities, %d rels)",
			stats.UniqueEntities, stats.TotalRelationships)
	}
	return p.centerText(title)
}

// renderFilterInfo renders current filter information.
func (p *ContextPanel) renderFilterInfo() string {
	switch p.filterMode {
	case FilterEntity:
		return fmt.Sprintf("Filter: %s entities", p.filterType.String())
	case FilterSearch:
		return fmt.Sprintf("Search: %s", p.filterQuery)
	default:
		return ""
	}
}

// renderSections renders all sections.
func (p *ContextPanel) renderSections() []string {
	var lines []string

	sections := []Section{
		SectionEntities,
		SectionRelationships,
		SectionTypedHoles,
		SectionDependencies,
		SectionTriples,
	}

	for _, section := range sections {
		state := p.sections[section]

		// Skip filtered sections
		if state.Filtered {
			continue
		}

		sectionLines := p.renderSection(section, state)
		lines = append(lines, sectionLines...)
		lines = append(lines, "") // Blank line between sections
	}

	return lines
}

// renderSection renders a single section.
func (p *ContextPanel) renderSection(section Section, state SectionState) []string {
	var lines []string

	// Section header
	header := p.renderSectionHeader(section, state)
	lines = append(lines, header)

	// Section content (if expanded)
	if state.Expanded {
		content := p.renderSectionContent(section)
		lines = append(lines, content...)
	}

	return lines
}

// renderSectionHeader renders a section header.
func (p *ContextPanel) renderSectionHeader(section Section, state SectionState) string {
	indicator := "▶"
	if state.Expanded {
		indicator = "▼"
	}

	active := ""
	if section == p.activeSection {
		active = " *"
	}

	count := p.getSectionCount(section)
	return fmt.Sprintf("%s %s (%d)%s", indicator, section.String(), count, active)
}

// getSectionCount returns the number of items in a section.
func (p *ContextPanel) getSectionCount(section Section) int {
	if p.analysis == nil {
		return 0
	}

	switch section {
	case SectionEntities:
		return len(p.filterEntities(p.analysis.Entities))
	case SectionRelationships:
		return len(p.analysis.Relationships)
	case SectionTypedHoles:
		return len(p.analysis.TypedHoles)
	case SectionDependencies:
		return len(p.analysis.Dependencies)
	case SectionTriples:
		return len(p.analysis.Triples)
	default:
		return 0
	}
}

// renderSectionContent renders the content of a section.
func (p *ContextPanel) renderSectionContent(section Section) []string {
	switch section {
	case SectionEntities:
		return p.renderEntities()
	case SectionRelationships:
		return p.renderRelationships()
	case SectionTypedHoles:
		return p.renderTypedHoles()
	case SectionDependencies:
		return p.renderDependencies()
	case SectionTriples:
		return p.renderTriples()
	default:
		return []string{}
	}
}

// renderEntities renders the entities section.
func (p *ContextPanel) renderEntities() []string {
	var lines []string

	entities := p.filterEntities(p.analysis.Entities)
	maxResults := p.config.MaxResults[SectionEntities]

	for i, entity := range entities {
		if i >= maxResults {
			lines = append(lines, fmt.Sprintf("  ... and %d more", len(entities)-i))
			break
		}

		line := fmt.Sprintf("  %s [%s]", entity.Text, entity.Type.String())
		if entity.Count > 1 {
			line += fmt.Sprintf(" ×%d", entity.Count)
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		lines = append(lines, "  (no entities)")
	}

	return lines
}

// renderRelationships renders the relationships section.
func (p *ContextPanel) renderRelationships() []string {
	var lines []string

	maxResults := p.config.MaxResults[SectionRelationships]
	relationships := p.analysis.Relationships

	for i, rel := range relationships {
		if i >= maxResults {
			lines = append(lines, fmt.Sprintf("  ... and %d more", len(relationships)-i))
			break
		}

		line := fmt.Sprintf("  %s → %s → %s", rel.Subject, rel.Predicate, rel.Object)
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		lines = append(lines, "  (no relationships)")
	}

	return lines
}

// renderTypedHoles renders the typed holes section.
func (p *ContextPanel) renderTypedHoles() []string {
	var lines []string

	maxResults := p.config.MaxResults[SectionTypedHoles]

	for i, hole := range p.enhancedHoles {
		if i >= maxResults {
			lines = append(lines, fmt.Sprintf("  ... and %d more", len(p.enhancedHoles)-i))
			break
		}

		// Format: ??Type [P:8 C:5]
		line := fmt.Sprintf("  ??%s [P:%d C:%d]",
			hole.Type, hole.Priority, hole.Complexity)

		if hole.Constraint != "" {
			line += fmt.Sprintf(" !%s", hole.Constraint)
		}

		lines = append(lines, line)
	}

	if len(lines) == 0 {
		lines = append(lines, "  (no typed holes)")
	}

	return lines
}

// renderDependencies renders the dependencies section.
func (p *ContextPanel) renderDependencies() []string {
	var lines []string

	maxResults := p.config.MaxResults[SectionDependencies]
	dependencies := p.analysis.Dependencies

	for i, dep := range dependencies {
		if i >= maxResults {
			lines = append(lines, fmt.Sprintf("  ... and %d more", len(dependencies)-i))
			break
		}

		line := fmt.Sprintf("  %s %s", dep.Type, dep.Target)
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		lines = append(lines, "  (no dependencies)")
	}

	return lines
}

// renderTriples renders the triples section.
func (p *ContextPanel) renderTriples() []string {
	var lines []string

	maxResults := p.config.MaxResults[SectionTriples]
	triples := p.analysis.Triples

	for i, triple := range triples {
		if i >= maxResults {
			lines = append(lines, fmt.Sprintf("  ... and %d more", len(triples)-i))
			break
		}

		line := fmt.Sprintf("  (%s, %s, %s)",
			triple.Subject, triple.Predicate, triple.Object)
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		lines = append(lines, "  (no triples)")
	}

	return lines
}

// filterEntities filters entities based on current filter mode.
func (p *ContextPanel) filterEntities(entities []semantic.Entity) []semantic.Entity {
	if p.filterMode == FilterNone {
		return entities
	}

	var filtered []semantic.Entity

	for _, entity := range entities {
		include := false

		switch p.filterMode {
		case FilterEntity:
			include = entity.Type == p.filterType
		case FilterSearch:
			include = strings.Contains(
				strings.ToLower(entity.Text),
				strings.ToLower(p.filterQuery),
			)
		}

		if include {
			filtered = append(filtered, entity)
		}
	}

	return filtered
}

// applyScrolling applies scroll offset to visible lines.
func (p *ContextPanel) applyScrolling(lines []string) []string {
	if p.scrollOffset >= len(lines) {
		return []string{}
	}

	end := p.scrollOffset + p.config.Height
	if end > len(lines) {
		end = len(lines)
	}

	return lines[p.scrollOffset:end]
}

// formatOutput formats output lines to fit panel width.
func (p *ContextPanel) formatOutput(lines []string) string {
	var formatted []string

	for _, line := range lines {
		// Truncate if too long
		if len(line) > p.config.Width {
			line = line[:p.config.Width-3] + "..."
		}

		// Pad if too short
		if len(line) < p.config.Width {
			line = line + strings.Repeat(" ", p.config.Width-len(line))
		}

		formatted = append(formatted, line)
	}

	return strings.Join(formatted, "\n")
}

// centerText centers text within the panel width.
func (p *ContextPanel) centerText(text string) string {
	if len(text) >= p.config.Width {
		return text[:p.config.Width]
	}

	padding := (p.config.Width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}
