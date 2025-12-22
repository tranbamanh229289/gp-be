package repository

import (
	"be/internal/domain/document"
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
	"context"
)

type HealthInsuranceRepository struct {
	db     *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewHealthInsuranceRepository(db *postgres.PostgresDB, logger *logger.ZapLogger) document.IHealthInsuranceRepository {
	return &HealthInsuranceRepository{
		db:     db,
		logger: logger,
	}
}

func (r *HealthInsuranceRepository) FindHealthInsuranceByPublicId(ctx context.Context, publicId string) (*document.HealthInsurance, error) {
	var entity document.HealthInsurance
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *HealthInsuranceRepository) FindHealthInsuranceByInsuranceNumber(ctx context.Context, insuranceNumber string) (*document.HealthInsurance, error) {
	var entity document.HealthInsurance
	if err := r.db.GetGormDB().WithContext(ctx).Where("insurance_number = ?", insuranceNumber).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *HealthInsuranceRepository) FindAllHealthInsurances(ctx context.Context) ([]*document.HealthInsurance, error) {
	var entities []*document.HealthInsurance
	if err := r.db.GetGormDB().WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *HealthInsuranceRepository) CreateHealthInsurance(ctx context.Context, entity *document.HealthInsurance) (*document.HealthInsurance, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *HealthInsuranceRepository) SaveHealthInsurance(ctx context.Context, entity *document.HealthInsurance) (*document.HealthInsurance, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *HealthInsuranceRepository) UpdateHealthInsurance(ctx context.Context, entity *document.HealthInsurance, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
