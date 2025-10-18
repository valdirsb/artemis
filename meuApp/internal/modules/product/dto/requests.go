package dto

// CreateProductRequest representa a requisição de criação de produto
type CreateProductRequest struct {
Name        string  `json:"name" validate:"required"`
Description string  `json:"description"`
CategoryID  string  `json:"category_id" validate:"required"`
Price       float64 `json:"price" validate:"required,gt=0"`
Stock       int     `json:"stock" validate:"required,gte=0"`
}

// UpdateProductRequest representa a requisição de atualização de produto
type UpdateProductRequest struct {
Name        *string  `json:"name,omitempty"`
Description *string  `json:"description,omitempty"`
CategoryID  *string  `json:"category_id,omitempty"`
Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
}

// UpdateStockRequest representa a requisição de atualização de estoque
type UpdateStockRequest struct {
Stock int `json:"stock" validate:"required,gte=0"`
}

// ProductFiltersRequest representa os filtros de listagem via query params
type ProductFiltersRequest struct {
CategoryID *string
MinPrice   *float64
MaxPrice   *float64
InStock    *bool
}
