package dtos

import (
	"be/internal/shared/constant"

	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Email string 		`json:"email" binding:"required,email"`
	Password string	`json:"password" binding:"required"`
	Name string			`json:"name" binding:"required"`
}

type LoginRequest struct {
	Email string		`json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRequest struct {
	Email string 		`json:"email" binding:"required"`
	Name string			`json:"name" binding:"required"`	
}

type UserResponse struct {
	ID string 			`json:"id" binding:"required,uuid"`
	Email string 		`json:"email" binding:"required"`
	Name string			`json:"name" binding:"required"`	
}

type TokenResponse struct {
	AccessToken string		`json:"access_token" binding:"required"`
	RefreshToken string		`json:"refresh_token" binding:"required"`
}

type Claims struct {
	ID string									`json:"id"`
	Email string							`json:"email"`
	Name  string							`json:"name"`
	Role 	constant.UserRole		`json:"role"`
	jwt.RegisteredClaims
}

