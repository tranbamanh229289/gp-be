package schema

import (
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type Schema struct {
	ID            uint        `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID      uuid.UUID   `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	IssuerID      uint        `gorm:"column:issuer_id;not null;index" json:"issuer_id" validate:"required,gt=0"`
	Type          string      `gorm:"column:type;type:varchar(255);not null;index" json:"type" validate:"required,min=3,max=255"`
	Version       string      `gorm:"column:version;type:varchar(64);not null" json:"version" validate:"required,semver"`
	Title         string      `gorm:"column:title;type:varchar(255)" json:"title,omitempty" validate:"omitempty,max=255"`
	Description   string      `gorm:"column:description;type:text" json:"description,omitempty" validate:"omitempty,max=2000"`
	JSONSchema    types.JSONB `gorm:"column:json_schema;type:jsonb;not null" json:"json_schema" validate:"required"`
	JSONLDContext types.JSONB `gorm:"column:jsonld_context;type:jsonb;not null" json:"jsonld_context" validate:"required"`
	JSONCID       string      `gorm:"column:json_cid;type:varchar(255);uniqueIndex" json:"json_cid" validate:"required,len=59,startswith=Qm"`
	JSONLDCID     string      `gorm:"column:jsonld_cid;type:varchar(255);uniqueIndex" json:"jsonld_cid" validate:"required,len=59,startswith=Qm"`
	Status        string      `gorm:"column:status;type:varchar(32);default:'active';index" json:"status" validate:"required,oneof=active revoked"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`
	RevokedAt     *time.Time  `gorm:"type:timestamptz;index" json:"revoked_at,omitempty" validate:"omitempty"`
}
