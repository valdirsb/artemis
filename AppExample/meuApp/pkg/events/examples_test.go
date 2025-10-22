package events_test

import (
	"context"
	"fmt"
	"testing"

	"meuApp/pkg/events"
)

// Exemplo de uso do sistema de eventos tipados
func ExampleTypedEventPublisher() {
	// Criar EventBus
	bus := events.NewEventBus()

	// Criar publisher tipado
	publisher := events.NewTypedEventPublisher(bus)

	// Registrar handler para UserCreatedEvent (logger nil usa default)
	err := events.SubscribeTyped(bus, events.EventTypeUserCreated,
		events.UserCreatedHandlerFunc(nil))
	if err != nil {
		panic(err)
	}

	// Publicar evento tipado
	ctx := context.Background()
	err = publisher.PublishUserCreated(ctx, events.UserCreatedEvent{
		UserID:   "user-123",
		Username: "john_doe",
		Email:    "john@example.com",
	})
	if err != nil {
		panic(err)
	}

	// Output: [INFO] Processing UserCreatedEvent for user: john_doe (john@example.com)
}

// Exemplo de handler customizado
func ExampleTypedEventPublisher_customHandler() {
	bus := events.NewEventBus()
	publisher := events.NewTypedEventPublisher(bus)

	// Handler customizado
	customHandler := func(ctx context.Context, data events.OrderCreatedEvent) error {
		// Lógica personalizada
		fmt.Printf("Pedido criado: %s\n", data.OrderID)
		return nil
	}

	// Registrar handler
	events.SubscribeTyped(bus, events.EventTypeOrderCreated, customHandler)

	// Publicar
	publisher.PublishOrderCreated(context.Background(), events.OrderCreatedEvent{
		OrderID:     "order-456",
		UserID:      "user-123",
		TotalAmount: 199.99,
		ItemCount:   3,
	})

	// Output: Pedido criado: order-456
}

// Teste de múltiplos handlers
func TestMultipleHandlers(t *testing.T) {
	bus := events.NewEventBus()
	publisher := events.NewTypedEventPublisher(bus)

	// Registrar múltiplos handlers para o mesmo evento
	handler1 := func(ctx context.Context, data events.LowStockEvent) error {
		t.Logf("Handler 1: Low stock for %s", data.ProductName)
		return nil
	}

	handler2 := func(ctx context.Context, data events.LowStockEvent) error {
		t.Logf("Handler 2: Notifying admin about %s", data.ProductName)
		return nil
	}

	events.SubscribeTyped(bus, events.EventTypeLowStock, handler1)
	events.SubscribeTyped(bus, events.EventTypeLowStock, handler2)

	// Ambos os handlers serão executados
	err := publisher.PublishLowStock(context.Background(), events.LowStockEvent{
		ProductID:    "prod-789",
		ProductName:  "Widget",
		CurrentStock: 5,
		Threshold:    10,
	})

	if err != nil {
		t.Errorf("Error publishing event: %v", err)
	}
}
