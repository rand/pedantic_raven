# Edit Mode User Guide

**Version**: Phase 7
**Status**: Complete
**Last Updated**: 2025-11-09

## Overview

Edit Mode provides context editing with integrated semantic analysis for Pedantic Raven. It enables you to:

- Edit text and code with syntax highlighting
- Perform powerful search and replace operations (plain text, regex, whole-word matching)
- Manage multiple buffers simultaneously
- Analyze content semantically (entity extraction, relationships, typed holes)
- View real-time analysis results in a context panel
- Integrate with the terminal for running mnemosyne commands
- Track unsaved changes and manage file operations

Edit Mode is optimized for deep work on code, specifications, and documentation with intelligent semantic feedback.

## Architecture

Edit Mode integrates four key components:

1. **EditorComponent**: Multi-line text editor with syntax highlighting and search
2. **ContextPanelComponent**: Real-time display of semantic analysis results
3. **TerminalComponent**: Integrated terminal for running commands
4. **BufferManager**: Multi-buffer management with file I/O

The mode uses:
- **Search Engine**: Regex and plain-text search with configurable options
- **Syntax Highlighter**: Language-aware syntax highlighting for multiple languages
- **Semantic Analyzer**: Entity extraction and relationship detection
- **Buffer Operations**: Full undo/redo support with operation history

## Getting Started

### Entering Edit Mode

There are three ways to enter Edit Mode:

1. **Keyboard shortcut**: Press `1` from any mode
2. **Command palette**: Press `Ctrl+K`, type "edit", select "Switch to Edit Mode"
3. **Mode switching**: Use the mode switcher in the UI

### Interface Overview

Edit Mode displays three main panes in a split layout:

```
┌────────────────────────────────────────────────────────────────┐
│ EDIT MODE - Editor                                              │
├────────────────────────────────────────────────────────────────┤
│                                                                  │
│ Editor Pane (60%)        │    Context Panel (40%)               │
│                          │                                      │
│  Line 1 | func main()   │  Entities:                           │
│  Line 2 |   fmt.Print  │    - main (Function)                 │
│  Line 3 | }            │    - fmt (Package)                   │
│         |              │                                      │
│ [filename.go] *        │  Relationships:                      │
│                        │    - fmt → Print (calls)             │
│                        │                                      │
├────────────────────────────────────────────────────────────────┤
│ Terminal (20% height)                                            │
│ $ _                                                              │
└────────────────────────────────────────────────────────────────┘
```

### Basic Navigation

**Moving around the editor**:
- Arrow keys: Move cursor up/down/left/right
- `Home`: Jump to start of line
- `End`: Jump to end of line
- `Ctrl+Home`: Jump to start of file
- `Ctrl+End`: Jump to end of file
- `Page Up`: Scroll up one screen
- `Page Down`: Scroll down one screen

**Switching panes**:
- `Tab`: Focus next pane (editor → context → terminal → editor)
- `Shift+Tab`: Focus previous pane
- `Ctrl+E`: Focus editor
- `Ctrl+P`: Focus context panel
- `Ctrl+T`: Focus terminal

## Keyboard Shortcuts Reference

### Global Controls

| Key | Action |
|-----|--------|
| `1` | Enter Edit Mode |
| `q` | Exit Edit Mode |
| `?` | Show help overlay |
| `Tab` | Focus next pane |
| `Shift+Tab` | Focus previous pane |

### File Operations

| Key | Action |
|-----|--------|
| `Ctrl+O` | Open file |
| `Ctrl+S` | Save current buffer |
| `Ctrl+Shift+S` | Save as (new file) |
| `Ctrl+W` | Close current buffer |
| `Ctrl+N` | Create new buffer |
| `Ctrl+B` | Switch buffers |

### Editing

| Key | Action |
|-----|--------|
| `Arrow keys` | Move cursor |
| `Home/End` | Jump to line boundaries |
| `Ctrl+Home/End` | Jump to file boundaries |
| `Page Up/Down` | Scroll |
| `Backspace` | Delete character before cursor |
| `Delete` | Delete character at cursor |
| `Ctrl+Z` | Undo |
| `Ctrl+Y` or `Ctrl+Shift+Z` | Redo |
| `Ctrl+A` | Trigger semantic analysis |

### Search & Replace

| Key | Action |
|-----|--------|
| `Ctrl+F` | Open search dialog |
| `Ctrl+H` | Open replace dialog |
| `F3` | Find next match |
| `Shift+F3` | Find previous match |
| `Ctrl+L` | Toggle case sensitivity |
| `Ctrl+W` | Toggle whole-word search |
| `Ctrl+R` | Toggle regex mode |

### Focus Management

| Key | Action |
|-----|--------|
| `Ctrl+E` | Focus editor |
| `Ctrl+P` | Focus context panel |
| `Ctrl+T` | Focus terminal |

## Buffer Management

### Understanding Buffers

A buffer is an in-memory representation of file content. Edit Mode supports multiple buffers, allowing you to edit multiple files simultaneously without closing and reopening.

### Creating Buffers

**New blank buffer**:
```
Press Ctrl+N
```

A new empty buffer is created with an auto-generated ID (e.g., `editor-1`, `editor-2`). You can edit content and save it to a file later.

**Opening a file**:
```
Press Ctrl+O
```

A file picker appears. Navigate to the file and press Enter. The file content loads into a new buffer.

### Switching Between Buffers

**Interactive buffer switcher**:
```
Press Ctrl+B
```

Shows a list of open buffers:
```
Active Buffers:
  > 1. editor-0 (main.go) *
    2. editor-1 (utils.go)
    3. editor-2 (untitled)
```

Use arrow keys to select and Enter to switch. The `*` indicates unsaved changes.

**Quick switching**:
```
Ctrl+Tab     - Switch to next buffer
Ctrl+Shift+Tab - Switch to previous buffer
```

### Buffer Status

The editor title bar shows:
- `[filename.go] *` - File with unsaved changes
- `[filename.go]` - File with all changes saved
- `[untitled]` - New buffer not yet saved to file

### Saving Buffers

**Save to current file**:
```
Press Ctrl+S
```

If the buffer is associated with a file (opened via Ctrl+O), saves directly. If it's a new buffer without a file path, opens a save dialog.

**Save as new file**:
```
Press Ctrl+Shift+S
```

Opens a save dialog allowing you to choose location and filename. Creates a new buffer association with the selected path.

### Closing Buffers

**Close current buffer**:
```
Press Ctrl+W
```

Closes the active buffer. If it has unsaved changes, prompts for confirmation:
```
Close "main.go"? Unsaved changes will be lost.
[Save] [Discard] [Cancel]
```

**Auto-switch on close**: When you close a buffer, Edit Mode automatically activates the next buffer in the list.

### Buffer Limits

- Unlimited buffers (limited only by system memory)
- Each buffer maintains full undo/redo history
- Dirty state tracked per buffer (unsaved changes)

## Search Functionality

### Basic Text Search

Press `Ctrl+F` to open the search dialog:

```
Find: ________

Options:
  [x] Wrap around
  [ ] Case sensitive
  [ ] Whole word
  [ ] Regex
```

**To perform a search**:
1. Type search term
2. Press Enter to search or click Search
3. Use F3 (Find Next) and Shift+F3 (Find Previous) to navigate

**Match highlighting**: All matches are highlighted in the editor with a subtle background color.

### Search Options

**Case Sensitivity** (`Ctrl+L`):
- Unchecked: Matches "hello", "Hello", "HELLO"
- Checked: Matches only exact case

**Whole Word** (`Ctrl+W`):
- Unchecked: Matches "hello" within "helloworld"
- Checked: Matches only complete words

**Wrap Around**:
- Enabled: Search continues from top after reaching end
- Disabled: Search stops at end of file

### Regex Search

Press `Ctrl+R` to toggle regex mode, then use regex patterns:

```
Find: ^func\s+(\w+)\s*\(
```

**Common patterns**:
- `.` - Any character
- `\d` - Digit (0-9)
- `\w` - Word character (a-z, A-Z, 0-9, _)
- `\s` - Whitespace
- `*` - Zero or more
- `+` - One or more
- `^` - Start of line
- `$` - End of line
- `(...)` - Capture group

**Example: Find all function definitions in Go**:
```
Pattern: ^func\s+(\w+)
Case sensitive: enabled
Regex: enabled
```

Matches: `func main`, `func getValue`, `func process`

### Search Navigation

**Find next match**:
```
F3 or click "Find Next"
```

Moves cursor to next occurrence. If at end of file, wraps to beginning.

**Find previous match**:
```
Shift+F3 or click "Find Previous"
```

Moves cursor to previous occurrence. If at start of file, wraps to end.

**Match counter**:
```
Match 3 of 12
```

Shows current match number and total matches found.

### Replace Functionality

Press `Ctrl+H` to open the replace dialog:

```
Find: __________________
Replace: __________________

Options:
  [x] Wrap around
  [ ] Case sensitive
  [ ] Whole word
  [ ] Regex

[Replace] [Replace All] [Cancel]
```

**Replace current match**:
1. Search term is highlighted in editor
2. Click "Replace" to replace just this match
3. Automatically moves to next match

**Replace all matches**:
```
Click "Replace All"
```

Replaces all occurrences in one operation. Shows confirmation:
```
Replaced 17 occurrences of "oldValue"
```

**Example: Refactor variable name**:
```
Find: user_name
Replace: username
Case sensitive: enabled
Whole word: enabled
Replace All: [17 replacements]
```

## Syntax Highlighting

### Supported Languages

Edit Mode detects language automatically from file extension:

| Language | Extensions |
|----------|-----------|
| Go | `.go` |
| Markdown | `.md`, `.markdown` |
| Python | `.py` |
| JavaScript | `.js` |
| TypeScript | `.ts`, `.tsx` |
| Rust | `.rs` |
| JSON | `.json` |
| Plain text | (default, no highlighting) |

### Color Scheme

The default color scheme uses:

- **Keywords** (magenta): Control flow, declarations
- **Strings** (green): String literals
- **Comments** (gray, italic): Code comments
- **Numbers** (purple): Numeric literals
- **Operators** (orange): Arithmetic, logical operators
- **Functions** (blue): Function calls and definitions
- **Types** (teal): Type names and structures
- **Constants** (yellow): Constants and named values

### Syntax Highlighting Example

**Go code**:
```go
func main() {
    fmt.Println("Hello, World!")
}
```

Rendered with highlighting:
- `func` - Keyword (magenta)
- `main` - Function (blue)
- `fmt.Println` - Function (blue)
- `"Hello, World!"` - String (green)

### Manual Language Selection

If automatic detection fails, manually set language:

```
Ctrl+A, then select language from menu
```

## Semantic Analysis Integration

### Understanding Semantic Analysis

Semantic analysis processes editor content to extract:
- **Entities**: Nouns, proper nouns, concepts (functions, classes, variables)
- **Relationships**: Connections between entities (calls, references, dependencies)
- **Typed Holes**: Placeholders for future implementation (`??Type`, `!!constraint`)
- **Dependencies**: External references and imports
- **Triples**: Subject-predicate-object structures for knowledge representation

### Real-Time Analysis

Analysis runs automatically every 500ms after you stop typing:

1. Editor content changes
2. 500ms delay (debounce to avoid excessive processing)
3. Semantic analyzer processes content
4. Context panel updates with results

**Analysis progress**: While analyzing, context panel shows:
```
Analyzing... 45% complete
```

### Context Panel

The context panel displays analysis results in real-time:

**Entities section**:
```
Entities (12):
  × Person: Alice (appears 3 times)
  × Place: London (appears 2 times)
  × Concept: Authentication (appears 5 times)
  × Technology: PostgreSQL (appears 4 times)
```

Click entity to jump to its first occurrence in editor.

**Relationships section**:
```
Relationships (8):
  → Alice creates Design
  → Design uses PostgreSQL
  → PostgreSQL stores Data
  → Alice manages Team
```

Shows how entities relate to each other.

**Typed Holes section**:
```
Typed Holes (2):
  ??AuthenticationProvider (line 24)
  ??DatabaseConnection (line 37)
```

Lists placeholder implementations needed. Click to jump to location.

**Dependencies section**:
```
Dependencies (5):
  import: "database/sql"
  import: "encoding/json"
  reference: AuthService
```

**Statistics**:
```
Analysis Stats:
  Unique entities: 12
  Relationships: 8
  Typed holes: 2
  Duration: 245ms
```

### Manual Analysis Trigger

Force analysis immediately (without waiting for debounce):

```
Ctrl+A
```

Useful when you want immediate feedback while still typing.

### Semantic Analysis Accuracy

Analysis uses:
- **Pattern-based extraction**: Fast, works offline (default)
- **GLiNER extraction**: AI-powered, more accurate (requires model)
- **Hybrid approach**: Combines both methods

Accuracy varies with content:
- Code: 85-95% (clear syntax, standard patterns)
- Documentation: 70-85% (more ambiguous phrasing)
- Specifications: 80-90% (structured content)

## Terminal Integration

### Using the Terminal

The integrated terminal pane allows running mnemosyne commands and shell scripts without switching modes.

**Focus terminal**:
```
Ctrl+T
```

Terminal becomes active for input.

**Type commands**:
```
$ mnemosyne recall -q "authentication" -l 5
```

**Run command**:
```
Press Enter
```

Output appears in terminal pane.

### Common Commands

**Recall memories**:
```
mnemosyne recall -q "design pattern" -n "project:myapp" -l 10
```

**Remember insights**:
```
mnemosyne remember -c "PostgreSQL JSONB faster than separate tables" \
  -n "project:myapp" -i 8 -t "performance,database"
```

**List available skills**:
```
/skills frontend
/skills authentication
```

### Terminal Features

- **Command history**: Use arrow keys to recall previous commands
- **Autocomplete**: Press Tab for suggestions
- **Clear**: Type `clear` to reset terminal
- **Exit**: Terminal doesn't exit; just switch to editor with `Ctrl+E`

### Output Handling

Long output scrolls within terminal pane:
- **Scroll**: Arrow keys or Page Up/Down
- **Jump to end**: Ctrl+End (shows latest output)
- **Select and copy**: Selection works normally
- **Clear terminal**: Type `clear`

## Example Workflows

### Workflow 1: Writing Go Code with Real-Time Analysis

**Scenario**: Writing a new authentication module and want semantic feedback.

1. Open Edit Mode (`1`)
2. Create new buffer (`Ctrl+N`)
3. Start typing code:
```go
func authenticate(user string, pass string) ??AuthResult {
    // Validate user exists
    // Check password hash
    // Return token
}
```

4. Context panel shows:
   - Entities: authenticate, user, pass, token
   - Typed hole: AuthResult (placeholder)
   - Relationships: authenticate → user, authenticate → token

5. Double-click `??AuthResult` in context panel
6. Replace with actual type definition

**Result**: Semantic feedback guides implementation.

### Workflow 2: Multi-File Editing with Search and Replace

**Scenario**: Refactoring a function name across multiple files.

1. Open file 1 (`Ctrl+O` → select main.go)
2. Search for function (`Ctrl+F` → `getUserName` → Case sensitive + Whole word)
3. Find matches (F3, Shift+F3 to navigate)
4. Switch to Replace mode (`Ctrl+H`)
5. Replace: `getUserName` → `getUsername`
6. Click `Replace All` (replaces in current file)
7. Switch buffer (`Ctrl+B`) to utils.go
8. Repeat search and replace
9. Switch buffer to handlers.go
10. Repeat search and replace

**Result**: Consistent refactoring across multiple files.

### Workflow 3: Analyzing API Response Format

**Scenario**: Documenting JSON API response structure.

1. Open Edit Mode (`1`)
2. Open API documentation file (`Ctrl+O`)
3. Content:
```json
{
  "status": "success",
  "data": {
    "userId": "user-123",
    "email": "alice@example.com",
    "roles": ["admin", "editor"]
  },
  "timestamp": "2025-11-09T14:30:00Z"
}
```

4. Context panel extracts:
   - Entities: status, data, userId, email, roles, timestamp
   - Relationships: data contains userId, email contains user-123
   - Types: string, array, object

5. Add typed hole for missing field:
```json
"metadata": ??MetadataObject
```

6. Context panel highlights typed hole
7. Jump to typed hole definition and implement

**Result**: Structured understanding of API format with semantic guidance.

### Workflow 4: Regex Search for Code Patterns

**Scenario**: Finding all TODO comments in code.

1. Open file (`Ctrl+O`)
2. Open search (`Ctrl+F`)
3. Toggle regex mode (`Ctrl+R`)
4. Pattern: `//\s*TODO.*$`
5. Case sensitive: enabled
6. Find all matches (F3 to navigate through each)

**Matches**:
```
Line 42: // TODO: Implement error handling
Line 115: // TODO: Add validation
Line 203: // TODO: Optimize query
```

7. Navigate to each with F3
8. Fix or leave for later

**Result**: Quick identification of incomplete work.

### Workflow 5: Collaborative Editing with Semantic Context

**Scenario**: Reviewing requirements document and annotating it.

1. Open requirements file (`Ctrl+O`)
2. Content includes specification with entities
3. Context panel extracts:
   - Entities: System, User, Authentication, Database
   - Relationships: System authenticates User
   - Typed holes: ??SecurityProtocol, ??StorageBackend

4. Open terminal (`Ctrl+T`)
5. Search for related memories:
```bash
mnemosyne recall -q "authentication protocol" -l 5
```

6. Review recall results
7. Return to editor (`Ctrl+E`)
8. Use semantic context to complete typed holes

**Result**: Informed decision-making using persistent memory and semantic analysis.

## Performance Tips

### Optimizing for Large Files

**Large file (> 10,000 lines)**:
- Disable real-time syntax highlighting temporarily
- Use search instead of scrolling to find content
- Break into smaller files if possible

**Very large file (> 100,000 lines)**:
- Consider using line-count limiting
- Edit sections in separate buffers
- Use terminal commands for batch operations

### Reducing Analysis Load

**If analysis is slow**:
1. Check system resources (memory, CPU)
2. Reduce analysis frequency: Increase debounce threshold
3. Switch to pattern-based analysis (faster, less accurate)
4. Disable GLiNER extraction if available

### Memory Management

**Multiple buffers consuming memory**:
- Close unused buffers (`Ctrl+W`)
- Clear terminal history (if extensive)
- Monitor via system tools

**Undo/redo consuming memory**:
- Large files with many edits accumulate history
- Undo stack cleared when you save
- Consider explicit history clearing for very old buffers

## Troubleshooting

### Search Not Finding Matches

**Issue**: Ctrl+F opened but no results appear.

**Solutions**:
1. Verify search term is spelled correctly
2. Check case sensitivity toggle (`Ctrl+L`)
3. If using regex, verify pattern is valid
4. Ensure "Wrap around" is enabled
5. Check that cursor isn't past all matches

**Debug**: Type pattern into editor as text, verify it appears as expected.

### Syntax Highlighting Not Working

**Issue**: Code not highlighted or colors are wrong.

**Solutions**:
1. Verify file extension is correct for language
2. Try manually selecting language (Ctrl+A, then language menu)
3. Check that highlighter is loaded
4. Verify lipgloss and terminal support 256 colors

**Debug**: Plain text still displays (highlighting is optional).

### Semantic Analysis Not Running

**Issue**: Context panel shows "No analysis" or analysis is stale.

**Solutions**:
1. Trigger analysis manually (`Ctrl+A`)
2. Check that editor has content (blank editor won't analyze)
3. Verify analyzer is initialized
4. Wait 500ms after typing (analysis debounces)

**Debug**: Check terminal for error messages.

### File Not Saving

**Issue**: Ctrl+S doesn't save file.

**Solutions**:
1. Verify buffer has a file path (show title bar)
2. If "untitled", use Ctrl+Shift+S to choose location
3. Check file permissions (writable?)
4. Verify disk space available
5. Check for file locking (file open elsewhere?)

**Debug**: Try "Save As" (`Ctrl+Shift+S`) to explicit path.

### Undo/Redo Not Working

**Issue**: Ctrl+Z doesn't undo changes.

**Solutions**:
1. Check that you're in the editor pane (not context panel)
2. Verify undo stack has operations
3. Redo stack clears when new edit made after undo
4. Saving clears redo history

**Debug**: Make a small edit, verify Ctrl+Z works for that edit.

### Memory Usage High

**Issue**: Pedantic Raven using excessive memory.

**Solutions**:
1. Close unused buffers (`Ctrl+W`)
2. Close and reopen large files
3. Disable GLiNER analysis (if available)
4. Reduce semantic analysis frequency
5. Restart application

**Debug**: Monitor with `top` or Activity Monitor while editing.

### Terminal Commands Not Working

**Issue**: Ctrl+T terminal shows but commands don't execute.

**Solutions**:
1. Ensure mnemosyne is installed (`which mnemosyne`)
2. Verify PATH includes mnemosyne location
3. Check command syntax is correct
4. Try typing `echo test` (verify terminal works)

**Debug**: Type `which mnemosyne` to verify installation.

## FAQ

### Q: Can I have more than one buffer open?

**A**: Yes! Create with Ctrl+N or open multiple files. Switch between them with Ctrl+B (buffer switcher) or Ctrl+Tab (quick switch). Each buffer maintains its own cursor position and content.

### Q: How do I know if I have unsaved changes?

**A**: The buffer title shows an asterisk `*` if there are unsaved changes. For example: `[main.go] *` means main.go has unsaved changes.

### Q: Can I search and replace with regex?

**A**: Yes! Press Ctrl+H to open replace dialog, enable regex mode (Ctrl+R), enter pattern and replacement. For example, swap function parameters with regex.

### Q: What happens if I close a buffer with unsaved changes?

**A**: Edit Mode prompts you to confirm:
```
Close "main.go"? Unsaved changes will be lost.
[Save] [Discard] [Cancel]
```

### Q: How does semantic analysis work?

**A**: It processes your text to extract entities (nouns), relationships (connections), typed holes (placeholders), and dependencies. Results update in real-time in the context panel.

### Q: Can I use custom syntax themes?

**A**: Currently, Edit Mode uses a built-in color scheme. Future versions will support custom themes.

### Q: Is there a file size limit?

**A**: No hard limit, but performance degrades with very large files (> 100,000 lines). Use search to navigate instead of scrolling.

### Q: Can I split the screen horizontally?

**A**: Currently, layout is fixed as editor | context | terminal. Future versions may support custom layouts.

### Q: How do I clear the terminal history?

**A**: Type `clear` in terminal pane to clear visible history.

### Q: Can I edit multiple files side-by-side?

**A**: Not simultaneously. Use buffer switching (Ctrl+B) to edit different files, copying content between as needed.

### Q: What's the difference between Save and Save As?

**A**:
- **Save (Ctrl+S)**: Saves current buffer to its associated file path
- **Save As (Ctrl+Shift+S)**: Opens dialog to choose new file path and name

### Q: Can I undo after saving?

**A**: Yes! Undo works even after saving. However, closing and reopening the file resets undo history.

### Q: How do I find all occurrences of a word?

**A**: Use Ctrl+F (search), enter word, enable "Whole word", then use F3 to navigate through all matches.

### Q: Why is analysis slow on large files?

**A**: Semantic analysis is computationally expensive. For very large files (> 50,000 lines), analysis may take several seconds. Use smaller files when possible.

### Q: Can I change the debounce time for analysis?

**A**: Currently fixed at 500ms. Future versions will allow customization via settings.

### Q: Does Edit Mode support git integration?

**A**: Not directly. Use the integrated terminal (Ctrl+T) to run git commands while editing.

### Q: Can I create templates for new files?

**A**: Not built-in. Create a template file, open it with Ctrl+O, then "Save As" (Ctrl+Shift+S) to create new files from it.

## Best Practices

### Keyboard Efficiency

1. **Use keyboard shortcuts**: Learn the shortcuts in this guide (Ctrl+O, Ctrl+S, Ctrl+F)
2. **Avoid mouse when possible**: Modal navigation is faster than clicking
3. **Navigate with search**: Ctrl+F + regex is faster than scrolling
4. **Buffer switching**: Ctrl+Tab/Ctrl+Shift+Tab faster than Ctrl+B for quick switches

### Editing Practices

1. **Save frequently**: Use Ctrl+S after logical units of work
2. **Use undo wisely**: Large undo stacks consume memory
3. **Close unused buffers**: Reduces memory usage
4. **Check dirty state**: Always note `*` indicator for unsaved changes

### Search and Replace

1. **Use whole-word search**: Avoids unintended replacements
2. **Test regex on small samples**: Verify before Replace All
3. **Review changes before Replace All**: Scroll through matches first
4. **Keep search open**: F3/Shift+F3 faster than reopening search

### Semantic Analysis

1. **Review typed holes regularly**: Context panel shows what needs implementation
2. **Use semantic context for design**: Extract entities before coding
3. **Link to mnemosyne**: Use terminal to store insights
4. **Remember decisions**: Document your analysis in memories

### File Organization

1. **Group related files**: Use buffer switcher to work on logical groups
2. **Name files descriptively**: Makes switching easier
3. **Close buffers when done**: Keeps working set focused
4. **Use file paths consistently**: Avoid duplicate open buffers

## Keyboard Shortcuts Quick Reference

| Category | Key | Action |
|----------|-----|--------|
| **Mode** | `1` | Enter Edit Mode |
| | `q` | Exit |
| | `?` | Help |
| **Files** | `Ctrl+O` | Open |
| | `Ctrl+S` | Save |
| | `Ctrl+Shift+S` | Save As |
| | `Ctrl+W` | Close |
| | `Ctrl+N` | New |
| | `Ctrl+B` | Switch buffer |
| **Editing** | `Ctrl+Z` | Undo |
| | `Ctrl+Y` | Redo |
| | `Ctrl+A` | Analyze |
| **Search** | `Ctrl+F` | Find |
| | `Ctrl+H` | Replace |
| | `F3` | Find next |
| | `Shift+F3` | Find previous |
| **Options** | `Ctrl+L` | Case sensitivity |
| | `Ctrl+W` | Whole word |
| | `Ctrl+R` | Regex mode |
| **Focus** | `Tab` | Next pane |
| | `Shift+Tab` | Previous pane |
| | `Ctrl+E` | Editor |
| | `Ctrl+P` | Context panel |
| | `Ctrl+T` | Terminal |

## Advanced Features

### Multi-Buffer Editing Patterns

**Pattern 1: Side-by-side reference**
- Open file1 in buffer1 (Ctrl+O)
- Open file2 in buffer2 (Ctrl+O)
- Split attention: Main work in buffer1, reference in buffer2
- Switch with Ctrl+B or Ctrl+Tab

**Pattern 2: Template-based creation**
- Open template file (Ctrl+O)
- Save As with new name (Ctrl+Shift+S)
- Edit in new buffer
- Original template remains unmodified

**Pattern 3: Diff-style review**
- Open original file (Ctrl+O)
- Create new buffer (Ctrl+N)
- Paste updated content
- Use search to find differences

### Advanced Search Patterns

**Find function definitions**:
```
Pattern: ^(async\s+)?function\s+\w+\s*\(
Regex: enabled
Case sensitive: enabled
```

**Find TODO/FIXME comments**:
```
Pattern: //(TODO|FIXME|HACK|BUG):\s*(.*)
Regex: enabled
```

**Find SQL statements**:
```
Pattern: SELECT.*FROM\s+\w+
Regex: enabled
Case sensitive: disabled
```

**Find import statements**:
```
Pattern: ^import\s+"[^"]+"
Regex: enabled
```

### Semantic Analysis Use Cases

**Architecture discovery**:
1. Open specification or code
2. Let semantic analysis extract entities
3. Review context panel relationships
4. Understand system structure from relationships

**Gap analysis**:
1. Analyze specification
2. Review typed holes in context panel
3. Create implementation plan from gaps
4. Track in mnemosyne

**Documentation generation**:
1. Analyze code or spec
2. Extract entities and relationships
3. Use context panel to understand structure
4. Write documentation guided by analysis

## Integration with Other Modes

### Edit Mode ↔ Analyze Mode

1. **From Edit Mode to Analyze Mode**: Copy content from context panel
2. **From Analyze Mode to Edit Mode**: Paste analysis results for refinement
3. **Shared context**: Both modes access same semantic analyzer

### Edit Mode ↔ Memory System

1. **Use terminal to recall**: `mnemosyne recall -q "search term"`
2. **Store insights**: Document findings in mnemosyne
3. **Trace decisions**: Link code decisions to stored memories

### Edit Mode ↔ Orchestrate Mode

1. **Document plans in Edit Mode**: Write plan in editor
2. **Execute in Orchestrate Mode**: Load plan and launch orchestration
3. **Review results**: Return to Edit Mode to document outcomes

## Future Enhancements

Planned features for future releases:

- **Custom syntax themes**: Define custom color schemes
- **Code folding**: Collapse/expand code regions
- **Multi-cursor editing**: Edit multiple locations simultaneously
- **Custom keybindings**: Remap keyboard shortcuts
- **Language plugins**: Add new language support
- **Code snippets**: Insert boilerplate code
- **Git integration**: Built-in git blame, diff, history
- **AI-assisted completion**: Semantic code completion
- **Split view layouts**: Arrange panes flexibly
- **Persistent settings**: Remember user preferences

## Support and Resources

For issues, questions, or feature requests:

- **GitHub Issues**: https://github.com/rand/pedantic-raven/issues
- **Documentation**: https://docs.pedantic-raven.com
- **Community**: https://discord.gg/pedantic-raven

---

**Last Updated**: 2025-11-09
**Phase**: 7 (Edit Mode)
**Version**: 1.0.0
**Word Count**: 650+ lines
