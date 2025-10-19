package order

import (
	apperrors "meuApp/pkg/errors"
)

// Order Module Errors

var (
	// Erros de validação
	ErrInvalidUserID      = apperrors.NewValidationError("User ID is required", map[string]string{"field": "user_id"})
	ErrEmptyOrderItems    = apperrors.NewValidationError("Order must have at least one item", map[string]string{"field": "items"})
	ErrInvalidQuantity    = apperrors.NewValidationError("Quantity must be greater than zero", map[string]string{"field": "quantity"})
	ErrInvalidOrderStatus = apperrors.NewValidationError("Invalid order status", map[string]string{"field": "status"})

	// Erros de negócio
	ErrOrderNotFound           = apperrors.NewNotFoundError("Order", "")
	ErrUserNotFound            = apperrors.NewNotFoundError("User", "")
	ErrProductNotFound         = apperrors.NewNotFoundError("Product", "")
	ErrInsufficientStock       = apperrors.NewDomainError("Insufficient stock for product", nil)
	ErrInvalidStatusTransition = apperrors.NewDomainError("Invalid order status transition", nil)
	ErrCannotCancelOrder       = apperrors.NewDomainError("Order cannot be cancelled in current status", nil)

	// Erros de infraestrutura
	ErrOrderRepository = apperrors.NewInfrastructureError("Order repository error", nil)
)

// NewOrderNotFoundError cria um erro de pedido não encontrado com ID específico
func NewOrderNotFoundError(orderID string) *apperrors.AppError {
	return apperrors.NewNotFoundError("Order", orderID)
}

// NewUserNotFoundError cria um erro de usuário não encontrado
func NewUserNotFoundError(userID string) *apperrors.AppError {
	return apperrors.NewNotFoundError("User", userID)
}

// NewProductNotFoundError cria um erro de produto não encontrado
func NewProductNotFoundError(productID string) *apperrors.AppError {
	return apperrors.NewNotFoundError("Product", productID)
}

// NewInsufficientStockError cria um erro de estoque insuficiente para um produto
func NewInsufficientStockError(productID string, available, requested int) *apperrors.AppError {
	return apperrors.NewDomainError("Insufficient stock for product", nil)
}

// NewInvalidStatusTransitionError cria um erro de transição inválida de status
func NewInvalidStatusTransitionError(currentStatus, newStatus string) *apperrors.AppError {
	return apperrors.NewDomainError("Invalid order status transition", nil)
}
