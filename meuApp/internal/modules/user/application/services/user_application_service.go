package services

import (
"context"

"meuApp/internal/modules/user/application/commands"
"meuApp/internal/modules/user/application/queries"
"meuApp/internal/modules/user/domain"
"meuApp/internal/modules/user/ports"
)

type UserApplicationService struct {
createUserHandler *commands.CreateUserHandler
updateUserHandler *commands.UpdateUserHandler
deleteUserHandler *commands.DeleteUserHandler
getUserHandler    *queries.GetUserHandler
listUsersHandler  *queries.ListUsersHandler
userService       ports.UserService
}

func NewUserApplicationService(
createUserHandler *commands.CreateUserHandler,
updateUserHandler *commands.UpdateUserHandler,
deleteUserHandler *commands.DeleteUserHandler,
getUserHandler *queries.GetUserHandler,
listUsersHandler *queries.ListUsersHandler,
userService ports.UserService,
) *UserApplicationService {
return &UserApplicationService{
createUserHandler: createUserHandler,
updateUserHandler: updateUserHandler,
deleteUserHandler: deleteUserHandler,
getUserHandler:    getUserHandler,
listUsersHandler:  listUsersHandler,
userService:       userService,
}
}

func (s *UserApplicationService) CreateUser(ctx context.Context, username, email, password string) (*domain.User, error) {
cmd := commands.CreateUserCommand{Username: username, Email: email, Password: password}
return s.createUserHandler.Handle(ctx, cmd)
}

func (s *UserApplicationService) UpdateUser(ctx context.Context, id string, username, email *string) (*domain.User, error) {
cmd := commands.UpdateUserCommand{UserID: id, Username: username, Email: email}
return s.updateUserHandler.Handle(ctx, cmd)
}

func (s *UserApplicationService) DeleteUser(ctx context.Context, id string) error {
cmd := commands.DeleteUserCommand{UserID: id}
return s.deleteUserHandler.Handle(ctx, cmd)
}

func (s *UserApplicationService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
query := queries.GetUserQuery{UserID: id}
return s.getUserHandler.Handle(ctx, query)
}

func (s *UserApplicationService) ListUsers(ctx context.Context, page, pageSize int) (*queries.ListUsersResult, error) {
query := queries.ListUsersQuery{Page: page, PageSize: pageSize}
return s.listUsersHandler.Handle(ctx, query)
}

func (s *UserApplicationService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
return s.userService.GetUserByEmail(ctx, email)
}

func (s *UserApplicationService) ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error) {
return s.userService.ValidateCredentials(ctx, email, password)
}
