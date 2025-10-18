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
	"meuApp/internal/routes"
	"meuApp/internal/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load application config
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load application config: %v", err)
	}

	// Bootstrap application with framework integration
	container, framework, err := bootstrap.FrameworkBootstrap("framework.yaml")
	if err != nil {
		log.Fatalf("Failed to bootstrap application: %v", err)
	}

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health check endpoint (always available)
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
			"features": framework.ListEnabledFeatures(),
		})
	})

	// Framework info endpoint
	router.GET("/framework/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"framework": framework.GetConfig().Framework,
			"features":  framework.ListEnabledFeatures(),
		})
	})

	// Register module routes conditionally

	routes.RegisterRoutes(router, container, framework.ListEnabledFeatures()["modules"])

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", appConfig.Port)
		log.Printf("Framework: %s v%s", framework.GetConfig().Framework.Name, framework.GetConfig().Framework.Version)
		log.Printf("Enabled features: %+v", framework.ListEnabledFeatures())

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Shutdown framework
	if err := framework.Shutdown(ctx); err != nil {
		log.Printf("Framework shutdown error: %v", err)
	}

	log.Println("Server exited")
}
