# Analyze Mode User Guide

## Overview

Analyze Mode is a powerful semantic analysis visualization system in Pedantic Raven that transforms raw semantic data into actionable intelligence. It provides interactive visualizations and insights into entities, relationships, and typed holes extracted from your documents.

### What is Analyze Mode?

Analyze Mode offers four distinct views for exploring your semantic data:

1. **Triple Graph**: Interactive force-directed visualization of entity-relationship networks
2. **Entity Frequency**: Statistical analysis of entity occurrences and type distributions
3. **Relationship Patterns**: Discovery and analysis of common relationship patterns
4. **Typed Holes**: Prioritization and dependency analysis for implementation planning

### Key Features

- **Interactive Navigation**: Pan, zoom, and select elements across all visualizations
- **Real-time Filtering**: Filter by entity type, confidence scores, and text search
- **Multi-format Export**: Generate reports in Markdown, HTML, or PDF
- **Keyboard-driven**: Efficient navigation with comprehensive keyboard shortcuts
- **Performance**: Handles datasets with 1000+ entities smoothly

### When to Use Analyze Mode

Use Analyze Mode when you need to:

- Understand the structure of your semantic knowledge graph
- Identify the most important entities and their relationships
- Discover patterns in how entities relate to each other
- Prioritize implementation work for typed holes
- Generate reports for documentation or sharing insights
- Explore entity co-occurrence and clustering

## Getting Started

### Entering Analyze Mode

From the main application, press `3` or use the mode switcher to enter Analyze Mode. The mode will load semantic analysis data from your current document or context.

### Quick Start Example

```
1. Enter Analyze Mode (press 3)
2. View the Triple Graph (default view)
3. Pan around with h/j/k/l keys
4. Zoom in with + or out with -
5. Select a node with Enter
6. Switch to Entity Frequency view (press 2)
7. Export a report (press e)
```

### Basic Navigation

- **Switch Views**: Press `Tab` to cycle through views, or `1-4` for direct selection
- **Navigate Content**: Use arrow keys or `h/j/k/l` (vim-style)
- **Select Items**: Press `Enter` to select/expand
- **Clear Selection**: Press `Esc`
- **Export**: Press `e` to export current view
- **Help**: Press `?` to show keyboard shortcuts

### Understanding the Views

Each view provides different perspectives on your semantic data:

**Triple Graph**: Shows the complete entity-relationship network as a force-directed graph. Nodes represent entities, edges represent relationships.

**Entity Frequency**: Displays statistical analysis of entity occurrences, including bar charts and type distributions.

**Relationship Patterns**: Identifies recurring relationship patterns (e.g., "Person works_at Organization") with strength scores.

**Typed Holes**: Lists typed holes prioritized by complexity and dependencies, helping you plan implementation order.

## View Modes

### 1. Triple Graph View

The Triple Graph provides an interactive visualization of your semantic knowledge graph using a force-directed layout.

#### Features

- **Nodes**: Entities color-coded by type
  - Person: Blue
  - Organization: Green
  - Place: Yellow
  - Technology: Red
  - Concept: Purple
  - Thing: Gray
- **Edges**: Relationships with predicate labels and confidence scores
- **Selection**: Click or press Enter to select nodes
- **Highlighting**: Connected nodes are highlighted when a node is selected
- **Details Panel**: Shows entity information in the top-right corner

#### Navigation

**Pan (Move the View)**:
- `h` or `Left Arrow`: Pan left
- `l` or `Right Arrow`: Pan right
- `k` or `Up Arrow`: Pan up
- `j` or `Down Arrow`: Pan down

**Zoom**:
- `+` or `=`: Zoom in
- `-` or `_`: Zoom out

**Selection**:
- `Enter` or `Space`: Select/cycle through nodes
- `Esc`: Clear selection

**View Controls**:
- `c`: Center view on graph
- `r`: Reset view (zoom and position)
- `s`: Stabilize layout (run more layout iterations)

#### Understanding the Graph

**Node Size**: Nodes with higher importance scores appear with star (★) prefix

**Edge Labels**: Show relationship predicates (e.g., "works_at", "creates")

**Confidence Scores**: Displayed on edges as decimal values (0.0-1.0)

**Connectivity**: Highly connected nodes are central to your knowledge graph

#### Statistics Display

The footer shows real-time statistics:
- **Nodes**: Total number of entities in view
- **Edges**: Total number of relationships
- **Layout**: Number of layout iterations performed
- **Zoom**: Current zoom level
- **Offset**: Current pan position

### 2. Entity Frequency View

Entity Frequency View provides statistical analysis of entity occurrences across your documents.

#### Components

**Top Entities Bar Chart**:
- Horizontal bar chart showing the most frequent entities
- Bars are scaled relative to the maximum frequency
- Shows entity text, visual bar, and count

**Entity Type Distribution**:
- Breakdown of entities by type (Person, Organization, Place, etc.)
- Shows count for each entity type
- Sorted by frequency (descending)

#### Interpreting Results

**High Frequency Entities**: These are central to your documents and likely represent key concepts, people, or technologies.

**Type Distribution**: Reveals the focus of your content:
- High Person/Organization counts: People-centric content
- High Technology counts: Technical documentation
- High Concept counts: Abstract or theoretical content

**Importance Scores**: Calculated from frequency (logarithmic) plus type-based bonus:
- Person/Organization: +2 bonus
- Technology/Concept: +1 bonus
- Place/Thing: No bonus

### 3. Relationship Patterns View

The Patterns View discovers and displays recurring relationship patterns in your semantic graph.

#### Pattern Structure

Each pattern represents:
- **Subject Type**: Type of the subject entity
- **Predicate**: The relationship verb/action
- **Object Type**: Type of the object entity
- **Occurrences**: How many times this pattern appears
- **Avg Confidence**: Average confidence across all instances
- **Strength**: Calculated pattern strength score

#### Pattern Display

**Compact Mode**:
```
1. [Person] → works_at → [Organization] | Count: 12 | Conf: 0.87 | Str: 0.045
```

**Expanded Mode**:
```
Pattern 1: [Person] → works_at → [Organization]
  Occurrences: 12  Avg Confidence: 0.87  Strength: 0.045
  Examples:
    • John → works_at → Acme Corp (conf: 0.85)
    • Alice → works_at → Tech Inc (conf: 0.89)
    • Bob → works_at → StartupXYZ (conf: 0.88)
```

#### Strength Calculation

Pattern strength is calculated as:
```
strength = (occurrences × avg_confidence × (1 + diversity_factor)) / total_relationships
```

Where `diversity_factor = number_of_examples / max_examples`

Higher strength indicates more reliable and common patterns.

#### Pattern Sorting

Patterns can be sorted by:
- **Strength** (default): Most reliable patterns first
- **Frequency**: Most common patterns first
- **Confidence**: Highest confidence patterns first
- **Predicate**: Alphabetical by relationship type

#### Use Cases

- **Documentation**: Identify common relationship types to document
- **Validation**: Verify that expected patterns exist
- **Discovery**: Find unexpected patterns that might indicate new insights
- **Templates**: Use patterns as templates for suggesting new relationships

### 4. Typed Holes View

The Typed Holes View helps prioritize implementation work by analyzing typed hole complexity and dependencies.

#### Priority List

Displays typed holes sorted by priority score:
```
1. [██████████] ??AuthService (Complexity: 8)
2. [█████████░] ??DatabaseLayer (Complexity: 7)
3. [████████░░] ??ConfigLoader (Complexity: 4)
```

Each entry shows:
- **Priority Bar**: Visual indicator (10 blocks max)
- **Name**: Typed hole identifier
- **Complexity**: Implementation complexity score (0-10)

#### Priority Scoring

Priority is calculated from:
```
priority = base_priority(5)
         + dependency_count × 2
         + constraint_complexity
         + min(mention_count, 3)
```

Higher priority holes should be implemented first.

#### Complexity Scoring

Complexity is calculated from:
```
complexity = base_complexity(3)
           + constraint_count
           + special_constraint_bonus
```

Special constraint bonuses:
- Thread-safe/concurrent: +2
- Async: +1
- Pure/idempotent: +1

#### Dependency Analysis

**Circular Dependencies**: Warning displayed if circular dependencies detected

**Dependency Chains**: Shows which holes depend on others

**Implementation Order**: Suggests topologically-sorted implementation order

**Critical Path**: Identifies the longest dependency chain by complexity

#### Viewing Details

- **Navigate**: Use `↑/↓` arrow keys to select holes
- **Expand**: Press `Enter` to see detailed information
- **Dependencies**: Press `d` to show dependency tree
- **Export**: Press `e` to export implementation roadmap

## Keyboard Shortcuts

### Global Navigation

| Key | Action |
|-----|--------|
| `Tab` | Cycle to next view |
| `Shift+Tab` | Cycle to previous view |
| `1` | Switch to Triple Graph view |
| `2` | Switch to Entity Frequency view |
| `3` | Switch to Relationship Patterns view |
| `4` | Switch to Typed Holes view |
| `?` | Show help overlay |
| `q` | Quit Analyze Mode (return to Edit Mode) |

### View Controls (Triple Graph)

| Key | Action |
|-----|--------|
| `h` or `←` | Pan left |
| `j` or `↓` | Pan down |
| `k` or `↑` | Pan up |
| `l` or `→` | Pan right |
| `+` or `=` | Zoom in |
| `-` or `_` | Zoom out |
| `Enter` | Select node (cycles through nodes) |
| `Space` | Select node (same as Enter) |
| `Esc` | Clear selection |
| `c` | Center view on graph |
| `r` | Reset view (zoom + position) |
| `s` | Stabilize layout (run 10 more iterations) |

### View Controls (Lists)

| Key | Action |
|-----|--------|
| `↑` or `k` | Navigate up |
| `↓` or `j` | Navigate down |
| `PgUp` | Page up |
| `PgDown` | Page down |
| `Home` | Go to top |
| `End` | Go to bottom |
| `Enter` | Expand/collapse item |

### Export Functionality

| Key | Action |
|-----|--------|
| `e` | Export report in current format |
| `Ctrl+E` | Cycle export format (MD → HTML → PDF) |

## Export Functionality

Analyze Mode supports exporting reports in three formats: Markdown, HTML, and PDF.

### Supported Formats

#### Markdown (`.md`)

- GitHub-flavored markdown
- Mermaid diagrams for graphs
- Tables for statistics
- Code blocks for typed holes
- Ideal for documentation and version control

#### HTML (`.html`)

- Responsive design
- Interactive charts (Chart.js/D3.js)
- Searchable tables
- Print-optimized CSS
- Ideal for sharing and presentation

#### PDF (`.pdf`)

- Professional layout (A4 page size)
- Embedded charts as images
- Table of contents with hyperlinks
- Syntax highlighting
- Ideal for formal reports and archival

### Export Process

1. **Select View**: Navigate to the view you want to export (or export all)
2. **Choose Format**: Press `Ctrl+E` to cycle through formats (status shown in footer)
3. **Export**: Press `e` to generate the report
4. **Location**: Reports are saved to `./analysis-reports/` by default

### Output Structure

All exports include:

1. **Executive Summary**: Key metrics and highlights
2. **Entity Analysis**: Frequency charts, top entities, type distribution
3. **Relationship Patterns**: Common patterns, strength scores, examples
4. **Typed Hole Report**: Priority queue, dependency tree, implementation roadmap
5. **Recommendations**: Suggested next steps (if applicable)

### Export Examples

**Markdown Export**:
```markdown
# Semantic Analysis Report
Generated: 2025-11-08 14:30:00

## Executive Summary
- Total Entities: 135
- Total Relationships: 87
- Typed Holes: 23
- Top Entity Types: Person (45), Organization (32), Technology (23)

## Entity Analysis
### Top 10 Entities
1. API (15 occurrences)
2. John (12 occurrences)
3. Database (10 occurrences)
...

## Relationship Patterns
### Pattern: [Person] works_at [Organization]
- Occurrences: 12
- Confidence: 0.87
- Examples:
  - John works_at Acme Corp
  - Alice works_at Tech Inc
...
```

**HTML Export** includes interactive charts and searchable tables.

**PDF Export** provides a professionally formatted document suitable for printing.

### Customizing Exports

Export settings can be configured in `config.toml`:

```toml
[analyze]
export_dir = "./analysis-reports"
pdf_template = "default"
html_theme = "light"  # or "dark"
```

## Tips and Best Practices

### Effective Graph Exploration

1. **Start with Reset**: Press `r` to reset the view before exploring
2. **Center Important Nodes**: Select a node and press `c` to center it
3. **Stabilize First**: Press `s` multiple times to stabilize the layout before exploring
4. **Use Selection**: Select nodes to see their connections highlighted
5. **Zoom for Detail**: Zoom in to read entity labels clearly

### Finding Insights

1. **High-Frequency Entities**: Look for entities that appear frequently - they're central to your content
2. **Strong Patterns**: Patterns with high strength scores represent reliable relationships
3. **Critical Path Holes**: Focus on holes in the critical path first
4. **Type Distribution**: Understand your content focus from entity type distribution

### Performance Tips

1. **Filter Large Graphs**: Use filters to focus on specific entity types or high-confidence relationships
2. **Limit View Size**: Large graphs (1000+ nodes) may benefit from filtering
3. **Export Incrementally**: Export individual views rather than full reports for faster generation

### Workflow Recommendations

**For Documentation Writers**:
1. Start with Entity Frequency to identify key concepts
2. Use Relationship Patterns to document common structures
3. Export to Markdown for inclusion in docs

**For Developers**:
1. Check Typed Holes view for implementation priorities
2. Review dependencies to plan work order
3. Export implementation roadmap to PDF

**For Analysts**:
1. Explore Triple Graph for network structure
2. Analyze patterns for insights
3. Export to HTML for interactive exploration

## Troubleshooting

### Graph Not Displaying

**Cause**: No semantic data loaded
**Solution**: Ensure you've analyzed a document in Edit Mode first

### Layout Unstable

**Cause**: Insufficient layout iterations
**Solution**: Press `s` multiple times to stabilize

### Performance Issues

**Cause**: Too many entities in view
**Solution**: Apply filters to reduce complexity

### Export Fails

**Cause**: Missing export directory or permissions
**Solution**: Check `export_dir` setting and file permissions

### Empty Views

**Cause**: No data matching current filters
**Solution**: Clear filters or adjust filter criteria

## Advanced Features

### Filter Syntax

Filters can be applied to:
- **Entity Types**: Show only specific types (Person, Organization, etc.)
- **Confidence**: Minimum confidence threshold (0.0-1.0)
- **Text Search**: Filter by entity text (case-insensitive)

### Custom Scoring

Importance and priority scores can be influenced by:
- Entity frequency (logarithmic scale)
- Entity type (type-specific bonuses)
- Dependency count (for typed holes)
- Constraint complexity (for typed holes)

### Pattern Clustering

Patterns with similar predicates are automatically clustered:
- Uses Levenshtein distance for similarity
- Threshold configurable (default: 0.7)
- Helps identify related relationship types

## Configuration

Analyze Mode behavior can be customized in `config.toml`:

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

## Next Steps

- **Explore Each View**: Spend time in each view to understand your data from different angles
- **Export Reports**: Generate reports to share insights with your team
- **Iterate**: Use insights from Analyze Mode to refine your documents
- **Plan Implementation**: Use Typed Holes view to prioritize development work

## See Also

- [API Documentation](./analyze-mode-api.md) - For developers integrating with Analyze Mode
- [Phase 6 Specification](./PHASE6_SPEC.md) - Technical details and architecture
- [Usage Guide](./USAGE.md) - General application usage
