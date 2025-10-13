package auth

import "context"

type UserRepository interface {
	FindById(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindAll(ctx context.Context)([]User, error)
	Save(ctx context.Context, user User) error
}