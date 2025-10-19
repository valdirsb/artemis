package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meuApp/internal/bootstrap"
	"meuApp/pkg/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "meuApp/docs" // Importar docs gerados pelo swag
)

// @title Artemis API
// @version 1.0
// @description API for Artemis - Clean Architecture Framework in Go
// @description Modular framework with Clean Architecture, Hexagonal Architecture, DDD and CQRS patterns

// @contact.name API Support
// @contact.email support@artemis.dev

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @schemes http https

// @tag.name users
// @tag.description User management endpoints

// @tag.name products
// @tag.description Product management endpoints

// @tag.name orders
// @tag.description Order management endpoints

func main() {
	log.Println("🚀 Starting meuApp with ModuleRegistry...")

	// Load application config
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load application config: %v", err)
	}

	// Bootstrap application with ModuleRegistry
	_, registry, framework, err := bootstrap.FrameworkBootstrapWithRegistry("framework.yaml")
	if err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		ctx := c.Request.Context()
		healthResults := framework.HealthCheck(ctx)

		status := "healthy"
		for _, err := range healthResults {
			if err != nil {
				status = "unhealthy"
				break
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   status,
			"services": healthResults,
			"registry": registry.Stats(),
		})
	})

	// Framework info endpoint
	router.GET("/framework/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"framework": framework.GetConfig().Framework,
			"modules":   registry.Stats(),
		})
	})

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("📚 Swagger UI available at http://localhost:8080/swagger/index.html")

	// Register all module routes from registry
	api := router.Group("/api/v1")
	registry.RegisterHTTPRoutes(api)
	log.Println("✅ All HTTP routes registered from ModuleRegistry")

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌐 Server starting on port %s", appConfig.Port)
		log.Printf("📦 Framework: %s v%s", framework.GetConfig().Framework.Name, framework.GetConfig().Framework.Version)
		log.Printf("📊 %s", registry.Stats())

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Shutdown framework
	log.Println("Shutting down framework...")
	if err := framework.Shutdown(ctx); err != nil {
		log.Printf("Framework shutdown error: %v", err)
	}

	log.Println("✅ Framework shutdown completed")
	log.Println("Server exited")
}
