package events

import (
	"context"
	"meuApp/pkg/contracts"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helpers

type testEventHandler struct {
	called     bool
	eventType  string
	payload    interface{}
	mu         sync.Mutex
	callCount  int
	shouldFail bool
}

func newTestEventHandler() *testEventHandler {
	return &testEventHandler{}
}

func (h *testEventHandler) Handle(ctx context.Context, event contracts.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.called = true
	h.callCount++
	h.eventType = event.Type
	h.payload = event.Payload

	if h.shouldFail {
		return assert.AnError
	}
	return nil
}

func (h *testEventHandler) WasCalled() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.called
}

func (h *testEventHandler) GetCallCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.callCount
}

// EventBus Tests

func TestNewEventBus(t *testing.T) {
	// Act
	bus := NewEventBus()

	// Assert
	assert.NotNil(t, bus)
	assert.NotNil(t, bus.handlers)
	assert.Equal(t, 0, len(bus.handlers))
}

func TestEventBus_Subscribe(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	handler := func(ctx context.Context, event contracts.Event) error {
		return nil
	}

	// Act
	err := bus.Subscribe("test.event", handler)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, len(bus.handlers))
	assert.Equal(t, 1, len(bus.handlers["test.event"]))
}

func TestEventBus_Subscribe_MultipleHandlers(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	handler1 := func(ctx context.Context, event contracts.Event) error { return nil }
	handler2 := func(ctx context.Context, event contracts.Event) error { return nil }
	handler3 := func(ctx context.Context, event contracts.Event) error { return nil }

	// Act
	err1 := bus.Subscribe("test.event", handler1)
	err2 := bus.Subscribe("test.event", handler2)
	err3 := bus.Subscribe("test.event", handler3)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	assert.Equal(t, 1, len(bus.handlers))               // 1 event type
	assert.Equal(t, 3, len(bus.handlers["test.event"])) // 3 handlers
}

func TestEventBus_Subscribe_DifferentEventTypes(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	handler1 := func(ctx context.Context, event contracts.Event) error { return nil }
	handler2 := func(ctx context.Context, event contracts.Event) error { return nil }

	// Act
	bus.Subscribe("event.type1", handler1)
	bus.Subscribe("event.type2", handler2)

	// Assert
	assert.Equal(t, 2, len(bus.handlers))
	assert.Equal(t, 1, len(bus.handlers["event.type1"]))
	assert.Equal(t, 1, len(bus.handlers["event.type2"]))
}

func TestEventBus_Publish_WithNoHandlers(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	ctx := context.Background()
	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"key": "value"},
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	require.NoError(t, err) // Should not error when no handlers exist
}

func TestEventBus_Publish_WithHandler(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	ctx := context.Background()
	handler := newTestEventHandler()

	bus.Subscribe("test.event", handler.Handle)

	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"key": "value"},
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	require.NoError(t, err)
	assert.True(t, handler.WasCalled())
	assert.Equal(t, "test.event", handler.eventType)
	assert.Equal(t, 1, handler.GetCallCount())
}

func TestEventBus_Publish_WithMultipleHandlers(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	ctx := context.Background()
	handler1 := newTestEventHandler()
	handler2 := newTestEventHandler()
	handler3 := newTestEventHandler()

	bus.Subscribe("test.event", handler1.Handle)
	bus.Subscribe("test.event", handler2.Handle)
	bus.Subscribe("test.event", handler3.Handle)

	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"data": "test"},
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	require.NoError(t, err)
	assert.True(t, handler1.WasCalled())
	assert.True(t, handler2.WasCalled())
	assert.True(t, handler3.WasCalled())
}

func TestEventBus_Publish_HandlerError_ContinuesExecution(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	ctx := context.Background()

	handler1 := newTestEventHandler()
	handler1.shouldFail = true // This handler will fail
	handler2 := newTestEventHandler()

	bus.Subscribe("test.event", handler1.Handle)
	bus.Subscribe("test.event", handler2.Handle)

	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{},
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	require.NoError(t, err) // Publish doesn't return error even if handler fails
	assert.True(t, handler1.WasCalled())
	assert.True(t, handler2.WasCalled()) // Handler2 should still be called
}

func TestEventBus_Publish_OnlyMatchingHandlersCalled(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	ctx := context.Background()

	handlerA := newTestEventHandler()
	handlerB := newTestEventHandler()

	bus.Subscribe("event.a", handlerA.Handle)
	bus.Subscribe("event.b", handlerB.Handle)

	event := contracts.Event{
		Type:      "event.a",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{},
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	require.NoError(t, err)
	assert.True(t, handlerA.WasCalled())
	assert.False(t, handlerB.WasCalled()) // Handler B should NOT be called
}

func TestEventBus_ConcurrentPublish(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	handler := newTestEventHandler()
	bus.Subscribe("test.event", handler.Handle)

	ctx := context.Background()
	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{},
	}

	// Act - Publish concurrently
	var wg sync.WaitGroup
	iterations := 10

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bus.Publish(ctx, event)
		}()
	}

	wg.Wait()

	// Assert
	assert.Equal(t, iterations, handler.GetCallCount())
}

func TestEventBus_ConcurrentSubscribe(t *testing.T) {
	// Arrange
	bus := NewEventBus()

	// Act - Subscribe concurrently
	var wg sync.WaitGroup
	iterations := 10

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			handler := func(ctx context.Context, event contracts.Event) error { return nil }
			bus.Subscribe("test.event", handler)
		}()
	}

	wg.Wait()

	// Assert
	assert.Equal(t, iterations, len(bus.handlers["test.event"]))
}

// Event Type Constants Tests

func TestEventTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"UserCreated", UserCreatedEventType, "user.created"},
		{"UserUpdated", UserUpdatedEventType, "user.updated"},
		{"UserDeleted", UserDeletedEventType, "user.deleted"},
		{"ProductCreated", ProductCreatedEventType, "product.created"},
		{"ProductUpdated", ProductUpdatedEventType, "product.updated"},
		{"ProductStockUpdated", ProductStockUpdatedEventType, "product.stock.updated"},
		{"OrderCreated", OrderCreatedEventType, "order.created"},
		{"OrderStatusUpdated", OrderStatusUpdatedEventType, "order.status.updated"},
		{"OrderCancelled", OrderCancelledEventType, "order.cancelled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

// Benchmark Tests

func BenchmarkEventBus_Subscribe(b *testing.B) {
	bus := NewEventBus()
	handler := func(ctx context.Context, event contracts.Event) error { return nil }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Subscribe("test.event", handler)
	}
}

func BenchmarkEventBus_Publish_NoHandlers(b *testing.B) {
	bus := NewEventBus()
	ctx := context.Background()
	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(ctx, event)
	}
}

func BenchmarkEventBus_Publish_OneHandler(b *testing.B) {
	bus := NewEventBus()
	handler := func(ctx context.Context, event contracts.Event) error { return nil }
	bus.Subscribe("test.event", handler)

	ctx := context.Background()
	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(ctx, event)
	}
}

func BenchmarkEventBus_Publish_MultipleHandlers(b *testing.B) {
	bus := NewEventBus()
	handler := func(ctx context.Context, event contracts.Event) error { return nil }

	// Subscribe 10 handlers
	for i := 0; i < 10; i++ {
		bus.Subscribe("test.event", handler)
	}

	ctx := context.Background()
	event := contracts.Event{
		Type:      "test.event",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(ctx, event)
	}
}
