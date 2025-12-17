package repository

import (
	"be/internal/domain/credential"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type SchemaRepository struct {
	db *postgres.PostgresDB
}

func NewSchemaRepository(db *postgres.PostgresDB) credential.ISchemaRepository {
	return &SchemaRepository{
		db: db,
	}
}

func (r *SchemaRepository) FindSchemaByPublicId(ctx context.Context, publicId string) (*credential.Schema, error) {
	var entity credential.Schema
	if err := r.db.GetGormDB().WithContext(ctx).First(&entity).Where("public_id = ?", publicId).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *SchemaRepository) FindAllSchemas(ctx context.Context) ([]*credential.Schema, error) {
	var entities []*credential.Schema
	if err := r.db.GetGormDB().WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *SchemaRepository) CreateSchema(ctx context.Context, entity *credential.Schema) (*credential.Schema, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *SchemaRepository) SaveSchema(ctx context.Context, entity *credential.Schema) (*credential.Schema, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
