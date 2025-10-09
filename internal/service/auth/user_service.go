package auth

import "be/internal/domain/auth"

type UserService struct {
	repo auth.UserRepository
}
func (us *UserService) Register() (string, error) {
	return "", nil
}

func (us *UserService) Login() (string, error) {
	return "", nil
}

func (us *UserService) GetProfile(id int64)(*auth.User, error){
	return us.repo.FindById(id)
}

func (us *UserService) GetAllUsers()([]*auth.User, error) {
	return us.repo.FindAll()
}