package credential

import (
	"context"
)

type IIdentityRepository interface {
	FindIdentityByPublicId(ctx context.Context, publicId string) (*Identity, error)
	FindIdentityByDID(ctx context.Context, did string) (*Identity, error)
	FindIdentityByPublicKey(ctx context.Context, publicKeyX, publicKeyY string) (*Identity, error)
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

type ICredentialRequestRepository interface {
	FindCredentialRequestByPublicId(ctx context.Context, publicId string) (*CredentialRequest, error)
	FindAllCredentialRequest(ctx context.Context) ([]*CredentialRequest, error)
	CreateCredentialRequest(ctx context.Context, entity *CredentialRequest) (*CredentialRequest, error)
	SaveCredentialRequest(ctx context.Context, entity *CredentialRequest) (*CredentialRequest, error)
	UpdateCredentialRequest(ctx context.Context, entity *CredentialRequest, changes map[string]interface{}) error
}
