package service

import (
	"be/config"
	"context"
)

type IAuthZkService interface {
}

type AuthZkService struct {
	config *config.Config
}

func NewAuthZkService(config *config.Config) *AuthZkService {
	return &AuthZkService{
		config: config,
	}
}

func Login(ctx *context.Context) {

}

func Logout(ctx *context.Context) {

}
