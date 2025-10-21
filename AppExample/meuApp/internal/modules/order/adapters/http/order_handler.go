package http

import (
	"net/http"
	"strconv"

	"meuApp/internal/modules/order/domain"
	"meuApp/internal/modules/order/dto"
	"meuApp/internal/modules/order/ports"
	"meuApp/pkg/adapters/http/middleware"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService ports.OrderService
}

func NewOrderHandler(orderService ports.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items for a specific user
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dto.CreateOrderRequest true "Order data"
// @Success 201 {object} dto.OrderResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User or Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converter DTO para ports.CreateOrderItem
	items := dto.ToCreateOrderItems(req.Items)

	createdOrder, err := h.orderService.CreateOrder(c.Request.Context(), req.UserID, items)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	response := dto.ToOrderResponse(createdOrder)
	middleware.RespondWithJSON(c.Writer, http.StatusCreated, response)
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Get an order by its unique identifier
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderResponse
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	response := dto.ToOrderResponse(order)
	middleware.RespondWithJSON(c.Writer, http.StatusOK, response)
}

// GetOrdersByUser godoc
// @Summary Get orders by user ID
// @Description Get all orders for a specific user with pagination
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Success 200 {object} dto.PaginatedOrderResponse
// @Failure 400 {object} map[string]string "Invalid pagination parameters"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/orders/user/{user_id} [get]
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userID := c.Param("user_id")

	// Parse pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	result, err := h.orderService.GetOrdersByUserIDPaginated(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	response := dto.ToPaginatedOrderResponse(result)
	middleware.RespondWithJSON(c.Writer, http.StatusOK, response)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order (pending, confirmed, shipped, delivered, cancelled)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body dto.UpdateOrderStatusRequest true "New status"
// @Success 200 {object} map[string]string "Order status updated successfully"
// @Failure 400 {object} map[string]string "Invalid request or status"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Converter string para domain.OrderStatus
	status := domain.OrderStatus(req.Status)

	if err := h.orderService.UpdateOrderStatus(c.Request.Context(), id, status); err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	middleware.RespondWithJSON(c.Writer, http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel an order if it's in pending or confirmed status
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order cancelled successfully"
// @Failure 400 {object} map[string]string "Cannot cancel order in current status"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.orderService.CancelOrder(c.Request.Context(), id); err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	middleware.RespondWithJSON(c.Writer, http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// ListOrders godoc
// @Summary List all orders
// @Description Get a paginated list of orders
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} dto.PaginatedOrderResponse
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	result, err := h.orderService.ListOrders(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToPaginatedOrderResponse(result)
	c.JSON(http.StatusOK, response)
}
