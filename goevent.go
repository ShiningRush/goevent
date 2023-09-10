package goevent

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
)

var defaultBus EventBus = NewInMemoryEventBus()

// SetEventBus is used for unit test
func SetEventBus(bus EventBus) {
	defaultBus = bus
}

// EventBus
type EventBus interface {
	Subscribe(handler EventHandler) error
	Publish(event Event)
	PublishSync(cxt context.Context, event Event)
	Close()
}

// EventHandler
type EventHandler interface {
	Topic() []string
	Handle(ctx context.Context, event Event)
}

// Event
type Event interface {
	Topic() []string
}

// NewInMemoryEventBus create a event-bus which work in memory
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{handlers: make(map[string][]EventHandler)}
}

// InMemoryEventBus
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.Mutex

	closed atomic.Bool
	wg     sync.WaitGroup
}

// Subscribe is thread-safe, but dont invoke Publish at same time.
func (b *InMemoryEventBus) Subscribe(handler EventHandler) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, tp := range handler.Topic() {
		b.handlers[tp] =
			append(b.handlers[tp], handler)
	}
	return nil
}

// Publish a event and DO NOT wait the handler executing
func (b *InMemoryEventBus) Publish(event Event) {
	if b.closed.Load() {
		panic("event bus is already closed")
	}

	matchedHandlers := b.getMatchHandlers(event.Topic())
	if matchedHandlers != nil {
		for i := range matchedHandlers {
			b.wg.Add(1)
			go func(idx int) {
				matchedHandlers[idx].Handle(context.TODO(), event)
				b.wg.Done()
			}(i)
		}
	}
}

// PublishSync a event and WILL wait the handler executing
func (b *InMemoryEventBus) PublishSync(ctx context.Context, event Event) {
	if b.closed.Load() {
		panic("event bus is already closed")
	}

	matchedHandlers := b.getMatchHandlers(event.Topic())
	if len(matchedHandlers) > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(matchedHandlers))
		for i := range matchedHandlers {
			go func(idx int) {
				matchedHandlers[idx].Handle(ctx, event)
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
}

// Close will wait all async goroutine completed to prevent lost changes
func (b *InMemoryEventBus) Close() {
	if !b.closed.Load() {
		b.closed.Store(true)
	}
	b.wg.Wait()
}

func (b *InMemoryEventBus) getMatchHandlers(topics []string) (matchedHandlers []EventHandler) {
	for _, tp := range topics {
		if h, ok := b.handlers[tp]; ok {
			matchedHandlers = append(matchedHandlers, h...)
		}
	}
	return
}

func (b *InMemoryEventBus) getEventTopic(eventType reflect.Type) reflect.Type {
	return eventType
}

// Subscribe is thread-safe, but dont call Publish at same time.
func Subscribe(handler EventHandler) error {
	return defaultBus.Subscribe(handler)
}

// Publish a event and DO NOT wait the handler executing
func Publish(event Event) {
	defaultBus.Publish(event)
}

// PublishSync a event and WILL wait the handler executing
func PublishSync(ctx context.Context, event Event) {
	defaultBus.PublishSync(ctx, event)
}

// Close will wait all async goroutine completed to prevent lost changes
func Close() {
	defaultBus.Close()
}
