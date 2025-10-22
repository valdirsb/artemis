package dto

// CreateOrderRequest representa a requisição de criação de pedido
type CreateOrderRequest struct {
UserID string            `json:"user_id" validate:"required"`
Items  []CreateOrderItem `json:"items" validate:"required,dive"`
}

// CreateOrderItem representa um item na criação de pedido
type CreateOrderItem struct {
ProductID string `json:"product_id" validate:"required"`
Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

// UpdateOrderStatusRequest representa a requisição de atualização de status
type UpdateOrderStatusRequest struct {
Status string `json:"status" validate:"required,oneof=pending confirmed shipped delivered cancelled"`
}
