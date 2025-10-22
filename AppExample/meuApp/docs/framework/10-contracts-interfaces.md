# 🔌 Contratos e Interfaces (Ports)

## 📋 Índice
- [O que são Ports?](#o-que-são-ports)
- [Contratos Principais](#contratos-principais)
- [Command Bus](#command-bus)
- [Query Bus](#query-bus)
- [Event Bus](#event-bus)
- [Repository Interfaces](#repository-interfaces)
- [Definindo Novos Contratos](#definindo-novos-contratos)
- [Boas Práticas](#boas-práticas)

---

## 🎯 O que são Ports?

**Ports** (Portas) são **interfaces** que definem **contratos** entre camadas da aplicação, implementando o padrão **Hexagonal Architecture**.

### Conceito Visual

```
┌─────────────────────────────────────────────────┐
│            HEXAGONAL ARCHITECTURE               │
│                                                 │
│         ┌─────────────────────────┐            │
│         │    APPLICATION CORE     │            │
│         │   (Business Logic)      │            │
│         │                         │            │
│         │  - Commands             │            │
│         │  - Queries              │            │
│         │  - Domain Logic         │            │
│         └───────┬─────────┬───────┘            │
│                 │         │                     │
│        ┌────────┘         └────────┐           │
│        │                           │            │
│   ┌────▼──────┐             ┌──────▼────┐     │
│   │   PORT    │             │   PORT    │     │
│   │(Interface)│             │(Interface)│     │
│   └────┬──────┘             └──────┬────┘     │
│        │                           │            │
│   ┌────▼──────┐             ┌──────▼────┐     │
│   │  ADAPTER  │             │  ADAPTER  │     │
│   │   (HTTP)  │             │   (DB)    │     │
│   └───────────┘             └───────────┘     │
│                                                 │
└─────────────────────────────────────────────────┘

Driving Adapters → Ports → Core → Ports → Driven Adapters
      (HTTP)        │             │        (Database)
      (gRPC)        │             │        (Email)
      (CLI)         ▼             ▼        (Queue)
```

### Por que usar Ports?

**Sem Ports (❌ Ruim):**
```go
// Application depende de implementação concreta
type UserService struct {
    db *gorm.DB  // ❌ Acoplamento direto
}

func (s *UserService) CreateUser(name, email string) error {
    return s.db.Create(&User{...}).Error  // ❌ SQL vazando
}
```

**Com Ports (✅ Bom):**
```go
// Application depende de interface (Port)
type UserService struct {
    repo UserRepository  // ✅ Interface (Port)
}

// Port (Contrato)
type UserRepository interface {
    Save(user *User) error
    FindByID(id string) (*User, error)
}

func (s *UserService) CreateUser(name, email string) error {
    user := &User{Name: name, Email: email}
    return s.repo.Save(user)  // ✅ Abstração
}

// Adapter implementa Port
type MySQLUserRepository struct {
    db *gorm.DB
}

func (r *MySQLUserRepository) Save(user *User) error {
    return r.db.Create(user).Error
}
```

**Benefícios:**
- ✅ **Dependency Inversion Principle** (SOLID)
- ✅ Application não conhece detalhes de infraestrutura
- ✅ Fácil trocar implementações (MySQL → PostgreSQL)
- ✅ Testes com mocks
- ✅ Múltiplas implementações do mesmo port

---

## 📦 Contratos Principais

### Localização no Projeto

```
pkg/contracts/
├── command_bus.go      # Command Bus interface
├── query_bus.go        # Query Bus interface
├── event_bus.go        # Event Bus interface
├── logger.go           # Logger interface
├── cache.go            # Cache interface
└── repository.go       # Repository base interfaces
```

---

## 🎯 Command Bus

### Interface

**`pkg/contracts/command_bus.go`**

```go
package contracts

import "context"

// Command representa um comando (write operation)
type Command interface {
    CommandName() string
}

// CommandHandler processa um comando
type CommandHandler interface {
    Handle(ctx context.Context, cmd Command) (interface{}, error)
}

// CommandBus despacha comandos para seus handlers
type CommandBus interface {
    // Execute executa um comando
    Execute(ctx context.Context, cmd Command) (interface{}, error)
    
    // Register registra um handler para um tipo de comando
    Register(cmdType string, handler CommandHandler) error
}
```

### Implementação

**`pkg/framework/command_bus.go`**

```go
package framework

import (
    "context"
    "fmt"
    "sync"
    "meuApp/pkg/contracts"
)

type commandBus struct {
    handlers map[string]contracts.CommandHandler
    mu       sync.RWMutex
}

func NewCommandBus() contracts.CommandBus {
    return &commandBus{
        handlers: make(map[string]contracts.CommandHandler),
    }
}

func (b *commandBus) Execute(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    cmdName := cmd.CommandName()
    
    b.mu.RLock()
    handler, exists := b.handlers[cmdName]
    b.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("no handler registered for command: %s", cmdName)
    }
    
    return handler.Handle(ctx, cmd)
}

func (b *commandBus) Register(
    cmdType string,
    handler contracts.CommandHandler,
) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    if _, exists := b.handlers[cmdType]; exists {
        return fmt.Errorf("handler already registered for: %s", cmdType)
    }
    
    b.handlers[cmdType] = handler
    return nil
}
```

### Uso

**Command:**
```go
type CreateUserCommand struct {
    Name     string
    Email    string
    Password string
}

func (c *CreateUserCommand) CommandName() string {
    return "CreateUserCommand"
}
```

**Handler:**
```go
type CreateUserHandler struct {
    repo UserRepository
}

func (h *CreateUserHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateUserCommand)
    
    user := &User{
        ID:       uuid.New().String(),
        Name:     createCmd.Name,
        Email:    createCmd.Email,
        Password: hashPassword(createCmd.Password),
    }
    
    if err := h.repo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    return user.ID, nil
}
```

**Execução:**
```go
cmd := &CreateUserCommand{
    Name:     "John Doe",
    Email:    "john@example.com",
    Password: "secret123",
}

userID, err := commandBus.Execute(ctx, cmd)
```

---

## 🔍 Query Bus

### Interface

**`pkg/contracts/query_bus.go`**

```go
package contracts

import "context"

// Query representa uma consulta (read operation)
type Query interface {
    QueryName() string
}

// QueryHandler processa uma query
type QueryHandler interface {
    Handle(ctx context.Context, qry Query) (interface{}, error)
}

// QueryBus despacha queries para seus handlers
type QueryBus interface {
    // Execute executa uma query
    Execute(ctx context.Context, qry Query) (interface{}, error)
    
    // Register registra um handler para um tipo de query
    Register(qryType string, handler QueryHandler) error
}
```

### Implementação

**`pkg/framework/query_bus.go`**

```go
package framework

import (
    "context"
    "fmt"
    "sync"
    "meuApp/pkg/contracts"
)

type queryBus struct {
    handlers map[string]contracts.QueryHandler
    mu       sync.RWMutex
}

func NewQueryBus() contracts.QueryBus {
    return &queryBus{
        handlers: make(map[string]contracts.QueryHandler),
    }
}

func (b *queryBus) Execute(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    qryName := qry.QueryName()
    
    b.mu.RLock()
    handler, exists := b.handlers[qryName]
    b.mu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("no handler registered for query: %s", qryName)
    }
    
    return handler.Handle(ctx, qry)
}

func (b *queryBus) Register(
    qryType string,
    handler contracts.QueryHandler,
) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    if _, exists := b.handlers[qryType]; exists {
        return fmt.Errorf("handler already registered for: %s", qryType)
    }
    
    b.handlers[qryType] = handler
    return nil
}
```

### Uso

**Query:**
```go
type GetUserByIDQuery struct {
    UserID string
}

func (q *GetUserByIDQuery) QueryName() string {
    return "GetUserByIDQuery"
}
```

**Handler:**
```go
type GetUserByIDHandler struct {
    repo UserRepository
}

func (h *GetUserByIDHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    getUserQuery := qry.(*GetUserByIDQuery)
    
    user, err := h.repo.FindByID(ctx, getUserQuery.UserID)
    if err != nil {
        return nil, err
    }
    
    // Retornar DTO
    return &UserDTO{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

---

## 📡 Event Bus

### Interface

**`pkg/contracts/event_bus.go`**

```go
package contracts

import "context"

// Event representa um evento de domínio
type Event interface {
    EventName() string
    OccurredAt() time.Time
}

// EventSubscriber escuta eventos
type EventSubscriber interface {
    Handle(ctx context.Context, event Event) error
    SubscribedTo() []string
}

// EventBus publica e distribui eventos
type EventBus interface {
    // Publish publica um evento
    Publish(ctx context.Context, event Event) error
    
    // Subscribe inscreve um subscriber
    Subscribe(subscriber EventSubscriber) error
    
    // Unsubscribe remove um subscriber
    Unsubscribe(eventName string, subscriber EventSubscriber) error
}
```

### Implementação

**`pkg/framework/event_bus.go`**

```go
package framework

import (
    "context"
    "sync"
    "meuApp/pkg/contracts"
)

type eventBus struct {
    subscribers map[string][]contracts.EventSubscriber
    mu          sync.RWMutex
}

func NewEventBus() contracts.EventBus {
    return &eventBus{
        subscribers: make(map[string][]contracts.EventSubscriber),
    }
}

func (b *eventBus) Publish(ctx context.Context, event contracts.Event) error {
    eventName := event.EventName()
    
    b.mu.RLock()
    subs := b.subscribers[eventName]
    b.mu.RUnlock()
    
    // Chamar subscribers assincronamente
    var wg sync.WaitGroup
    for _, sub := range subs {
        wg.Add(1)
        go func(subscriber contracts.EventSubscriber) {
            defer wg.Done()
            subscriber.Handle(ctx, event)
        }(sub)
    }
    
    wg.Wait()
    return nil
}

func (b *eventBus) Subscribe(subscriber contracts.EventSubscriber) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    for _, eventName := range subscriber.SubscribedTo() {
        b.subscribers[eventName] = append(
            b.subscribers[eventName],
            subscriber,
        )
    }
    
    return nil
}
```

### Uso

**Event:**
```go
type UserCreatedEvent struct {
    UserID    string
    Email     string
    occurredAt time.Time
}

func (e *UserCreatedEvent) EventName() string {
    return "UserCreated"
}

func (e *UserCreatedEvent) OccurredAt() time.Time {
    return e.occurredAt
}
```

**Subscriber:**
```go
type SendWelcomeEmailSubscriber struct {
    emailService EmailService
}

func (s *SendWelcomeEmailSubscriber) Handle(
    ctx context.Context,
    event contracts.Event,
) error {
    userCreated := event.(*UserCreatedEvent)
    
    return s.emailService.SendWelcomeEmail(
        userCreated.Email,
        "Welcome to our platform!",
    )
}

func (s *SendWelcomeEmailSubscriber) SubscribedTo() []string {
    return []string{"UserCreated"}
}
```

---

## 💾 Repository Interfaces

### Interface Base

**`pkg/contracts/repository.go`**

```go
package contracts

import "context"

// Repository define operações básicas de persistência
type Repository[T any] interface {
    // Create
    Save(ctx context.Context, entity T) error
    
    // Read
    FindByID(ctx context.Context, id string) (T, error)
    FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error)
    
    // Update
    Update(ctx context.Context, entity T) error
    
    // Delete
    Delete(ctx context.Context, id string) error
    
    // Exists
    Exists(ctx context.Context, id string) (bool, error)
}
```

### Interfaces Específicas

**`internal/modules/user/domain/repository/user_repository.go`**

```go
package repository

import (
    "context"
    "meuApp/internal/modules/user/domain/entities"
)

// UserRepository define operações específicas para User
type UserRepository interface {
    // CRUD
    Save(ctx context.Context, user *entities.User) error
    FindByID(ctx context.Context, id string) (*entities.User, error)
    FindAll(ctx context.Context, page, pageSize int) ([]*entities.User, int64, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id string) error
    
    // Queries específicas
    FindByEmail(ctx context.Context, email string) (*entities.User, error)
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    FindActiveUsers(ctx context.Context) ([]*entities.User, error)
    CountByStatus(ctx context.Context, active bool) (int64, error)
}
```

---

## 🆕 Definindo Novos Contratos

### Exemplo: Cache Interface

**`pkg/contracts/cache.go`**

```go
package contracts

import "time"

// Cache define operações de cache
type Cache interface {
    // Get recupera valor do cache
    Get(key string) (interface{}, bool)
    
    // Set armazena valor no cache com TTL
    Set(key string, value interface{}, ttl time.Duration) error
    
    // Delete remove valor do cache
    Delete(key string) error
    
    // Clear limpa todo o cache
    Clear() error
    
    // Exists verifica se chave existe
    Exists(key string) bool
}
```

### Implementações

**Redis Adapter:**
```go
type RedisCache struct {
    client *redis.Client
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    return c.client.Set(context.Background(), key, value, ttl).Err()
}

func (c *RedisCache) Get(key string) (interface{}, bool) {
    val, err := c.client.Get(context.Background(), key).Result()
    if err == redis.Nil {
        return nil, false
    }
    return val, true
}
```

**In-Memory Adapter:**
```go
type InMemoryCache struct {
    data map[string]cacheItem
    mu   sync.RWMutex
}

type cacheItem struct {
    value      interface{}
    expiration time.Time
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.data[key] = cacheItem{
        value:      value,
        expiration: time.Now().Add(ttl),
    }
    return nil
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    item, exists := c.data[key]
    if !exists || time.Now().After(item.expiration) {
        return nil, false
    }
    
    return item.value, true
}
```

---

## 🎯 Boas Práticas

### 1. Interfaces Pequenas (ISP - SOLID)

```go
// ✅ BOM: Interfaces pequenas e focadas
type UserReader interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserWriter interface {
    Save(ctx context.Context, user *User) error
}

// Cliente usa apenas o que precisa
type GetUserHandler struct {
    reader UserReader  // Não precisa de Write
}

// ❌ RUIM: Interface grande
type UserRepository interface {
    Save(...)
    Update(...)
    Delete(...)
    FindByID(...)
    FindByEmail(...)
    FindAll(...)
    FindActiveUsers(...)
    CountUsers(...)
    // ... 20 métodos
}
```

### 2. Retorne Erros de Domínio

```go
// ✅ Erros específicos do domínio
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrDuplicateEmail   = errors.New("email already exists")
)

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

// Handler pode tratar especificamente
user, err := repo.FindByID(ctx, id)
if err == repository.ErrUserNotFound {
    return nil, domain.ErrResourceNotFound
}
```

### 3. Use Context

```go
// ✅ Sempre aceite context
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
}

// ❌ Sem context (impossibilita timeout/cancelamento)
type UserRepository interface {
    Save(user *User) error
    FindByID(id string) (*User, error)
}
```

### 4. Interfaces no Package do Consumidor

```go
// ✅ BOM: Interface onde é usada (application layer)
package application

type UserRepository interface {
    Save(ctx context.Context, user *User) error
}

type CreateUserHandler struct {
    repo UserRepository  // Application define o contrato
}

// Implementação na infrastructure
package repository

type MySQLUserRepository struct { ... }

func (r *MySQLUserRepository) Save(...) { ... }

// ❌ RUIM: Interface na implementação
package repository

type UserRepository interface { ... }
type MySQLUserRepository struct { ... }
```

### 5. Documente Contratos

```go
// UserRepository define operações de persistência para User.
//
// Todas as operações devem ser thread-safe.
// Erros específicos devem ser retornados quando apropriado:
//   - ErrUserNotFound quando usuário não existe
//   - ErrDuplicateEmail quando email já está em uso
type UserRepository interface {
    // Save persiste um novo usuário.
    // Retorna ErrDuplicateEmail se email já existe.
    Save(ctx context.Context, user *entities.User) error
    
    // FindByID busca usuário por ID.
    // Retorna ErrUserNotFound se não encontrado.
    FindByID(ctx context.Context, id string) (*entities.User, error)
}
```

### 6. Use Generics Quando Apropriado (Go 1.18+)

```go
// Interface genérica
type Repository[T any] interface {
    Save(ctx context.Context, entity T) error
    FindByID(ctx context.Context, id string) (T, error)
    FindAll(ctx context.Context) ([]T, error)
    Delete(ctx context.Context, id string) error
}

// Uso específico
type UserRepository Repository[*entities.User]
type ProductRepository Repository[*entities.Product]
```

---

## 🧪 Testando com Mocks

### Mock Manual

```go
type MockUserRepository struct {
    SaveFunc    func(ctx context.Context, user *User) error
    FindByIDFunc func(ctx context.Context, id string) (*User, error)
}

func (m *MockUserRepository) Save(ctx context.Context, user *User) error {
    if m.SaveFunc != nil {
        return m.SaveFunc(ctx, user)
    }
    return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    if m.FindByIDFunc != nil {
        return m.FindByIDFunc(ctx, id)
    }
    return nil, ErrUserNotFound
}

// Teste
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{
        SaveFunc: func(ctx context.Context, user *User) error {
            assert.Equal(t, "John", user.Name)
            return nil
        },
    }
    
    handler := NewCreateUserHandler(mockRepo)
    err := handler.Handle(ctx, &CreateUserCommand{Name: "John"})
    
    assert.NoError(t, err)
}
```

### Mock com Testify

```go
import "github.com/stretchr/testify/mock"

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Teste
func TestCreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    handler := NewCreateUserHandler(mockRepo)
    err := handler.Handle(ctx, &CreateUserCommand{Name: "John"})
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

---

## 📚 Resumo dos Principais Ports

| Port | Propósito | Localização |
|------|-----------|-------------|
| **CommandBus** | Executa comandos (write) | `pkg/contracts/command_bus.go` |
| **QueryBus** | Executa queries (read) | `pkg/contracts/query_bus.go` |
| **EventBus** | Publica/subscreve eventos | `pkg/contracts/event_bus.go` |
| **Repository** | Persistência de dados | `pkg/contracts/repository.go` |
| **Logger** | Logging | `pkg/contracts/logger.go` |
| **Cache** | Cache de dados | `pkg/contracts/cache.go` |

---

## 📚 Próximos Passos

- **[Adapters](09-adapters.md)** - Implementações dos ports
- **[Repositórios](11-repositories.md)** - Repository pattern detalhado
- **[Dependency Injection](07-dependency-injection.md)** - Como injetar ports

---

**[⬅️ Adapters](09-adapters.md)** | **[Índice](README.md)** | **[Repositórios ➡️](11-repositories.md)**
