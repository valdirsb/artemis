package modules

import (
	"fmt"

	"meuApp/internal/modules/product/adapters/grpc"
	"meuApp/internal/modules/product/adapters/http"
	"meuApp/internal/modules/product/application/commands"
	"meuApp/internal/modules/product/application/queries"
	"meuApp/internal/modules/product/application/services"
	"meuApp/internal/modules/product/repository"
	"meuApp/pkg/container"
	"meuApp/pkg/contracts"
	"meuApp/pkg/events"

	"github.com/gin-gonic/gin"
	grpclib "google.golang.org/grpc"
	"gorm.io/gorm"
)

// ProductModule implementa o módulo de produtos com auto-registro
type ProductModule struct {
	db       *gorm.DB
	eventBus *events.EventBus
	logger   contracts.Logger
}

// NewProductModule cria uma nova instância do módulo Product
func NewProductModule(db *gorm.DB, eventBus *events.EventBus, logger contracts.Logger) *ProductModule {
	return &ProductModule{
		db:       db,
		eventBus: eventBus,
		logger:   logger,
	}
}

// Name retorna o nome do módulo
func (m *ProductModule) Name() string {
	return "product"
}

// Register registra todos os componentes do módulo no registry
func (m *ProductModule) Register(registry *container.ModuleRegistry) error {
	fmt.Printf("🔧 Registering module: %s\n", m.Name())

	// 1. Registrar Repository
	productRepo := repository.NewMySQLProductRepository(m.db)
	registry.RegisterRepository("product", productRepo)

	// 2. Criar Command Handlers
	createProductHandler := commands.NewCreateProductHandler(productRepo, m.eventBus, m.logger)
	updateProductHandler := commands.NewUpdateProductHandler(productRepo, m.logger)
	deleteProductHandler := commands.NewDeleteProductHandler(productRepo, m.logger)
	updateStockHandler := commands.NewUpdateStockHandler(productRepo, m.eventBus, m.logger)

	// 3. Criar Query Handlers
	getProductHandler := queries.NewGetProductHandler(productRepo, m.logger)
	listProductsHandler := queries.NewListProductsHandler(productRepo, m.logger)

	// 4. Registrar Application Service
	productAppService := services.NewProductApplicationService(
		createProductHandler,
		updateProductHandler,
		deleteProductHandler,
		updateStockHandler,
		getProductHandler,
		listProductsHandler,
	)
	registry.RegisterApplicationService("product", productAppService)

	// 5. Registrar HTTP Handler
	httpHandler := http.NewProductHandler(productAppService)
	registry.RegisterHTTPHandler("product", &productHTTPHandlerAdapter{handler: httpHandler})

	// 6. Registrar gRPC Handler
	grpcHandler := grpc.NewGRPCHandler(productAppService)
	registry.RegisterGRPCService(&productGRPCServiceAdapter{handler: grpcHandler})

	fmt.Printf("✅ Module %s registered successfully\n", m.Name())
	return nil
}

// productHTTPHandlerAdapter adapta o HTTP handler para a interface do registry
type productHTTPHandlerAdapter struct {
	handler *http.ProductHandler
}

func (a *productHTTPHandlerAdapter) RegisterRoutes(router *gin.RouterGroup) {
	products := router.Group("/products")
	{
		products.POST("", a.handler.CreateProduct)
		products.GET("/:id", a.handler.GetProduct)
		products.PUT("/:id", a.handler.UpdateProduct)
		products.DELETE("/:id", a.handler.DeleteProduct)
		products.GET("", a.handler.GetProducts)
		products.PATCH("/:id/stock", a.handler.UpdateStock)
	}
}

// productGRPCServiceAdapter adapta o gRPC handler para a interface do registry
type productGRPCServiceAdapter struct {
	handler *grpc.ProductGRPCHandler
}

func (a *productGRPCServiceAdapter) RegisterService(server *grpclib.Server) {
	a.handler.RegisterWithServer(server)
}
