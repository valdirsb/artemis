package repository

import (
	"context"
	"fmt"

	"meuApp/internal/modules/product/domain"
	"meuApp/internal/modules/product/ports"

	"gorm.io/gorm"
)

// MySQLProductRepository implementa ProductRepository usando MySQL
type MySQLProductRepository struct {
	db *gorm.DB
}

// NewMySQLProductRepository cria uma nova instância do repository
func NewMySQLProductRepository(db *gorm.DB) ports.ProductRepository {
	return &MySQLProductRepository{db: db}
}

// Create cria um novo produto
func (r *MySQLProductRepository) Create(ctx context.Context, product *domain.Product) error {
	productModel := &ProductModel{}
	productModel.FromDomain(product)

	if err := r.db.WithContext(ctx).Create(productModel).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetByID busca um produto por ID
func (r *MySQLProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var productModel ProductModel
	if err := r.db.WithContext(ctx).First(&productModel, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return productModel.ToDomain(), nil
}

// Update atualiza um produto existente
func (r *MySQLProductRepository) Update(ctx context.Context, product *domain.Product) error {
	productModel := &ProductModel{}
	productModel.FromDomain(product)

	if err := r.db.WithContext(ctx).Save(productModel).Error; err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// Delete remove um produto
func (r *MySQLProductRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&ProductModel{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// List lista produtos com filtros
func (r *MySQLProductRepository) List(ctx context.Context, filters ports.ProductFilters) ([]*domain.Product, error) {
	query := r.db.WithContext(ctx).Model(&ProductModel{})

	// Aplicar filtros
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}

	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	if filters.InStock != nil && *filters.InStock {
		query = query.Where("stock > 0")
	}

	var productModels []ProductModel
	if err := query.Find(&productModels).Error; err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Converter models para domain
	products := make([]*domain.Product, len(productModels))
	for i, model := range productModels {
		products[i] = model.ToDomain()
	}

	return products, nil
}
