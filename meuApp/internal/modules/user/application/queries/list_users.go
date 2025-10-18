package queries

import (
"context"

"meuApp/internal/modules/user/domain"
"meuApp/internal/modules/user/ports"
"meuApp/pkg/contracts"
)

type ListUsersQuery struct {
Page     int
PageSize int
}

type ListUsersResult struct {
Users      []*domain.User
Total      int64
Page       int
PageSize   int
TotalPages int
}

type ListUsersHandler struct {
userRepo ports.UserRepository
logger   contracts.Logger
}

func NewListUsersHandler(userRepo ports.UserRepository, logger contracts.Logger) *ListUsersHandler {
return &ListUsersHandler{userRepo: userRepo, logger: logger}
}

func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) (*ListUsersResult, error) {
if query.Page < 1 {
query.Page = 1
}
if query.PageSize < 1 {
query.PageSize = 10
}
if query.PageSize > 100 {
query.PageSize = 100
}

offset := (query.Page - 1) * query.PageSize

users, err := h.userRepo.List(ctx, offset, query.PageSize)
if err != nil {
h.logger.Error("Failed to list users", contracts.Field{Key: "error", Value: err})
return nil, err
}

total := int64(len(users))
totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))

return &ListUsersResult{
Users:      users,
Total:      total,
Page:       query.Page,
PageSize:   query.PageSize,
TotalPages: totalPages,
}, nil
}
