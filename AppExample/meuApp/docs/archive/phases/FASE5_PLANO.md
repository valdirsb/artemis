# 🚀 FASE 5 - EVENT BUS COM GENERICS

> **Data de Início:** 18 de Outubro de 2025  
> **Status:** 🟡 Em Progresso  
> **Objetivo:** Modernizar sistema de eventos com Type-Safety e Generics

---

## 🎯 OBJETIVOS

Refatorar o sistema de eventos atual para usar **Generics do Go 1.18+**, trazendo:
- ✅ Type-safety em tempo de compilação
- ✅ Menos código boilerplate
- ✅ Melhor IntelliSense/autocomplete
- ✅ Eventos fortemente tipados
- ✅ Handlers específicos por tipo de evento

---

## 📋 PLANO DE IMPLEMENTAÇÃO

### 1️⃣ Sistema Base de Eventos (Alta Prioridade)
- [ ] Criar `pkg/events/event.go`
  - Estrutura genérica `Event[T]`
  - Interface `EventHandler[T]`
  - Métodos utilitários

- [ ] Criar `pkg/events/bus.go`
  - Implementação `EventBus` com generics
  - Registro de handlers tipados
  - Publicação async/sync
  - Suporte a múltiplos subscribers

### 2️⃣ Definição de Eventos por Módulo

**User Module:**
- [ ] `UserCreatedEvent` - Após criação de usuário
- [ ] `UserUpdatedEvent` - Após atualização
- [ ] `UserDeletedEvent` - Após exclusão

**Product Module:**
- [ ] `ProductCreatedEvent` - Após criação
- [ ] `ProductUpdatedEvent` - Após atualização
- [ ] `LowStockEvent` - Quando estoque < 10

**Order Module:**
- [ ] `OrderCreatedEvent` - Após criação de pedido
- [ ] `OrderStatusChangedEvent` - Mudança de status
- [ ] `OrderCancelledEvent` - Cancelamento

### 3️⃣ Event Handlers
- [ ] Email handler (envio de emails)
- [ ] Notification handler (notificações)
- [ ] Audit log handler (auditoria)
- [ ] Analytics handler (métricas)

### 4️⃣ Integração
- [ ] Atualizar Commands para publicar eventos tipados
- [ ] Atualizar bootstrap
- [ ] Remover sistema antigo de eventos

---

## 🏗️ ARQUITETURA PROPOSTA

### Antes (Sistema Atual):
```go
// Evento genérico sem tipo
type Event struct {
    Type      string
    Payload   interface{}  // ❌ Sem type-safety
    Timestamp time.Time
}

// Publicação manual
eventBus.Publish(ctx, Event{
    Type: "user.created",
    Payload: map[string]interface{}{  // ❌ Propenso a erros
        "user_id": userID,
        "email": email,
    },
})
```

### Depois (Com Generics):
```go
// Evento tipado
type UserCreatedEvent struct {
    UserID string
    Email  string
    Name   string
}

// Publicação type-safe
eventBus.Publish(ctx, Event[UserCreatedEvent]{
    Data: UserCreatedEvent{  // ✅ Type-safe!
        UserID: user.ID,
        Email:  user.Email,
        Name:   user.Name,
    },
    Timestamp: time.Now(),
})

// Handler tipado
type EmailHandler struct{}

func (h *EmailHandler) Handle(ctx context.Context, event Event[UserCreatedEvent]) error {
    // ✅ Acesso direto aos campos tipados
    return sendEmail(event.Data.Email, event.Data.Name)
}
```

---

## 📊 BENEFÍCIOS ESPERADOS

| Aspecto | Antes | Depois |
|---------|-------|--------|
| Type-Safety | ❌ Runtime errors | ✅ Compile-time |
| Autocomplete | ❌ Limitado | ✅ Completo |
| Refactoring | ❌ Difícil | ✅ Seguro |
| Performance | 🟡 Reflection | ✅ Direto |
| Testabilidade | 🟡 Média | ✅ Alta |
| Manutenção | 🟡 Média | ✅ Fácil |

---

## 🎓 CONCEITOS APLICADOS

- **Generics (Go 1.18+)** - Type parameters para reutilização
- **Observer Pattern** - Pub/Sub desacoplado
- **Event-Driven Architecture** - Comunicação assíncrona
- **Type Constraints** - Garantias em tempo de compilação
- **Async Processing** - Goroutines para handlers

---

## 📈 ESTIMATIVA

- **Complexidade:** Média-Alta
- **Tempo Estimado:** 2-3 horas
- **Impacto:** Alto (melhora significativa)
- **Breaking Changes:** Médio (migração controlada)

---

## ✅ CHECKLIST DE TAREFAS

- [ ] 1. Criar estruturas base com generics
- [ ] 2. Implementar EventBus genérico
- [ ] 3. Definir eventos do User Module
- [ ] 4. Definir eventos do Product Module
- [ ] 5. Definir eventos do Order Module
- [ ] 6. Criar event handlers de exemplo
- [ ] 7. Migrar publicação de eventos
- [ ] 8. Atualizar bootstrap
- [ ] 9. Testar e validar

---

**Vamos começar! 🚀**
