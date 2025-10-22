package events

import (
	"context"
	"encoding/json"
	"meuApp/pkg/contracts"
	"time"
)

// TypedEventPublisher é um wrapper type-safe para publicação de eventos
type TypedEventPublisher struct {
	bus *EventBus
}

// NewTypedEventPublisher cria um novo publisher tipado
func NewTypedEventPublisher(bus *EventBus) *TypedEventPublisher {
	return &TypedEventPublisher{bus: bus}
}

// PublishUserCreated publica um evento de usuário criado
func (p *TypedEventPublisher) PublishUserCreated(ctx context.Context, data UserCreatedEvent) error {
	event := contracts.Event{
		Type:      EventTypeUserCreated,
		Timestamp: time.Now(),
		Payload: contracts.UserCreatedEvent{
			UserID: data.UserID,
			Email:  data.Email,
		},
	}
	return p.bus.Publish(ctx, event)
}

// PublishUserDeleted publica um evento de usuário deletado
func (p *TypedEventPublisher) PublishUserDeleted(ctx context.Context, data UserDeletedEvent) error {
	event := contracts.Event{
		Type:      EventTypeUserDeleted,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"user_id": data.UserID,
			"email":   data.Email,
		},
	}
	return p.bus.Publish(ctx, event)
}

// PublishProductCreated publica um evento de produto criado
func (p *TypedEventPublisher) PublishProductCreated(ctx context.Context, data ProductCreatedEvent) error {
	event := contracts.Event{
		Type:      EventTypeProductCreated,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"product_id":  data.ProductID,
			"name":        data.Name,
			"category_id": data.CategoryID,
			"price":       data.Price,
			"stock":       data.Stock,
		},
	}
	return p.bus.Publish(ctx, event)
}

// PublishLowStock publica um alerta de estoque baixo
func (p *TypedEventPublisher) PublishLowStock(ctx context.Context, data LowStockEvent) error {
	event := contracts.Event{
		Type:      EventTypeLowStock,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"product_id":    data.ProductID,
			"product_name":  data.ProductName,
			"current_stock": data.CurrentStock,
			"threshold":     data.Threshold,
		},
	}
	return p.bus.Publish(ctx, event)
}

// PublishOrderCreated publica um evento de pedido criado
func (p *TypedEventPublisher) PublishOrderCreated(ctx context.Context, data OrderCreatedEvent) error {
	event := contracts.Event{
		Type:      EventTypeOrderCreated,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id":     data.OrderID,
			"user_id":      data.UserID,
			"total_amount": data.TotalAmount,
			"item_count":   data.ItemCount,
		},
	}
	return p.bus.Publish(ctx, event)
}

// PublishOrderStatusChanged publica um evento de mudança de status
func (p *TypedEventPublisher) PublishOrderStatusChanged(ctx context.Context, data OrderStatusChangedEvent) error {
	event := contracts.Event{
		Type:      EventTypeOrderStatusChanged,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id":        data.OrderID,
			"previous_status": data.PreviousStatus,
			"new_status":      data.NewStatus,
			"changed_by":      data.ChangedBy,
		},
	}
	return p.bus.Publish(ctx, event)
}

// PublishOrderCancelled publica um evento de pedido cancelado
func (p *TypedEventPublisher) PublishOrderCancelled(ctx context.Context, data OrderCancelledEvent) error {
	event := contracts.Event{
		Type:      EventTypeOrderCancelled,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id":      data.OrderID,
			"user_id":       data.UserID,
			"reason":        data.Reason,
			"refund_amount": data.RefundAmount,
		},
	}
	return p.bus.Publish(ctx, event)
}

// TypedEventHandler é um wrapper para handlers tipados
type TypedEventHandler[T any] struct {
	eventType string
	handler   func(ctx context.Context, data T) error
}

// NewTypedHandler cria um novo handler tipado
func NewTypedHandler[T any](eventType string, handler func(ctx context.Context, data T) error) *TypedEventHandler[T] {
	return &TypedEventHandler[T]{
		eventType: eventType,
		handler:   handler,
	}
}

// AsContractHandler converte para contracts.EventHandler
func (h *TypedEventHandler[T]) AsContractHandler() contracts.EventHandler {
	return func(ctx context.Context, event contracts.Event) error {
		// Deserializar payload para o tipo correto
		var data T

		// Converter payload para JSON e depois para o tipo T
		jsonData, err := json.Marshal(event.Payload)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(jsonData, &data); err != nil {
			return err
		}

		return h.handler(ctx, data)
	}
}

// SubscribeTyped registra um handler tipado no EventBus
func SubscribeTyped[T any](bus *EventBus, eventType string, handler func(ctx context.Context, data T) error) error {
	typedHandler := NewTypedHandler(eventType, handler)
	return bus.Subscribe(eventType, typedHandler.AsContractHandler())
}
