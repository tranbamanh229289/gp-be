package schema

import (
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type Schema struct {
	ID            uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	PublicID      uuid.UUID   `gorm:"type:uuid;uniqueIndex" json:"public_id"`
	IssuerID      uint        `gorm:"index" json:"issuer_id"`
	Type          string      `gorm:"type:varchar(255);index" json:"type"`
	Version       string      `gorm:"type:varchar(64);index:,composite:type_version" json:"version"`
	JSONSchema    types.JSONB `gorm:"type:jsonb" json:"json_schema"`
	JSONLDContext types.JSONB `gorm:"type:jsonb" json:"jsonld_context"`
	JSONCID       string      `gorm:"type:varchar(255);uniqueIndex" json:"json_cid"`
	JSONLDCID     string      `gorm:"type:varchar(255)" json:"jsonld_cid"`
	Status        string      `gorm:"type:varchar(32);index;default:'active'" json:"status"` // active, revoked
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	RevokeAt      *time.Time  `gorm:"index" json:"revoke_at,omitempty"`
}
