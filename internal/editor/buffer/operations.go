package buffer

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Insert implements Buffer.
func (b *SimpleBuffer) Insert(pos Position, text string) (*Operation, error) {
	// Validate position
	if pos.Line < 0 || pos.Line >= len(b.lines) {
		return nil, fmt.Errorf("line %d out of range [0, %d)", pos.Line, len(b.lines))
	}

	line := b.lines[pos.Line]
	if pos.Column < 0 || pos.Column > len(line) {
		return nil, fmt.Errorf("column %d out of range [0, %d]", pos.Column, len(line))
	}

	// Create operation
	op := &Operation{
		ID:        uuid.New().String(),
		Type:      OpInsert,
		Position:  pos,
		Text:      text,
		Timestamp: time.Now(),
		BufferID:  b.id,
	}

	// Apply the operation
	if err := b.applyInsert(op); err != nil {
		return nil, err
	}

	// Add to history
	b.history.Add(op)
	b.dirty = true

	return op, nil
}

// Delete implements Buffer.
func (b *SimpleBuffer) Delete(from, to Position) (*Operation, error) {
	// Validate positions
	if from.Line < 0 || from.Line >= len(b.lines) {
		return nil, fmt.Errorf("from line %d out of range", from.Line)
	}
	if to.Line < 0 || to.Line >= len(b.lines) {
		return nil, fmt.Errorf("to line %d out of range", to.Line)
	}

	// Ensure from <= to
	if from.Line > to.Line || (from.Line == to.Line && from.Column > to.Column) {
		from, to = to, from
	}

	// Extract deleted text
	deleted := b.extractText(from, to)

	// Create operation
	op := &Operation{
		ID:        uuid.New().String(),
		Type:      OpDelete,
		Position:  from,
		End:       to,
		Deleted:   deleted,
		Timestamp: time.Now(),
		BufferID:  b.id,
	}

	// Apply the operation
	if err := b.applyDelete(op); err != nil {
		return nil, err
	}

	// Add to history
	b.history.Add(op)
	b.dirty = true

	return op, nil
}

// Replace implements Buffer.
func (b *SimpleBuffer) Replace(from, to Position, text string) (*Operation, error) {
	// Validate positions
	if from.Line < 0 || from.Line >= len(b.lines) {
		return nil, fmt.Errorf("from line %d out of range", from.Line)
	}
	if to.Line < 0 || to.Line >= len(b.lines) {
		return nil, fmt.Errorf("to line %d out of range", to.Line)
	}

	// Ensure from <= to
	if from.Line > to.Line || (from.Line == to.Line && from.Column > to.Column) {
		from, to = to, from
	}

	// Extract deleted text
	deleted := b.extractText(from, to)

	// Create operation
	op := &Operation{
		ID:        uuid.New().String(),
		Type:      OpReplace,
		Position:  from,
		End:       to,
		Text:      text,
		Deleted:   deleted,
		Timestamp: time.Now(),
		BufferID:  b.id,
	}

	// Apply the operation
	if err := b.applyReplace(op); err != nil {
		return nil, err
	}

	// Add to history
	b.history.Add(op)
	b.dirty = true

	return op, nil
}

// Apply implements Buffer.
func (b *SimpleBuffer) Apply(op *Operation) error {
	switch op.Type {
	case OpInsert:
		return b.applyInsert(op)
	case OpDelete:
		return b.applyDelete(op)
	case OpReplace:
		return b.applyReplace(op)
	default:
		return fmt.Errorf("unknown operation type: %v", op.Type)
	}
}

// Undo implements Buffer.
func (b *SimpleBuffer) Undo() bool {
	if !b.history.CanUndo() {
		return false
	}

	op := b.history.Undo()
	if op == nil {
		return false
	}

	// Reverse the operation
	switch op.Type {
	case OpInsert:
		// Undo insert by deleting the inserted text
		end := b.calculateEndPosition(op.Position, op.Text)
		b.deleteText(op.Position, end)

	case OpDelete:
		// Undo delete by reinserting the deleted text
		b.insertText(op.Position, op.Deleted)

	case OpReplace:
		// Undo replace by deleting new text and reinserting old text
		end := b.calculateEndPosition(op.Position, op.Text)
		b.deleteText(op.Position, end)
		b.insertText(op.Position, op.Deleted)
	}

	b.dirty = true
	return true
}

// Redo implements Buffer.
func (b *SimpleBuffer) Redo() bool {
	if !b.history.CanRedo() {
		return false
	}

	op := b.history.Redo()
	if op == nil {
		return false
	}

	// Reapply the operation
	switch op.Type {
	case OpInsert:
		b.insertText(op.Position, op.Text)

	case OpDelete:
		b.deleteText(op.Position, op.End)

	case OpReplace:
		b.deleteText(op.Position, op.End)
		b.insertText(op.Position, op.Text)
	}

	b.dirty = true
	return true
}

// CanUndo implements Buffer.
func (b *SimpleBuffer) CanUndo() bool {
	return b.history.CanUndo()
}

// CanRedo implements Buffer.
func (b *SimpleBuffer) CanRedo() bool {
	return b.history.CanRedo()
}

// --- Internal helper methods ---

// applyInsert applies an insert operation to the buffer.
func (b *SimpleBuffer) applyInsert(op *Operation) error {
	return b.insertText(op.Position, op.Text)
}

// applyDelete applies a delete operation to the buffer.
func (b *SimpleBuffer) applyDelete(op *Operation) error {
	return b.deleteText(op.Position, op.End)
}

// applyReplace applies a replace operation to the buffer.
func (b *SimpleBuffer) applyReplace(op *Operation) error {
	// Delete old text
	if err := b.deleteText(op.Position, op.End); err != nil {
		return err
	}

	// Insert new text
	return b.insertText(op.Position, op.Text)
}

// insertText inserts text at the given position.
func (b *SimpleBuffer) insertText(pos Position, text string) error {
	if pos.Line < 0 || pos.Line >= len(b.lines) {
		return fmt.Errorf("line %d out of range", pos.Line)
	}

	line := b.lines[pos.Line]
	if pos.Column < 0 || pos.Column > len(line) {
		return fmt.Errorf("column %d out of range", pos.Column)
	}

	// Handle multi-line inserts
	if strings.Contains(text, "\n") {
		lines := strings.Split(text, "\n")

		// First line: insert at position
		firstLine := line[:pos.Column] + lines[0]

		// Last line: append remainder of original line
		lastLine := lines[len(lines)-1] + line[pos.Column:]

		// Middle lines: use as-is
		newLines := make([]string, 0, len(lines))
		newLines = append(newLines, firstLine)
		if len(lines) > 2 {
			newLines = append(newLines, lines[1:len(lines)-1]...)
		}
		newLines = append(newLines, lastLine)

		// Replace line and insert new lines
		b.lines = append(
			append(b.lines[:pos.Line], newLines...),
			b.lines[pos.Line+1:]...,
		)
	} else {
		// Single-line insert
		newLine := line[:pos.Column] + text + line[pos.Column:]
		b.lines[pos.Line] = newLine
	}

	return nil
}

// deleteText deletes text from start to end position.
func (b *SimpleBuffer) deleteText(from, to Position) error {
	if from.Line < 0 || from.Line >= len(b.lines) {
		return fmt.Errorf("from line %d out of range", from.Line)
	}
	if to.Line < 0 || to.Line >= len(b.lines) {
		return fmt.Errorf("to line %d out of range", to.Line)
	}

	// Same line deletion
	if from.Line == to.Line {
		line := b.lines[from.Line]
		newLine := line[:from.Column] + line[to.Column:]
		b.lines[from.Line] = newLine
		return nil
	}

	// Multi-line deletion
	startLine := b.lines[from.Line][:from.Column]
	endLine := b.lines[to.Line][to.Column:]
	mergedLine := startLine + endLine

	// Remove lines
	b.lines = append(
		append(b.lines[:from.Line], mergedLine),
		b.lines[to.Line+1:]...,
	)

	return nil
}

// extractText extracts text from start to end position.
func (b *SimpleBuffer) extractText(from, to Position) string {
	if from.Line == to.Line {
		// Single line
		line := b.lines[from.Line]
		return line[from.Column:to.Column]
	}

	// Multi-line
	var result strings.Builder

	// First line
	result.WriteString(b.lines[from.Line][from.Column:])
	result.WriteString("\n")

	// Middle lines
	for i := from.Line + 1; i < to.Line; i++ {
		result.WriteString(b.lines[i])
		result.WriteString("\n")
	}

	// Last line
	result.WriteString(b.lines[to.Line][:to.Column])

	return result.String()
}

// calculateEndPosition calculates the end position after inserting text.
func (b *SimpleBuffer) calculateEndPosition(start Position, text string) Position {
	if !strings.Contains(text, "\n") {
		// Single line
		return Position{
			Line:   start.Line,
			Column: start.Column + len(text),
		}
	}

	// Multi-line
	lines := strings.Split(text, "\n")
	return Position{
		Line:   start.Line + len(lines) - 1,
		Column: len(lines[len(lines)-1]),
	}
}
