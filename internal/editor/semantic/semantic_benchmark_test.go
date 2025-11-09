package semantic

import (
	"context"
	"strings"
	"testing"
)

// BenchmarkTokenize benchmarks text tokenization.
func BenchmarkTokenize(b *testing.B) {
	tokenizer := NewTokenizer()
	texts := []string{
		"Simple text.",
		"The quick brown fox jumps over the lazy dog. This is a benchmark test.",
		strings.Repeat("This is a longer text with many words. ", 10),
		strings.Repeat("Even longer text for stress testing tokenization. ", 50),
	}

	for i, text := range texts {
		b.Run(string(rune(i))+"_"+string(rune(len(text)))+"chars", func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				_ = tokenizer.Tokenize(text)
			}
		})
	}
}

// BenchmarkPatternExtractorSimple benchmarks pattern-based entity extraction.
func BenchmarkPatternExtractorSimple(b *testing.B) {
	extractor := NewPatternExtractor()
	ctx := context.Background()
	text := "John Smith works at Microsoft in Seattle. He develops software using Go and Rust."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = extractor.ExtractEntities(ctx, text, nil)
	}
}

// BenchmarkPatternExtractorComplexText benchmarks extraction on complex text.
func BenchmarkPatternExtractorComplexText(b *testing.B) {
	extractor := NewPatternExtractor()
	ctx := context.Background()

	textSizes := []struct {
		name string
		text string
	}{
		{"100words", generateTestText(100)},
		{"500words", generateTestText(500)},
		{"1000words", generateTestText(1000)},
		{"2000words", generateTestText(2000)},
	}

	for _, size := range textSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = extractor.ExtractEntities(ctx, size.text, nil)
			}
		})
	}
}

// BenchmarkPatternExtractorWithTypeFilter benchmarks filtered extraction.
func BenchmarkPatternExtractorWithTypeFilter(b *testing.B) {
	extractor := NewPatternExtractor()
	ctx := context.Background()
	text := generateTestText(500)

	filters := []struct {
		name  string
		types []string
	}{
		{"AllTypes", nil},
		{"PersonOnly", []string{"person"}},
		{"TechOnly", []string{"technology"}},
		{"PersonAndOrg", []string{"person", "organization"}},
	}

	for _, filter := range filters {
		b.Run(filter.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = extractor.ExtractEntities(ctx, text, filter.types)
			}
		})
	}
}

// BenchmarkEntityClassification benchmarks entity type classification.
func BenchmarkEntityClassification(b *testing.B) {
	classifier := NewEntityClassifier()
	testWords := []string{
		"John", "Microsoft", "Seattle", "Python", "algorithm",
		"database", "AWS", "Tokyo", "software", "blockchain",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := testWords[i%len(testWords)]
		ctx := &ClassificationContext{}
		_ = classifier.ClassifyEntity(word, ctx)
	}
}

// BenchmarkMultiWordEntityExtraction benchmarks multi-word entity detection.
func BenchmarkMultiWordEntityExtraction(b *testing.B) {
	tokenizer := NewTokenizer()
	text := "John Smith from New York City works at Red Hat in San Francisco. " +
		"He studied Computer Science at Massachusetts Institute of Technology."
	tokens := tokenizer.Tokenize(text)

	maxLengths := []int{2, 3, 4}
	for _, maxLen := range maxLengths {
		b.Run(string(rune(maxLen))+"words", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = ExtractMultiWordEntities(tokens, maxLen)
			}
		})
	}
}

// BenchmarkHybridExtractor benchmarks the hybrid extractor fallback logic.
func BenchmarkHybridExtractor(b *testing.B) {
	primary := NewPatternExtractor()
	fallback := NewPatternExtractor()
	hybrid := NewHybridExtractor(primary, fallback, true)

	ctx := context.Background()
	text := generateTestText(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = hybrid.ExtractEntities(ctx, text, nil)
	}
}

// BenchmarkAnalyzerStreamingAnalysis benchmarks full streaming analysis.
func BenchmarkAnalyzerStreamingAnalysis(b *testing.B) {
	textSizes := []int{100, 500, 1000}

	for _, size := range textSizes {
		b.Run(string(rune(size))+"words", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				analyzer := NewAnalyzer()
				text := generateTestText(size)
				b.StartTimer()

				updateChan := analyzer.Analyze(text)
				// Drain the channel
				for range updateChan {
				}

				b.StopTimer()
				analyzer.Stop()
				b.StartTimer()
			}
		})
	}
}

// BenchmarkRelationshipExtraction benchmarks relationship detection.
func BenchmarkRelationshipExtraction(b *testing.B) {
	// This would benchmark relationship extraction if we have access to the function
	// For now, we'll benchmark entity pairs that could form relationships
	extractor := NewPatternExtractor()
	ctx := context.Background()
	text := "Alice works with Bob. Charlie manages Dave. Eve leads the team."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entities, _ := extractor.ExtractEntities(ctx, text, nil)
		// Count potential relationships (n*(n-1)/2 pairs)
		_ = len(entities) * (len(entities) - 1) / 2
	}
}

// BenchmarkTypedHoleDetection benchmarks typed hole pattern detection.
func BenchmarkTypedHoleDetection(b *testing.B) {
	texts := []string{
		"Need to implement ??Function here",
		"TODO: Add ??Interface for handling requests",
		strings.Repeat("Some code here. ??Type placeholder. More code. ", 10),
	}

	for i, text := range texts {
		b.Run(string(rune(i)), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				// Simple pattern matching for typed holes
				_ = strings.Count(text, "??")
			}
		})
	}
}

// BenchmarkEntityDeduplication benchmarks entity deduplication logic.
func BenchmarkEntityDeduplication(b *testing.B) {
	// Create entities with duplicates
	entities := make([]Entity, 1000)
	for i := 0; i < 1000; i++ {
		entities[i] = Entity{
			Text:  string(rune('A' + (i % 26))),
			Type:  EntityTypeConcept,
			Count: 1,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityMap := make(map[string]*Entity)
		for i := range entities {
			e := &entities[i]
			key := strings.ToLower(e.Text)
			if existing, ok := entityMap[key]; ok {
				existing.Count++
			} else {
				entityMap[key] = e
			}
		}
	}
}

// BenchmarkContextBuilding benchmarks classification context creation.
func BenchmarkContextBuilding(b *testing.B) {
	tokenizer := NewTokenizer()
	extractor := NewPatternExtractor()
	text := generateTestText(500)
	tokens := tokenizer.Tokenize(text)

	if len(tokens) == 0 {
		b.Skip("No tokens generated")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token := tokens[i%len(tokens)]
		_ = extractor.buildContext(text, token, tokens)
	}
}

// BenchmarkMemoryAllocation benchmarks memory usage during extraction.
func BenchmarkMemoryAllocation(b *testing.B) {
	extractor := NewPatternExtractor()
	ctx := context.Background()
	text := generateTestText(500)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = extractor.ExtractEntities(ctx, text, nil)
	}
}

// Helper function to generate test text with entities
func generateTestText(wordCount int) string {
	entities := []string{
		"John Smith", "Microsoft", "Seattle", "Python", "Alice Johnson",
		"Google", "New York", "JavaScript", "Bob Williams", "Amazon",
		"San Francisco", "Go", "Charlie Brown", "Apple", "Boston",
		"Rust", "Diana Prince", "Meta", "Chicago", "TypeScript",
	}

	words := []string{
		"works", "at", "in", "uses", "develops", "with", "from",
		"creates", "manages", "leads", "designs", "implements",
		"the", "a", "and", "or", "but", "for", "to", "of",
	}

	var result strings.Builder
	for i := 0; i < wordCount; i++ {
		if i%5 == 0 {
			result.WriteString(entities[i%len(entities)])
		} else {
			result.WriteString(words[i%len(words)])
		}
		result.WriteString(" ")
	}

	return result.String()
}
