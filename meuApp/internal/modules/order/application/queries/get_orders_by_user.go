package queries

import (
"context"
"fmt"

"meuApp/internal/modules/order/domain"
"meuApp/internal/modules/order/ports"
"meuApp/pkg/contracts"
)

// GetOrdersByUserQuery representa a query para buscar pedidos de um usuário
type GetOrdersByUserQuery struct {
UserID string
}

// GetOrdersByUserHandler processa a busca de pedidos por usuário
type GetOrdersByUserHandler struct {
orderRepo ports.OrderRepository
logger    contracts.Logger
}

// NewGetOrdersByUserHandler cria uma nova instância do handler
func NewGetOrdersByUserHandler(
orderRepo ports.OrderRepository,
logger contracts.Logger,
) *GetOrdersByUserHandler {
return &GetOrdersByUserHandler{
orderRepo: orderRepo,
logger:    logger,
}
}

// Handle executa a query de busca de pedidos por usuário
func (h *GetOrdersByUserHandler) Handle(ctx context.Context, query GetOrdersByUserQuery) ([]*domain.Order, error) {
h.logger.Info(fmt.Sprintf("Getting orders for user: %s", query.UserID))

orders, err := h.orderRepo.GetByUserID(ctx, query.UserID)
if err != nil {
h.logger.Error(fmt.Sprintf("Failed to get orders for user %s: %v", query.UserID, err))
return nil, fmt.Errorf("failed to get orders: %w", err)
}

h.logger.Info(fmt.Sprintf("Found %d orders for user %s", len(orders), query.UserID))
return orders, nil
}
