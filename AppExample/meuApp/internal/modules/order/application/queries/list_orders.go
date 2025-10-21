package queries

import (
	"context"

	"meuApp/internal/modules/order/ports"
	"meuApp/pkg/contracts"
)

// ListOrdersQuery representa a query para listar pedidos
type ListOrdersQuery struct {
	Page     int
	PageSize int
}

// ListOrdersHandler processa a listagem de pedidos
type ListOrdersHandler struct {
	orderRepo ports.OrderRepository
	logger    contracts.Logger
}

// NewListOrdersHandler cria uma nova instância do handler
func NewListOrdersHandler(
	orderRepo ports.OrderRepository,
	logger contracts.Logger,
) *ListOrdersHandler {
	return &ListOrdersHandler{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

// Handle executa a query de listagem de pedidos com paginação
func (h *ListOrdersHandler) Handle(ctx context.Context, query ListOrdersQuery) (*ports.PaginatedOrderResult, error) {
	h.logger.Info("Listing orders with pagination")

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	result, err := h.orderRepo.ListPaginated(ctx, query.Page, query.PageSize)
	if err != nil {
		h.logger.Error("Failed to list orders", contracts.Field{Key: "error", Value: err})
		return nil, err
	}

	h.logger.Info("Orders listed successfully")
	return result, nil
}
