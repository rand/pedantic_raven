package semantic

import (
	"context"
	"sync"
	"time"
)

// StreamingAnalyzer performs streaming semantic analysis.
type StreamingAnalyzer struct {
	mu            sync.RWMutex
	tokenizer     Tokenizer
	extractor     EntityExtractor
	analysis      *Analysis
	running       bool
	cancel        context.CancelFunc
	updateChan    chan AnalysisUpdate
	channelClosed *sync.Once
}

// NewAnalyzer creates a new streaming analyzer with the default pattern-based extractor.
func NewAnalyzer() Analyzer {
	return &StreamingAnalyzer{
		tokenizer: NewTokenizer(),
		extractor: NewPatternExtractor(),
	}
}

// NewAnalyzerWithExtractor creates a new streaming analyzer with a custom extractor.
func NewAnalyzerWithExtractor(extractor EntityExtractor) Analyzer {
	return &StreamingAnalyzer{
		tokenizer: NewTokenizer(),
		extractor: extractor,
	}
}

// Analyze implements Analyzer.
func (a *StreamingAnalyzer) Analyze(content string) <-chan AnalysisUpdate {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Cancel any existing analysis
	if a.cancel != nil {
		a.cancel()
	}

	// Create new context and channel
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.channelClosed = &sync.Once{} // Reset for each new analysis
	a.updateChan = make(chan AnalysisUpdate, 10)
	a.running = true

	// Initialize analysis
	a.analysis = &Analysis{
		Content:   content,
		Timestamp: time.Now(),
	}

	// Start analysis in background
	go a.performAnalysis(ctx, content)

	return a.updateChan
}

// Stop implements Analyzer.
func (a *StreamingAnalyzer) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}

	a.running = false
}

// Results implements Analyzer.
func (a *StreamingAnalyzer) Results() *Analysis {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.analysis == nil {
		return &Analysis{}
	}

	// Return a copy
	return a.analysis
}

// IsRunning implements Analyzer.
func (a *StreamingAnalyzer) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.running
}

// performAnalysis performs the actual analysis.
func (a *StreamingAnalyzer) performAnalysis(ctx context.Context, content string) {
	startTime := time.Now()

	// Capture the update channel for this analysis run
	a.mu.RLock()
	updateChan := a.updateChan
	a.mu.RUnlock()

	defer func() {
		a.mu.Lock()
		a.running = false
		a.analysis.Duration = time.Since(startTime)
		closedOnce := a.channelClosed
		a.mu.Unlock()

		// Ensure channel is closed only once, even if called concurrently
		closedOnce.Do(func() {
			close(updateChan)
		})
	}()

	// Step 1: Tokenization (10% of progress)
	a.sendUpdate(updateChan, AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.1,
	})

	tokens := a.tokenizer.Tokenize(content)

	if ctx.Err() != nil {
		return
	}

	// Step 2: Extract entities (30%)
	a.sendUpdate(updateChan, AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.3,
	})

	entities := a.extractEntities(ctx, content)
	a.mu.Lock()
	a.analysis.Entities = entities
	a.mu.Unlock()

	for _, entity := range entities {
		a.sendUpdate(updateChan, AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.3,
			Data:     entity,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 3: Extract relationships (50%)
	a.sendUpdate(updateChan, AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.5,
	})

	relationships := a.extractRelationships(tokens)
	a.mu.Lock()
	a.analysis.Relationships = relationships
	a.mu.Unlock()

	for _, rel := range relationships {
		a.sendUpdate(updateChan, AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.5,
			Data:     rel,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 4: Extract typed holes (70%)
	a.sendUpdate(updateChan, AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.7,
	})

	typedHoles := ExtractTypedHoles(tokens)
	a.mu.Lock()
	a.analysis.TypedHoles = typedHoles
	a.mu.Unlock()

	for _, hole := range typedHoles {
		a.sendUpdate(updateChan, AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.7,
			Data:     hole,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 5: Extract dependencies (85%)
	a.sendUpdate(updateChan, AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.85,
	})

	dependencies := ExtractDependencies(content)
	a.mu.Lock()
	a.analysis.Dependencies = dependencies
	a.mu.Unlock()

	for _, dep := range dependencies {
		a.sendUpdate(updateChan, AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.85,
			Data:     dep,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 6: Generate triples (100%)
	a.sendUpdate(updateChan, AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.95,
	})

	triples := a.generateTriples(relationships)
	a.mu.Lock()
	a.analysis.Triples = triples
	a.mu.Unlock()

	// Complete - use blocking send for final update
	updateChan <- AnalysisUpdate{
		Type:     UpdateComplete,
		Progress: 1.0,
	}
}

// extractEntities extracts entities using the configured extractor.
func (a *StreamingAnalyzer) extractEntities(ctx context.Context, content string) []Entity {
	// Use extractor to extract all entity types
	entities, err := a.extractor.ExtractEntities(ctx, content, []string{})
	if err != nil {
		// If extractor fails, return empty list
		// (errors are logged but don't fail the entire analysis)
		return []Entity{}
	}

	return entities
}

// extractRelationships extracts relationships from tokens.
func (a *StreamingAnalyzer) extractRelationships(tokens []Token) []Relationship {
	var relationships []Relationship

	// Filter out non-word tokens for pattern matching
	var words []Token
	for _, tok := range tokens {
		if tok.Type != TokenWhitespace && tok.Type != TokenNewline && tok.Type != TokenPunctuation {
			words = append(words, tok)
		}
	}

	// Simple pattern: Entity Verb Entity
	for i := 0; i < len(words)-2; i++ {
		if isEntityToken(words[i]) && words[i+1].Type == TokenVerb && isEntityToken(words[i+2]) {
			relationships = append(relationships, Relationship{
				Subject:   words[i].Text,
				Predicate: words[i+1].Text,
				Object:    words[i+2].Text,
				Span:      words[i].Span,
			})
		}
	}

	return relationships
}

// generateTriples generates triples from relationships.
func (a *StreamingAnalyzer) generateTriples(relationships []Relationship) []Triple {
	var triples []Triple

	for _, rel := range relationships {
		triples = append(triples, Triple{
			Subject:   rel.Subject,
			Predicate: rel.Predicate,
			Object:    rel.Object,
			Source:    rel.Span,
		})
	}

	return triples
}

// sendUpdate sends an update to the channel (non-blocking).
func (a *StreamingAnalyzer) sendUpdate(updateChan chan<- AnalysisUpdate, update AnalysisUpdate) {
	select {
	case updateChan <- update:
	default:
		// Channel full, drop update
	}
}

// isEntityToken checks if a token represents an entity.
func isEntityToken(token Token) bool {
	return token.Type == TokenCapitalizedWord || token.Type == TokenProperNoun || token.Type == TokenWord
}
