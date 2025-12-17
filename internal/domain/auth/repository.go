package auth

import "context"

type IUserRepository interface {
	FindUserById(ctx context.Context, id int64) (*User, error)
	FindUserByPublicId(ctx context.Context, id string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindAllUsers(ctx context.Context) ([]*User, error)
	SaveUser(ctx context.Context, user *User) (*User, error)
}
