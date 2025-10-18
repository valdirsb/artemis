package routes

import (
	"meuApp/pkg/container"
	"meuApp/pkg/contracts"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *container.Container, modules []string) {

	for _, v := range modules {
		switch v {
		case "user":
			registerUserRoutes(router, container)
		case "product":
			registerProductRoutes(router, container)
		case "order":
			registerOrderRoutes(router, container)
		}
	}

}

// registerUserRoutes registers user module routes
func registerUserRoutes(router *gin.Engine, container *container.Container) {
	handler := container.MustGet("userHandler").(contracts.UserHandler)
	group := router.Group("/api/v1/users")
	{
		group.POST("/", handler.CreateUser)
		group.GET("/:id", handler.GetUser)
		group.PUT("/:id", handler.UpdateUser)
		group.DELETE("/:id", handler.DeleteUser)
		group.POST("/login", handler.ValidateUser)
	}
}

// registerProductRoutes registers product module routes
func registerProductRoutes(router *gin.Engine, container *container.Container) {
	handler := container.MustGet("productHandler").(contracts.ProductHandler)
	group := router.Group("/api/v1/products")
	{

		group.POST("/", handler.CreateProduct)
		group.GET("/", handler.GetProducts)
		group.GET("/:id", handler.GetProduct)
		group.PUT("/:id", handler.UpdateProduct)
		group.DELETE("/:id", handler.DeleteProduct)
		group.PUT("/:id/stock", handler.UpdateStock)
	}
}

// registerOrderRoutes registers order module routes
func registerOrderRoutes(router *gin.Engine, container *container.Container) {
	handler := container.MustGet("orderHandler").(contracts.OrderHandler)
	group := router.Group("/api/v1/orders")
	{
		group.POST("/", handler.CreateOrder)
		group.GET("/:id", handler.GetOrder)
		group.PUT("/:id/status", handler.UpdateOrderStatus)
		group.POST("/:id/cancel", handler.CancelOrder)
		group.GET("/user/:user_id", handler.GetOrdersByUser)

	}
}
