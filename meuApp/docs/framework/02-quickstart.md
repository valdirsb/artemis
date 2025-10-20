# 🚀 Guia de Início Rápido

## 📋 Índice
- [Configuração Inicial](#configuração-inicial)
- [Criando Seu Primeiro Módulo](#criando-seu-primeiro-módulo)
- [Implementando Funcionalidades](#implementando-funcionalidades)
- [Testando](#testando)
- [Executando a Aplicação](#executando-a-aplicação)

---

## ⚙️ Configuração Inicial

### Pré-requisitos

```bash
✅ Go 1.21 ou superior
✅ MySQL 8.0+ (ou outro banco suportado)
✅ Git
✅ Make (opcional, mas recomendado)
```

### 1. Clone o Framework

```bash
# Clone o repositório
git clone https://github.com/sua-agencia/artemis.git meu-projeto
cd meu-projeto

# Instale as dependências
go mod download
```

### 2. Configure as Variáveis de Ambiente

```bash
# Copie o arquivo de exemplo
cp .env.example .env

# Edite com suas configurações
nano .env
```

**.env** exemplo:
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=secret
DB_DATABASE=artemis

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# JWT
JWT_SECRET=your-super-secret-key-change-this
JWT_EXPIRATION=24h

# Environment
ENVIRONMENT=development
DEBUG_MODE=true
```

### 3. Configure o Framework

Edite `framework.yaml` para habilitar os recursos que você precisa:

```yaml
framework:
  name: "meu-projeto"
  version: "1.0.0"

core:
  http_server: true
  dependency_injection: true
  event_system: true
  logging: true
  database_migrations: true

database:
  mysql: true        # ou postgresql, mongodb, sqlite
  
protocols:
  http: true         # REST API
  grpc: true         # gRPC API
  
modules:
  user: true
  product: true
  order: true
  # Adicione seus módulos aqui
```

### 4. Execute as Migrações

```bash
# Via Makefile (recomendado)
make migrate

# Ou diretamente
go run main.go migrate
```

### 5. Verifique a Instalação

```bash
# Execute a aplicação
go run main.go

# Em outro terminal, teste o health check
curl http://localhost:8080/health
```

Resposta esperada:
```json
{
  "status": "healthy",
  "services": {
    "database": null,
    "eventbus": null
  },
  "registry": {
    "http_handlers": 3,
    "grpc_services": 3,
    "repositories": 3,
    "app_services": 3
  }
}
```

🎉 **Pronto! Seu framework está configurado e rodando!**

---

## 📦 Criando Seu Primeiro Módulo

Vamos criar um módulo **Blog** completo do zero.

### Estrutura do Módulo

```
internal/modules/blog/
├── domain/                    # Camada de Domínio
│   ├── entities/
│   │   └── post.go
│   ├── value_objects/
│   │   └── content.go
│   └── events/
│       └── post_created.go
│
├── application/               # Camada de Aplicação
│   ├── commands/             # Write operations
│   │   ├── create_post.go
│   │   ├── update_post.go
│   │   └── delete_post.go
│   ├── queries/              # Read operations
│   │   ├── get_post.go
│   │   └── list_posts.go
│   └── services/
│       └── post_service.go
│
├── adapters/                 # Camada de Apresentação
│   ├── http/
│   │   └── handler.go
│   └── grpc/
│       └── service.go
│
├── repository/               # Camada de Infraestrutura
│   └── mysql_repository.go
│
├── dto/                      # Data Transfer Objects
│   ├── request.go
│   └── response.go
│
└── blog_module.go           # Módulo principal
```

### Passo 1: Criar a Entidade de Domínio

**`internal/modules/blog/domain/entities/post.go`**

```go
package entities

import (
    "time"
)

// Post é a entidade principal do domínio
type Post struct {
    ID        string
    Title     string
    Content   string
    AuthorID  string
    Status    PostStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

type PostStatus string

const (
    PostStatusDraft     PostStatus = "draft"
    PostStatusPublished PostStatus = "published"
    PostStatusArchived  PostStatus = "archived"
)

// NewPost cria um novo post com validações
func NewPost(title, content, authorID string) (*Post, error) {
    if title == "" {
        return nil, errors.New("title cannot be empty")
    }
    if content == "" {
        return nil, errors.New("content cannot be empty")
    }
    if authorID == "" {
        return nil, errors.New("author ID cannot be empty")
    }

    return &Post{
        ID:        generateID(),
        Title:     title,
        Content:   content,
        AuthorID:  authorID,
        Status:    PostStatusDraft,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}

// Publish publica o post
func (p *Post) Publish() error {
    if p.Status == PostStatusPublished {
        return errors.New("post is already published")
    }
    p.Status = PostStatusPublished
    p.UpdatedAt = time.Now()
    return nil
}

// Update atualiza o post
func (p *Post) Update(title, content string) error {
    if title != "" {
        p.Title = title
    }
    if content != "" {
        p.Content = content
    }
    p.UpdatedAt = time.Now()
    return nil
}
```

### Passo 2: Criar o Repository (Interface)

**`internal/modules/blog/repository/repository.go`**

```go
package repository

import (
    "context"
    "meuApp/internal/modules/blog/domain/entities"
)

// PostRepository define as operações de persistência
type PostRepository interface {
    Save(ctx context.Context, post *entities.Post) error
    FindByID(ctx context.Context, id string) (*entities.Post, error)
    FindAll(ctx context.Context, page, pageSize int) ([]*entities.Post, int64, error)
    Update(ctx context.Context, post *entities.Post) error
    Delete(ctx context.Context, id string) error
}
```

### Passo 3: Implementar o Repository (MySQL)

**`internal/modules/blog/repository/mysql_repository.go`**

```go
package repository

import (
    "context"
    "meuApp/internal/modules/blog/domain/entities"
    "gorm.io/gorm"
)

type MySQLPostRepository struct {
    db *gorm.DB
}

func NewMySQLPostRepository(db *gorm.DB) *MySQLPostRepository {
    return &MySQLPostRepository{db: db}
}

// PostModel é o modelo do banco de dados
type PostModel struct {
    ID        string `gorm:"primaryKey"`
    Title     string
    Content   string
    AuthorID  string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (PostModel) TableName() string {
    return "posts"
}

func (r *MySQLPostRepository) Save(ctx context.Context, post *entities.Post) error {
    model := &PostModel{
        ID:        post.ID,
        Title:     post.Title,
        Content:   post.Content,
        AuthorID:  post.AuthorID,
        Status:    string(post.Status),
        CreatedAt: post.CreatedAt,
        UpdatedAt: post.UpdatedAt,
    }
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *MySQLPostRepository) FindByID(ctx context.Context, id string) (*entities.Post, error) {
    var model PostModel
    if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return r.toEntity(&model), nil
}

func (r *MySQLPostRepository) FindAll(ctx context.Context, page, pageSize int) ([]*entities.Post, int64, error) {
    var models []PostModel
    var total int64

    offset := (page - 1) * pageSize

    if err := r.db.WithContext(ctx).Model(&PostModel{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := r.db.WithContext(ctx).
        Offset(offset).
        Limit(pageSize).
        Find(&models).Error; err != nil {
        return nil, 0, err
    }

    posts := make([]*entities.Post, len(models))
    for i, model := range models {
        posts[i] = r.toEntity(&model)
    }

    return posts, total, nil
}

func (r *MySQLPostRepository) toEntity(model *PostModel) *entities.Post {
    return &entities.Post{
        ID:        model.ID,
        Title:     model.Title,
        Content:   model.Content,
        AuthorID:  model.AuthorID,
        Status:    entities.PostStatus(model.Status),
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
    }
}
```

### Passo 4: Criar Command Handler

**`internal/modules/blog/application/commands/create_post.go`**

```go
package commands

import (
    "context"
    "meuApp/internal/modules/blog/domain/entities"
    "meuApp/internal/modules/blog/repository"
    "meuApp/pkg/events"
)

type CreatePostCommand struct {
    Title    string
    Content  string
    AuthorID string
}

type CreatePostHandler struct {
    repo     repository.PostRepository
    eventBus *events.EventBus
}

func NewCreatePostHandler(repo repository.PostRepository, eventBus *events.EventBus) *CreatePostHandler {
    return &CreatePostHandler{
        repo:     repo,
        eventBus: eventBus,
    }
}

func (h *CreatePostHandler) Handle(ctx context.Context, cmd *CreatePostCommand) (*entities.Post, error) {
    // 1. Criar entidade de domínio
    post, err := entities.NewPost(cmd.Title, cmd.Content, cmd.AuthorID)
    if err != nil {
        return nil, err
    }

    // 2. Salvar no repositório
    if err := h.repo.Save(ctx, post); err != nil {
        return nil, err
    }

    // 3. Publicar evento
    h.eventBus.Publish("post.created", map[string]interface{}{
        "post_id":   post.ID,
        "author_id": post.AuthorID,
        "title":     post.Title,
    })

    return post, nil
}
```

### Passo 5: Criar Query Handler

**`internal/modules/blog/application/queries/list_posts.go`**

```go
package queries

import (
    "context"
    "meuApp/internal/modules/blog/domain/entities"
    "meuApp/internal/modules/blog/repository"
)

type ListPostsQuery struct {
    Page     int
    PageSize int
}

type ListPostsHandler struct {
    repo repository.PostRepository
}

func NewListPostsHandler(repo repository.PostRepository) *ListPostsHandler {
    return &ListPostsHandler{repo: repo}
}

type ListPostsResponse struct {
    Posts []*entities.Post
    Total int64
    Page  int
}

func (h *ListPostsHandler) Handle(ctx context.Context, query *ListPostsQuery) (*ListPostsResponse, error) {
    posts, total, err := h.repo.FindAll(ctx, query.Page, query.PageSize)
    if err != nil {
        return nil, err
    }

    return &ListPostsResponse{
        Posts: posts,
        Total: total,
        Page:  query.Page,
    }, nil
}
```

### Passo 6: Criar Application Service

**`internal/modules/blog/application/services/post_service.go`**

```go
package services

import (
    "context"
    "meuApp/internal/modules/blog/application/commands"
    "meuApp/internal/modules/blog/application/queries"
    "meuApp/internal/modules/blog/domain/entities"
)

type PostApplicationService struct {
    createHandler *commands.CreatePostHandler
    listHandler   *queries.ListPostsHandler
}

func NewPostApplicationService(
    createHandler *commands.CreatePostHandler,
    listHandler *queries.ListPostsHandler,
) *PostApplicationService {
    return &PostApplicationService{
        createHandler: createHandler,
        listHandler:   listHandler,
    }
}

func (s *PostApplicationService) CreatePost(ctx context.Context, cmd *commands.CreatePostCommand) (*entities.Post, error) {
    return s.createHandler.Handle(ctx, cmd)
}

func (s *PostApplicationService) ListPosts(ctx context.Context, query *queries.ListPostsQuery) (*queries.ListPostsResponse, error) {
    return s.listHandler.Handle(ctx, query)
}
```

### Passo 7: Criar HTTP Handler

**`internal/modules/blog/adapters/http/handler.go`**

```go
package http

import (
    "net/http"
    "strconv"
    
    "meuApp/internal/modules/blog/application/commands"
    "meuApp/internal/modules/blog/application/queries"
    "meuApp/internal/modules/blog/application/services"
    
    "github.com/gin-gonic/gin"
)

type PostHTTPHandler struct {
    service *services.PostApplicationService
}

func NewPostHTTPHandler(service *services.PostApplicationService) *PostHTTPHandler {
    return &PostHTTPHandler{service: service}
}

func (h *PostHTTPHandler) RegisterRoutes(router *gin.RouterGroup) {
    posts := router.Group("/posts")
    {
        posts.POST("", h.CreatePost)
        posts.GET("", h.ListPosts)
    }
}

type CreatePostRequest struct {
    Title    string `json:"title" binding:"required"`
    Content  string `json:"content" binding:"required"`
    AuthorID string `json:"author_id" binding:"required"`
}

func (h *PostHTTPHandler) CreatePost(c *gin.Context) {
    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    cmd := &commands.CreatePostCommand{
        Title:    req.Title,
        Content:  req.Content,
        AuthorID: req.AuthorID,
    }

    post, err := h.service.CreatePost(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, post)
}

func (h *PostHTTPHandler) ListPosts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

    query := &queries.ListPostsQuery{
        Page:     page,
        PageSize: pageSize,
    }

    result, err := h.service.ListPosts(c.Request.Context(), query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, result)
}
```

### Passo 8: Criar o Módulo

**`internal/modules/blog_module.go`**

```go
package modules

import (
    "fmt"
    
    "meuApp/internal/modules/blog/adapters/http"
    "meuApp/internal/modules/blog/application/commands"
    "meuApp/internal/modules/blog/application/queries"
    "meuApp/internal/modules/blog/application/services"
    "meuApp/internal/modules/blog/repository"
    "meuApp/pkg/container"
    "meuApp/pkg/events"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type BlogModule struct {
    db       *gorm.DB
    eventBus *events.EventBus
}

func NewBlogModule(db *gorm.DB, eventBus *events.EventBus) *BlogModule {
    return &BlogModule{
        db:       db,
        eventBus: eventBus,
    }
}

func (m *BlogModule) Name() string {
    return "blog"
}

func (m *BlogModule) Register(registry *container.ModuleRegistry) error {
    fmt.Printf("🔧 Registering module: %s\n", m.Name())

    // 1. Registrar Repository
    postRepo := repository.NewMySQLPostRepository(m.db)
    registry.RegisterRepository("post", postRepo)

    // 2. Criar Command Handlers
    createPostHandler := commands.NewCreatePostHandler(postRepo, m.eventBus)

    // 3. Criar Query Handlers
    listPostsHandler := queries.NewListPostsHandler(postRepo)

    // 4. Registrar Application Service
    postService := services.NewPostApplicationService(createPostHandler, listPostsHandler)
    registry.RegisterApplicationService("post", postService)

    // 5. Registrar HTTP Handler
    httpHandler := http.NewPostHTTPHandler(postService)
    registry.RegisterHTTPHandler("post", &postHTTPHandlerAdapter{handler: httpHandler})

    fmt.Printf("✅ Module %s registered successfully\n", m.Name())
    return nil
}

// Adapter para interface do registry
type postHTTPHandlerAdapter struct {
    handler *http.PostHTTPHandler
}

func (a *postHTTPHandlerAdapter) RegisterRoutes(router *gin.RouterGroup) {
    a.handler.RegisterRoutes(router)
}
```

### Passo 9: Registrar o Módulo no Bootstrap

**`internal/bootstrap/bootstrap_registry.go`**

```go
// Adicione o módulo na função registerModules
func registerModules(registry *container.ModuleRegistry, db *gorm.DB, eventBus *events.EventBus, logger contracts.Logger) error {
    log.Println("📦 Registering modules...")

    applicationModules := []frameworkInterfaces.Module{
        modules.NewUserModule(db, eventBus),
        modules.NewProductModule(db, eventBus, logger),
        modules.NewOrderModule(db, eventBus, logger),
        modules.NewBlogModule(db, eventBus), // ← SEU NOVO MÓDULO
    }

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

### Passo 10: Atualizar framework.yaml

```yaml
modules:
  user: true
  product: true
  order: true
  blog: true  # ← Adicione seu módulo
```

---

## 🧪 Testando

### Criar Teste Unitário

**`internal/modules/blog/application/commands/create_post_test.go`**

```go
package commands_test

import (
    "context"
    "testing"
    
    "meuApp/internal/modules/blog/application/commands"
    "meuApp/internal/modules/blog/repository"
    "meuApp/pkg/events"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockPostRepository struct {
    mock.Mock
}

func (m *MockPostRepository) Save(ctx context.Context, post *entities.Post) error {
    args := m.Called(ctx, post)
    return args.Error(0)
}

func TestCreatePostHandler_Handle(t *testing.T) {
    // Arrange
    mockRepo := new(MockPostRepository)
    eventBus := events.NewEventBus()
    handler := commands.NewCreatePostHandler(mockRepo, eventBus)

    cmd := &commands.CreatePostCommand{
        Title:    "Test Post",
        Content:  "Test Content",
        AuthorID: "author123",
    }

    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

    // Act
    post, err := handler.Handle(context.Background(), cmd)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, post)
    assert.Equal(t, "Test Post", post.Title)
    mockRepo.AssertExpectations(t)
}
```

### Executar Testes

```bash
# Testar módulo específico
go test ./internal/modules/blog/...

# Testar tudo
go test ./...

# Com coverage
go test -cover ./...
```

---

## ▶️ Executando a Aplicação

```bash
# Executar em modo desenvolvimento
go run main.go

# Ou com hot reload (instale air primeiro: go install github.com/cosmtrek/air@latest)
air
```

### Testando os Endpoints

```bash
# Criar um post
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Meu Primeiro Post",
    "content": "Conteúdo do post",
    "author_id": "user123"
  }'

# Listar posts
curl http://localhost:8080/api/v1/posts?page=1&page_size=10
```

---

## 🎉 Parabéns!

Você criou seu primeiro módulo completo no Artemis Framework! 

### 📚 Próximos Passos

1. **[Estrutura do Projeto](03-project-structure.md)** - Entenda melhor a organização
2. **[Sistema de Módulos](06-modules-system.md)** - Conceitos avançados de módulos
3. **[Implementando Commands](13-implementing-commands.md)** - Aprofunde em Commands
4. **[Trabalhando com Eventos](15-working-with-events.md)** - Sistema de eventos

---

**[⬅️ Visão Geral](01-overview.md)** | **[Índice](README.md)** | **[Estrutura do Projeto ➡️](03-project-structure.md)**
