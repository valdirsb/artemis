# 📡 Sistema de Eventos

## 📋 Índice
- [O que são Eventos?](#o-que-são-eventos)
- [Event Bus](#event-bus)
- [Tipos de Eventos](#tipos-de-eventos)
- [Publishers](#publishers)
- [Subscribers](#subscribers)
- [Exemplos Práticos](#exemplos-práticos)

---

## 🎯 O que são Eventos?

Eventos representam **algo que aconteceu** no sistema. Eles são usados para desacoplar módulos e implementar side effects de forma assíncrona.

### Características dos Eventos

✅ **Imutáveis** - Representam fatos do passado  
✅ **Nomeados no passado** - `UserCreated`, não `CreateUser`  
✅ **Assíncronos** - Não bloqueiam a operação principal  
✅ **Desacoplados** - Publisher não conhece subscribers  

### Benefícios

```
┌──────────────┐                    ┌──────────────┐
│ User Module  │                    │Email Service │
│              │    user.created    │              │
│ CreateUser() ├───────────────────►│ SendWelcome()│
│              │                    │              │
└──────────────┘                    └──────────────┘
                                    
                                    ┌──────────────┐
                                    │Audit Service │
                    user.created    │              │
                   ────────────────►│ LogCreation()│
                                    │              │
                                    └──────────────┘

User Module não sabe que Email e Audit existem!
```

---

## 🚌 Event Bus

O **Event Bus** é o componente central que gerencia publicação e subscrição de eventos.

### Estrutura Básica

```go
// pkg/events/event_bus.go
type EventBus struct {
    subscribers map[string][]EventHandler
    mu          sync.RWMutex
}

type EventHandler func(event Event) error

type Event struct {
    Name      string
    Data      interface{}
    OccurredAt time.Time
}
```

### Criação

```go
// No bootstrap
eventBus := events.NewEventBus()
container.Register("eventbus", eventBus)
```

### Operações Básicas

```go
// Publicar evento
eventBus.Publish("user.created", UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
})

// Subscrever a evento
eventBus.Subscribe("user.created", func(event Event) error {
    // Processar evento
    return nil
})
```

---

## 📬 Tipos de Eventos

### 1. Domain Events (Eventos de Domínio)

Representam mudanças no estado do domínio.

```go
// internal/modules/user/domain/events/user_created.go
type UserCreatedEvent struct {
    UserID    string
    Email     string
    Name      string
    OccurredAt time.Time
}

func NewUserCreatedEvent(user *User) *UserCreatedEvent {
    return &UserCreatedEvent{
        UserID:    user.ID,
        Email:     user.Email,
        Name:      user.Name,
        OccurredAt: time.Now(),
    }
}
```

**Exemplos:**
```go
// User Module
type UserCreatedEvent struct {...}
type UserUpdatedEvent struct {...}
type UserDeletedEvent struct {...}
type UserActivatedEvent struct {...}

// Order Module
type OrderCreatedEvent struct {...}
type OrderCancelledEvent struct {...}
type OrderCompletedEvent struct {...}
type OrderPaidEvent struct {...}

// Product Module
type ProductCreatedEvent struct {...}
type StockUpdatedEvent struct {...}
type PriceChangedEvent struct {...}
```

### 2. Integration Events (Eventos de Integração)

Para comunicação entre bounded contexts ou sistemas externos.

```go
type OrderPlacedIntegrationEvent struct {
    OrderID      string
    UserID       string
    TotalAmount  float64
    Items        []OrderItem
    PlacedAt     time.Time
    
    // Metadata para integração
    CorrelationID string
    Source        string
    Version       string
}
```

### 3. System Events (Eventos de Sistema)

Eventos relacionados à infraestrutura.

```go
type DatabaseConnectionLostEvent struct {...}
type CacheClearedEvent struct {...}
type ServiceStartedEvent struct {...}
```

---

## 📤 Publishers

### Publicar Evento Simples

```go
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    // 1. Criar usuário
    user := &User{...}
    h.repo.Save(user)
    
    // 2. Publicar evento
    h.eventBus.Publish("user.created", map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
        "name":    user.Name,
    })
    
    return nil
}
```

### Publicar com Evento Tipado

```go
// Melhor: Use structs tipadas
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    user := &User{...}
    h.repo.Save(user)
    
    event := events.UserCreatedEvent{
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
        OccurredAt: time.Now(),
    }
    
    h.eventBus.Publish("user.created", event)
    return nil
}
```

### TypedEventPublisher

```go
// pkg/events/typed_publisher.go
type TypedEventPublisher struct {
    eventBus *EventBus
}

func (p *TypedEventPublisher) PublishUserCreated(event *UserCreatedEvent) error {
    return p.eventBus.Publish("user.created", event)
}

// Uso
publisher.PublishUserCreated(&UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
})
```

### Publicar Múltiplos Eventos

```go
func (h *CompleteOrderHandler) Handle(cmd *CompleteOrderCommand) error {
    order := &Order{...}
    
    // Operação principal
    if err := h.repo.Update(order); err != nil {
        return err
    }
    
    // Publicar eventos
    events := []struct{
        name string
        data interface{}
    }{
        {"order.completed", OrderCompletedEvent{...}},
        {"inventory.reserved", InventoryReservedEvent{...}},
        {"payment.processed", PaymentProcessedEvent{...}},
    }
    
    for _, event := range events {
        h.eventBus.Publish(event.name, event.data)
    }
    
    return nil
}
```

---

## 📥 Subscribers

### Subscriber Básico

```go
// internal/modules/user/adapters/subscribers/email_subscriber.go
type UserEmailSubscriber struct {
    emailService EmailService
    logger       Logger
}

func NewUserEmailSubscriber(emailService EmailService, logger Logger) *UserEmailSubscriber {
    return &UserEmailSubscriber{
        emailService: emailService,
        logger:       logger,
    }
}

func (s *UserEmailSubscriber) Subscribe(eventBus *events.EventBus) error {
    return eventBus.Subscribe("user.created", s.handleUserCreated)
}

func (s *UserEmailSubscriber) handleUserCreated(event events.Event) error {
    data, ok := event.Data.(UserCreatedEvent)
    if !ok {
        return errors.New("invalid event data")
    }
    
    s.logger.Info("Sending welcome email", data.Email)
    
    if err := s.emailService.SendWelcome(data.Email, data.Name); err != nil {
        s.logger.Error("Failed to send email", err)
        return err
    }
    
    return nil
}
```

### Registrar Subscribers no Módulo

```go
// user_module.go
func (m *UserModule) Register(registry *container.ModuleRegistry) error {
    // ... outros registros
    
    // Criar e registrar subscribers
    emailSubscriber := subscribers.NewUserEmailSubscriber(emailService, logger)
    auditSubscriber := subscribers.NewUserAuditSubscriber(auditService, logger)
    
    // Registrar no event bus
    emailSubscriber.Subscribe(m.eventBus)
    auditSubscriber.Subscribe(m.eventBus)
    
    // Registrar no registry (opcional, para tracking)
    registry.RegisterEventSubscriber(emailSubscriber)
    registry.RegisterEventSubscriber(auditSubscriber)
    
    return nil
}
```

### Subscriber com Múltiplos Eventos

```go
type NotificationSubscriber struct {
    notificationService NotificationService
}

func (s *NotificationSubscriber) Subscribe(eventBus *events.EventBus) error {
    // Subscrever múltiplos eventos
    eventBus.Subscribe("user.created", s.onUserCreated)
    eventBus.Subscribe("order.completed", s.onOrderCompleted)
    eventBus.Subscribe("payment.received", s.onPaymentReceived)
    return nil
}

func (s *NotificationSubscriber) onUserCreated(event events.Event) error {
    return s.notificationService.Send("Welcome to our platform!")
}

func (s *NotificationSubscriber) onOrderCompleted(event events.Event) error {
    return s.notificationService.Send("Your order is complete!")
}

func (s *NotificationSubscriber) onPaymentReceived(event events.Event) error {
    return s.notificationService.Send("Payment received!")
}
```

---

## 💡 Exemplos Práticos

### Exemplo 1: Sistema de Auditoria

```go
// audit_subscriber.go
type AuditSubscriber struct {
    auditRepo AuditRepository
    logger    Logger
}

func (s *AuditSubscriber) Subscribe(eventBus *events.EventBus) error {
    // Auditar todas as mudanças
    eventBus.Subscribe("user.created", s.audit)
    eventBus.Subscribe("user.updated", s.audit)
    eventBus.Subscribe("user.deleted", s.audit)
    eventBus.Subscribe("order.created", s.audit)
    // ... outros eventos
    return nil
}

func (s *AuditSubscriber) audit(event events.Event) error {
    audit := &AuditLog{
        EventName:  event.Name,
        EventData:  event.Data,
        OccurredAt: event.OccurredAt,
        User:       extractUser(event),
    }
    
    if err := s.auditRepo.Save(audit); err != nil {
        s.logger.Error("Failed to audit event", err)
        return err
    }
    
    return nil
}
```

### Exemplo 2: Cache Invalidation

```go
type CacheInvalidationSubscriber struct {
    cache Cache
}

func (s *CacheInvalidationSubscriber) Subscribe(eventBus *events.EventBus) error {
    eventBus.Subscribe("user.updated", s.invalidateUserCache)
    eventBus.Subscribe("user.deleted", s.invalidateUserCache)
    return nil
}

func (s *CacheInvalidationSubscriber) invalidateUserCache(event events.Event) error {
    data := event.Data.(UserEvent)
    return s.cache.Delete(fmt.Sprintf("user:%s", data.UserID))
}
```

### Exemplo 3: Workflow Complexo

```go
// Quando um pedido é criado, precisa:
// 1. Reservar estoque
// 2. Processar pagamento
// 3. Enviar email
// 4. Notificar vendedor

type OrderWorkflowSubscriber struct {
    inventoryService  InventoryService
    paymentService    PaymentService
    emailService      EmailService
    notificationService NotificationService
}

func (s *OrderWorkflowSubscriber) Subscribe(eventBus *events.EventBus) error {
    return eventBus.Subscribe("order.created", s.handleOrderCreated)
}

func (s *OrderWorkflowSubscriber) handleOrderCreated(event events.Event) error {
    data := event.Data.(OrderCreatedEvent)
    
    // 1. Reservar estoque
    if err := s.inventoryService.Reserve(data.Items); err != nil {
        return fmt.Errorf("failed to reserve inventory: %w", err)
    }
    
    // 2. Processar pagamento
    if err := s.paymentService.Process(data.OrderID, data.Amount); err != nil {
        s.inventoryService.Release(data.Items) // Rollback
        return fmt.Errorf("failed to process payment: %w", err)
    }
    
    // 3. Enviar email (não falha se der erro)
    if err := s.emailService.SendOrderConfirmation(data.UserEmail); err != nil {
        log.Printf("Warning: failed to send email: %v", err)
    }
    
    // 4. Notificar vendedor
    if err := s.notificationService.NotifySeller(data.SellerID, data.OrderID); err != nil {
        log.Printf("Warning: failed to notify seller: %v", err)
    }
    
    return nil
}
```

### Exemplo 4: Saga Pattern

```go
// Orquestrar transação distribuída usando eventos
type OrderSagaSubscriber struct {
    orderService     OrderService
    paymentService   PaymentService
    inventoryService InventoryService
    eventBus         *events.EventBus
}

func (s *OrderSagaSubscriber) Subscribe(eventBus *events.EventBus) error {
    eventBus.Subscribe("order.created", s.onOrderCreated)
    eventBus.Subscribe("payment.completed", s.onPaymentCompleted)
    eventBus.Subscribe("payment.failed", s.onPaymentFailed)
    eventBus.Subscribe("inventory.reserved", s.onInventoryReserved)
    eventBus.Subscribe("inventory.failed", s.onInventoryFailed)
    return nil
}

func (s *OrderSagaSubscriber) onOrderCreated(event events.Event) error {
    data := event.Data.(OrderCreatedEvent)
    
    // Próximo passo: processar pagamento
    s.eventBus.Publish("process.payment", ProcessPaymentEvent{
        OrderID: data.OrderID,
        Amount:  data.Amount,
    })
    return nil
}

func (s *OrderSagaSubscriber) onPaymentCompleted(event events.Event) error {
    data := event.Data.(PaymentCompletedEvent)
    
    // Próximo passo: reservar estoque
    s.eventBus.Publish("reserve.inventory", ReserveInventoryEvent{
        OrderID: data.OrderID,
    })
    return nil
}

func (s *OrderSagaSubscriber) onPaymentFailed(event events.Event) error {
    data := event.Data.(PaymentFailedEvent)
    
    // Compensação: cancelar pedido
    return s.orderService.Cancel(data.OrderID, "payment failed")
}
```

---

## 🎯 Padrões e Melhores Práticas

### 1. Nomeação de Eventos

```go
✅ user.created       // Domínio.Ação no passado
✅ order.completed
✅ payment.processed

❌ createUser         // Não use imperativos
❌ UserCreate         // Não use presente
❌ user_created       // Use . não _
```

### 2. Eventos Devem Ser Imutáveis

```go
// ✅ Struct imutável
type UserCreatedEvent struct {
    userID    string  // private
    email     string
    occurredAt time.Time
}

func (e *UserCreatedEvent) UserID() string {
    return e.userID  // Getter, sem setter
}
```

### 3. Não Falhe a Operação por Causa de Eventos

```go
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    user := &User{...}
    
    // Operação principal - DEVE ter sucesso
    if err := h.repo.Save(user); err != nil {
        return err  // ✅ Retorna erro
    }
    
    // Eventos são side effects - NÃO devem falhar a operação
    if err := h.eventBus.Publish("user.created", event); err != nil {
        h.logger.Error("Failed to publish event", err)  // ✅ Apenas loga
        // NÃO retorna erro!
    }
    
    return nil
}
```

### 4. Use Eventos para Desacoplamento

```go
// ❌ Acoplado
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    user := &User{...}
    h.repo.Save(user)
    
    // Handler sabe sobre email, audit, notification...
    h.emailService.SendWelcome(user.Email)
    h.auditService.Log("user.created", user)
    h.notificationService.Notify(user)
}

// ✅ Desacoplado
func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    user := &User{...}
    h.repo.Save(user)
    
    // Handler só publica evento
    h.eventBus.Publish("user.created", UserCreatedEvent{...})
    
    // Outros serviços subscrevem
}
```

### 5. Idempotência de Subscribers

```go
// Subscribers devem ser idempotentes
func (s *EmailSubscriber) handleUserCreated(event events.Event) error {
    // Verificar se já processou
    if s.emailRepo.WasSent(event.UserID, "welcome") {
        return nil  // Já enviado, não faz nada
    }
    
    // Enviar email
    if err := s.emailService.SendWelcome(event.Email); err != nil {
        return err
    }
    
    // Marcar como enviado
    s.emailRepo.MarkAsSent(event.UserID, "welcome")
    return nil
}
```

---

## 🚀 Event Bus Avançado

### Com Retry

```go
type EventBusWithRetry struct {
    eventBus   *EventBus
    maxRetries int
}

func (eb *EventBusWithRetry) Publish(name string, data interface{}) error {
    for i := 0; i < eb.maxRetries; i++ {
        if err := eb.eventBus.Publish(name, data); err != nil {
            if i == eb.maxRetries-1 {
                return err
            }
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }
        return nil
    }
    return nil
}
```

### Com Persistência (Event Store)

```go
type PersistentEventBus struct {
    eventStore EventStore
    eventBus   *EventBus
}

func (peb *PersistentEventBus) Publish(name string, data interface{}) error {
    // 1. Persistir evento
    event := Event{
        ID:        generateID(),
        Name:      name,
        Data:      data,
        OccurredAt: time.Now(),
    }
    
    if err := peb.eventStore.Save(event); err != nil {
        return err
    }
    
    // 2. Publicar no bus
    return peb.eventBus.Publish(name, data)
}
```

---

## 📚 Próximos Passos

- **[Implementando Commands](13-implementing-commands.md)**
- **[Trabalhando com Eventos - Exemplos](15-working-with-events.md)**
- **[Boas Práticas](23-best-practices.md)**

---

**[⬅️ Índice](README.md)** | **[Boas Práticas ➡️](23-best-practices.md)**
