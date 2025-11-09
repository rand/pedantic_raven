package semantic

import (
	"context"
	"strings"
	"testing"
	"time"
)

// --- Tokenizer Tests ---

func TestTokenizerBasicWords(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "hello world"

	tokens := tokenizer.Tokenize(content)

	// Filter out whitespace
	words := []Token{}
	for _, tok := range tokens {
		if tok.Type == TokenWord {
			words = append(words, tok)
		}
	}

	if len(words) != 2 {
		t.Fatalf("Expected 2 words, got %d", len(words))
	}

	if words[0].Text != "hello" || words[1].Text != "world" {
		t.Errorf("Expected 'hello' and 'world', got '%s' and '%s'", words[0].Text, words[1].Text)
	}
}

func TestTokenizerCapitalizedWords(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "Hello World"

	tokens := tokenizer.Tokenize(content)

	// Find capitalized words
	capWords := []Token{}
	for _, tok := range tokens {
		if tok.Type == TokenCapitalizedWord || tok.Type == TokenProperNoun {
			capWords = append(capWords, tok)
		}
	}

	if len(capWords) != 2 {
		t.Fatalf("Expected 2 capitalized words, got %d", len(capWords))
	}

	if capWords[0].Text != "Hello" || capWords[1].Text != "World" {
		t.Errorf("Expected 'Hello' and 'World', got '%s' and '%s'", capWords[0].Text, capWords[1].Text)
	}
}

func TestTokenizerProperNouns(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "HTTP API JSON"

	tokens := tokenizer.Tokenize(content)

	// Find proper nouns (all caps)
	properNouns := []Token{}
	for _, tok := range tokens {
		if tok.Type == TokenProperNoun {
			properNouns = append(properNouns, tok)
		}
	}

	if len(properNouns) != 3 {
		t.Fatalf("Expected 3 proper nouns, got %d", len(properNouns))
	}
}

func TestTokenizerVerbs(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "User creates document"

	tokens := tokenizer.Tokenize(content)

	// Find verb
	var verb *Token
	for i := range tokens {
		if tokens[i].Type == TokenVerb {
			verb = &tokens[i]
			break
		}
	}

	if verb == nil {
		t.Fatal("Expected to find a verb token")
	}

	if verb.Text != "creates" {
		t.Errorf("Expected 'creates', got '%s'", verb.Text)
	}
}

func TestTokenizerTypedHoles(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "Implement ??Function here"

	tokens := tokenizer.Tokenize(content)

	// Find typed hole
	var hole *Token
	for i := range tokens {
		if tokens[i].Type == TokenTypeHole {
			hole = &tokens[i]
			break
		}
	}

	if hole == nil {
		t.Fatal("Expected to find a typed hole token")
	}

	if hole.Text != "??Function" {
		t.Errorf("Expected '??Function', got '%s'", hole.Text)
	}

	if hole.Value != "Function" {
		t.Errorf("Expected hole value 'Function', got '%s'", hole.Value)
	}
}

func TestTokenizerConstraintHoles(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "Must satisfy !!performance constraint"

	tokens := tokenizer.Tokenize(content)

	// Find constraint hole
	var hole *Token
	for i := range tokens {
		if tokens[i].Type == TokenConstraintHole {
			hole = &tokens[i]
			break
		}
	}

	if hole == nil {
		t.Fatal("Expected to find a constraint hole token")
	}

	if hole.Text != "!!performance" {
		t.Errorf("Expected '!!performance', got '%s'", hole.Text)
	}

	if hole.Value != "performance" {
		t.Errorf("Expected hole value 'performance', got '%s'", hole.Value)
	}
}

func TestTokenizerNumbers(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "Value is 42 or 3.14"

	tokens := tokenizer.Tokenize(content)

	// Find numbers
	numbers := []Token{}
	for _, tok := range tokens {
		if tok.Type == TokenNumber {
			numbers = append(numbers, tok)
		}
	}

	if len(numbers) != 2 {
		t.Fatalf("Expected 2 numbers, got %d", len(numbers))
	}

	if numbers[0].Text != "42" || numbers[1].Text != "3.14" {
		t.Errorf("Expected '42' and '3.14', got '%s' and '%s'", numbers[0].Text, numbers[1].Text)
	}
}

func TestExtractTypedHoles(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "Need ??Interface and !!constraint here"

	tokens := tokenizer.Tokenize(content)
	holes := ExtractTypedHoles(tokens)

	if len(holes) != 2 {
		t.Fatalf("Expected 2 typed holes, got %d", len(holes))
	}

	// First should be type hole
	if holes[0].Type != "Interface" {
		t.Errorf("Expected type 'Interface', got '%s'", holes[0].Type)
	}

	// Second should be constraint hole
	if holes[1].Constraint != "constraint" {
		t.Errorf("Expected constraint 'constraint', got '%s'", holes[1].Constraint)
	}
}

func TestExtractDependencies(t *testing.T) {
	content := `
		import "github.com/user/package"
		import "fmt"
	`

	deps := ExtractDependencies(content)

	if len(deps) < 2 {
		t.Fatalf("Expected at least 2 dependencies, got %d", len(deps))
	}

	// Check that we found the imports
	found := false
	for _, dep := range deps {
		if dep.Target == "github.com/user/package" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'github.com/user/package' dependency")
	}
}

func TestExtractDependenciesMultiplePatterns(t *testing.T) {
	content := `
		import "golang"
		require('nodejs')
		from 'python'
		#include <cpp>
		use rust;
	`

	deps := ExtractDependencies(content)

	// Should detect multiple dependency patterns
	if len(deps) < 3 {
		t.Fatalf("Expected at least 3 dependencies from different patterns, got %d", len(deps))
	}
}

// --- Entity Extraction Tests ---

func TestExtractEntities(t *testing.T) {
	analyzer := NewAnalyzer().(*StreamingAnalyzer)

	content := "The User creates a Document in the System"

	entities := analyzer.extractEntities(context.Background(), content)

	if len(entities) == 0 {
		t.Fatal("Expected to find entities")
	}

	// Check for specific entities
	entityTexts := make(map[string]bool)
	for _, e := range entities {
		entityTexts[e.Text] = true
	}

	if !entityTexts["User"] {
		t.Error("Expected to find 'User' entity")
	}

	if !entityTexts["Document"] {
		t.Error("Expected to find 'Document' entity")
	}

	if !entityTexts["System"] {
		t.Error("Expected to find 'System' entity")
	}
}

func TestEntityCounting(t *testing.T) {
	analyzer := NewAnalyzer().(*StreamingAnalyzer)

	content := "User creates User and User modifies Document"

	entities := analyzer.extractEntities(context.Background(), content)

	// Find User entity
	var userEntity *Entity
	for i := range entities {
		if entities[i].Text == "User" {
			userEntity = &entities[i]
			break
		}
	}

	if userEntity == nil {
		t.Fatal("Expected to find 'User' entity")
	}

	if userEntity.Count != 3 {
		t.Errorf("Expected User count of 3, got %d", userEntity.Count)
	}
}

// --- Relationship Extraction Tests ---

func TestExtractRelationships(t *testing.T) {
	analyzer := NewAnalyzer().(*StreamingAnalyzer)
	tokenizer := NewTokenizer()

	content := "User creates Document"
	tokens := tokenizer.Tokenize(content)

	relationships := analyzer.extractRelationships(tokens)

	if len(relationships) == 0 {
		t.Fatal("Expected to find at least one relationship")
	}

	rel := relationships[0]

	if rel.Subject != "User" {
		t.Errorf("Expected subject 'User', got '%s'", rel.Subject)
	}

	if rel.Predicate != "creates" {
		t.Errorf("Expected predicate 'creates', got '%s'", rel.Predicate)
	}

	if rel.Object != "Document" {
		t.Errorf("Expected object 'Document', got '%s'", rel.Object)
	}
}

func TestExtractMultipleRelationships(t *testing.T) {
	analyzer := NewAnalyzer().(*StreamingAnalyzer)
	tokenizer := NewTokenizer()

	content := "User creates Document. System validates Request."
	tokens := tokenizer.Tokenize(content)

	relationships := analyzer.extractRelationships(tokens)

	if len(relationships) < 2 {
		t.Fatalf("Expected at least 2 relationships, got %d", len(relationships))
	}
}

// --- Triple Generation Tests ---

func TestGenerateTriples(t *testing.T) {
	analyzer := NewAnalyzer().(*StreamingAnalyzer)

	relationships := []Relationship{
		{
			Subject:   "User",
			Predicate: "creates",
			Object:    "Document",
			Span:      Span{Start: 0, End: 20, Line: 0},
		},
	}

	triples := analyzer.generateTriples(relationships)

	if len(triples) != 1 {
		t.Fatalf("Expected 1 triple, got %d", len(triples))
	}

	triple := triples[0]

	if triple.Subject != "User" {
		t.Errorf("Expected subject 'User', got '%s'", triple.Subject)
	}

	if triple.Predicate != "creates" {
		t.Errorf("Expected predicate 'creates', got '%s'", triple.Predicate)
	}

	if triple.Object != "Document" {
		t.Errorf("Expected object 'Document', got '%s'", triple.Object)
	}
}

// --- Streaming Analyzer Tests ---

func TestAnalyzerAnalyze(t *testing.T) {
	analyzer := NewAnalyzer()

	content := "User creates Document and System validates Request"

	updateChan := analyzer.Analyze(content)

	// Collect all updates
	updates := []AnalysisUpdate{}
	for update := range updateChan {
		updates = append(updates, update)
	}

	if len(updates) == 0 {
		t.Fatal("Expected to receive updates")
	}

	// Check for completion update
	foundComplete := false
	for _, update := range updates {
		if update.Type == UpdateComplete {
			foundComplete = true
			if update.Progress != 1.0 {
				t.Errorf("Expected complete progress of 1.0, got %f", update.Progress)
			}
		}
	}

	if !foundComplete {
		t.Error("Expected to receive UpdateComplete")
	}

	// Verify results
	results := analyzer.Results()

	if results == nil {
		t.Fatal("Expected analysis results")
	}

	if results.Content != content {
		t.Errorf("Expected content to match input")
	}

	if len(results.Entities) == 0 {
		t.Error("Expected to find entities")
	}

	if len(results.Relationships) == 0 {
		t.Error("Expected to find relationships")
	}

	if len(results.Triples) == 0 {
		t.Error("Expected to find triples")
	}
}

func TestAnalyzerCancel(t *testing.T) {
	analyzer := NewAnalyzer()

	// Start analysis with long content
	content := strings.Repeat("User creates Document. ", 1000)
	updateChan := analyzer.Analyze(content)

	// Immediately stop
	time.Sleep(10 * time.Millisecond)
	analyzer.Stop()

	// Drain channel
	for range updateChan {
	}

	// Should not be running
	if analyzer.IsRunning() {
		t.Error("Analyzer should not be running after Stop()")
	}
}

func TestAnalyzerProgressUpdates(t *testing.T) {
	analyzer := NewAnalyzer()

	content := "User creates Document"

	updateChan := analyzer.Analyze(content)

	// Track progress values
	progressValues := []float32{}
	for update := range updateChan {
		if update.Type == UpdateIncremental || update.Type == UpdateComplete {
			progressValues = append(progressValues, update.Progress)
		}
	}

	// Should have multiple progress updates
	if len(progressValues) < 3 {
		t.Errorf("Expected multiple progress updates, got %d", len(progressValues))
	}

	// Progress should increase
	for i := 1; i < len(progressValues); i++ {
		if progressValues[i] < progressValues[i-1] {
			t.Error("Progress should not decrease")
		}
	}

	// Final progress should be 1.0
	if progressValues[len(progressValues)-1] != 1.0 {
		t.Errorf("Expected final progress of 1.0, got %f", progressValues[len(progressValues)-1])
	}
}

func TestAnalyzerConcurrentAnalysis(t *testing.T) {
	analyzer := NewAnalyzer()

	content1 := "User creates Document"
	content2 := "System validates Request"

	// Start first analysis
	chan1 := analyzer.Analyze(content1)

	// Start second analysis (should cancel first)
	time.Sleep(10 * time.Millisecond)
	chan2 := analyzer.Analyze(content2)

	// Drain both channels
	go func() {
		for range chan1 {
		}
	}()

	updates := []AnalysisUpdate{}
	for update := range chan2 {
		updates = append(updates, update)
	}

	// Should have completed second analysis
	results := analyzer.Results()
	if results.Content != content2 {
		t.Errorf("Expected final content to be from second analysis")
	}
}

func TestAnalyzerResults(t *testing.T) {
	analyzer := NewAnalyzer()

	// Before any analysis
	results := analyzer.Results()
	if results == nil {
		t.Error("Results should not be nil even before analysis")
	}

	// After analysis
	content := "User creates Document"
	updateChan := analyzer.Analyze(content)
	for range updateChan {
	}

	results = analyzer.Results()
	if results.Content != content {
		t.Error("Results should contain analyzed content")
	}

	if results.Duration == 0 {
		t.Error("Results should have duration")
	}

	if results.Timestamp.IsZero() {
		t.Error("Results should have timestamp")
	}
}

// --- Statistics Tests ---

func TestAnalysisStatistics(t *testing.T) {
	analyzer := NewAnalyzer()

	content := "User creates Document and User modifies File"

	updateChan := analyzer.Analyze(content)
	for range updateChan {
	}

	results := analyzer.Results()
	stats := results.GetStatistics()

	if stats.UniqueEntities == 0 {
		t.Error("Expected unique entities count")
	}

	if stats.TotalRelationships == 0 {
		t.Error("Expected relationships count")
	}

	if stats.TotalTriples == 0 {
		t.Error("Expected triples count")
	}

	if stats.AnalysisDuration == 0 {
		t.Error("Expected analysis duration")
	}
}

