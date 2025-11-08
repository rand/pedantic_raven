package search

import (
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/buffer"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("Expected engine to be created")
	}
}

func TestSearchEmpty(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "Hello world")

	result, err := engine.Search(buf, "", DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Matches) != 0 {
		t.Error("Expected no matches for empty query")
	}
}

func TestSearchCaseSensitive(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "Hello HELLO hello")

	// Case insensitive (default)
	opts := DefaultSearchOptions()
	result, err := engine.Search(buf, "hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Matches) != 3 {
		t.Errorf("Expected 3 matches (case insensitive), got %d", len(result.Matches))
	}

	// Case sensitive
	opts.CaseSensitive = true
	result, err = engine.Search(buf, "hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Matches) != 1 {
		t.Errorf("Expected 1 match (case sensitive), got %d", len(result.Matches))
	}
}

func TestSearchWholeWord(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello helloworld world")

	// Without whole word
	opts := DefaultSearchOptions()
	result, err := engine.Search(buf, "hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Matches) != 2 {
		t.Errorf("Expected 2 matches (substring), got %d", len(result.Matches))
	}

	// With whole word
	opts.WholeWord = true
	result, err = engine.Search(buf, "hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Matches) != 1 {
		t.Errorf("Expected 1 match (whole word), got %d", len(result.Matches))
	}
}

func TestSearchMultiLine(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "line one\nline two\nline one again")

	result, err := engine.Search(buf, "line", DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Matches) != 3 {
		t.Errorf("Expected 3 matches across lines, got %d", len(result.Matches))
	}

	// Verify positions
	if result.Matches[0].Start.Line != 0 {
		t.Errorf("Expected first match on line 0, got %d", result.Matches[0].Start.Line)
	}
	if result.Matches[1].Start.Line != 1 {
		t.Errorf("Expected second match on line 1, got %d", result.Matches[1].Start.Line)
	}
	if result.Matches[2].Start.Line != 2 {
		t.Errorf("Expected third match on line 2, got %d", result.Matches[2].Start.Line)
	}
}

func TestSearchRegex(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "func test123 func test456")

	opts := DefaultSearchOptions()
	opts.Regex = true

	result, err := engine.Search(buf, "func test[0-9]+", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Matches) != 2 {
		t.Errorf("Expected 2 regex matches, got %d", len(result.Matches))
	}

	if result.Matches[0].Text != "func test123" {
		t.Errorf("Expected 'func test123', got %s", result.Matches[0].Text)
	}
}

func TestSearchRegexInvalid(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "test")

	opts := DefaultSearchOptions()
	opts.Regex = true

	_, err := engine.Search(buf, "[invalid(regex", opts)
	if err == nil {
		t.Error("Expected error for invalid regex")
	}
}

func TestFindNext(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello world hello universe hello")

	opts := DefaultSearchOptions()

	// Find first match (after position before first match)
	match, err := engine.FindNext(buf, "hello", buffer.Position{Line: 0, Column: -1}, opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match == nil {
		t.Fatal("Expected to find match")
	}
	// First "hello" is at column 0, but FindNext finds match AFTER position
	// Since we're at -1 (before 0), we should find the one at column 0
	if match.Start.Column != 0 {
		t.Errorf("Expected first match at column 0, got %d", match.Start.Column)
	}

	// Find second match (after first match)
	match, err = engine.FindNext(buf, "hello", buffer.Position{Line: 0, Column: 5}, opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match == nil {
		t.Fatal("Expected to find match")
	}
	if match.Start.Column != 12 {
		t.Errorf("Expected second match at column 12, got %d", match.Start.Column)
	}
}

func TestFindNextWrapAround(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello world")

	opts := DefaultSearchOptions()
	opts.WrapAround = true

	// Search from end should wrap to beginning
	match, err := engine.FindNext(buf, "hello", buffer.Position{Line: 0, Column: 10}, opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match == nil {
		t.Fatal("Expected to find match (wrapped)")
	}
	if match.Start.Column != 0 {
		t.Errorf("Expected wrapped match at column 0, got %d", match.Start.Column)
	}
}

func TestFindPrevious(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello world hello universe hello")

	opts := DefaultSearchOptions()

	// Find from end
	match, err := engine.FindPrevious(buf, "hello", buffer.Position{Line: 0, Column: 32}, opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match == nil {
		t.Fatal("Expected to find match")
	}
	if match.Start.Column != 27 {
		t.Errorf("Expected match at column 27, got %d", match.Start.Column)
	}
}

func TestFindPreviousWrapAround(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello world")

	opts := DefaultSearchOptions()
	opts.WrapAround = true

	// Search from beginning should wrap to end
	match, err := engine.FindPrevious(buf, "world", buffer.Position{Line: 0, Column: 0}, opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if match == nil {
		t.Fatal("Expected to find match (wrapped)")
	}
	if match.Start.Column != 6 {
		t.Errorf("Expected wrapped match at column 6, got %d", match.Start.Column)
	}
}

func TestReplace(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello world")

	// Find match
	result, _ := engine.Search(buf, "world", DefaultSearchOptions())
	if len(result.Matches) == 0 {
		t.Fatal("Expected to find match")
	}

	// Replace
	err := engine.Replace(buf, result.Matches[0], "universe")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content := buf.Content()
	if content != "hello universe" {
		t.Errorf("Expected 'hello universe', got %s", content)
	}
}

func TestReplaceAll(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello world hello universe hello")

	count, err := engine.ReplaceAll(buf, "hello", "hi", DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 replacements, got %d", count)
	}

	content := buf.Content()
	if content != "hi world hi universe hi" {
		t.Errorf("Expected 'hi world hi universe hi', got %s", content)
	}
}

func TestReplaceAllCaseSensitive(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "Hello hello HELLO")

	opts := DefaultSearchOptions()
	opts.CaseSensitive = true

	count, err := engine.ReplaceAll(buf, "hello", "hi", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 replacement (case sensitive), got %d", count)
	}

	content := buf.Content()
	if content != "Hello hi HELLO" {
		t.Errorf("Expected 'Hello hi HELLO', got %s", content)
	}
}

func TestReplaceWithUndo(t *testing.T) {
	engine := NewEngine()
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "hello world")

	// Find and replace
	result, _ := engine.Search(buf, "world", DefaultSearchOptions())
	engine.Replace(buf, result.Matches[0], "universe")

	// Undo should work
	if !buf.CanUndo() {
		t.Error("Expected undo to be available")
	}

	buf.Undo()
	content := buf.Content()
	if content != "hello world" {
		t.Errorf("Expected 'hello world' after undo, got %s", content)
	}
}

func TestOffsetToPosition(t *testing.T) {
	engine := &SimpleEngine{}
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "line 1\nline 2\nline 3")

	// Test various offsets
	tests := []struct {
		offset int
		line   int
		column int
	}{
		{0, 0, 0},
		{5, 0, 5},
		{7, 1, 0},  // After first newline
		{14, 2, 0}, // After second newline
	}

	for _, test := range tests {
		pos := engine.offsetToPosition(buf, test.offset)
		if pos.Line != test.line || pos.Column != test.column {
			t.Errorf("Offset %d: expected (%d,%d), got (%d,%d)",
				test.offset, test.line, test.column, pos.Line, pos.Column)
		}
	}
}

func TestPositionToOffset(t *testing.T) {
	engine := &SimpleEngine{}
	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "line 1\nline 2\nline 3")

	// Test various positions
	tests := []struct {
		line   int
		column int
		offset int
	}{
		{0, 0, 0},
		{0, 5, 5},
		{1, 0, 7},  // Start of second line
		{2, 0, 14}, // Start of third line
	}

	for _, test := range tests {
		pos := buffer.Position{Line: test.line, Column: test.column}
		offset := engine.positionToOffset(buf, pos)
		if offset != test.offset {
			t.Errorf("Position (%d,%d): expected offset %d, got %d",
				test.line, test.column, test.offset, offset)
		}
	}
}
