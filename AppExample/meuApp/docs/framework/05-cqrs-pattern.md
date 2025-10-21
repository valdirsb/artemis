# ⚡ Padrão CQRS (Command Query Responsibility Segregation)

## 📋 Índice
- [O que é CQRS?](#o-que-é-cqrs)
- [Por que Usar CQRS?](#por-que-usar-cqrs)
- [Commands (Escrita)](#commands-escrita)
- [Queries (Leitura)](#queries-leitura)
- [Implementação no Artemis](#implementação-no-artemis)
- [Exemplos Práticos](#exemplos-práticos)

---

## 🎯 O que é CQRS?

**CQRS** é um padrão arquitetural que separa operações de **leitura** (Queries) de operações de **escrita** (Commands).

### Modelo Tradicional (CRUD)

```
         ┌──────────────┐
         │   Service    │
         │              │
         │  • Create    │
         │  • Read      │
         │  • Update    │
         │  • Delete    │
         └──────┬───────┘
                │
         ┌──────▼───────┐
         │  Repository  │
         │   (Database) │
         └──────────────┘
```

### Modelo CQRS

```
         ┌────────────────────────────────────┐
         │       Application Layer            │
         ├──────────────┬─────────────────────┤
         │   COMMANDS   │      QUERIES        │
         │   (Write)    │      (Read)         │
         ├──────────────┼─────────────────────┤
         │ • Create     │  • Get              │
         │ • Update     │  • List             │
         │ • Delete     │  • Search           │
         └──────┬───────┴──────────┬──────────┘
                │                  │
         ┌──────▼───────┐   ┌─────▼──────────┐
         │ Write Model  │   │  Read Model    │
         │  (Complex)   │   │  (Optimized)   │
         └──────────────┘   └────────────────┘
```

### Princípio Fundamental

> **"Métodos devem retornar um resultado OU mudar estado, nunca ambos."**  
> — Bertrand Meyer (Command-Query Separation)

```go
// ❌ RUIM: Mistura query e command
func (s *UserService) CreateUser(name string) (*User, error) {
    user := &User{Name: name}
    s.repo.Save(user)  // Muda estado
    return user, nil   // E retorna resultado
}

// ✅ BOM: Separado
// Command - apenas muda estado
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    user := &User{Name: cmd.Name}
    return h.repo.Save(user)
}

// Query - apenas lê
func (h *GetUserHandler) Handle(query *GetUserQuery) (*User, error) {
    return h.repo.FindByID(query.UserID)
}
```

---

## 💪 Por que Usar CQRS?

### 1. **Clareza de Código**

```go
// Fica óbvio o que cada operação faz
commands.CreateUser()   // Escreve
commands.UpdateUser()   // Escreve
commands.DeleteUser()   // Escreve

queries.GetUser()       // Lê
queries.ListUsers()     // Lê
queries.SearchUsers()   // Lê
```

### 2. **Otimização Independente**

```go
// Write pode usar validações complexas
type CreateOrderHandler struct {
    repo              OrderRepository
    inventoryService  InventoryService
    paymentService    PaymentService
    emailService      EmailService
}

// Read pode usar views otimizadas
type ListOrdersHandler struct {
    readRepo  OrderReadRepository  // Pode ser um cache, view materializada, etc
}
```

### 3. **Escalabilidade**

```
┌─────────────┐         ┌─────────────┐
│   Writes    │         │    Reads    │
│  (Commands) │         │  (Queries)  │
│             │         │             │
│  📝 Menos   │         │  👀 Muito   │
│  frequente  │         │  frequente  │
└──────┬──────┘         └──────┬──────┘
       │                       │
┌──────▼──────┐         ┌──────▼──────┐
│   MySQL     │   sync  │    Redis    │
│  (Master)   ├────────►│   (Cache)   │
└─────────────┘         └─────────────┘

Escala verticalmente    Escala horizontalmente
```

### 4. **Segurança**

```go
// Diferentes permissões para read/write
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    if !cmd.Actor.HasPermission("user.create") {
        return errors.ErrUnauthorized
    }
    // ...
}

func (h *ListUsersHandler) Handle(query *ListUsersQuery) ([]*User, error) {
    if !query.Actor.HasPermission("user.read") {
        return nil, errors.ErrUnauthorized
    }
    // ...
}
```

### 5. **Auditoria e Rastreabilidade**

```go
// Commands são events em potencial
type CreateUserCommand struct {
    Name  string
    Email string
    Actor string  // Quem executou
    When  time.Time
}

// Fácil auditar todas as mudanças
auditLog.Record("user.created", command)
```

---

## 📝 Commands (Escrita)

### Estrutura de um Command

```go
// 1. Command DTO
type CreateUserCommand struct {
    Name     string
    Email    string
    Password string
}

// 2. Command Handler
type CreateUserHandler struct {
    repo         UserRepository
    hasher       PasswordHasher
    emailService EmailService
    eventBus     *events.EventBus
    logger       Logger
}

func NewCreateUserHandler(
    repo UserRepository,
    hasher PasswordHasher,
    emailService EmailService,
    eventBus *events.EventBus,
    logger Logger,
) *CreateUserHandler {
    return &CreateUserHandler{
        repo:         repo,
        hasher:       hasher,
        emailService: emailService,
        eventBus:     eventBus,
        logger:       logger,
    }
}

// 3. Handle method
func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) error {
    // Validação
    if err := h.validate(cmd); err != nil {
        return err
    }

    // Criar entidade
    hashedPassword, _ := h.hasher.Hash(cmd.Password)
    user := &User{
        ID:       generateID(),
        Name:     cmd.Name,
        Email:    cmd.Email,
        Password: hashedPassword,
    }

    // Persistir
    if err := h.repo.Save(ctx, user); err != nil {
        h.logger.Error("Failed to save user", err)
        return err
    }

    // Side effects
    h.emailService.SendWelcome(user.Email)
    h.eventBus.Publish("user.created", UserCreatedEvent{
        UserID: user.ID,
        Email:  user.Email,
    })

    h.logger.Info("User created successfully", user.ID)
    return nil
}
```

### Características dos Commands

✅ **Modificam estado** - CREATE, UPDATE, DELETE  
✅ **Retornam erro ou void** - Não retornam dados  
✅ **Validam regras de negócio** - Business rules  
✅ **Publicam eventos** - Domain events  
✅ **São transacionais** - All or nothing  

### Exemplos de Commands

```go
// User Module
type CreateUserCommand struct {...}
type UpdateUserCommand struct {...}
type DeleteUserCommand struct {...}
type ActivateUserCommand struct {...}
type DeactivateUserCommand struct {...}
type ChangePasswordCommand struct {...}

// Order Module
type CreateOrderCommand struct {...}
type CancelOrderCommand struct {...}
type CompleteOrderCommand struct {...}
type AddItemToOrderCommand struct {...}

// Product Module
type CreateProductCommand struct {...}
type UpdateStockCommand struct {...}
type SetProductPriceCommand struct {...}
```

---

## 🔍 Queries (Leitura)

### Estrutura de uma Query

```go
// 1. Query DTO
type GetUserQuery struct {
    UserID string
}

type ListUsersQuery struct {
    Page     int
    PageSize int
    Status   string
    SortBy   string
}

// 2. Query Handler
type GetUserHandler struct {
    repo   UserRepository
    logger Logger
}

func NewGetUserHandler(repo UserRepository, logger Logger) *GetUserHandler {
    return &GetUserHandler{
        repo:   repo,
        logger: logger,
    }
}

// 3. Handle method - retorna dados
func (h *GetUserHandler) Handle(ctx context.Context, query *GetUserQuery) (*UserDTO, error) {
    user, err := h.repo.FindByID(ctx, query.UserID)
    if err != nil {
        h.logger.Error("User not found", query.UserID)
        return nil, err
    }

    // Converte para DTO
    return &UserDTO{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    }, nil
}
```

### Características das Queries

✅ **Não modificam estado** - Apenas leitura  
✅ **Retornam dados** - DTOs, não entidades  
✅ **Podem usar cache** - Otimização de leitura  
✅ **Podem usar views otimizadas** - Denormalizadas  
✅ **Idempotentes** - Mesmo resultado sempre  

### Exemplos de Queries

```go
// User Module
type GetUserQuery struct {...}
type ListUsersQuery struct {...}
type SearchUsersQuery struct {...}
type GetUserByEmailQuery struct {...}
type GetUserStatisticsQuery struct {...}

// Order Module
type GetOrderQuery struct {...}
type ListOrdersQuery struct {...}
type GetOrdersByUserQuery struct {...}
type GetOrderSummaryQuery struct {...}

// Product Module
type GetProductQuery struct {...}
type ListProductsQuery struct {...}
type SearchProductsQuery struct {...}
type GetProductsWithLowStockQuery struct {...}
```

---

## 🏗️ Implementação no Artemis

### Organização de Pastas

```
internal/modules/user/
└── application/
    ├── commands/           # 📝 WRITE
    │   ├── create_user.go
    │   ├── update_user.go
    │   ├── delete_user.go
    │   └── change_password.go
    │
    ├── queries/            # 👀 READ
    │   ├── get_user.go
    │   ├── list_users.go
    │   ├── search_users.go
    │   └── get_user_by_email.go
    │
    └── services/
        └── user_service.go  # Orquestra commands e queries
```

### Application Service

O Application Service orquestra commands e queries:

```go
type UserApplicationService struct {
    // Command Handlers
    createUserHandler       *commands.CreateUserHandler
    updateUserHandler       *commands.UpdateUserHandler
    deleteUserHandler       *commands.DeleteUserHandler
    changePasswordHandler   *commands.ChangePasswordHandler
    
    // Query Handlers
    getUserHandler          *queries.GetUserHandler
    listUsersHandler        *queries.ListUsersHandler
    searchUsersHandler      *queries.SearchUsersHandler
    getUserByEmailHandler   *queries.GetUserByEmailHandler
}

// Commands (Write)
func (s *UserApplicationService) CreateUser(ctx context.Context, cmd *commands.CreateUserCommand) error {
    return s.createUserHandler.Handle(ctx, cmd)
}

// Queries (Read)
func (s *UserApplicationService) GetUser(ctx context.Context, query *queries.GetUserQuery) (*UserDTO, error) {
    return s.getUserHandler.Handle(ctx, query)
}
```

### Fluxo Completo

```
HTTP Request
    │
    ▼
┌─────────────────┐
│  HTTP Handler   │  ← Adapter Layer
└────────┬────────┘
         │ Converte Request → Command/Query
         ▼
┌─────────────────┐
│ App Service     │  ← Application Layer
└────────┬────────┘
         │ Delega para Handler específico
         ▼
┌─────────────────┐
│ Command/Query   │  ← Application Layer
│    Handler      │
└────────┬────────┘
         │ Usa Repository
         ▼
┌─────────────────┐
│   Repository    │  ← Infrastructure Layer
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Database     │
└─────────────────┘
```

---

## 💡 Exemplos Práticos

### Exemplo 1: Criar Usuário (Command)

```go
// 1. HTTP Request
POST /api/v1/users
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123"
}

// 2. HTTP Handler converte para Command
func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.BindJSON(&req)
    
    cmd := &commands.CreateUserCommand{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }
    
    err := h.service.CreateUser(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, gin.H{"message": "User created"})
}

// 3. Command Handler executa
func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) error {
    // Validar
    // Criar entidade
    // Salvar
    // Publicar evento
    // Log
}
```

### Exemplo 2: Listar Usuários (Query)

```go
// 1. HTTP Request
GET /api/v1/users?page=1&page_size=20&status=active

// 2. HTTP Handler converte para Query
func (h *UserHTTPHandler) ListUsers(c *gin.Context) {
    query := &queries.ListUsersQuery{
        Page:     c.Query("page"),
        PageSize: c.Query("page_size"),
        Status:   c.Query("status"),
    }
    
    result, err := h.service.ListUsers(c.Request.Context(), query)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, result)
}

// 3. Query Handler executa
func (h *ListUsersHandler) Handle(ctx context.Context, query *ListUsersQuery) (*ListUsersResponse, error) {
    users, total, err := h.repo.FindAll(ctx, query.Page, query.PageSize, query.Status)
    if err != nil {
        return nil, err
    }
    
    return &ListUsersResponse{
        Users: users,
        Total: total,
        Page:  query.Page,
    }, nil
}
```

### Exemplo 3: Atualizar Pedido (Command com Validações)

```go
type CompleteOrderCommand struct {
    OrderID   string
    ActorID   string
}

type CompleteOrderHandler struct {
    orderRepo    OrderRepository
    paymentSvc   PaymentService
    inventorySvc InventoryService
    eventBus     *events.EventBus
}

func (h *CompleteOrderHandler) Handle(ctx context.Context, cmd *CompleteOrderCommand) error {
    // 1. Buscar pedido
    order, err := h.orderRepo.FindByID(ctx, cmd.OrderID)
    if err != nil {
        return err
    }
    
    // 2. Validar regras de negócio
    if order.Status != OrderStatusPending {
        return errors.New("order is not pending")
    }
    
    if !order.IsPaid() {
        return errors.New("order is not paid")
    }
    
    // 3. Verificar estoque
    for _, item := range order.Items {
        available, err := h.inventorySvc.CheckStock(item.ProductID, item.Quantity)
        if err != nil || !available {
            return errors.New("insufficient stock")
        }
    }
    
    // 4. Executar transação
    return h.orderRepo.Transaction(ctx, func(tx *sql.Tx) error {
        // Atualizar status
        order.Complete()
        if err := h.orderRepo.Update(ctx, order); err != nil {
            return err
        }
        
        // Atualizar estoque
        for _, item := range order.Items {
            if err := h.inventorySvc.DeductStock(item.ProductID, item.Quantity); err != nil {
                return err
            }
        }
        
        // Publicar evento
        h.eventBus.Publish("order.completed", OrderCompletedEvent{
            OrderID: order.ID,
            UserID:  order.UserID,
        })
        
        return nil
    })
}
```

---

## 🎯 CQRS Avançado

### Event Sourcing + CQRS

```go
// Commands geram eventos
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    event := UserCreatedEvent{
        ID:    generateID(),
        Name:  cmd.Name,
        Email: cmd.Email,
    }
    
    // Salva o evento (não o estado)
    return h.eventStore.Append("users", event)
}

// Queries leem de projeções
func (h *GetUserHandler) Handle(query *GetUserQuery) (*User, error) {
    // Lê de uma view materializada
    return h.readModel.GetUser(query.UserID)
}
```

### Diferentes Read Models

```go
// Command sempre usa o mesmo
type WriteUserRepository interface {
    Save(user *User) error
}

// Queries podem usar diferentes fontes
type UserReadRepository interface {
    FindByID(id string) (*User, error)
}

type UserCacheRepository interface {
    FindByID(id string) (*User, error)
}

type UserSearchRepository interface {
    Search(term string) ([]*User, error)
}
```

---

## 📚 Referências

- **[Implementando Commands](13-implementing-commands.md)**
- **[Implementando Queries](14-implementing-queries.md)**
- **[Trabalhando com Eventos](15-working-with-events.md)**

---

**[⬅️ Arquitetura](04-architecture.md)** | **[Índice](README.md)** | **[Sistema de Módulos ➡️](06-modules-system.md)**
