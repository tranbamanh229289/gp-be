package repository

import (
	"be/internal/domain/schema"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
	"context"
)

type IdentityRepository struct {
	db *postgres.PostgresDB
}

func NewIdentityRepository(db *postgres.PostgresDB) schema.IIdentityRepository {
	return &IdentityRepository{
		db: db,
	}
}

func (r *IdentityRepository) FindIdentityByPublicId(ctx context.Context, publicId string) (*schema.Identity, error) {
	var identity schema.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_id = ?", publicId).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindIdentityByPublicKey(ctx context.Context, publicKeyX string, publicKeyY string) (*schema.Identity, error) {
	var identity schema.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Where("public_key_x = ? AND public_key_y = ?", publicKeyX, publicKeyY).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindIdentityByDID(ctx context.Context, did string) (*schema.Identity, error) {
	var identity schema.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Where("did = ?", did).First(&identity).Error; err != nil {
		return nil, err
	}
	return &identity, nil
}

func (r *IdentityRepository) FindIdentityByRole(ctx context.Context, role string) ([]*schema.Identity, error) {
	var identities []*schema.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Where("role = ?", role).Find(&identities).Error; err != nil {
		return nil, err
	}
	return identities, nil
}

func (r *IdentityRepository) FindAllIdentities(ctx context.Context) ([]*schema.Identity, error) {
	var identities []*schema.Identity
	if err := r.db.GetGormDB().WithContext(ctx).Find(&identities).Error; err != nil {
		return nil, err
	}
	return identities, nil
}

func (r *IdentityRepository) CreateIdentity(ctx context.Context, entity *schema.Identity) (*schema.Identity, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *IdentityRepository) UpdateIdentity(ctx context.Context, entity *schema.Identity, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
