# ✅ Fase 6: Auto-registro - Order Module

> **Módulo mais complexo:** Implementa dependências cross-module

---

## 🎯 Order Module - Implementado com Sucesso

### Características Especiais

O Order Module é o mais complexo porque:
1. **Depende de UserRepository** - Para validar usuários
2. **Depende de ProductRepository** - Para validar produtos e atualizar estoque
3. **Demonstra injeção de dependências cross-module**

### Estrutura Implementada

```go
type OrderModule struct {
    db       *gorm.DB
    eventBus *events.EventBus
    logger   contracts.Logger
}

func (m *OrderModule) Register(registry *container.ModuleRegistry) error {
    // 1. Registrar próprio repository
    orderRepo := repository.NewMySQLOrderRepository(m.db)
    
    // 2. ⭐ Obter dependências de outros módulos
    userRepo := registry.GetRepository("user")
    productRepo := registry.GetRepository("product")
    
    // 3. Criar handlers com todas as dependências
    createOrderHandler := commands.NewCreateOrderHandler(
        orderRepo, userRepo, productRepo, eventBus, logger
    )
    
    // ... registrar o resto
}
```

### 🔗 Dependências Cross-Module

**Ordem de Registro Importante:**
```
1. User Module    → Registra UserRepository
2. Product Module → Registra ProductRepository  
3. Order Module   → Usa ambos os repositórios acima ✅
```

Se tentar registrar Order antes de User/Product, vai dar erro!

### Componentes Registrados

- ✅ **Repository:** `MySQLOrderRepository`
- ✅ **Command Handlers:** 
  - CreateOrder (com UserRepo + ProductRepo)
  - UpdateOrderStatus
  - CancelOrder
- ✅ **Query Handlers:**
  - GetOrder
  - GetOrdersByUser
- ✅ **Application Service:** `OrderApplicationService`
- ✅ **HTTP Handler** com rotas:
  - `POST /orders` - Criar pedido
  - `GET /orders/:id` - Buscar pedido
  - `PUT /orders/:id/status` - Atualizar status
  - `POST /orders/:id/cancel` - Cancelar pedido
  - `GET /orders/user/:user_id` - Buscar pedidos do usuário
- ✅ **gRPC Handler**

### 🎨 Pattern: Dependency Resolution

```go
// Obter repository de outro módulo
userRepoInterface, err := registry.GetRepository("user")
if err != nil {
    return fmt.Errorf("failed to get user repository: %w", err)
}

// Type assertion para interface correta
userRepo, ok := userRepoInterface.(userPorts.UserRepository)
if !ok {
    return fmt.Errorf("user repository has wrong type")
}

// Usar na criação do handler
createOrderHandler := commands.NewCreateOrderHandler(
    orderRepo,
    userRepo,      // ← Cross-module dependency
    productRepo,   // ← Cross-module dependency
    m.eventBus,
    m.logger,
)
```

### 💡 Benefícios da Abordagem

1. **Desacoplamento:** Módulos não se importam diretamente
2. **Flexibilidade:** Fácil trocar implementações
3. **Testabilidade:** Pode injetar mocks facilmente
4. **Ordem de Inicialização:** Registry garante que dependências existam
5. **Type Safety:** Type assertions garantem tipos corretos

### ⚠️ Tratamento de Erros

```go
// Verificar se dependência existe
if err != nil {
    return fmt.Errorf("failed to get user repository: %w", err)
}

// Verificar se tipo está correto
if !ok {
    return fmt.Errorf("user repository has wrong type")
}
```

Isso previne:
- Pânico em runtime
- Type mismatch
- Módulos registrados fora de ordem

---

## 📊 Status Geral - Fase 6

### Progresso: 87% (7/8 tarefas)

| Módulo | Status | Dependências | Complexidade |
|--------|--------|--------------|--------------|
| User | ✅ | Nenhuma | Baixa |
| Product | ✅ | Nenhuma | Baixa |
| Order | ✅ | User + Product | **Alta** |

**Próximo:** Refatorar Bootstrap! 🚀

---

## 🎯 Comparação: Antes vs Depois

### Antes (Bootstrap Manual)

```go
// ❌ Código repetitivo e propenso a erros
func InitializeOrder() {
    // Obter UserRepository de algum lugar...
    userRepo := GetUserRepositorySomehow()
    
    // Obter ProductRepository de outro lugar...
    productRepo := GetProductRepositorySomehow()
    
    // Criar todos os handlers manualmente...
    createHandler := commands.NewCreateOrderHandler(...)
    updateHandler := commands.NewUpdateOrderHandler(...)
    // ... 20+ linhas
    
    // Criar application service
    orderService := services.NewOrderApplicationService(...)
    
    // Criar HTTP handler
    httpHandler := http.NewOrderHandler(...)
    
    // Registrar rotas manualmente
    router.POST("/orders", httpHandler.CreateOrder)
    // ... mais rotas
}
```

### Depois (Auto-registro)

```go
// ✅ Simples, limpo e escalável
orderModule := modules.NewOrderModule(db, eventBus, logger)
if err := orderModule.Register(registry); err != nil {
    log.Fatal(err)
}

// Pronto! Tudo registrado automaticamente ✨
```

**Redução:** ~95% menos código! 🎉

---

## ✅ Checklist de Validação

- [x] Order Module compila sem erros
- [x] Dependências cross-module resolvidas
- [x] Type assertions implementadas
- [x] Error handling robusto
- [x] HTTP routes configuradas
- [x] gRPC service registrado
- [x] Repository registrado
- [x] Application service registrado
- [ ] Bootstrap refatorado (próximo!)
- [ ] Testes de integração

---

**Próximo passo:** Refatorar `bootstrap.go` para usar o Registry! 🚀
