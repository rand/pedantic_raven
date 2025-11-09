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

// TestBrokerGlobalSubscriberCount verifies GlobalSubscriberCount returns correct count
func TestBrokerGlobalSubscriberCount(t *testing.T) {
	broker := NewBroker(10)
	defer broker.Clear()

	if count := broker.GlobalSubscriberCount(); count != 0 {
		t.Errorf("Expected 0 global subscribers initially, got %d", count)
	}

	ch1 := broker.SubscribeAll()
	if count := broker.GlobalSubscriberCount(); count != 1 {
		t.Errorf("Expected 1 global subscriber, got %d", count)
	}

	ch2 := broker.SubscribeAll()
	if count := broker.GlobalSubscriberCount(); count != 2 {
		t.Errorf("Expected 2 global subscribers, got %d", count)
	}

	broker.Unsubscribe(ch1)
	if count := broker.GlobalSubscriberCount(); count != 1 {
		t.Errorf("Expected 1 global subscriber after unsubscribe, got %d", count)
	}

	broker.Unsubscribe(ch2)
	if count := broker.GlobalSubscriberCount(); count != 0 {
		t.Errorf("Expected 0 global subscribers after unsubscribe, got %d", count)
	}
}

// TestBrokerPublishWithZeroSubscribers verifies publish handles zero subscribers gracefully
func TestBrokerPublishWithZeroSubscribers(t *testing.T) {
	broker := NewBroker(10)
	defer broker.Clear()

	// Publish with no subscribers - should not panic
	done := make(chan bool)
	go func() {
		broker.PublishMemory(MemoryCreated, MemoryData{MemoryID: "mem-1"})
		broker.PublishAgent(AgentStarted, AgentData{AgentID: "agent-1"})
		broker.PublishUI(ModeChanged, UIData{Mode: "edit"})
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Publish with zero subscribers blocked unexpectedly")
	}
}

// TestBrokerUnsubscribeNonExistent verifies unsubscribing non-existent channel is safe
func TestBrokerUnsubscribeNonExistent(t *testing.T) {
	broker := NewBroker(10)
	defer broker.Clear()

	// Create a channel but don't subscribe it
	unrelatedCh := make(chan Event, 10)
	defer close(unrelatedCh)

	// Should not panic or affect broker state
	broker.Unsubscribe(unrelatedCh)

	// Verify broker still works
	eventCh := broker.Subscribe(MemoryCreated)
	broker.PublishMemory(MemoryCreated, MemoryData{MemoryID: "mem-1"})

	select {
	case event := <-eventCh:
		if event.Type != MemoryCreated {
			t.Errorf("Expected MemoryCreated, got %v", event.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Failed to receive event after unsubscribing non-existent channel")
	}
	broker.Unsubscribe(eventCh)
}

// TestBrokerEventDeliveryOrder verifies events are delivered in publish order
func TestBrokerEventDeliveryOrder(t *testing.T) {
	broker := NewBroker(100)
	defer broker.Clear()

	eventCh := broker.Subscribe(MemoryUpdated)
	defer broker.Unsubscribe(eventCh)

	// Publish events in sequence
	ids := []string{"mem-1", "mem-2", "mem-3", "mem-4", "mem-5"}
	for _, id := range ids {
		broker.PublishMemory(MemoryUpdated, MemoryData{MemoryID: id})
	}

	// Receive events and verify order
	for _, expectedID := range ids {
		select {
		case event := <-eventCh:
			data := event.Data.(MemoryData)
			if data.MemoryID != expectedID {
				t.Errorf("Expected %s, got %s", expectedID, data.MemoryID)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timeout waiting for event with ID %s", expectedID)
		}
	}
}

// TestBrokerConcurrentPublishSubscribe tests concurrent publish/subscribe operations
func TestBrokerConcurrentPublishSubscribe(t *testing.T) {
	t.Parallel()

	broker := NewBroker(100)
	defer broker.Clear()

	done := make(chan int, 3)
	numPublishers := 3
	numEvents := 20

	// Multiple publishers
	for i := 0; i < numPublishers; i++ {
		go func(publisherID int) {
			for j := 0; j < numEvents; j++ {
				broker.PublishMemory(MemoryCreated, MemoryData{
					MemoryID: "mem-" + string(rune('A'+publisherID)) + "-" + string(rune('0'+j%10)),
				})
			}
			done <- 1
		}(i)
	}

	// Single subscriber collecting all events
	eventCh := broker.Subscribe(MemoryCreated)
	receivedCount := 0
	timeout := time.After(2 * time.Second)

	go func() {
		for receivedCount < numPublishers*numEvents {
			select {
			case <-eventCh:
				receivedCount++
			case <-timeout:
				return
			}
		}
		done <- 1
	}()

	// Wait for publishers
	for i := 0; i < numPublishers; i++ {
		<-done
	}

	// Wait for subscriber to collect all events
	<-done

	broker.Unsubscribe(eventCh)

	if receivedCount != numPublishers*numEvents {
		t.Errorf("Expected %d events, received %d", numPublishers*numEvents, receivedCount)
	}
}

// TestBrokerConcurrentSubscribeUnsubscribe tests concurrent subscribe/unsubscribe
func TestBrokerConcurrentSubscribeUnsubscribe(t *testing.T) {
	t.Parallel()

	broker := NewBroker(50)
	defer broker.Clear()

	done := make(chan bool, 10)

	// Multiple subscribers/unsubscribers
	for i := 0; i < 10; i++ {
		go func() {
			ch := broker.Subscribe(AgentStarted)
			time.Sleep(10 * time.Millisecond)
			broker.Unsubscribe(ch)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify broker still works
	finalCount := broker.SubscriberCount(AgentStarted)
	if finalCount != 0 {
		t.Errorf("Expected 0 subscribers after concurrent operations, got %d", finalCount)
	}
}

// TestBrokerMultipleEventTypes tests publishing multiple event types concurrently
func TestBrokerMultipleEventTypes(t *testing.T) {
	t.Parallel()

	broker := NewBroker(50)
	defer broker.Clear()

	memCh := broker.Subscribe(MemoryCreated)
	agentCh := broker.Subscribe(AgentStarted)
	uiCh := broker.Subscribe(ModeChanged)
	defer broker.Unsubscribe(memCh)
	defer broker.Unsubscribe(agentCh)
	defer broker.Unsubscribe(uiCh)

	// Publish different types
	broker.PublishMemory(MemoryCreated, MemoryData{MemoryID: "mem-1"})
	broker.PublishAgent(AgentStarted, AgentData{AgentID: "agent-1"})
	broker.PublishUI(ModeChanged, UIData{Mode: "edit"})

	timeout := time.After(100 * time.Millisecond)

	receivedMem := false
	receivedAgent := false
	receivedUI := false

	for !(receivedMem && receivedAgent && receivedUI) {
		select {
		case <-memCh:
			receivedMem = true
		case <-agentCh:
			receivedAgent = true
		case <-uiCh:
			receivedUI = true
		case <-timeout:
			t.Fatalf("Timeout: mem=%v, agent=%v, ui=%v", receivedMem, receivedAgent, receivedUI)
		}
	}
}

// TestBrokerSubscriberIsolation verifies one subscriber's buffer doesn't affect others
func TestBrokerSubscriberIsolation(t *testing.T) {
	broker := NewBroker(1) // Small buffer to test isolation
	defer broker.Clear()

	fastCh := broker.Subscribe(ProposalGenerated)
	slowCh := broker.Subscribe(ProposalGenerated)
	defer broker.Unsubscribe(fastCh)
	defer broker.Unsubscribe(slowCh)

	// Publish two events
	broker.PublishProposal(ProposalGenerated, ProposalData{ProposalID: "prop-1"})
	broker.PublishProposal(ProposalGenerated, ProposalData{ProposalID: "prop-2"})

	// Fast subscriber reads first event immediately
	select {
	case event := <-fastCh:
		if event.Data.(ProposalData).ProposalID != "prop-1" {
			t.Error("Fast subscriber should receive prop-1")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Fast subscriber timeout")
	}

	// Fast subscriber reads second event (buffer full, but event may have been dropped)
	// Slow subscriber should still work independently
	select {
	case event := <-slowCh:
		data := event.Data.(ProposalData)
		if data.ProposalID != "prop-1" && data.ProposalID != "prop-2" {
			t.Errorf("Slow subscriber received unexpected event: %s", data.ProposalID)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Slow subscriber timeout")
	}
}

// TestBrokerClearCloses all channels
func TestBrokerClearClosesChannels(t *testing.T) {
	broker := NewBroker(10)

	typeCh1 := broker.Subscribe(MemoryCreated)
	typeCh2 := broker.Subscribe(AgentStarted)
	globalCh := broker.SubscribeAll()

	broker.Clear()

	// All channels should be closed
	select {
	case _, ok := <-typeCh1:
		if ok {
			t.Error("typeCh1 should be closed after Clear()")
		}
	default:
		// Expected - channel is closed
	}

	select {
	case _, ok := <-typeCh2:
		if ok {
			t.Error("typeCh2 should be closed after Clear()")
		}
	default:
		// Expected - channel is closed
	}

	select {
	case _, ok := <-globalCh:
		if ok {
			t.Error("globalCh should be closed after Clear()")
		}
	default:
		// Expected - channel is closed
	}

	// Verify broker state is reset
	if broker.SubscriberCount(MemoryCreated) != 0 {
		t.Error("Expected 0 subscribers after Clear()")
	}
	if broker.GlobalSubscriberCount() != 0 {
		t.Error("Expected 0 global subscribers after Clear()")
	}
}

// TestBrokerClearThenReuse verifies broker can be reused after Clear
func TestBrokerClearThenReuse(t *testing.T) {
	broker := NewBroker(10)

	ch1 := broker.Subscribe(FileChanged)
	broker.Clear()
	broker.Unsubscribe(ch1) // Should not panic

	// Should work normally after Clear
	ch2 := broker.Subscribe(FileChanged)
	broker.PublishFile(FileChanged, FileData{Path: "/test"})

	select {
	case event := <-ch2:
		if event.Type != FileChanged {
			t.Errorf("Expected FileChanged, got %v", event.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout receiving event after Clear")
	}
	broker.Unsubscribe(ch2)
}

// TestPublishConvenienceMethods tests all convenience publish methods
func TestPublishConvenienceMethods(t *testing.T) {
	t.Parallel()

	broker := NewBroker(10)
	defer broker.Clear()

	tests := []struct {
		name     string
		eventType EventType
		publish  func()
	}{
		{
			name:      "PublishUser",
			eventType: UserJoined,
			publish: func() {
				broker.PublishUser(UserJoined, UserData{UserID: "user-1"})
			},
		},
		{
			name:      "PublishDiagnostic",
			eventType: DiagnosticsUpdated,
			publish: func() {
				broker.PublishDiagnostic(DiagnosticsUpdated, DiagnosticData{BufferID: "buf-1"})
			},
		},
		{
			name:      "PublishServer",
			eventType: ServerConnected,
			publish: func() {
				broker.PublishServer(ServerConnected, ServerData{ServerURL: "http://localhost"})
			},
		},
		{
			name:      "PublishError",
			eventType: ErrorOccurred,
			publish: func() {
				broker.PublishError(ErrorData{Context: "test"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := broker.Subscribe(tt.eventType)
			defer broker.Unsubscribe(ch)

			tt.publish()

			select {
			case event := <-ch:
				if event.Type != tt.eventType {
					t.Errorf("Expected %v, got %v", tt.eventType, event.Type)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout waiting for event")
			}
		})
	}
}

// TestEventTypeStringCoverage tests all EventType string representations
func TestEventTypeStringCoverage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		eventType EventType
		expected  string
	}{
		// Semantic analysis events
		{SemanticAnalysisStarted, "SemanticAnalysisStarted"},
		{SemanticAnalysisProgress, "SemanticAnalysisProgress"},
		{SemanticAnalysisComplete, "SemanticAnalysisComplete"},
		{SemanticAnalysisFailed, "SemanticAnalysisFailed"},
		// Memory events
		{MemoryRecalled, "MemoryRecalled"},
		{MemoryCreated, "MemoryCreated"},
		{MemoryUpdated, "MemoryUpdated"},
		{MemoryDeleted, "MemoryDeleted"},
		{MemorySearchStarted, "MemorySearchStarted"},
		{MemorySearchComplete, "MemorySearchComplete"},
		// Agent events
		{AgentStarted, "AgentStarted"},
		{AgentProgress, "AgentProgress"},
		{AgentCompleted, "AgentCompleted"},
		{AgentFailed, "AgentFailed"},
		{AgentPaused, "AgentPaused"},
		{AgentResumed, "AgentResumed"},
		// User collaboration events
		{UserJoined, "UserJoined"},
		{UserLeft, "UserLeft"},
		{CursorMoved, "CursorMoved"},
		{EditMade, "EditMade"},
		{ChatMessageReceived, "ChatMessageReceived"},
		// Proposal events
		{ProposalGenerated, "ProposalGenerated"},
		{ProposalAccepted, "ProposalAccepted"},
		{ProposalRejected, "ProposalRejected"},
		{ProposalModified, "ProposalModified"},
		// Diagnostic events
		{DiagnosticsUpdated, "DiagnosticsUpdated"},
		{ValidationStarted, "ValidationStarted"},
		{ValidationComplete, "ValidationComplete"},
		// File system events
		{FileChanged, "FileChanged"},
		{FileOpened, "FileOpened"},
		{FileClosed, "FileClosed"},
		{FileHistoryUpdated, "FileHistoryUpdated"},
		// UI events
		{ModeChanged, "ModeChanged"},
		{LayoutChanged, "LayoutChanged"},
		{PaneFocused, "PaneFocused"},
		{OverlayOpened, "OverlayOpened"},
		{OverlayClosed, "OverlayClosed"},
		// System events
		{ServerConnected, "ServerConnected"},
		{ServerDisconnected, "ServerDisconnected"},
		{ServerReconnecting, "ServerReconnecting"},
		{ErrorOccurred, "ErrorOccurred"},
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

// TestUnknownEventTypeString tests unknown event type returns "Unknown"
func TestUnknownEventTypeString(t *testing.T) {
	unknownType := EventType(999)
	result := unknownType.String()
	if result != "Unknown" {
		t.Errorf("Expected 'Unknown', got %s", result)
	}
}
