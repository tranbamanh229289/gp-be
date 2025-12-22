package schema

import (
	"context"
)

type ISchemaRepository interface {
	FindSchemaByPublicId(ctx context.Context, publicId string) (*Schema, error)
	FindAllSchemas(ctx context.Context) ([]*Schema, error)
	CreateSchema(ctx context.Context, entity *Schema) (*Schema, error)
	SaveSchema(ctx context.Context, entity *Schema) (*Schema, error)
	UpdateSchema(ctx context.Context, entity *Schema, changes map[string]interface{}) error
}
