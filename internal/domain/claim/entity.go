package claim

import (
	"be/internal/domain/credential"
	"be/internal/shared/types"
	"time"

	"github.com/google/uuid"
)

type StateTransition struct {
	ID          uint
	PublicID    uuid.UUID
	IdentityID  uint
	OldState    types.Hash
	NewState    types.Hash
	TxHash      string
	BlockNumber int64
	Timestamp   *time.Time
	IsGenesis   bool
	CreatedAt   time.Time
	Identity    *credential.Identity
}
