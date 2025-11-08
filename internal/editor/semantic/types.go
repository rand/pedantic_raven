// Package semantic provides semantic analysis and triple extraction.
//
// The semantic analyzer processes text content to extract:
// - Entities (nouns, proper nouns, concepts)
// - Relationships (verb phrases connecting entities)
// - Typed holes (??Type placeholders for future implementation)
// - Dependencies (imports, references, requirements)
// - Triples (subject-predicate-object structures)
//
// Analysis is performed in a streaming fashion with progressive updates.
package semantic

import (
	"time"
)

// EntityType classifies entities.
type EntityType int

const (
	EntityUnknown EntityType = iota
	EntityPerson
	EntityPlace
	EntityThing
	EntityConcept
	EntityOrganization
	EntityTechnology
)

// String returns the string representation of the entity type.
func (t EntityType) String() string {
	switch t {
	case EntityPerson:
		return "Person"
	case EntityPlace:
		return "Place"
	case EntityThing:
		return "Thing"
	case EntityConcept:
		return "Concept"
	case EntityOrganization:
		return "Organization"
	case EntityTechnology:
		return "Technology"
	default:
		return "Unknown"
	}
}

// Span represents a text range.
type Span struct {
	Start int // Start offset in content
	End   int // End offset in content
	Line  int // Line number (0-indexed)
}

// Entity represents a detected entity.
type Entity struct {
	Text  string     // Entity text
	Type  EntityType // Entity type
	Span  Span       // Location in content
	Count int        // Number of occurrences
}

// Relationship represents a connection between entities.
type Relationship struct {
	Subject   string // Subject entity
	Predicate string // Relationship verb/action
	Object    string // Object entity
	Span      Span   // Location in content
}

// TypedHole represents a placeholder for future implementation.
//
// Typed holes follow the pattern:
// - ??Type: A hole of specific type needing implementation
// - !!constraint: A constraint that must be satisfied
type TypedHole struct {
	Type       string // Hole type (e.g., "Function", "Interface")
	Constraint string // Optional constraint
	Span       Span   // Location in content
}

// Dependency represents a reference to external resources.
type Dependency struct {
	Type   string // Dependency type (import, require, reference)
	Target string // Dependency target
	Span   Span   // Location in content
}

// Triple represents a subject-predicate-object structure.
type Triple struct {
	Subject   string // Subject
	Predicate string // Predicate (relationship)
	Object    string // Object
	Source    Span   // Source location
}

// Analysis holds the complete semantic analysis results.
type Analysis struct {
	Content       string         // Original content
	Entities      []Entity       // Extracted entities
	Relationships []Relationship // Detected relationships
	TypedHoles    []TypedHole    // Typed holes
	Dependencies  []Dependency   // Dependencies
	Triples       []Triple       // Extracted triples
	Timestamp     time.Time      // Analysis timestamp
	Duration      time.Duration  // Analysis duration
}

// UpdateType indicates the type of analysis update.
type UpdateType int

const (
	UpdateIncremental UpdateType = iota // Incremental progress
	UpdateComplete                      // Analysis complete
	UpdateError                         // Error occurred
)

// String returns the string representation of the update type.
func (t UpdateType) String() string {
	switch t {
	case UpdateIncremental:
		return "Incremental"
	case UpdateComplete:
		return "Complete"
	case UpdateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// AnalysisUpdate represents a streaming analysis update.
type AnalysisUpdate struct {
	Type     UpdateType  // Update type
	Progress float32     // Progress (0.0 to 1.0)
	Data     interface{} // Update data (Entity, Relationship, etc.)
	Error    error       // Error if Type is UpdateError
}

// Analyzer performs semantic analysis on text content.
type Analyzer interface {
	// Analyze performs streaming semantic analysis.
	// Returns a channel of updates that will be closed when complete.
	Analyze(content string) <-chan AnalysisUpdate

	// Stop cancels the current analysis.
	Stop()

	// Results returns the current analysis results.
	// May be partial if analysis is still in progress.
	Results() *Analysis

	// IsRunning returns true if analysis is in progress.
	IsRunning() bool
}

// Token represents a lexical token from the content.
type Token struct {
	Type  TokenType // Token type
	Text  string    // Token text
	Span  Span      // Location in content
	Value string    // Normalized value
}

// TokenType classifies tokens.
type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenWord
	TokenNumber
	TokenPunctuation
	TokenWhitespace
	TokenNewline
	TokenCapitalizedWord // Word starting with capital
	TokenProperNoun      // Proper noun (all caps or Title Case)
	TokenVerb            // Verb (action word)
	TokenTypeHole        // ??Type pattern
	TokenConstraintHole  // !!constraint pattern
)

// String returns the string representation of the token type.
func (t TokenType) String() string {
	switch t {
	case TokenWord:
		return "Word"
	case TokenNumber:
		return "Number"
	case TokenPunctuation:
		return "Punctuation"
	case TokenWhitespace:
		return "Whitespace"
	case TokenNewline:
		return "Newline"
	case TokenCapitalizedWord:
		return "CapitalizedWord"
	case TokenProperNoun:
		return "ProperNoun"
	case TokenVerb:
		return "Verb"
	case TokenTypeHole:
		return "TypeHole"
	case TokenConstraintHole:
		return "ConstraintHole"
	default:
		return "Unknown"
	}
}

// Tokenizer converts text into tokens.
type Tokenizer interface {
	// Tokenize converts content into a stream of tokens.
	Tokenize(content string) []Token
}

// Statistics holds analysis statistics.
type Statistics struct {
	TotalTokens       int
	UniqueEntities    int
	TotalRelationships int
	TotalTypedHoles   int
	TotalDependencies int
	TotalTriples      int
	AnalysisDuration  time.Duration
}

// GetStatistics computes statistics for an analysis.
func (a *Analysis) GetStatistics() Statistics {
	uniqueEntities := make(map[string]bool)
	for _, entity := range a.Entities {
		uniqueEntities[entity.Text] = true
	}

	return Statistics{
		UniqueEntities:    len(uniqueEntities),
		TotalRelationships: len(a.Relationships),
		TotalTypedHoles:   len(a.TypedHoles),
		TotalDependencies: len(a.Dependencies),
		TotalTriples:      len(a.Triples),
		AnalysisDuration:  a.Duration,
	}
}
