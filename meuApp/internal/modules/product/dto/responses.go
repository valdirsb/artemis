package dto

import "time"

// ProductResponse representa a resposta com dados de um produto
type ProductResponse struct {
ID          string    `json:"id"`
Name        string    `json:"name"`
Description string    `json:"description"`
Price       float64   `json:"price"`
Stock       int       `json:"stock"`
CategoryID  string    `json:"category_id"`
CreatedAt   time.Time `json:"created_at"`
UpdatedAt   time.Time `json:"updated_at"`
}

// ProductListResponse representa a resposta com lista de produtos
type ProductListResponse struct {
Products []ProductResponse `json:"products"`
Total    int               `json:"total"`
}
