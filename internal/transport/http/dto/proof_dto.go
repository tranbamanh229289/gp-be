package dto

import (
	"time"
)

type ProofRequestCreatedRequestDto struct {
	CircuitType       string
	VerifierID        string
	SchemaID          string
	Query             map[string]interface{}
	AllowedIssuersDID []string
	Challenge         string
	Status            string
	ExpireAt          time.Time
}

type ProofRequestCreatedResponseDto struct {
	PublicID          string
	VerifierID        string
	SchemaID          string
	CircuitType       string
	AllowedIssuersDID []string
	Challenge         string
	Status            string
	QRCodeData        string
	CallbackURL       string
	ExpireAt          time.Time
}
