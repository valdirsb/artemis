package commands

import (
	"context"
	"fmt"

	"meuApp/internal/modules/product/domain"
	"meuApp/internal/modules/product/ports"
	"meuApp/pkg/contracts"
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

	// Gerar ID (pode ser substituído por UUID generator)
	productID := fmt.Sprintf("prod_%d", ctx.Value("timestamp"))

	// Criar aggregate do produto
	product, err := domain.NewProduct(
		productID,
		cmd.Name,
		cmd.Description,
		cmd.CategoryID,
		cmd.Price,
		cmd.Stock,
	)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to create product aggregate: %v", err))
		return nil, fmt.Errorf("invalid product data: %w", err)
	}

	// Salvar no repositório
	if err := h.productRepo.Create(ctx, product); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to save product: %v", err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Publicar evento (assíncrono)
	go func() {
		event := contracts.Event{
			Type: "product.created",
			Payload: map[string]interface{}{
				"product_id":  product.ID,
				"name":        product.Name,
				"price":       product.Price,
				"stock":       product.Stock,
				"category_id": product.CategoryID,
			},
			Timestamp: product.CreatedAt,
		}
		h.eventBus.Publish(context.Background(), event)
		h.logger.Info(fmt.Sprintf("Event published: product.created for %s", product.ID))
	}()

	h.logger.Info(fmt.Sprintf("Product created successfully: %s", product.ID))
	return product, nil
}
