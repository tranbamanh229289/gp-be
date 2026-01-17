package schema

import (
	"context"
)

type IIdentityRepository interface {
	FindIdentityByPublicId(ctx context.Context, publicId string) (*Identity, error)
	FindIdentityByDID(ctx context.Context, did string) (*Identity, error)
	FindIdentityByPublicKey(ctx context.Context, publicKeyX, publicKeyY string) (*Identity, error)
	FindIdentityByRole(ctx context.Context, role string) ([]*Identity, error)
	FindAllIdentities(ctx context.Context) ([]*Identity, error)
	CreateIdentity(ctx context.Context, entity *Identity) (*Identity, error)
	UpdateIdentity(ctx context.Context, entity *Identity, changes map[string]interface{}) error
}

type ISchemaRepository interface {
	FindSchemaByPublicId(ctx context.Context, publicId string) (*Schema, error)
	FindSchemaByHash(ctx context.Context, hash string) (*Schema, error)
	FindSchemaByContextURL(ctx context.Context, hash string) (*Schema, error)
	FindAllSchemas(ctx context.Context) ([]*Schema, error)
	CreateSchema(ctx context.Context, schema *Schema) (*Schema, error)
	UpdateSchema(ctx context.Context, entity *Schema, changes map[string]interface{}) error
}

type ISchemaAttributeRepository interface {
	CreateSchemaAttributes(ctx context.Context, entities []*SchemaAttribute) ([]*SchemaAttribute, error)
	FindSchemaAttributesBySchemaId(ctx context.Context, schemaId uint) ([]*SchemaAttribute, error)
	UpdateAttributesBySchemaId(ctx context.Context, schemaId uint, change map[string]interface{}) error
}
