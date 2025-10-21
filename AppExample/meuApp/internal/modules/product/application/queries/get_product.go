package queries

import (
	"context"
	"fmt"

	"meuApp/internal/modules/product"
	"meuApp/internal/modules/product/domain"
	"meuApp/internal/modules/product/ports"
	"meuApp/pkg/contracts"
	apperrors "meuApp/pkg/errors"
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

	if query.ID == "" {
		return nil, apperrors.NewValidationError("product ID cannot be empty", map[string]string{"field": "id"})
	}

	foundProduct, err := h.productRepo.GetByID(ctx, query.ID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to get product: %s", query.ID))
		return nil, apperrors.NewInfrastructureError("failed to retrieve product", err)
	}

	if foundProduct == nil {
		return nil, product.NewProductNotFoundError(query.ID)
	}

	return foundProduct, nil
}
