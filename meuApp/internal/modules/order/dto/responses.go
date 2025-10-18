package dto

import "time"

// OrderResponse representa a resposta com dados de um pedido
type OrderResponse struct {
ID        string          `json:"id"`
UserID    string          `json:"user_id"`
Items     []OrderItemResponse `json:"items"`
Status    string          `json:"status"`
Total     float64         `json:"total"`
CreatedAt time.Time       `json:"created_at"`
UpdatedAt time.Time       `json:"updated_at"`
}

// OrderItemResponse representa um item do pedido na resposta
type OrderItemResponse struct {
ProductID string  `json:"product_id"`
Quantity  int     `json:"quantity"`
Price     float64 `json:"price"`
}

// OrderListResponse representa a resposta com lista de pedidos
type OrderListResponse struct {
Orders []OrderResponse `json:"orders"`
Total  int             `json:"total"`
}

// MessageResponse representa uma resposta simples com mensagem
type MessageResponse struct {
Message string `json:"message"`
}
