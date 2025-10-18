package bootstrap

import (
	"context"
	"log"
	orderHandler "meuApp/internal/modules/order/handler"
	orderPorts "meuApp/internal/modules/order/ports"
	productGRPC "meuApp/internal/modules/product/adapters/grpc"
	productHTTP "meuApp/internal/modules/product/adapters/http"
	productPorts "meuApp/internal/modules/product/ports"
	userGRPC "meuApp/internal/modules/user/adapters/grpc"
	userHTTP "meuApp/internal/modules/user/adapters/http"
	userPorts "meuApp/internal/modules/user/ports"
	"meuApp/pkg/container"
	"meuApp/pkg/framework"
	"meuApp/pkg/framework/providers"
	grpcProvider "meuApp/pkg/framework/providers/grpc"
)

// registerModuleHandlers registers HTTP handlers for enabled modules
func registerModuleHandlers(c *container.Container, fw *framework.Framework) {

	if framework.IsEnabled("modules", "user") {
		c.RegisterSingleton("userHandler", func() interface{} {
			// Usar Application Service em vez do Service tradicional
			userAppService := c.MustGet("userApplicationService").(userPorts.UserService)
			return userHTTP.NewUserHTTPHandler(userAppService)
		})
	}

	if framework.IsEnabled("modules", "product") {
		c.RegisterSingleton("productHandler", func() interface{} {
			// Usar Application Service em vez do Service tradicional
			productAppService := c.MustGet("productApplicationService").(productPorts.ProductService)
			return productHTTP.NewProductHandler(productAppService)
		})
	}

	if framework.IsEnabled("modules", "order") {
		c.RegisterSingleton("orderHandler", func() interface{} {
			orderSvc := c.MustGet("orderService").(orderPorts.OrderService)
			return orderHandler.NewOrderHandler(orderSvc)
		})
	}
}

// registerGRPCServices registers gRPC services with the gRPC provider
func registerGRPCServices(c *container.Container, fw *framework.Framework) {
	// Get gRPC provider from framework
	grpcProv := fw.GetProvider("grpc")
	if grpcProv == nil {
		log.Printf("gRPC provider not found")
		return
	}

	grpcProvider, ok := grpcProv.(*grpcProvider.GRPCProvider)
	if !ok {
		log.Printf("Invalid gRPC provider type")
		return
	}

	// Register User gRPC Service
	if framework.IsEnabled("modules", "user") {
		// Usar Application Service em vez do Service tradicional
		userAppSvc := c.MustGet("userApplicationService").(userPorts.UserService)
		userRepo := c.MustGet("userRepository").(userPorts.UserRepository)
		userGRPCService := userGRPC.NewUserGRPCHandler(userAppSvc, userRepo)
		grpcProvider.RegisterService(userGRPCService)
		log.Printf("✅ User gRPC service registered")
	}

	// Register Product gRPC Service
	if framework.IsEnabled("modules", "product") {
		// Usar Application Service em vez do Service tradicional
		productAppSvc := c.MustGet("productApplicationService").(productPorts.ProductService)
		productGRPCService := productGRPC.NewGRPCHandler(productAppSvc)
		grpcProvider.RegisterService(productGRPCService)
		log.Printf("✅ Product gRPC service registered")
	}

	// Register Order gRPC Service
	if framework.IsEnabled("modules", "order") {
		orderSvc := c.MustGet("orderService").(orderPorts.OrderService)
		orderGRPCService := orderHandler.NewOrderGRPCHandler(orderSvc)
		grpcProvider.RegisterService(orderGRPCService)
		log.Printf("✅ Order gRPC service registered")
	}

	// Agora inicializar o provider com todos os serviços registrados
	ctx := context.Background()
	deps := providers.Dependencies{
		Config: fw.GetConfig(),
	}
	if err := grpcProvider.Initialize(ctx, deps); err != nil {
		log.Printf("❌ Failed to initialize gRPC provider: %v", err)
		return
	}

	log.Printf("✅ gRPC server started successfully on port %s", grpcProvider.GetPort())
}
