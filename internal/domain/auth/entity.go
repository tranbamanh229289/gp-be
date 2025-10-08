package auth

import "time"

type User struct {
	ID int64 						`gorm:"primaryKey"`
	Name string					`gorm:"type:varchar(100);not null"`	
	Email string				`gorm:"type:varchar(100);not null"`	
	Role UserRole				`gorm:"type:varchar(20);not null;default:'user'"`	
	CreatedAt time.Time	`gorm:"autoCreateTime"`
}

type UserRole string

const (
	UserRoleUser UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)