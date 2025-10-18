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
	userService := s.container.MustGet("userService").(userPorts.UserService)

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

	// Query Handlers
	getUserHandler := userQueries.NewGetUserHandler(
		userRepo,
		logger,
	)

	listUsersHandler := userQueries.NewListUsersHandler(
		userRepo,
		logger,
	)

	// Application Service
	userApplicationService := userServices.NewUserApplicationService(
		createUserHandler,
		updateUserHandler,
		deleteUserHandler,
		getUserHandler,
		listUsersHandler,
		userService, // Para operações que ainda não foram migradas
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

}
