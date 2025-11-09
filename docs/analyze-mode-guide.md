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

**Filter Examples**:
```
Filter entity type: "Organization"
Filter confidence: 0.75
Filter text: "Python"
Combined: "Person" OR "Technology" with confidence >= 0.8
```

### Custom Scoring

Importance and priority scores can be influenced by:
- Entity frequency (logarithmic scale)
- Entity type (type-specific bonuses)
- Dependency count (for typed holes)
- Constraint complexity (for typed holes)

**Scoring Formula**:
```
Entity Importance = log10(frequency) × 3 + type_bonus

Priority Score = base_priority(5)
               + dependency_count × 2
               + constraint_complexity
               + min(mention_count, 3)

Complexity Score = base_complexity(3)
                 + constraint_count
                 + special_constraint_bonus
```

### Pattern Clustering

Patterns with similar predicates are automatically clustered:
- Uses Levenshtein distance for similarity
- Threshold configurable (default: 0.7)
- Helps identify related relationship types

**Example**: Predicates like "creates", "created", "creates_new" will be clustered together.

### Graph Layout Algorithm

The force-directed layout uses:
- **Repulsive forces**: Nodes push each other away (avoids overlaps)
- **Attractive forces**: Edges pull connected nodes together
- **Damping**: Reduces oscillations and stabilizes layout (default: 0.8)
- **Iterations**: Each iteration refines positions (50-200 recommended)

**Performance Tips**:
- Increase iterations for better layout (press `s` multiple times)
- Reduce max entities if layout is slow (use filters)
- Adjust zoom for better visibility of labels

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

# Triple Graph settings
force_strength = 0.5
repulsion_strength = 100.0
damping_factor = 0.8
```

## Common Workflows

### Workflow 1: Document Structure Analysis

**Goal**: Understand how concepts relate in a document

1. Load document in Edit Mode
2. Run semantic analysis
3. Enter Analyze Mode (press `3`)
4. View Triple Graph (default)
5. Press `r` to reset and stabilize layout
6. Select key entities and explore connections
7. Switch to Relationship Patterns view
8. Export to Markdown for documentation

**Key Metrics to Watch**:
- Highly connected nodes (important concepts)
- Relationship patterns (document structure)
- Confidence scores (relationship strength)

### Workflow 2: Implementation Planning

**Goal**: Prioritize development work using typed holes

1. Analyze document with typed holes
2. Enter Analyze Mode (press `3`)
3. Switch to Typed Holes view (press `4`)
4. Review complexity scores and dependencies
5. Check for circular dependencies (warning displayed)
6. Use implementation order suggestion
7. Export implementation roadmap (press `e`)
8. Use roadmap to plan sprints

**Key Decisions**:
- Start with high-priority holes
- Resolve dependencies before dependent holes
- Watch for critical path holes

### Workflow 3: Pattern Discovery and Validation

**Goal**: Identify and validate relationship patterns

1. Analyze multi-document corpus
2. View Relationship Patterns (press `3`)
3. Sort by strength (default, strongest first)
4. Review examples for each pattern
5. Look for unexpected patterns (new insights)
6. Filter by predicate type (press `f`)
7. Export to HTML for team review
8. Use patterns to establish coding conventions

**Success Indicators**:
- Strong patterns (strength > 0.05) indicate common structures
- Consistent confidence scores mean reliable patterns
- Multiple diverse examples show pattern robustness

### Workflow 4: Knowledge Graph Exploration

**Goal**: Understand entity relationships and clustering

1. Load large document or corpus
2. View Triple Graph
3. Apply entity type filters (e.g., "Person" only)
4. Press `r` to reset view
5. Press `s` multiple times to stabilize
6. Select central nodes (high degree)
7. Use zoom to explore clusters
8. Identify communities (groups of connected entities)

**Interpretation**:
- Clustered entities indicate communities of related concepts
- Bridge nodes (connecting clusters) are important
- Isolated nodes may indicate missing relationships

## FAQ (Frequently Asked Questions)

**Q: Why is the triple graph layout unstable?**
A: The force-directed layout needs iterations to stabilize. Press `s` several times (each adds 10 iterations). The physics simulation converges over time.

**Q: How do I export only specific entity types?**
A: Use the entity type filter (press `f`) before exporting. The export will include only filtered entities and their relationships.

**Q: Can I combine multiple filters?**
A: Filters are AND-ed together currently. Confidence filter + entity type filter means: "entity type X AND confidence >= threshold"

**Q: What does confidence score mean?**
A: Confidence (0.0-1.0) indicates the NER model's certainty that this is a real entity/relationship. Higher = more certain. Typical threshold: 0.5-0.7.

**Q: How are importance scores calculated?**
A: Importance = log10(frequency) × 3 + type_bonus. Frequently occurring entities have higher importance. Type bonuses: Person/Org +2, Tech/Concept +1, Place/Thing +0.

**Q: Can I share analysis results with my team?**
A: Yes! Export to HTML for interactive viewing in a web browser. HTML exports include searchable tables and interactive charts.

**Q: What's the difference between Pattern Strength and Confidence?**
A: Confidence is per-relationship (varies). Strength is pattern-level aggregate. High-strength pattern = common, reliable relationship type.

**Q: Why do some entities not appear in the graph?**
A: Check filters (entity type, confidence threshold). Low-confidence entities may be filtered. Use `clear filters` to see all.

**Q: Can I edit the graph manually?**
A: Currently graphs are read-only (derived from analysis). To change entities/relationships, modify source documents and re-analyze.

**Q: How do I handle large graphs (1000+ entities)?**
A: Apply aggressive filters: high confidence threshold (0.8+), specific entity types, text search. Export individual views instead of full reports.

**Q: What export format should I use for presentations?**
A: HTML for interactive presentations. PDF for static, printable documents. Markdown for version control integration.

## Performance Optimization

### For Large Datasets (1000+ entities)

1. **Apply Aggressive Filters**:
   - Confidence threshold: 0.75+
   - Entity types: Select 2-3 most relevant
   - Text search: Focus on keywords

2. **Optimize Rendering**:
   - Reduce layout iterations (press `s` fewer times)
   - Use zoom selectively
   - Export to separate files per view

3. **Manage Memory**:
   - Process large documents in batches
   - Export results incrementally
   - Clear old analysis data periodically

### For Complex Relationship Graphs

1. **Visualize Subgraphs**:
   - Select a node and view connected components
   - Use text search to focus on specific patterns
   - Filter by relationship type

2. **Export Strategy**:
   - Export Patterns view separately (lightweight)
   - Use Markdown for tabular data
   - Generate HTML for interactive exploration

## Integration with mnemosyne

Analyze Mode can integrate with mnemosyne memory system:

**Storing Analysis Results** (future feature):
```bash
mnemosyne remember -c "Document analysis: 45 entities, 12 key patterns" \
  -n "project:docs" -i 8 -t "analysis,documentation"
```

**Recalling Previous Analyses**:
```bash
mnemosyne recall -q "entity relationships" -n "project:docs" -l 5
```

**Memory Evolution**:
```bash
mnemosyne evolve  # Consolidate related memories, decay old analysis
```

Currently, analysis results are stored in `./analysis-reports/` directory. Future versions will integrate with mnemosyne for persistent, queryable analysis memory.

## Troubleshooting Guide

### Issue: Graph Does Not Display

**Symptoms**: Blank canvas, "No graph data to display" message

**Possible Causes**:
- No semantic analysis performed
- Filters too restrictive (zero entities match)
- Graph loading failed silently

**Solutions**:
1. Ensure document has been analyzed (Edit Mode → run analysis)
2. Clear filters: press `f` → select "clear all filters"
3. Check filter settings (confidence, entity types)
4. Try switching views (Tab) and switching back

### Issue: Layout Unstable or Oscillating

**Symptoms**: Nodes jump around, don't settle in place

**Causes**:
- Insufficient layout iterations
- High repulsion forces in large graphs
- Natural chaos in force-directed layouts (normal initially)

**Solutions**:
1. Press `s` multiple times to run more iterations (10 per press)
2. Press `r` to reset and restart from initial layout
3. Reduce entity count with filters
4. Be patient (small oscillations are normal)

### Issue: Performance Degradation

**Symptoms**: Slow panning, zooming, view switching

**Causes**:
- Too many entities (1000+)
- Complex relationship graph
- Insufficient system resources

**Solutions**:
1. Apply filters to reduce entity count to 100-300
2. Switch to lighter views (Entity Frequency instead of Triple Graph)
3. Close other applications
4. Export and analyze incrementally

### Issue: Export Takes Too Long

**Symptoms**: Long wait time when exporting, no feedback

**Causes**:
- Large number of entities/relationships
- PDF generation (most CPU-intensive)
- File I/O bottleneck

**Solutions**:
1. Export to Markdown first (fastest)
2. Reduce entity count with filters
3. Export individual views, not full report
4. Check disk space and permissions

### Issue: Text Search Not Finding Entities

**Symptoms**: Filter applied but no results shown

**Causes**:
- Typos in search term
- Case sensitivity issues
- Partial matches not supported

**Solutions**:
1. Double-check spelling
2. Try partial matches (e.g., "Acme" instead of "Acme Corp")
3. Use filters instead of text search for type-based filtering

### Issue: Confidence Scores Seem Inconsistent

**Symptoms**: Similar entities/relationships have different confidence scores

**Causes**:
- Different source documents have different confidence
- NER model confidence varies by context
- Aggregation from multiple analyses

**Solutions**:
1. This is expected behavior
2. Focus on patterns and clusters, not individual scores
3. Use confidence threshold filter to focus on high-confidence data

## Next Steps

- **Explore Each View**: Spend time in each view to understand your data from different angles
- **Export Reports**: Generate reports to share insights with your team
- **Iterate**: Use insights from Analyze Mode to refine your documents
- **Plan Implementation**: Use Typed Holes view to prioritize development work
- **Try Different Filters**: Experiment with entity type and confidence filters
- **Share Results**: Export to HTML and share with stakeholders

## See Also

- [API Documentation](./analyze-mode-api.md) - For developers integrating with Analyze Mode
- [Orchestrate Mode Guide](./orchestrate-mode-guide.md) - For multi-agent coordination
- [Usage Guide](./USAGE.md) - General application usage
- [mnemosyne Documentation](https://github.com/steveyegge/mnemosyne) - Memory and orchestration system
