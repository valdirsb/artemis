package bootstrap

import (
	orderService "meuApp/internal/modules/order/service"
	productService "meuApp/internal/modules/product/service"
	userService "meuApp/internal/modules/user/service"
	"meuApp/pkg/contracts"
)

func (s *Bootstrap) StartServices() {

	eventPublisher := s.container.MustGet("eventbus").(contracts.EventPublisher)
	passwordHasher := s.container.MustGet("passwordHasher").(contracts.PasswordHasher)
	emailService := s.container.MustGet("emailService").(contracts.EmailService)
	tokenGenerator := s.container.MustGet("tokenGenerator").(contracts.TokenGenerator)
	logger := s.container.MustGet("logger").(contracts.Logger)

	// User
	userRepo := s.container.MustGet("userRepository").(contracts.UserRepository)
	s.Services["userService"] = userService.NewUserService(userRepo, passwordHasher, emailService, tokenGenerator, eventPublisher, logger)

	// Product
	productRepo := s.container.MustGet("productRepository").(contracts.ProductRepository)
	s.Services["productService"] = productService.NewProductService(productRepo, eventPublisher)

	// Order
	orderRepo := s.container.MustGet("orderRepository").(contracts.OrderRepository)
	productSvc := s.Services["productService"].(contracts.ProductService)
	userSvc := s.Services["userService"].(contracts.UserService)
	s.Services["orderService"] = orderService.NewOrderService(orderRepo, productSvc, userSvc, eventPublisher)

}
