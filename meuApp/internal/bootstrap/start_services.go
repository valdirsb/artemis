package bootstrap

import (
	// orderService "meuApp/internal/modules/order/service" // TODO: Refatorar Order module
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
	// TODO: Refatorar Order service para usar ports dos módulos
	orderRepo := s.container.MustGet("orderRepository").(contracts.OrderRepository)
	// productSvc := s.Services["productService"].(productPorts.ProductService)
	// userSvc := s.Services["userService"].(userPorts.UserService)
	// Temporariamente desabilitado até refatorar Order module
	_ = orderRepo
	// s.Services["orderService"] = orderService.NewOrderService(orderRepo, productSvc, userSvc, eventPublisher)

}
