package repository

import (
	"be/internal/domain/document"
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
	"context"
)

type PassportRepository struct {
	db     *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewPassportRepository(db *postgres.PostgresDB, logger *logger.ZapLogger) document.IPassportRepository {
	return &PassportRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PassportRepository) FindPassportByPublicId(ctx context.Context, publicId string) (*document.Passport, error) {
	var entity document.Passport
	if err := r.db.GetGormDB().WithContext(ctx).First(&entity).Where("public_id = ?", publicId).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PassportRepository) FindPassportByPassportNumber(ctx context.Context, passportNumber string) (*document.Passport, error) {
	var entity document.Passport
	if err := r.db.GetGormDB().WithContext(ctx).First(&entity).Where("passport_number = ?", passportNumber).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PassportRepository) FindAllPassports(ctx context.Context) ([]*document.Passport, error) {
	var entities []*document.Passport
	if err := r.db.GetGormDB().WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *PassportRepository) CreatePassport(ctx context.Context, entity *document.Passport) (*document.Passport, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *PassportRepository) SavePassport(ctx context.Context, entity *document.Passport) (*document.Passport, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *PassportRepository) UpdatePassport(ctx context.Context, entity *document.Passport, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
