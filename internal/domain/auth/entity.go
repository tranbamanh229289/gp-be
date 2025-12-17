package auth

import (
	"be/internal/shared/constant"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64             `gorm:"primaryKey" json:"id"`
	PublicID  uuid.UUID         `gorm:"type:uuid;uniqueIndex" json:"public_id"`
	Name      string            `gorm:"type:varchar(100);not null" json:"name"`
	Email     string            `gorm:"type:varchar(100);not null" json:"email"`
	Password  string            `gorm:"type:varchar(100);not null" json:"password"`
	Role      constant.UserRole `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	CreatedAt time.Time         `gorm:"autoCreateTime" json:"created_at"`
}
