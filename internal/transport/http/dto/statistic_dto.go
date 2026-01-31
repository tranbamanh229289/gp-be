package dto

type IssuerStatisticResponseDto struct {
	DocumentNum          int64 `json:"documentNum"`
	SchemaNum            int64 `json:"schemaNum"`
	CredentialRequestNum int64 `json:"credentialRequestNum"`
	CredentialIssuedNum  int64 `json:"credentialIssuedNum"`
}

type HolderStatisticResponseDto struct {
	CredentialRequestNum    int64 `json:"credentialRequestNum"`
	VerifiableCredentialNum int64 `json:"verifiableCredentialNum"`
	ProofSubmissionNum      int64 `json:"proofSubmissionNum"`
	ProofAcceptedNum        int64 `json:"proofAcceptedNum"`
}

type VerifierStatisticResponseDto struct {
	ProofRequestNum  int64 `json:"proofRequestNum"`
	ProofSubmission  int64 `json:"proofSubmissionNum"`
	ProofAcceptedNum int64 `json:"proofAcceptedNum"`
}
