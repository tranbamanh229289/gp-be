package proof

import (
	"context"
)

type IProofRepository interface {
	FindProofRequestByPublicId(ctx context.Context, id string) (*ProofRequest, error)
	FindProofRequestByThreadId(ctx context.Context, threadId string) (*ProofRequest, error)
	FindAllProofRequests(ctx context.Context) ([]*ProofRequest, error)
	FindAllProofRequestsByVerifierDID(ctx context.Context, did string) ([]*ProofRequest, error)
	CreateProofRequest(ctx context.Context, entity *ProofRequest) (*ProofRequest, error)
	UpdateProofRequest(ctx context.Context, entity *ProofRequest, changes map[string]interface{}) error

	FindProofSubmissionByPublicId(ctx context.Context, id string) (*ProofSubmission, error)
	FindAllProofSubmissions(ctx context.Context) ([]*ProofSubmission, error)
	FindAllProofSubmissionsByHolderDID(ctx context.Context, did string) ([]*ProofSubmission, error)
	FindAllProofSubmissionsByVerifierDID(ctx context.Context, did string) ([]*ProofSubmission, error)
	CreateProofSubmission(ctx context.Context, entity *ProofSubmission) (*ProofSubmission, error)
	UpdateProofSubmission(ctx context.Context, entity *ProofSubmission, changes map[string]interface{}) error
}
