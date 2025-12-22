package repository

import (
	"be/internal/domain/document"
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
	"context"
)

type DriverLicenseRepository struct {
	db     *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewDriverLicenseRepository(db *postgres.PostgresDB, logger *logger.ZapLogger) document.IDriverLicenseRepository {
	return &DriverLicenseRepository{
		db:     db,
		logger: logger,
	}
}

func (r *DriverLicenseRepository) FindDriverLicenseByPublicId(ctx context.Context, publicId string) (*document.DriverLicense, error) {
	var entity document.DriverLicense
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *DriverLicenseRepository) FindDriverLicenseByLicenseId(ctx context.Context, licenseNumber string) (*document.DriverLicense, error) {
	var entity document.DriverLicense
	if err := r.db.GetGormDB().WithContext(ctx).Where("license_number = ?", licenseNumber).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *DriverLicenseRepository) FindAllDriverLicenses(ctx context.Context) ([]*document.DriverLicense, error) {
	var entities []*document.DriverLicense
	if err := r.db.GetGormDB().WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *DriverLicenseRepository) CreateDriverLicense(ctx context.Context, entity *document.DriverLicense) (*document.DriverLicense, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *DriverLicenseRepository) SaveDriverLicense(ctx context.Context, entity *document.DriverLicense) (*document.DriverLicense, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *DriverLicenseRepository) UpdateDriverLicense(ctx context.Context, entity *document.DriverLicense, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
