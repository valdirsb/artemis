# ⚠️ Error Handling

## 📋 Índice
- [Visão Geral](#visão-geral)
- [AppError Structure](#apperror-structure)
- [Tipos de Erros](#tipos-de-erros)
- [Criando Erros](#criando-erros)
- [Error Wrapping](#error-wrapping)
- [HTTP Error Handling](#http-error-handling)
- [gRPC Error Handling](#grpc-error-handling)
- [Logging de Erros](#logging-de-erros)
- [Exemplos Práticos](#exemplos-práticos)

---

## 🎯 Visão Geral

O Artemis Framework usa um **sistema de erros estruturado** baseado em `AppError`, que permite:

- ✅ **Classificação clara** de erros por tipo
- ✅ **HTTP Status Codes** automáticos
- ✅ **Context e detalhes** estruturados
- ✅ **Error wrapping** com contexto adicional
- ✅ **Logging consistente**

```
┌────────────────────────────────────────────────┐
│         Domain Layer Error                     │
│   • Business rule violations                   │
│   • Invariant violations                       │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│      Application Layer Error                   │
│   • Wrap domain errors                         │
│   • Add context                                │
│   • Handle validation                          │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│       Infrastructure Layer Error               │
│   • Convert to HTTP/gRPC                       │
│   • Map status codes                           │
│   • Format response                            │
└────────────────────────────────────────────────┘
```

---

## 🏗️ AppError Structure

### Estrutura Base

**`pkg/errors/errors.go`**

```go
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
    Err        error             `json:"-"` // Erro original (não exposto no JSON)
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
```

---

## 📝 Tipos de Erros

### 1. Domain Errors

Violações de regras de negócio e invariantes de domínio.

```go
// NewDomainError cria um erro de domínio
func NewDomainError(message string, err error) *AppError {
    return &AppError{
        Type:       ErrorTypeDomain,
        Message:    message,
        StatusCode: http.StatusUnprocessableEntity,
        Err:        err,
    }
}
```

**Exemplo:**

```go
// internal/modules/product/domain/product.go
func (p *Product) DecreaseStock(quantity int) error {
    if quantity <= 0 {
        return errors.NewDomainError("quantity must be positive", nil)
    }
    if p.stock < quantity {
        return errors.NewDomainError(
            fmt.Sprintf("insufficient stock: have %d, need %d", p.stock, quantity),
            nil,
        )
    }
    p.stock -= quantity
    return nil
}
```

### 2. Validation Errors

Erros de validação de entrada.

```go
// NewValidationError cria um erro de validação
func NewValidationError(message string, details map[string]string) *AppError {
    return &AppError{
        Type:       ErrorTypeValidation,
        Message:    message,
        Details:    details,
        StatusCode: http.StatusBadRequest,
    }
}
```

**Exemplo:**

```go
// internal/modules/user/errors.go
var (
    ErrInvalidEmail = errors.NewValidationError(
        "Invalid email format",
        map[string]string{"field": "email"},
    )
    
    ErrInvalidPassword = errors.NewValidationError(
        "Password must be at least 6 characters",
        map[string]string{"field": "password"},
    )
)
```

### 3. Not Found Errors

Recurso não encontrado.

```go
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
```

**Exemplo:**

```go
// internal/modules/user/errors.go
func NewUserNotFoundError(userID string) *errors.AppError {
    return errors.NewNotFoundError("User", userID)
}

// Uso
user, err := userRepo.GetByID(ctx, userID)
if err != nil {
    return nil, NewUserNotFoundError(userID)
}
```

### 4. Conflict Errors

Conflitos de estado ou duplicação.

```go
// NewConflictError cria um erro de conflito
func NewConflictError(message string, details map[string]string) *AppError {
    return &AppError{
        Type:       ErrorTypeConflict,
        Message:    message,
        Details:    details,
        StatusCode: http.StatusConflict,
    }
}
```

**Exemplo:**

```go
// internal/modules/user/errors.go
func NewEmailAlreadyExistsError(email string) *errors.AppError {
    return errors.NewConflictError(
        "Email already registered",
        map[string]string{
            "field": "email",
            "value": email,
        },
    )
}
```

### 5. Authorization Errors

```go
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
```

**Exemplo:**

```go
var (
    ErrInvalidCredentials = errors.NewUnauthorizedError("Invalid email or password")
    ErrAccessDenied       = errors.NewForbiddenError("You don't have permission to access this resource")
)
```

### 6. Infrastructure Errors

Erros de infraestrutura (database, cache, external APIs).

```go
// NewInfrastructureError cria um erro de infraestrutura
func NewInfrastructureError(message string, err error) *AppError {
    return &AppError{
        Type:       ErrorTypeInfrastructure,
        Message:    message,
        StatusCode: http.StatusInternalServerError,
        Err:        err,
    }
}
```

**Exemplo:**

```go
// internal/modules/user/repository/user_repository.go
func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
    model := toUserModel(user)
    
    if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
        return errors.NewInfrastructureError(
            "Failed to create user in database",
            err,
        )
    }
    
    return nil
}
```

### 7. Internal Errors

Erros internos não classificados.

```go
// NewInternalError cria um erro interno genérico
func NewInternalError(message string, err error) *AppError {
    return &AppError{
        Type:       ErrorTypeInternal,
        Message:    message,
        StatusCode: http.StatusInternalServerError,
        Err:        err,
    }
}
```

---

## 🎨 Criando Erros

### Erros de Módulo

**`internal/modules/user/errors.go`**

```go
package user

import (
    apperrors "meuApp/pkg/errors"
)

// User Module Errors

var (
    // Validation errors
    ErrInvalidEmail    = apperrors.NewValidationError("Invalid email format", map[string]string{"field": "email"})
    ErrInvalidPassword = apperrors.NewValidationError("Password must be at least 6 characters", map[string]string{"field": "password"})
    ErrInvalidName     = apperrors.NewValidationError("Name is required", map[string]string{"field": "name"})

    // Business errors
    ErrUserNotFound       = apperrors.NewNotFoundError("User", "")
    ErrEmailAlreadyExists = apperrors.NewConflictError("Email already registered", map[string]string{"field": "email"})
    ErrInvalidCredentials = apperrors.NewUnauthorizedError("Invalid email or password")

    // Infrastructure errors
    ErrUserRepository = apperrors.NewInfrastructureError("User repository error", nil)
)

// NewUserNotFoundError cria um erro de usuário não encontrado com ID específico
func NewUserNotFoundError(userID string) *apperrors.AppError {
    return apperrors.NewNotFoundError("User", userID)
}

// NewEmailAlreadyExistsError cria um erro de email duplicado
func NewEmailAlreadyExistsError(email string) *apperrors.AppError {
    return apperrors.NewConflictError("Email already registered", map[string]string{
        "field": "email",
        "value": email,
    })
}
```

### Erros de Domínio

**`internal/modules/product/domain/product.go`**

```go
package domain

import (
    "errors"
    "fmt"
)

var (
    ErrInsufficientStock = errors.New("insufficient stock")
    ErrInvalidQuantity   = errors.New("invalid quantity")
    ErrInvalidPrice      = errors.New("invalid price")
)

func (p *Product) DecreaseStock(quantity int) error {
    if quantity <= 0 {
        return fmt.Errorf("%w: quantity must be positive", ErrInvalidQuantity)
    }
    
    if p.stock < quantity {
        return fmt.Errorf("%w: have %d, need %d", ErrInsufficientStock, p.stock, quantity)
    }
    
    p.stock -= quantity
    p.updatedAt = time.Now()
    return nil
}
```

---

## 🔄 Error Wrapping

### WrapError

Adiciona contexto a erros existentes.

```go
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
```

**Uso:**

```go
// internal/modules/user/application/commands/create_user.go
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
    // Criar usuário
    user, err := domain.NewUser(cmd.Username, cmd.Email)
    if err != nil {
        return nil, errors.WrapError(err, "failed to create user entity")
    }

    // Salvar no repositório
    if err := h.userRepo.Create(ctx, user); err != nil {
        return nil, errors.WrapError(err, "failed to save user")
    }

    return user, nil
}
```

### errors.Is e errors.As

Go 1.13+ error wrapping.

```go
import (
    "errors"
    "fmt"
)

// Verificar tipo de erro
if errors.Is(err, domain.ErrInsufficientStock) {
    // Tratar erro de estoque insuficiente
    return errors.NewDomainError("Cannot process order: insufficient stock", err)
}

// Extrair AppError
var appErr *errors.AppError
if errors.As(err, &appErr) {
    // Trabalhar com AppError
    log.Printf("Error type: %s, Message: %s", appErr.Type, appErr.Message)
}
```

---

## 🌐 HTTP Error Handling

### Middleware de Erro

**`pkg/adapters/http/middleware/error_handler.go`**

```go
package middleware

import (
    "encoding/json"
    "log"
    "net/http"

    apperrors "meuApp/pkg/errors"
)

// ErrorResponse representa a estrutura de resposta de erro HTTP
type ErrorResponse struct {
    Error   string            `json:"error"`
    Type    string            `json:"type,omitempty"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"`
}

// ErrorHandler é um middleware que captura panics e trata erros
func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                respondWithError(w, apperrors.NewInternalError("Unexpected error occurred", nil))
            }
        }()

        next.ServeHTTP(w, r)
    })
}

// RespondWithAppError envia uma resposta de erro baseada em AppError
func RespondWithAppError(w http.ResponseWriter, err error) {
    if err == nil {
        return
    }

    // Tenta converter para AppError
    if appErr, ok := apperrors.AsAppError(err); ok {
        respondWithError(w, appErr)
        return
    }

    // Erro genérico
    respondWithError(w, apperrors.NewInternalError("An unexpected error occurred", err))
}

// respondWithError envia a resposta de erro HTTP
func respondWithError(w http.ResponseWriter, appErr *apperrors.AppError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(appErr.HTTPStatusCode())

    response := ErrorResponse{
        Error:   string(appErr.Type),
        Type:    string(appErr.Type),
        Message: appErr.Message,
        Details: appErr.Details,
    }

    // Log do erro (não expor detalhes internos)
    if appErr.Err != nil {
        log.Printf("Error [%s] %s: %v", appErr.Type, appErr.Message, appErr.Err)
    } else {
        log.Printf("Error [%s] %s", appErr.Type, appErr.Message)
    }

    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Failed to encode error response: %v", err)
    }
}
```

### Uso em Handlers

**`internal/modules/user/adapters/http/user_http_handler.go`**

```go
func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        middleware.RespondWithAppError(c.Writer, 
            errors.NewValidationError("Invalid request body", nil))
        return
    }

    user, err := h.userService.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password)
    if err != nil {
        middleware.RespondWithAppError(c.Writer, err)
        return
    }

    response := dto.ToUserResponse(user)
    middleware.RespondWithJSON(c.Writer, http.StatusCreated, response)
}
```

### Respostas de Erro

**Validation Error (400):**

```json
{
  "error": "VALIDATION_ERROR",
  "type": "VALIDATION_ERROR",
  "message": "Invalid email format",
  "details": {
    "field": "email"
  }
}
```

**Not Found (404):**

```json
{
  "error": "NOT_FOUND",
  "type": "NOT_FOUND",
  "message": "User not found",
  "details": {
    "resource": "User",
    "identifier": "user-123"
  }
}
```

**Conflict (409):**

```json
{
  "error": "CONFLICT",
  "type": "CONFLICT",
  "message": "Email already registered",
  "details": {
    "field": "email",
    "value": "user@example.com"
  }
}
```

---

## 🔌 gRPC Error Handling

### Conversão de Erros para gRPC

**`pkg/adapters/grpc/error_converter.go`**

```go
package grpc

import (
    "errors"

    apperrors "meuApp/pkg/errors"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// ToGRPCError converte AppError para gRPC status
func ToGRPCError(err error) error {
    if err == nil {
        return nil
    }

    // Extrair AppError
    var appErr *apperrors.AppError
    if !errors.As(err, &appErr) {
        return status.Error(codes.Internal, "internal server error")
    }

    // Mapear ErrorType para gRPC codes
    code := mapErrorTypeToGRPCCode(appErr.Type)
    return status.Error(code, appErr.Message)
}

func mapErrorTypeToGRPCCode(errType apperrors.ErrorType) codes.Code {
    switch errType {
    case apperrors.ErrorTypeValidation:
        return codes.InvalidArgument
    case apperrors.ErrorTypeNotFound:
        return codes.NotFound
    case apperrors.ErrorTypeConflict:
        return codes.AlreadyExists
    case apperrors.ErrorTypeUnauthorized:
        return codes.Unauthenticated
    case apperrors.ErrorTypeForbidden:
        return codes.PermissionDenied
    case apperrors.ErrorTypeDomain:
        return codes.FailedPrecondition
    case apperrors.ErrorTypeInfrastructure:
        return codes.Unavailable
    default:
        return codes.Internal
    }
}
```

### Uso em gRPC Handlers

**`internal/modules/product/adapters/grpc/product_grpc_handler.go`**

```go
package grpc

import (
    "context"

    "meuApp/internal/modules/product/application/commands"
    pb "meuApp/proto"
    grpcerrors "meuApp/pkg/adapters/grpc"
    "meuApp/pkg/contracts"
)

type ProductGRPCHandler struct {
    pb.UnimplementedProductServiceServer
    commandBus contracts.CommandBus
}

func (h *ProductGRPCHandler) CreateProduct(
    ctx context.Context,
    req *pb.CreateProductRequest,
) (*pb.ProductResponse, error) {
    cmd := commands.CreateProductCommand{
        Name:        req.Name,
        Description: req.Description,
        CategoryID:  req.CategoryId,
        Price:       req.Price,
        Stock:       int(req.Stock),
    }

    result, err := h.commandBus.Execute(ctx, cmd)
    if err != nil {
        return nil, grpcerrors.ToGRPCError(err)
    }

    product := result.(*domain.Product)
    return toProductResponse(product), nil
}
```

---

## 📊 Logging de Erros

### Structured Logging

**`pkg/adapters/logger/logger.go`**

```go
package logger

import (
    "log"

    apperrors "meuApp/pkg/errors"
)

type Logger struct {
    // Implementação do logger
}

// LogError registra um erro com contexto
func (l *Logger) LogError(err error, context map[string]interface{}) {
    if err == nil {
        return
    }

    // Se é AppError, extrair mais informações
    if appErr, ok := apperrors.AsAppError(err); ok {
        log.Printf("[%s] %s | Details: %v | Context: %v | Cause: %v",
            appErr.Type,
            appErr.Message,
            appErr.Details,
            context,
            appErr.Err,
        )
        return
    }

    // Erro genérico
    log.Printf("[ERROR] %v | Context: %v", err, context)
}
```

### Uso em Commands

```go
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
    user, err := h.userRepo.Create(ctx, newUser)
    if err != nil {
        h.logger.LogError(err, map[string]interface{}{
            "command":  "CreateUser",
            "email":    cmd.Email,
            "username": cmd.Username,
        })
        return nil, errors.WrapError(err, "failed to create user")
    }

    return user, nil
}
```

---

## 📚 Exemplos Práticos

### Exemplo 1: Command Handler Completo

**`internal/modules/order/application/commands/create_order.go`**

```go
package commands

import (
    "context"
    "fmt"

    "meuApp/internal/modules/order/domain"
    "meuApp/internal/modules/order/ports"
    "meuApp/pkg/contracts"
    "meuApp/pkg/errors"
)

type CreateOrderCommand struct {
    UserID string
    Items  []OrderItemCommand
}

type CreateOrderHandler struct {
    orderRepo   ports.OrderRepository
    productRepo ports.ProductRepository
    logger      contracts.Logger
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*domain.Order, error) {
    // 1. Validar entrada
    if err := h.validateCommand(cmd); err != nil {
        return nil, err
    }

    // 2. Buscar produtos e validar estoque
    products, err := h.validateAndFetchProducts(ctx, cmd.Items)
    if err != nil {
        return nil, err
    }

    // 3. Criar entidade de domínio
    order, err := domain.NewOrder(cmd.UserID, products)
    if err != nil {
        h.logger.LogError(err, map[string]interface{}{
            "command": "CreateOrder",
            "user_id": cmd.UserID,
        })
        return nil, errors.WrapError(err, "failed to create order entity")
    }

    // 4. Salvar ordem
    if err := h.orderRepo.Create(ctx, order); err != nil {
        h.logger.LogError(err, map[string]interface{}{
            "command":  "CreateOrder",
            "order_id": order.ID(),
        })
        return nil, errors.WrapError(err, "failed to save order")
    }

    return order, nil
}

func (h *CreateOrderHandler) validateCommand(cmd CreateOrderCommand) error {
    if cmd.UserID == "" {
        return errors.NewValidationError("user_id is required", map[string]string{
            "field": "user_id",
        })
    }

    if len(cmd.Items) == 0 {
        return errors.NewValidationError("at least one item is required", map[string]string{
            "field": "items",
        })
    }

    return nil
}

func (h *CreateOrderHandler) validateAndFetchProducts(
    ctx context.Context,
    items []OrderItemCommand,
) ([]domain.OrderItem, error) {
    var orderItems []domain.OrderItem

    for i, item := range items {
        // Buscar produto
        product, err := h.productRepo.GetByID(ctx, item.ProductID)
        if err != nil {
            return nil, errors.NewNotFoundError("Product", item.ProductID)
        }

        // Validar estoque
        if product.Stock() < item.Quantity {
            return nil, errors.NewDomainError(
                fmt.Sprintf("insufficient stock for product %s", item.ProductID),
                nil,
            )
        }

        orderItems = append(orderItems, domain.OrderItem{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     product.Price(),
        })
    }

    return orderItems, nil
}
```

### Exemplo 2: Repository com Error Handling

**`internal/modules/user/repository/user_repository.go`**

```go
package repository

import (
    "context"
    "errors"

    "meuApp/internal/modules/user/domain"
    apperrors "meuApp/pkg/errors"

    "gorm.io/gorm"
)

type mysqlUserRepository struct {
    db *gorm.DB
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
    var model UserModel

    err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperrors.NewNotFoundError("User", id)
        }
        return nil, apperrors.NewInfrastructureError("failed to get user from database", err)
    }

    return model.ToDomain(), nil
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
    model := toUserModel(user)

    if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
        // Verificar constraint de email único
        if isDuplicateKeyError(err) {
            return apperrors.NewConflictError(
                "email already registered",
                map[string]string{"field": "email"},
            )
        }
        return apperrors.NewInfrastructureError("failed to create user in database", err)
    }

    return nil
}

func isDuplicateKeyError(err error) bool {
    // MySQL error 1062 = duplicate key
    // Implementar verificação específica do driver
    return strings.Contains(err.Error(), "Duplicate entry")
}
```

### Exemplo 3: Error Recovery

**`pkg/adapters/http/middleware/recovery.go`**

```go
package middleware

import (
    "fmt"
    "log"
    "net/http"
    "runtime/debug"

    apperrors "meuApp/pkg/errors"

    "github.com/gin-gonic/gin"
)

// Recovery middleware recupera de panics
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log stack trace
                log.Printf("Panic recovered: %v\n%s", err, debug.Stack())

                // Responder com erro interno
                appErr := apperrors.NewInternalError(
                    "An unexpected error occurred",
                    fmt.Errorf("%v", err),
                )

                RespondWithAppError(c.Writer, appErr)
                c.Abort()
            }
        }()

        c.Next()
    }
}
```

---

## 🎯 Best Practices

### ✅ DO

1. **Use tipos de erro específicos**
   ```go
   return errors.NewNotFoundError("User", userID)
   ```

2. **Adicione contexto aos erros**
   ```go
   return errors.WrapError(err, "failed to create user")
   ```

3. **Log erros com contexto**
   ```go
   logger.LogError(err, map[string]interface{}{
       "command": "CreateUser",
       "email": cmd.Email,
   })
   ```

4. **Use errors.Is para verificar tipos**
   ```go
   if errors.Is(err, domain.ErrInsufficientStock) {
       // Handle specifically
   }
   ```

### ❌ DON'T

1. **Não exponha detalhes internos**
   ```go
   // ❌ Ruim
   return fmt.Errorf("database error: %v", dbErr)
   
   // ✅ Bom
   return errors.NewInfrastructureError("failed to save user", dbErr)
   ```

2. **Não ignore erros**
   ```go
   // ❌ Ruim
   user, _ := userRepo.GetByID(ctx, id)
   
   // ✅ Bom
   user, err := userRepo.GetByID(ctx, id)
   if err != nil {
       return nil, err
   }
   ```

3. **Não retorne nil sem verificar**
   ```go
   // ❌ Ruim
   if err != nil {
       return nil, nil  // Pode causar nil pointer
   }
   
   // ✅ Bom
   if err != nil {
       return nil, errors.WrapError(err, "operation failed")
   }
   ```

---

## 📚 Próximos Passos

- **[Testing](21-testing-strategy.md)** - Testar error handling
- **[Logging](20-logging.md)** - Logging estruturado
- **[Best Practices](23-best-practices.md)** - Boas práticas gerais

---

**[⬅️ Validation](18-validation.md)** | **[Índice](README.md)** | **[Testing ➡️](21-testing-strategy.md)**
