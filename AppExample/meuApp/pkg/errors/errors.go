package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType representa o tipo de erro
type ErrorType string

const (
	ErrorTypeDomain         ErrorType = "DOMAIN_ERROR"
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeConflict       ErrorType = "CONFLICT"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
	ErrorTypeInfrastructure ErrorType = "INFRASTRUCTURE_ERROR"
	ErrorTypeInternal       ErrorType = "INTERNAL_ERROR"
)

// AppError é a estrutura base para todos os erros da aplicação
type AppError struct {
	Type       ErrorType         `json:"type"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
	StatusCode int               `json:"-"`
	Err        error             `json:"-"`
}

// Error implementa a interface error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap retorna o erro original (para errors.Is e errors.As)
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatusCode retorna o código HTTP apropriado
func (e *AppError) HTTPStatusCode() int {
	if e.StatusCode != 0 {
		return e.StatusCode
	}
	// Status code padrão baseado no tipo
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeInfrastructure, ErrorTypeInternal:
		return http.StatusInternalServerError
	case ErrorTypeDomain:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// NewDomainError cria um erro de domínio
func NewDomainError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeDomain,
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
		Err:        err,
	}
}

// NewValidationError cria um erro de validação
func NewValidationError(message string, details map[string]string) *AppError {
	return &AppError{
		Type:       ErrorTypeValidation,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError cria um erro de recurso não encontrado
func NewNotFoundError(resource string, identifier string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Details: map[string]string{
			"resource":   resource,
			"identifier": identifier,
		},
		StatusCode: http.StatusNotFound,
	}
}

// NewConflictError cria um erro de conflito
func NewConflictError(message string, details map[string]string) *AppError {
	return &AppError{
		Type:       ErrorTypeConflict,
		Message:    message,
		Details:    details,
		StatusCode: http.StatusConflict,
	}
}

// NewUnauthorizedError cria um erro de não autorizado
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError cria um erro de acesso negado
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:       ErrorTypeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewInfrastructureError cria um erro de infraestrutura
func NewInfrastructureError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeInfrastructure,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewInternalError cria um erro interno genérico
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// IsAppError verifica se um erro é do tipo AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError converte um erro para AppError (se possível)
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}

// WrapError envolve um erro genérico em AppError
func WrapError(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	// Se já é um AppError, adiciona contexto
	if appErr, ok := AsAppError(err); ok {
		return &AppError{
			Type:       appErr.Type,
			Message:    fmt.Sprintf("%s: %s", message, appErr.Message),
			Details:    appErr.Details,
			StatusCode: appErr.StatusCode,
			Err:        appErr.Err,
		}
	}

	// Cria um novo AppError
	return NewInternalError(message, err)
}

// Erros comuns pré-definidos
var (
	ErrInvalidInput       = NewValidationError("Invalid input", nil)
	ErrUnauthorized       = NewUnauthorizedError("Unauthorized access")
	ErrForbidden          = NewForbiddenError("Access forbidden")
	ErrInternalServer     = NewInternalError("Internal server error", nil)
	ErrDatabaseConnection = NewInfrastructureError("Database connection error", nil)
)
