package commands

import (
"context"
"errors"
"time"

"meuApp/internal/modules/user/ports"
"meuApp/pkg/contracts"
"meuApp/pkg/events"
)

type DeleteUserCommand struct {
UserID string
}

type DeleteUserHandler struct {
userRepo       ports.UserRepository
eventPublisher contracts.EventPublisher
logger         contracts.Logger
}

func NewDeleteUserHandler(
userRepo ports.UserRepository,
eventPublisher contracts.EventPublisher,
logger contracts.Logger,
) *DeleteUserHandler {
return &DeleteUserHandler{
userRepo:       userRepo,
eventPublisher: eventPublisher,
logger:         logger,
}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
h.logger.Info("Executing DeleteUserCommand", contracts.Field{Key: "user_id", Value: cmd.UserID})

if cmd.UserID == "" {
return errors.New("user ID cannot be empty")
}

existingUser, err := h.userRepo.GetByID(ctx, cmd.UserID)
if err != nil {
return err
}

if existingUser == nil {
return errors.New("user not found")
}

if err := h.userRepo.Delete(ctx, cmd.UserID); err != nil {
h.logger.Error("Failed to delete user", contracts.Field{Key: "error", Value: err})
return errors.New("failed to delete user")
}

event := contracts.Event{
Type:      events.UserDeletedEventType,
Timestamp: time.Now(),
Payload:   map[string]string{"user_id": cmd.UserID},
}

if err := h.eventPublisher.Publish(ctx, event); err != nil {
h.logger.Warn("Failed to publish user deleted event", contracts.Field{Key: "error", Value: err})
}

h.logger.Info("User deleted successfully", contracts.Field{Key: "user_id", Value: cmd.UserID})
return nil
}
