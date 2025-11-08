# Pedantic Raven - Usage Guide

## Building

```bash
# Build the application
go build -o pedantic_raven .

# Or run directly
go run .
```

## Running

```bash
# Run the compiled binary
./pedantic_raven
```

## Interface Overview

```
┌─────────────────────────────────────────────────────────┐
│ 🐦 Pedantic Raven - Phase 2                            │
│ Mode: edit | 120x40                                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────────┐  ┌──────────────────┐           │
│  │                  │  │  Context Panel   │           │
│  │     Editor       │  │                  │           │
│  │                  │  │  Entities:       │           │
│  │                  │  │    User [Person] │           │
│  │                  │  │    Doc [Thing]   │           │
│  │                  │  │                  │           │
│  │                  │  │  Relationships:  │           │
│  │                  │  │    User→creates  │           │
│  │                  │  │      →Document   │           │
│  └──────────────────┘  └──────────────────┘           │
│                                                         │
│  ┌────────────────────────────────────────┐           │
│  │  Terminal                               │           │
│  │  $ mnemosyne recall "context"           │           │
│  │                                         │           │
│  └────────────────────────────────────────┘           │
│                                                         │
├─────────────────────────────────────────────────────────┤
│ Keys: 1,2,3=modes | Tab=focus | Ctrl+K=palette |       │
│       ?=about | Ctrl+C=quit                             │
└─────────────────────────────────────────────────────────┘
```

## Keyboard Shortcuts

### Global Keys
- `Ctrl+C` - Quit application
- `Ctrl+K` - Open command palette
- `?` - Show about dialog
- `1` - Switch to Edit mode
- `2` - Switch to Explore mode (placeholder)
- `3` - Switch to Analyze mode (placeholder)
- `Tab` - Focus next pane

### Edit Mode

#### Editor Pane (when focused)
- `a-z`, `0-9`, etc. - Type characters
- `Backspace` - Delete character
- `Enter` - New line
- `Ctrl+Z` - Undo
- `Ctrl+Y` or `Ctrl+Shift+Z` - Redo
- `Ctrl+O` - Open file (shows file picker)
- `Ctrl+S` - Save file
- `Ctrl+Shift+S` - Save as (shows file picker)
- `Ctrl+F` - Search (shows search overlay)
- `Ctrl+H` - Find and replace (shows replace overlay)
- `F3` - Find next match
- `Shift+F3` - Find previous match

#### Context Panel (when focused)
- `j` or `Down` - Scroll down one line
- `k` or `Up` - Scroll up one line
- `PgDn` - Scroll down 10 lines
- `PgUp` - Scroll up 10 lines
- `Home` - Scroll to top
- `End` - Scroll to bottom
- `Tab` - Next section
- `Shift+Tab` - Previous section
- `Enter` - Toggle section expand/collapse

#### Terminal Pane (when focused)
- `a-z`, `0-9`, etc. - Type command
- `Backspace` - Delete character
- `Enter` - Execute command
- `Up` - Previous command in history
- `Down` - Next command in history
- `Left` - Move cursor left
- `Right` - Move cursor right
- `Home` - Move cursor to start
- `End` - Move cursor to end
- `PgUp` - Scroll output up
- `PgDn` - Scroll output down

## Edit Mode Features

### Automatic Semantic Analysis

When you type in the editor, semantic analysis triggers automatically (with 500ms debounce):

1. **Entities** - Detects nouns and proper nouns, classifies as:
   - Person (names, "user", "admin")
   - Place (locations, "server", "database")
   - Thing (objects, "document", "file")
   - Concept (abstract ideas, "security", "authentication")
   - Organization (companies, "GitHub", "Corp")
   - Technology (tools, "Python", "Docker")

2. **Relationships** - Extracts verb connections:
   - Subject-Predicate-Object patterns
   - Example: "User creates Document" → `User --creates--> Document`

3. **Typed Holes** - Detects placeholders:
   - `??Type` - Implementation needed (e.g., `??Function`, `??Handler`)
   - `!!constraint` - Requirement (e.g., `!!thread-safe`, `!!async`)
   - Shows priority (0-10) and complexity (0-10)

4. **Dependencies** - Finds imports/requires:
   - `import X from Y`
   - `require('module')`
   - `use crate::module`

5. **Triples** - RDF-style semantic triples:
   - Subject-Predicate-Object structures for knowledge graphs

### Terminal Commands

#### Built-in Commands (start with `:`)
- `:clear` - Clear terminal output
- `:help` - Show command help
- `:history` - Show command history
- `:exit` or `:quit` - Exit terminal

#### Mnemosyne Commands
Any command starting with `mnemosyne` will be executed:
```bash
mnemosyne recall -q "search query"
mnemosyne remember -c "content to remember"
```

#### Shell Commands
Any other command is executed as a shell command:
```bash
ls -la
git status
echo "hello world"
```

## Search and Replace

### Opening Search
Press `Ctrl+F` to open the search overlay:
- Type your search query
- Toggle search options:
  - `Ctrl+C`: Case sensitive
  - `Ctrl+W`: Whole word
  - `Ctrl+R`: Regex mode
- `Enter`: Perform search
- `F3`: Find next match
- `Shift+F3`: Find previous match
- `Esc`: Close search

### Opening Replace
Press `Ctrl+H` to open the find and replace overlay:
- Type your search query in the first field
- `Tab`: Switch to replacement field
- Type your replacement text
- `Enter`: Replace current match
- `Ctrl+A`: Replace all matches
- `F3`/`Shift+F3`: Navigate matches
- `Esc`: Close replace

### Search Options

**Case Sensitive** (`Ctrl+C`):
- Unchecked: "hello" matches "hello", "Hello", "HELLO"
- Checked: "hello" only matches "hello"

**Whole Word** (`Ctrl+W`):
- Unchecked: "hello" matches "hello" and "helloworld"
- Checked: "hello" only matches "hello" (not "helloworld")

**Regex** (`Ctrl+R`):
- Unchecked: Plain text search
- Checked: Regular expression search (e.g., `func test[0-9]+`)

### Match Navigation
- The overlay shows "Match X of Y" when search is active
- F3 navigates forward through matches (wraps around)
- Shift+F3 navigates backward through matches (wraps around)
- Cursor automatically moves to each match

### Replace Operations
- **Replace**: Replaces the current highlighted match
- **Replace All**: Replaces all matches at once
- All replace operations support undo (`Ctrl+Z`)

## Command Palette

Press `Ctrl+K` to open the command palette, then type to filter commands:

Available commands:
- `Switch to Edit Mode` (1)
- `Switch to Explore Mode` (2)
- `Switch to Analyze Mode` (3)
- `About Pedantic Raven` (?)

Use arrow keys to navigate, `Enter` to execute, `Esc` to cancel.

## Context Panel Sections

### Entities
Shows all detected entities with:
- Entity name
- Type in brackets: `[Person]`, `[Thing]`, etc.
- Occurrence count: `(3)` means appeared 3 times

### Relationships
Shows subject-predicate-object connections:
```
User → creates → Document
Admin → manages → System
```

### Typed Holes
Shows implementation placeholders with priority and complexity:
```
??Function [P:7 C:3]
??Handler [P:8 C:5]
```
- P = Priority (higher = more important)
- C = Complexity (higher = more complex)

### Dependencies
Shows external references:
```
import: react (from: react)
require: express (from: express)
```

### Triples
Shows RDF-style semantic triples:
```
(User, creates, Document)
(System, processes, Request)
```

## Tips

1. **Focus Management**: Use `Tab` to cycle through editor, context panel, and terminal
2. **Semantic Analysis**: Type naturally - the analyzer understands English text and code comments
3. **Typed Holes**: Use `??Type` to mark areas needing implementation
4. **Terminal History**: Press `Up` to recall previous commands
5. **Section Navigation**: In context panel, use `Tab` to jump between sections
6. **Search and Replace**: Use `Ctrl+F` for quick search, `Ctrl+H` for replace, `F3`/`Shift+F3` to navigate
7. **Undo/Redo**: All operations support undo (`Ctrl+Z`) and redo (`Ctrl+Y`)

## Troubleshooting

### Analysis Not Updating
- Ensure you're focused on the editor pane (border should be highlighted)
- Wait 500ms after typing (debounce period)
- Check for content - empty editor won't trigger analysis

### Terminal Not Responding
- Ensure terminal pane is focused (use `Tab` to switch)
- Check that command is typed correctly
- Built-in commands need `:` prefix

### Application Won't Start
- Check Go version: `go version` (requires Go 1.21+)
- Rebuild: `go build -o pedantic_raven .`
- Check for port conflicts if using network features

## Example Session

```
1. Start application: ./pedantic_raven
2. Focus on editor (it starts focused)
3. Type: "User creates Document and Admin reviews it"
4. Wait 500ms for analysis
5. Press Tab to focus context panel
6. See entities: User [Person], Admin [Person], Document [Thing]
7. See relationships: User → creates → Document, Admin → reviews → it
8. Press Tab to focus terminal
9. Type: :help
10. Press Enter
11. See built-in command help
12. Type: ls
13. Press Enter
14. See directory listing
15. Press Up to recall command
16. Press Ctrl+C to quit
```

## Development Mode

For development, you can run tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/editor/semantic/... -v

# Check test coverage
go test ./... -cover
```

## Architecture Notes

- **Modes**: Application supports multiple modes (Edit, Explore, Analyze)
- **Components**: Each mode has its own layout and components
- **Streaming**: Semantic analysis runs asynchronously with progress updates
- **Bubble Tea**: Built on Charm's Bubble Tea framework for TUI
- **Lipgloss**: Uses Lipgloss for styling and layout

See `docs/PHASE2_SUMMARY.md` for detailed architecture documentation.
