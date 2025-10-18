package services

import (
"context"

"meuApp/internal/modules/product/application/commands"
"meuApp/internal/modules/product/application/queries"
"meuApp/internal/modules/product/domain"
"meuApp/internal/modules/product/ports"
)

// ProductApplicationService é o serviço de aplicação que orquestra os use cases
type ProductApplicationService struct {
createHandler      *commands.CreateProductHandler
updateHandler      *commands.UpdateProductHandler
deleteHandler      *commands.DeleteProductHandler
updateStockHandler *commands.UpdateStockHandler
getHandler         *queries.GetProductHandler
listHandler        *queries.ListProductsHandler
}

// NewProductApplicationService cria uma nova instância do serviço de aplicação
func NewProductApplicationService(
createHandler *commands.CreateProductHandler,
updateHandler *commands.UpdateProductHandler,
deleteHandler *commands.DeleteProductHandler,
updateStockHandler *commands.UpdateStockHandler,
getHandler *queries.GetProductHandler,
listHandler *queries.ListProductsHandler,
) *ProductApplicationService {
return &ProductApplicationService{
createHandler:      createHandler,
updateHandler:      updateHandler,
deleteHandler:      deleteHandler,
updateStockHandler: updateStockHandler,
getHandler:         getHandler,
listHandler:        listHandler,
}
}

// CreateProduct cria um novo produto
func (s *ProductApplicationService) CreateProduct(ctx context.Context, name, description, categoryID string, price float64, stock int) (*domain.Product, error) {
cmd := commands.CreateProductCommand{
Name:        name,
Description: description,
CategoryID:  categoryID,
Price:       price,
Stock:       stock,
}
return s.createHandler.Handle(ctx, cmd)
}

// UpdateProduct atualiza um produto existente
func (s *ProductApplicationService) UpdateProduct(ctx context.Context, id string, name, description, categoryID *string, price *float64, stock *int) (*domain.Product, error) {
cmd := commands.UpdateProductCommand{
ID:          id,
Name:        name,
Description: description,
CategoryID:  categoryID,
Price:       price,
}
return s.updateHandler.Handle(ctx, cmd)
}

// DeleteProduct deleta um produto
func (s *ProductApplicationService) DeleteProduct(ctx context.Context, id string) error {
cmd := commands.DeleteProductCommand{ID: id}
return s.deleteHandler.Handle(ctx, cmd)
}

// UpdateStock atualiza o estoque de um produto
func (s *ProductApplicationService) UpdateStock(ctx context.Context, id string, quantity int) error {
cmd := commands.UpdateStockCommand{
ID:       id,
Quantity: quantity,
}
return s.updateStockHandler.Handle(ctx, cmd)
}

// GetProductByID busca um produto por ID
func (s *ProductApplicationService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
query := queries.GetProductQuery{ID: id}
return s.getHandler.Handle(ctx, query)
}

// ListProducts lista produtos com filtros
func (s *ProductApplicationService) ListProducts(ctx context.Context, filters ports.ProductFilters) ([]*domain.Product, error) {
query := queries.ListProductsQuery{
CategoryID: filters.CategoryID,
MinPrice:   filters.MinPrice,
MaxPrice:   filters.MaxPrice,
InStock:    filters.InStock,
}
return s.listHandler.Handle(ctx, query)
}
