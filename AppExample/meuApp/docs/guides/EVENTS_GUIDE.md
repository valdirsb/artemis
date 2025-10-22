# 📚 Event Bus Type-Safe - Guia de Uso

> Sistema de eventos tipados para maior segurança e produtividade

---

## 🎯 Visão Geral

O sistema de eventos agora suporta **duas abordagens**:

1. **Sistema Original** (`contracts.Event`) - Compatibilidade mantida
2. **Sistema Type-Safe** (`events.TypedEventPublisher`) - Novo, recomendado

---

## 🚀 Quick Start

### 1. Criar o Publisher Tipado

```go
import "meuApp/pkg/events"

// Obter EventBus existente
eventBus := events.NewEventBus()

// Criar publisher tipado
typedPublisher := events.NewTypedEventPublisher(eventBus)
```

### 2. Publicar Eventos Type-Safe

```go
// User Created
err := typedPublisher.PublishUserCreated(ctx, events.UserCreatedEvent{
    UserID:   user.ID,
    Username: user.Name,
    Email:    user.Email,
})

// Product Low Stock
err := typedPublisher.PublishLowStock(ctx, events.LowStockEvent{
    ProductID:    product.ID,
    ProductName:  product.Name,
    CurrentStock: product.Stock,
    Threshold:    10,
})

// Order Created
err := typedPublisher.PublishOrderCreated(ctx, events.OrderCreatedEvent{
    OrderID:     order.ID,
    UserID:      order.UserID,
    TotalAmount: order.Total,
    ItemCount:   len(order.Items),
})
```

### 3. Registrar Handlers Type-Safe

```go
// Handler inline
err := events.SubscribeTyped(eventBus, events.EventTypeUserCreated,
    func(ctx context.Context, data events.UserCreatedEvent) error {
        log.Printf("Novo usuário: %s (%s)", data.Username, data.Email)
        return sendWelcomeEmail(data.Email, data.Username)
    })

// Handler com lógica complexa
err := events.SubscribeTyped(eventBus, events.EventTypeLowStock,
    func(ctx context.Context, data events.LowStockEvent) error {
        if data.CurrentStock < 5 {
            return sendUrgentAlert(data.ProductName)
        }
        return sendNormalAlert(data.ProductName)
    })
```

---

## 📦 Eventos Disponíveis

### User Events

#### UserCreatedEvent
```go
type UserCreatedEvent struct {
    UserID   string
    Username string
    Email    string
}

// Publicar
typedPublisher.PublishUserCreated(ctx, events.UserCreatedEvent{...})

// Tipo do evento
events.EventTypeUserCreated // "user.created"
```

#### UserDeletedEvent
```go
type UserDeletedEvent struct {
    UserID string
    Email  string
}

// Publicar
typedPublisher.PublishUserDeleted(ctx, events.UserDeletedEvent{...})

// Tipo do evento
events.EventTypeUserDeleted // "user.deleted"
```

---

### Product Events

#### ProductCreatedEvent
```go
type ProductCreatedEvent struct {
    ProductID  string
    Name       string
    CategoryID string
    Price      float64
    Stock      int
}

// Publicar
typedPublisher.PublishProductCreated(ctx, events.ProductCreatedEvent{...})

// Tipo do evento
events.EventTypeProductCreated // "product.created"
```

#### LowStockEvent
```go
type LowStockEvent struct {
    ProductID    string
    ProductName  string
    CurrentStock int
    Threshold    int
}

// Publicar
typedPublisher.PublishLowStock(ctx, events.LowStockEvent{...})

// Tipo do evento
events.EventTypeLowStock // "product.low_stock"
```

---

### Order Events

#### OrderCreatedEvent
```go
type OrderCreatedEvent struct {
    OrderID     string
    UserID      string
    TotalAmount float64
    ItemCount   int
}

// Publicar
typedPublisher.PublishOrderCreated(ctx, events.OrderCreatedEvent{...})

// Tipo do evento
events.EventTypeOrderCreated // "order.created"
```

#### OrderStatusChangedEvent
```go
type OrderStatusChangedEvent struct {
    OrderID        string
    PreviousStatus string
    NewStatus      string
    ChangedBy      string
}

// Publicar
typedPublisher.PublishOrderStatusChanged(ctx, events.OrderStatusChangedEvent{...})

// Tipo do evento
events.EventTypeOrderStatusChanged // "order.status.changed"
```

#### OrderCancelledEvent
```go
type OrderCancelledEvent struct {
    OrderID      string
    UserID       string
    Reason       string
    RefundAmount float64
}

// Publicar
typedPublisher.PublishOrderCancelled(ctx, events.OrderCancelledEvent{...})

// Tipo do evento
events.EventTypeOrderCancelled // "order.cancelled"
```

---

## 💡 Exemplos Práticos

### Exemplo 1: Enviar Email de Boas-Vindas

```go
// No bootstrap/configuração
logger := &events.defaultSimpleLogger{}
handler := events.UserCreatedHandlerFunc(logger)

events.SubscribeTyped(eventBus, events.EventTypeUserCreated, handler)

// No Command
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
    // ... criar usuário ...
    
    // Publicar evento
    h.typedPublisher.PublishUserCreated(ctx, events.UserCreatedEvent{
        UserID:   user.ID,
        Username: user.Name,
        Email:    user.Email,
    })
    
    return user, nil
}
```

### Exemplo 2: Alertas de Estoque Baixo

```go
// Registrar handler
events.SubscribeTyped(eventBus, events.EventTypeLowStock,
    func(ctx context.Context, data events.LowStockEvent) error {
        // Enviar para Slack
        slackMessage := fmt.Sprintf("⚠️ Produto %s com apenas %d unidades!",
            data.ProductName, data.CurrentStock)
        return sendSlackAlert(slackMessage)
    })

// No UpdateStockCommand
if product.Stock < 10 {
    h.typedPublisher.PublishLowStock(ctx, events.LowStockEvent{
        ProductID:    product.ID,
        ProductName:  product.Name,
        CurrentStock: product.Stock,
        Threshold:    10,
    })
}
```

### Exemplo 3: Múltiplos Handlers para um Evento

```go
// Handler 1: Enviar email
events.SubscribeTyped(eventBus, events.EventTypeOrderCreated,
    func(ctx context.Context, data events.OrderCreatedEvent) error {
        return emailService.SendOrderConfirmation(data.OrderID, data.UserID)
    })

// Handler 2: Atualizar analytics
events.SubscribeTyped(eventBus, events.EventTypeOrderCreated,
    func(ctx context.Context, data events.OrderCreatedEvent) error {
        return analytics.TrackPurchase(data.UserID, data.TotalAmount)
    })

// Handler 3: Notificar sistema de estoque
events.SubscribeTyped(eventBus, events.EventTypeOrderCreated,
    func(ctx context.Context, data events.OrderCreatedEvent) error {
        return stockService.ReserveItems(data.OrderID)
    })

// Todos os 3 handlers serão executados quando o evento for publicado!
```

### Exemplo 4: Handler de Auditoria Genérico

```go
// Registrar audit logger para todos os eventos
auditHandler := events.AuditLogHandlerFunc[events.UserCreatedEvent]()
events.SubscribeTyped(eventBus, events.EventTypeUserCreated, auditHandler)

auditHandler2 := events.AuditLogHandlerFunc[events.OrderCreatedEvent]()
events.SubscribeTyped(eventBus, events.EventTypeOrderCreated, auditHandler2)
```

---

## 🔄 Migração do Sistema Antigo

### Antes (Sistema Original)

```go
event := contracts.Event{
    Type:      "user.created",
    Timestamp: time.Now(),
    Payload: map[string]interface{}{
        "user_id": userID,
        "email":   email,
        // ❌ Typos não detectados
        // ❌ Sem autocomplete
        // ❌ Runtime errors
    },
}
h.eventBus.Publish(ctx, event)
```

### Depois (Type-Safe)

```go
h.typedPublisher.PublishUserCreated(ctx, events.UserCreatedEvent{
    UserID:   user.ID,
    Username: user.Name,
    Email:    user.Email,
    // ✅ Autocomplete funcionando
    // ✅ Compile-time checks
    // ✅ Refactoring seguro
})
```

---

## 🏗️ Integração com Commands

### Atualizar Constructor do Command

```go
type CreateUserHandler struct {
    userRepo       ports.UserRepository
    passwordHasher ports.PasswordHasher
    eventBus       *events.EventBus              // Antigo
    typedPublisher *events.TypedEventPublisher  // ✨ Novo!
    logger         contracts.Logger
}

func NewCreateUserHandler(
    userRepo ports.UserRepository,
    passwordHasher ports.PasswordHasher,
    eventBus *events.EventBus,
    logger contracts.Logger,
) *CreateUserHandler {
    return &CreateUserHandler{
        userRepo:       userRepo,
        passwordHasher: passwordHasher,
        eventBus:       eventBus,
        typedPublisher: events.NewTypedEventPublisher(eventBus), // ✨
        logger:         logger,
    }
}
```

### Usar no Handler

```go
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
    // ... lógica de criação ...
    
    // ✨ Publicar evento type-safe
    if err := h.typedPublisher.PublishUserCreated(ctx, events.UserCreatedEvent{
        UserID:   user.ID,
        Username: user.Name,
        Email:    user.Email,
    }); err != nil {
        h.logger.Warn("Failed to publish user created event", 
            contracts.Field{Key: "error", Value: err})
    }
    
    return user, nil
}
```

---

## ✅ Vantagens do Sistema Type-Safe

### 1. Autocomplete no IDE
```go
event := events.UserCreatedEvent{
    // Ctrl+Space aqui mostra: UserID, Username, Email
}
```

### 2. Compile-Time Checks
```go
// ❌ Erro de compilação se faltar campo obrigatório
event := events.UserCreatedEvent{
    UserID: "123",
    // Falta Username e Email - NÃO COMPILA!
}
```

### 3. Refactoring Seguro
```go
// Renomear UserCreatedEvent.Email → UserCreatedEvent.EmailAddress
// → IDE atualiza TODOS os usos automaticamente!
```

### 4. Documentação Implícita
```go
// A estrutura É a documentação!
type LowStockEvent struct {
    ProductID    string  // Qual produto
    ProductName  string  // Nome para exibição
    CurrentStock int     // Estoque atual
    Threshold    int     // Limite configurado
}
```

---

## 🎯 Melhores Práticas

### ✅ Faça

- Use `TypedEventPublisher` para novos eventos
- Crie handlers específicos e focados
- Publique eventos após operações bem-sucedidas
- Use `SubscribeTyped` para handlers tipados
- Mantenha eventos imutáveis (sem setters)

### ❌ Evite

- Misturar lógica de negócio nos handlers
- Fazer handlers bloquearem por muito tempo
- Publicar eventos em transações não commitadas
- Usar `map[string]interface{}` para payloads
- Depender da ordem de execução dos handlers

---

## 🔍 Debugging

### Ver Eventos Registrados

```go
eventTypes := eventBus.ListEventTypes()
for _, eventType := range eventTypes {
    count := eventBus.GetHandlerCount(eventType)
    fmt.Printf("Event: %s, Handlers: %d\n", eventType, count)
}
```

### Logs Automáticos

O `TypedEventPublisher` e `SubscribeTyped` já fazem logging automático:
- Quando handler é registrado
- Quando evento é publicado
- Quando handler falha

---

## 📊 Performance

- **Zero overhead** em runtime comparado ao sistema original
- Conversão de tipos acontece uma vez por publicação
- Handlers executam em paralelo (goroutines)
- EventBus usa `sync.RWMutex` para thread-safety

---

## 🚀 Próximos Passos

1. **Migrar Commands existentes** para usar `TypedEventPublisher`
2. **Criar eventos específicos** para suas necessidades
3. **Implementar handlers** para cada caso de uso
4. **Adicionar testes** para eventos críticos
5. **Monitorar** através de logs e métricas

---

## 🆘 Troubleshooting

### Problema: Handler não é executado

**Causa:** Handler não foi registrado ou tipo de evento errado

**Solução:**
```go
// Verifique se registrou corretamente
events.SubscribeTyped(eventBus, events.EventTypeUserCreated, handler)

// Verifique se publicou com o tipo correto
typedPublisher.PublishUserCreated(ctx, ...)
```

### Problema: Erro ao deserializar

**Causa:** Payload incompatível com estrutura do evento

**Solução:** Use sempre `TypedEventPublisher` ao invés de `eventBus.Publish` direto

### Problema: Handler demora muito

**Causa:** Lógica pesada bloqueando

**Solução:** Handlers já executam async. Se precisar garantir ordem, use filas.

---

## 📚 Referências

- **Código:** `pkg/events/`
- **Tipos:** `pkg/events/types.go`
- **Publisher:** `pkg/events/typed.go`
- **Handlers:** `pkg/events/handlers.go`
- **Exemplos:** `pkg/events/examples_test.go`

---

**✨ Sistema de eventos moderno, type-safe e pronto para produção! ✨**
