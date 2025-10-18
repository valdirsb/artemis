package commands

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
logger         contracts.Logger
}

func NewCreateUserHandler(
userRepo ports.UserRepository,
passwordHasher ports.PasswordHasher,
emailService ports.EmailService,
eventPublisher contracts.EventPublisher,
logger contracts.Logger,
) *CreateUserHandler {
return &CreateUserHandler{
userRepo:       userRepo,
passwordHasher: passwordHasher,
emailService:   emailService,
eventPublisher: eventPublisher,
logger:         logger,
}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
h.logger.Info("Executing CreateUserCommand", contracts.Field{Key: "email", Value: cmd.Email})

existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
if err == nil && existingUser != nil {
return nil, errors.New("user with this email already exists")
}

userID := uuid.New().String()

user, err := domain.NewUser(userID, cmd.Username, cmd.Email)
if err != nil {
h.logger.Error("Failed to create user domain entity", contracts.Field{Key: "error", Value: err})
return nil, err
}

hashedPassword, err := h.passwordHasher.Hash(cmd.Password)
if err != nil {
h.logger.Error("Failed to hash password", contracts.Field{Key: "error", Value: err})
return nil, errors.New("failed to process password")
}

userAggregate := domain.NewUserAggregate(user)
userAggregate.SetPassword(hashedPassword)

if err := userAggregate.IsValid(); err != nil {
return nil, err
}

if err := h.userRepo.Create(ctx, userAggregate.GetUser()); err != nil {
h.logger.Error("Failed to create user in repository", contracts.Field{Key: "error", Value: err})
return nil, errors.New("failed to create user")
}

event := contracts.Event{
Type:      events.UserCreatedEventType,
Timestamp: time.Now(),
Payload: contracts.UserCreatedEvent{
UserID: userID,
Email:  cmd.Email,
},
}

if err := h.eventPublisher.Publish(ctx, event); err != nil {
h.logger.Warn("Failed to publish user created event", contracts.Field{Key: "error", Value: err})
}

go func() {
if err := h.emailService.SendWelcomeEmail(context.Background(), userID, cmd.Email); err != nil {
h.logger.Warn("Failed to send welcome email", contracts.Field{Key: "error", Value: err})
}
}()

h.logger.Info("User created successfully", contracts.Field{Key: "user_id", Value: userID})
return userAggregate.GetUser(), nil
}
