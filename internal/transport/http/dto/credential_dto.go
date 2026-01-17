package dto

import (
	"be/internal/domain/credential"
	"be/internal/shared/constant"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-schema-processor/v2/verifiable"
)

type CredentialRequestUpdatedRequestDto struct {
	Status string `json:"status"`
}

type CredentialRequestResponseDto struct {
	PublicID     string                           `json:"id"`
	ThreadID     string                           `json:"threadId"`
	HolderDID    string                           `json:"holderDID"`
	HolderName   string                           `json:"holderName"`
	IssuerDID    string                           `json:"issuerDID"`
	IssuerName   string                           `json:"issuerName"`
	SchemaID     string                           `json:"schemaId"`
	SchemaTitle  string                           `json:"schemaTitle"`
	SchemaURL    string                           `json:"schemaURL"`
	ContextURL   string                           `json:"contextURL"`
	SchemaType   string                           `json:"schemaType"`
	SchemaHash   string                           `json:"schemaHash"`
	DocumentType constant.DocumentType            `json:"documentType"`
	IsMerklized  bool                             `json:"isMerklized"`
	Status       constant.CredentialRequestStatus `json:"status"`
	Expiration   int64                            `json:"expiration"`
	CreatedTime  *int64                           `json:"createdTime"`
	ExpiresTime  *int64                           `json:"expiresTime"`
}

func ToCredentialRequestResponseDto(credentialRequest *credential.CredentialRequest) *CredentialRequestResponseDto {
	return &CredentialRequestResponseDto{
		PublicID:     credentialRequest.PublicID.String(),
		ThreadID:     credentialRequest.ThreadID,
		HolderDID:    credentialRequest.HolderDID,
		HolderName:   credentialRequest.Holder.Name,
		IssuerDID:    credentialRequest.IssuerDID,
		IssuerName:   credentialRequest.Issuer.Name,
		SchemaID:     credentialRequest.Schema.PublicID.String(),
		SchemaTitle:  credentialRequest.Schema.Title,
		SchemaURL:    credentialRequest.Schema.SchemaURL,
		ContextURL:   credentialRequest.Schema.ContextURL,
		SchemaType:   credentialRequest.Schema.Type,
		SchemaHash:   credentialRequest.Schema.Hash,
		DocumentType: credentialRequest.Schema.DocumentType,
		IsMerklized:  credentialRequest.Schema.IsMerklized,
		Status:       credentialRequest.Status,
		Expiration:   credentialRequest.Expiration,
		CreatedTime:  credentialRequest.CreatedTime,
		ExpiresTime:  credentialRequest.ExpiresTime,
	}
}

type IssueVerifiableCredentialRequestDto struct {
	IsMerklized       bool                        `json:"isMerklized"`
	CredentialStatus  verifiable.CredentialStatus `json:"credentialStatus"`
	CredentialSubject map[string]interface{}      `json:"credentialSubject"`
	Signature         string                      `json:"signature"`
}

type VerifiableUpdatedRequestDto struct {
	Status string `json:"status"`
}

type VerifiableCredentialResponseDto struct {
	PublicID          string                              `json:"id"`
	CredentialID      string                              `json:"credentialId"`
	HolderDID         string                              `json:"holderDID"`
	HolderName        string                              `json:"holderName"`
	IssuerDID         string                              `json:"issuerDID"`
	IssuerName        string                              `json:"issuerName"`
	SchemaID          string                              `json:"schemaId"`
	SchemaURL         string                              `json:"schemaURL"`
	ContextURL        string                              `json:"contextURL"`
	SchemaType        string                              `json:"schemaType"`
	DocumentType      constant.DocumentType               `json:"documentType"`
	CredentialSubject map[string]interface{}              `json:"credentialSubject"`
	ClaimSubject      string                              `json:"claimSubject"`
	ClaimHi           string                              `json:"claimHi"`
	ClaimHv           string                              `json:"claimHv"`
	ClaimHex          string                              `json:"claimHex"`
	ClaimMTP          merkletree.Proof                    `json:"claimMTP"`
	RevNonce          uint64                              `json:"revNonce"`
	AuthClaimHex      string                              `json:"authClaimHex"`
	AuthClaimMTP      merkletree.Proof                    `json:"authClaimMTP"`
	Signature         SignatureDto                        `json:"signature"`
	IssuerState       string                              `json:"issuerState"`
	ClaimsTreeRoot    string                              `json:"claimsTreeRoot"`
	RevTreeRoot       string                              `json:"revTreeRoot"`
	RootsTreeRoot     string                              `json:"rootsTreeRoot"`
	IssuanceDate      *time.Time                          `json:"issuanceDate,omitempty"`
	ExpirationDate    *time.Time                          `json:"expirationDate,omitempty"`
	Status            constant.VerifiableCredentialStatus `json:"status"`
}

func ToVerifiableCredentialResponseDto(vc *credential.VerifiableCredential) *VerifiableCredentialResponseDto {
	var claimProof merkletree.Proof
	var authProof merkletree.Proof
	err := claimProof.UnmarshalJSON(vc.ClaimMTP)
	if err != nil {
		fmt.Println("claim proof unmarshal false %w", err)
	}
	err = authProof.UnmarshalJSON(vc.AuthClaimMTP)
	if err != nil {
		fmt.Println("auth claim proof unmarshal false %w", err)
	}
	signature, err := DecodeSignatureString(vc.Signature)
	if err != nil {
		fmt.Println("decode error %w", err)
	}

	return &VerifiableCredentialResponseDto{
		PublicID:          vc.PublicID.String(),
		CredentialID:      vc.CredentialID,
		HolderDID:         vc.HolderDID,
		HolderName:        vc.Holder.Name,
		IssuerDID:         vc.IssuerDID,
		IssuerName:        vc.Issuer.Name,
		SchemaID:          vc.Schema.PublicID.String(),
		SchemaURL:         vc.Schema.SchemaURL,
		ContextURL:        vc.Schema.ContextURL,
		SchemaType:        vc.Schema.Type,
		DocumentType:      vc.Schema.DocumentType,
		CredentialSubject: vc.CredentialSubject,
		ClaimSubject:      vc.ClaimSubject,
		ClaimHi:           vc.ClaimHi,
		ClaimHv:           vc.ClaimHv,
		ClaimHex:          vc.ClaimHex,
		ClaimMTP:          claimProof,
		RevNonce:          vc.RevNonce,
		AuthClaimHex:      vc.AuthClaimHex,
		AuthClaimMTP:      authProof,
		IssuerState:       vc.IssuerState,
		ClaimsTreeRoot:    vc.ClaimsTreeRoot,
		RevTreeRoot:       vc.RevTreeRoot,
		RootsTreeRoot:     vc.RootsTreeRoot,
		ExpirationDate:    vc.ExpirationDate,
		IssuanceDate:      vc.IssuanceDate,
		Status:            vc.Status,
		Signature: SignatureDto{
			SignatureS:   signature.S.String(),
			SignatureR8X: signature.R8.X.String(),
			SignatureR8Y: signature.R8.Y.String(),
		},
	}
}

type SignatureDto struct {
	SignatureS   string `json:"signatureR8S"`
	SignatureR8X string `json:"signatureR8X"`
	SignatureR8Y string `json:"signatureR8Y"`
}

func DecodeSignatureString(sigString string) (*babyjub.Signature, error) {
	signCompBytes, err := hex.DecodeString(sigString)
	if err != nil {
		return nil, err
	}
	var compSig babyjub.SignatureComp
	copy(compSig[:], signCompBytes)
	signature, err := compSig.Decompress()
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func ToW3CCredential(vc *credential.VerifiableCredential) *verifiable.W3CCredential {
	var authIncProof merkletree.Proof
	var incProof merkletree.Proof

	err := authIncProof.UnmarshalJSON(vc.AuthClaimMTP)
	if err != nil {
		fmt.Println("Unmarshal auth claim proof error: %w", err)
	}
	err = incProof.UnmarshalJSON(vc.ClaimMTP)
	if err != nil {
		fmt.Println("Unmarshal core claim proof error: %w", err)
	}

	iden3SparseMerkleProof := &verifiable.Iden3SparseMerkleTreeProof{
		Type: verifiable.Iden3SparseMerkleTreeProofType,
		IssuerData: verifiable.IssuerData{
			ID: vc.IssuerDID,
			State: verifiable.State{
				Value:              &vc.IssuerState,
				ClaimsTreeRoot:     &vc.ClaimsTreeRoot,
				RevocationTreeRoot: &vc.RevTreeRoot,
				RootOfRoots:        &vc.RootsTreeRoot,
				Status:             string(vc.Status),
			},
			AuthCoreClaim:    vc.AuthClaimHex,
			MTP:              &authIncProof,
			CredentialStatus: string(vc.Status),
		},
		CoreClaim: vc.ClaimHex,
		MTP:       &incProof,
	}
	bjjSignatureProof := &verifiable.BJJSignatureProof2021{
		Type: verifiable.BJJSignatureProofType,
		IssuerData: verifiable.IssuerData{
			ID: vc.IssuerDID,
			State: verifiable.State{
				Value:              &vc.IssuerState,
				ClaimsTreeRoot:     &vc.ClaimsTreeRoot,
				RevocationTreeRoot: &vc.RevTreeRoot,
				RootOfRoots:        &vc.RootsTreeRoot,
				Status:             string(constant.VerifiableCredentialIssuedStatus),
			},
			AuthCoreClaim:    vc.AuthClaimHex,
			MTP:              &authIncProof,
			CredentialStatus: string(vc.Status),
		},
		CoreClaim: vc.ClaimHex,
		Signature: vc.Signature}

	return &verifiable.W3CCredential{
		ID: vc.CredentialID,
		Context: []string{
			verifiable.JSONLDSchemaW3CCredential2018,
			verifiable.JSONLDSchemaIden3Credential,
			vc.Schema.ContextURL,
		},
		Type: []string{
			verifiable.TypeW3CVerifiableCredential,
			vc.Schema.Type,
		},
		IssuanceDate: vc.IssuanceDate,
		Expiration:   vc.ExpirationDate,

		Issuer: vc.IssuerDID,
		CredentialSchema: verifiable.CredentialSchema{
			ID:   vc.Schema.SchemaURL,
			Type: verifiable.JSONSchema2023,
		},
		CredentialSubject: vc.CredentialSubject,
		Proof:             []verifiable.CredentialProof{iden3SparseMerkleProof, bjjSignatureProof},
	}
}
