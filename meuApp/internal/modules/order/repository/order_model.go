package repository

import (
	"time"

	"meuApp/internal/modules/order/domain"
)

// OrderModel representa a estrutura da tabela orders no banco
type OrderModel struct {
	ID        string           `gorm:"primaryKey;size:36"`
	UserID    string           `gorm:"size:36;not null;index"`
	Status    string           `gorm:"size:20;not null"`
	Total     float64          `gorm:"type:decimal(10,2);not null"`
	Items     []OrderItemModel `gorm:"foreignKey:OrderID"`
	CreatedAt time.Time        `gorm:"autoCreateTime"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime"`
}

// TableName especifica o nome da tabela
func (OrderModel) TableName() string {
	return "orders"
}

// OrderItemModel representa a estrutura da tabela order_items no banco
type OrderItemModel struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   string  `gorm:"size:36;not null;index"`
	ProductID string  `gorm:"size:36;not null"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
}

// TableName especifica o nome da tabela
func (OrderItemModel) TableName() string {
	return "order_items"
}

// ToDomain converte OrderModel para domain.Order
func (o *OrderModel) ToDomain() *domain.Order {
	items := make([]domain.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	return &domain.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		Items:     items,
		Status:    domain.OrderStatus(o.Status),
		Total:     o.Total,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

// FromDomain converte domain.Order para OrderModel
func (o *OrderModel) FromDomain(order *domain.Order) {
	o.ID = order.ID
	o.UserID = order.UserID
	o.Status = string(order.Status)
	o.Total = order.Total
	o.CreatedAt = order.CreatedAt
	o.UpdatedAt = order.UpdatedAt

	// Converter items
	o.Items = make([]OrderItemModel, len(order.Items))
	for i, item := range order.Items {
		o.Items[i] = OrderItemModel{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
}
