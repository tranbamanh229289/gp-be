package dto

import (
	"be/internal/shared/constant"

	"github.com/golang-jwt/jwt/v5"
)

type ZKClaims struct {
	ID    string                `json:"id"`
	Name  string                `json:"name"`
	DID   string                `json:"did"`
	State string                `json:"state"`
	Role  constant.IdentityRole `json:"role"`
	jwt.RegisteredClaims
}

type ZKLoginResponseDto struct {
	Claims       ZKClaims     `json:"claims"`
	PublicKey    PublicKeyDto `json:"publicKey"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}

type PublicKeyDto struct {
	X string `json:"x"`
	Y string `json:"y"`
}
type RefreshTokenResponseDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
