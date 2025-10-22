package queries

import (
	"context"

	"meuApp/internal/modules/user"
	"meuApp/internal/modules/user/domain"
	"meuApp/internal/modules/user/ports"
	"meuApp/pkg/contracts"
	apperrors "meuApp/pkg/errors"
)

type GetUserQuery struct {
	UserID string
}

type GetUserHandler struct {
	userRepo ports.UserRepository
	logger   contracts.Logger
}

func NewGetUserHandler(userRepo ports.UserRepository, logger contracts.Logger) *GetUserHandler {
	return &GetUserHandler{userRepo: userRepo, logger: logger}
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*domain.User, error) {
	if query.UserID == "" {
		return nil, apperrors.NewValidationError("user ID cannot be empty", map[string]string{"field": "user_id"})
	}

	foundUser, err := h.userRepo.GetByID(ctx, query.UserID)
	if err != nil {
		h.logger.Error("Failed to get user by ID",
			contracts.Field{Key: "user_id", Value: query.UserID},
			contracts.Field{Key: "error", Value: err})
		return nil, apperrors.NewInfrastructureError("failed to retrieve user", err)
	}

	if foundUser == nil {
		return nil, user.NewUserNotFoundError(query.UserID)
	}

	return foundUser, nil
}
