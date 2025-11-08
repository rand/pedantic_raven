// Package buffer provides text buffer management with editing operations.
//
// The buffer system supports:
// - Multiple buffer instances
// - CRDT-based editing operations
// - Undo/redo with full history
// - Cursor and selection management
// - Line-based text operations
// - Dirty state tracking (unsaved changes)
package buffer

import (
	"strings"
	"time"
)

// BufferID uniquely identifies a buffer.
type BufferID string

// OpType identifies the type of buffer operation.
type OpType int

const (
	OpInsert OpType = iota
	OpDelete
	OpReplace
)

// String returns the string representation of the operation type.
func (t OpType) String() string {
	switch t {
	case OpInsert:
		return "Insert"
	case OpDelete:
		return "Delete"
	case OpReplace:
		return "Replace"
	default:
		return "Unknown"
	}
}

// Position represents a location in the buffer.
type Position struct {
	Line   int
	Column int
}

// Span represents a range in the buffer.
type Span struct {
	Start Position
	End   Position
}

// Operation represents a buffer editing operation.
//
// Operations are the fundamental unit of change in the buffer system.
// They support undo/redo and can be serialized for CRDT synchronization.
type Operation struct {
	ID        string    // Unique operation ID
	Type      OpType    // Operation type
	Position  Position  // Start position
	End       Position  // End position (for delete/replace)
	Text      string    // Text to insert/replace
	Deleted   string    // Text that was deleted (for undo)
	Timestamp time.Time // Operation timestamp
	BufferID  BufferID  // Buffer this operation belongs to
}

// Buffer represents a text buffer with editing capabilities.
//
// Buffers maintain their own content, cursor position, and edit history.
// They support standard text operations (insert, delete, replace) with
// full undo/redo capabilities.
type Buffer interface {
	// ID returns the unique identifier for this buffer.
	ID() BufferID

	// Content returns the full text content of the buffer.
	Content() string

	// Lines returns the buffer content split into lines.
	Lines() []string

	// LineCount returns the number of lines in the buffer.
	LineCount() int

	// Line returns the content of a specific line (0-indexed).
	// Returns empty string if line index is out of range.
	Line(index int) string

	// Insert adds text at the given position.
	// Returns the operation that was performed.
	Insert(pos Position, text string) (*Operation, error)

	// Delete removes text from start to end position.
	// Returns the operation that was performed.
	Delete(from, to Position) (*Operation, error)

	// Replace replaces text from start to end with new text.
	// Returns the operation that was performed.
	Replace(from, to Position, text string) (*Operation, error)

	// Undo reverses the last operation.
	// Returns true if an operation was undone, false if nothing to undo.
	Undo() bool

	// Redo reapplies the last undone operation.
	// Returns true if an operation was redone, false if nothing to redo.
	Redo() bool

	// CanUndo returns true if there are operations to undo.
	CanUndo() bool

	// CanRedo returns true if there are operations to redo.
	CanRedo() bool

	// Apply applies an operation to the buffer.
	// Used for CRDT synchronization and remote operations.
	Apply(op *Operation) error

	// Cursor returns the current cursor position.
	Cursor() Position

	// SetCursor moves the cursor to the given position.
	// The position will be clamped to valid buffer bounds.
	SetCursor(pos Position)

	// IsDirty returns true if the buffer has unsaved changes.
	IsDirty() bool

	// MarkClean marks the buffer as having no unsaved changes.
	// Typically called after saving.
	MarkClean()

	// Path returns the file path associated with this buffer (if any).
	Path() string

	// SetPath sets the file path for this buffer.
	SetPath(path string)

	// Clear removes all content from the buffer.
	Clear()
}

// SimpleBuffer is a basic in-memory implementation of Buffer.
//
// This implementation uses a simple line-based representation
// and maintains a full undo/redo history.
type SimpleBuffer struct {
	id        BufferID
	lines     []string
	cursor    Position
	dirty     bool
	path      string
	history   *History
	undoStack []*Operation
	redoStack []*Operation
}

// NewBuffer creates a new empty buffer with the given ID.
func NewBuffer(id BufferID) *SimpleBuffer {
	return &SimpleBuffer{
		id:        id,
		lines:     []string{""},
		cursor:    Position{Line: 0, Column: 0},
		dirty:     false,
		path:      "",
		history:   NewHistory(),
		undoStack: make([]*Operation, 0),
		redoStack: make([]*Operation, 0),
	}
}

// NewBufferFromContent creates a new buffer with the given content.
func NewBufferFromContent(id BufferID, content string) *SimpleBuffer {
	buf := NewBuffer(id)
	buf.SetContent(content)
	buf.MarkClean() // Initial content is considered clean
	return buf
}

// SetContent replaces the buffer content.
// This does not create an undo operation.
func (b *SimpleBuffer) SetContent(content string) {
	if content == "" {
		b.lines = []string{""}
	} else {
		b.lines = strings.Split(content, "\n")
	}
	b.cursor = Position{Line: 0, Column: 0}
	b.dirty = true
}

// ID implements Buffer.
func (b *SimpleBuffer) ID() BufferID {
	return b.id
}

// Content implements Buffer.
func (b *SimpleBuffer) Content() string {
	return strings.Join(b.lines, "\n")
}

// Lines implements Buffer.
func (b *SimpleBuffer) Lines() []string {
	// Return a copy to prevent external modification
	linesCopy := make([]string, len(b.lines))
	copy(linesCopy, b.lines)
	return linesCopy
}

// LineCount implements Buffer.
func (b *SimpleBuffer) LineCount() int {
	return len(b.lines)
}

// Line implements Buffer.
func (b *SimpleBuffer) Line(index int) string {
	if index < 0 || index >= len(b.lines) {
		return ""
	}
	return b.lines[index]
}

// Cursor implements Buffer.
func (b *SimpleBuffer) Cursor() Position {
	return b.cursor
}

// SetCursor implements Buffer.
func (b *SimpleBuffer) SetCursor(pos Position) {
	// Clamp to valid bounds
	if pos.Line < 0 {
		pos.Line = 0
	}
	if pos.Line >= len(b.lines) {
		pos.Line = len(b.lines) - 1
	}

	lineLen := len(b.lines[pos.Line])
	if pos.Column < 0 {
		pos.Column = 0
	}
	if pos.Column > lineLen {
		pos.Column = lineLen
	}

	b.cursor = pos
}

// IsDirty implements Buffer.
func (b *SimpleBuffer) IsDirty() bool {
	return b.dirty
}

// MarkClean implements Buffer.
func (b *SimpleBuffer) MarkClean() {
	b.dirty = false
}

// Path implements Buffer.
func (b *SimpleBuffer) Path() string {
	return b.path
}

// SetPath implements Buffer.
func (b *SimpleBuffer) SetPath(path string) {
	b.path = path
}

// Clear implements Buffer.
func (b *SimpleBuffer) Clear() {
	b.lines = []string{""}
	b.cursor = Position{Line: 0, Column: 0}
	b.dirty = true
	b.undoStack = make([]*Operation, 0)
	b.redoStack = make([]*Operation, 0)
}

// History manages undo/redo operations for a buffer.
type History struct {
	operations []*Operation
	position   int // Current position in history
}

// NewHistory creates a new empty history.
func NewHistory() *History {
	return &History{
		operations: make([]*Operation, 0),
		position:   -1,
	}
}

// Add adds an operation to the history.
// This truncates any redo history.
func (h *History) Add(op *Operation) {
	// Truncate redo history
	if h.position < len(h.operations)-1 {
		h.operations = h.operations[:h.position+1]
	}

	h.operations = append(h.operations, op)
	h.position = len(h.operations) - 1
}

// CanUndo returns true if there are operations to undo.
func (h *History) CanUndo() bool {
	return h.position >= 0
}

// CanRedo returns true if there are operations to redo.
func (h *History) CanRedo() bool {
	return h.position < len(h.operations)-1
}

// Undo returns the operation to undo.
func (h *History) Undo() *Operation {
	if !h.CanUndo() {
		return nil
	}

	op := h.operations[h.position]
	h.position--
	return op
}

// Redo returns the operation to redo.
func (h *History) Redo() *Operation {
	if !h.CanRedo() {
		return nil
	}

	h.position++
	return h.operations[h.position]
}

// Count returns the total number of operations in history.
func (h *History) Count() int {
	return len(h.operations)
}

// Position returns the current position in history.
func (h *History) Position() int {
	return h.position
}
