package repository

import (
	"be/internal/domain/document"
	"be/internal/infrastructure/database/postgres"
	"be/pkg/logger"
	"context"
)

type CitizenIdentityRepository struct {
	db     *postgres.PostgresDB
	logger *logger.ZapLogger
}

func NewCitizenIdentityRepository(db *postgres.PostgresDB, logger *logger.ZapLogger) document.ICitizenIdentityRepository {
	return &CitizenIdentityRepository{
		db:     db,
		logger: logger,
	}
}

func (r *CitizenIdentityRepository) FindCitizenIdentityByPublicId(ctx context.Context, publicId string) (*document.CitizenIdentity, error) {
	var entity document.CitizenIdentity
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *CitizenIdentityRepository) FindCitizenIdentityByIdNumber(ctx context.Context, idNumber string) (*document.CitizenIdentity, error) {
	var entity document.CitizenIdentity
	if err := r.db.GetGormDB().WithContext(ctx).Where("id_number = ?", idNumber).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *CitizenIdentityRepository) FindCitizenIdentityByHolderDID(ctx context.Context, holderDID string) (*document.CitizenIdentity, error) {
	var entity document.CitizenIdentity
	if err := r.db.GetGormDB().WithContext(ctx).Where("holder_did = ?", holderDID).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *CitizenIdentityRepository) FindAllCitizenIdentities(ctx context.Context) ([]*document.CitizenIdentity, error) {
	var entities []*document.CitizenIdentity
	if err := r.db.GetGormDB().WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *CitizenIdentityRepository) CreateCitizenIdentity(ctx context.Context, entity *document.CitizenIdentity) (*document.CitizenIdentity, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CitizenIdentityRepository) SaveCitizenIdentity(ctx context.Context, entity *document.CitizenIdentity) (*document.CitizenIdentity, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CitizenIdentityRepository) UpdateCitizenIdentity(ctx context.Context, entity *document.CitizenIdentity, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
