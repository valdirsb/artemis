package commands

import (
	"context"
	"fmt"

	"meuApp/internal/modules/order"
	"meuApp/internal/modules/order/domain"
	"meuApp/internal/modules/order/ports"
	productPorts "meuApp/internal/modules/product/ports"
	userPorts "meuApp/internal/modules/user/ports"
	"meuApp/pkg/contracts"
	apperrors "meuApp/pkg/errors"
)

// CreateOrderCommand representa o comando para criar um pedido
type CreateOrderCommand struct {
	UserID string
	Items  []CreateOrderItemCommand
}

// CreateOrderItemCommand representa um item do pedido
type CreateOrderItemCommand struct {
	ProductID string
	Quantity  int
}

// CreateOrderHandler processa a criação de pedidos
type CreateOrderHandler struct {
	orderRepo   ports.OrderRepository
	userRepo    userPorts.UserRepository
	productRepo productPorts.ProductRepository
	eventBus    contracts.EventPublisher
	logger      contracts.Logger
}

// NewCreateOrderHandler cria uma nova instância do handler
func NewCreateOrderHandler(
	orderRepo ports.OrderRepository,
	userRepo userPorts.UserRepository,
	productRepo productPorts.ProductRepository,
	eventBus contracts.EventPublisher,
	logger contracts.Logger,
) *CreateOrderHandler {
	return &CreateOrderHandler{
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
		eventBus:    eventBus,
		logger:      logger,
	}
}

// Handle executa o comando de criação de pedido
func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*domain.Order, error) {
	h.logger.Info(fmt.Sprintf("Creating order for user: %s", cmd.UserID))

	// Validações de entrada
	if cmd.UserID == "" {
		return nil, order.ErrInvalidUserID
	}
	if len(cmd.Items) == 0 {
		return nil, order.ErrEmptyOrderItems
	}

	// Validar se o usuário existe
	user, err := h.userRepo.GetByID(ctx, cmd.UserID)
	if err != nil || user == nil {
		h.logger.Error(fmt.Sprintf("User not found: %s", cmd.UserID))
		return nil, order.NewUserNotFoundError(cmd.UserID)
	}

	// Validar produtos e calcular preços
	orderItems := make([]domain.OrderItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		if item.Quantity <= 0 {
			return nil, order.ErrInvalidQuantity
		}

		product, err := h.productRepo.GetByID(ctx, item.ProductID)
		if err != nil || product == nil {
			h.logger.Error(fmt.Sprintf("Product not found: %s", item.ProductID))
			return nil, order.NewProductNotFoundError(item.ProductID)
		}

		// Verificar estoque disponível
		if product.Stock < item.Quantity {
			h.logger.Warn(fmt.Sprintf("Insufficient stock for product %s: available=%d, requested=%d",
				item.ProductID, product.Stock, item.Quantity))
			return nil, order.NewInsufficientStockError(item.ProductID, product.Stock, item.Quantity)
		}

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	// Gerar ID (pode ser substituído por UUID generator)
	orderID := fmt.Sprintf("order_%d", ctx.Value("timestamp"))

	// Criar aggregate do pedido
	newOrder, err := domain.NewOrder(orderID, cmd.UserID, orderItems)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to create order aggregate: %v", err))
		return nil, apperrors.WrapError(err, "invalid order data")
	}

	// Salvar no repositório
	if err := h.orderRepo.Create(ctx, newOrder); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to save order: %v", err))
		return nil, apperrors.NewInfrastructureError("failed to create order", err)
	}

	// Publicar evento (assíncrono)
	go func() {
		event := contracts.Event{
			Type: "order.created",
			Payload: map[string]interface{}{
				"order_id": newOrder.ID,
				"user_id":  newOrder.UserID,
				"total":    newOrder.Total,
				"items":    len(newOrder.Items),
			},
			Timestamp: newOrder.CreatedAt,
		}
		h.eventBus.Publish(context.Background(), event)
		h.logger.Info(fmt.Sprintf("Event published: order.created for %s", newOrder.ID))
	}()

	h.logger.Info(fmt.Sprintf("Order created successfully: %s (total: %.2f)", newOrder.ID, newOrder.Total))
	return newOrder, nil
}
