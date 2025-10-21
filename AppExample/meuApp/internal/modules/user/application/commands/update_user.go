package commands

import (
"context"
"errors"

"meuApp/internal/modules/user/domain"
"meuApp/internal/modules/user/ports"
"meuApp/pkg/contracts"
)

type UpdateUserCommand struct {
UserID   string
Username *string
Email    *string
}

type UpdateUserHandler struct {
userRepo ports.UserRepository
logger   contracts.Logger
}

func NewUpdateUserHandler(
userRepo ports.UserRepository,
logger contracts.Logger,
) *UpdateUserHandler {
return &UpdateUserHandler{
userRepo: userRepo,
logger:   logger,
}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*domain.User, error) {
h.logger.Info("Executing UpdateUserCommand", contracts.Field{Key: "user_id", Value: cmd.UserID})

existingUser, err := h.userRepo.GetByID(ctx, cmd.UserID)
if err != nil {
return nil, err
}

if existingUser == nil {
return nil, errors.New("user not found")
}

userAggregate := domain.NewUserAggregate(existingUser)

if cmd.Email != nil {
if err := userAggregate.UpdateEmail(*cmd.Email); err != nil {
return nil, err
}
}

if cmd.Username != nil {
if err := userAggregate.UpdateUsername(*cmd.Username); err != nil {
return nil, err
}
}

if err := userAggregate.IsValid(); err != nil {
return nil, err
}

if err := h.userRepo.Update(ctx, userAggregate.GetUser()); err != nil {
h.logger.Error("Failed to update user in repository", contracts.Field{Key: "error", Value: err})
return nil, errors.New("failed to update user")
}

h.logger.Info("User updated successfully", contracts.Field{Key: "user_id", Value: cmd.UserID})
return userAggregate.GetUser(), nil
}
