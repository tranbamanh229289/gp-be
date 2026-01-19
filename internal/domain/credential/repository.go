package credential

import (
	"context"
)

type IVerifiableCredentialRepository interface {
	FindVerifiableCredentialByPublicId(ctx context.Context, publicId string) (*VerifiableCredential, error)
	FindVerifiableCredentialByCredentialId(ctx context.Context, id string) (*VerifiableCredential, error)
	FindAllVerifiableCredentialsByHolderDID(ctx context.Context, did string) ([]*VerifiableCredential, error)
	FindAllVerifiableCredentialsByIssuerDID(ctx context.Context, did string) ([]*VerifiableCredential, error)
	CreateVerifiableCredential(ctx context.Context, entity *VerifiableCredential) (*VerifiableCredential, error)
	SaveVerifiableCredential(ctx context.Context, entity *VerifiableCredential) (*VerifiableCredential, error)
	UpdateVerifiableCredential(ctx context.Context, entity *VerifiableCredential, changes map[string]interface{}) error
}

type ICredentialRequestRepository interface {
	FindCredentialRequestByPublicId(ctx context.Context, publicId string) (*CredentialRequest, error)
	FindCredentialRequestByThreadId(ctx context.Context, threadId string) (*CredentialRequest, error)
	FindAllCredentialRequests(ctx context.Context) ([]*CredentialRequest, error)
	FindAllCredentialRequestsByIssuerDID(ctx context.Context, did string) ([]*CredentialRequest, error)
	FindAllCredentialRequestsByHolderDID(ctx context.Context, did string) ([]*CredentialRequest, error)
	CreateCredentialRequest(ctx context.Context, entity *CredentialRequest) (*CredentialRequest, error)
	SaveCredentialRequest(ctx context.Context, entity *CredentialRequest) (*CredentialRequest, error)
	UpdateCredentialRequest(ctx context.Context, entity *CredentialRequest, changes map[string]interface{}) error
}
