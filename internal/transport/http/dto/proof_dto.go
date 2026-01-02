package dto

import (
	"be/internal/domain/proof"
	"time"
)

type ProofRequestUpdatedRequestDto struct {
	Status string
}
type ProofRequestResponseDto struct {
	PublicID       string                 `json:"id"`
	ThreadID       string                 `json:"threadId"`
	VerifierDID    string                 `json:"verifierDID"`
	CallbackURL    string                 `json:"callbackURL"`
	Reason         string                 `json:"reason"`
	Message        string                 `json:"message"`
	ScopeID        uint32                 `json:"scopeId"`
	CircuitID      string                 `json:"circuitId"`
	Params         map[string]interface{} `json:"params"`
	Query          map[string]interface{} `json:"query"`
	AllowedIssuers []string               `json:"allowedIssuers"`
	Status         string                 `json:"status"`
	ExpiresTime    *int64                 `json:"expiresTime"`
	CreatedTime    *int64                 `json:"createdTime"`
}

func ToProofRequestResponseDto(entity *proof.ProofRequest) *ProofRequestResponseDto {
	return &ProofRequestResponseDto{
		PublicID:       entity.PublicID.String(),
		ThreadID:       entity.ThreadID,
		VerifierDID:    entity.VerifierDID,
		CallbackURL:    entity.CallbackURL,
		Reason:         entity.Reason,
		Message:        entity.Message,
		ScopeID:        entity.ScopeID,
		CircuitID:      entity.CircuitID,
		Params:         entity.Params,
		Query:          entity.Query,
		AllowedIssuers: entity.AllowedIssuers,
		Status:         entity.Status,
		ExpiresTime:    entity.ExpiresTime,
		CreatedTime:    entity.CreatedTime,
	}
}

type ProofVerificationResponseDto struct {
	Status     string    `json:"status"`
	Reason     string    `json:"reason"`
	Message    string    `json:"message"`
	ThreadID   string    `json:"threadId"`
	HolderDID  string    `json:"holderDID"`
	HolderName string    `json:"holderName"`
	VerifiedAt time.Time `json:"verifiedAt"`
}
