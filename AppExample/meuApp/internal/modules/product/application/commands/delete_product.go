package commands

import (
	"context"
	"fmt"

	"meuApp/internal/modules/product/ports"
	"meuApp/pkg/contracts"
)

// DeleteProductCommand representa o comando para deletar um produto
type DeleteProductCommand struct {
	ID string
}

// DeleteProductHandler processa a exclusão de produtos
type DeleteProductHandler struct {
	productRepo ports.ProductRepository
	logger      contracts.Logger
}

// NewDeleteProductHandler cria uma nova instância do handler
func NewDeleteProductHandler(
	productRepo ports.ProductRepository,
	logger contracts.Logger,
) *DeleteProductHandler {
	return &DeleteProductHandler{
		productRepo: productRepo,
		logger:      logger,
	}
}

// Handle executa o comando de exclusão de produto
func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
	h.logger.Info(fmt.Sprintf("Deleting product: %s", cmd.ID))

	// Verificar se o produto existe
	_, err := h.productRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Product not found: %s", cmd.ID))
		return fmt.Errorf("product not found: %w", err)
	}

	// Deletar produto
	if err := h.productRepo.Delete(ctx, cmd.ID); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to delete product: %v", err))
		return fmt.Errorf("failed to delete product: %w", err)
	}

	h.logger.Info(fmt.Sprintf("Product deleted successfully: %s", cmd.ID))
	return nil
}
