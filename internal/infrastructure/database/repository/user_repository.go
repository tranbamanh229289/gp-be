package repository

import (
	"be/internal/domain/auth"
	"be/pkg/logger"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	logger *logger.ZapLogger
}

func NewUserRepository(db *gorm.DB, logger *logger.ZapLogger) auth.IUserRepository{
	return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByPublicId(ctx context.Context, id string) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).Where("public_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context)([]*auth.User, error){
	var users []*auth.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	
	return users, nil
}

func (r *UserRepository) Save(ctx context.Context, user *auth.User) (*auth.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

