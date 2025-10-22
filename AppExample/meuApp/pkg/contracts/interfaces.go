package contracts

import (
	"context"
	"time"
)

// Event Publisher para comunicação assíncrona entre módulos
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler) error
}

type EventHandler func(ctx context.Context, event Event) error

// Events para comunicação entre módulos

type Event struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}
