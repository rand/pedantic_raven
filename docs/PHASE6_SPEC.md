# Phase 6: Analyze Mode - Technical Specification

**Version**: 1.0
**Status**: Active
**Timeline**: 2-3 weeks
**Start Date**: 2025-11-09
**Epic**: `pedantic_raven-nd9`

---

## Executive Summary

Phase 6 implements **Analyze Mode**, a statistical analysis and visualization layer for semantic data extracted from context documents. This mode provides insights into entity patterns, relationship networks, and typed hole complexity, enabling users to understand and navigate their semantic knowledge graph.

**Key Innovation**: Transform raw semantic analysis into actionable intelligence through interactive visualizations and data mining.

---

## Objectives

### Primary Goals
1. **Entity Analysis**: Visualize entity frequency, types, and co-occurrence patterns
2. **Relationship Mining**: Discover and analyze relationship patterns across documents
3. **Typed Hole Intelligence**: Prioritize implementation work based on complexity and dependencies
4. **Interactive Exploration**: Provide zoom, pan, filter capabilities for all visualizations
5. **Export Capabilities**: Generate reports in multiple formats (PDF, HTML, Markdown)

### Success Criteria
- All visualizations render in <100ms for datasets <1000 entities
- Interactive filters update in <50ms
- Export generates complete reports in <2s
- All keyboard shortcuts work consistently
- Zero regressions in existing modes

---

## Architecture

### Component Structure

```
internal/modes/analyze.go           # Analyze Mode orchestrator
internal/analyze/
  ├── model.go                      # Main analyze model
  ├── triple_graph.go               # Triple visualization component
  ├── entity_analysis.go            # Entity frequency analysis
  ├── relationship_mining.go        # Relationship pattern detection
  ├── hole_prioritization.go        # Typed hole scoring
  ├── export.go                     # Report generation
  ├── filters.go                    # Interactive filtering
  └── stats.go                      # Statistical calculations

internal/analyze/visualizations/
  ├── bar_chart.go                  # Terminal-based bar charts
  ├── word_cloud.go                 # ASCII word cloud
  ├── dependency_tree.go            # Tree visualization
  └── graph_layout.go               # Graph layout algorithms

internal/analyze/export/
  ├── pdf.go                        # PDF generation
  ├── html.go                       # HTML report generation
  └── markdown.go                   # Markdown report generation
```

### Data Flow

```
[Edit Mode] → Semantic Analyzer → Analysis Data
     ↓
[Analyze Mode]
     ├→ Triple Graph Component → [Visualization]
     ├→ Entity Analysis → [Bar Charts, Word Clouds]
     ├→ Relationship Mining → [Pattern Tables]
     ├→ Hole Prioritization → [Dependency Trees]
     └→ Export → [PDF/HTML/Markdown Reports]
```

### Integration Points

- **Semantic Analyzer** (`internal/editor/semantic/`): Source of entity/relationship data
- **Graph Visualization** (`internal/memorygraph/`): Reuse layout algorithms
- **Mode Registry** (`internal/modes/`): Register Analyze Mode
- **Event System** (`internal/app/events/`): Mode switching, data updates

---

## Component 6.1: Triple Graph Visualization

**Task**: `pedantic_raven-nd9.1`
**Priority**: 1
**Estimated Effort**: 5 days

### Requirements

**Functional**:
- Display entity-relationship graph with subject-predicate-object triples
- Interactive navigation (zoom, pan, node selection)
- Filter by entity type, relationship type, confidence score
- Highlight connected components
- Show/hide entity metadata (type, frequency, importance)

**Visual Design**:
```
┌─────────────────────────────────────────────────────────┐
│ Triple Graph View                         [Filter] [?] │
├─────────────────────────────────────────────────────────┤
│                                                         │
│      [Person: John]                                     │
│            │                                            │
│            │ works_at (0.9)                             │
│            ↓                                            │
│      [Org: Acme Corp] ──────creates──────→ [Product]   │
│            │                                    │       │
│            │ located_in (0.85)                  │       │
│            ↓                                    │       │
│      [Place: NYC]                              │       │
│                                                ↓       │
│                                          [Tech: API]    │
│                                                         │
│ Entities: 5  Relationships: 4  Avg Confidence: 0.88    │
│                                                         │
│ [h/j/k/l] Pan  [+/-] Zoom  [Enter] Select  [f] Filter  │
└─────────────────────────────────────────────────────────┘
```

**Technical Details**:
- Reuse force-directed layout from `internal/memorygraph/layout.go`
- Add relationship edges with labels and confidence scores
- Implement node clustering by entity type
- Use color coding: Person (blue), Org (green), Place (yellow), Tech (red)

**Performance**:
- Handle 1000+ entities without lag
- Sub-100ms re-layout after filter changes
- Spatial indexing for click detection

### Test Coverage

- Graph rendering with various entity counts (10, 100, 1000 entities)
- Filter updates (by type, confidence, keyword)
- Pan/zoom interactions
- Node selection and metadata display
- Layout convergence and stability
- Performance benchmarks

---

## Component 6.2: Entity Frequency Analysis

**Task**: `pedantic_raven-nd9.2`
**Priority**: 1
**Estimated Effort**: 3 days

### Requirements

**Functional**:
- Bar chart showing entity frequency by type
- Word cloud visualization for most frequent entities
- Filter by entity type (Person, Org, Place, Thing, Concept, Tech)
- Sort by frequency, alphabetical, importance
- Show entity details on selection

**Visual Design**:

**Bar Chart**:
```
┌────────────────────────────────────────────────┐
│ Entity Frequency by Type                      │
├────────────────────────────────────────────────┤
│                                                │
│ Person      ████████████████ 45              │
│ Org         ████████████ 32                   │
│ Technology  ████████ 23                       │
│ Place       █████ 15                          │
│ Concept     ████ 12                           │
│ Thing       ██ 8                              │
│                                                │
│ Total Entities: 135                            │
│ [s] Sort  [f] Filter  [Enter] Details         │
└────────────────────────────────────────────────┘
```

**Word Cloud** (ASCII-based):
```
┌────────────────────────────────────────────────┐
│ Most Frequent Entities                         │
├────────────────────────────────────────────────┤
│                                                │
│           API                 DATABASE         │
│     John          Server                       │
│                         React                  │
│   Acme Corp                        GraphQL     │
│              New York                          │
│        Python                                  │
│                  Authentication                │
│                                                │
│ Font size = log(frequency)                     │
└────────────────────────────────────────────────┘
```

**Technical Details**:
- Calculate entity frequency from semantic analysis results
- Implement logarithmic font scaling for word cloud
- Support GLiNER custom entity types if enabled
- Cache calculations for performance

**Statistical Metrics**:
- Total entity count
- Unique entities
- Average entities per document
- Entity type distribution (%)
- Top 10 most frequent entities

### Test Coverage

- Frequency calculation with various datasets
- Bar chart rendering at different terminal widths
- Word cloud text placement algorithm
- Filtering and sorting operations
- Empty state handling

---

## Component 6.3: Relationship Mining

**Task**: `pedantic_raven-nd9.3`
**Priority**: 1
**Estimated Effort**: 4 days

### Requirements

**Functional**:
- Detect common relationship patterns (e.g., "X works_at Y", "A creates B")
- Score relationship strength based on frequency and confidence
- Identify relationship clusters
- Suggest new potential relationships
- Interactive pattern exploration

**Visual Design**:
```
┌─────────────────────────────────────────────────────────┐
│ Relationship Patterns                                   │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ Pattern: [Person] works_at [Organization]              │
│   Occurrences: 12  Avg Confidence: 0.87                │
│   Examples:                                             │
│     • John works_at Acme Corp                           │
│     • Alice works_at Tech Inc                           │
│     • Bob works_at StartupXYZ                           │
│                                                         │
│ Pattern: [Organization] creates [Product]               │
│   Occurrences: 8  Avg Confidence: 0.92                 │
│   Examples:                                             │
│     • Acme Corp creates API Gateway                     │
│     • Tech Inc creates Mobile App                       │
│                                                         │
│ Pattern: [Technology] used_in [Product]                 │
│   Occurrences: 15  Avg Confidence: 0.81                │
│                                                         │
│ [↑/↓] Navigate  [Enter] Expand  [f] Filter             │
└─────────────────────────────────────────────────────────┘
```

**Algorithms**:
1. **Pattern Detection**: Group relationships by (subject type, predicate, object type)
2. **Strength Scoring**: `score = (frequency × avg_confidence) / total_relationships`
3. **Clustering**: Use cosine similarity on relationship embedding
4. **Suggestion**: Infer missing relationships from pattern templates

**Technical Details**:
- Pattern matching using relationship triples
- Confidence aggregation (mean, median, weighted)
- Support for custom GLiNER entity types
- Export pattern library for reuse

### Test Coverage

- Pattern detection with various relationship sets
- Strength scoring algorithm verification
- Clustering with different similarity thresholds
- Suggestion generation accuracy
- Performance with 1000+ relationships

---

## Component 6.4: Typed Hole Prioritization

**Task**: `pedantic_raven-nd9.4`
**Priority**: 1
**Estimated Effort**: 3 days

### Requirements

**Functional**:
- Calculate complexity score for each typed hole
- Build dependency tree showing hole relationships
- Prioritize holes by: complexity, dependencies, criticality
- Visualize dependency graph
- Suggest implementation order

**Visual Design**:

**Priority List**:
```
┌─────────────────────────────────────────────────────────┐
│ Typed Hole Priority Queue                              │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ 1. ??AuthService (Priority: 10, Complexity: 8)         │
│    • Required by: UserHandler, AdminHandler            │
│    • Constraints: thread-safe, async                   │
│    • Estimated effort: 8-12 hours                      │
│                                                         │
│ 2. ??DatabaseLayer (Priority: 9, Complexity: 7)        │
│    • Required by: AuthService, DataService             │
│    • Constraints: atomic, concurrent                   │
│    • Estimated effort: 6-10 hours                      │
│                                                         │
│ 3. ??ConfigLoader (Priority: 8, Complexity: 4)         │
│    • Required by: AuthService, DatabaseLayer           │
│    • Constraints: immutable                            │
│    • Estimated effort: 2-4 hours                       │
│                                                         │
│ [↑/↓] Navigate  [Enter] Details  [d] Show Deps         │
└─────────────────────────────────────────────────────────┘
```

**Dependency Tree**:
```
┌─────────────────────────────────────────────────────────┐
│ Typed Hole Dependencies                                 │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ ??ConfigLoader (Complexity: 4)                          │
│   ├── ??DatabaseLayer (Complexity: 7)                   │
│   │     ├── ??AuthService (Complexity: 8)               │
│   │     │     └── ??UserHandler (Complexity: 5)         │
│   │     └── ??DataService (Complexity: 6)               │
│   └── ??Logger (Complexity: 3)                          │
│                                                         │
│ Critical Path: ConfigLoader → DatabaseLayer →           │
│                AuthService → UserHandler                │
│ Total Complexity: 24                                    │
│                                                         │
│ [↑/↓] Navigate  [Enter] Expand  [e] Export             │
└─────────────────────────────────────────────────────────┘
```

**Scoring Algorithm**:
```go
priority = base_priority(5)
         + dependency_count * 2
         + constraint_complexity
         + min(mention_count, 3)

complexity = base_complexity(3)
           + constraint_count
           + special_constraint_bonus
           // thread/concurrent: +2, async: +1, pure/idempotent: +1
```

**Technical Details**:
- Reuse existing priority/complexity functions from `internal/editor/semantic/holes.go`
- Build dependency DAG from hole references
- Implement topological sort for implementation order
- Detect circular dependencies

### Test Coverage

- Priority calculation with various hole types
- Complexity scoring with different constraints
- Dependency tree construction
- Topological sort correctness
- Circular dependency detection

---

## Component 6.5: Export Analysis Reports

**Task**: `pedantic_raven-nd9.5`
**Priority**: 2
**Estimated Effort**: 4 days

### Requirements

**Functional**:
- Generate PDF reports with charts and statistics
- Generate HTML reports with interactive visualizations
- Generate Markdown summaries
- Include all analysis sections (entities, relationships, holes)
- Customizable report templates

**Report Sections**:
1. **Executive Summary**: Key metrics and highlights
2. **Entity Analysis**: Frequency charts, top entities
3. **Relationship Patterns**: Common patterns, strength scores
4. **Typed Hole Report**: Priority queue, dependency tree
5. **Recommendations**: Suggested next steps

**Technical Details**:

**PDF Generation** (using `github.com/jung-kurt/gofpdf`):
- A4 page size, professional layout
- Embedded charts as images (PNG)
- Table of contents with hyperlinks
- Syntax highlighting for code samples

**HTML Generation**:
- Responsive design with Bootstrap/Tailwind
- Interactive charts using Chart.js or D3.js
- Searchable tables
- Print-optimized CSS

**Markdown Generation**:
- GitHub-flavored markdown
- Mermaid diagrams for graphs
- Tables for statistics
- Code blocks for typed holes

### Test Coverage

- PDF generation with various data sizes
- HTML template rendering
- Markdown formatting
- Image embedding in PDF
- Report completeness verification

---

## Component 6.6: Analyze Mode UI Integration

**Task**: `pedantic_raven-nd9.6`
**Priority**: 1
**Estimated Effort**: 3 days

### Requirements

**Functional**:
- Implement Analyze Mode with Bubble Tea
- Multi-view layout (Triple Graph, Entity Analysis, Relationships, Holes)
- View switcher (Tab to cycle, number keys 1-4 for direct)
- Global filter panel
- Help overlay with keyboard shortcuts

**Layout Design**:
```
┌─────────────────────────────────────────────────────────┐
│ ANALYZE MODE                    [1][2][3][4]    [?][q] │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ [Active View: Triple Graph / Entity Freq / Relations / │
│  Typed Holes]                                           │
│                                                         │
│                                                         │
│  (View-specific content here)                           │
│                                                         │
│                                                         │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ Filters: [Type: All] [Min Confidence: 0.0]             │
│ Stats: Entities: 135 | Relations: 87 | Holes: 23       │
└─────────────────────────────────────────────────────────┘
```

**Keyboard Shortcuts**:
```
Global:
  q       Quit Analyze Mode (return to Edit Mode)
  ?       Show help overlay
  Tab     Cycle through views
  1-4     Direct view selection
  f       Focus filter panel
  e       Export current view

View Navigation:
  h/j/k/l Pan (graph views)
  +/-     Zoom (graph views)
  ↑/↓     Navigate lists
  Enter   Select item / expand details
  Esc     Clear selection

Filters:
  t       Filter by type
  c       Filter by confidence
  /       Search/filter text
  Esc     Clear filters
```

**Technical Details**:
- Implement as new `Mode` in `internal/modes/analyze.go`
- Reuse lipgloss styling from existing modes
- Event-driven view updates
- Stateful filter persistence across views

### Test Coverage

- Mode initialization and lifecycle
- View switching (Tab, number keys)
- Keyboard shortcut handling
- Filter application across views
- Help overlay display
- State persistence

---

## Testing Strategy

### Unit Tests

**Target**: Each component independently tested

**Coverage Goals**:
- Component logic: 85%+
- Visualization rendering: 70%+
- Export functions: 80%+
- Filter operations: 90%+

**Test Categories**:
- Data processing (entity aggregation, pattern detection)
- Visualization rendering (charts, graphs, trees)
- Export generation (PDF, HTML, Markdown)
- Filter logic (type, confidence, text search)
- Keyboard input handling

### Integration Tests

**Scenarios**:
1. Load semantic data → Display all views → Export report
2. Apply filters → Verify all views update consistently
3. Navigate between views → Verify state preservation
4. Resize terminal → Verify responsive layout
5. Large dataset (1000+ entities) → Verify performance

### E2E Tests

**User Journeys**:
1. **Entity Analysis**: Enter Analyze Mode → View entity frequency → Filter by type → Export PDF
2. **Relationship Discovery**: Switch to Relationships view → Sort by strength → Expand pattern → Export HTML
3. **Hole Prioritization**: Switch to Holes view → View dependency tree → Export Markdown

### Performance Benchmarks

**Targets**:
- View rendering: <100ms (1000 entities)
- Filter application: <50ms
- Export generation: <2s (PDF), <1s (HTML/Markdown)
- Memory usage: <50MB for 1000-entity dataset

---

## Implementation Timeline

### Week 1: Core Components (Days 1-5)

**Day 1-2**: Triple Graph Visualization
- Implement graph data model
- Port/adapt force-directed layout
- Add relationship edge rendering
- Basic pan/zoom controls

**Day 3**: Entity Frequency Analysis
- Frequency calculation
- Bar chart component
- Word cloud algorithm
- Sorting/filtering

**Day 4-5**: Relationship Mining
- Pattern detection algorithm
- Strength scoring
- Clustering implementation
- Pattern display component

### Week 2: Advanced Features (Days 6-10)

**Day 6-7**: Typed Hole Prioritization
- Dependency tree construction
- Priority scoring
- Visualization components
- Implementation suggestions

**Day 8-9**: Export Reports
- PDF generation setup
- HTML template creation
- Markdown formatter
- Report assembly logic

**Day 10**: UI Integration
- Analyze Mode implementation
- View switcher
- Global filters
- Help overlay

### Week 3: Polish & Testing (Days 11-15)

**Day 11-12**: Testing
- Unit tests for all components
- Integration test suite
- E2E test scenarios
- Performance benchmarks

**Day 13**: Bug fixes and refinement
- Address test failures
- Performance optimization
- UI polish

**Day 14**: Documentation
- User guide for Analyze Mode
- API documentation
- Example reports

**Day 15**: Final review and release
- Code review
- Integration with main
- Update ROADMAP
- Create release notes

---

## Dependencies

### Internal
- `internal/editor/semantic/` - Source of analysis data
- `internal/memorygraph/layout.go` - Graph layout algorithms
- `internal/modes/` - Mode registry
- `internal/app/events/` - Event system

### External
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/jung-kurt/gofpdf` - PDF generation
- `github.com/gomarkdown/markdown` - Markdown processing
- Standard library: `html/template`, `text/template`

### Optional
- `github.com/wcharczuk/go-chart` - Chart generation (if raster charts needed)
- `gonum.org/v1/gonum` - Statistical calculations

---

## Risk Mitigation

### Technical Risks

**Risk**: Complex visualization performance with large datasets
**Mitigation**: Implement spatial indexing, viewport culling, lazy rendering

**Risk**: PDF generation library limitations
**Mitigation**: Use well-tested library (gofpdf), fallback to HTML→PDF if needed

**Risk**: Terminal size constraints for visualizations
**Mitigation**: Responsive design, scrollable viewports, zoom controls

### Schedule Risks

**Risk**: Export functionality takes longer than estimated
**Mitigation**: Prioritize Markdown export (simplest), defer PDF to later if needed

**Risk**: Graph layout algorithm needs significant adaptation
**Mitigation**: Reuse existing memorygraph code with minimal changes

---

## Success Metrics

### Functionality
- ✅ All 6 components implemented and tested
- ✅ All keyboard shortcuts working
- ✅ Export generates valid PDF/HTML/Markdown
- ✅ Handles datasets up to 1000 entities

### Performance
- ✅ View rendering <100ms
- ✅ Filter updates <50ms
- ✅ Export <2s
- ✅ Memory usage <50MB

### Quality
- ✅ 100+ new tests written
- ✅ All tests passing
- ✅ Zero regressions in existing modes
- ✅ User documentation complete

### User Experience
- ✅ Intuitive navigation
- ✅ Responsive to terminal resize
- ✅ Clear visual hierarchy
- ✅ Helpful error messages

---

## Future Enhancements (Post-Phase 6)

1. **Real-time Analysis**: Watch file changes, update analysis automatically
2. **Comparison Mode**: Compare analysis between document versions
3. **AI Insights**: Use LLM to generate natural language insights
4. **Collaborative Analysis**: Share analysis with team (Phase 8 integration)
5. **Custom Metrics**: User-defined analysis functions
6. **Export Templates**: Customizable report layouts

---

## Appendices

### A. Data Structures

```go
// AnalysisData aggregates semantic analysis results
type AnalysisData struct {
    Entities      []Entity
    Relationships []Relationship
    TypedHoles    []TypedHole
    Triples       []Triple
    Stats         Statistics
}

// Statistics holds computed metrics
type Statistics struct {
    TotalEntities      int
    UniqueEntities     int
    TotalRelationships int
    AvgConfidence      float64
    EntityTypeDistribution map[string]int
    TopEntities        []EntityFrequency
}

// EntityFrequency tracks entity occurrence
type EntityFrequency struct {
    Text      string
    Type      string
    Count     int
    Importance int
}

// RelationshipPattern represents a discovered pattern
type RelationshipPattern struct {
    SubjectType string
    Predicate   string
    ObjectType  string
    Occurrences int
    AvgConfidence float64
    Examples    []Relationship
    Strength    float64
}

// HolePriority wraps typed hole with calculated scores
type HolePriority struct {
    Hole         TypedHole
    Priority     int
    Complexity   int
    Dependencies []string
    Effort       string // "2-4 hours", "1-2 days", etc.
}
```

### B. Configuration

```toml
[analyze]
# Default view on mode entry
default_view = "triple_graph"  # or "entity_freq", "relationships", "holes"

# Performance tuning
max_entities_in_graph = 1000
layout_iterations = 50
render_timeout_ms = 100

# Export settings
export_dir = "./analysis-reports"
pdf_template = "default"
html_theme = "light"  # or "dark"

# Filters
default_min_confidence = 0.5
default_entity_types = ["all"]
```

---

**Document Version**: 1.0
**Last Updated**: 2025-11-09
**Author**: Development Team
