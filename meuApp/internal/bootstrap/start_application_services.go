package bootstrap

import (
	userCommands "meuApp/internal/modules/user/application/commands"
	userQueries "meuApp/internal/modules/user/application/queries"
	userServices "meuApp/internal/modules/user/application/services"
	userPorts "meuApp/internal/modules/user/ports"

	productCommands "meuApp/internal/modules/product/application/commands"
	productQueries "meuApp/internal/modules/product/application/queries"
	productServices "meuApp/internal/modules/product/application/services"
	productPorts "meuApp/internal/modules/product/ports"

	orderCommands "meuApp/internal/modules/order/application/commands"
	orderQueries "meuApp/internal/modules/order/application/queries"
	orderServices "meuApp/internal/modules/order/application/services"
	orderPorts "meuApp/internal/modules/order/ports"

	"meuApp/pkg/contracts"
)

// StartApplicationServices registra os Application Services (CQRS)
func (s *Bootstrap) StartApplicationServices() {

	// ===== USER MODULE =====

	logger := s.container.MustGet("logger").(contracts.Logger)
	eventPublisher := s.container.MustGet("eventbus").(contracts.EventPublisher)

	// Dependencies do User Module
	userRepo := s.container.MustGet("userRepository").(userPorts.UserRepository)
	passwordHasher := s.container.MustGet("passwordHasher").(userPorts.PasswordHasher)
	emailService := s.container.MustGet("emailService").(userPorts.EmailService)

	// Command Handlers
	createUserHandler := userCommands.NewCreateUserHandler(
		userRepo,
		passwordHasher,
		emailService,
		eventPublisher,
		logger,
	)

	updateUserHandler := userCommands.NewUpdateUserHandler(
		userRepo,
		logger,
	)

	deleteUserHandler := userCommands.NewDeleteUserHandler(
		userRepo,
		eventPublisher,
		logger,
	)

	validateCredentialsHandler := userCommands.NewValidateCredentialsHandler(
		userRepo,
		passwordHasher,
		logger,
	)

	// Query Handlers
	getUserHandler := userQueries.NewGetUserHandler(
		userRepo,
		logger,
	)

	listUsersHandler := userQueries.NewListUsersHandler(
		userRepo,
		logger,
	)

	getUserByEmailHandler := userQueries.NewGetUserByEmailHandler(
		userRepo,
		logger,
	)

	// Application Service
	userApplicationService := userServices.NewUserApplicationService(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
		validateCredentialsHandler,
		getUserHandler,
		listUsersHandler,
		getUserByEmailHandler,
	)

	// Registrar no container
	s.Services["userApplicationService"] = userApplicationService

	// ===== PRODUCT MODULE =====

	// Dependencies do Product Module
	productRepo := s.container.MustGet("productRepository").(productPorts.ProductRepository)

	// Command Handlers
	createProductHandler := productCommands.NewCreateProductHandler(
		productRepo,
		eventPublisher,
		logger,
	)

	updateProductHandler := productCommands.NewUpdateProductHandler(
		productRepo,
		logger,
	)

	deleteProductHandler := productCommands.NewDeleteProductHandler(
		productRepo,
		logger,
	)

	updateStockHandler := productCommands.NewUpdateStockHandler(
		productRepo,
		eventPublisher,
		logger,
	)

	// Query Handlers
	getProductHandler := productQueries.NewGetProductHandler(
		productRepo,
		logger,
	)

	listProductsHandler := productQueries.NewListProductsHandler(
		productRepo,
		logger,
	)

	// Application Service
	productApplicationService := productServices.NewProductApplicationService(
		createProductHandler,
		updateProductHandler,
		deleteProductHandler,
		updateStockHandler,
		getProductHandler,
		listProductsHandler,
	)

	// Registrar no container
	s.Services["productApplicationService"] = productApplicationService

	// ===== ORDER MODULE =====

	// Dependencies do Order Module
	orderRepo := s.container.MustGet("orderRepository").(orderPorts.OrderRepository)

	// Command Handlers
	createOrderHandler := orderCommands.NewCreateOrderHandler(
		orderRepo,
		userRepo,
		productRepo,
		eventPublisher,
		logger,
	)

	updateOrderStatusHandler := orderCommands.NewUpdateOrderStatusHandler(
		orderRepo,
		eventPublisher,
		logger,
	)

	cancelOrderHandler := orderCommands.NewCancelOrderHandler(
		orderRepo,
		eventPublisher,
		logger,
	)

	// Query Handlers
	getOrderHandler := orderQueries.NewGetOrderHandler(
		orderRepo,
		logger,
	)

	getOrdersByUserHandler := orderQueries.NewGetOrdersByUserHandler(
		orderRepo,
		logger,
	)

	// Application Service
	orderApplicationService := orderServices.NewOrderApplicationService(
		createOrderHandler,
		updateOrderStatusHandler,
		cancelOrderHandler,
		getOrderHandler,
		getOrdersByUserHandler,
	)

	// Registrar no container
	s.Services["orderApplicationService"] = orderApplicationService

}
