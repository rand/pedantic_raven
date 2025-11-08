package semantic

import (
	"context"
	"strings"
	"sync"
	"time"
)

// StreamingAnalyzer performs streaming semantic analysis.
type StreamingAnalyzer struct {
	mu         sync.RWMutex
	tokenizer  Tokenizer
	classifier *EntityClassifier
	analysis   *Analysis
	running    bool
	cancel     context.CancelFunc
	updateChan chan AnalysisUpdate
}

// NewAnalyzer creates a new streaming analyzer.
func NewAnalyzer() Analyzer {
	return &StreamingAnalyzer{
		tokenizer:  NewTokenizer(),
		classifier: NewEntityClassifier(),
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
		a.mu.Unlock()
		close(updateChan)
	}()

	// Step 1: Tokenization (10% of progress)
	a.sendUpdate(AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.1,
	})

	tokens := a.tokenizer.Tokenize(content)

	if ctx.Err() != nil {
		return
	}

	// Step 2: Extract entities (30%)
	a.sendUpdate(AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.3,
	})

	entities := a.extractEntities(tokens)
	a.mu.Lock()
	a.analysis.Entities = entities
	a.mu.Unlock()

	for _, entity := range entities {
		a.sendUpdate(AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.3,
			Data:     entity,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 3: Extract relationships (50%)
	a.sendUpdate(AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.5,
	})

	relationships := a.extractRelationships(tokens)
	a.mu.Lock()
	a.analysis.Relationships = relationships
	a.mu.Unlock()

	for _, rel := range relationships {
		a.sendUpdate(AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.5,
			Data:     rel,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 4: Extract typed holes (70%)
	a.sendUpdate(AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.7,
	})

	typedHoles := ExtractTypedHoles(tokens)
	a.mu.Lock()
	a.analysis.TypedHoles = typedHoles
	a.mu.Unlock()

	for _, hole := range typedHoles {
		a.sendUpdate(AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.7,
			Data:     hole,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 5: Extract dependencies (85%)
	a.sendUpdate(AnalysisUpdate{
		Type:     UpdateIncremental,
		Progress: 0.85,
	})

	dependencies := ExtractDependencies(content)
	a.mu.Lock()
	a.analysis.Dependencies = dependencies
	a.mu.Unlock()

	for _, dep := range dependencies {
		a.sendUpdate(AnalysisUpdate{
			Type:     UpdateIncremental,
			Progress: 0.85,
			Data:     dep,
		})
	}

	if ctx.Err() != nil {
		return
	}

	// Step 6: Generate triples (100%)
	a.sendUpdate(AnalysisUpdate{
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

// extractEntities extracts entities from tokens with enhanced classification.
func (a *StreamingAnalyzer) extractEntities(tokens []Token) []Entity {
	entityMap := make(map[string]*Entity)

	// Build context for each potential entity
	for i, token := range tokens {
		if token.Type == TokenCapitalizedWord || token.Type == TokenProperNoun {
			key := token.Value

			if existing, ok := entityMap[key]; ok {
				existing.Count++
			} else {
				// Build classification context
				context := a.buildContext(tokens, i)

				// Use enhanced classifier
				entityType := a.classifier.ClassifyEntity(token.Text, context)

				entityMap[key] = &Entity{
					Text:  token.Text,
					Type:  entityType,
					Span:  token.Span,
					Count: 1,
				}
			}
		}
	}

	// Also extract multi-word entities
	multiWordEntities := ExtractMultiWordEntities(tokens, 3)
	for _, mwe := range multiWordEntities {
		key := strings.ToLower(mwe.Text)
		if existing, ok := entityMap[key]; ok {
			existing.Count += mwe.Count
		} else {
			entityMap[key] = &Entity{
				Text:  mwe.Text,
				Type:  mwe.Type,
				Span:  mwe.Span,
				Count: mwe.Count,
			}
		}
	}

	// Convert map to slice
	var entities []Entity
	for _, entity := range entityMap {
		entities = append(entities, *entity)
	}

	return entities
}

// buildContext creates a classification context for a token.
func (a *StreamingAnalyzer) buildContext(tokens []Token, index int) *ClassificationContext {
	context := &ClassificationContext{
		PrecedingWords: []string{},
		FollowingWords: []string{},
	}

	// Collect preceding words (up to 3)
	for i := index - 1; i >= 0 && len(context.PrecedingWords) < 3; i-- {
		if isWordToken(tokens[i]) {
			context.PrecedingWords = append([]string{tokens[i].Text}, context.PrecedingWords...)
		}
	}

	// Collect following words (up to 3)
	for i := index + 1; i < len(tokens) && len(context.FollowingWords) < 3; i++ {
		if isWordToken(tokens[i]) {
			context.FollowingWords = append(context.FollowingWords, tokens[i].Text)
		}
	}

	return context
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
func (a *StreamingAnalyzer) sendUpdate(update AnalysisUpdate) {
	select {
	case a.updateChan <- update:
	default:
		// Channel full, drop update
	}
}

// isEntityToken checks if a token represents an entity.
func isEntityToken(token Token) bool {
	return token.Type == TokenCapitalizedWord || token.Type == TokenProperNoun || token.Type == TokenWord
}
