package syntax

import (
	"testing"

	"github.com/rand/pedantic-raven/internal/editor/buffer"
)

// --- TokenType Tests ---

func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TokenKeyword, "keyword"},
		{TokenString, "string"},
		{TokenComment, "comment"},
		{TokenNumber, "number"},
		{TokenOperator, "operator"},
		{TokenIdentifier, "identifier"},
		{TokenFunction, "function"},
		{TokenType, "type"},
		{TokenConstant, "constant"},
		{TokenPunctuation, "punctuation"},
		{TokenWhitespace, "whitespace"},
		{TokenNone, "none"},
	}

	for _, test := range tests {
		result := test.tokenType.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

// --- Language Tests ---

func TestLanguageString(t *testing.T) {
	tests := []struct {
		lang     Language
		expected string
	}{
		{LangGo, "go"},
		{LangMarkdown, "markdown"},
		{LangPython, "python"},
		{LangJavaScript, "javascript"},
		{LangTypeScript, "typescript"},
		{LangRust, "rust"},
		{LangJSON, "json"},
		{LangNone, "none"},
	}

	for _, test := range tests {
		result := test.lang.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

// --- Go Tokenizer Tests ---

func TestGoTokenizerKeywords(t *testing.T) {
	tokenizer := NewGoTokenizer()
	line := "func main() { var x int }"

	tokens := tokenizer.Tokenize(line, 0)

	// Find keyword tokens
	keywords := []string{}
	for _, token := range tokens {
		if token.Type == TokenKeyword {
			keywords = append(keywords, token.Text)
		}
	}

	expectedKeywords := []string{"func", "var"}
	if len(keywords) != len(expectedKeywords) {
		t.Errorf("Expected %d keywords, got %d", len(expectedKeywords), len(keywords))
	}

	for i, kw := range expectedKeywords {
		if i >= len(keywords) || keywords[i] != kw {
			t.Errorf("Expected keyword %s at position %d, got %v", kw, i, keywords)
			break
		}
	}
}

func TestGoTokenizerStrings(t *testing.T) {
	tokenizer := NewGoTokenizer()

	tests := []struct {
		line     string
		expected string
	}{
		{`"hello world"`, `"hello world"`},
		{`"escaped \"quote\""`, `"escaped \"quote\""`},
		{"`raw string`", "`raw string`"},
		{"'a'", "'a'"},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.line, 0)

		found := false
		for _, token := range tokens {
			if token.Type == TokenString {
				if token.Text != test.expected {
					t.Errorf("Expected string %s, got %s", test.expected, token.Text)
				}
				found = true
				break
			}
		}

		if !found {
			t.Errorf("No string token found for: %s", test.line)
		}
	}
}

func TestGoTokenizerComments(t *testing.T) {
	tokenizer := NewGoTokenizer()

	tests := []struct {
		line     string
		expected string
	}{
		{"// single line comment", "// single line comment"},
		{"/* block comment */", "/* block comment */"},
		{"code // comment", "// comment"},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.line, 0)

		found := false
		for _, token := range tokens {
			if token.Type == TokenComment {
				if token.Text != test.expected {
					t.Errorf("Expected comment %s, got %s", test.expected, token.Text)
				}
				found = true
				break
			}
		}

		if !found {
			t.Errorf("No comment token found for: %s", test.line)
		}
	}
}

func TestGoTokenizerNumbers(t *testing.T) {
	tokenizer := NewGoTokenizer()

	tests := []string{"123", "3.14", "0x1F", "1e10"}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test, 0)

		found := false
		for _, token := range tokens {
			if token.Type == TokenNumber {
				if token.Text != test {
					t.Errorf("Expected number %s, got %s", test, token.Text)
				}
				found = true
				break
			}
		}

		if !found {
			t.Errorf("No number token found for: %s", test)
		}
	}
}

func TestGoTokenizerFunctions(t *testing.T) {
	tokenizer := NewGoTokenizer()
	line := "fmt.Println(x)"

	tokens := tokenizer.Tokenize(line, 0)

	// Find function token
	found := false
	for _, token := range tokens {
		if token.Type == TokenFunction && token.Text == "Println" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find function token for 'Println'")
	}
}

func TestGoTokenizerTypes(t *testing.T) {
	tokenizer := NewGoTokenizer()
	line := "var x int"

	tokens := tokenizer.Tokenize(line, 0)

	// Find type token
	found := false
	for _, token := range tokens {
		if token.Type == TokenType && token.Text == "int" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find type token for 'int'")
	}
}

func TestGoTokenizerConstants(t *testing.T) {
	tokenizer := NewGoTokenizer()
	line := "if x == nil { return true }"

	tokens := tokenizer.Tokenize(line, 0)

	// Find constant tokens
	constants := []string{}
	for _, token := range tokens {
		if token.Type == TokenConstant {
			constants = append(constants, token.Text)
		}
	}

	expectedConstants := []string{"nil", "true"}
	if len(constants) != len(expectedConstants) {
		t.Errorf("Expected %d constants, got %d: %v", len(expectedConstants), len(constants), constants)
	}
}

func TestGoTokenizerOperators(t *testing.T) {
	tokenizer := NewGoTokenizer()
	line := "x := y + z"

	tokens := tokenizer.Tokenize(line, 0)

	// Find operator tokens
	operators := []string{}
	for _, token := range tokens {
		if token.Type == TokenOperator {
			operators = append(operators, token.Text)
		}
	}

	if len(operators) < 2 {
		t.Errorf("Expected at least 2 operators, got %d: %v", len(operators), operators)
	}
}

// --- Markdown Tokenizer Tests ---

func TestMarkdownTokenizerHeaders(t *testing.T) {
	tokenizer := NewMarkdownTokenizer()

	tests := []string{"# Header 1", "## Header 2", "### Header 3"}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test, 0)

		if len(tokens) == 0 {
			t.Errorf("No tokens found for: %s", test)
			continue
		}

		// First token should be keyword (header marker)
		if tokens[0].Type != TokenKeyword {
			t.Errorf("Expected keyword token for header marker, got %v", tokens[0].Type)
		}
	}
}

func TestMarkdownTokenizerCodeBlocks(t *testing.T) {
	tokenizer := NewMarkdownTokenizer()

	tests := []string{"```go", "    code", "\tcode"}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test, 0)

		if len(tokens) == 0 {
			t.Errorf("No tokens found for: %s", test)
			continue
		}

		// Should be treated as comment
		if tokens[0].Type != TokenComment {
			t.Errorf("Expected comment token for code block, got %v", tokens[0].Type)
		}
	}
}

func TestMarkdownTokenizerInlineCode(t *testing.T) {
	tokenizer := NewMarkdownTokenizer()
	line := "This is `code` here"

	tokens := tokenizer.Tokenize(line, 0)

	// Find inline code token
	found := false
	for _, token := range tokens {
		if token.Type == TokenString && token.Text == "`code`" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find inline code token")
	}
}

func TestMarkdownTokenizerBold(t *testing.T) {
	tokenizer := NewMarkdownTokenizer()
	line := "This is **bold** text"

	tokens := tokenizer.Tokenize(line, 0)

	// Find bold token
	found := false
	for _, token := range tokens {
		if token.Type == TokenKeyword && token.Text == "**bold**" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find bold token")
	}
}

func TestMarkdownTokenizerItalic(t *testing.T) {
	tokenizer := NewMarkdownTokenizer()
	line := "This is *italic* text"

	tokens := tokenizer.Tokenize(line, 0)

	// Find italic token
	found := false
	for _, token := range tokens {
		if token.Type == TokenConstant && token.Text == "*italic*" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find italic token")
	}
}

func TestMarkdownTokenizerLinks(t *testing.T) {
	tokenizer := NewMarkdownTokenizer()
	line := "Check [this link](https://example.com)"

	tokens := tokenizer.Tokenize(line, 0)

	// Find link token
	found := false
	for _, token := range tokens {
		if token.Type == TokenFunction && token.Text == "[this link](https://example.com)" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find link token")
	}
}

func TestMarkdownTokenizerLists(t *testing.T) {
	tokenizer := NewMarkdownTokenizer()

	tests := []string{"- Item 1", "* Item 2", "+ Item 3"}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test, 0)

		if len(tokens) == 0 {
			t.Errorf("No tokens found for: %s", test)
			continue
		}

		// First token should be operator (list marker)
		if tokens[0].Type != TokenOperator {
			t.Errorf("Expected operator token for list marker, got %v", tokens[0].Type)
		}
	}
}

// --- Language Detector Tests ---

func TestDetectLanguageByExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected Language
	}{
		{"main.go", LangGo},
		{"README.md", LangMarkdown},
		{"script.py", LangPython},
		{"app.js", LangJavaScript},
		{"component.tsx", LangTypeScript},
		{"lib.rs", LangRust},
		{"config.json", LangJSON},
		{"unknown.xyz", LangNone},
	}

	for _, test := range tests {
		result := DetectLanguage(test.path)
		if result != test.expected {
			t.Errorf("For path %s: expected %v, got %v", test.path, test.expected, result)
		}
	}
}

func TestDetectLanguageByFilename(t *testing.T) {
	tests := []struct {
		path     string
		expected Language
	}{
		{"go.mod", LangGo},
		{"go.sum", LangGo},
		{"README", LangMarkdown},
		{"package.json", LangJSON},
	}

	for _, test := range tests {
		result := DetectLanguage(test.path)
		if result != test.expected {
			t.Errorf("For path %s: expected %v, got %v", test.path, test.expected, result)
		}
	}
}

func TestDetectLanguageFromContent(t *testing.T) {
	tests := []struct {
		content  string
		expected Language
	}{
		{"#!/usr/bin/env python", LangPython},
		{"package main\n\nfunc main() {}", LangGo},
		{"# Header\n\n## Subheader", LangMarkdown},
		{`{"key": "value"}`, LangJSON},
		{"", LangNone},
	}

	for _, test := range tests {
		result := DetectLanguageFromContent(test.content)
		if result != test.expected {
			t.Errorf("For content: expected %v, got %v", test.expected, result)
		}
	}
}

// --- Highlighter Tests ---

func TestHighlighterLine(t *testing.T) {
	scheme := DefaultStyleScheme()
	highlighter := NewHighlighter(LangGo, scheme)

	line := "func main() {}"
	highlighted := highlighter.HighlightLine(line, 0)

	// Should contain some styling (ANSI codes)
	if highlighted == line {
		t.Error("Expected styled output, got plain text")
	}

	// Result should be longer due to ANSI codes
	if len(highlighted) <= len(line) {
		t.Errorf("Expected highlighted line to be longer, got %d <= %d", len(highlighted), len(line))
	}
}

func TestHighlighterBuffer(t *testing.T) {
	scheme := DefaultStyleScheme()
	highlighter := NewHighlighter(LangGo, scheme)

	buf := buffer.NewBuffer("test")
	buf.Insert(buffer.Position{Line: 0, Column: 0}, "package main\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}")

	highlighted := highlighter.HighlightBuffer(buf)

	if len(highlighted) != buf.LineCount() {
		t.Errorf("Expected %d highlighted lines, got %d", buf.LineCount(), len(highlighted))
	}

	// First line should be highlighted
	if highlighted[0] == "package main" {
		t.Error("Expected first line to be highlighted")
	}
}

func TestHighlighterNoTokenizer(t *testing.T) {
	scheme := DefaultStyleScheme()
	highlighter := NewHighlighter(LangNone, scheme)

	line := "plain text"
	highlighted := highlighter.HighlightLine(line, 0)

	// With no tokenizer, should return original text
	if highlighted != line {
		t.Errorf("Expected plain text, got %s", highlighted)
	}
}

// --- GetTokenizer Tests ---

func TestGetTokenizer(t *testing.T) {
	tests := []struct {
		lang     Language
		expected bool // whether tokenizer should exist
	}{
		{LangGo, true},
		{LangMarkdown, true},
		{LangPython, false},    // Not implemented yet
		{LangJavaScript, false}, // Not implemented yet
		{LangNone, false},
	}

	for _, test := range tests {
		tokenizer := GetTokenizer(test.lang)
		exists := tokenizer != nil

		if exists != test.expected {
			t.Errorf("For language %v: expected tokenizer exists=%v, got %v", test.lang, test.expected, exists)
		}

		if exists && tokenizer.Language() != test.lang {
			t.Errorf("Tokenizer returned wrong language: expected %v, got %v", test.lang, tokenizer.Language())
		}
	}
}
