package credential

import (
	"be/internal/shared/types"
	"context"
)

type IIdentityRepository interface {
	FindIdentityByPublicId(ctx context.Context, publicId string) (*Identity, error)
	FindIdentityByDID(ctx context.Context, did types.DID) (*Identity, error)
	FindAllIdentities(ctx context.Context) ([]*Identity, error)
	CreateIdentity(ctx context.Context, entity *Identity) (*Identity, error)
	UpdateCitizenIdentity(ctx context.Context, entity *Identity, changes map[string]interface{}) error
}

type ICredentialRepository interface {
	FindCredentialByPublicId(ctx context.Context, publicId string) (*Credential, error)
	FindAllCredential(ctx context.Context) ([]*Credential, error)
	CreateCredential(ctx context.Context, entity *Credential) (*Credential, error)
	SaveCredential(ctx context.Context, entity *Credential) (*Credential, error)
	UpdateCredential(ctx context.Context, entity *Credential, changes map[string]interface{}) error
}

type ISchemaRepository interface {
	FindSchemaByPublicId(ctx context.Context, publicId string) (*Schema, error)
	FindAllSchemas(ctx context.Context) ([]*Schema, error)
	CreateSchema(ctx context.Context, entity *Schema) (*Schema, error)
	SaveSchema(ctx context.Context, entity *Schema) (*Schema, error)
}
