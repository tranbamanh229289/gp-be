package repository

import (
	"be/config"
	"be/internal/domain/schema"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
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

func (r *SchemaAttributeRepository) FindSchemaAttributesBySchemaId(ctx context.Context, schemaId uint) ([]*schema.SchemaAttribute, error) {
	var schemaAttributes []*schema.SchemaAttribute
	if err := r.db.GetGormDB().WithContext(ctx).Where("schema_id = ?", schemaId).Find(&schemaAttributes).Error; err != nil {
		return nil, err
	}
	return schemaAttributes, nil
}

func (r *SchemaAttributeRepository) CreateSchemaAttributes(ctx context.Context, entities []*schema.SchemaAttribute) ([]*schema.SchemaAttribute, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *SchemaAttributeRepository) UpdateAttributesBySchemaId(ctx context.Context, schemaId uint, change map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(&schema.SchemaAttribute{}).Where("schema_id = ?", schemaId).Updates(change).Error; err != nil {
		return err
	}
	return nil
}
