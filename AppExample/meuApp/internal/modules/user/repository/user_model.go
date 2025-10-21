package repository

import (
	"time"

	"meuApp/internal/modules/user/domain"
)

// UserModel representa a estrutura da tabela users no banco
type UserModel struct {
	ID        string    `gorm:"primaryKey;size:36"`
	Username  string    `gorm:"uniqueIndex;size:50;not null"`
	Email     string    `gorm:"uniqueIndex;size:100;not null"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName especifica o nome da tabela
func (UserModel) TableName() string {
	return "users"
}

// ToDomain converte UserModel para domain.User
func (u *UserModel) ToDomain() *domain.User {
	return &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromDomain converte domain.User para UserModel
func (u *UserModel) FromDomain(user *domain.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.Email = user.Email
	u.Password = user.Password
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}
