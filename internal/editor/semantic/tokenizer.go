package semantic

import (
	"regexp"
	"strings"
	"unicode"
)

// SimpleTokenizer is a basic tokenizer implementation.
type SimpleTokenizer struct{}

// NewTokenizer creates a new simple tokenizer.
func NewTokenizer() Tokenizer {
	return &SimpleTokenizer{}
}

// Tokenize implements Tokenizer.
func (t *SimpleTokenizer) Tokenize(content string) []Token {
	var tokens []Token
	line := 0
	offset := 0

	// Process line by line
	lines := strings.Split(content, "\n")

	for lineNum, lineText := range lines {
		lineOffset := 0

		for lineOffset < len(lineText) {
			// Skip whitespace
			if unicode.IsSpace(rune(lineText[lineOffset])) && lineText[lineOffset] != '\n' {
				start := lineOffset
				for lineOffset < len(lineText) && unicode.IsSpace(rune(lineText[lineOffset])) && lineText[lineOffset] != '\n' {
					lineOffset++
				}

				tokens = append(tokens, Token{
					Type: TokenWhitespace,
					Text: lineText[start:lineOffset],
					Span: Span{
						Start: offset + start,
						End:   offset + lineOffset,
						Line:  lineNum,
					},
				})
				continue
			}

			// Typed hole: ??Type
			if lineOffset+1 < len(lineText) && lineText[lineOffset:lineOffset+2] == "??" {
				start := lineOffset
				lineOffset += 2

				// Extract type name
				for lineOffset < len(lineText) && (unicode.IsLetter(rune(lineText[lineOffset])) || unicode.IsDigit(rune(lineText[lineOffset]))) {
					lineOffset++
				}

				tokens = append(tokens, Token{
					Type:  TokenTypeHole,
					Text:  lineText[start:lineOffset],
					Value: lineText[start+2 : lineOffset], // Type without ??
					Span: Span{
						Start: offset + start,
						End:   offset + lineOffset,
						Line:  lineNum,
					},
				})
				continue
			}

			// Constraint hole: !!constraint
			if lineOffset+1 < len(lineText) && lineText[lineOffset:lineOffset+2] == "!!" {
				start := lineOffset
				lineOffset += 2

				// Extract constraint
				for lineOffset < len(lineText) && !unicode.IsSpace(rune(lineText[lineOffset])) {
					lineOffset++
				}

				tokens = append(tokens, Token{
					Type:  TokenConstraintHole,
					Text:  lineText[start:lineOffset],
					Value: lineText[start+2 : lineOffset], // Constraint without !!
					Span: Span{
						Start: offset + start,
						End:   offset + lineOffset,
						Line:  lineNum,
					},
				})
				continue
			}

			// Numbers
			if unicode.IsDigit(rune(lineText[lineOffset])) {
				start := lineOffset
				for lineOffset < len(lineText) && (unicode.IsDigit(rune(lineText[lineOffset])) || lineText[lineOffset] == '.') {
					lineOffset++
				}

				tokens = append(tokens, Token{
					Type: TokenNumber,
					Text: lineText[start:lineOffset],
					Span: Span{
						Start: offset + start,
						End:   offset + lineOffset,
						Line:  lineNum,
					},
				})
				continue
			}

			// Words (including capitalized and proper nouns)
			if unicode.IsLetter(rune(lineText[lineOffset])) {
				start := lineOffset
				for lineOffset < len(lineText) && (unicode.IsLetter(rune(lineText[lineOffset])) || lineText[lineOffset] == '-' || lineText[lineOffset] == '_') {
					lineOffset++
				}

				word := lineText[start:lineOffset]
				tokenType := classifyWord(word)

				tokens = append(tokens, Token{
					Type:  tokenType,
					Text:  word,
					Value: strings.ToLower(word),
					Span: Span{
						Start: offset + start,
						End:   offset + lineOffset,
						Line:  lineNum,
					},
				})
				continue
			}

			// Punctuation
			if !unicode.IsSpace(rune(lineText[lineOffset])) {
				tokens = append(tokens, Token{
					Type: TokenPunctuation,
					Text: string(lineText[lineOffset]),
					Span: Span{
						Start: offset + lineOffset,
						End:   offset + lineOffset + 1,
						Line:  lineNum,
					},
				})
				lineOffset++
				continue
			}

			// Fallback: skip character
			lineOffset++
		}

		// Add newline token
		if lineNum < len(lines)-1 {
			tokens = append(tokens, Token{
				Type: TokenNewline,
				Text: "\n",
				Span: Span{
					Start: offset + len(lineText),
					End:   offset + len(lineText) + 1,
					Line:  lineNum,
				},
			})
		}

		offset += len(lineText) + 1 // +1 for newline
		line++
	}

	return tokens
}

// classifyWord determines the token type for a word.
func classifyWord(word string) TokenType {
	if len(word) == 0 {
		return TokenUnknown
	}

	// Check for all uppercase (acronym or proper noun)
	if isAllUpper(word) && len(word) > 1 {
		return TokenProperNoun
	}

	// Check for Title Case (first letter upper, rest lower)
	if unicode.IsUpper(rune(word[0])) {
		if isTitleCase(word) {
			return TokenProperNoun
		}
		return TokenCapitalizedWord
	}

	// Check if it's a common verb
	if isCommonVerb(word) {
		return TokenVerb
	}

	// Default to regular word
	return TokenWord
}

// isAllUpper checks if all letters in a word are uppercase.
func isAllUpper(word string) bool {
	for _, ch := range word {
		if unicode.IsLetter(ch) && !unicode.IsUpper(ch) {
			return false
		}
	}
	return true
}

// isTitleCase checks if a word is in Title Case.
func isTitleCase(word string) bool {
	if len(word) < 2 {
		return false
	}

	if !unicode.IsUpper(rune(word[0])) {
		return false
	}

	// Rest should be lowercase
	for i := 1; i < len(word); i++ {
		if unicode.IsUpper(rune(word[i])) {
			return false
		}
	}

	return true
}

// isCommonVerb checks if a word is a common verb.
// This is a simple heuristic and can be expanded.
func isCommonVerb(word string) bool {
	commonVerbs := map[string]bool{
		"is": true, "are": true, "was": true, "were": true,
		"has": true, "have": true, "had": true,
		"does": true, "do": true, "did": true,
		"can": true, "could": true, "will": true, "would": true,
		"should": true, "may": true, "might": true, "must": true,
		"gets": true, "get": true, "got": true,
		"makes": true, "make": true, "made": true,
		"takes": true, "take": true, "took": true,
		"gives": true, "give": true, "gave": true,
		"uses": true, "use": true, "used": true,
		"creates": true, "create": true, "created": true,
		"provides": true, "provide": true, "provided": true,
		"implements": true, "implement": true, "implemented": true,
		"returns": true, "return": true, "returned": true,
		"calls": true, "call": true, "called": true,
		"runs": true, "run": true, "ran": true,
		"executes": true, "execute": true, "executed": true,
		"validates": true, "validate": true, "validated": true,
		"processes": true, "process": true, "processed": true,
		"handles": true, "handle": true, "handled": true,
		"manages": true, "manage": true, "managed": true,
	}

	return commonVerbs[strings.ToLower(word)]
}

// ExtractWords extracts just the words from tokens.
func ExtractWords(tokens []Token) []string {
	var words []string
	for _, token := range tokens {
		if token.Type == TokenWord || token.Type == TokenCapitalizedWord || token.Type == TokenProperNoun {
			words = append(words, token.Text)
		}
	}
	return words
}

// ExtractTypedHoles extracts typed holes from tokens.
func ExtractTypedHoles(tokens []Token) []TypedHole {
	var holes []TypedHole
	for _, token := range tokens {
		if token.Type == TokenTypeHole {
			holes = append(holes, TypedHole{
				Type: token.Value,
				Span: token.Span,
			})
		} else if token.Type == TokenConstraintHole {
			holes = append(holes, TypedHole{
				Constraint: token.Value,
				Span:       token.Span,
			})
		}
	}
	return holes
}

// Common dependency patterns
var dependencyPatterns = []*regexp.Regexp{
	regexp.MustCompile(`import\s+"([^"]+)"`),      // Go imports
	regexp.MustCompile(`import\s+(\S+)`),          // Generic imports
	regexp.MustCompile(`require\s+['"]([^'"]+)`),  // Node.js requires
	regexp.MustCompile(`from\s+['"]([^'"]+)`),     // Python/ES6 imports
	regexp.MustCompile(`#include\s+<([^>]+)>`),    // C/C++ includes
	regexp.MustCompile(`use\s+(\S+);`),            // Rust use statements
}

// ExtractDependencies extracts dependencies from content.
func ExtractDependencies(content string) []Dependency {
	var deps []Dependency

	for _, pattern := range dependencyPatterns {
		matches := pattern.FindAllStringSubmatchIndex(content, -1)
		for _, match := range matches {
			if len(match) >= 4 {
				target := content[match[2]:match[3]]
				deps = append(deps, Dependency{
					Type:   "import",
					Target: target,
					Span: Span{
						Start: match[0],
						End:   match[1],
					},
				})
			}
		}
	}

	return deps
}
