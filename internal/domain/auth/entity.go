package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID int64 						`gorm:"primaryKey" json:"id"`
	PublicID uuid.UUID  `gorm:"type:uuid" json:"public_id"`
	Name string					`gorm:"type:varchar(100);not null" json:"name"`	
	Email string				`gorm:"type:varchar(100);not null" json:"email"`	
	Role UserRole				`gorm:"type:varchar(20);not null;default:'user'" json:"role"`	
	CreatedAt time.Time	`gorm:"autoCreateTime" json:"create_at"`
}

type UserRole string

const (
	UserRoleUser UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)