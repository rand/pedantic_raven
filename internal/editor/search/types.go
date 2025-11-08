// Package search provides text search and replace functionality for the editor.
package search

import (
	"github.com/rand/pedantic-raven/internal/editor/buffer"
)

// SearchOptions configures search behavior.
type SearchOptions struct {
	CaseSensitive bool
	WholeWord     bool
	Regex         bool
	WrapAround    bool
}

// DefaultSearchOptions returns the default search configuration.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		CaseSensitive: false,
		WholeWord:     false,
		Regex:         false,
		WrapAround:    true,
	}
}

// Match represents a single search match in the buffer.
type Match struct {
	Start  buffer.Position // Start position of the match
	End    buffer.Position // End position of the match
	Text   string          // Matched text
	Index  int             // Match index (0-based)
	Length int             // Length of matched text
}

// SearchResult contains all matches found for a query.
type SearchResult struct {
	Query   string
	Matches []Match
	Options SearchOptions
}

// ReplaceOptions configures replace behavior.
type ReplaceOptions struct {
	ReplaceAll bool
}

// Engine provides search and replace functionality for buffers.
type Engine interface {
	// Search finds all occurrences of the query in the buffer.
	Search(buf buffer.Buffer, query string, opts SearchOptions) (*SearchResult, error)

	// FindNext finds the next match after the given position.
	FindNext(buf buffer.Buffer, query string, after buffer.Position, opts SearchOptions) (*Match, error)

	// FindPrevious finds the previous match before the given position.
	FindPrevious(buf buffer.Buffer, query string, before buffer.Position, opts SearchOptions) (*Match, error)

	// Replace replaces a specific match with the replacement text.
	Replace(buf buffer.Buffer, match Match, replacement string) error

	// ReplaceAll replaces all matches with the replacement text.
	ReplaceAll(buf buffer.Buffer, query string, replacement string, opts SearchOptions) (int, error)
}
