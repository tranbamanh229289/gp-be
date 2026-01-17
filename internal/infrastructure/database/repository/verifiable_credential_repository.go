package repository

import (
	"be/config"
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
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

func (r *VerifiableCredentialRepository) FindAllVerifiableCredentialsByHolderDID(ctx context.Context, did string) ([]*credential.VerifiableCredential, error) {
	var entities []*credential.VerifiableCredential

	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Where("holder_did = ?", did).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *VerifiableCredentialRepository) FindAllVerifiableCredentialsByIssuerDID(ctx context.Context, did string) ([]*credential.VerifiableCredential, error) {
	var entities []*credential.VerifiableCredential

	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Where("issuer_did = ?", did).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *VerifiableCredentialRepository) CreateVerifiableCredential(ctx context.Context, entity *credential.VerifiableCredential) (*credential.VerifiableCredential, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *VerifiableCredentialRepository) SaveVerifiableCredential(ctx context.Context, entity *credential.VerifiableCredential) (*credential.VerifiableCredential, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *VerifiableCredentialRepository) UpdateVerifiableCredential(ctx context.Context, entity *credential.VerifiableCredential, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
