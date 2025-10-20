# 🔌 Adapters (Adaptadores)

## 📋 Índice
- [O que são Adapters?](#o-que-são-adapters)
- [Tipos de Adapters](#tipos-de-adapters)
- [HTTP Adapter](#http-adapter)
- [gRPC Adapter](#grpc-adapter)
- [Database Adapter](#database-adapter)
- [Criando Novos Adapters](#criando-novos-adapters)
- [Boas Práticas](#boas-práticas)

---

## 🎯 O que são Adapters?

**Adapters** (Adaptadores) são componentes que **conectam o núcleo da aplicação com o mundo externo**, implementando o padrão **Hexagonal Architecture** (Ports & Adapters).

### Conceito Visual

```
┌─────────────────────────────────────────────────────┐
│                  MUNDO EXTERNO                      │
│                                                     │
│  HTTP    gRPC    Database    Queue    Email        │
│  REST    Proto   MySQL       RabbitMQ SendGrid     │
└────┬───────┬────────┬──────────┬─────────┬─────────┘
     │       │        │          │         │
     ▼       ▼        ▼          ▼         ▼
┌────────────────────────────────────────────────────┐
│              ADAPTERS LAYER                        │
│                                                    │
│  HTTPAdapter  gRPCAdapter  MySQLAdapter           │
│  QueueAdapter EmailAdapter                        │
│                                                    │
│  ↕ Traduzem entre protocolos externos e domínio   │
└────────────────────────────────────────────────────┘
     ▲       ▲        ▲
     │       │        │
┌────┴───────┴────────┴────────────────────────────┐
│         APPLICATION LAYER (CORE)                 │
│                                                  │
│  Commands, Queries, Services                    │
│  ↕ Lógica de negócio pura                       │
└──────────────────────────────────────────────────┘
     ▲
     │
┌────┴──────────────────────────────────────────────┐
│              DOMAIN LAYER                         │
│                                                   │
│  Entities, Value Objects, Domain Events          │
└───────────────────────────────────────────────────┘
```

### Por que usar Adapters?

**Sem Adapters (❌ Ruim):**
```go
// Lógica de negócio misturada com HTTP
func CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // SQL direto no handler
    _, err := db.Exec("INSERT INTO users...")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "success"})
}
```

**Com Adapters (✅ Bom):**
```go
// Adapter traduz HTTP → Application
type HTTPUserAdapter struct {
    commandBus *CommandBus
}

func (a *HTTPUserAdapter) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse(err))
        return
    }
    
    // Usa application layer (lógica pura)
    cmd := &CreateUserCommand{
        Name:  req.Name,
        Email: req.Email,
    }
    
    result, err := a.commandBus.Execute(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(500, ErrorResponse(err))
        return
    }
    
    c.JSON(201, SuccessResponse(result))
}
```

**Benefícios:**
- ✅ Lógica de negócio isolada de protocolos externos
- ✅ Fácil trocar HTTP por gRPC
- ✅ Testável sem dependências externas
- ✅ Reutilização de lógica entre adapters

---

## 📦 Tipos de Adapters

### 1. **Driving Adapters** (Primários)
Iniciam a comunicação → Entram na aplicação

```
HTTP Request → HTTP Adapter → Application
gRPC Call    → gRPC Adapter → Application
CLI Command  → CLI Adapter  → Application
```

### 2. **Driven Adapters** (Secundários)
São chamados pela aplicação → Saem para mundo externo

```
Application → Repository Adapter → Database
Application → Email Adapter      → SMTP Server
Application → Queue Adapter      → RabbitMQ
```

---

## 🌐 HTTP Adapter

### Estrutura

```
internal/modules/user/adapters/http/
├── handler.go          # Handlers HTTP
├── dto.go              # Request/Response DTOs
├── routes.go           # Registro de rotas
└── middleware.go       # Middlewares específicos
```

### Exemplo Completo

**`adapters/http/handler.go`**

```go
package http

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "meuApp/internal/modules/user/application/commands"
    "meuApp/internal/modules/user/application/queries"
    "meuApp/pkg/contracts"
)

type UserHTTPHandler struct {
    commandBus contracts.CommandBus
    queryBus   contracts.QueryBus
}

func NewUserHTTPHandler(
    commandBus contracts.CommandBus,
    queryBus contracts.QueryBus,
) *UserHTTPHandler {
    return &UserHTTPHandler{
        commandBus: commandBus,
        queryBus:   queryBus,
    }
}

// CreateUser cria um novo usuário
// @Summary Criar usuário
// @Description Cria um novo usuário no sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "Dados do usuário"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // 1. Parse request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:   "Invalid request body",
            Details: err.Error(),
        })
        return
    }
    
    // 2. Validação
    if err := req.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:   "Validation failed",
            Details: err.Error(),
        })
        return
    }
    
    // 3. Criar comando
    cmd := &commands.CreateUserCommand{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }
    
    // 4. Executar via Command Bus
    result, err := h.commandBus.Execute(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: err.Error(),
        })
        return
    }
    
    // 5. Retornar response
    userID := result.(string)
    c.JSON(http.StatusCreated, CreateUserResponse{
        ID:      userID,
        Message: "User created successfully",
    })
}

// GetUser retorna um usuário por ID
func (h *UserHTTPHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    // 1. Criar query
    query := &queries.GetUserByIDQuery{
        UserID: userID,
    }
    
    // 2. Executar via Query Bus
    result, err := h.queryBus.Execute(c.Request.Context(), query)
    if err != nil {
        if err == queries.ErrUserNotFound {
            c.JSON(http.StatusNotFound, ErrorResponse{
                Error: "User not found",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: err.Error(),
        })
        return
    }
    
    // 3. Converter para DTO
    user := result.(*queries.UserDTO)
    c.JSON(http.StatusOK, UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Active:    user.Active,
        CreatedAt: user.CreatedAt,
    })
}

// ListUsers lista usuários com paginação
func (h *UserHTTPHandler) ListUsers(c *gin.Context) {
    // 1. Parse query params
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "10")
    
    pageInt, _ := strconv.Atoi(page)
    pageSizeInt, _ := strconv.Atoi(pageSize)
    
    // 2. Criar query
    query := &queries.ListUsersQuery{
        Page:     pageInt,
        PageSize: pageSizeInt,
    }
    
    // 3. Executar
    result, err := h.queryBus.Execute(c.Request.Context(), query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: err.Error(),
        })
        return
    }
    
    // 4. Retornar com paginação
    paginatedResult := result.(*queries.PaginatedUsersDTO)
    c.JSON(http.StatusOK, ListUsersResponse{
        Users:      paginatedResult.Users,
        Total:      paginatedResult.Total,
        Page:       paginatedResult.Page,
        PageSize:   paginatedResult.PageSize,
        TotalPages: (paginatedResult.Total + int64(pageSizeInt) - 1) / int64(pageSizeInt),
    })
}

// DeleteUser remove um usuário
func (h *UserHTTPHandler) DeleteUser(c *gin.Context) {
    userID := c.Param("id")
    
    cmd := &commands.DeleteUserCommand{
        UserID: userID,
    }
    
    _, err := h.commandBus.Execute(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusNoContent, nil)
}
```

**`adapters/http/dto.go`**

```go
package http

import (
    "errors"
    "regexp"
    "time"
)

// Request DTOs
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

func (r *CreateUserRequest) Validate() error {
    if len(r.Name) < 3 {
        return errors.New("name must be at least 3 characters")
    }
    
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    if !emailRegex.MatchString(r.Email) {
        return errors.New("invalid email format")
    }
    
    if len(r.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    
    return nil
}

// Response DTOs
type CreateUserResponse struct {
    ID      string `json:"id"`
    Message string `json:"message"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
}

type ListUsersResponse struct {
    Users      []UserResponse `json:"users"`
    Total      int64          `json:"total"`
    Page       int            `json:"page"`
    PageSize   int            `json:"page_size"`
    TotalPages int64          `json:"total_pages"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}
```

**`adapters/http/routes.go`**

```go
package http

import "github.com/gin-gonic/gin"

func (h *UserHTTPHandler) RegisterRoutes(router *gin.Engine) {
    users := router.Group("/api/v1/users")
    {
        users.POST("", h.CreateUser)
        users.GET("/:id", h.GetUser)
        users.GET("", h.ListUsers)
        users.PUT("/:id", h.UpdateUser)
        users.DELETE("/:id", h.DeleteUser)
    }
}
```

---

## 🔄 gRPC Adapter

### Estrutura

```
internal/modules/user/adapters/grpc/
├── handler.go          # Implementação do serviço gRPC
├── mapper.go           # Proto ↔ Domain mapping
└── interceptor.go      # Interceptors específicos
```

### Exemplo Completo

**`proto/user.proto`**

```protobuf
syntax = "proto3";

package user;
option go_package = "meuApp/pkg/proto";

service UserService {
  rpc CreateUser (CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser (GetUserRequest) returns (GetUserResponse);
  rpc ListUsers (ListUsersRequest) returns (ListUsersResponse);
  rpc DeleteUser (DeleteUserRequest) returns (DeleteUserResponse);
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  string password = 3;
}

message CreateUserResponse {
  string id = 1;
  string message = 2;
}

message GetUserRequest {
  string id = 1;
}

message GetUserResponse {
  string id = 1;
  string name = 2;
  string email = 3;
  bool active = 4;
  string created_at = 5;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListUsersResponse {
  repeated GetUserResponse users = 1;
  int64 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message DeleteUserRequest {
  string id = 1;
}

message DeleteUserResponse {
  string message = 1;
}
```

**`adapters/grpc/handler.go`**

```go
package grpc

import (
    "context"
    pb "meuApp/pkg/proto"
    "meuApp/internal/modules/user/application/commands"
    "meuApp/internal/modules/user/application/queries"
    "meuApp/pkg/contracts"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type UserGRPCHandler struct {
    pb.UnimplementedUserServiceServer
    commandBus contracts.CommandBus
    queryBus   contracts.QueryBus
}

func NewUserGRPCHandler(
    commandBus contracts.CommandBus,
    queryBus contracts.QueryBus,
) *UserGRPCHandler {
    return &UserGRPCHandler{
        commandBus: commandBus,
        queryBus:   queryBus,
    }
}

// CreateUser implementa UserServiceServer
func (h *UserGRPCHandler) CreateUser(
    ctx context.Context,
    req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
    // 1. Validação
    if err := validateCreateUserRequest(req); err != nil {
        return nil, status.Errorf(codes.InvalidArgument, err.Error())
    }
    
    // 2. Criar comando
    cmd := &commands.CreateUserCommand{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }
    
    // 3. Executar
    result, err := h.commandBus.Execute(ctx, cmd)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    
    // 4. Retornar response
    userID := result.(string)
    return &pb.CreateUserResponse{
        Id:      userID,
        Message: "User created successfully",
    }, nil
}

// GetUser busca usuário por ID
func (h *UserGRPCHandler) GetUser(
    ctx context.Context,
    req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
    query := &queries.GetUserByIDQuery{
        UserID: req.Id,
    }
    
    result, err := h.queryBus.Execute(ctx, query)
    if err != nil {
        if err == queries.ErrUserNotFound {
            return nil, status.Errorf(codes.NotFound, "User not found")
        }
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    
    user := result.(*queries.UserDTO)
    return &pb.GetUserResponse{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Active:    user.Active,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
    }, nil
}

// ListUsers lista usuários com paginação
func (h *UserGRPCHandler) ListUsers(
    ctx context.Context,
    req *pb.ListUsersRequest,
) (*pb.ListUsersResponse, error) {
    query := &queries.ListUsersQuery{
        Page:     int(req.Page),
        PageSize: int(req.PageSize),
    }
    
    result, err := h.queryBus.Execute(ctx, query)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }
    
    paginatedResult := result.(*queries.PaginatedUsersDTO)
    
    // Converter para proto
    users := make([]*pb.GetUserResponse, len(paginatedResult.Users))
    for i, user := range paginatedResult.Users {
        users[i] = &pb.GetUserResponse{
            Id:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            Active:    user.Active,
            CreatedAt: user.CreatedAt.Format(time.RFC3339),
        }
    }
    
    return &pb.ListUsersResponse{
        Users:    users,
        Total:    paginatedResult.Total,
        Page:     int32(paginatedResult.Page),
        PageSize: int32(paginatedResult.PageSize),
    }, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    if req.Email == "" {
        return errors.New("email is required")
    }
    if len(req.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}
```

---

## 💾 Database Adapter

### Estrutura

```
pkg/adapters/database/
├── mysql.go            # MySQL adapter
├── postgres.go         # PostgreSQL adapter
├── mongodb.go          # MongoDB adapter
└── connection_pool.go  # Pool de conexões
```

### MySQL Adapter

**`pkg/adapters/database/mysql.go`**

```go
package database

import (
    "fmt"
    "time"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type MySQLConfig struct {
    Host            string
    Port            int
    Database        string
    Username        string
    Password        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    LogLevel        logger.LogLevel
}

type MySQLAdapter struct {
    db     *gorm.DB
    config MySQLConfig
}

func NewMySQLAdapter(config MySQLConfig) (*MySQLAdapter, error) {
    // DSN
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        config.Username,
        config.Password,
        config.Host,
        config.Port,
        config.Database,
    )
    
    // GORM Config
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(config.LogLevel),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    }
    
    // Connect
    db, err := gorm.Open(mysql.Open(dsn), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
    }
    
    // Connection Pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxOpenConns(config.MaxOpenConns)
    sqlDB.SetMaxIdleConns(config.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
    
    // Health Check
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping MySQL: %w", err)
    }
    
    return &MySQLAdapter{
        db:     db,
        config: config,
    }, nil
}

func (a *MySQLAdapter) DB() *gorm.DB {
    return a.db
}

func (a *MySQLAdapter) Close() error {
    sqlDB, err := a.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func (a *MySQLAdapter) HealthCheck() error {
    sqlDB, err := a.db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}

// Transaction helpers
func (a *MySQLAdapter) Transaction(fn func(*gorm.DB) error) error {
    return a.db.Transaction(fn)
}

func (a *MySQLAdapter) BeginTx() *gorm.DB {
    return a.db.Begin()
}
```

---

## 🆕 Criando Novos Adapters

### Exemplo: Email Adapter

**`pkg/adapters/email/smtp_adapter.go`**

```go
package email

import (
    "fmt"
    "net/smtp"
)

// Interface (Port)
type EmailSender interface {
    SendEmail(to, subject, body string) error
}

// SMTP Adapter
type SMTPAdapter struct {
    host     string
    port     int
    username string
    password string
    from     string
}

func NewSMTPAdapter(host string, port int, username, password, from string) *SMTPAdapter {
    return &SMTPAdapter{
        host:     host,
        port:     port,
        username: username,
        password: password,
        from:     from,
    }
}

func (a *SMTPAdapter) SendEmail(to, subject, body string) error {
    // Configurar autenticação
    auth := smtp.PlainAuth("", a.username, a.password, a.host)
    
    // Montar mensagem
    msg := []byte(fmt.Sprintf(
        "From: %s\r\n"+
            "To: %s\r\n"+
            "Subject: %s\r\n"+
            "\r\n"+
            "%s\r\n",
        a.from, to, subject, body,
    ))
    
    // Enviar
    addr := fmt.Sprintf("%s:%d", a.host, a.port)
    return smtp.SendMail(addr, auth, a.from, []string{to}, msg)
}
```

### Exemplo: Queue Adapter (RabbitMQ)

**`pkg/adapters/queue/rabbitmq_adapter.go`**

```go
package queue

import (
    "github.com/streadway/amqp"
)

type QueuePublisher interface {
    Publish(queue string, message []byte) error
}

type RabbitMQAdapter struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQAdapter(url string) (*RabbitMQAdapter, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }
    
    channel, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    
    return &RabbitMQAdapter{
        conn:    conn,
        channel: channel,
    }, nil
}

func (a *RabbitMQAdapter) Publish(queue string, message []byte) error {
    // Declarar fila
    _, err := a.channel.QueueDeclare(
        queue, // name
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        return err
    }
    
    // Publicar mensagem
    return a.channel.Publish(
        "",    // exchange
        queue, // routing key
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        message,
        },
    )
}

func (a *RabbitMQAdapter) Close() error {
    if err := a.channel.Close(); err != nil {
        return err
    }
    return a.conn.Close()
}
```

---

## 🎯 Boas Práticas

### 1. Separe Concerns

```go
// ✅ BOM: Adapter só traduz protocolos
type HTTPAdapter struct {
    commandBus CommandBus // Delega lógica
}

func (a *HTTPAdapter) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.ShouldBindJSON(&req) // HTTP concern
    
    cmd := &CreateUserCommand{...} // Domain concern
    a.commandBus.Execute(ctx, cmd)  // Application concern
}

// ❌ RUIM: Lógica no adapter
type HTTPAdapter struct {
    repo UserRepository
}

func (a *HTTPAdapter) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.ShouldBindJSON(&req)
    
    // ❌ Lógica de negócio no adapter
    if !isValidEmail(req.Email) {
        c.JSON(400, "invalid email")
        return
    }
    
    user := &User{...}
    a.repo.Save(user) // ❌ Acesso direto ao repo
}
```

### 2. Use DTOs Específicos

```go
// ✅ DTO específico para HTTP
type CreateUserHTTPRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// ✅ DTO específico para gRPC (gerado automaticamente)
type CreateUserGRPCRequest struct {
    Name     string
    Email    string
    Password string
}

// ✅ Command interno (domain)
type CreateUserCommand struct {
    Name     string
    Email    string
    Password string
}
```

### 3. Trate Erros Adequadamente

```go
func (h *HTTPAdapter) CreateUser(c *gin.Context) {
    result, err := h.commandBus.Execute(ctx, cmd)
    
    if err != nil {
        // Traduzir erro de domínio para HTTP
        switch err {
        case domain.ErrUserAlreadyExists:
            c.JSON(http.StatusConflict, ErrorResponse{
                Error: "User already exists",
            })
        case domain.ErrInvalidEmail:
            c.JSON(http.StatusBadRequest, ErrorResponse{
                Error: "Invalid email format",
            })
        default:
            c.JSON(http.StatusInternalServerError, ErrorResponse{
                Error: "Internal server error",
            })
        }
        return
    }
    
    c.JSON(http.StatusCreated, result)
}
```

### 4. Use Middleware para Concerns Transversais

```go
// Middleware de autenticação
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        
        // Validar token
        userID, err := validateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }
        
        c.Set("user_id", userID)
        c.Next()
    }
}

// Usar no router
router.POST("/users", AuthMiddleware(), handler.CreateUser)
```

### 5. Implemente Health Checks

```go
type HealthCheckAdapter struct {
    db    *MySQLAdapter
    cache *RedisAdapter
}

func (h *HealthCheckAdapter) Check(c *gin.Context) {
    checks := map[string]string{
        "database": "healthy",
        "cache":    "healthy",
    }
    
    // Check database
    if err := h.db.HealthCheck(); err != nil {
        checks["database"] = "unhealthy: " + err.Error()
    }
    
    // Check cache
    if err := h.cache.Ping(); err != nil {
        checks["cache"] = "unhealthy: " + err.Error()
    }
    
    // Status geral
    status := "healthy"
    for _, check := range checks {
        if check != "healthy" {
            status = "unhealthy"
            break
        }
    }
    
    statusCode := http.StatusOK
    if status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    c.JSON(statusCode, gin.H{
        "status": status,
        "checks": checks,
    })
}
```

---

## 📚 Próximos Passos

- **[Contratos e Interfaces](10-contracts-interfaces.md)** - Como definir ports
- **[Repositórios](11-repositories.md)** - Database adapters específicos
- **[Criando Módulos](12-creating-modules.md)** - Como registrar adapters

---

**[⬅️ Eventos](08-events-system.md)** | **[Índice](README.md)** | **[Contratos ➡️](10-contracts-interfaces.md)**
