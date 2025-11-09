package integration

import (
	"fmt"
	"strings"
)

// SampleMarkdownContent provides sample markdown content for testing.
func SampleMarkdownContent() string {
	return `# Project Documentation

## Overview
Alice Chen and Bob Rodriguez are working on a new project called ProjectAlpha.
The project aims to improve system performance by 50% within Q4 2024.

## Architecture

### Components
- API Server: Handles incoming requests
- Database: PostgreSQL for data storage
- Cache: Redis for performance optimization
- Message Queue: RabbitMQ for async processing

### Design Patterns
- ?? Implement Circuit Breaker pattern for resilience
- Microservices architecture
- Event-driven design
- !! Ensure backward compatibility with API v1

## Tasks

### Phase 1: Foundation
- [ ] Design API specifications
- [ ] Setup CI/CD pipeline
- [ ] Create database schema

### Phase 2: Implementation
- [ ] Implement core services
- [ ] Add authentication layer
- [ ] Setup monitoring and logging

### Phase 3: Testing
- [ ] Unit tests (80%+ coverage)
- [ ] Integration tests
- [ ] Load testing

## References
- @alice.chen for API design
- @bob.rodriguez for infrastructure
- #ProjectAlpha in project management tool
- Documentation in /docs/ARCHITECTURE.md
`
}

// SampleCodeContent provides sample code content.
func SampleCodeContent() string {
	return `package main

import (
	"fmt"
	"log"
)

// User represents a user in the system
type User struct {
	ID    int
	Name  string
	Email string
}

// ProcessUser processes a user record
// ?? Need to implement validation logic
func ProcessUser(u User) error {
	if u.Name == "" {
		return fmt.Errorf("user name required")
	}

	// TODO: Add database persistence
	log.Printf("Processing user: %v", u.Name)
	return nil
}

// ?? Implement batch processing for users
func BatchProcessUsers(users []User) error {
	for _, u := range users {
		if err := ProcessUser(u); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	users := []User{
		{1, "Alice", "alice@example.com"},
		{2, "Bob", "bob@example.com"},
		{3, "Charlie", "charlie@example.com"},
	}

	// !! Must handle errors gracefully
	if err := BatchProcessUsers(users); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Processing complete")
}
`
}

// SampleWorkPlan provides a sample work plan.
func SampleWorkPlan() string {
	return `# Q4 2024 Work Plan

## Sprint 1: Foundation (Week 1-2)

### Backend Infrastructure
- [ ] Setup Kubernetes cluster
- [ ] Configure CI/CD pipeline
- [ ] Setup monitoring stack (Prometheus + Grafana)

### Frontend Setup
- [ ] Create React project
- [ ] Setup component library
- [ ] Configure routing

### DevOps
- [ ] Setup staging environment
- [ ] Configure SSL/TLS
- [ ] Setup backup strategy

## Sprint 2: Core Features (Week 3-4)

### Authentication
- [ ] Implement OAuth2 flow
- [ ] Create login page
- [ ] Setup user management

### API Development
- [ ] Design RESTful endpoints
- [ ] Implement CRUD operations
- [ ] Add input validation

### Database
- [ ] Create schema
- [ ] Setup replication
- [ ] Configure backups

## Sprint 3: Testing & QA (Week 5-6)

### Testing
- [ ] Write unit tests (80%+ coverage)
- [ ] Write integration tests
- [ ] Setup end-to-end tests

### Performance
- [ ] Load testing
- [ ] Optimization
- [ ] Benchmarking

### Security
- [ ] Security audit
- [ ] Penetration testing
- [ ] Fix vulnerabilities

## Sprint 4: Deployment (Week 7-8)

### Pre-Production
- [ ] Staging deployment
- [ ] User acceptance testing
- [ ] Performance validation

### Production
- [ ] Production deployment
- [ ] Monitoring setup
- [ ] Incident response plan

### Maintenance
- [ ] Documentation
- [ ] Training
- [ ] Support setup

## Key Personnel
- Alice Chen: Product Manager
- Bob Rodriguez: Tech Lead
- Charlie Zhang: Senior Developer
- Diana Park: QA Lead
- Eve Johnson: DevOps Engineer

## Risks & Mitigations
- ?? Complexity: May require additional resources
- !! Timeline: Must account for review cycles
- Integration: Coordinate with third-party services
`
}

// SampleEntityContent provides content rich in entities for analysis.
func SampleEntityContent() string {
	return `# Business Meeting Notes

## Attendees
- John Smith (CEO) from Microsoft
- Sarah Johnson (CTO) from Google
- Michael Chen (VP Engineering) from Apple
- Lisa Anderson (Product Director) from Amazon
- David Park (Director) from Meta

## Companies Represented
- Microsoft
- Google
- Apple
- Amazon
- Meta
- IBM

## Meeting Topics

### Product Strategy
The team discussed strategies for improving cloud services.
John Smith emphasized the importance of security.
Sarah Johnson presented Google's approach to scalability.

### Market Opportunities
Michael Chen identified opportunities in enterprise market.
Lisa Anderson discussed consumer segment potential.
David Park presented Meta's vision for metaverse integration.

### Technology Stack
- Primary: Python and Go
- Frontend: React and Vue.js
- Infrastructure: Kubernetes and Docker
- Databases: PostgreSQL and MongoDB

### Locations Mentioned
- San Francisco (HQ discussions)
- Seattle (engineering team)
- New York (marketing team)
- London (European operations)
- Tokyo (Asia-Pacific hub)
- Sydney (APAC operations)

### Timeline
- Q1 2024: Market research
- Q2 2024: Product development
- Q3 2024: Beta testing
- Q4 2024: General availability

## Action Items
- John Smith: Prepare board presentation
- Sarah Johnson: Setup technical working group
- Michael Chen: Design scalability architecture
- Lisa Anderson: Create marketing plan
- David Park: Define integration requirements
`
}

// GenerateEntityRichContent creates content with specified number of entities.
func GenerateEntityRichContent(entityCount int) string {
	entities := []string{
		"Alice", "Bob", "Charlie", "Diana", "Eve",
		"Frank", "Grace", "Henry", "Iris", "Jack",
	}

	companies := []string{
		"Acme Corp", "TechFlow Inc", "DataSys Ltd",
		"CloudBridge Co", "InnovateTech", "FutureWare",
		"QuantumSoft", "VelocityAI", "SynergyTech", "NexusCore",
	}

	locations := []string{
		"San Francisco", "New York", "London",
		"Tokyo", "Berlin", "Singapore", "Toronto", "Sydney",
	}

	var sb strings.Builder
	sb.WriteString("# Entity Rich Content\n\n")

	for i := 0; i < entityCount; i++ {
		entity := entities[i%len(entities)]
		company := companies[i%len(companies)]
		location := locations[i%len(locations)]

		sb.WriteString(fmt.Sprintf(
			"%s works at %s in %s. ",
			entity, company, location,
		))

		if i%5 == 0 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// GenerateLargeContent creates large content for performance testing.
func GenerateLargeContent(lines int) string {
	var sb strings.Builder

	for i := 0; i < lines; i++ {
		sb.WriteString(fmt.Sprintf(
			"Line %d: Lorem ipsum dolor sit amet, consectetur adipiscing elit. "+
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\n",
			i+1,
		))
	}

	return sb.String()
}

// GenerateNestedContent creates nested structure for hierarchy testing.
func GenerateNestedContent(depth int, itemsPerLevel int) string {
	var sb strings.Builder
	sb.WriteString("# Nested Structure\n\n")

	var buildNested func(int, int, *strings.Builder)
	buildNested = func(currentDepth, maxDepth int, b *strings.Builder) {
		if currentDepth >= maxDepth {
			return
		}

		for i := 0; i < itemsPerLevel; i++ {
			indent := strings.Repeat("  ", currentDepth)
			b.WriteString(fmt.Sprintf("%s- Item %d.%d\n", indent, currentDepth, i))
			buildNested(currentDepth+1, maxDepth, b)
		}
	}

	buildNested(0, depth, &sb)
	return sb.String()
}

// GenerateCodeWithTypedHoles creates code content with typed holes.
func GenerateCodeWithTypedHoles() string {
	return `package main

import (
	"fmt"
)

// DataProcessor handles data processing
type DataProcessor struct {
	// ?? Implementation details
}

// ?? Create interface for data transformers
func (dp *DataProcessor) Transform(data string) string {
	// ?? Implement transformation logic
	return data
}

// ?? Add error handling wrapper
func SafeTransform(data string) (string, error) {
	// !! Must handle nil input
	processor := &DataProcessor{}
	return processor.Transform(data), nil
}

func main() {
	input := "test data"
	// ?? Add validation before processing
	result, _ := SafeTransform(input)
	fmt.Println(result)
}
`
}

// SampleSmallContent provides minimal test content.
func SampleSmallContent() string {
	return "Hello World"
}

// SampleEmptyContent provides empty test content.
func SampleEmptyContent() string {
	return ""
}

// SampleSpecialCharacters provides content with special characters.
func SampleSpecialCharacters() string {
	return `# Special Characters Test

## Symbols
@ # $ % ^ & * ( ) _ + = - [ ] { } | ; : ' " < > , . ? /

## Unicode
- Greek: Α Β Γ Δ Ε Ζ Η Θ
- Cyrillic: А Б В Г Д Е Ж З
- Arabic: ا ب ج د ه و ز ح
- Chinese: 中文 日本語 한국어

## Math Symbols
∑ ∫ √ ∛ ∞ ± × ÷ ≤ ≥ ≠ ≈

## Emojis
😀 🎉 🚀 ⚡ 💡 ❌ ✓ ⚠️

## Code Snippets
` + "```python" + `
def hello():
    print("Hello, World!")
` + "```" + `

## Tables
| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| A        | B        | C        |
| 1        | 2        | 3        |
`
}
