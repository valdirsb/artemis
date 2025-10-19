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

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	if err := h.orderService.CancelOrder(c.Request.Context(), id); err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	middleware.RespondWithJSON(c.Writer, http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
