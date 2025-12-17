package proof

import (
	"context"
)

type IProofRepository interface {
	GetProofRequestByPublicId(ctx context.Context, publicId string) (*ProofRequest, error)
	GetProofRequestByRequestId(ctx context.Context, requestId string) (*ProofRequest, error)
	CreateProofRequest(ctx context.Context, entity *ProofRequest) (*ProofRequest, error)
	GetProofResponseByPublicId(ctx context.Context, publicId string) (*ProofResponse, error)
	CreateProofResponse(ctx context.Context, entity *ProofResponse) (*ProofResponse, error)
}
