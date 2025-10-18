package bootstrap

import (
	orderPorts "meuApp/internal/modules/order/ports"
	orderService "meuApp/internal/modules/order/service"
	productPorts "meuApp/internal/modules/product/ports"
	productService "meuApp/internal/modules/product/service"
	userPorts "meuApp/internal/modules/user/ports"
	userService "meuApp/internal/modules/user/service"
	"meuApp/pkg/contracts"
)

func (s *Bootstrap) StartServices() {

	eventPublisher := s.container.MustGet("eventbus").(contracts.EventPublisher)
	passwordHasher := s.container.MustGet("passwordHasher").(userPorts.PasswordHasher)
	emailService := s.container.MustGet("emailService").(userPorts.EmailService)
	tokenGenerator := s.container.MustGet("tokenGenerator").(userPorts.TokenGenerator)
	logger := s.container.MustGet("logger").(contracts.Logger)

	// User
	userRepo := s.container.MustGet("userRepository").(userPorts.UserRepository)
	s.Services["userService"] = userService.NewUserService(userRepo, passwordHasher, emailService, tokenGenerator, eventPublisher, logger)

	// Product
	productRepo := s.container.MustGet("productRepository").(productPorts.ProductRepository)
	s.Services["productService"] = productService.NewProductService(productRepo, eventPublisher)

	// Order
	orderRepo := s.container.MustGet("orderRepository").(orderPorts.OrderRepository)
	productSvc := s.Services["productService"].(productPorts.ProductService)
	userSvc := s.Services["userService"].(userPorts.UserService)
	s.Services["orderService"] = orderService.NewOrderService(orderRepo, productSvc, userSvc, eventPublisher)

}
