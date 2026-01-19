package proof

import (
	"be/internal/domain/schema"
	"be/internal/shared/constant"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProofRequest struct {
	ID                       uint                        `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID                 uuid.UUID                   `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	ThreadID                 string                      `gorm:"column:thread_id;type:varchar(255);not null" json:"thread_id" validate:"required"`
	VerifierDID              string                      `gorm:"column:verifier_did;type:varchar(255);not null;index" json:"verifier_did" validate:"required,startswith=did:"`
	CallbackURL              string                      `gorm:"column:callback_url;type:text;not null" json:"callback_url" validate:"required,url"`
	Reason                   string                      `gorm:"column:reason;type:text" json:"reason,omitempty" validate:"omitempty,max=500"`
	Message                  string                      `gorm:"column:message;type:text" json:"message,omitempty" validate:"omitempty,max=1000"`
	ScopeID                  uint32                      `gorm:"column:scope_id;not null" json:"scope_id" validate:"required"`
	CircuitID                string                      `gorm:"column:circuit_id;type:varchar(100);not null" json:"circuit_id" validate:"required"`
	AllowedIssuers           datatypes.JSONSlice[string] `gorm:"column:allowed_issuers_did;type:json;not null" json:"allowed_issuers_did" validate:"required"`
	CredentialSubject        datatypes.JSONMap           `gorm:"column:credential_subject;type:jsonb" json:"credential_subject,omitempty" validate:"omitempty"`
	SchemaID                 uint                        `gorm:"column:schema_id;not null" json:"schema_id" validate:"required"`
	ProofType                string                      `gorm:"column:proof_type;type:varchar(100)" json:"proof_type,omitempty" validate:"omitempty"`
	SkipClaimRevocationCheck bool                        `gorm:"column:skip_claim_revocation_check;type:boolean;default:false" json:"skip_claim_revocation_check,omitempty" validate:"omitempty"`
	GroupID                  int                         `gorm:"column:group_id" json:"group_id,omitempty" validate:"omitempty"`
	NullifierSession         string                      `gorm:"column:nullifier_session;type:text" json:"nullifier_session,omitempty" validate:"omitempty"`
	Status                   constant.ProofRequestStatus `gorm:"column:status;type:varchar(50);default:'active'" json:"status" validate:"required"`
	CreatedTime              *int64                      `gorm:"column:created_time;type:bigint" json:"created_time,omitempty" validate:"omitempty"`
	ExpiresTime              *int64                      `gorm:"column:expires_time;type:bigint" json:"expires_time,omitempty" validate:"omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`

	ProofResponses []*ProofResponse `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"proof_responses,omitempty"`
	Schema         *schema.Schema   `gorm:"foreignKey:SchemaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"schema,omitempty"`
	Verifier       *schema.Identity `gorm:"foreignKey:VerifierDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"verifier,omitempty"`
}

type ProofResponse struct {
	ID           uint                         `gorm:"primaryKey;autoIncrement" json:"id,omitempty" validate:"-"`
	PublicID     uuid.UUID                    `gorm:"column:public_id;type:uuid;uniqueIndex;default:gen_random_uuid()" json:"public_id" validate:"required"`
	RequestID    uint                         `gorm:"column:request_id;not null;index" json:"request_id" validate:"required"`
	HolderDID    string                       `gorm:"column:holder_did;type:varchar(255);not null;index" json:"holder_did" validate:"required,startswith=did:"`
	Status       constant.ProofResponseStatus `gorm:"column:status;type:varchar(50);default:'pending';index" json:"status" validate:"required"`
	ProofRequest *ProofRequest                `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"request,omitempty"`
	Holder       *schema.Identity             `gorm:"foreignKey:HolderDID;references:DID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"holder,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" validate:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" validate:"-"`
}
