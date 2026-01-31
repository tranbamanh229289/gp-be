package dto

import "math/big"

type CredentialAtomicQueryV3InputRequestDto struct {
	ScopeID        big.Int `json:"scopeId"`
	ProofRequestID string  `json:"proofRequestId"`
	CredentialID   string  `json:"credentialId"`
}
