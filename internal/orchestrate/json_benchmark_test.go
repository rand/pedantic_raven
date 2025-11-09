package orchestrate

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// BenchmarkWorkPlanJSONParsing benchmarks parsing WorkPlan JSON.
func BenchmarkWorkPlanJSONParsing(b *testing.B) {
	sizes := []int{5, 10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dtasks", size), func(b *testing.B) {
			jsonData := generateWorkPlanJSON(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var plan WorkPlan
				_ = json.Unmarshal(jsonData, &plan)
			}
		})
	}
}

// BenchmarkWorkPlanValidation benchmarks work plan validation.
func BenchmarkWorkPlanValidation(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dtasks", size), func(b *testing.B) {
			plan := createTestWorkPlan(size, size*2)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = plan.Validate()
			}
		})
	}
}

// BenchmarkTaskJSONMarshaling benchmarks marshaling Task to JSON.
func BenchmarkTaskJSONMarshaling(b *testing.B) {
	task := Task{
		ID:           "task-1",
		Description:  "Test task description",
		Dependencies: []string{"dep-1", "dep-2"},
		Type:         TaskTypeParallel,
		Agent:        AgentExecutor,
		Priority:     5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(task)
	}
}

// BenchmarkTaskJSONUnmarshaling benchmarks unmarshaling Task from JSON.
func BenchmarkTaskJSONUnmarshaling(b *testing.B) {
	jsonData := []byte(`{
		"id": "task-1",
		"description": "Test task description",
		"dependencies": ["dep-1", "dep-2"],
		"type": 0,
		"agent": 3,
		"priority": 5
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var task Task
		_ = json.Unmarshal(jsonData, &task)
	}
}

// BenchmarkWorkPlanJSONMarshaling benchmarks marshaling WorkPlan to JSON.
func BenchmarkWorkPlanJSONMarshaling(b *testing.B) {
	sizes := []int{10, 50, 100}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dtasks", size), func(b *testing.B) {
			plan := createTestWorkPlan(size, size*2)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = json.Marshal(plan)
			}
		})
	}
}

// BenchmarkTaskValidation benchmarks individual task validation.
func BenchmarkTaskValidation(b *testing.B) {
	task := Task{
		ID:           "task-1",
		Description:  "Test task",
		Dependencies: []string{"dep-1", "dep-2", "dep-3"},
		Type:         TaskTypeSequential,
		Priority:     7,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = task.Validate()
	}
}

// BenchmarkCyclicDependencyDetection benchmarks detecting circular dependencies.
func BenchmarkCyclicDependencyDetection(b *testing.B) {
	// Create a plan with potential cycles
	plan := &WorkPlan{
		Version:     "1.0",
		ProjectName: "test",
		Tasks: []Task{
			{ID: "task-1", Description: "Task 1", Dependencies: []string{"task-2"}},
			{ID: "task-2", Description: "Task 2", Dependencies: []string{"task-3"}},
			{ID: "task-3", Description: "Task 3", Dependencies: []string{"task-4"}},
			{ID: "task-4", Description: "Task 4", Dependencies: []string{"task-1"}}, // Cycle
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = plan.Validate()
	}
}

// BenchmarkLargeWorkPlanParsing benchmarks parsing large work plans.
func BenchmarkLargeWorkPlanParsing(b *testing.B) {
	jsonData := generateWorkPlanJSON(500)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var plan WorkPlan
		_ = json.Unmarshal(jsonData, &plan)
	}
}

// BenchmarkDeepDependencyTree benchmarks validation with deep dependency trees.
func BenchmarkDeepDependencyTree(b *testing.B) {
	// Create a deep dependency chain
	numTasks := 50
	plan := &WorkPlan{
		Version:     "1.0",
		ProjectName: "test",
		Tasks:       make([]Task, numTasks),
	}

	// Create linear dependency chain
	for i := 0; i < numTasks; i++ {
		deps := []string{}
		if i > 0 {
			deps = append(deps, fmt.Sprintf("task-%d", i-1))
		}
		plan.Tasks[i] = Task{
			ID:           fmt.Sprintf("task-%d", i),
			Description:  fmt.Sprintf("Task %d", i),
			Dependencies: deps,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = plan.Validate()
	}
}

// BenchmarkAgentEventJSONMarshaling benchmarks marshaling AgentEvent to JSON.
func BenchmarkAgentEventJSONMarshaling(b *testing.B) {
	event := AgentEvent{
		Timestamp:   "2024-01-01T12:00:00Z",
		AgentType:   AgentExecutor,
		EventType:   EventProgress,
		TaskID:      "task-123",
		Message:     "Task in progress",
		Metadata:    map[string]interface{}{"progress": 0.5},
		FromAgentID: "agent-1",
		ToAgentID:   "agent-2",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(event)
	}
}

// BenchmarkAgentEventJSONUnmarshaling benchmarks unmarshaling AgentEvent from JSON.
func BenchmarkAgentEventJSONUnmarshaling(b *testing.B) {
	jsonData := []byte(`{
		"timestamp": "2024-01-01T12:00:00Z",
		"agent_type": 3,
		"event_type": 1,
		"task_id": "task-123",
		"message": "Task in progress",
		"metadata": {"progress": 0.5},
		"from_agent_id": "agent-1",
		"to_agent_id": "agent-2"
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var event AgentEvent
		_ = json.Unmarshal(jsonData, &event)
	}
}

// BenchmarkSessionStateJSONMarshaling benchmarks marshaling SessionState to JSON.
func BenchmarkSessionStateJSONMarshaling(b *testing.B) {
	state := SessionState{
		SessionID:      "session-123",
		StartTime:      "2024-01-01T12:00:00Z",
		Status:         "active",
		TasksCompleted: 15,
		TasksTotal:     50,
		ActiveAgents:   []string{"agent-1", "agent-2", "agent-3"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(state)
	}
}

// BenchmarkMemoryAllocationJSON benchmarks memory allocations during JSON operations.
func BenchmarkMemoryAllocationJSON(b *testing.B) {
	plan := createTestWorkPlan(50, 100)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(plan)
		var parsed WorkPlan
		_ = json.Unmarshal(data, &parsed)
	}
}

// BenchmarkTaskTypeStringConversion benchmarks TaskType string conversion.
func BenchmarkTaskTypeStringConversion(b *testing.B) {
	taskType := TaskTypeSequential

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = taskType.String()
	}
}

// BenchmarkTaskStatusStringConversion benchmarks TaskStatus string conversion.
func BenchmarkTaskStatusStringConversion(b *testing.B) {
	status := TaskStatusActive

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = status.String()
	}
}

// BenchmarkAgentTypeStringConversion benchmarks AgentType string conversion.
func BenchmarkAgentTypeStringConversion(b *testing.B) {
	agentType := AgentExecutor

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = agentType.String()
	}
}

// Helper functions

func generateWorkPlanJSON(numTasks int) []byte {
	var tasks []string
	for i := 0; i < numTasks; i++ {
		deps := []string{}
		if i > 0 {
			for j := 0; j < min(3, i); j++ {
				deps = append(deps, fmt.Sprintf("\"task-%d\"", i-j-1))
			}
		}
		depsStr := strings.Join(deps, ",")

		task := fmt.Sprintf(`{
			"id": "task-%d",
			"description": "Task %d description",
			"dependencies": [%s],
			"type": %d,
			"priority": %d
		}`, i, i, depsStr, i%3, (i%10)+1)
		tasks = append(tasks, task)
	}

	json := fmt.Sprintf(`{
		"version": "1.0",
		"project_name": "benchmark-test",
		"tasks": [%s]
	}`, strings.Join(tasks, ","))

	return []byte(json)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
