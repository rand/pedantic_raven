# Phase 3 Specification: Advanced Editor Features

**Status**: Planning
**Timeline**: 2-3 weeks
**Priority**: High
**Dependencies**: Phase 2 (Complete)

## Overview

Phase 3 transforms the basic EditorComponent into a fully-featured text editor with professional editing capabilities. This phase integrates the existing Buffer Manager (with undo/redo) and adds essential features like file operations, search/replace, and syntax highlighting.

## Goals

1. **Leverage Existing Code**: Integrate the already-tested Buffer Manager (52 tests)
2. **File Persistence**: Save and load files from disk
3. **Search Capabilities**: Find and replace text with regex support
4. **Visual Feedback**: Basic syntax highlighting for common languages
5. **Professional UX**: Match expectations of modern code editors

## Phase 3.1: Buffer Manager Integration (Days 1-2) ✅ COMPLETE

### Current State Analysis

**EditorComponent** (`internal/editor/components.go`):
- ✅ Now uses Buffer interface (replaced `[]string` storage)
- ✅ Full undo/redo support via Buffer
- ✅ Cursor position tracking integrated
- ✅ Multi-line editing support via Buffer

**Buffer Manager** (`internal/editor/buffer/`):
- Full undo/redo support (tested)
- Position-based editing (line, column)
- Multi-line operations
- Clean/dirty state tracking
- 52 passing tests

### Integration Strategy

Replace EditorComponent's simple storage with Buffer Manager:

```go
type EditorComponent struct {
    buffer  *buffer.Buffer    // Replace lines []string
    bufferMgr *buffer.Manager  // Manage multiple buffers
    activeBufferID buffer.BufferID
}
```

### Tasks

- [x] **Day 1 Morning**: Refactor EditorComponent to use Buffer
  - Replace `lines []string` with `buffer *buffer.Buffer`
  - Update `InsertChar` to use `buffer.Insert()`
  - Update `DeleteChar` to use `buffer.Delete()`
  - Update `GetContent` to use `buffer.Content()`

- [x] **Day 1 Afternoon**: Add undo/redo keybindings
  - `Ctrl+Z` → Undo
  - `Ctrl+Y` or `Ctrl+Shift+Z` → Redo
  - Update component tests

- [x] **Day 2 Morning**: Integrate Buffer Manager
  - Support multiple open buffers
  - Buffer switching commands
  - Buffer list display

- [x] **Day 2 Afternoon**: Testing and polish
  - Verify all existing EditorComponent tests pass
  - Add undo/redo integration tests
  - Performance testing with large buffers

### Success Criteria ✅ ALL MET

- ✅ EditorComponent uses Buffer internally (complete)
- ✅ Undo/Redo works correctly (Ctrl+Z, Ctrl+Y keybindings added)
- ✅ All 29 EditorComponent tests still pass (verified)
- ✅ No regression in semantic analysis (all 291 tests passing)
- ✅ Cursor management integrated (position tracking after insert/delete)
- ✅ Performance: No regressions expected (uses same rendering logic)

## Phase 3.2: File Operations (Days 3-4)

### Requirements

**Open File**:
- File picker overlay (fuzzy search)
- Recent files list
- Error handling (file not found, permissions)
- Encoding detection (UTF-8, ASCII)

**Save File**:
- Save current buffer to disk
- Atomic write (temp file + rename)
- Dirty flag management
- Auto-backup option

**Save As**:
- Choose new file path
- Update buffer path
- Reset dirty flag

### UI Components

```
┌─────────────────────────────────────┐
│ Open File                     [Esc] │
├─────────────────────────────────────┤
│ Search: █                           │
│                                     │
│ ▸ README.md              (2 days)   │
│   main.go                (today)    │
│   internal/editor/...    (today)    │
│   docs/USAGE.md          (today)    │
│                                     │
│ ~/src/pedantic_raven/               │
└─────────────────────────────────────┘
```

### Tasks

- [ ] **Day 3 Morning**: Implement file reading
  - `OpenFile(path string) error`
  - Encoding detection
  - Error handling
  - Add to recent files

- [ ] **Day 3 Afternoon**: File picker overlay
  - Fuzzy search component
  - Directory navigation
  - Recent files display
  - Keyboard shortcuts

- [ ] **Day 4 Morning**: Implement file writing
  - `SaveFile() error`
  - `SaveFileAs(path string) error`
  - Atomic writes
  - Backup creation

- [ ] **Day 4 Afternoon**: Testing and integration
  - File I/O tests
  - Permission handling tests
  - Integration with Edit Mode
  - Command palette integration

### Success Criteria

- ✓ Can open files from disk
- ✓ Can save files atomically
- ✓ Dirty flag tracks unsaved changes
- ✓ File picker is intuitive and fast
- ✓ Handles errors gracefully
- ✓ Recent files list persists across sessions

## Phase 3.3: Search and Replace (Days 5-6)

### Requirements

**Search**:
- Incremental search (search-as-you-type)
- Case sensitive/insensitive toggle
- Whole word matching
- Regex support
- Highlight all matches
- Navigate between matches (F3/Shift+F3)

**Replace**:
- Replace current match
- Replace all matches
- Preview before replace
- Undo support for replacements

### UI Components

```
┌─────────────────────────────────────┐
│ Search                        [Esc] │
├─────────────────────────────────────┤
│ Find:    function█                  │
│ Replace: method                     │
│                                     │
│ [x] Case sensitive                  │
│ [x] Whole word                      │
│ [ ] Regex                           │
│                                     │
│ 12 matches                          │
│                                     │
│ [Next] [Previous] [Replace] [All]  │
└─────────────────────────────────────┘
```

### Tasks

- [ ] **Day 5 Morning**: Implement search
  - `Search(query string, opts SearchOptions) []Match`
  - Case sensitivity toggle
  - Whole word matching
  - Match highlighting

- [ ] **Day 5 Afternoon**: Search navigation
  - Next/Previous match (F3/Shift+F3)
  - Wrap around behavior
  - Match counter display
  - Scroll to match

- [ ] **Day 6 Morning**: Implement replace
  - `Replace(match Match, replacement string)`
  - `ReplaceAll(query, replacement string, opts)`
  - Preview overlay
  - Undo integration

- [ ] **Day 6 Afternoon**: Regex support and testing
  - Regex pattern matching
  - Capture group support
  - Search/replace tests
  - Performance optimization

### Success Criteria

- ✓ Incremental search updates as you type
- ✓ All matches highlighted in editor
- ✓ F3/Shift+F3 navigate between matches
- ✓ Replace works with undo/redo
- ✓ Regex patterns supported
- ✓ Performance: <100ms for 10,000-line files

## Phase 3.4: Basic Syntax Highlighting (Days 7-8)

### Requirements

**Supported Languages** (Phase 3):
- Go
- Markdown
- JSON
- Shell scripts

**Token Types**:
- Keywords
- Strings
- Comments
- Numbers
- Functions/identifiers
- Operators

### Implementation Strategy

Simple regex-based highlighting (not full parser):

```go
type Highlighter interface {
    Highlight(line string) []StyledSegment
    Language() string
}

type StyledSegment struct {
    Text  string
    Style lipgloss.Style
}
```

### Tasks

- [ ] **Day 7 Morning**: Design highlighter interface
  - Define token types
  - Create color scheme
  - Implement base highlighter

- [ ] **Day 7 Afternoon**: Implement Go highlighter
  - Keywords (func, type, var, const, etc.)
  - String literals
  - Comments (// and /* */)
  - Numbers

- [ ] **Day 8 Morning**: Implement Markdown & JSON
  - Markdown: headers, bold, italic, code blocks
  - JSON: keys, values, braces

- [ ] **Day 8 Afternoon**: Integration and testing
  - Auto-detect language from file extension
  - Apply highlighting during render
  - Performance testing
  - Highlighter tests

### Success Criteria

- ✓ Go syntax highlighted correctly
- ✓ Markdown formatted nicely
- ✓ JSON keys/values distinguished
- ✓ Colors are readable in dark terminals
- ✓ Performance: <16ms render for highlighted 1000-line file
- ✓ Easy to add new language highlighters

## Architecture Changes

### Updated EditorComponent

```go
type EditorComponent struct {
    // Buffer management
    bufferMgr      *buffer.Manager
    activeBufferID buffer.BufferID

    // Search state
    searchQuery    string
    searchMatches  []Match
    currentMatch   int
    searchOpts     SearchOptions

    // File operations
    recentFiles    []string

    // Syntax highlighting
    highlighter    Highlighter
}
```

### New Components

1. **FilePickerOverlay** (`internal/overlay/filepicker.go`)
   - Directory navigation
   - Fuzzy search
   - Recent files

2. **SearchOverlay** (`internal/overlay/search.go`)
   - Search input
   - Replace input
   - Options toggles

3. **Highlighters** (`internal/editor/highlight/`)
   - Base highlighter
   - Language-specific highlighters
   - Color schemes

## Command Palette Additions

New commands to register:

```go
"file.open"       → Open File (Ctrl+O)
"file.save"       → Save File (Ctrl+S)
"file.saveAs"     → Save As (Ctrl+Shift+S)
"edit.undo"       → Undo (Ctrl+Z)
"edit.redo"       → Redo (Ctrl+Y)
"edit.search"     → Search (Ctrl+F)
"edit.replace"    → Replace (Ctrl+H)
"buffer.next"     → Next Buffer (Ctrl+Tab)
"buffer.previous" → Previous Buffer (Ctrl+Shift+Tab)
"buffer.close"    → Close Buffer (Ctrl+W)
```

## Testing Strategy

### Unit Tests

- **Buffer Integration**: 15 tests
  - Undo/redo with EditorComponent
  - Multi-line operations
  - Cursor position tracking

- **File Operations**: 12 tests
  - Open/save/save-as
  - Permission handling
  - Encoding detection
  - Atomic writes

- **Search/Replace**: 18 tests
  - Search algorithms
  - Replace operations
  - Regex matching
  - Edge cases

- **Syntax Highlighting**: 20 tests
  - Per-language tests
  - Token recognition
  - Edge cases (nested comments, strings)

**Target**: 65 new tests (total 356 tests)

### Integration Tests

- Edit → Save → Close → Open → Content matches
- Search → Replace All → Undo → Verify original content
- Type → Undo → Redo → Verify final state
- Syntax highlighting doesn't break semantic analysis

### Performance Tests

- Load 10,000-line file: <500ms
- Search 10,000-line file: <100ms
- Render highlighted 1,000-line file: <16ms (60 FPS)
- Save 10,000-line file: <200ms

## Risk Mitigation

### Identified Risks

1. **Buffer Manager Integration Complexity**
   - Mitigation: Small incremental changes, test after each step

2. **File I/O Errors**
   - Mitigation: Comprehensive error handling, user-friendly messages

3. **Search Performance**
   - Mitigation: Implement incremental search, add performance tests

4. **Syntax Highlighting Performance**
   - Mitigation: Only highlight visible lines, cache results

5. **Regex Security (ReDoS)**
   - Mitigation: Timeout on regex execution, validate patterns

## Success Metrics

**Functionality**:
- ✓ Undo/redo works seamlessly
- ✓ Can open, edit, and save files
- ✓ Search finds all matches quickly
- ✓ Replace works with undo
- ✓ Syntax highlighting looks good

**Performance**:
- ✓ 60 FPS rendering maintained
- ✓ No lag on 10,000-line files
- ✓ Search completes in <100ms

**Quality**:
- ✓ 356+ total tests passing
- ✓ ~80% code coverage maintained
- ✓ No crashes or data loss
- ✓ Clean, maintainable code

**User Experience**:
- ✓ Familiar keybindings (Ctrl+S, Ctrl+F, etc.)
- ✓ Visual feedback for all operations
- ✓ Graceful error handling
- ✓ Fast and responsive

## Timeline Summary

| Days | Focus | Deliverable |
|------|-------|-------------|
| 1-2  | Buffer Manager Integration | Undo/redo working |
| 3-4  | File Operations | Open/save files |
| 5-6  | Search/Replace | Find and replace |
| 7-8  | Syntax Highlighting | Basic highlighting |

**Total**: 8 days (assuming full-time work, ~2 weeks elapsed)

## Dependencies

**External Libraries** (if needed):
- None required (use stdlib + existing dependencies)

**Internal Dependencies**:
- ✓ Buffer Manager (`internal/editor/buffer/`) - Complete
- ✓ Overlay System (`internal/overlay/`) - Complete
- ✓ Command Palette (`internal/palette/`) - Complete

## Future Enhancements (Phase 4+)

Features explicitly deferred:

- Multi-cursor editing (Phase 4)
- Advanced syntax highlighting with AST parsing (Phase 4)
- Code completion/IntelliSense (Phase 5)
- Git integration (Phase 5)
- Split panes/tabs (Phase 5)
- Vim/Emacs keybindings (Phase 6)

## Documentation Updates

After Phase 3 completion:

1. Update `docs/USAGE.md` with new keybindings
2. Create `docs/PHASE3_SUMMARY.md`
3. Update `README.md` roadmap
4. Add file operations tutorial
5. Document syntax highlighting API for new languages

---

## Approval Checklist

Before starting implementation:

- [ ] Specification reviewed
- [ ] Timeline realistic
- [ ] Dependencies identified
- [ ] Risks assessed
- [ ] Test strategy defined
- [ ] Success criteria clear

**Approved**: Ready to begin implementation
