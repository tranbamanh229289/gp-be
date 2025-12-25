package dto

import (
	"be/internal/domain/credential"
)

type IdentityCreatedRequestDto struct {
	PublicKeyX string `json:"publicKeyX"`
	PublicKeyY string `json:"publicKeyY"`
	Name       string `json:"name"`
	Role       string `json:"role"`
}

type IdentityResponseDto struct {
	PublicID   string `json:"publicID"`
	PublicKeyX string `json:"publicKeyX"`
	PublicKeyY string `json:"publicKeyY"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	DID        string `json:"did"`
	State      string `json:"state"`
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
