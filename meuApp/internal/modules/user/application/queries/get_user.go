package queries

import (
"context"
"errors"

"meuApp/internal/modules/user/domain"
"meuApp/internal/modules/user/ports"
"meuApp/pkg/contracts"
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
return nil, errors.New("user ID cannot be empty")
}

user, err := h.userRepo.GetByID(ctx, query.UserID)
if err != nil {
h.logger.Error("Failed to get user by ID",
contracts.Field{Key: "user_id", Value: query.UserID},
contracts.Field{Key: "error", Value: err})
return nil, err
}

if user == nil {
return nil, errors.New("user not found")
}

return user, nil
}
