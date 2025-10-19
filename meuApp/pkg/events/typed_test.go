package events

import (
	"context"
	"meuApp/pkg/contracts"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TypedEventPublisher Tests

func TestNewTypedEventPublisher(t *testing.T) {
	// Arrange
	bus := NewEventBus()

	// Act
	publisher := NewTypedEventPublisher(bus)

	// Assert
	assert.NotNil(t, publisher)
	assert.Equal(t, bus, publisher.bus)
}

func TestTypedEventPublisher_PublishUserCreated(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeUserCreated, handler)

	data := UserCreatedEvent{
		UserID: "user-123",
		Email:  "test@example.com",
	}

	// Act
	err := publisher.PublishUserCreated(ctx, data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, EventTypeUserCreated, receivedEvent.Type)
	assert.NotZero(t, receivedEvent.Timestamp)

	// Verify payload
	payload, ok := receivedEvent.Payload.(contracts.UserCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, "user-123", payload.UserID)
	assert.Equal(t, "test@example.com", payload.Email)
}

func TestTypedEventPublisher_PublishUserDeleted(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeUserDeleted, handler)

	data := UserDeletedEvent{
		UserID: "user-456",
		Email:  "deleted@example.com",
	}

	// Act
	err := publisher.PublishUserDeleted(ctx, data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, EventTypeUserDeleted, receivedEvent.Type)

	// Verify payload
	payload, ok := receivedEvent.Payload.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "user-456", payload["user_id"])
	assert.Equal(t, "deleted@example.com", payload["email"])
}

func TestTypedEventPublisher_PublishProductCreated(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeProductCreated, handler)

	data := ProductCreatedEvent{
		ProductID:  "prod-789",
		Name:       "Test Product",
		CategoryID: "cat-1",
		Price:      99.99,
		Stock:      50,
	}

	// Act
	err := publisher.PublishProductCreated(ctx, data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, EventTypeProductCreated, receivedEvent.Type)

	// Verify payload
	payload, ok := receivedEvent.Payload.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "prod-789", payload["product_id"])
	assert.Equal(t, "Test Product", payload["name"])
	assert.Equal(t, "cat-1", payload["category_id"])
	assert.Equal(t, 99.99, payload["price"])
	assert.Equal(t, 50, payload["stock"])
}

func TestTypedEventPublisher_PublishLowStock(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeLowStock, handler)

	data := LowStockEvent{
		ProductID:    "prod-999",
		ProductName:  "Low Stock Item",
		CurrentStock: 5,
		Threshold:    10,
	}

	// Act
	err := publisher.PublishLowStock(ctx, data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, EventTypeLowStock, receivedEvent.Type)

	// Verify payload
	payload, ok := receivedEvent.Payload.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "prod-999", payload["product_id"])
	assert.Equal(t, "Low Stock Item", payload["product_name"])
	assert.Equal(t, 5, payload["current_stock"])
	assert.Equal(t, 10, payload["threshold"])
}

func TestTypedEventPublisher_PublishOrderCreated(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeOrderCreated, handler)

	data := OrderCreatedEvent{
		OrderID:     "order-123",
		UserID:      "user-456",
		TotalAmount: 299.99,
		ItemCount:   3,
	}

	// Act
	err := publisher.PublishOrderCreated(ctx, data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, EventTypeOrderCreated, receivedEvent.Type)

	// Verify payload
	payload, ok := receivedEvent.Payload.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "order-123", payload["order_id"])
	assert.Equal(t, "user-456", payload["user_id"])
	assert.Equal(t, 299.99, payload["total_amount"])
	assert.Equal(t, 3, payload["item_count"])
}

func TestTypedEventPublisher_PublishOrderStatusChanged(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeOrderStatusChanged, handler)

	data := OrderStatusChangedEvent{
		OrderID:        "order-789",
		PreviousStatus: "pending",
		NewStatus:      "shipped",
		ChangedBy:      "user-001",
	}

	// Act
	err := publisher.PublishOrderStatusChanged(ctx, data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, EventTypeOrderStatusChanged, receivedEvent.Type)

	// Verify payload
	payload, ok := receivedEvent.Payload.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "order-789", payload["order_id"])
	assert.Equal(t, "pending", payload["previous_status"])
	assert.Equal(t, "shipped", payload["new_status"])
}

func TestTypedEventPublisher_PublishOrderCancelled(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeOrderCancelled, handler)

	data := OrderCancelledEvent{
		OrderID:      "order-999",
		UserID:       "user-555",
		Reason:       "Customer request",
		RefundAmount: 299.99,
	}

	// Act
	err := publisher.PublishOrderCancelled(ctx, data)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, EventTypeOrderCancelled, receivedEvent.Type)

	// Verify payload
	payload, ok := receivedEvent.Payload.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "order-999", payload["order_id"])
	assert.Equal(t, "user-555", payload["user_id"])
	assert.Equal(t, "Customer request", payload["reason"])
	assert.Equal(t, 299.99, payload["refund_amount"])
}

func TestTypedEventPublisher_MultiplePublishes(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	callCount := 0
	handler := func(ctx context.Context, event contracts.Event) error {
		callCount++
		return nil
	}

	bus.Subscribe(EventTypeUserCreated, handler)
	bus.Subscribe(EventTypeProductCreated, handler)
	bus.Subscribe(EventTypeOrderCreated, handler)

	// Act
	publisher.PublishUserCreated(ctx, UserCreatedEvent{UserID: "u1", Email: "u1@test.com"})
	publisher.PublishProductCreated(ctx, ProductCreatedEvent{ProductID: "p1", Name: "Product 1", CategoryID: "c1", Price: 10, Stock: 100})
	publisher.PublishOrderCreated(ctx, OrderCreatedEvent{OrderID: "o1", UserID: "u1", TotalAmount: 100, ItemCount: 1})

	// Assert
	assert.Equal(t, 3, callCount)
}

func TestTypedEventPublisher_SubscribeTyped(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	ctx := context.Background()

	receivedCount := 0
	handlerFunc := func(ctx context.Context, data map[string]interface{}) error {
		receivedCount++
		return nil
	}

	SubscribeTyped(bus, EventTypeProductCreated, handlerFunc)

	event := contracts.Event{
		Type:      EventTypeProductCreated,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"product_id": "test-123",
			"name":       "Test Product",
		},
	}

	// Act
	err := bus.Publish(ctx, event)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, receivedCount)
}

func TestTypedEventPublisher_TimestampIsSet(t *testing.T) {
	// Arrange
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	var receivedEvent contracts.Event
	handler := func(ctx context.Context, event contracts.Event) error {
		receivedEvent = event
		return nil
	}
	bus.Subscribe(EventTypeUserCreated, handler)

	before := time.Now()
	data := UserCreatedEvent{UserID: "u1", Email: "test@test.com"}

	// Act
	publisher.PublishUserCreated(ctx, data)
	after := time.Now()

	// Assert
	assert.True(t, receivedEvent.Timestamp.After(before) || receivedEvent.Timestamp.Equal(before))
	assert.True(t, receivedEvent.Timestamp.Before(after) || receivedEvent.Timestamp.Equal(after))
}

// Benchmark Tests

func BenchmarkTypedEventPublisher_PublishUserCreated(b *testing.B) {
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	handler := func(ctx context.Context, event contracts.Event) error { return nil }
	bus.Subscribe(EventTypeUserCreated, handler)

	data := UserCreatedEvent{UserID: "u1", Email: "test@test.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher.PublishUserCreated(ctx, data)
	}
}

func BenchmarkTypedEventPublisher_PublishOrderCreated(b *testing.B) {
	bus := NewEventBus()
	publisher := NewTypedEventPublisher(bus)
	ctx := context.Background()

	handler := func(ctx context.Context, event contracts.Event) error { return nil }
	bus.Subscribe(EventTypeOrderCreated, handler)

	data := OrderCreatedEvent{
		OrderID:     "o1",
		UserID:      "u1",
		TotalAmount: 100.0,
		ItemCount:   5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		publisher.PublishOrderCreated(ctx, data)
	}
}
