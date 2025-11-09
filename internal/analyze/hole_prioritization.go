package analyze

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rand/pedantic-raven/internal/editor/semantic"
)

// HoleDependency represents a dependency relationship between typed holes.
type HoleDependency struct {
	From         string   // ID of dependent hole
	To           string   // ID of dependency
	Relationship string   // Type of relationship (requires, extends, implements)
	Strength     int      // Strength of dependency (1-10)
}

// HoleNode represents a node in the dependency tree.
type HoleNode struct {
	ID           string                   // Unique identifier
	Hole         semantic.EnhancedTypedHole // Enhanced hole data
	Children     []*HoleNode               // Child dependencies
	Parents      []*HoleNode               // Parent dependencies
	Depth        int                       // Depth in dependency tree
	CriticalPath bool                      // Whether this is on critical path
}

// HoleAnalysis aggregates all typed hole analysis data.
type HoleAnalysis struct {
	Holes            []semantic.EnhancedTypedHole // All enhanced typed holes
	Dependencies     []HoleDependency             // Dependency relationships
	DependencyTree   *HoleNode                    // Root of dependency tree
	ImplementOrder   []semantic.EnhancedTypedHole // Recommended implementation order
	CriticalPath     []semantic.EnhancedTypedHole // Critical path holes
	TotalComplexity  int                          // Sum of all hole complexities
	AveragePriority  float64                      // Average priority
	CircularDeps     [][]string                   // Circular dependency chains
}

// AnalyzeTypedHoles performs comprehensive typed hole analysis.
func AnalyzeTypedHoles(analysis *semantic.Analysis) *HoleAnalysis {
	if analysis == nil || len(analysis.TypedHoles) == 0 {
		return &HoleAnalysis{
			Holes:           []semantic.EnhancedTypedHole{},
			Dependencies:    []HoleDependency{},
			ImplementOrder:  []semantic.EnhancedTypedHole{},
			CriticalPath:    []semantic.EnhancedTypedHole{},
			CircularDeps:    [][]string{},
		}
	}

	// Create prioritizer to get enhanced holes
	prioritizer := semantic.NewHolePrioritizer(analysis.TypedHoles, analysis.Relationships)
	enhancedHoles := prioritizer.Holes

	// Build dependency relationships
	dependencies := buildDependencies(enhancedHoles, analysis.Relationships)

	// Build dependency tree
	tree := buildDependencyTree(enhancedHoles, dependencies)

	// Detect circular dependencies
	circular := detectCircularDependencies(enhancedHoles, dependencies)

	// Calculate implementation order (topological sort with priority)
	implementOrder := calculateImplementationOrder(enhancedHoles, dependencies)

	// Identify critical path
	criticalPath := identifyCriticalPath(tree, implementOrder)

	// Calculate aggregate metrics
	totalComplexity := 0
	totalPriority := 0
	for _, hole := range enhancedHoles {
		totalComplexity += hole.Complexity
		totalPriority += hole.Priority
	}

	avgPriority := 0.0
	if len(enhancedHoles) > 0 {
		avgPriority = float64(totalPriority) / float64(len(enhancedHoles))
	}

	return &HoleAnalysis{
		Holes:           enhancedHoles,
		Dependencies:    dependencies,
		DependencyTree:  tree,
		ImplementOrder:  implementOrder,
		CriticalPath:    criticalPath,
		TotalComplexity: totalComplexity,
		AveragePriority: avgPriority,
		CircularDeps:    circular,
	}
}

// buildDependencies extracts dependency relationships from holes and relationships.
func buildDependencies(holes []semantic.EnhancedTypedHole, relationships []semantic.Relationship) []HoleDependency {
	deps := []HoleDependency{}

	// Initialize hole metadata
	for i := range holes {
		holes[i].RelatedHoles = []string{}
		holes[i].Dependencies = []string{}
	}

	// Analyze relationships for dependencies
	for _, rel := range relationships {
		// Look for dependency patterns in relationships
		// Pattern 1: "X requires Y"
		if strings.Contains(strings.ToLower(rel.Predicate), "require") {
			// Find holes mentioned in subject and object
			for i, hole1 := range holes {
				for j, hole2 := range holes {
					if i != j {
						if strings.Contains(rel.Subject, hole1.Type) && strings.Contains(rel.Object, hole2.Type) {
							id1 := fmt.Sprintf("%s_%d", hole1.Type, i)
							id2 := fmt.Sprintf("%s_%d", hole2.Type, j)

							deps = append(deps, HoleDependency{
								From:         id1,
								To:           id2,
								Relationship: "requires",
								Strength:     8,
							})

							// Track in enhanced holes
							if !contains(holes[i].Dependencies, id2) {
								holes[i].Dependencies = append(holes[i].Dependencies, id2)
							}
							if !contains(holes[i].RelatedHoles, id2) {
								holes[i].RelatedHoles = append(holes[i].RelatedHoles, id2)
							}
						}
					}
				}
			}
		}

		// Pattern 2: "X extends Y" or "X implements Y"
		if strings.Contains(strings.ToLower(rel.Predicate), "extend") ||
		   strings.Contains(strings.ToLower(rel.Predicate), "implement") {
			relType := "extends"
			if strings.Contains(strings.ToLower(rel.Predicate), "implement") {
				relType = "implements"
			}

			for i, hole1 := range holes {
				for j, hole2 := range holes {
					if i != j {
						if strings.Contains(rel.Subject, hole1.Type) && strings.Contains(rel.Object, hole2.Type) {
							id1 := fmt.Sprintf("%s_%d", hole1.Type, i)
							id2 := fmt.Sprintf("%s_%d", hole2.Type, j)

							deps = append(deps, HoleDependency{
								From:         id1,
								To:           id2,
								Relationship: relType,
								Strength:     7,
							})

							// Track in enhanced holes
							if !contains(holes[i].Dependencies, id2) {
								holes[i].Dependencies = append(holes[i].Dependencies, id2)
							}
							if !contains(holes[i].RelatedHoles, id2) {
								holes[i].RelatedHoles = append(holes[i].RelatedHoles, id2)
							}
						}
					}
				}
			}
		}

		// Pattern 3: General associations (weaker dependency)
		for i, hole1 := range holes {
			for j, hole2 := range holes {
				if i != j {
					if (strings.Contains(rel.Subject, hole1.Type) && strings.Contains(rel.Object, hole2.Type)) ||
					   (strings.Contains(rel.Subject, hole2.Type) && strings.Contains(rel.Object, hole1.Type)) {
						id1 := fmt.Sprintf("%s_%d", hole1.Type, i)
						id2 := fmt.Sprintf("%s_%d", hole2.Type, j)

						// Only add if not already added as stronger dependency
						if !hasDependency(deps, id1, id2) {
							deps = append(deps, HoleDependency{
								From:         id1,
								To:           id2,
								Relationship: "related",
								Strength:     3,
							})

							// Track related holes
							if !contains(holes[i].RelatedHoles, id2) {
								holes[i].RelatedHoles = append(holes[i].RelatedHoles, id2)
							}
						}
					}
				}
			}
		}
	}

	return deps
}

// buildDependencyTree constructs a tree representation of hole dependencies.
func buildDependencyTree(holes []semantic.EnhancedTypedHole, deps []HoleDependency) *HoleNode {
	// Create nodes for all holes
	nodes := make(map[string]*HoleNode)
	for i, hole := range holes {
		id := fmt.Sprintf("%s_%d", hole.Type, i)
		nodes[id] = &HoleNode{
			ID:       id,
			Hole:     hole,
			Children: []*HoleNode{},
			Parents:  []*HoleNode{},
			Depth:    0,
		}
	}

	// Build parent-child relationships from dependencies
	for _, dep := range deps {
		if dep.Relationship == "requires" || dep.Relationship == "extends" || dep.Relationship == "implements" {
			if fromNode, ok := nodes[dep.From]; ok {
				if toNode, ok := nodes[dep.To]; ok {
					// "From requires To" means To is a parent of From
					fromNode.Parents = append(fromNode.Parents, toNode)
					toNode.Children = append(toNode.Children, fromNode)
				}
			}
		}
	}

	// Calculate depth for each node
	var calculateDepth func(*HoleNode, int, map[string]bool)
	calculateDepth = func(node *HoleNode, depth int, visited map[string]bool) {
		if visited[node.ID] {
			return // Avoid cycles
		}
		visited[node.ID] = true

		if depth > node.Depth {
			node.Depth = depth
		}

		for _, child := range node.Children {
			calculateDepth(child, depth+1, visited)
		}
	}

	// Find root nodes (no parents) and calculate depths
	for _, node := range nodes {
		if len(node.Parents) == 0 {
			calculateDepth(node, 0, make(map[string]bool))
		}
	}

	// Create virtual root connecting all actual roots
	root := &HoleNode{
		ID:       "root",
		Children: []*HoleNode{},
		Depth:    -1,
	}

	for _, node := range nodes {
		if len(node.Parents) == 0 {
			root.Children = append(root.Children, node)
		}
	}

	return root
}

// detectCircularDependencies finds circular dependency chains.
func detectCircularDependencies(holes []semantic.EnhancedTypedHole, deps []HoleDependency) [][]string {
	// Build adjacency list
	graph := make(map[string][]string)
	for i, hole := range holes {
		id := fmt.Sprintf("%s_%d", hole.Type, i)
		graph[id] = []string{}
	}

	for _, dep := range deps {
		if dep.Relationship == "requires" || dep.Relationship == "extends" || dep.Relationship == "implements" {
			graph[dep.From] = append(graph[dep.From], dep.To)
		}
	}

	// DFS to detect cycles
	cycles := [][]string{}
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := []string{}

	var dfs func(string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				// Found cycle - extract it from path
				cycleStart := -1
				for i, n := range path {
					if n == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := make([]string, len(path)-cycleStart)
					copy(cycle, path[cycleStart:])
					cycles = append(cycles, cycle)
				}
				return true
			}
		}

		path = path[:len(path)-1]
		recStack[node] = false
		return false
	}

	for id := range graph {
		if !visited[id] {
			dfs(id)
		}
	}

	return cycles
}

// calculateImplementationOrder determines the recommended order for implementing holes.
// Uses topological sort combined with priority scoring.
func calculateImplementationOrder(holes []semantic.EnhancedTypedHole, deps []HoleDependency) []semantic.EnhancedTypedHole {
	// Build in-degree map
	inDegree := make(map[string]int)
	adjList := make(map[string][]string)
	holeMap := make(map[string]semantic.EnhancedTypedHole)

	for i, hole := range holes {
		id := fmt.Sprintf("%s_%d", hole.Type, i)
		inDegree[id] = 0
		adjList[id] = []string{}
		holeMap[id] = hole
	}

	// Count in-degrees from strong dependencies
	for _, dep := range deps {
		if dep.Relationship == "requires" || dep.Relationship == "extends" || dep.Relationship == "implements" {
			inDegree[dep.From]++
			adjList[dep.To] = append(adjList[dep.To], dep.From)
		}
	}

	// Topological sort with priority
	result := []semantic.EnhancedTypedHole{}
	queue := []string{}

	// Start with nodes that have no dependencies
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}

	// Sort queue by priority/complexity ratio (highest first)
	sort.Slice(queue, func(i, j int) bool {
		hole1 := holeMap[queue[i]]
		hole2 := holeMap[queue[j]]
		score1 := float64(hole1.Priority) / float64(max(hole1.Complexity, 1))
		score2 := float64(hole2.Priority) / float64(max(hole2.Complexity, 1))
		return score2 < score1
	})

	for len(queue) > 0 {
		// Pop from queue
		current := queue[0]
		queue = queue[1:]

		result = append(result, holeMap[current])

		// Process neighbors
		for _, neighbor := range adjList[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)

				// Re-sort queue by priority
				sort.Slice(queue, func(i, j int) bool {
					hole1 := holeMap[queue[i]]
					hole2 := holeMap[queue[j]]
					score1 := float64(hole1.Priority) / float64(max(hole1.Complexity, 1))
					score2 := float64(hole2.Priority) / float64(max(hole2.Complexity, 1))
					return score2 < score1
				})
			}
		}
	}

	return result
}

// identifyCriticalPath finds the longest path through the dependency tree.
func identifyCriticalPath(tree *HoleNode, implementOrder []semantic.EnhancedTypedHole) []semantic.EnhancedTypedHole {
	if tree == nil || len(tree.Children) == 0 {
		return []semantic.EnhancedTypedHole{}
	}

	// Find longest path by complexity
	var findLongestPath func(*HoleNode, []semantic.EnhancedTypedHole) []semantic.EnhancedTypedHole
	findLongestPath = func(node *HoleNode, currentPath []semantic.EnhancedTypedHole) []semantic.EnhancedTypedHole {
		if node.ID == "root" {
			// Process children of root
			longestPath := []semantic.EnhancedTypedHole{}
			for _, child := range node.Children {
				path := findLongestPath(child, currentPath)
				if calculatePathComplexity(path) > calculatePathComplexity(longestPath) {
					longestPath = path
				}
			}
			return longestPath
		}

		newPath := append(currentPath, node.Hole)

		if len(node.Children) == 0 {
			return newPath
		}

		// Recursively find longest path through children
		longestPath := newPath
		for _, child := range node.Children {
			path := findLongestPath(child, newPath)
			if calculatePathComplexity(path) > calculatePathComplexity(longestPath) {
				longestPath = path
			}
		}

		return longestPath
	}

	path := findLongestPath(tree, []semantic.EnhancedTypedHole{})

	// Mark nodes on critical path
	for _, pathNode := range path {
		markCriticalPath(tree, pathNode)
	}

	return path
}

// calculatePathComplexity sums complexity along a path.
func calculatePathComplexity(path []semantic.EnhancedTypedHole) int {
	total := 0
	for _, hole := range path {
		total += hole.Complexity
	}
	return total
}

// markCriticalPath marks nodes on the critical path.
func markCriticalPath(node *HoleNode, target semantic.EnhancedTypedHole) bool {
	if node.ID == "root" {
		for _, child := range node.Children {
			if markCriticalPath(child, target) {
				return true
			}
		}
		return false
	}

	if node.Hole.Type == target.Type && node.Hole.Priority == target.Priority {
		node.CriticalPath = true
		return true
	}

	for _, child := range node.Children {
		if markCriticalPath(child, target) {
			node.CriticalPath = true
			return true
		}
	}

	return false
}

// GenerateImplementationRoadmap creates a structured roadmap for implementing typed holes.
func GenerateImplementationRoadmap(analysis *HoleAnalysis) string {
	var sb strings.Builder

	sb.WriteString("=== Typed Hole Implementation Roadmap ===\n\n")

	if len(analysis.CircularDeps) > 0 {
		sb.WriteString("WARNING: Circular dependencies detected!\n")
		for i, cycle := range analysis.CircularDeps {
			sb.WriteString(fmt.Sprintf("  Cycle %d: %s\n", i+1, strings.Join(cycle, " -> ")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Total Holes: %d\n", len(analysis.Holes)))
	sb.WriteString(fmt.Sprintf("Total Complexity: %d\n", analysis.TotalComplexity))
	sb.WriteString(fmt.Sprintf("Average Priority: %.1f\n\n", analysis.AveragePriority))

	// Group by milestones (every 3-4 holes or complexity threshold)
	milestones := groupIntoMilestones(analysis.ImplementOrder)

	for i, milestone := range milestones {
		sb.WriteString(fmt.Sprintf("Milestone %d:\n", i+1))
		for j, hole := range milestone {
			sb.WriteString(fmt.Sprintf("  %d. %s (Priority: %d, Complexity: %d)\n",
				j+1, hole.Type, hole.Priority, hole.Complexity))
			if hole.Constraint != "" {
				sb.WriteString(fmt.Sprintf("     Constraints: %s\n", hole.Constraint))
			}
			if hole.SuggestedImpl != "" {
				sb.WriteString(fmt.Sprintf("     Suggestion: %s\n", hole.SuggestedImpl))
			}
			if len(hole.Dependencies) > 0 {
				sb.WriteString(fmt.Sprintf("     Dependencies: %d\n", len(hole.Dependencies)))
			}
		}
		sb.WriteString("\n")
	}

	if len(analysis.CriticalPath) > 0 {
		sb.WriteString("Critical Path:\n")
		for i, hole := range analysis.CriticalPath {
			sb.WriteString(fmt.Sprintf("  %d. %s (Complexity: %d)\n", i+1, hole.Type, hole.Complexity))
		}
		totalCritical := calculatePathComplexity(analysis.CriticalPath)
		sb.WriteString(fmt.Sprintf("\nCritical Path Total Complexity: %d\n", totalCritical))
	}

	return sb.String()
}

// groupIntoMilestones groups holes into implementation milestones.
func groupIntoMilestones(holes []semantic.EnhancedTypedHole) [][]semantic.EnhancedTypedHole {
	const complexityThreshold = 25
	const countThreshold = 4

	milestones := [][]semantic.EnhancedTypedHole{}
	currentMilestone := []semantic.EnhancedTypedHole{}
	currentComplexity := 0

	for _, hole := range holes {
		currentMilestone = append(currentMilestone, hole)
		currentComplexity += hole.Complexity

		if currentComplexity >= complexityThreshold || len(currentMilestone) >= countThreshold {
			milestones = append(milestones, currentMilestone)
			currentMilestone = []semantic.EnhancedTypedHole{}
			currentComplexity = 0
		}
	}

	// Add remaining holes
	if len(currentMilestone) > 0 {
		milestones = append(milestones, currentMilestone)
	}

	return milestones
}

// Helper functions

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func hasDependency(deps []HoleDependency, from, to string) bool {
	for _, dep := range deps {
		if dep.From == from && dep.To == to {
			return true
		}
	}
	return false
}
