package http

import (
	"net/http"
	"strconv"

	"meuApp/internal/modules/product/dto"
	"meuApp/internal/modules/product/ports"
	"meuApp/pkg/adapters/http/middleware"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService ports.ProductService
}

func NewProductHandler(productService ports.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProduct, err := h.productService.CreateProduct(c.Request.Context(), req.Name, req.Description, req.CategoryID, req.Price, req.Stock)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	response := dto.ToProductResponse(createdProduct)
	middleware.RespondWithJSON(c.Writer, http.StatusCreated, response)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	response := dto.ToProductResponse(product)
	middleware.RespondWithJSON(c.Writer, http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProduct, err := h.productService.UpdateProduct(c.Request.Context(), id, req.Name, req.Description, req.CategoryID, req.Price, req.Stock)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	response := dto.ToProductResponse(updatedProduct)
	middleware.RespondWithJSON(c.Writer, http.StatusOK, response)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse query parameters for filters
	filters := ports.ProductFilters{}

	if categoryID := c.Query("category_id"); categoryID != "" {
		filters.CategoryID = &categoryID
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	if inStockStr := c.Query("in_stock"); inStockStr != "" {
		inStock := inStockStr == "true"
		filters.InStock = &inStock
	}

	products, err := h.productService.ListProducts(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToProductResponseList(products)
	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.UpdateStock(c.Request.Context(), id, req.Stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}
