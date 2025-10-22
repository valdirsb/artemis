package commands

import (
"context"
"fmt"

"meuApp/internal/modules/order/domain"
"meuApp/internal/modules/order/ports"
"meuApp/pkg/contracts"
)

// UpdateOrderStatusCommand representa o comando para atualizar status de pedido
type UpdateOrderStatusCommand struct {
ID     string
Status domain.OrderStatus
}

// UpdateOrderStatusHandler processa a atualização de status
type UpdateOrderStatusHandler struct {
orderRepo ports.OrderRepository
eventBus  contracts.EventPublisher
logger    contracts.Logger
}

// NewUpdateOrderStatusHandler cria uma nova instância do handler
func NewUpdateOrderStatusHandler(
orderRepo ports.OrderRepository,
eventBus contracts.EventPublisher,
logger contracts.Logger,
) *UpdateOrderStatusHandler {
return &UpdateOrderStatusHandler{
orderRepo: orderRepo,
eventBus:  eventBus,
logger:    logger,
}
}

// Handle executa o comando de atualização de status
func (h *UpdateOrderStatusHandler) Handle(ctx context.Context, cmd UpdateOrderStatusCommand) error {
h.logger.Info(fmt.Sprintf("Updating order status: %s to %s", cmd.ID, cmd.Status))

// Buscar pedido existente
order, err := h.orderRepo.GetByID(ctx, cmd.ID)
if err != nil {
h.logger.Error(fmt.Sprintf("Order not found: %s", cmd.ID))
return fmt.Errorf("order not found: %w", err)
}

// Criar aggregate para aplicar mudanças
aggregate := domain.NewOrderAggregate(order)

// Atualizar status
if err := aggregate.UpdateStatus(cmd.Status); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update status: %v", err))
return fmt.Errorf("invalid status transition: %w", err)
}

// Persistir mudanças
updatedOrder := aggregate.GetOrder()
if err := h.orderRepo.Update(ctx, updatedOrder); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update order: %v", err))
return fmt.Errorf("failed to update order: %w", err)
}

// Publicar evento
go func() {
event := contracts.Event{
Type: "order.status_updated",
Payload: map[string]interface{}{
"order_id": updatedOrder.ID,
"status":   string(updatedOrder.Status),
"user_id":  updatedOrder.UserID,
},
Timestamp: updatedOrder.UpdatedAt,
}
h.eventBus.Publish(context.Background(), event)
}()

h.logger.Info(fmt.Sprintf("Order status updated successfully: %s", cmd.ID))
return nil
}
