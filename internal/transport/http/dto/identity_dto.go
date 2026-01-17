package dto

import (
	"be/internal/domain/schema"
	"be/internal/shared/constant"
)

type IdentityCreatedRequestDto struct {
	PublicKeyX string                `json:"publicKeyX"`
	PublicKeyY string                `json:"publicKeyY"`
	Name       string                `json:"name"`
	Role       constant.IdentityRole `json:"role"`
}

type IdentityResponseDto struct {
	PublicID   string                `json:"id"`
	PublicKeyX string                `json:"publicKeyX"`
	PublicKeyY string                `json:"publicKeyY"`
	Name       string                `json:"name"`
	Role       constant.IdentityRole `json:"role"`
	DID        string                `json:"did"`
	State      string                `json:"state"`
}

func ToIdentityResponseDto(entity *schema.Identity) *IdentityResponseDto {
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
