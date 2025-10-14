package service

import (
	"be/internal/domain/auth"
	"be/internal/infrastructure/database/repository"
	"be/pkg/logger"
	"context"
)

type IAuthService interface {
	Register(ctx context.Context, email, password, name string) (string, error)
	Login(ctx context.Context, email, password string) (string error)
	RefreshToken(ctx context.Context) (string, error)
	GetProfile(ctx context.Context, id string) (auth.User, error)
	UpdateProfile(ctx context.Context, user auth.User) (auth.User, error)
	GetAllUsers(ctx context.Context) ([]auth.User, error)
}

type AuthService struct {
	userRepository auth.IUserRepository
	logger *logger.ZapLogger
}

func NewAuthService(userRepo repository.UserRepository, logger *logger.ZapLogger) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		logger: logger,
	}
}

func (s *AuthService) GetAllUsers(ctx context.Context)([]auth.User, error) {
	return s.userRepository.FindAll(ctx)
}

func (s *AuthService) GetProfile(ctx context.Context, id string) (auth.User, error) {
	return s.userRepository.FindById(ctx, id)
}

func(s *AuthService) UpdateProfile(ctx context.Context, user auth.User) (auth.User, error) {

}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (string, error) {

}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

func (s *AuthService) RefreshToken(ctx context.Context) (string, error) {
	return "", nil
}


func GetToken() (string, error){
	
}
