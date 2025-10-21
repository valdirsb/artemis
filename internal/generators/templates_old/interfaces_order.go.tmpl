package contracts

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Interfaces para comunicação entre módulos (Ports)

// OrderService define operações de negócio relacionadas a pedidos
type OrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error)
	GetOrderByID(ctx context.Context, id string) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status OrderStatus) error
	CancelOrder(ctx context.Context, id string) error
}

// Repository Interfaces (Adapters)

// OrderRepository define a interface para persistência de pedidos
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByUserID(ctx context.Context, userID string) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id string) error
}

// Domain Models

// Order representa o modelo de domínio do pedido
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Status    OrderStatus `json:"status"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Request/Response DTOs

type CreateOrderRequest struct {
	UserID string            `json:"user_id" validate:"required"`
	Items  []CreateOrderItem `json:"items" validate:"required,dive"`
}

type CreateOrderItem struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

// Events para comunicação entre módulos

type OrderCreatedEvent struct {
	OrderID string  `json:"order_id"`
	UserID  string  `json:"user_id"`
	Total   float64 `json:"total"`
}

// Handler Interfaces

type OrderHandler interface {
	CreateOrder(ctx *gin.Context)
	GetOrder(ctx *gin.Context)
	GetOrdersByUser(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
}
