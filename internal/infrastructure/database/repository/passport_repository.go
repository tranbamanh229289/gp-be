package repository

import (
	"be/internal/domain/document"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
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
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PassportRepository) FindPassportByPassportNumber(ctx context.Context, passportNumber string) (*document.Passport, error) {
	var entity document.Passport
	if err := r.db.GetGormDB().WithContext(ctx).Where("passport_number = ?", passportNumber).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PassportRepository) FindPassportByHolderDID(ctx context.Context, holderDID string) (*document.Passport, error) {
	var entity document.Passport
	if err := r.db.GetGormDB().WithContext(ctx).Where("holder_did = ?", holderDID).First(&entity).Error; err != nil {
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
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *PassportRepository) SavePassport(ctx context.Context, entity *document.Passport) (*document.Passport, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *PassportRepository) UpdatePassport(ctx context.Context, entity *document.Passport, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
