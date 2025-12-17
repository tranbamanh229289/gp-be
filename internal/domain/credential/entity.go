package credential

import (
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID          uuid.UUID `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	DID               types.DID `gorm:"type:varchar(255);uniqueIndex;not null" json:"did" validate:"required"`
	Type              string
	State             types.Hash
	PublicKeyX        string    `json:"public_key_x"`
	PublicKeyY        string    `json:"public_key_y"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdateAt          time.Time
	IssuedCredentials []Credential `gorm:"foreignKey:IssuerID" json:"issued_credentials,omitempty"`
	HeldCredentials   []Credential `gorm:"foreignKey:HolderID" json:"held_credentials,omitempty"`
}

type Schema struct {
	ID         uint
	PublicID   uuid.UUID
	SchemaHash string
	SchemaType string
	SchemaURL  string
	SchemaJSON types.JSONB
	IPFSCID    string
	CreateAt   time.Time
}

type Credential struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID        uuid.UUID `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	IssuerID        uint      `gorm:"not null;index" json:"issuer_id"`
	HolderID        uint      `gorm:"not null;index" json:"holder_id"`
	SchemaID        uint
	ClaimData       types.JSONB
	ClaimHi         types.Hash
	ClaimHv         types.Hash
	ClaimSubject    types.JSONB
	RevocationNonce uint64
	ExpirationDate  *time.Time
	ProofType       string
	IssuerState     types.Hash
	SignatureProof  types.JSONB
	MTPProof        types.JSONB
	Status          string     `gorm:"type:varchar(30);default:'Active'" json:"status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	RevokedAt       *time.Time `json:"revoked_at,omitempty"`
	Issuer          *Identity  `gorm:"foreignKey:IssuerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"issuer,omitempty"`
	Holder          *Identity  `gorm:"foreignKey:HolderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"holder,omitempty"`
	Schema          *Schema
}
