package syntax

import (
	"strings"
	"unicode"
)

// GoTokenizer implements syntax highlighting for Go language.
type GoTokenizer struct{}

// NewGoTokenizer creates a new Go tokenizer.
func NewGoTokenizer() *GoTokenizer {
	return &GoTokenizer{}
}

// Language implements Tokenizer.
func (g *GoTokenizer) Language() Language {
	return LangGo
}

// Go keywords
var goKeywords = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"defer":       true,
	"else":        true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

// Go builtin types
var goTypes = map[string]bool{
	"bool":       true,
	"byte":       true,
	"complex64":  true,
	"complex128": true,
	"error":      true,
	"float32":    true,
	"float64":    true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"rune":       true,
	"string":     true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
}

// Go builtin constants
var goConstants = map[string]bool{
	"true":  true,
	"false": true,
	"nil":   true,
	"iota":  true,
}

// Tokenize implements Tokenizer.
func (g *GoTokenizer) Tokenize(line string, lineNum int) []Token {
	var tokens []Token
	i := 0

	for i < len(line) {
		// Skip whitespace
		if unicode.IsSpace(rune(line[i])) {
			start := i
			for i < len(line) && unicode.IsSpace(rune(line[i])) {
				i++
			}
			tokens = append(tokens, Token{
				Type:  TokenWhitespace,
				Start: start,
				End:   i,
				Text:  line[start:i],
				Line:  lineNum,
			})
			continue
		}

		// Single-line comment
		if i+1 < len(line) && line[i:i+2] == "//" {
			tokens = append(tokens, Token{
				Type:  TokenComment,
				Start: i,
				End:   len(line),
				Text:  line[i:],
				Line:  lineNum,
			})
			break
		}

		// Multi-line comment start (we only handle single lines, so treat as comment)
		if i+1 < len(line) && line[i:i+2] == "/*" {
			end := strings.Index(line[i+2:], "*/")
			if end != -1 {
				end = i + 2 + end + 2
			} else {
				end = len(line)
			}
			tokens = append(tokens, Token{
				Type:  TokenComment,
				Start: i,
				End:   end,
				Text:  line[i:end],
				Line:  lineNum,
			})
			i = end
			continue
		}

		// String literals (double quotes)
		if line[i] == '"' {
			start := i
			i++ // Skip opening quote
			for i < len(line) && line[i] != '"' {
				if line[i] == '\\' && i+1 < len(line) {
					i += 2 // Skip escaped character
				} else {
					i++
				}
			}
			if i < len(line) {
				i++ // Skip closing quote
			}
			tokens = append(tokens, Token{
				Type:  TokenString,
				Start: start,
				End:   i,
				Text:  line[start:i],
				Line:  lineNum,
			})
			continue
		}

		// Raw string literals (backticks)
		if line[i] == '`' {
			start := i
			i++ // Skip opening backtick
			for i < len(line) && line[i] != '`' {
				i++
			}
			if i < len(line) {
				i++ // Skip closing backtick
			}
			tokens = append(tokens, Token{
				Type:  TokenString,
				Start: start,
				End:   i,
				Text:  line[start:i],
				Line:  lineNum,
			})
			continue
		}

		// Rune literals
		if line[i] == '\'' {
			start := i
			i++ // Skip opening quote
			for i < len(line) && line[i] != '\'' {
				if line[i] == '\\' && i+1 < len(line) {
					i += 2 // Skip escaped character
				} else {
					i++
				}
			}
			if i < len(line) {
				i++ // Skip closing quote
			}
			tokens = append(tokens, Token{
				Type:  TokenString,
				Start: start,
				End:   i,
				Text:  line[start:i],
				Line:  lineNum,
			})
			continue
		}

		// Numbers
		if unicode.IsDigit(rune(line[i])) {
			start := i
			// Handle hex numbers (0x...)
			if i+1 < len(line) && line[i] == '0' && (line[i+1] == 'x' || line[i+1] == 'X') {
				i += 2
				for i < len(line) && (unicode.IsDigit(rune(line[i])) || (line[i] >= 'a' && line[i] <= 'f') || (line[i] >= 'A' && line[i] <= 'F')) {
					i++
				}
			} else {
				// Handle decimal/float numbers
				for i < len(line) && (unicode.IsDigit(rune(line[i])) || line[i] == '.' || line[i] == 'e' || line[i] == 'E') {
					i++
				}
			}
			tokens = append(tokens, Token{
				Type:  TokenNumber,
				Start: start,
				End:   i,
				Text:  line[start:i],
				Line:  lineNum,
			})
			continue
		}

		// Identifiers and keywords
		if unicode.IsLetter(rune(line[i])) || line[i] == '_' {
			start := i
			for i < len(line) && (unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i])) || line[i] == '_') {
				i++
			}
			word := line[start:i]

			// Check if it's followed by '(' to identify function calls
			nextNonSpace := i
			for nextNonSpace < len(line) && unicode.IsSpace(rune(line[nextNonSpace])) {
				nextNonSpace++
			}
			isFunction := nextNonSpace < len(line) && line[nextNonSpace] == '('

			// Determine token type
			tokenType := TokenIdentifier
			if goKeywords[word] {
				tokenType = TokenKeyword
			} else if goTypes[word] {
				tokenType = TokenTypeName
			} else if goConstants[word] {
				tokenType = TokenConstant
			} else if isFunction {
				tokenType = TokenFunction
			} else if unicode.IsUpper(rune(word[0])) {
				// Exported identifiers might be types
				tokenType = TokenTypeName
			}

			tokens = append(tokens, Token{
				Type:  tokenType,
				Start: start,
				End:   i,
				Text:  word,
				Line:  lineNum,
			})
			continue
		}

		// Operators and punctuation
		if isOperatorOrPunctuation(line[i]) {
			start := i
			isMultiChar := false
			// Handle multi-character operators
			if i+1 < len(line) && isMultiCharOperator(line[i:i+2]) {
				i += 2
				isMultiChar = true
			} else {
				i++
			}

			tokenType := TokenOperator
			// Only treat as punctuation if it's a single character AND is in punctuation set
			if !isMultiChar && isPunctuation(line[start]) {
				tokenType = TokenPunctuation
			}

			tokens = append(tokens, Token{
				Type:  tokenType,
				Start: start,
				End:   i,
				Text:  line[start:i],
				Line:  lineNum,
			})
			continue
		}

		// Unknown character, treat as punctuation
		tokens = append(tokens, Token{
			Type:  TokenPunctuation,
			Start: i,
			End:   i + 1,
			Text:  string(line[i]),
			Line:  lineNum,
		})
		i++
	}

	return tokens
}

// isOperatorOrPunctuation checks if a character is an operator or punctuation.
func isOperatorOrPunctuation(c byte) bool {
	return strings.ContainsRune("+-*/%&|^!<>=:;,.()[]{}#", rune(c))
}

// isPunctuation checks if a character is punctuation (not an operator).
func isPunctuation(c byte) bool {
	return strings.ContainsRune(";,.()[]{}:", rune(c))
}

// isMultiCharOperator checks if a two-character sequence is a multi-char operator.
func isMultiCharOperator(s string) bool {
	multiChar := map[string]bool{
		"==": true,
		"!=": true,
		"<=": true,
		">=": true,
		"&&": true,
		"||": true,
		"++": true,
		"--": true,
		"<<": true,
		">>": true,
		"&^": true,
		":=": true,
		"->": true,
		"<-": true,
	}
	return multiChar[s]
}
