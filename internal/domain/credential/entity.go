package credential

import (
	"be/internal/domain/schema"
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID   uuid.UUID  `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	PublicKeyX string     `gorm:"type:varchar(255)" json:"public_key_x,omitempty"`
	PublicKeyY string     `gorm:"type:varchar(255)" json:"public_key_y,omitempty"`
	Name       string     `gorm:"type:varchar(255)" json:"name,omitempty"`
	Role       string     `gorm:"type:varchar(100);not null" json:"role" validate:"required,oneof=holder issuer verifier"`
	DID        types.DID  `gorm:"type:varchar(255);uniqueIndex;not null" json:"did" validate:"required"`
	State      types.Hash `gorm:"type:varchar(255)" json:"state,omitempty"`
	ClaimsMTID uint64     `json:"claims_mt_id,omitempty"`
	RevMTID    uint64     `json:"rev_mt_id,omitempty"`
	RootsMTID  uint64     `json:"roots_mt_id,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`

	// Relationships
	IssuedCredentials  []Credential `gorm:"foreignKey:IssuerID" json:"issued_credentials,omitempty"`
	HeldCredentials    []Credential `gorm:"foreignKey:HolderID" json:"held_credentials,omitempty"`
	CredentialRequests []Credential `gorm:"foreignKey:HolderID" json:"credential_requests,omitempty"`
}
type Credential struct {
	ID              uint        `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID        uuid.UUID   `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	IssuerID        uint        `gorm:"not null;index" json:"issuer_id" validate:"required"`
	HolderID        uint        `gorm:"not null;index" json:"holder_id" validate:"required"`
	SchemaID        uint        `gorm:"index" json:"schema_id,omitempty"`
	ClaimData       types.JSONB `gorm:"type:jsonb" json:"claim_data,omitempty"`
	ClaimHi         types.Hash  `gorm:"type:varchar(255)" json:"claim_hi,omitempty"`
	ClaimHv         types.Hash  `gorm:"type:varchar(255)" json:"claim_hv,omitempty"`
	ClaimSubject    types.JSONB `gorm:"type:jsonb" json:"claim_subject,omitempty"`
	RevocationNonce uint64      `json:"revocation_nonce,omitempty"`
	ExpirationDate  *time.Time  `json:"expiration_date,omitempty"`
	ProofType       string      `gorm:"type:varchar(100)" json:"proof_type,omitempty"`
	IssuerState     types.Hash  `gorm:"type:varchar(255)" json:"issuer_state,omitempty"`
	SignatureProof  types.JSONB `gorm:"type:jsonb" json:"signature_proof,omitempty"`
	MTPProof        types.JSONB `gorm:"type:jsonb" json:"mtp_proof,omitempty"`
	Status          string      `gorm:"type:varchar(30);default:'Active'" json:"status" validate:"required,oneof=pending approved rejected"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RejectedAt      *time.Time  `gorm:"type:timestamptz" json:"rejected_at,omitempty" validate:"omitempty"`

	// Relationships
	Issuer *Identity      `gorm:"foreignKey:IssuerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"issuer,omitempty"`
	Holder *Identity      `gorm:"foreignKey:HolderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	Schema *schema.Schema `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
}

type CredentialRequest struct {
	ID             uint        `gorm:"primaryKey;autoIncrement" json:"id" validate:"omitempty"`
	PublicID       uuid.UUID   `gorm:"type:uuid;uniqueIndex;not null" json:"public_id" validate:"required,uuid4"`
	HolderID       uint        `gorm:"not null;index" json:"holder_id" validate:"required"`
	SchemaID       uint        `gorm:"not null" json:"schema_id" validate:"required,gt=0"`
	CredentialData types.JSONB `gorm:"type:jsonb;not null" json:"credential_data" validate:"required"`
	Status         string      `gorm:"type:varchar(50);default:'pending';index" json:"status" validate:"required,oneof=pending approved rejected"`
	CreatedAt      time.Time   `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RejectedAt     *time.Time  `gorm:"type:timestamptz" json:"rejected_at,omitempty" validate:"omitempty"`

	// Relationships
	Holder *Identity      `gorm:"foreignKey:HolderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	Schema *schema.Schema `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
}
