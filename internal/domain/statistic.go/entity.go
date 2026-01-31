package statistic

import (
	"be/internal/domain/schema"
	"time"
)

type IssuerStatistic struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	IssuerDID string `gorm:"column:issuer_did;type:varchar(255);not null;uniqueIndex:idx_issuer_did" json:"issuer_did" validate:"required,startswith=did:"`

	DocumentNum          int64 `gorm:"default:0" json:"document_num"`
	SchemaNum            int64 `gorm:"default:0" json:"schema_num"`
	CredentialRequestNum int64 `gorm:"default:0" json:"credential_request_num"`
	CredentialIssuedNum  int64 `gorm:"default:0" json:"credential_issued_num"`

	Issuer *schema.Identity `gorm:"foreignKey:IssuerDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"issuer,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
}

type HolderStatistic struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	HolderDID string `gorm:"column:holder_did;type:varchar(255);not null;uniqueIndex:idx_holder_did" json:"holder_did" validate:"required,startswith=did:"`

	CredentialRequestNum    int64 `gorm:"default:0" json:"credential_request_num"`
	VerifiableCredentialNum int64 `gorm:"default:0" json:"verifiable_credential_num"`
	ProofSubmissionNum      int64 `gorm:"default:0" json:"proof_submission_num"`
	ProofAcceptedNum        int64 `gorm:"default:0" json:"proof_accepted_num"`

	Holder *schema.Identity `gorm:"foreignKey:HolderDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"holder,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
}

type VerifierStatistic struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	VerifierDID string `gorm:"column:verifier_did;type:varchar(255);not null;uniqueIndex:idx_verifier_did" json:"verifier_did" validate:"required,startswith=did:"`

	ProofRequestNum    int64 `gorm:"default:0" json:"proof_request_num"`
	ProofSubmissionNum int64 `gorm:"default:0" json:"proof_submission_num"`
	ProofAcceptedNum   int64 `gorm:"default:0" json:"proof_accepted_num"`

	Verifier *schema.Identity `gorm:"foreignKey:VerifierDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"verifier,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
}
