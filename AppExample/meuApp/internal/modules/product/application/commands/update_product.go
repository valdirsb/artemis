package commands

import (
"context"
"fmt"
"time"

"meuApp/internal/modules/product/domain"
"meuApp/internal/modules/product/ports"
"meuApp/pkg/contracts"
)

// UpdateProductCommand representa o comando para atualizar um produto
type UpdateProductCommand struct {
ID          string
Name        *string
Description *string
CategoryID  *string
Price       *float64
}

// UpdateProductHandler processa a atualização de produtos
type UpdateProductHandler struct {
productRepo ports.ProductRepository
logger      contracts.Logger
}

// NewUpdateProductHandler cria uma nova instância do handler
func NewUpdateProductHandler(
productRepo ports.ProductRepository,
logger contracts.Logger,
) *UpdateProductHandler {
return &UpdateProductHandler{
productRepo: productRepo,
logger:      logger,
}
}

// Handle executa o comando de atualização de produto
func (h *UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProductCommand) (*domain.Product, error) {
h.logger.Info(fmt.Sprintf("Updating product: %s", cmd.ID))

// Buscar produto existente
product, err := h.productRepo.GetByID(ctx, cmd.ID)
if err != nil {
h.logger.Error(fmt.Sprintf("Product not found: %s", cmd.ID))
return nil, fmt.Errorf("product not found: %w", err)
}

// Criar aggregate para aplicar mudanças
aggregate := domain.NewProductAggregate(product)

// Aplicar mudanças via métodos do aggregate
if cmd.Name != nil {
if err := aggregate.UpdateName(*cmd.Name); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update name: %v", err))
return nil, fmt.Errorf("invalid name: %w", err)
}
}

if cmd.Description != nil {
if err := aggregate.UpdateDescription(*cmd.Description); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update description: %v", err))
return nil, fmt.Errorf("invalid description: %w", err)
}
}

if cmd.Price != nil {
if err := aggregate.UpdatePrice(*cmd.Price); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update price: %v", err))
return nil, fmt.Errorf("invalid price: %w", err)
}
}

if cmd.CategoryID != nil {
// Atualizar categoria diretamente (aggregate não tem método específico)
product.CategoryID = *cmd.CategoryID
product.UpdatedAt = time.Now()
}

// Persistir mudanças
updatedProduct := aggregate.GetProduct()
if err := h.productRepo.Update(ctx, updatedProduct); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update product: %v", err))
return nil, fmt.Errorf("failed to update product: %w", err)
}

h.logger.Info(fmt.Sprintf("Product updated successfully: %s", cmd.ID))
return updatedProduct, nil
}
