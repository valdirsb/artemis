package grpc

import (
	"context"
	"fmt"
	"time"

	"meuApp/internal/modules/product/dto"
	"meuApp/internal/modules/product/ports"
	pb "meuApp/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler implements the gRPC ProductService
type GRPCHandler struct {
	pb.UnimplementedProductServiceServer
	productService ports.ProductService
}

// NewGRPCHandler creates a new gRPC product service
func NewGRPCHandler(productService ports.ProductService) *GRPCHandler {
	return &GRPCHandler{
		productService: productService,
	}
}

// RegisterWithServer registers the service with the gRPC server
func (s *GRPCHandler) RegisterWithServer(server *grpc.Server) {
	pb.RegisterProductServiceServer(server, s)
}

// CreateProduct creates a new product
func (s *GRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be greater than 0")
	}

	product, err := s.productService.CreateProduct(ctx, req.Name, req.Description, "default", req.Price, int(req.Quantity))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create product: %v", err))
	}

	productResponse := dto.ToProductResponse(product)
	protoProduct := &pb.Product{
		Id:          productResponse.ID,
		Name:        productResponse.Name,
		Description: productResponse.Description,
		Price:       productResponse.Price,
		Quantity:    int32(productResponse.Stock),
		CreatedAt:   productResponse.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   productResponse.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.CreateProductResponse{
		Product: protoProduct,
		Message: "Product created successfully",
	}, nil
}

// GetProduct gets a product by ID
func (s *GRPCHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	product, err := s.productService.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("product not found: %v", err))
	}

	productResponse := dto.ToProductResponse(product)
	protoProduct := &pb.Product{
		Id:          productResponse.ID,
		Name:        productResponse.Name,
		Description: productResponse.Description,
		Price:       productResponse.Price,
		Quantity:    int32(productResponse.Stock),
		CreatedAt:   productResponse.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   productResponse.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.GetProductResponse{
		Product: protoProduct,
		Message: "Product retrieved successfully",
	}, nil
}

// UpdateProduct updates a product
func (s *GRPCHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	name := req.Name
	description := req.Description
	price := req.Price
	quantity := int(req.Quantity)
	category := "default"

	product, err := s.productService.UpdateProduct(ctx, req.Id, &name, &description, &category, &price, &quantity)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update product: %v", err))
	}

	productResponse := dto.ToProductResponse(product)
	protoProduct := &pb.Product{
		Id:          productResponse.ID,
		Name:        productResponse.Name,
		Description: productResponse.Description,
		Price:       productResponse.Price,
		Quantity:    int32(productResponse.Stock),
		CreatedAt:   productResponse.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   productResponse.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.UpdateProductResponse{
		Product: protoProduct,
		Message: "Product updated successfully",
	}, nil
}

// DeleteProduct deletes a product
func (s *GRPCHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.productService.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete product: %v", err))
	}

	return &pb.DeleteProductResponse{

		Message: "Product deleted successfully",
	}, nil
}

// ListProducts lists products (simplified - no filters)
func (s *GRPCHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	filters := ports.ProductFilters{}
	products, err := s.productService.ListProducts(ctx, filters)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list products: %v", err))
	}

	var protoProducts []*pb.Product
	for _, product := range products {
		productResponse := dto.ToProductResponse(product)
		protoProducts = append(protoProducts, &pb.Product{
			Id:          productResponse.ID,
			Name:        productResponse.Name,
			Description: productResponse.Description,
			Price:       productResponse.Price,
			Quantity:    int32(productResponse.Stock),
			CreatedAt:   productResponse.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   productResponse.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(len(protoProducts)),
		Message:  "Products retrieved successfully",
	}, nil
}
