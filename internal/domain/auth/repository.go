package auth

type UserRepository interface {
	FindById(id int64) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll()([]*User, error)
	Create(user *User) error
	Update(user *User) error
}