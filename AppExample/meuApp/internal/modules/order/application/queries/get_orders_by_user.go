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
	UserID   string
	Page     int
	PageSize int
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

// Handle executa a query de busca de pedidos por usuário (sem paginação - mantido para compatibilidade)
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

// HandlePaginated executa a query de busca de pedidos por usuário com paginação
func (h *GetOrdersByUserHandler) HandlePaginated(ctx context.Context, query GetOrdersByUserQuery) (*ports.PaginatedOrderResult, error) {
	h.logger.Info(fmt.Sprintf("Getting paginated orders for user: %s (page: %d, pageSize: %d)", query.UserID, query.Page, query.PageSize))

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	result, err := h.orderRepo.GetByUserIDPaginated(ctx, query.UserID, query.Page, query.PageSize)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to get paginated orders for user %s: %v", query.UserID, err))
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	h.logger.Info(fmt.Sprintf("Found %d orders for user %s (page %d of %d)", len(result.Items), query.UserID, result.Page, result.TotalPages))
	return result, nil
}
