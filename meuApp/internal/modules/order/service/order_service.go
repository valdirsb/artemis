package service

import (
	"context"
	"errors"
	"time"

	"meuApp/internal/modules/order/domain"
	"meuApp/internal/modules/order/ports"
	productDomain "meuApp/internal/modules/product/domain"
	productPorts "meuApp/internal/modules/product/ports"
	userPorts "meuApp/internal/modules/user/ports"
	"meuApp/pkg/contracts"
	"meuApp/pkg/events"

	"github.com/google/uuid"
)

// OrderService implementa a lógica de negócio do módulo de pedidos
type OrderService struct {
	orderRepo      ports.OrderRepository
	productService productPorts.ProductService // Para validar produtos e verificar estoque
	userService    userPorts.UserService       // Para validar usuários
	eventPublisher contracts.EventPublisher
}

// NewOrderService cria uma nova instância do serviço de pedidos
func NewOrderService(
	orderRepo ports.OrderRepository,
	productService productPorts.ProductService,
	userService userPorts.UserService,
	eventPublisher contracts.EventPublisher,
) ports.OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		productService: productService,
		userService:    userService,
		eventPublisher: eventPublisher,
	}
}

// CreateOrder cria um novo pedido
func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []ports.CreateOrderItem) (*domain.Order, error) {
	// Validar se o usuário existe
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("invalid user ID")
	}

	// Cache de produtos para evitar múltiplas consultas
	type Product = productDomain.Product
	productCache := make(map[string]*Product)

	// Validar produtos, calcular preços e cachear produtos
	orderItems := make([]domain.OrderItem, len(items))
	for i, item := range items {
		// Buscar produto (apenas uma vez por produto único)
		product, exists := productCache[item.ProductID]
		if !exists {
			product, err = s.productService.GetProductByID(ctx, item.ProductID)
			if err != nil || product == nil {
				return nil, errors.New("invalid product ID: " + item.ProductID)
			}
			// Cachear produto para futuras consultas
			productCache[item.ProductID] = product
		}

		// Verificar se há estoque suficiente
		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		orderItems[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Usar o preço atual do produto
		}
	}

	// Gerar ID único para o pedido
	orderID := uuid.New().String()

	// Criar entidade de domínio
	order, err := domain.NewOrder(orderID, userID, orderItems)
	if err != nil {
		return nil, err
	}

	// Criar aggregate e validar
	orderAggregate := domain.NewOrderAggregate(order)
	if err := orderAggregate.IsValid(); err != nil {
		return nil, err
	}

	// Reduzir estoque dos produtos (usar cache para evitar consultas desnecessárias)
	updatedProducts := make(map[string]int) // Controla quantidade consumida por produto

	for _, item := range orderItems {
		// Usar produto do cache
		product := productCache[item.ProductID]

		// Somar quantidade total se o mesmo produto aparece múltiplas vezes no pedido
		totalQuantityUsed := updatedProducts[item.ProductID] + item.Quantity

		// Calcular novo estoque (estoque atual - quantidade do item)
		newStock := product.Stock - item.Quantity
		if newStock < 0 {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		// Atualizar cache local apenas se necessário
		if updatedProducts[item.ProductID] == 0 || totalQuantityUsed != updatedProducts[item.ProductID] {

			// Atualizar cache com novo estoque
			product.Stock = newStock
			updatedProducts[item.ProductID] = totalQuantityUsed
		}
	}

	// Aplicar todas as atualizações de estoque no banco
	for k, v := range productCache {
		if err := s.productService.UpdateStock(ctx, k, v.Stock); err != nil {
			return nil, errors.New("failed to update stock for product: " + k)
		}
	}

	// Persistir pedido
	orderToSave := orderAggregate.GetOrder()
	if err := s.orderRepo.Create(ctx, orderToSave); err != nil {
		// Se falhar, reverter estoque usando o cache (mais eficiente)
		for productID, quantityUsed := range updatedProducts {
			if product := productCache[productID]; product != nil {
				restoredStock := product.Stock + quantityUsed
				s.productService.UpdateStock(ctx, productID, restoredStock) // Reverter
			}
		}
		return nil, errors.New("failed to create order")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.OrderCreatedEventType,
		Timestamp: time.Now(),
		Payload: contracts.OrderCreatedEvent{
			OrderID: orderID,
			UserID:  userID,
			Total:   orderToSave.Total,
		},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		// Log do erro mas não falhar a operação
	}

	return orderToSave, nil
}

// GetOrderByID obtém um pedido por ID
func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("order ID cannot be empty")
	}

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	return order, nil
}

// GetOrdersByUserID obtém todos os pedidos de um usuário
func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	// Validar se o usuário existe
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("invalid user ID")
	}

	orders, err := s.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateOrderStatus atualiza o status de um pedido
func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	if id == "" {
		return errors.New("order ID cannot be empty")
	}

	// Buscar pedido existente
	existingOrder, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingOrder == nil {
		return errors.New("order not found")
	}

	// Criar aggregate e atualizar status
	orderAggregate := domain.NewOrderAggregate(existingOrder)
	if err := orderAggregate.UpdateStatus(status); err != nil {
		return err
	}

	// Persistir alterações
	updatedOrder := orderAggregate.GetOrder()
	if err := s.orderRepo.Update(ctx, updatedOrder); err != nil {
		return errors.New("failed to update order")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.OrderStatusUpdatedEventType,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id":   id,
			"new_status": string(status),
			"old_status": string(existingOrder.Status),
		},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		// Log do erro mas não falhar a operação
	}

	return nil
}

// CancelOrder cancela um pedido
func (s *OrderService) CancelOrder(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order ID cannot be empty")
	}

	// Buscar pedido existente
	existingOrder, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingOrder == nil {
		return errors.New("order not found")
	}

	// Criar aggregate e cancelar
	orderAggregate := domain.NewOrderAggregate(existingOrder)
	if err := orderAggregate.Cancel(); err != nil {
		return err
	}

	// Restaurar estoque dos produtos (se o pedido ainda estava pendente/confirmado)
	if existingOrder.Status == domain.OrderStatusPending || existingOrder.Status == domain.OrderStatusConfirmed {
		for _, item := range existingOrder.Items {
			// Buscar produto para ter o estoque atual e restaurar
			if product, err := s.productService.GetProductByID(ctx, item.ProductID); err == nil && product != nil {
				restoredStock := product.Stock + item.Quantity
				if err := s.productService.UpdateStock(ctx, item.ProductID, restoredStock); err != nil {
					// Log do erro mas continuar o cancelamento
				}
			}
		}
	}

	// Persistir alterações
	cancelledOrder := orderAggregate.GetOrder()
	if err := s.orderRepo.Update(ctx, cancelledOrder); err != nil {
		return errors.New("failed to cancel order")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.OrderCancelledEventType,
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"order_id": id,
			"user_id":  existingOrder.UserID,
			"total":    existingOrder.Total,
		},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		// Log do erro mas não falhar a operação
	}

	return nil
}
