package service

import (
	"be/config"
	"context"
)

type ICredentialService interface {
}
type CredentialService struct {
	config *config.Config
}

func NewCredentialService(config *config.Config) ICredentialService {
	return &CredentialService{
		config: config,
	}
}

func IssueCredential(ctx context.Context) {

}
