# 🎨 Sistema de Módulos

## 📋 Índice
- [O que são Módulos?](#o-que-são-módulos)
- [Anatomia de um Módulo](#anatomia-de-um-módulo)
- [Auto-Registro](#auto-registro)
- [Ciclo de Vida](#ciclo-de-vida)
- [Dependências entre Módulos](#dependências-entre-módulos)
- [Módulos Existentes](#módulos-existentes)

---

## 🎯 O que são Módulos?

Um **módulo** no Artemis Framework é uma unidade funcional completa e autocontida que encapsula um bounded context do domínio.

### Conceito

```
┌─────────────────────────────────────────────────┐
│              MÓDULO USER                        │
│                                                 │
│  ┌─────────────┐  ┌──────────────┐            │
│  │   Domain    │  │ Application  │            │
│  │             │  │              │            │
│  │ • User      │  │ • Commands   │            │
│  │ • Email     │  │ • Queries    │            │
│  │ • Events    │  │ • Services   │            │
│  └─────────────┘  └──────────────┘            │
│                                                 │
│  ┌─────────────┐  ┌──────────────┐            │
│  │  Adapters   │  │ Repository   │            │
│  │             │  │              │            │
│  │ • HTTP      │  │ • MySQL      │            │
│  │ • gRPC      │  │ • Cache      │            │
│  └─────────────┘  └──────────────┘            │
│                                                 │
│  AUTO-REGISTRA TUDO NO FRAMEWORK               │
└─────────────────────────────────────────────────┘
```

### Características de um Módulo

✅ **Autocontido** - Tudo que precisa está dentro do módulo  
✅ **Auto-registrável** - Se registra automaticamente no framework  
✅ **Independente** - Não depende diretamente de outros módulos  
✅ **Completo** - Domain, Application, Adapters, Infrastructure  
✅ **Testável** - Pode ser testado isoladamente  

### Benefícios

```go
// Antes: Código espalhado, difícil manter
services/
  user_service.go
  product_service.go
repositories/
  user_repository.go
  product_repository.go
handlers/
  user_handler.go
  product_handler.go

// Depois: Módulos organizados e claros
modules/
  user/              # Tudo sobre User
    domain/
    application/
    adapters/
    repository/
  product/           # Tudo sobre Product
    domain/
    application/
    adapters/
    repository/
```

---

## 🏗️ Anatomia de um Módulo

### Estrutura Completa

```
internal/modules/user/
│
├── user_module.go              # ⚙️ MÓDULO PRINCIPAL
│
├── domain/                     # 🎯 DOMÍNIO
│   ├── entities/
│   │   ├── user.go
│   │   └── user_test.go
│   ├── value_objects/
│   │   ├── email.go
│   │   └── password.go
│   ├── events/
│   │   ├── user_created.go
│   │   └── user_updated.go
│   └── services/
│       └── user_domain_service.go
│
├── application/                # 💼 APLICAÇÃO
│   ├── commands/
│   │   ├── create_user.go
│   │   ├── update_user.go
│   │   ├── delete_user.go
│   │   └── *_test.go
│   ├── queries/
│   │   ├── get_user.go
│   │   ├── list_users.go
│   │   └── *_test.go
│   └── services/
│       └── user_application_service.go
│
├── adapters/                   # 🔌 ADAPTADORES
│   ├── http/
│   │   ├── handler.go
│   │   ├── routes.go
│   │   └── middleware.go
│   ├── grpc/
│   │   └── service.go
│   └── subscribers/
│       ├── email_subscriber.go
│       └── audit_subscriber.go
│
├── repository/                 # 💾 PERSISTÊNCIA
│   ├── repository.go          # Interface
│   ├── mysql_repository.go    # Implementação
│   └── cache_repository.go    # Decorator
│
├── dto/                        # 📦 DTOs
│   ├── request.go
│   ├── response.go
│   └── mapper.go
│
└── errors.go                   # ⚠️ ERROS
```

### Arquivo Principal do Módulo

**`user_module.go`**

```go
package modules

import (
    "meuApp/internal/modules/user/adapters/http"
    "meuApp/internal/modules/user/application/commands"
    "meuApp/internal/modules/user/application/queries"
    "meuApp/internal/modules/user/application/services"
    "meuApp/internal/modules/user/repository"
    "meuApp/pkg/container"
    "meuApp/pkg/events"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// UserModule implementa o módulo de usuários
type UserModule struct {
    db       *gorm.DB
    eventBus *events.EventBus
}

// NewUserModule cria uma nova instância
func NewUserModule(db *gorm.DB, eventBus *events.EventBus) *UserModule {
    return &UserModule{
        db:       db,
        eventBus: eventBus,
    }
}

// Name retorna o nome do módulo
func (m *UserModule) Name() string {
    return "user"
}

// Register registra todos os componentes no ModuleRegistry
func (m *UserModule) Register(registry *container.ModuleRegistry) error {
    fmt.Printf("🔧 Registering module: %s\n", m.Name())

    // 1. Registrar Repository
    userRepo := repository.NewMySQLUserRepository(m.db)
    registry.RegisterRepository("user", userRepo)

    // 2. Registrar Services auxiliares
    passwordHasher := adapters.NewArgon2PasswordHasher()
    emailService := adapters.NewEmailService()
    logger := adapters.NewStructuredLogger()

    registry.Container().Register("user.password_hasher", passwordHasher)
    registry.Container().Register("user.email_service", emailService)
    registry.Container().Register("user.logger", logger)

    // 3. Criar Command Handlers
    createUserHandler := commands.NewCreateUserHandler(
        userRepo, passwordHasher, emailService, m.eventBus, logger,
    )
    updateUserHandler := commands.NewUpdateUserHandler(userRepo, logger)
    deleteUserHandler := commands.NewDeleteUserHandler(userRepo, m.eventBus, logger)

    // 4. Criar Query Handlers
    getUserHandler := queries.NewGetUserHandler(userRepo, logger)
    listUsersHandler := queries.NewListUsersHandler(userRepo, logger)

    // 5. Criar Application Service
    userAppService := services.NewUserApplicationService(
        createUserHandler,
        updateUserHandler,
        deleteUserHandler,
        getUserHandler,
        listUsersHandler,
    )
    registry.RegisterApplicationService("user", userAppService)

    // 6. Registrar HTTP Handler
    httpHandler := http.NewUserHTTPHandler(userAppService)
    registry.RegisterHTTPHandler("user", httpHandler)

    // 7. Registrar gRPC Service
    grpcService := grpc.NewUserGRPCService(userAppService)
    registry.RegisterGRPCService(grpcService)

    // 8. Registrar Event Subscribers
    emailSubscriber := subscribers.NewUserEmailSubscriber(emailService, logger)
    emailSubscriber.Subscribe(m.eventBus)

    fmt.Printf("✅ Module %s registered successfully\n", m.Name())
    return nil
}
```

---

## 🔄 Auto-Registro

### Interface do Módulo

Todos os módulos implementam a interface `Module`:

```go
// pkg/framework/interfaces/module.go
package interfaces

import "meuApp/pkg/container"

// Module representa um módulo da aplicação
type Module interface {
    // Name retorna o nome do módulo
    Name() string
    
    // Register registra o módulo no ModuleRegistry
    Register(registry *container.ModuleRegistry) error
}
```

### Fluxo de Auto-Registro

```
1. Bootstrap cria ModuleRegistry
         │
         ▼
2. Cria instâncias dos módulos
         │
         ├─► UserModule
         ├─► ProductModule
         └─► OrderModule
         │
         ▼
3. Cada módulo se registra
         │
         ├─► user.Register(registry)
         │     ├─► Registra Repository
         │     ├─► Registra Handlers
         │     ├─► Registra Services
         │     └─► Registra HTTP/gRPC
         │
         ├─► product.Register(registry)
         │
         └─► order.Register(registry)
         │
         ▼
4. Framework está pronto
```

### Código de Bootstrap

**`internal/bootstrap/bootstrap_registry.go`**

```go
func registerModules(
    registry *container.ModuleRegistry,
    db *gorm.DB,
    eventBus *events.EventBus,
    logger contracts.Logger,
) error {
    log.Println("📦 Registering modules...")

    // Lista de módulos (ORDEM IMPORTA!)
    applicationModules := []frameworkInterfaces.Module{
        modules.NewUserModule(db, eventBus),           // 1º - Independente
        modules.NewProductModule(db, eventBus, logger), // 2º - Independente
        modules.NewOrderModule(db, eventBus, logger),   // 3º - Depende de User + Product
    }

    // Registrar cada módulo
    for _, module := range applicationModules {
        log.Printf("  → Registering module: %s", module.Name())
        if err := module.Register(registry); err != nil {
            return fmt.Errorf("failed to register module %s: %w", module.Name(), err)
        }
    }

    log.Printf("✅ All %d modules registered successfully", len(applicationModules))
    return nil
}
```

### Benefícios do Auto-Registro

✅ **Menos Boilerplate** - Não precisa registrar cada componente manualmente  
✅ **Configuração Centralizada** - Tudo em um lugar  
✅ **Fácil Adicionar/Remover** - Só adiciona/remove da lista  
✅ **Type Safe** - Compilador valida  
✅ **Explícito** - Fica claro o que cada módulo precisa  

---

## ⏱️ Ciclo de Vida

### Fases de um Módulo

```
┌─────────────────────────────────────────────────┐
│  1. CRIAÇÃO (Constructor)                      │
│     NewUserModule(db, eventBus)                │
│     - Recebe dependências externas             │
│     - Não faz nada pesado aqui                 │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│  2. REGISTRO (Register)                        │
│     module.Register(registry)                  │
│     - Cria e registra componentes              │
│     - Conecta dependências                     │
│     - Configura subscribers                    │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│  3. OPERAÇÃO (Runtime)                         │
│     - Processa requisições                     │
│     - Executa comandos/queries                 │
│     - Publica/consome eventos                  │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│  4. SHUTDOWN (Cleanup)                         │
│     - Fecha conexões                           │
│     - Flush de caches                          │
│     - Finaliza workers                         │
└─────────────────────────────────────────────────┘
```

### Exemplo de Ciclo Completo

```go
// 1. CRIAÇÃO
func main() {
    db := connectDB()
    eventBus := events.NewEventBus()
    
    userModule := modules.NewUserModule(db, eventBus)
    // Apenas cria, não inicializa nada pesado
}

// 2. REGISTRO
func bootstrap() {
    registry := container.NewModuleRegistry(container)
    
    // Módulo se registra
    userModule.Register(registry)
    // Aqui todos os componentes são criados e conectados
}

// 3. OPERAÇÃO
func handleRequest() {
    // Módulo está operacional
    result, err := userService.CreateUser(cmd)
}

// 4. SHUTDOWN
func shutdown() {
    // Limpeza (se necessário)
    userModule.Cleanup() // Opcional
}
```

---

## 🔗 Dependências entre Módulos

### Regra de Ouro

> **Módulos NÃO devem depender diretamente uns dos outros!**

### ❌ Errado - Acoplamento Direto

```go
// ❌ OrderModule importando UserModule diretamente
import "meuApp/internal/modules/user/domain/entities"

type OrderModule struct {
    userService *user.UserService  // ❌ Acoplamento!
}
```

### ✅ Correto - Comunicação via Eventos

```go
// ✅ OrderModule publica evento
func (h *CreateOrderHandler) Handle(cmd *CreateOrderCommand) error {
    order := &Order{...}
    h.repo.Save(order)
    
    // Publica evento (não conhece quem vai processar)
    h.eventBus.Publish("order.created", OrderCreatedEvent{
        OrderID: order.ID,
        UserID:  order.UserID,
    })
}

// ✅ UserModule subscreve ao evento
type UserStatsSubscriber struct {
    userRepo UserRepository
}

func (s *UserStatsSubscriber) Subscribe(eventBus *events.EventBus) {
    eventBus.Subscribe("order.created", s.onOrderCreated)
}

func (s *UserStatsSubscriber) onOrderCreated(event Event) error {
    data := event.Data.(OrderCreatedEvent)
    // Atualiza estatísticas do usuário
    return s.userRepo.IncrementOrderCount(data.UserID)
}
```

### Tipos de Dependências

#### 1. Dependências de Infraestrutura (OK)

```go
// ✅ OK: Módulos compartilham infraestrutura
type UserModule struct {
    db       *gorm.DB        // Compartilhado
    eventBus *events.EventBus // Compartilhado
    logger   Logger          // Compartilhado
}

type ProductModule struct {
    db       *gorm.DB        // Mesmo DB
    eventBus *events.EventBus // Mesmo EventBus
    logger   Logger          // Mesmo Logger
}
```

#### 2. Dependências de Domínio (❌ EVITAR)

```go
// ❌ EVITAR: Módulo dependendo de entidade de outro módulo
import "meuApp/internal/modules/user/domain/entities"

func (h *CreateOrderHandler) Handle(cmd *CreateOrderCommand) error {
    user := entities.User{} // ❌ Não faça isso!
}

// ✅ MELHOR: Use IDs e eventos
func (h *CreateOrderHandler) Handle(cmd *CreateOrderCommand) error {
    order := &Order{
        UserID: cmd.UserID, // Só armazena o ID
    }
    
    // Valida se user existe via serviço
    if !h.userExists(cmd.UserID) {
        return errors.New("user not found")
    }
}
```

#### 3. Dependências de Ordem (CUIDADO)

Alguns módulos precisam ser registrados antes de outros:

```go
// ORDEM IMPORTA!
applicationModules := []frameworkInterfaces.Module{
    modules.NewUserModule(db, eventBus),      // 1º - Base
    modules.NewProductModule(db, eventBus),   // 2º - Base
    modules.NewOrderModule(db, eventBus),     // 3º - Depende de User + Product
    modules.NewPaymentModule(db, eventBus),   // 4º - Depende de Order
}
```

### Padrão Anti-Corruption Layer

Para casos onde você PRECISA acessar outro módulo:

```go
// Anti-Corruption Layer - Interface local
type UserService interface {
    UserExists(userID string) (bool, error)
    GetUserEmail(userID string) (string, error)
}

// Adapter que chama o outro módulo
type UserServiceAdapter struct {
    userRepo UserRepository // Via DI
}

func (a *UserServiceAdapter) UserExists(userID string) (bool, error) {
    _, err := a.userRepo.FindByID(userID)
    return err == nil, nil
}

// OrderModule usa a interface, não o módulo diretamente
type CreateOrderHandler struct {
    orderRepo   OrderRepository
    userService UserService  // Interface!
}
```

---

## 📦 Módulos Existentes

### User Module

**Responsabilidade**: Gerenciamento de usuários e autenticação

```
Entidades:
  - User

Commands:
  - CreateUser
  - UpdateUser
  - DeleteUser
  - ValidateCredentials

Queries:
  - GetUser
  - ListUsers
  - GetUserByEmail

Eventos:
  - user.created
  - user.updated
  - user.deleted
```

### Product Module

**Responsabilidade**: Catálogo de produtos

```
Entidades:
  - Product

Commands:
  - CreateProduct
  - UpdateProduct
  - DeleteProduct
  - UpdateStock

Queries:
  - GetProduct
  - ListProducts
  - SearchProducts

Eventos:
  - product.created
  - product.updated
  - stock.updated
```

### Order Module

**Responsabilidade**: Gerenciamento de pedidos

```
Entidades:
  - Order
  - OrderItem

Commands:
  - CreateOrder
  - CancelOrder
  - CompleteOrder

Queries:
  - GetOrder
  - ListOrders
  - GetOrdersByUser

Eventos:
  - order.created
  - order.cancelled
  - order.completed

Dependências:
  - Precisa validar UserID (via evento ou ACL)
  - Precisa validar ProductID (via evento ou ACL)
```

---

## 🎯 Criando um Novo Módulo

### Checklist

- [ ] 1. Criar estrutura de pastas
- [ ] 2. Definir entidades de domínio
- [ ] 3. Criar repository (interface + implementação)
- [ ] 4. Implementar commands
- [ ] 5. Implementar queries
- [ ] 6. Criar application service
- [ ] 7. Criar HTTP handler
- [ ] 8. Criar gRPC service (opcional)
- [ ] 9. Criar arquivo do módulo (`module.go`)
- [ ] 10. Registrar no bootstrap
- [ ] 11. Adicionar ao `framework.yaml`
- [ ] 12. Escrever testes

### Template Rápido

```bash
# Estrutura mínima
mkdir -p internal/modules/mymodule/{domain/entities,application/{commands,queries,services},adapters/http,repository,dto}

# Criar arquivo do módulo
touch internal/modules/mymodule_module.go

# Adicionar ao framework.yaml
echo "  mymodule: true" >> framework.yaml
```

👉 **[Tutorial completo em 12-creating-modules.md](12-creating-modules.md)**

---

## 🎨 Boas Práticas

### 1. Um Módulo = Um Bounded Context

```go
✅ Módulo "User" - Autenticação, perfil, permissões
✅ Módulo "Product" - Catálogo, estoque, categorias
✅ Módulo "Order" - Pedidos, checkout, pagamento

❌ Módulo "Everything" - Tudo misturado
```

### 2. Módulos Devem Ser Pequenos

```go
✅ Módulo "Blog" com Post, Comment, Tag
❌ Módulo "Content" com Blog, Wiki, Forum, News, etc
```

### 3. Comunicação Via Eventos

```go
✅ Publish("order.created", event)
❌ userModule.UpdateStats(userID)
```

### 4. Não Exponha Entidades

```go
✅ Retorna DTO
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
}

❌ Retorna Entity
return &entities.User{}
```

### 5. Testes Isolados

```go
✅ Testa módulo sem depender de outros
func TestUserModule(t *testing.T) {
    mockDB := setupMockDB()
    module := NewUserModule(mockDB, mockEventBus)
    // Testa isoladamente
}
```

---

## 📚 Próximos Passos

- **[Dependency Injection](07-dependency-injection.md)** - Como funciona o DI
- **[Criando um Módulo](12-creating-modules.md)** - Tutorial completo
- **[Sistema de Eventos](08-events-system.md)** - Comunicação entre módulos

---

**[⬅️ CQRS](05-cqrs-pattern.md)** | **[Índice](README.md)** | **[DI Container ➡️](07-dependency-injection.md)**
