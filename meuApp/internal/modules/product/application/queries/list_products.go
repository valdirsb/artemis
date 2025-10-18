package queries

import (
"context"
"fmt"

"meuApp/internal/modules/product/domain"
"meuApp/internal/modules/product/ports"
"meuApp/pkg/contracts"
)

// ListProductsQuery representa a query para listar produtos
type ListProductsQuery struct {
CategoryID *string
MinPrice   *float64
MaxPrice   *float64
InStock    *bool
}

// ListProductsHandler processa a listagem de produtos
type ListProductsHandler struct {
productRepo ports.ProductRepository
logger      contracts.Logger
}

// NewListProductsHandler cria uma nova instância do handler
func NewListProductsHandler(
productRepo ports.ProductRepository,
logger contracts.Logger,
) *ListProductsHandler {
return &ListProductsHandler{
productRepo: productRepo,
logger:      logger,
}
}

// Handle executa a query de listagem de produtos
func (h *ListProductsHandler) Handle(ctx context.Context, query ListProductsQuery) ([]*domain.Product, error) {
h.logger.Info("Listing products with filters")

filters := ports.ProductFilters{
CategoryID: query.CategoryID,
MinPrice:   query.MinPrice,
MaxPrice:   query.MaxPrice,
InStock:    query.InStock,
}

products, err := h.productRepo.List(ctx, filters)
if err != nil {
h.logger.Error(fmt.Sprintf("Failed to list products: %v", err))
return nil, fmt.Errorf("failed to list products: %w", err)
}

h.logger.Info(fmt.Sprintf("Found %d products", len(products)))
return products, nil
}
