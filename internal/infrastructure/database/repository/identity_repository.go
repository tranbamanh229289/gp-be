package repository

import (
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
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
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindIdentityByPublicKey(ctx context.Context, publicKeyX string, publicKeyY string) (*credential.Identity, error) {
	var identity credential.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_key_x = ? AND public_key_y = ?", publicKeyX, publicKeyY).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindIdentityByDID(ctx context.Context, did string) (*credential.Identity, error) {
	var identity credential.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Where("did = ?", did).First(&identity).Error; err != nil {
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
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
