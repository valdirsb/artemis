package ports

import (
	"context"
	"meuApp/internal/modules/product/domain"
)

// ===== PRIMARY PORTS (Application Services) =====

// ProductService define as operações de negócio do módulo de produto
type ProductService interface {
	CreateProduct(ctx context.Context, name, description, categoryID string, price float64, stock int) (*domain.Product, error)
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id string, name, description, categoryID *string, price *float64, stock *int) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	ListProducts(ctx context.Context, filters ProductFilters) ([]*domain.Product, error)
	ListProductsPaginated(ctx context.Context, filters ProductFilters) (*PaginatedResult, error)
	UpdateStock(ctx context.Context, id string, quantity int) error
}

// ===== SECONDARY PORTS (Adapters/Dependencies) =====

// ProductRepository define a interface para persistência de produtos
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ProductFilters) ([]*domain.Product, error)
	ListPaginated(ctx context.Context, filters ProductFilters) (*PaginatedResult, error)
}

// ProductFilters define os filtros para listagem de produtos
type ProductFilters struct {
	CategoryID *string
	MinPrice   *float64
	MaxPrice   *float64
	InStock    *bool
	Page       int
	PageSize   int
}

// PaginatedResult representa um resultado paginado
type PaginatedResult struct {
	Items      []*domain.Product
	TotalItems int64
	Page       int
	PageSize   int
	TotalPages int
}
