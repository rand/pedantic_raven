package mnemosyne

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/rand/pedantic-raven/internal/mnemosyne/pb/mnemosyne/v1"
)

// BenchmarkRecallOptionsCreation benchmarks creating recall options.
func BenchmarkRecallOptionsCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RecallOptions{
			Query:      "test query",
			Namespace:  GlobalNamespace(),
			MaxResults: 10,
		}
	}
}

// BenchmarkStoreMemoryOptionsCreation benchmarks creating store memory options.
func BenchmarkStoreMemoryOptionsCreation(b *testing.B) {
	importance := uint32(8)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = StoreMemoryOptions{
			Content:    "test content",
			Namespace:  ProjectNamespace("test-project"),
			Importance: &importance,
			Tags:       []string{"test", "benchmark"},
		}
	}
}

// BenchmarkNamespaceCreation benchmarks namespace creation functions.
func BenchmarkNamespaceCreation(b *testing.B) {
	b.Run("Global", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GlobalNamespace()
		}
	})

	b.Run("Project", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = ProjectNamespace("test-project")
		}
	})

	b.Run("Session", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = SessionNamespace("test-project", "session-123")
		}
	})
}

// BenchmarkRecallValidation benchmarks recall request validation.
func BenchmarkRecallValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail with ErrNotConnected, but exercises validation logic
		_, _ = client.Recall(ctx, RecallOptions{
			Query:      "test query",
			Namespace:  GlobalNamespace(),
			MaxResults: 10,
		})
	}
}

// BenchmarkStoreMemoryValidation benchmarks store memory validation.
func BenchmarkStoreMemoryValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()
	importance := uint32(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.StoreMemory(ctx, StoreMemoryOptions{
			Content:    "test content",
			Namespace:  GlobalNamespace(),
			Importance: &importance,
		})
	}
}

// BenchmarkListMemoriesOptionsCreation benchmarks list options creation.
func BenchmarkListMemoriesOptionsCreation(b *testing.B) {
	minImportance := uint32(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ListMemoriesOptions{
			Namespace:     ProjectNamespace("test"),
			Tags:          []string{"important", "work"},
			MinImportance: &minImportance,
			MaxResults:    50,
		}
	}
}

// BenchmarkGraphTraverseOptionsCreation benchmarks graph traverse options.
func BenchmarkGraphTraverseOptionsCreation(b *testing.B) {
	minStrength := float32(0.5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GraphTraverseOptions{
			SeedIDs:         []string{"mem1", "mem2", "mem3"},
			MaxHops:         3,
			MinLinkStrength: &minStrength,
		}
	}
}

// BenchmarkSemanticSearchOptionsCreation benchmarks semantic search options.
func BenchmarkSemanticSearchOptionsCreation(b *testing.B) {
	embedding := make([]float32, 768) // Typical embedding size
	for i := range embedding {
		embedding[i] = 0.1
	}
	minImportance := uint32(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SemanticSearchOptions{
			Embedding:     embedding,
			Namespace:     GlobalNamespace(),
			MaxResults:    20,
			MinImportance: &minImportance,
		}
	}
}

// BenchmarkUpdateMemoryOptionsCreation benchmarks update memory options.
func BenchmarkUpdateMemoryOptionsCreation(b *testing.B) {
	content := "updated content"
	importance := uint32(7)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = UpdateMemoryOptions{
			MemoryID:   "mem-123",
			Content:    &content,
			Importance: &importance,
			AddTags:    []string{"new-tag"},
			RemoveTags: []string{"old-tag"},
		}
	}
}

// BenchmarkClientConfigValidation benchmarks config validation.
func BenchmarkClientConfigValidation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg := DefaultConfig()
		_, _ = NewClient(cfg)
	}
}

// BenchmarkMemoryAllocationRecall benchmarks memory allocations during recall.
func BenchmarkMemoryAllocationRecall(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = client.Recall(ctx, RecallOptions{
			Query:      fmt.Sprintf("test query %d", i),
			Namespace:  GlobalNamespace(),
			MaxResults: 10,
		})
	}
}

// BenchmarkMemoryAllocationStore benchmarks memory allocations during store.
func BenchmarkMemoryAllocationStore(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()
	importance := uint32(5)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = client.StoreMemory(ctx, StoreMemoryOptions{
			Content:    fmt.Sprintf("test content %d", i),
			Namespace:  GlobalNamespace(),
			Importance: &importance,
		})
	}
}

// BenchmarkTagCreation benchmarks creating tag slices.
func BenchmarkTagCreation(b *testing.B) {
	tagCounts := []int{1, 5, 10, 20}

	for _, count := range tagCounts {
		b.Run(fmt.Sprintf("%dtags", count), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tags := make([]string, count)
				for j := 0; j < count; j++ {
					tags[j] = fmt.Sprintf("tag-%d", j)
				}
			}
		})
	}
}

// BenchmarkGetMemoryValidation benchmarks get memory validation.
func BenchmarkGetMemoryValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetMemory(ctx, "mem-123")
	}
}

// BenchmarkDeleteMemoryValidation benchmarks delete memory validation.
func BenchmarkDeleteMemoryValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.DeleteMemory(ctx, "mem-123")
	}
}

// BenchmarkListMemoriesValidation benchmarks list memories validation.
func BenchmarkListMemoriesValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.ListMemories(ctx, ListMemoriesOptions{
			Namespace:  GlobalNamespace(),
			MaxResults: 100,
		})
	}
}

// BenchmarkGetContextValidation benchmarks get context validation.
func BenchmarkGetContextValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetContext(ctx, GetContextOptions{
			MemoryIDs:      []string{"mem-1", "mem-2"},
			IncludeLinks:   true,
			MaxLinkedDepth: 2,
		})
	}
}

// BenchmarkGraphTraverseValidation benchmarks graph traverse validation.
func BenchmarkGraphTraverseValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = client.GraphTraverse(ctx, GraphTraverseOptions{
			SeedIDs: []string{"mem-1", "mem-2", "mem-3"},
			MaxHops: 2,
		})
	}
}

// BenchmarkSemanticSearchValidation benchmarks semantic search validation.
func BenchmarkSemanticSearchValidation(b *testing.B) {
	client, _ := NewClient(DefaultConfig())
	ctx := context.Background()
	embedding := make([]float32, 768)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.SemanticSearch(ctx, SemanticSearchOptions{
			Embedding:  embedding,
			Namespace:  GlobalNamespace(),
			MaxResults: 10,
		})
	}
}

// BenchmarkEmbeddingCreation benchmarks creating embedding vectors.
func BenchmarkEmbeddingCreation(b *testing.B) {
	sizes := []int{384, 768, 1536}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("dim%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				embedding := make([]float32, size)
				for j := range embedding {
					embedding[j] = float32(j) * 0.001
				}
			}
		})
	}
}

// BenchmarkMultipleMemoryIDs benchmarks creating memory ID slices.
func BenchmarkMultipleMemoryIDs(b *testing.B) {
	counts := []int{1, 5, 10, 50}

	for _, count := range counts {
		b.Run(fmt.Sprintf("%dids", count), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ids := make([]string, count)
				for j := 0; j < count; j++ {
					ids[j] = fmt.Sprintf("mem-%d", j)
				}
			}
		})
	}
}
