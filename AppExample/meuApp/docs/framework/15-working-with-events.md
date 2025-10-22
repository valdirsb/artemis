# 📡 Trabalhando com Eventos - Guia Prático

## 📋 Índice
- [Introdução](#introdução)
- [Cenário 1: Auditoria](#cenário-1-auditoria)
- [Cenário 2: Notificações](#cenário-2-notificações)
- [Cenário 3: Cache Invalidation](#cenário-3-cache-invalidation)
- [Cenário 4: Workflow Complexo](#cenário-4-workflow-complexo)
- [Cenário 5: Saga Pattern](#cenário-5-saga-pattern)
- [Cenário 6: Event Sourcing](#cenário-6-event-sourcing)
- [Testes](#testes)
- [Troubleshooting](#troubleshooting)

---

## 🎯 Introdução

Este guia apresenta **exemplos práticos** de uso de eventos no Artemis Framework, desde casos simples até padrões avançados.

### Fluxo Básico

```
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│   Command    │─────▶│  Publisher   │─────▶│  Event Bus   │
│   Handler    │      │  (emite)     │      │              │
└──────────────┘      └──────────────┘      └──────┬───────┘
                                                    │
                                                    ├─────────┐
                                                    │         │
                                            ┌───────▼──┐  ┌──▼──────┐
                                            │Subscriber│  │Subscriber│
                                            │    1     │  │    2    │
                                            └──────────┘  └─────────┘
```

---

## 📝 Cenário 1: Auditoria

**Objetivo:** Registrar todas as ações de usuários para auditoria.

### 1. Definir Evento

**`internal/modules/user/domain/events/user_action_event.go`**

```go
package events

import "time"

type UserActionEvent struct {
    ActionID   string
    UserID     string
    Action     string // "CREATE", "UPDATE", "DELETE"
    EntityType string // "USER", "PRODUCT", "ORDER"
    EntityID   string
    Details    map[string]interface{}
    IPAddress  string
    UserAgent  string
    occurredAt time.Time
}

func NewUserActionEvent(
    userID, action, entityType, entityID string,
    details map[string]interface{},
    ipAddress, userAgent string,
) *UserActionEvent {
    return &UserActionEvent{
        ActionID:   uuid.New().String(),
        UserID:     userID,
        Action:     action,
        EntityType: entityType,
        EntityID:   entityID,
        Details:    details,
        IPAddress:  ipAddress,
        UserAgent:  userAgent,
        occurredAt: time.Now(),
    }
}

func (e *UserActionEvent) EventName() string {
    return "UserAction"
}

func (e *UserActionEvent) OccurredAt() time.Time {
    return e.occurredAt
}
```

### 2. Publicar no Command Handler

**`internal/modules/user/application/commands/create_user_handler.go`**

```go
package commands

type CreateUserHandler struct {
    repo     UserRepository
    eventBus contracts.EventBus
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd contracts.Command) (interface{}, error) {
    createCmd := cmd.(*CreateUserCommand)
    
    // 1. Criar usuário
    user := &entities.User{
        ID:       uuid.New().String(),
        Name:     createCmd.Name,
        Email:    createCmd.Email,
        Password: hashPassword(createCmd.Password),
    }
    
    if err := h.repo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 2. Publicar evento de auditoria
    event := events.NewUserActionEvent(
        createCmd.ActorID,          // Quem executou
        "CREATE",                    // Ação
        "USER",                      // Tipo de entidade
        user.ID,                     // ID da entidade
        map[string]interface{}{      // Detalhes
            "name":  user.Name,
            "email": user.Email,
        },
        createCmd.IPAddress,         // IP do requisitante
        createCmd.UserAgent,         // User Agent
    )
    
    h.eventBus.Publish(ctx, event)
    
    return user.ID, nil
}
```

### 3. Criar Subscriber de Auditoria

**`internal/modules/audit/adapters/subscribers/audit_subscriber.go`**

```go
package subscribers

import (
    "context"
    "meuApp/internal/modules/audit/domain/entities"
    "meuApp/internal/modules/audit/repository"
    "meuApp/pkg/contracts"
)

type AuditSubscriber struct {
    auditRepo repository.AuditLogRepository
    logger    contracts.Logger
}

func NewAuditSubscriber(
    auditRepo repository.AuditLogRepository,
    logger contracts.Logger,
) *AuditSubscriber {
    return &AuditSubscriber{
        auditRepo: auditRepo,
        logger:    logger,
    }
}

func (s *AuditSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    userAction := event.(*events.UserActionEvent)
    
    // Criar registro de auditoria
    auditLog := &entities.AuditLog{
        ID:         uuid.New().String(),
        UserID:     userAction.UserID,
        Action:     userAction.Action,
        EntityType: userAction.EntityType,
        EntityID:   userAction.EntityID,
        Details:    userAction.Details,
        IPAddress:  userAction.IPAddress,
        UserAgent:  userAction.UserAgent,
        CreatedAt:  userAction.OccurredAt(),
    }
    
    // Persistir
    if err := s.auditRepo.Save(ctx, auditLog); err != nil {
        s.logger.Error("Failed to save audit log", map[string]interface{}{
            "error":     err.Error(),
            "action_id": userAction.ActionID,
        })
        return err
    }
    
    s.logger.Info("Audit log saved", map[string]interface{}{
        "action_id": userAction.ActionID,
        "user_id":   userAction.UserID,
        "action":    userAction.Action,
    })
    
    return nil
}

func (s *AuditSubscriber) SubscribedTo() []string {
    return []string{"UserAction"}
}
```

### 4. Registrar Subscriber

**`internal/modules/audit/audit_module.go`**

```go
func (m *AuditModule) RegisterEventSubscribers(eventBus contracts.EventBus) error {
    subscriber := subscribers.NewAuditSubscriber(
        m.auditRepo,
        m.logger,
    )
    
    return eventBus.Subscribe(subscriber)
}
```

---

## 📧 Cenário 2: Notificações

**Objetivo:** Enviar emails/SMS quando eventos importantes acontecerem.

### 1. Evento de Usuário Criado

**`internal/modules/user/domain/events/user_created_event.go`**

```go
package events

type UserCreatedEvent struct {
    UserID    string
    Name      string
    Email     string
    occurredAt time.Time
}

func NewUserCreatedEvent(userID, name, email string) *UserCreatedEvent {
    return &UserCreatedEvent{
        UserID:     userID,
        Name:       name,
        Email:      email,
        occurredAt: time.Now(),
    }
}

func (e *UserCreatedEvent) EventName() string {
    return "UserCreated"
}

func (e *UserCreatedEvent) OccurredAt() time.Time {
    return e.occurredAt
}
```

### 2. Múltiplos Subscribers

**A) Welcome Email Subscriber:**

```go
package subscribers

type SendWelcomeEmailSubscriber struct {
    emailService EmailService
    logger       contracts.Logger
}

func (s *SendWelcomeEmailSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    userCreated := event.(*events.UserCreatedEvent)
    
    // Template de email
    emailBody := fmt.Sprintf(`
        <h1>Bem-vindo, %s!</h1>
        <p>Obrigado por se cadastrar em nossa plataforma.</p>
        <p>Seu email de cadastro é: %s</p>
    `, userCreated.Name, userCreated.Email)
    
    // Enviar email
    err := s.emailService.Send(
        userCreated.Email,
        "Bem-vindo à plataforma!",
        emailBody,
    )
    
    if err != nil {
        s.logger.Error("Failed to send welcome email", map[string]interface{}{
            "user_id": userCreated.UserID,
            "error":   err.Error(),
        })
        return err
    }
    
    s.logger.Info("Welcome email sent", map[string]interface{}{
        "user_id": userCreated.UserID,
        "email":   userCreated.Email,
    })
    
    return nil
}

func (s *SendWelcomeEmailSubscriber) SubscribedTo() []string {
    return []string{"UserCreated"}
}
```

**B) Admin Notification Subscriber:**

```go
type NotifyAdminSubscriber struct {
    emailService EmailService
    adminEmail   string
}

func (s *NotifyAdminSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    userCreated := event.(*events.UserCreatedEvent)
    
    // Notificar admin
    emailBody := fmt.Sprintf(`
        Novo usuário cadastrado:
        - ID: %s
        - Nome: %s
        - Email: %s
        - Data: %s
    `,
        userCreated.UserID,
        userCreated.Name,
        userCreated.Email,
        userCreated.OccurredAt().Format("02/01/2006 15:04"),
    )
    
    return s.emailService.Send(
        s.adminEmail,
        "Novo usuário cadastrado",
        emailBody,
    )
}

func (s *NotifyAdminSubscriber) SubscribedTo() []string {
    return []string{"UserCreated"}
}
```

**C) Analytics Subscriber:**

```go
type TrackUserSignupSubscriber struct {
    analyticsService AnalyticsService
}

func (s *TrackUserSignupSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    userCreated := event.(*events.UserCreatedEvent)
    
    // Enviar para analytics (Google Analytics, Mixpanel, etc.)
    return s.analyticsService.Track("user_signup", map[string]interface{}{
        "user_id": userCreated.UserID,
        "email":   userCreated.Email,
        "name":    userCreated.Name,
    })
}

func (s *TrackUserSignupSubscriber) SubscribedTo() []string {
    return []string{"UserCreated"}
}
```

---

## 🗑️ Cenário 3: Cache Invalidation

**Objetivo:** Invalidar cache quando dados forem alterados.

### 1. Evento de Produto Atualizado

**`internal/modules/product/domain/events/product_updated_event.go`**

```go
package events

type ProductUpdatedEvent struct {
    ProductID   string
    Name        string
    Price       float64
    UpdatedBy   string
    occurredAt  time.Time
}

func NewProductUpdatedEvent(productID, name string, price float64, updatedBy string) *ProductUpdatedEvent {
    return &ProductUpdatedEvent{
        ProductID:  productID,
        Name:       name,
        Price:      price,
        UpdatedBy:  updatedBy,
        occurredAt: time.Now(),
    }
}

func (e *ProductUpdatedEvent) EventName() string {
    return "ProductUpdated"
}

func (e *ProductUpdatedEvent) OccurredAt() time.Time {
    return e.occurredAt
}
```

### 2. Cache Invalidation Subscriber

```go
package subscribers

type CacheInvalidationSubscriber struct {
    cache  contracts.Cache
    logger contracts.Logger
}

func NewCacheInvalidationSubscriber(cache contracts.Cache, logger contracts.Logger) *CacheInvalidationSubscriber {
    return &CacheInvalidationSubscriber{
        cache:  cache,
        logger: logger,
    }
}

func (s *CacheInvalidationSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    productUpdated := event.(*events.ProductUpdatedEvent)
    
    // Invalidar cache do produto específico
    productKey := fmt.Sprintf("product:%s", productUpdated.ProductID)
    s.cache.Delete(productKey)
    
    // Invalidar cache de listagem
    s.cache.Delete("products:list")
    
    s.logger.Info("Cache invalidated", map[string]interface{}{
        "product_id": productUpdated.ProductID,
        "keys":       []string{productKey, "products:list"},
    })
    
    return nil
}

func (s *CacheInvalidationSubscriber) SubscribedTo() []string {
    return []string{
        "ProductCreated",
        "ProductUpdated",
        "ProductDeleted",
    }
}
```

### 3. Query Handler com Cache

```go
type GetProductByIDHandler struct {
    repo  ProductRepository
    cache contracts.Cache
}

func (h *GetProductByIDHandler) Handle(ctx context.Context, qry contracts.Query) (interface{}, error) {
    query := qry.(*GetProductByIDQuery)
    
    // 1. Tentar cache primeiro
    cacheKey := fmt.Sprintf("product:%s", query.ProductID)
    if cached, found := h.cache.Get(cacheKey); found {
        return cached.(*ProductDTO), nil
    }
    
    // 2. Buscar do banco
    product, err := h.repo.FindByID(ctx, query.ProductID)
    if err != nil {
        return nil, err
    }
    
    // 3. Converter para DTO
    dto := &ProductDTO{
        ID:    product.ID,
        Name:  product.Name,
        Price: product.Price,
    }
    
    // 4. Armazenar em cache (5 minutos)
    h.cache.Set(cacheKey, dto, 5*time.Minute)
    
    return dto, nil
}
```

---

## 🔄 Cenário 4: Workflow Complexo

**Objetivo:** Orquestrar múltiplas operações quando um pedido é criado.

### 1. Evento de Pedido Criado

```go
package events

type OrderCreatedEvent struct {
    OrderID    string
    UserID     string
    Items      []OrderItem
    TotalPrice float64
    occurredAt time.Time
}

func (e *OrderCreatedEvent) EventName() string {
    return "OrderCreated"
}
```

### 2. Múltiplos Subscribers Orquestrados

**A) Reservar Estoque:**

```go
type ReserveStockSubscriber struct {
    inventoryService InventoryService
    eventBus         contracts.EventBus
}

func (s *ReserveStockSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    orderCreated := event.(*events.OrderCreatedEvent)
    
    // Reservar estoque para cada item
    for _, item := range orderCreated.Items {
        err := s.inventoryService.ReserveStock(ctx, item.ProductID, item.Quantity)
        if err != nil {
            // Publicar evento de falha
            s.eventBus.Publish(ctx, &events.OrderStockReservationFailedEvent{
                OrderID:   orderCreated.OrderID,
                ProductID: item.ProductID,
                Reason:    err.Error(),
            })
            return err
        }
    }
    
    // Publicar sucesso
    s.eventBus.Publish(ctx, &events.OrderStockReservedEvent{
        OrderID: orderCreated.OrderID,
    })
    
    return nil
}

func (s *ReserveStockSubscriber) SubscribedTo() []string {
    return []string{"OrderCreated"}
}
```

**B) Processar Pagamento:**

```go
type ProcessPaymentSubscriber struct {
    paymentService PaymentService
    eventBus       contracts.EventBus
}

func (s *ProcessPaymentSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    // Aguardar estoque ser reservado
    stockReserved := event.(*events.OrderStockReservedEvent)
    
    // Processar pagamento
    paymentID, err := s.paymentService.ProcessPayment(ctx, stockReserved.OrderID)
    if err != nil {
        s.eventBus.Publish(ctx, &events.OrderPaymentFailedEvent{
            OrderID: stockReserved.OrderID,
            Reason:  err.Error(),
        })
        return err
    }
    
    // Publicar sucesso
    s.eventBus.Publish(ctx, &events.OrderPaymentProcessedEvent{
        OrderID:   stockReserved.OrderID,
        PaymentID: paymentID,
    })
    
    return nil
}

func (s *ProcessPaymentSubscriber) SubscribedTo() []string {
    return []string{"OrderStockReserved"}
}
```

**C) Enviar para Logística:**

```go
type SendToShippingSubscriber struct {
    shippingService ShippingService
    eventBus        contracts.EventBus
}

func (s *SendToShippingSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    paymentProcessed := event.(*events.OrderPaymentProcessedEvent)
    
    // Criar pedido de envio
    trackingCode, err := s.shippingService.CreateShipment(ctx, paymentProcessed.OrderID)
    if err != nil {
        return err
    }
    
    // Publicar evento de envio criado
    s.eventBus.Publish(ctx, &events.OrderShipmentCreatedEvent{
        OrderID:      paymentProcessed.OrderID,
        TrackingCode: trackingCode,
    })
    
    return nil
}

func (s *SendToShippingSubscriber) SubscribedTo() []string {
    return []string{"OrderPaymentProcessed"}
}
```

**D) Notificar Cliente:**

```go
type NotifyCustomerSubscriber struct {
    emailService EmailService
}

func (s *NotifyCustomerSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    shipmentCreated := event.(*events.OrderShipmentCreatedEvent)
    
    // Buscar dados do pedido
    order, _ := s.orderRepo.FindByID(ctx, shipmentCreated.OrderID)
    user, _ := s.userRepo.FindByID(ctx, order.UserID)
    
    // Enviar email
    emailBody := fmt.Sprintf(`
        Pedido #%s enviado!
        Código de rastreamento: %s
        Acompanhe em: https://rastreio.com/%s
    `, order.ID, shipmentCreated.TrackingCode, shipmentCreated.TrackingCode)
    
    return s.emailService.Send(user.Email, "Pedido enviado!", emailBody)
}

func (s *NotifyCustomerSubscriber) SubscribedTo() []string {
    return []string{"OrderShipmentCreated"}
}
```

### Fluxo Completo

```
OrderCreated
    │
    ├──▶ ReserveStock
    │        │
    │        ▼
    │    OrderStockReserved
    │        │
    │        ▼
    │    ProcessPayment
    │        │
    │        ▼
    │    OrderPaymentProcessed
    │        │
    │        ▼
    │    SendToShipping
    │        │
    │        ▼
    │    OrderShipmentCreated
    │        │
    │        ▼
    └───▶ NotifyCustomer
```

---

## 🎭 Cenário 5: Saga Pattern

**Objetivo:** Implementar transações distribuídas com compensação.

### Saga Orchestrator

```go
package sagas

type OrderSaga struct {
    eventBus         contracts.EventBus
    orderRepo        OrderRepository
    inventoryService InventoryService
    paymentService   PaymentService
    logger           contracts.Logger
}

func (s *OrderSaga) Handle(ctx context.Context, event contracts.Event) error {
    switch e := event.(type) {
    case *events.OrderCreatedEvent:
        return s.handleOrderCreated(ctx, e)
    case *events.OrderStockReservationFailedEvent:
        return s.compensateOrderCreation(ctx, e)
    case *events.OrderPaymentFailedEvent:
        return s.compensateStockReservation(ctx, e)
    }
    return nil
}

func (s *OrderSaga) handleOrderCreated(ctx context.Context, event *events.OrderCreatedEvent) error {
    // Step 1: Reservar estoque
    for _, item := range event.Items {
        err := s.inventoryService.ReserveStock(ctx, item.ProductID, item.Quantity)
        if err != nil {
            // Falhou - publicar evento de compensação
            s.eventBus.Publish(ctx, &events.OrderStockReservationFailedEvent{
                OrderID: event.OrderID,
                Reason:  err.Error(),
            })
            return err
        }
    }
    
    // Step 2: Processar pagamento
    _, err := s.paymentService.ProcessPayment(ctx, event.OrderID, event.TotalPrice)
    if err != nil {
        // Falhou - publicar evento de compensação
        s.eventBus.Publish(ctx, &events.OrderPaymentFailedEvent{
            OrderID: event.OrderID,
            Reason:  err.Error(),
        })
        return err
    }
    
    // Sucesso - atualizar status do pedido
    order, _ := s.orderRepo.FindByID(ctx, event.OrderID)
    order.Status = "CONFIRMED"
    s.orderRepo.Update(ctx, order)
    
    return nil
}

func (s *OrderSaga) compensateOrderCreation(ctx context.Context, event *events.OrderStockReservationFailedEvent) error {
    s.logger.Warn("Compensating order creation", map[string]interface{}{
        "order_id": event.OrderID,
        "reason":   event.Reason,
    })
    
    // Cancelar pedido
    order, _ := s.orderRepo.FindByID(ctx, event.OrderID)
    order.Status = "CANCELLED"
    order.CancellationReason = "Stock reservation failed: " + event.Reason
    
    return s.orderRepo.Update(ctx, order)
}

func (s *OrderSaga) compensateStockReservation(ctx context.Context, event *events.OrderPaymentFailedEvent) error {
    s.logger.Warn("Compensating stock reservation", map[string]interface{}{
        "order_id": event.OrderID,
        "reason":   event.Reason,
    })
    
    // Liberar estoque reservado
    order, _ := s.orderRepo.FindByID(ctx, event.OrderID)
    for _, item := range order.Items {
        s.inventoryService.ReleaseStock(ctx, item.ProductID, item.Quantity)
    }
    
    // Cancelar pedido
    order.Status = "CANCELLED"
    order.CancellationReason = "Payment failed: " + event.Reason
    
    return s.orderRepo.Update(ctx, order)
}

func (s *OrderSaga) SubscribedTo() []string {
    return []string{
        "OrderCreated",
        "OrderStockReservationFailed",
        "OrderPaymentFailed",
    }
}
```

---

## 🧪 Testes

### 1. Teste Unitário de Subscriber

```go
func TestAuditSubscriber_Handle(t *testing.T) {
    // Arrange
    mockRepo := new(MockAuditRepository)
    mockLogger := new(MockLogger)
    
    subscriber := subscribers.NewAuditSubscriber(mockRepo, mockLogger)
    
    event := &events.UserActionEvent{
        ActionID:   "action-1",
        UserID:     "user-1",
        Action:     "CREATE",
        EntityType: "USER",
        EntityID:   "new-user-1",
    }
    
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    mockLogger.On("Info", mock.Anything, mock.Anything).Return()
    
    // Act
    err := subscriber.Handle(context.Background(), event)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
    mockLogger.AssertExpectations(t)
}
```

### 2. Teste de Integração de Eventos

```go
func TestEventFlow_UserCreation(t *testing.T) {
    // Setup
    eventBus := framework.NewEventBus()
    
    emailsSent := []string{}
    emailSubscriber := &MockEmailSubscriber{
        HandleFunc: func(ctx context.Context, event contracts.Event) error {
            userCreated := event.(*events.UserCreatedEvent)
            emailsSent = append(emailsSent, userCreated.Email)
            return nil
        },
    }
    
    eventBus.Subscribe(emailSubscriber)
    
    // Test
    event := &events.UserCreatedEvent{
        UserID: "user-1",
        Email:  "test@example.com",
    }
    
    err := eventBus.Publish(context.Background(), event)
    
    // Assert
    assert.NoError(t, err)
    assert.Len(t, emailsSent, 1)
    assert.Equal(t, "test@example.com", emailsSent[0])
}
```

---

## 🔍 Troubleshooting

### Problema 1: Eventos não são processados

**Sintoma:**
```
Publico evento, mas subscriber não é chamado
```

**Solução:**
```go
// Verificar se subscriber foi registrado
eventBus.Subscribe(mySubscriber)

// Verificar se EventName() corresponde
func (e *MyEvent) EventName() string {
    return "MyEvent"  // Deve ser exato
}

func (s *MySubscriber) SubscribedTo() []string {
    return []string{"MyEvent"}  // Deve ser igual
}
```

### Problema 2: Deadlock

**Sintoma:**
```
Application trava ao publicar evento
```

**Solução:**
```go
// ❌ NÃO faça isso (subscriber publica mesmo evento)
func (s *MySubscriber) Handle(ctx context.Context, event contracts.Event) error {
    // ...
    s.eventBus.Publish(ctx, event)  // ❌ Loop infinito!
}

// ✅ Publique evento diferente ou use flag
func (s *MySubscriber) Handle(ctx context.Context, event contracts.Event) error {
    // Publicar evento diferente
    s.eventBus.Publish(ctx, &DifferentEvent{...})
}
```

### Problema 3: Erros não são tratados

**Sintoma:**
```
Subscriber falha mas aplicação continua
```

**Solução:**
```go
// Implementar error handling e retry
func (s *MySubscriber) Handle(ctx context.Context, event contracts.Event) error {
    maxRetries := 3
    
    for i := 0; i < maxRetries; i++ {
        err := s.doWork(ctx, event)
        if err == nil {
            return nil
        }
        
        s.logger.Warn("Retry", map[string]interface{}{
            "attempt": i + 1,
            "error":   err.Error(),
        })
        
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    // Log e retornar erro após todas as tentativas
    s.logger.Error("Failed after retries", map[string]interface{}{
        "event": event.EventName(),
    })
    
    return errors.New("max retries exceeded")
}
```

---

## 📚 Próximos Passos

- **[Sistema de Eventos](08-events-system.md)** - Conceitos fundamentais
- **[Boas Práticas](23-best-practices.md)** - Padrões recomendados
- **[Criando Módulos](12-creating-modules.md)** - Como integrar eventos

---

**[⬅️ Queries](14-implementing-queries.md)** | **[Índice](README.md)** | **[Configuração ➡️](16-configuration-bootstrap.md)**
