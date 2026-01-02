package schema

import (
	"time"

	"gorm.io/datatypes"

	"github.com/google/uuid"
)

type Identity struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID   uuid.UUID `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	PublicKeyX string    `gorm:"column:public_key_x;type:varchar(255);not null" json:"public_key_x" validate:"required,max=255"`
	PublicKeyY string    `gorm:"column:public_key_y;type:varchar(255);not null" json:"public_key_y" validate:"required,max=255"`
	Name       string    `gorm:"column:name;type:varchar(255)" json:"name,omitempty" validate:"omitempty,max=255"`
	Role       string    `gorm:"column:role;type:varchar(100);not null" json:"role" validate:"required,oneof=holder issuer verifier"`
	DID        string    `gorm:"column:did;type:varchar(255);uniqueIndex;not null" json:"did" validate:"required,startswith=did:"`
	State      string    `gorm:"column:state;type:varchar(255);not null" json:"state" validate:"required,len=64"`
	ClaimsMTID uint64    `gorm:"column:claims_mt_id;index;not null" json:"claims_mt_id" validate:"required"`
	RevMTID    uint64    `gorm:"column:rev_mt_id;index;not null" json:"rev_mt_id" validate:"required"`
	RootsMTID  uint64    `gorm:"column:roots_mt_id;index;not null" json:"roots_mt_id" validate:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
}

type Schema struct {
	ID            uint              `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID      uuid.UUID         `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	IssuerDID     string            `gorm:"column:issuer_did;type:varchar(255);index;not null" json:"issuer_did" validate:"required,startswith=did:"`
	Hash          string            `gorm:"column:hash;type:varchar(128);index;not null" json:"hash" validate:"required"`
	Type          string            `gorm:"column:type;type:varchar(255);not null" json:"type" validate:"required,min=3,max=255"`
	Version       string            `gorm:"column:version;type:varchar(64);not null" json:"version" validate:"required,semver"`
	Title         string            `gorm:"column:title;type:varchar(255);not null" json:"title" validate:"required,max=255"`
	Description   string            `gorm:"column:description;type:text; not null" json:"description" validate:"required,max=2000"`
	IsMerklized   bool              `gorm:"column:is_merklized;default:false" json:"is_merklized"`
	JSONSchema    datatypes.JSONMap `gorm:"column:json_schema;type:jsonb;not null" json:"json_schema" validate:"required"`
	JSONLDContext datatypes.JSONMap `gorm:"column:jsonld_context;type:jsonb;not null" json:"jsonld_context" validate:"required"`
	SchemaURL     string            `gorm:"column:schema_url;type:varchar(255);uniqueIndex" json:"schema_url" validate:"omitempty"`
	ContextURL    string            `gorm:"column:context_url;type:varchar(255);uniqueIndex" json:"context_url" validate:"omitempty"`
	Status        string            `gorm:"column:status;type:varchar(32);default:'active'" json:"status" validate:"required,oneof=active revoked"`
	CreatedAt     time.Time         `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt     time.Time         `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RevokedAt     *time.Time        `gorm:"type:timestamptz;index" json:"revoked_at,omitempty" validate:"omitempty"`

	Issuer           *Identity          `gorm:"foreignKey:IssuerDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"issuer,omitempty"`
	SchemaAttributes []*SchemaAttribute `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema_attributes,omitempty"`
}

type SchemaAttribute struct {
	ID          uint              `gorm:"primaryKey;autoIncrement" json:"id" swagger:"-"`
	PublicID    uuid.UUID         `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	SchemaID    uint              `gorm:"column:schema_id;index;not null" json:"-" validate:"required"`
	Name        string            `gorm:"column:name;type:varchar(128);not null" json:"name" validate:"required,max=128"`
	Title       string            `gorm:"column:title;type:varchar(255);not null" json:"title" validate:"required,max=255"`
	Type        string            `gorm:"column:type;type:varchar(64);not null" json:"type" validate:"required,oneof=string number integer boolean object array null"`
	Description string            `gorm:"column:description;type:text;not null" json:"description" validate:"required,max=1000"`
	Required    bool              `gorm:"column:required;default:false" json:"required"`
	Slot        string            `gorm:"column:slot;type:varchar(64)" json:"slot,omitempty" validate:"oneof=indexSlotA indexSlotB dataSlotA dataSlotB"`
	Format      string            `gorm:"type:varchar(64)" json:"format,omitempty" validate:"omitempty,oneof=date date-time time uri email duration ipv4 ipv6 hostname"`
	Pattern     string            `gorm:"type:varchar(255)" json:"pattern,omitempty"`
	MinLength   *int              `gorm:"column:min_length" json:"min_length,omitempty"`
	MaxLength   *int              `gorm:"column:max_length" json:"max_length,omitempty"`
	Minimum     *float64          `gorm:"column:minimum" json:"minimum,omitempty"`
	Maximum     *float64          `gorm:"column:maximum" json:"maximum,omitempty"`
	Enum        datatypes.JSONMap `gorm:"type:jsonb" json:"enum,omitempty"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RevokedAt   *time.Time        `gorm:"type:timestamptz;index" json:"revoked_at,omitempty" validate:"omitempty"`
	Schema      *Schema           `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}
