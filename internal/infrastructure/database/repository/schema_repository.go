package repository

import (
	"be/internal/domain/schema"
	"be/internal/infrastructure/database/postgres"
	"context"

	shell "github.com/ipfs/go-ipfs-api"
)

type SchemaRepository struct {
	db        *postgres.PostgresDB
	ipfsShell *shell.Shell
}

func NewSchemaRepository(db *postgres.PostgresDB) schema.ISchemaRepository {
	return &SchemaRepository{
		db: db,
	}
}

func (r *SchemaRepository) FindSchemaByPublicId(ctx context.Context, publicId string) (*schema.Schema, error) {
	var entity schema.Schema
	if err := r.db.GetGormDB().WithContext(ctx).Preload("SchemaAttributes").Where("public_id = ?", publicId).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SchemaRepository) FindSchemaByHash(ctx context.Context, hash string) (*schema.Schema, error) {
	var entity schema.Schema
	if err := r.db.GetGormDB().WithContext(ctx).Where("hash = ?", hash).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SchemaRepository) FindSchemaByContextURL(ctx context.Context, url string) (*schema.Schema, error) {
	var entity schema.Schema
	if err := r.db.GetGormDB().WithContext(ctx).Where("context_url = ?", url).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SchemaRepository) FindAllSchemas(ctx context.Context) ([]*schema.Schema, error) {
	var entities []*schema.Schema
	if err := r.db.GetGormDB().WithContext(ctx).Preload("SchemaAttributes").Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *SchemaRepository) CreateSchema(ctx context.Context, entity *schema.Schema) (*schema.Schema, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *SchemaRepository) SaveSchema(ctx context.Context, entity *schema.Schema) (*schema.Schema, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *SchemaRepository) UpdateSchema(ctx context.Context, entity *schema.Schema, changes map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
