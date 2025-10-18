package service

import (
	"context"
	"fmt"
	"time"

	"meuApp/internal/modules/product/domain"
	"meuApp/internal/modules/product/ports"
	"meuApp/pkg/contracts"

	"github.com/google/uuid"
)

// ProductService implementa a interface ports.ProductService
type ProductService struct {
	repo           ports.ProductRepository
	eventPublisher contracts.EventPublisher
}

// NewProductService cria uma nova instância do ProductService
func NewProductService(
	repo ports.ProductRepository,
	eventPublisher contracts.EventPublisher,
) ports.ProductService {
	return &ProductService{
		repo:           repo,
		eventPublisher: eventPublisher,
	}
}

// CreateProduct cria um novo produto
func (s *ProductService) CreateProduct(ctx context.Context, name, description, categoryID string, price float64, stock int) (*domain.Product, error) {
	// Gerar ID único
	productID := uuid.New().String()

	// Criar produto com validações de domínio
	product, err := domain.NewProduct(productID, name, description, categoryID, price, stock)
	if err != nil {
		return nil, fmt.Errorf("invalid product data: %w", err)
	}

	// Salvar no banco de dados
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Publicar evento de produto criado
	event := contracts.Event{
		Type:      "ProductCreatedEventType",
		Timestamp: time.Now(),
		Payload: contracts.ProductCreatedEvent{
			ProductID:  product.ID,
			Name:       product.Name,
			CategoryID: product.CategoryID,
			Price:      product.Price,
		},
	}

	// Ignorar erros de evento para não falhar a operação
	_ = s.eventPublisher.Publish(ctx, event)

	return product, nil
}

// GetProductByID busca um produto por ID
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// UpdateProduct atualiza um produto existente
func (s *ProductService) UpdateProduct(ctx context.Context, id string, name, description, categoryID *string, price *float64, stock *int) (*domain.Product, error) {
	// Buscar produto existente
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Criar aggregate para validações e atualizações
	aggregate := domain.NewProductAggregate(existingProduct)

	// Aplicar atualizações
	if name != nil {
		if err := aggregate.UpdateName(*name); err != nil {
			return nil, fmt.Errorf("invalid name: %w", err)
		}
	}

	if description != nil {
		if err := aggregate.UpdateDescription(*description); err != nil {
			return nil, fmt.Errorf("invalid description: %w", err)
		}
	}

	// CategoryID não pode ser atualizado (imutável por regra de negócio)
	// Se precisar mudar categoria, deve-se criar um novo produto

	if price != nil {
		if err := aggregate.UpdatePrice(*price); err != nil {
			return nil, fmt.Errorf("invalid price: %w", err)
		}
	}

	if stock != nil {
		if err := aggregate.UpdateStock(*stock); err != nil {
			return nil, fmt.Errorf("invalid stock: %w", err)
		}
	}

	// Salvar alterações
	updatedProduct := aggregate.GetProduct()
	if err := s.repo.Update(ctx, updatedProduct); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}

// DeleteProduct remove um produto
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// ListProducts lista produtos com filtros
func (s *ProductService) ListProducts(ctx context.Context, filters ports.ProductFilters) ([]*domain.Product, error) {
	products, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

// UpdateStock atualiza o estoque de um produto
func (s *ProductService) UpdateStock(ctx context.Context, id string, quantity int) error {
	// Buscar produto existente
	existingProduct, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Criar aggregate e atualizar estoque
	aggregate := domain.NewProductAggregate(existingProduct)

	if err := aggregate.UpdateStock(quantity); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Salvar alterações
	updatedProduct := aggregate.GetProduct()
	if err := s.repo.Update(ctx, updatedProduct); err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	// Publicar evento de estoque atualizado
	event := contracts.Event{
		Type:      "ProductStockUpdatedEventType",
		Timestamp: time.Now(),
		Payload: contracts.ProductStockUpdatedEvent{
			ProductID: id,
			NewStock:  quantity,
		},
	}

	// Ignorar erros de evento para não falhar a operação
	_ = s.eventPublisher.Publish(ctx, event)

	return nil
}
