# 💉 Dependency Injection e Container

## 📋 Índice
- [O que é Dependency Injection?](#o-que-é-dependency-injection)
- [Container DI](#container-di)
- [ModuleRegistry](#moduleregistry)
- [Padrões de Uso](#padrões-de-uso)
- [Resolução de Dependências](#resolução-de-dependências)

---

## 🎯 O que é Dependency Injection?

**Dependency Injection (DI)** é um padrão onde as dependências de um objeto são fornecidas externamente, em vez de serem criadas internamente.

### Sem DI (❌ Ruim)

```go
type UserService struct {
    repo UserRepository
}

func NewUserService() *UserService {
    return &UserService{
        repo: NewMySQLUserRepository(), // ❌ Acoplamento forte!
    }
}
```

**Problemas:**
- ❌ Acoplado a implementação específica (MySQL)
- ❌ Difícil de testar (não pode mockar)
- ❌ Difícil de trocar implementação
- ❌ Viola Dependency Inversion Principle

### Com DI (✅ Bom)

```go
type UserService struct {
    repo UserRepository // Interface!
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{
        repo: repo, // ✅ Injetado de fora
    }
}

// Uso em produção
mysqlRepo := NewMySQLUserRepository(db)
service := NewUserService(mysqlRepo)

// Uso em testes
mockRepo := &MockUserRepository{}
service := NewUserService(mockRepo)
```

**Benefícios:**
- ✅ Baixo acoplamento
- ✅ Fácil de testar
- ✅ Fácil de trocar implementação
- ✅ Segue SOLID

---

## 📦 Container DI

O **Container** é um objeto que gerencia a criação e resolução de dependências.

### Estrutura Básica

**`pkg/container/container.go`**

```go
package container

import (
    "fmt"
    "reflect"
    "sync"
)

// Container é um simples DI container
type Container struct {
    services map[string]interface{}
    mu       sync.RWMutex
}

// NewContainer cria uma nova instância do container
func NewContainer() *Container {
    return &Container{
        services: make(map[string]interface{}),
    }
}

// Register registra um serviço no container
func (c *Container) Register(name string, service interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.services[name] = service
}

// RegisterSingleton registra um singleton no container
func (c *Container) RegisterSingleton(name string, factory func() interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Lazy initialization
    c.services[name] = &singleton{
        factory: factory,
        once:    &sync.Once{},
    }
}

// Get obtém um serviço do container
func (c *Container) Get(name string) (interface{}, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    service, exists := c.services[name]
    if !exists {
        return nil, fmt.Errorf("service '%s' not found", name)
    }

    // Se for um singleton, inicializa se necessário
    if s, ok := service.(*singleton); ok {
        s.once.Do(func() {
            s.instance = s.factory()
        })
        return s.instance, nil
    }

    return service, nil
}

// MustGet obtém um serviço ou entra em pânico
func (c *Container) MustGet(name string) interface{} {
    service, err := c.Get(name)
    if err != nil {
        panic(err)
    }
    return service
}

// GetAs obtém um serviço com cast de tipo
func (c *Container) GetAs(name string, target interface{}) error {
    service, err := c.Get(name)
    if err != nil {
        return err
    }

    serviceValue := reflect.ValueOf(service)
    targetValue := reflect.ValueOf(target)

    if targetValue.Kind() != reflect.Ptr {
        return fmt.Errorf("target must be a pointer")
    }

    targetElem := targetValue.Elem()
    if !serviceValue.Type().AssignableTo(targetElem.Type()) {
        return fmt.Errorf("cannot assign %s to %s", 
            serviceValue.Type(), targetElem.Type())
    }

    targetElem.Set(serviceValue)
    return nil
}

type singleton struct {
    factory  func() interface{}
    instance interface{}
    once     *sync.Once
}
```

### Uso Básico

```go
// 1. Criar container
container := container.NewContainer()

// 2. Registrar serviços
container.Register("logger", logger)
container.Register("database", db)

// 3. Registrar singletons
container.RegisterSingleton("cache", func() interface{} {
    return redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
})

// 4. Recuperar serviços
logger, err := container.Get("logger")
db, err := container.Get("database")

// 5. Ou com type safety
var logger Logger
container.GetAs("logger", &logger)
```

### Padrão Singleton

```go
// Registra factory, não instância
container.RegisterSingleton("expensive_service", func() interface{} {
    // Só executado na primeira vez que Get é chamado
    conn := connectToExpensiveResource()
    return NewExpensiveService(conn)
})

// Primeira vez: cria a instância
svc1 := container.MustGet("expensive_service")

// Próximas vezes: retorna mesma instância
svc2 := container.MustGet("expensive_service")

// svc1 == svc2 (mesma instância!)
```

---

## 🗂️ ModuleRegistry

O **ModuleRegistry** estende o Container para gerenciar módulos completos.

### Estrutura

**`pkg/container/registry.go`**

```go
package container

import (
    "context"
    "fmt"
    "sync"

    "github.com/gin-gonic/gin"
    "google.golang.org/grpc"
)

// ModuleRegistry é o registro central para componentes dos módulos
type ModuleRegistry struct {
    container *Container

    // Handlers HTTP
    httpHandlers map[string]HTTPHandler

    // Handlers gRPC
    grpcServices []GRPCServiceRegistrar

    // Repositórios
    repositories map[string]interface{}

    // Application Services
    appServices map[string]interface{}

    // Event Subscribers
    eventSubscribers []EventSubscriber

    mu sync.RWMutex
}

// HTTPHandler representa um handler HTTP
type HTTPHandler interface {
    RegisterRoutes(router *gin.RouterGroup)
}

// GRPCServiceRegistrar representa um serviço gRPC
type GRPCServiceRegistrar interface {
    RegisterService(server *grpc.Server)
}

// EventSubscriber representa um subscriber de eventos
type EventSubscriber interface {
    Subscribe(eventBus interface{}) error
}

// NewModuleRegistry cria um novo registry de módulos
func NewModuleRegistry(container *Container) *ModuleRegistry {
    return &ModuleRegistry{
        container:        container,
        httpHandlers:     make(map[string]HTTPHandler),
        grpcServices:     make([]GRPCServiceRegistrar, 0),
        repositories:     make(map[string]interface{}),
        appServices:      make(map[string]interface{}),
        eventSubscribers: make([]EventSubscriber, 0),
    }
}

// Container retorna o container DI
func (r *ModuleRegistry) Container() *Container {
    return r.container
}

// RegisterHTTPHandler registra um handler HTTP
func (r *ModuleRegistry) RegisterHTTPHandler(name string, handler HTTPHandler) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.httpHandlers[name] = handler
    r.container.Register(fmt.Sprintf("http.handler.%s", name), handler)
}

// RegisterGRPCService registra um serviço gRPC
func (r *ModuleRegistry) RegisterGRPCService(service GRPCServiceRegistrar) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.grpcServices = append(r.grpcServices, service)
}

// RegisterRepository registra um repositório
func (r *ModuleRegistry) RegisterRepository(name string, repo interface{}) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.repositories[name] = repo
    r.container.Register(fmt.Sprintf("repository.%s", name), repo)
}

// RegisterApplicationService registra um application service
func (r *ModuleRegistry) RegisterApplicationService(name string, service interface{}) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.appServices[name] = service
    r.container.Register(fmt.Sprintf("app.service.%s", name), service)
}

// RegisterEventSubscriber registra um event subscriber
func (r *ModuleRegistry) RegisterEventSubscriber(subscriber EventSubscriber) {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.eventSubscribers = append(r.eventSubscribers, subscriber)
}

// RegisterHTTPRoutes registra todas as rotas HTTP
func (r *ModuleRegistry) RegisterHTTPRoutes(router *gin.RouterGroup) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for name, handler := range r.httpHandlers {
        log.Printf("  → Registering HTTP routes for: %s", name)
        handler.RegisterRoutes(router)
    }
}

// RegisterGRPCServices registra todos os serviços gRPC
func (r *ModuleRegistry) RegisterGRPCServices(server *grpc.Server) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for _, service := range r.grpcServices {
        service.RegisterService(server)
    }
}

// Stats retorna estatísticas do registry
func (r *ModuleRegistry) Stats() RegistryStats {
    r.mu.RLock()
    defer r.mu.RUnlock()

    return RegistryStats{
        HTTPHandlers:     len(r.httpHandlers),
        GRPCServices:     len(r.grpcServices),
        Repositories:     len(r.repositories),
        AppServices:      len(r.appServices),
        EventSubscribers: len(r.eventSubscribers),
    }
}

type RegistryStats struct {
    HTTPHandlers     int
    GRPCServices     int
    Repositories     int
    AppServices      int
    EventSubscribers int
}

func (s RegistryStats) String() string {
    return fmt.Sprintf(
        "HTTP: %d, gRPC: %d, Repos: %d, Services: %d, Subscribers: %d",
        s.HTTPHandlers, s.GRPCServices, s.Repositories, 
        s.AppServices, s.EventSubscribers,
    )
}
```

### Uso no Módulo

```go
func (m *UserModule) Register(registry *container.ModuleRegistry) error {
    // 1. Registrar componentes simples
    logger := NewLogger()
    registry.Container().Register("user.logger", logger)

    // 2. Registrar repository
    userRepo := NewMySQLUserRepository(m.db)
    registry.RegisterRepository("user", userRepo)

    // 3. Criar handlers com dependências
    createHandler := NewCreateUserHandler(userRepo, logger)
    listHandler := NewListUsersHandler(userRepo, logger)

    // 4. Registrar application service
    userService := NewUserApplicationService(createHandler, listHandler)
    registry.RegisterApplicationService("user", userService)

    // 5. Registrar HTTP handler
    httpHandler := NewUserHTTPHandler(userService)
    registry.RegisterHTTPHandler("user", httpHandler)

    return nil
}
```

---

## 🎯 Padrões de Uso

### 1. Constructor Injection (Recomendado)

```go
// ✅ Dependências injetadas no constructor
type CreateUserHandler struct {
    repo      UserRepository
    hasher    PasswordHasher
    emailSvc  EmailService
    eventBus  *events.EventBus
    logger    Logger
}

func NewCreateUserHandler(
    repo UserRepository,
    hasher PasswordHasher,
    emailSvc EmailService,
    eventBus *events.EventBus,
    logger Logger,
) *CreateUserHandler {
    return &CreateUserHandler{
        repo:     repo,
        hasher:   hasher,
        emailSvc: emailSvc,
        eventBus: eventBus,
        logger:   logger,
    }
}
```

**Benefícios:**
- ✅ Dependências explícitas
- ✅ Imutável após criação
- ✅ Fácil de testar
- ✅ Type safe

### 2. Service Locator (Evitar)

```go
// ❌ Anti-pattern: Service Locator
type CreateUserHandler struct {
    container *Container
}

func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    // ❌ Busca dependências do container
    repo := h.container.MustGet("repository.user").(UserRepository)
    logger := h.container.MustGet("logger").(Logger)
    
    // Lógica...
}
```

**Problemas:**
- ❌ Dependências ocultas
- ❌ Falha em runtime, não compile time
- ❌ Difícil de testar
- ❌ Viola princípios SOLID

### 3. Optional Dependencies

```go
type UserService struct {
    repo   UserRepository  // Obrigatório
    cache  Cache          // Opcional
    logger Logger         // Opcional
}

func NewUserService(repo UserRepository, opts ...Option) *UserService {
    s := &UserService{
        repo:   repo,
        logger: &NoOpLogger{}, // Default
    }
    
    for _, opt := range opts {
        opt(s)
    }
    
    return s
}

type Option func(*UserService)

func WithCache(cache Cache) Option {
    return func(s *UserService) {
        s.cache = cache
    }
}

func WithLogger(logger Logger) Option {
    return func(s *UserService) {
        s.logger = logger
    }
}

// Uso
service := NewUserService(
    repo,
    WithCache(redisCache),
    WithLogger(logger),
)
```

---

## 🔍 Resolução de Dependências

### Grafo de Dependências

```
┌─────────────────┐
│ UserHTTPHandler │
└────────┬────────┘
         │
         ▼
┌────────────────────┐
│ UserAppService     │
└────────┬───────────┘
         │
         ├──────────────┬──────────────┐
         ▼              ▼              ▼
┌────────────────┐ ┌──────────┐ ┌──────────┐
│CreateUserHandler│ │GetUser   │ │ListUsers │
│                │ │Handler   │ │Handler   │
└────────┬───────┘ └────┬─────┘ └────┬─────┘
         │              │            │
         ├──────────────┴────────────┘
         ▼
┌────────────────┐
│ UserRepository │
└────────────────┘
```

### Ordem de Criação

```go
// 1. Criar dependências de nível mais baixo
userRepo := NewMySQLUserRepository(db)
logger := NewLogger()
hasher := NewPasswordHasher()

// 2. Criar handlers
createHandler := NewCreateUserHandler(userRepo, hasher, logger)
getHandler := NewGetUserHandler(userRepo, logger)
listHandler := NewListUsersHandler(userRepo, logger)

// 3. Criar service
userService := NewUserApplicationService(
    createHandler,
    getHandler,
    listHandler,
)

// 4. Criar handler HTTP
httpHandler := NewUserHTTPHandler(userService)
```

### Dependências Circulares (❌ EVITAR)

```go
// ❌ EVITAR: Dependência circular
type UserService struct {
    orderService *OrderService
}

type OrderService struct {
    userService *UserService  // Circular!
}
```

**Solução: Use eventos ou interfaces**

```go
// ✅ Solução: Interface
type UserProvider interface {
    GetUser(id string) (*User, error)
}

type OrderService struct {
    userProvider UserProvider  // Interface, não implementação
}

// ✅ Ou use eventos
type OrderService struct {
    eventBus *events.EventBus
}

func (s *OrderService) CreateOrder() {
    // Publica evento em vez de chamar UserService
    s.eventBus.Publish("order.created", event)
}
```

---

## 🎨 Exemplos Práticos

### Exemplo 1: Setup Completo

```go
// main.go
func main() {
    // 1. Criar container
    container := container.NewContainer()

    // 2. Registrar infraestrutura
    db := connectDB()
    container.Register("database", db)
    
    eventBus := events.NewEventBus()
    container.Register("eventbus", eventBus)
    
    logger := NewLogger()
    container.Register("logger", logger)

    // 3. Criar ModuleRegistry
    registry := container.NewModuleRegistry(container)

    // 4. Registrar módulos
    userModule := modules.NewUserModule(db, eventBus)
    userModule.Register(registry)

    productModule := modules.NewProductModule(db, eventBus, logger)
    productModule.Register(registry)

    // 5. Setup HTTP server
    router := gin.New()
    api := router.Group("/api/v1")
    registry.RegisterHTTPRoutes(api)

    // 6. Iniciar servidor
    router.Run(":8080")
}
```

### Exemplo 2: Teste com Mocks

```go
func TestCreateUserHandler(t *testing.T) {
    // 1. Criar mocks
    mockRepo := &MockUserRepository{}
    mockHasher := &MockPasswordHasher{}
    mockEmail := &MockEmailService{}
    mockEventBus := events.NewEventBus()
    mockLogger := &MockLogger{}

    // 2. Injetar mocks
    handler := NewCreateUserHandler(
        mockRepo,
        mockHasher,
        mockEmail,
        mockEventBus,
        mockLogger,
    )

    // 3. Setup expectations
    mockRepo.On("Save", mock.Anything).Return(nil)
    mockHasher.On("Hash", "password123").Return("hashed", nil)
    mockEmail.On("SendWelcome", mock.Anything).Return(nil)

    // 4. Executar
    cmd := &CreateUserCommand{
        Name:     "John",
        Email:    "john@example.com",
        Password: "password123",
    }
    err := handler.Handle(context.Background(), cmd)

    // 5. Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Exemplo 3: Múltiplas Implementações

```go
// Interface
type NotificationService interface {
    Send(to, message string) error
}

// Implementação 1: Email
type EmailNotificationService struct {
    smtp *smtp.Client
}

// Implementação 2: SMS
type SMSNotificationService struct {
    twilio *twilio.Client
}

// Implementação 3: Push
type PushNotificationService struct {
    fcm *fcm.Client
}

// Composite que usa todas
type CompositeNotificationService struct {
    services []NotificationService
}

func (c *CompositeNotificationService) Send(to, message string) error {
    for _, svc := range c.services {
        if err := svc.Send(to, message); err != nil {
            log.Printf("Failed to send via %T: %v", svc, err)
        }
    }
    return nil
}

// Setup
emailSvc := NewEmailNotificationService(smtp)
smsSvc := NewSMSNotificationService(twilio)
pushSvc := NewPushNotificationService(fcm)

compositeSvc := &CompositeNotificationService{
    services: []NotificationService{emailSvc, smsSvc, pushSvc},
}

container.Register("notification", compositeSvc)
```

---

## 🎯 Boas Práticas

### 1. Prefira Constructor Injection

```go
// ✅ Bom
func NewService(repo Repository, logger Logger) *Service {
    return &Service{repo: repo, logger: logger}
}

// ❌ Ruim
func NewService() *Service {
    return &Service{
        repo: getRepoFromSomewhere(),
    }
}
```

### 2. Use Interfaces

```go
// ✅ Dependa de interface
type Service struct {
    repo UserRepository  // Interface
}

// ❌ Não dependa de implementação
type Service struct {
    repo *MySQLUserRepository  // Implementação concreta
}
```

### 3. Minimize Dependências

```go
// ❌ Muitas dependências (God Object)
func NewService(
    repo1, repo2, repo3, repo4, repo5 Repository,
    svc1, svc2, svc3 Service,
    logger Logger,
    cache Cache,
    // ... 10 mais
) *Service

// ✅ Poucas dependências focadas
func NewService(repo Repository, logger Logger) *Service
```

### 4. Evite Service Locator

```go
// ❌ Service Locator
func (s *Service) DoSomething() {
    repo := s.container.Get("repository")
}

// ✅ Constructor Injection
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}
```

### 5. Documente Dependências

```go
// CreateUserHandler handles user creation.
// Dependencies:
//   - UserRepository: for persistence
//   - PasswordHasher: for password hashing
//   - EmailService: for welcome emails
//   - EventBus: for publishing domain events
//   - Logger: for logging
type CreateUserHandler struct {
    repo     UserRepository
    hasher   PasswordHasher
    emailSvc EmailService
    eventBus *events.EventBus
    logger   Logger
}
```

---

## 📚 Próximos Passos

- **[Sistema de Módulos](06-modules-system.md)** - Como módulos usam DI
- **[Adapters](09-adapters.md)** - Implementações concretas
- **[Boas Práticas](23-best-practices.md)** - Padrões recomendados

---

**[⬅️ Sistema de Módulos](06-modules-system.md)** | **[Índice](README.md)** | **[Sistema de Eventos ➡️](08-events-system.md)**
