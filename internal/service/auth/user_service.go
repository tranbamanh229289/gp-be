package auth

import (
	"be/internal/domain/auth"
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository auth.UserRepository
}

func NewUserService(userRepo auth.UserRepository) *UserService {
	return &UserService{
		userRepository: auth.UserRepository,
	}
}

func (s *UserService) Register(ctx context.Context, email, password, name string) (string, error) {
	if _, err = s.userRepository.FindByEmail(ctx, email), err != nil {

	}

	hashPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		
	}

	user := auth.User{
		PublicID: uuid.New().String(),
		Email: email,
		Password: string(hashPassword),
		Name: name,
	}

	user, err = s.userRepository.Save(ctx, user)
	if err != nil {

	}
}

func (s *UserService) Login() (string, error) {
	return "", nil
}

func (s *UserService) GetProfile(ctx context.Context, id string)(auth.User, error){
	return s.userRepository.FindById(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context)([]auth.User, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *UserService) GetToken() (string, error){
	
}