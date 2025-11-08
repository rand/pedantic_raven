package syntax

import (
	"strings"
	"unicode"
)

// MarkdownTokenizer implements syntax highlighting for Markdown.
type MarkdownTokenizer struct{}

// NewMarkdownTokenizer creates a new Markdown tokenizer.
func NewMarkdownTokenizer() *MarkdownTokenizer {
	return &MarkdownTokenizer{}
}

// Language implements Tokenizer.
func (m *MarkdownTokenizer) Language() Language {
	return LangMarkdown
}

// Tokenize implements Tokenizer.
func (m *MarkdownTokenizer) Tokenize(line string, lineNum int) []Token {
	var tokens []Token

	// Empty line
	if strings.TrimSpace(line) == "" {
		return tokens
	}

	i := 0

	// Headers (# ## ###)
	if line[0] == '#' {
		headerLevel := 0
		for i < len(line) && line[i] == '#' {
			headerLevel++
			i++
		}
		if i < len(line) && (line[i] == ' ' || line[i] == '\t') {
			// Valid header
			tokens = append(tokens, Token{
				Type:  TokenKeyword,
				Start: 0,
				End:   i,
				Text:  line[0:i],
				Line:  lineNum,
			})
			// Rest of line is the header text
			if i < len(line) {
				tokens = append(tokens, Token{
					Type:  TokenConstant, // Use constant for header text
					Start: i,
					End:   len(line),
					Text:  line[i:],
					Line:  lineNum,
				})
			}
			return tokens
		}
	}

	// Code blocks (``` or indented)
	if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
		tokens = append(tokens, Token{
			Type:  TokenComment, // Use comment color for code blocks
			Start: 0,
			End:   len(line),
			Text:  line,
			Line:  lineNum,
		})
		return tokens
	}

	// Lists (- * +)
	trimmed := strings.TrimSpace(line)
	if len(trimmed) > 0 && (trimmed[0] == '-' || trimmed[0] == '*' || trimmed[0] == '+') {
		if len(trimmed) > 1 && (trimmed[1] == ' ' || trimmed[1] == '\t') {
			// List marker
			markerEnd := strings.IndexAny(line, "-*+") + 1
			tokens = append(tokens, Token{
				Type:  TokenOperator,
				Start: 0,
				End:   markerEnd,
				Text:  line[0:markerEnd],
				Line:  lineNum,
			})
			i = markerEnd
		}
	}

	// Process inline elements
	for i < len(line) {
		// Bold (**text**)
		if i+1 < len(line) && line[i:i+2] == "**" {
			start := i
			i += 2
			end := strings.Index(line[i:], "**")
			if end != -1 {
				end = i + end + 2
				tokens = append(tokens, Token{
					Type:  TokenKeyword,
					Start: start,
					End:   end,
					Text:  line[start:end],
					Line:  lineNum,
				})
				i = end
				continue
			}
		}

		// Italic (*text* or _text_)
		if (line[i] == '*' || line[i] == '_') && i+1 < len(line) && !unicode.IsSpace(rune(line[i+1])) {
			marker := line[i]
			start := i
			i++
			for i < len(line) && line[i] != byte(marker) {
				i++
			}
			if i < len(line) {
				i++ // Include closing marker
				tokens = append(tokens, Token{
					Type:  TokenConstant,
					Start: start,
					End:   i,
					Text:  line[start:i],
					Line:  lineNum,
				})
				continue
			}
		}

		// Inline code (`code`)
		if line[i] == '`' {
			start := i
			i++
			for i < len(line) && line[i] != '`' {
				i++
			}
			if i < len(line) {
				i++ // Include closing backtick
				tokens = append(tokens, Token{
					Type:  TokenString,
					Start: start,
					End:   i,
					Text:  line[start:i],
					Line:  lineNum,
				})
				continue
			}
		}

		// Links [text](url)
		if line[i] == '[' {
			start := i
			i++
			// Find closing ]
			for i < len(line) && line[i] != ']' {
				i++
			}
			if i < len(line) && i+1 < len(line) && line[i+1] == '(' {
				i++ // Skip ]
				i++ // Skip (
				// Find closing )
				for i < len(line) && line[i] != ')' {
					i++
				}
				if i < len(line) {
					i++ // Include closing )
					tokens = append(tokens, Token{
						Type:  TokenFunction, // Use function color for links
						Start: start,
						End:   i,
						Text:  line[start:i],
						Line:  lineNum,
					})
					continue
				}
			}
		}

		// Regular text
		start := i
		for i < len(line) && line[i] != '*' && line[i] != '_' && line[i] != '`' && line[i] != '[' {
			i++
		}
		if i > start {
			tokens = append(tokens, Token{
				Type:  TokenIdentifier,
				Start: start,
				End:   i,
				Text:  line[start:i],
				Line:  lineNum,
			})
		}
	}

	return tokens
}
