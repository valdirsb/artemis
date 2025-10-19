# 📦 Como Criar um Novo Módulo

> **Guia passo-a-passo** para desenvolvedores criarem novos módulos no framework Artemis com auto-registro completo

## 📋 Índice

- [📦 Como Criar um Novo Módulo](#-como-criar-um-novo-módulo)
  - [📋 Índice](#-índice)
  - [🎯 Pré-requisitos](#-pré-requisitos)
  - [🏗️ Visão Geral](#️-visão-geral)
    - [O que vamos construir:](#o-que-vamos-construir)
    - [Funcionalidades:](#funcionalidades)
  - [📁 Passo 1: Estrutura de Diretórios](#-passo-1-estrutura-de-diretórios)
  - [🎯 Passo 2: Domain Layer](#-passo-2-domain-layer)
    - [2.1 Criar a Entidade de Domínio](#21-criar-a-entidade-de-domínio)
    - [2.2 Criar Erros do Domínio](#22-criar-erros-do-domínio)
  - [🔌 Passo 3: Ports (Interfaces)](#-passo-3-ports-interfaces)
  - [📄 Passo 4: DTOs](#-passo-4-dtos)
    - [4.1 Request DTOs](#41-request-dtos)
    - [4.2 Response DTOs](#42-response-dtos)
    - [4.3 Mappers](#43-mappers)
  - [💾 Passo 5: Repository](#-passo-5-repository)
    - [5.1 Model (GORM)](#51-model-gorm)
    - [5.2 Repository Implementation](#52-repository-implementation)
  - [✍️ Passo 6: Commands](#️-passo-6-commands)
    - [6.1 Create Category Command](#61-create-category-command)
    - [6.2 Update Category Command](#62-update-category-command)
    - [6.3 Delete Category Command](#63-delete-category-command)
  - [🔍 Passo 7: Queries](#-passo-7-queries)
    - [7.1 Get Category Query](#71-get-category-query)
    - [7.2 List Categories Query](#72-list-categories-query)
  - [🎯 Passo 8: Application Service](#-passo-8-application-service)
  - [🌐 Passo 9: HTTP Handler](#-passo-9-http-handler)
  - [⚡ Passo 10: gRPC Service](#-passo-10-grpc-service)
    - [10.1 Definir Proto](#101-definir-proto)
    - [10.2 Gerar código gRPC](#102-gerar-código-grpc)
    - [10.3 Implementar Service](#103-implementar-service)
  - [🔌 Passo 11: Auto-Registro](#-passo-11-auto-registro)
  - [🚀 Passo 12: Integração no Bootstrap](#-passo-12-integração-no-bootstrap)
  - [✅ Checklist](#-checklist)
  - [🎉 Conclusão](#-conclusão)
    - [Próximos Passos:](#próximos-passos)
    - [Recursos Adicionais:](#recursos-adicionais)

---

## 🎯 Pré-requisitos

Antes de começar, certifique-se de ter:

- ✅ Go 1.24+ instalado
- ✅ Entendimento básico de Clean Architecture
- ✅ Conhecimento de CQRS (Command Query Responsibility Segregation)
- ✅ Familiaridade com Dependency Injection

**Leitura recomendada:**
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Entenda a arquitetura
- [EVENTS_GUIDE.md](../EVENTS_GUIDE.md) - Sistema de eventos

---

## 🏗️ Visão Geral

Vamos criar um módulo **Category** (Categoria de Produtos) como exemplo.

### O que vamos construir:

```
internal/modules/
├── category_module.go          # ← Auto-registro
└── category/
    ├── domain/
    │   └── category.go         # ← Entidade
    ├── ports/
    │   └── ports.go            # ← Interfaces
    ├── dto/
    │   ├── requests.go         # ← Request DTOs
    │   ├── responses.go        # ← Response DTOs
    │   └── mapper.go           # ← Conversões
    ├── application/
    │   ├── commands/
    │   │   ├── create_category.go
    │   │   ├── update_category.go
    │   │   └── delete_category.go
    │   ├── queries/
    │   │   ├── get_category.go
    │   │   └── list_categories.go
    │   └── services/
    │       └── category_application_service.go
    ├── adapters/
    │   ├── http/
    │   │   └── category_handler.go
    │   └── grpc/
    │       └── category_service.go
    ├── repository/
    │   ├── category_model.go
    │   └── category_repository.go
    └── errors.go
```

### Funcionalidades:
- ✅ CRUD completo (Create, Read, Update, Delete)
- ✅ HTTP REST API
- ✅ gRPC API
- ✅ Event Bus integration
- ✅ Auto-registro completo

---

## 📁 Passo 1: Estrutura de Diretórios

Crie a estrutura base:

```bash
cd internal/modules
mkdir -p category/{domain,ports,dto,application/{commands,queries,services},adapters/{http,grpc},repository}
touch category/errors.go
touch category_module.go
```

---

## 🎯 Passo 2: Domain Layer

### 2.1 Criar a Entidade de Domínio

**Arquivo:** `internal/modules/category/domain/category.go`

```go
package domain

import (
	"errors"
	"time"
)

// Category representa uma categoria de produtos no domínio
type Category struct {
	ID          string
	Name        string
	Description string
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCategory cria uma nova categoria com validações
func NewCategory(name, description string) (*Category, error) {
	category := &Category{
		Name:        name,
		Description: description,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	return category, nil
}

// Validate valida as regras de negócio da categoria
func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("category name cannot be empty")
	}

	if len(c.Name) < 3 {
		return errors.New("category name must have at least 3 characters")
	}

	if len(c.Name) > 100 {
		return errors.New("category name must have at most 100 characters")
	}

	return nil
}

// Deactivate desativa a categoria
func (c *Category) Deactivate() {
	c.Active = false
	c.UpdatedAt = time.Now()
}

// Activate ativa a categoria
func (c *Category) Activate() {
	c.Active = true
	c.UpdatedAt = time.Now()
}

// Update atualiza os dados da categoria
func (c *Category) Update(name, description string) error {
	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now()

	return c.Validate()
}
```

### 2.2 Criar Erros do Domínio

**Arquivo:** `internal/modules/category/errors.go`

```go
package category

import "meuApp/pkg/errors"

var (
	// Domain errors
	ErrCategoryNotFound      = errors.NewNotFoundError("CATEGORY_NOT_FOUND", "category not found")
	ErrCategoryAlreadyExists = errors.NewConflictError("CATEGORY_ALREADY_EXISTS", map[string]string{"field": "category"})
	ErrInvalidCategoryData   = errors.NewValidationError("INVALID_CATEGORY_DATA", map[string]string{"field": "data"})
	ErrCategoryNameRequired  = errors.NewValidationError("CATEGORY_NAME_REQUIRED", map[string]string{"field": "name"})
	ErrCategoryInactive      = errors.NewDomainError("CATEGORY_INACTIVE", nil)
)
```

---

## 🔌 Passo 3: Ports (Interfaces)

**Arquivo:** `internal/modules/category/ports/ports.go`

```go
package ports

import (
	"context"
	"meuApp/internal/modules/category/domain"
	"meuApp/internal/modules/category/dto"
)

// ========== PRIMARY PORTS (Use Cases - O que o módulo OFERECE) ==========

// CreateCategoryUseCase define o contrato para criar categoria
type CreateCategoryUseCase interface {
	Execute(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
}

// UpdateCategoryUseCase define o contrato para atualizar categoria
type UpdateCategoryUseCase interface {
	Execute(ctx context.Context, id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
}

// DeleteCategoryUseCase define o contrato para deletar categoria
type DeleteCategoryUseCase interface {
	Execute(ctx context.Context, id string) error
}

// GetCategoryQuery define o contrato para buscar categoria por ID
type GetCategoryQuery interface {
	Execute(ctx context.Context, id string) (*dto.CategoryResponse, error)
}

// ListCategoriesQuery define o contrato para listar categorias
type ListCategoriesQuery interface {
	Execute(ctx context.Context, activeOnly bool) ([]*dto.CategoryResponse, error)
}

// ========== SECONDARY PORTS (O que o módulo PRECISA) ==========

// CategoryRepository define o contrato para persistência
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	FindByName(ctx context.Context, name string) (*domain.Category, error)
	FindAll(ctx context.Context) ([]*domain.Category, error)
	FindAllActive(ctx context.Context) ([]*domain.Category, error)
}
```

---

## 📄 Passo 4: DTOs

### 4.1 Request DTOs

**Arquivo:** `internal/modules/category/dto/requests.go`

```go
package dto

// CreateCategoryRequest representa a requisição para criar categoria
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// UpdateCategoryRequest representa a requisição para atualizar categoria
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=500"`
	Active      *bool  `json:"active"` // Pointer para diferenciar null de false
}
```

### 4.2 Response DTOs

**Arquivo:** `internal/modules/category/dto/responses.go`

```go
package dto

import "time"

// CategoryResponse representa a resposta com dados da categoria
type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryListResponse representa a lista de categorias
type CategoryListResponse struct {
	Categories []*CategoryResponse `json:"categories"`
	Total      int                 `json:"total"`
}
```

### 4.3 Mappers

**Arquivo:** `internal/modules/category/dto/mapper.go`

```go
package dto

import "meuApp/internal/modules/category/domain"

// ToCategoryResponse converte domain.Category → CategoryResponse
func ToCategoryResponse(category *domain.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Active:      category.Active,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// ToCategoryListResponse converte []domain.Category → CategoryListResponse
func ToCategoryListResponse(categories []*domain.Category) *CategoryListResponse {
	responses := make([]*CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = ToCategoryResponse(category)
	}

	return &CategoryListResponse{
		Categories: responses,
		Total:      len(responses),
	}
}
```

---

## 💾 Passo 5: Repository

### 5.1 Model (GORM)

**Arquivo:** `internal/modules/category/repository/category_model.go`

```go
package repository

import (
	"meuApp/internal/modules/category/domain"
	"time"
)

// CategoryModel representa a categoria no banco de dados
type CategoryModel struct {
	ID          string    `gorm:"type:char(36);primaryKey"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string    `gorm:"type:varchar(500)"`
	Active      bool      `gorm:"default:true;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// TableName sobrescreve o nome da tabela
func (CategoryModel) TableName() string {
	return "categories"
}

// FromDomain converte domain.Category → CategoryModel
func FromDomain(category *domain.Category) *CategoryModel {
	return &CategoryModel{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Active:      category.Active,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// ToDomain converte CategoryModel → domain.Category
func (m *CategoryModel) ToDomain() *domain.Category {
	return &domain.Category{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
```

### 5.2 Repository Implementation

**Arquivo:** `internal/modules/category/repository/category_repository.go`

```go
package repository

import (
	"context"
	"meuApp/internal/modules/category"
	"meuApp/internal/modules/category/domain"
	"meuApp/internal/modules/category/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository cria uma nova instância do repositório
func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, cat *domain.Category) error {
	cat.ID = uuid.New().String()
	model := FromDomain(cat)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}

	*cat = *model.ToDomain()
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, cat *domain.Category) error {
	model := FromDomain(cat)

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return category.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&CategoryModel{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return category.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	var model CategoryModel

	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, category.ErrCategoryNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *categoryRepository) FindByName(ctx context.Context, name string) (*domain.Category, error) {
	var model CategoryModel

	if err := r.db.WithContext(ctx).First(&model, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, category.ErrCategoryNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]*domain.Category, error) {
	var models []CategoryModel

	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(models))
	for i, model := range models {
		categories[i] = model.ToDomain()
	}

	return categories, nil
}

func (r *categoryRepository) FindAllActive(ctx context.Context) ([]*domain.Category, error) {
	var models []CategoryModel

	if err := r.db.WithContext(ctx).Where("active = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(models))
	for i, model := range models {
		categories[i] = model.ToDomain()
	}

	return categories, nil
}
```

---

## ✍️ Passo 6: Commands

### 6.1 Create Category Command

**Arquivo:** `internal/modules/category/application/commands/create_category.go`

```go
package commands

import (
	"context"
	"fmt"

	"meuApp/internal/modules/category"
	"meuApp/internal/modules/category/domain"
	"meuApp/internal/modules/category/dto"
	"meuApp/internal/modules/category/ports"
	"meuApp/pkg/contracts"
)

type CreateCategoryHandler struct {
	repo     ports.CategoryRepository
	eventBus contracts.EventPublisher
}

func NewCreateCategoryHandler(
	repo ports.CategoryRepository,
	eventBus contracts.EventPublisher,
) *CreateCategoryHandler {
	return &CreateCategoryHandler{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (h *CreateCategoryHandler) Execute(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// 1. Validar se já existe categoria com mesmo nome
	existing, _ := h.repo.FindByName(ctx, req.Name)
	if existing != nil {
		return nil, category.ErrCategoryAlreadyExists
	}

	// 2. Criar entidade de domínio
	cat, err := domain.NewCategory(req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// 3. Persistir
	if err := h.repo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("failed to save category: %w", err)
	}

	// // 4. Publicar evento
	// event := events.CategoryCreatedEvent{
	//     CategoryID:  cat.ID,
	//     Name:        cat.Name,
	//     Description: cat.Description,
	//     CreatedAt:   cat.CreatedAt,
	// }
	// h.eventBus.Publish("category.created", event)

	// 5. Retornar resposta
	return dto.ToCategoryResponse(cat), nil
}
```

### 6.2 Update Category Command

**Arquivo:** `internal/modules/category/application/commands/update_category.go`

```go
package commands

import (
	"context"
	"fmt"

	"meuApp/internal/modules/category/dto"
	"meuApp/internal/modules/category/ports"
	"meuApp/pkg/contracts"
)

type UpdateCategoryHandler struct {
	repo   ports.CategoryRepository
	logger contracts.Logger
}

func NewUpdateCategoryHandler(
	repo ports.CategoryRepository,
	logger contracts.Logger,
) *UpdateCategoryHandler {
	return &UpdateCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

func (h *UpdateCategoryHandler) Execute(ctx context.Context, id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	// 1. Buscar categoria existente
	cat, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Atualizar dados
	if err := cat.Update(req.Name, req.Description); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// 3. Atualizar status se fornecido
	if req.Active != nil {
		if *req.Active {
			cat.Activate()
		} else {
			cat.Deactivate()
		}
	}

	// 4. Persistir
	if err := h.repo.Update(ctx, cat); err != nil {
		return nil, fmt.Errorf("failed to save category: %w", err)
	}

	h.logger.Info(fmt.Sprintf("Category updated successfully: %s", id))

	// 5. Retornar resposta
	return dto.ToCategoryResponse(cat), nil
}
```

### 6.3 Delete Category Command

**Arquivo:** `internal/modules/category/application/commands/delete_category.go`

```go
package commands

import (
	"context"
	"fmt"

	"meuApp/internal/modules/category/ports"
	"meuApp/pkg/contracts"
)

type DeleteCategoryHandler struct {
	repo   ports.CategoryRepository
	logger contracts.Logger
}

func NewDeleteCategoryHandler(
	repo ports.CategoryRepository,
	logger contracts.Logger,
) *DeleteCategoryHandler {
	return &DeleteCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

func (h *DeleteCategoryHandler) Execute(ctx context.Context, id string) error {
	// 1. Verificar se existe
	_, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Deletar
	if err := h.repo.Delete(ctx, id); err != nil {
		return err
	}

	h.logger.Info(fmt.Sprintf("Category deleted successfully: %s", id))
	return nil
}
```

---

## 🔍 Passo 7: Queries

### 7.1 Get Category Query

**Arquivo:** `internal/modules/category/application/queries/get_category.go`

```go
package queries

import (
	"context"

	"meuApp/internal/modules/category/dto"
	"meuApp/internal/modules/category/ports"
)

type GetCategoryHandler struct {
	repo ports.CategoryRepository
}

func NewGetCategoryHandler(repo ports.CategoryRepository) *GetCategoryHandler {
	return &GetCategoryHandler{repo: repo}
}

func (h *GetCategoryHandler) Execute(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	category, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.ToCategoryResponse(category), nil
}
```

### 7.2 List Categories Query

**Arquivo:** `internal/modules/category/application/queries/list_categories.go`

```go
package queries

import (
	"context"

	"meuApp/internal/modules/category/domain"
	"meuApp/internal/modules/category/dto"
	"meuApp/internal/modules/category/ports"
)

type ListCategoriesHandler struct {
	repo ports.CategoryRepository
}

func NewListCategoriesHandler(repo ports.CategoryRepository) *ListCategoriesHandler {
	return &ListCategoriesHandler{repo: repo}
}

func (h *ListCategoriesHandler) Execute(ctx context.Context, activeOnly bool) ([]*dto.CategoryResponse, error) {
	var categories []*domain.Category
	var err error

	if activeOnly {
		categories, err = h.repo.FindAllActive(ctx)
	} else {
		categories, err = h.repo.FindAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	responses := make([]*dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dto.ToCategoryResponse(cat)
	}

	return responses, nil
}
```

---

## 🎯 Passo 8: Application Service

**Arquivo:** `internal/modules/category/application/services/category_application_service.go`

```go
package services

import (
	"context"

	"meuApp/internal/modules/category/application/commands"
	"meuApp/internal/modules/category/application/queries"
	"meuApp/internal/modules/category/dto"
)

// CategoryApplicationService orquestra todos os use cases
type CategoryApplicationService struct {
	createCmd *commands.CreateCategoryHandler
	updateCmd *commands.UpdateCategoryHandler
	deleteCmd *commands.DeleteCategoryHandler
	getQuery  *queries.GetCategoryHandler
	listQuery *queries.ListCategoriesHandler
}

func NewCategoryApplicationService(
	createCmd *commands.CreateCategoryHandler,
	updateCmd *commands.UpdateCategoryHandler,
	deleteCmd *commands.DeleteCategoryHandler,
	getQuery *queries.GetCategoryHandler,
	listQuery *queries.ListCategoriesHandler,
) *CategoryApplicationService {
	return &CategoryApplicationService{
		createCmd: createCmd,
		updateCmd: updateCmd,
		deleteCmd: deleteCmd,
		getQuery:  getQuery,
		listQuery: listQuery,
	}
}

// Commands
func (s *CategoryApplicationService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	return s.createCmd.Execute(ctx, req)
}

func (s *CategoryApplicationService) UpdateCategory(ctx context.Context, id string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	return s.updateCmd.Execute(ctx, id, req)
}

func (s *CategoryApplicationService) DeleteCategory(ctx context.Context, id string) error {
	return s.deleteCmd.Execute(ctx, id)
}

// Queries
func (s *CategoryApplicationService) GetCategory(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	return s.getQuery.Execute(ctx, id)
}

func (s *CategoryApplicationService) ListCategories(ctx context.Context, activeOnly bool) ([]*dto.CategoryResponse, error) {
	return s.listQuery.Execute(ctx, activeOnly)
}
```

---

## 🌐 Passo 9: HTTP Handler

**Arquivo:** `internal/modules/category/adapters/http/category_handler.go`

```go
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meuApp/internal/modules/category/application/services"
	"meuApp/internal/modules/category/dto"
	"meuApp/pkg/errors"
)

type CategoryHandler struct {
	appService *services.CategoryApplicationService
}

func NewCategoryHandler(appService *services.CategoryApplicationService) *CategoryHandler {
	return &CategoryHandler{appService: appService}
}

// RegisterRoutes registra as rotas HTTP
func (h *CategoryHandler) RegisterRoutes(router gin.IRouter) {

	categories := router.Group("/categories")
	{
		categories.POST("", h.CreateCategory)
		categories.GET("/:id", h.GetCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
		categories.GET("", h.ListCategories)
	}

}

// CreateCategory godoc
// @Summary Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.appService.CreateCategory(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatusCode(), gin.H{"error": appErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetCategory godoc
// @Summary Get category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")

	response, err := h.appService.GetCategory(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatusCode(), gin.H{"error": appErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCategory godoc
// @Summary Update category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body dto.UpdateCategoryRequest true "Category data"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.appService.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatusCode(), gin.H{"error": appErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteCategory godoc
// @Summary Delete category
// @Tags categories
// @Param id path string true "Category ID"
// @Success 204
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.appService.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.HTTPStatusCode(), gin.H{"error": appErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListCategories godoc
// @Summary List all categories
// @Tags categories
// @Produce json
// @Param active_only query bool false "Show only active categories"
// @Success 200 {array} dto.CategoryResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"

	response, err := h.appService.ListCategories(c.Request.Context(), activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": response, "total": len(response)})
}
```

---

## ⚡ Passo 10: gRPC Service

### 10.1 Definir Proto

**Arquivo:** `proto/category.proto`

```protobuf
syntax = "proto3";

package category;
option go_package = "meuApp/pkg/proto";

import "google/protobuf/timestamp.proto";

service CategoryService {
  rpc CreateCategory(CreateCategoryRequest) returns (CategoryResponse);
  rpc GetCategory(GetCategoryRequest) returns (CategoryResponse);
  rpc UpdateCategory(UpdateCategoryRequest) returns (CategoryResponse);
  rpc DeleteCategory(DeleteCategoryRequest) returns (DeleteCategoryResponse);
  rpc ListCategories(ListCategoriesRequest) returns (ListCategoriesResponse);
}

message CreateCategoryRequest {
  string name = 1;
  string description = 2;
}

message UpdateCategoryRequest {
  string id = 1;
  string name = 2;
  string description = 3;
  optional bool active = 4;
}

message GetCategoryRequest {
  string id = 1;
}

message DeleteCategoryRequest {
  string id = 1;
}

message DeleteCategoryResponse {
  bool success = 1;
}

message ListCategoriesRequest {
  bool active_only = 1;
}

message CategoryResponse {
  string id = 1;
  string name = 2;
  string description = 3;
  bool active = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}

message ListCategoriesResponse {
  repeated CategoryResponse categories = 1;
  int32 total = 2;
}
```

### 10.2 Gerar código gRPC

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/category.proto
```

### 10.3 Implementar Service

**Arquivo:** `internal/modules/category/adapters/grpc/category_service.go`

```go
package grpc

import (
	"context"

	"meuApp/internal/modules/category/application/services"
	"meuApp/internal/modules/category/dto"
	pb "meuApp/pkg/proto"

	"google.golang.org/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryGRPCHandler struct {
	pb.UnimplementedCategoryServiceServer
	appService *services.CategoryApplicationService
}

func NewCategoryGRPCService(appService *services.CategoryApplicationService) *CategoryGRPCHandler {
	return &CategoryGRPCHandler{appService: appService}
}

// RegisterWithServer registers the service with the gRPC server
func (s *CategoryGRPCHandler) RegisterWithServer(server *grpc.Server) {
	pb.RegisterCategoryServiceServer(server, s)
}

func (s *CategoryGRPCHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	dtoReq := dto.CreateCategoryRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	response, err := s.appService.CreateCategory(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &pb.CategoryResponse{
		Id:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		Active:      response.Active,
		CreatedAt:   timestamppb.New(response.CreatedAt),
		UpdatedAt:   timestamppb.New(response.UpdatedAt),
	}, nil
}

func (s *CategoryGRPCHandler) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
	response, err := s.appService.GetCategory(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.CategoryResponse{
		Id:          response.ID,
		Name:        response.Name,
		Description: response.Description,
		Active:      response.Active,
		CreatedAt:   timestamppb.New(response.CreatedAt),
		UpdatedAt:   timestamppb.New(response.UpdatedAt),
	}, nil
}

// ... implementar outros métodos

```

---

## 🔌 Passo 11: Auto-Registro

**Arquivo:** `internal/modules/category_module.go`

```go
package modules

import (
	"fmt"

	"gorm.io/gorm"

	"meuApp/internal/modules/category/adapters/grpc"
	categoryHTTP "meuApp/internal/modules/category/adapters/http"
	"meuApp/internal/modules/category/application/commands"
	"meuApp/internal/modules/category/application/queries"
	"meuApp/internal/modules/category/application/services"
	"meuApp/internal/modules/category/repository"
	"meuApp/pkg/container"
	"meuApp/pkg/contracts"

	"github.com/gin-gonic/gin"
	grpclib "google.golang.org/grpc"
)

type CategoryModule struct {
	db       *gorm.DB
	eventBus contracts.EventPublisher
	logger   contracts.Logger
}

func NewCategoryModule(db *gorm.DB, eventBus contracts.EventPublisher, logger contracts.Logger) *CategoryModule {
	return &CategoryModule{
		db:       db,
		eventBus: eventBus,
		logger:   logger,
	}
}

func (m *CategoryModule) Name() string {
	return "category"
}

func (m *CategoryModule) Register(registry *container.ModuleRegistry) error {
	fmt.Printf("🔧 Registering module: %s\n", m.Name())

	// 1. Repository
	categoryRepo := repository.NewCategoryRepository(m.db)
	registry.RegisterRepository("category", categoryRepo)

	// 2. Command Handlers
	createCmd := commands.NewCreateCategoryHandler(categoryRepo, m.eventBus)
	updateCmd := commands.NewUpdateCategoryHandler(categoryRepo, m.logger)
	deleteCmd := commands.NewDeleteCategoryHandler(categoryRepo, m.logger)

	// 3. Query Handlers
	getQuery := queries.NewGetCategoryHandler(categoryRepo)
	listQuery := queries.NewListCategoriesHandler(categoryRepo)

	// 4. Application Service
	appService := services.NewCategoryApplicationService(
		createCmd,
		updateCmd,
		deleteCmd,
		getQuery,
		listQuery,
	)
	registry.RegisterApplicationService("category", appService)

	// 5. HTTP Handler
	httpHandler := categoryHTTP.NewCategoryHandler(appService)
	registry.RegisterHTTPHandler("category", &categoryHTTPHandlerAdapter{handler: httpHandler})

	// 6. gRPC Service
	grpcHandler := grpc.NewCategoryGRPCService(appService)
	registry.RegisterGRPCService(&categoryGRPCServiceAdapter{handler: grpcHandler})

	fmt.Println("✅ Module category registered successfully")
	return nil
}

// categoryHTTPHandlerAdapter adapta o HTTP handler para a interface do registry
type categoryHTTPHandlerAdapter struct {
	handler *categoryHTTP.CategoryHandler
}

func (a *categoryHTTPHandlerAdapter) RegisterRoutes(router *gin.RouterGroup) {

	a.handler.RegisterRoutes(router)
}

// categoryGRPCServiceAdapter adapta o gRPC handler para a interface do registry
type categoryGRPCServiceAdapter struct {
	handler *grpc.CategoryGRPCHandler
}

func (a *categoryGRPCServiceAdapter) RegisterService(server *grpclib.Server) {
	a.handler.RegisterWithServer(server)
}
```

---

## 🚀 Passo 12: Integração no Bootstrap

**Arquivo:** `internal/bootstrap/bootstrap_registry.go`

```go
// Adicionar no registerModules()
func registerModules(registry *container.ModuleRegistry, db *gorm.DB, eventBus contracts.EventPublisher, logger contracts.Logger) error {

    applicationModules := []frameworkInterfaces.Module{
		// ... módulos existentes
		modules.NewCategoryModule(db, eventBus, logger),
	}
    // ... 
}
```

---

## ✅ Checklist

Antes de finalizar, verifique se tudo está completo:

- [ ] **Domain Layer**
  - [ ] Entidade criada (`domain/category.go`)
  - [ ] Validações implementadas
  - [ ] Erros do domínio definidos (`errors.go`)

- [ ] **Ports**
  - [ ] Interfaces de use cases definidas
  - [ ] Interface do repositório definida

- [ ] **DTOs**
  - [ ] Request DTOs criados
  - [ ] Response DTOs criados
  - [ ] Mappers implementados

- [ ] **Repository**
  - [ ] Model GORM criado
  - [ ] Repository implementado
  - [ ] Todos os métodos testados

- [ ] **Application**
  - [ ] Commands criados (Create, Update, Delete)
  - [ ] Queries criados (Get, List)
  - [ ] Application Service criado

- [ ] **Adapters**
  - [ ] HTTP Handler implementado
  - [ ] Rotas registradas
  - [ ] gRPC Service implementado
  - [ ] Proto definido e gerado

- [ ] **Auto-Registro**
  - [ ] Module file criado
  - [ ] Interface Module implementada
  - [ ] Register() completo
  - [ ] Integrado no bootstrap

- [ ] **Testes**
  - [ ] Compilação OK (`go build`)
  - [ ] Servidor inicia sem erros
  - [ ] Endpoints HTTP funcionando
  - [ ] gRPC funcionando

---

## 🎉 Conclusão

Parabéns! Você criou um módulo completo seguindo todos os padrões do Artemis Framework.

### Próximos Passos:

1. **Testar o módulo:**
   ```bash
   # HTTP
   curl -X POST http://localhost:8080/api/v1/categories \
     -H "Content-Type: application/json" \
     -d '{"name":"Electronics","description":"Electronic products"}'
   
   # gRPC
   grpcurl -plaintext -d '{"name":"Electronics"}' \
     localhost:50051 category.CategoryService/CreateCategory
   ```

2. **Adicionar testes unitários**
3. **Adicionar eventos personalizados**
4. **Documentar no Swagger**

### Recursos Adicionais:

- [ARCHITECTURE.md](./ARCHITECTURE.md)
- [EVENTS_GUIDE.md](../EVENTS_GUIDE.md)
- [API.md](./API.md)

---

**Happy Coding! 🚀**
