package dto

import (
	"be/internal/domain/credential"
	"time"
)

type IdentityCreatedRequestDto struct {
	PublicKeyX string
	PublicKeyY string
	Type       string
}

type IdentityCreatedResponseDto struct {
	PublicID   string
	DID        string
	Type       string
	State      string
	PublicKeyX string
	PublicKeyY string
}

func ToIdentityCreatedResponseDto(entity *credential.Identity) *IdentityCreatedResponseDto {
	return &IdentityCreatedResponseDto{
		PublicID:   entity.PublicID.String(),
		DID:        string(entity.DID),
		Type:       entity.Type,
		State:      string(entity.State),
		PublicKeyX: entity.PublicKeyX,
		PublicKeyY: entity.PublicKeyY,
	}
}

type CredentialRequestCreatedRequestDto struct {
	IssuerID    string
	HolderID    string
	SchemaID    string
	ProofType   string
	Expirations time.Time
}

type CredentialRequestCreatedResponseDto struct {
}
