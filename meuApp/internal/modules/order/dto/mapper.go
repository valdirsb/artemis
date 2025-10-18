package dto

import (
"meuApp/internal/modules/order/domain"
"meuApp/internal/modules/order/ports"
)

// ToOrderResponse converte domain.Order para OrderResponse
func ToOrderResponse(order *domain.Order) OrderResponse {
items := make([]OrderItemResponse, len(order.Items))
for i, item := range order.Items {
items[i] = OrderItemResponse{
ProductID: item.ProductID,
Quantity:  item.Quantity,
Price:     item.Price,
}
}

return OrderResponse{
ID:        order.ID,
UserID:    order.UserID,
Items:     items,
Status:    string(order.Status),
Total:     order.Total,
CreatedAt: order.CreatedAt,
UpdatedAt: order.UpdatedAt,
}
}

// ToOrderResponseList converte lista de domain.Order para OrderListResponse
func ToOrderResponseList(orders []*domain.Order) OrderListResponse {
orderResponses := make([]OrderResponse, len(orders))
for i, order := range orders {
orderResponses[i] = ToOrderResponse(order)
}

return OrderListResponse{
Orders: orderResponses,
Total:  len(orders),
}
}

// ToCreateOrderItems converte DTOs para ports.CreateOrderItem
func ToCreateOrderItems(items []CreateOrderItem) []ports.CreateOrderItem {
result := make([]ports.CreateOrderItem, len(items))
for i, item := range items {
result[i] = ports.CreateOrderItem{
ProductID: item.ProductID,
Quantity:  item.Quantity,
}
}
return result
}
