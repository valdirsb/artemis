package grpc

import (
	"context"
	"fmt"
	"time"

	"meuApp/internal/modules/user/dto"
	"meuApp/internal/modules/user/ports"
	pb "meuApp/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserGRPCHandler implements the gRPC UserService
type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userService ports.UserService
	userRepo    ports.UserRepository
}

// NewUserGRPCHandler creates a new gRPC user service
func NewUserGRPCHandler(userService ports.UserService, userRepo ports.UserRepository) *UserGRPCHandler {
	return &UserGRPCHandler{
		userService: userService,
		userRepo:    userRepo,
	}
}

// RegisterWithServer registers the service with the gRPC server
func (s *UserGRPCHandler) RegisterWithServer(server *grpc.Server) {
	pb.RegisterUserServiceServer(server, s)
}

// CreateUser creates a new user
func (s *UserGRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Call service with individual parameters
	user, err := s.userService.CreateUser(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create user: %v", err))
	}

	// Convert to proto message
	userResponse := dto.ToUserResponse(user)
	protoUser := &pb.User{
		Id:        userResponse.ID,
		Name:      userResponse.Username,
		Email:     userResponse.Email,
		CreatedAt: userResponse.CreatedAt.Format(time.RFC3339),
		UpdatedAt: userResponse.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.CreateUserResponse{
		User:    protoUser,
		Message: "User created successfully",
	}, nil
}

// GetUser gets a user by ID
func (s *UserGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found: %v", err))
	}

	protoUser := &pb.User{
		Id:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.GetUserResponse{
		User:    protoUser,
		Message: "User retrieved successfully",
	}, nil
}

// UpdateUser updates a user
func (s *UserGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	name := req.Name
	email := req.Email

	user, err := s.userService.UpdateUser(ctx, req.Id, &name, &email)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	userResponse := dto.ToUserResponse(user)
	protoUser := &pb.User{
		Id:        userResponse.ID,
		Name:      userResponse.Username,
		Email:     userResponse.Email,
		CreatedAt: userResponse.CreatedAt.Format(time.RFC3339),
		UpdatedAt: userResponse.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.UpdateUserResponse{
		User:    protoUser,
		Message: "User updated successfully",
	}, nil
}

// DeleteUser deletes a user
func (s *UserGRPCHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete user: %v", err))
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

// ListUsers lists users with pagination
func (s *UserGRPCHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Para simplificar, vamos retornar apenas uma resposta vazia por enquanto
	// Em uma implementação real, você criaria um método no repository para paginação
	return &pb.ListUsersResponse{
		Users:    []*pb.User{},
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  "Users retrieved successfully",
	}, nil
}

// GetUserByEmail gets a user by email
func (s *UserGRPCHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found: %v", err))
	}

	userResponse := dto.ToUserResponse(user)
	protoUser := &pb.User{
		Id:        userResponse.ID,
		Name:      userResponse.Username,
		Email:     userResponse.Email,
		CreatedAt: userResponse.CreatedAt.Format(time.RFC3339),
		UpdatedAt: userResponse.UpdatedAt.Format(time.RFC3339),
	}

	return &pb.GetUserResponse{
		User:    protoUser,
		Message: "User retrieved successfully",
	}, nil
}
