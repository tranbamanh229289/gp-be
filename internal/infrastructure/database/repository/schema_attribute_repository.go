package repository

import (
	"be/config"
	"be/internal/domain/schema"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type SchemaAttributeRepository struct {
	config *config.Config
	db     *postgres.PostgresDB
}

func NewSchemaAttributeRepository(config *config.Config, db *postgres.PostgresDB) schema.ISchemaAttributeRepository {
	return &SchemaAttributeRepository{
		config: config,
		db:     db,
	}
}

func (r *SchemaAttributeRepository) FindSchemaAttributesBySchemaID(ctx context.Context, schemaID uint) ([]*schema.SchemaAttribute, error) {
	var schemaAttributes []*schema.SchemaAttribute
	if err := r.db.GetGormDB().WithContext(ctx).Where("schema_id = ?", schemaID).Find(&schemaAttributes).Error; err != nil {
		return nil, err
	}
	return schemaAttributes, nil
}

func (r *SchemaAttributeRepository) CreateSchemaAttributes(ctx context.Context, entities []*schema.SchemaAttribute) ([]*schema.SchemaAttribute, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *SchemaAttributeRepository) UpdateAttributesBySchemaID(ctx context.Context, schemaID uint, change map[string]interface{}) error {
	if err := r.db.GetGormDB().WithContext(ctx).Model(&schema.SchemaAttribute{}).Where("schema_id = ?", schemaID).Updates(change).Error; err != nil {
		return err
	}
	return nil
}
