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

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with name, description, category, price and stock
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products [post]
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

// GetProduct godoc
// @Summary Get product by ID
// @Description Get a product by its unique identifier
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products/{id} [get]
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

// UpdateProduct godoc
// @Summary Update product
// @Description Update product information
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dto.UpdateProductRequest true "Product data"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products/{id} [put]
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

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetProducts godoc
// @Summary List all products
// @Description Get a list of products with optional filters and pagination
// @Tags products
// @Accept json
// @Produce json
// @Param category_id query string false "Filter by category ID"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param in_stock query boolean false "Filter only products in stock"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10)"
// @Success 200 {object} dto.PaginatedProductResponse
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products [get]
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

	// Parse pagination parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	filters.Page = page

	pageSize := 10
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}
	filters.PageSize = pageSize

	result, err := h.productService.ListProductsPaginated(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ToPaginatedProductResponse(result)
	c.JSON(http.StatusOK, response)
}

// UpdateStock godoc
// @Summary Update product stock
// @Description Update the stock quantity of a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param stock body dto.UpdateStockRequest true "New stock quantity"
// @Success 200 {object} map[string]string "Stock updated successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products/{id}/stock [patch]
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
