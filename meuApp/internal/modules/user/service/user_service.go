package service

import (
	"context"
	"errors"
	"time"

	"meuApp/internal/modules/user/domain"
	"meuApp/internal/modules/user/ports"
	"meuApp/pkg/contracts"
	"meuApp/pkg/events"

	"github.com/google/uuid"
)

// UserService implementa a lógica de negócio do módulo de usuário
type UserService struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
	emailService   ports.EmailService
	tokenGenerator ports.TokenGenerator
	eventPublisher contracts.EventPublisher
	logger         contracts.Logger
}

// NewUserService cria uma nova instância do serviço de usuário
func NewUserService(
	userRepo ports.UserRepository,
	passwordHasher ports.PasswordHasher,
	emailService ports.EmailService,
	tokenGenerator ports.TokenGenerator,
	eventPublisher contracts.EventPublisher,
	logger contracts.Logger,
) ports.UserService {
	return &UserService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		emailService:   emailService,
		tokenGenerator: tokenGenerator,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// CreateUser cria um novo usuário
func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	s.logger.Info("Creating new user", contracts.Field{Key: "email", Value: email})

	// Verificar se o email já existe
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Gerar ID único
	userID := uuid.New().String()

	// Criar entidade de domínio
	user, err := domain.NewUser(userID, username, email)
	if err != nil {
		s.logger.Error("Failed to create user domain entity", contracts.Field{Key: "error", Value: err})
		return nil, err
	}

	// Hash da senha
	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		s.logger.Error("Failed to hash password", contracts.Field{Key: "error", Value: err})
		return nil, errors.New("failed to process password")
	}

	// Criar aggregate e definir senha
	userAggregate := domain.NewUserAggregate(user)
	userAggregate.SetPassword(hashedPassword)

	// Validar aggregate
	if err := userAggregate.IsValid(); err != nil {
		return nil, err
	}

	// Persistir usuário
	if err := s.userRepo.Create(ctx, userAggregate.GetUser()); err != nil {
		s.logger.Error("Failed to create user in repository", contracts.Field{Key: "error", Value: err})
		return nil, errors.New("failed to create user")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.UserCreatedEventType,
		Timestamp: time.Now(),
		Payload: contracts.UserCreatedEvent{
			UserID: userID,
			Email:  email,
		},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		s.logger.Warn("Failed to publish user created event", contracts.Field{Key: "error", Value: err})
	}

	// Enviar email de boas-vindas (assíncrono)
	go func() {
		if err := s.emailService.SendWelcomeEmail(context.Background(), userID, email); err != nil {
			s.logger.Warn("Failed to send welcome email", contracts.Field{Key: "error", Value: err})
		}
	}()

	s.logger.Info("User created successfully", contracts.Field{Key: "user_id", Value: userID})
	return userAggregate.GetUser(), nil
}

// GetUserByID obtém um usuário por ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get user by ID",
			contracts.Field{Key: "user_id", Value: id},
			contracts.Field{Key: "error", Value: err})
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByEmail obtém um usuário por email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get user by email",
			contracts.Field{Key: "email", Value: email},
			contracts.Field{Key: "error", Value: err})
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUser atualiza um usuário
func (s *UserService) UpdateUser(ctx context.Context, id string, username, email *string) (*domain.User, error) {
	// Buscar usuário existente
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// Criar aggregate
	userAggregate := domain.NewUserAggregate(existingUser)

	// Aplicar atualizações
	if email != nil {
		if err := userAggregate.UpdateEmail(*email); err != nil {
			return nil, err
		}
	}

	if username != nil {
		if err := userAggregate.UpdateUsername(*username); err != nil {
			return nil, err
		}
	}

	// Validar aggregate
	if err := userAggregate.IsValid(); err != nil {
		return nil, err
	}

	// Persistir alterações
	if err := s.userRepo.Update(ctx, userAggregate.GetUser()); err != nil {
		s.logger.Error("Failed to update user in repository", contracts.Field{Key: "error", Value: err})
		return nil, errors.New("failed to update user")
	}

	s.logger.Info("User updated successfully", contracts.Field{Key: "user_id", Value: id})
	return userAggregate.GetUser(), nil
}

// DeleteUser remove um usuário
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID cannot be empty")
	}

	// Verificar se o usuário existe
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	// Deletar usuário
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete user", contracts.Field{Key: "error", Value: err})
		return errors.New("failed to delete user")
	}

	// Publicar evento
	event := contracts.Event{
		Type:      events.UserDeletedEventType,
		Timestamp: time.Now(),
		Payload:   map[string]string{"user_id": id},
	}

	if err := s.eventPublisher.Publish(ctx, event); err != nil {
		s.logger.Warn("Failed to publish user deleted event", contracts.Field{Key: "error", Value: err})
	}

	s.logger.Info("User deleted successfully", contracts.Field{Key: "user_id", Value: id})
	return nil
}

// ValidateCredentials valida credenciais de usuário
func (s *UserService) ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Buscar usuário por email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get user by email", contracts.Field{Key: "error", Value: err})
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verificar senha
	if !s.passwordHasher.Verify(password, user.Password) {
		s.logger.Warn("Invalid password attempt", contracts.Field{Key: "email", Value: email})
		return nil, errors.New("invalid credentials")
	}

	s.logger.Info("User validated successfully", contracts.Field{Key: "user_id", Value: user.ID})
	return user, nil
}
