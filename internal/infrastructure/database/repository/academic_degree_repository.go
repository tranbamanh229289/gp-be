package repository

import (
	"be/internal/domain/document"
	credential "be/internal/domain/document"
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
	"context"
)

type AcademicDegreeRepository struct {
	db     *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewAcademicDegreeRepository(db *postgres.PostgresDB, logger *logger.ZapLogger) document.IAcademicDegreeRepository {
	return &AcademicDegreeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *AcademicDegreeRepository) FindAcademicDegreeByPublicId(ctx context.Context, publicId string) (*document.AcademicDegree, error) {
	var entity credential.AcademicDegree
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *AcademicDegreeRepository) FindAcademicDegreeByDegreeNumber(ctx context.Context, degreeNumber string) (*document.AcademicDegree, error) {
	var entity credential.AcademicDegree
	if err := r.db.GetGormDB().WithContext(ctx).Where("degree_number = ?", degreeNumber).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *AcademicDegreeRepository) FindAllAcademicDegrees(ctx context.Context) ([]*document.AcademicDegree, error) {
	var entities []*credential.AcademicDegree
	if err := r.db.GetGormDB().WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *AcademicDegreeRepository) CreateAcademicDegree(ctx context.Context, entity *document.AcademicDegree) (*document.AcademicDegree, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *AcademicDegreeRepository) SaveAcademicDegree(ctx context.Context, entity *document.AcademicDegree) (*document.AcademicDegree, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *AcademicDegreeRepository) UpdateAcademicDegree(ctx context.Context, entity *credential.AcademicDegree, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
