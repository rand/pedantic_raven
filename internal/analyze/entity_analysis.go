package analyze

import (
	"math"
	"sort"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// EntityFrequency represents frequency data for an entity.
type EntityFrequency struct {
	Text       string                // Entity text
	Type       semantic.EntityType   // Entity type
	Count      int                   // Total occurrences
	Importance int                   // Importance score (0-10)
}

// CalculateEntityFrequency computes entity frequency from semantic analysis.
// Returns a slice of EntityFrequency sorted by count (descending).
func CalculateEntityFrequency(analysis *semantic.Analysis) []EntityFrequency {
	if analysis == nil || len(analysis.Entities) == 0 {
		return []EntityFrequency{}
	}

	// Aggregate entities by text (case-insensitive)
	entityMap := make(map[string]*EntityFrequency)

	for _, entity := range analysis.Entities {
		key := entity.Text
		if existing, ok := entityMap[key]; ok {
			// Update count
			existing.Count += entity.Count
		} else {
			// Create new entry
			entityMap[key] = &EntityFrequency{
				Text:       entity.Text,
				Type:       entity.Type,
				Count:      entity.Count,
				Importance: calculateImportance(entity.Count, entity.Type),
			}
		}
	}

	// Convert map to slice
	result := make([]EntityFrequency, 0, len(entityMap))
	for _, ef := range entityMap {
		result = append(result, *ef)
	}

	// Sort by count descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}

// calculateImportance computes importance score based on frequency and type.
// Returns a score from 0-10.
func calculateImportance(count int, entityType semantic.EntityType) int {
	// Base importance from frequency (logarithmic scale)
	// 1 occurrence = 1, 10 occurrences = 2, 100 occurrences = 3, etc.
	freqScore := 0
	if count > 0 {
		freqScore = int(math.Log10(float64(count)) * 3)
	}

	// Type-based bonus
	typeBonus := 0
	switch entityType {
	case semantic.EntityPerson, semantic.EntityOrganization:
		typeBonus = 2 // Higher importance for people and orgs
	case semantic.EntityTechnology, semantic.EntityConcept:
		typeBonus = 1 // Medium importance for tech and concepts
	case semantic.EntityPlace, semantic.EntityThing:
		typeBonus = 0 // Lower importance for places and things
	}

	importance := freqScore + typeBonus
	if importance > 10 {
		importance = 10
	}
	if importance < 0 {
		importance = 0
	}

	return importance
}

// FrequencyList is a sortable slice of EntityFrequency.
type FrequencyList []EntityFrequency

// SortByFrequency sorts entities by count (descending).
func (fl FrequencyList) SortByFrequency() {
	sort.Slice(fl, func(i, j int) bool {
		if fl[i].Count != fl[j].Count {
			return fl[i].Count > fl[j].Count
		}
		// Secondary sort by text (alphabetical)
		return fl[i].Text < fl[j].Text
	})
}

// SortByType sorts entities by type, then by count within each type.
func (fl FrequencyList) SortByType() {
	sort.Slice(fl, func(i, j int) bool {
		if fl[i].Type != fl[j].Type {
			return fl[i].Type < fl[j].Type
		}
		// Within same type, sort by count
		if fl[i].Count != fl[j].Count {
			return fl[i].Count > fl[j].Count
		}
		// Tertiary sort by text
		return fl[i].Text < fl[j].Text
	})
}

// SortAlphabetical sorts entities alphabetically by text.
func (fl FrequencyList) SortAlphabetical() {
	sort.Slice(fl, func(i, j int) bool {
		return fl[i].Text < fl[j].Text
	})
}

// FilterByType filters entities by type.
func (fl FrequencyList) FilterByType(entityType semantic.EntityType) FrequencyList {
	result := make(FrequencyList, 0)
	for _, ef := range fl {
		if ef.Type == entityType {
			result = append(result, ef)
		}
	}
	return result
}

// FilterByMinCount filters entities with count >= minCount.
func (fl FrequencyList) FilterByMinCount(minCount int) FrequencyList {
	result := make(FrequencyList, 0)
	for _, ef := range fl {
		if ef.Count >= minCount {
			result = append(result, ef)
		}
	}
	return result
}

// TopN returns the top N entities by count.
func (fl FrequencyList) TopN(n int) FrequencyList {
	if n >= len(fl) {
		return fl
	}
	if n < 0 {
		n = 0
	}
	return fl[:n]
}

// BarChartData represents data for bar chart visualization.
type BarChartData struct {
	TypeCounts map[semantic.EntityType]int // Count per entity type
	MaxCount   int                         // Maximum count (for scaling)
	TotalCount int                         // Total entity count
}

// CalculateBarChartData prepares data for bar chart rendering.
func CalculateBarChartData(frequencies []EntityFrequency) BarChartData {
	data := BarChartData{
		TypeCounts: make(map[semantic.EntityType]int),
	}

	for _, ef := range frequencies {
		data.TypeCounts[ef.Type] += ef.Count
		data.TotalCount += ef.Count

		if ef.Count > data.MaxCount {
			data.MaxCount = ef.Count
		}
	}

	return data
}

// GetTypeCounts returns counts aggregated by entity type.
func GetTypeCounts(frequencies []EntityFrequency) map[semantic.EntityType]int {
	counts := make(map[semantic.EntityType]int)
	for _, ef := range frequencies {
		counts[ef.Type]++
	}
	return counts
}

// GetTypeFrequencies returns the sum of frequencies by entity type.
func GetTypeFrequencies(frequencies []EntityFrequency) map[semantic.EntityType]int {
	sums := make(map[semantic.EntityType]int)
	for _, ef := range frequencies {
		sums[ef.Type] += ef.Count
	}
	return sums
}

// SortedTypes returns entity types sorted by total frequency (descending).
func SortedTypes(frequencies []EntityFrequency) []semantic.EntityType {
	typeFreqs := GetTypeFrequencies(frequencies)

	// Convert to slice for sorting
	type typeFreqPair struct {
		entityType semantic.EntityType
		freq       int
	}

	pairs := make([]typeFreqPair, 0, len(typeFreqs))
	for t, f := range typeFreqs {
		pairs = append(pairs, typeFreqPair{t, f})
	}

	// Sort by frequency descending
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].freq > pairs[j].freq
	})

	// Extract sorted types
	result := make([]semantic.EntityType, len(pairs))
	for i, p := range pairs {
		result[i] = p.entityType
	}

	return result
}
