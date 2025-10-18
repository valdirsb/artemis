package queries

import (
"context"
"fmt"

"meuApp/internal/modules/product/domain"
"meuApp/internal/modules/product/ports"
"meuApp/pkg/contracts"
)

// GetProductQuery representa a query para buscar um produto
type GetProductQuery struct {
ID string
}

// GetProductHandler processa a busca de produto por ID
type GetProductHandler struct {
productRepo ports.ProductRepository
logger      contracts.Logger
}

// NewGetProductHandler cria uma nova instância do handler
func NewGetProductHandler(
productRepo ports.ProductRepository,
logger contracts.Logger,
) *GetProductHandler {
return &GetProductHandler{
productRepo: productRepo,
logger:      logger,
}
}

// Handle executa a query de busca de produto
func (h *GetProductHandler) Handle(ctx context.Context, query GetProductQuery) (*domain.Product, error) {
h.logger.Info(fmt.Sprintf("Getting product: %s", query.ID))

product, err := h.productRepo.GetByID(ctx, query.ID)
if err != nil {
h.logger.Error(fmt.Sprintf("Product not found: %s", query.ID))
return nil, fmt.Errorf("product not found: %w", err)
}

return product, nil
}
