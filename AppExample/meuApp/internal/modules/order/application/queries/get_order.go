package queries

import (
"context"
"fmt"

"meuApp/internal/modules/order/domain"
"meuApp/internal/modules/order/ports"
"meuApp/pkg/contracts"
)

// GetOrderQuery representa a query para buscar um pedido
type GetOrderQuery struct {
ID string
}

// GetOrderHandler processa a busca de pedido por ID
type GetOrderHandler struct {
orderRepo ports.OrderRepository
logger    contracts.Logger
}

// NewGetOrderHandler cria uma nova instância do handler
func NewGetOrderHandler(
orderRepo ports.OrderRepository,
logger contracts.Logger,
) *GetOrderHandler {
return &GetOrderHandler{
orderRepo: orderRepo,
logger:    logger,
}
}

// Handle executa a query de busca de pedido
func (h *GetOrderHandler) Handle(ctx context.Context, query GetOrderQuery) (*domain.Order, error) {
h.logger.Info(fmt.Sprintf("Getting order: %s", query.ID))

order, err := h.orderRepo.GetByID(ctx, query.ID)
if err != nil {
h.logger.Error(fmt.Sprintf("Order not found: %s", query.ID))
return nil, fmt.Errorf("order not found: %w", err)
}

return order, nil
}
