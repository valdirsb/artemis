package commands

import (
	"context"
	"fmt"

	"meuApp/internal/modules/product"
	"meuApp/internal/modules/product/domain"
	"meuApp/internal/modules/product/ports"
	"meuApp/pkg/contracts"
	apperrors "meuApp/pkg/errors"
)

// CreateProductCommand representa o comando para criar um produto
type CreateProductCommand struct {
	Name        string
	Description string
	CategoryID  string
	Price       float64
	Stock       int
}

// CreateProductHandler processa a criação de produtos
type CreateProductHandler struct {
	productRepo ports.ProductRepository
	eventBus    contracts.EventPublisher
	logger      contracts.Logger
}

// NewCreateProductHandler cria uma nova instância do handler
func NewCreateProductHandler(
	productRepo ports.ProductRepository,
	eventBus contracts.EventPublisher,
	logger contracts.Logger,
) *CreateProductHandler {
	return &CreateProductHandler{
		productRepo: productRepo,
		eventBus:    eventBus,
		logger:      logger,
	}
}

// Handle executa o comando de criação de produto
func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*domain.Product, error) {
	h.logger.Info(fmt.Sprintf("Creating product: %s", cmd.Name))

	// Validações de entrada
	if cmd.Name == "" {
		return nil, product.ErrInvalidName
	}
	if cmd.Description == "" {
		return nil, product.ErrInvalidDescription
	}
	if cmd.Price <= 0 {
		return nil, product.ErrInvalidPrice
	}
	if cmd.Stock < 0 {
		return nil, product.ErrInvalidStock
	}

	// Gerar ID (pode ser substituído por UUID generator)
	productID := fmt.Sprintf("prod_%d", ctx.Value("timestamp"))

	// Criar aggregate do produto
	newProduct, err := domain.NewProduct(
		productID,
		cmd.Name,
		cmd.Description,
		cmd.CategoryID,
		cmd.Price,
		cmd.Stock,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to create product aggregate: %v", err))
		return nil, apperrors.WrapError(err, "invalid product data")
	}

	// Salvar no repositório
	if err := h.productRepo.Create(ctx, newProduct); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to save product: %v", err))
		return nil, apperrors.NewInfrastructureError("failed to create product", err)
	}

	// Publicar evento (assíncrono)
	go func() {
		event := contracts.Event{
			Type: "product.created",
			Payload: map[string]interface{}{
				"product_id":  newProduct.ID,
				"name":        newProduct.Name,
				"price":       newProduct.Price,
				"stock":       newProduct.Stock,
				"category_id": newProduct.CategoryID,
			},
			Timestamp: newProduct.CreatedAt,
		}
		h.eventBus.Publish(context.Background(), event)
		h.logger.Info(fmt.Sprintf("Event published: product.created for %s", newProduct.ID))
	}()

	h.logger.Info(fmt.Sprintf("Product created successfully: %s", newProduct.ID))
	return newProduct, nil
}
