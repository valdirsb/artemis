package commands

import (
	"context"

	"meuApp/internal/modules/user"
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
	h.logger.Info("Validating credentials", contracts.Field{Key: "email", Value: cmd.Email})

	// Buscar usuário por email
	foundUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err != nil || foundUser == nil {
		h.logger.Warn("Invalid login attempt - user not found", contracts.Field{Key: "email", Value: cmd.Email})
		return nil, user.ErrInvalidCredentials
	}

	// Verificar senha
	if !h.passwordHasher.Verify(cmd.Password, foundUser.Password) {
		h.logger.Warn("Invalid login attempt - wrong password", contracts.Field{Key: "email", Value: cmd.Email})
		return nil, user.ErrInvalidCredentials
	}

	h.logger.Info("Credentials validated successfully", contracts.Field{Key: "user_id", Value: foundUser.ID})
	return foundUser, nil
}
