package semantic

import (
	"strings"
)

// EntityClassifier provides enhanced entity classification.
type EntityClassifier struct {
	// Knowledge bases for classification
	personIndicators     map[string]bool
	placeIndicators      map[string]bool
	organizationSuffixes map[string]bool
	technologyTerms      map[string]bool
}

// NewEntityClassifier creates a new entity classifier.
func NewEntityClassifier() *EntityClassifier {
	return &EntityClassifier{
		personIndicators: map[string]bool{
			// Titles
			"Dr": true, "Mr": true, "Mrs": true, "Ms": true, "Miss": true,
			"Prof": true, "Professor": true, "Sir": true, "Dame": true,
			"Lord": true, "Lady": true, "Captain": true, "Major": true,
			// Roles
			"Developer": true, "Engineer": true, "Manager": true,
			"Director": true, "President": true, "CEO": true,
			"Designer": true, "Analyst": true, "Administrator": true,
			"User": true, "Client": true, "Customer": true,
		},
		placeIndicators: map[string]bool{
			// Types
			"City": true, "Country": true, "State": true, "Region": true,
			"District": true, "Province": true, "Territory": true,
			// Locations
			"Street": true, "Avenue": true, "Road": true, "Lane": true,
			"Building": true, "Office": true, "Campus": true,
			"Server": true, "Database": true, "Repository": true,
		},
		organizationSuffixes: map[string]bool{
			"Inc": true, "Corp": true, "LLC": true, "Ltd": true,
			"Company": true, "Corporation": true, "Organization": true,
			"Foundation": true, "Institute": true, "Agency": true,
			"Department": true, "Division": true, "Team": true,
		},
		technologyTerms: map[string]bool{
			// Protocols
			"HTTP": true, "HTTPS": true, "FTP": true, "SSH": true,
			"TCP": true, "UDP": true, "IP": true, "DNS": true,
			// Formats
			"JSON": true, "XML": true, "YAML": true, "CSV": true,
			"HTML": true, "CSS": true, "Markdown": true,
			// Databases
			"SQL": true, "NoSQL": true, "PostgreSQL": true, "MySQL": true,
			"MongoDB": true, "Redis": true, "Elasticsearch": true,
			// APIs
			"REST": true, "GraphQL": true, "gRPC": true, "SOAP": true,
			// Languages
			"Go": true, "Python": true, "JavaScript": true, "TypeScript": true,
			"Rust": true, "Java": true, "C": true, "Ruby": true,
			// Frameworks
			"React": true, "Angular": true, "Vue": true, "Django": true,
			"Flask": true, "Express": true, "Rails": true,
			// Infrastructure
			"Docker": true, "Kubernetes": true, "AWS": true, "Azure": true,
			"GCP": true, "CI": true, "CD": true,
			// Concepts
			"API": true, "SDK": true, "CLI": true, "GUI": true,
			"IDE": true, "URL": true, "URI": true, "UUID": true,
		},
	}
}

// ClassifyEntity determines the entity type with context awareness.
func (c *EntityClassifier) ClassifyEntity(text string, context *ClassificationContext) EntityType {
	// Check technology terms first (most specific)
	if c.technologyTerms[text] {
		return EntityTechnology
	}

	// Check for organization suffixes
	for suffix := range c.organizationSuffixes {
		if strings.HasSuffix(text, suffix) || strings.Contains(text, suffix) {
			return EntityOrganization
		}
	}

	// Check person indicators
	if c.personIndicators[text] {
		return EntityPerson
	}

	// Check place indicators
	if c.placeIndicators[text] {
		return EntityPlace
	}

	// Context-based classification
	if context != nil {
		// Check preceding words
		if len(context.PrecedingWords) > 0 {
			prev := context.PrecedingWords[len(context.PrecedingWords)-1]

			// "User X", "Developer Y" -> Person
			if c.personIndicators[prev] {
				return EntityPerson
			}

			// "City X", "Server Y" -> Place
			if c.placeIndicators[prev] {
				return EntityPlace
			}
		}

		// Check following words
		if len(context.FollowingWords) > 0 {
			next := context.FollowingWords[0]

			// "X Server", "Y Database" -> Place/Technology
			if c.placeIndicators[next] || c.technologyTerms[next] {
				return EntityPlace
			}
		}
	}

	// Check for capitalization patterns
	if isAllUpper(text) {
		// All caps usually indicates acronym or technology
		return EntityTechnology
	}

	if isTitleCase(text) {
		// Title case often indicates proper noun
		// Default to concept unless we have more context
		return EntityConcept
	}

	// Default classification based on word characteristics
	if len(text) > 0 && isUpperCase(text[0]) {
		return EntityConcept
	}

	return EntityUnknown
}

// isUpperCase checks if a byte represents an uppercase letter.
func isUpperCase(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

// ClassificationContext provides surrounding context for entity classification.
type ClassificationContext struct {
	PrecedingWords []string // Words before the entity
	FollowingWords []string // Words after the entity
	Sentence       string   // Full sentence containing the entity
	Document       string   // Full document (optional)
}


// MultiWordEntity represents an entity spanning multiple words.
type MultiWordEntity struct {
	Words []string   // Component words
	Text  string     // Full text
	Type  EntityType // Entity type
	Span  Span       // Location in content
	Count int        // Occurrences
}

// ExtractMultiWordEntities extracts entities that span multiple words.
func ExtractMultiWordEntities(tokens []Token, maxWords int) []MultiWordEntity {
	var entities []MultiWordEntity
	entityMap := make(map[string]*MultiWordEntity)

	if maxWords < 2 {
		maxWords = 2
	}
	if maxWords > 5 {
		maxWords = 5 // Reasonable upper limit
	}

	// Filter to get only word tokens with their original positions
	type wordWithPos struct {
		token Token
		pos   int
	}
	var wordTokens []wordWithPos
	for i, tok := range tokens {
		if isWordToken(tok) {
			wordTokens = append(wordTokens, wordWithPos{token: tok, pos: i})
		}
	}

	// Extract multi-word sequences from filtered words
	for i := 0; i < len(wordTokens); i++ {
		// Try different window sizes
		for windowSize := 2; windowSize <= maxWords && i+windowSize-1 < len(wordTokens); windowSize++ {
			words := []string{}
			for j := 0; j < windowSize; j++ {
				words = append(words, wordTokens[i+j].token.Text)
			}

			// Check if this forms a valid entity pattern
			if isValidMultiWordEntity(words) {
				text := strings.Join(words, " ")
				key := strings.ToLower(text)

				if existing, ok := entityMap[key]; ok {
					existing.Count++
				} else {
					entityMap[key] = &MultiWordEntity{
						Words: words,
						Text:  text,
						Type:  classifyMultiWordEntity(words),
						Span:  wordTokens[i].token.Span,
						Count: 1,
					}
				}
			}
		}
	}

	// Convert map to slice
	for _, entity := range entityMap {
		entities = append(entities, *entity)
	}

	return entities
}

// isWordToken checks if a token is a word-like token.
func isWordToken(token Token) bool {
	return token.Type == TokenWord ||
		token.Type == TokenCapitalizedWord ||
		token.Type == TokenProperNoun
}

// isValidMultiWordEntity checks if a word sequence forms a valid multi-word entity.
func isValidMultiWordEntity(words []string) bool {
	if len(words) < 2 {
		return false
	}

	// At least one word should be capitalized
	hasCapital := false
	for _, word := range words {
		if len(word) > 0 && isUpperCase(word[0]) {
			hasCapital = true
			break
		}
	}

	if !hasCapital {
		return false
	}

	// Common patterns
	// Pattern: Title + Noun (e.g., "User Account", "Document Service")
	// Pattern: Adjective + Noun (e.g., "Primary Key", "Foreign Key")
	// Pattern: Noun + Noun (e.g., "Database Server", "File System")

	return true
}

// classifyMultiWordEntity classifies a multi-word entity.
func classifyMultiWordEntity(words []string) EntityType {
	classifier := NewEntityClassifier()

	// Check all words for technology terms (highest priority)
	for _, word := range words {
		if classifier.technologyTerms[word] {
			return EntityTechnology
		}
	}

	// Check last word (often carries the type)
	if len(words) > 0 {
		lastWord := words[len(words)-1]

		// Person indicators
		if classifier.personIndicators[lastWord] {
			return EntityPerson
		}

		// Place indicators
		if classifier.placeIndicators[lastWord] {
			return EntityPlace
		}

		// Organization suffixes
		for suffix := range classifier.organizationSuffixes {
			if strings.Contains(lastWord, suffix) {
				return EntityOrganization
			}
		}
	}

	// Check first word for person indicators
	if len(words) > 0 {
		firstWord := words[0]

		if classifier.personIndicators[firstWord] {
			return EntityPerson
		}
	}

	// Default to concept for multi-word entities
	return EntityConcept
}

// RelationshipPattern represents a pattern for relationship extraction.
type RelationshipPattern struct {
	Name        string   // Pattern name
	Pattern     []string // Token type pattern (e.g., ["entity", "verb", "entity"])
	Confidence  float64  // Pattern confidence (0.0 to 1.0)
	Description string   // Pattern description
}

// DefaultRelationshipPatterns returns common relationship patterns.
func DefaultRelationshipPatterns() []RelationshipPattern {
	return []RelationshipPattern{
		{
			Name:        "entity-verb-entity",
			Pattern:     []string{"entity", "verb", "entity"},
			Confidence:  0.9,
			Description: "Simple subject-predicate-object pattern",
		},
		{
			Name:        "entity-verb-prep-entity",
			Pattern:     []string{"entity", "verb", "prep", "entity"},
			Confidence:  0.85,
			Description: "Subject-verb-preposition-object (e.g., 'User belongs to Group')",
		},
		{
			Name:        "entity-has-entity",
			Pattern:     []string{"entity", "has", "entity"},
			Confidence:  0.95,
			Description: "Possession or composition relationship",
		},
		{
			Name:        "entity-is-entity",
			Pattern:     []string{"entity", "is", "entity"},
			Confidence:  0.9,
			Description: "Identity or type relationship",
		},
	}
}

// RelationshipWithConfidence represents a relationship with confidence score.
type RelationshipWithConfidence struct {
	Relationship Relationship
	Confidence   float64
	Pattern      string
}
