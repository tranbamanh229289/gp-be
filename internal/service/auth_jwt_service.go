package service

import (
	"be/config"
	"be/internal/domain/auth"
	"be/internal/infrastructure/cache/redis"
	"be/internal/shared/constant"
	"be/internal/transport/http/dto"
	"be/pkg/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthJWTService interface {
	GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error)
	GetProfile(ctx context.Context, id string) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, id string, user *dto.UserRequest) (*dto.UserResponse, error)
	Register(ctx context.Context, email, password, name string) (string, string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	RefreshToken(ctx context.Context, tokenString string) (string, string, error)
	VerifyToken(tokenString string, tokenType constant.TokenType) (*dto.Claims, error)
}

type AuthJWTService struct {
	logger   *logger.ZapLogger
	config   *config.Config
	redis    *redis.RedisCache
	userRepo auth.IUserRepository
}

func NewAuthJWTService(config *config.Config, logger *logger.ZapLogger, redis *redis.RedisCache, userRepo auth.IUserRepository) IAuthJWTService {
	return &AuthJWTService{
		config:   config,
		logger:   logger,
		redis:    redis,
		userRepo: userRepo,
	}
}

func (s *AuthJWTService) GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.FindAllUsers(ctx)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &constant.UserNotFound
	}
	var resp []*dto.UserResponse
	for _, u := range users {
		resp = append(resp, &dto.UserResponse{
			ID:    u.PublicID.String(),
			Name:  u.Name,
			Email: u.Email,
		})
	}

	return resp, nil
}

func (s *AuthJWTService) GetProfile(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByPublicId(ctx, id)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &constant.UserNotFound
	}
	return &dto.UserResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthJWTService) UpdateProfile(ctx context.Context, id string, userRequest *dto.UserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindUserByPublicId(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &constant.UserNotFound
		}
		return nil, &constant.InternalServer
	}

	publicId, err := uuid.Parse(id)
	if err != nil {
		return nil, &constant.InternalServer
	}

	userUpdated, err := s.userRepo.SaveUser(ctx, &auth.User{ID: user.ID, PublicID: publicId, Name: userRequest.Name, Email: userRequest.Email})

	if err != nil {
		return nil, &constant.InternalServer
	}
	return &dto.UserResponse{Name: userUpdated.Name, Email: userUpdated.Email}, nil
}

func (s *AuthJWTService) Register(ctx context.Context, email, password, name string) (string, string, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", &constant.InternalServer
	}

	if user != nil {
		return "", "", &constant.UserExisted
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", &constant.InternalServer
	}

	userCreated, err := s.userRepo.SaveUser(ctx, &auth.User{
		PublicID: uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     constant.UserRoleUser,
	})

	if err != nil {
		return "", "", &constant.InternalServer
	}

	accessToken, err := s.GetToken(&dto.Claims{
		ID:    userCreated.PublicID.String(),
		Email: email,
		Name:  name,
	}, constant.AccessToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate access token %s", err))
		return "", "", &constant.InternalServer
	}

	refreshToken, err := s.GetToken(&dto.Claims{
		ID:    userCreated.PublicID.String(),
		Email: email,
		Name:  name,
	}, constant.RefreshToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate refresh token %s", err))
		return "", "", &constant.InternalServer
	}

	return accessToken, refreshToken, nil
}

func (s *AuthJWTService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", &constant.UserNotFound
		}
		return "", "", &constant.InternalServer
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", &constant.InternalServer
	}

	if string(hashedPassword) != user.Password {
		return "", "", &constant.Unauthorized
	}

	accessToken, err := s.GetToken(&dto.Claims{
		ID:    user.PublicID.String(),
		Email: email,
		Name:  user.Name,
	}, constant.AccessToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate access token %s", err))
		return "", "", &constant.InternalServer
	}

	refreshToken, err := s.GetToken(&dto.Claims{
		ID:    user.PublicID.String(),
		Email: email,
		Name:  user.Name,
	}, constant.RefreshToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate refresh token %s", err))
		return "", "", &constant.InternalServer
	}

	return accessToken, refreshToken, nil
}

func (s *AuthJWTService) RefreshToken(ctx context.Context, tokenString string) (string, string, error) {
	claims, err := s.VerifyToken(tokenString, constant.RefreshToken)
	if err != nil {
		return "", "", &constant.InvalidToken
	}

	accessToken, err := s.GetToken(claims, constant.AccessToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate access token %s", err))
		return "", "", &constant.InternalServer
	}

	refreshToken, err := s.GetToken(claims, constant.RefreshToken)

	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to generate refresh token %s", err))
		return "", "", &constant.InternalServer
	}
	return accessToken, refreshToken, nil
}

func (s *AuthJWTService) GetToken(claims *dto.Claims, tokenType constant.TokenType) (string, error) {
	now := time.Now()
	expiration := time.Now()
	if tokenType == constant.AccessToken {
		expiration = expiration.Add(s.config.JWT.AccessTokenTTL)
	} else {
		expiration = expiration.Add(s.config.JWT.RefreshTokenTTL)
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiration),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.config.JWT.Secret)

	if err != nil {
		return "", err
	}

	tokenRedisKey := string(tokenType) + "_" + claims.ID

	if err := s.redis.Set(context.Background(), tokenRedisKey, tokenString, s.config.JWT.AccessTokenTTL); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthJWTService) VerifyToken(tokenString string, tokenType constant.TokenType) (*dto.Claims, error) {
	var claims dto.Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	tokenRedisKey := string(tokenType) + "_" + claims.ID
	tokenRedisValue, err := s.redis.Get(context.Background(), tokenRedisKey).Result()
	if err != nil {
		return nil, err
	}

	if tokenRedisValue != tokenString {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}
