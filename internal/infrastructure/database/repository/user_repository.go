package repository

import (
	"be/internal/domain/auth"
	"be/pkg/logger"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	logger *logger.ZapLogger
}

func NewUserRepository(db *gorm.DB, logger *logger.ZapLogger) (*UserRepository){
	return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			r.logger.Warn("User not found", zap.Int64("id", id))
			return nil, err
		}
		r.logger.Error("Failed to find user by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	r.logger.Info("Found user by id", zap.Int64("id", id))
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			r.logger.Warn("User not found", zap.String("email", email))
			return nil, err
		}
		r.logger.Error("Failed to find user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	r.logger.Info("Found user by email", zap.String("email", email))
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context)([]*auth.User, error){
	var users []*auth.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	r.logger.Info("Found all users")
	return users, nil
}

func (r *UserRepository) Save(ctx context.Context, user *auth.User) (*auth.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Error("Failed to save user", zap.Error(err))
		return nil, err
	}
	r.logger.Info("User saved successfully", zap.Int64("id", user.ID))
	return user, nil
}

