// Package main demonstrates the usage of Analyze Mode in Pedantic Raven.
// This is a standalone example showing how to:
// - Create and configure an analyze mode instance
// - Load sample semantic data
// - Switch between different views
// - Apply filters and perform analysis
// - Export reports in multiple formats
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rand/pedantic-raven/internal/analyze"
	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

func main() {
	fmt.Println("Pedantic Raven - Analyze Mode Demo")
	fmt.Println("===================================")

	// Step 1: Create sample semantic analysis data
	// In a real application, this would come from the semantic analyzer
	analysis := createSampleAnalysis()
	fmt.Printf("Created sample analysis with %d entities and %d relationships\n\n",
		len(analysis.Entities), len(analysis.Relationships))

	// Step 2: Create and configure analyze mode
	mode := analyze.NewAnalyzeMode()
	mode.SetSize(120, 40) // Terminal size: 120 columns × 40 rows
	fmt.Println("Initialized Analyze Mode with 120×40 viewport")

	// Step 3: Set the analysis data
	mode.SetAnalysis(analysis)
	fmt.Println("Loaded analysis data into Analyze Mode")

	// Step 4: Demonstrate entity frequency analysis
	demonstrateEntityAnalysis(analysis)

	// Step 5: Demonstrate relationship pattern mining
	demonstratePatternMining(analysis)

	// Step 6: Demonstrate typed hole prioritization
	demonstrateHoleAnalysis(analysis)

	// Step 7: Demonstrate filtering
	demonstrateFiltering(analysis)

	// Step 8: Demonstrate export functionality
	demonstrateExport(mode, analysis)

	fmt.Println("\nDemo completed successfully!")
}

// createSampleAnalysis creates sample semantic analysis data for demonstration.
func createSampleAnalysis() *semantic.Analysis {
	return &semantic.Analysis{
		// Sample entities representing a software project
		Entities: []semantic.Entity{
			// People
			{Text: "John", Type: semantic.EntityPerson, Count: 12},
			{Text: "Alice", Type: semantic.EntityPerson, Count: 8},
			{Text: "Bob", Type: semantic.EntityPerson, Count: 6},

			// Organizations
			{Text: "Acme Corp", Type: semantic.EntityOrganization, Count: 10},
			{Text: "Tech Inc", Type: semantic.EntityOrganization, Count: 7},
			{Text: "StartupXYZ", Type: semantic.EntityOrganization, Count: 4},

			// Technologies
			{Text: "API", Type: semantic.EntityTechnology, Count: 15},
			{Text: "Database", Type: semantic.EntityTechnology, Count: 11},
			{Text: "React", Type: semantic.EntityTechnology, Count: 9},
			{Text: "Python", Type: semantic.EntityTechnology, Count: 8},
			{Text: "GraphQL", Type: semantic.EntityTechnology, Count: 7},
			{Text: "PostgreSQL", Type: semantic.EntityTechnology, Count: 6},

			// Places
			{Text: "New York", Type: semantic.EntityPlace, Count: 5},
			{Text: "San Francisco", Type: semantic.EntityPlace, Count: 4},

			// Concepts
			{Text: "Authentication", Type: semantic.EntityConcept, Count: 9},
			{Text: "Microservices", Type: semantic.EntityConcept, Count: 6},
			{Text: "API Design", Type: semantic.EntityConcept, Count: 5},

			// Things
			{Text: "Server", Type: semantic.EntityThing, Count: 8},
			{Text: "Container", Type: semantic.EntityThing, Count: 5},
			{Text: "Cache", Type: semantic.EntityThing, Count: 4},
		},

		// Sample relationships showing common patterns
		Relationships: []semantic.Relationship{
			// Person → works_at → Organization pattern
			{Subject: "John", Predicate: "works_at", Object: "Acme Corp"},
			{Subject: "Alice", Predicate: "works_at", Object: "Tech Inc"},
			{Subject: "Bob", Predicate: "works_at", Object: "StartupXYZ"},

			// Organization → creates → Technology pattern
			{Subject: "Acme Corp", Predicate: "creates", Object: "API"},
			{Subject: "Tech Inc", Predicate: "creates", Object: "Database"},
			{Subject: "StartupXYZ", Predicate: "develops", Object: "React"},

			// Technology → uses → Technology pattern
			{Subject: "API", Predicate: "uses", Object: "Database"},
			{Subject: "React", Predicate: "uses", Object: "GraphQL"},
			{Subject: "API", Predicate: "uses", Object: "Authentication"},
			{Subject: "Database", Predicate: "stores", Object: "Cache"},

			// Person → develops → Technology pattern
			{Subject: "John", Predicate: "develops", Object: "API"},
			{Subject: "Alice", Predicate: "develops", Object: "Database"},
			{Subject: "Bob", Predicate: "develops", Object: "React"},

			// Organization → located_in → Place pattern
			{Subject: "Acme Corp", Predicate: "located_in", Object: "New York"},
			{Subject: "Tech Inc", Predicate: "located_in", Object: "San Francisco"},
			{Subject: "StartupXYZ", Predicate: "located_in", Object: "San Francisco"},

			// Technology → implements → Concept pattern
			{Subject: "API", Predicate: "implements", Object: "Authentication"},
			{Subject: "Database", Predicate: "implements", Object: "Microservices"},
			{Subject: "GraphQL", Predicate: "follows", Object: "API Design"},
		},

		// Sample typed holes for implementation planning
		TypedHoles: []semantic.TypedHole{
			{
				Type:       "AuthService",
				Constraint: "thread-safe, async",
			},
			{
				Type:       "DatabaseLayer",
				Constraint: "atomic, concurrent",
			},
			{
				Type:       "ConfigLoader",
				Constraint: "immutable",
			},
			{
				Type:       "APIGateway",
				Constraint: "async",
			},
			{
				Type:       "CacheManager",
				Constraint: "thread-safe",
			},
		},
	}
}

// demonstrateEntityAnalysis shows entity frequency analysis features.
func demonstrateEntityAnalysis(analysis *semantic.Analysis) {
	fmt.Println("=== Entity Frequency Analysis ===")

	// Calculate entity frequencies
	freqs := analyze.CalculateEntityFrequency(analysis)
	fmt.Printf("Total unique entities: %d\n", len(freqs))

	// Create frequency list for manipulation
	list := analyze.FrequencyList(freqs)

	// Sort by frequency (descending)
	list.SortByFrequency()
	fmt.Println("\nTop 5 entities by frequency:")
	for i, ef := range list.TopN(5) {
		fmt.Printf("%d. %s (%s): %d occurrences (importance: %d/10)\n",
			i+1, ef.Text, ef.Type, ef.Count, ef.Importance)
	}

	// Get entities by type
	peopleList := list.FilterByType(semantic.EntityPerson)
	fmt.Printf("\nPeople entities: %d\n", len(peopleList))
	for _, ef := range peopleList {
		fmt.Printf("  - %s: %d occurrences\n", ef.Text, ef.Count)
	}

	// Get type distribution
	typeCounts := analyze.GetTypeCounts(freqs)
	fmt.Println("\nEntity type distribution:")
	for entityType, count := range typeCounts {
		fmt.Printf("  %s: %d\n", entityType, count)
	}

	// Calculate bar chart data
	barData := analyze.CalculateBarChartData(freqs)
	fmt.Printf("\nTotal entity occurrences: %d\n", barData.TotalCount)
	fmt.Printf("Maximum single entity count: %d\n", barData.MaxCount)

	fmt.Println()
}

// demonstratePatternMining shows relationship pattern mining features.
func demonstratePatternMining(analysis *semantic.Analysis) {
	fmt.Println("=== Relationship Pattern Mining ===")

	// Mine patterns with default options
	patterns := analyze.MinePatterns(analysis)
	fmt.Printf("Discovered %d relationship patterns\n", len(patterns))

	// Display top patterns
	fmt.Println("\nTop 5 patterns by strength:")
	for i, pattern := range patterns {
		if i >= 5 {
			break
		}
		fmt.Printf("\nPattern %d: [%s] → %s → [%s]\n",
			i+1, pattern.SubjectType, pattern.Predicate, pattern.ObjectType)
		fmt.Printf("  Occurrences: %d\n", pattern.Occurrences)
		fmt.Printf("  Avg Confidence: %.2f\n", pattern.AvgConfidence)
		fmt.Printf("  Strength: %.3f\n", pattern.Strength)

		if len(pattern.Examples) > 0 {
			fmt.Println("  Examples:")
			for j, ex := range pattern.Examples {
				if j >= 3 { // Limit to 3 examples
					break
				}
				fmt.Printf("    • %s → %s → %s (conf: %.2f)\n",
					ex.Subject, ex.Predicate, ex.Object, ex.Confidence)
			}
		}
	}

	// Mine with custom options
	fmt.Println("\n--- Mining with custom options ---")
	opts := analyze.DefaultMiningOptions()
	opts.MinOccurrences = 2      // At least 2 occurrences
	opts.MinConfidence = 0.7     // High confidence only
	opts.MaxExamples = 5         // Store up to 5 examples
	customPatterns := analyze.MinePatternsWithOptions(analysis, opts)
	fmt.Printf("Found %d patterns with custom filters\n", len(customPatterns))

	// Calculate pattern statistics
	stats := analyze.CalculatePatternStats(patterns)
	fmt.Printf("\nPattern Statistics:\n")
	fmt.Printf("  Total patterns: %d\n", stats.TotalPatterns)
	fmt.Printf("  Unique predicates: %d\n", stats.UniquePredicates)
	fmt.Printf("  Avg occurrences per pattern: %.1f\n", stats.AvgOccurrences)
	fmt.Printf("  Avg confidence: %.2f\n", stats.AvgConfidence)

	if len(stats.TopPredicates) > 0 {
		fmt.Println("  Top predicates:")
		for i, pred := range stats.TopPredicates {
			fmt.Printf("    %d. %s\n", i+1, pred)
		}
	}

	// Cluster patterns
	clusters := analyze.ClusterPatterns(patterns, 0.7)
	fmt.Printf("\nClustered into %d groups\n", len(clusters))
	for i, cluster := range clusters {
		if i >= 3 { // Show first 3 clusters
			break
		}
		fmt.Printf("\nCluster %d: %s\n", i+1, cluster.ClusterLabel)
		fmt.Printf("  Similar predicates: %v\n", cluster.Predicates)
		fmt.Printf("  Pattern count: %d\n", len(cluster.Patterns))
		fmt.Printf("  Cluster strength: %.3f\n", cluster.Strength)
	}

	fmt.Println()
}

// demonstrateHoleAnalysis shows typed hole prioritization features.
func demonstrateHoleAnalysis(analysis *semantic.Analysis) {
	fmt.Println("=== Typed Hole Prioritization ===")

	// Analyze typed holes
	holeAnalysis := analyze.AnalyzeTypedHoles(analysis)
	fmt.Printf("Found %d typed holes\n", len(holeAnalysis.Holes))

	// Display priority order
	fmt.Println("\nRecommended implementation order:")
	for i, hole := range holeAnalysis.ImplementOrder {
		if i >= 5 { // Show top 5
			break
		}
		// Create priority bar (10 blocks max)
		priorityBar := ""
		for j := 0; j < 10; j++ {
			if j < hole.Priority {
				priorityBar += "█"
			} else {
				priorityBar += "░"
			}
		}
		fmt.Printf("%d. [%s] ??%s (Complexity: %d)\n",
			i+1, priorityBar, hole.Type, hole.Complexity)

		if hole.Constraint != "" {
			fmt.Printf("   Constraints: %s\n", hole.Constraint)
		}
		if len(hole.Dependencies) > 0 {
			fmt.Printf("   Dependencies: %d\n", len(hole.Dependencies))
		}
	}

	// Show statistics
	fmt.Printf("\nTyped Hole Statistics:\n")
	fmt.Printf("  Total complexity: %d\n", holeAnalysis.TotalComplexity)
	fmt.Printf("  Average priority: %.1f\n", holeAnalysis.AveragePriority)

	// Show critical path
	if len(holeAnalysis.CriticalPath) > 0 {
		fmt.Println("\nCritical Path:")
		for i, hole := range holeAnalysis.CriticalPath {
			fmt.Printf("  %d. ??%s (Complexity: %d)\n", i+1, hole.Type, hole.Complexity)
		}
	}

	// Check for circular dependencies
	if len(holeAnalysis.CircularDeps) > 0 {
		fmt.Printf("\n⚠️  Warning: %d circular dependencies detected!\n", len(holeAnalysis.CircularDeps))
		for i, cycle := range holeAnalysis.CircularDeps {
			fmt.Printf("  Cycle %d: %v\n", i+1, cycle)
		}
	} else {
		fmt.Println("\n✓ No circular dependencies detected")
	}

	// Generate implementation roadmap
	fmt.Println("\n--- Implementation Roadmap ---")
	roadmap := analyze.GenerateImplementationRoadmap(holeAnalysis)
	fmt.Println(roadmap)

	fmt.Println()
}

// demonstrateFiltering shows graph filtering features.
func demonstrateFiltering(analysis *semantic.Analysis) {
	fmt.Println("=== Graph Filtering ===")

	// Build full graph
	graph := analyze.BuildFromAnalysis(analysis)
	fmt.Printf("Full graph: %d nodes, %d edges\n", graph.NodeCount(), graph.EdgeCount())

	// Filter by entity types (People and Organizations only)
	fmt.Println("\n--- Filter: People and Organizations only ---")
	filter := analyze.Filter{
		EntityTypes: map[semantic.EntityType]bool{
			semantic.EntityPerson:       true,
			semantic.EntityOrganization: true,
		},
	}
	filtered := graph.ApplyFilter(filter)
	fmt.Printf("Filtered graph: %d nodes, %d edges\n", filtered.NodeCount(), filtered.EdgeCount())

	// Filter by importance
	fmt.Println("\n--- Filter: High importance entities (7+) ---")
	importanceFilter := analyze.Filter{
		MinImportance: 7,
	}
	highImportance := graph.ApplyFilter(importanceFilter)
	fmt.Printf("High importance graph: %d nodes, %d edges\n",
		highImportance.NodeCount(), highImportance.EdgeCount())

	// Filter by search term
	fmt.Println("\n--- Filter: Search for 'tech' ---")
	searchFilter := analyze.Filter{
		SearchTerm: "tech",
	}
	searchResults := graph.ApplyFilter(searchFilter)
	fmt.Printf("Search results: %d nodes, %d edges\n",
		searchResults.NodeCount(), searchResults.EdgeCount())

	// Combined filter
	fmt.Println("\n--- Combined Filter: Technology type + minimum importance ---")
	combinedFilter := analyze.Filter{
		EntityTypes: map[semantic.EntityType]bool{
			semantic.EntityTechnology: true,
		},
		MinImportance: 5,
		MinConfidence: 0.6,
	}
	combined := graph.ApplyFilter(combinedFilter)
	fmt.Printf("Combined filter graph: %d nodes, %d edges\n",
		combined.NodeCount(), combined.EdgeCount())

	// Display filtered nodes
	if combined.NodeCount() > 0 {
		fmt.Println("\nFiltered entities:")
		for id, node := range combined.Nodes {
			fmt.Printf("  - %s: %d occurrences, importance: %d/10\n",
				id, node.Frequency, node.Importance)
		}
	}

	fmt.Println()
}

// demonstrateExport shows export functionality (simplified for demo).
func demonstrateExport(mode *analyze.AnalyzeMode, analysis *semantic.Analysis) {
	fmt.Println("=== Export Demonstration ===")

	// Note: Full export requires the export package which may have dependencies
	// This demonstrates the concepts without actual file I/O

	// Calculate all analysis components
	entityFreqs := analyze.CalculateEntityFrequency(analysis)
	patterns := analyze.MinePatterns(analysis)
	holeAnalysis := analyze.AnalyzeTypedHoles(analysis)

	fmt.Println("\nAnalysis data prepared for export:")
	fmt.Printf("  - %d entity frequencies\n", len(entityFreqs))
	fmt.Printf("  - %d relationship patterns\n", len(patterns))
	fmt.Printf("  - %d typed holes\n", len(holeAnalysis.Holes))

	// Show what would be exported
	fmt.Println("\nExport formats available:")
	fmt.Println("  1. Markdown (.md) - GitHub-flavored with Mermaid diagrams")
	fmt.Println("  2. HTML (.html) - Interactive with Chart.js visualizations")
	fmt.Println("  3. PDF (.pdf) - Professional layout with embedded charts")

	// Generate example export filename
	timestamp := time.Now().Format("20060102_150405")
	fmt.Printf("\nExample export filenames:\n")
	fmt.Printf("  - analysis_%s.md\n", timestamp)
	fmt.Printf("  - analysis_%s.html\n", timestamp)
	fmt.Printf("  - analysis_%s.pdf\n", timestamp)

	// Simulate export directory check
	exportDir := "./analysis-reports"
	if _, err := os.Stat(exportDir); os.IsNotExist(err) {
		fmt.Printf("\nNote: Export directory '%s' would be created\n", exportDir)
	} else {
		fmt.Printf("\nExport directory: %s\n", exportDir)
	}

	fmt.Println("\n✓ Export demonstration complete")
	fmt.Println("  (In a real application, files would be written to disk)")

	fmt.Println()
}

// Helper function to demonstrate graph layout (optional)
func demonstrateLayout() {
	fmt.Println("=== Graph Layout (Advanced) ===")
	fmt.Println("\nForce-directed layout configuration:")
	fmt.Printf("  - Repulsion strength: %.1f\n", analyze.RepulsionStrength)
	fmt.Printf("  - Attraction strength: %.2f\n", analyze.AttractionStrength)
	fmt.Printf("  - Max force: %.1f\n", analyze.MaxForce)
	fmt.Printf("  - Ideal distance: %.1f\n", analyze.IdealDistance)

	fmt.Println("\nLayout process:")
	fmt.Println("  1. Initialize nodes in circle")
	fmt.Println("  2. Apply repulsive forces (all node pairs)")
	fmt.Println("  3. Apply attractive forces (along edges)")
	fmt.Println("  4. Update positions with damping (0.8)")
	fmt.Println("  5. Repeat for 50-100 iterations")

	fmt.Println()
}
