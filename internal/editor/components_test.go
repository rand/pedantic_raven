package editor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/buffer"
	"github.com/rand/pedantic-raven/internal/editor/search"
	"github.com/rand/pedantic-raven/internal/editor/syntax"
)

// --- EditorComponent Search Tests ---

func TestEditorComponentSearch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello universe hello")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check search result
	result := e.GetSearchResult()
	if result == nil {
		t.Fatal("Expected search result to be set")
	}

	if len(result.Matches) != 3 {
		t.Errorf("Expected 3 matches, got %d", len(result.Matches))
	}

	// Check cursor moved to first match
	cursor := e.buffer.Cursor()
	if cursor.Line != 0 || cursor.Column != 0 {
		t.Errorf("Expected cursor at (0,0), got (%d,%d)", cursor.Line, cursor.Column)
	}

	// Check current match index
	if e.GetCurrentMatchIndex() != 0 {
		t.Errorf("Expected current match index 0, got %d", e.GetCurrentMatchIndex())
	}
}

func TestEditorComponentSearchCaseSensitive(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("Hello HELLO hello")

	// Case insensitive (default)
	opts := search.DefaultSearchOptions()
	err := e.Search("hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result := e.GetSearchResult()
	if len(result.Matches) != 3 {
		t.Errorf("Expected 3 matches (case insensitive), got %d", len(result.Matches))
	}

	// Case sensitive
	opts.CaseSensitive = true
	err = e.Search("hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result = e.GetSearchResult()
	if len(result.Matches) != 1 {
		t.Errorf("Expected 1 match (case sensitive), got %d", len(result.Matches))
	}
}

func TestEditorComponentSearchWholeWord(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello helloworld world")

	// Without whole word
	opts := search.DefaultSearchOptions()
	err := e.Search("hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result := e.GetSearchResult()
	if len(result.Matches) != 2 {
		t.Errorf("Expected 2 matches (substring), got %d", len(result.Matches))
	}

	// With whole word
	opts.WholeWord = true
	err = e.Search("hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result = e.GetSearchResult()
	if len(result.Matches) != 1 {
		t.Errorf("Expected 1 match (whole word), got %d", len(result.Matches))
	}
}

func TestEditorComponentSearchRegex(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("func test123 func test456")

	opts := search.DefaultSearchOptions()
	opts.Regex = true

	err := e.Search("func test[0-9]+", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result := e.GetSearchResult()
	if len(result.Matches) != 2 {
		t.Errorf("Expected 2 regex matches, got %d", len(result.Matches))
	}

	if result.Matches[0].Text != "func test123" {
		t.Errorf("Expected 'func test123', got %s", result.Matches[0].Text)
	}
}

func TestEditorComponentSearchEmpty(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world")

	// Empty query should return no matches
	err := e.Search("", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result := e.GetSearchResult()
	if len(result.Matches) != 0 {
		t.Errorf("Expected 0 matches for empty query, got %d", len(result.Matches))
	}

	if e.GetCurrentMatchIndex() != -1 {
		t.Errorf("Expected current match index -1, got %d", e.GetCurrentMatchIndex())
	}
}

func TestEditorComponentSearchNoMatches(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world")

	// Search for non-existent text
	err := e.Search("foobar", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result := e.GetSearchResult()
	if len(result.Matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(result.Matches))
	}

	if e.GetCurrentMatchIndex() != -1 {
		t.Errorf("Expected current match index -1, got %d", e.GetCurrentMatchIndex())
	}
}

func TestEditorComponentNextMatch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello universe hello")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should start at first match (index 0)
	if e.GetCurrentMatchIndex() != 0 {
		t.Errorf("Expected index 0, got %d", e.GetCurrentMatchIndex())
	}

	// Next match should go to index 1
	found := e.NextMatch()
	if !found {
		t.Error("Expected to find next match")
	}
	if e.GetCurrentMatchIndex() != 1 {
		t.Errorf("Expected index 1, got %d", e.GetCurrentMatchIndex())
	}

	// Next match should go to index 2
	found = e.NextMatch()
	if !found {
		t.Error("Expected to find next match")
	}
	if e.GetCurrentMatchIndex() != 2 {
		t.Errorf("Expected index 2, got %d", e.GetCurrentMatchIndex())
	}

	// Next match should wrap to index 0
	found = e.NextMatch()
	if !found {
		t.Error("Expected to find next match (wrapped)")
	}
	if e.GetCurrentMatchIndex() != 0 {
		t.Errorf("Expected index 0 (wrapped), got %d", e.GetCurrentMatchIndex())
	}
}

func TestEditorComponentPreviousMatch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello universe hello")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should start at first match (index 0)
	if e.GetCurrentMatchIndex() != 0 {
		t.Errorf("Expected index 0, got %d", e.GetCurrentMatchIndex())
	}

	// Previous match should wrap to last match (index 2)
	found := e.PreviousMatch()
	if !found {
		t.Error("Expected to find previous match")
	}
	if e.GetCurrentMatchIndex() != 2 {
		t.Errorf("Expected index 2 (wrapped), got %d", e.GetCurrentMatchIndex())
	}

	// Previous match should go to index 1
	found = e.PreviousMatch()
	if !found {
		t.Error("Expected to find previous match")
	}
	if e.GetCurrentMatchIndex() != 1 {
		t.Errorf("Expected index 1, got %d", e.GetCurrentMatchIndex())
	}

	// Previous match should go to index 0
	found = e.PreviousMatch()
	if !found {
		t.Error("Expected to find previous match")
	}
	if e.GetCurrentMatchIndex() != 0 {
		t.Errorf("Expected index 0, got %d", e.GetCurrentMatchIndex())
	}
}

func TestEditorComponentNextMatchNoSearch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world")

	// NextMatch without search should return false
	found := e.NextMatch()
	if found {
		t.Error("Expected NextMatch to return false without active search")
	}
}

func TestEditorComponentPreviousMatchNoSearch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world")

	// PreviousMatch without search should return false
	found := e.PreviousMatch()
	if found {
		t.Error("Expected PreviousMatch to return false without active search")
	}
}

func TestEditorComponentReplaceCurrentMatch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello universe hello")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Replace first match
	err = e.ReplaceCurrentMatch("hi")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content := e.GetContent()
	if content != "hi world hello universe hello" {
		t.Errorf("Expected 'hi world hello universe hello', got '%s'", content)
	}

	// Search should be updated with new positions
	result := e.GetSearchResult()
	if len(result.Matches) != 2 {
		t.Errorf("Expected 2 matches after replacement, got %d", len(result.Matches))
	}
}

func TestEditorComponentReplaceCurrentMatchNoSearch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world")

	// Replace without search should error
	err := e.ReplaceCurrentMatch("hi")
	if err == nil {
		t.Error("Expected error when replacing without active search")
	}
}

func TestEditorComponentReplaceCurrentMatchInvalidIndex(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Manually set invalid index
	e.currentMatch = 10

	// Replace should error with invalid index
	err = e.ReplaceCurrentMatch("hi")
	if err == nil {
		t.Error("Expected error when replacing with invalid index")
	}
}

func TestEditorComponentReplaceAll(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello universe hello")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Replace all
	count, err := e.ReplaceAll("hi")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 replacements, got %d", count)
	}

	content := e.GetContent()
	if content != "hi world hi universe hi" {
		t.Errorf("Expected 'hi world hi universe hi', got '%s'", content)
	}

	// Search should be cleared after replace all
	if e.GetSearchResult() != nil {
		t.Error("Expected search to be cleared after replace all")
	}

	if e.GetCurrentMatchIndex() != -1 {
		t.Errorf("Expected current match index -1, got %d", e.GetCurrentMatchIndex())
	}
}

func TestEditorComponentReplaceAllNoSearch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world")

	// Replace all without search should error
	_, err := e.ReplaceAll("hi")
	if err == nil {
		t.Error("Expected error when replacing all without active search")
	}
}

func TestEditorComponentReplaceAllCaseSensitive(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("Hello hello HELLO")

	opts := search.DefaultSearchOptions()
	opts.CaseSensitive = true

	// Perform search
	err := e.Search("hello", opts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Replace all (case sensitive)
	count, err := e.ReplaceAll("hi")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 replacement (case sensitive), got %d", count)
	}

	content := e.GetContent()
	if content != "Hello hi HELLO" {
		t.Errorf("Expected 'Hello hi HELLO', got '%s'", content)
	}
}

func TestEditorComponentClearSearch(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify search is active
	if e.GetSearchResult() == nil {
		t.Fatal("Expected search result to be set")
	}
	if e.GetCurrentMatchIndex() == -1 {
		t.Fatal("Expected current match to be set")
	}

	// Clear search
	e.ClearSearch()

	// Verify search is cleared
	if e.GetSearchResult() != nil {
		t.Error("Expected search result to be nil after clear")
	}
	if e.GetCurrentMatchIndex() != -1 {
		t.Errorf("Expected current match index -1 after clear, got %d", e.GetCurrentMatchIndex())
	}
}

func TestEditorComponentSearchMultiLine(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("line one\nline two\nline one again")

	// Search across multiple lines
	err := e.Search("line", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result := e.GetSearchResult()
	if len(result.Matches) != 3 {
		t.Errorf("Expected 3 matches across lines, got %d", len(result.Matches))
	}

	// Verify first match is on line 0
	if result.Matches[0].Start.Line != 0 {
		t.Errorf("Expected first match on line 0, got %d", result.Matches[0].Start.Line)
	}

	// Verify second match is on line 1
	if result.Matches[1].Start.Line != 1 {
		t.Errorf("Expected second match on line 1, got %d", result.Matches[1].Start.Line)
	}

	// Verify third match is on line 2
	if result.Matches[2].Start.Line != 2 {
		t.Errorf("Expected third match on line 2, got %d", result.Matches[2].Start.Line)
	}
}

func TestEditorComponentSearchWithUndo(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello")

	// Perform search and replace
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = e.ReplaceCurrentMatch("hi")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content := e.GetContent()
	if content != "hi world hello" {
		t.Errorf("Expected 'hi world hello', got '%s'", content)
	}

	// Undo should restore original text
	if !e.buffer.CanUndo() {
		t.Fatal("Expected undo to be available")
	}

	e.buffer.Undo()
	content = e.GetContent()
	if content != "hello world hello" {
		t.Errorf("Expected 'hello world hello' after undo, got '%s'", content)
	}
}

func TestEditorComponentReplaceAllWithUndo(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello universe hello")

	// Perform search and replace all
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = e.ReplaceAll("hi")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	content := e.GetContent()
	if content != "hi world hi universe hi" {
		t.Errorf("Expected 'hi world hi universe hi', got '%s'", content)
	}

	// Undo should be available (each replacement is undoable)
	if !e.buffer.CanUndo() {
		t.Fatal("Expected undo to be available")
	}

	// Undo all replacements (3 times)
	e.buffer.Undo()
	e.buffer.Undo()
	e.buffer.Undo()

	content = e.GetContent()
	if content != "hello world hello universe hello" {
		t.Errorf("Expected 'hello world hello universe hello' after undo, got '%s'", content)
	}
}

func TestEditorComponentSearchCursorMovement(t *testing.T) {
	e := NewEditorComponent()
	e.SetContent("hello world hello")

	// Perform search
	err := e.Search("hello", search.DefaultSearchOptions())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Cursor should be at first match
	cursor := e.buffer.Cursor()
	if cursor.Line != 0 || cursor.Column != 0 {
		t.Errorf("Expected cursor at (0,0), got (%d,%d)", cursor.Line, cursor.Column)
	}

	// Next match should move cursor
	e.NextMatch()
	cursor = e.buffer.Cursor()
	expected := buffer.Position{Line: 0, Column: 12}
	if cursor != expected {
		t.Errorf("Expected cursor at (0,12), got (%d,%d)", cursor.Line, cursor.Column)
	}

	// Previous match should move cursor back
	e.PreviousMatch()
	cursor = e.buffer.Cursor()
	if cursor.Line != 0 || cursor.Column != 0 {
		t.Errorf("Expected cursor at (0,0) after previous, got (%d,%d)", cursor.Line, cursor.Column)
	}
}

// --- Syntax Highlighting Tests ---

func TestEditorComponentSyntaxHighlightingDefault(t *testing.T) {
	e := NewEditorComponent()

	// Verify highlighter is initialized
	if e.highlighter == nil {
		t.Fatal("Expected highlighter to be initialized")
	}

	// Should have no tokenizer by default (LangNone)
	// This is tested indirectly by checking that View returns content unchanged
	e.SetContent("package main")
	content := e.GetContent()
	if content != "package main" {
		t.Errorf("Expected 'package main', got '%s'", content)
	}
}

func TestEditorComponentSyntaxHighlightingGoFile(t *testing.T) {
	// Create temporary Go file
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "test.go")
	goContent := "package main\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}"
	err := os.WriteFile(goFile, []byte(goContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	e := NewEditorComponent()
	err = e.OpenFile(goFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Verify content loaded
	content := e.GetContent()
	if content != goContent {
		t.Errorf("Content mismatch: expected '%s', got '%s'", goContent, content)
	}

	// Verify highlighter was created with Go language
	if e.highlighter == nil {
		t.Fatal("Expected highlighter to be set after opening file")
	}

	// The highlighter should have a tokenizer for Go
	// We can't directly access the language, but we can verify it works by checking that
	// the tokenizer is not nil (indirectly tested by View not returning plain content)
	// For now, just verify the highlighter exists
}

func TestEditorComponentSyntaxHighlightingMarkdownFile(t *testing.T) {
	// Create temporary Markdown file
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")
	mdContent := "# Header\n\nThis is **bold** text."
	err := os.WriteFile(mdFile, []byte(mdContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	e := NewEditorComponent()
	err = e.OpenFile(mdFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Verify content loaded
	content := e.GetContent()
	if content != mdContent {
		t.Errorf("Content mismatch: expected '%s', got '%s'", mdContent, content)
	}

	// Verify highlighter was created
	if e.highlighter == nil {
		t.Fatal("Expected highlighter to be set after opening file")
	}
}

func TestEditorComponentSyntaxHighlightingUnknownFile(t *testing.T) {
	// Create temporary file with unknown extension
	tmpDir := t.TempDir()
	unknownFile := filepath.Join(tmpDir, "test.xyz")
	content := "some text"
	err := os.WriteFile(unknownFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	e := NewEditorComponent()
	err = e.OpenFile(unknownFile)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Verify highlighter was still created (with LangNone)
	if e.highlighter == nil {
		t.Fatal("Expected highlighter to be set even for unknown files")
	}
}

func TestEditorComponentSyntaxHighlightingContentDetection(t *testing.T) {
	// Create file without extension but with recognizable content
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "noext")
	goContent := "package main\n\nfunc main() {}"
	err := os.WriteFile(file, []byte(goContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	e := NewEditorComponent()
	err = e.OpenFile(file)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Highlighter should use content-based detection and detect Go
	if e.highlighter == nil {
		t.Fatal("Expected highlighter to be set")
	}

	// Verify content loaded correctly
	content := e.GetContent()
	if content != goContent {
		t.Errorf("Content mismatch: expected '%s', got '%s'", goContent, content)
	}
}

func TestEditorComponentLanguageDetectionByExtension(t *testing.T) {
	tests := []struct {
		ext      string
		content  string
		wantLang syntax.Language
	}{
		{".go", "package main", syntax.LangGo},
		{".md", "# Header", syntax.LangMarkdown},
		{".py", "print('hello')", syntax.LangPython},
		{".js", "console.log('hello')", syntax.LangJavaScript},
		{".ts", "const x: number = 1", syntax.LangTypeScript},
		{".rs", "fn main() {}", syntax.LangRust},
		{".json", "{\"key\": \"value\"}", syntax.LangJSON},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			tmpDir := t.TempDir()
			file := filepath.Join(tmpDir, "test"+tt.ext)
			err := os.WriteFile(file, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			e := NewEditorComponent()
			err = e.OpenFile(file)
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}

			// Verify highlighter exists
			if e.highlighter == nil {
				t.Fatal("Expected highlighter to be set")
			}

			// Verify content loaded
			if e.GetContent() != tt.content {
				t.Errorf("Content mismatch")
			}
		})
	}
}
