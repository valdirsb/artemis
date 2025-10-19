package services

import (
	"context"

	"meuApp/internal/modules/order/application/commands"
	"meuApp/internal/modules/order/application/queries"
	"meuApp/internal/modules/order/domain"
	"meuApp/internal/modules/order/ports"
)

// OrderApplicationService é o serviço de aplicação que orquestra os use cases
type OrderApplicationService struct {
	createHandler          *commands.CreateOrderHandler
	updateStatusHandler    *commands.UpdateOrderStatusHandler
	cancelHandler          *commands.CancelOrderHandler
	getOrderHandler        *queries.GetOrderHandler
	getOrdersByUserHandler *queries.GetOrdersByUserHandler
}

// NewOrderApplicationService cria uma nova instância do serviço de aplicação
func NewOrderApplicationService(
	createHandler *commands.CreateOrderHandler,
	updateStatusHandler *commands.UpdateOrderStatusHandler,
	cancelHandler *commands.CancelOrderHandler,
	getOrderHandler *queries.GetOrderHandler,
	getOrdersByUserHandler *queries.GetOrdersByUserHandler,
) *OrderApplicationService {
	return &OrderApplicationService{
		createHandler:          createHandler,
		updateStatusHandler:    updateStatusHandler,
		cancelHandler:          cancelHandler,
		getOrderHandler:        getOrderHandler,
		getOrdersByUserHandler: getOrdersByUserHandler,
	}
}

// CreateOrder cria um novo pedido
func (s *OrderApplicationService) CreateOrder(ctx context.Context, userID string, items []ports.CreateOrderItem) (*domain.Order, error) {
	cmdItems := make([]commands.CreateOrderItemCommand, len(items))
	for i, item := range items {
		cmdItems[i] = commands.CreateOrderItemCommand{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	cmd := commands.CreateOrderCommand{
		UserID: userID,
		Items:  cmdItems,
	}
	return s.createHandler.Handle(ctx, cmd)
}

// GetOrderByID busca um pedido por ID
func (s *OrderApplicationService) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	query := queries.GetOrderQuery{ID: id}
	return s.getOrderHandler.Handle(ctx, query)
}

// GetOrdersByUserID busca pedidos de um usuário
func (s *OrderApplicationService) GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	query := queries.GetOrdersByUserQuery{UserID: userID}
	return s.getOrdersByUserHandler.Handle(ctx, query)
}

// UpdateOrderStatus atualiza o status de um pedido
func (s *OrderApplicationService) UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	cmd := commands.UpdateOrderStatusCommand{
		ID:     id,
		Status: status,
	}
	return s.updateStatusHandler.Handle(ctx, cmd)
}

// CancelOrder cancela um pedido
func (s *OrderApplicationService) CancelOrder(ctx context.Context, id string) error {
	cmd := commands.CancelOrderCommand{ID: id}
	return s.cancelHandler.Handle(ctx, cmd)
}
