# Explore Mode User Guide

**Version**: Phase 9
**Status**: Complete
**Last Updated**: 2025-11-09

## Overview

Explore Mode provides a comprehensive memory workspace for browsing, searching, and navigating your mnemosyne knowledge base. It offers three coordinated views:

- **Memory List**: Browse all memories with filtering and search capabilities
- **Memory Detail**: View, edit, create, and delete individual memories with full metadata
- **Memory Graph**: Visualize relationships and connections between memories using force-directed layout

Explore Mode is ideal for:
- Building and maintaining knowledge bases
- Project documentation and reference material
- Research notes and literature reviews
- Code reference libraries with interconnected examples
- Personal journals with theme tracking across entries
- Creating semantic networks of related concepts

### Key Capabilities

- **Full CRUD Operations**: Create, read, update, and delete memories
- **Advanced Search**: Plain text and regex-based search with real-time results
- **Powerful Filtering**: Filter by tags (AND/OR logic), importance levels, and namespaces
- **Link Management**: Create, view, and navigate links between memories
- **Graph Visualization**: Interactive force-directed graph of memory relationships
- **Navigation History**: Back/forward navigation with breadcrumb trails
- **Offline Support**: Queue operations when disconnected, sync when restored
- **Multi-layout**: Switch between list+detail and full-screen graph views
- **Real-time Updates**: Live synchronization with mnemosyne server

## Getting Started

### Entering Explore Mode

There are three ways to enter Explore Mode:

1. **Keyboard shortcut**: Press `3` from any mode
2. **Command palette**: Press `Ctrl+K`, type "explore", select "Switch to Explore Mode"
3. **Mode switcher**: Use the mode switcher UI at the top of the screen

### Interface Overview

Explore Mode displays two primary layouts depending on your needs:

**Standard Layout (List + Detail)**:

```
┌────────────────────────────────────────────────────────────────┐
│ EXPLORE MODE - Memory Workspace                                 │
├────────────────────────────────────────────────────────────────┤
│ Breadcrumb: root > concept-a > detail-a1                        │
│                                                                  │
│  Memory List (40%)      │    Memory Detail (60%)                │
│                         │                                       │
│  [1] Authentication     │  ID: mem-1                           │
│  [2] Caching Layer      │  Content:                            │
│  [3] Database Schema    │  Architecture decision: Using event  │
│  [4] API Design         │  sourcing for audit trail            │
│                         │                                       │
│  1-4 of 100             │  Importance: 8 ⭐⭐⭐⭐⭐            │
│                         │  Tags: architecture, patterns        │
│  [Search] [Filters]     │  Namespace: project:myapp           │
│                         │  Created: 2024-01-01                │
│                         │  Updated: 2024-01-02                │
│                         │                                      │
│                         │  Links:                              │
│                         │    • mem-2 (references)              │
│                         │    • mem-4 (extends)                 │
├────────────────────────────────────────────────────────────────┤
│ Status: Connected | 100 memories | Filter: None                 │
└────────────────────────────────────────────────────────────────┘
```

**Graph Layout (Full Screen)**:

```
┌────────────────────────────────────────────────────────────────┐
│ EXPLORE MODE - Memory Graph                                     │
├────────────────────────────────────────────────────────────────┤
│                                                                  │
│                        [mem-1]---mem-2                          │
│                        /   |   \                                │
│                   mem-3    |    mem-4                           │
│                            |                                    │
│                         [mem-5]---mem-6                         │
│                                                                  │
│                    (Force-directed layout)                      │
│                                                                  │
├────────────────────────────────────────────────────────────────┤
│ Selected: mem-1 | Nodes: 10 | Edges: 15 | Layout: Ready        │
└────────────────────────────────────────────────────────────────┘
```

### Two Layout Modes

**Standard Layout** (Default):
- Memory List on the left (40% width)
- Memory Detail on the right (60% width)
- Best for simultaneous browsing and detailed reading
- Focus toggles between list and detail with Tab key
- Press `g` to switch to graph mode

**Graph Layout** (Full Screen):
- Full-screen force-directed visualization
- Nodes represent memories, edges represent links
- Best for understanding relationships and patterns
- Press `g` to switch back to standard layout

### Basic Navigation

- **Arrow keys** or **j/k**: Move up/down in focused component
- **Enter**: Select item (in list) or navigate to link (in detail)
- **Tab**: Switch focus between list and detail (standard layout only)
- **g**: Toggle between standard and graph layouts
- **Esc**: Close help overlay, clear error messages
- **?**: Show keyboard shortcuts help

## Keyboard Shortcuts Reference

### Navigation Shortcuts

| Key | Action |
|-----|--------|
| `↑` or `k` | Move up in list |
| `↓` or `j` | Move down in list |
| `Enter` | Select memory in list, navigate to linked memory |
| `Page Up` | Scroll list up (jump by 5 items) |
| `Page Down` | Scroll list down (jump by 5 items) |
| `Home` | Jump to top of list |
| `End` | Jump to bottom of list |
| `Alt+Left` or `Ctrl+[` | Navigate back in history |
| `Alt+Right` or `Ctrl+]` | Navigate forward in history |

### Display and Layout

| Key | Action |
|-----|--------|
| `Tab` | Toggle focus between list and detail (standard layout) |
| `Shift+Tab` | Previous pane (when applicable) |
| `g` | Toggle between standard and graph layouts |
| `m` | Toggle metadata display (in memory detail) |
| `?` | Show keyboard shortcuts help |
| `Esc` | Hide help, clear selections |

### Memory Operations

| Key | Action |
|-----|--------|
| `n` | Create new memory |
| `e` | Edit selected memory |
| `Delete` | Delete selected memory (with confirmation) |
| `c` | Create link from current memory |
| `x` | Delete link from current memory (with confirmation) |

### Search and Filtering

| Key | Action |
|-----|--------|
| `/` | Activate search mode (plain text by default) |
| `Ctrl+R` | Toggle regex mode in search (when searching) |
| `Ctrl+M` | Cycle search mode (plain, regex, fuzzy) |
| `Ctrl+F` | Filter by importance level |
| `Ctrl+T` | Filter by tags (AND/OR logic) |
| `Ctrl+X` | Clear all filters |
| `Escape` (in search) | Cancel search, return to normal mode |

### Graph Navigation

| Key | Action |
|-----|--------|
| `h` or `←` | Pan left |
| `j` or `↓` | Pan down |
| `k` or `↑` | Pan up |
| `l` or `→` | Pan right |
| `+` or `=` | Zoom in |
| `-` or `_` | Zoom out |
| `0` | Reset view (center and default zoom) |
| `Tab` | Select next node (cycles through all nodes) |
| `Enter` | Navigate to selected node's memory |
| `r` | Re-run layout algorithm (force-directed) |
| `Space` | Single step of layout algorithm |
| `c` | Center view on selected node |

### Mode Switching

| Key | Action |
|-----|--------|
| `1` | Switch to Edit Mode |
| `2` | Switch to Analyze Mode |
| `3` | Switch to Explore Mode |
| `4` | Switch to Orchestrate Mode |
| `Tab` (at top level) | Cycle to next mode |
| `Shift+Tab` (at top level) | Cycle to previous mode |

### Other Controls

| Key | Action |
|-----|--------|
| `Cmd+P` or `Ctrl+K` | Open command palette |
| `r` | Refresh/reload memories from server |
| `q` | Quit application (from top level) |

## Memory List View

The Memory List displays all memories in your knowledge base with filtering, search, and sorting capabilities.

### Browsing Memories

Navigate the list with arrow keys or vim-style j/k keys:
- Each item shows a preview of the memory
- Selected memory is highlighted
- Status bar shows current position (e.g., "1-4 of 100")

### List Preview Format

Each memory in the list shows:

```
[Priority] ID - Content Snippet
     Tags: tag1, tag2
     Importance: 8 | Updated: 2024-01-02
```

Example:
```
[1] mem-3 - Security review: JWT token validation needs improvement
     Tags: security, auth, jwt
     Importance: 9 | Updated: 2024-01-06
```

### Selection and Navigation

- **Select current memory**: Press `Enter` to view full details in the detail pane
- **Jump to top**: Press `Home` or `gg` (go to top)
- **Jump to bottom**: Press `End` or `G`
- **Page through**: Press `Page Up`/`Page Down` or `Ctrl+U`/`Ctrl+D`
- **Quick navigation**: Type number keys for quick jump (when available)

### Pagination and Scrolling

The list supports two scrolling modes:
- **Line-by-line**: Arrow keys move one item at a time
- **Page-by-page**: `Page Up`/`Page Down` or `Ctrl+U`/`Ctrl+D` jump by 5 items

For large datasets (100+ memories), pagination buttons appear at the bottom:
```
[< Previous] [1] [2] [3] [4] [5] [Next >]
```

### Sorting Options

Default sort is by **importance descending** (most important first). Press `s` to cycle through sort modes:

1. **By Importance** (descending): Most important memories first
2. **By Date Updated** (newest first): Recently modified memories first
3. **By Date Created** (newest first): Newly created memories first
4. **Alphabetically**: By content (A-Z)

Current sort mode displays in the status bar:
```
Sort: Importance ↓ | Filter: tags=api
```

## Memory Detail View

The Memory Detail view displays complete information about a selected memory and allows viewing, editing, creating, and deleting memories.

### Viewing Memory Details

When you select a memory from the list, the detail view shows:

```
┌─────────────────────────────────────────────────────────────┐
│ mem-1: Architecture Decision                                 │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ Content:                                                     │
│ Architecture decision: Using event sourcing for audit trail  │
│ This approach enables us to maintain a complete history of   │
│ all state changes in the system, which is crucial for       │
│ compliance and debugging.                                    │
│                                                               │
│ Importance: 8 ⭐⭐⭐⭐⭐⭐⭐⭐              (1-10 scale)    │
│                                                               │
│ Tags: architecture, patterns, event-sourcing                │
│                                                               │
│ Namespace: project:myapp                                    │
│                                                               │
│ Timestamps:                                                 │
│   Created: 2024-01-01 12:00:00 UTC                         │
│   Updated: 2024-01-02 15:30:45 UTC                         │
│                                                               │
│ Links (2):                                                  │
│   ├─ mem-2 (references): Performance optimization...       │
│   └─ mem-4 (extends): Database schema: Created users...    │
│                                                               │
│ [Edit] [New] [Delete] [New Link]                           │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Viewing Memory Fields

The detail view displays:

- **Memory ID**: Unique identifier assigned by mnemosyne
- **Content**: Full text (up to 10,000 characters)
- **Importance**: 1-10 scale (displayed with star icons)
- **Tags**: Comma-separated labels (max 20 tags)
- **Namespace**: Either "global" or "project:projectname"
- **Created At**: UTC timestamp (read-only)
- **Updated At**: UTC timestamp (read-only)
- **Links**: List of related memories with link type and strength

### Editing Memories

Press `e` to enter edit mode. The detail view transforms into an editor:

```
┌─────────────────────────────────────────────────────────────┐
│ EDIT: mem-1                                                  │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ Content (Tab to next field):                                │
│ ┌───────────────────────────────────────────────────────┐   │
│ │ Architecture decision: Using event sourcing...        │   │
│ │ Explains the rationale and benefits of event sourcing.│   │
│ │ Max: 10,000 characters (850/10000)                    │   │
│ └───────────────────────────────────────────────────────┘   │
│                                                               │
│ Tags (comma-separated, max 20):                             │
│ architecture, patterns, event-sourcing                       │
│                                                               │
│ Importance (1-10):                                          │
│ 8                                                            │
│                                                               │
│ Namespace (project:name or global):                         │
│ project:myapp                                               │
│                                                               │
│ [Save with Ctrl+S] [Cancel with Esc]                       │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Edit Mode Keyboard**:
- **Tab**: Move to next field
- **Shift+Tab**: Move to previous field
- **Ctrl+S**: Save changes
- **Esc**: Cancel editing (no changes saved)

**Field Validation Rules**:
- **Content**: Required, 1-10,000 characters
- **Importance**: Required, integer 1-10
- **Tags**: Optional, up to 20 tags, alphanumeric + hyphens + underscores
- **Namespace**: Required format "project:name" or "global"

Validation errors appear inline:
```
Tags: api, rest, design
❌ Error: Contains 3 tags (valid)

Tags: api, rest-design, my_tag
✓ Valid: 3 tags
```

### Creating New Memories

Press `n` to create a new memory. A form appears with empty fields:

```
┌─────────────────────────────────────────────────────────────┐
│ NEW MEMORY                                                   │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ Content (required, max 10,000 chars):                       │
│ ┌───────────────────────────────────────────────────────┐   │
│ │                                                       │   │
│ │ (0/10000)                                             │   │
│ └───────────────────────────────────────────────────────┘   │
│                                                               │
│ Tags (optional, comma-separated, max 20):                   │
│                                                               │
│ Importance (required, 1-10):                                │
│ (empty)                                                      │
│                                                               │
│ Namespace (required):                                        │
│ global (default)                                             │
│                                                               │
│ [Save with Ctrl+S] [Cancel with Esc]                       │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Default Values**:
- Content: Empty (required field)
- Tags: Empty (optional)
- Importance: 5 (default, can be changed)
- Namespace: "global" (can be changed to "project:name")

After saving, the new memory is created on the server and added to the list.

### Deleting Memories

Press `Delete` to delete the current memory. A confirmation dialog appears:

```
┌─────────────────────────────────────────────────────────────┐
│ CONFIRM DELETE                                               │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ Delete memory "mem-1: Architecture Decision"?               │
│                                                               │
│ This action cannot be undone.                               │
│                                                               │
│ Links to this memory (1) will be orphaned:                  │
│   • mem-2 → mem-1                                           │
│                                                               │
│ [Y] Delete  [N] Cancel  [?] Help                            │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

- **Y**: Confirm deletion
- **N**: Cancel deletion
- **?**: Show help for deletion

Deleted memories cannot be recovered. All links to the deleted memory are marked as orphaned.

## Creating New Memories

Creating new memories is straightforward and supports both quick and detailed entry modes.

### Quick Create

Press `n` to open the new memory form:

1. Type memory content (required)
2. Press Tab to move to tags field
3. Enter tags separated by commas (optional)
4. Press Tab to move to importance field
5. Enter importance 1-10 (required)
6. Press Tab to move to namespace field
7. Select or type namespace (required)
8. Press Ctrl+S to save

The memory is created immediately on the server.

### Detailed Create

Use the same form but with more careful consideration:

**Content Guidelines**:
- Be specific and concrete
- Include relevant context
- 50-500 characters is typical (but up to 10,000 allowed)
- Examples: "Use PostgreSQL JSONB for flexibility", "Implement rate limiting before launch"

**Importance Guidelines**:
- **1-3**: Nice-to-have, reference material
- **4-6**: Useful knowledge, moderately important
- **7-8**: Important, affects decision-making
- **9-10**: Critical, blocks work or affects compliance

**Tags Guidelines**:
- Use lowercase, hyphen-separated: `api-design`, `performance`, `security`
- Related memories should share tags for easy filtering
- Max 20 tags per memory
- Examples: `authentication`, `database`, `testing`, `refactoring`

**Namespace Guidelines**:
- **global**: Cross-project knowledge, reusable patterns
- **project:name**: Project-specific decisions, implementation details

## Search Functionality

Search allows you to find memories by content, with support for plain text and regex patterns.

### Activating Search

Press `/` to activate search mode. A search input appears at the bottom of the list:

```
┌────────────────────────────────────────────────┐
│ Memory List                                     │
├────────────────────────────────────────────────┤
│ [1] Authentication                             │
│ [2] Caching Layer                              │
│ [3] Database Schema                            │
│ [4] API Design                                 │
│                                                 │
│                                                 │
│ Search: [authentication__________]  (4 results)│
│ Mode: Plain Text | Ctrl+R: Toggle Regex        │
└────────────────────────────────────────────────┘
```

### Plain Text Search

Default search mode - finds memories containing your search text:

- **Query**: `authentication`
- **Matches**: Memories with "authentication" anywhere in content, tags, or metadata
- **Case**: Insensitive by default
- **Matching**: Substring matching (partial matches count)

Example results:
```
[1] mem-3 - Security review: JWT token validation needs improvement
     Tags: security, auth, jwt
```

Found because content contains "token" and tags contain "auth".

### Regex Search

Press `Ctrl+R` while searching to toggle regex mode:

```
Search: [/auth.*token/__________] (2 results)
Mode: Regex | Ctrl+R: Toggle Plain | Help: ?
```

**Regex Examples**:
- `auth.*token` - "auth" followed by any text, then "token"
- `^API` - Memories starting with "API"
- `\buser\b` - Word boundary for exact "user" match
- `[0-9]{3}` - Exactly three consecutive digits
- `(json|xml)` - Either "json" or "xml"
- `\b[A-Z]{2,}\b` - Acronyms (2+ uppercase letters)

**Regex Tips**:
- `.` matches any character except newline
- `.*` matches zero or more of any character
- `\b` is word boundary (useful for exact matches)
- `[abc]` matches 'a', 'b', or 'c'
- `{n,m}` matches between n and m times

### Search Results

Search results update in real-time as you type:

```
Search: [authentication___________] (4 results)
       1. mem-3: Security review: JWT token validation needs improvement
       2. mem-5: API Design: RESTful endpoints for user management
       3. mem-7: User authentication system requirements
       4. mem-12: Two-factor authentication implementation
```

### Search Navigation

- **Tab** or **↓**: Jump to next result
- **Shift+Tab** or **↑**: Jump to previous result
- **Enter**: Open current result in detail view
- **Esc**: Close search, return to full list
- **Ctrl+M**: Cycle through search modes (plain, regex, fuzzy)

### Search History

Successful searches are saved to history. Press **↑/↓** in search mode to navigate history:

```
Search: [authentication___________]
History: authentication, api, database, schema, caching
         ^ Most recent
```

### Search Performance

- Client-side search is fast for <1,000 memories
- For 1,000+ memories, search delegates to server
- Search debouncing prevents excessive requests (100ms delay)

## Filtering Memories

Filtering allows you to narrow memory list by tags, importance, and namespace. Multiple filters work together (AND logic by default).

### Active Filters Display

When filters are active, they appear as pills at the top of the memory list:

```
Memory List (with active filters)

[🏷️ api,rest] [⭐ 8-10] [🌍 project:myapp] [×]

[1] mem-1 - RESTful API design patterns
[2] mem-4 - REST authentication scheme
```

Click the `[×]` button or press `Ctrl+X` to clear all filters.

### Tag Filtering

Press `Ctrl+T` to open the tag filter menu:

```
┌──────────────────────────────────────┐
│ FILTER BY TAGS                        │
├──────────────────────────────────────┤
│                                      │
│ Available Tags (6):                  │
│ ☐ api                                │
│ ☐ authentication                     │
│ ☐ caching                            │
│ ☐ database                           │
│ ☐ performance                        │
│ ☐ security                           │
│                                      │
│ Logic: ⦿ AND  ○ OR                   │
│ (AND = must have ALL selected)       │
│ (OR = must have AT LEAST ONE)        │
│                                      │
│ [Apply] [Clear] [Cancel]             │
│                                      │
└──────────────────────────────────────┘
```

**AND Logic** (default):
- Select: `api` AND `rest`
- Matches: Memories with BOTH tags
- Result: 3 memories with both "api" and "rest"

**OR Logic**:
- Select: `api` OR `rest` or `database`
- Matches: Memories with ANY of those tags
- Result: 12 memories with at least one of the tags

### Importance Filtering

Press `Ctrl+F` to open importance filter menu:

```
┌──────────────────────────────────────┐
│ FILTER BY IMPORTANCE                  │
├──────────────────────────────────────┤
│                                      │
│ ○ All (no filter)                    │
│ ○ Critical (9-10)  - 2 memories      │
│ ○ High (8-10)      - 5 memories      │
│ ⦿ Medium (5-7)     - 18 memories     │
│ ○ Low (1-4)        - 7 memories      │
│                                      │
│ Custom Range:                        │
│ Min: [_] Max: [_]  [Apply]          │
│                                      │
│ [Apply] [Clear] [Cancel]             │
│                                      │
└──────────────────────────────────────┘
```

**Quick Filters**:
- **Critical (9-10)**: Only the most important memories
- **High (8-10)**: Important and critical memories
- **Medium (5-7)**: Moderately important (default focus area)
- **Low (1-4)**: Reference material and nice-to-haves

**Custom Range**:
- Enter min: `6`, max: `8` to filter importance 6, 7, 8
- Leave max empty for "6 and above"
- Leave min empty for "up to X"

### Namespace Filtering

Press `Ctrl+N` to filter by namespace:

```
┌──────────────────────────────────────┐
│ FILTER BY NAMESPACE                   │
├──────────────────────────────────────┤
│                                      │
│ ☐ global             - 15 memories   │
│ ☐ project:myapp      - 32 memories   │
│ ☐ project:research   - 8 memories    │
│ ☐ project:learning   - 5 memories    │
│                                      │
│ [Apply] [Clear] [Cancel]             │
│                                      │
└──────────────────────────────────────┘
```

### Combined Filtering

Filters work together with AND logic across filter types:

```
Active Filters:
  Tags: api AND rest (AND logic)
  AND Importance: 8-10
  AND Namespace: project:myapp

Result: 2 memories matching all criteria
```

Example: Find high-priority REST API memories in myapp project:
1. Press `Ctrl+T`, select "api" and "rest" with AND logic
2. Press `Ctrl+F`, select "High (8-10)"
3. Press `Ctrl+N`, select "project:myapp"
4. Result: 2 memories

### Clearing Filters

**Clear all filters**: Press `Ctrl+X`

**Clear individual filter**: Press `Ctrl+T` (tags), `Ctrl+F` (importance), or `Ctrl+N` (namespace), then click "Clear"

## Link Management

Links connect memories to create semantic relationships. Create links to document how memories relate to each other.

### Creating Links

Press `c` to create a link from the current memory:

```
┌─────────────────────────────────────────────────────────────┐
│ CREATE LINK FROM mem-1                                       │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ Target Memory (search or select):                           │
│ [_________________________]  (15 matches)                    │
│   1. mem-2: Performance optimization: Added caching...      │
│   2. mem-4: Database schema: Created users table...         │
│   3. mem-5: API design: RESTful endpoints...                │
│                                                               │
│ Link Type:                                                  │
│ ⦿ references   (general reference)                          │
│ ○ extends      (builds upon, adds detail)                   │
│ ○ contradicts  (conflicts with)                             │
│ ○ depends-on   (requires, prerequisite)                     │
│ ○ related      (loosely related)                            │
│                                                               │
│ Link Strength (0.0-1.0):                                    │
│ [████████░░] 0.85                                            │
│                                                               │
│ [Create] [Cancel]                                           │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

**Creating Steps**:
1. Search for target memory by typing memory ID or content
2. Select from suggestions
3. Choose link type (references, extends, contradicts, depends-on, related)
4. Set link strength (optional, default 0.8)
5. Press Enter or click [Create]

### Link Types

- **references**: General reference, mentions, or cites
  - Example: "Design document references implementation guide"

- **extends**: Builds upon, adds detail, or expands
  - Example: "API authentication extends JWT discussion"

- **contradicts**: Conflicts with, challenges, or opposes
  - Example: "Event sourcing contradicts CQRS in this context"

- **depends-on**: Requires, has prerequisite, or needs
  - Example: "Database design depends-on schema decisions"

- **related**: Loosely related, tangentially connected
  - Example: "Performance tips related to caching strategies"

### Link Strength

Link strength is a value 0.0-1.0 indicating relationship confidence:

- **0.0-0.3**: Weak link (tangential connection)
- **0.4-0.6**: Medium link (some relevance)
- **0.7-1.0**: Strong link (direct relationship)

In the graph, strong links appear as solid lines, weak links as dashed lines.

### Viewing Links

Links appear in the Memory Detail view:

```
Links (3):
  ├─ mem-2 (references)  Strength: 0.9
  │  Performance optimization: Added caching layer...
  ├─ mem-4 (extends)     Strength: 0.8
  │  Database schema: Created users table with indexes
  └─ mem-7 (related)     Strength: 0.5
     API design: RESTful endpoints for user management
```

Press Enter on a link to navigate to that memory.

### Deleting Links

Press `x` to delete a link from the current memory:

```
┌─────────────────────────────────────────────────────────────┐
│ CONFIRM DELETE LINK                                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│ Delete link from mem-1 to mem-2?                            │
│ Type: references | Strength: 0.9                            │
│                                                               │
│ This action cannot be undone.                               │
│                                                               │
│ [Y] Delete  [N] Cancel                                      │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

- **Y**: Confirm deletion
- **N**: Cancel deletion

## Link Navigation

Navigate between memories using their links, creating exploratory paths through your knowledge base.

### Clicking Links

In Memory Detail view, links are clickable:

```
Links (2):
  ├─ mem-2 (references)  [Click to navigate]
  │  Performance optimization: Added caching layer...
  └─ mem-4 (extends)     [Click to navigate]
     Database schema: Created users table with indexes
```

Press Enter on a link to navigate to that memory. The current location is added to navigation history.

### Navigation History

Back/forward navigation tracks your path through memories:

**Keyboard**:
- `Alt+Left` or `Ctrl+[`: Navigate back to previous memory
- `Alt+Right` or `Ctrl+]`: Navigate forward to next memory

**State**:
- Navigation history is per-session (lost when closing Explore Mode)
- Maximum history depth is 100 memories
- Cycle navigation with arrow keys on valid targets

### Breadcrumb Trail

A breadcrumb trail displays your navigation path:

```
Breadcrumb: root > concept-a > detail-a1
```

Shows the chain of memories you've navigated through. The rightmost item (bold) is your current location.

**Navigation Path Example**:
```
Start at mem-1 (root)
  ↓ [Click link to mem-2]
Navigate to mem-2
  ↓ [Click link to mem-5]
Navigate to mem-5
  ↓ [Press Alt+Left]
Back to mem-2
  ↓ [Press Alt+Right]
Forward to mem-5
```

### Breadcrumb Limitations

- Maximum 5 levels displayed (deeper paths show "...")
- Most recent navigation is always visible
- Breadcrumb is for reference only (doesn't affect navigation)

Example with deep navigation:
```
Before: a > b > c > d > e > f > g > h
After (with ellipsis): ... > e > f > g > h
```

## Graph Visualization

The Memory Graph provides an interactive force-directed visualization of your memory network.

### Switching to Graph View

Press `g` to toggle between standard layout and full-screen graph:

```
Before (Standard Layout):
┌──────────────┬──────────────┐
│ Memory List  │ Memory Detail│
└──────────────┴──────────────┘

After (Graph Layout):
┌──────────────────────────────┐
│   Force-Directed Graph        │
│   (Full screen)               │
└──────────────────────────────┘
```

### Navigation Controls

**Panning** (move the view):
- `h` or `←`: Pan left
- `j` or `↓`: Pan down
- `k` or `↑`: Pan up
- `l` or `→`: Pan right

**Zooming**:
- `+` or `=`: Zoom in (closer look)
- `-` or `_`: Zoom out (see more context)
- `0`: Reset view (center and return to default zoom)

**Node Selection**:
- `Tab`: Select next node (cycles through all nodes)
- `Shift+Tab`: Select previous node
- `Enter`: Navigate to selected node's memory
- `Esc`: Deselect current node

### Graph Elements

**Nodes** (circles representing memories):
- Blue node: Standard memory
- Gold node: High importance (8-10)
- Gray node: Low importance (1-3)
- Bold border: Currently selected node
- Node label: Memory ID (first 8 characters)

**Edges** (lines representing links):
- Solid line: Strong link (strength > 0.7)
- Dashed line: Weak link (strength ≤ 0.7)
- Arrow direction: Link direction (A → B)
- Line thickness: Proportional to link strength

**Layout**:
- Force-directed: Nodes repel each other, linked nodes attract
- Auto-layout: Layout runs in background
- Manual: Nodes can drift as forces balance

### Layout Control

**Re-running Layout**:
- `r`: Re-run force-directed layout algorithm
- `Space`: Single step of layout algorithm

Useful after:
- Adding new memories and links
- Manual node repositioning
- Layout converges to suboptimal configuration

Example session:
```
1. Navigate to graph view (g)
2. Pan around (h/j/k/l)
3. Select node (Tab)
4. Navigate to memory (Enter)
5. Create new link (c)
6. Return to graph (g)
7. Re-layout (r)
8. New link should be visible
```

### Node Interaction

**Viewing Node Details**:
- `Tab` to select node
- Selected node shows ID in status bar
- Linked nodes are highlighted

**Navigating to Node**:
- `Tab` to select
- `Enter` to view selected node's memory in detail panel
- Status shows current selection

**Expanding/Collapsing** (where applicable):
- `e`: Expand selected node (show all links)
- `x`: Collapse selected node (hide all links)
- `c`: Center view on selected node

### Graph Performance

Graph rendering performance depends on memory count:

- **0-50 memories**: Instant rendering, smooth animation
- **50-200 memories**: Fast rendering, layout runs in background
- **200-500 memories**: Good performance, slight latency
- **500+ memories**: Reduced node detail, faster rendering

For large graphs, consider:
- Using filters to show relevant subset
- Focusing on high-importance memories
- Searching for specific tags or namespaces

## Offline Mode

Explore Mode supports offline operation, automatically queuing changes when disconnected and syncing when connection is restored.

### How Offline Mode Works

Offline mode is automatic:

1. **Detection**: When mnemosyne server disconnects, offline mode activates
2. **Status**: UI shows "Offline" indicator in status bar
3. **Queueing**: Create, update, delete, and link operations are queued locally
4. **Display**: Queued changes appear immediately in UI (optimistic update)
5. **Sync**: When connection restores, queue processes automatically

### Offline Indicator

Status bar shows connection state:

```
When online:
Status: Connected | 100 memories | Mode: Standard

When offline:
Status: ⚠️ Offline (Queued: 3 operations) | 100 memories | Mode: Standard

When syncing:
Status: Syncing (2 of 3 operations...) | 100 memories | Mode: Standard
```

### Operations in Offline Mode

**Fully Supported**:
- ✓ Create new memory
- ✓ Update memory content, tags, importance, namespace
- ✓ Delete memory
- ✓ Create link
- ✓ Delete link
- ✓ Navigation (history, links)
- ✓ Filtering and search (on loaded data)

**Limited/Unsupported**:
- ✗ Load new memories from server (only previously loaded)
- ✗ Search queries (searches local cache only)
- ✗ Graph visualization (uses loaded data)
- ⚠️ Refresh/reload (will fail, waits for connection)

### Operation Queuing

Operations are queued in FIFO order:

```
Queue State:
[1] Create mem-X (pending)
[2] Update mem-1 (pending)
[3] Delete mem-3 (pending)

When connection restored:
[1] Create mem-X (success)
[2] Update mem-1 (success)
[3] Delete mem-3 (success)

All operations complete, queue empties
```

### Sync Behavior

When connection is restored:

1. **Automatic Detection**: App detects server availability
2. **Queue Processing**: Operations sync in original order
3. **Conflict Resolution**:
   - Conflicts detected (e.g., memory modified on server)
   - User is prompted for resolution
   - Options: Keep local, use server version, or merge
4. **Completion**: Status updates when all operations sync
5. **Refresh**: List and graph refresh with server data

### Manual Sync

Force sync before automatic detection:

- Press `r` to refresh/reload
- If offline, operation is queued
- If online but connection is slow, retry with `r` again

### Limitations in Offline Mode

- **Search**: Limited to memories already loaded in memory
- **Filters**: Work on loaded data (may be incomplete)
- **Graph**: Shows only loaded memories and links
- **New Data**: Cannot load additional memories from server

Workaround: Ensure you load the memories you need before disconnecting.

## Example Workflows

### Workflow 1: Building a Knowledge Base

**Goal**: Create and organize a knowledge base for API design patterns

**Steps**:

1. **Enter Explore Mode** (`3`)
2. **Create foundational concept** (`n`)
   - Content: "REST API Design: Core principles and best practices"
   - Tags: `api, rest, design`
   - Importance: 9
   - Save (`Ctrl+S`)

3. **Create related memory** (`n`)
   - Content: "Authentication: JWT implementation for REST APIs"
   - Tags: `api, auth, jwt, rest`
   - Importance: 8
   - Save (`Ctrl+S`)

4. **Create another memory** (`n`)
   - Content: "Versioning: API versioning strategies"
   - Tags: `api, versioning, design`
   - Importance: 7
   - Save (`Ctrl+S`)

5. **Link the memories**
   - Select first memory (list focus)
   - Press `Enter` to view in detail
   - Press `c` to create link
   - Search for second memory, select "depends-on"
   - Create link (`Enter`)

6. **View the graph** (`g`)
   - See three nodes connected
   - Re-layout if needed (`r`)
   - Pan/zoom to explore (`h/j/k/l`, `+/-`)

7. **Add more details**
   - Press `Tab` to focus list
   - Press `j` to move down
   - Select another memory
   - Press `e` to edit, add more content
   - Press `Ctrl+S` to save

### Workflow 2: Project Documentation

**Goal**: Document a project's architecture decisions and rationale

**Steps**:

1. **Enter Explore Mode** (`3`)

2. **Set project namespace**
   - Create memories with `Namespace: project:myapp`
   - This groups all project memories together

3. **Document decisions**
   - Create "Database Selection: PostgreSQL"
   - Create "Caching Layer: Redis"
   - Create "Authentication: OAuth2"

4. **Link decisions to rationale**
   - Each decision links to a memory explaining the rationale
   - Example: "PostgreSQL" → "Rationale: JSONB support" (link: `extends`)

5. **Track trade-offs**
   - Create memories for alternative choices
   - Link with `contradicts` for conflicting approaches
   - Importance reflects impact on project

6. **Filter by project**
   - Press `Ctrl+N` to filter namespace
   - Select `project:myapp`
   - View all project-specific knowledge

7. **Maintain in graph**
   - Press `g` to see project architecture
   - Add new decisions and link to existing ones
   - Graph shows decision dependencies

### Workflow 3: Research Notes

**Goal**: Organize research notes with cross-references

**Steps**:

1. **Create literature entries**
   - "Paper: Smith et al. - Event Sourcing Patterns"
   - Tags: `research, event-sourcing, paper`
   - Importance: 8

2. **Add findings**
   - "Finding: Event sourcing improves audit trails"
   - Links to paper with `references`
   - Tags: `findings, event-sourcing`

3. **Add synthesis notes**
   - "Analysis: Event sourcing trade-offs"
   - Links findings together
   - Tags: `synthesis, event-sourcing, analysis`

4. **Build theme clusters**
   - Filter by tag `event-sourcing` (Ctrl+T)
   - See all related research
   - Links connect papers, findings, analysis

5. **Export and present**
   - Graph view shows research network
   - High-importance findings are prominent nodes
   - Links show relationships between concepts

### Workflow 4: Code Reference Library

**Goal**: Build a library of code snippets and patterns

**Steps**:

1. **Create pattern entries**
   - "Pattern: Builder pattern for complex objects"
   - Include code example in content
   - Tags: `pattern, builder, design-patterns`
   - Importance: 7

2. **Add implementations**
   - "Implementation: Builder in Python"
   - Links to pattern with `extends`
   - Tags: `python, implementation, builder`

3. **Add use cases**
   - "Use case: Building API requests"
   - Links to pattern with `related`
   - Tags: `api, use-case, builder`

4. **Cross-reference languages**
   - Filter by `python` + `pattern` tags (Ctrl+T with AND)
   - See pattern implementations in Python
   - Links show how patterns apply across languages

5. **Discover patterns**
   - Search for common keywords (/)
   - Find related patterns
   - Graph shows pattern relationships

### Workflow 5: Personal Journal with Themes

**Goal**: Keep a personal journal with theme tracking

**Steps**:

1. **Daily entries**
   - Create memory for each day
   - Tags: `personal, 2024-11-09, reflection`
   - Namespace: `global`
   - Importance: varies by significance

2. **Identify themes**
   - As you write, tag entries with recurring themes
   - Examples: `learning`, `challenge`, `growth`, `family`

3. **Link related entries**
   - Connect entries discussing same theme
   - Links: `related` for thematic connections
   - Links: `extends` for follow-ups

4. **Track themes over time**
   - Filter by theme tag (Ctrl+T)
   - See all entries touching that theme
   - Sort by date to see evolution

5. **Reflect on patterns**
   - View theme cluster in graph
   - Which themes are most important?
   - How do themes interconnect?

6. **Annual review**
   - High importance (8-10) entries are key moments
   - Links show how moments connect
   - Graph reveals personal growth patterns

## Troubleshooting

### Common Issues and Solutions

**Memory not loading**
- **Problem**: Detail view shows "Loading..." indefinitely
- **Causes**:
  - mnemosyne server is offline
  - Network connection dropped
  - Memory ID is invalid
- **Solutions**:
  1. Check status bar: "Connected" vs "Offline"
  2. Try refreshing (press `r`)
  3. Check network connectivity
  4. Try selecting a different memory

**Search not working**
- **Problem**: Search returns no results or wrong results
- **Causes**:
  - Invalid regex pattern (in regex mode)
  - Search query too specific
  - Memories don't contain search text
- **Solutions**:
  1. Clear filters first (Ctrl+X)
  2. Try plain text search (Ctrl+R to toggle mode)
  3. Check tag filters aren't excluding results
  4. Use broader search term

**Graph not rendering**
- **Problem**: Graph view is blank or shows "Initializing..."
- **Causes**:
  - Large dataset (500+ memories)
  - Memory is low
  - Graph algorithm is still computing
- **Solutions**:
  1. Wait a moment (layout algorithm is running)
  2. Filter memories first to reduce dataset
  3. Try simpler zoom (press `0` to reset)
  4. Return to standard layout (press `g`)

**Validation errors on save**
- **Problem**: Cannot create or update memory, validation error shown
- **Common errors**:
  - Content is empty or too long (max 10,000 chars)
  - Importance not 1-10
  - Tags exceed 20 or use invalid characters
  - Namespace not "global" or "project:name" format
- **Solutions**:
  1. Read error message carefully
  2. Adjust field to meet requirements
  3. Check field character count
  4. Try different tags format (comma-separated, alphanumeric + hyphens)

**Offline operations not syncing**
- **Problem**: Queued operations not syncing when connection restored
- **Causes**:
  - Connection restored but server unreachable
  - Queue contains conflicting operations
  - Memory was deleted on server
- **Solutions**:
  1. Check status bar shows "Connected" (not "Offline")
  2. Manual refresh (press `r`)
  3. Check mnemosyne server is running
  4. Wait 5-10 seconds for auto-sync
  5. If still stuck, restart Explore Mode

**Performance issues**
- **Problem**: Slow navigation, lag in list or graph
- **Causes**:
  - Large dataset (1000+ memories)
  - Complex search/filter operations
  - Graph with many links
- **Solutions**:
  1. Apply filters to reduce dataset
  2. Clear search and filters (Ctrl+X)
  3. Use standard layout instead of graph
  4. If graph is slow, zoom out first (`-` key)
  5. Consider archiving old memories

## FAQ

**Q: How many memories can I store?**
A: Technically unlimited. Performance is excellent up to 1000 memories. Beyond that, consider filtering/searching or archiving old memories.

**Q: What's the difference between importance levels?**
A: Importance (1-10) helps prioritize knowledge. 1-3 is reference material, 4-6 is useful, 7-8 is important, 9-10 is critical. Higher importance memories appear first in sorted lists.

**Q: Can I export memories?**
A: Not yet built-in, but memories are stored in mnemosyne server which supports export. Contact mnemosyne team for export options.

**Q: How do I backup my data?**
A: Backups are managed by mnemosyne server. Contact your system administrator for backup procedures.

**Q: Can I share memories with others?**
A: Sharing depends on mnemosyne server configuration. Check with your administrator about multi-user support.

**Q: What happens if I delete a linked memory?**
A: The memory is deleted. Links to it become orphaned (one-way broken links). The other memories still exist but reference a deleted memory. Consider archiving instead of deleting.

**Q: How does search ranking work?**
A: Simple text matching - no ranking yet. Results appear in memory order from server. Exact matches count same as partial matches.

**Q: Can I customize the graph layout?**
A: Graph uses force-directed layout (physics simulation). You can re-run layout (r), pan (h/j/k/l), and zoom (+/-) but cannot customize algorithm parameters.

**Q: What's the maximum link depth?**
A: No hard limit. Navigation history shows up to 100 memories. Breadcrumb shows up to 5 levels.

**Q: How do I migrate from another system?**
A: Export from old system (if supported) and import into mnemosyne via CLI or API. Contact mnemosyne team for migration tools.

**Q: Can I use markdown in memory content?**
A: Content is plain text only currently. No markdown formatting support yet.

**Q: How often does the graph refresh?**
A: Graph reflects loaded memory data. After creating links or memories, press `r` to re-layout and show new content.

**Q: What search modes are available?**
A: Plain text (default), Regex, and Fuzzy (experimental). Cycle modes with Ctrl+M or toggle Regex with Ctrl+R.

**Q: Can I sort memories by multiple fields?**
A: Not simultaneously. Use single sort (Importance, Date, Alphabetical) or filter to narrow dataset first.

**Q: How do I recover deleted memories?**
A: Deleted memories cannot be recovered from UI. Check with mnemosyne team if backups are available.

## Best Practices

### Organization Strategies

**Use Consistent Tagging**:
- Establish tag conventions for your team
- Example: `api-design`, `db-schema`, `security-audit`
- Avoid duplicate tags with different cases
- Review tag list periodically to consolidate

**Namespace for Context**:
- `global`: Cross-project knowledge, reusable patterns
- `project:name`: Project-specific decisions and details
- Keep namespaces consistent across teams

**Importance Scoring**:
- Be honest about importance - not everything is critical
- High importance (8-10) should be uncommon
- Medium (5-7) is most common
- Low (1-4) for reference material

### Link Usage Patterns

**Strong Links** (0.7-1.0):
- Direct dependencies
- Core relationships
- Critical connections

**Weak Links** (0.3-0.6):
- Tangential references
- Loose associations
- Optional reading

**Link Types**:
- `depends-on` for prerequisites
- `extends` for elaboration
- `contradicts` for conflicts
- `references` for citations
- `related` for loose connections

### Search Optimization

**Effective Queries**:
- Use specific terms: `authentication` better than `auth`
- Combine with tags: Search "jwt" with tag filter "security"
- Regex for pattern: `auth.*token` to find specific patterns

**Search History**:
- Reuse recent searches
- Save complex queries as memory content
- Document common search patterns

### Performance Tips

**For Large Datasets** (500+ memories):
- Use filters to reduce scope
- Filter by namespace when working on project
- Apply tag filters before searching
- Avoid full graph for all memories (use filter first)

**For Graph Visualization**:
- Keep node count under 200 for smooth interaction
- Filter by high-importance memories to see key nodes
- Use tag filters to show specific relationship clusters

## Integration with Other Modes

### With Edit Mode

Complement Explore Mode by editing raw content:

1. **Explore**: Browse memories, find what to edit (Explore Mode)
2. **Edit**: Switch to Edit Mode (1), edit the memory ID directly
3. **Save**: Save in Edit Mode
4. **Return**: Back to Explore Mode (3), refresh to see changes

Workflow example:
```
[In Explore Mode]
Find memory mem-1, view it, notice typo

[Switch to Edit Mode] (press 1)
Open file for mem-1
Fix the content
Save (Ctrl+S)

[Return to Explore Mode] (press 3)
Refresh (r) to see updated memory
```

### With Analyze Mode

Understand semantic relationships through analysis:

1. **Explore**: Find relevant memories (Explore Mode)
2. **Analyze**: Switch to Analyze Mode (2) to see semantic graph
3. **Compare**: Compare entities and relationships
4. **Return**: Back to Explore Mode with insights

Workflow example:
```
[In Explore Mode]
Tag filter: architecture, patterns
See 5 key architecture memories

[Switch to Analyze Mode] (press 2)
View semantic relationships in triple graph
Identify key entities and patterns

[Return to Explore Mode] (press 3)
Create links based on semantic insights
Tag memories with discovered themes
```

### With Orchestrate Mode

Plan work based on memory knowledge:

1. **Explore**: Find and review implementation memories
2. **Plan**: Switch to Orchestrate Mode (4) to create work plan
3. **Reference**: Link plan to relevant memories via documentation
4. **Execute**: Run plan, document results back in Explore Mode

Workflow example:
```
[In Explore Mode]
Review authentication memories
Understand requirements (mem-1), dependencies (mem-2), security (mem-3)

[Switch to Orchestrate Mode] (press 4)
Create plan with tasks based on memory requirements
Reference memory IDs in task descriptions

[Execute Plan]
Complete tasks

[Return to Explore Mode] (press 3)
Create memories documenting implementation decisions
Link to original requirement memories
```

## Summary

Explore Mode is a powerful workspace for building and navigating your knowledge base. Key capabilities:

- **Three coordinated views**: List, Detail, Graph
- **Full CRUD**: Create, read, update, delete memories and links
- **Advanced search**: Plain text and regex searching
- **Smart filtering**: By tags, importance, namespace (AND/OR logic)
- **Link management**: Create, view, delete relationships
- **Graph visualization**: Interactive force-directed layout
- **Navigation**: History tracking, breadcrumbs, back/forward
- **Offline support**: Queue operations, sync when connected
- **Real-time**: Live updates with mnemosyne server

**Common Tasks**:
- Browse memories: List + search + filter
- Understand relationships: View links, explore graph
- Build knowledge: Create memories, link concepts
- Maintain knowledge: Edit, update, organize
- Review and reflect: Filter by tags/themes, track evolution

For detailed information on specific features, refer to the relevant section above or press `?` in Explore Mode for quick help.
