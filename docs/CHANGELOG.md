# Pedantic Raven Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [Unreleased]

### Phase 4.3: Memory UI Components (In Progress)
- Memory List Component implementation
- Memory Detail View implementation

---

## [Phase 4.1] - 2025-01-XX

### Added
- **mnemosyne RPC Client** - Complete gRPC client library
  - Connection management with configurable timeouts and retry settings
  - Health check operations and server statistics retrieval
  - Resource cleanup with Close() method
- **CRUD Operations** - Full memory lifecycle management
  - StoreMemory: Create memories with optional LLM enrichment
  - GetMemory: Retrieve single memory by ID
  - UpdateMemory: Partial updates with tag operations
  - DeleteMemory: Remove memories
  - ListMemories: Query with filters (namespace, type, tags, importance, pagination)
- **Search Operations** - Advanced memory retrieval
  - Recall: Hybrid search combining semantic + FTS + graph connectivity
  - SemanticSearch: Pure embedding-based search (supports 768d and 1536d vectors)
  - GraphTraverse: Multi-hop graph traversal from seed nodes
  - GetContext: Retrieve memories with linked context
- **Streaming Support** - Progressive result delivery
  - RecallStream: Stream search results as found
  - ListMemoriesStream: Stream memories in batches
  - StoreMemoryStream: Stream with progress updates ("enriching", "embedding", "indexing")
- **Error Handling** - Comprehensive gRPC status mapping
  - Domain-specific errors (NotFound, InvalidArgument, AlreadyExists, etc.)
  - Helper functions: IsNotFound(), IsInvalidArgument(), IsUnavailable()
  - Error wrapping with operation context
- **Namespace Helpers** - Convenient namespace constructors
  - GlobalNamespace() - Global scope
  - ProjectNamespace(name) - Project scope
  - SessionNamespace(project, sessionID) - Session scope

### Technical Details
- Added 66 comprehensive tests (22 client, 20 error, 24 memory operation tests)
- Total project tests: 461 (up from 424)
- Code added: 2,202 lines (excluding 182KB generated protobuf code)
- Protocol Buffer integration with go_package option
- Makefile targets for protobuf code generation

---

## [Phase 3] - 2025-01-XX

### Phase 3.4: Basic Syntax Highlighting (Complete)

#### Added
- **Token-Based Highlighting System** (4 files, 1,421 lines)
  - 12 token types (Keyword, String, Comment, Number, Operator, Identifier, Function, TypeName, Constant, Punctuation, Whitespace)
  - Extensible Tokenizer interface for language-specific implementations
  - DefaultStyleScheme with 256-color palette
- **Go Language Tokenizer** - Full tokenization support
  - 23 keywords (func, type, var, const, if, for, etc.)
  - 19 built-in types (int, string, bool, error, etc.)
  - Constants (true, false, nil, iota)
  - String literals (double quotes, backticks, runes)
  - Comments (// and /* */)
  - Numbers (decimal, hex 0x1F, float with exponent)
  - Function detection (identifier followed by '(')
  - Multi-char operators (:=, ==, !=, <=, >=, &&, ||)
- **Markdown Tokenizer** - Formatting support
  - Headers (# ## ###) with level detection
  - Code blocks (``` or indented)
  - Lists (- * + markers)
  - Inline code (`code`)
  - Bold (**text**) and italic (*text* or _text_)
  - Links ([text](url))
- **Language Detection** - Automatic language identification
  - By file extension (.go, .md, .py, .js, .ts, .rs, .json)
  - By filename (go.mod, go.sum, README, package.json)
  - By content (shebang detection, pattern matching)
- **EditorComponent Integration** - Automatic highlighting
  - Language detection in OpenFile()
  - Line-by-line highlighting in View()
  - Default color scheme

#### Tests
- Added 31 syntax tests (24 tokenizer + 7 integration tests)
- Total Phase 3 tests: 133 new tests
- Coverage: Token recognition, language detection, edge cases

---

### Phase 3.3: Search and Replace (Complete)

#### Added
- **Search Engine** (3 files, 1,144 lines)
  - Literal and regex pattern matching
  - Case sensitive/insensitive toggle
  - Whole word matching option
  - Multi-line search support
  - Match position tracking with wraparound
- **Replace Operations**
  - Replace current match
  - Replace all matches in single operation
  - Undo integration for all replacements
  - Position re-calculation after each replacement
- **Search Navigation**
  - NextMatch() and PreviousMatch() methods
  - Match counter display
  - ClearSearch() to reset state

#### Tests
- Added 56 tests (35 engine + 21 component integration)
- Coverage: Literal search, regex patterns, whole word, multi-line, replace operations, undo integration

---

### Phase 3.2: File Operations (Complete)

#### Added
- **File I/O Operations**
  - OpenFile(path string) - Read files with UTF-8 encoding
  - SaveFile() - Save to existing path with atomic writes
  - SaveFileAs(path string) - Save to new path and update state
  - GetFilePath() and IsDirty() methods
- **Atomic Write Pattern** - Safe file saves
  - Write to temporary file (.tmp)
  - Rename to final destination
  - Prevents data corruption on failure
- **Error Handling** - Comprehensive validation
  - File not found errors
  - Permission errors
  - Path management
  - Dirty flag tracking via Buffer integration

#### Tests
- Added 12 file operation tests
- Coverage: File reading, atomic saves, error handling, dirty flags, path management

---

### Phase 3.1: Buffer Manager Integration (Complete)

#### Added
- **Buffer Interface Integration** - Replaced `[]string` with Buffer
  - Full undo/redo support via Buffer implementation
  - Multi-line operations
  - Clean/dirty state tracking
  - Position-based editing
- **Keybindings**
  - Ctrl+Z - Undo
  - Ctrl+Y or Ctrl+Shift+Z - Redo
- **Cursor Position Management** - Track position after insert/delete operations

#### Tests
- Added 6 integration tests
- Verified all 29 existing EditorComponent tests still pass
- No regression in semantic analysis (all 291 tests passing)

---

## [Phase 2] - 2025-01-XX

### Added
- **Semantic Analyzer** (5 files, 1,620 lines) - Real-time NLP-style analysis
  - Tokenization with 7 token types
  - Entity extraction with 6 types (Person, Place, Thing, Concept, Organization, Technology)
  - Multi-word entity recognition with occurrence counting
  - Relationship detection (subject-predicate-object patterns)
  - Confidence scoring for relationships
  - Typed holes parsing (??Type and !!constraint markers)
  - Priority (0-10) and complexity (0-10) calculation
  - Implementation suggestions based on constraints
  - Dependency detection (imports, requires, references)
  - RDF-style triple generation
  - Streaming updates with cancellation support
- **Context Panel** (3 files, 1,139 lines) - Semantic results display
  - 5 sections: Entities, Relationships, Typed Holes, Dependencies, Triples
  - Rich formatting with type labels, occurrence counts, priority/complexity markers
  - Filtering by entity type or search query
  - Section expand/collapse navigation
  - Responsive auto-sizing
- **Integrated Terminal** (3 files, 1,258 lines) - Command execution
  - 3 command types: Built-in, Mnemosyne, Shell
  - Built-in commands: `:help`, `:clear`, `:history`, `:exit`
  - Line editing with full cursor management
  - Command history (100 entries) with up/down navigation
  - Scrollable output buffer (1000 line retention)
  - Execution tracking (stdout/stderr, exit codes, duration)
- **Edit Mode** (3 files, 957 lines) - Complete editing environment
  - EditorComponent: Text editing with line management
  - ContextPanelComponent: Wrapped panel with keyboard navigation
  - TerminalComponent: Wrapped terminal with command input
  - Auto-triggered semantic analysis (500ms debounce)
  - Component coordination and lifecycle management

### Tests
- Added 204 tests across all components
- Total: 291 tests passing
- Coverage: Semantic analysis (63 tests), Context panel (25 tests), Terminal (38 tests), Edit mode (29 tests), Integration (49 tests)

---

## [Phase 1] - 2025-01-11

### Added
- **PubSub Event System** (3 files, 543 lines)
  - Thread-safe event broker with non-blocking publish
  - 40+ domain event types with typed data structures
  - Global and type-specific subscriptions
  - RWMutex for thread safety
- **Multi-Pane Layout Engine** (3 files, 661 lines)
  - Composite pattern for hierarchical pane composition
  - 5 layout modes: Focus, Standard, Analysis, Compact, Custom
  - Responsive design (auto-switches to Compact for terminals <120x30)
  - Focus management with navigation (next, previous, by ID)
  - Component registry
- **Mode Registry** (2 files, 285 lines)
  - 5 application modes: Edit, Explore, Analyze, Orchestrate, Collaborate
  - Lifecycle hooks: Init(), OnEnter(), OnExit()
  - Mode switching with previous mode tracking
  - BaseMode implementation with layout engine
  - Keybinding documentation per mode
- **Overlay System** (3 files, 568 lines)
  - Modal and non-modal overlay support
  - Stack-based overlay management
  - 3 position strategies: Center, Cursor (with bounds), Custom
  - Built-in dialogs: ConfirmDialog (Yes/No), MessageDialog (OK)
  - Input blocking for modal overlays
  - Dismissal via Esc, programmatic, or click outside
- **Command Palette** (3 files, 435 lines)
  - Command registry with categories (File, Edit, View, Mode, Memory, Orchestrate, Help)
  - Fuzzy search with scoring algorithm
    - Exact name match: +100
    - Name contains query: +50
    - Description contains query: +20
    - Category match: +10
    - Subsequence match: +30
  - Navigation: Up/Down, Ctrl+P/N, Enter to execute
  - Query editing: Typing, Backspace, Ctrl+U to clear

### Tests
- Total: 87 tests passing
- Breakdown: Events (7), Layout (19), Modes (17), Overlay (24), Palette (20)
- Coverage: All foundation components

### Architecture
- Adopted Elm Architecture via Bubble Tea framework
- Established design patterns: Composite, Observer, Registry, Strategy
- Created foundational component structure
- Set up internal package organization

---

## Technical Summary

### Overall Statistics (Current)
- **Total Tests**: 763 tests (762 passing, 1 failing in memorygraph)
- **Total Code**: ~33,368 lines of Go code
- **Test Coverage**: ~64% average
- **Phases Completed**: 1, 2, 3, 4.1 (4 of 8 phases)

### Technology Stack
- **Language**: Go 1.25+
- **TUI Framework**: Bubble Tea v1.2.6 (Elm Architecture)
- **Styling**: Lipgloss v1.0.0
- **Components**: Bubbles v0.21.0
- **RPC**: gRPC + Protocol Buffers
- **Architecture**: Model-Update-View pattern with event-driven communication

### Key Architecture Patterns
1. **Elm Architecture** - Immutable state, pure functions, command-based side effects
2. **Composite Pattern** - Hierarchical pane composition (LeafPane, SplitPane)
3. **Observer Pattern** - PubSub event broker for decoupled components
4. **Registry Pattern** - Mode registry, command registry, component registry
5. **Strategy Pattern** - Position strategies, layout modes, command execution

---

## Links

- [Project Repository](https://github.com/rand/pedantic_raven)
- [mnemosyne Memory System](https://github.com/rand/mnemosyne)
- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [Roadmap](../ROADMAP.md)
- [Specification](../spec.md)

---

**Changelog Maintained By**: Development Team
**Last Updated**: 2025-11-08
