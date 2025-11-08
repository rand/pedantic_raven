package semantic

import (
	"strings"
	"testing"
)

// --- Entity Classifier Tests ---

func TestEntityClassifierTechnologyTerms(t *testing.T) {
	classifier := NewEntityClassifier()

	tests := []struct {
		text     string
		expected EntityType
	}{
		{"HTTP", EntityTechnology},
		{"JSON", EntityTechnology},
		{"GraphQL", EntityTechnology},
		{"PostgreSQL", EntityTechnology},
		{"Docker", EntityTechnology},
		{"React", EntityTechnology},
	}

	for _, tt := range tests {
		result := classifier.ClassifyEntity(tt.text, nil)
		if result != tt.expected {
			t.Errorf("Expected %s to be classified as %v, got %v",
				tt.text, tt.expected, result)
		}
	}
}

func TestEntityClassifierPersonIndicators(t *testing.T) {
	classifier := NewEntityClassifier()

	tests := []struct {
		text     string
		expected EntityType
	}{
		{"Dr", EntityPerson},
		{"Developer", EntityPerson},
		{"Manager", EntityPerson},
		{"User", EntityPerson},
		{"CEO", EntityPerson},
	}

	for _, tt := range tests {
		result := classifier.ClassifyEntity(tt.text, nil)
		if result != tt.expected {
			t.Errorf("Expected %s to be classified as %v, got %v",
				tt.text, tt.expected, result)
		}
	}
}

func TestEntityClassifierPlaceIndicators(t *testing.T) {
	classifier := NewEntityClassifier()

	tests := []struct {
		text     string
		expected EntityType
	}{
		{"Server", EntityPlace},
		{"Database", EntityPlace},
		{"City", EntityPlace},
		{"Building", EntityPlace},
	}

	for _, tt := range tests {
		result := classifier.ClassifyEntity(tt.text, nil)
		if result != tt.expected {
			t.Errorf("Expected %s to be classified as %v, got %v",
				tt.text, tt.expected, result)
		}
	}
}

func TestEntityClassifierOrganizationSuffixes(t *testing.T) {
	classifier := NewEntityClassifier()

	tests := []struct {
		text     string
		expected EntityType
	}{
		{"CompanyInc", EntityOrganization},
		{"CorporationLLC", EntityOrganization},
		{"Department", EntityOrganization},
		{"Organization", EntityOrganization},
	}

	for _, tt := range tests {
		result := classifier.ClassifyEntity(tt.text, nil)
		if result != tt.expected {
			t.Errorf("Expected %s to be classified as %v, got %v",
				tt.text, tt.expected, result)
		}
	}
}

func TestEntityClassifierWithContext(t *testing.T) {
	classifier := NewEntityClassifier()

	// Context: "User Alice creates..."
	context := &ClassificationContext{
		PrecedingWords: []string{"User"},
		FollowingWords: []string{"creates"},
	}

	result := classifier.ClassifyEntity("Alice", context)
	if result != EntityPerson {
		t.Errorf("Expected 'Alice' with 'User' context to be Person, got %v", result)
	}
}

func TestEntityClassifierContextualPlace(t *testing.T) {
	classifier := NewEntityClassifier()

	// Context: "Database Production..."
	context := &ClassificationContext{
		PrecedingWords: []string{},
		FollowingWords: []string{"Database"},
	}

	result := classifier.ClassifyEntity("Production", context)
	if result != EntityPlace {
		t.Errorf("Expected 'Production' with 'Database' following to be Place, got %v", result)
	}
}

func TestEntityClassifierAllCaps(t *testing.T) {
	classifier := NewEntityClassifier()

	// Unknown all-caps words should default to Technology
	result := classifier.ClassifyEntity("XYZ", nil)
	if result != EntityTechnology {
		t.Errorf("Expected all-caps unknown 'XYZ' to be Technology, got %v", result)
	}
}

func TestEntityClassifierTitleCase(t *testing.T) {
	classifier := NewEntityClassifier()

	// Unknown Title Case words should default to Concept
	result := classifier.ClassifyEntity("Something", nil)
	if result != EntityConcept {
		t.Errorf("Expected Title Case unknown 'Something' to be Concept, got %v", result)
	}
}

// --- Multi-Word Entity Tests ---

func TestExtractMultiWordEntitiesBasic(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "User Account System"

	tokens := tokenizer.Tokenize(content)
	entities := ExtractMultiWordEntities(tokens, 3)

	if len(entities) == 0 {
		t.Fatal("Expected to find multi-word entities")
	}

	// Should find "User Account" and possibly "Account System"
	found := false
	for _, entity := range entities {
		if entity.Text == "User Account" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find 'User Account' multi-word entity")
	}
}

func TestExtractMultiWordEntitiesWithCount(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "User Account and User Account"

	tokens := tokenizer.Tokenize(content)
	entities := ExtractMultiWordEntities(tokens, 3)

	// Find "User Account" entity
	var userAccount *MultiWordEntity
	for i := range entities {
		if entities[i].Text == "User Account" {
			userAccount = &entities[i]
			break
		}
	}

	if userAccount == nil {
		t.Fatal("Expected to find 'User Account' entity")
	}

	if userAccount.Count != 2 {
		t.Errorf("Expected 'User Account' count of 2, got %d", userAccount.Count)
	}
}

func TestExtractMultiWordEntitiesMaxWords(t *testing.T) {
	tokenizer := NewTokenizer()
	content := "Database Management System Server"

	tokens := tokenizer.Tokenize(content)

	// Test with maxWords = 2
	entities := ExtractMultiWordEntities(tokens, 2)

	for _, entity := range entities {
		if len(entity.Words) > 2 {
			t.Errorf("Expected max 2 words, got %d in '%s'",
				len(entity.Words), entity.Text)
		}
	}
}

func TestMultiWordEntityClassification(t *testing.T) {
	tests := []struct {
		words    []string
		expected EntityType
	}{
		{[]string{"HTTP", "Server"}, EntityTechnology},
		{[]string{"User", "Account"}, EntityPerson},
		{[]string{"Database", "Server"}, EntityPlace},
		{[]string{"Company", "Inc"}, EntityOrganization},
		{[]string{"Primary", "Key"}, EntityConcept},
	}

	for _, tt := range tests {
		result := classifyMultiWordEntity(tt.words)
		if result != tt.expected {
			t.Errorf("Expected %v to be classified as %v, got %v",
				tt.words, tt.expected, result)
		}
	}
}

func TestIsValidMultiWordEntity(t *testing.T) {
	tests := []struct {
		words    []string
		expected bool
	}{
		{[]string{"User", "Account"}, true},     // Capitalized
		{[]string{"user", "account"}, false},    // Not capitalized
		{[]string{"User"}, false},               // Single word
		{[]string{"HTTP", "Server"}, true},      // All caps + Title
		{[]string{"the", "system"}, false},      // No capitalization
	}

	for _, tt := range tests {
		result := isValidMultiWordEntity(tt.words)
		if result != tt.expected {
			t.Errorf("isValidMultiWordEntity(%v) = %v, expected %v",
				tt.words, result, tt.expected)
		}
	}
}

// --- Integration Tests with Analyzer ---

func TestAnalyzerEnhancedEntityExtraction(t *testing.T) {
	analyzer := NewAnalyzer()

	content := "The Developer manages the HTTP Server and the Database System"

	updateChan := analyzer.Analyze(content)

	// Drain updates
	for range updateChan {
	}

	results := analyzer.Results()

	if len(results.Entities) == 0 {
		t.Fatal("Expected to find entities")
	}

	// Check for specific entity types
	foundDeveloper := false
	foundHTTP := false
	foundServer := false

	for _, entity := range results.Entities {
		if entity.Text == "Developer" && entity.Type == EntityPerson {
			foundDeveloper = true
		}
		if entity.Text == "HTTP" && entity.Type == EntityTechnology {
			foundHTTP = true
		}
		if entity.Text == "Server" && entity.Type == EntityPlace {
			foundServer = true
		}
	}

	if !foundDeveloper {
		t.Error("Expected to find Developer (Person)")
	}

	if !foundHTTP {
		t.Error("Expected to find HTTP (Technology)")
	}

	if !foundServer {
		t.Error("Expected to find Server (Place)")
	}
}

func TestAnalyzerMultiWordEntities(t *testing.T) {
	analyzer := NewAnalyzer()

	content := "User Account connects to Database Server"

	updateChan := analyzer.Analyze(content)

	// Drain updates
	for range updateChan {
	}

	results := analyzer.Results()

	// Should find multi-word entities
	foundMultiWord := false
	for _, entity := range results.Entities {
		if len(entity.Text) > 1 && entity.Text != strings.ToLower(entity.Text) {
			// Has spaces or is multi-word
			if entity.Text == "User Account" || entity.Text == "Database Server" {
				foundMultiWord = true
				break
			}
		}
	}

	if !foundMultiWord {
		t.Error("Expected to find multi-word entities")
	}
}

func TestAnalyzerContextAwareClassification(t *testing.T) {
	analyzer := NewAnalyzer()

	// "Alice" after "User" should be classified as Person
	content := "User Alice creates Document"

	updateChan := analyzer.Analyze(content)

	// Drain updates
	for range updateChan {
	}

	results := analyzer.Results()

	// Check Alice classification
	for _, entity := range results.Entities {
		if entity.Text == "Alice" {
			if entity.Type != EntityPerson && entity.Type != EntityConcept {
				t.Errorf("Expected Alice to be Person or Concept, got %v", entity.Type)
			}
			return
		}
	}

	t.Error("Expected to find Alice entity")
}

// --- Relationship Pattern Tests ---

func TestDefaultRelationshipPatterns(t *testing.T) {
	patterns := DefaultRelationshipPatterns()

	if len(patterns) == 0 {
		t.Fatal("Expected default relationship patterns")
	}

	// Check for basic patterns
	foundEntityVerbEntity := false
	for _, pattern := range patterns {
		if pattern.Name == "entity-verb-entity" {
			foundEntityVerbEntity = true
			if pattern.Confidence == 0 {
				t.Error("Expected pattern to have confidence score")
			}
		}
	}

	if !foundEntityVerbEntity {
		t.Error("Expected to find entity-verb-entity pattern")
	}
}

func TestRelationshipPatternConfidence(t *testing.T) {
	patterns := DefaultRelationshipPatterns()

	for _, pattern := range patterns {
		if pattern.Confidence < 0 || pattern.Confidence > 1 {
			t.Errorf("Pattern %s has invalid confidence: %f",
				pattern.Name, pattern.Confidence)
		}
	}
}
