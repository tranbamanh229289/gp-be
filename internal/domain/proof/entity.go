package proof

import (
	"time"

	"be/internal/domain/credential"
	"be/internal/domain/schema"
	"be/internal/shared/types"

	"github.com/google/uuid"
)

type ProofRequest struct {
	ID                uint                 `gorm:"primaryKey" json:"id"`
	PublicID          uuid.UUID            `gorm:"type:uuid;uniqueIndex" json:"public_id"`
	VerifierID        uint                 `gorm:"not null;index" json:"verifier_id"`
	SchemaID          uint                 `gorm:"not null" json:"schema_id"`
	RequestData       types.JSONB          `gorm:"type:jsonb" json:"request_data"`
	AllowedIssuersDID []string             `gorm:"type:text[]" json:"allowed_issuers_did,omitempty"`
	Challenge         string               `gorm:"type:varchar(255)" json:"challenge,omitempty"`
	Status            string               `gorm:"type:varchar(50);default:'pending';index" json:"status" validate:"required,oneof=pending completed failed expired cancelled"`
	QRCodeData        string               `gorm:"type:text" json:"qr_code_data,omitempty"`
	CallbackURL       string               `gorm:"type:text" json:"callback_url,omitempty"`
	ExpireAt          time.Time            `gorm:"index" json:"expire_at"`
	CreatedAt         time.Time            `gorm:"autoCreateTime" json:"created_at"`
	Verifier          *credential.Identity `gorm:"foreignKey:VerifierID" json:"verifier,omitempty"`
	Schema            *schema.Schema       `gorm:"foreignKey:SchemaID" json:"schema,omitempty"`
	Responses         []*ProofResponse     `gorm:"foreignKey:RequestID" json:"responses,omitempty"`
}

type ProofResponse struct {
	ID            uint                 `gorm:"primaryKey" json:"id"`
	PublicID      uuid.UUID            `gorm:"type:uuid;uniqueIndex" json:"public_id"`
	RequestID     uint                 `gorm:"not null;index" json:"request_id"`
	HolderID      uint                 `gorm:"not null;index" json:"holder_id"`
	ProofData     types.JSONB          `gorm:"type:jsonb" json:"proof_data"`
	PublicSignals types.JSONB          `gorm:"type:jsonb" json:"public_signals"`
	IsVerified    bool                 `gorm:"default:false" json:"is_verified"`
	VerifiedAt    *time.Time           `json:"verified_at,omitempty"`
	CreatedAt     time.Time            `gorm:"autoCreateTime" json:"created_at"`
	Request       ProofRequest         `gorm:"foreignKey:RequestID" json:"request,omitempty"`
	Holder        *credential.Identity `gorm:"foreignKey:HolderID" json:"holder,omitempty"`
}
