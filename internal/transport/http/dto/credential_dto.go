package dto

type CredentialRequestCreatedRequestDto struct {
	HolderID       string
	SchemaID       string
	CredentialData map[string]interface{}
}

type CredentialRequestResponseDto struct {
	PublicID       string
	HolderID       string
	SchemaID       string
	CredentialData map[string]interface{}
	Status         string
}

type CredentialIssuedCreatedRequestDto struct {
	RequestID string
}
