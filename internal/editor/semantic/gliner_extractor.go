package semantic

import (
	"context"
	"strings"

	"github.com/rand/pedantic-raven/internal/gliner"
)

// GLiNERExtractor implements EntityExtractor using the GLiNER ML model.
//
// GLiNER (Generalist Model for Named Entity Recognition) provides:
// - Zero-shot NER: Extract ANY entity types defined at runtime
// - High accuracy: ~85-95% (outperforms ChatGPT on NER benchmarks)
// - Context-aware: Understands surrounding text for disambiguation
// - Custom types: Support domain-specific entity types
//
// Requires:
// - GLiNER service running (Python FastAPI service)
// - Network connection to service URL
// - ~1GB RAM for service (340M parameter model)
//
// Performance:
// - Extraction: 100-300ms per request (vs <1ms for pattern)
// - Accuracy: Much higher than pattern matching
// - Resource usage: Runs in separate Python process
type GLiNERExtractor struct {
	client         *gliner.Client
	defaultTypes   []string // Default entity types if none specified
	scoreThreshold float64  // Minimum confidence score (0.0-1.0)
}

// NewGLiNERExtractor creates a new GLiNER-based entity extractor.
//
// Parameters:
//   - client: GLiNER client (configured with service URL, timeouts, etc.)
//   - defaultTypes: Default entity types to extract if none specified
//   - scoreThreshold: Minimum confidence score (0.0-1.0), typical: 0.3
func NewGLiNERExtractor(client *gliner.Client, defaultTypes []string, scoreThreshold float64) *GLiNERExtractor {
	if scoreThreshold == 0 {
		scoreThreshold = 0.3 // Default threshold
	}

	if defaultTypes == nil || len(defaultTypes) == 0 {
		// Sensible defaults for context engineering
		defaultTypes = []string{
			"person",
			"organization",
			"location",
			"technology",
			"concept",
			"product",
		}
	}

	return &GLiNERExtractor{
		client:         client,
		defaultTypes:   defaultTypes,
		scoreThreshold: scoreThreshold,
	}
}

// ExtractEntities implements EntityExtractor.
func (g *GLiNERExtractor) ExtractEntities(ctx context.Context, text string, entityTypes []string) ([]Entity, error) {
	// Use default types if none specified
	if len(entityTypes) == 0 {
		entityTypes = g.defaultTypes
	}

	// Call GLiNER service
	glinerEntities, err := g.client.ExtractEntities(ctx, text, entityTypes, g.scoreThreshold)
	if err != nil {
		// Check if service is disabled
		if gliner.IsDisabled(err) {
			return nil, ErrExtractorDisabled
		}
		// Check if service is unavailable
		if gliner.IsUnavailable(err) {
			return nil, ErrExtractorUnavailable
		}
		return nil, err
	}

	// Convert GLiNER entities to semantic entities
	entities := make([]Entity, 0, len(glinerEntities))
	entityCounts := make(map[string]int) // Track entity occurrences

	for _, ge := range glinerEntities {
		// Map GLiNER label to EntityType
		entityType := g.mapLabelToType(ge.Label)

		// Create entity
		entity := Entity{
			Text: ge.Text,
			Type: entityType,
			Span: Span{
				Start: ge.Start,
				End:   ge.End,
			},
			Count: 1, // Will be aggregated below
		}

		// Track occurrences (normalize by lowercase)
		key := strings.ToLower(ge.Text)
		entityCounts[key]++

		entities = append(entities, entity)
	}

	// Update counts for duplicate entities
	for i := range entities {
		key := strings.ToLower(entities[i].Text)
		entities[i].Count = entityCounts[key]
	}

	// Deduplicate entities (keep first occurrence of each unique text)
	seen := make(map[string]bool)
	uniqueEntities := make([]Entity, 0, len(entities))

	for _, entity := range entities {
		key := strings.ToLower(entity.Text)
		if !seen[key] {
			uniqueEntities = append(uniqueEntities, entity)
			seen[key] = true
		}
	}

	return uniqueEntities, nil
}

// Name implements EntityExtractor.
func (g *GLiNERExtractor) Name() string {
	return "GLiNER"
}

// IsAvailable implements EntityExtractor.
func (g *GLiNERExtractor) IsAvailable(ctx context.Context) bool {
	if !g.client.IsEnabled() {
		return false
	}

	// Check service availability
	err := g.client.CheckAvailability(ctx)
	return err == nil
}

// SetScoreThreshold updates the confidence threshold.
func (g *GLiNERExtractor) SetScoreThreshold(threshold float64) {
	if threshold >= 0.0 && threshold <= 1.0 {
		g.scoreThreshold = threshold
	}
}

// SetDefaultTypes updates the default entity types.
func (g *GLiNERExtractor) SetDefaultTypes(types []string) {
	g.defaultTypes = types
}

// mapLabelToType maps GLiNER label strings to EntityType enum.
//
// GLiNER uses lowercase labels (e.g., "person", "organization"),
// while EntityType uses PascalCase enums.
func (g *GLiNERExtractor) mapLabelToType(label string) EntityType {
	// Normalize label
	labelLower := strings.ToLower(strings.TrimSpace(label))

	// Direct mappings
	switch labelLower {
	case "person", "people", "individual":
		return EntityPerson
	case "organization", "org", "company", "corporation":
		return EntityOrganization
	case "location", "place", "city", "country", "region":
		return EntityPlace
	case "technology", "tech", "software", "hardware", "tool":
		return EntityTechnology
	case "concept", "idea", "notion", "principle":
		return EntityConcept
	case "thing", "object", "item", "product":
		return EntityThing
	default:
		// Custom types or unknown → map to Concept
		// (Custom types like "api_endpoint" → Concept)
		return EntityConcept
	}
}
