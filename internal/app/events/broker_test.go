package events

import (
	"testing"
	"time"
)

func TestBrokerBasicPubSub(t *testing.T) {
	broker := NewBroker(10)
	defer broker.Clear()

	// Subscribe to semantic analysis events
	eventCh := broker.Subscribe(SemanticAnalysisComplete)
	defer broker.Unsubscribe(eventCh)

	// Publish an event
	data := SemanticAnalysisData{
		BufferID: "test-buffer",
		Progress: 1.0,
	}
	broker.PublishSemanticAnalysis(SemanticAnalysisComplete, data)

	// Receive the event
	select {
	case event := <-eventCh:
		if event.Type != SemanticAnalysisComplete {
			t.Errorf("Expected event type %v, got %v", SemanticAnalysisComplete, event.Type)
		}

		eventData, ok := event.Data.(SemanticAnalysisData)
		if !ok {
			t.Fatalf("Expected SemanticAnalysisData, got %T", event.Data)
		}

		if eventData.BufferID != "test-buffer" {
			t.Errorf("Expected BufferID 'test-buffer', got '%s'", eventData.BufferID)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for event")
	}
}

func TestBrokerMultipleSubscribers(t *testing.T) {
	broker := NewBroker(10)
	defer broker.Clear()

	// Create multiple subscribers for the same event type
	eventCh1 := broker.Subscribe(MemoryCreated)
	defer broker.Unsubscribe(eventCh1)

	eventCh2 := broker.Subscribe(MemoryCreated)
	defer broker.Unsubscribe(eventCh2)

	// Publish an event
	data := MemoryData{
		MemoryID: "mem-123",
	}
	broker.PublishMemory(MemoryCreated, data)

	// Both subscribers should receive the event
	received := 0
	timeout := time.After(100 * time.Millisecond)

	for received < 2 {
		select {
		case event := <-eventCh1:
			if event.Type != MemoryCreated {
				t.Errorf("Subscriber 1: Expected event type %v, got %v", MemoryCreated, event.Type)
			}
			received++
		case event := <-eventCh2:
			if event.Type != MemoryCreated {
				t.Errorf("Subscriber 2: Expected event type %v, got %v", MemoryCreated, event.Type)
			}
			received++
		case <-timeout:
			t.Fatalf("Timeout: only received %d/2 events", received)
		}
	}
}

func TestBrokerGlobalSubscription(t *testing.T) {
	broker := NewBroker(10)
	defer broker.Clear()

	// Subscribe to all events
	eventCh := broker.SubscribeAll()
	defer broker.Unsubscribe(eventCh)

	// Publish events of different types
	broker.PublishMemory(MemoryCreated, MemoryData{MemoryID: "mem-1"})
	broker.PublishAgent(AgentStarted, AgentData{AgentID: "agent-1"})
	broker.PublishUI(ModeChanged, UIData{Mode: "edit"})

	// Should receive all 3 events
	received := 0
	timeout := time.After(100 * time.Millisecond)

	for received < 3 {
		select {
		case <-eventCh:
			received++
		case <-timeout:
			t.Fatalf("Timeout: only received %d/3 events", received)
		}
	}
}

func TestBrokerUnsubscribe(t *testing.T) {
	broker := NewBroker(10)
	defer broker.Clear()

	// Subscribe then immediately unsubscribe
	eventCh := broker.Subscribe(FileChanged)

	count := broker.SubscriberCount(FileChanged)
	if count != 1 {
		t.Errorf("Expected 1 subscriber, got %d", count)
	}

	broker.Unsubscribe(eventCh)

	count = broker.SubscriberCount(FileChanged)
	if count != 0 {
		t.Errorf("Expected 0 subscribers after unsubscribe, got %d", count)
	}

	// Publish event - should not panic even though channel is closed
	broker.PublishFile(FileChanged, FileData{Path: "/test/file.txt"})
}

func TestBrokerEventTypeString(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  string
	}{
		{SemanticAnalysisStarted, "SemanticAnalysisStarted"},
		{MemoryCreated, "MemoryCreated"},
		{AgentCompleted, "AgentCompleted"},
		{ModeChanged, "ModeChanged"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.eventType.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestBrokerNonBlockingPublish(t *testing.T) {
	broker := NewBroker(1) // Small buffer to test backpressure
	defer broker.Clear()

	eventCh := broker.Subscribe(ProposalGenerated)
	defer broker.Unsubscribe(eventCh)

	// Fill the buffer
	broker.PublishProposal(ProposalGenerated, ProposalData{ProposalID: "prop-1"})

	// This should not block even though buffer is full
	done := make(chan bool)
	go func() {
		broker.PublishProposal(ProposalGenerated, ProposalData{ProposalID: "prop-2"})
		done <- true
	}()

	select {
	case <-done:
		// Success - publish didn't block
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Publish blocked when it should have dropped the event")
	}

	// Drain the channel to verify first event was delivered
	select {
	case event := <-eventCh:
		data := event.Data.(ProposalData)
		if data.ProposalID != "prop-1" {
			t.Errorf("Expected prop-1, got %s", data.ProposalID)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for first event")
	}
}

func TestNewEvent(t *testing.T) {
	data := MemoryData{MemoryID: "mem-456"}
	event := NewEvent(MemoryUpdated, data)

	if event.Type != MemoryUpdated {
		t.Errorf("Expected event type %v, got %v", MemoryUpdated, event.Type)
	}

	if event.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}

	eventData, ok := event.Data.(MemoryData)
	if !ok {
		t.Fatalf("Expected MemoryData, got %T", event.Data)
	}

	if eventData.MemoryID != "mem-456" {
		t.Errorf("Expected MemoryID 'mem-456', got '%s'", eventData.MemoryID)
	}
}
