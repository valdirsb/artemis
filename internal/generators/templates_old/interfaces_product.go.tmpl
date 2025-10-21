package contracts

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Interfaces para comunicação entre módulos (Ports)

// ProductService define operações de negócio relacionadas a produtos
type ProductService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	UpdateProduct(ctx context.Context, id string, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id string) error
	GetProducts(ctx context.Context, filters ProductFilters) ([]*Product, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
}

// Repository Interfaces (Adapters)

// ProductRepository define a interface para persistência de produtos
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ProductFilters) ([]*Product, error)
}

// Domain Models

// Product representa o modelo de domínio do produto
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request/Response DTOs

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	CategoryID  string  `json:"category_id" validate:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
	CategoryID  *string  `json:"category_id,omitempty"`
}

type ProductFilters struct {
	CategoryID *string
	MinPrice   *float64
	MaxPrice   *float64
	Name       *string
	Limit      int
	Offset     int
}

// Events para comunicação entre módulos

type ProductCreatedEvent struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	CategoryID string  `json:"category_id"`
	Price      float64 `json:"price"`
}

type ProductStockUpdatedEvent struct {
	ProductID string `json:"product_id"`
	NewStock  int    `json:"new_stock"`
}

// Handler Interfaces

type ProductHandler interface {
	CreateProduct(ctx *gin.Context)
	GetProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	GetProducts(ctx *gin.Context)
	UpdateStock(ctx *gin.Context)
}
