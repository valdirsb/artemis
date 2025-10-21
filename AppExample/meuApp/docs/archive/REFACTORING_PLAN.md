# 🚀 Plano de Ação - Refatoração da Arquitetura

> **Projeto:** meuApp - Artemis Framework  
> **Data de Criação:** 18/10/2025  
> **Objetivo:** Melhorar a arquitetura seguindo melhores práticas de Clean Architecture, Hexagonal Architecture e DDD

---

## 📊 Resumo Executivo

**Avaliação Atual:** 7.2/10  
**Meta:** 9/10  
**Tempo Estimado:** 3-4 semanas  
**Impacto:** Alto - Melhor manutenibilidade, testabilidade e escalabilidade

---

## 🎯 Fase 1: Reorganização de Estrutura (ALTA PRIORIDADE)

**Objetivo:** Corrigir a organização de diretórios e eliminar duplicações  
**Tempo Estimado:** 3-5 dias  
**Impacto:** 🔴 Alto

### 1.1 Mover Database Models para Repositórios

**Status:** ⬜ Não Iniciado

**Problema Atual:**
```
internal/shared/database/database.go
├── UserModel      ❌
├── ProductModel   ❌
└── OrderModel     ❌
```

**Solução:**
```
internal/modules/user/repository/
├── user_repository.go
└── user_model.go  ✅

internal/modules/product/repository/
├── product_repository.go
└── product_model.go  ✅

internal/modules/order/repository/
├── order_repository.go
└── order_model.go  ✅
```

#### Checklist:

- [ ] **1.1.1** Criar `internal/modules/user/repository/user_model.go`
  - [ ] Mover `UserModel` de `database.go`
  - [ ] Mover métodos `ToContract()` e `FromContract()`
  - [ ] Adicionar imports necessários

- [ ] **1.1.2** Criar `internal/modules/product/repository/product_model.go`
  - [ ] Mover `ProductModel` de `database.go`
  - [ ] Mover métodos de conversão
  - [ ] Atualizar imports

- [ ] **1.1.3** Criar `internal/modules/order/repository/order_model.go`
  - [ ] Mover `OrderModel` de `database.go`
  - [ ] Mover `OrderItemModel`
  - [ ] Mover métodos de conversão

- [ ] **1.1.4** Atualizar `database.go`
  - [ ] Remover todos os models
  - [ ] Manter apenas funções de conexão e migração
  - [ ] Atualizar função `AutoMigrate()`

- [ ] **1.1.5** Atualizar imports nos repositórios
  - [ ] `user_repository.go`
  - [ ] `product_repository.go`
  - [ ] `order_repository.go`

- [ ] **1.1.6** Testar compilação e execução
  - [ ] `go build`
  - [ ] Verificar testes unitários

---

### 1.2 Reorganizar `internal/shared` → `pkg/adapters`

**Status:** ⬜ Não Iniciado

**Problema Atual:**
```
internal/shared/
├── config/      ❌ Não é "shared" interno
├── database/    ❌ É infraestrutura
├── logger/      ❌ É adapter
└── middleware/  ❌ É HTTP-specific
```

**Nova Estrutura:**
```
pkg/
├── adapters/
│   ├── database/
│   │   ├── mysql/
│   │   │   ├── connection.go
│   │   │   └── migrations.go
│   │   └── postgres/ (futuro)
│   ├── logger/
│   │   ├── logger.go
│   │   └── zap_logger.go (exemplo)
│   └── http/
│       └── middleware/
│           ├── auth.go
│           ├── cors.go
│           └── rate_limit.go
└── config/
    ├── config.go
    └── loader.go
```

#### Checklist:

- [ ] **1.2.1** Criar nova estrutura de diretórios
  - [ ] `mkdir -p pkg/adapters/database/mysql`
  - [ ] `mkdir -p pkg/adapters/logger`
  - [ ] `mkdir -p pkg/adapters/http/middleware`
  - [ ] `mkdir -p pkg/config`

- [ ] **1.2.2** Mover `config`
  - [ ] `mv internal/shared/config/* pkg/config/`
  - [ ] Atualizar package names
  - [ ] Atualizar imports em todos os arquivos

- [ ] **1.2.3** Mover `database`
  - [ ] Mover `database.go` → `pkg/adapters/database/mysql/connection.go`
  - [ ] Criar `pkg/adapters/database/mysql/migrations.go`
  - [ ] Separar lógica de conexão de migrações
  - [ ] Atualizar imports

- [ ] **1.2.4** Mover `logger`
  - [ ] `mv internal/shared/logger/* pkg/adapters/logger/`
  - [ ] Implementar interface `Logger` em `pkg/ports/`
  - [ ] Atualizar imports

- [ ] **1.2.5** Mover `middleware`
  - [ ] `mv internal/shared/middleware/* pkg/adapters/http/middleware/`
  - [ ] Atualizar imports
  - [ ] Atualizar referências em `routes.go`

- [ ] **1.2.6** Remover `internal/shared`
  - [ ] Verificar que está vazio
  - [ ] `rm -rf internal/shared`

- [ ] **1.2.7** Atualizar todos os imports
  - [ ] `main.go`
  - [ ] `bootstrap/*.go`
  - [ ] `routes/routes.go`
  - [ ] Todos os módulos

- [ ] **1.2.8** Testar compilação
  - [ ] `go mod tidy`
  - [ ] `go build`
  - [ ] `make test` (se existir)

---

## 🔧 Fase 2: Refatoração de Interfaces e Contratos (ALTA PRIORIDADE)

**Objetivo:** Eliminar duplicação de interfaces e organizar contratos  
**Tempo Estimado:** 2-3 dias  
**Impacto:** 🔴 Alto

### 2.1 Remover Duplicação de Interfaces

**Status:** ⬜ Não Iniciado

**Problema:** Interfaces duplicadas em `pkg/contracts` e `internal/modules/*/ports`

#### Checklist:

- [ ] **2.1.1** Analisar interfaces duplicadas
  - [ ] Listar todas as interfaces em `pkg/contracts/interfaces_*.go`
  - [ ] Listar todas as interfaces em `internal/modules/*/ports/`
  - [ ] Criar matriz de duplicação

- [ ] **2.1.2** Decidir localização definitiva
  - [ ] Regra: Interfaces de domínio → `internal/modules/*/ports/`
  - [ ] Regra: Interfaces compartilhadas → `pkg/ports/`
  - [ ] Regra: DTOs → `pkg/dto/`

- [ ] **2.1.3** Criar `pkg/ports/` (interfaces compartilhadas)
  - [ ] Criar estrutura
  - [ ] Mover `Logger` interface
  - [ ] Mover `EventPublisher` interface

- [ ] **2.1.4** Limpar `pkg/contracts/interfaces_user.go`
  - [ ] Remover `UserService` (manter em `internal/modules/user/ports/`)
  - [ ] Remover `UserRepository` (manter em `internal/modules/user/ports/`)
  - [ ] Remover `UserHandler` (não deveria existir)
  - [ ] Mover `PasswordHasher` → `internal/modules/user/ports/`

- [ ] **2.1.5** Repetir para Product
  - [ ] Limpar `pkg/contracts/interfaces_product.go`
  - [ ] Manter interfaces apenas em `internal/modules/product/ports/`

- [ ] **2.1.6** Repetir para Order
  - [ ] Limpar `pkg/contracts/interfaces_order.go`
  - [ ] Manter interfaces apenas em `internal/modules/order/ports/`

- [ ] **2.1.7** Atualizar imports
  - [ ] Atualizar services
  - [ ] Atualizar handlers
  - [ ] Atualizar bootstrap

- [ ] **2.1.8** Testar compilação
  - [ ] `go build`
  - [ ] Verificar não há imports circulares

---

### 2.2 Reorganizar DTOs e Contratos

**Status:** ⬜ Não Iniciado

**Nova Estrutura:**
```
pkg/
├── dto/
│   ├── user_dto.go      # User, CreateUserRequest, UpdateUserRequest
│   ├── product_dto.go   # Product, CreateProductRequest, etc
│   └── order_dto.go     # Order, CreateOrderRequest, etc
├── ports/
│   ├── logger.go        # Logger interface
│   └── events.go        # EventPublisher interface
└── events/
    ├── user_events.go   # UserCreatedEvent, etc
    ├── product_events.go
    └── order_events.go
```

#### Checklist:

- [ ] **2.2.1** Criar estrutura de DTOs
  - [ ] `mkdir -p pkg/dto`
  - [ ] `mkdir -p pkg/events`
  - [ ] `mkdir -p pkg/ports`

- [ ] **2.2.2** Criar `pkg/dto/user_dto.go`
  - [ ] Mover `User` struct
  - [ ] Mover `CreateUserRequest`
  - [ ] Mover `UpdateUserRequest`
  - [ ] Adicionar validations tags

- [ ] **2.2.3** Criar `pkg/dto/product_dto.go`
  - [ ] Mover structs de Product
  - [ ] Mover requests/responses

- [ ] **2.2.4** Criar `pkg/dto/order_dto.go`
  - [ ] Mover structs de Order
  - [ ] Mover requests/responses

- [ ] **2.2.5** Criar `pkg/events/user_events.go`
  - [ ] Mover `UserCreatedEvent`
  - [ ] Adicionar novos eventos (UserUpdated, UserDeleted)

- [ ] **2.2.6** Criar `pkg/events/product_events.go`
  - [ ] Mover eventos de Product
  - [ ] Padronizar estrutura

- [ ] **2.2.7** Criar `pkg/events/order_events.go`
  - [ ] Mover eventos de Order
  - [ ] Padronizar estrutura

- [ ] **2.2.8** Atualizar todos os imports
  - [ ] Services
  - [ ] Handlers
  - [ ] Repositories

- [ ] **2.2.9** Remover arquivos antigos
  - [ ] Limpar `pkg/contracts/`
  - [ ] Manter apenas arquivo base se necessário

- [ ] **2.2.10** Testar
  - [ ] `go build`
  - [ ] Verificar testes

---

### 2.3 Remover Dependências de Framework das Interfaces

**Status:** ⬜ Não Iniciado

**Problema:** Handlers com dependência de `*gin.Context`

```go
// ❌ ERRADO
type UserHandler interface {
    CreateUser(ctx *gin.Context)
}

// ✅ CORRETO - Usar adapter pattern
type UserHTTPHandler struct {
    userService ports.UserService
}

func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
    // Adapter: gin.Context → UseCase Input
}
```

#### Checklist:

- [ ] **2.3.1** Remover interface `UserHandler` de contracts
  - [ ] Deletar de `pkg/contracts/interfaces_user.go`

- [ ] **2.3.2** Atualizar `user_handler.go`
  - [ ] Remover implementação de interface
  - [ ] Manter apenas struct concreta
  - [ ] Adicionar comentário explicativo

- [ ] **2.3.3** Repetir para Product e Order
  - [ ] Atualizar `product_handler.go`
  - [ ] Atualizar `order_handler.go`

- [ ] **2.3.4** Atualizar `routes/routes.go`
  - [ ] Remover casts de interface
  - [ ] Usar tipos concretos
  - [ ] Simplificar registro de rotas

- [ ] **2.3.5** Atualizar bootstrap
  - [ ] Ajustar registro de handlers
  - [ ] Remover referências à interface

---

## 🏗️ Fase 3: Implementação da Camada de Application (MÉDIA PRIORIDADE)

**Objetivo:** Separar lógica de aplicação (use cases) da lógica de domínio  
**Tempo Estimado:** 5-7 dias  
**Impacto:** 🟡 Médio-Alto

### 3.1 Criar Estrutura de Use Cases

**Status:** ⬜ Não Iniciado

**Nova Estrutura:**
```
internal/modules/user/
├── domain/              # Entities, Value Objects, Domain Services
│   ├── user.go
│   ├── user_aggregate.go
│   └── user_validations.go
├── application/         # Use Cases, Commands, Queries
│   ├── commands/
│   │   ├── create_user.go
│   │   ├── update_user.go
│   │   └── delete_user.go
│   ├── queries/
│   │   ├── get_user.go
│   │   └── list_users.go
│   └── services/
│       └── user_application_service.go
├── adapters/
│   ├── http/
│   │   └── user_http_handler.go
│   ├── grpc/
│   │   └── user_grpc_handler.go
│   └── repository/
│       ├── user_repository.go
│       └── user_model.go
└── ports/
    ├── repositories.go
    └── services.go
```

#### Checklist:

- [ ] **3.1.1** Criar estrutura de diretórios para User
  - [ ] `mkdir -p internal/modules/user/application/commands`
  - [ ] `mkdir -p internal/modules/user/application/queries`
  - [ ] `mkdir -p internal/modules/user/application/services`
  - [ ] `mkdir -p internal/modules/user/adapters/http`
  - [ ] `mkdir -p internal/modules/user/adapters/grpc`
  - [ ] `mkdir -p internal/modules/user/adapters/repository`

- [ ] **3.1.2** Criar Command: CreateUser
  - [ ] Criar `create_user_command.go`
  - [ ] Definir struct `CreateUserCommand`
  - [ ] Implementar `Handle(ctx, cmd)` method
  - [ ] Adicionar validações de negócio

- [ ] **3.1.3** Criar Command: UpdateUser
  - [ ] Criar `update_user_command.go`
  - [ ] Definir struct e handler

- [ ] **3.1.4** Criar Command: DeleteUser
  - [ ] Criar `delete_user_command.go`
  - [ ] Definir struct e handler

- [ ] **3.1.5** Criar Query: GetUser
  - [ ] Criar `get_user_query.go`
  - [ ] Definir struct `GetUserQuery`
  - [ ] Implementar `Handle(ctx, query)`

- [ ] **3.1.6** Criar Query: ListUsers (novo)
  - [ ] Criar `list_users_query.go`
  - [ ] Adicionar paginação
  - [ ] Adicionar filtros

- [ ] **3.1.7** Criar Application Service
  - [ ] Criar `user_application_service.go`
  - [ ] Orquestrar commands e queries
  - [ ] Gerenciar transações
  - [ ] Publicar eventos

- [ ] **3.1.8** Mover lógica do UserService atual
  - [ ] Extrair regras de negócio → Domain
  - [ ] Extrair orquestração → Application
  - [ ] Manter apenas coordenação

- [ ] **3.1.9** Reorganizar handlers
  - [ ] Mover `handler/user_handler.go` → `adapters/http/user_http_handler.go`
  - [ ] Mover `handler/user_grpc_handler.go` → `adapters/grpc/user_grpc_handler.go`

- [ ] **3.1.10** Mover repository
  - [ ] Mover `repository/*` → `adapters/repository/`

- [ ] **3.1.11** Atualizar imports e bootstrap

- [ ] **3.1.12** Testar módulo User completo

---

### 3.2 Replicar para Product

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **3.2.1** Criar estrutura de diretórios
- [ ] **3.2.2** Criar Commands (Create, Update, Delete, UpdateStock)
- [ ] **3.2.3** Criar Queries (Get, List)
- [ ] **3.2.4** Criar Application Service
- [ ] **3.2.5** Reorganizar handlers e repository
- [ ] **3.2.6** Atualizar bootstrap
- [ ] **3.2.7** Testar

---

### 3.3 Replicar para Order

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **3.3.1** Criar estrutura de diretórios
- [ ] **3.3.2** Criar Commands (Create, UpdateStatus, Cancel)
- [ ] **3.3.3** Criar Queries (Get, ListByUser)
- [ ] **3.3.4** Criar Application Service
- [ ] **3.3.5** Reorganizar handlers e repository
- [ ] **3.3.6** Atualizar bootstrap
- [ ] **3.3.7** Testar

---

## ⚠️ Fase 4: Sistema de Erros Tipados (MÉDIA PRIORIDADE)

**Objetivo:** Implementar erros de domínio padronizados e rastreáveis  
**Tempo Estimado:** 2-3 dias  
**Impacto:** 🟡 Médio

### 4.1 Criar Sistema de Erros Base

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **4.1.1** Criar `pkg/errors/domain_errors.go`
  ```go
  type DomainError struct {
      Code    string
      Message string
      Cause   error
      Fields  map[string]interface{}
  }
  ```

- [ ] **4.1.2** Implementar métodos
  - [ ] `Error() string`
  - [ ] `Unwrap() error`
  - [ ] `Is(error) bool`
  - [ ] `WithField(key, value)`

- [ ] **4.1.3** Criar categorias de erros
  - [ ] `ValidationError`
  - [ ] `NotFoundError`
  - [ ] `ConflictError`
  - [ ] `UnauthorizedError`
  - [ ] `InternalError`

- [ ] **4.1.4** Criar factory functions
  ```go
  func NewValidationError(message string) *DomainError
  func NewNotFoundError(entity, id string) *DomainError
  func NewConflictError(message string) *DomainError
  ```

- [ ] **4.1.5** Criar `pkg/errors/error_codes.go`
  ```go
  const (
      ErrCodeUserNotFound = "USER_NOT_FOUND"
      ErrCodeInvalidEmail = "INVALID_EMAIL"
      ErrCodeUserAlreadyExists = "USER_ALREADY_EXISTS"
      // ... etc
  )
  ```

---

### 4.2 Implementar Erros de Domínio por Módulo

**Status:** ⬜ Não Iniciado

#### Checklist - User Module:

- [ ] **4.2.1** Criar `internal/modules/user/domain/errors.go`
  ```go
  var (
      ErrUserNotFound = errors.NewNotFoundError("user", "")
      ErrInvalidEmail = errors.NewValidationError("invalid email format")
      ErrUserAlreadyExists = errors.NewConflictError("user already exists")
      ErrWeakPassword = errors.NewValidationError("password too weak")
  )
  ```

- [ ] **4.2.2** Atualizar validações em `user.go`
  - [ ] Substituir `errors.New()` por erros tipados
  - [ ] Adicionar contexto aos erros

- [ ] **4.2.3** Atualizar `user_service.go`
  - [ ] Usar erros de domínio
  - [ ] Propagar erros corretamente

- [ ] **4.2.4** Atualizar handlers HTTP
  - [ ] Mapear erros de domínio → HTTP status codes
  - [ ] Retornar JSON estruturado

- [ ] **4.2.5** Criar middleware de erro
  - [ ] `pkg/adapters/http/middleware/error_handler.go`
  - [ ] Converter DomainError → HTTP Response
  - [ ] Logging de erros

#### Checklist - Product Module:

- [ ] **4.2.6** Criar erros de domínio para Product
- [ ] **4.2.7** Atualizar validações
- [ ] **4.2.8** Atualizar service e handlers

#### Checklist - Order Module:

- [ ] **4.2.9** Criar erros de domínio para Order
- [ ] **4.2.10** Atualizar validações
- [ ] **4.2.11** Atualizar service e handlers

---

### 4.3 Melhorar Observabilidade de Erros

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **4.3.1** Adicionar stack traces
  - [ ] Integrar com `github.com/pkg/errors` ou similar
  - [ ] Capturar stack trace em erros

- [ ] **4.3.2** Adicionar request ID tracking
  - [ ] Middleware para gerar request ID
  - [ ] Propagar em context
  - [ ] Incluir em logs de erro

- [ ] **4.3.3** Implementar error reporting
  - [ ] Interface para error reporter (Sentry, etc)
  - [ ] Integrar no middleware

- [ ] **4.3.4** Criar testes para erros
  - [ ] Testar mapeamento erro → HTTP status
  - [ ] Testar serialização JSON

---

## 🎯 Fase 5: Melhorias no Event Bus (MÉDIA PRIORIDADE)

**Objetivo:** Implementar Event Bus tipado e robusto  
**Tempo Estimado:** 3-4 dias  
**Impacto:** 🟡 Médio

### 5.1 Refatorar Event Bus com Generics

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **5.1.1** Criar `pkg/events/event.go` (novo)
  ```go
  type Event[T any] struct {
      ID        string
      Type      string
      AggregateID string
      Payload   T
      Timestamp time.Time
      Metadata  map[string]string
  }
  ```

- [ ] **5.1.2** Criar `pkg/events/event_bus.go` (refatorado)
  ```go
  type EventBus interface {
      Publish[T any](ctx context.Context, event Event[T]) error
      Subscribe[T any](eventType string, handler EventHandler[T]) error
  }
  
  type EventHandler[T any] func(ctx context.Context, event Event[T]) error
  ```

- [ ] **5.1.3** Implementar InMemoryEventBus
  - [ ] Suporte a generics
  - [ ] Goroutines para handlers assíncronos
  - [ ] Error handling para handlers

- [ ] **5.1.4** Adicionar retry mechanism
  - [ ] Configurar max retries
  - [ ] Exponential backoff
  - [ ] Dead letter queue

- [ ] **5.1.5** Adicionar event persistence (opcional)
  - [ ] Interface `EventStore`
  - [ ] Implementação em memória
  - [ ] Preparar para DB persistence

---

### 5.2 Migrar Eventos Existentes

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **5.2.1** Atualizar `UserCreatedEvent`
  ```go
  type UserCreatedEvent struct {
      UserID    string
      Username  string
      Email     string
      CreatedAt time.Time
  }
  
  // Uso:
  event := events.Event[UserCreatedEvent]{
      Type: "user.created",
      Payload: UserCreatedEvent{...},
  }
  ```

- [ ] **5.2.2** Criar novos eventos User
  - [ ] `UserUpdatedEvent`
  - [ ] `UserDeletedEvent`
  - [ ] `UserEmailChangedEvent`

- [ ] **5.2.3** Atualizar eventos Product
  - [ ] `ProductCreatedEvent`
  - [ ] `ProductUpdatedEvent`
  - [ ] `StockUpdatedEvent`

- [ ] **5.2.4** Atualizar eventos Order
  - [ ] `OrderCreatedEvent`
  - [ ] `OrderStatusChangedEvent`
  - [ ] `OrderCancelledEvent`

- [ ] **5.2.5** Atualizar publishers nos services
  - [ ] User service
  - [ ] Product service
  - [ ] Order service

- [ ] **5.2.6** Atualizar subscribers
  - [ ] Criar subscribers para integração entre módulos
  - [ ] Exemplo: Order subscribe UserCreated

---

### 5.3 Adicionar Event Sourcing (Opcional/Futuro)

**Status:** ⬜ Não Iniciado (Futuro)

#### Checklist:

- [ ] **5.3.1** Estudar viabilidade
- [ ] **5.3.2** Definir aggregates para event sourcing
- [ ] **5.3.3** Implementar event store
- [ ] **5.3.4** Implementar projections
- [ ] **5.3.5** Testar reconstrução de estado

---

## 🔌 Fase 6: Auto-registro de Módulos (BAIXA PRIORIDADE)

**Objetivo:** Simplificar bootstrap com registro automático  
**Tempo Estimado:** 2-3 dias  
**Impacto:** 🟢 Baixo-Médio

### 6.1 Criar Sistema de Registry

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **6.1.1** Criar `pkg/registry/module_registry.go`
  ```go
  type Module interface {
      Name() string
      Initialize(container *container.Container) error
      RegisterHandlers(router *gin.Engine) error
      Shutdown(ctx context.Context) error
  }
  
  type Registry struct {
      modules map[string]Module
  }
  ```

- [ ] **6.1.2** Implementar métodos
  - [ ] `Register(module Module)`
  - [ ] `GetModule(name string)`
  - [ ] `InitializeAll(container)`
  - [ ] `ShutdownAll(ctx)`

- [ ] **6.1.3** Criar registry global
  ```go
  var globalRegistry = NewRegistry()
  
  func Register(module Module) {
      globalRegistry.Register(module)
  }
  ```

---

### 6.2 Implementar Module Interface por Módulo

**Status:** ⬜ Não Iniciado

#### Checklist - User Module:

- [ ] **6.2.1** Criar `internal/modules/user/module.go`
  ```go
  type UserModule struct {
      container *container.Container
  }
  
  func (m *UserModule) Name() string {
      return "user"
  }
  
  func (m *UserModule) Initialize(c *container.Container) error {
      // Register repositories
      // Register services
      // Register handlers
      return nil
  }
  
  func (m *UserModule) RegisterHandlers(router *gin.Engine) error {
      // Register routes
      return nil
  }
  ```

- [ ] **6.2.2** Criar `init()` function
  ```go
  func init() {
      registry.Register(&UserModule{})
  }
  ```

- [ ] **6.2.3** Testar auto-registro

#### Checklist - Product Module:

- [ ] **6.2.4** Criar `module.go`
- [ ] **6.2.5** Implementar interface
- [ ] **6.2.6** Adicionar auto-registro

#### Checklist - Order Module:

- [ ] **6.2.7** Criar `module.go`
- [ ] **6.2.8** Implementar interface
- [ ] **6.2.9** Adicionar auto-registro

---

### 6.3 Simplificar Bootstrap

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **6.3.1** Refatorar `bootstrap.go`
  ```go
  func Bootstrap(configPath string) (*container.Container, error) {
      c := container.NewContainer()
      
      // Initialize core infrastructure
      initializeCore(c)
      
      // Auto-initialize all registered modules
      if err := registry.InitializeAll(c); err != nil {
          return nil, err
      }
      
      return c, nil
  }
  ```

- [ ] **6.3.2** Remover código manual
  - [ ] Remover `registerModuleHandlers()`
  - [ ] Remover `start_*.go` files individuais
  - [ ] Simplificar lógica

- [ ] **6.3.3** Atualizar `routes.go`
  ```go
  func RegisterRoutes(router *gin.Engine) {
      registry.RegisterAllHandlers(router)
  }
  ```

- [ ] **6.3.4** Testar bootstrap completo

- [ ] **6.3.5** Atualizar `main.go` se necessário

---

## 📚 Fase 7: Melhorias Adicionais (BAIXA PRIORIDADE)

**Tempo Estimado:** Contínuo  
**Impacto:** 🟢 Baixo

### 7.1 Documentação

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **7.1.1** Atualizar README.md
  - [ ] Documentar nova arquitetura
  - [ ] Adicionar diagramas
  - [ ] Explicar decisões de design

- [ ] **7.1.2** Criar ARCHITECTURE.md
  - [ ] Documentar camadas
  - [ ] Explicar fluxo de dados
  - [ ] Padrões utilizados

- [ ] **7.1.3** Criar ADRs (Architecture Decision Records)
  - [ ] ADR 001: Hexagonal Architecture
  - [ ] ADR 002: DDD Tactical Patterns
  - [ ] ADR 003: Event-Driven Communication

- [ ] **7.1.4** Documentar APIs
  - [ ] Adicionar Swagger/OpenAPI
  - [ ] Documentar endpoints
  - [ ] Exemplos de uso

---

### 7.2 Testes

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **7.2.1** Unit Tests
  - [ ] Domain entities (>90% coverage)
  - [ ] Application services (>80% coverage)
  - [ ] Repositories (mock) (>80% coverage)

- [ ] **7.2.2** Integration Tests
  - [ ] Database integration
  - [ ] Event Bus integration
  - [ ] API endpoints

- [ ] **7.2.3** E2E Tests
  - [ ] User flow completo
  - [ ] Product flow completo
  - [ ] Order flow completo

- [ ] **7.2.4** Setup CI/CD
  - [ ] GitHub Actions / GitLab CI
  - [ ] Run tests automaticamente
  - [ ] Coverage report

---

### 7.3 Observabilidade

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **7.3.1** Logging Estruturado
  - [ ] Implementar com Zap/Zerolog
  - [ ] Padronizar mensagens
  - [ ] Adicionar contexto

- [ ] **7.3.2** Métricas
  - [ ] Integrar Prometheus
  - [ ] Métricas de negócio
  - [ ] Métricas de infraestrutura

- [ ] **7.3.3** Tracing
  - [ ] Integrar OpenTelemetry
  - [ ] Distributed tracing
  - [ ] Correlation IDs

- [ ] **7.3.4** Health Checks
  - [ ] Melhorar endpoint /health
  - [ ] Liveness probe
  - [ ] Readiness probe

---

### 7.4 Performance e Segurança

**Status:** ⬜ Não Iniciado

#### Checklist:

- [ ] **7.4.1** Caching
  - [ ] Integrar Redis
  - [ ] Cache de queries frequentes
  - [ ] Cache invalidation strategy

- [ ] **7.4.2** Rate Limiting
  - [ ] Implementar middleware
  - [ ] Por IP / Por usuário
  - [ ] Configurável

- [ ] **7.4.3** Autenticação/Autorização
  - [ ] JWT implementation
  - [ ] Refresh tokens
  - [ ] RBAC (Role-Based Access Control)

- [ ] **7.4.4** Validação de Input
  - [ ] Validação em todos os endpoints
  - [ ] Sanitização
  - [ ] Rate limiting

---

## 📈 Métricas de Sucesso

### KPIs:

- [ ] **Cobertura de testes:** > 80%
- [ ] **Duplicação de código:** < 3%
- [ ] **Complexidade ciclomática:** < 10 por função
- [ ] **Tempo de build:** < 30s
- [ ] **Tempo de startup:** < 5s
- [ ] **Separação de concerns:** 100% (sem vazamentos entre camadas)

### Revisões:

- [ ] **Revisão Fase 1:** ___/___/___
- [ ] **Revisão Fase 2:** ___/___/___
- [ ] **Revisão Fase 3:** ___/___/___
- [ ] **Revisão Fase 4:** ___/___/___
- [ ] **Revisão Fase 5:** ___/___/___
- [ ] **Revisão Fase 6:** ___/___/___
- [ ] **Revisão Final:** ___/___/___

---

## 🚨 Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Breaking changes em produção | Média | Alto | Feature flags, rollback plan |
| Aumento temporário de bugs | Alta | Médio | Testes extensivos, QA |
| Resistência da equipe | Baixa | Médio | Documentação, treinamento |
| Atraso no cronograma | Média | Médio | Priorização, MVP approach |

---

## 📞 Suporte

**Documentação de Referência:**
- Clean Architecture - Robert C. Martin
- Domain-Driven Design - Eric Evans
- Hexagonal Architecture - Alistair Cockburn
- Go Best Practices

**Ferramentas Recomendadas:**
- golangci-lint (linting)
- gocyclo (complexity analysis)
- go-callvis (dependency visualization)
- goconvey (testing)

---

**Última atualização:** 18/10/2025  
**Versão:** 1.0  
**Status Geral:** ⬜ Não Iniciado (0%)

