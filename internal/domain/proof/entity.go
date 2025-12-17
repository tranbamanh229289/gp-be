package proof

import (
	"be/internal/domain/credential"
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type ProofRequest struct {
	ID                uint
	PublicID          uuid.UUID
	VerifierID        uint
	SchemaID          uint
	CircuitType       string
	RequestData       types.JSONB
	AllowedIssuersDID []string
	Challenge         string
	Status            string
	QRCodeData        string
	CallbackURL       string
	ExpireAt          time.Time
	CreatedAt         time.Time
	Verifier          *credential.Identity
	Schema            *credential.Schema
	Responses         []*ProofResponse
}

type ProofResponse struct {
	ID            uint
	PublicID      uuid.UUID
	RequestID     uint
	HolderID      uint
	ProofData     types.JSONB
	PublicSignals types.JSONB
	IsVerified    bool
	VerifiedAt    time.Time
	CreatedAt     time.Time
	Request       ProofRequest
	Holder        *credential.Identity
}
