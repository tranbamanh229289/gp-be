package dto

import (
	"be/internal/domain/proof"
	"be/internal/shared/constant"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/iden3comm/v2/packers"
	"github.com/iden3/iden3comm/v2/protocol"
)

type ProofRequestUpdatedRequestDto struct {
	Status constant.ProofRequestStatus
}

func ToAuthorizationRequest(pr *proof.ProofRequest) protocol.AuthorizationRequestMessage {
	query := make(map[string]interface{})
	params := make(map[string]interface{})
	query["allowedIssuers"] = pr.AllowedIssuers
	query["context"] = pr.Schema.ContextURL
	query["type"] = pr.Schema.Type
	query["credentialSubject"] = pr.CredentialSubject
	query["proofType"] = pr.ProofType
	query["skipClaimRevocationCheck"] = pr.SkipClaimRevocationCheck
	query["groupId"] = pr.GroupID
	params["nullifierSessionId"] = pr.NullifierSession

	return protocol.AuthorizationRequestMessage{
		ID:       pr.ThreadID,
		From:     pr.VerifierDID,
		Typ:      packers.MediaTypePlainMessage,
		Type:     protocol.AuthorizationRequestMessageType,
		ThreadID: pr.ThreadID,
		Body: protocol.AuthorizationRequestMessageBody{
			CallbackURL: pr.CallbackURL,
			Reason:      pr.Reason,
			Message:     pr.Message,
			Scope: []protocol.ZeroKnowledgeProofRequest{
				{
					ID:        pr.ScopeID,
					CircuitID: pr.CircuitID,
					Query:     query,
					Params:    params,
				},
			},
		},

		ExpiresTime: pr.ExpiresTime,
		CreatedTime: pr.CreatedTime,
	}
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

type ProofSubmissionResponseDto struct {
	PublicID     string                         `json:"id"`
	RequestID    string                         `json:"requestId"`
	ThreadID     string                         `json:"threadId"`
	HolderDID    string                         `json:"holderDID"`
	HolderName   string                         `json:"holderName"`
	VerifierDID  string                         `json:"verifierDID"`
	VerifierName string                         `json:"verifierName"`
	CircuitID    string                         `json:"circuitId"`
	ScopeID      uint32                         `json:"scopeId"`
	Message      string                         `json:"message"`
	ZKProof      types.ZKProof                  `json:"zkProof"`
	CreatedTime  *int64                         `json:"createdTime"`
	ExpiresTime  *int64                         `json:"expiresTime"`
	Status       constant.ProofSubmissionStatus `json:"status"`
	VerifiedDate *time.Time                     `json:"verifiedDate"`
}

func ToAuthorizationResponse(ps *proof.ProofSubmission) protocol.AuthorizationResponseMessage {
	var zkProof types.ZKProof
	err := json.Unmarshal(ps.ZKProof, &zkProof)
	if err != nil {
		fmt.Println("zk proof unmarshal false %w", err)
	}

	return protocol.AuthorizationResponseMessage{
		ID:       ps.ThreadID,
		Typ:      packers.MediaTypePlainMessage,
		Type:     protocol.AuthorizationResponseMessageType,
		ThreadID: ps.ThreadID,
		From:     ps.HolderDID,
		To:       ps.ProofRequest.VerifierDID,
		Body: protocol.AuthorizationMessageResponseBody{
			Message: ps.Message,
			Scope: []protocol.ZeroKnowledgeProofResponse{
				{
					ID:        ps.ScopeID,
					CircuitID: ps.CircuitID,
					ZKProof:   zkProof,
				},
			},
		},
		CreatedTime: ps.CreatedTime,
		ExpiresTime: ps.ExpiresTime,
	}
}

func ToProofSubmissionResponseDto(entity *proof.ProofSubmission) *ProofSubmissionResponseDto {
	var zkProof types.ZKProof
	err := json.Unmarshal(entity.ZKProof, &zkProof)
	if err != nil {
		fmt.Println("zk proof unmarshal false %w", err)
	}
	return &ProofSubmissionResponseDto{
		PublicID:     entity.PublicID.String(),
		RequestID:    entity.ProofRequest.PublicID.String(),
		ThreadID:     entity.ThreadID,
		HolderDID:    entity.HolderDID,
		HolderName:   entity.Holder.Name,
		VerifierDID:  entity.ProofRequest.VerifierDID,
		VerifierName: entity.ProofRequest.Verifier.Name,
		Message:      entity.Message,
		ScopeID:      entity.ScopeID,
		CircuitID:    entity.CircuitID,
		ZKProof:      zkProof,
		CreatedTime:  entity.CreatedTime,
		ExpiresTime:  entity.ExpiresTime,
		VerifiedDate: entity.VerifiedDate,
		Status:       entity.Status,
	}
}
