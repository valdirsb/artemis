package handler

import (
"context"
"fmt"
"time"

"meuApp/internal/modules/order/domain"
"meuApp/internal/modules/order/dto"
"meuApp/internal/modules/order/ports"
pb "meuApp/pkg/proto"

"google.golang.org/grpc"
"google.golang.org/grpc/codes"
"google.golang.org/grpc/status"
)

// OrderGRPCHandler implements the gRPC OrderService
type OrderGRPCHandler struct {
pb.UnimplementedOrderServiceServer
orderService ports.OrderService
}

// NewOrderGRPCHandler creates a new gRPC order handler
func NewOrderGRPCHandler(orderService ports.OrderService) *OrderGRPCHandler {
return &OrderGRPCHandler{
orderService: orderService,
}
}

// RegisterWithServer registers the service with the gRPC server
func (s *OrderGRPCHandler) RegisterWithServer(server *grpc.Server) {
pb.RegisterOrderServiceServer(server, s)
}

// CreateOrder creates a new order
func (s *OrderGRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
if req.UserId == "" {
return nil, status.Error(codes.InvalidArgument, "user_id is required")
}

var items []ports.CreateOrderItem
for _, v := range req.GetItems() {
if v.ProductId == "" {
return nil, status.Error(codes.InvalidArgument, "product_id is required")
}
if v.Quantity <= 0 {
return nil, status.Error(codes.InvalidArgument, "quantity must be greater than 0")
}
items = append(items, ports.CreateOrderItem{
ProductID: v.ProductId,
Quantity:  int(v.Quantity),
})
}

order, err := s.orderService.CreateOrder(ctx, req.UserId, items)
if err != nil {
return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create order: %v", err))
}

orderResponse := dto.ToOrderResponse(order)
protoOrder := &pb.Order{
Id:         orderResponse.ID,
UserId:     orderResponse.UserID,
TotalPrice: orderResponse.Total,
Status:     orderResponse.Status,
CreatedAt:  orderResponse.CreatedAt.Format(time.RFC3339),
UpdatedAt:  orderResponse.UpdatedAt.Format(time.RFC3339),
}

for _, item := range orderResponse.Items {
protoOrder.Items = append(protoOrder.Items, &pb.OrderItem{
ProductId: item.ProductID,
Quantity:  int32(item.Quantity),
Price:     item.Price,
})
}

return &pb.CreateOrderResponse{
Order:   protoOrder,
Message: "Order created successfully",
}, nil
}

// GetOrder gets an order by ID
func (s *OrderGRPCHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
if req.Id == "" {
return nil, status.Error(codes.InvalidArgument, "id is required")
}

order, err := s.orderService.GetOrderByID(ctx, req.Id)
if err != nil {
return nil, status.Error(codes.NotFound, fmt.Sprintf("order not found: %v", err))
}

orderResponse := dto.ToOrderResponse(order)
protoOrder := &pb.Order{
Id:         orderResponse.ID,
UserId:     orderResponse.UserID,
TotalPrice: orderResponse.Total,
Status:     orderResponse.Status,
CreatedAt:  orderResponse.CreatedAt.Format(time.RFC3339),
UpdatedAt:  orderResponse.UpdatedAt.Format(time.RFC3339),
}

for _, item := range orderResponse.Items {
protoOrder.Items = append(protoOrder.Items, &pb.OrderItem{
ProductId: item.ProductID,
Quantity:  int32(item.Quantity),
Price:     item.Price,
})
}

return &pb.GetOrderResponse{
Order:   protoOrder,
Message: "Order retrieved successfully",
}, nil
}

// UpdateOrderStatus updates order status
func (s *OrderGRPCHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
if req.Id == "" {
return nil, status.Error(codes.InvalidArgument, "id is required")
}

orderStatus := domain.OrderStatus(req.Status)
err := s.orderService.UpdateOrderStatus(ctx, req.Id, orderStatus)
if err != nil {
return &pb.UpdateOrderStatusResponse{
Message: fmt.Sprintf("failed to update order status: %v", err),
}, nil
}

return &pb.UpdateOrderStatusResponse{
Message: "Order status updated successfully",
}, nil
}
