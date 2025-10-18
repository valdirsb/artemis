package repository

import (
	"context"
	"fmt"

	"meuApp/internal/modules/order/domain"
	"meuApp/internal/modules/order/ports"

	"gorm.io/gorm"
)

// mysqlOrderRepository implementa a interface OrderRepository usando MySQL/GORM
type mysqlOrderRepository struct {
	db *gorm.DB
}

// NewMySQLOrderRepository cria uma nova instância do repositório MySQL
func NewMySQLOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &mysqlOrderRepository{
		db: db,
	}
}

// Create cria um novo pedido no banco de dados
func (r *mysqlOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	orderModel := &OrderModel{}
	orderModel.FromDomain(order)

	// Usar transação para garantir consistência entre order e order_items
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Criar o pedido
		if err := tx.Create(orderModel).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// GetByID busca um pedido pelo ID
func (r *mysqlOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var orderModel OrderModel

	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&orderModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	return orderModel.ToDomain(), nil
}

// GetByUserID busca todos os pedidos de um usuário
func (r *mysqlOrderRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	var orderModels []OrderModel

	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Find(&orderModels).Error; err != nil {
		return nil, fmt.Errorf("failed to get orders by user ID: %w", err)
	}

	orders := make([]*domain.Order, len(orderModels))
	for i, model := range orderModels {
		orders[i] = model.ToDomain()
	}

	return orders, nil
}

// Update atualiza um pedido existente
func (r *mysqlOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	orderModel := &OrderModel{}
	orderModel.FromDomain(order)

	// Usar transação para atualizar order e order_items
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Atualizar o pedido principal (sem os items)
		result := tx.Model(&OrderModel{}).Where("id = ?", order.ID).Updates(map[string]interface{}{
			"status":     orderModel.Status,
			"total":      orderModel.Total,
			"updated_at": orderModel.UpdatedAt,
		})

		if result.Error != nil {
			return fmt.Errorf("failed to update order: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("order not found")
		}

		// Remover items existentes
		if err := tx.Where("order_id = ?", order.ID).Delete(&OrderItemModel{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing order items: %w", err)
		}

		// Inserir novos items
		if len(orderModel.Items) > 0 {
			if err := tx.Create(&orderModel.Items).Error; err != nil {
				return fmt.Errorf("failed to create order items: %w", err)
			}
		}

		return nil
	})

	return err
}

// Delete remove um pedido do banco de dados
func (r *mysqlOrderRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remover items primeiro (foreign key constraint)
		if err := tx.Where("order_id = ?", id).Delete(&OrderItemModel{}).Error; err != nil {
			return fmt.Errorf("failed to delete order items: %w", err)
		}

		// Remover o pedido
		result := tx.Where("id = ?", id).Delete(&OrderModel{})
		if result.Error != nil {
			return fmt.Errorf("failed to delete order: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("order not found")
		}

		return nil
	})

	return err
}
