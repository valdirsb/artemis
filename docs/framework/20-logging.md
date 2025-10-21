# 📝 Sistema de Logging

<div align="center">

**Guia Completo de Logging Estruturado no Artemis Framework**

[Conceitos](#-conceitos) • [Configuração](#-configuração) • [Uso Prático](#-uso-prático) • [Melhores Práticas](#-melhores-práticas)

</div>

---

## 📚 Índice

1. [Visão Geral](#-visão-geral)
2. [Conceitos Fundamentais](#-conceitos-fundamentais)
3. [Estrutura do Logger](#-estrutura-do-logger)
4. [Níveis de Log](#-níveis-de-log)
5. [Logging Estruturado](#-logging-estruturado)
6. [Context Logging](#-context-logging)
7. [Uso em Diferentes Camadas](#-uso-em-diferentes-camadas)
8. [Integração com Ferramentas](#-integração-com-ferramentas)
9. [Performance e Boas Práticas](#-performance-e-boas-práticas)
10. [Troubleshooting](#-troubleshooting)

---

## 🎯 Visão Geral

O sistema de logging do Artemis Framework é baseado em **logging estruturado**, oferecendo:

✅ **Níveis de Log** - Debug, Info, Warn, Error, Fatal  
✅ **Campos Estruturados** - Logs com contexto rico  
✅ **Type-Safe** - Interface definida por contratos  
✅ **Performance** - Zero allocations em hot paths  
✅ **Rastreabilidade** - Correlation IDs e contexto  
✅ **Integração** - Compatível com ELK, Grafana, Datadog  

### Por Que Logging Estruturado?

**❌ Logging Tradicional (String-based):**
```go
log.Printf("User %s created order %s with total $%.2f", userID, orderID, total)
```

**Problemas:**
- Difícil de fazer queries
- Não há padronização
- Parse complexo para ferramentas
- Sem type-safety

**✅ Logging Estruturado (Field-based):**
```go
logger.Info("Order created",
    Field{Key: "user_id", Value: userID},
    Field{Key: "order_id", Value: orderID},
    Field{Key: "total", Value: total},
)
```

**Benefícios:**
- Queries fáceis: `user_id:123`
- Formato padronizado (JSON)
- Ferramentas podem indexar
- Type-safe com autocomplete

---

## 🔍 Conceitos Fundamentais

### 1. Interface Logger

O framework define uma interface padrão em `pkg/contracts/infrastructure.go`:

```go
package contracts

type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    With(fields ...Field) Logger
}

type Field struct {
    Key   string
    Value interface{}
}
```

**Responsabilidades:**
- Define contrato de logging
- Permite múltiplas implementações (stdout, file, cloud)
- Garante consistência em todo o framework

### 2. Structured Fields

Campos estruturados carregam metadados contextuais:

```go
type Field struct {
    Key   string        // Nome do campo (ex: "user_id")
    Value interface{}   // Valor (qualquer tipo)
}
```

**Exemplos:**
```go
Field{Key: "user_id", Value: "usr_123"}
Field{Key: "duration_ms", Value: 45}
Field{Key: "success", Value: true}
Field{Key: "error", Value: err}
```

### 3. Log Levels

Hierarquia de níveis (do mais detalhado ao mais crítico):

```
DEBUG → INFO → WARN → ERROR → FATAL
```

**Quando usar cada nível:**
- **DEBUG**: Informações de desenvolvimento/debug
- **INFO**: Eventos normais do sistema
- **WARN**: Situações anormais mas recuperáveis
- **ERROR**: Erros que afetam operações
- **FATAL**: Erros irrecuperáveis (app termina)

---

## 🏗️ Estrutura do Logger

### Implementação Básica

**`pkg/adapters/logger/structured_logger.go`**

```go
package logger

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
    
    "meuApp/pkg/contracts"
)

type StructuredLogger struct {
    level      LogLevel
    output     *os.File
    contextFields []contracts.Field
}

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

type LogEntry struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
}

func NewStructuredLogger(level LogLevel) contracts.Logger {
    return &StructuredLogger{
        level:         level,
        output:        os.Stdout,
        contextFields: []contracts.Field{},
    }
}

func (l *StructuredLogger) Debug(msg string, fields ...contracts.Field) {
    if l.level <= DEBUG {
        l.log("DEBUG", msg, fields...)
    }
}

func (l *StructuredLogger) Info(msg string, fields ...contracts.Field) {
    if l.level <= INFO {
        l.log("INFO", msg, fields...)
    }
}

func (l *StructuredLogger) Warn(msg string, fields ...contracts.Field) {
    if l.level <= WARN {
        l.log("WARN", msg, fields...)
    }
}

func (l *StructuredLogger) Error(msg string, fields ...contracts.Field) {
    if l.level <= ERROR {
        l.log("ERROR", msg, fields...)
    }
}

func (l *StructuredLogger) Fatal(msg string, fields ...contracts.Field) {
    l.log("FATAL", msg, fields...)
    os.Exit(1)
}

func (l *StructuredLogger) With(fields ...contracts.Field) contracts.Logger {
    // Cria novo logger com campos contextuais
    return &StructuredLogger{
        level:         l.level,
        output:        l.output,
        contextFields: append(l.contextFields, fields...),
    }
}

func (l *StructuredLogger) log(level, msg string, fields ...contracts.Field) {
    // Combina campos contextuais + campos do log
    allFields := append(l.contextFields, fields...)
    
    // Converte para map
    fieldMap := make(map[string]interface{})
    for _, f := range allFields {
        fieldMap[f.Key] = f.Value
    }
    
    // Cria entry
    entry := LogEntry{
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Level:     level,
        Message:   msg,
        Fields:    fieldMap,
    }
    
    // Serializa para JSON
    data, err := json.Marshal(entry)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to marshal log: %v\n", err)
        return
    }
    
    // Escreve no output
    fmt.Fprintln(l.output, string(data))
}
```

**Saída JSON:**
```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "level": "INFO",
  "message": "User created successfully",
  "fields": {
    "user_id": "usr_123",
    "email": "user@example.com",
    "duration_ms": 45
  }
}
```

---

## 📊 Níveis de Log

### 1. DEBUG

**Quando usar:**
- Desenvolvimento e troubleshooting
- Valores de variáveis durante execução
- Fluxo detalhado de código

**Exemplo:**
```go
func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) error {
    h.logger.Debug("Starting user creation",
        contracts.Field{Key: "name", Value: cmd.Name},
        contracts.Field{Key: "email", Value: cmd.Email},
    )
    
    // Validação
    if err := h.validate(cmd); err != nil {
        h.logger.Debug("Validation failed",
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err
    }
    
    // Hash password
    h.logger.Debug("Hashing password")
    hashedPwd, _ := h.hasher.Hash(cmd.Password)
    
    // Criar usuário
    user := &domain.User{
        ID:       generateID(),
        Name:     cmd.Name,
        Email:    cmd.Email,
        Password: hashedPwd,
    }
    
    h.logger.Debug("Saving user to database",
        contracts.Field{Key: "user_id", Value: user.ID},
    )
    
    if err := h.repo.Save(ctx, user); err != nil {
        return err
    }
    
    h.logger.Debug("User creation completed",
        contracts.Field{Key: "user_id", Value: user.ID},
    )
    
    return nil
}
```

### 2. INFO

**Quando usar:**
- Eventos normais do sistema
- Requisições HTTP
- Operações bem-sucedidas
- Métricas de negócio

**Exemplo:**
```go
func (h *CreateOrderHandler) Handle(ctx context.Context, cmd *CreateOrderCommand) error {
    startTime := time.Now()
    
    // Criar pedido
    order, err := h.buildOrder(cmd)
    if err != nil {
        return err
    }
    
    // Salvar
    if err := h.repo.Save(ctx, order); err != nil {
        return err
    }
    
    // Log de sucesso
    h.logger.Info("Order created successfully",
        contracts.Field{Key: "order_id", Value: order.ID},
        contracts.Field{Key: "user_id", Value: cmd.UserID},
        contracts.Field{Key: "total", Value: order.Total},
        contracts.Field{Key: "items_count", Value: len(order.Items)},
        contracts.Field{Key: "duration_ms", Value: time.Since(startTime).Milliseconds()},
    )
    
    return nil
}
```

### 3. WARN

**Quando usar:**
- Situações anormais mas recuperáveis
- Limites se aproximando (quota, rate limit)
- Dados inconsistentes mas não críticos
- Fallback para valores padrão

**Exemplo:**
```go
func (s *EmailService) SendWelcomeEmail(userID, email string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Tentar enviar email
    if err := s.client.Send(ctx, email, welcomeTemplate); err != nil {
        // Email falhou mas não é crítico
        s.logger.Warn("Failed to send welcome email",
            contracts.Field{Key: "user_id", Value: userID},
            contracts.Field{Key: "email", Value: email},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err // Retorna erro mas não bloqueia criação do usuário
    }
    
    return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
    val, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        // Cache miss - esperado, mas pode indicar problema se muito frequente
        r.logger.Warn("Cache miss",
            contracts.Field{Key: "key", Value: key},
        )
        return nil, ErrCacheMiss
    }
    
    if err != nil {
        // Erro de conexão - fallback para DB
        r.logger.Warn("Redis connection error, falling back to database",
            contracts.Field{Key: "key", Value: key},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return nil, err
    }
    
    return val, nil
}
```

### 4. ERROR

**Quando usar:**
- Erros que impedem operações
- Falhas de integração
- Violations de regras de negócio
- Exceptions não esperadas

**Exemplo:**
```go
func (h *ProcessPaymentHandler) Handle(ctx context.Context, cmd *ProcessPaymentCommand) error {
    // Buscar pedido
    order, err := h.orderRepo.FindByID(ctx, cmd.OrderID)
    if err != nil {
        h.logger.Error("Failed to find order",
            contracts.Field{Key: "order_id", Value: cmd.OrderID},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err
    }
    
    // Processar pagamento
    result, err := h.paymentGateway.ProcessPayment(ctx, &PaymentRequest{
        OrderID: order.ID,
        Amount:  order.Total,
        Method:  cmd.PaymentMethod,
    })
    
    if err != nil {
        h.logger.Error("Payment processing failed",
            contracts.Field{Key: "order_id", Value: order.ID},
            contracts.Field{Key: "amount", Value: order.Total},
            contracts.Field{Key: "method", Value: cmd.PaymentMethod},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err
    }
    
    h.logger.Info("Payment processed successfully",
        contracts.Field{Key: "order_id", Value: order.ID},
        contracts.Field{Key: "transaction_id", Value: result.TransactionID},
        contracts.Field{Key: "amount", Value: order.Total},
    )
    
    return nil
}
```

### 5. FATAL

**Quando usar:**
- Erros de inicialização
- Configurações críticas inválidas
- Dependências essenciais indisponíveis
- **CUIDADO**: App termina após log!

**Exemplo:**
```go
func main() {
    // Carregar configuração
    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Fatal("Failed to load configuration",
            contracts.Field{Key: "error", Value: err.Error()},
        )
        // App termina aqui
    }
    
    // Conectar ao banco
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        logger.Fatal("Failed to connect to database",
            contracts.Field{Key: "url", Value: cfg.DatabaseURL},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        // App termina aqui
    }
    
    logger.Info("Application started successfully")
}
```

---

## 🔗 Logging Estruturado

### Campos Comuns

**Padronize campos para facilitar queries:**

```go
// ✅ BOM: Nomes consistentes
logger.Info("User logged in",
    contracts.Field{Key: "user_id", Value: userID},
    contracts.Field{Key: "ip_address", Value: ip},
    contracts.Field{Key: "user_agent", Value: ua},
)

logger.Info("Order created",
    contracts.Field{Key: "user_id", Value: userID},  // Mesmo nome!
    contracts.Field{Key: "order_id", Value: orderID},
)

// ❌ RUIM: Nomes inconsistentes
logger.Info("User logged in",
    contracts.Field{Key: "userId", Value: userID},  // camelCase
)

logger.Info("Order created",
    contracts.Field{Key: "user", Value: userID},  // Nome diferente
)
```

### Campos Recomendados

**Identifiers:**
- `user_id` - ID do usuário
- `order_id` - ID do pedido
- `request_id` - ID da requisição
- `correlation_id` - ID de correlação (trace)
- `session_id` - ID da sessão

**Metrics:**
- `duration_ms` - Duração em milissegundos
- `status_code` - Código HTTP
- `count` - Quantidade de items
- `size_bytes` - Tamanho em bytes

**Context:**
- `method` - Método HTTP (GET, POST)
- `path` - Path da requisição
- `ip_address` - IP do cliente
- `user_agent` - User agent
- `environment` - Ambiente (dev, prod)

**Business:**
- `amount` - Valor monetário
- `currency` - Moeda (BRL, USD)
- `quantity` - Quantidade
- `category` - Categoria

### Exemplo Completo

```go
func (h *HTTPHandler) CreateOrder(c *gin.Context) {
    requestID := c.GetHeader("X-Request-ID")
    startTime := time.Now()
    
    var req dto.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("Invalid request body",
            contracts.Field{Key: "request_id", Value: requestID},
            contracts.Field{Key: "method", Value: c.Request.Method},
            contracts.Field{Key: "path", Value: c.Request.URL.Path},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    cmd := &CreateOrderCommand{
        UserID: req.UserID,
        Items:  req.Items,
    }
    
    orderID, err := h.commandBus.Execute(c.Request.Context(), cmd)
    if err != nil {
        h.logger.Error("Failed to create order",
            contracts.Field{Key: "request_id", Value: requestID},
            contracts.Field{Key: "user_id", Value: req.UserID},
            contracts.Field{Key: "items_count", Value: len(req.Items)},
            contracts.Field{Key: "error", Value: err.Error()},
            contracts.Field{Key: "duration_ms", Value: time.Since(startTime).Milliseconds()},
        )
        c.JSON(500, gin.H{"error": "failed to create order"})
        return
    }
    
    h.logger.Info("Order created successfully",
        contracts.Field{Key: "request_id", Value: requestID},
        contracts.Field{Key: "user_id", Value: req.UserID},
        contracts.Field{Key: "order_id", Value: orderID},
        contracts.Field{Key: "items_count", Value: len(req.Items)},
        contracts.Field{Key: "duration_ms", Value: time.Since(startTime).Milliseconds()},
        contracts.Field{Key: "status_code", Value: 201},
    )
    
    c.JSON(201, gin.H{"order_id": orderID})
}
```

**Saída JSON:**
```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "level": "INFO",
  "message": "Order created successfully",
  "fields": {
    "request_id": "req_abc123",
    "user_id": "usr_456",
    "order_id": "ord_789",
    "items_count": 3,
    "duration_ms": 125,
    "status_code": 201
  }
}
```

**Query no ELK:**
```
user_id:usr_456 AND level:ERROR
duration_ms:>1000
status_code:500
```

---

## 🎯 Context Logging

### Logger com Contexto

Use `With()` para criar loggers com campos contextuais:

```go
func (h *UserModule) RegisterHTTPRoutes(router *gin.Engine) {
    // Logger com contexto do módulo
    moduleLogger := h.logger.With(
        contracts.Field{Key: "module", Value: "user"},
        contracts.Field{Key: "layer", Value: "http"},
    )
    
    handler := NewHTTPHandler(h.commandBus, h.queryBus, moduleLogger)
    
    router.POST("/users", handler.CreateUser)
    router.GET("/users/:id", handler.GetUser)
}

func (h *HTTPHandler) CreateUser(c *gin.Context) {
    // Logger já tem campos contextuais (module, layer)
    h.logger.Info("Creating user")
    // Output: {"module": "user", "layer": "http", "message": "Creating user"}
    
    // Adicionar mais contexto
    requestLogger := h.logger.With(
        contracts.Field{Key: "request_id", Value: c.GetHeader("X-Request-ID")},
        contracts.Field{Key: "ip_address", Value: c.ClientIP()},
    )
    
    requestLogger.Info("Processing request")
    // Output: {"module": "user", "layer": "http", "request_id": "...", "ip_address": "...", ...}
}
```

### Correlation ID (Distributed Tracing)

**Middleware para injetar Correlation ID:**

```go
func CorrelationIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Pegar correlation ID do header ou gerar novo
        correlationID := c.GetHeader("X-Correlation-ID")
        if correlationID == "" {
            correlationID = uuid.New().String()
        }
        
        // Adicionar ao context
        ctx := context.WithValue(c.Request.Context(), "correlation_id", correlationID)
        c.Request = c.Request.WithContext(ctx)
        
        // Adicionar ao response header
        c.Header("X-Correlation-ID", correlationID)
        
        c.Next()
    }
}
```

**Usar Correlation ID em logs:**

```go
func (h *Handler) CreateUser(c *gin.Context) {
    correlationID := c.Request.Context().Value("correlation_id").(string)
    
    logger := h.logger.With(
        contracts.Field{Key: "correlation_id", Value: correlationID},
    )
    
    logger.Info("Creating user")
    
    // Passar logger para camadas inferiores
    cmd := &CreateUserCommand{...}
    err := h.commandHandler.Handle(c.Request.Context(), cmd, logger)
}
```

**Benefícios:**
- Rastrear requisição através de múltiplos serviços
- Correlacionar logs de diferentes componentes
- Debugging de sistemas distribuídos

---

## 🏢 Uso em Diferentes Camadas

### 1. HTTP Handlers (Adapters)

```go
func (h *HTTPHandler) CreateProduct(c *gin.Context) {
    requestID := c.GetHeader("X-Request-ID")
    startTime := time.Now()
    
    var req dto.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Warn("Invalid request payload",
            contracts.Field{Key: "request_id", Value: requestID},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        c.JSON(400, gin.H{"error": "invalid payload"})
        return
    }
    
    h.logger.Info("Processing product creation request",
        contracts.Field{Key: "request_id", Value: requestID},
        contracts.Field{Key: "sku", Value: req.SKU},
    )
    
    cmd := &CreateProductCommand{
        SKU:   req.SKU,
        Name:  req.Name,
        Price: req.Price,
    }
    
    productID, err := h.commandBus.Execute(c.Request.Context(), cmd)
    if err != nil {
        h.logger.Error("Failed to create product",
            contracts.Field{Key: "request_id", Value: requestID},
            contracts.Field{Key: "sku", Value: req.SKU},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        c.JSON(500, gin.H{"error": "internal error"})
        return
    }
    
    h.logger.Info("Product created successfully",
        contracts.Field{Key: "request_id", Value: requestID},
        contracts.Field{Key: "product_id", Value: productID},
        contracts.Field{Key: "duration_ms", Value: time.Since(startTime).Milliseconds()},
    )
    
    c.JSON(201, gin.H{"product_id": productID})
}
```

### 2. Command Handlers (Application)

```go
func (h *CreateUserHandler) Handle(ctx context.Context, cmd *CreateUserCommand) error {
    h.logger.Info("Handling CreateUserCommand",
        contracts.Field{Key: "name", Value: cmd.Name},
        contracts.Field{Key: "email", Value: cmd.Email},
    )
    
    // Validar
    if err := h.validate(cmd); err != nil {
        h.logger.Warn("Validation failed",
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err
    }
    
    // Verificar duplicação
    exists, err := h.repo.ExistsByEmail(ctx, cmd.Email)
    if err != nil {
        h.logger.Error("Failed to check email uniqueness",
            contracts.Field{Key: "email", Value: cmd.Email},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err
    }
    
    if exists {
        h.logger.Warn("Email already in use",
            contracts.Field{Key: "email", Value: cmd.Email},
        )
        return ErrEmailAlreadyInUse
    }
    
    // Criar usuário
    user := &domain.User{
        ID:    generateID(),
        Name:  cmd.Name,
        Email: cmd.Email,
    }
    
    if err := h.repo.Save(ctx, user); err != nil {
        h.logger.Error("Failed to save user",
            contracts.Field{Key: "user_id", Value: user.ID},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err
    }
    
    h.logger.Info("User created successfully",
        contracts.Field{Key: "user_id", Value: user.ID},
        contracts.Field{Key: "email", Value: user.Email},
    )
    
    return nil
}
```

### 3. Repositories (Infrastructure)

```go
func (r *MySQLUserRepository) Save(ctx context.Context, user *domain.User) error {
    startTime := time.Now()
    
    query := `
        INSERT INTO users (id, name, email, password, created_at)
        VALUES (?, ?, ?, ?, ?)
    `
    
    result, err := r.db.ExecContext(ctx, query,
        user.ID, user.Name, user.Email, user.Password, user.CreatedAt,
    )
    
    if err != nil {
        r.logger.Error("Failed to insert user",
            contracts.Field{Key: "user_id", Value: user.ID},
            contracts.Field{Key: "error", Value: err.Error()},
            contracts.Field{Key: "query_time_ms", Value: time.Since(startTime).Milliseconds()},
        )
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    
    r.logger.Debug("User saved to database",
        contracts.Field{Key: "user_id", Value: user.ID},
        contracts.Field{Key: "rows_affected", Value: rowsAffected},
        contracts.Field{Key: "query_time_ms", Value: time.Since(startTime).Milliseconds()},
    )
    
    return nil
}
```

### 4. Event Subscribers

```go
func (s *EmailSubscriber) Handle(ctx context.Context, event contracts.Event) error {
    userCreated := event.(*events.UserCreatedEvent)
    
    s.logger.Info("Handling UserCreatedEvent",
        contracts.Field{Key: "user_id", Value: userCreated.UserID},
        contracts.Field{Key: "email", Value: userCreated.Email},
    )
    
    if err := s.emailService.SendWelcomeEmail(userCreated.Email); err != nil {
        s.logger.Error("Failed to send welcome email",
            contracts.Field{Key: "user_id", Value: userCreated.UserID},
            contracts.Field{Key: "email", Value: userCreated.Email},
            contracts.Field{Key: "error", Value: err.Error()},
        )
        return err
    }
    
    s.logger.Info("Welcome email sent",
        contracts.Field{Key: "user_id", Value: userCreated.UserID},
    )
    
    return nil
}
```

---

## 🔌 Integração com Ferramentas

### 1. ELK Stack (Elasticsearch, Logstash, Kibana)

**Filebeat config (`filebeat.yml`):**

```yaml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/artemis/*.log
    json.keys_under_root: true
    json.add_error_key: true

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "artemis-logs-%{+yyyy.MM.dd}"

setup.kibana:
  host: "localhost:5601"
```

**Queries no Kibana:**
```
level:ERROR
user_id:usr_123 AND level:ERROR
duration_ms:>1000
module:order AND level:INFO
```

### 2. Grafana Loki

**Promtail config (`promtail.yml`):**

```yaml
server:
  http_listen_port: 9080

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://localhost:3100/loki/api/v1/push

scrape_configs:
  - job_name: artemis
    static_configs:
      - targets:
          - localhost
        labels:
          job: artemis
          __path__: /var/log/artemis/*.log
    pipeline_stages:
      - json:
          expressions:
            level: level
            message: message
            user_id: fields.user_id
      - labels:
          level:
          user_id:
```

**Queries no Grafana:**
```
{job="artemis"} |= "ERROR"
{job="artemis", user_id="usr_123"}
{job="artemis"} | json | duration_ms > 1000
```

### 3. Datadog

**Implementação custom:**

```go
package logger

import (
    "github.com/DataDog/datadog-go/statsd"
    "meuApp/pkg/contracts"
)

type DatadogLogger struct {
    client *statsd.Client
    // ...
}

func (l *DatadogLogger) Error(msg string, fields ...contracts.Field) {
    // Log para stdout (coletado por agent)
    l.baseLogger.Error(msg, fields...)
    
    // Enviar métrica
    l.client.Incr("errors.count", []string{
        "level:error",
        "module:" + l.getField(fields, "module"),
    }, 1)
}
```

### 4. CloudWatch (AWS)

```go
package logger

import (
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type CloudWatchLogger struct {
    client    *cloudwatchlogs.CloudWatchLogs
    groupName string
    streamName string
}

func (l *CloudWatchLogger) Info(msg string, fields ...contracts.Field) {
    logEvent := &cloudwatchlogs.InputLogEvent{
        Message:   aws.String(l.formatMessage(msg, fields)),
        Timestamp: aws.Int64(time.Now().UnixNano() / 1000000),
    }
    
    _, err := l.client.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
        LogGroupName:  aws.String(l.groupName),
        LogStreamName: aws.String(l.streamName),
        LogEvents:     []*cloudwatchlogs.InputLogEvent{logEvent},
    })
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to send to CloudWatch: %v\n", err)
    }
}
```

---

## ⚡ Performance e Boas Práticas

### 1. Lazy Evaluation

**❌ Evite computação desnecessária:**

```go
// ❌ RUIM: JSON serializado mesmo se DEBUG desabilitado
logger.Debug("User data",
    contracts.Field{Key: "user", Value: json.Marshal(user)}, // Sempre executa!
)
```

**✅ Use closures para lazy eval:**

```go
// ✅ BOM: Só serializa se DEBUG estiver ativo
if logger.IsDebugEnabled() {
    userData, _ := json.Marshal(user)
    logger.Debug("User data",
        contracts.Field{Key: "user", Value: string(userData)},
    )
}
```

### 2. Evite Alocações Excessivas

```go
// ❌ RUIM: Múltiplas alocações
logger.Info("Order created",
    contracts.Field{Key: "order_id", Value: order.ID},
    contracts.Field{Key: "user_id", Value: order.UserID},
    contracts.Field{Key: "total", Value: fmt.Sprintf("%.2f", order.Total)}, // Aloca string
)

// ✅ BOM: Menos alocações
logger.Info("Order created",
    contracts.Field{Key: "order_id", Value: order.ID},
    contracts.Field{Key: "user_id", Value: order.UserID},
    contracts.Field{Key: "total", Value: order.Total}, // Número direto
)
```

### 3. Bufferize Logs (Batch)

```go
type BufferedLogger struct {
    buffer chan LogEntry
    flush  time.Duration
}

func (l *BufferedLogger) Start() {
    ticker := time.NewTicker(l.flush)
    go func() {
        batch := make([]LogEntry, 0, 100)
        
        for {
            select {
            case entry := <-l.buffer:
                batch = append(batch, entry)
                
                if len(batch) >= 100 {
                    l.flushBatch(batch)
                    batch = batch[:0]
                }
                
            case <-ticker.C:
                if len(batch) > 0 {
                    l.flushBatch(batch)
                    batch = batch[:0]
                }
            }
        }
    }()
}
```

### 4. Sampling (Rate Limiting)

Para logs muito frequentes:

```go
type SampledLogger struct {
    baseLogger contracts.Logger
    sampleRate float64  // 0.1 = 10% dos logs
}

func (l *SampledLogger) Debug(msg string, fields ...contracts.Field) {
    if rand.Float64() < l.sampleRate {
        l.baseLogger.Debug(msg, fields...)
    }
}
```

### 5. Não Logue Informações Sensíveis

```go
// ❌ RUIM: Senha em plain text
logger.Info("User created",
    contracts.Field{Key: "email", Value: user.Email},
    contracts.Field{Key: "password", Value: user.Password}, // 🚨 NUNCA!
)

// ✅ BOM: Sem dados sensíveis
logger.Info("User created",
    contracts.Field{Key: "user_id", Value: user.ID},
    contracts.Field{Key: "email", Value: user.Email},
    // Senha omitida
)

// ✅ BOM: Hash se necessário
logger.Debug("Password validation",
    contracts.Field{Key: "password_hash", Value: fmt.Sprintf("%.8s...", user.Password)}, // Primeiros 8 chars
)
```

### 6. Use Níveis Apropriados

```go
// ❌ RUIM: Tudo em INFO
logger.Info("Starting request processing")
logger.Info("Validating input")
logger.Info("Querying database")
logger.Info("Processing results")
logger.Info("Request completed")

// ✅ BOM: Níveis apropriados
logger.Info("Request received")         // Evento importante
logger.Debug("Validating input")        // Detalhe de implementação
logger.Debug("Querying database")       // Detalhe de implementação
logger.Debug("Processing results")      // Detalhe de implementação
logger.Info("Request completed")        // Evento importante
```

### 7. Contextualize Erros

```go
// ❌ RUIM: Erro sem contexto
logger.Error("Database error", 
    contracts.Field{Key: "error", Value: err.Error()},
)

// ✅ BOM: Erro com contexto rico
logger.Error("Failed to save order",
    contracts.Field{Key: "operation", Value: "save_order"},
    contracts.Field{Key: "order_id", Value: order.ID},
    contracts.Field{Key: "user_id", Value: order.UserID},
    contracts.Field{Key: "table", Value: "orders"},
    contracts.Field{Key: "error", Value: err.Error()},
    contracts.Field{Key: "error_type", Value: fmt.Sprintf("%T", err)},
)
```

---

## 🔍 Troubleshooting

### Problema 1: Logs não aparecem

**Causa:** Nível de log muito alto

**Solução:**
```go
// Verifique o nível configurado
logger := logger.NewStructuredLogger(logger.DEBUG) // DEBUG mostra tudo

// Ou via variável de ambiente
logLevel := os.Getenv("LOG_LEVEL") // "DEBUG", "INFO", "WARN", "ERROR"
```

### Problema 2: Logs desordenados

**Causa:** Concorrência sem sincronização

**Solução:**
```go
type SyncedLogger struct {
    mu     sync.Mutex
    logger contracts.Logger
}

func (l *SyncedLogger) Info(msg string, fields ...contracts.Field) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logger.Info(msg, fields...)
}
```

### Problema 3: Performance degradada

**Causa:** Logs síncronos bloqueando goroutines

**Solução:**
```go
type AsyncLogger struct {
    buffer chan LogEntry
}

func (l *AsyncLogger) Info(msg string, fields ...contracts.Field) {
    // Não bloqueia - envia para channel
    select {
    case l.buffer <- LogEntry{Level: "INFO", Message: msg, Fields: fields}:
    default:
        // Buffer cheio - descarta ou bloqueia
    }
}
```

### Problema 4: Logs muito grandes

**Causa:** Objetos completos sendo logados

**Solução:**
```go
// ❌ Evite:
logger.Info("User", contracts.Field{Key: "user", Value: user}) // Serializa tudo

// ✅ Prefira:
logger.Info("User",
    contracts.Field{Key: "user_id", Value: user.ID},
    contracts.Field{Key: "email", Value: user.Email},
    // Só campos relevantes
)
```

---

## 📋 Checklist de Logging

### ✅ Durante Desenvolvimento

- [ ] Usar níveis apropriados (DEBUG, INFO, WARN, ERROR)
- [ ] Incluir campos contextuais (IDs, timestamps)
- [ ] Adicionar correlation IDs para rastreamento
- [ ] Não logar dados sensíveis (senhas, tokens)
- [ ] Incluir duração de operações importantes
- [ ] Logar início e fim de operações críticas

### ✅ Em Produção

- [ ] Configurar nível INFO ou superior
- [ ] Rotacionar logs para evitar disco cheio
- [ ] Enviar logs para sistema centralizado (ELK, Loki)
- [ ] Configurar alertas para erros críticos
- [ ] Monitorar volume de logs (pode indicar problemas)
- [ ] Implementar sampling se volume for alto

### ✅ Para Debugging

- [ ] Usar DEBUG temporariamente
- [ ] Adicionar timestamps precisos
- [ ] Incluir stack traces para panics
- [ ] Logar payloads de requisições (sem dados sensíveis)
- [ ] Adicionar breadcrumbs em fluxos complexos

---

## 🎓 Resumo

### O Que Aprendemos

✅ **Logging estruturado** é superior a string-based  
✅ **Níveis de log** comunicam severidade  
✅ **Campos estruturados** facilitam queries  
✅ **Context logging** mantém contexto através de camadas  
✅ **Correlation IDs** rastreiam requisições distribuídas  
✅ **Performance** requer atenção (lazy eval, buffers)  
✅ **Integração** com ferramentas é essencial para produção  

### Padrões Essenciais

```go
// 1. Logger com contexto
moduleLogger := logger.With(
    contracts.Field{Key: "module", Value: "user"},
)

// 2. Campos padronizados
logger.Info("Operation completed",
    contracts.Field{Key: "user_id", Value: userID},
    contracts.Field{Key: "duration_ms", Value: duration},
)

// 3. Correlation ID
logger := baseLogger.With(
    contracts.Field{Key: "correlation_id", Value: correlationID},
)

// 4. Error context
logger.Error("Operation failed",
    contracts.Field{Key: "operation", Value: "save_order"},
    contracts.Field{Key: "order_id", Value: orderID},
    contracts.Field{Key: "error", Value: err.Error()},
)
```

---

## 📚 Próximos Passos

- **[22. Mocks e Fixtures →](22-mocks-fixtures.md)** - Criando mocks para testes
- **[21. Estratégia de Testes ←](21-testing-strategy.md)** - Testes com logging
- **[19. Error Handling ←](19-error-handling.md)** - Integração logging + errors

---

<div align="center">

**[⬆️ Voltar ao Topo](#-sistema-de-logging)**

</div>
