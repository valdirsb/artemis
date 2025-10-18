package ports

import (
	"context"
	"meuApp/internal/modules/user/domain"
)

// ===== PRIMARY PORTS (Application Services) =====

// UserService define as operações de negócio do módulo de usuário
// Esta é a interface principal que expõe as funcionalidades do módulo
type UserService interface {
	CreateUser(ctx context.Context, username, email, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, username, email *string) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error)
}

// ===== SECONDARY PORTS (Adapters/Dependencies) =====

// UserRepository define a interface para persistência de usuários
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

// PasswordHasher define a interface para hashing de senhas
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

// EmailService define a interface para envio de emails (porta opcional)
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, userID, email string) error
	SendPasswordResetEmail(ctx context.Context, userID, email, token string) error
}

// TokenGenerator define a interface para geração de tokens JWT (porta opcional)
type TokenGenerator interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (string, error) // retorna userID
}
