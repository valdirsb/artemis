# 🏗️ Criando Módulos - Guia Completo

## 📋 Índice
- [🏗️ Criando Módulos - Guia Completo](#️-criando-módulos---guia-completo)
  - [📋 Índice](#-índice)
  - [🎯 Introdução](#-introdução)
    - [O que vamos criar?](#o-que-vamos-criar)
    - [Arquitetura do Módulo](#arquitetura-do-módulo)
  - [📁 Passo 1: Estrutura de Diretórios](#-passo-1-estrutura-de-diretórios)
    - [Criar estrutura completa](#criar-estrutura-completa)
    - [Estrutura final](#estrutura-final)
  - [🎨 Passo 2: Domain Layer](#-passo-2-domain-layer)
    - [2.1 Entidade Category](#21-entidade-category)
    - [2.2 Domain Events](#22-domain-events)
  - [⚡ Passo 3: Application Layer](#-passo-3-application-layer)
    - [3.1 Commands](#31-commands)
    - [3.2 Queries](#32-queries)
  - [💾 Passo 4: Repository](#-passo-4-repository)
    - [4.1 Interface](#41-interface)
    - [4.2 Implementação MySQL](#42-implementação-mysql)
  - [🔌 Passo 5: Adapters](#-passo-5-adapters)
    - [5.1 HTTP Adapter](#51-http-adapter)
  - [🎯 Passo 6: Registro do Módulo](#-passo-6-registro-do-módulo)
    - [Registrar no Bootstrap](#registrar-no-bootstrap)
  - [🧪 Passo 7: Testes](#-passo-7-testes)
    - [7.1 Teste Unitário do Handler](#71-teste-unitário-do-handler)
    - [7.2 Teste de Integração](#72-teste-de-integração)
  - [✅ Checklist de Criação de Módulo](#-checklist-de-criação-de-módulo)
  - [📚 Próximos Passos](#-próximos-passos)

---

## 🎯 Introdução

Este guia apresenta um **tutorial passo a passo** para criar um novo módulo no Artemis Framework, desde a estrutura de diretórios até testes.

### O que vamos criar?

Um módulo **Category** (Categorias) para um e-commerce, com funcionalidades:
- ✅ Criar categoria
- ✅ Listar categorias
- ✅ Buscar categoria por ID
- ✅ Atualizar categoria
- ✅ Deletar categoria
- ✅ Listar produtos de uma categoria

### Arquitetura do Módulo

```
┌─────────────────────────────────────────────────┐
│              HTTP/gRPC Adapters                 │
│         (Controllers, Handlers)                 │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│           Application Layer                     │
│   Commands, Queries, Handlers, Services        │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│              Domain Layer                       │
│   Entities, Value Objects, Domain Events       │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│           Infrastructure Layer                  │
│      Repository (MySQL, PostgreSQL, etc.)      │
└─────────────────────────────────────────────────┘
```

---

## 📁 Passo 1: Estrutura de Diretórios

### Criar estrutura completa

```bash
# Criar diretórios do módulo
mkdir -p internal/modules/category/{domain,application,adapters,repository,dto}
mkdir -p internal/modules/category/domain/{entities,events,valueobjects}
mkdir -p internal/modules/category/application/{commands,queries,services}
mkdir -p internal/modules/category/adapters/{http,grpc,subscribers}
```

### Estrutura final

```
internal/modules/category/
├── domain/
│   ├── entities/
│   │   └── category.go
│   ├── events/
│   │   ├── category_created_event.go
│   │   ├── category_updated_event.go
│   │   └── category_deleted_event.go
│   └── valueobjects/
│       └── category_status.go
├── application/
│   ├── commands/
│   │   ├── create_category_command.go
│   │   ├── create_category_handler.go
│   │   ├── update_category_command.go
│   │   ├── update_category_handler.go
│   │   ├── delete_category_command.go
│   │   └── delete_category_handler.go
│   ├── queries/
│   │   ├── get_category_query.go
│   │   ├── get_category_handler.go
│   │   ├── list_categories_query.go
│   │   └── list_categories_handler.go
│   └── services/
│       └── category_service.go
├── repository/
│   ├── category_repository.go          # Interface
│   └── mysql_category_repository.go    # Implementação
├── adapters/
│   ├── http/
│   │   ├── category_handler.go
│   │   ├── dto.go
│   │   └── routes.go
│   ├── grpc/
│   │   └── category_grpc_handler.go
│   └── subscribers/
│       └── category_event_subscriber.go
├── dto/
│   └── category_dto.go
└── category_module.go                   # Módulo principal
```

---

## 🎨 Passo 2: Domain Layer

### 2.1 Entidade Category

**`domain/entities/category.go`**

```go
package entities

import (
    "errors"
    "time"
)

type Category struct {
    ID          string
    Name        string
    Slug        string
    Description string
    ParentID    *string    // Categoria pai (opcional)
    Active      bool
    SortOrder   int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// NewCategory cria uma nova categoria
func NewCategory(name, slug, description string, parentID *string) (*Category, error) {
    if err := validateCategoryData(name, slug); err != nil {
        return nil, err
    }
    
    return &Category{
        ID:          generateID(),
        Name:        name,
        Slug:        slug,
        Description: description,
        ParentID:    parentID,
        Active:      true,
        SortOrder:   0,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }, nil
}

// Update atualiza dados da categoria
func (c *Category) Update(name, slug, description string) error {
    if err := validateCategoryData(name, slug); err != nil {
        return err
    }
    
    c.Name = name
    c.Slug = slug
    c.Description = description
    c.UpdatedAt = time.Now()
    
    return nil
}

// Activate ativa a categoria
func (c *Category) Activate() {
    c.Active = true
    c.UpdatedAt = time.Now()
}

// Deactivate desativa a categoria
func (c *Category) Deactivate() {
    c.Active = false
    c.UpdatedAt = time.Now()
}

// SetSortOrder define ordem de exibição
func (c *Category) SetSortOrder(order int) {
    c.SortOrder = order
    c.UpdatedAt = time.Now()
}

// Validações
func validateCategoryData(name, slug string) error {
    if name == "" {
        return errors.New("category name is required")
    }
    if len(name) < 3 {
        return errors.New("category name must be at least 3 characters")
    }
    if slug == "" {
        return errors.New("category slug is required")
    }
    if !isValidSlug(slug) {
        return errors.New("category slug must contain only lowercase letters, numbers and hyphens")
    }
    return nil
}

func isValidSlug(slug string) bool {
    // Regex: ^[a-z0-9-]+$
    for _, char := range slug {
        if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
            return false
        }
    }
    return true
}

func generateID() string {
    return uuid.New().String()
}
```

### 2.2 Domain Events

**`domain/events/category_created_event.go`**

```go
package events

import "time"

type CategoryCreatedEvent struct {
    CategoryID  string
    Name        string
    Slug        string
    Description string
    CreatedBy   string
    occurredAt  time.Time
}

func NewCategoryCreatedEvent(categoryID, name, slug, description, createdBy string) *CategoryCreatedEvent {
    return &CategoryCreatedEvent{
        CategoryID:  categoryID,
        Name:        name,
        Slug:        slug,
        Description: description,
        CreatedBy:   createdBy,
        occurredAt:  time.Now(),
    }
}

func (e *CategoryCreatedEvent) EventName() string {
    return "CategoryCreated"
}

func (e *CategoryCreatedEvent) OccurredAt() time.Time {
    return e.occurredAt
}
```

**`domain/events/category_updated_event.go`**

```go
package events

import "time"

type CategoryUpdatedEvent struct {
    CategoryID  string
    Name        string
    Slug        string
    UpdatedBy   string
    occurredAt  time.Time
}

func NewCategoryUpdatedEvent(categoryID, name, slug, updatedBy string) *CategoryUpdatedEvent {
    return &CategoryUpdatedEvent{
        CategoryID: categoryID,
        Name:       name,
        Slug:       slug,
        UpdatedBy:  updatedBy,
        occurredAt: time.Now(),
    }
}

func (e *CategoryUpdatedEvent) EventName() string {
    return "CategoryUpdated"
}

func (e *CategoryUpdatedEvent) OccurredAt() time.Time {
    return e.occurredAt
}
```

---

## ⚡ Passo 3: Application Layer

### 3.1 Commands

**`application/commands/create_category_command.go`**

```go
package commands

type CreateCategoryCommand struct {
    Name        string
    Slug        string
    Description string
    ParentID    *string
    CreatedBy   string
}

func (c *CreateCategoryCommand) CommandName() string {
    return "CreateCategoryCommand"
}
```

**`application/commands/create_category_handler.go`**

```go
package commands

import (
    "context"
    "meuApp/internal/modules/category/domain/entities"
    "meuApp/internal/modules/category/domain/events"
    "meuApp/internal/modules/category/repository"
    "meuApp/pkg/contracts"
)

type CreateCategoryHandler struct {
    repo     repository.CategoryRepository
    eventBus contracts.EventBus
}

func NewCreateCategoryHandler(
    repo repository.CategoryRepository,
    eventBus contracts.EventBus,
) *CreateCategoryHandler {
    return &CreateCategoryHandler{
        repo:     repo,
        eventBus: eventBus,
    }
}

func (h *CreateCategoryHandler) Handle(
    ctx context.Context,
    cmd contracts.Command,
) (interface{}, error) {
    createCmd := cmd.(*CreateCategoryCommand)
    
    // 1. Verificar se slug já existe
    exists, err := h.repo.ExistsBySlug(ctx, createCmd.Slug)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("category slug already exists")
    }
    
    // 2. Criar entidade
    category, err := entities.NewCategory(
        createCmd.Name,
        createCmd.Slug,
        createCmd.Description,
        createCmd.ParentID,
    )
    if err != nil {
        return nil, err
    }
    
    // 3. Persistir
    if err := h.repo.Save(ctx, category); err != nil {
        return nil, err
    }
    
    // 4. Publicar evento
    event := events.NewCategoryCreatedEvent(
        category.ID,
        category.Name,
        category.Slug,
        category.Description,
        createCmd.CreatedBy,
    )
    h.eventBus.Publish(ctx, event)
    
    return category.ID, nil
}
```

### 3.2 Queries

**`application/queries/get_category_query.go`**

```go
package queries

type GetCategoryQuery struct {
    CategoryID string
}

func (q *GetCategoryQuery) QueryName() string {
    return "GetCategoryQuery"
}
```

**`application/queries/get_category_handler.go`**

```go
package queries

import (
    "context"
    "meuApp/internal/modules/category/dto"
    "meuApp/internal/modules/category/repository"
    "meuApp/pkg/contracts"
)

type GetCategoryHandler struct {
    repo repository.CategoryRepository
}

func NewGetCategoryHandler(repo repository.CategoryRepository) *GetCategoryHandler {
    return &GetCategoryHandler{repo: repo}
}

func (h *GetCategoryHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*GetCategoryQuery)
    
    // Buscar categoria
    category, err := h.repo.FindByID(ctx, query.CategoryID)
    if err != nil {
        return nil, err
    }
    
    // Converter para DTO
    return &dto.CategoryDTO{
        ID:          category.ID,
        Name:        category.Name,
        Slug:        category.Slug,
        Description: category.Description,
        ParentID:    category.ParentID,
        Active:      category.Active,
        SortOrder:   category.SortOrder,
        CreatedAt:   category.CreatedAt,
        UpdatedAt:   category.UpdatedAt,
    }, nil
}
```

**`application/queries/list_categories_query.go`**

```go
package queries

type ListCategoriesQuery struct {
    Page       int
    PageSize   int
    ParentID   *string  // Filtro por categoria pai
    ActiveOnly bool     // Apenas ativas
}

func (q *ListCategoriesQuery) QueryName() string {
    return "ListCategoriesQuery"
}
```

**`application/queries/list_categories_handler.go`**

```go
package queries

import (
    "context"
    "meuApp/internal/modules/category/dto"
    "meuApp/internal/modules/category/repository"
    "meuApp/pkg/contracts"
)

type ListCategoriesHandler struct {
    repo repository.CategoryRepository
}

func NewListCategoriesHandler(repo repository.CategoryRepository) *ListCategoriesHandler {
    return &ListCategoriesHandler{repo: repo}
}

func (h *ListCategoriesHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*ListCategoriesQuery)
    
    // Buscar categorias
    categories, total, err := h.repo.FindAll(
        ctx,
        query.Page,
        query.PageSize,
        query.ParentID,
        query.ActiveOnly,
    )
    if err != nil {
        return nil, err
    }
    
    // Converter para DTOs
    dtos := make([]*dto.CategoryDTO, len(categories))
    for i, cat := range categories {
        dtos[i] = &dto.CategoryDTO{
            ID:          cat.ID,
            Name:        cat.Name,
            Slug:        cat.Slug,
            Description: cat.Description,
            ParentID:    cat.ParentID,
            Active:      cat.Active,
            SortOrder:   cat.SortOrder,
            CreatedAt:   cat.CreatedAt,
            UpdatedAt:   cat.UpdatedAt,
        }
    }
    
    return &dto.PaginatedCategoriesDTO{
        Categories: dtos,
        Total:      total,
        Page:       query.Page,
        PageSize:   query.PageSize,
    }, nil
}
```

---

## 💾 Passo 4: Repository

### 4.1 Interface

**`repository/category_repository.go`**

```go
package repository

import (
    "context"
    "meuApp/internal/modules/category/domain/entities"
)

type CategoryRepository interface {
    // CRUD
    Save(ctx context.Context, category *entities.Category) error
    FindByID(ctx context.Context, id string) (*entities.Category, error)
    FindBySlug(ctx context.Context, slug string) (*entities.Category, error)
    FindAll(ctx context.Context, page, pageSize int, parentID *string, activeOnly bool) ([]*entities.Category, int64, error)
    Update(ctx context.Context, category *entities.Category) error
    Delete(ctx context.Context, id string) error
    
    // Queries específicas
    ExistsBySlug(ctx context.Context, slug string) (bool, error)
    FindByParentID(ctx context.Context, parentID string) ([]*entities.Category, error)
    CountByParentID(ctx context.Context, parentID string) (int64, error)
}
```

### 4.2 Implementação MySQL

**`repository/mysql_category_repository.go`**

```go
package repository

import (
    "context"
    "meuApp/internal/modules/category/domain/entities"
    "gorm.io/gorm"
)

type MySQLCategoryRepository struct {
    db *gorm.DB
}

func NewMySQLCategoryRepository(db *gorm.DB) *MySQLCategoryRepository {
    return &MySQLCategoryRepository{db: db}
}

// CategoryModel para banco de dados
type CategoryModel struct {
    ID          string  `gorm:"primaryKey"`
    Name        string  `gorm:"index"`
    Slug        string  `gorm:"uniqueIndex"`
    Description string  `gorm:"type:text"`
    ParentID    *string `gorm:"index"`
    Active      bool    `gorm:"index"`
    SortOrder   int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (CategoryModel) TableName() string {
    return "categories"
}

// Save persiste nova categoria
func (r *MySQLCategoryRepository) Save(ctx context.Context, category *entities.Category) error {
    model := r.toModel(category)
    return r.db.WithContext(ctx).Create(model).Error
}

// FindByID busca categoria por ID
func (r *MySQLCategoryRepository) FindByID(ctx context.Context, id string) (*entities.Category, error) {
    var model CategoryModel
    
    err := r.db.WithContext(ctx).
        Where("id = ?", id).
        First(&model).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("category not found")
    }
    if err != nil {
        return nil, err
    }
    
    return r.toEntity(&model), nil
}

// FindBySlug busca categoria por slug
func (r *MySQLCategoryRepository) FindBySlug(ctx context.Context, slug string) (*entities.Category, error) {
    var model CategoryModel
    
    err := r.db.WithContext(ctx).
        Where("slug = ?", slug).
        First(&model).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("category not found")
    }
    if err != nil {
        return nil, err
    }
    
    return r.toEntity(&model), nil
}

// FindAll lista categorias com filtros
func (r *MySQLCategoryRepository) FindAll(
    ctx context.Context,
    page, pageSize int,
    parentID *string,
    activeOnly bool,
) ([]*entities.Category, int64, error) {
    var models []CategoryModel
    var total int64
    
    query := r.db.WithContext(ctx).Model(&CategoryModel{})
    
    // Filtros
    if parentID != nil {
        query = query.Where("parent_id = ?", *parentID)
    }
    if activeOnly {
        query = query.Where("active = ?", true)
    }
    
    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Query paginada
    offset := (page - 1) * pageSize
    err := query.
        Offset(offset).
        Limit(pageSize).
        Order("sort_order ASC, name ASC").
        Find(&models).Error
    
    if err != nil {
        return nil, 0, err
    }
    
    // Converter para entities
    categories := make([]*entities.Category, len(models))
    for i, model := range models {
        categories[i] = r.toEntity(&model)
    }
    
    return categories, total, nil
}

// Update atualiza categoria
func (r *MySQLCategoryRepository) Update(ctx context.Context, category *entities.Category) error {
    model := r.toModel(category)
    return r.db.WithContext(ctx).Save(model).Error
}

// Delete remove categoria
func (r *MySQLCategoryRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&CategoryModel{}).Error
}

// ExistsBySlug verifica se slug já existe
func (r *MySQLCategoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&CategoryModel{}).
        Where("slug = ?", slug).
        Count(&count).Error
    
    return count > 0, err
}

// FindByParentID busca categorias filhas
func (r *MySQLCategoryRepository) FindByParentID(ctx context.Context, parentID string) ([]*entities.Category, error) {
    var models []CategoryModel
    
    err := r.db.WithContext(ctx).
        Where("parent_id = ?", parentID).
        Order("sort_order ASC").
        Find(&models).Error
    
    if err != nil {
        return nil, err
    }
    
    categories := make([]*entities.Category, len(models))
    for i, model := range models {
        categories[i] = r.toEntity(&model)
    }
    
    return categories, nil
}

// CountByParentID conta categorias filhas
func (r *MySQLCategoryRepository) CountByParentID(ctx context.Context, parentID string) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&CategoryModel{}).
        Where("parent_id = ?", parentID).
        Count(&count).Error
    
    return count, err
}

// Mappers
func (r *MySQLCategoryRepository) toEntity(model *CategoryModel) *entities.Category {
    return &entities.Category{
        ID:          model.ID,
        Name:        model.Name,
        Slug:        model.Slug,
        Description: model.Description,
        ParentID:    model.ParentID,
        Active:      model.Active,
        SortOrder:   model.SortOrder,
        CreatedAt:   model.CreatedAt,
        UpdatedAt:   model.UpdatedAt,
    }
}

func (r *MySQLCategoryRepository) toModel(category *entities.Category) *CategoryModel {
    return &CategoryModel{
        ID:          category.ID,
        Name:        category.Name,
        Slug:        category.Slug,
        Description: category.Description,
        ParentID:    category.ParentID,
        Active:      category.Active,
        SortOrder:   category.SortOrder,
        CreatedAt:   category.CreatedAt,
        UpdatedAt:   category.UpdatedAt,
    }
}
```

---

## 🔌 Passo 5: Adapters

### 5.1 HTTP Adapter

**`adapters/http/category_handler.go`**

```go
package http

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "meuApp/internal/modules/category/application/commands"
    "meuApp/internal/modules/category/application/queries"
    "meuApp/pkg/contracts"
)

type CategoryHTTPHandler struct {
    commandBus contracts.CommandBus
    queryBus   contracts.QueryBus
}

func NewCategoryHTTPHandler(
    commandBus contracts.CommandBus,
    queryBus contracts.QueryBus,
) *CategoryHTTPHandler {
    return &CategoryHTTPHandler{
        commandBus: commandBus,
        queryBus:   queryBus,
    }
}

// CreateCategory cria nova categoria
func (h *CategoryHTTPHandler) CreateCategory(c *gin.Context) {
    var req CreateCategoryRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    cmd := &commands.CreateCategoryCommand{
        Name:        req.Name,
        Slug:        req.Slug,
        Description: req.Description,
        ParentID:    req.ParentID,
        CreatedBy:   getUserIDFromContext(c),
    }
    
    categoryID, err := h.commandBus.Execute(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "id":      categoryID,
        "message": "Category created successfully",
    })
}

// GetCategory retorna categoria por ID
func (h *CategoryHTTPHandler) GetCategory(c *gin.Context) {
    categoryID := c.Param("id")
    
    query := &queries.GetCategoryQuery{
        CategoryID: categoryID,
    }
    
    result, err := h.queryBus.Execute(c.Request.Context(), query)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
        return
    }
    
    c.JSON(http.StatusOK, result)
}

// ListCategories lista categorias com paginação
func (h *CategoryHTTPHandler) ListCategories(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    activeOnly := c.Query("active_only") == "true"
    parentID := c.Query("parent_id")
    
    query := &queries.ListCategoriesQuery{
        Page:       page,
        PageSize:   pageSize,
        ActiveOnly: activeOnly,
    }
    
    if parentID != "" {
        query.ParentID = &parentID
    }
    
    result, err := h.queryBus.Execute(c.Request.Context(), query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, result)
}

func getUserIDFromContext(c *gin.Context) string {
    userID, exists := c.Get("user_id")
    if !exists {
        return "system"
    }
    return userID.(string)
}
```

**`adapters/http/dto.go`**

```go
package http

type CreateCategoryRequest struct {
    Name        string  `json:"name" binding:"required"`
    Slug        string  `json:"slug" binding:"required"`
    Description string  `json:"description"`
    ParentID    *string `json:"parent_id"`
}

type UpdateCategoryRequest struct {
    Name        string `json:"name" binding:"required"`
    Slug        string `json:"slug" binding:"required"`
    Description string `json:"description"`
}
```

**`adapters/http/routes.go`**

```go
package http

import "github.com/gin-gonic/gin"

func (h *CategoryHTTPHandler) RegisterRoutes(router *gin.Engine) {
    categories := router.Group("/api/v1/categories")
    {
        categories.POST("", h.CreateCategory)
        categories.GET("/:id", h.GetCategory)
        categories.GET("", h.ListCategories)
        categories.PUT("/:id", h.UpdateCategory)
        categories.DELETE("/:id", h.DeleteCategory)
    }
}
```

---

## 🎯 Passo 6: Registro do Módulo

**`category_module.go`**

```go
package category

import (
    "meuApp/internal/modules/category/adapters/http"
    "meuApp/internal/modules/category/application/commands"
    "meuApp/internal/modules/category/application/queries"
    "meuApp/internal/modules/category/repository"
    "meuApp/pkg/contracts"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type CategoryModule struct {
    db         *gorm.DB
    commandBus contracts.CommandBus
    queryBus   contracts.QueryBus
    eventBus   contracts.EventBus
    
    // Repository
    categoryRepo repository.CategoryRepository
    
    // Handlers
    httpHandler *http.CategoryHTTPHandler
}

func NewCategoryModule(
    db *gorm.DB,
    commandBus contracts.CommandBus,
    queryBus contracts.QueryBus,
    eventBus contracts.EventBus,
) *CategoryModule {
    return &CategoryModule{
        db:         db,
        commandBus: commandBus,
        queryBus:   queryBus,
        eventBus:   eventBus,
    }
}

func (m *CategoryModule) Name() string {
    return "CategoryModule"
}

func (m *CategoryModule) Initialize() error {
    // 1. Criar repository
    m.categoryRepo = repository.NewMySQLCategoryRepository(m.db)
    
    // 2. Registrar command handlers
    m.commandBus.Register(
        "CreateCategoryCommand",
        commands.NewCreateCategoryHandler(m.categoryRepo, m.eventBus),
    )
    
    m.commandBus.Register(
        "UpdateCategoryCommand",
        commands.NewUpdateCategoryHandler(m.categoryRepo, m.eventBus),
    )
    
    m.commandBus.Register(
        "DeleteCategoryCommand",
        commands.NewDeleteCategoryHandler(m.categoryRepo, m.eventBus),
    )
    
    // 3. Registrar query handlers
    m.queryBus.Register(
        "GetCategoryQuery",
        queries.NewGetCategoryHandler(m.categoryRepo),
    )
    
    m.queryBus.Register(
        "ListCategoriesQuery",
        queries.NewListCategoriesHandler(m.categoryRepo),
    )
    
    // 4. Criar HTTP handler
    m.httpHandler = http.NewCategoryHTTPHandler(m.commandBus, m.queryBus)
    
    return nil
}

func (m *CategoryModule) RegisterHTTPRoutes(router *gin.Engine) error {
    m.httpHandler.RegisterRoutes(router)
    return nil
}

func (m *CategoryModule) RegisterEventSubscribers(eventBus contracts.EventBus) error {
    // Registrar subscribers se necessário
    return nil
}
```

### Registrar no Bootstrap

**`internal/bootstrap/bootstrap_registry.go`**

```go
package bootstrap

import (
    "meuApp/internal/modules/category"
    // ... outros imports
)

func RegisterModules(registry *container.ModuleRegistry) error {
    // ... outros módulos
    
    // Category Module
    categoryModule := category.NewCategoryModule(
        registry.GetDB(),
        registry.GetCommandBus(),
        registry.GetQueryBus(),
        registry.GetEventBus(),
    )
    registry.Register(categoryModule)
    
    return nil
}
```

---

## 🧪 Passo 7: Testes

### 7.1 Teste Unitário do Handler

**`application/commands/create_category_handler_test.go`**

```go
package commands_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "meuApp/internal/modules/category/application/commands"
)

type MockCategoryRepository struct {
    mock.Mock
}

func (m *MockCategoryRepository) Save(ctx context.Context, category interface{}) error {
    args := m.Called(ctx, category)
    return args.Error(0)
}

func (m *MockCategoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
    args := m.Called(ctx, slug)
    return args.Bool(0), args.Error(1)
}

func TestCreateCategoryHandler_Handle(t *testing.T) {
    // Arrange
    mockRepo := new(MockCategoryRepository)
    mockEventBus := new(MockEventBus)
    
    handler := commands.NewCreateCategoryHandler(mockRepo, mockEventBus)
    
    cmd := &commands.CreateCategoryCommand{
        Name:        "Electronics",
        Slug:        "electronics",
        Description: "Electronic products",
        CreatedBy:   "user-123",
    }
    
    mockRepo.On("ExistsBySlug", mock.Anything, "electronics").Return(false, nil)
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    mockEventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
    
    // Act
    result, err := handler.Handle(context.Background(), cmd)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
    mockRepo.AssertExpectations(t)
    mockEventBus.AssertExpectations(t)
}
```

### 7.2 Teste de Integração

**`repository/mysql_category_repository_test.go`**

```go
package repository_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "meuApp/internal/modules/category/domain/entities"
    "meuApp/internal/modules/category/repository"
)

func TestMySQLCategoryRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup
    db := setupTestDB(t)
    defer db.Close()
    
    repo := repository.NewMySQLCategoryRepository(db)
    ctx := context.Background()
    
    t.Run("Save and FindByID", func(t *testing.T) {
        // Create
        category, _ := entities.NewCategory(
            "Electronics",
            "electronics",
            "Electronic products",
            nil,
        )
        
        err := repo.Save(ctx, category)
        assert.NoError(t, err)
        
        // Find
        found, err := repo.FindByID(ctx, category.ID)
        assert.NoError(t, err)
        assert.Equal(t, "Electronics", found.Name)
        assert.Equal(t, "electronics", found.Slug)
    })
}
```

---

## ✅ Checklist de Criação de Módulo

```markdown
## Domain Layer
- [ ] Criar entidade principal
- [ ] Adicionar métodos de negócio
- [ ] Criar value objects (se necessário)
- [ ] Definir eventos de domínio
- [ ] Validações de domínio

## Application Layer
- [ ] Criar commands (write operations)
- [ ] Criar command handlers
- [ ] Criar queries (read operations)
- [ ] Criar query handlers
- [ ] Criar DTOs

## Repository
- [ ] Definir interface do repository
- [ ] Implementar repository (MySQL/PostgreSQL)
- [ ] Criar mappers (Entity ↔ Model)
- [ ] Implementar queries específicas

## Adapters
- [ ] Criar HTTP handlers
- [ ] Criar request/response DTOs
- [ ] Registrar rotas
- [ ] Criar gRPC handlers (se necessário)
- [ ] Criar event subscribers (se necessário)

## Module Registration
- [ ] Criar arquivo do módulo
- [ ] Implementar Initialize()
- [ ] Registrar commands
- [ ] Registrar queries
- [ ] Registrar rotas HTTP
- [ ] Registrar no bootstrap

## Tests
- [ ] Testes unitários de handlers
- [ ] Testes de repository
- [ ] Testes de integração HTTP
- [ ] Testes de eventos

## Database
- [ ] Criar migration
- [ ] Criar seeders (se necessário)
- [ ] Adicionar índices

## Documentation
- [ ] Documentar API endpoints
- [ ] Atualizar README
- [ ] Adicionar exemplos
```

---

## 📚 Próximos Passos

- **[Commands](13-implementing-commands.md)** - Detalhes de write operations
- **[Queries](14-implementing-queries.md)** - Detalhes de read operations
- **[Testes](21-testing-strategy.md)** - Estratégia completa de testes

---

**[⬅️ Repositórios](11-repositories.md)** | **[Índice](README.md)** | **[Commands ➡️](13-implementing-commands.md)**
