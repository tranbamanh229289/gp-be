package dto

import (
	"be/internal/domain/credential"
)

type IdentityCreatedRequestDto struct {
	PublicKeyX string
	PublicKeyY string
	Name       string
	Role       string
}

type IdentityResponseDto struct {
	PublicID   string
	PublicKeyX string
	PublicKeyY string
	Name       string
	Role       string
	DID        string
	State      string
}

func ToIdentityResponseDto(entity *credential.Identity) *IdentityResponseDto {
	return &IdentityResponseDto{
		PublicID:   entity.PublicID.String(),
		PublicKeyX: entity.PublicKeyX,
		PublicKeyY: entity.PublicKeyY,
		Role:       entity.Role,
		Name:       entity.Name,
		DID:        string(entity.DID),
		State:      string(entity.State),
	}
}
