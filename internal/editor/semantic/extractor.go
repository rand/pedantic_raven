package semantic

import (
	"context"
)

// EntityExtractor provides an interface for entity extraction strategies.
//
// This interface allows different entity extraction implementations:
// - PatternExtractor: Uses pattern matching and keyword dictionaries
// - GLiNERExtractor: Uses GLiNER ML model for zero-shot NER
//
// The extractor is pluggable, allowing graceful degradation when ML services
// are unavailable while maintaining a consistent API.
type EntityExtractor interface {
	// ExtractEntities extracts entities from text.
	//
	// Parameters:
	//   - ctx: Context for cancellation
	//   - text: Text to analyze
	//   - entityTypes: Desired entity types (e.g., ["person", "organization"])
	//
	// Returns:
	//   - List of extracted entities
	//   - Error if extraction fails
	ExtractEntities(ctx context.Context, text string, entityTypes []string) ([]Entity, error)

	// Name returns the extractor name (e.g., "Pattern", "GLiNER")
	Name() string

	// IsAvailable returns true if the extractor is ready to use.
	// For PatternExtractor, always returns true.
	// For GLiNERExtractor, checks if service is reachable.
	IsAvailable(ctx context.Context) bool
}

// HybridExtractor combines multiple extractors with fallback logic.
//
// It attempts extraction with the primary extractor first, falling back
// to the secondary extractor if the primary fails or is unavailable.
type HybridExtractor struct {
	primary   EntityExtractor // Preferred extractor (e.g., GLiNER)
	fallback  EntityExtractor // Fallback extractor (e.g., Pattern)
	useFallback bool          // Whether fallback is enabled
}

// NewHybridExtractor creates a new hybrid extractor.
//
// Parameters:
//   - primary: Primary extractor to try first
//   - fallback: Fallback extractor if primary fails
//   - useFallback: Whether to use fallback (if false, only uses primary)
func NewHybridExtractor(primary, fallback EntityExtractor, useFallback bool) *HybridExtractor {
	return &HybridExtractor{
		primary:     primary,
		fallback:    fallback,
		useFallback: useFallback,
	}
}

// ExtractEntities implements EntityExtractor.
func (h *HybridExtractor) ExtractEntities(ctx context.Context, text string, entityTypes []string) ([]Entity, error) {
	// Try primary extractor if available
	if h.primary != nil && h.primary.IsAvailable(ctx) {
		entities, err := h.primary.ExtractEntities(ctx, text, entityTypes)
		if err == nil {
			return entities, nil
		}
		// Primary failed, try fallback if enabled
	}

	// Use fallback if enabled and available
	if h.useFallback && h.fallback != nil && h.fallback.IsAvailable(ctx) {
		return h.fallback.ExtractEntities(ctx, text, entityTypes)
	}

	// No extractors available
	return nil, ErrNoExtractorAvailable
}

// Name implements EntityExtractor.
func (h *HybridExtractor) Name() string {
	if h.primary != nil && h.primary.IsAvailable(context.Background()) {
		return h.primary.Name()
	}
	if h.useFallback && h.fallback != nil {
		return h.fallback.Name() + " (Fallback)"
	}
	return "Hybrid"
}

// IsAvailable implements EntityExtractor.
func (h *HybridExtractor) IsAvailable(ctx context.Context) bool {
	// Available if either extractor is available
	if h.primary != nil && h.primary.IsAvailable(ctx) {
		return true
	}
	if h.useFallback && h.fallback != nil && h.fallback.IsAvailable(ctx) {
		return true
	}
	return false
}

// SetPrimary updates the primary extractor.
func (h *HybridExtractor) SetPrimary(extractor EntityExtractor) {
	h.primary = extractor
}

// SetFallback updates the fallback extractor.
func (h *HybridExtractor) SetFallback(extractor EntityExtractor) {
	h.fallback = extractor
}

// EnableFallback enables or disables fallback extraction.
func (h *HybridExtractor) EnableFallback(enabled bool) {
	h.useFallback = enabled
}
