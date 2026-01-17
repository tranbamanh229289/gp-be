package dto

import (
	"be/internal/domain/proof"
	"be/internal/shared/constant"
	"time"
)

type ProofRequestUpdatedRequestDto struct {
	Status constant.ProofRequestStatus
}
type ProofRequestResponseDto struct {
	PublicID                 string                      `json:"id"`
	ThreadID                 string                      `json:"threadId"`
	VerifierDID              string                      `json:"verifierDID"`
	VerifierName             string                      `json:"verifierName"`
	SchemaID                 string                      `json:"schemaId"`
	CallbackURL              string                      `json:"callbackURL"`
	Reason                   string                      `json:"reason"`
	Message                  string                      `json:"message"`
	ScopeID                  uint32                      `json:"scopeId"`
	CircuitID                string                      `json:"circuitId"`
	AllowedIssuers           []string                    `json:"allowedIssuers"`
	CredentialSubject        map[string]interface{}      `json:"credentialSubject"`
	Context                  string                      `json:"context"`
	Type                     string                      `json:"type"`
	ProofType                string                      `json:"proofType"`
	SkipClaimRevocationCheck bool                        `json:"skipClaimRevocationCheck"`
	GroupID                  int                         `json:"groupId"`
	NullifierSession         string                      `json:"nullifierSession"`
	Status                   constant.ProofRequestStatus `json:"status"`
	ExpiresTime              *int64                      `json:"expiresTime"`
	CreatedTime              *int64                      `json:"createdTime"`
}

func ToProofRequestResponseDto(entity *proof.ProofRequest) *ProofRequestResponseDto {
	return &ProofRequestResponseDto{
		PublicID:                 entity.PublicID.String(),
		ThreadID:                 entity.ThreadID,
		VerifierDID:              entity.VerifierDID,
		VerifierName:             entity.Verifier.Name,
		SchemaID:                 entity.Schema.PublicID.String(),
		CallbackURL:              entity.CallbackURL,
		Reason:                   entity.Reason,
		Message:                  entity.Message,
		ScopeID:                  entity.ScopeID,
		CircuitID:                entity.CircuitID,
		AllowedIssuers:           entity.AllowedIssuers,
		CredentialSubject:        entity.CredentialSubject,
		Context:                  entity.Schema.ContextURL,
		Type:                     entity.Schema.Type,
		ProofType:                entity.ProofType,
		SkipClaimRevocationCheck: entity.SkipClaimRevocationCheck,
		GroupID:                  entity.GroupID,
		NullifierSession:         entity.NullifierSession,
		Status:                   entity.Status,
		ExpiresTime:              entity.ExpiresTime,
		CreatedTime:              entity.CreatedTime,
	}
}

type ProofVerificationResponseDto struct {
	Status     constant.ProofResponseStatus `json:"status"`
	Reason     string                       `json:"reason"`
	Message    string                       `json:"message"`
	ThreadID   string                       `json:"threadId"`
	HolderDID  string                       `json:"holderDID"`
	HolderName string                       `json:"holderName"`
	VerifiedAt time.Time                    `json:"verifiedAt"`
}
