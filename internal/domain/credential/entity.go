package credential

import (
	"be/internal/domain/schema"
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID   uuid.UUID `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	PublicKeyX string    `gorm:"column:public_key_x;type:varchar(255)" json:"public_key_x,omitempty" validate:"omitempty,max=255"`
	PublicKeyY string    `gorm:"column:public_key_y;type:varchar(255)" json:"public_key_y,omitempty" validate:"omitempty,max=255"`
	Name       string    `gorm:"column:name;type:varchar(255)" json:"name,omitempty" validate:"omitempty,max=255"`
	Role       string    `gorm:"column:role;type:varchar(100);not null;index" json:"role" validate:"required,oneof=holder issuer verifier"`
	DID        string    `gorm:"column:did;type:varchar(255);uniqueIndex;not null" json:"did" validate:"required,startswith=did:"`
	State      string    `gorm:"column:state;type:varchar(255)" json:"state,omitempty" validate:"omitempty,len=64"`
	ClaimsMTID uint64    `gorm:"column:claims_mt_id;index" json:"claims_mt_id,omitempty" validate:"omitempty"`
	RevMTID    uint64    `gorm:"column:rev_mt_id;index" json:"rev_mt_id,omitempty" validate:"omitempty"`
	RootsMTID  uint64    `gorm:"column:roots_mt_id;index" json:"roots_mt_id,omitempty" validate:"omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`

	IssuedCredentials  []VerifiableCredential `gorm:"foreignKey:IssuerID" json:"issued_credentials,omitempty"`
	HeldCredentials    []VerifiableCredential `gorm:"foreignKey:HolderID" json:"held_credentials,omitempty"`
	CredentialRequests []CredentialRequest    `gorm:"foreignKey:HolderID" json:"credential_requests,omitempty"`
}

type VerifiableCredential struct {
	ID              uint        `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID        uuid.UUID   `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	IssuerID        uint        `gorm:"column:issuer_id;not null;index" json:"issuer_id" validate:"required,gt=0"`
	HolderID        uint        `gorm:"column:holder_id;not null;index" json:"holder_id" validate:"required,gt=0"`
	SchemaID        uint        `gorm:"column:schema_id;index" json:"schema_id,omitempty" validate:"omitempty,gt=0"`
	ClaimData       types.JSONB `gorm:"column:claim_data;type:jsonb" json:"claim_data,omitempty" validate:"omitempty"`
	ClaimHi         string      `gorm:"column:claim_hi;type:varchar(255)" json:"claim_hi,omitempty" validate:"omitempty,len=64"`
	ClaimHv         string      `gorm:"column:claim_hv;type:varchar(255)" json:"claim_hv,omitempty" validate:"omitempty,len=64"`
	ClaimSubject    types.JSONB `gorm:"column:claim_subject;type:jsonb" json:"claim_subject,omitempty" validate:"omitempty"`
	RevocationNonce uint64      `gorm:"column:revocation_nonce" json:"revocation_nonce,omitempty" validate:"omitempty"`
	ExpirationDate  *time.Time  `gorm:"column:expiration_date;type:date" json:"expiration_date,omitempty" validate:"omitempty"`
	ProofType       string      `gorm:"column:proof_type;type:varchar(100)" json:"proof_type,omitempty" validate:"omitempty,oneof=BjjSignature2021 Iden3SparseMerkleTreeProof"`
	IssuerState     string      `gorm:"column:issuer_proof;type:varchar(255)" json:"issuer_state,omitempty" validate:"omitempty,len=64"`
	SignatureProof  types.JSONB `gorm:"column:signature_proof;type:jsonb" json:"signature_proof,omitempty" validate:"omitempty"`
	MTPProof        types.JSONB `gorm:"column:mtp_proof;type:jsonb" json:"mtp_proof,omitempty" validate:"omitempty"`
	Status          string      `gorm:"column:status;type:varchar(30);default:'pending';index" json:"status" validate:"required,oneof=pending issued revoked expired"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RevokedAt       *time.Time  `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`

	Issuer *Identity      `gorm:"foreignKey:IssuerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"issuer,omitempty"`
	Holder *Identity      `gorm:"foreignKey:HolderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	Schema *schema.Schema `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"schema,omitempty"`
}

type CredentialRequest struct {
	ID             uint        `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID   `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	HolderID       uint        `gorm:"column:public_id;not null;index" json:"holder_id" validate:"required,gt=0"`
	SchemaID       uint        `gorm:"column:schema_id;not null;index" json:"schema_id" validate:"required,gt=0"`
	CredentialData types.JSONB `gorm:"column:credential_data;type:jsonb;not null" json:"credential_data" validate:"required"`
	Message        string      `gorm:"column:message;type:text" json:"message,omitempty" validate:"omitempty,max=1000"`
	Status         string      `gorm:"column:status;type:varchar(30);default:'pending';index" json:"status" validate:"required,oneof=pending approved rejected"`
	CreatedAt      time.Time   `gorm:"autoCreateTime" json:"created_at" validate:"-"`

	Holder *Identity      `gorm:"foreignKey:HolderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	Schema *schema.Schema `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
}
