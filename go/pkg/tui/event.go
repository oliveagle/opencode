package tui

import (
	"sync"
)

// EventType defines the type of event
type EventType int

const (
	EventFileCreated EventType = iota
	EventFileModified
	EventFileDeleted
	EventFileRenamed
	EventProjectLoaded
	EventEditorSaved
)

// Event represents an application event
type Event struct {
	Type    EventType
	Path    string
	Content string
	Data    interface{}
}

// EventBus manages event subscriptions and publishing
type EventBus struct {
	subscribers map[EventType][]chan Event
	mu          sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]chan Event),
	}
}

// Subscribe returns a channel for the given event type
func (eb *EventBus) Subscribe(eventType EventType) <-chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event, 10)
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
	return ch
}

// Publish sends an event to all subscribers
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	subscribers := eb.subscribers[event.Type]
	eb.mu.RUnlock()

	for _, ch := range subscribers {
		select {
		case ch <- event:
		default:
			// Channel full, skip
		}
	}
}

// Unsubscribe removes a channel from subscribers
func (eb *EventBus) Unsubscribe(eventType EventType, ch <-chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscribers := eb.subscribers[eventType]
	for i, sub := range subscribers {
		if sub == ch {
			// Remove by replacing with last and truncating
			subscribers[i] = subscribers[len(subscribers)-1]
			eb.subscribers[eventType] = subscribers[:len(subscribers)-1]
			break
		}
	}
}
