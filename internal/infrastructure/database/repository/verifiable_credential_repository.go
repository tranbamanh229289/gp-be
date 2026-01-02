package repository

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type VerifiableCredentialRepository struct {
	db     *postgres.PostgresDB
	config *config.Config
}

func NewVerifiableCredentialRepository(db *postgres.PostgresDB, config *config.Config) credential.IVerifiableCredentialRepository {
	return &VerifiableCredentialRepository{
		db:     db,
		config: config,
	}
}

func (r *VerifiableCredentialRepository) FindVerifiableCredentialByPublicId(ctx context.Context, publicId string) (*credential.VerifiableCredential, error) {
	var entity credential.VerifiableCredential
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *VerifiableCredentialRepository) FindAllVerifiableCredential(ctx context.Context) ([]*credential.VerifiableCredential, error) {
	var entities []*credential.VerifiableCredential
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *VerifiableCredentialRepository) CreateVerifiableCredential(ctx context.Context, entity *credential.VerifiableCredential) (*credential.VerifiableCredential, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *VerifiableCredentialRepository) SaveVerifiableCredential(ctx context.Context, entity *credential.VerifiableCredential) (*credential.VerifiableCredential, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *VerifiableCredentialRepository) UpdateVerifiableCredential(ctx context.Context, entity *credential.VerifiableCredential, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
