# Phase 2: Edit Mode - Specification

**Duration**: 3-4 weeks
**Goal**: Rich context editing with semantic analysis (ICS parity with streaming improvements)

---

## Overview

Phase 2 implements Edit Mode - the primary interface for creating and editing context documents with real-time semantic analysis. This brings the core ICS capabilities into the new Pedantic Raven architecture with improved streaming and interaction patterns.

## Objectives

1. **Multi-buffer editing** with CRDT-based synchronization
2. **Streaming semantic analysis** with real-time triple extraction
3. **Context panel** showing notes, triples, and dependencies
4. **Terminal integration** for command execution

---

## Components

### 1. Buffer Manager (`internal/editor/buffer/`)

**Purpose**: Manage multiple text buffers with editing operations and synchronization

**Features**:
- Multiple buffer support
- CRDT-based editing operations
- Undo/redo with history
- Cursor and selection management
- Line-based operations
- Dirty tracking (unsaved changes)

**Key Types**:
```go
type Buffer interface {
    ID() BufferID
    Content() string
    Lines() []string
    Insert(pos Position, text string) Operation
    Delete(from, to Position) Operation
    Undo() Operation
    Redo() Operation
    Apply(op Operation) error
    Cursor() Position
    SetCursor(pos Position)
    IsDirty() bool
    Save() error
}

type BufferManager struct {
    buffers  map[BufferID]Buffer
    active   BufferID
    history  *History
}
```

**Operations**:
- Create/Open/Close buffers
- Switch active buffer
- Edit operations (insert, delete, replace)
- Undo/redo across buffers
- Save/load from disk

**Tests**:
- Buffer creation and content management
- Insert/delete operations
- Undo/redo functionality
- Multi-buffer switching
- Dirty state tracking
- CRDT consistency

---

### 2. Semantic Analyzer (`internal/editor/semantic/`)

**Purpose**: Stream-based semantic analysis with triple extraction

**Features**:
- Streaming analysis (progressive results)
- Entity extraction
- Relationship mapping
- Typed hole detection
- Dependency analysis
- Real-time updates

**Key Types**:
```go
type Analyzer interface {
    Analyze(content string) <-chan AnalysisUpdate
    Stop()
    Results() *Analysis
}

type Analysis struct {
    Entities     []Entity
    Relationships []Relationship
    TypedHoles   []TypedHole
    Dependencies []Dependency
    Triples      []Triple
}

type AnalysisUpdate struct {
    Type     UpdateType  // Incremental, Complete, Error
    Progress float32
    Data     interface{}
    Error    error
}
```

**Analysis Pipeline**:
1. Tokenization
2. Entity extraction (nouns, proper nouns)
3. Relationship detection (verb phrases)
4. Typed hole identification (`??Type`, `!!constraint`)
5. Dependency mapping
6. Triple generation (subject-predicate-object)

**Tests**:
- Entity extraction accuracy
- Relationship detection
- Typed hole parsing
- Streaming updates
- Progress reporting
- Error handling

---

### 3. Context Panel (`internal/editor/context/`)

**Purpose**: Display semantic analysis results and memory context

**Features**:
- Memory notes list
- Triple viewer with filtering
- Dependency graph (text-based for now)
- Search and navigation
- Expand/collapse sections

**Key Types**:
```go
type ContextPanel struct {
    analysis  *Analysis
    notes     []MemoryNote
    filter    Filter
    expanded  map[string]bool
    selected  int
}

type Filter struct {
    EntityType   string
    Relationship string
    SearchQuery  string
}
```

**Sections**:
- **Entities**: Extracted entities with counts
- **Relationships**: Detected relationships
- **Typed Holes**: Holes needing implementation
- **Dependencies**: Dependency tree
- **Triples**: Subject-predicate-object triples
- **Memory Notes**: Related notes from mnemosyne (Phase 3)

**Tests**:
- Panel rendering
- Filtering logic
- Navigation (up/down, expand/collapse)
- Selection handling
- Search functionality

---

### 4. Terminal Component (`internal/editor/terminal/`)

**Purpose**: Embedded terminal for command execution

**Features**:
- Command input with history
- Output streaming
- Scrollback buffer
- Command execution (shell, mnemosyne CLI)
- Output parsing and highlighting

**Key Types**:
```go
type Terminal struct {
    history   []string
    current   string
    output    []string
    scrollPos int
    maxLines  int
}

type TerminalCommand struct {
    Command string
    Output  string
    Error   error
    ExitCode int
}
```

**Commands**:
- Shell commands (`ls`, `git status`, etc.)
- mnemosyne CLI (`mnemosyne recall`, `mnemosyne remember`)
- Built-in commands (`:clear`, `:help`)

**Tests**:
- Command input and history
- Output capture
- Scrolling behavior
- Command execution
- Error handling

---

### 5. Edit Mode Integration (`internal/modes/edit/`)

**Purpose**: Integrate all components into Edit mode

**Features**:
- Layout: Editor (60%) + Context Panel (20%) + Terminal (20%)
- Key bindings: Insert/Normal modes (vim-like)
- Focus management between panes
- Real-time semantic analysis on edits
- Command palette commands

**Key Types**:
```go
type EditMode struct {
    *modes.BaseMode
    bufferManager   *buffer.BufferManager
    analyzer        *semantic.Analyzer
    contextPanel    *context.ContextPanel
    terminal        *terminal.Terminal
    mode            EditorMode  // Normal, Insert
}

type EditorMode int
const (
    NormalMode EditorMode = iota
    InsertMode
)
```

**Key Bindings**:
- `i`: Enter insert mode
- `Esc`: Exit insert mode
- `h/j/k/l`: Navigation (normal mode)
- `x/dd`: Delete (normal mode)
- `u`: Undo
- `Ctrl+R`: Redo
- `Ctrl+S`: Save
- `Ctrl+P`: Command palette
- `:`: Terminal command

---

## Implementation Plan

### Week 1: Buffer Manager
**Days 1-2**: Buffer interface and implementation
- Create Buffer type with content management
- Implement insert/delete operations
- Add cursor position tracking

**Days 3-4**: BufferManager and multi-buffer support
- Buffer registry and switching
- Active buffer management
- Buffer lifecycle (create, open, close)

**Days 5-7**: Undo/redo and history
- History data structure
- Undo/redo operations
- CRDT consistency
- Tests (aim for 15-20 tests)

### Week 2: Semantic Analysis
**Days 1-2**: Analyzer interface and tokenization
- Streaming analyzer framework
- Tokenization pipeline
- Progress reporting

**Days 3-4**: Entity and relationship extraction
- Entity detection (nouns, proper nouns)
- Relationship detection (verb phrases)
- Triple generation

**Days 5-7**: Typed holes and dependencies
- Typed hole parsing (`??Type`, `!!constraint`)
- Dependency analysis
- Tests (aim for 15-20 tests)

### Week 3: Context Panel and Terminal
**Days 1-3**: Context panel
- Panel layout and rendering
- Section display (entities, relationships, etc.)
- Filtering and search
- Tests (aim for 10-15 tests)

**Days 4-7**: Terminal integration
- Terminal component
- Command execution
- Output streaming
- History management
- Tests (aim for 10-15 tests)

### Week 4: Edit Mode Integration and Polish
**Days 1-3**: Edit mode
- Integrate all components
- Layout and focus management
- Key bindings (normal/insert modes)
- Tests (aim for 10-15 tests)

**Days 4-5**: Integration testing
- End-to-end editing workflow
- Performance testing
- Bug fixes

**Days 6-7**: Documentation and polish
- Update README and docs
- Code comments and examples
- Performance optimization

---

## Testing Strategy

### Unit Tests
- Each component tested in isolation
- Mock dependencies
- Cover edge cases
- Aim for 70%+ coverage

### Integration Tests
- Buffer + Analyzer integration
- Context panel + Analysis results
- Terminal + Command execution
- Edit mode + All components

### End-to-End Tests
- Full editing workflow
- Save/load buffers
- Semantic analysis pipeline
- Command execution

**Target**: 60-80 tests total

---

## Data Structures

### Buffer Operations (CRDT)
```go
type Operation struct {
    ID        OperationID
    Type      OpType  // Insert, Delete
    Position  Position
    Text      string
    Timestamp time.Time
    BufferID  BufferID
}

type Position struct {
    Line   int
    Column int
}
```

### Semantic Analysis Results
```go
type Entity struct {
    Text     string
    Type     EntityType  // Person, Place, Thing, Concept
    Span     Span
    Count    int
}

type Relationship struct {
    Subject   Entity
    Predicate string
    Object    Entity
    Span      Span
}

type TypedHole struct {
    Type       string
    Constraint string
    Span       Span
}

type Triple struct {
    Subject   string
    Predicate string
    Object    string
}
```

---

## Performance Targets

- **Buffer operations**: <1ms per operation
- **Semantic analysis**: <500ms for 1000 lines
- **Streaming updates**: 10-30 updates/sec
- **Render time**: <16ms (60 FPS)
- **Memory**: <50MB for 10 open buffers

---

## Dependencies

### New Dependencies
```bash
# No new external dependencies required
# All components built with stdlib and existing deps
```

### Bubble Tea Components
- Use existing layout engine
- Use existing overlay system
- Integrate with event broker

---

## Exit Criteria

Phase 2 is complete when:

- [ ] Buffer manager with multi-buffer support working
- [ ] Undo/redo functionality implemented
- [ ] Streaming semantic analysis running
- [ ] Entity and triple extraction accurate
- [ ] Typed hole detection working
- [ ] Context panel displaying all sections
- [ ] Terminal component executing commands
- [ ] Edit mode integrating all components
- [ ] All tests passing (60-80 tests)
- [ ] Performance targets met
- [ ] Documentation updated

---

## Risks and Mitigations

### Risk 1: CRDT Complexity
**Mitigation**: Start with simple last-write-wins for Phase 2, add full CRDT in later phases

### Risk 2: Semantic Analysis Accuracy
**Mitigation**: Start with simple regex-based extraction, improve with NLP in later phases

### Risk 3: Performance Issues
**Mitigation**: Profile early, optimize hot paths, use streaming/chunking

### Risk 4: Integration Complexity
**Mitigation**: Test components independently first, integrate incrementally

---

## Success Metrics

- **Functionality**: All features working as specified
- **Performance**: Meets or exceeds targets
- **Quality**: 70%+ test coverage
- **Usability**: Smooth editing experience
- **Documentation**: Comprehensive docs and examples

---

## Future Enhancements (Phase 3+)

- Mnemosyne RPC integration for memory notes
- Syntax highlighting
- Code completion
- LSP integration
- Collaborative editing (Phase 7)
- Advanced NLP for semantic analysis
- Full CRDT with conflict resolution

---

**Next**: Begin Week 1 - Buffer Manager implementation
