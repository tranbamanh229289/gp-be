package gist

import (
	"be/internal/domain/schema"
	"time"

	"github.com/google/uuid"
)

type StateTransition struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID    uuid.UUID  `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	IdentityID  uint       `gorm:"column:identity_id;not null;index" json:"identity_id" validate:"required,gt=0"`
	OldState    string     `gorm:"column:old_state;type:varchar(255);not null" json:"old_state" validate:"required,len=64"`
	NewState    string     `gorm:"column:new_state;varchar(255);not null" json:"new_state" validate:"required,len=64"`
	TxHash      string     `gorm:"column:tx_hash;type:varchar(100);index" json:"tx_hash,omitempty" validate:"omitempty,len=66,startswith=0x"`
	BlockNumber int64      `gorm:"column:block_number;index" json:"block_number,omitempty" validate:"omitempty,gte=0"`
	Timestamp   *time.Time `gorm:"column:time_stamp;type:timestamptz;index" json:"timestamp,omitempty" validate:"omitempty"`
	IsGenesis   bool       `gorm:"column:is_genesis;not null;default:false" json:"is_genesis" validate:"-"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty" validate:"-"`

	// Relationship
	Identity *schema.Identity `gorm:"foreignKey:IdentityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"identity,omitempty"`
}
