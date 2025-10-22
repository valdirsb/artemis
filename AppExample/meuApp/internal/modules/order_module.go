package modules

import (
	"fmt"

	"meuApp/internal/modules/order/adapters/grpc"
	"meuApp/internal/modules/order/adapters/http"
	"meuApp/internal/modules/order/application/commands"
	"meuApp/internal/modules/order/application/queries"
	"meuApp/internal/modules/order/application/services"
	"meuApp/internal/modules/order/repository"
	productPorts "meuApp/internal/modules/product/ports"
	userPorts "meuApp/internal/modules/user/ports"
	"meuApp/pkg/container"
	"meuApp/pkg/contracts"
	"meuApp/pkg/events"

	"github.com/gin-gonic/gin"
	grpclib "google.golang.org/grpc"
	"gorm.io/gorm"
)

// OrderModule implementa o módulo de pedidos com auto-registro
// Este módulo tem dependências cross-module (User e Product)
type OrderModule struct {
	db       *gorm.DB
	eventBus *events.EventBus
	logger   contracts.Logger
}

// NewOrderModule cria uma nova instância do módulo Order
func NewOrderModule(db *gorm.DB, eventBus *events.EventBus, logger contracts.Logger) *OrderModule {
	return &OrderModule{
		db:       db,
		eventBus: eventBus,
		logger:   logger,
	}
}

// Name retorna o nome do módulo
func (m *OrderModule) Name() string {
	return "order"
}

// Register registra todos os componentes do módulo no registry
// Este método requer que os módulos User e Product já tenham sido registrados
func (m *OrderModule) Register(registry *container.ModuleRegistry) error {
	fmt.Printf("🔧 Registering module: %s\n", m.Name())

	// 1. Registrar Repository
	orderRepo := repository.NewMySQLOrderRepository(m.db)
	registry.RegisterRepository("order", orderRepo)

	// 2. Obter dependências cross-module
	// User Repository (necessário para validar usuário)
	userRepoInterface, err := registry.GetRepository("user")
	if err != nil {
		return fmt.Errorf("failed to get user repository: %w", err)
	}
	userRepo, ok := userRepoInterface.(userPorts.UserRepository)
	if !ok {
		return fmt.Errorf("user repository has wrong type")
	}

	// Product Repository (necessário para validar produtos e atualizar estoque)
	productRepoInterface, err := registry.GetRepository("product")
	if err != nil {
		return fmt.Errorf("failed to get product repository: %w", err)
	}
	productRepo, ok := productRepoInterface.(productPorts.ProductRepository)
	if !ok {
		return fmt.Errorf("product repository has wrong type")
	}

	// 3. Criar Command Handlers
	createOrderHandler := commands.NewCreateOrderHandler(
		orderRepo,
		userRepo,
		productRepo,
		m.eventBus,
		m.logger,
	)
	updateStatusHandler := commands.NewUpdateOrderStatusHandler(orderRepo, m.eventBus, m.logger)
	cancelOrderHandler := commands.NewCancelOrderHandler(orderRepo, m.eventBus, m.logger)

	// 4. Criar Query Handlers
	getOrderHandler := queries.NewGetOrderHandler(orderRepo, m.logger)
	getOrdersByUserHandler := queries.NewGetOrdersByUserHandler(orderRepo, m.logger)
	listOrdersHandler := queries.NewListOrdersHandler(orderRepo, m.logger)

	// 5. Registrar Application Service
	orderAppService := services.NewOrderApplicationService(
		createOrderHandler,
		updateStatusHandler,
		cancelOrderHandler,
		getOrderHandler,
		getOrdersByUserHandler,
		listOrdersHandler,
	)
	registry.RegisterApplicationService("order", orderAppService)

	// 6. Registrar HTTP Handler
	httpHandler := http.NewOrderHandler(orderAppService)
	registry.RegisterHTTPHandler("order", &orderHTTPHandlerAdapter{handler: httpHandler})

	// 7. Registrar gRPC Handler
	grpcHandler := grpc.NewOrderGRPCHandler(orderAppService)
	registry.RegisterGRPCService(&orderGRPCServiceAdapter{handler: grpcHandler})

	fmt.Printf("✅ Module %s registered successfully (with cross-module dependencies)\n", m.Name())
	return nil
}

// orderHTTPHandlerAdapter adapta o HTTP handler para a interface do registry
type orderHTTPHandlerAdapter struct {
	handler *http.OrderHandler
}

func (a *orderHTTPHandlerAdapter) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		orders.POST("", a.handler.CreateOrder)
		orders.GET("/:id", a.handler.GetOrder)
		orders.PUT("/:id/status", a.handler.UpdateOrderStatus)
		orders.POST("/:id/cancel", a.handler.CancelOrder)
		orders.GET("/user/:user_id", a.handler.GetOrdersByUser)
	}
}

// orderGRPCServiceAdapter adapta o gRPC handler para a interface do registry
type orderGRPCServiceAdapter struct {
	handler *grpc.OrderGRPCHandler
}

func (a *orderGRPCServiceAdapter) RegisterService(server *grpclib.Server) {
	a.handler.RegisterWithServer(server)
}
