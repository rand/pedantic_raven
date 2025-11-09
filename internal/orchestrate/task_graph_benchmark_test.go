package orchestrate

import (
	"fmt"
	"testing"
)

// BenchmarkTaskGraphCreation benchmarks creating task graphs of various sizes.
func BenchmarkTaskGraphCreation(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dnodes", size), func(b *testing.B) {
			plan := createBenchmarkWorkPlan(size, size*2)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = NewTaskGraph(plan, 80, 40)
			}
		})
	}
}

// BenchmarkTaskGraphLayout benchmarks force-directed layout iterations.
func BenchmarkTaskGraphLayout(b *testing.B) {
	sizes := []struct {
		name  string
		nodes int
		edges int
	}{
		{"10nodes", 10, 15},
		{"50nodes", 50, 80},
		{"100nodes", 100, 200},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			plan := createBenchmarkWorkPlan(size.nodes, size.edges)
			graph, err := NewTaskGraph(plan, 80, 40)
			if err != nil {
				b.Fatalf("Failed to create task graph: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.applyForceIteration()
			}
		})
	}
}

// BenchmarkTaskGraphStabilize benchmarks complete layout stabilization.
func BenchmarkTaskGraphStabilize(b *testing.B) {
	sizes := []struct {
		name       string
		nodes      int
		edges      int
		iterations int
	}{
		{"10nodes_50iter", 10, 15, 50},
		{"50nodes_50iter", 50, 80, 50},
		{"100nodes_50iter", 100, 200, 50},
		{"100nodes_100iter", 100, 200, 100},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				plan := createBenchmarkWorkPlan(size.nodes, size.edges)
				graph, err := NewTaskGraph(plan, 80, 40)
				if err != nil {
					b.Fatalf("Failed to create task graph: %v", err)
				}
				b.StartTimer()

				graph.stabilize(size.iterations)
			}
		})
	}
}

// BenchmarkTaskGraphRepulsion benchmarks repulsion force calculation.
func BenchmarkTaskGraphRepulsion(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dnodes", size), func(b *testing.B) {
			plan := createBenchmarkWorkPlan(size, size*2)
			graph, err := NewTaskGraph(plan, 80, 40)
			if err != nil {
				b.Fatalf("Failed to create task graph: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.applyRepulsion()
			}
		})
	}
}

// BenchmarkTaskGraphAttraction benchmarks attraction force calculation.
func BenchmarkTaskGraphAttraction(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dnodes", size), func(b *testing.B) {
			plan := createBenchmarkWorkPlan(size, size*2)
			graph, err := NewTaskGraph(plan, 80, 40)
			if err != nil {
				b.Fatalf("Failed to create task graph: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.applyAttraction()
			}
		})
	}
}

// BenchmarkTaskGraphUpdatePositions benchmarks position updates.
func BenchmarkTaskGraphUpdatePositions(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dnodes", size), func(b *testing.B) {
			plan := createBenchmarkWorkPlan(size, size*2)
			graph, err := NewTaskGraph(plan, 80, 40)
			if err != nil {
				b.Fatalf("Failed to create task graph: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				graph.updatePositions()
			}
		})
	}
}

// BenchmarkTaskGraphUpdateStatus benchmarks status updates.
func BenchmarkTaskGraphUpdateStatus(b *testing.B) {
	plan := createBenchmarkWorkPlan(100, 200)
	graph, err := NewTaskGraph(plan, 80, 40)
	if err != nil {
		b.Fatalf("Failed to create task graph: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskID := fmt.Sprintf("task-%d", i%100)
		graph.UpdateStatus(taskID, TaskStatusActive)
	}
}

// BenchmarkTaskGraphSelectNode benchmarks node selection.
func BenchmarkTaskGraphSelectNode(b *testing.B) {
	plan := createBenchmarkWorkPlan(100, 200)
	graph, err := NewTaskGraph(plan, 80, 40)
	if err != nil {
		b.Fatalf("Failed to create task graph: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskID := fmt.Sprintf("task-%d", i%100)
		graph.SelectNode(taskID)
	}
}

// BenchmarkTaskGraphGetBounds benchmarks bounding box calculation.
func BenchmarkTaskGraphGetBounds(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dnodes", size), func(b *testing.B) {
			plan := createBenchmarkWorkPlan(size, size*2)
			graph, err := NewTaskGraph(plan, 80, 40)
			if err != nil {
				b.Fatalf("Failed to create task graph: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _, _ = graph.getBounds()
			}
		})
	}
}

// BenchmarkTaskGraphResize benchmarks viewport resizing.
func BenchmarkTaskGraphResize(b *testing.B) {
	plan := createBenchmarkWorkPlan(50, 100)
	graph, err := NewTaskGraph(plan, 80, 40)
	if err != nil {
		b.Fatalf("Failed to create task graph: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Resize(100, 50)
	}
}

