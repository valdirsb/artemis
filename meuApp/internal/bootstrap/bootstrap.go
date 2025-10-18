package bootstrap

import (
	"context"
	"fmt"
	"log"

	"meuApp/internal/modules/user/adapters"
	"meuApp/pkg/adapters/database/mysql"
	"meuApp/pkg/config"
	"meuApp/pkg/container"
	"meuApp/pkg/events"
	"meuApp/pkg/framework"
	grpcProvider "meuApp/pkg/framework/providers/grpc"
)

type Bootstrap struct {
	Repositories map[string]interface{}
	Services     map[string]interface{}
	container    *container.Container
}

func CreateBootstrap(c *container.Container) *Bootstrap {

	return &Bootstrap{
		Repositories: make(map[string]interface{}),
		Services:     make(map[string]interface{}),
		container:    c,
	}
}

// initializeGRPCProvider initializes the gRPC provider
func initializeGRPCProvider(fw *framework.Framework) error {
	// Load app config to get gRPC port
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load app config: %v", err)
	}

	grpcProvider := grpcProvider.NewGRPCProvider(appConfig.GRPCPort)
	fw.RegisterProvider("grpc", grpcProvider)

	// Note: não inicializamos ainda, apenas registramos
	// A inicialização será feita após o registro dos serviços

	log.Printf("✅ gRPC provider registered on port %s", appConfig.GRPCPort)
	return nil
}

// FrameworkBootstrap configures the application with framework integration
func FrameworkBootstrap(configPath string) (*container.Container, *framework.Framework, error) {
	// Create DI container
	c := container.NewContainer()

	// Initialize framework
	fw, err := framework.NewFramework(configPath, c)
	if err != nil {
		return nil, nil, err
	}

	// Register framework in container for access by other services
	c.RegisterSingleton("framework", func() interface{} {
		return fw
	})

	// Initialize framework
	ctx := context.Background()
	if err := fw.Initialize(ctx); err != nil {
		return nil, nil, err
	}

	// Register traditional services based on enabled modules
	registerCoreInfrastructure(c, fw)

	registerServices(c)

	// Register handlers for enabled modules
	registerModuleHandlers(c, fw)

	// Initialize and register gRPC services if enabled
	if framework.IsEnabled("protocols", "grpc") {
		if err := initializeGRPCProvider(fw); err != nil {
			return nil, nil, err
		}
		registerGRPCServices(c, fw)
	}

	return c, fw, nil
}

// registerCoreInfrastructure registers core infrastructure services
func registerCoreInfrastructure(c *container.Container, fw *framework.Framework) {
	// Database Connection (if enabled)
	if framework.IsEnabled("database", "mysql") {
		c.RegisterSingleton("database", func() interface{} {
			cfg, err := config.LoadConfig()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			dbConfig := &mysql.DatabaseConfig{
				Host:     cfg.DBHost,
				Port:     cfg.DBPort,
				Username: cfg.DBUsername,
				Password: cfg.DBPassword,
				Database: cfg.DBDatabase,
			}
			db, err := mysql.Connect(dbConfig)
			if err != nil {
				log.Fatalf("Failed to connect to database: %v", err)
			}

			// Executar migrações
			if err := mysql.AutoMigrate(db); err != nil {
				log.Fatalf("Failed to run database migrations: %v", err)
			}

			return db
		})
	}

	// Event Bus (if enabled)
	if framework.IsEnabled("core", "event_system") {
		c.RegisterSingleton("eventbus", func() interface{} {
			return events.NewEventBus()
		})
	}

	// Password Hasher
	c.RegisterSingleton("passwordHasher", func() interface{} {
		return adapters.NewArgon2PasswordHasher()
	})

	// Mock Services (would be replaced with real implementations)
	c.RegisterSingleton("emailService", func() interface{} {
		return &MockEmailService{}
	})

	c.RegisterSingleton("tokenGenerator", func() interface{} {
		return &MockTokenGenerator{}
	})

	c.RegisterSingleton("logger", func() interface{} {
		return &SimpleLogger{}
	})
}

func registerServices(c *container.Container) {

	start := CreateBootstrap(c)

	start.StartRepositories()

	for name, repo := range start.Repositories {
		c.RegisterSingleton(name, func() interface{} {
			return repo
		})
	}

	start.StartServices()

	for name, service := range start.Services {
		c.RegisterSingleton(name, func() interface{} {
			return service
		})
	}

	// Register Application Services (CQRS)
	start.StartApplicationServices()

	for name, service := range start.Services {
		c.RegisterSingleton(name, func() interface{} {
			return service
		})
	}
}
