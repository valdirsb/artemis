# ⚡ Implementando Commands - Write Operations

## 📋 Índice
- [O que são Commands?](#o-que-são-commands)
- [Anatomia de um Command](#anatomia-de-um-command)
- [Command Handlers](#command-handlers)
- [Validação](#validação)
- [Transações](#transações)
- [Error Handling](#error-handling)
- [Exemplos Práticos](#exemplos-práticos)
- [Boas Práticas](#boas-práticas)

---

## 🎯 O que são Commands?

**Commands** são objetos que representam **intenções de modificar o estado** do sistema (write operations: CREATE, UPDATE, DELETE).

### Command vs Query

```
┌─────────────────────────────────────────────────┐
│                    CQRS                         │
├─────────────────────┬───────────────────────────┤
│     COMMANDS        │        QUERIES            │
│   (Write Side)      │      (Read Side)          │
├─────────────────────┼───────────────────────────┤
│ • Modificam estado  │ • Apenas leem             │
│ • Retornam void/ID  │ • Retornam dados          │
│ • Validação pesada  │ • Validação leve          │
│ • Transações        │ • Sem transações          │
│ • Eventos           │ • Sem eventos             │
│ • Async possível    │ • Sempre sync             │
└─────────────────────┴───────────────────────────┘
```

### Fluxo de um Command

```
HTTP Request (POST/PUT/DELETE)
        │
        ▼
┌───────────────┐
│  HTTP Adapter │ Valida input
└───────┬───────┘
        │
        ▼
┌───────────────┐
│    Command    │ DTO com dados
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ Command Bus   │ Despacha comando
└───────┬───────┘
        │
        ▼
┌───────────────┐
│Command Handler│ Lógica de negócio
│               │ 1. Valida
│               │ 2. Busca entidade
│               │ 3. Aplica mudança
│               │ 4. Persiste
│               │ 5. Publica evento
└───────┬───────┘
        │
        ▼
┌───────────────┐
│   Repository  │ Persiste no DB
└───────┬───────┘
        │
        ▼
┌───────────────┐
│   Event Bus   │ Notifica sistema
└───────────────┘
```

---

## 🏗️ Anatomia de um Command

### Estrutura Básica

```go
package commands

// Command é um DTO imutável
type CreateUserCommand struct {
    // Dados necessários para a operação
    Name     string
    Email    string
    Password string
    
    // Metadados
    CreatedBy string
    IPAddress string
    UserAgent string
}

// CommandName identifica o comando
func (c *CreateUserCommand) CommandName() string {
    return "CreateUserCommand"
}
```

### Princípios de um Command

1. **Imutável** - Dados não mudam após criação
2. **Específico** - Um comando = uma intenção
3. **Validável** - Pode ser validado antes de executar
4. **Rastreável** - Contém metadados de auditoria

### Exemplo Completo

```go
package commands

import (
    "errors"
    "time"
)

type CreateProductCommand struct {
    // Dados do produto
    Name        string
    Description string
    Price       float64
    SKU         string
    CategoryID  string
    Stock       int
    
    // Metadados
    CreatedBy string
    Timestamp time.Time
}

func NewCreateProductCommand(
    name, description, sku, categoryID string,
    price float64,
    stock int,
    createdBy string,
) *CreateProductCommand {
    return &CreateProductCommand{
        Name:        name,
        Description: description,
        Price:       price,
        SKU:         sku,
        CategoryID:  categoryID,
        Stock:       stock,
        CreatedBy:   createdBy,
        Timestamp:   time.Now(),
    }
}

func (c *CreateProductCommand) CommandName() string {
    return "CreateProductCommand"
}

// Validate valida dados do comando
func (c *CreateProductCommand) Validate() error {
    if c.Name == "" {
        return errors.New("product name is required")
    }
    if len(c.Name) < 3 {
        return errors.New("product name must be at least 3 characters")
    }
    if c.Price <= 0 {
        return errors.New("product price must be greater than zero")
    }
    if c.SKU == "" {
        return errors.New("product SKU is required")
    }
    if c.Stock < 0 {
        return errors.New("product stock cannot be negative")
    }
    return nil
}
```

---

## 🔧 Command Handlers

### Interface do Handler

```go
package contracts

import "context"

type CommandHandler interface {
    Handle(ctx context.Context, cmd Command) (interface{}, error)
}
```

### Handler Básico

```go
package commands

import (
    "context"
    "meuApp/internal/modules/product/domain/entities"
    "meuApp/internal/modules/product/repository"
    "meuApp/pkg/contracts"
)

type CreateProductHandler struct {
    repo repository.ProductRepository
}

func NewCreateProductHandler(repo repository.ProductRepository) *CreateProductHandler {
    return &CreateProductHandler{repo: repo}
}

func (h *CreateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    // 1. Type assertion
    createCmd := cmd.(*CreateProductCommand)
    
    // 2. Validar comando
    if err := createCmd.Validate(); err != nil {
        return nil, err
    }
    
    // 3. Criar entidade de domínio
    product, err := entities.NewProduct(
        createCmd.Name,
        createCmd.Description,
        createCmd.Price,
        createCmd.SKU,
    )
    if err != nil {
        return nil, err
    }
    
    // 4. Persistir
    if err := h.repo.Save(ctx, product); err != nil {
        return nil, err
    }
    
    // 5. Retornar ID
    return product.ID, nil
}
```

### Handler com Eventos

```go
type CreateProductHandler struct {
    repo     repository.ProductRepository
    eventBus contracts.EventBus
}

func NewCreateProductHandler(
    repo repository.ProductRepository,
    eventBus contracts.EventBus,
) *CreateProductHandler {
    return &CreateProductHandler{
        repo:     repo,
        eventBus: eventBus,
    }
}

func (h *CreateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateProductCommand)
    
    // Validação
    if err := createCmd.Validate(); err != nil {
        return nil, err
    }
    
    // Criar produto
    product, err := entities.NewProduct(
        createCmd.Name,
        createCmd.Description,
        createCmd.Price,
        createCmd.SKU,
    )
    if err != nil {
        return nil, err
    }
    
    // Persistir
    if err := h.repo.Save(ctx, product); err != nil {
        return nil, err
    }
    
    // Publicar evento
    event := events.NewProductCreatedEvent(
        product.ID,
        product.Name,
        product.Price,
        createCmd.CreatedBy,
    )
    h.eventBus.Publish(ctx, event)
    
    return product.ID, nil
}
```

### Handler com Múltiplas Dependências

```go
type CreateOrderHandler struct {
    orderRepo    repository.OrderRepository
    productRepo  repository.ProductRepository
    userRepo     repository.UserRepository
    eventBus     contracts.EventBus
    logger       contracts.Logger
}

func NewCreateOrderHandler(
    orderRepo repository.OrderRepository,
    productRepo repository.ProductRepository,
    userRepo repository.UserRepository,
    eventBus contracts.EventBus,
    logger contracts.Logger,
) *CreateOrderHandler {
    return &CreateOrderHandler{
        orderRepo:   orderRepo,
        productRepo: productRepo,
        userRepo:    userRepo,
        eventBus:    eventBus,
        logger:      logger,
    }
}

func (h *CreateOrderHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateOrderCommand)
    
    h.logger.Info("Creating order", map[string]interface{}{
        "user_id": createCmd.UserID,
        "items":   len(createCmd.Items),
    })
    
    // 1. Validar usuário existe
    user, err := h.userRepo.FindByID(ctx, createCmd.UserID)
    if err != nil {
        return nil, errors.New("user not found")
    }
    
    // 2. Validar produtos e calcular total
    var total float64
    for _, item := range createCmd.Items {
        product, err := h.productRepo.FindByID(ctx, item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product %s not found", item.ProductID)
        }
        
        if product.Stock < item.Quantity {
            return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
        }
        
        total += product.Price * float64(item.Quantity)
    }
    
    // 3. Criar pedido
    order, err := entities.NewOrder(user.ID, createCmd.Items)
    if err != nil {
        return nil, err
    }
    
    // 4. Persistir
    if err := h.orderRepo.Save(ctx, order); err != nil {
        return nil, err
    }
    
    // 5. Publicar evento
    event := events.NewOrderCreatedEvent(order.ID, user.ID, total)
    h.eventBus.Publish(ctx, event)
    
    h.logger.Info("Order created", map[string]interface{}{
        "order_id": order.ID,
        "total":    total,
    })
    
    return order.ID, nil
}
```

---

## ✅ Validação

### Validação no Command

```go
type CreateProductCommand struct {
    Name  string
    Price float64
    SKU   string
}

func (c *CreateProductCommand) Validate() error {
    var errs []string
    
    if c.Name == "" {
        errs = append(errs, "name is required")
    }
    
    if c.Price <= 0 {
        errs = append(errs, "price must be positive")
    }
    
    if c.SKU == "" {
        errs = append(errs, "SKU is required")
    }
    
    if len(errs) > 0 {
        return &ValidationError{Errors: errs}
    }
    
    return nil
}

// ValidationError agrupa múltiplos erros
type ValidationError struct {
    Errors []string
}

func (e *ValidationError) Error() string {
    return strings.Join(e.Errors, "; ")
}
```

### Validação com go-playground/validator

```go
import "github.com/go-playground/validator/v10"

type CreateProductCommand struct {
    Name  string  `validate:"required,min=3,max=100"`
    Price float64 `validate:"required,gt=0"`
    SKU   string  `validate:"required,alphanum"`
    Email string  `validate:"required,email"`
}

var validate = validator.New()

func (c *CreateProductCommand) Validate() error {
    return validate.Struct(c)
}
```

### Validação de Regras de Negócio

```go
func (h *CreateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateProductCommand)
    
    // 1. Validação de estrutura
    if err := createCmd.Validate(); err != nil {
        return nil, err
    }
    
    // 2. Validação de regras de negócio
    
    // Verificar se SKU já existe
    exists, err := h.repo.ExistsBySKU(ctx, createCmd.SKU)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("product with this SKU already exists")
    }
    
    // Verificar se categoria existe
    category, err := h.categoryRepo.FindByID(ctx, createCmd.CategoryID)
    if err != nil {
        return nil, errors.New("category not found")
    }
    if !category.Active {
        return nil, errors.New("category is not active")
    }
    
    // Continuar com criação...
}
```

---

## 🔄 Transações

### Transação Simples

```go
func (h *CreateOrderHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateOrderCommand)
    
    // Iniciar transação
    tx, err := h.db.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback() // Rollback se não commitado
    
    // Criar pedido
    order := &Order{...}
    if err := h.orderRepo.SaveWithTx(ctx, tx, order); err != nil {
        return nil, err
    }
    
    // Atualizar estoque
    for _, item := range createCmd.Items {
        if err := h.productRepo.DecrementStockWithTx(ctx, tx, item.ProductID, item.Quantity); err != nil {
            return nil, err
        }
    }
    
    // Commit
    if err := tx.Commit(); err != nil {
        return nil, err
    }
    
    return order.ID, nil
}
```

### Transação com GORM

```go
func (h *CreateOrderHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateOrderCommand)
    
    var orderID string
    
    // Transaction do GORM
    err := h.db.Transaction(func(tx *gorm.DB) error {
        // 1. Criar pedido
        order := &Order{...}
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        orderID = order.ID
        
        // 2. Criar itens do pedido
        for _, item := range createCmd.Items {
            orderItem := &OrderItem{
                OrderID:   order.ID,
                ProductID: item.ProductID,
                Quantity:  item.Quantity,
            }
            if err := tx.Create(orderItem).Error; err != nil {
                return err
            }
        }
        
        // 3. Atualizar estoque
        for _, item := range createCmd.Items {
            if err := tx.Model(&Product{}).
                Where("id = ?", item.ProductID).
                Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
                return err
            }
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return orderID, nil
}
```

### Unit of Work Pattern

```go
type UnitOfWork struct {
    db *gorm.DB
    tx *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
    return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) Begin() error {
    uow.tx = uow.db.Begin()
    return uow.tx.Error
}

func (uow *UnitOfWork) Commit() error {
    return uow.tx.Commit().Error
}

func (uow *UnitOfWork) Rollback() error {
    return uow.tx.Rollback().Error
}

// Repositórios transacionais
func (uow *UnitOfWork) OrderRepo() repository.OrderRepository {
    return repository.NewMySQLOrderRepository(uow.tx)
}

func (uow *UnitOfWork) ProductRepo() repository.ProductRepository {
    return repository.NewMySQLProductRepository(uow.tx)
}

// Uso no Handler
func (h *CreateOrderHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateOrderCommand)
    
    uow := NewUnitOfWork(h.db)
    if err := uow.Begin(); err != nil {
        return nil, err
    }
    defer uow.Rollback()
    
    // Usar repositórios transacionais
    order := &Order{...}
    if err := uow.OrderRepo().Save(ctx, order); err != nil {
        return nil, err
    }
    
    for _, item := range createCmd.Items {
        if err := uow.ProductRepo().DecrementStock(ctx, item.ProductID, item.Quantity); err != nil {
            return nil, err
        }
    }
    
    if err := uow.Commit(); err != nil {
        return nil, err
    }
    
    return order.ID, nil
}
```

---

## 🚨 Error Handling

### Tipos de Erros

```go
package errors

import "errors"

// Domain Errors
var (
    ErrNotFound          = errors.New("resource not found")
    ErrAlreadyExists     = errors.New("resource already exists")
    ErrInvalidInput      = errors.New("invalid input")
    ErrUnauthorized      = errors.New("unauthorized")
    ErrForbidden         = errors.New("forbidden")
    ErrInsufficientStock = errors.New("insufficient stock")
)

// DomainError encapsula erro com contexto
type DomainError struct {
    Code    string
    Message string
    Cause   error
}

func (e *DomainError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (cause: %v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewDomainError(code, message string, cause error) *DomainError {
    return &DomainError{
        Code:    code,
        Message: message,
        Cause:   cause,
    }
}
```

### Tratamento no Handler

```go
func (h *CreateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateProductCommand)
    
    // Validação
    if err := createCmd.Validate(); err != nil {
        return nil, NewDomainError("VALIDATION_ERROR", "Invalid product data", err)
    }
    
    // Verificar duplicata
    exists, err := h.repo.ExistsBySKU(ctx, createCmd.SKU)
    if err != nil {
        return nil, NewDomainError("DATABASE_ERROR", "Failed to check SKU", err)
    }
    if exists {
        return nil, NewDomainError("DUPLICATE_SKU", "Product with this SKU already exists", nil)
    }
    
    // Criar produto
    product, err := entities.NewProduct(createCmd.Name, createCmd.Price, createCmd.SKU)
    if err != nil {
        return nil, NewDomainError("DOMAIN_ERROR", "Failed to create product entity", err)
    }
    
    // Persistir
    if err := h.repo.Save(ctx, product); err != nil {
        return nil, NewDomainError("DATABASE_ERROR", "Failed to save product", err)
    }
    
    return product.ID, nil
}
```

### Tratamento no HTTP Adapter

```go
func (h *ProductHTTPHandler) CreateProduct(c *gin.Context) {
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    cmd := &CreateProductCommand{...}
    result, err := h.commandBus.Execute(c.Request.Context(), cmd)
    
    if err != nil {
        // Tratar erros de domínio
        var domainErr *errors.DomainError
        if errors.As(err, &domainErr) {
            switch domainErr.Code {
            case "VALIDATION_ERROR":
                c.JSON(http.StatusBadRequest, gin.H{"error": domainErr.Message})
            case "DUPLICATE_SKU":
                c.JSON(http.StatusConflict, gin.H{"error": domainErr.Message})
            case "DOMAIN_ERROR":
                c.JSON(http.StatusUnprocessableEntity, gin.H{"error": domainErr.Message})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
            }
            return
        }
        
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"id": result})
}
```

---

## 📝 Exemplos Práticos

### Exemplo 1: Update Command

```go
// Command
type UpdateProductCommand struct {
    ProductID   string
    Name        string
    Description string
    Price       float64
    UpdatedBy   string
}

func (c *UpdateProductCommand) CommandName() string {
    return "UpdateProductCommand"
}

// Handler
type UpdateProductHandler struct {
    repo     repository.ProductRepository
    eventBus contracts.EventBus
}

func (h *UpdateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    updateCmd := cmd.(*UpdateProductCommand)
    
    // 1. Buscar produto
    product, err := h.repo.FindByID(ctx, updateCmd.ProductID)
    if err != nil {
        return nil, errors.ErrNotFound
    }
    
    // 2. Atualizar dados
    if err := product.Update(
        updateCmd.Name,
        updateCmd.Description,
        updateCmd.Price,
    ); err != nil {
        return nil, err
    }
    
    // 3. Persistir
    if err := h.repo.Update(ctx, product); err != nil {
        return nil, err
    }
    
    // 4. Publicar evento
    event := events.NewProductUpdatedEvent(
        product.ID,
        product.Name,
        product.Price,
        updateCmd.UpdatedBy,
    )
    h.eventBus.Publish(ctx, event)
    
    return nil, nil
}
```

### Exemplo 2: Delete Command (Soft Delete)

```go
// Command
type DeleteProductCommand struct {
    ProductID string
    DeletedBy string
    Reason    string
}

func (c *DeleteProductCommand) CommandName() string {
    return "DeleteProductCommand"
}

// Handler
type DeleteProductHandler struct {
    repo     repository.ProductRepository
    eventBus contracts.EventBus
}

func (h *DeleteProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    deleteCmd := cmd.(*DeleteProductCommand)
    
    // 1. Buscar produto
    product, err := h.repo.FindByID(ctx, deleteCmd.ProductID)
    if err != nil {
        return nil, errors.ErrNotFound
    }
    
    // 2. Verificar se pode deletar
    if product.HasActiveOrders() {
        return nil, errors.New("cannot delete product with active orders")
    }
    
    // 3. Soft delete
    product.MarkAsDeleted(deleteCmd.DeletedBy, deleteCmd.Reason)
    
    // 4. Persistir
    if err := h.repo.Update(ctx, product); err != nil {
        return nil, err
    }
    
    // 5. Publicar evento
    event := events.NewProductDeletedEvent(
        product.ID,
        deleteCmd.DeletedBy,
        deleteCmd.Reason,
    )
    h.eventBus.Publish(ctx, event)
    
    return nil, nil
}
```

### Exemplo 3: Batch Command

```go
// Command para operação em lote
type UpdateProductPricesCommand struct {
    Updates   []PriceUpdate
    UpdatedBy string
}

type PriceUpdate struct {
    ProductID string
    NewPrice  float64
}

func (c *UpdateProductPricesCommand) CommandName() string {
    return "UpdateProductPricesCommand"
}

// Handler
type UpdateProductPricesHandler struct {
    repo     repository.ProductRepository
    eventBus contracts.EventBus
    logger   contracts.Logger
}

func (h *UpdateProductPricesHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    updateCmd := cmd.(*UpdateProductPricesCommand)
    
    h.logger.Info("Updating prices", map[string]interface{}{
        "count": len(updateCmd.Updates),
    })
    
    // Usar transação para operações em lote
    err := h.db.Transaction(func(tx *gorm.DB) error {
        for _, update := range updateCmd.Updates {
            // Buscar produto
            product, err := h.repo.FindByIDWithTx(ctx, tx, update.ProductID)
            if err != nil {
                h.logger.Warn("Product not found", map[string]interface{}{
                    "product_id": update.ProductID,
                })
                continue
            }
            
            // Atualizar preço
            oldPrice := product.Price
            product.UpdatePrice(update.NewPrice)
            
            // Persistir
            if err := h.repo.UpdateWithTx(ctx, tx, product); err != nil {
                return err
            }
            
            // Publicar evento
            event := events.NewProductPriceChangedEvent(
                product.ID,
                oldPrice,
                update.NewPrice,
                updateCmd.UpdatedBy,
            )
            h.eventBus.Publish(ctx, event)
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    h.logger.Info("Prices updated", map[string]interface{}{
        "count": len(updateCmd.Updates),
    })
    
    return map[string]int{"updated": len(updateCmd.Updates)}, nil
}
```

---

## 🎯 Boas Práticas

### 1. Um Command = Uma Intenção

```go
// ✅ BOM: Comandos específicos
type CreateUserCommand struct { ... }
type UpdateUserEmailCommand struct { ... }
type ActivateUserCommand struct { ... }

// ❌ RUIM: Comando genérico
type UserCommand struct {
    Action string // "create", "update", "delete"
    Data   map[string]interface{}
}
```

### 2. Commands Imutáveis

```go
// ✅ BOM: Struct imutável
type CreateProductCommand struct {
    name  string // privado, não pode ser alterado
    price float64
}

func NewCreateProductCommand(name string, price float64) *CreateProductCommand {
    return &CreateProductCommand{name: name, price: price}
}

func (c *CreateProductCommand) Name() string { return c.name }
func (c *CreateProductCommand) Price() float64 { return c.price }

// ❌ RUIM: Campos públicos mutáveis
type CreateProductCommand struct {
    Name  string
    Price float64
}
```

### 3. Validação em Camadas

```go
func (h *CreateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateProductCommand)
    
    // 1️⃣ Validação estrutural (feita no HTTP adapter)
    // 2️⃣ Validação de regras de negócio
    if err := createCmd.Validate(); err != nil {
        return nil, err
    }
    
    // 3️⃣ Validação de domínio
    product, err := entities.NewProduct(...)
    if err != nil {
        return nil, err
    }
    
    // Continuar...
}
```

### 4. Retorne Apenas o Essencial

```go
// ✅ BOM: Retorna apenas ID
func (h *CreateProductHandler) Handle(...) (interface{}, error) {
    // ...
    return product.ID, nil
}

// ❌ RUIM: Retorna entidade completa
func (h *CreateProductHandler) Handle(...) (interface{}, error) {
    // ...
    return product, nil // ❌ Vaza detalhes de domínio
}
```

### 5. Idempotência

```go
func (h *CreateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateProductCommand)
    
    // Verificar se já existe (idempotência)
    existingProduct, err := h.repo.FindBySKU(ctx, createCmd.SKU)
    if err == nil {
        // Já existe, retornar ID existente
        return existingProduct.ID, nil
    }
    
    // Criar novo produto...
}
```

### 6. Logging e Observabilidade

```go
func (h *CreateProductHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateProductCommand)
    
    h.logger.Info("Creating product", map[string]interface{}{
        "name":  createCmd.Name,
        "price": createCmd.Price,
        "sku":   createCmd.SKU,
    })
    
    product, err := entities.NewProduct(...)
    if err != nil {
        h.logger.Error("Failed to create product entity", map[string]interface{}{
            "error": err.Error(),
            "sku":   createCmd.SKU,
        })
        return nil, err
    }
    
    if err := h.repo.Save(ctx, product); err != nil {
        h.logger.Error("Failed to save product", map[string]interface{}{
            "error":      err.Error(),
            "product_id": product.ID,
        })
        return nil, err
    }
    
    h.logger.Info("Product created", map[string]interface{}{
        "product_id": product.ID,
        "sku":        product.SKU,
    })
    
    return product.ID, nil
}
```

---

## 📚 Próximos Passos

- **[Queries](14-implementing-queries.md)** - Read operations
- **[Eventos](15-working-with-events.md)** - Event-driven architecture
- **[Testes](21-testing-strategy.md)** - Testando commands

---

**[⬅️ Criando Módulos](12-creating-modules.md)** | **[Índice](README.md)** | **[Queries ➡️](14-implementing-queries.md)**
