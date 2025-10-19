package http

import (
	"net/http"

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
// @Description Get all orders for a specific user
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} dto.OrderResponse
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/orders/user/{user_id} [get]
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userID := c.Param("user_id")

	orders, err := h.orderService.GetOrdersByUserID(c.Request.Context(), userID)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	response := dto.ToOrderResponseList(orders)
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
