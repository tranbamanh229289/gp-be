package auth

import "context"

type IUserRepository interface {
	FindById(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context)([]*User, error)
	Save(ctx context.Context, user *User) (*User, error)
}