package credential

import (
	"be/internal/domain/schema"
	"time"

	"gorm.io/datatypes"

	"github.com/google/uuid"
)

type VerifiableCredential struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID       uuid.UUID  `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	HolderDID      string     `gorm:"column:holder_did;type:varchar(255);index;not null" json:"holder_did" validate:"required,startswith=did:"`
	IssuerDID      string     `gorm:"column:issuer_did;type:varchar(255);index;not null" json:"issuer_did" validate:"required,startswith=did:"`
	SchemaID       uint       `gorm:"column:schema_id;index;not null" json:"schema_id" validate:"required,gt=0"`
	ClaimSubject   string     `gorm:"column:claim_subject;type:varchar(255)" json:"claim_subject,omitempty" validate:"omitempty,startswith=did:"`
	ClaimHi        string     `gorm:"column:claim_hi;type:varchar(255);not null" json:"claim_hi" validate:"required,len=64"`
	ClaimHv        string     `gorm:"column:claim_hv;type:varchar(255);not null" json:"claim_hv" validate:"required,len=64"`
	IssuerState    string     `gorm:"column:issuer_state;type:varchar(255);not null" json:"issuer_state" validate:"required,len=64"`
	RevNonce       uint64     `gorm:"column:rev_nonce;type:bigint;not null" json:"rev_nonce" validate:"required"`
	IssuanceDate   *time.Time `gorm:"column:issuance_date;type:timestamptz" json:"issuance_date,omitempty" validate:"omitempty"`
	ExpirationDate *time.Time `gorm:"column:expiration_date;type:timestamptz" json:"expiration_date,omitempty" validate:"omitempty"`
	ProofType      string     `gorm:"column:proof_type;type:varchar(100)" json:"proof_type,omitempty" validate:"omitempty,oneof=BjjSignature2021 Iden3SparseMerkleTreeProof"`
	Status         string     `gorm:"column:status;type:varchar(30);default:'issued'" json:"status" validate:"required,oneof=notSigned issued revoked expired"`
	Signature      string     `gorm:"column:signature;type:varchar(255)" json:"signature,omitempty" validate:"omitempty"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RevokedAt *time.Time `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`

	Issuer *schema.Identity `gorm:"foreignKey:IssuerDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"issuer,omitempty"`
	Holder *schema.Identity `gorm:"foreignKey:HolderDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	Schema *schema.Schema   `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
}

type CredentialRequest struct {
	ID          uint              `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID    uuid.UUID         `gorm:"column:public_id;type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"public_id" validate:"required"`
	RequestID   string            `gorm:"column:request_id;type:varchar(128);uniqueIndex;not null" json:"request_id" validate:"required"`
	HolderDID   string            `gorm:"column:holder_did;type:varchar(255);not null;index" json:"holder_did" validate:"required,startswith=did:"`
	IssuerDID   string            `gorm:"column:issuer_did;type:varchar(255);not null;index" json:"issuer_did" validate:"required,startswith=did:"`
	SchemaID    uint              `gorm:"column:schema_id;not null;index" json:"schema_id" validate:"required,gt=0"`
	SchemaHash  string            `gorm:"column:schema_hash;type:varchar(128);index;not null" json:"schema_hash" validate:"required"`
	Data        datatypes.JSONMap `gorm:"column:data;type:jsonb;not null;default:'{}'::jsonb" json:"data" validate:"required"`
	Status      string            `gorm:"column:status;type:varchar(20);not null;default:'pending'" json:"status" validate:"required,oneof=pending approved rejected"`
	Expiration  int64             `gorm:"column:expiration;type:bigint" json:"expiration"`
	CreatedTime *int64            `gorm:"column:created_time;type:bigint" json:"created_time,omitempty" validate:"omitempty"`
	ExpiresTime *int64            `gorm:"column:expires_time;type:bigint" json:"expires_time,omitempty" validate:"omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`

	Holder *schema.Identity `gorm:"foreignKey:HolderDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	Issuer *schema.Identity `gorm:"foreignKey:IssuerDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"issuer,omitempty"`
	Schema *schema.Schema   `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
}
