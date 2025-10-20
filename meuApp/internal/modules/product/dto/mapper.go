package dto

import (
	"meuApp/internal/modules/product/domain"
	"meuApp/internal/modules/product/ports"
)

// ToProductResponse converte domain.Product para ProductResponse
func ToProductResponse(product *domain.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// ToProductResponseList converte lista de domain.Product para ProductListResponse
func ToProductResponseList(products []*domain.Product) ProductListResponse {
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = ToProductResponse(product)
	}

	return ProductListResponse{
		Products: productResponses,
		Total:    len(products),
	}
}

// ToPaginatedProductResponse converte resultado paginado para PaginatedProductResponse
func ToPaginatedProductResponse(result *ports.PaginatedResult) PaginatedProductResponse {
	productResponses := make([]ProductResponse, len(result.Items))
	for i, product := range result.Items {
		productResponses[i] = ToProductResponse(product)
	}

	return PaginatedProductResponse{
		Products:   productResponses,
		TotalItems: result.TotalItems,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}

// ToProductFilters converte ProductFiltersRequest para ports.ProductFilters
func ToProductFilters(req ProductFiltersRequest) ports.ProductFilters {
	return ports.ProductFilters{
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		InStock:    req.InStock,
	}
}
