package dto

import (
	"be/internal/domain/credential"
	"time"

	"github.com/iden3/go-schema-processor/v2/verifiable"
)

type CredentialRequestUpdatedRequestDto struct {
	Status string `json:"status"`
}

type CredentialRequestResponseDto struct {
	PublicID    string `json:"id"`
	RequestID   string `json:"requestId"`
	HolderDID   string `json:"holderDID"`
	HolderName  string `json:"holderName"`
	IssuerDID   string `json:"issuerDID"`
	IssuerName  string `json:"issuerName"`
	SchemaURL   string `json:"schemaURL"`
	SchemaType  string `json:"schemaType"`
	Expiration  int64  `json:"expiration"`
	CreatedTime *int64 `json:"createdTime"`
	ExpiresTime *int64 `json:"expiresTime"`
	Status      string `json:"status"`
}

func ToCredentialRequestResponseDto(credentialRequest *credential.CredentialRequest) *CredentialRequestResponseDto {

	return &CredentialRequestResponseDto{
		PublicID:    credentialRequest.PublicID.String(),
		RequestID:   credentialRequest.RequestID,
		HolderDID:   credentialRequest.HolderDID,
		IssuerDID:   credentialRequest.IssuerDID,
		Status:      credentialRequest.Status,
		Expiration:  credentialRequest.Expiration,
		CreatedTime: credentialRequest.CreatedTime,
		ExpiresTime: credentialRequest.ExpiresTime,
		SchemaURL:   credentialRequest.Schema.ContextURL,
		SchemaType:  credentialRequest.Schema.Type,
	}
}

type IssueVerifiableCredentialRequestDto struct {
	ProofType verifiable.ProofType `json:"proofType"`
}
type VerifiableUpdatedRequestDto struct {
	Status string `json:"status"`
}

type SignCredentialRequestDto struct {
	Signature string `json:"signature"`
}

type VerifiableCredentialResponseDto struct {
	PublicID       string     `json:"id"`
	HolderDID      string     `json:"holderDID"`
	IssuerDID      string     `json:"issuerDID"`
	SchemaURL      string     `json:"schemaURL"`
	SchemaType     string     `json:"schemaType"`
	ClaimSubject   string     `json:"claimSubject"`
	ClaimHi        string     `json:"claimHi"`
	ClaimHv        string     `json:"claimHv"`
	RevNonce       uint64     `json:"revNonce"`
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
	IssuanceDate   *time.Time `json:"issuanceDate,omitempty"`
	IssuerState    string     `json:"issuerState"`
	ProofType      string     `json:"proofType"`
	Status         string     `json:"status"`
	Signature      string     `json:"signature"`
}

func ToVerifiableCredentialResponseDto(vc *credential.VerifiableCredential) *VerifiableCredentialResponseDto {
	return &VerifiableCredentialResponseDto{
		PublicID:       vc.PublicID.String(),
		HolderDID:      vc.HolderDID,
		IssuerDID:      vc.IssuerDID,
		SchemaURL:      vc.Schema.ContextURL,
		SchemaType:     vc.Schema.Type,
		ClaimSubject:   vc.ClaimSubject,
		ClaimHi:        vc.ClaimHi,
		ClaimHv:        vc.ClaimHv,
		RevNonce:       vc.RevNonce,
		ExpirationDate: vc.ExpirationDate,
		IssuanceDate:   vc.IssuanceDate,
		IssuerState:    vc.IssuerState,
		ProofType:      vc.ProofType,
		Status:         vc.Status,
		Signature:      vc.Signature,
	}
}
