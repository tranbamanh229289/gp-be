package repository

import (
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type CredentialRequestRepository struct {
	db *postgres.PostgresDB
}

func NewCredentialRequestRepository(
	db *postgres.PostgresDB,
) credential.ICredentialRequestRepository {
	return &CredentialRequestRepository{
		db: db,
	}
}

func (r *CredentialRequestRepository) FindCredentialRequestByPublicId(ctx context.Context, publicId string) (*credential.CredentialRequest, error) {
	var credentialRequest *credential.CredentialRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Where("public_id = ?", publicId).First(credentialRequest).Error; err != nil {
		return nil, err
	}
	return credentialRequest, nil
}

func (r *CredentialRequestRepository) FindAllCredentialRequest(ctx context.Context) ([]*credential.CredentialRequest, error) {
	var credentialRequests []*credential.CredentialRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Find(credentialRequests).Error; err != nil {
		return nil, err
	}
	return credentialRequests, nil
}

func (r *CredentialRequestRepository) CreateCredentialRequest(ctx context.Context, entity *credential.CredentialRequest) (*credential.CredentialRequest, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CredentialRequestRepository) SaveCredentialRequest(ctx context.Context, entity *credential.CredentialRequest) (*credential.CredentialRequest, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CredentialRequestRepository) UpdateCredentialRequest(ctx context.Context, entity *credential.CredentialRequest, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
