package events

import (
	"sync"
)

// Broker coordinates pub/sub event distribution across the application.
//
// Components publish events to the broker, which distributes them to all subscribers.
// This enables reactive updates and decoupled communication between components.
//
// The broker is thread-safe and can be used from multiple goroutines.
type Broker struct {
	mu          sync.RWMutex
	subscribers map[EventType][]chan Event
	globalSubs  []chan Event // Subscribers for all event types
	bufferSize  int
}

// NewBroker creates a new event broker with the specified buffer size per subscriber.
//
// The buffer size determines how many events can be queued before blocking.
// A larger buffer reduces backpressure but increases memory usage.
// Typical values: 10-100 for UI events, 1000+ for high-frequency events.
func NewBroker(bufferSize int) *Broker {
	return &Broker{
		subscribers: make(map[EventType][]chan Event),
		globalSubs:  make([]chan Event, 0),
		bufferSize:  bufferSize,
	}
}

// Publish sends an event to all subscribers.
//
// This is a non-blocking operation. If a subscriber's channel is full,
// the event will be dropped for that subscriber to prevent blocking.
//
// The publisher does not wait for subscribers to process the event.
func (b *Broker) Publish(event Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Send to type-specific subscribers
	if subs, ok := b.subscribers[event.Type]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
				// Event delivered
			default:
				// Channel full, drop event to prevent blocking
				// TODO: Add metrics/logging for dropped events
			}
		}
	}

	// Send to global subscribers (listening to all events)
	for _, ch := range b.globalSubs {
		select {
		case ch <- event:
			// Event delivered
		default:
			// Channel full, drop event
		}
	}
}

// Subscribe returns a channel that receives events of the specified type.
//
// The caller is responsible for reading from the channel to prevent blocking.
// When done, call Unsubscribe to clean up resources.
//
// Example:
//
//	eventCh := broker.Subscribe(SemanticAnalysisComplete)
//	defer broker.Unsubscribe(eventCh)
//
//	for event := range eventCh {
//	    // Handle event
//	}
func (b *Broker) Subscribe(eventType EventType) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, b.bufferSize)
	b.subscribers[eventType] = append(b.subscribers[eventType], ch)
	return ch
}

// SubscribeAll returns a channel that receives all events regardless of type.
//
// This is useful for logging, debugging, or implementing centralized event handlers.
// Use with caution as it creates high volume of events.
func (b *Broker) SubscribeAll() <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, b.bufferSize)
	b.globalSubs = append(b.globalSubs, ch)
	return ch
}

// Unsubscribe removes a subscriber channel and closes it.
//
// After unsubscribing, the channel will be closed and no more events will be sent to it.
// It's safe to call this multiple times with the same channel.
func (b *Broker) Unsubscribe(ch <-chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Remove from type-specific subscribers
	for eventType, subs := range b.subscribers {
		for i, sub := range subs {
			if sub == ch {
				// Remove this subscriber (preserve order)
				b.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
				close(sub)
				return
			}
		}
	}

	// Remove from global subscribers
	for i, sub := range b.globalSubs {
		if sub == ch {
			b.globalSubs = append(b.globalSubs[:i], b.globalSubs[i+1:]...)
			close(sub)
			return
		}
	}
}

// SubscriberCount returns the number of active subscribers for a given event type.
//
// This is useful for debugging and monitoring event flow.
func (b *Broker) SubscriberCount(eventType EventType) int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.subscribers[eventType])
}

// GlobalSubscriberCount returns the number of global subscribers (listening to all events).
func (b *Broker) GlobalSubscriberCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.globalSubs)
}

// Clear removes all subscribers and closes all channels.
//
// This is typically called during shutdown to clean up resources.
func (b *Broker) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Close all type-specific subscriber channels
	for _, subs := range b.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
	b.subscribers = make(map[EventType][]chan Event)

	// Close all global subscriber channels
	for _, ch := range b.globalSubs {
		close(ch)
	}
	b.globalSubs = make([]chan Event, 0)
}

// --- Convenience Methods ---

// PublishSemanticAnalysis publishes a semantic analysis event.
func (b *Broker) PublishSemanticAnalysis(eventType EventType, data SemanticAnalysisData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishMemory publishes a memory event.
func (b *Broker) PublishMemory(eventType EventType, data MemoryData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishAgent publishes an agent orchestration event.
func (b *Broker) PublishAgent(eventType EventType, data AgentData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishUser publishes a collaboration event.
func (b *Broker) PublishUser(eventType EventType, data UserData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishProposal publishes a proposal event.
func (b *Broker) PublishProposal(eventType EventType, data ProposalData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishDiagnostic publishes a diagnostic event.
func (b *Broker) PublishDiagnostic(eventType EventType, data DiagnosticData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishFile publishes a file system event.
func (b *Broker) PublishFile(eventType EventType, data FileData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishUI publishes a UI event.
func (b *Broker) PublishUI(eventType EventType, data UIData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishServer publishes a server connection event.
func (b *Broker) PublishServer(eventType EventType, data ServerData) {
	b.Publish(NewEvent(eventType, data))
}

// PublishError publishes an error event.
func (b *Broker) PublishError(data ErrorData) {
	b.Publish(NewEvent(ErrorOccurred, data))
}
