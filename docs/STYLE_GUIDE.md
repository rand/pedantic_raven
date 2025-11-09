# Pedantic Raven Style Guide

Version: 1.0
Last Updated: November 2024

This document provides comprehensive styling guidelines for Pedantic Raven's terminal UI. All UI elements should follow these standards to ensure visual consistency, accessibility, and excellent user experience across all five modes (Edit, Analyze, Orchestrate, Explore, Dashboard).

## Table of Contents

1. [Color Palette](#color-palette)
2. [Typography & Text Styles](#typography--text-styles)
3. [Spacing & Layout](#spacing--layout)
4. [Component Patterns](#component-patterns)
5. [Error Messages](#error-messages)
6. [Loading States](#loading-states)
7. [Help System](#help-system)
8. [Accessibility](#accessibility)
9. [Implementation Examples](#implementation-examples)

---

## Color Palette

All colors are defined using lipgloss color codes for consistency. The palette is organized by semantic meaning.

### Primary Colors (Core UI)

| Color | Code | Usage | Example |
|-------|------|-------|---------|
| Bright Blue | `39` | Primary accents, headers, focus states | Memory importance (6), primary emphasis |
| Light Gray | `252` | Default text, body content | Normal list items, standard text |
| Dark Gray | `235` | Backgrounds, base layer | Header/footer backgrounds |
| Medium Gray | `243` | Secondary text, metadata | Timestamps, file sizes, counts |
| Darker Gray | `240` | Subtle borders, inactive state | Unfocused borders, disabled elements |

### Semantic Colors (Meaning & Status)

| Semantic | Code | Hex | Usage | Context |
|----------|------|-----|-------|---------|
| Success/Green | `34` or `118` | `#00AA00` | Positive actions, confirmed state | Memory importance (6), successful operations |
| Warning/Yellow | `226` or `229` | `#FFAA00` | Caution, intermediate priority | Memory importance (7), alerts |
| Error/Red | `196` | `#FF0000` | Error messages, failures, critical | Memory importance (9), errors, failed tasks |
| Danger/Orange | `208` | | Mid-level warnings | Memory importance (8), warnings |
| Cyan/Info | `81` | `#00AAFF` | Information, links, help | Memory importance (5), informational text |
| Magenta/Purple | `170` | `#6600FF` | Focused state, special emphasis | Focused borders, orchestrate mode |
| White | `15` or `#FFFFFF` | `#FFFFFF` | High contrast text, emphasis | Headers on dark backgrounds |

### Memory Importance Color Mapping

The importance scale (1-10) uses a dedicated color palette:

```go
importanceColors = map[uint32]lipgloss.Color{
    9: lipgloss.Color("196"), // Red (critical)
    8: lipgloss.Color("208"), // Orange (very high)
    7: lipgloss.Color("226"), // Yellow (high)
    6: lipgloss.Color("118"), // Green (moderate)
    5: lipgloss.Color("81"),  // Cyan (normal)
    4: lipgloss.Color("12"),  // Blue (low) - if implemented
    3: lipgloss.Color("7"),   // Gray (very low) - if implemented
    2: lipgloss.Color("8"),   // Minimal
    1: lipgloss.Color("0"),   // Minimal
}
```

**Usage Rules:**
- Use memory importance colors only for representing actual importance values
- Never mix color semantics (e.g., don't use red for non-error states)
- Maintain sufficient contrast for accessibility (WCAG AA compliant)

---

## Typography & Text Styles

### Style Hierarchy

All text styles are built using lipgloss with clear hierarchy:

#### Headers (Level 1 - Section Titles)

```go
headerStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("39")).      // Bright blue
    Background(lipgloss.Color("235")).     // Dark background
    Padding(0, 1)
```

**When to use:** Section titles, mode names, major headings
**Example:** "Memory Workspace", "Edit Mode", "Analysis Results"

#### Subheaders (Level 2 - Component Titles)

```go
subheaderStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("170"))      // Magenta
```

**When to use:** Subsection titles, component names, pane labels
**Example:** "Linked Memories", "Task Queue", "Agent Status"

#### Normal Text (Body)

```go
normalStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("252"))      // Light gray
```

**When to use:** Default content, list items, descriptions
**Example:** Memory content, task descriptions, regular output

#### Metadata Text (Secondary)

```go
metaStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("243"))      // Medium gray
```

**When to use:** Supplementary information, timestamps, counts
**Example:** "Updated 2h ago", "5 items", "Confidence: 0.95"

#### Error Text

```go
errorStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("196")).     // Red
    Bold(true)
```

**When to use:** Error messages, failures, alerts
**Example:** "Error: Failed to load memory", "Invalid input"

#### Success Text

```go
successStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("34")).      // Green
    Bold(true)
```

**When to use:** Success confirmations, positive outcomes
**Example:** "✓ Memory created", "Analysis complete"

#### Warning Text

```go
warningStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("226")).     // Yellow
    Bold(true)
```

**When to use:** Cautions, warnings, attention-needed items
**Example:** "⚠ Unsaved changes", "Connection lost"

#### Info Text

```go
infoStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("81"))       // Cyan
```

**When to use:** Informational messages, hints, explanations
**Example:** "Tip: Use Ctrl+S to save", "Loading..."

#### Loading/Spinning Text

```go
loadingStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("39")).      // Bright blue
    Bold(true)
```

**When to use:** Loading indicators, pending states
**Example:** "Loading...", "Processing..."

#### Code/Monospace Text

For code snippets and technical content, use monospace with appropriate colors:

```go
codeStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("118"))      // Green
```

**When to use:** File paths, command syntax, technical identifiers
**Example:** `/home/user/.config`, `ctrl+s`

#### Focus Indicator Text

```go
focusedStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("170"))      // Magenta
    Bold(true)
```

**When to use:** Highlighting current selection, focused elements
**Example:** Highlighted list item text, selected command

#### Emphasis/Highlight Text

```go
emphasisStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("229")).     // Light yellow
    Bold(true)
```

**When to use:** Highlighting important content, search matches
**Example:** Search term highlights, important keywords

### Font Attributes

- **Bold**: Use for headers, errors, warnings, success messages, and emphasis
- **Italic**: Use sparingly for empty states, supplementary notes
- **Underline**: Avoid in terminal UI (limited support across terminals)
- **Reverse**: Avoid (use background colors instead)

---

## Spacing & Layout

### Padding Standards

Lipgloss padding follows `Padding(vertical, horizontal)` format:

```go
// No padding - used for compact layouts
style.Padding(0)

// Minimal padding - used for tight components
style.Padding(0, 1)        // 0 vertical, 1 horizontal

// Standard padding - default for components
style.Padding(1)           // 1 all sides (1 vertical, 1 horizontal)

// Generous padding - used for dialogs, overlays
style.Padding(1, 2)        // 1 vertical, 2 horizontal

// Spacious padding - used for major sections
style.Padding(2, 2)        // 2 all sides
```

### Margins Between Elements

- **Between list items:** 0 lines (compact, separated by horizontal rule if needed)
- **Between sections:** 1 blank line
- **Between panes:** 0 lines (separated by border)
- **Before/after headers:** 0 lines above, 1 line below
- **Before/after footers:** 1 line above, 0 lines below

### Border Styles

```go
// Rounded borders - default for component panes
border := lipgloss.RoundedBorder()

// Solid borders - use for emphasis
border := lipgloss.NormalBorder()

// No border - for minimal layouts
border := lipgloss.NoBorder()
```

### Border Colors

```go
// Focused state - Magenta
BorderForeground(lipgloss.Color("170"))

// Normal state - Dark gray
BorderForeground(lipgloss.Color("240"))

// Error state - Red
BorderForeground(lipgloss.Color("196"))

// Info state - Cyan
BorderForeground(lipgloss.Color("81"))
```

### Indentation

- **Nested content:** 2 spaces per level
- **Menu items:** 2 spaces before text
- **Error/info icons:** 2 spaces before text
- **List items:** 2 spaces before text (accommodates icon + space)

---

## Component Patterns

### Boxes (Bordered Containers)

Standard box pattern for grouping related content:

```go
boxStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("240")).
    Padding(1)

// Apply to content
box := boxStyle.Render(content)
```

**Usage:** Panes, panels, sections
**Example:** Editor pane, context panel, agent status

### Selected/Focused Box

```go
focusedBoxStyle = boxStyle.
    BorderForeground(lipgloss.Color("170"))  // Magenta when focused
```

### Error Box

```go
errorBoxStyle = boxStyle.
    BorderForeground(lipgloss.Color("196")).  // Red border
    Render(content)
```

### List Items

Standard list pattern:

```
│ > Item Name                          │
│   └─ Metadata or description        │
│                                      │
│   Item Name (not selected)           │
│   └─ Metadata or description        │
```

**Selection indicator:** `> ` (chevron + space)
**Unselected prefix:** `  ` (two spaces)
**Indent for details:** 3 spaces

### Tables/Data Grids

For structured data display:

```
┌─ Header 1 ─ Header 2 ─ Header 3 ─┐
├──────────────────────────────────┤
│ Row 1    │ Data     │ Data        │
│ Row 2    │ Data     │ Data        │
└──────────────────────────────────┘
```

**Rules:**
- Headers: Bold, blue (color 39) background
- Separator: use ├─ and ┼─ characters
- Columns: separated by `│`
- Selected row: highlighted or marked with `> `

### Progress Bars

```go
progressStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("34"))  // Green
```

**Pattern:**
```
[████████░░░░] 60% Complete
```

### Buttons/Actions

In help text and overlays:

```
[Enter] Confirm    [Esc] Cancel    [Tab] Next
```

**Pattern:** `[Key] Action` separated by spaces or pipes

### Status Indicators

Use single-character indicators with color:

```go
// Success
successIndicator = "✓"  // or "•" for bullet

// Error
errorIndicator = "✗"   // or "!" or "×"

// Warning
warningIndicator = "⚠"  // or "!"

// Info
infoIndicator = "ℹ"   // or "i" or "→"

// Loading
loadingIndicator = "◐" or cycle through "◐◓◑◒"
```

### Empty States

When no content is available:

```
(No memories loaded)
```

**Pattern:** `(Text in parentheses)` with italic gray style

---

## Error Messages

All error messages must follow consistent patterns for clarity and user assistance.

### Error Message Structure

Every error message should follow this format:

```
[ICON] Error Category: Specific problem description.
       Action suggestion (if applicable).
```

### Error Categories

| Category | Icon | Color | Typical Causes |
|----------|------|-------|---------|
| Validation | ✗ | Red | Invalid input, missing fields |
| Connection | ✗ | Red | Network issues, server unavailable |
| Permission | ✗ | Red | Access denied, authentication failed |
| Resource | ✗ | Red | File not found, quota exceeded |
| Operation | ✗ | Red | Function failed, unexpected state |
| Config | ⚠ | Orange | Configuration error, invalid settings |

### Error Message Examples

#### Input Validation Errors

**Bad:**
```
error: content required
```

**Good:**
```
✗ Validation Error: Content is required.
  Tip: Enter at least 1 character before saving.
```

**Code:**
```go
errorMsg := errorStyle.Render("✗ Validation Error: Content is required.") + "\n" +
            metaStyle.Render("  Tip: Enter at least 1 character before saving.")
```

#### Connection Errors

**Bad:**
```
error connecting to server
```

**Good:**
```
✗ Connection Error: Cannot reach mnemosyne server at localhost:50051.
  Try: Check if server is running (./mnemosyne-server).
       Or configure different host/port in config.toml.
```

**Code:**
```go
msg := fmt.Sprintf("✗ Connection Error: Cannot reach mnemosyne server at %s:%d.\n", host, port) +
       "Try: Check if server is running (./mnemosyne-server).\n" +
       "     Or configure different host/port in config.toml."
return errorStyle.Render(msg)
```

#### File System Errors

**Bad:**
```
file not found
```

**Good:**
```
✗ File Error: Cannot open /path/to/file.txt (not found).
  Try: Check the file path and permissions.
       Or use file picker to browse (Ctrl+O).
```

#### Timeout Errors

**Bad:**
```
timeout waiting for response
```

**Good:**
```
✗ Timeout Error: Operation took too long (>30s).
  Try: Check network connection.
       Or try again in a moment.
       Or increase timeout in config.
```

#### API/Service Errors

**Bad:**
```
extraction failed
```

**Good:**
```
✗ Extraction Error: GLiNER service failed to process text.
  Possible causes:
    • Model not loaded: Check server status (./gliner-server)
    • Invalid text: Empty or malformed input
    • Rate limit: Server overloaded, wait and retry
  Try: Reload model (Ctrl+Shift+R) or check logs.
```

### Error Message Best Practices

1. **Be specific**: "Content is required" not "Invalid input"
2. **Be actionable**: "Check if server is running" not "Server unavailable"
3. **Provide context**: Include what failed and where
4. **Use appropriate icon**: ✗ for errors, ⚠ for warnings
5. **Suggest recovery**: "Try" or "You can" with concrete actions
6. **Keep it concise**: 1-3 lines maximum (+ optional tips)

### Standard Error List

Below are 50+ standardized error messages for common scenarios:

#### Validation Errors (10 messages)

```
✗ Validation: Content is required. Enter at least 1 character.
✗ Validation: Content too long (max 10000 chars). Truncate and try again.
✗ Validation: Importance must be 1-10. Enter a number in this range.
✗ Validation: Too many tags (max 20). Remove some and retry.
✗ Validation: Namespace is required. Use format: project:name.
✗ Validation: Link strength must be 0.0-1.0. Use decimal value.
✗ Validation: Cannot link memory to itself. Choose different target.
✗ Validation: Invalid tag format. Use alphanumeric and hyphens only.
✗ Validation: Search query is empty. Enter search terms.
✗ Validation: Duplicate tag. This tag already exists.
```

#### Connection Errors (10 messages)

```
✗ Connection: Cannot reach mnemosyne at localhost:50051. Check if running.
✗ Connection: Cannot reach GLiNER at localhost:8000. Check if running.
✗ Connection: Host cannot be empty. Check config.toml.
✗ Connection: Port must be 1-65535. Got invalid value.
✗ Connection: Timeout connecting to server (30s). Check network.
✗ Connection: Request timeout. Server not responding. Retry or check logs.
✗ Connection: Connection reset by server. Reconnecting...
✗ Connection: DNS resolution failed. Check hostname.
✗ Connection: Connection refused. Is server listening?
✗ Connection: SSL certificate error. Check certificate or use http.
```

#### Memory/Data Errors (10 messages)

```
✗ Memory: Not found (ID: abc123). May have been deleted.
✗ Memory: Cannot create duplicate. This memory already exists.
✗ Memory: Creation failed on server. Check logs for details.
✗ Memory: Update failed. Memory may have changed. Refresh (Ctrl+R).
✗ Memory: Delete failed. Memory may not exist. Refresh (Ctrl+R).
✗ Memory: No memories found. Try adjusting search filters.
✗ Memory: Source memory not found. Refresh and retry.
✗ Memory: Target memory not found. Verify target exists.
✗ Memory: Link not found. May have been deleted. Refresh.
✗ Memory: Sync failed at operation 5. Check logs and retry.
```

#### File System Errors (10 messages)

```
✗ File: Cannot open /path/to/file (not found). Check path.
✗ File: Cannot read file. Permission denied. Check permissions.
✗ File: Cannot write file. Check disk space and permissions.
✗ File: File is too large (>10MB). Choose smaller file.
✗ File: Unsupported file format. Use .txt, .md, or .json.
✗ File: Encoding error. File must be UTF-8 encoded.
✗ File: Directory not found. Check path exists.
✗ File: Cannot create directory. Check permissions.
✗ File: File locked. Close file in other applications.
✗ File: Filename invalid. Use alphanumeric, dash, underscore only.
```

#### Service Errors (10+ messages)

```
✗ Service: GLiNER unavailable. Check gliner-server (./gliner-server).
✗ Service: Model not loaded. Reload with Ctrl+Shift+R.
✗ Service: Extraction failed. Text may be invalid or too long.
✗ Service: Rate limited. Wait 60s before retrying.
✗ Service: Semantic analyzer not configured. Check config.yaml.
✗ Service: Analysis timeout (>60s). Text may be too large.
✗ Service: Triple extraction failed. Invalid entity/relation format.
✗ Service: Pattern detection failed. Graph may be corrupted.
✗ Service: Export failed. Check file permissions.
✗ Service: Import failed. File format invalid.
✗ Service: Authentication failed. Check credentials.
✗ Service: Authorization denied. Insufficient permissions.
```

### Error Display in Each Mode

#### Edit Mode

```
┌─ Editor ────────────────────┐
│ [editor content]             │
└──────────────────────────────┘

✗ Error: Content is required.
```

**Position:** Below editor, in main area

#### Analyze Mode

```
┌─ Analysis Results ──────────────────┐
│ ✗ Error: Analysis failed.           │
│    Try: Check if GLiNER is running. │
└─────────────────────────────────────┘
```

**Position:** In content area where results would be

#### Orchestrate Mode

```
┌─ Orchestrate ──────────────────────────┐
│ Status: ✗ Error: Task execution failed │
│ ✗ Agent Worker-1: Timeout after 30s    │
└────────────────────────────────────────┘
```

**Position:** In dashboard or agent status area

#### Explore Mode

```
┌─ Memories ────────────────────┐
│ ✗ Error: Cannot load list.    │
│    Tip: Refresh (Ctrl+R)      │
└───────────────────────────────┘
```

**Position:** In content area

---

## Loading States

All asynchronous operations must show clear loading feedback.

### Loading Indicators

#### Text-Based (Simple)

```
Loading...
Processing...
Analyzing...
```

**Style:**
```go
loadingStyle.Render("Loading...")
```

#### Spinner (Animated)

```
◐ Loading...
◓ Loading...
◑ Loading...
◒ Loading...
```

**Pattern:** Cycle every 250ms through frames

#### Progress Bar (With Percentage)

```
Progress: [████████░░░░░░░░░░] 40%
```

**Code:**
```go
func renderProgressBar(current, total int, width int) string {
    percent := float64(current) / float64(total)
    filled := int(float64(width) * percent)
    bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
    pct := int(percent * 100)
    return fmt.Sprintf("[%s] %d%%", bar, pct)
}
```

#### Detailed Progress (Multi-step)

```
Loading memories...
  ├─ Connecting to server... ✓
  ├─ Fetching metadata... ○
  └─ Loading full content...

Estimated time: 2s
```

**Pattern:** Use ├─ and └─ for tree structure, show status per step

### Loading State Rules

1. **Always show immediately** - Don't wait more than 100ms
2. **Provide context** - "Loading memories...", not just "Loading..."
3. **Show duration for slow ops** - "Loading (5s)..." when >2s expected
4. **Use spinners** - Animate for visual feedback
5. **Provide abort option** - "Press Ctrl+C to cancel"
6. **Show progress** - When duration unknown, use spinner; if known, use percent

### Timeout Handling

If operation takes longer than expected:

```
Loading memories (10s)...
This is taking longer than expected.
Press Ctrl+C to cancel.
```

---

## Help System

Help overlays provide context and guidance for all modes.

### Help Overlay Pattern

All help overlays follow this consistent format:

```
╭─ Keyboard Shortcuts ──────────────────────────╮
│ Editor Navigation                            │
│  Arrows      Move cursor                     │
│  Ctrl+Home   Beginning of buffer             │
│  Ctrl+End    End of buffer                   │
│                                              │
│ Edit Actions                                 │
│  Ctrl+A      Select all                      │
│  Ctrl+X/C/V  Cut/copy/paste                  │
│  Ctrl+Z/Y    Undo/redo                       │
│                                              │
│ Mode Navigation                              │
│  Tab         Next mode                       │
│  Shift+Tab   Previous mode                   │
│                                              │
│ [Esc] Close  [Space] Scroll  [?] Toggle     │
╰────────────────────────────────────────────╯
```

**Style:**
```go
helpBoxStyle = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("81")).    // Cyan
    Background(lipgloss.Color("235")).        // Dark background
    Padding(1, 2)
```

### Help Content Structure

Every mode's help should include:

1. **Navigation** - How to move through the mode
2. **Actions** - Primary operations available
3. **Special Keys** - Ctrl/Shift combinations
4. **Mode Switching** - How to switch modes
5. **General** - Esc to close, ? to toggle help
6. **Tips** - Useful hints and shortcuts

### Help for Each Mode

#### Edit Mode Help

```
┌─ Edit Mode Help ──────────────────────────────┐
│                                               │
│ NAVIGATION                                   │
│  ↑↓←→ or vim keys    Cursor movement         │
│  Ctrl+Home/End       Start/end of document   │
│  Page Up/Down        Scroll by page          │
│                                              │
│ EDITING                                      │
│  Type                Insert text             │
│  Backspace/Delete    Remove character        │
│  Ctrl+A              Select all              │
│  Ctrl+X/C/V          Cut/copy/paste          │
│  Ctrl+Z/Y            Undo/redo               │
│                                              │
│ FILE OPERATIONS                              │
│  Ctrl+O              Open file               │
│  Ctrl+S              Save context            │
│                                              │
│ ANALYSIS                                     │
│  Ctrl+E              Analyze with semantic   │
│  Ctrl+Shift+T        Extract triples         │
│                                              │
│ MODE SWITCHING                               │
│  Tab / Shift+Tab     Next/previous mode      │
│  1-4                 Jump to mode number     │
│                                              │
│ OTHER                                        │
│  Cmd+P               Command palette         │
│  Esc                 Quit / Clear selection  │
│  ?                   Toggle this help        │
│                                              │
└───────────────────────────────────────────────┘
```

#### Analyze Mode Help

```
┌─ Analyze Mode Help ───────────────────────────┐
│                                               │
│ VIEW SWITCHING                               │
│  1 / 2 / 3 / 4      Switch between views    │
│  t                  Triple graph             │
│  e                  Entity frequency         │
│  p                  Patterns                 │
│  h                  Typed holes              │
│                                              │
│ INTERACTION                                  │
│  ↑↓←→               Navigate                 │
│  Enter              Select/drill down        │
│  Esc                Clear selection          │
│                                              │
│ EXPORT                                       │
│  Ctrl+E             Export as markdown       │
│  Ctrl+Shift+E       Export as HTML           │
│  Ctrl+Alt+E         Export as PDF            │
│                                              │
│ MODE SWITCHING                               │
│  Tab / Shift+Tab    Next/previous mode      │
│  1-4                Jump to mode number     │
│                                              │
│ OTHER                                        │
│  Cmd+P              Command palette         │
│  ?                  Toggle this help        │
│                                              │
└───────────────────────────────────────────────┘
```

#### Orchestrate Mode Help

```
┌─ Orchestrate Mode Help ──────────────────────┐
│                                               │
│ VIEW SWITCHING                               │
│  1 / 2 / 3 / 4      Plan/Dashboard/Graph/Log│
│  p                  Plan editor              │
│  d                  Dashboard                │
│  g                  Task graph               │
│  l                  Agent logs               │
│                                              │
│ ORCHESTRATION CONTROL                       │
│  Space              Pause/resume             │
│  Ctrl+Enter         Launch orchestration     │
│  Ctrl+C             Stop orchestration       │
│  r                  Restart                  │
│                                              │
│ AGENT INTERACTION                            │
│  ↑↓                 Select agent             │
│  Enter              View agent details       │
│  k                  Kill agent               │
│  d                  View agent logs          │
│                                              │
│ PLAN EDITING                                 │
│  e                  Edit plan                │
│  Ctrl+S             Save plan                │
│  Ctrl+L             Load plan                │
│                                              │
│ OTHER                                        │
│  Tab / Shift+Tab    Next/previous mode      │
│  Cmd+P              Command palette         │
│  ?                  Toggle this help        │
│                                              │
└───────────────────────────────────────────────┘
```

#### Explore Mode Help

```
┌─ Explore Mode Help ───────────────────────────┐
│                                               │
│ NAVIGATION                                   │
│  ↑↓                 Select memory            │
│  ← / →              Toggle layout / focus    │
│  Enter              View memory details      │
│  /                  Search memories         │
│                                              │
│ DISPLAY                                      │
│  Tab                Toggle list/detail/graph│
│  g                  Show graph view          │
│  l                  Show list view           │
│  d                  Show detail view         │
│                                              │
│ MEMORY OPERATIONS                            │
│  n                  New memory               │
│  e                  Edit memory              │
│  Delete             Delete memory            │
│  c                  Create link              │
│  x                  Delete link              │
│                                              │
│ SEARCH & FILTER                              │
│  /                  Start search             │
│  Ctrl+F             Filter by importance    │
│  Ctrl+T             Filter by tags          │
│  Ctrl+X             Clear filters           │
│                                              │
│ MODE SWITCHING                               │
│  Tab / Shift+Tab    Next/previous mode      │
│  1-4                Jump to mode number     │
│                                              │
│ OTHER                                        │
│  Cmd+P              Command palette         │
│  ?                  Toggle this help        │
│                                              │
└───────────────────────────────────────────────┘
```

#### Dashboard Mode Help

```
┌─ Dashboard Mode Help ─────────────────────────┐
│                                               │
│ METRICS & MONITORING                         │
│  ↑↓                 Scroll metrics           │
│  m                  Show/hide metrics        │
│  s                  Show/hide status         │
│  h                  Show/hide history        │
│                                              │
│ TASK MANAGEMENT                              │
│  ↑↓←→               Navigate tasks           │
│  Enter              View task details        │
│  Space              Pause/resume task        │
│  x                  Cancel task              │
│                                              │
│ CONFIGURATION                                │
│  c                  Configure dashboard      │
│  r                  Refresh all              │
│  u                  Update interval          │
│                                              │
│ EXPORT                                       │
│  Ctrl+E             Export metrics           │
│  Ctrl+S             Save snapshot            │
│                                              │
│ MODE SWITCHING                               │
│  Tab / Shift+Tab    Next/previous mode      │
│  1-4                Jump to mode number     │
│                                              │
│ OTHER                                        │
│  Cmd+P              Command palette         │
│  ?                  Toggle this help        │
│                                              │
└───────────────────────────────────────────────┘
```

### Help Text Color Scheme

```go
helpHeaderStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("81"))          // Cyan

helpKeyStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("226"))         // Yellow

helpDescStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("252"))         // Light gray

helpTipStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("243")).        // Dark gray
    Italic(true)
```

---

## Accessibility

### Color Contrast

All text must meet WCAG AA standards (minimum 4.5:1 contrast ratio).

**Safe Combinations:**
- Light text (252, 15) on dark backgrounds (235, 237, 0)
- Dark text (0, 8) on light backgrounds (15, 252)
- Bright colors (39, 196, 226) on dark backgrounds (235, 237)

**Avoid:**
- Gray on gray
- Similar brightness colors
- Red+green combinations (for colorblind users)

### Visual Indicators

Use multiple indicators beyond color alone:

- **Selected items**: Color + chevron (>) + bold
- **Errors**: Color + icon (✗) + text
- **Success**: Color + icon (✓) + text
- **Loading**: Animation + text, not just color

### Keyboard Navigation

Every mode must support:

- **Arrow keys**: Up/down for lists, left/right for panes
- **Tab/Shift+Tab**: Mode switching and focus movement
- **Esc**: Cancel/close dialogs and overlays
- **Enter**: Confirm/execute
- **Alt+Key**: Mode numbering (Alt+1 for Edit, etc.)
- **Ctrl+Key**: Common operations (Save, Search, etc.)

---

## Implementation Examples

### Example 1: Error Box with Styling

```go
package mymode

import "github.com/charmbracelet/lipgloss"

var (
    errorBoxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("196")).
        Padding(1).
        Width(60)

    errorTitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("196"))

    errorTextStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("252"))

    errorTipStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("243")).
        Italic(true)
)

func renderErrorBox(title, message, tip string) string {
    var b strings.Builder

    b.WriteString(errorTitleStyle.Render("✗ " + title))
    b.WriteString("\n")
    b.WriteString(errorTextStyle.Render(message))

    if tip != "" {
        b.WriteString("\n")
        b.WriteString(errorTipStyle.Render("Tip: " + tip))
    }

    return errorBoxStyle.Render(b.String())
}
```

### Example 2: Progress Display

```go
func renderProgress(current, total int) string {
    const width = 20
    percent := float64(current) / float64(total) * 100
    filled := int(float64(width) * float64(current) / float64(total))

    bar := "["
    bar += strings.Repeat("█", filled)
    bar += strings.Repeat("░", width-filled)
    bar += fmt.Sprintf("] %d%% (%d/%d)", int(percent), current, total)

    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("34"))

    return style.Render(bar)
}
```

### Example 3: List Item with Importance

```go
func renderMemoryListItem(mem *Memory, selected bool) string {
    var b strings.Builder

    // Selection indicator
    if selected {
        b.WriteString("> ")
    } else {
        b.WriteString("  ")
    }

    // Importance color
    importanceColors := map[int]lipgloss.Color{
        9: lipgloss.Color("196"), // Red
        8: lipgloss.Color("208"), // Orange
        7: lipgloss.Color("226"), // Yellow
        6: lipgloss.Color("118"), // Green
        5: lipgloss.Color("81"),  // Cyan
    }

    color := importanceColors[mem.Importance]
    if color == "" {
        color = lipgloss.Color("243")
    }

    style := lipgloss.NewStyle().
        Foreground(color).
        Bold(true)

    b.WriteString(style.Render("[" + strconv.Itoa(mem.Importance) + "]"))
    b.WriteString(" ")
    b.WriteString(mem.Name)

    if selected {
        // Highlight selection
        return lipgloss.NewStyle().
            Background(lipgloss.Color("237")).
            Render(b.String())
    }

    return b.String()
}
```

### Example 4: Help Overlay

```go
func renderHelpOverlay() string {
    var b strings.Builder

    b.WriteString("╭─ Keyboard Shortcuts ─────────────────────╮\n")

    sections := []struct {
        title string
        items []struct{ key, desc string }
    }{
        {
            title: "Navigation",
            items: []struct{ key, desc string }{
                {"↑↓", "Move cursor"},
                {"Ctrl+Home", "Start of buffer"},
            },
        },
        {
            title: "Editing",
            items: []struct{ key, desc string }{
                {"Ctrl+A", "Select all"},
                {"Ctrl+Z", "Undo"},
            },
        },
    }

    for _, section := range sections {
        b.WriteString(fmt.Sprintf("│ %s\n", section.title))
        for _, item := range section.items {
            b.WriteString(fmt.Sprintf("│  %-15s %s\n", item.key, item.desc))
        }
        b.WriteString("│\n")
    }

    b.WriteString("╰─────────────────────────────────────────╯\n")

    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("81")).
        Padding(1, 2).
        Render(b.String())
}
```

---

## Summary

This style guide ensures consistency across all five modes of Pedantic Raven:

- **Edit Mode**: Focus on editing with semantic analysis
- **Analyze Mode**: Triple graph and semantic analysis visualization
- **Orchestrate Mode**: Multi-agent orchestration with real-time dashboard
- **Explore Mode**: Memory workspace with list, detail, and graph views
- **Dashboard Mode**: Real-time metrics and monitoring

All implementations must follow these guidelines for:
- Color consistency using defined palette
- Typography hierarchy with semantic meaning
- Proper spacing and layout standards
- Clear error messages with actionable suggestions
- Loading states that provide feedback
- Help systems accessible via `?` key
- Accessibility standards (contrast, keyboard navigation)

When in doubt, follow the principle: **clarity before beauty, accessibility before polish**.
