package commands

import (
	"context"

	"meuApp/internal/modules/user"
	"meuApp/internal/modules/user/domain"
	"meuApp/internal/modules/user/ports"
	"meuApp/pkg/contracts"
	apperrors "meuApp/pkg/errors"
	"meuApp/pkg/events"

	"github.com/google/uuid"
)

type CreateUserCommand struct {
	Username string
	Email    string
	Password string
}

type CreateUserHandler struct {
	userRepo       ports.UserRepository
	passwordHasher ports.PasswordHasher
	emailService   ports.EmailService
	eventPublisher contracts.EventPublisher
	typedPublisher *events.TypedEventPublisher
	logger         contracts.Logger
}

func NewCreateUserHandler(
	userRepo ports.UserRepository,
	passwordHasher ports.PasswordHasher,
	emailService ports.EmailService,
	eventPublisher contracts.EventPublisher,
	logger contracts.Logger,
) *CreateUserHandler {
	// Criar typed publisher a partir do eventPublisher (que é um EventBus)
	var typedPub *events.TypedEventPublisher
	if eventBus, ok := eventPublisher.(*events.EventBus); ok {
		typedPub = events.NewTypedEventPublisher(eventBus)
	}

	return &CreateUserHandler{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		emailService:   emailService,
		eventPublisher: eventPublisher,
		typedPublisher: typedPub,
		logger:         logger,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
	h.logger.Info("Executing CreateUserCommand", contracts.Field{Key: "email", Value: cmd.Email})

	// Validação de entrada
	if cmd.Email == "" {
		return nil, user.ErrInvalidEmail
	}
	if len(cmd.Password) < 6 {
		return nil, user.ErrInvalidPassword
	}
	if cmd.Username == "" {
		return nil, user.ErrInvalidName
	}

	// Verificar se email já existe
	existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
	if err == nil && existingUser != nil {
		return nil, user.NewEmailAlreadyExistsError(cmd.Email)
	}

	userID := uuid.New().String()

	// Criar entidade de domínio
	domainUser, err := domain.NewUser(userID, cmd.Username, cmd.Email)
	if err != nil {
		h.logger.Error("Failed to create user domain entity", contracts.Field{Key: "error", Value: err})
		return nil, apperrors.WrapError(err, "invalid user data")
	}

	// Hash da senha
	hashedPassword, err := h.passwordHasher.Hash(cmd.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", contracts.Field{Key: "error", Value: err})
		return nil, apperrors.NewInfrastructureError("failed to process password", err)
	}

	// Criar aggregate
	userAggregate := domain.NewUserAggregate(domainUser)
	userAggregate.SetPassword(hashedPassword)

	if err := userAggregate.IsValid(); err != nil {
		return nil, apperrors.WrapError(err, "invalid user aggregate")
	}

	// Persistir
	if err := h.userRepo.Create(ctx, userAggregate.GetUser()); err != nil {
		h.logger.Error("Failed to create user in repository", contracts.Field{Key: "error", Value: err})
		return nil, apperrors.NewInfrastructureError("failed to create user", err)
	}

	// ✨ Publicar evento type-safe
	if h.typedPublisher != nil {
		if err := h.typedPublisher.PublishUserCreated(ctx, events.UserCreatedEvent{
			UserID:   userID,
			Username: cmd.Username,
			Email:    cmd.Email,
		}); err != nil {
			h.logger.Warn("Failed to publish user created event", contracts.Field{Key: "error", Value: err})
		}
	}

	// Enviar email de boas-vindas (async)
	go func() {
		if err := h.emailService.SendWelcomeEmail(context.Background(), userID, cmd.Email); err != nil {
			h.logger.Warn("Failed to send welcome email", contracts.Field{Key: "error", Value: err})
		}
	}()

	h.logger.Info("User created successfully", contracts.Field{Key: "user_id", Value: userID})
	return userAggregate.GetUser(), nil
}
