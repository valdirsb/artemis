package queries

import (
"context"
"fmt"

"meuApp/internal/modules/user/domain"
"meuApp/internal/modules/user/ports"
"meuApp/pkg/contracts"
)

// GetUserByEmailQuery representa a query para buscar um usuário por email
type GetUserByEmailQuery struct {
Email string
}

// GetUserByEmailHandler processa a busca de usuário por email
type GetUserByEmailHandler struct {
userRepo ports.UserRepository
logger   contracts.Logger
}

// NewGetUserByEmailHandler cria uma nova instância do handler
func NewGetUserByEmailHandler(
userRepo ports.UserRepository,
logger contracts.Logger,
) *GetUserByEmailHandler {
return &GetUserByEmailHandler{
userRepo: userRepo,
logger:   logger,
}
}

// Handle executa a query de busca de usuário por email
func (h *GetUserByEmailHandler) Handle(ctx context.Context, query GetUserByEmailQuery) (*domain.User, error) {
h.logger.Info(fmt.Sprintf("Getting user by email: %s", query.Email))

user, err := h.userRepo.GetByEmail(ctx, query.Email)
if err != nil {
h.logger.Error(fmt.Sprintf("User not found with email: %s", query.Email))
return nil, fmt.Errorf("user not found: %w", err)
}

return user, nil
}
