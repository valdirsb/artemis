package domain

import "context"

type UsersRepository interface {
	Create(ctx context.Context, entity *Users) error
}
