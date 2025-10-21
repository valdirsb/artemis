package commands

import (
"context"
"fmt"

"meuApp/internal/modules/order/domain"
"meuApp/internal/modules/order/ports"
"meuApp/pkg/contracts"
)

// CancelOrderCommand representa o comando para cancelar um pedido
type CancelOrderCommand struct {
ID string
}

// CancelOrderHandler processa o cancelamento de pedidos
type CancelOrderHandler struct {
orderRepo ports.OrderRepository
eventBus  contracts.EventPublisher
logger    contracts.Logger
}

// NewCancelOrderHandler cria uma nova instância do handler
func NewCancelOrderHandler(
orderRepo ports.OrderRepository,
eventBus contracts.EventPublisher,
logger contracts.Logger,
) *CancelOrderHandler {
return &CancelOrderHandler{
orderRepo: orderRepo,
eventBus:  eventBus,
logger:    logger,
}
}

// Handle executa o comando de cancelamento de pedido
func (h *CancelOrderHandler) Handle(ctx context.Context, cmd CancelOrderCommand) error {
h.logger.Info(fmt.Sprintf("Cancelling order: %s", cmd.ID))

// Buscar pedido existente
order, err := h.orderRepo.GetByID(ctx, cmd.ID)
if err != nil {
h.logger.Error(fmt.Sprintf("Order not found: %s", cmd.ID))
return fmt.Errorf("order not found: %w", err)
}

// Criar aggregate para aplicar mudanças
aggregate := domain.NewOrderAggregate(order)

// Cancelar pedido
if err := aggregate.Cancel(); err != nil {
h.logger.Error(fmt.Sprintf("Failed to cancel order: %v", err))
return fmt.Errorf("cannot cancel order: %w", err)
}

// Persistir mudanças
cancelledOrder := aggregate.GetOrder()
if err := h.orderRepo.Update(ctx, cancelledOrder); err != nil {
h.logger.Error(fmt.Sprintf("Failed to update order: %v", err))
return fmt.Errorf("failed to cancel order: %w", err)
}

// Publicar evento
go func() {
event := contracts.Event{
Type: "order.cancelled",
Payload: map[string]interface{}{
"order_id": cancelledOrder.ID,
"user_id":  cancelledOrder.UserID,
"total":    cancelledOrder.Total,
},
Timestamp: cancelledOrder.UpdatedAt,
}
h.eventBus.Publish(context.Background(), event)
h.logger.Info(fmt.Sprintf("Event published: order.cancelled for %s", cancelledOrder.ID))
}()

h.logger.Info(fmt.Sprintf("Order cancelled successfully: %s", cmd.ID))
return nil
}
