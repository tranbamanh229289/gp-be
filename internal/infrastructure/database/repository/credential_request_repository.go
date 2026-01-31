package repository

import (
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
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
	var credentialRequest credential.CredentialRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Preload("Issuer").Preload("Holder").Where("public_id = ?", publicId).First(&credentialRequest).Error; err != nil {
		return nil, err
	}

	return &credentialRequest, nil
}

func (r *CredentialRequestRepository) FindCredentialRequestByThreadId(ctx context.Context, threadId string) (*credential.CredentialRequest, error) {
	var credentialRequest credential.CredentialRequest
	if err := r.db.GetGormDB().WithContext(ctx).Where("thread_id = ?", threadId).First(&credentialRequest).Error; err != nil {
		return nil, err
	}
	return &credentialRequest, nil
}

func (r *CredentialRequestRepository) FindAllCredentialRequests(ctx context.Context) ([]*credential.CredentialRequest, error) {
	var credentialRequests []*credential.CredentialRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Preload("Issuer").Preload("Holder").Find(&credentialRequests).Error; err != nil {
		return nil, err
	}
	return credentialRequests, nil
}

func (r *CredentialRequestRepository) FindAllCredentialRequestsByHolderDID(ctx context.Context, did string) ([]*credential.CredentialRequest, error) {
	var credentialRequests []*credential.CredentialRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Preload("Issuer").Preload("Holder").Where("holder_did", did).Find(&credentialRequests).Error; err != nil {
		return nil, err
	}
	return credentialRequests, nil
}

func (r *CredentialRequestRepository) FindAllCredentialRequestsByIssuerDID(ctx context.Context, did string) ([]*credential.CredentialRequest, error) {
	var credentialRequests []*credential.CredentialRequest
	if err := r.db.GetGormDB().WithContext(ctx).Preload("Schema").Preload("Issuer").Preload("Holder").Where("issuer_did", did).Find(&credentialRequests).Error; err != nil {
		return nil, err
	}
	return credentialRequests, nil
}

func (r *CredentialRequestRepository) CreateCredentialRequest(ctx context.Context, entity *credential.CredentialRequest) (*credential.CredentialRequest, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CredentialRequestRepository) SaveCredentialRequest(ctx context.Context, entity *credential.CredentialRequest) (*credential.CredentialRequest, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CredentialRequestRepository) UpdateCredentialRequest(ctx context.Context, entity *credential.CredentialRequest, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
