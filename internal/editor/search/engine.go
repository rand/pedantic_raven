package search

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/rand/pedantic-raven/internal/editor/buffer"
)

// SimpleEngine is a basic search engine implementation.
type SimpleEngine struct{}

// NewEngine creates a new search engine.
func NewEngine() Engine {
	return &SimpleEngine{}
}

// Search implements Engine.
func (e *SimpleEngine) Search(buf buffer.Buffer, query string, opts SearchOptions) (*SearchResult, error) {
	if query == "" {
		return &SearchResult{Query: query, Matches: []Match{}, Options: opts}, nil
	}

	var matches []Match
	content := buf.Content()

	if opts.Regex {
		// Regex search
		var re *regexp.Regexp
		var err error
		if opts.CaseSensitive {
			re, err = regexp.Compile(query)
		} else {
			re, err = regexp.Compile("(?i)" + query)
		}
		if err != nil {
			return nil, fmt.Errorf("invalid regex: %w", err)
		}

		matches = e.findRegexMatches(buf, re, opts)
	} else {
		// Plain text search
		matches = e.findTextMatches(buf, content, query, opts)
	}

	return &SearchResult{
		Query:   query,
		Matches: matches,
		Options: opts,
	}, nil
}

// findTextMatches finds all plain text matches.
func (e *SimpleEngine) findTextMatches(buf buffer.Buffer, content, query string, opts SearchOptions) []Match {
	var matches []Match
	searchQuery := query
	searchContent := content

	// Apply case sensitivity
	if !opts.CaseSensitive {
		searchQuery = strings.ToLower(query)
		searchContent = strings.ToLower(content)
	}

	// Find all matches
	index := 0
	offset := 0
	for {
		pos := strings.Index(searchContent[offset:], searchQuery)
		if pos == -1 {
			break
		}

		actualPos := offset + pos

		// Check whole word matching
		if opts.WholeWord {
			if !e.isWholeWordMatch(content, actualPos, len(query)) {
				offset = actualPos + 1
				continue
			}
		}

		// Convert byte offset to line/column position
		start := e.offsetToPosition(buf, actualPos)
		end := e.offsetToPosition(buf, actualPos+len(query))

		matches = append(matches, Match{
			Start:  start,
			End:    end,
			Text:   content[actualPos : actualPos+len(query)],
			Index:  index,
			Length: len(query),
		})

		index++
		offset = actualPos + 1
	}

	return matches
}

// findRegexMatches finds all regex matches.
func (e *SimpleEngine) findRegexMatches(buf buffer.Buffer, re *regexp.Regexp, opts SearchOptions) []Match {
	var matches []Match
	content := buf.Content()

	allMatches := re.FindAllStringIndex(content, -1)
	for i, match := range allMatches {
		start := e.offsetToPosition(buf, match[0])
		end := e.offsetToPosition(buf, match[1])

		matches = append(matches, Match{
			Start:  start,
			End:    end,
			Text:   content[match[0]:match[1]],
			Index:  i,
			Length: match[1] - match[0],
		})
	}

	return matches
}

// isWholeWordMatch checks if a match is a whole word.
func (e *SimpleEngine) isWholeWordMatch(content string, pos, length int) bool {
	// Check character before match
	if pos > 0 {
		charBefore := rune(content[pos-1])
		if unicode.IsLetter(charBefore) || unicode.IsDigit(charBefore) || charBefore == '_' {
			return false
		}
	}

	// Check character after match
	if pos+length < len(content) {
		charAfter := rune(content[pos+length])
		if unicode.IsLetter(charAfter) || unicode.IsDigit(charAfter) || charAfter == '_' {
			return false
		}
	}

	return true
}

// offsetToPosition converts a byte offset to a buffer position.
func (e *SimpleEngine) offsetToPosition(buf buffer.Buffer, offset int) buffer.Position {
	content := buf.Content()
	if offset < 0 {
		offset = 0
	}
	if offset > len(content) {
		offset = len(content)
	}

	line := 0
	column := 0
	currentOffset := 0

	for i, char := range content {
		if currentOffset >= offset {
			break
		}
		if char == '\n' {
			line++
			column = 0
		} else {
			column++
		}
		currentOffset = i + 1
	}

	return buffer.Position{Line: line, Column: column}
}

// positionToOffset converts a buffer position to a byte offset.
func (e *SimpleEngine) positionToOffset(buf buffer.Buffer, pos buffer.Position) int {
	offset := 0
	for i := 0; i < pos.Line && i < buf.LineCount(); i++ {
		offset += len(buf.Line(i)) + 1 // +1 for newline
	}
	if pos.Line < buf.LineCount() {
		line := buf.Line(pos.Line)
		if pos.Column < len(line) {
			offset += pos.Column
		} else {
			offset += len(line)
		}
	}
	return offset
}

// FindNext implements Engine.
func (e *SimpleEngine) FindNext(buf buffer.Buffer, query string, after buffer.Position, opts SearchOptions) (*Match, error) {
	result, err := e.Search(buf, query, opts)
	if err != nil {
		return nil, err
	}

	if len(result.Matches) == 0 {
		return nil, nil
	}

	afterOffset := e.positionToOffset(buf, after)

	// Find first match after the position
	for i := range result.Matches {
		matchOffset := e.positionToOffset(buf, result.Matches[i].Start)
		if matchOffset > afterOffset {
			return &result.Matches[i], nil
		}
	}

	// Wrap around if enabled
	if opts.WrapAround && len(result.Matches) > 0 {
		return &result.Matches[0], nil
	}

	return nil, nil
}

// FindPrevious implements Engine.
func (e *SimpleEngine) FindPrevious(buf buffer.Buffer, query string, before buffer.Position, opts SearchOptions) (*Match, error) {
	result, err := e.Search(buf, query, opts)
	if err != nil {
		return nil, err
	}

	if len(result.Matches) == 0 {
		return nil, nil
	}

	beforeOffset := e.positionToOffset(buf, before)

	// Find last match before the position
	for i := len(result.Matches) - 1; i >= 0; i-- {
		matchOffset := e.positionToOffset(buf, result.Matches[i].Start)
		if matchOffset < beforeOffset {
			return &result.Matches[i], nil
		}
	}

	// Wrap around if enabled
	if opts.WrapAround && len(result.Matches) > 0 {
		return &result.Matches[len(result.Matches)-1], nil
	}

	return nil, nil
}

// Replace implements Engine.
func (e *SimpleEngine) Replace(buf buffer.Buffer, match Match, replacement string) error {
	// Use buffer's Replace method to support undo/redo
	_, err := buf.Replace(match.Start, match.End, replacement)
	return err
}

// ReplaceAll implements Engine.
func (e *SimpleEngine) ReplaceAll(buf buffer.Buffer, query string, replacement string, opts SearchOptions) (int, error) {
	result, err := e.Search(buf, query, opts)
	if err != nil {
		return 0, err
	}

	// Replace in reverse order to maintain position validity
	count := 0
	for i := len(result.Matches) - 1; i >= 0; i-- {
		match := result.Matches[i]
		err := e.Replace(buf, match, replacement)
		if err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}
