package buffer

import (
	"testing"
)

// --- Buffer Creation Tests ---

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer("test")

	if buf.ID() != "test" {
		t.Errorf("Expected ID 'test', got '%s'", buf.ID())
	}

	if buf.Content() != "" {
		t.Errorf("Expected empty content, got '%s'", buf.Content())
	}

	if buf.LineCount() != 1 {
		t.Errorf("Expected 1 line, got %d", buf.LineCount())
	}

	if buf.IsDirty() {
		t.Error("New buffer should not be dirty")
	}
}

func TestNewBufferFromContent(t *testing.T) {
	content := "Hello\nWorld"
	buf := NewBufferFromContent("test", content)

	if buf.Content() != content {
		t.Errorf("Expected content '%s', got '%s'", content, buf.Content())
	}

	if buf.LineCount() != 2 {
		t.Errorf("Expected 2 lines, got %d", buf.LineCount())
	}

	if buf.IsDirty() {
		t.Error("Buffer created from content should start clean")
	}
}

// --- Insert Operations ---

func TestInsertSingleLine(t *testing.T) {
	buf := NewBuffer("test")

	op, err := buf.Insert(Position{Line: 0, Column: 0}, "Hello")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if op.Type != OpInsert {
		t.Errorf("Expected OpInsert, got %v", op.Type)
	}

	if buf.Content() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", buf.Content())
	}

	if !buf.IsDirty() {
		t.Error("Buffer should be dirty after insert")
	}
}

func TestInsertMultiLine(t *testing.T) {
	buf := NewBuffer("test")

	_, err := buf.Insert(Position{Line: 0, Column: 0}, "Line 1\nLine 2\nLine 3")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if buf.LineCount() != 3 {
		t.Errorf("Expected 3 lines, got %d", buf.LineCount())
	}

	if buf.Line(0) != "Line 1" {
		t.Errorf("Expected 'Line 1', got '%s'", buf.Line(0))
	}

	if buf.Line(1) != "Line 2" {
		t.Errorf("Expected 'Line 2', got '%s'", buf.Line(1))
	}

	if buf.Line(2) != "Line 3" {
		t.Errorf("Expected 'Line 3', got '%s'", buf.Line(2))
	}
}

func TestInsertMiddle(t *testing.T) {
	buf := NewBufferFromContent("test", "HelloWorld")

	_, err := buf.Insert(Position{Line: 0, Column: 5}, " ")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if buf.Content() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", buf.Content())
	}
}

func TestInsertOutOfBounds(t *testing.T) {
	buf := NewBuffer("test")

	_, err := buf.Insert(Position{Line: 10, Column: 0}, "text")
	if err == nil {
		t.Error("Expected error for out of bounds insert")
	}

	_, err = buf.Insert(Position{Line: 0, Column: 100}, "text")
	if err == nil {
		t.Error("Expected error for out of bounds column")
	}
}

// --- Delete Operations ---

func TestDeleteSingleLine(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello World")

	_, err := buf.Delete(Position{Line: 0, Column: 6}, Position{Line: 0, Column: 11})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if buf.Content() != "Hello " {
		t.Errorf("Expected 'Hello ', got '%s'", buf.Content())
	}
}

func TestDeleteMultiLine(t *testing.T) {
	buf := NewBufferFromContent("test", "Line 1\nLine 2\nLine 3")

	// Delete from "Lin|e 1" to "Line| 3" (column 4 is the space)
	// Result: "Lin" + " 3" = "Lin 3"
	_, err := buf.Delete(Position{Line: 0, Column: 3}, Position{Line: 2, Column: 4})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if buf.LineCount() != 1 {
		t.Errorf("Expected 1 line, got %d", buf.LineCount())
	}

	expected := "Lin 3"
	if buf.Content() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.Content())
	}
}

func TestDeleteReversedPositions(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello World")

	// Delete with reversed positions (should swap automatically)
	_, err := buf.Delete(Position{Line: 0, Column: 11}, Position{Line: 0, Column: 6})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if buf.Content() != "Hello " {
		t.Errorf("Expected 'Hello ', got '%s'", buf.Content())
	}
}

// --- Replace Operations ---

func TestReplaceSingleLine(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello World")

	_, err := buf.Replace(Position{Line: 0, Column: 6}, Position{Line: 0, Column: 11}, "Universe")
	if err != nil {
		t.Fatalf("Replace failed: %v", err)
	}

	if buf.Content() != "Hello Universe" {
		t.Errorf("Expected 'Hello Universe', got '%s'", buf.Content())
	}
}

func TestReplaceMultiLine(t *testing.T) {
	buf := NewBufferFromContent("test", "Line 1\nLine 2\nLine 3")

	// Replace from "Line | 1" (after "Line ") to "Line | 3" (after "Line ")
	// Result: "Line " + " Replacement" + "3" = "Line  Replacement3"
	_, err := buf.Replace(
		Position{Line: 0, Column: 5},
		Position{Line: 2, Column: 5},
		" Replacement",
	)
	if err != nil {
		t.Fatalf("Replace failed: %v", err)
	}

	expected := "Line  Replacement3"
	if buf.Content() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.Content())
	}
}

// --- Undo/Redo Tests ---

func TestUndoInsert(t *testing.T) {
	buf := NewBuffer("test")

	// Insert text
	_, err := buf.Insert(Position{Line: 0, Column: 0}, "Hello")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if buf.Content() != "Hello" {
		t.Fatalf("Expected 'Hello', got '%s'", buf.Content())
	}

	// Undo
	if !buf.Undo() {
		t.Fatal("Undo should succeed")
	}

	if buf.Content() != "" {
		t.Errorf("Expected empty content after undo, got '%s'", buf.Content())
	}
}

func TestUndoDelete(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello World")

	// Delete text
	_, err := buf.Delete(Position{Line: 0, Column: 6}, Position{Line: 0, Column: 11})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if buf.Content() != "Hello " {
		t.Fatalf("Expected 'Hello ', got '%s'", buf.Content())
	}

	// Undo
	if !buf.Undo() {
		t.Fatal("Undo should succeed")
	}

	if buf.Content() != "Hello World" {
		t.Errorf("Expected 'Hello World' after undo, got '%s'", buf.Content())
	}
}

func TestUndoReplace(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello World")

	// Replace text
	_, err := buf.Replace(Position{Line: 0, Column: 6}, Position{Line: 0, Column: 11}, "Universe")
	if err != nil {
		t.Fatalf("Replace failed: %v", err)
	}

	if buf.Content() != "Hello Universe" {
		t.Fatalf("Expected 'Hello Universe', got '%s'", buf.Content())
	}

	// Undo
	if !buf.Undo() {
		t.Fatal("Undo should succeed")
	}

	if buf.Content() != "Hello World" {
		t.Errorf("Expected 'Hello World' after undo, got '%s'", buf.Content())
	}
}

func TestRedoInsert(t *testing.T) {
	buf := NewBuffer("test")

	// Insert text
	buf.Insert(Position{Line: 0, Column: 0}, "Hello")

	// Undo
	buf.Undo()

	// Redo
	if !buf.Redo() {
		t.Fatal("Redo should succeed")
	}

	if buf.Content() != "Hello" {
		t.Errorf("Expected 'Hello' after redo, got '%s'", buf.Content())
	}
}

func TestMultipleUndoRedo(t *testing.T) {
	buf := NewBuffer("test")

	// Perform multiple operations
	buf.Insert(Position{Line: 0, Column: 0}, "A")
	buf.Insert(Position{Line: 0, Column: 1}, "B")
	buf.Insert(Position{Line: 0, Column: 2}, "C")

	if buf.Content() != "ABC" {
		t.Fatalf("Expected 'ABC', got '%s'", buf.Content())
	}

	// Undo all
	buf.Undo()
	if buf.Content() != "AB" {
		t.Errorf("Expected 'AB' after undo 1, got '%s'", buf.Content())
	}

	buf.Undo()
	if buf.Content() != "A" {
		t.Errorf("Expected 'A' after undo 2, got '%s'", buf.Content())
	}

	buf.Undo()
	if buf.Content() != "" {
		t.Errorf("Expected empty after undo 3, got '%s'", buf.Content())
	}

	// Redo all
	buf.Redo()
	if buf.Content() != "A" {
		t.Errorf("Expected 'A' after redo 1, got '%s'", buf.Content())
	}

	buf.Redo()
	if buf.Content() != "AB" {
		t.Errorf("Expected 'AB' after redo 2, got '%s'", buf.Content())
	}

	buf.Redo()
	if buf.Content() != "ABC" {
		t.Errorf("Expected 'ABC' after redo 3, got '%s'", buf.Content())
	}
}

func TestUndoRedoBounds(t *testing.T) {
	buf := NewBuffer("test")

	// Undo on empty history
	if buf.Undo() {
		t.Error("Undo should fail on empty history")
	}

	// Redo on empty redo stack
	if buf.Redo() {
		t.Error("Redo should fail on empty redo stack")
	}

	// Insert something
	buf.Insert(Position{Line: 0, Column: 0}, "Test")

	// Redo without undo
	if buf.Redo() {
		t.Error("Redo should fail without undo")
	}
}

func TestUndoTruncatesRedo(t *testing.T) {
	buf := NewBuffer("test")

	// Perform operations
	buf.Insert(Position{Line: 0, Column: 0}, "A")
	buf.Insert(Position{Line: 0, Column: 1}, "B")

	// Undo once
	buf.Undo()

	if buf.Content() != "A" {
		t.Fatalf("Expected 'A', got '%s'", buf.Content())
	}

	// Perform new operation (should truncate redo history)
	buf.Insert(Position{Line: 0, Column: 1}, "C")

	if buf.Content() != "AC" {
		t.Fatalf("Expected 'AC', got '%s'", buf.Content())
	}

	// Redo should fail (history truncated)
	if buf.Redo() {
		t.Error("Redo should fail after new operation")
	}
}

// --- Cursor Management ---

func TestCursorPosition(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello World")

	// Initial cursor at 0,0
	cursor := buf.Cursor()
	if cursor.Line != 0 || cursor.Column != 0 {
		t.Errorf("Expected cursor at (0,0), got (%d,%d)", cursor.Line, cursor.Column)
	}

	// Set cursor
	buf.SetCursor(Position{Line: 0, Column: 5})
	cursor = buf.Cursor()
	if cursor.Line != 0 || cursor.Column != 5 {
		t.Errorf("Expected cursor at (0,5), got (%d,%d)", cursor.Line, cursor.Column)
	}
}

func TestCursorClamping(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello\nWorld")

	// Set cursor beyond line bounds
	buf.SetCursor(Position{Line: 10, Column: 0})
	cursor := buf.Cursor()
	if cursor.Line != 1 {
		t.Errorf("Expected cursor line 1, got %d", cursor.Line)
	}

	// Set cursor beyond column bounds
	buf.SetCursor(Position{Line: 0, Column: 100})
	cursor = buf.Cursor()
	if cursor.Column != 5 {
		t.Errorf("Expected cursor column 5, got %d", cursor.Column)
	}

	// Negative values
	buf.SetCursor(Position{Line: -1, Column: -1})
	cursor = buf.Cursor()
	if cursor.Line != 0 || cursor.Column != 0 {
		t.Errorf("Expected cursor at (0,0), got (%d,%d)", cursor.Line, cursor.Column)
	}
}

// --- Lines and Content ---

func TestLines(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	buf := NewBufferFromContent("test", content)

	lines := buf.Lines()

	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines, got %d", len(lines))
	}

	if lines[0] != "Line 1" || lines[1] != "Line 2" || lines[2] != "Line 3" {
		t.Errorf("Lines content incorrect: %v", lines)
	}

	// Verify it's a copy (modification doesn't affect buffer)
	lines[0] = "Modified"
	if buf.Line(0) != "Line 1" {
		t.Error("Lines() should return a copy")
	}
}

func TestLine(t *testing.T) {
	buf := NewBufferFromContent("test", "Line 1\nLine 2")

	// Valid lines
	if buf.Line(0) != "Line 1" {
		t.Errorf("Expected 'Line 1', got '%s'", buf.Line(0))
	}

	if buf.Line(1) != "Line 2" {
		t.Errorf("Expected 'Line 2', got '%s'", buf.Line(1))
	}

	// Out of bounds
	if buf.Line(-1) != "" {
		t.Error("Expected empty string for negative line")
	}

	if buf.Line(10) != "" {
		t.Error("Expected empty string for out of bounds line")
	}
}

// --- Dirty State ---

func TestDirtyState(t *testing.T) {
	buf := NewBuffer("test")

	if buf.IsDirty() {
		t.Error("New buffer should not be dirty")
	}

	// Insert makes it dirty
	buf.Insert(Position{Line: 0, Column: 0}, "Text")
	if !buf.IsDirty() {
		t.Error("Buffer should be dirty after insert")
	}

	// MarkClean
	buf.MarkClean()
	if buf.IsDirty() {
		t.Error("Buffer should not be dirty after MarkClean")
	}
}

// --- Path Management ---

func TestPath(t *testing.T) {
	buf := NewBuffer("test")

	if buf.Path() != "" {
		t.Errorf("Expected empty path, got '%s'", buf.Path())
	}

	buf.SetPath("/tmp/test.txt")
	if buf.Path() != "/tmp/test.txt" {
		t.Errorf("Expected '/tmp/test.txt', got '%s'", buf.Path())
	}
}

// --- Clear ---

func TestClear(t *testing.T) {
	buf := NewBufferFromContent("test", "Hello\nWorld")

	buf.Clear()

	if buf.Content() != "" {
		t.Errorf("Expected empty content after clear, got '%s'", buf.Content())
	}

	if buf.LineCount() != 1 {
		t.Errorf("Expected 1 line after clear, got %d", buf.LineCount())
	}

	cursor := buf.Cursor()
	if cursor.Line != 0 || cursor.Column != 0 {
		t.Errorf("Expected cursor at (0,0) after clear, got (%d,%d)", cursor.Line, cursor.Column)
	}
}

// --- History Tests ---

func TestHistoryCanUndoRedo(t *testing.T) {
	buf := NewBuffer("test")

	if buf.CanUndo() {
		t.Error("Should not be able to undo empty buffer")
	}

	if buf.CanRedo() {
		t.Error("Should not be able to redo empty buffer")
	}

	buf.Insert(Position{Line: 0, Column: 0}, "Test")

	if !buf.CanUndo() {
		t.Error("Should be able to undo after insert")
	}

	if buf.CanRedo() {
		t.Error("Should not be able to redo before undo")
	}

	buf.Undo()

	if buf.CanUndo() {
		t.Error("Should not be able to undo after single undo")
	}

	if !buf.CanRedo() {
		t.Error("Should be able to redo after undo")
	}
}

// --- Complex Scenarios ---

func TestComplexEditingScenario(t *testing.T) {
	buf := NewBuffer("test")

	// Build up text with multiple operations
	buf.Insert(Position{Line: 0, Column: 0}, "The quick brown fox")
	buf.Insert(Position{Line: 0, Column: 19}, " jumps")
	buf.Insert(Position{Line: 0, Column: 25}, "\nover the lazy dog")

	expected := "The quick brown fox jumps\nover the lazy dog"
	if buf.Content() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.Content())
	}

	// Replace "lazy" with "sleeping"
	buf.Replace(Position{Line: 1, Column: 9}, Position{Line: 1, Column: 13}, "sleeping")

	expected = "The quick brown fox jumps\nover the sleeping dog"
	if buf.Content() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, buf.Content())
	}

	// Undo the replace
	buf.Undo()

	expected = "The quick brown fox jumps\nover the lazy dog"
	if buf.Content() != expected {
		t.Errorf("Expected '%s' after undo, got '%s'", expected, buf.Content())
	}
}

// --- OpType Tests ---

func TestOpTypeString(t *testing.T) {
	tests := []struct {
		opType   OpType
		expected string
	}{
		{OpInsert, "Insert"},
		{OpDelete, "Delete"},
		{OpReplace, "Replace"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.opType.String() != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, tt.opType.String())
			}
		})
	}
}
