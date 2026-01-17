package proof

import (
	"context"
)

type IProofRepository interface {
	FindProofRequestByPublicId(ctx context.Context, id string) (*ProofRequest, error)
	FindProofRequestByThreadId(ctx context.Context, threadID string) (*ProofRequest, error)
	FindAllProofRequests(ctx context.Context) ([]*ProofRequest, error)
	FindAllProofRequestsByVerifierDID(ctx context.Context, did string) ([]*ProofRequest, error)
	CreateProofRequest(ctx context.Context, entity *ProofRequest) (*ProofRequest, error)
	UpdateProofRequest(ctx context.Context, entity *ProofRequest, changes map[string]interface{}) error
	FindProofResponseByPublicId(ctx context.Context, id string) (*ProofResponse, error)
	CreateProofResponse(ctx context.Context, entity *ProofResponse) (*ProofResponse, error)
	FindAllProofResponses(ctx context.Context) ([]*ProofResponse, error)
}
