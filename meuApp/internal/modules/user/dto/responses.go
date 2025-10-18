package dto

import "time"

// UserResponse representa a resposta com dados do usuário
type UserResponse struct {
ID        string    `json:"id"`
Username  string    `json:"username"`
Email     string    `json:"email"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse representa a resposta do login
type LoginResponse struct {
User         UserResponse `json:"user"`
AccessToken  string       `json:"access_token"`
RefreshToken string       `json:"refresh_token,omitempty"`
ExpiresIn    int64        `json:"expires_in"`
}

// MessageResponse representa uma resposta genérica com mensagem
type MessageResponse struct {
Message string `json:"message"`
}
