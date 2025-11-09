// Package analyze provides relationship pattern mining and analysis.
package analyze

import (
	"fmt"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
	"sort"
	"strings"
)

// RelationshipPattern represents a discovered pattern in the relationship graph.
type RelationshipPattern struct {
	SubjectType   semantic.EntityType // Type of subject entities
	Predicate     string              // Relationship predicate (verb/action)
	ObjectType    semantic.EntityType // Type of object entities
	Occurrences   int                 // Number of times this pattern occurs
	AvgConfidence float64             // Average confidence across all instances
	Examples      []PatternExample    // Example instances of this pattern
	Strength      float64             // Calculated pattern strength (0.0-1.0)
}

// PatternExample represents a concrete instance of a pattern.
type PatternExample struct {
	Subject    string  // Subject entity text
	Predicate  string  // Predicate text
	Object     string  // Object entity text
	Confidence float64 // Confidence for this instance
}

// PatternCluster represents a group of similar relationship patterns.
type PatternCluster struct {
	Predicates   []string             // List of similar predicates
	Patterns     []RelationshipPattern // Patterns in this cluster
	ClusterLabel string               // Representative label for the cluster
	Strength     float64              // Combined cluster strength
}

// MiningOptions configures the pattern mining process.
type MiningOptions struct {
	MinOccurrences   int     // Minimum occurrences to be considered a pattern (default: 2)
	MinConfidence    float64 // Minimum average confidence (default: 0.5)
	MaxExamples      int     // Maximum examples to store per pattern (default: 3)
	ClusterThreshold float64 // Similarity threshold for clustering (default: 0.7)
}

// DefaultMiningOptions returns default mining options.
func DefaultMiningOptions() MiningOptions {
	return MiningOptions{
		MinOccurrences:   2,
		MinConfidence:    0.5,
		MaxExamples:      3,
		ClusterThreshold: 0.7,
	}
}

// MinePatterns discovers relationship patterns from semantic analysis results.
//
// The algorithm:
// 1. Group relationships by (subject_type, predicate, object_type)
// 2. Calculate average confidence for each pattern
// 3. Compute pattern strength using occurrences and confidence
// 4. Filter by minimum occurrences and confidence
// 5. Sort by strength (descending)
func MinePatterns(analysis *semantic.Analysis) []RelationshipPattern {
	return MinePatternsWithOptions(analysis, DefaultMiningOptions())
}

// MinePatternsWithOptions mines patterns with custom options.
func MinePatternsWithOptions(analysis *semantic.Analysis, opts MiningOptions) []RelationshipPattern {
	if analysis == nil || len(analysis.Relationships) == 0 {
		return []RelationshipPattern{}
	}

	// Build entity type lookup from analysis
	entityTypes := buildEntityTypeLookup(analysis)

	// Group relationships by pattern key
	patternMap := make(map[string]*RelationshipPattern)

	for _, rel := range analysis.Relationships {
		// Look up entity types for subject and object
		subjectType := entityTypes[rel.Subject]
		objectType := entityTypes[rel.Object]

		// Create pattern key: subjectType|predicate|objectType
		key := fmt.Sprintf("%d|%s|%d", subjectType, normalizePredicateForKey(rel.Predicate), objectType)

		// Get or create pattern
		pattern, exists := patternMap[key]
		if !exists {
			pattern = &RelationshipPattern{
				SubjectType: subjectType,
				Predicate:   normalizePredicateForKey(rel.Predicate),
				ObjectType:  objectType,
				Examples:    make([]PatternExample, 0, opts.MaxExamples),
			}
			patternMap[key] = pattern
		}

		// Update pattern statistics
		pattern.Occurrences++

		// Add example if we haven't reached max
		if len(pattern.Examples) < opts.MaxExamples {
			pattern.Examples = append(pattern.Examples, PatternExample{
				Subject:    rel.Subject,
				Predicate:  rel.Predicate,
				Object:     rel.Object,
				Confidence: 0.8, // Default confidence (could be enhanced with actual NER confidence)
			})
		}
	}

	// Calculate average confidence and strength for each pattern
	totalRelationships := float64(len(analysis.Relationships))
	for _, pattern := range patternMap {
		// Calculate average confidence from examples
		if len(pattern.Examples) > 0 {
			sum := 0.0
			for _, ex := range pattern.Examples {
				sum += ex.Confidence
			}
			pattern.AvgConfidence = sum / float64(len(pattern.Examples))
		} else {
			pattern.AvgConfidence = 0.8 // Default
		}

		// Calculate strength: (occurrences × avg_confidence) / total_relationships
		// Also factor in diversity (how many unique examples we have)
		diversityFactor := float64(len(pattern.Examples)) / float64(opts.MaxExamples)
		if diversityFactor > 1.0 {
			diversityFactor = 1.0
		}

		pattern.Strength = (float64(pattern.Occurrences) * pattern.AvgConfidence * (1.0 + diversityFactor)) / totalRelationships
	}

	// Convert map to slice and filter
	patterns := make([]RelationshipPattern, 0, len(patternMap))
	for _, pattern := range patternMap {
		// Apply filters
		if pattern.Occurrences < opts.MinOccurrences {
			continue
		}
		if pattern.AvgConfidence < opts.MinConfidence {
			continue
		}

		patterns = append(patterns, *pattern)
	}

	// Sort by strength (descending)
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Strength > patterns[j].Strength
	})

	return patterns
}

// ClusterPatterns groups patterns with similar predicates together.
//
// The algorithm uses simple string similarity (Levenshtein-like):
// - Compare each predicate to find similar ones
// - Group predicates with similarity >= threshold
// - Create cluster labels from most common predicate
func ClusterPatterns(patterns []RelationshipPattern, threshold float64) []PatternCluster {
	if len(patterns) == 0 {
		return []PatternCluster{}
	}

	// Build predicate groups
	predicateGroups := make(map[string][]string)
	processed := make(map[string]bool)

	for _, pattern := range patterns {
		pred := pattern.Predicate

		if processed[pred] {
			continue
		}

		// Find similar predicates
		group := []string{pred}
		processed[pred] = true

		for _, otherPattern := range patterns {
			otherPred := otherPattern.Predicate
			if processed[otherPred] {
				continue
			}

			// Calculate similarity
			similarity := calculatePredicateSimilarity(pred, otherPred)
			if similarity >= threshold {
				group = append(group, otherPred)
				processed[otherPred] = true
			}
		}

		// Use the shortest predicate as cluster label
		label := findShortestString(group)
		predicateGroups[label] = group
	}

	// Create clusters from groups
	clusters := make([]PatternCluster, 0, len(predicateGroups))
	for label, predicates := range predicateGroups {
		cluster := PatternCluster{
			Predicates:   predicates,
			ClusterLabel: label,
			Patterns:     make([]RelationshipPattern, 0),
		}

		// Add patterns that match any predicate in this cluster
		totalStrength := 0.0
		for _, pattern := range patterns {
			for _, pred := range predicates {
				if pattern.Predicate == pred {
					cluster.Patterns = append(cluster.Patterns, pattern)
					totalStrength += pattern.Strength
					break
				}
			}
		}

		cluster.Strength = totalStrength / float64(len(cluster.Patterns))
		clusters = append(clusters, cluster)
	}

	// Sort clusters by strength
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Strength > clusters[j].Strength
	})

	return clusters
}

// PatternStats holds summary statistics about discovered patterns.
type PatternStats struct {
	TotalPatterns    int     // Total number of patterns discovered
	TotalClusters    int     // Total number of clusters
	AvgOccurrences   float64 // Average occurrences per pattern
	AvgConfidence    float64 // Average confidence across all patterns
	TopPredicates    []string // Most common predicates
	UniquePredicates int     // Number of unique predicates
}

// CalculatePatternStats computes summary statistics for patterns.
func CalculatePatternStats(patterns []RelationshipPattern) PatternStats {
	if len(patterns) == 0 {
		return PatternStats{}
	}

	// Calculate averages
	totalOccurrences := 0
	totalConfidence := 0.0
	predicateCounts := make(map[string]int)

	for _, pattern := range patterns {
		totalOccurrences += pattern.Occurrences
		totalConfidence += pattern.AvgConfidence
		predicateCounts[pattern.Predicate]++
	}

	// Find top predicates
	type predicateCount struct {
		predicate string
		count     int
	}
	predicates := make([]predicateCount, 0, len(predicateCounts))
	for pred, count := range predicateCounts {
		predicates = append(predicates, predicateCount{pred, count})
	}
	sort.Slice(predicates, func(i, j int) bool {
		return predicates[i].count > predicates[j].count
	})

	topPredicates := make([]string, 0, min(5, len(predicates)))
	for i := 0; i < min(5, len(predicates)); i++ {
		topPredicates = append(topPredicates, predicates[i].predicate)
	}

	return PatternStats{
		TotalPatterns:    len(patterns),
		AvgOccurrences:   float64(totalOccurrences) / float64(len(patterns)),
		AvgConfidence:    totalConfidence / float64(len(patterns)),
		TopPredicates:    topPredicates,
		UniquePredicates: len(predicateCounts),
	}
}

// Helper functions

// buildEntityTypeLookup creates a map from entity text to entity type.
func buildEntityTypeLookup(analysis *semantic.Analysis) map[string]semantic.EntityType {
	lookup := make(map[string]semantic.EntityType)
	for _, entity := range analysis.Entities {
		lookup[entity.Text] = entity.Type
	}
	return lookup
}

// normalizePredicateForKey normalizes a predicate for use as a map key.
// Converts to lowercase and replaces spaces with underscores.
func normalizePredicateForKey(predicate string) string {
	normalized := strings.ToLower(strings.TrimSpace(predicate))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	return normalized
}

// calculatePredicateSimilarity computes similarity between two predicates.
// Uses a combination of substring matching and edit distance.
func calculatePredicateSimilarity(pred1, pred2 string) float64 {
	if pred1 == pred2 {
		return 1.0
	}

	// Normalize
	p1 := strings.ToLower(pred1)
	p2 := strings.ToLower(pred2)

	// Check for exact match after normalization
	if p1 == p2 {
		return 1.0
	}

	// Check for substring matches
	if strings.Contains(p1, p2) || strings.Contains(p2, p1) {
		return 0.85
	}

	// Calculate Levenshtein distance
	distance := levenshteinDistance(p1, p2)
	maxLen := max(len(p1), len(p2))

	if maxLen == 0 {
		return 1.0
	}

	// Convert distance to similarity (0.0 to 1.0)
	similarity := 1.0 - (float64(distance) / float64(maxLen))
	return similarity
}

// levenshteinDistance calculates the Levenshtein edit distance between two strings.
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create distance matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill in the rest of the matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				min(
					matrix[i][j-1]+1,  // insertion
					matrix[i-1][j-1]+cost, // substitution
				),
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// findShortestString returns the shortest string from a slice.
func findShortestString(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	shortest := strs[0]
	for _, s := range strs[1:] {
		if len(s) < len(shortest) {
			shortest = s
		}
	}
	return shortest
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
