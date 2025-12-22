package claim

import (
	"be/internal/domain/credential"
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type StateTransition struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID    uuid.UUID  `gorm:"type:uuid;uniqueIndex;default:uuid_generate_v4()" json:"public_id"`
	IdentityID  uint       `gorm:"not null;index" json:"identity_id" validate:"required"`
	OldState    types.Hash `gorm:"type:varchar(255);not null" json:"old_state"`
	NewState    types.Hash `gorm:"type:varchar(255);not null" json:"new_state"`
	TxHash      string     `gorm:"type:varchar(100);index" json:"tx_hash,omitempty"`
	BlockNumber int64      `gorm:"index" json:"block_number,omitempty"`
	Timestamp   *time.Time `gorm:"index" json:"timestamp,omitempty"`
	IsGenesis   bool       `gorm:"default:false" json:"is_genesis"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`

	// Relationship
	Identity *credential.Identity `gorm:"foreignKey:IdentityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"identity,omitempty"`
}
