package product

import (
	apperrors "meuApp/pkg/errors"
)

// Product Module Errors

var (
	// Erros de validação
	ErrInvalidName        = apperrors.NewValidationError("Product name is required", map[string]string{"field": "name"})
	ErrInvalidDescription = apperrors.NewValidationError("Product description is required", map[string]string{"field": "description"})
	ErrInvalidPrice       = apperrors.NewValidationError("Price must be greater than zero", map[string]string{"field": "price"})
	ErrInvalidStock       = apperrors.NewValidationError("Stock cannot be negative", map[string]string{"field": "stock"})

	// Erros de negócio
	ErrProductNotFound   = apperrors.NewNotFoundError("Product", "")
	ErrInsufficientStock = apperrors.NewDomainError("Insufficient stock available", nil)
	ErrLowStock          = apperrors.NewDomainError("Stock level is low", nil)

	// Erros de infraestrutura
	ErrProductRepository = apperrors.NewInfrastructureError("Product repository error", nil)
)

// NewProductNotFoundError cria um erro de produto não encontrado com ID específico
func NewProductNotFoundError(productID string) *apperrors.AppError {
	return apperrors.NewNotFoundError("Product", productID)
}

// NewInsufficientStockError cria um erro de estoque insuficiente
func NewInsufficientStockError(productID string, available, requested int) *apperrors.AppError {
	return apperrors.NewDomainError("Insufficient stock available", nil)
}
