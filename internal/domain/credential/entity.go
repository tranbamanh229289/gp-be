package credential

import (
	"be/internal/domain/schema"
	"be/internal/shared/constant"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CredentialRequest struct {
	ID          uint                             `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID    uuid.UUID                        `gorm:"column:public_id;type:uuid;uniqueIndex;not null;default:gen_random_uuid()" json:"public_id" validate:"required"`
	ThreadID    string                           `gorm:"column:thread_id;type:varchar(255);not null;index" json:"thread_id" validate:"required"`
	HolderDID   string                           `gorm:"column:holder_did;type:varchar(255);not null;index" json:"holder_did" validate:"required,startswith=did:"`
	IssuerDID   string                           `gorm:"column:issuer_did;type:varchar(255);not null;index" json:"issuer_did" validate:"required,startswith=did:"`
	SchemaID    uint                             `gorm:"column:schema_id;not null;index" json:"schema_id" validate:"required,gt=0"`
	SchemaHash  string                           `gorm:"column:schema_hash;type:varchar(128);not null" json:"schema_hash" validate:"required"`
	Status      constant.CredentialRequestStatus `gorm:"column:status;type:varchar(20);not null;default:'pending'" json:"status" validate:"required"`
	Expiration  int64                            `gorm:"column:expiration;type:bigint" json:"expiration"`
	CreatedTime *int64                           `gorm:"column:created_time;type:bigint" json:"created_time,omitempty" validate:"omitempty"`
	ExpiresTime *int64                           `gorm:"column:expires_time;type:bigint" json:"expires_time,omitempty" validate:"omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`

	Holder               *schema.Identity      `gorm:"foreignKey:HolderDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	Issuer               *schema.Identity      `gorm:"foreignKey:IssuerDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"issuer,omitempty"`
	Schema               *schema.Schema        `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
	VerifiableCredential *VerifiableCredential `gorm:"foreignKey:CRID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"verifiable_credential,omitempty"`
}

type VerifiableCredential struct {
	ID                uint                                `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID          uuid.UUID                           `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	CRID              uint                                `gorm:"column:crid;index;not null" json:"crid" validate:"required,gt=0"`
	HolderDID         string                              `gorm:"column:holder_did;type:varchar(255);index;not null" json:"holder_did" validate:"required,startswith=did:"`
	IssuerDID         string                              `gorm:"column:issuer_did;type:varchar(255);index;not null" json:"issuer_did" validate:"required,startswith=did:"`
	SchemaID          uint                                `gorm:"column:schema_id;index;not null" json:"schema_id" validate:"required,gt=0"`
	SchemaHash        string                              `gorm:"column:schema_hash;type:varchar(128);not null" json:"schema_hash" validate:"required"`
	CredentialID      string                              `gorm:"column:credential_id;type:varchar(255)" json:"credential_id,omitempty"`
	CredentialSubject datatypes.JSONMap                   `gorm:"column:credential_subject;not null" json:"credential_subject" validate:"required"`
	ClaimSubject      string                              `gorm:"column:claim_subject;type:varchar(255);not null" json:"claim_subject" validate:"required"`
	ClaimHi           string                              `gorm:"column:claim_hi;type:varchar(255);not null" json:"claim_hi" validate:"required,len=255"`
	ClaimHv           string                              `gorm:"column:claim_hv;type:varchar(255);not null" json:"claim_hv" validate:"required,len=255"`
	ClaimHex          string                              `gorm:"column:claim_hex;type:text;not null" json:"claim_hex" validate:"required"`
	ClaimMTP          []byte                              `gorm:"column:claim_mtp;type:bytea;not null" json:"claim_mtp" validate:"required"`
	RevNonce          uint64                              `gorm:"column:rev_nonce;type:bigint;not null" json:"rev_nonce" validate:"required"`
	AuthClaimHex      string                              `gorm:"column:auth_claim_hex;type:text;not null" json:"auth_claim_hex" validate:"required"`
	AuthClaimMTP      []byte                              `gorm:"column:auth_claim_mtp;type:bytea;not null" json:"auth_claim_mtp" validate:"required"`
	Signature         string                              `gorm:"column:signature;type:varchar(66);not null" json:"signature" validate:"required"`
	IssuerState       string                              `gorm:"column:issuer_state;type:varchar(66);not null" json:"issuer_state" validate:"required"`
	ClaimsTreeRoot    string                              `gorm:"column:claims_tree_root;type:char(66);not null" json:"claims_tree_root" validate:"required"`
	RevTreeRoot       string                              `gorm:"column:rev_tree_root;type:char(66);not null" json:"rev_tree_root" validate:"required"`
	RootsTreeRoot     string                              `gorm:"column:roots_tree_root;type:char(66);not null" json:"roots_tree_root" validate:"required"`
	IssuanceDate      *time.Time                          `gorm:"column:issuance_date;type:timestamptz" json:"issuance_date,omitempty" validate:"omitempty"`
	ExpirationDate    *time.Time                          `gorm:"column:expiration_date;type:timestamptz" json:"expiration_date,omitempty" validate:"omitempty"`
	Status            constant.VerifiableCredentialStatus `gorm:"column:status;type:varchar(30);default:'issued'" json:"status" validate:"required"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RevokedAt *time.Time `gorm:"type:timestamptz" json:"revoked_at,omitempty" validate:"omitempty"`

	Issuer            *schema.Identity   `gorm:"foreignKey:IssuerDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"issuer,omitempty"`
	Holder            *schema.Identity   `gorm:"foreignKey:HolderDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`
	CredentialRequest *CredentialRequest `gorm:"foreignKey:CRID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"credential_request,omitempty"`
	Schema            *schema.Schema     `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
}
