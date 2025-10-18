package ports

import (
"context"
"meuApp/internal/modules/order/domain"
)

// ===== PRIMARY PORTS (Application Services) =====

// OrderService define as operações de negócio do módulo de pedido
type OrderService interface {
CreateOrder(ctx context.Context, userID string, items []CreateOrderItem) (*domain.Order, error)
GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) error
CancelOrder(ctx context.Context, id string) error
}

// ===== SECONDARY PORTS (Adapters/Dependencies) =====

// OrderRepository define a interface para persistência de pedidos
type OrderRepository interface {
Create(ctx context.Context, order *domain.Order) error
GetByID(ctx context.Context, id string) (*domain.Order, error)
GetByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
Update(ctx context.Context, order *domain.Order) error
Delete(ctx context.Context, id string) error
}

// CreateOrderItem representa um item na criação de pedido
type CreateOrderItem struct {
ProductID string
Quantity  int
}
