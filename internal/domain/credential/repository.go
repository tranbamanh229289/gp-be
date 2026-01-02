package credential

import (
	"context"
)

type IVerifiableCredentialRepository interface {
	FindVerifiableCredentialByPublicId(ctx context.Context, publicId string) (*VerifiableCredential, error)
	FindAllVerifiableCredential(ctx context.Context) ([]*VerifiableCredential, error)
	CreateVerifiableCredential(ctx context.Context, entity *VerifiableCredential) (*VerifiableCredential, error)
	SaveVerifiableCredential(ctx context.Context, entity *VerifiableCredential) (*VerifiableCredential, error)
	UpdateVerifiableCredential(ctx context.Context, entity *VerifiableCredential, changes map[string]interface{}) error
}

type ICredentialRequestRepository interface {
	FindCredentialRequestByPublicId(ctx context.Context, publicId string) (*CredentialRequest, error)
	FindAllCredentialRequest(ctx context.Context) ([]*CredentialRequest, error)
	CreateCredentialRequest(ctx context.Context, entity *CredentialRequest) (*CredentialRequest, error)
	SaveCredentialRequest(ctx context.Context, entity *CredentialRequest) (*CredentialRequest, error)
	UpdateCredentialRequest(ctx context.Context, entity *CredentialRequest, changes map[string]interface{}) error
}
