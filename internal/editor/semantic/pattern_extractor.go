package semantic

import (
	"context"
	"strings"
)

// PatternExtractor implements EntityExtractor using pattern matching and keyword dictionaries.
//
// This is the original entity extraction logic from Pedantic Raven, which uses:
// - Keyword dictionaries for person/place/organization/technology
// - Context-aware classification (surrounding words)
// - Capitalization patterns (Title Case, ALL CAPS)
// - Multi-word entity detection
//
// Strengths:
// - Fast (<1ms extraction)
// - No external dependencies
// - Predictable behavior
//
// Limitations:
// - Limited to predefined patterns
// - ~60-70% accuracy
// - Cannot handle ambiguous cases well
type PatternExtractor struct {
	classifier *EntityClassifier
	tokenizer  Tokenizer
}

// NewPatternExtractor creates a new pattern-based entity extractor.
func NewPatternExtractor() *PatternExtractor {
	return &PatternExtractor{
		classifier: NewEntityClassifier(),
		tokenizer:  NewTokenizer(),
	}
}

// ExtractEntities implements EntityExtractor.
func (p *PatternExtractor) ExtractEntities(ctx context.Context, text string, entityTypes []string) ([]Entity, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Tokenize text
	tokens := p.tokenizer.Tokenize(text)

	// Extract multi-word entities first
	multiWordEntities := ExtractMultiWordEntities(tokens, 3)

	// Extract single-word entities
	entityMap := make(map[string]*Entity)

	for _, token := range tokens {
		// Only process word-like tokens
		if !isWordToken(token) {
			continue
		}

		// Build classification context
		ctx := p.buildContext(text, token, tokens)

		// Classify entity
		entityType := p.classifier.ClassifyEntity(token.Text, ctx)

		// Filter by requested entity types if specified
		if len(entityTypes) > 0 && !p.matchesEntityType(entityType, entityTypes) {
			continue
		}

		// Skip unknown entities unless explicitly requested
		if entityType == EntityUnknown && !p.containsType(entityTypes, "unknown") {
			continue
		}

		// Add or update entity
		key := strings.ToLower(token.Text)
		if existing, ok := entityMap[key]; ok {
			existing.Count++
		} else {
			entityMap[key] = &Entity{
				Text:  token.Text,
				Type:  entityType,
				Span:  token.Span,
				Count: 1,
			}
		}
	}

	// Add multi-word entities
	for _, mwe := range multiWordEntities {
		// Filter by requested types
		if len(entityTypes) > 0 && !p.matchesEntityType(mwe.Type, entityTypes) {
			continue
		}

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
	entities := make([]Entity, 0, len(entityMap))
	for _, entity := range entityMap {
		entities = append(entities, *entity)
	}

	return entities, nil
}

// Name implements EntityExtractor.
func (p *PatternExtractor) Name() string {
	return "Pattern"
}

// IsAvailable implements EntityExtractor.
func (p *PatternExtractor) IsAvailable(ctx context.Context) bool {
	// Pattern extractor is always available
	return true
}

// buildContext builds classification context for a token.
func (p *PatternExtractor) buildContext(text string, token Token, allTokens []Token) *ClassificationContext {
	ctx := &ClassificationContext{
		Document: text,
	}

	// Find token position
	tokenPos := -1
	for i, t := range allTokens {
		if t.Span.Start == token.Span.Start && t.Span.End == token.Span.End {
			tokenPos = i
			break
		}
	}

	if tokenPos == -1 {
		return ctx
	}

	// Get preceding words (up to 3)
	for i := tokenPos - 1; i >= 0 && len(ctx.PrecedingWords) < 3; i-- {
		if isWordToken(allTokens[i]) {
			ctx.PrecedingWords = append([]string{allTokens[i].Text}, ctx.PrecedingWords...)
		}
	}

	// Get following words (up to 3)
	for i := tokenPos + 1; i < len(allTokens) && len(ctx.FollowingWords) < 3; i++ {
		if isWordToken(allTokens[i]) {
			ctx.FollowingWords = append(ctx.FollowingWords, allTokens[i].Text)
		}
	}

	return ctx
}

// matchesEntityType checks if an entity type matches any requested types.
func (p *PatternExtractor) matchesEntityType(entityType EntityType, requestedTypes []string) bool {
	entityTypeStr := strings.ToLower(entityType.String())

	for _, reqType := range requestedTypes {
		reqTypeLower := strings.ToLower(reqType)

		// Exact match
		if entityTypeStr == reqTypeLower {
			return true
		}

		// Partial match (e.g., "tech" matches "technology")
		if strings.Contains(entityTypeStr, reqTypeLower) || strings.Contains(reqTypeLower, entityTypeStr) {
			return true
		}
	}

	return false
}

// containsType checks if a type list contains a specific type.
func (p *PatternExtractor) containsType(types []string, target string) bool {
	for _, t := range types {
		if strings.EqualFold(t, target) {
			return true
		}
	}
	return false
}
