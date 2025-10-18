package commands

import (
"context"
"fmt"

"meuApp/internal/modules/user/domain"
"meuApp/internal/modules/user/ports"
"meuApp/pkg/contracts"
)

// ValidateCredentialsCommand representa o comando para validar credenciais
type ValidateCredentialsCommand struct {
Email    string
Password string
}

// ValidateCredentialsHandler processa a validação de credenciais
type ValidateCredentialsHandler struct {
userRepo       ports.UserRepository
passwordHasher ports.PasswordHasher
logger         contracts.Logger
}

// NewValidateCredentialsHandler cria uma nova instância do handler
func NewValidateCredentialsHandler(
userRepo ports.UserRepository,
passwordHasher ports.PasswordHasher,
logger contracts.Logger,
) *ValidateCredentialsHandler {
return &ValidateCredentialsHandler{
userRepo:       userRepo,
passwordHasher: passwordHasher,
logger:         logger,
}
}

// Handle executa o comando de validação de credenciais
func (h *ValidateCredentialsHandler) Handle(ctx context.Context, cmd ValidateCredentialsCommand) (*domain.User, error) {
h.logger.Info(fmt.Sprintf("Validating credentials for email: %s", cmd.Email))

// Buscar usuário por email
user, err := h.userRepo.GetByEmail(ctx, cmd.Email)
if err != nil {
h.logger.Error(fmt.Sprintf("User not found with email: %s", cmd.Email))
return nil, fmt.Errorf("invalid credentials")
}

// Verificar senha
if !h.passwordHasher.Verify(cmd.Password, user.Password) {
h.logger.Warn(fmt.Sprintf("Invalid password for user: %s", cmd.Email))
return nil, fmt.Errorf("invalid credentials")
}

h.logger.Info(fmt.Sprintf("Credentials validated successfully for user: %s", user.ID))
return user, nil
}
