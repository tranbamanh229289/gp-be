package repository

import (
	"be/internal/domain/auth"
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
	"context"
)

type UserRepository struct {
	db     *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewUserRepository(db *postgres.PostgresDB, logger *logger.ZapLogger) auth.IUserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (r *UserRepository) FindUserById(ctx context.Context, id int64) (*auth.User, error) {
	var user auth.User
	if err := r.db.GetGormDB().WithContext(ctx).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindUserByPublicId(ctx context.Context, id string) (*auth.User, error) {
	var user auth.User
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	var user auth.User
	if err := r.db.GetGormDB().WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAllUsers(ctx context.Context) ([]*auth.User, error) {
	var users []*auth.User
	if err := r.db.GetGormDB().WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
