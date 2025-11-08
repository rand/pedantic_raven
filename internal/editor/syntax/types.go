// Package syntax provides syntax highlighting for various programming languages.
package syntax

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/rand/pedantic-raven/internal/editor/buffer"
)

// TokenType represents the type of a syntax token.
type TokenType int

const (
	TokenNone TokenType = iota
	TokenKeyword
	TokenString
	TokenComment
	TokenNumber
	TokenOperator
	TokenIdentifier
	TokenFunction
	TokenTypeName
	TokenConstant
	TokenPunctuation
	TokenWhitespace
)

// String returns the string representation of a token type.
func (t TokenType) String() string {
	switch t {
	case TokenKeyword:
		return "keyword"
	case TokenString:
		return "string"
	case TokenComment:
		return "comment"
	case TokenNumber:
		return "number"
	case TokenOperator:
		return "operator"
	case TokenIdentifier:
		return "identifier"
	case TokenFunction:
		return "function"
	case TokenTypeName:
		return "type"
	case TokenConstant:
		return "constant"
	case TokenPunctuation:
		return "punctuation"
	case TokenWhitespace:
		return "whitespace"
	default:
		return "none"
	}
}

// Token represents a syntax token with position and type.
type Token struct {
	Type   TokenType
	Start  int // Byte offset in line
	End    int // Byte offset in line (exclusive)
	Text   string
	Line   int // Line number
}

// Language represents a supported programming language.
type Language int

const (
	LangNone Language = iota
	LangGo
	LangMarkdown
	LangPython
	LangJavaScript
	LangTypeScript
	LangRust
	LangJSON
)

// String returns the string representation of a language.
func (l Language) String() string {
	switch l {
	case LangGo:
		return "go"
	case LangMarkdown:
		return "markdown"
	case LangPython:
		return "python"
	case LangJavaScript:
		return "javascript"
	case LangTypeScript:
		return "typescript"
	case LangRust:
		return "rust"
	case LangJSON:
		return "json"
	default:
		return "none"
	}
}

// Tokenizer is the interface that all language tokenizers must implement.
type Tokenizer interface {
	// Tokenize tokenizes a single line of text.
	// Returns a slice of tokens for that line.
	Tokenize(line string, lineNum int) []Token

	// Language returns the language this tokenizer handles.
	Language() Language
}

// StyleScheme defines color styles for token types.
type StyleScheme struct {
	Keyword     lipgloss.Style
	String      lipgloss.Style
	Comment     lipgloss.Style
	Number      lipgloss.Style
	Operator    lipgloss.Style
	Identifier  lipgloss.Style
	Function    lipgloss.Style
	Type        lipgloss.Style
	Constant    lipgloss.Style
	Punctuation lipgloss.Style
	Default     lipgloss.Style
}

// DefaultStyleScheme returns a default color scheme for syntax highlighting.
func DefaultStyleScheme() StyleScheme {
	return StyleScheme{
		Keyword:     lipgloss.NewStyle().Foreground(lipgloss.Color("205")), // Magenta
		String:      lipgloss.NewStyle().Foreground(lipgloss.Color("107")), // Green
		Comment:     lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Italic(true), // Gray
		Number:      lipgloss.NewStyle().Foreground(lipgloss.Color("141")), // Purple
		Operator:    lipgloss.NewStyle().Foreground(lipgloss.Color("208")), // Orange
		Identifier:  lipgloss.NewStyle().Foreground(lipgloss.Color("252")), // Light gray
		Function:    lipgloss.NewStyle().Foreground(lipgloss.Color("75")),  // Blue
		Type:        lipgloss.NewStyle().Foreground(lipgloss.Color("114")), // Teal
		Constant:    lipgloss.NewStyle().Foreground(lipgloss.Color("179")), // Yellow
		Punctuation: lipgloss.NewStyle().Foreground(lipgloss.Color("246")), // Medium gray
		Default:     lipgloss.NewStyle().Foreground(lipgloss.Color("252")), // Light gray
	}
}

// Highlighter applies syntax highlighting to buffer content.
type Highlighter struct {
	tokenizer Tokenizer
	scheme    StyleScheme
	language  Language
}

// NewHighlighter creates a new syntax highlighter.
func NewHighlighter(lang Language, scheme StyleScheme) *Highlighter {
	tokenizer := GetTokenizer(lang)
	return &Highlighter{
		tokenizer: tokenizer,
		scheme:    scheme,
		language:  lang,
	}
}

// HighlightLine applies syntax highlighting to a single line.
// Returns the styled line as a string.
func (h *Highlighter) HighlightLine(line string, lineNum int) string {
	if h.tokenizer == nil {
		return line
	}

	tokens := h.tokenizer.Tokenize(line, lineNum)
	if len(tokens) == 0 {
		return line
	}

	// Build highlighted string from tokens
	result := ""
	lastEnd := 0

	for _, token := range tokens {
		// Add any unhighlighted text before this token
		if token.Start > lastEnd {
			result += line[lastEnd:token.Start]
		}

		// Apply style to token
		style := h.styleForToken(token.Type)
		result += style.Render(token.Text)

		lastEnd = token.End
	}

	// Add any remaining text
	if lastEnd < len(line) {
		result += line[lastEnd:]
	}

	return result
}

// HighlightBuffer applies syntax highlighting to a buffer.
// Returns a slice of highlighted lines.
func (h *Highlighter) HighlightBuffer(buf buffer.Buffer) []string {
	lines := buf.Lines()
	highlighted := make([]string, len(lines))

	for i, line := range lines {
		highlighted[i] = h.HighlightLine(line, i)
	}

	return highlighted
}

// styleForToken returns the appropriate style for a token type.
func (h *Highlighter) styleForToken(tokenType TokenType) lipgloss.Style {
	switch tokenType {
	case TokenKeyword:
		return h.scheme.Keyword
	case TokenString:
		return h.scheme.String
	case TokenComment:
		return h.scheme.Comment
	case TokenNumber:
		return h.scheme.Number
	case TokenOperator:
		return h.scheme.Operator
	case TokenIdentifier:
		return h.scheme.Identifier
	case TokenFunction:
		return h.scheme.Function
	case TokenTypeName:
		return h.scheme.Type
	case TokenConstant:
		return h.scheme.Constant
	case TokenPunctuation:
		return h.scheme.Punctuation
	default:
		return h.scheme.Default
	}
}

// GetTokenizer returns the appropriate tokenizer for a language.
func GetTokenizer(lang Language) Tokenizer {
	switch lang {
	case LangGo:
		return NewGoTokenizer()
	case LangMarkdown:
		return NewMarkdownTokenizer()
	default:
		return nil
	}
}
