package events

import (
	"sync"
	"testing"
)

// BenchmarkBrokerPublish benchmarks event publishing.
func BenchmarkBrokerPublish(b *testing.B) {
	bufferSizes := []int{10, 100, 1000}

	for _, bufSize := range bufferSizes {
		b.Run(string(rune(bufSize))+"buffer", func(b *testing.B) {
			broker := NewBroker(bufSize)
			event := NewEvent(SemanticAnalysisComplete, SemanticAnalysisData{
				BufferID: "test-buffer",
				Progress: 1.0,
			})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				broker.Publish(event)
			}
		})
	}
}

// BenchmarkBrokerPublishWithSubscribers benchmarks publishing with active subscribers.
func BenchmarkBrokerPublishWithSubscribers(b *testing.B) {
	subscriberCounts := []int{1, 5, 10, 50}

	for _, subCount := range subscriberCounts {
		b.Run(string(rune(subCount))+"subscribers", func(b *testing.B) {
			broker := NewBroker(100)
			event := NewEvent(SemanticAnalysisComplete, SemanticAnalysisData{
				BufferID: "test-buffer",
				Progress: 1.0,
			})

			// Create subscribers
			channels := make([]<-chan Event, subCount)
			for i := 0; i < subCount; i++ {
				channels[i] = broker.Subscribe(SemanticAnalysisComplete)
			}

			// Start goroutines to drain channels
			var wg sync.WaitGroup
			for _, ch := range channels {
				wg.Add(1)
				go func(c <-chan Event) {
					defer wg.Done()
					for range c {
						// Drain events
					}
				}(ch)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				broker.Publish(event)
			}
			b.StopTimer()

			// Cleanup
			broker.Clear()
			wg.Wait()
		})
	}
}

// BenchmarkBrokerSubscribe benchmarks subscription creation.
func BenchmarkBrokerSubscribe(b *testing.B) {
	broker := NewBroker(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = broker.Subscribe(SemanticAnalysisComplete)
	}
}

// BenchmarkBrokerUnsubscribe benchmarks unsubscription.
func BenchmarkBrokerUnsubscribe(b *testing.B) {
	broker := NewBroker(100)

	// Pre-create channels to unsubscribe
	channels := make([]<-chan Event, b.N)
	for i := 0; i < b.N; i++ {
		channels[i] = broker.Subscribe(SemanticAnalysisComplete)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		broker.Unsubscribe(channels[i])
	}
}

// BenchmarkBrokerSubscribeAll benchmarks global subscription.
func BenchmarkBrokerSubscribeAll(b *testing.B) {
	broker := NewBroker(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = broker.SubscribeAll()
	}
}

// BenchmarkBrokerPublishConcurrent benchmarks concurrent publishing.
func BenchmarkBrokerPublishConcurrent(b *testing.B) {
	concurrencyLevels := []int{1, 2, 4, 8}

	for _, concurrency := range concurrencyLevels {
		b.Run(string(rune(concurrency))+"goroutines", func(b *testing.B) {
			broker := NewBroker(1000)
			event := NewEvent(SemanticAnalysisComplete, SemanticAnalysisData{
				BufferID: "test-buffer",
				Progress: 1.0,
			})

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					broker.Publish(event)
				}
			})
		})
	}
}

// BenchmarkBrokerPublishMultipleEventTypes benchmarks publishing to different event types.
func BenchmarkBrokerPublishMultipleEventTypes(b *testing.B) {
	broker := NewBroker(100)
	events := []Event{
		NewEvent(SemanticAnalysisComplete, SemanticAnalysisData{BufferID: "buf1", Progress: 1.0}),
		NewEvent(MemoryRecalled, MemoryData{MemoryID: "mem1"}),
		NewEvent(AgentStarted, AgentData{AgentID: "agent1"}),
		NewEvent(FileChanged, FileData{Path: "/test/file.go"}),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		broker.Publish(events[i%len(events)])
	}
}

// BenchmarkBrokerSubscriberCount benchmarks counting subscribers.
func BenchmarkBrokerSubscriberCount(b *testing.B) {
	broker := NewBroker(100)

	// Create some subscribers
	for i := 0; i < 10; i++ {
		broker.Subscribe(SemanticAnalysisComplete)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = broker.SubscriberCount(SemanticAnalysisComplete)
	}
}

// BenchmarkBrokerGlobalSubscriberCount benchmarks counting global subscribers.
func BenchmarkBrokerGlobalSubscriberCount(b *testing.B) {
	broker := NewBroker(100)

	// Create some global subscribers
	for i := 0; i < 10; i++ {
		broker.SubscribeAll()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = broker.GlobalSubscriberCount()
	}
}

// BenchmarkBrokerClear benchmarks clearing all subscribers.
func BenchmarkBrokerClear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		broker := NewBroker(100)
		// Create many subscribers
		for j := 0; j < 100; j++ {
			broker.Subscribe(SemanticAnalysisComplete)
			broker.SubscribeAll()
		}
		b.StartTimer()

		broker.Clear()
	}
}

// BenchmarkBrokerPublishSemanticAnalysis benchmarks convenience method.
func BenchmarkBrokerPublishSemanticAnalysis(b *testing.B) {
	broker := NewBroker(100)
	data := SemanticAnalysisData{
		BufferID: "test-buffer",
		Progress: 1.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		broker.PublishSemanticAnalysis(SemanticAnalysisComplete, data)
	}
}

// BenchmarkBrokerHighThroughput benchmarks high-throughput scenario.
func BenchmarkBrokerHighThroughput(b *testing.B) {
	broker := NewBroker(10000)
	event := NewEvent(SemanticAnalysisComplete, SemanticAnalysisData{
		BufferID: "test-buffer",
		Progress: 1.0,
	})

	// Create subscribers that actively drain
	numSubs := 10
	channels := make([]<-chan Event, numSubs)
	for i := 0; i < numSubs; i++ {
		channels[i] = broker.Subscribe(SemanticAnalysisComplete)
	}

	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan Event) {
			defer wg.Done()
			for range c {
				// Process events
			}
		}(ch)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		broker.Publish(event)
	}
	b.StopTimer()

	broker.Clear()
	wg.Wait()
}

// BenchmarkBrokerMemoryAllocation benchmarks memory allocations.
func BenchmarkBrokerMemoryAllocation(b *testing.B) {
	broker := NewBroker(100)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		event := NewEvent(SemanticAnalysisComplete, SemanticAnalysisData{
			BufferID: "test-buffer",
			Progress: 1.0,
		})
		broker.Publish(event)
	}
}
