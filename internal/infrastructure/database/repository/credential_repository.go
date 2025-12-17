package repository

import (
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type CredentialRepository struct {
	db *postgres.PostgresDB
}

func NewCredentialRepository(db *postgres.PostgresDB) credential.ICredentialRepository {
	return &CredentialRepository{
		db: db,
	}
}

func (r *CredentialRepository) FindCredentialByPublicId(ctx context.Context, publicId string) (*credential.Credential, error) {
	var credential *credential.Credential
	if err := r.db.GetGormDB().WithContext(ctx).First(credential).Where("public_id = ?", publicId).Error; err != nil {
		return nil, err
	}
	return credential, nil
}

func (r *CredentialRepository) FindAllCredential(ctx context.Context) ([]*credential.Credential, error) {
	var credentials []*credential.Credential
	if err := r.db.GetGormDB().WithContext(ctx).Find(credentials).Error; err != nil {
		return nil, err
	}
	return credentials, nil
}

func (r *CredentialRepository) CreateCredential(ctx context.Context, entity *credential.Credential) (*credential.Credential, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CredentialRepository) SaveCredential(ctx context.Context, entity *credential.Credential) (*credential.Credential, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CredentialRepository) UpdateCredential(ctx context.Context, entity *credential.Credential, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
