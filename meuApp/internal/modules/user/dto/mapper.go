package dto

import (
	"meuApp/internal/modules/user/domain"
	"meuApp/internal/modules/user/ports"
)

// ToUserResponse converte domain.User para UserResponse (sem senha)
func ToUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserResponseList converte uma lista de domain.User para UserResponse
func ToUserResponseList(users []*domain.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}

// ToPaginatedUserResponse converte resultado paginado para PaginatedUserResponse
func ToPaginatedUserResponse(result *ports.PaginatedUserResult) *PaginatedUserResponse {
	userResponses := make([]UserResponse, len(result.Items))
	for i, user := range result.Items {
		userResponses[i] = *ToUserResponse(user)
	}

	return &PaginatedUserResponse{
		Users:      userResponses,
		TotalItems: result.TotalItems,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}
}
