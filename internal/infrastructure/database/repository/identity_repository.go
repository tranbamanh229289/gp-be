package repository

import (
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/types"
	"context"
)

type IdentityRepository struct {
	db *postgres.PostgresDB
}

func NewIdentityRepository(db *postgres.PostgresDB) credential.IIdentityRepository {
	return &IdentityRepository{
		db: db,
	}
}

func (r *IdentityRepository) FindIdentityByPublicId(ctx context.Context, publicId string) (*credential.Identity, error) {
	var identity credential.Identity
	if err := r.db.GetGormDB().WithContext(ctx).First(&identity).Where("public_id = ?", publicId).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindIdentityByDID(ctx context.Context, did types.DID) (*credential.Identity, error) {
	var identity credential.Identity
	if err := r.db.GetGormDB().WithContext(ctx).First(&identity).Where("did = ?", did).Error; err != nil {
		return nil, err
	}
	return &identity, nil

}

func (r *IdentityRepository) FindAllIdentities(ctx context.Context) ([]*credential.Identity, error) {
	var identities []*credential.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Find(&identities).Error; err != nil {
		return nil, err
	}
	return identities, nil
}

func (r *IdentityRepository) CreateIdentity(ctx context.Context, entity *credential.Identity) (*credential.Identity, error) {
	if err := r.db.GetGormDB().Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *IdentityRepository) UpdateCitizenIdentity(ctx context.Context, entity *credential.Identity, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
