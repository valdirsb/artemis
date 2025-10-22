package commands

import (
"context"
"fmt"

"meuApp/internal/modules/product/domain"
"meuApp/internal/modules/product/ports"
"meuApp/pkg/contracts"
)

// UpdateStockCommand representa o comando para atualizar estoque
type UpdateStockCommand struct {
ID       string
Quantity int
}

// UpdateStockHandler processa a atualização de estoque
type UpdateStockHandler struct {
productRepo ports.ProductRepository
eventBus    contracts.EventPublisher
logger      contracts.Logger
}

// NewUpdateStockHandler cria uma nova instância do handler
func NewUpdateStockHandler(
productRepo ports.ProductRepository,
eventBus contracts.EventPublisher,
logger contracts.Logger,
) *UpdateStockHandler {
return &UpdateStockHandler{
productRepo: productRepo,
eventBus:    eventBus,
logger:      logger,
}
}

// Handle executa o comando de atualização de estoque
func (h *UpdateStockHandler) Handle(ctx context.Context, cmd UpdateStockCommand) error {
h.logger.Info(fmt.Sprintf("Updating stock for product: %s, quantity: %d", cmd.ID, cmd.Quantity))

// Buscar produto existente
product, err := h.productRepo.GetByID(ctx, cmd.ID)
if err != nil {
h.logger.Error(fmt.Sprintf("Product not found: %s", cmd.ID))
return fmt.Errorf("product not found: %w", err)
}

// Criar aggregate para aplicar mudanças
aggregate := domain.NewProductAggregate(product)

// Atualizar estoque
if err := aggregate.UpdateStock(cmd.Quantity); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update stock: %v", err))
return fmt.Errorf("invalid stock quantity: %w", err)
}

// Persistir mudanças
updatedProduct := aggregate.GetProduct()
if err := h.productRepo.Update(ctx, updatedProduct); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update product stock: %v", err))
return fmt.Errorf("failed to update stock: %w", err)
}

// Publicar evento se estoque estiver baixo
if updatedProduct.Stock < 10 {
go func() {
event := contracts.Event{
Type: "product.low_stock",
Payload: map[string]interface{}{
"product_id": updatedProduct.ID,
"stock":      updatedProduct.Stock,
"name":       updatedProduct.Name,
},
Timestamp: updatedProduct.UpdatedAt,
}
h.eventBus.Publish(context.Background(), event)
h.logger.Warn(fmt.Sprintf("Low stock alert: %s (stock: %d)", updatedProduct.ID, updatedProduct.Stock))
}()
}

h.logger.Info(fmt.Sprintf("Stock updated successfully for product: %s", cmd.ID))
return nil
}
